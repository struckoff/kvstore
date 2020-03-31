package node

import (
	"github.com/pkg/errors"
	balancer "github.com/struckoff/SFCFramework"
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

func (n *InternalNode) ID() string                  { return n.id }
func (n *InternalNode) Power() balancer.Power       { return n.p }
func (n *InternalNode) Capacity() balancer.Capacity { return n.c }

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

// Return value for a given key from local storage
func (n *InternalNode) Receive(key string) ([]byte, error) {
	var body []byte
	err := n.db.View(func(tx *bolt.Tx) error {
		bc := tx.Bucket(mainBucket)
		if bc == nil {
			return nil
		}
		body = bc.Get([]byte(key))
		return nil
	})
	return body, err
}
func (n *InternalNode) Explore() ([]string, error) {
	var res []string
	err := n.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(mainBucket)
		if b == nil {
			return errors.New("bucket not found")
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
		ID:         n.id,
		Address:    n.address,
		RPCAddress: n.rpcaddress,
		Power:      n.p.Get(),
		Capacity:   n.p.Get(),
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
