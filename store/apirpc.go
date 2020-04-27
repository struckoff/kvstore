package store

import (
	"context"
	"github.com/struckoff/kvstore/router/nodes"
	"log"
	"net"

	"github.com/struckoff/kvstore/router/rpcapi"
	"google.golang.org/grpc"
)

func (inn *LocalNode) RunRPCServer(conf *Config) error {
	addy, err := net.ResolveTCPAddr("tcp", conf.RPCAddress)
	if err != nil {
		return err
	}
	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		return err
	}
	inn.rpcserver = grpc.NewServer()
	rpcapi.RegisterRPCNodeServer(inn.rpcserver, inn)

	if err := inn.rpcserver.Serve(inbound); err != nil {
		return err
	}
	return nil
}

func (inn *LocalNode) RPCStore(ctx context.Context, in *rpcapi.KeyValue) (*rpcapi.Empty, error) {
	if err := inn.Store(in.Key, in.Value); err != nil {
		return nil, err
	}
	return &rpcapi.Empty{}, nil
}

func (inn *LocalNode) RPCStorePairs(ctx context.Context, in *rpcapi.KeyValues) (*rpcapi.Empty, error) {
	log.Println("Receive keys")
	if err := inn.StorePairs(in.KVs); err != nil {
		return nil, err
	}
	return &rpcapi.Empty{}, nil
}

func (inn *LocalNode) RPCReceive(ctx context.Context, in *rpcapi.KeyReq) (*rpcapi.KeyValue, error) {
	b, err := inn.Receive(in.Key)
	if err != nil {
		return nil, err
	}
	return &rpcapi.KeyValue{Key: in.Key, Value: b}, nil
}

func (inn *LocalNode) RPCRemove(ctx context.Context, in *rpcapi.KeyReq) (*rpcapi.Empty, error) {
	if err := inn.Remove(in.Key); err != nil {
		return nil, err
	}
	return &rpcapi.Empty{}, nil
}

func (inn *LocalNode) RPCExplore(ctx context.Context, in *rpcapi.Empty) (*rpcapi.ExploreRes, error) {
	keys, err := inn.Explore()
	if err != nil {
		return nil, err
	}
	return &rpcapi.ExploreRes{Keys: keys}, nil
}
func (inn *LocalNode) RPCMeta(ctx context.Context, in *rpcapi.Empty) (*rpcapi.NodeMeta, error) {
	meta := &rpcapi.NodeMeta{
		ID:         inn.ID(),
		Address:    inn.HTTPAddress(),
		RPCAddress: inn.RPCAddress(),
		Power:      inn.Power().Get(),
		Capacity:   inn.Capacity().Get(),
	}
	return meta, nil
}

func (inn *LocalNode) RPCMove(ctx context.Context, in *rpcapi.MoveReq) (*rpcapi.Empty, error) {
	var en nodes.Node
	var err error

	res := make(map[nodes.Node][]string)
	for _, kl := range in.KL {
		if inn.kvr != nil {
			en, err = inn.kvr.GetNode(kl.Node.ID)
		} else {
			en, err = nodes.NewExternalNode(kl.Node, nil)
		}
		if err != nil {
			return nil, err
		}
		res[en] = kl.Keys
	}
	if err := inn.Move(res); err != nil {
		return nil, err
	}

	return &rpcapi.Empty{}, nil
}
