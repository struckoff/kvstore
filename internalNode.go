package kvstore

import (
	balancer "github.com/struckoff/SFCFramework"
	bolt "go.etcd.io/bbolt"
	"io"
	"io/ioutil"
	"sync"
)

var mainBucket = []byte("pairs")

// InternalNode represents local node
type InternalNode struct {
	mu      sync.RWMutex
	id      string
	address string
	p       Power
	c       Capacity
	db      *bolt.DB
}

func (n *InternalNode) ID() string                  { return n.id }
func (n *InternalNode) Power() balancer.Power       { return n.p }
func (n *InternalNode) Capacity() balancer.Capacity { return n.c }

// Store value for a given key in local storage
func (n *InternalNode) Store(key string, body io.Reader) error {
	return n.db.Update(func(tx *bolt.Tx) error {
		bc, err := tx.CreateBucketIfNotExists(mainBucket)
		if err != nil {
			return err
		}
		b, err := ioutil.ReadAll(body)
		if err != nil {
			return err
		}
		return bc.Put([]byte(key), b)
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
func (n *InternalNode) Explore() ([]byte, error) {
	return nil, nil
}

// Return meta information about the node
func (n *InternalNode) Meta() NodeMeta {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return NodeMeta{
		ID:       n.id,
		Address:  n.address,
		Power:    n.p.Get(),
		Capacity: n.p.Get(),
	}
}

func NewInternalNode(id, address string, p float64, c float64, db *bolt.DB) *InternalNode {
	return &InternalNode{
		id:      id,
		address: address,
		p:       NewPower(p),
		c:       NewCapacity(c),
		db:      db,
	}
}
