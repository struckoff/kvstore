package main

import "log"
import bolt "go.etcd.io/bbolt"

const addr = "0.0.0.0:9191"
const name = "node-9191"
const dbpath = "bl.db"

var mainBucket = []byte("pairs")
var mainNode Node

func main() {
	db, err := bolt.Open(dbpath, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mainNode = NewInternalNode(name, addr, 1, 1, db)
	if err := RunServer(addr); err != nil {
		panic(err)
	}
}

//! STUB
func GetNode(key string) Node {
	return mainNode
}
