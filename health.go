package kvstore

import (
	"github.com/struckoff/kvstore/rpcapi"
	"golang.org/x/net/context"
)

type HealthCheck int

func (h HealthCheck) Check(context.Context, *rpcapi.HealthCheckRequest) (*rpcapi.HealthCheckResponse, error) {
	return nil, nil
}

func (h HealthCheck) Watch(*rpcapi.HealthCheckRequest, rpcapi.Health_WatchServer) error {
	return nil
}
