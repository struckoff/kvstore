![Go](https://github.com/struckoff/kvstore/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/struckoff/kvstore)](https://goreportcard.com/report/github.com/struckoff/kvstore)

# KV Store

## Description
Proof of work for [github.com/struckoff/sfcframework](http://github.com/struckoff/sfcframework)

Distributed key-value storage of geo data. It uses longitude and latitude extracted from key to find appropriate node on the space-filling curve. 
When new node appeared in the system, router automatically run redistributing process.
It also works when each node responsible for distribution due to stateless mechanism of sfcframework which 
provides hashing functionality with same results on each node. 


## Router
/router service provides data balancing logic for the system.

It also provides HTTP API.
Router could be run with [store](#store) as one service or separate.

### HTTP API
```
GET /nodes
```
    
Get list of nodes.
```
POST /put/:key
```

####Example
```zsh
echo "test-data" | http POST ":9190/put/{\"Lon\":-4,\"Lat\":-20}"      
HTTP/1.1 200 OK
Content-Length: 2
Content-Type: text/plain; charset=utf-8
Date: Thu, 08 Oct 2020 22:24:00 GMT

OK

```

Store binary data by given key.
```
GET /get/:keys
```
####Example
```zsh
http GET ":9190/get/{\"Lon\":-4,\"Lat\":-20}"
        
HTTP/1.1 200 OK
Content-Length: 59
Content-Type: text/plain; charset=utf-8

[
    {
        "Key": "{\"Lon\":-4,\"Lat\":-20}",
        "Value": "test-data"
    }
]

```
Return data by given slash separated keys.
 
```
GET /list
```   
Return list of keys.

### Config
#### config.json
```json
{
  "Address": "0.0.0.0:9190",
  "RPCAddress": "127.0.0.1:9290",
  "Balancer":{
    "State": true,
    "Mode": "SFC",
    "DataMode": "geo",
    "SFC": {
      "Dimensions":2,
      "Size":256,
      "Curve": "morton"
    },
    "NodeHash": "geosfc"
  }
}
```

## Store
Provides gRPC interface for bbolt DB to store data at the local file system.  

### Config
#### Mode
 - standalone - will connect nodes between each other without external service. 
 - kvrouter - will try to use [router](#router) as discovery service.
 - consul - will use consul as discovery service.
 
 #### config.json
```json
{
  "Mode": "consul",
  "DBPath": "/data/data.db",
  "Power": 1,
  "Capacity": 100000,
  "Balancer":{
    "Mode": "SFC",
    "DataMode": "geo",
    "SFC": {
      "Dimensions":2,
      "Size":64,
      "Curve": "morton"
    },
    "NodeHash": "geosfc"
  },
  "Health": {
    "CheckInterval": "10s",
    "CheckTimeout": "10s",
    "DeregisterCriticalServiceAfter": "10s"
  },
  "Geo":{
    "Latitude":0.0,
    "Longitude":0.0
  },
  "KVRouter": {
    "Address": "127.0.0.1:9290"
  },
  "Consul": {
    "Service": "kvstore"
  }
}
```
### Environment
Configuration is also possible with environment variables.

#### Example
```sh
KVSTORE_NAME=node-0
KVSTORE_MODE=kvrouter
KVSTORE_CAPACITY=1500
KVSTORE_KVROUTER_ADDRESS=localhost:9290
KVSTORE_RPC_ADDRESS=localhost:9293
KVSTORE_INFLUX_ADDRESS=http://127.0.0.1:8086
KVSTORE_GEO_LONGITUDE=-180
KVSTORE_GEO_LATITUDE=-90
KVSTORE_RPC_LATENCY=150ms
KVSTORE_DBPATH=/var/lib/kvstore/data/node-0/data.db
```
