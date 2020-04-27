package balanceradapter

import (
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"github.com/pkg/errors"
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/kvstore/router/nodes"
)

type ConsistentRing struct {
	ring *consistent.Consistent
}

func NewConsistentBalancer(conf *consistent.Config) *ConsistentRing {
	conf.Hasher = hasher{}
	return &ConsistentRing{
		ring: consistent.New(nil, *conf),
	}
}

func (c ConsistentRing) AddNode(n nodes.Node) error {
	c.ring.Add(RingNode{n})
	return nil
}

func (c ConsistentRing) RemoveNode(id string) error {
	c.ring.Remove(id)
	return nil
}

func (c ConsistentRing) SetNodes(ns []nodes.Node) error {
	ms := c.ring.GetMembers()
	for _, m := range ms {
		c.ring.Remove(m.String())
	}
	for _, n := range ns {
		c.ring.Add(RingNode{n})
	}
	return nil
}

func (c ConsistentRing) LocateData(di balancer.DataItem) (nodes.Node, error) {
	m := c.ring.LocateKey([]byte(di.ID()))
	if rn, ok := m.(RingNode); ok {
		return rn, nil
	}
	return nil, errors.New("wrong node type")
}

func (c ConsistentRing) Nodes() ([]nodes.Node, error) {
	var ok bool
	ms := c.ring.GetMembers()
	ns := make([]nodes.Node, len(ms))
	for iter := range ms {
		if ns[iter], ok = ms[iter].(RingNode); !ok {
			return nil, errors.New("wrong node type")
		}
	}
	return ns, nil
}

func (c ConsistentRing) GetNode(id string) (nodes.Node, error) {
	ms := c.ring.GetMembers()
	for _, m := range ms {
		if m.String() == id {
			if rn, ok := m.(RingNode); ok {
				return rn, nil
			}
			return nil, errors.New("wrong node type")
		}
	}
	return nil, errors.New("node not found")
}

// consistent package doesn't provide a default hashing function.
// You should provide a proper one to distribute keys/members uniformly.
type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	// you should use a proper hash function for uniformity.
	return xxhash.Sum64(data)
}

type RingNode struct {
	nodes.Node
}

func (rn RingNode) String() string {
	return rn.ID()
}

func (rn RingNode) Self() nodes.Node {
	return rn.Node
}
