package main

import (
	bolt "go.etcd.io/bbolt"
	"io"
	"io/ioutil"
)

type InternalNode struct {
	id      string
	address string
	p       Power
	c       Capacity
	db      *bolt.DB
}

func (n *InternalNode) Power() Power       { return n.p }
func (n *InternalNode) Capacity() Capacity { return n.c }
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

func NewInternalNode(id, address string, p float64, c float64, db *bolt.DB) *InternalNode {
	return &InternalNode{
		id,
		address,
		NewPower(p),
		NewCapacity(c),
		db,
	}
}
