package balanceradapter

import (
	"github.com/pkg/errors"
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/SFCFramework/curve"
	balancertransform "github.com/struckoff/SFCFramework/transform"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/nodes"
	"github.com/struckoff/kvstore/router/optimizer"
	"github.com/struckoff/kvstore/router/transform"
	"log"
	"math"
	"sync"
)

type SFC struct {
	bal *balancer.Balancer
	m   sync.RWMutex
}

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

func (sb *SFC) AddNode(n nodes.Node) error {
	//sb.m.Lock()
	//defer sb.m.Unlock()
	//if len(sb.bal.Space().CellGroups()) > 0 {
	//	return sb.bal.AddNode(n, false)
	//}
	return sb.bal.AddNode(n, true)

}

func (sb *SFC) RemoveNode(id string) error {
	//sb.m.Lock()
	//defer sb.m.Unlock()
	return sb.bal.RemoveNode(id, true)
}

func (sb *SFC) SetNodes(ns []nodes.Node) error {
	//sb.m.Lock()
	//defer sb.m.Unlock()
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

func (sb *SFC) LocateData(di balancer.DataItem) (nodes.Node, uint64, error) {
	//sb.m.RLock()
	//defer sb.m.RUnlock()
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

func (sb *SFC) AddData(di balancer.DataItem) (nodes.Node, uint64, error) {
	//sb.m.Lock()
	//defer sb.m.Unlock()
	nb, cid, err := sb.bal.LocateData(di)
	if err != nil {
		return nil, 0, err
	}

	ok, err := sb.checkNodeCapacity(nb, di)
	if err != nil {
		return nil, 0, err
	}
	if !ok {
		ncID, err := sb.findBetterCell(di, cid)
		if err != nil {
			return nil, 0, err
		}
		nb, _, err = sb.bal.RelocateData(di, ncID)
		if err != nil {
			return nil, 0, err
		}
	} else {
		nb, cid, err = sb.bal.AddData(di)
		if err != nil {
			return nil, 0, err
		}
	}
	n, ok := nb.(nodes.Node)
	if !ok {
		return nil, 0, errors.New("wrong node type")
	}

	return n, cid, nil
}

func (sb *SFC) RemoveData(di balancer.DataItem) error {
	//sb.m.Lock()
	//defer sb.m.Unlock()
	return sb.bal.RemoveData(di)
}

func (sb *SFC) checkNodeCapacity(n balancer.Node, di balancer.DataItem) (bool, error) {
	cgs := sb.bal.Space().CellGroups()
	c, err := n.Capacity().Get()
	if err != nil {
		return false, err
	}
	nf := true
	for iter := range cgs {
		if cgs[iter].Node().ID() == n.ID() {
			nf = false
			diff := c - float64(cgs[iter].TotalLoad()) - float64(di.Size())
			//diff := c - float64(di.Size())
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

func (sb *SFC) findBetterCell(di balancer.DataItem, cid uint64) (uint64, error) {
	dis := math.MaxInt64
	ncID := cid
	//var cg *balancer.CellGroup

	cgs := sb.bal.Space().CellGroups()
	for iter := range cgs {
		l := cgs[iter].TotalLoad()
		c, err := cgs[iter].Node().Capacity().Get()
		if err != nil {
			log.Println(err)
			continue
		}
		dc := c - float64(l) - float64(di.Size())
		//dc := c - float64(di.Size())
		if dc < 0 {
			continue
		}
		if cgs[iter].Range().Len <= 0 {
			continue
		}

		// find closest cell to filled group
		if cgs[iter].Range().Max <= cid {
			if lft := int(cid) - int(cgs[iter].Range().Max); lft < dis {
				dis = lft
				ncID = cgs[iter].Range().Max - 1 //closest cell in available group
			}
		} else if cgs[iter].Range().Min > cid {
			if rght := int(cgs[iter].Range().Min) - int(cid); rght < dis {
				dis = rght
				ncID = cgs[iter].Range().Min //closest cell in available group
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

func (sb *SFC) Nodes() ([]nodes.Node, error) {
	//sb.m.RLock()
	//defer sb.m.RUnlock()
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
	//sb.m.RLock()
	//defer sb.m.RUnlock()
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
	//sb.m.RLock()
	//defer sb.m.RUnlock()
	return sb.bal.SFC()
}

func (sb *SFC) Optimize() error {
	//sb.m.Lock()
	//defer sb.m.Unlock()
	return sb.bal.Optimize()
}

func (sb *SFC) Reset() error {
	//sb.m.Lock()
	//defer sb.m.Unlock()
	for _, cg := range sb.bal.Space().CellGroups() {
		cg.Truncate()
	}
	return nil
}
