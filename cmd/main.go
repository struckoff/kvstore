package main

import (
	"encoding/json"
	"flag"
	"github.com/pkg/errors"
	"github.com/struckoff/kvrouter/balancer_adapter"
	kvrouter "github.com/struckoff/kvrouter/router"
	"github.com/struckoff/kvstore"
	bolt "go.etcd.io/bbolt"
	"os"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var conf kvstore.Config
	var inn *kvstore.InternalNode
	// If config implies use of consul, consul agent name  will be  used as name.
	// Otherwise, hostname will be used instead.
	errCh := make(chan error)

	cfgPath := flag.String("c", "config.json", "path to config file")
	flag.Parse()
	configFile, err := os.Open(*cfgPath)
	if err != nil {
		return errors.Wrap(err, "failed to open config file")
	}
	defer configFile.Close()
	if err := json.NewDecoder(configFile).Decode(&conf); err != nil {
		return errors.Wrap(err, "failed to parse config file")
	}

	if err := conf.Prepare(); err != nil {
		return err
	}

	//Initialize database
	db, err := bolt.Open(conf.DBpath, 0600, nil)
	if err != nil {
		return errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	switch conf.Mode {
	case kvstore.StandaloneMode, kvstore.ConsulMode:
		bal, err := balancer_adapter.NewSFCBalancer(conf.Balancer)
		if err != nil {
			return err
		}
		kvr, err := kvrouter.NewRouter(bal)
		if err != nil {
			return errors.Wrap(err, "failed to initialize router")
		}
		//Initialize local node1
		inn = kvstore.NewInternalNode(&conf, db, kvr)

		//Run API servers
		go func(errCh chan error, conf *kvstore.Config) {
			if err := inn.RunHTTPServer(conf.Address); err != nil {
				errCh <- errors.Wrap(err, "failed to run HTTP server")
				return
			}
		}(errCh, &conf)
	case kvstore.KvrouterMode:
		inn = kvstore.NewInternalNode(&conf, db, nil)
	}

	go func(errCh chan error, conf *kvstore.Config) {
		if err := inn.RunRPCServer(conf); err != nil {
			errCh <- errors.Wrap(err, "failed to run RPC server")
			return
		}
	}(errCh, &conf)

	//Run discovery connection
	go func(errCh chan error, inn *kvstore.InternalNode, conf *kvstore.Config) {
		ds := discoveryService(conf.Mode, inn)
		if err := ds(conf); err != nil {
			errCh <- errors.Wrap(err, "failed to run discovery")
			return
		}
	}(errCh, inn, &conf)

	return <-errCh
}

func discoveryService(mode kvstore.DiscoverMode, inn *kvstore.InternalNode) func(conf *kvstore.Config) error {
	switch mode {
	case kvstore.KvrouterMode:
		return inn.RunKVRouter
	case kvstore.ConsulMode:
		return inn.RunConsul
	}
	return nil
}
