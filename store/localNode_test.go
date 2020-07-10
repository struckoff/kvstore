package store

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/stretchr/testify/assert"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/nodes"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"testing"

	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/kvstore/router/rpcapi"
	bolt "go.etcd.io/bbolt"
)

func TestNewInternalNode(t *testing.T) {
	type args struct {
		id    string
		addr  string
		raddr string
		p     float64
		c     float64
		db    *bolt.DB
		bConf config.BalancerConfig
	}
	tests := []struct {
		name string
		args args
		want *LocalNode
	}{
		{
			name: "test",
			args: args{
				id:    "test_id",
				addr:  "test_addr",
				raddr: "test_raddr",
				p:     1.1,
				c:     2.3,
				db:    &bolt.DB{},
				bConf: config.BalancerConfig{
					Latency: config.Duration{
						Duration: 0,
					},
				},
			},
			want: &LocalNode{
				id:         "test_id",
				address:    "test_addr",
				rpcaddress: "test_raddr",
				p:          nodes.NewPower(1.1),
				c:          nodes.NewCapacity(2.3),
				db:         &bolt.DB{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := &Config{
				Name:       &tt.args.id,
				Address:    tt.args.addr,
				RPCAddress: tt.args.raddr,
				Power:      tt.args.p,
				Capacity:   tt.args.c,
				Balancer:   &tt.args.bConf,
			}
			metrics := make(chan *influxdb2.Point)
			defer close(metrics)
			go func(metrics chan *influxdb2.Point) {
				for m := range metrics {
					if m == nil {
						fmt.Println("closed")
						return
					}
					continue
				}
			}(metrics)
			got, _ := NewLocalNode(conf, tt.args.db, nil, metrics)
			assert.Equal(t, tt.want.id, got.id)
			assert.Equal(t, tt.want.rpcserver, got.rpcserver)
			assert.Equal(t, tt.want.consul, got.consul)
			assert.Equal(t, tt.want.kvr, got.kvr)
			assert.Equal(t, tt.want.db, got.db)
			//if !reflect.DeepEqual(, ) {
			//	t.Errorf("NewLocalNode() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestInternalNode_Meta(t *testing.T) {
	type fields struct {
		id         string
		address    string
		rpcaddress string
		p          nodes.Power
		c          nodes.Capacity
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
				c:          nodes.NewCapacity(2.3),
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
			got := n.Meta()
			assert.Equal(t, tt.want, got)
			//if got := n.Meta(); !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("Meta() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestInternalNode_Capacity(t *testing.T) {
	type fields struct {
		id         string
		address    string
		rpcaddress string
		p          nodes.Power
		c          nodes.Capacity
		db         *bolt.DB
	}
	tests := []struct {
		name   string
		fields fields
		want   balancer.Capacity
	}{
		{
			name: "test",
			fields: fields{
				id:         "test_id",
				address:    "test_addr",
				rpcaddress: "test_raddr",
				p:          nodes.NewPower(1.1),
				c:          nodes.NewCapacity(2.3),
				db:         nil,
			},
			want: nodes.NewCapacity(2.3),
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
			if got := n.Capacity(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Capacity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInternalNode_HTTPAddress(t *testing.T) {
	type fields struct {
		id         string
		address    string
		rpcaddress string
		p          nodes.Power
		c          nodes.Capacity
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
				c:          nodes.NewCapacity(2.3),
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
			if got := n.HTTPAddress(); got != tt.want {
				t.Errorf("HTTPAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInternalNode_ID(t *testing.T) {
	type fields struct {
		id         string
		address    string
		rpcaddress string
		p          nodes.Power
		c          nodes.Capacity
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
				c:          nodes.NewCapacity(2.3),
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
		c          nodes.Capacity
		db         *bolt.DB
	}
	tests := []struct {
		name   string
		fields fields
		want   balancer.Power
	}{
		{
			name: "test",
			fields: fields{
				id:         "test_id",
				address:    "test_addr",
				rpcaddress: "test_raddr",
				p:          nodes.NewPower(1.1),
				c:          nodes.NewCapacity(2.3),
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
		c          nodes.Capacity
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
				c:          nodes.NewCapacity(2.3),
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
	type kv struct {
		key string
		val []byte
	}
	type args struct {
		kvs        []kv
		removeKeys []string
	}
	tests := []struct {
		name                    string
		args                    args
		wantErr                 bool
		wantBeforeRemove        *rpcapi.KeyValues
		wantExploreBeforeRemove []string
		wantAfterRemove         *rpcapi.KeyValues
		wantExploreAfterRemove  []string
	}{
		{
			name: "test",
			args: args{
				kvs: []kv{
					{"t0", []byte("t0val")},
					{"t1", []byte("t1val")},
					{"t2", []byte("t2val")},
					{"t3", []byte("t3val")},
					{"t4", []byte("t4val")},
				},
				removeKeys: []string{"t0", "t2", "t3"},
			},
			wantErr: false,
			wantBeforeRemove: &rpcapi.KeyValues{
				KVs: []*rpcapi.KeyValue{
					{Key: "t0", Value: "t0val", Found: true},
					{Key: "t1", Value: "t1val", Found: true},
					{Key: "t2", Value: "t2val", Found: true},
					{Key: "t3", Value: "t3val", Found: true},
					{Key: "t4", Value: "t4val", Found: true},
				},
			},
			wantExploreBeforeRemove: []string{"t0", "t1", "t2", "t3", "t4"},
			wantAfterRemove: &rpcapi.KeyValues{
				KVs: []*rpcapi.KeyValue{
					{Key: "t1", Value: "t1val", Found: true},
					{Key: "t4", Value: "t4val", Found: true},
				},
			},
			wantExploreAfterRemove: []string{"t1", "t4"},
		},
		{
			name: "remove not existing key",
			args: args{
				kvs: []kv{
					{"t0", []byte("t0val")},
					{"t1", []byte("t1val")},
					{"t2", []byte("t2val")},
					{"t3", []byte("t3val")},
					{"t4", []byte("t4val")},
				},
				removeKeys: []string{"k0", "k2", "k3"},
			},
			wantErr: false,
			wantBeforeRemove: &rpcapi.KeyValues{
				KVs: []*rpcapi.KeyValue{
					{Key: "t0", Value: "t0val", Found: true},
					{Key: "t1", Value: "t1val", Found: true},
					{Key: "t2", Value: "t2val", Found: true},
					{Key: "t3", Value: "t3val", Found: true},
					{Key: "t4", Value: "t4val", Found: true},
				},
			},
			wantExploreBeforeRemove: []string{"t0", "t1", "t2", "t3", "t4"},
			wantAfterRemove: &rpcapi.KeyValues{
				KVs: []*rpcapi.KeyValue{
					{Key: "t0", Value: "t0val", Found: true},
					{Key: "t1", Value: "t1val", Found: true},
					{Key: "t2", Value: "t2val", Found: true},
					{Key: "t3", Value: "t3val", Found: true},
					{Key: "t4", Value: "t4val", Found: true},
				},
			},
			wantExploreAfterRemove: []string{"t0", "t1", "t2", "t3", "t4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbpath := tempfile()
			defer os.Remove(dbpath)
			db, err := bolt.Open(dbpath, 0666, nil)
			if err != nil {
				panic(err)
			}
			metrics := make(chan *influxdb2.Point)
			defer close(metrics)
			go func(metrics chan *influxdb2.Point) {
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
			}
			for _, kv := range tt.args.kvs {
				if err := n.Store(kv.key, kv.val); (err != nil) != tt.wantErr {
					t.Errorf("Store() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			explore, err := n.Explore()
			if (err != nil) != tt.wantErr {
				t.Errorf("Before remove: Explore() error = %v, wantErr %v", err, tt.wantErr)
			}

			sort.Strings(explore)
			sort.Strings(tt.wantExploreBeforeRemove)
			if !reflect.DeepEqual(explore, tt.wantExploreBeforeRemove) {
				t.Errorf("Before remove: Explore() = %v, want %v", explore, tt.wantExploreBeforeRemove)
			}

			keys := make([]string, len(tt.wantBeforeRemove.KVs))
			for iter := range tt.wantBeforeRemove.KVs {
				keys[iter] = tt.wantBeforeRemove.KVs[iter].Key
			}

			kvs, err := n.Receive(keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("Before remove: Receive error = %w", err)
			}

			assert.Equal(t, tt.wantBeforeRemove, kvs)

			for _, key := range tt.args.removeKeys {
				if err := n.Remove([]string{key}); (err != nil) != tt.wantErr {
					t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			explore, err = n.Explore()
			if (err != nil) != tt.wantErr {
				t.Errorf("After remove: Explore() error = %v, wantErr %v", err, tt.wantErr)
			}

			sort.Strings(explore)
			sort.Strings(tt.wantExploreAfterRemove)
			assert.Equal(t, tt.wantExploreAfterRemove, explore)

			kvs, err = n.Receive(keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("After remove: Receive error = %w", err)
			}

			keys = make([]string, len(tt.wantAfterRemove.KVs))
			for iter := range tt.wantAfterRemove.KVs {
				keys[iter] = tt.wantAfterRemove.KVs[iter].Key
			}

			kvs, err = n.Receive(keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("Before remove: Receive error = %w", err)
			}

			assert.Equal(t, tt.wantAfterRemove, kvs)
		})
	}
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
