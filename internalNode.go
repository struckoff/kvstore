package kvstore

import (
	"github.com/pkg/errors"
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/kvstore/rpcapi"
	bolt "go.etcd.io/bbolt"
	"sync"
)

var mainBucket = []byte("pairs")

// InternalNode represents local node
type InternalNode struct {
	mu         sync.RWMutex
	id         string
	address    string
	rpcaddress string
	p          Power
	c          Capacity
	db         *bolt.DB
}

func (n *InternalNode) SetID(id string) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.id = id
}

func (n *InternalNode) ID() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.id
}

func (n *InternalNode) RPCAddress() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.rpcaddress
}

func (n *InternalNode) HTTPAddress() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.address
}

func (n *InternalNode) Power() balancer.Power {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.p
}
func (n *InternalNode) Capacity() balancer.Capacity {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.c
}

// Store value for a given key in local storage
func (n *InternalNode) Store(key string, body []byte) error {
	return n.db.Update(func(tx *bolt.Tx) error {
		bc, err := tx.CreateBucketIfNotExists(mainBucket)
		if err != nil {
			return err
		}
		return bc.Put([]byte(key), body)
	})
}

func (n *InternalNode) StorePairs(pairs []*rpcapi.KeyValue) error {
	err := n.db.Update(func(tx *bolt.Tx) error {
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
func (n *InternalNode) Receive(key string) ([]byte, error) {
	var body []byte
	err := n.db.View(func(tx *bolt.Tx) error {
		bc := tx.Bucket(mainBucket)
		if bc == nil {
			return errors.New("unable to receive value, bucket not found")
		}
		body = bc.Get([]byte(key))
		return nil
	})
	return body, err
}
func (n *InternalNode) Remove(key string) error {
	err := n.db.Update(func(tx *bolt.Tx) error {
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

func (n *InternalNode) Move(keys []string, en Node) error {
	err := n.db.Update(func(tx *bolt.Tx) error {
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
		return err
	}
	return nil
}

func (n *InternalNode) Explore() ([]string, error) {
	res := make([]string, 0)
	err := n.db.View(func(tx *bolt.Tx) error {
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
func (n *InternalNode) Meta() NodeMeta {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return NodeMeta{
		ID:         n.ID(),
		Address:    n.HTTPAddress(),
		RPCAddress: n.RPCAddress(),
		Power:      n.Power().Get(),
		Capacity:   n.Capacity().Get(),
	}
}

func NewInternalNode(id, addr, raddr string, p float64, c float64, db *bolt.DB) *InternalNode {
	return &InternalNode{
		id:         id,
		address:    addr,
		rpcaddress: raddr,
		p:          NewPower(p),
		c:          NewCapacity(c),
		db:         db,
	}
}