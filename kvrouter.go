package kvstore

import (
	"context"
	"github.com/pkg/errors"
	"github.com/struckoff/kvrouter/rpcapi"
	"google.golang.org/grpc"
	"log"
	"time"
)

func (inn *InternalNode) RunKVRouter(conf *Config) error {
	c, err := kvClient(conf.KVRouter.Address)
	if err != nil {
		return errors.Wrap(err, "failed to initialize kvrouter client")
	}
	inn.kvrAgent = c
	if err := inn.kvrouterAnnounce(conf); err != nil {
		return errors.Wrap(err, "unable to run announce node in kvrouter")
	}
	return nil
}

func kvClient(addr string) (rpcapi.RPCBalancerClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure()) // TODO Make it secure
	if err != nil {
		return nil, err
	}
	c := rpcapi.NewRPCBalancerClient(conn)
	return c, nil
}

func (inn *InternalNode) kvrouterAnnounce(conf *Config) error {
	checkInterval, err := time.ParseDuration(conf.Health.CheckInterval)
	if err != nil {
		return errors.Wrap(err, "failed to parse health check interval")
	}
	checkTimeout, err := time.ParseDuration(conf.Health.CheckTimeout)
	if err != nil {
		return errors.Wrap(err, "failed to parse health check timeout")
	}

	meta := inn.Meta()
	meta.Check = &rpcapi.HealthCheck{
		Timeout:                        (checkInterval + checkTimeout).String(),
		DeregisterCriticalServiceAfter: conf.Health.DeregisterCriticalServiceAfter,
	}
	if _, err := inn.kvrAgent.RPCRegister(context.TODO(), &meta); err != nil {
		return errors.Wrap(err, "failed to register node in kvrouter")
	}

	// Run TTL updater
	go inn.updateTTLKVRoute(checkInterval)

	return nil
}
func (inn *InternalNode) updateTTLKVRoute(interval time.Duration) {
	ticker := time.NewTicker(interval)
	p := &rpcapi.Ping{
		NodeID: inn.ID(),
	}
	defer ticker.Stop()
	for range ticker.C {
		if _, err := inn.kvrAgent.RPCHeartbeat(context.TODO(), p); err != nil {
			log.Printf("err=\"Check failed\" msg=\"%s\"", err.Error())
		}
	}
}
