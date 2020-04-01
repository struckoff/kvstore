package main

import (
	"encoding/json"
	"flag"
	"github.com/struckoff/kvstore"
	"github.com/struckoff/kvstore/balancer_adapter"
	bolt "go.etcd.io/bbolt"
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
	var conf kvstore.Config

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

	//Initialize database
	db, err := bolt.Open(conf.DBpath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//Initialize balancer
	bal, err := balancer_adapter.NewSFCBalancer(conf)
	if err != nil {
		return err
	}

	//Initialize local node
	mainNode := kvstore.NewInternalNode(conf.Name, conf.Address, conf.RPCAddress, conf.Power, conf.Capacity, db)
	addy, err := net.ResolveTCPAddr("tcp", conf.RPCAddress)
	if err != nil {
		return err
	}
	inbound, err := net.ListenTCP("tcp", addy)
	if err != nil {
		return err
	}
	rpcServer := grpc.NewServer()
	h, err := kvstore.NewHost(mainNode, bal, rpcServer)
	if err != nil {
		return err
	}

	//if len(conf.Entrypoints) > 0 {
	//	if err := h.Lookup(conf.Entrypoints); err != nil {
	//		return err
	//	}
	//}

	//Run API servers
	errCh := make(chan error)
	go func(errCh chan error) {
		if err := h.RunServer(conf.Address); err != nil {
			errCh <- err
			return
		}
	}(errCh)
	go func(errCh chan error) {
		if err := rpcServer.Serve(inbound); err != nil {
			errCh <- err
			return
		}
	}(errCh)
	go func(errCh chan error) {
		if err := h.RunConsul(&conf); err != nil {
			errCh <- err
			return
		}
	}(errCh)
	return <-errCh
}
