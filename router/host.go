package router

import (
	"context"
	"github.com/struckoff/kvstore/logger"
	"github.com/struckoff/kvstore/router/balanceradapter"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/dataitem"
	"github.com/struckoff/kvstore/router/nodehasher"
	"github.com/struckoff/kvstore/router/nodes"
	"github.com/struckoff/sfcframework/curve"
	"go.uber.org/zap"
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
	rpcndf, err := dataitem.GetDataItemFromRpcFunc(conf.Balancer.DataMode)
	if err != nil {
		return nil, err
	}
	kvr, err := NewRouter(bal, hr, ndf, rpcndf, conf.Balancer)
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

func (h *Host) RPCRegister(_ context.Context, in *rpcapi.NodeMeta) (*rpcapi.Empty, error) {
	en, err := nodes.NewExternalNode(in, h.kvr.Hasher())
	if err != nil {
		logger.Logger().Error("Error registering node", zap.Error(err))
		return nil, errors.Wrap(err, "failed to create external node")
	}

	onDead := onDeadHandler(en.ID())
	onRemove := h.onRemoveHandler(en.ID())
	check, err := ttl.NewTTLCheck(in.Check, onDead, onRemove)
	if err != nil {
		logger.Logger().Error("Error registering check", zap.String("Node", en.ID()), zap.Error(err))
		return nil, errors.Wrap(err, "failed to create ttl check")
	}
	h.checks.Store(en.ID(), check)
	if err := h.kvr.AddNode(en); err != nil {
		logger.Logger().Error("Error adding node", zap.String("Node", en.ID()), zap.Error(err))
		return nil, errors.Wrap(err, "failed to addNode")
	}
	return &rpcapi.Empty{}, nil
}
func (h *Host) RPCHeartbeat(_ context.Context, in *rpcapi.Ping) (*rpcapi.Empty, error) {
	if ok := h.checks.Update(in.NodeID); !ok {
		return nil, errors.Errorf("unable to find check for node(%s)", in.NodeID)
	}
	return &rpcapi.Empty{}, nil
}
func (h *Host) RunHTTPServer(addr string) error {
	r := h.kvr.HTTPHandler()
	logger.Logger().Info("Run server", zap.String("address", addr))
	l := LatencyMiddleware(r, h.httplatency)
	if err := http.ListenAndServe(addr, l); err != nil {
		return err
	}
	return nil
}

func onDeadHandler(nodeID string) func() {
	return func() {
		logger.Logger().Warn("node seems to be dead", zap.String("Node", nodeID))
	}
}
func (h *Host) onRemoveHandler(nodeID string) func() {
	return func() {
		if err := h.kvr.RemoveNode(nodeID); err != nil {
			logger.Logger().Error("Error removing node", zap.String("Node", nodeID), zap.Error(err))
			return
		}
		if err := h.kvr.Optimize(); err != nil {
			logger.Logger().Error("Error redistributing keys", zap.Error(err))
			return
		}
		h.checks.Delete(nodeID)
		logger.Logger().Info("node removed", zap.String("Node", nodeID))
	}
}
