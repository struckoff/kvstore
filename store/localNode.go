package store

import (
	"context"
	"github.com/influxdata/influxdb-client-go/api/write"
	"github.com/struckoff/kvstore/logger"
	"github.com/struckoff/kvstore/router/dataitem"
	"github.com/struckoff/kvstore/router/nodes"
	"github.com/struckoff/sfcframework/node"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router"
	"github.com/struckoff/kvstore/router/rpcapi"
	bolt "go.etcd.io/bbolt"
	"google.golang.org/grpc"
)

var mainBucket = []byte("pairs")

// Return new instance LocalNode.
func NewLocalNode(conf *Config, ndf dataitem.NewDataItemFunc, db *bolt.DB, kvr *router.Router, metrics chan<- *write.Point) (*LocalNode, error) {
	lwID := int64(0)
	ln := &LocalNode{
		id:          *conf.Name,
		address:     conf.Address,
		rpcaddress:  conf.RPCAddress,
		p:           nodes.NewPower(conf.Power),
		c:           NewCapacity(conf.Capacity),
		db:          db,
		kvr:         kvr,
		geo:         conf.Geo,
		rpclatency:  conf.Latency.Duration,
		httplatency: conf.Balancer.Latency.Duration,
		lwID:        &lwID,
		metrics:     metrics,
		ndf:         ndf,
	}
	if ln.kvr != nil {
		h, err := ln.kvr.Hasher().Sum(ln.meta(context.Background()))
		if err != nil {
			return nil, errors.Wrap(err, "failed to calculate hash sum")
		}
		ln.h = h
	}
	if conf.Mode == ConsulMode {
		consul, err := consulapi.NewClient(&conf.Consul.Config)
		if err != nil {
			return nil, err
		}
		ln.consul = consul
	}
	dis, err := ln.Explore(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to explore local node")
	}
	ln.c.Add(-float64(len(dis))) // reduce capacity
	return ln, nil
}

// LocalNode represents local node
type LocalNode struct {
	//mu          sync.RWMutex
	id          string
	address     string
	rpcaddress  string
	rpcserver   *grpc.Server
	p           nodes.Power
	c           Capacity
	db          *bolt.DB
	kvr         *router.Router
	consul      *consulapi.Client
	kvrAgent    rpcapi.RPCBalancerClient
	geo         *rpcapi.GeoData
	h           uint64
	rpclatency  time.Duration
	httplatency time.Duration
	lwID        *int64
	metrics     chan<- *write.Point
	ndf         dataitem.NewDataItemFunc
}

func (inn *LocalNode) RunHTTPServer(addr string) error {
	h := inn.kvr.HTTPHandler()
	l := router.LatencyMiddleware(h, inn.httplatency)
	logger.Logger().Info("HTTP server listening", zap.String("address", addr))
	if err := http.ListenAndServe(addr, l); err != nil {
		return err
	}
	return nil
}

//ID returns the node's ID
func (inn *LocalNode) ID() string {
	//inn.mu.RLock()
	//defer inn.mu.RUnlock()
	return inn.id
}

//RPCAddress returns the node's rpc address
func (inn *LocalNode) RPCAddress() string {
	//inn.mu.RLock()
	//defer inn.mu.RUnlock()
	return inn.rpcaddress
}

//HTTPAddress returns the node's http address
func (inn *LocalNode) HTTPAddress() string {
	//inn.mu.RLock()
	//defer inn.mu.RUnlock()
	return inn.address
}

//Power returns the node's power
func (inn *LocalNode) Power() node.Power {
	//inn.mu.RLock()
	//defer inn.mu.RUnlock()
	return inn.p
}

//Capacity returns the node's capacity
func (inn *LocalNode) Capacity() nodes.Capacity {
	//inn.mu.RLock()
	//defer inn.mu.RUnlock()
	return &inn.c
}

func (inn *LocalNode) Hash() uint64 {
	//inn.mu.RLock()
	//defer inn.mu.RUnlock()
	return inn.h
}

func (inn *LocalNode) SetHash(h uint64) {
	//inn.mu.Lock()
	//defer inn.mu.Unlock()
	inn.h = h
}

// Store value for a given key in local storage
func (inn *LocalNode) Store(ctx context.Context, kv *rpcapi.KeyValue) (*rpcapi.DataItem, error) {
	err := inn.db.Update(func(tx *bolt.Tx) error {
		bc, err := tx.CreateBucketIfNotExists(mainBucket)
		if err != nil {
			return err
		}
		return bc.Put([]byte(kv.Key.ID), []byte(kv.Value))
	})
	if err != nil {
		return nil, err
	}
	inn.c.Add(-1) // reduce capacity
	di := kv.Key
	di.Size = uint64(len([]byte(kv.Value)))
	return di, nil
}

// Store KV pairs in local storage
func (inn *LocalNode) StorePairs(ctx context.Context, pairs []*rpcapi.KeyValue) ([]*rpcapi.DataItem, error) {
	cp := 0.0
	res := make([]*rpcapi.DataItem, len(pairs))
	err := inn.db.Batch(func(tx *bolt.Tx) error {
		bc, err := tx.CreateBucketIfNotExists(mainBucket)
		if err != nil {
			return err
		}
		for i := range pairs {
			b := []byte(pairs[i].Value)
			size := len(b)
			di := pairs[i].Key
			di.Size = uint64(size)
			res[i] = di
			if err := bc.Put([]byte(di.ID), b); err != nil {
				return errors.Wrap(err, "failed to store pair")
			}
			cp++
		}
		return nil
	})
	inn.c.Add(-cp) // reduce capacity
	return res, err
}

// Return value for a given key from local storage
func (inn *LocalNode) Receive(ctx context.Context, dis []*rpcapi.DataItem) (*rpcapi.KeyValues, error) {
	kvs := &rpcapi.KeyValues{
		KVs: make([]*rpcapi.KeyValue, len(dis)),
	}
	err := inn.db.View(func(tx *bolt.Tx) error {
		bc := tx.Bucket(mainBucket)
		if bc == nil {
			return errors.New("unable to receive value, bucket not found")
		}
		for i := range dis {
			val := bc.Get([]byte(dis[i].ID))
			ok := val != nil
			kvs.KVs[i] = &rpcapi.KeyValue{
				Key:   dis[i],
				Value: val,
				Found: ok,
			}
		}
		return nil
	})
	return kvs, err
}

// Remove value for a given key
func (inn *LocalNode) Remove(ctx context.Context, rdis []*rpcapi.DataItem) (dis []*rpcapi.DataItem, err error) {
	err = inn.db.Update(func(tx *bolt.Tx) error {
		bc := tx.Bucket(mainBucket)
		if bc == nil {
			return nil
		}
		for i := range rdis {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				size := uint64(len(bc.Get([]byte(rdis[i].ID))))
				if err := bc.Delete([]byte(rdis[i].ID)); err != nil {
					return err
				}

				rdis[i].Size = size
				dis = append(dis, rdis[i])
			}
		}
		return nil
	})
	//inn.c.Add(cp) // increase capacity

	if err != nil {
		return nil, errors.Wrap(err, "failed to remove key")
	}
	return
}

// Move values for a given keys to another node
func (inn *LocalNode) Move(ctx context.Context, nk map[nodes.Node][]*rpcapi.DataItem) error {
	var wg sync.WaitGroup
	for en, dis := range nk {
		if len(dis) == 0 {
			continue
		}
		wg.Add(1)
		go func(ctx context.Context, en nodes.Node, dis []*rpcapi.DataItem, wg *sync.WaitGroup) {
			defer wg.Done()
			pairs := make([]*rpcapi.KeyValue, len(dis))
			err := inn.db.View(func(tx *bolt.Tx) error {
				bc := tx.Bucket(mainBucket)
				if bc == nil {
					return nil
				}
				for i := range dis {
					body := bc.Get(dis[i].ID)
					pairs[i] = &rpcapi.KeyValue{Key: dis[i], Value: body}
				}
				return nil
			})
			if err != nil {
				logger.Logger().Error("error moving keys", zap.Error(err), zap.String("Node", en.ID()))
				return
			}
			_, err = en.StorePairs(ctx, pairs)
			if err != nil {
				logger.Logger().Error("error moving keys", zap.Error(err), zap.String("Node", en.ID()))
				return
			}
			err = inn.db.Batch(func(tx *bolt.Tx) error {
				bc := tx.Bucket(mainBucket)
				if bc == nil {
					return nil
				}
				for i := range dis {
					if err := bc.Delete(dis[i].ID); err != nil {
						return errors.Wrap(err, "failed to delete keys")
					}
				}
				return nil
			})
			if err != nil {
				logger.Logger().Error("error moving keys", zap.Error(err), zap.String("Node", en.ID()))
				return
			}
			logger.Logger().Debug("keys relocated to node", zap.Int("Amount", len(dis)), zap.String("Node", en.ID()))
		}(ctx, en, dis, &wg)
	}
	wg.Wait()
	return nil
}

// Explore returns the list of keys in local storage.
func (inn *LocalNode) Explore(ctx context.Context) ([]*rpcapi.DataItem, error) {
	res := make([]*rpcapi.DataItem, 0)
	err := inn.db.View(func(tx *bolt.Tx) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			b := tx.Bucket(mainBucket)
			if b == nil {
				//return errors.New("bucket not found")
				return nil
			}
			err := b.ForEach(func(k, v []byte) error {
				di, err := inn.ndf(string(k), uint64(len(v)))
				if err != nil {
					return err
				}
				res = append(res, di.RPCApi())
				return nil
			})
			return err
		}
	})
	return res, err
}

// Return meta information about the node
func (inn *LocalNode) Meta(ctx context.Context) *rpcapi.NodeMeta {
	//inn.mu.RLock()
	//defer inn.mu.RUnlock()
	return inn.meta(ctx)
}

func (inn *LocalNode) meta(_ context.Context) *rpcapi.NodeMeta {
	cp, _ := inn.c.Get()
	return &rpcapi.NodeMeta{
		ID:         inn.ID(),
		Address:    inn.HTTPAddress(),
		RPCAddress: inn.RPCAddress(),
		Power:      inn.Power().Get(),
		Capacity:   cp,
		Geo:        inn.geo,
	}
}
