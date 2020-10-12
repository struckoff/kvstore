package store

import (
	"context"
	"errors"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/influxdata/influxdb-client-go/api/write"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/struckoff/kvstore/mocks"
	"github.com/struckoff/kvstore/router"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/dataitem"
	"github.com/struckoff/kvstore/router/nodes"
	"github.com/struckoff/kvstore/router/rpcapi"
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestInternalNode_Meta(t *testing.T) {
	type fields struct {
		id         string
		address    string
		rpcaddress string
		p          nodes.Power
		c          Capacity
		db         *bolt.DB
	}
	tests := []struct {
		name   string
		fields fields
		want   *rpcapi.NodeMeta
	}{
		{
			name: "test",
			fields: fields{
				id:         "test_id",
				address:    "test_addr",
				rpcaddress: "test_raddr",
				p:          nodes.NewPower(1.1),
				c:          NewCapacity(2.3),
				db:         nil,
			},
			want: &rpcapi.NodeMeta{
				ID:         "test_id",
				Address:    "test_addr",
				RPCAddress: "test_raddr",
				Power:      1.1,
				Capacity:   2.3,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &LocalNode{
				id:         tt.fields.id,
				address:    tt.fields.address,
				rpcaddress: tt.fields.rpcaddress,
				p:          tt.fields.p,
				c:          tt.fields.c,
				db:         tt.fields.db,
			}
			got := n.Meta(context.Background())
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInternalNode_Capacity(t *testing.T) {
	type fields struct {
		c Capacity
	}
	tests := []struct {
		name   string
		fields fields
		want   Capacity
	}{
		{
			name: "test",
			fields: fields{
				c: NewCapacity(2.3),
			},
			want: NewCapacity(2.3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &LocalNode{
				c: tt.fields.c,
			}

			got := n.Capacity()
			assert.Equal(t, &tt.want, got)
		})
	}
}

func TestInternalNode_HTTPAddress(t *testing.T) {
	type fields struct {
		id         string
		address    string
		rpcaddress string
		p          nodes.Power
		c          Capacity
		db         *bolt.DB
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test",
			fields: fields{
				id:         "test_id",
				address:    "test_addr",
				rpcaddress: "test_raddr",
				p:          nodes.NewPower(1.1),
				c:          NewCapacity(2.3),
				db:         nil,
			},
			want: "test_addr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &LocalNode{
				id:         tt.fields.id,
				address:    tt.fields.address,
				rpcaddress: tt.fields.rpcaddress,
				p:          tt.fields.p,
				c:          tt.fields.c,
				db:         tt.fields.db,
			}
			got := n.HTTPAddress()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInternalNode_ID(t *testing.T) {
	type fields struct {
		id         string
		address    string
		rpcaddress string
		p          nodes.Power
		c          Capacity
		db         *bolt.DB
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test",
			fields: fields{
				id:         "test_id",
				address:    "test_addr",
				rpcaddress: "test_raddr",
				p:          nodes.NewPower(1.1),
				c:          NewCapacity(2.3),
				db:         nil,
			},
			want: "test_id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &LocalNode{
				id:         tt.fields.id,
				address:    tt.fields.address,
				rpcaddress: tt.fields.rpcaddress,
				p:          tt.fields.p,
				c:          tt.fields.c,
				db:         tt.fields.db,
			}
			if got := n.ID(); got != tt.want {
				t.Errorf("ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInternalNode_Power(t *testing.T) {
	type fields struct {
		id         string
		address    string
		rpcaddress string
		p          nodes.Power
		c          Capacity
		db         *bolt.DB
	}
	tests := []struct {
		name   string
		fields fields
		want   nodes.Power
	}{
		{
			name: "test",
			fields: fields{
				id:         "test_id",
				address:    "test_addr",
				rpcaddress: "test_raddr",
				p:          nodes.NewPower(1.1),
				c:          NewCapacity(2.3),
				db:         nil,
			},
			want: nodes.NewPower(1.1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &LocalNode{
				id:         tt.fields.id,
				address:    tt.fields.address,
				rpcaddress: tt.fields.rpcaddress,
				p:          tt.fields.p,
				c:          tt.fields.c,
				db:         tt.fields.db,
			}
			if got := n.Power(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Power() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInternalNode_RPCAddress(t *testing.T) {
	type fields struct {
		id         string
		address    string
		rpcaddress string
		p          nodes.Power
		c          Capacity
		db         *bolt.DB
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test",
			fields: fields{
				id:         "test_id",
				address:    "test_addr",
				rpcaddress: "test_raddr",
				p:          nodes.NewPower(1.1),
				c:          NewCapacity(2.3),
				db:         nil,
			},
			want: "test_raddr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &LocalNode{
				id:         tt.fields.id,
				address:    tt.fields.address,
				rpcaddress: tt.fields.rpcaddress,
				p:          tt.fields.p,
				c:          tt.fields.c,
				db:         tt.fields.db,
			}
			if got := n.RPCAddress(); got != tt.want {
				t.Errorf("RPCAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInternalNode_StoreExploreRemove(t *testing.T) {
	type args struct {
		kvs        []*rpcapi.KeyValue
		removeKeys []*rpcapi.DataItem
	}
	tests := []struct {
		name                    string
		args                    args
		wantErr                 bool
		wantBeforeRemove        *rpcapi.KeyValues
		wantExploreBeforeRemove []*rpcapi.DataItem
		wantAfterRemove         *rpcapi.KeyValues
		wantExploreAfterRemove  []*rpcapi.DataItem
	}{
		{
			name: "test",
			args: args{
				kvs: []*rpcapi.KeyValue{
					{Key: &rpcapi.DataItem{ID: []byte("t0"), Size: 5}, Value: []byte("t0val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t1"), Size: 5}, Value: []byte("t1val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t2"), Size: 5}, Value: []byte("t2val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t3"), Size: 5}, Value: []byte("t3val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t4"), Size: 5}, Value: []byte("t4val"), Found: true},
				},
				removeKeys: []*rpcapi.DataItem{{ID: []byte("t0")}, {ID: []byte("t2")}, {ID: []byte("t3")}},
			},
			wantErr: false,
			wantBeforeRemove: &rpcapi.KeyValues{
				KVs: []*rpcapi.KeyValue{
					{Key: &rpcapi.DataItem{ID: []byte("t0")}, Value: []byte("t0val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t1")}, Value: []byte("t1val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t2")}, Value: []byte("t2val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t3")}, Value: []byte("t3val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t4")}, Value: []byte("t4val"), Found: true},
				},
			},
			wantExploreBeforeRemove: []*rpcapi.DataItem{
				{ID: []byte("t0"), Size: 5},
				{ID: []byte("t1"), Size: 5},
				{ID: []byte("t2"), Size: 5},
				{ID: []byte("t3"), Size: 5},
				{ID: []byte("t4"), Size: 5},
			},
			wantAfterRemove: &rpcapi.KeyValues{
				KVs: []*rpcapi.KeyValue{
					{Key: &rpcapi.DataItem{ID: []byte("t1"), Size: 5}, Value: []byte("t1val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t4"), Size: 5}, Value: []byte("t4val"), Found: true},
				},
			},
			wantExploreAfterRemove: []*rpcapi.DataItem{
				{ID: []byte("t1"), Size: 5},
				{ID: []byte("t4"), Size: 5},
			},
		},
		{
			name: "remove not existing key",
			args: args{
				kvs: []*rpcapi.KeyValue{
					{Key: &rpcapi.DataItem{ID: []byte("t0")}, Value: []byte("t0val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t1")}, Value: []byte("t1val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t2")}, Value: []byte("t2val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t3")}, Value: []byte("t3val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t4")}, Value: []byte("t4val"), Found: true},
				},
				removeKeys: []*rpcapi.DataItem{{ID: []byte("k0")}, {ID: []byte("k2")}, {ID: []byte("k3")}},
			},
			wantErr: false,
			wantBeforeRemove: &rpcapi.KeyValues{
				KVs: []*rpcapi.KeyValue{
					{Key: &rpcapi.DataItem{ID: []byte("t0"), Size: 5}, Value: []byte("t0val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t1"), Size: 5}, Value: []byte("t1val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t2"), Size: 5}, Value: []byte("t2val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t3"), Size: 5}, Value: []byte("t3val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t4"), Size: 5}, Value: []byte("t4val"), Found: true},
				},
			},
			wantExploreBeforeRemove: []*rpcapi.DataItem{
				{ID: []byte("t0"), Size: 5},
				{ID: []byte("t1"), Size: 5},
				{ID: []byte("t2"), Size: 5},
				{ID: []byte("t3"), Size: 5},
				{ID: []byte("t4"), Size: 5},
			},
			wantAfterRemove: &rpcapi.KeyValues{
				KVs: []*rpcapi.KeyValue{
					{Key: &rpcapi.DataItem{ID: []byte("t0")}, Value: []byte("t0val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t1")}, Value: []byte("t1val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t2")}, Value: []byte("t2val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t3")}, Value: []byte("t3val"), Found: true},
					{Key: &rpcapi.DataItem{ID: []byte("t4")}, Value: []byte("t4val"), Found: true},
				},
			},
			wantExploreAfterRemove: []*rpcapi.DataItem{
				{ID: []byte("t0"), Size: 5},
				{ID: []byte("t1"), Size: 5},
				{ID: []byte("t2"), Size: 5},
				{ID: []byte("t3"), Size: 5},
				{ID: []byte("t4"), Size: 5},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbpath := tempfile()
			defer func(dbpath string) {
				if err := os.Remove(dbpath); err != nil {
					t.Fatal(err)
				}
			}(dbpath)
			db, err := bolt.Open(dbpath, 0666, nil)
			if err != nil {
				t.Fatal(err)
			}
			metrics := make(chan *write.Point)
			defer close(metrics)
			go func(metrics chan *write.Point) {
				for m := range metrics {
					if m == nil {
						return
					}
					continue
				}
			}(metrics)
			n := &LocalNode{
				db:      db,
				metrics: metrics,
				ndf: func(s string, u uint64) (dataitem.DataItem, error) {
					di := &mocks.DataItem{}
					di.On("ID").Return(s)
					di.On("Size").Return(u)
					di.On("RPCApi").Return(&rpcapi.DataItem{
						ID:   []byte(s),
						Size: u,
					})
					return di, nil
				},
			}
			for _, kv := range tt.args.kvs {
				if _, err := n.Store(context.Background(), kv); (err != nil) != tt.wantErr {
					t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			explore, err := n.Explore(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("Before remove: Explore() error = %v, wantErr %v", err, tt.wantErr)
			}

			//sort.Strings(explore)
			//sort.Strings(tt.wantExploreBeforeRemove)
			assert.Equal(t, explore, tt.wantExploreBeforeRemove)

			keys := make([]*rpcapi.DataItem, len(tt.wantBeforeRemove.KVs))
			for i := range tt.wantBeforeRemove.KVs {
				keys[i] = tt.wantBeforeRemove.KVs[i].Key
			}

			kvs, err := n.Receive(context.Background(), keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("Before remove: Receive error = %w", err)
			}

			assert.Equal(t, tt.wantBeforeRemove, kvs)

			_, err = n.Remove(context.Background(), tt.args.removeKeys)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			explore, err = n.Explore(context.Background())
			if (err != nil) != tt.wantErr {
				t.Errorf("After remove: Explore() error = %v, wantErr %v", err, tt.wantErr)
			}

			//sort.Strings(explore)
			//sort.Strings(tt.wantExploreAfterRemove)
			assert.Equal(t, tt.wantExploreAfterRemove, explore)

			kvs, err = n.Receive(context.Background(), keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("After remove: Receive error = %w", err)
			}

			keys = make([]*rpcapi.DataItem, len(tt.wantAfterRemove.KVs))
			for i := range tt.wantAfterRemove.KVs {
				keys[i] = tt.wantAfterRemove.KVs[i].Key
			}

			kvs, err = n.Receive(context.Background(), keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("Before remove: Receive error = %w", err)
			}

			assert.Equal(t, tt.wantAfterRemove, kvs)
		})
	}
}

func TestNewLocalNode_DBErr(t *testing.T) {
	lnn := &LocalNode{db: &bolt.DB{}}
	_, err := lnn.Store(context.Background(), &rpcapi.KeyValue{Key: &rpcapi.DataItem{ID: []byte("test-key")}})
	assert.Error(t, err)

	_, err = lnn.Receive(context.Background(), []*rpcapi.DataItem{{ID: []byte("test-key")}})
	assert.Error(t, err)

	_, err = lnn.Remove(context.Background(), []*rpcapi.DataItem{{ID: []byte("test-key")}})
	assert.Error(t, err)

	_, err = lnn.Explore(context.Background())
	assert.Error(t, err)
}

func TestNewLocalNode_KVR(t *testing.T) {
	type want struct {
		err bool
		n   *LocalNode
	}
	type args struct {
		hashsum uint64
		hasherr error
		conf    *Config
		metrics chan<- *write.Point
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test",
			args: args{
				hasherr: nil,
				hashsum: 42,
				conf: &Config{
					Name:          stringprt("test-node"),
					Address:       "main-address",
					RPCAddress:    "rpc-address",
					InfluxAddress: "influx-address",
					Geo: &rpcapi.GeoData{
						Longitude: 1,
						Latitude:  2,
					},
					Mode:     KvrouterMode,
					Power:    0,
					Capacity: 0,
					DBpath:   "",
					Health: HealthConfig{
						CheckInterval:                  "",
						CheckTimeout:                   "",
						DeregisterCriticalServiceAfter: "",
					},
					KVRouter: &KVRouterConfig{
						Address: "",
					},
					Balancer: &config.BalancerConfig{
						Mode: 0,
						SFC: &config.SFCConfig{
							Dimensions: 0,
							Size:       0,
							Curve: config.CurveType{
								CurveType: 0,
							},
						},
						NodeHash: 0,
						DataMode: 0,
						Latency: config.Duration{
							Duration: 0,
						},
						State: false,
					},
					Latency: Duration{
						Duration: 0,
					},
				},
				metrics: nil,
			},
			want: want{
				err: false,
				n: &LocalNode{
					id: "test-node",
					geo: &rpcapi.GeoData{
						Longitude: 1,
						Latitude:  2,
					},
					lwID:       int64prt(0),
					address:    "main-address",
					rpcaddress: "rpc-address",
					h:          uint64(42),
				},
			},
		},
		{
			name: "hash err",
			args: args{
				hasherr: errors.New("hash test err"),
				conf: &Config{
					Name:          stringprt("test-node"),
					Address:       "main-address",
					RPCAddress:    "rpc-address",
					InfluxAddress: "influx-address",
					Geo: &rpcapi.GeoData{
						Longitude: 1,
						Latitude:  2,
					},
					Mode:     KvrouterMode,
					Power:    0,
					Capacity: 0,
					DBpath:   "",
					Health: HealthConfig{
						CheckInterval:                  "",
						CheckTimeout:                   "",
						DeregisterCriticalServiceAfter: "",
					},
					KVRouter: &KVRouterConfig{
						Address: "",
					},
					Balancer: &config.BalancerConfig{
						Mode: 0,
						SFC: &config.SFCConfig{
							Dimensions: 0,
							Size:       0,
							Curve: config.CurveType{
								CurveType: 0,
							},
						},
						NodeHash: 0,
						DataMode: 0,
						Latency: config.Duration{
							Duration: 0,
						},
						State: false,
					},
					Latency: Duration{
						Duration: 0,
					},
				},
				metrics: nil,
			},
			want: want{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbpath := tempfile()
			defer func(dbpath string) {
				if err := os.Remove(dbpath); err != nil {
					t.Fatal(err)
				}
			}(dbpath)
			db, err := bolt.Open(dbpath, 0666, nil)
			if err != nil {
				t.Fatal(err)
			}

			h := &mocks.Hasher{}
			h.On("Sum", mock.Anything).Return(tt.args.hashsum, tt.args.hasherr)
			kvr, err := router.NewRouter(nil, h, nil, nil, nil)

			if tt.want.n != nil {
				tt.want.n.kvr = kvr
				tt.want.n.db = db
			}
			got, err := NewLocalNode(tt.args.conf, nil, db, kvr, tt.args.metrics)
			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.n, got)
			}

			_, err = NewLocalNode(tt.args.conf, nil, &bolt.DB{}, kvr, tt.args.metrics)
			assert.Error(t, err)
		})
	}
}

func TestNewLocalNode_Consul(t *testing.T) {
	type want struct {
		err bool
		n   *LocalNode
	}
	type args struct {
		hashsum uint64
		hasherr error
		conf    *Config
		metrics chan<- *write.Point
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "test",
			args: args{
				hasherr: nil,
				hashsum: 42,
				conf: &Config{
					Name:          stringprt("test-node"),
					Address:       "main-address",
					RPCAddress:    "rpc-address",
					InfluxAddress: "influx-address",
					Geo: &rpcapi.GeoData{
						Longitude: 1,
						Latitude:  2,
					},
					Mode:     ConsulMode,
					Power:    0,
					Capacity: 0,
					DBpath:   "",
					Health: HealthConfig{
						CheckInterval:                  "",
						CheckTimeout:                   "",
						DeregisterCriticalServiceAfter: "",
					},
					Consul: &ConfigConsul{
						Config:   consulapi.Config{},
						Service:  "test-service",
						KVFolder: "test-folder",
					},
					Balancer: &config.BalancerConfig{
						Mode: 0,
						SFC: &config.SFCConfig{
							Dimensions: 0,
							Size:       0,
							Curve: config.CurveType{
								CurveType: 0,
							},
						},
						NodeHash: 0,
						DataMode: 0,
						Latency: config.Duration{
							Duration: 0,
						},
						State: false,
					},
					Latency: Duration{
						Duration: 0,
					},
				},
				metrics: nil,
			},
			want: want{
				err: false,
				n: &LocalNode{
					id: "test-node",
					geo: &rpcapi.GeoData{
						Longitude: 1,
						Latitude:  2,
					},
					lwID:       int64prt(0),
					address:    "main-address",
					rpcaddress: "rpc-address",
					h:          uint64(42),
				},
			},
		},
		{
			name: "hash err",
			args: args{
				hasherr: errors.New("hash test err"),
				conf: &Config{
					Name:          stringprt("test-node"),
					Address:       "main-address",
					RPCAddress:    "rpc-address",
					InfluxAddress: "influx-address",
					Geo: &rpcapi.GeoData{
						Longitude: 1,
						Latitude:  2,
					},
					Mode:     ConsulMode,
					Power:    0,
					Capacity: 0,
					DBpath:   "",
					Health: HealthConfig{
						CheckInterval:                  "",
						CheckTimeout:                   "",
						DeregisterCriticalServiceAfter: "",
					},
					KVRouter: &KVRouterConfig{
						Address: "",
					},
					Balancer: &config.BalancerConfig{
						Mode: 0,
						SFC: &config.SFCConfig{
							Dimensions: 0,
							Size:       0,
							Curve: config.CurveType{
								CurveType: 0,
							},
						},
						NodeHash: 0,
						DataMode: 0,
						Latency: config.Duration{
							Duration: 0,
						},
						State: false,
					},
					Latency: Duration{
						Duration: 0,
					},
				},
				metrics: nil,
			},
			want: want{
				err: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbpath := tempfile()
			defer func(dbpath string) {
				if err := os.Remove(dbpath); err != nil {
					t.Fatal(err)
				}
			}(dbpath)
			db, err := bolt.Open(dbpath, 0666, nil)
			if err != nil {
				t.Fatal(err)
			}

			h := &mocks.Hasher{}
			h.On("Sum", mock.Anything).Return(tt.args.hashsum, tt.args.hasherr)
			kvr, err := router.NewRouter(nil, h, nil, nil, nil)

			if tt.want.n != nil {
				tt.want.n.kvr = kvr
				tt.want.n.db = db
				tt.want.n.consul, _ = consulapi.NewClient(&tt.args.conf.Consul.Config)
			}

			got, err := NewLocalNode(tt.args.conf, nil, db, kvr, tt.args.metrics)
			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.n, got)
			}

			_, err = NewLocalNode(tt.args.conf, nil, &bolt.DB{}, kvr, tt.args.metrics)
			assert.Error(t, err)
		})
	}
}

func TestLocalNode_Hash(t *testing.T) {
	inn := &LocalNode{
		h: 42,
	}
	got := inn.Hash()
	assert.Equal(t, 42, int(got))
}

func TestLocalNode_SetHash(t *testing.T) {
	inn := &LocalNode{
		h: 0,
	}
	inn.SetHash(42)
	assert.Equal(t, 42, int(inn.h))
}

func TestLocalNode_Move(t *testing.T) {
	dbpath := tempfile()
	defer func(dbpath string) {
		if err := os.Remove(dbpath); err != nil {
			t.Fatal(err)
		}
	}(dbpath)
	db, err := bolt.Open(dbpath, 0666, nil)
	if err != nil {
		t.Fatal(err)
	}

	en := &mocks.Node{}
	en.On("StorePairs", mock.Anything, mock.Anything).Return(nil, nil)
	en.On("ID", mock.Anything).Return("test-en")

	en2 := &mocks.Node{}
	en2.On("StorePairs", mock.Anything, mock.Anything).Return(nil, nil)
	en2.On("ID", mock.Anything).Return("test-en")

	enErr := &mocks.Node{}
	enErr.On("StorePairs", mock.Anything, mock.Anything).Return(nil, errors.New("test err"))
	enErr.On("ID", mock.Anything).Return("test-en")

	nk := map[nodes.Node][]*rpcapi.DataItem{
		en:    {{ID: []byte("test-key")}},
		en2:   {},
		enErr: {{ID: []byte("test-key")}},
	}

	inn := &LocalNode{db: db}

	err = inn.Move(context.Background(), nk)
	assert.NoError(t, err)

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(mainBucket)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}

	err = inn.Move(context.Background(), nk)
	assert.NoError(t, err)

	inn = &LocalNode{db: &bolt.DB{}}
	err = inn.Move(context.Background(), nk)
	assert.NoError(t, err)
}

func TestLocalNode_StorePairs(t *testing.T) {
	dbpath := tempfile()
	defer func(dbpath string) {
		if err := os.Remove(dbpath); err != nil {
			t.Fatal(err)
		}
	}(dbpath)
	db, err := bolt.Open(dbpath, 0666, nil)
	if err != nil {
		t.Fatal(err)
	}

	pairs := []*rpcapi.KeyValue{
		{
			Key:   &rpcapi.DataItem{ID: []byte("test-key")},
			Value: []byte("test-val"),
			Found: true,
		},
	}

	inn := &LocalNode{db: db}

	_, err = inn.StorePairs(context.Background(), pairs)
	assert.NoError(t, err)

	inn = &LocalNode{db: &bolt.DB{}}

	_, err = inn.StorePairs(context.Background(), pairs)
	assert.Error(t, err)
}

func tempfile() string {
	f, err := ioutil.TempFile("testdata", "bolt-")
	if err != nil {
		panic(err)
	}
	if err := f.Close(); err != nil {
		panic(err)
	}
	if err := os.Remove(f.Name()); err != nil {
		panic(err)
	}
	return f.Name()
}

func stringprt(s string) *string {
	return &s
}

func int64prt(i int64) *int64 {
	return &i
}

func IgnoreRegistry(sliceString []string) []string {
	countString := make(map[string]bool)
	var uniq []string
	for _, v := range sliceString {
		if _, value := countString[strings.ToUpper(v)]; !value {
			countString[strings.ToUpper(v)] = true
			uniq = append(uniq, v)
		}

	}
	return uniq
}

func IgnoreRegistryAlt(sliceString []string) []string {
	countString := make(map[string]struct{})
	for i := range sliceString {
		countString[strings.ToUpper(sliceString[i])] = struct{}{}
	}
	uniq := make([]string, len(countString))
	i := 0
	for key := range countString {
		uniq[i] = key
		i++
	}
	return uniq
}

func IgnoreRegistryFold(sliceString []string) []string {
	var uniq []string
	for i := range sliceString {
		f := true
		for j := range uniq {
			if strings.EqualFold(sliceString[i], uniq[j]) {
				f = false
				break
			}
		}
		if f {
			uniq = append(uniq, sliceString[i])
		}
	}
	return uniq
}

func genKeys() (base []string) {
	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			base = append(base, fmt.Sprintf("key-%d", j))
		}
	}
	return
}

var base = genKeys()

func BenchmarkLocalNode_IgnoreRegistry(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Println(IgnoreRegistry(base))
	}
}

func BenchmarkLocalNode_IgnoreRegistryAlt(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Println(IgnoreRegistryAlt(base))
	}
}

func BenchmarkLocalNode_IgnoreRegistryFold(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		log.Println(IgnoreRegistryFold(base))
	}
}
