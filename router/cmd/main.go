package main

import (
	"encoding/json"
	"flag"
	"github.com/kelseyhightower/envconfig"
	"github.com/struckoff/kvstore/router"
	"github.com/struckoff/kvstore/router/rpcapi"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var conf router.Config
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
	go func(errCh chan error, h *router.Host, conf *router.Config) {
		if err := RunRPCServer(h, conf); err != nil {
			errCh <- err
			return
		}
	}(errCh, h, &conf)

	return <-errCh
}

func RunRPCServer(h rpcapi.RPCBalancerServer, conf *router.Config) error {
	addy, err := net.ResolveTCPAddr("tcp", conf.RPCAddress)
	if err != nil {
		return err
	}
	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	rpcapi.RegisterRPCBalancerServer(s, h)

	log.Printf("RUN RPC Server [%s]", conf.RPCAddress)
	if err := s.Serve(inbound); err != nil {
		return err
	}
	return nil
}
