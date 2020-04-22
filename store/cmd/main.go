package main

import (
	"flag"
	"github.com/pkg/errors"
	"github.com/struckoff/SFCFramework/curve"
	"github.com/struckoff/kvstore/router"
	"github.com/struckoff/kvstore/router/nodehasher"
	"github.com/struckoff/kvstore/store"
	bolt "go.etcd.io/bbolt"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var conf store.Config
	var err error
	var inn *store.LocalNode
	// If config implies use of consul, consul agent name  will be  used as name.
	// Otherwise, hostname will be used instead.
	errCh := make(chan error)

	cfgPath := flag.String("c", "config.json", "path to config file")
	flag.Parse()
	conf, err = store.ReadConfig(*cfgPath)
	if err != nil {
		return err
	}

	//Initialize database
	db, err := bolt.Open(conf.DBpath, 0600, nil)
	if err != nil {
		return errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	switch conf.Mode {
	case store.StandaloneMode, store.ConsulMode:
		bal, err := router.NewSFCBalancer(conf.Balancer)
		if err != nil {
			return err
		}
		var hr nodehasher.Hasher
		switch conf.Balancer.NodeHash {
		case router.GeoSfc:
			sfc, err := curve.NewCurve(conf.Balancer.Curve.CurveType, 2, bal.SFC().Bits())
			if err != nil {
				return errors.Wrap(err, "failed to create curve")
			}
			hr = nodehasher.NewGeoSfc(sfc)
		case router.XXHash:
			hr = nodehasher.NewXXHash()
		default:
			return errors.New("invalid node hasher")
		}
		kvr, err := router.NewRouter(bal, hr)
		if err != nil {
			return errors.Wrap(err, "failed to initialize router")
		}
		//Initialize local node1
		inn, err = store.NewLocalNode(&conf, db, kvr)
		if err != nil {
			return err
		}

		//Run API servers
		go func(errCh chan error, conf *store.Config) {
			if err := inn.RunHTTPServer(conf.Address); err != nil {
				errCh <- errors.Wrap(err, "failed to run HTTP server")
				return
			}
		}(errCh, &conf)
	case store.KvrouterMode:
		inn, err = store.NewLocalNode(&conf, db, nil)
		if err != nil {
			return err
		}
	}

	go func(errCh chan error, conf *store.Config) {
		if err := inn.RunRPCServer(conf); err != nil {
			errCh <- errors.Wrap(err, "failed to run RPC server")
			return
		}
	}(errCh, &conf)

	//Run discovery connection
	go func(errCh chan error, inn *store.LocalNode, conf *store.Config) {
		ds := discoveryService(conf.Mode, inn)
		if err := ds(conf); err != nil {
			errCh <- errors.Wrap(err, "failed to run discovery")
			return
		}
	}(errCh, inn, &conf)

	return <-errCh
}

func discoveryService(mode store.DiscoverMode, inn *store.LocalNode) func(conf *store.Config) error {
	switch mode {
	case store.KvrouterMode:
		return inn.RunRouter
	case store.ConsulMode:
		return inn.RunConsul
	}
	return nil
}
