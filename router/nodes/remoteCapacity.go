package nodes

import (
	"context"
	"github.com/struckoff/kvstore/router/rpcapi"
)

type RemoteCapacity struct {
	rc rpcapi.RPCCapacityClient
}

func (c *RemoteCapacity) Get() (float64, error) {
	cp, err := c.rc.RPCGet(context.TODO(), &rpcapi.Empty{})
	if err != nil {
		return 0, err
	}
	return cp.Capacity, nil
}

func NewCapacity(rc rpcapi.RPCCapacityClient) RemoteCapacity {
	return RemoteCapacity{rc: rc}
}
