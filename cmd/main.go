package main

import (
	"encoding/json"
	"flag"
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/SFCFramework/optimizer"
	"github.com/struckoff/SFCFramework/transform"
	"github.com/struckoff/kvstore"
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
	cfgPath := flag.String("c", "config.json", "path to config file")
	flag.Parse()
	var conf kvstore.Config
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
	bal, err := balancer.NewBalancer(conf.Curve.CurveType, conf.Dimensions, conf.Size, transform.KVTransform, optimizer.RangeOptimizer, nil)
	if err != nil {
		return err
	}

	//Initialize local node
	mainNode := kvstore.NewInternalNode(conf.Name, conf.Address, conf.Power, conf.Capacity, db)
	h, err := kvstore.NewHost(mainNode, bal)
	if err != nil {
		return err
	}

	if len(conf.Entrypoints) > 0 {
		if err := h.Lookup(conf.Entrypoints); err != nil {
			return err
		}
	}

	//Run API server
	if err := h.RunServer(conf.Address); err != nil {
		return err
	}
	return nil
}
