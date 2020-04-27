package router

import (
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/balanceradapter"
	"github.com/struckoff/kvstore/router/dataitem"
	"github.com/struckoff/kvstore/router/nodehasher"
	"github.com/struckoff/kvstore/router/nodes"
)

// Router represents bounding of network api with kvrouter lib and local node
type Router struct {
	bal    balanceradapter.Balancer
	hasher nodehasher.Hasher
	ndf    dataitem.NewDataItemFunc
}

func NewRouter(bal balanceradapter.Balancer, h nodehasher.Hasher, ndf dataitem.NewDataItemFunc) (*Router, error) {
	r := &Router{
		bal:    bal,
		hasher: h,
		ndf:    ndf,
	}
	return r, nil
}

// AddNode adds node to kvrouter
func (h *Router) AddNode(n nodes.Node) error {
	return h.bal.AddNode(n)
}

// RemoveNode removes node from kvrouter
func (h *Router) RemoveNode(id string) error {
	return h.bal.RemoveNode(id)
}

// Returns node from kvrouter by given key.
func (h *Router) LocateKey(key string) (nodes.Node, error) {
	//di := DataItem(key)
	di, err := h.ndf(key)
	if err != nil {
		return nil, err
	}
	nb, err := h.bal.LocateData(di)
	if err != nil {
		return nil, err
	}
	n, ok := nb.(nodes.Node)
	if !ok {
		return nil, errors.New("wrong node type")
	}
	return n, nil
}

// GetNodes - returns a list of nodes in the balancer
func (h *Router) GetNodes() ([]nodes.Node, error) {
	return h.bal.Nodes()
}

// SetNodes - removes all nodes from the balancer and set a new ones
func (h *Router) SetNodes(ns []nodes.Node) error {
	return h.bal.SetNodes(ns)
}

// GetNode returns the node with the given ID
func (h *Router) GetNode(id string) (nodes.Node, error) {
	return h.bal.GetNode(id)
}

func (h *Router) Hasher() nodehasher.Hasher {
	return h.hasher
}
