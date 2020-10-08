package main

import (
	"encoding/json"
	"flag"
	"github.com/kelseyhightower/envconfig"
	"github.com/struckoff/kvstore/logger"
	"github.com/struckoff/kvstore/router"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/rpcapi"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"os"
	"time"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var conf config.Config
	errCh := make(chan error)
	// If config implies use of consul, consul agent name  will be  used as name.
	// Otherwise, hostname will be used instead.
	cfgPath := flag.String("c", "config.json", "path to config file")
	flag.Parse()
	configFile, err := os.Open(*cfgPath)
	if err != nil {
		return err
	}
	defer configFile.Close()
	if err := json.NewDecoder(configFile).Decode(&conf); err != nil {
		return err
	}

	if err := envconfig.Process("KVROUTER", &conf); err != nil {
		return err
	}

	h, err := router.NewHost(&conf)
	if err != nil {
		return err
	}

	//Run API servers
	go func(errCh chan error) {
		if err := h.RunHTTPServer(conf.Address); err != nil {
			errCh <- err
			return
		}
	}(errCh)
	go func(errCh chan error, h *router.Host, conf *config.Config) {
		if err := RunRPCServer(h, conf); err != nil {
			errCh <- err
			return
		}
	}(errCh, h, &conf)

	return <-errCh
}

func RunRPCServer(h rpcapi.RPCBalancerServer, conf *config.Config) error {
	addy, err := net.ResolveTCPAddr("tcp", conf.RPCAddress)
	if err != nil {
		return err
	}
	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		return err
	}
	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 5 * time.Minute,
		}),
	)
	rpcapi.RegisterRPCBalancerServer(s, h)

	logger.Logger().Info("RUN RPC Server", zap.String("address", conf.RPCAddress))
	if err := s.Serve(inbound); err != nil {
		return err
	}
	return nil
}
