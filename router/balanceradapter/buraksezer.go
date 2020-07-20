package balanceradapter

import (
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"github.com/pkg/errors"
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/kvstore/router/nodes"
)

type BuraksezerRing struct {
	conf consistent.Config
	ring *consistent.Consistent
}

func NewBuraksezerRingBalancer(conf *consistent.Config) *BuraksezerRing {
	conf.Hasher = hasher{}
	cr := &BuraksezerRing{
		conf: *conf,
		ring: consistent.New(nil, *conf),
	}
	return cr
}

func (c BuraksezerRing) AddNode(n nodes.Node) error {
	rn := RingNode{n}
	if c.ring == nil {
		c.ring = consistent.New([]consistent.Member{rn, rn}, c.conf)
		return nil
	}
	c.ring.Add(rn)
	return nil
}

func (c BuraksezerRing) RemoveNode(id string) error {
	c.ring.Remove(id)
	return nil
}

type myMember string

func (m myMember) String() string {
	return string(m)
}

func (c BuraksezerRing) SetNodes(ns []nodes.Node) error {
	ms := make([]consistent.Member, len(ns))
	for iter := range ns {
		ms[iter] = RingNode{ns[iter]}
	}
	//ms = append(ms, RingNode{ns[0]})

	c.ring = consistent.New(nil, c.conf)
	c.ring.Add(RingNode{ns[0]})
	return nil
}

func (c BuraksezerRing) LocateData(di balancer.DataItem) (nodes.Node, error) {
	m := c.ring.LocateKey([]byte(di.ID()))
	if rn, ok := m.(RingNode); ok {
		return rn, nil
	}
	return nil, errors.New("wrong node type")
}

func (c BuraksezerRing) AddData(di balancer.DataItem) (nodes.Node, error) {
	return c.LocateData(di)
}

func (c BuraksezerRing) Nodes() ([]nodes.Node, error) {
	var ok bool
	ms := c.ring.GetMembers()
	ns := make([]nodes.Node, len(ms))
	for iter := range ms {
		if ns[iter], ok = ms[iter].(RingNode); !ok {
			return nil, errors.New("wrong node type")
		}
	}
	return ns, nil
}

func (c BuraksezerRing) GetNode(id string) (nodes.Node, error) {
	ms := c.ring.GetMembers()
	for _, m := range ms {
		if m.String() == id {
			if rn, ok := m.(RingNode); ok {
				return rn, nil
			}
			return nil, errors.New("wrong node type")
		}
	}
	return nil, errors.New("node not found")
}

// consistent package doesn't provide a default hashing function.
// You should provide a proper one to distribute keys/members uniformly.
type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	// you should use a proper hash function for uniformity.
	return xxhash.Sum64(data)
}

type RingNode struct {
	nodes.Node
}

func (rn RingNode) String() string {
	return rn.ID()
}
