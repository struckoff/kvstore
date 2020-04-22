package router

import (
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/kvstore/router/rpcapi"
	"sort"
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

type mockNode struct {
	rpcapi.NodeMeta
	id string
	kv map[string][]byte
}

func (m mockNode) ID() string {
	return m.id
}

func (m mockNode) Power() balancer.Power {
	return Power{m.NodeMeta.Power}
}

func (m mockNode) Capacity() balancer.Capacity {
	return Power{m.NodeMeta.Capacity}
}

func (m mockNode) Hash() uint64 {
	panic("implement me")
}

func (m mockNode) Store(key string, body []byte) error {
	m.kv[key] = body
	return nil
}

func (m mockNode) StorePairs(values []*rpcapi.KeyValue) error {
	panic("implement me")
}

func (m mockNode) Receive(key string) ([]byte, error) {
	return m.kv[key], nil
}

func (m mockNode) Remove(key string) error {
	delete(m.kv, key)
	return nil
}

func (m mockNode) Explore() ([]string, error) {
	keys := make([]string, 0, len(m.kv))
	for key := range m.kv {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys, nil
}

func (m mockNode) Meta() *rpcapi.NodeMeta {
	return &rpcapi.NodeMeta{
		ID:         m.id,
		Address:    m.Address,
		RPCAddress: m.RPCAddress,
		Power:      m.Power().Get(),
		Capacity:   m.Capacity().Get(),
		Check:      m.Check,
		Geo:        m.Geo,
	}
}

func (m mockNode) Move(m2 map[Node][]string) error {
	panic("implement me")
}
