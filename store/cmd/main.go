package main

import (
	"flag"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api/write"
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/logger"
	"github.com/struckoff/kvstore/router"
	"github.com/struckoff/kvstore/router/balanceradapter"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/dataitem"
	"github.com/struckoff/kvstore/router/nodehasher"
	"github.com/struckoff/kvstore/store"
	"github.com/struckoff/sfcframework/curve"
	bolt "go.etcd.io/bbolt"
	"go.uber.org/zap"
	"log"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//log.SetOutput(ioutil.Discard)
	if err := run(); err != nil {
		panic(err)
	}
}

func run() (err error) {
	var conf store.Config
	var inn *store.LocalNode
	// If config implies use of consul, consul agent name will be used as name.
	// Otherwise, hostname will be used instead.
	errCh := make(chan error)

	cfgPath := flag.String("c", "config.json", "path to config file")
	flag.Parse()
	conf, err = store.ReadConfig(*cfgPath)
	if err != nil {
		return errors.Wrap(err, "failed to parse config")
	}

	//Initialize database
	db, err := bolt.Open(conf.DBpath, 0600, nil)
	if err != nil {
		return errors.Wrap(err, "failed to open database")
	}
	defer db.Close()

	//Run metrics client
	metrics := make(chan *write.Point, 100)
	go func(points <-chan *write.Point) {
		client := influxdb2.NewClientWithOptions(
			conf.InfluxAddress,
			"",
			influxdb2.DefaultOptions().SetBatchSize(20))
		defer client.Close()
		writeApi := client.WriteApi("", "kvstore/autogen")
		defer writeApi.Flush()
		balancermode, err := conf.Balancer.Mode.String()
		if err != nil {
			logger.Logger().Error("unable to run metrics client", zap.Error(err))
			return
		}
		for p := range points {
			p.AddTag("balancermode", balancermode)
			writeApi.WritePoint(p)
		}
	}(metrics)

	// create balancer
	switch conf.Mode {
	case store.StandaloneMode, store.ConsulMode:
		var bal balanceradapter.Balancer
		var hr nodehasher.Hasher
		switch conf.Balancer.Mode {
		case config.ConsistentMode:
			bal = balanceradapter.NewConsistentBalancer()
		case config.SFCMode:
			bal, err = balanceradapter.NewSFCBalancer(conf.Balancer)
			if err != nil {
				return err
			}
		}

		// create node hasher
		switch conf.Balancer.NodeHash {
		case config.GeoSfc:
			sb := bal.(*balanceradapter.SFC)
			sfc, err := curve.NewCurve(conf.Balancer.SFC.Curve.CurveType, 2, sb.SFC().Bits())
			if err != nil {
				return errors.Wrap(err, "failed to create curve")
			}
			hr = nodehasher.NewGeoSfc(sfc)
		case config.XXHash:
			hr = nodehasher.NewXXHash()
		default:
			return errors.New("invalid node hasher")
		}
		ndf, err := dataitem.GetDataItemFunc(conf.Balancer.DataMode)
		if err != nil {
			return err
		}
		kvr, err := router.NewRouter(bal, hr, ndf, conf.Balancer)
		if err != nil {
			return errors.Wrap(err, "failed to initialize router")
		}
		//Initialize local node
		inn, err = store.NewLocalNode(&conf, db, kvr, metrics)
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
		inn, err = store.NewLocalNode(&conf, db, nil, metrics)
		if err != nil {
			return err
		}
	default:
		return errors.New("wrong node node")
	}

	if err := inn.RunRPCServer(&conf, errCh); err != nil {
		//errCh <- errors.Wrap(err, "failed to run RPC server")
		return errors.Wrap(err, "failed to run RPC server")
	}

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
