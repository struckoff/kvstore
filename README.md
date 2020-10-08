#KV Store

## Descriptor
Proof of work for [github.com/struckoff/sfcframework](http://github.com/struckoff/sfcframework)

Distributed key-value storage of geo data. 

##Router
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
    
Store binary data by given key.
```
GET /get/:keys
```
    
Return data by given slash separated keys.
 
```
GET /list
```   
Return list of keys.

###Config
####Example
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

##Store
Provides gRPC interface for bbolt DB to store data at the local file system.  

###Config
#### Mode
 - standalone - will connect nodes between each other without external service. 
 - kvrouter - will try to use [router](#router) as discovery service.
 - consul - will use consul as discovery service.
 
 ####Example
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
