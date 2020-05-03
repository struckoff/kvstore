package nodes

import (
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/kvstore/router/rpcapi"
)

type Node interface {
	balancer.Node
	Store(string, []byte) error                  // Save value for a given key
	StorePairs([]*rpcapi.KeyValue) error         // Save multiple key-value pairs
	Receive([]string) (*rpcapi.KeyValues, error) // Return value for a given key
	Remove([]string) error                       // Remove value for a given key
	Explore() ([]string, error)                  // Return all keys in a cluster
	Meta() *rpcapi.NodeMeta                      // Return information about cluster units
	Move(map[Node][]string) error                // Move kv pairs to another node
}
