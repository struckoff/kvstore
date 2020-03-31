package kvstore

import (
	"errors"
	balancer "github.com/struckoff/SFCFramework"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

// ExternalNode represents compunction API with cluster unit
// It also contains meta information
type ExternalNode struct {
	mu      sync.RWMutex
	id      string
	address string
	p       Power
	c       Capacity
}

func (n *ExternalNode) ID() string                  { return n.id }
func (n *ExternalNode) Power() balancer.Power       { return n.p }
func (n *ExternalNode) Capacity() balancer.Capacity { return n.c }

//Save value for a given key on the remote node
func (n *ExternalNode) Store(key string, body io.Reader) error {
	log.Printf("Store key(%s) on %s", key, n.id)
	p := strings.Join([]string{"http:/", n.address, "kv", key}, "/")
	r, err := http.Post(p, "application/text", body)
	if err != nil {
		return err
	}
	if r.StatusCode >= 400 {
		return errors.New(r.Status)
	}
	return nil
}

//Receive value for a given key from the remote node
func (n *ExternalNode) Receive(key string) ([]byte, error) {
	log.Printf("Receive key(%s) from %s", key, n.id)
	p := strings.Join([]string{"http:/", n.address, "kv", key}, "/")
	r, err := http.Get(p)
	if err != nil {
		return nil, err
	}
	if r.StatusCode >= 400 {
		return nil, errors.New(r.Status)
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	return b, err
}

//TODO: implement
func (n *ExternalNode) Explore() ([]byte, error) { return nil, nil }

// Return meta information about the node
func (n *ExternalNode) Meta() NodeMeta {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return NodeMeta{
		ID:       n.id,
		Address:  n.address,
		Power:    n.p.Get(),
		Capacity: n.p.Get(),
	}
}

func NewExternalNode(meta NodeMeta) *ExternalNode {
	return &ExternalNode{
		id:      meta.ID,
		address: meta.Address,
		p:       NewPower(meta.Power),
		c:       NewCapacity(meta.Capacity),
	}
}
