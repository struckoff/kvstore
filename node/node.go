package node

import (
	balancer "github.com/struckoff/SFCFramework"
	"io"
)

type Node interface {
	ID() string
	Power() balancer.Power
	Capacity() balancer.Capacity
	Store(key string, body io.Reader) error // Save value for a given key
	Receive(key string) ([]byte, error)     // Return value for a given key
	Explore() ([]byte, error)               // Return all keys in a cluster
	Meta() NodeMeta                         //Return information about cluster units
}

// NodeMeta represents node meta information with exposed fields
// for marshaling and unmarshaling
type NodeMeta struct {
	ID       string
	Address  string
	Power    float64
	Capacity float64
}
