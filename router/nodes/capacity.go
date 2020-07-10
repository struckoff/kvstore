package nodes

import (
	"context"
	"github.com/struckoff/kvstore/router/rpcapi"
)

type Capacity struct {
	rc rpcapi.RPCCapacityClient
}

func (c *Capacity) Get() (float64, error) {
	cp, err := c.rc.RPCGet(context.TODO(), &rpcapi.Empty{})
	if err != nil {
		return 0, err
	}
	return cp.Capacity, nil
}

func NewCapacity(rc rpcapi.RPCCapacityClient) Capacity {
	return Capacity{rc: rc}
}
