module github.com/struckoff/kvstore

go 1.14

require (
	github.com/armon/go-metrics v0.3.3 // indirect
	github.com/golang/protobuf v1.3.5 // indirect
	github.com/hashicorp/consul/api v1.4.0
	github.com/hashicorp/go-hclog v0.12.2 // indirect
	github.com/hashicorp/go-immutable-radix v1.2.0 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/serf v0.9.0 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/mattn/go-colorable v0.1.6 // indirect
	github.com/mitchellh/mapstructure v1.2.2 // indirect
	github.com/pkg/errors v0.9.1
	github.com/struckoff/SFCFramework v0.0.0-20200405132449-c125da4b1018
	go.etcd.io/bbolt v1.3.4
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e // indirect
	golang.org/x/sys v0.0.0-20200406155108-e3b113bbe6a4 // indirect
	google.golang.org/genproto v0.0.0-20200407120235-9eb9bb161a06 // indirect
	google.golang.org/grpc v1.28.1
)
replace github.com/struckoff/kvrouter => /home/struckoff/Projects/Go/src/github.com/struckoff/kvrouter
replace github.com/struckoff/kvrouter/rpcapi => /home/struckoff/Projects/Go/src/github.com/struckoff/kvrouter/rpacapi
replace github.com/struckoff/SFCFramework => /home/struckoff/Projects/Go/src/github.com/struckoff/SFCFramework