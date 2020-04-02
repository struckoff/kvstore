package main

import (
	"encoding/json"
	"flag"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/struckoff/kvstore"
	"github.com/struckoff/kvstore/balancer_adapter"
	bolt "go.etcd.io/bbolt"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var conf kvstore.Config
	var name string
	errCh := make(chan error)

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

	// Initialize consul client id config allows
	var consul *consulapi.Client
	if conf.Consul == nil {
		name, err = os.Hostname()
		if err != nil {
			return err
		}
	} else {
		consul, err = consulapi.NewClient(&conf.Consul.Config)
		if err != nil {
			return err
		}
		name, err = consul.Agent().NodeName()
		if err != nil {
			return err
		}
	}

	//Initialize local node
	mainNode := kvstore.NewInternalNode(name, conf.Address, conf.RPCAddress, conf.Power, conf.Capacity, db)

	h, err := kvstore.NewHost(mainNode, bal, consul)
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
	go func(errCh chan error) {
		if err := h.RunRPCServer(&conf); err != nil {
			errCh <- err
			return
		}
	}(errCh)
	if conf.Consul != nil {
		go func(errCh chan error) {
			if err := h.RunConsul(&conf); err != nil {
				errCh <- err
				return
			}
		}(errCh)
	}
	return <-errCh
}
