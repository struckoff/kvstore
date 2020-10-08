package balanceradapter

import (
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/nodes"
	"github.com/struckoff/kvstore/router/optimizer"
	"github.com/struckoff/kvstore/router/transform"
	balancer "github.com/struckoff/sfcframework"
	"github.com/struckoff/sfcframework/curve"
	balancertransform "github.com/struckoff/sfcframework/transform"
	"math"
)

//SFC - Space-filling curve adapter
type SFC struct {
	bal *balancer.Balancer
}

//NewSFCBalancer create new space-filling curve powered balancer adapter instance
func NewSFCBalancer(conf *config.BalancerConfig) (*SFC, error) {
	var tf balancer.TransformFunc
	var of balancer.OptimizerFunc
	switch conf.DataMode {
	case config.GeoData:
		tf = transform.SpaceTransform
	case config.KVData:
		tf = balancertransform.KVTransform
	default:
		return nil, errors.New("wrong data mode")
	}
	if conf.State {
		of = optimizer.CapPowerOptimizer
	} else {
		of = optimizer.RangeOptimizer
	}
	bal, err := balancer.NewBalancer(
		conf.SFC.Curve.CurveType,
		conf.SFC.Dimensions,
		conf.SFC.Size,
		tf,
		of,
		nil)
	if err != nil {
		return nil, err
	}

	return &SFC{bal: bal}, nil
}

//AddNode to the balancer
func (sb *SFC) AddNode(n nodes.Node) error {
	return sb.bal.AddNode(n, true)

}

//RemoveNode from the balancer
func (sb *SFC) RemoveNode(id string) error {
	return sb.bal.RemoveNode(id, true)
}

//SetNodes wipes all nodes from the balancer and set a new ones
func (sb *SFC) SetNodes(ns []nodes.Node) error {
	for _, cell := range sb.bal.Space().Cells() {
		cell.Truncate()
	}
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

//LocateData on the space-filling curve
func (sb *SFC) LocateData(di balancer.DataItem) (nodes.Node, uint64, error) {
	nb, cid, err := sb.bal.LocateData(di)
	if err != nil {
		return nil, 0, err
	}
	n, ok := nb.(nodes.Node)
	if !ok {
		return nil, 0, errors.New("wrong node type")
	}
	return n, cid, nil
}

//AddData adds data to the balancer considering capacity of nodes.
func (sb *SFC) AddData(cID uint64, di balancer.DataItem) error {
	err := sb.bal.AddData(cID, di)
	if err != nil {
		return err
	}
	//
	//ok, err := sb.checkNodeCapacity(nb, di)
	//if err != nil {
	//	return nil, 0, err
	//}
	//if !ok {
	//	ncID, err := sb.findBetterCell(di, cid)
	//	if err != nil {
	//		return nil, 0, err
	//	}
	//	nb, _, err = sb.bal.RelocateData(di, ncID)
	//	if err != nil {
	//		return nil, 0, err
	//	}
	//} else {
	//	nb, cid, err = sb.bal.AddData(di)
	//	if err != nil {
	//		return nil, 0, err
	//	}
	//}
	//n, ok := nb.(nodes.Node)
	//if !ok {
	//	return nil, 0, errors.New("wrong node type")
	//}

	return nil
}

//RemoveData from the balancer.
func (sb *SFC) RemoveData(di balancer.DataItem) error {
	return sb.bal.RemoveData(di)
}

//checkNodeCapacity - returns true if the node is able to receive data.
func (sb *SFC) checkNodeCapacity(n nodes.Node, di balancer.DataItem) (bool, error) {
	cgs := sb.bal.Space().CellGroups()
	c, err := n.Capacity().Get()
	if err != nil {
		return false, err
	}
	nf := true
	for i := range cgs {
		if cgs[i].Node().ID() == n.ID() {
			nf = false
			diff := c - float64(cgs[i].TotalLoad()) - float64(di.Size())
			if diff >= 0 {
				return true, err
			}
			return false, nil
		}
	}
	if nf {
		return false, errors.New("cell group not found")
	}
	return false, nil
}

//findBetterCell returns the nearest cell id which ois able to receive data.
func (sb *SFC) findBetterCell(di balancer.DataItem, cid uint64) (uint64, error) {
	dis := math.MaxInt64
	ncID := cid

	cgs := sb.bal.Space().CellGroups()
	for i := range cgs {
		l := cgs[i].TotalLoad()
		n, ok := cgs[i].Node().(nodes.Node)
		if !ok {
			return 0, errors.New("unable to cast node")
		}
		c, err := n.Capacity().Get()
		if err != nil {
			continue
		}
		dc := c - float64(l) - float64(di.Size())
		if dc < 0 {
			continue
		}
		if cgs[i].Range().Len <= 0 {
			continue
		}

		// find closest cell to filled group
		if cgs[i].Range().Max <= cid {
			if lft := int(cid) - int(cgs[i].Range().Max); lft < dis {
				dis = lft
				ncID = cgs[i].Range().Max - 1 //closest cell in available group
			}
		} else if cgs[i].Range().Min > cid {
			if rght := int(cgs[i].Range().Min) - int(cid); rght < dis {
				dis = rght
				ncID = cgs[i].Range().Min //closest cell in available group
			}
		}
	}

	if dis == math.MaxInt64 {
		return 0, errors.New("out of capacity")
	}
	if cid == ncID {
		return 0, errors.New("appropriate cell not found")
	}

	return ncID, nil
}

//Nodes discovers nodes in the balancer
func (sb *SFC) Nodes() ([]nodes.Node, error) {
	nbs := sb.bal.Nodes()
	ns := make([]nodes.Node, len(nbs))
	for i := range nbs {
		n, ok := nbs[i].(nodes.Node)
		if !ok {
			return nil, errors.New("wrong node type")
		}
		ns[i] = n
	}
	return ns, nil
}

//GetNode returns node by given ID.
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

//SFC returns space-filling curve instance from the balancer.
func (sb *SFC) SFC() curve.Curve {
	return sb.bal.SFC()
}

//Optimize runs the optimization of the balancer space.
func (sb *SFC) Optimize() error {
	return sb.bal.Optimize()
}

//Reset wipes out the balancer.
func (sb *SFC) Reset() error {
	for _, cg := range sb.bal.Space().CellGroups() {
		cg.Truncate()
	}
	return nil
}
