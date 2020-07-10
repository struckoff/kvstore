package optimizer

import (
	"github.com/pkg/errors"
	balancer "github.com/struckoff/SFCFramework"
	"math"
	"sort"
)

func StaticCapacityStateful(s *balancer.Space) (res []*balancer.CellGroup, err error) {
	//TODO: reduce Capacity calls

	cells := s.Cells()
	totalPower := s.TotalPower()
	//totalFree := s.TotalCapacity() - float64(s.TotalLoad())
	cgs := s.CellGroups()
	if len(cgs) == 0 {
		return res, nil
	}
	var max, min uint64

	caps := make([]float64, len(cgs))
	for iter := range cgs {
		caps[iter], err = cgs[iter].Node().Capacity().Get()
		if err != nil {
			return nil, err
		}
	}

	sort.Slice(cgs, func(i, j int) bool {
		capI, _ := cgs[i].Node().Capacity().Get()
		capJ, _ := cgs[j].Node().Capacity().Get()
		return (capI - float64(cgs[i].TotalLoad())) < (capJ - float64(cgs[j].TotalLoad()))
	})

	for iter := 0; iter < len(cgs); iter++ {
		min = max
		p := cgs[iter].Node().Power().Get() / totalPower
		f, err := cgs[iter].Node().Capacity().Get()
		if err != nil {
			return nil, err
		}
		max = min + uint64(math.Round(float64(s.Capacity())*p))

		for citer := 0; citer < len(cells); citer++ {
			if cells[citer].ID() > max {
				break
			}
			if cells[citer].ID() >= min {
				f -= float64(cells[citer].Load())
				if f <= 0 {
					c := citer - 1
					if c < 0 {
						c = 0
					}
					max = cells[citer].ID()
					break
				}
				cgs[iter].AddCell(cells[citer], true)
			}
		}
		if err := cgs[iter].SetRange(min, max); err != nil {
			return nil, errors.Wrap(err, "power range optimizer error")
		}
	}

	if max < s.Capacity() {
		if err := cgs[len(cgs)-1].SetRange(min, s.Capacity()+1); err != nil {
			return nil, errors.Wrap(err, "range optimizer error")
		}
		for citer := range cells {
			if cells[citer].ID() >= max {
				cgs[len(cgs)-1].AddCell(cells[citer], true)
			}
		}
	}

	//for _, cg := range cgs {
	//	log.Print(cg.ID(), ":", cg.TotalLoad(), cg.Range().Min, cg.Range().Max, cg.Range().Len, len(cells))
	//}

	return cgs, nil
}

func DynamicCapacityStateful(s *balancer.Space) (res []*balancer.CellGroup, err error) {
	//TODO: reduce Capacity calls

	cells := s.Cells()
	totalPower := s.TotalPower()
	//totalFree := s.TotalCapacity() - float64(s.TotalLoad())
	cgs := s.CellGroups()
	if len(cgs) == 0 {
		return res, nil
	}
	var max, min uint64

	caps := make([]float64, len(cgs))
	for iter := range cgs {
		caps[iter], err = cgs[iter].Node().Capacity().Get()
		if err != nil {
			return nil, err
		}
	}

	sort.Slice(cgs, func(i, j int) bool {
		capI, err := cgs[i].Node().Capacity().Get()
		if err != nil {
			return false
		}
		capJ, err := cgs[j].Node().Capacity().Get()
		if err != nil {
			return false
		}
		return capI < capJ
	})

	for iter := 0; iter < len(cgs); iter++ {
		min = max
		p := cgs[iter].Node().Power().Get() / totalPower
		f, err := cgs[iter].Node().Capacity().Get()
		if err != nil {
			return nil, err
		}
		max = min + uint64(math.Round(float64(s.Capacity())*p))

		for citer := 0; citer < len(cells); citer++ {
			if cells[citer].ID() > max {
				break
			}
			if cells[citer].ID() >= min {
				f -= float64(cells[citer].Load())
				if f <= 0 {
					c := citer - 1
					if c < 0 {
						c = 0
					}
					max = cells[citer].ID()
					break
				}
				cgs[iter].AddCell(cells[citer], true)
			}
		}
		if err := cgs[iter].SetRange(min, max); err != nil {
			return nil, errors.Wrap(err, "power range optimizer error")
		}
	}

	if max < s.Capacity() {
		if err := cgs[len(cgs)-1].SetRange(min, s.Capacity()+1); err != nil {
			return nil, errors.Wrap(err, "range optimizer error")
		}
		for citer := range cells {
			if cells[citer].ID() >= max {
				cgs[len(cgs)-1].AddCell(cells[citer], true)
			}
		}
	}

	return cgs, nil
}
