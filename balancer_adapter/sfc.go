package balancer_adapter

import (
	"github.com/pkg/errors"
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/SFCFramework/optimizer"
	"github.com/struckoff/SFCFramework/transform"
	"github.com/struckoff/kvstore"
)

// DataItem represents string key as balancer item
type DataItem string

func (di DataItem) ID() string   { return string(di) }
func (di DataItem) Size() uint64 { return 1 }
func (di DataItem) Values() []interface{} {
	vals := make([]interface{}, 1)
	vals[0] = string(di)
	return vals
}

type SFCBalancer struct {
	bal *balancer.Balancer
}

func NewSFCBalancer(conf kvstore.Config) (*SFCBalancer, error) {
	bal, err := balancer.NewBalancer(
		conf.Curve.CurveType,
		conf.Dimensions,
		conf.Size,
		transform.KVTransform,
		optimizer.RangeOptimizer,
		nil)
	if err != nil {
		return nil, err
	}
	return &SFCBalancer{bal: bal}, nil
}

func (sb *SFCBalancer) AddNode(n kvstore.Node) error {
	return sb.bal.AddNode(n)
}

func (sb *SFCBalancer) RemoveNode(id string) error {
	return sb.bal.RemoveNode(id)
}

func (sb *SFCBalancer) LocateKey(key string) (kvstore.Node, error) {
	di := DataItem(key)
	nb, err := sb.bal.LocateData(di)
	if err != nil {
		return nil, err
	}
	n, ok := nb.(kvstore.Node)
	if !ok {
		return nil, errors.New("wrong node type")
	}
	return n, nil
}

func (sb *SFCBalancer) Nodes() ([]kvstore.Node, error) {
	nbs := sb.bal.Nodes()
	ns := make([]kvstore.Node, len(nbs))
	for iter := range nbs {
		n, ok := nbs[iter].(kvstore.Node)
		if !ok {
			return nil, errors.New("wrong node type")
		}
		ns[iter] = n
	}
	return ns, nil
}
