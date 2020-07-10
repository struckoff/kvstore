package store

import (
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/struckoff/kvstore/router/nodes"
	"log"
	"net/http"
	"sync"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/kvstore/router"
	"github.com/struckoff/kvstore/router/rpcapi"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc"
)

var mainBucket = []byte("pairs")

// Return new instance LocalNode.
func NewLocalNode(conf *Config, db *bolt.DB, kvr *router.Router, metrics chan<- *influxdb2.Point) (*LocalNode, error) {
	lwID := int64(0)
	ln := &LocalNode{
		id:          *conf.Name,
		address:     conf.Address,
		rpcaddress:  conf.RPCAddress,
		p:           nodes.NewPower(conf.Power),
		c:           NewCapacity(conf.Capacity),
		db:          db,
		kvr:         kvr,
		geo:         conf.Geo,
		rpclatency:  conf.Latency.Duration,
		httplatency: conf.Balancer.Latency.Duration,
		lwID:        &lwID,
		metrics:     metrics,
	}
	if ln.kvr != nil {
		h, err := ln.kvr.Hasher().Sum(ln.meta())
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate hash sum")
		}
		ln.h = h
	}
	if conf.Mode == ConsulMode {
		consul, err := consulapi.NewClient(&conf.Consul.Config)
		if err != nil {
			return nil, err
		}
		ln.consul = consul
	}
	keys, err := ln.Explore()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to explore local node")
	}
	ln.c.Add(-float64(len(keys))) // reduce capacity
	return ln, nil
}

// LocalNode represents local node
type LocalNode struct {
	mu          sync.RWMutex
	id          string
	address     string
	rpcaddress  string
	rpcserver   *grpc.Server
	p           nodes.Power
	c           Capacity
	db          *bolt.DB
	kvr         *router.Router
	consul      *consulapi.Client
	kvrAgent    rpcapi.RPCBalancerClient
	geo         *rpcapi.GeoData
	h           uint64
	rpclatency  time.Duration
	httplatency time.Duration
	lwID        *int64
	metrics     chan<- *influxdb2.Point
}

func (inn *LocalNode) RunHTTPServer(addr string) error {
	h := inn.kvr.HTTPHandler()
	l := router.LatencyMiddleware(h, inn.httplatency)
	log.Printf("HTTP server listening on %s", addr)
	if err := http.ListenAndServe(addr, l); err != nil {
		return err
	}
	return nil
}

//ID returns the node's ID
func (inn *LocalNode) ID() string {
	inn.mu.RLock()
	defer inn.mu.RUnlock()
	return inn.id
}

//RPCAddress returns the node's rpc address
func (inn *LocalNode) RPCAddress() string {
	inn.mu.RLock()
	defer inn.mu.RUnlock()
	return inn.rpcaddress
}

//HTTPAddress returns the node's http address
func (inn *LocalNode) HTTPAddress() string {
	inn.mu.RLock()
	defer inn.mu.RUnlock()
	return inn.address
}

//Power returns the node's power
func (inn *LocalNode) Power() balancer.Power {
	inn.mu.RLock()
	defer inn.mu.RUnlock()
	return inn.p
}

//Capacity returns the node's capacity
func (inn *LocalNode) Capacity() balancer.Capacity {
	inn.mu.RLock()
	defer inn.mu.RUnlock()
	return &inn.c
}

func (inn *LocalNode) Hash() uint64 {
	inn.mu.RLock()
	defer inn.mu.RUnlock()
	return inn.h
}

func (inn *LocalNode) SetHash(h uint64) {
	inn.mu.Lock()
	defer inn.mu.Unlock()
	inn.h = h
}

// Store value for a given key in local storage
func (inn *LocalNode) Store(key string, body []byte) error {
	err := inn.db.Update(func(tx *bolt.Tx) error {
		bc, err := tx.CreateBucketIfNotExists(mainBucket)
		if err != nil {
			return err
		}
		return bc.Put([]byte(key), body)
	})
	if err != nil {
		return err
	}
	inn.c.Add(-1) // reduce capacity
	return nil
}

// Store KV pairs in local storage
func (inn *LocalNode) StorePairs(pairs []*rpcapi.KeyValue) error {
	cp := 0.0
	err := inn.db.Update(func(tx *bolt.Tx) error {
		bc, err := tx.CreateBucketIfNotExists(mainBucket)
		if err != nil {
			return err
		}
		for iter := range pairs {
			if err := bc.Put([]byte(pairs[iter].Key), []byte(pairs[iter].Value)); err != nil {
				return errors.Wrap(err, "failed to store pair")
			}
			cp++
		}
		return nil
	})
	inn.c.Add(-cp) // reduce capacity
	return err
}

// Return value for a given key from local storage
func (inn *LocalNode) Receive(keys []string) (*rpcapi.KeyValues, error) {
	kvs := &rpcapi.KeyValues{
		KVs: make([]*rpcapi.KeyValue, len(keys)),
	}
	err := inn.db.View(func(tx *bolt.Tx) error {
		bc := tx.Bucket(mainBucket)
		if bc == nil {
			return errors.New("unable to receive value, bucket not found")
		}
		for iter := range keys {
			val := bc.Get([]byte(keys[iter]))
			ok := val != nil
			kvs.KVs[iter] = &rpcapi.KeyValue{
				Key:   keys[iter],
				Value: string(val),
				Found: ok,
			}
		}
		return nil
	})
	return kvs, err
}

// Remove value for a given key
func (inn *LocalNode) Remove(keys []string) error {
	cp := 0.0
	err := inn.db.Update(func(tx *bolt.Tx) error {
		bc := tx.Bucket(mainBucket)
		if bc == nil {
			return nil
		}
		for iter := range keys {
			if err := bc.Delete([]byte(keys[iter])); err != nil {
				return err
			}
			cp++
		}
		return nil
	})
	inn.c.Add(cp) // increase capacity

	if err != nil {
		return errors.Wrap(err, "failed to remove key")
	}
	return nil
}

// Move values for a given keys to another node
func (inn *LocalNode) Move(nk map[nodes.Node][]string) error {
	var wg sync.WaitGroup
	for en, keys := range nk {
		if len(keys) == 0 {
			continue
		}
		wg.Add(1)
		go func(en nodes.Node, keys []string, wg *sync.WaitGroup) {
			defer wg.Done()
			err := inn.db.Update(func(tx *bolt.Tx) error {
				bc := tx.Bucket(mainBucket)
				if bc == nil {
					return nil
				}
				pairs := make([]*rpcapi.KeyValue, len(keys))
				for iter := range keys {
					body := bc.Get([]byte(keys[iter]))
					pairs[iter] = &rpcapi.KeyValue{Key: keys[iter], Value: string(body)}
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
			log.Printf("%d keys relocated to node(%s)", len(keys), en.ID())
		}(en, keys, &wg)
	}
	wg.Wait()
	return nil
}

// Explore returns the list of keys in local storage.
func (inn *LocalNode) Explore() ([]string, error) {
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
func (inn *LocalNode) Meta() *rpcapi.NodeMeta {
	inn.mu.RLock()
	defer inn.mu.RUnlock()
	return inn.meta()
}

func (inn *LocalNode) meta() *rpcapi.NodeMeta {
	cp, err := inn.c.Get()
	if err != nil {
		return nil
	}
	return &rpcapi.NodeMeta{
		ID:         inn.ID(),
		Address:    inn.HTTPAddress(),
		RPCAddress: inn.RPCAddress(),
		Power:      inn.Power().Get(),
		Capacity:   cp,
		Geo:        inn.geo,
	}
}

//func MakeHasher(conf *router.BalancerConfig) Hasher{
//	switch conf.NodeHash:
//
//}
