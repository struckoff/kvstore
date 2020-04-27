package store

import (
	"fmt"
	"github.com/struckoff/kvstore/router/nodes"
	"log"
	"strconv"
	"strings"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	consulwatch "github.com/hashicorp/consul/api/watch"
	"github.com/pkg/errors"
)

// RunRouter - register service in consul
// Function register TTL check and sends heartbeat each Config.Health.CheckInterval
func (inn *LocalNode) RunConsul(conf *Config) error {
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
		CheckID: checkID,
		Name:    checkID,
		Status:  "passing",
		//TCP:      conf.RPCAddress,
		//Interval: conf.ConfigConsul.CheckInterval,
		//Timeout:  conf.ConfigConsul.CheckTimeout,
		//AliasNode:                      conf.Name,
		//AliasService:                   conf.ConfigConsul.Service,
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
		Meta: map[string]string{
			"power": checkID,
		},
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
	entries, ok := data.([]*consulapi.ServiceEntry)
	//fmt.Println(id, len(entries))
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
		log.Println(err.Error())
	}
	locations, err := inn.keysLocations()
	if err != nil {
		log.Println(err.Error())
		return
	}
	inn.Move(locations)
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
		log.Printf("unable to connect to node %s(%s)", id, addr)
		nCh <- nil
		return
	}
	nCh <- en
	log.Printf("registered node %s(%s)", id, addr)
}

//keysLocations - returns keys which should be moved to another node
func (inn *LocalNode) keysLocations() (map[nodes.Node][]string, error) {
	res := make(map[nodes.Node][]string)
	keys, err := inn.Explore()
	if err != nil {
		return nil, err
	}
	for iter := range keys {
		n, err := inn.kvr.LocateKey(keys[iter])
		if err != nil {
			return nil, err
		}
		if inn.ID() != n.ID() {
			res[n] = append(res[n], keys[iter])
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
			log.Printf("err=\"Check failed\" msg=\"%s\"", err.Error())
		}
	}
}
