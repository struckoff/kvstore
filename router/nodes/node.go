package nodes

import (
	"context"
	"github.com/struckoff/kvstore/router/rpcapi"
	balancernode "github.com/struckoff/sfcframework/node"
)

type Node interface {
	balancernode.Node
	Capacity() Capacity
	Store(context.Context, *rpcapi.KeyValue) (*rpcapi.DataItem, error)          // Save value for a given key
	StorePairs(context.Context, []*rpcapi.KeyValue) ([]*rpcapi.DataItem, error) // Save multiple key-value pairs
	Receive(context.Context, []*rpcapi.DataItem) (*rpcapi.KeyValues, error)     // Return value for a given key
	Remove(context.Context, []*rpcapi.DataItem) ([]*rpcapi.DataItem, error)     // Remove value for a given key
	Explore(context.Context) ([]*rpcapi.DataItem, error)                        // Return all keys in a cluster
	Meta(context.Context) *rpcapi.NodeMeta                                      // Return information about cluster units
	Move(context.Context, map[Node][]*rpcapi.DataItem) error                    // Move kv pairs to another node
}
