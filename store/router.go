package store

import (
	"context"
	"github.com/struckoff/kvstore/logger"
	"go.uber.org/zap"
	"time"

	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/rpcapi"
	"google.golang.org/grpc"
)

// RunRouter - resister node in remote KVRouter.
// Run goroutine which sends heartbeat each Config.Health.CheckInterval
func (inn *LocalNode) RunRouter(conf *Config) error {
	logger.Logger().Info("trying to connect to kvrouter", zap.String("kvrouter address", conf.KVRouter.Address))
	c, err := client(conf.KVRouter.Address)
	if err != nil {
		return errors.Wrap(err, "failed to initialize kvrouter client")
	}
	inn.kvrAgent = c
	if err := inn.routerAnnounce(conf); err != nil {
		return errors.Wrap(err, "unable to run announce node in kvrouter")
	}
	logger.Logger().Info("connected to kvrouter", zap.String("kvrouter address", conf.KVRouter.Address))
	return nil
}

// client returns rpc client for kvrouter
func client(addr string) (rpcapi.RPCBalancerClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure()) // TODO Make it secure
	if err != nil {
		return nil, err
	}
	c := rpcapi.NewRPCBalancerClient(conn)
	return c, nil
}

// routerAnnounce - register node in kvrouter
func (inn *LocalNode) routerAnnounce(conf *Config) error {
	checkInterval, err := time.ParseDuration(conf.Health.CheckInterval)
	if err != nil {
		return errors.Wrap(err, "failed to parse health check interval")
	}
	checkTimeout, err := time.ParseDuration(conf.Health.CheckTimeout)
	if err != nil {
		return errors.Wrap(err, "failed to parse health check timeout")
	}

	meta := inn.Meta(context.Background())
	meta.Check = &rpcapi.HealthCheck{
		Timeout:                        (checkInterval + checkTimeout).String(),
		DeregisterCriticalServiceAfter: conf.Health.DeregisterCriticalServiceAfter,
	}
	if _, err := inn.kvrAgent.RPCRegister(context.TODO(), meta); err != nil {
		return errors.Wrap(err, "failed to register node in kvrouter")
	}

	// Run TTL updater
	go inn.updateTTLRoute(checkInterval)

	return nil
}

//heartbeat
func (inn *LocalNode) updateTTLRoute(interval time.Duration) {
	ticker := time.NewTicker(interval)
	p := &rpcapi.Ping{
		NodeID: inn.ID(),
	}
	defer ticker.Stop()
	for range ticker.C {
		if _, err := inn.kvrAgent.RPCHeartbeat(context.TODO(), p); err != nil {
			logger.Logger().Warn("check failed", zap.Error(err))
		}
	}
}
