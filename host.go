package kvstore

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	consulwatch "github.com/hashicorp/consul/api/watch"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"strings"
	"time"
)

// Host represents bounding of network api with balancer lib and local node
type Host struct {
	bal       Balancer
	n         *InternalNode
	rpcserver *grpc.Server
	consul    *consulapi.Client
}

func NewHost(n *InternalNode, bal Balancer, c *consulapi.Client) (*Host, error) {
	if err := bal.AddNode(n); err != nil {
		return nil, err
	}
	h := &Host{
		bal:    bal,
		n:      n,
		consul: c,
	}
	return h, nil
}

// AddNode adds node to balancer
func (h *Host) AddNode(n Node) error {
	return h.bal.AddNode(n)
}

// RemoveNode removes node from balancer
func (h *Host) RemoveNode(id string) error {
	return h.bal.RemoveNode(id)
}

// Returns node from balancer by given key.
func (h *Host) GetNode(key string) (Node, error) {
	//di := DataItem(key)
	nb, err := h.bal.LocateKey(key)
	if err != nil {
		return nil, err
	}
	n, ok := nb.(Node)
	if !ok {
		return nil, errors.New("wrong node type")
	}
	return n, nil
}

func (h *Host) RunConsul(conf *Config) error {
	if err := h.consulAnnounce(conf); err != nil {
		return errors.Wrap(err, "unable to run announce node in consul")
	}
	if err := h.consulWatch(conf); err != nil {
		return errors.Wrap(err, "unable to run consul watcher")
	}
	return nil
}

func (h *Host) consulAnnounce(conf *Config) (err error) {
	checkID := h.n.ID() + "_ttl"

	addrParts := strings.Split(conf.RPCAddress, ":")
	if len(addrParts) < 2 {
		return errors.New("address format should be HOST:PORT")
	}
	port, err := strconv.ParseInt(addrParts[1], 10, 64)
	if err != nil {
		return err
	}

	checkInterval, err := time.ParseDuration(conf.Consul.CheckInterval)
	if err != nil {
		return err
	}
	checkTimeout, err := time.ParseDuration(conf.Consul.CheckTimeout)
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
		DeregisterCriticalServiceAfter: conf.Consul.DeregisterCriticalServiceAfter,
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

	if err := h.consul.Agent().ServiceRegister(service); err != nil {
		return err
	}

	// Run TTL updater
	go h.updateTTL(checkInterval, checkID)

	return nil
}
func (h *Host) consulWatch(conf *Config) error {
	filter := map[string]interface{}{
		"type":    "service",
		"service": conf.Consul.Service,
	}

	pl, err := consulwatch.Parse(filter)
	if err != nil {
		return err
	}
	pl.Handler = h.serviceHandler
	return pl.RunWithConfig(conf.Consul.Address, &conf.Consul.Config)
}
func (h *Host) serviceHandler(id uint64, data interface{}) {
	entries, ok := data.([]*consulapi.ServiceEntry)
	//fmt.Println(id, len(entries))
	if !ok {
		return
	}
	for _, entry := range entries {
		if entry.Node.Node == h.n.ID() {
			continue
		}
		addr := fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port)
		go h.registerExternalNode(entry.Node.Node, addr)
	}
}
func (h *Host) registerExternalNode(id, addr string) {
	en, err := NewExternalNode(addr)
	if err != nil {
		log.Printf("unable to connect to node %s(%s)", id, addr)
		return
	}
	if err := h.bal.AddNode(en); err != nil {
		log.Printf("unable to add node %s(%s) to balancer", id, addr)
		return
	}
	log.Printf("registered node %s(%s) to balancer", id, addr)
}
func (h *Host) updateTTL(interval time.Duration, checkID string) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		if err := h.consul.Agent().PassTTL(checkID, ""); err != nil {
			log.Printf("err=\"Check failed\" msg=\"%s\"", err.Error())
		}
	}
}
