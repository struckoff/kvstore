package kvstore

import (
	consulapi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	balancer "github.com/struckoff/SFCFramework"
	kvrouter "github.com/struckoff/kvrouter/router"
	"github.com/struckoff/kvrouter/rpcapi"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"sync"
)

var mainBucket = []byte("pairs")

// InternalNode represents local node
type InternalNode struct {
	mu         sync.RWMutex
	id         string
	address    string
	rpcaddress string
	rpcserver  *grpc.Server
	p          kvrouter.Power
	c          kvrouter.Capacity
	db         *bolt.DB
	kvr        *kvrouter.Router
	consul     *consulapi.Client
	kvrAgent   rpcapi.RPCBalancerClient
}

func (inn *InternalNode) RunHTTPServer(addr string) error {
	h := inn.kvr.HTTPHandler()
	if err := http.ListenAndServe(addr, h); err != nil {
		return err
	}
	return nil
}

//ID returns the node's ID
func (inn *InternalNode) ID() string {
	inn.mu.RLock()
	defer inn.mu.RUnlock()
	return inn.id
}

//RPCAddress returns the node's rpc address
func (inn *InternalNode) RPCAddress() string {
	inn.mu.RLock()
	defer inn.mu.RUnlock()
	return inn.rpcaddress
}

//HTTPAddress returns the node's http address
func (inn *InternalNode) HTTPAddress() string {
	inn.mu.RLock()
	defer inn.mu.RUnlock()
	return inn.address
}

//Power returns the node's power
func (inn *InternalNode) Power() balancer.Power {
	inn.mu.RLock()
	defer inn.mu.RUnlock()
	return inn.p
}

//Capacity returns the node's capacity
func (inn *InternalNode) Capacity() balancer.Capacity {
	inn.mu.RLock()
	defer inn.mu.RUnlock()
	return inn.c
}

// Store value for a given key in local storage
func (inn *InternalNode) Store(key string, body []byte) error {
	return inn.db.Update(func(tx *bolt.Tx) error {
		bc, err := tx.CreateBucketIfNotExists(mainBucket)
		if err != nil {
			return err
		}
		return bc.Put([]byte(key), body)
	})
}

// Store KV pairs in local storage
func (inn *InternalNode) StorePairs(pairs []*rpcapi.KeyValue) error {
	err := inn.db.Update(func(tx *bolt.Tx) error {
		bc, err := tx.CreateBucketIfNotExists(mainBucket)
		if err != nil {
			return err
		}
		for iter := range pairs {
			if err := bc.Put([]byte(pairs[iter].Key), pairs[iter].Value); err != nil {
				return errors.Wrap(err, "failed to store pair")
			}
		}
		return nil
	})
	return err
}

// Return value for a given key from local storage
func (inn *InternalNode) Receive(key string) ([]byte, error) {
	var body []byte
	err := inn.db.View(func(tx *bolt.Tx) error {
		bc := tx.Bucket(mainBucket)
		if bc == nil {
			return errors.New("unable to receive value, bucket not found")
		}
		body = bc.Get([]byte(key))
		return nil
	})
	return body, err
}

// Remove value for a given key
func (inn *InternalNode) Remove(key string) error {
	err := inn.db.Update(func(tx *bolt.Tx) error {
		bc := tx.Bucket(mainBucket)
		if bc == nil {
			return nil
		}
		return bc.Delete([]byte(key))
	})
	if err != nil {
		return errors.Wrap(err, "failed to remove key")
	}
	return nil
}

// Move values for a given keys to another node
func (inn *InternalNode) Move(nk map[kvrouter.Node][]string) error {
	var wg sync.WaitGroup
	for en, keys := range nk {
		if len(keys) == 0 {
			continue
		}
		wg.Add(1)
		go func(en kvrouter.Node, keys []string, wg *sync.WaitGroup) {
			defer wg.Done()
			err := inn.db.Update(func(tx *bolt.Tx) error {
				bc := tx.Bucket(mainBucket)
				if bc == nil {
					return nil
				}
				pairs := make([]*rpcapi.KeyValue, len(keys))
				for iter := range keys {
					body := bc.Get([]byte(keys[iter]))
					pairs[iter] = &rpcapi.KeyValue{Key: keys[iter], Value: body}
				}
				if err := en.StorePairs(pairs); err != nil {
					return errors.Wrap(err, "failed to move keys and values to another node")
				}

				for iter := range keys {
					if err := bc.Delete([]byte(keys[iter])); err != nil {
						return errors.Wrap(err, "failed to delete keys")
					}
				}
				return nil
			})
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("%d keys relocated to node(%s)", len(keys),en.ID())
		}(en, keys, &wg)
	}
	wg.Wait()
	return nil
}

// Explore returns the list of keys in local storage.
func (inn *InternalNode) Explore() ([]string, error) {
	res := make([]string, 0)
	err := inn.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(mainBucket)
		if b == nil {
			//return errors.New("bucket not found")
			return nil
		}
		err := b.ForEach(func(k, v []byte) error {
			res = append(res, string(k))
			return nil
		})
		return err
	})
	return res, err
}

// Return meta information about the node
func (inn *InternalNode) Meta() rpcapi.NodeMeta {
	inn.mu.RLock()
	defer inn.mu.RUnlock()
	return rpcapi.NodeMeta{
		ID:         inn.ID(),
		Address:    inn.HTTPAddress(),
		RPCAddress: inn.RPCAddress(),
		Power:      inn.Power().Get(),
		Capacity:   inn.Capacity().Get(),
	}
}

// Return new instance InternalNode.
func NewInternalNode(conf *Config, db *bolt.DB, kvr *kvrouter.Router) *InternalNode {
	return &InternalNode{
		id:         *conf.Name,
		address:    conf.Address,
		rpcaddress: conf.RPCAddress,
		p:          kvrouter.NewPower(conf.Power),
		c:          kvrouter.NewCapacity(conf.Capacity),
		db:         db,
		kvr:        kvr,
	}
}
