package nodes

import (
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/kvstore/router/rpcapi"
)

type Node interface {
	balancer.Node
	Store(key string, body []byte) error // Save value for a given key
	StorePairs([]*rpcapi.KeyValue) error // Save multiple key-value pairs
	Receive(key string) ([]byte, error)  // Return value for a given key
	Remove(key string) error             // Remove value for a given key
	Explore() ([]string, error)          // Return all keys in a cluster
	Meta() *rpcapi.NodeMeta              // Return information about cluster units
	Move(map[Node][]string) error        // Move kv pairs to another node
}
