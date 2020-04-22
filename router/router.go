package router

import (
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/nodehasher"
)

// Router represents bounding of network api with kvrouter lib and local node
type Router struct {
	bal    Balancer
	hasher nodehasher.Hasher
}

func NewRouter(bal Balancer, h nodehasher.Hasher) (*Router, error) {
	r := &Router{
		bal:    bal,
		hasher: h,
	}
	return r, nil
}

// AddNode adds node to kvrouter
func (h *Router) AddNode(n Node) error {
	return h.bal.AddNode(n)
}

// RemoveNode removes node from kvrouter
func (h *Router) RemoveNode(id string) error {
	return h.bal.RemoveNode(id)
}

// Returns node from kvrouter by given key.
func (h *Router) LocateKey(key string) (Node, error) {
	//di := DataItem(key)
	nb, err := h.bal.LocateKey(key)
	if err != nil {
		return nil, err
	}
	n, ok := nb.(Node)
	if !ok {
		return nil, errors.New("wrong node type")
	}
	return n, nil
}

// GetNodes - returns a list of nodes in the balancer
func (h *Router) GetNodes() ([]Node, error) {
	return h.bal.Nodes()
}

// SetNodes - removes all nodes from the balancer and set a new ones
func (h *Router) SetNodes(ns []Node) error {
	return h.bal.SetNodes(ns)
}

// GetNode returns the node with the given ID
func (h *Router) GetNode(id string) (Node, error) {
	return h.bal.GetNode(id)
}

func (h *Router) Hasher() nodehasher.Hasher {
	return h.hasher
}
