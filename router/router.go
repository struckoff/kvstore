package router

import (
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/balanceradapter"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/dataitem"
	"github.com/struckoff/kvstore/router/nodehasher"
	"github.com/struckoff/kvstore/router/nodes"
	"log"
	"sync"
)

// Router represents bounding of network api with kvrouter lib and local node
type Router struct {
	bal    balanceradapter.Balancer
	hasher nodehasher.Hasher
	ndf    dataitem.NewDataItemFunc
	conf   *config.BalancerConfig
}

func NewRouter(bal balanceradapter.Balancer, h nodehasher.Hasher, ndf dataitem.NewDataItemFunc, conf *config.BalancerConfig) (*Router, error) {
	r := &Router{
		bal:    bal,
		hasher: h,
		ndf:    ndf,
		conf:   conf,
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

func (h *Router) AddData(key string) (nodes.Node, error) {
	//di := DataItem(key)
	di, err := h.ndf(key)
	if err != nil {
		return nil, err
	}
	nb, err := h.bal.AddData(di)
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

func (h *Router) redistributeKeys() error {
	var wg sync.WaitGroup
	ns, err := h.GetNodes()
	if err != nil {
		return err
	}
	for _, n := range ns {
		go func(n nodes.Node, wg *sync.WaitGroup) {
			res := make(map[nodes.Node][]string)
			keys, err := n.Explore()
			if err != nil {
				log.Printf("failed to explore node(%s): %s", n.ID(), err.Error())
				return
			}
			for iter := range keys {
				en, err := h.LocateKey(keys[iter])
				if err != nil {
					log.Printf("failed to locate key(%s): %s", keys[iter], err.Error())
					continue
				}
				if en.ID() != n.ID() {
					res[en] = append(res[en], keys[iter])
				}
			}
			n.Move(res)
		}(n, &wg)
	}
	wg.Wait()
	return nil
}
