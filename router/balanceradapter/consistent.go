package balanceradapter

import (
	"github.com/lafikl/consistent"
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/nodes"
	balancer "github.com/struckoff/sfcframework"
	"sync"
)

//Consistent hash ring adapter
type Consistent struct {
	ring  *consistent.Consistent
	nodes sync.Map
}

//NewConsistentBalancer create new Consistent balancer adapter instance
func NewConsistentBalancer() *Consistent {
	return &Consistent{
		ring: consistent.New(),
	}
}

//AddNode places new node on the ring
func (c *Consistent) AddNode(n nodes.Node) error {
	c.nodes.Store(n.ID(), n)
	c.ring.Add(n.ID())
	return nil
}

//RemoveNode from the ring
func (c *Consistent) RemoveNode(id string) error {
	c.ring.Remove(id)
	c.nodes.Delete(id)
	return nil
}

//SetNodes wipes all nodes from the balancer and set a new ones
func (c *Consistent) SetNodes(ns []nodes.Node) error {
	c.ring = consistent.New()
	c.nodes = sync.Map{}
	for _, n := range ns {
		c.nodes.Store(n.ID(), n)
		c.ring.Add(n.ID())
	}
	return nil
}

//LocateData on the ring
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

//AddData alias LocateData
func (c *Consistent) AddData(cID uint64, di balancer.DataItem) error {
	return nil
}

//RemoveData stub implementation
func (c *Consistent) RemoveData(di balancer.DataItem) error {
	return nil
}

//Nodes discovers nodes in the balancer
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

//GetNode returns node by given ID.
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

//Optimize stub implementation
func (c *Consistent) Optimize() error {
	return nil
}

//Reset stub implementation
func (c *Consistent) Reset() error {
	return nil
}
