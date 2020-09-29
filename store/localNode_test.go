package store

import (
	"errors"
	"github.com/influxdata/influxdb-client-go/api/write"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/struckoff/kvstore/mocks"
	"github.com/struckoff/kvstore/router"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/nodes"
	"github.com/struckoff/kvstore/router/rpcapi"
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
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
			got := n.Meta()
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
			for i := range tt.wantBeforeRemove.KVs {
				keys[i] = tt.wantBeforeRemove.KVs[i].Key
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
			for i := range tt.wantAfterRemove.KVs {
				keys[i] = tt.wantAfterRemove.KVs[i].Key
			}

			kvs, err = n.Receive(keys)
			if (err != nil) != tt.wantErr {
				t.Errorf("Before remove: Receive error = %w", err)
			}

			assert.Equal(t, tt.wantAfterRemove, kvs)
		})
	}
}

func TestNewLocalNode_DBErr(t *testing.T) {
	lnn := &LocalNode{db: &bolt.DB{}}
	err := lnn.Store("test-key", nil)
	assert.Error(t, err)

	_, err = lnn.Receive([]string{"test-key"})
	assert.Error(t, err)

	err = lnn.Remove([]string{"test-key"})
	assert.Error(t, err)

	_, err = lnn.Explore()
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
					Mode:     0,
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
					Mode:     0,
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
			defer os.Remove(dbpath)
			db, err := bolt.Open(dbpath, 0666, nil)
			if err != nil {
				t.Fatal(err)
			}

			h := &mocks.Hasher{}
			h.On("Sum", mock.Anything).Return(tt.args.hashsum, tt.args.hasherr)
			kvr, err := router.NewRouter(nil, h, nil, nil)

			if tt.want.n != nil {
				tt.want.n.kvr = kvr
				tt.want.n.db = db
			}
			got, err := NewLocalNode(tt.args.conf, db, kvr, tt.args.metrics)
			if tt.want.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.n, got)
			}

			_, err = NewLocalNode(tt.args.conf, &bolt.DB{}, kvr, tt.args.metrics)
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
	defer os.Remove(dbpath)
	db, err := bolt.Open(dbpath, 0666, nil)
	if err != nil {
		t.Fatal(err)
	}

	en := &mocks.Node{}
	en.On("StorePairs", mock.Anything).Return(nil)
	en.On("ID", mock.Anything).Return("test-en")

	en2 := &mocks.Node{}
	en2.On("StorePairs", mock.Anything).Return(nil)
	en2.On("ID", mock.Anything).Return("test-en")

	enErr := &mocks.Node{}
	enErr.On("StorePairs", mock.Anything).Return(errors.New("test err"))
	enErr.On("ID", mock.Anything).Return("test-en")

	nk := map[nodes.Node][]string{
		en:    {"test-key"},
		en2:   {},
		enErr: {"test-key"},
	}

	inn := &LocalNode{db: db}

	err = inn.Move(nk)
	assert.NoError(t, err)

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(mainBucket)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}

	err = inn.Move(nk)
	assert.NoError(t, err)

	inn = &LocalNode{db: &bolt.DB{}}
	err = inn.Move(nk)
	assert.NoError(t, err)
}

func TestLocalNode_StorePairs(t *testing.T) {
	dbpath := tempfile()
	defer os.Remove(dbpath)
	db, err := bolt.Open(dbpath, 0666, nil)
	if err != nil {
		t.Fatal(err)
	}

	pairs := []*rpcapi.KeyValue{
		{
			Key:   "test-key",
			Value: "test-val",
			Found: true,
		},
	}

	inn := &LocalNode{db: db}

	err = inn.StorePairs(pairs)
	assert.NoError(t, err)

	inn = &LocalNode{db: &bolt.DB{}}

	err = inn.StorePairs(pairs)
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
