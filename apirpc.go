package kvstore

import (
	"context"
	"github.com/struckoff/kvrouter/rpcapi"
	"google.golang.org/grpc"
	"log"
	"net"
)

func (n *InternalNode) RunRPCServer(conf *Config) error {
	addy, err := net.ResolveTCPAddr("tcp", conf.RPCAddress)
	if err != nil {
		return err
	}
	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		return err
	}
	n.rpcserver = grpc.NewServer()
	rpcapi.RegisterRPCListenerServer(n.rpcserver, n)

	if err := n.rpcserver.Serve(inbound); err != nil {
		return err
	}
	return nil
}

func (n *InternalNode) RPCStore(ctx context.Context, in *rpcapi.KeyValue) (*rpcapi.Empty, error) {
	if err := n.Store(in.Key, in.Value); err != nil {
		return nil, err
	}
	return &rpcapi.Empty{}, nil
}

func (n *InternalNode) RPCStorePairs(ctx context.Context, in *rpcapi.KeyValues) (*rpcapi.Empty, error) {
	log.Println("Receive keys")
	if err := n.StorePairs(in.KVs); err != nil {
		return nil, err
	}
	return &rpcapi.Empty{}, nil
}

func (n *InternalNode) RPCReceive(ctx context.Context, in *rpcapi.KeyReq) (*rpcapi.KeyValue, error) {
	b, err := n.Receive(in.Key)
	if err != nil {
		return nil, err
	}
	return &rpcapi.KeyValue{Key: in.Key, Value: b}, nil
}

func (n *InternalNode) RPCRemove(ctx context.Context, in *rpcapi.KeyReq) (*rpcapi.Empty, error) {
	if err := n.Remove(in.Key); err != nil {
		return nil, err
	}
	return &rpcapi.Empty{}, nil
}

func (n *InternalNode) RPCExplore(ctx context.Context, in *rpcapi.Empty) (*rpcapi.ExploreRes, error) {
	keys, err := n.Explore()
	if err != nil {
		return nil, err
	}
	return &rpcapi.ExploreRes{Keys: keys}, nil
}
func (n *InternalNode) RPCMeta(ctx context.Context, in *rpcapi.Empty) (*rpcapi.NodeMeta, error) {
	meta := &rpcapi.NodeMeta{
		ID:         n.ID(),
		Address:    n.HTTPAddress(),
		RPCAddress: n.RPCAddress(),
		Power:      n.Power().Get(),
		Capacity:   n.Capacity().Get(),
	}
	return meta, nil
}
