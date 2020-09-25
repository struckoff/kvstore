package balanceradapter

import (
	"github.com/lafikl/consistent"
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/nodes"
	balancer "github.com/struckoff/sfcframework"
	"sync"
)

type Consistent struct {
	ring  *consistent.Consistent
	nodes sync.Map
}

func NewConsistentBalancer() *Consistent {
	return &Consistent{
		ring: consistent.New(),
	}
}

func (c *Consistent) AddNode(n nodes.Node) error {
	c.nodes.Store(n.ID(), n)
	c.ring.Add(n.ID())
	return nil
}

func (c *Consistent) RemoveNode(id string) error {
	c.ring.Remove(id)
	c.nodes.Delete(id)
	return nil
}

func (c *Consistent) SetNodes(ns []nodes.Node) error {
	c.ring = consistent.New()
	c.nodes = sync.Map{}
	for _, n := range ns {
		c.nodes.Store(n.ID(), n)
		c.ring.Add(n.ID())
	}
	return nil
}

func (c *Consistent) LocateData(di balancer.DataItem) (nodes.Node, uint64, error) {
	name, err := c.ring.GetLeast(di.ID())
	if err != nil {
		return nil, 0, err
	}
	ni, ok := c.nodes.Load(name)
	if !ok {
		return nil, 0, errors.New("node not found")
	}
	n, ok := ni.(nodes.Node)
	if !ok {
		return nil, 0, errors.New("wrong node type")
	}
	return n, 0, nil
}

func (c *Consistent) AddData(di balancer.DataItem) (nodes.Node, uint64, error) {
	return c.LocateData(di)
}

func (c *Consistent) RemoveData(di balancer.DataItem) error {
	return nil
}

func (c *Consistent) Nodes() ([]nodes.Node, error) {
	names := c.ring.Hosts()
	ns := make([]nodes.Node, len(names))
	for i := range names {
		ni, ok := c.nodes.Load(names[i])
		if !ok {
			return nil, errors.New("node not found")
		}
		n, ok := ni.(nodes.Node)
		if !ok {
			return nil, errors.New("wrong node type")
		}
		ns[i] = n
	}
	return ns, nil
}

func (c *Consistent) GetNode(id string) (nodes.Node, error) {
	ni, ok := c.nodes.Load(id)
	if !ok {
		return nil, errors.New("node not found")
	}
	n, ok := ni.(nodes.Node)
	if !ok {
		return nil, errors.New("wrong node type")
	}
	return n, nil
}

func (c *Consistent) Optimize() error {
	return nil
}

func (c *Consistent) Reset() error {
	return nil
}
