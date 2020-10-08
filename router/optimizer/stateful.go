package optimizer

import (
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/nodes"
	balancer "github.com/struckoff/sfcframework"
	"math"
	"sort"
)

func CapPowerOptimizer(s *balancer.Space) (cgs []*balancer.CellGroup, err error) {
	//defer func() {
	//	//cgs := s.CellGroups()
	//	for _, cg := range cgs {
	//		log.Println("CapPowerOptimizer",
	//			cg.Node().ID(),
	//			cg.Node().Hash(), ": ",
	//			cg.Range().Min,
	//			cg.Range().Max,
	//			cg.Range().Len,
	//			cg.TotalLoad(),
	//			len(cg.Cells()),
	//			s.TotalLoad(),
	//		)
	//	}
	//}()

	cgs, err = powerRangeOptimizer(s)
	if err != nil {
		return nil, err
	}
	sort.Slice(cgs, func(i, j int) bool { return cgs[i].Node().Hash() < cgs[j].Node().Hash() })
	if len(cgs) >= 2 {
		for cgIdx := 0; cgIdx < len(cgs)-1; cgIdx++ {
			if err := equalizer(cgs[cgIdx], cgs[cgIdx+1], s); err != nil {
				return nil, err
			}
		}
	}

	for i := range cgs {
		n, ok := cgs[i].Node().(nodes.Node)
		if !ok {
			return nil, errors.New("out of capacity")
		}
		c, err := n.Capacity().Get()
		if err != nil {
			return nil, err
		}
		if float64(cgs[i].TotalLoad()) > c {
			return nil, errors.New("out of capacity")
		}
	}

	return cgs, nil
}

func powerRangeOptimizer(s *balancer.Space) (res []*balancer.CellGroup, err error) {
	totalPower := s.TotalPower()
	cgs := s.CellGroups()
	if len(cgs) == 0 {
		return res, nil
	}
	var max, min uint64

	sort.Slice(cgs, func(i, j int) bool { return cgs[i].Node().Hash() < cgs[j].Node().Hash() })

	for i := 0; i < len(cgs); i++ {
		min = max
		p := cgs[i].Node().Power().Get() / totalPower
		max = min + uint64(math.Ceil(float64(s.Capacity())*p))
		if max > s.Capacity()+1 {
			max = s.Capacity() + 1
		}
		if err := cgs[i].SetRange(min, max); err != nil {
			return nil, errors.Wrap(err, "range power optimizer error")
		}
		s.FillCellGroup(cgs[i])
	}
	if max < s.Capacity() {
		if err := cgs[len(cgs)-1].SetRange(min, s.Capacity()+1); err != nil {
			return nil, errors.Wrap(err, "range power optimizer error")
		}
		s.FillCellGroup(cgs[len(cgs)-1])
	}
	return cgs, nil
}

func equalizer(lcg, rcg *balancer.CellGroup, s *balancer.Space) error {
	if lcg.Range().Max != rcg.Range().Min {
		return errors.New("wrong group pair")
	}

	n, ok := lcg.Node().(nodes.Node)
	if !ok {
		return errors.New("unable cast node to nodes.Node")
	}
	lc, err := n.Capacity().Get()
	if err != nil {
		return err
	}
	nbf := true
	for float64(lcg.TotalLoad()) > lc {
		nbf = false
		if err := lcg.SetRange(lcg.Range().Min, lcg.Range().Max-1); err != nil {
			return err
		}
		s.FillCellGroup(lcg)
		if err := rcg.SetRange(rcg.Range().Min-1, rcg.Range().Max); err != nil {
			return err
		}
		s.FillCellGroup(rcg)
	}
	if nbf {
		n, ok := lcg.Node().(nodes.Node)
		if !ok {
			return errors.New("unable cast node to nodes.Node")
		}
		rc, err := n.Capacity().Get()
		if err != nil {
			return err
		}
		//for float64(rcg.TotalLoad()) > rc && float64(lcg.TotalLoad()) <= lc {
		for float64(rcg.TotalLoad()) > rc {
			if err := lcg.SetRange(lcg.Range().Min, lcg.Range().Max+1); err != nil {
				return err
			}
			s.FillCellGroup(lcg)
			if err := rcg.SetRange(rcg.Range().Min+1, rcg.Range().Max); err != nil {
				return err
			}
			s.FillCellGroup(rcg)
		}
	}
	return nil
}
