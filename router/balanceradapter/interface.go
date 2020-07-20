package balanceradapter

import (
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/kvstore/router/nodes"
)

type Balancer interface {
	AddNode(n nodes.Node) error                          // Add node to the balancer
	RemoveNode(id string) error                          // Remove node from the balancer
	SetNodes(ns []nodes.Node) error                      // Remove all nodes from the balancer and set a new ones
	LocateData(di balancer.DataItem) (nodes.Node, error) // Return the node for the given key
	AddData(di balancer.DataItem) (nodes.Node, error)    // Return the node for the given key
	Nodes() ([]nodes.Node, error)                        // Return list of nodes in the balancer
	GetNode(id string) (nodes.Node, error)               // Return the node with the given id
	Optimize() error                                     //Force redistribution from an outside
}
