package optimizer

import (
	balancer "github.com/struckoff/sfcframework"
	balanceroptimizer "github.com/struckoff/sfcframework/optimizer"
	"log"
)

func RangeOptimizer(s *balancer.Space) (res []*balancer.CellGroup, err error) {
	defer func() {
		cgs := s.CellGroups()
		for _, cg := range cgs {
			log.Println("RangeOptimizer", cg.Node().ID(), ": ", cg.Range().Min, cg.Range().Max, cg.Range().Len, cg.TotalLoad(), len(cg.Cells()))
		}
	}()
	return balanceroptimizer.RangeOptimizer(s)
}
