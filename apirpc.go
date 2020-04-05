package kvstore

import (
	"context"
	"github.com/struckoff/kvstore/rpcapi"
	"google.golang.org/grpc"
	"log"
	"net"
)

func (h *Host) RunRPCServer(conf *Config) error {
	addy, err := net.ResolveTCPAddr("tcp", conf.RPCAddress)
	if err != nil {
		return err
	}
	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		return err
	}
	h.rpcserver = grpc.NewServer()
	rpcapi.RegisterRPCListenerServer(h.rpcserver, h)

	if err := h.rpcserver.Serve(inbound); err != nil {
		return err
	}
	return nil
}

func (h *Host) RPCStore(ctx context.Context, in *rpcapi.KeyValue) (*rpcapi.Empty, error) {
	if err := h.n.Store(in.Key, in.Value); err != nil {
		return nil, err
	}
	return &rpcapi.Empty{}, nil
}

func (h *Host) RPCStorePairs(ctx context.Context, in *rpcapi.KeyValues) (*rpcapi.Empty, error) {
	log.Println("Receive keys")
	if err := h.n.StorePairs(in.KVs); err != nil {
		return nil, err
	}
	return &rpcapi.Empty{}, nil
}

func (h *Host) RPCReceive(ctx context.Context, in *rpcapi.KeyReq) (*rpcapi.KeyValue, error) {
	b, err := h.n.Receive(in.Key)
	if err != nil {
		return nil, err
	}
	return &rpcapi.KeyValue{Key: in.Key, Value: b}, nil
}

func (h *Host) RPCRemove(ctx context.Context, in *rpcapi.KeyReq) (*rpcapi.Empty, error) {
	if err := h.n.Remove(in.Key); err != nil {
		return nil, err
	}
	return &rpcapi.Empty{}, nil
}

func (h *Host) RPCExplore(ctx context.Context, in *rpcapi.Empty) (*rpcapi.ExploreRes, error) {
	keys, err := h.n.Explore()
	if err != nil {
		return nil, err
	}
	return &rpcapi.ExploreRes{Keys: keys}, nil
}
func (h *Host) RPCMeta(ctx context.Context, in *rpcapi.Empty) (*rpcapi.NodeMeta, error) {
	meta := &rpcapi.NodeMeta{
		ID:         h.n.ID(),
		Address:    h.n.HTTPAddress(),
		RPCAddress: h.n.RPCAddress(),
		Power:      h.n.Power().Get(),
		Capacity:   h.n.Capacity().Get(),
	}
	return meta, nil
}
