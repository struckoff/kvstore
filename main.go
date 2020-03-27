package main

import (
	"encoding/json"
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/SFCFramework/optimizer"
	"github.com/struckoff/SFCFramework/transform"
	"github.com/struckoff/kvstore/config"
	"github.com/struckoff/kvstore/host"
	"github.com/struckoff/kvstore/node"
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
	var conf config.Config
	configP := "config.json"
	if len(os.Args)>1{
		configP = os.Args[1]
	}
	configFile, err := os.Open(configP)
    defer configFile.Close()
    if err != nil {
       return err
    }
    if err := json.NewDecoder(configFile).Decode(&conf); err != nil{
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
	mainNode := node.NewInternalNode(conf.Name, conf.Address, conf.Power, conf.Capacity, db)
	h, err := host.NewHost(mainNode, bal)
	if err != nil {
		return err
	}

	if len(conf.Entrypoints) > 0{
		if err := h.Lookup(conf.Entrypoints); err != nil{
			return err
		}
	}

	//Run API server
	if err := h.RunServer(conf.Address); err != nil {
		return err
	}
	return nil
}
