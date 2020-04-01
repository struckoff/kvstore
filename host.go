package kvstore

import (
	consulapi "github.com/hashicorp/consul/api"
	consulwatch "github.com/hashicorp/consul/api/watch"
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/rpcapi"
	"google.golang.org/grpc"
	"log"
)

// Host represents bounding of network api with balancer lib and local node
type Host struct {
	bal       Balancer
	n         *InternalNode
	rpcserver *grpc.Server
}

func NewHost(n *InternalNode, bal Balancer, gs *grpc.Server) (*Host, error) {
	if err := bal.AddNode(n); err != nil {
		return nil, err
	}
	h := &Host{bal, n, gs}
	rpcapi.RegisterRPCListenerServer(gs, h)
	hch := new(HealthCheck)
	rpcapi.RegisterHealthServer(gs, hch)
	return h, nil
}

// Lookup tries to connect to the remote node using addresses in the given list.
// Function ends on first success.
// If all attempts fail it will return an error.
//func (h *Host) Lookup(eps []string) error {
//	for _, addr := range eps {
//		p := strings.Join([]string{addr, "node"}, "/")
//		m := h.n.Meta()
//		buf := bytes.NewBuffer(nil)
//		if err := json.NewEncoder(buf).Encode(m); err != nil {
//			log.Println("[ERROR]", err)
//			continue
//		}
//		r, err := http.Post(p, "application", buf)
//		if err != nil {
//			log.Println("[ERROR]", err)
//			continue
//		}
//		if r.StatusCode >= 400 {
//			log.Println("[ERROR]", r.Status)
//			continue
//		}
//		var metas []NodeMeta
//		if err := json.NewDecoder(r.Body).Decode(&metas); err != nil {
//			log.Println("[ERROR]", err)
//			continue
//		}
//		for _, meta := range metas {
//			en, err := NewExternalNode(meta)
//			if err != nil {
//				return err
//			}
//			if err := h.bal.AddNode(en); err != nil {
//				return err
//			}
//		}
//		return nil
//	}
//	return errors.New("unable to connect to nodes")
//}

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

func (h *Host) consulAnnounce(conf *Config) error {
	consul, err := consulapi.NewClient(&conf.Consul.Config)
	if err != nil {
		return err
	}

	service, _, err := consul.Agent().Service(conf.Consul.Service, nil)
	if err != nil {
		return err
	}
	//hcd := consulapi.HealthCheckDefinition{
	//	//HTTP:   conf.Address + "/list",
	//	//Header: nil,
	//	//Method: "GET",
	//	//Body:                                   "",
	//	//TLSSkipVerify:                          false,
	//	TCP:              conf.RPCAddress,
	//	IntervalDuration: time.Second * 3,
	//	TimeoutDuration:  time.Second * 3,
	//	//DeregisterCriticalServiceAfterDuration: 0,
	//	//Interval: *consulapi.NewReadableDuration(time.Second),
	//	//Timeout:  *consulapi.NewReadableDuration(2 * time.Second),
	//	//DeregisterCriticalServiceAfter:         0,
	//}
	//
	//hc := &consulapi.HealthCheck{
	//	Node:        conf.Name,
	//	CheckID:     conf.Name + "-hc",
	//	Name:        "base node hc",
	//	ServiceID:   service.Service,
	//	ServiceName: conf.Consul.Service,
	//	Definition:  hcd,
	//}

	ac := &consulapi.AgentCheckRegistration{
		ID:        conf.Name + "_check",
		Name:      conf.Name,
		ServiceID: service.Service,
		AgentServiceCheck: consulapi.AgentServiceCheck{
			Name:     conf.Name,
			Status:   "passing",
			TCP:      conf.RPCAddress,
			Interval: conf.Consul.CheckInterval,
			Timeout:  conf.Consul.CheckTimeout,
		},
	}

	//acnode := &consulapi.AgentCheck{
	//	Node:    conf.Name,
	//	CheckID: conf.Name,
	//	Name:    conf.Name,
	//	//Status:  "passing",
	//	//Notes:       "",
	//	//Output:      "",
	//	Type:        "TCP",
	//	ServiceID:   conf.Consul.Service,
	//	ServiceName: conf.Consul.Service,
	//	Definition:  hcd,
	//	Namespace:   conf.Consul.Namespace,
	//}

	if err := consul.Agent().CheckRegister(ac); err != nil {
		return err
	}

	reg := &consulapi.CatalogRegistration{
		Node:       conf.Name,
		Address:    conf.RPCAddress,
		Datacenter: conf.Consul.Datacenter,
		Service:    service,
		//Check:      acnode,
		//Checks:     consulapi.HealthChecks{hc},
		//NodeMeta: map[string]string{
		//	"external-node":  "true",
		//	"external-probe": "true",
		//},
	}
	_, err = consul.Catalog().Register(reg, nil)
	if err != nil {
		return err
	}
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
	if !ok {
		return
	}
	for _, entry := range entries {
		if entry.Node.Node == h.n.id {
			continue
		}
		go h.registerExternalNode(entry.Node.Node, entry.Node.Address)
	}
}

func (h *Host) registerExternalNode(id, addr string) {
	en, err := NewExternalNode(addr)
	if err != nil {
		//log.Printf("unable to connect to node %s(%s)", id, addr)
		return
	}
	if err := h.bal.AddNode(en); err != nil {
		log.Printf("unable to add node %s(%s) to balancer", id, addr)
		return
	}
	log.Printf("registered node %s(%s) to balancer", id, addr)
}
