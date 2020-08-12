package router

import (
	"context"
	"github.com/struckoff/SFCFramework/curve"
	"github.com/struckoff/kvstore/router/balanceradapter"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/dataitem"
	"github.com/struckoff/kvstore/router/nodehasher"
	"github.com/struckoff/kvstore/router/nodes"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/rpcapi"
	"github.com/struckoff/kvstore/router/ttl"
)

func NewHost(conf *config.Config) (*Host, error) {
	var err error
	var bal balanceradapter.Balancer
	var hr nodehasher.Hasher

	switch conf.Balancer.Mode {
	case config.ConsistentMode:
		bal = balanceradapter.NewConsistentBalancer()
	case config.SFCMode:
		bal, err = balanceradapter.NewSFCBalancer(conf.Balancer)
		if err != nil {
			return nil, err
		}
	}

	switch conf.Balancer.NodeHash {
	case config.GeoSfc:
		sb := bal.(*balanceradapter.SFC)
		sfc, err := curve.NewCurve(conf.Balancer.SFC.Curve.CurveType, 2, sb.SFC().Bits())
		if err != nil {
			return nil, errors.Wrap(err, "failed to create curve")
		}
		hr = nodehasher.NewGeoSfc(sfc)
	case config.XXHash:
		hr = nodehasher.NewXXHash()
	default:
		return nil, errors.New("invalid node hasher")
	}

	ndf, err := dataitem.GetDataItemFunc(conf.Balancer.DataMode)
	if err != nil {
		return nil, err
	}
	kvr, err := NewRouter(bal, hr, ndf, conf.Balancer)
	if err != nil {
		return nil, err
	}
	h := &Host{
		kvr:         kvr,
		checks:      ttl.NewChecksMap(),
		httplatency: conf.Balancer.Latency.Duration,
	}
	return h, nil
}

type Host struct {
	kvr         *Router
	checks      *ttl.ChecksMap
	httplatency time.Duration
}

func (h *Host) RPCRegister(ctx context.Context, in *rpcapi.NodeMeta) (*rpcapi.Empty, error) {
	en, err := nodes.NewExternalNode(in, h.kvr.Hasher())
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "failed to create external node")
	}

	onDead := onDeadHandler(en.ID())
	onRemove := h.onRemoveHandler(en.ID())
	check, err := ttl.NewTTLCheck(in.Check, onDead, onRemove)
	if err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "failed to create ttl check")
	}
	h.checks.Store(en.ID(), check)
	if err := h.kvr.AddNode(en); err != nil {
		log.Println(err)
		return nil, errors.Wrap(err, "failed to addNode")
	}
	return &rpcapi.Empty{}, nil
}
func (h *Host) RPCHeartbeat(ctx context.Context, in *rpcapi.Ping) (*rpcapi.Empty, error) {
	if ok := h.checks.Update(in.NodeID); !ok {
		return nil, errors.Errorf("unable to find check for node(%s)", in.NodeID)
	}
	return &rpcapi.Empty{}, nil
}
func (h *Host) RunHTTPServer(addr string) error {
	r := h.kvr.HTTPHandler()
	log.Printf("Run server [%s]", addr)
	l := LatencyMiddleware(r, h.httplatency)
	if err := http.ListenAndServe(addr, l); err != nil {
		return err
	}
	return nil
}

func onDeadHandler(nodeID string) func() {
	return func() {
		log.Printf("node(%s) seems to be dead", nodeID)
	}
}
func (h *Host) onRemoveHandler(nodeID string) func() {
	return func() {
		if err := h.kvr.RemoveNode(nodeID); err != nil {
			log.Printf("Error removing node(%s): %s", nodeID, err.Error())
			return
		}
		if err := h.kvr.Optimize(); err != nil {
			log.Printf("Error redistributing keys: %s", err.Error())
			return
		}
		h.checks.Delete(nodeID)
		log.Printf("node(%s) removed", nodeID)
	}
}
