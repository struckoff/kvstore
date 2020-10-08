package nodes

import (
	"github.com/struckoff/kvstore/router/rpcapi"
	balancernode "github.com/struckoff/sfcframework/node"
)

type Node interface {
	balancernode.Node
	Capacity() Capacity
	Store(*rpcapi.KeyValue) (*rpcapi.DataItem, error)          // Save value for a given key
	StorePairs([]*rpcapi.KeyValue) ([]*rpcapi.DataItem, error) // Save multiple key-value pairs
	Receive([]*rpcapi.DataItem) (*rpcapi.KeyValues, error)     // Return value for a given key
	Remove([]*rpcapi.DataItem) ([]*rpcapi.DataItem, error)     // Remove value for a given key
	Explore() ([]*rpcapi.DataItem, error)                      // Return all keys in a cluster
	Meta() *rpcapi.NodeMeta                                    // Return information about cluster units
	Move(map[Node][]*rpcapi.DataItem) error                    // Move kv pairs to another node
}
