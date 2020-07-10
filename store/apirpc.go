package store

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/struckoff/kvstore/router/nodes"
	"log"
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
	inn.rpcserver = grpc.NewServer()
	rpcapi.RegisterRPCNodeServer(inn.rpcserver, inn)
	rpcapi.RegisterRPCCapacityServer(inn.rpcserver, &inn.c)

	go func(errCh chan<- error) {
		errCh <- inn.rpcserver.Serve(inbound)
	}(errCh)
	//if err := checkTCP(inbound.Addr().String()); err != nil {
	//	return err
	//}
	log.Printf("RPC server listening on %s", inbound.Addr().String())

	//if err := inn.rpcserver.Serve(inbound); err != nil {
	//	return err
	//}

	return nil
}

//func checkTCP(addr string) error {
//	_, err := net.Dial("tcp", addr)
//	if err != nil {
//		return err
//	}
//	return nil
//}

func (inn *LocalNode) RPCStore(ctx context.Context, in *rpcapi.KeyValue) (r *rpcapi.Empty, err error) {
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
	if err := inn.Store(in.Key, []byte(in.Value)); err != nil {
		return nil, err
	}
	return &rpcapi.Empty{}, nil
}

func (inn *LocalNode) RPCStorePairs(ctx context.Context, in *rpcapi.KeyValues) (*rpcapi.Empty, error) {
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

	log.Println("Receive keys")
	time.Sleep(inn.rpclatency)
	if err := inn.StorePairs(in.KVs); err != nil {
		return nil, err
	}
	return &rpcapi.Empty{}, nil
}

func (inn *LocalNode) RPCReceive(ctx context.Context, in *rpcapi.KeyReq) (*rpcapi.KeyValues, error) {
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
	return inn.Receive(in.Keys)
}

func (inn *LocalNode) RPCRemove(ctx context.Context, in *rpcapi.KeyReq) (*rpcapi.Empty, error) {
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
	if err := inn.Remove(in.Keys); err != nil {
		return nil, err
	}
	return &rpcapi.Empty{}, nil
}

func (inn *LocalNode) RPCExplore(ctx context.Context, in *rpcapi.Empty) (*rpcapi.ExploreRes, error) {
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
	keys, err := inn.Explore()
	if err != nil {
		return nil, err
	}
	return &rpcapi.ExploreRes{Keys: keys}, nil
}
func (inn *LocalNode) RPCMeta(ctx context.Context, in *rpcapi.Empty) (*rpcapi.NodeMeta, error) {
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
	meta := inn.meta()
	return meta, nil
}

func (inn *LocalNode) RPCMove(ctx context.Context, in *rpcapi.MoveReq) (*rpcapi.Empty, error) {
	start := time.Now()
	defer func() {
		end := time.Since(start)
		now := time.Now()

		p := influxdb2.NewPoint("grpc",
			map[string]string{
				"id":     inn.id,
				"method": "move",
			},
			map[string]interface{}{
				"duration_ms": end.Milliseconds(),
			},
			now)
		inn.metrics <- p
	}()

	time.Sleep(inn.rpclatency)

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
