module github.com/struckoff/kvstore

go 1.14

require (
	github.com/buraksezer/consistent v0.0.0-20191006190839-693edf70fd72
	github.com/cespare/xxhash v1.1.0
	github.com/golang/protobuf v1.4.2
	github.com/hashicorp/consul/api v1.5.0
	github.com/influxdata/influxdb-client-go v1.3.0
	github.com/julienschmidt/httprouter v1.3.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/lafikl/consistent v0.0.0-20190331123054-b5c3ef09639f
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1
	github.com/struckoff/SFCFramework v0.0.0-20200811232013-93d2eb2c5003
	go.etcd.io/bbolt v1.3.5
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	google.golang.org/grpc v1.30.0
)

//replace github.com/struckoff/SFCFramework => /home/struckoff/Projects/Go/src/github.com/struckoff/SFCFramework
