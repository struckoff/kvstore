package store

import (
	"context"
	"github.com/influxdata/influxdb-client-go"
	"github.com/struckoff/kvstore/logger"
	"github.com/struckoff/kvstore/router/nodes"
	"go.uber.org/zap"
	"google.golang.org/grpc/keepalive"
	"net"
	"time"

	"github.com/struckoff/kvstore/router/rpcapi"
	"google.golang.org/grpc"
)

func (inn *LocalNode) RunRPCServer(conf *Config, errCh chan<- error) error {
	addy, err := net.ResolveTCPAddr("tcp", conf.RPCAddress)
	if err != nil {
		return err
	}
	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		return err
	}
	inn.rpcserver = grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
	)
	rpcapi.RegisterRPCNodeServer(inn.rpcserver, inn)
	rpcapi.RegisterRPCCapacityServer(inn.rpcserver, &inn.c)

	go func(errCh chan<- error) {
		errCh <- inn.rpcserver.Serve(inbound)
	}(errCh)
	logger.Logger().Info("RPC server listening", zap.String("Address", inbound.Addr().String()))
	return nil
}

func (inn *LocalNode) RPCStore(ctx context.Context, in *rpcapi.KeyValue) (r *rpcapi.DataItem, err error) {
	start := time.Now()
	defer func() {
		end := time.Since(start)
		now := time.Now()

		p := influxdb2.NewPoint("grpc",
			map[string]string{
				"id":     inn.id,
				"method": "store",
			},
			map[string]interface{}{
				"duration_ms": end.Milliseconds(),
			},
			now)
		inn.metrics <- p
	}()

	time.Sleep(inn.rpclatency)
	return inn.Store(ctx, in)
}

func (inn *LocalNode) RPCStorePairs(ctx context.Context, in *rpcapi.KeyValues) (*rpcapi.DataItems, error) {
	start := time.Now()
	defer func() {
		end := time.Since(start)
		now := time.Now()

		p := influxdb2.NewPoint("grpc",
			map[string]string{
				"id":     inn.id,
				"method": "store_pairs",
			},
			map[string]interface{}{
				"duration_ms": end.Milliseconds(),
			},
			now)
		inn.metrics <- p
	}()

	time.Sleep(inn.rpclatency)
	dis, err := inn.StorePairs(ctx, in.KVs)
	if err != nil {
		return nil, err
	}
	return &rpcapi.DataItems{DIs: dis}, nil
}

func (inn *LocalNode) RPCReceive(ctx context.Context, in *rpcapi.DataItems) (*rpcapi.KeyValues, error) {
	start := time.Now()
	defer func() {
		end := time.Since(start)
		now := time.Now()

		p := influxdb2.NewPoint("grpc",
			map[string]string{
				"id":     inn.id,
				"method": "receive",
			},
			map[string]interface{}{
				"duration_ms": end.Milliseconds(),
			},
			now)
		inn.metrics <- p
	}()

	time.Sleep(inn.rpclatency)
	return inn.Receive(ctx, in.DIs)
}

func (inn *LocalNode) RPCRemove(ctx context.Context, in *rpcapi.DataItems) (*rpcapi.DataItems, error) {
	start := time.Now()
	defer func() {
		end := time.Since(start)
		now := time.Now()

		p := influxdb2.NewPoint("grpc",
			map[string]string{
				"id":     inn.id,
				"method": "remove",
			},
			map[string]interface{}{
				"duration_ms": end.Milliseconds(),
			},
			now)
		inn.metrics <- p
	}()

	time.Sleep(inn.rpclatency)
	dis, err := inn.Remove(ctx, in.DIs)
	if err != nil {
		return nil, err
	}
	ds := &rpcapi.DataItems{DIs: dis}
	return ds, nil
}

func (inn *LocalNode) RPCExplore(ctx context.Context, _ *rpcapi.Empty) (*rpcapi.DataItems, error) {
	start := time.Now()
	defer func() {
		end := time.Since(start)
		now := time.Now()

		p := influxdb2.NewPoint("grpc",
			map[string]string{
				"id":     inn.id,
				"method": "explore",
			},
			map[string]interface{}{
				"duration_ms": end.Milliseconds(),
			},
			now)
		inn.metrics <- p
	}()

	time.Sleep(inn.rpclatency)
	dis, err := inn.Explore(ctx)
	if err != nil {
		return nil, err
	}
	return &rpcapi.DataItems{DIs: dis}, nil
}
func (inn *LocalNode) RPCMeta(ctx context.Context, _ *rpcapi.Empty) (*rpcapi.NodeMeta, error) {
	start := time.Now()
	defer func() {
		end := time.Since(start)
		now := time.Now()

		p := influxdb2.NewPoint("grpc",
			map[string]string{
				"id":     inn.id,
				"method": "meta",
			},
			map[string]interface{}{
				"duration_ms": end.Milliseconds(),
			},
			now)
		inn.metrics <- p
	}()

	time.Sleep(inn.rpclatency)
	meta := inn.meta(ctx)
	return meta, nil
}

func (inn *LocalNode) RPCMove(ctx context.Context, in *rpcapi.MoveReq) (*rpcapi.Empty, error) {
	//start := time.Now()
	//defer func() {
	//	end := time.Since(start)
	//	now := time.Now()
	//
	//	p := influxdb2.NewPoint("grpc",
	//		map[string]string{
	//			"id":     inn.id,
	//			"method": "move",
	//		},
	//		map[string]interface{}{
	//			"duration_ms": end.Milliseconds(),
	//		},
	//		now)
	//	inn.metrics <- p
	//}()

	time.Sleep(inn.rpclatency)

	var en nodes.Node
	var err error

	res := make(map[nodes.Node][]*rpcapi.DataItem)
	for _, kl := range in.KLs {
		if inn.kvr != nil {
			en, err = inn.kvr.GetNode(kl.Node.ID)
		} else {
			en, err = nodes.NewExternalNode(kl.Node, nil)
		}
		if err != nil {
			return nil, err
		}
		res[en] = kl.Keys.DIs
	}
	if err := inn.Move(ctx, res); err != nil {
		return nil, err
	}

	return &rpcapi.Empty{}, nil
}
