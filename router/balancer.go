package router

import (
	"errors"

	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/SFCFramework/optimizer"
	"github.com/struckoff/SFCFramework/transform"
)

type Balancer interface {
	AddNode(n Node) error               // Add node to the balancer
	RemoveNode(id string) error         // Remove node from the balancer
	SetNodes(ns []Node) error           // Remove all nodes from the balancer and set a new ones
	LocateKey(key string) (Node, error) // Return the node for the given key
	Nodes() ([]Node, error)             // Return list of nodes in the balancer
	GetNode(id string) (Node, error)    // Return the node with the given id
}

type SFCBalancer struct {
	bal *balancer.Balancer
}

func NewSFCBalancer(conf *BalancerConfig) (*SFCBalancer, error) {
	bal, err := balancer.NewBalancer(
		conf.Curve.CurveType,
		conf.Dimensions,
		conf.Size,
		transform.SpaceTransform,
		optimizer.RangeOptimizer,
		nil)
	if err != nil {
		return nil, err
	}
	return &SFCBalancer{bal: bal}, nil
}

func (sb *SFCBalancer) AddNode(n Node) error {
	return sb.bal.AddNode(n, true)
}

func (sb *SFCBalancer) RemoveNode(id string) error {
	return sb.bal.RemoveNode(id)
}

func (sb *SFCBalancer) SetNodes(ns []Node) error {
	sb.bal.Space().SetGroups(nil)
	for _, node := range ns {
		if err := sb.bal.AddNode(node, false); err != nil {
			return err
		}
	}
	if err := sb.bal.Optimize(); err != nil {
		return err
	}
	return nil
}

func (sb *SFCBalancer) LocateKey(key string) (Node, error) {
	di, err := NewSpaceDataItem(key)
	if err != nil {
		return nil, err
	}
	nb, err := sb.bal.LocateData(di)
	if err != nil {
		return nil, err
	}
	n, ok := nb.(Node)
	if !ok {
		return nil, errors.New("wrong node type")
	}
	return n, nil
}

func (sb *SFCBalancer) Nodes() ([]Node, error) {
	nbs := sb.bal.Nodes()
	ns := make([]Node, len(nbs))
	for iter := range nbs {
		n, ok := nbs[iter].(Node)
		if !ok {
			return nil, errors.New("wrong node type")
		}
		ns[iter] = n
	}
	return ns, nil
}

func (sb *SFCBalancer) GetNode(id string) (Node, error) {
	nb, ok := sb.bal.GetNode(id)
	if !ok {
		return nil, errors.New("node not found")
	}
	n, ok := nb.(Node)
	if !ok {
		return nil, errors.New("wrong node type")
	}
	return n, nil
}
