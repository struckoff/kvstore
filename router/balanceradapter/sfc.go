package balanceradapter

import (
	"github.com/pkg/errors"
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/SFCFramework/curve"
	"github.com/struckoff/SFCFramework/optimizer"
	"github.com/struckoff/SFCFramework/transform"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/nodes"
)

type SFC struct {
	bal *balancer.Balancer
}

func NewSFCBalancer(conf *config.BalancerConfig) (*SFC, error) {
	var tf balancer.TransformFunc
	switch conf.DataMode {
	case config.GeoData:
		tf = transform.SpaceTransform
	case config.KVData:
		tf = transform.KVTransform
	default:
		return nil, errors.New("wrong data mode")
	}
	bal, err := balancer.NewBalancer(
		conf.SFC.Curve.CurveType,
		conf.SFC.Dimensions,
		conf.SFC.Size,
		tf,
		optimizer.RangeOptimizer,
		nil)
	if err != nil {
		return nil, err
	}
	return &SFC{bal: bal}, nil
}

func (sb *SFC) AddNode(n nodes.Node) error {
	return sb.bal.AddNode(n, true)
}

func (sb *SFC) RemoveNode(id string) error {
	return sb.bal.RemoveNode(id)
}

func (sb *SFC) SetNodes(ns []nodes.Node) error {
	sb.bal.Space().SetGroups(nil)
	for _, node := range ns {
		if err := sb.bal.AddNode(node, false); err != nil {
			return err
		}
	}
	if err := sb.bal.Optimize(); err != nil {
		return err
	}
	return nil
}

func (sb *SFC) LocateData(di balancer.DataItem) (nodes.Node, error) {
	nb, err := sb.bal.LocateData(di)
	if err != nil {
		return nil, err
	}
	n, ok := nb.(nodes.Node)
	if !ok {
		return nil, errors.New("wrong node type")
	}
	return n, nil
}

func (sb *SFC) Nodes() ([]nodes.Node, error) {
	nbs := sb.bal.Nodes()
	ns := make([]nodes.Node, len(nbs))
	for iter := range nbs {
		n, ok := nbs[iter].(nodes.Node)
		if !ok {
			return nil, errors.New("wrong node type")
		}
		ns[iter] = n
	}
	return ns, nil
}

func (sb *SFC) GetNode(id string) (nodes.Node, error) {
	nb, ok := sb.bal.GetNode(id)
	if !ok {
		return nil, errors.New("node not found")
	}
	n, ok := nb.(nodes.Node)
	if !ok {
		return nil, errors.New("wrong node type")
	}
	return n, nil
}

func (sb *SFC) SFC() curve.Curve {
	return sb.bal.SFC()
}
