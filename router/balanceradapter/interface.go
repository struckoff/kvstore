package balanceradapter

import (
	"github.com/struckoff/kvstore/router/nodes"
	balancer "github.com/struckoff/sfcframework"
)

type Balancer interface {
	AddNode(n nodes.Node) error                                  // Add node to the balancer
	RemoveNode(id string) error                                  // Remove node from the balancer
	SetNodes(ns []nodes.Node) error                              // Remove all nodes from the balancer and set a new ones
	LocateData(di balancer.DataItem) (nodes.Node, uint64, error) // Return the node for the given key
	AddData(di balancer.DataItem) (nodes.Node, uint64, error)    // Return the node for the given key
	RemoveData(di balancer.DataItem) error
	Nodes() ([]nodes.Node, error)          // Return list of nodes in the balancer
	GetNode(id string) (nodes.Node, error) // Return the node with the given id
	Optimize() error                       //Force redistribution from an outside
	Reset() error                          //Force balancer load reset
}
