package kvstore

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	consulwatch "github.com/hashicorp/consul/api/watch"
	"github.com/pkg/errors"
	kvrouter "github.com/struckoff/kvrouter/router"
	"log"
	"strconv"
	"strings"
	"time"
)

func (inn *InternalNode) RunConsul(conf *Config) error {
	if err := inn.consulAnnounce(conf); err != nil {
		return errors.Wrap(err, "unable to run announce node in consul")
	}
	if err := inn.consulWatch(conf); err != nil {
		return errors.Wrap(err, "unable to run consul watcher")
	}
	return nil
}

func (inn *InternalNode) consulAnnounce(conf *Config) (err error) {
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
func (inn *InternalNode) consulWatch(conf *Config) error {
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
func (inn *InternalNode) serviceHandler(id uint64, data interface{}) {
	nCh := make(chan kvrouter.Node)
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
	ns := make([]kvrouter.Node, 0, count)
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
	inn.relocate(locations)
}
func (inn *InternalNode) registerExternalNode(id, addr string, nCh chan<- kvrouter.Node) {
	if id == inn.ID() {
		nCh <- inn
		return
	}
	en, err := kvrouter.NewExternalNodeByAddr(addr)
	if err != nil {
		log.Printf("unable to connect to node %s(%s)", id, addr)
		nCh <- nil
		return
	}
	nCh <- en
	log.Printf("registered node %s(%s)", id, addr)
}
func (inn *InternalNode) keysLocations() (map[kvrouter.Node][]string, error) {
	res := make(map[kvrouter.Node][]string)
	keys, err := inn.Explore()
	if err != nil {
		return nil, err
	}
	for iter := range keys {
		n, err := inn.kvr.GetNode(keys[iter])
		if err != nil {
			return nil, err
		}
		if inn.ID() != n.ID() {
			res[n] = append(res[n], keys[iter])
		}
	}
	return res, nil
}
func (inn *InternalNode) relocate(locations map[kvrouter.Node][]string) {
	for n, keys := range locations {
		go func(n kvrouter.Node, keys []string) {
			if err := inn.Move(keys, n); err != nil {
				log.Println(err.Error())
				return
			}
		}(n, keys)
	}
}
func (inn *InternalNode) updateTTLConsul(interval time.Duration, checkID string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		if err := inn.consul.Agent().PassTTL(checkID, ""); err != nil {
			log.Printf("err=\"Check failed\" msg=\"%s\"", err.Error())
		}
	}
}
