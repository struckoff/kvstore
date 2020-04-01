package kvstore

import (
	balancer "github.com/struckoff/SFCFramework"
)

type Node interface {
	balancer.Node
	Store(key string, body []byte) error // Save value for a given key
	Receive(key string) ([]byte, error)  // Return value for a given key
	Explore() ([]string, error)          // Return all keys in a cluster
	Meta() NodeMeta                      // Return information about cluster units
}

// NodeMeta represents node meta information with exposed fields
// for marshaling and unmarshaling
type NodeMeta struct {
	ID         string
	Address    string
	RPCAddress string
	Power      float64
	Capacity   float64
}
