package store

import (
	"context"
	"fmt"
	"github.com/struckoff/kvstore/logger"
	"github.com/struckoff/kvstore/router/nodes"
	"github.com/struckoff/kvstore/router/rpcapi"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	consulwatch "github.com/hashicorp/consul/api/watch"
	"github.com/pkg/errors"
)

// RunRouter - register service in consul
// Function register TTL check and sends heartbeat each Config.Health.CheckInterval
func (inn *LocalNode) RunConsul(conf *Config) error {
	logger.Logger().Info("announcing in consul", zap.String("Address", conf.Address))
	if err := inn.consulAnnounce(conf); err != nil {
		return errors.Wrap(err, "unable to run announce node in consul")
	}
	if err := inn.consulWatch(conf); err != nil {
		return errors.Wrap(err, "unable to run consul watcher")
	}
	return nil
}

// consulAnnounce - register node and TTL check in consul
func (inn *LocalNode) consulAnnounce(conf *Config) (err error) {
	checkID := inn.ID() + "_ttl"

	addrParts := strings.Split(conf.RPCAddress, ":")
	if len(addrParts) < 2 {
		return errors.New("address format should be HOST:PORT")
	}
	port, err := strconv.ParseInt(addrParts[1], 10, 64)
	if err != nil {
		return err
	}

	checkInterval, err := time.ParseDuration(conf.Health.CheckInterval)
	if err != nil {
		return err
	}
	checkTimeout, err := time.ParseDuration(conf.Health.CheckTimeout)
	if err != nil {
		return err
	}

	// Create heartbeat check
	acc := consulapi.AgentServiceCheck{
		CheckID:                        checkID,
		Name:                           checkID,
		Status:                         consulapi.HealthCritical,
		DeregisterCriticalServiceAfter: conf.Health.DeregisterCriticalServiceAfter,
		TTL:                            (checkInterval + checkTimeout).String(),
	}

	service := &consulapi.AgentServiceRegistration{
		ID:   conf.Consul.Service,
		Name: conf.Consul.Service,
		//Tags:              nil,
		Port:      int(port),
		Address:   addrParts[0],
		Check:     &acc,
		Namespace: conf.Consul.Namespace,
	}

	if err := inn.consul.Agent().ServiceRegister(service); err != nil {
		return err
	}

	// Run TTL updater
	go inn.updateTTLConsul(checkInterval, checkID)

	return nil
}

// serviceWatch - cosul service watcher
func (inn *LocalNode) consulWatch(conf *Config) error {
	filter := map[string]interface{}{
		"type":    "service",
		"service": conf.Consul.Service,
	}

	pl, err := consulwatch.Parse(filter)
	if err != nil {
		return err
	}
	pl.Handler = inn.serviceHandler
	return pl.RunWithConfig(conf.Consul.Address, &conf.Consul.Config)
}

// serviceHandler - callback which calls on changes in service inside consul.
func (inn *LocalNode) serviceHandler(id uint64, data interface{}) {
	nCh := make(chan nodes.Node)
	defer close(nCh)

	lwID := atomic.LoadInt64(inn.lwID)
	if lwID >= int64(id) {
		logger.Logger().Warn("event ID less than last processed ID", zap.Uint64("event ID", id), zap.Int64("last processed ID", lwID))
		return
	}

	entries, ok := data.([]*consulapi.ServiceEntry)
	if !ok {
		return
	}
	for _, entry := range entries {
		addr := fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port)
		go inn.registerExternalNode(entry.Node.Node, addr, nCh)
	}
	count := len(entries)
	ns := make([]nodes.Node, 0, count)
	for node := range nCh {
		count--
		if node != nil {
			ns = append(ns, node)
		}
		if count == 0 {
			break
		}
	}
	if err := inn.kvr.SetNodes(ns); err != nil {
		logger.Logger().Error("unable to set nodes", zap.Error(err))
	}
	locations, err := inn.keysLocations(context.Background())
	if err != nil {
		logger.Logger().Error("unable to get keys locations", zap.Error(err))
		return
	}
	if err := inn.Move(context.Background(), locations); err != nil {
		logger.Logger().Error("local node move failed", zap.Error(err))
		return
	}

	if swapped := atomic.CompareAndSwapInt64(inn.lwID, lwID, int64(id)); !swapped {
		logger.Logger().Warn("last watcher ID was not swapped")
		inn.serviceHandler(id, entries)
		return
	}
}

//registerExternalNode - register nodes from consul in local kvrouter.
// Function gather node information from nodes by RPC.
func (inn *LocalNode) registerExternalNode(id, addr string, nCh chan<- nodes.Node) {
	if id == inn.ID() {
		nCh <- inn
		return
	}
	en, err := nodes.NewExternalNodeByAddr(addr, inn.kvr.Hasher())
	if err != nil {
		logger.Logger().Info("unable to connect to node", zap.String("Node", id), zap.String("Address", addr), zap.Error(err))
		nCh <- nil
		return
	}
	nCh <- en
	logger.Logger().Info("registered node", zap.String("Node", id), zap.String("Address", addr))
}

//keysLocations - returns keys which should be moved to another node
func (inn *LocalNode) keysLocations(ctx context.Context) (map[nodes.Node][]*rpcapi.DataItem, error) {
	res := make(map[nodes.Node][]*rpcapi.DataItem)
	dis, err := inn.Explore(ctx)
	if err != nil {
		return nil, err
	}
	for i := range dis {
		n, _, err := inn.kvr.LocateKey(dis[i])
		if err != nil {
			return nil, err
		}
		if inn.ID() != n.ID() {
			res[n] = append(res[n], dis[i])
		}
	}
	return res, nil
}

// heartbeat
func (inn *LocalNode) updateTTLConsul(interval time.Duration, checkID string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		if err := inn.consul.Agent().PassTTL(checkID, ""); err != nil {
			logger.Logger().Warn("check failed", zap.Error(err))
		}
	}
}
