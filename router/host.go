package router

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/rpcapi"
	"github.com/struckoff/kvstore/router/ttl"
)

type Host struct {
	kvr    *Router
	checks *ttl.ChecksMap
}

func (h *Host) RPCRegister(ctx context.Context, in *rpcapi.NodeMeta) (*rpcapi.Empty, error) {
	en, err := NewExternalNode(in)
	if err != nil {
		return nil, err
	}

	onDead := deadHandler(en.ID())
	onRemove := h.removeHandler(en.ID())
	check, err := ttl.NewTTLCheck(in.Check, onDead, onRemove)
	if err != nil {
		return nil, err
	}
	h.checks.Store(en.ID(), check)
	if err := h.kvr.AddNode(en); err != nil {
		return nil, err
	}
	log.Printf("node(%s) registered", en.ID())
	if err := h.redistributeKeys(); err != nil {
		return nil, err
	}
	return &rpcapi.Empty{}, nil
}
func (h *Host) RPCHeartbeat(ctx context.Context, in *rpcapi.Ping) (*rpcapi.Empty, error) {
	if ok := h.checks.Update(in.NodeID); !ok {
		return nil, errors.Errorf("unable to found check for node(%s)", in.NodeID)
	}
	return &rpcapi.Empty{}, nil
}
func (h *Host) RunHTTPServer(addr string) error {
	r := h.kvr.HTTPHandler()
	log.Printf("Run server [%s]", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		return err
	}
	return nil
}

func deadHandler(nodeID string) func() {
	return func() {
		log.Printf("node(%s) seems to be dead", nodeID)
	}
}
func (h *Host) removeHandler(nodeID string) func() {
	return func() {
		if err := h.kvr.RemoveNode(nodeID); err != nil {
			log.Printf("Error removing node(%s): %s", nodeID, err.Error())
			return
		}
		if err := h.redistributeKeys(); err != nil {
			log.Printf("Error redistributing keys: %s", err.Error())
			return
		}
		h.checks.Delete(nodeID)
		log.Printf("node(%s) removed", nodeID)
	}
}

func (h *Host) redistributeKeys() error {
	var wg sync.WaitGroup
	ns, err := h.kvr.GetNodes()
	if err != nil {
		return err
	}
	for _, n := range ns {
		go func(n Node, wg *sync.WaitGroup) {
			res := make(map[Node][]string)
			keys, err := n.Explore()
			if err != nil {
				log.Printf("failed to explore node(%s): %s", n.ID(), err.Error())
				return
			}
			for iter := range keys {
				en, err := h.kvr.LocateKey(keys[iter])
				if err != nil {
					log.Printf("failed to locate key(%s): %s", keys[iter], err.Error())
					continue
				}
				if en.ID() != n.ID() {
					res[en] = append(res[en], keys[iter])
				}
			}
			n.Move(res)
		}(n, &wg)
	}
	wg.Wait()
	return nil
}

func NewHost(conf *Config) (*Host, error) {
	bal, err := NewSFCBalancer(conf.Balancer)
	if err != nil {
		return nil, err
	}
	kvr, err := NewRouter(bal)
	if err != nil {
		return nil, err
	}
	h := &Host{
		kvr:    kvr,
		checks: ttl.NewChecksMap(),
	}
	return h, nil
}
