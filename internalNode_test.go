package kvstore

import (
	balancer "github.com/struckoff/SFCFramework"
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestNewInternalNode(t *testing.T) {
	type args struct {
		id    string
		addr  string
		raddr string
		p     float64
		c     float64
		db    *bolt.DB
	}
	tests := []struct {
		name string
		args args
		want *InternalNode
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
			},
			want: &InternalNode{
				id:         "test_id",
				address:    "test_addr",
				rpcaddress: "test_raddr",
				p: Power{
					p: 1.1,
				},
				c: Capacity{
					c: 2.3,
				},
				db: &bolt.DB{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewInternalNode(tt.args.id, tt.args.addr, tt.args.raddr, tt.args.p, tt.args.c, tt.args.db)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewInternalNode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInternalNode_Meta(t *testing.T) {
	type fields struct {
		id         string
		address    string
		rpcaddress string
		p          Power
		c          Capacity
		db         *bolt.DB
	}
	tests := []struct {
		name   string
		fields fields
		want   NodeMeta
	}{
		{
			name: "test",
			fields: fields{
				id:         "test_id",
				address:    "test_addr",
				rpcaddress: "test_raddr",
				p:          Power{1.1},
				c:          Capacity{2.3},
				db:         nil,
			},
			want: NodeMeta{
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
			n := &InternalNode{
				id:         tt.fields.id,
				address:    tt.fields.address,
				rpcaddress: tt.fields.rpcaddress,
				p:          tt.fields.p,
				c:          tt.fields.c,
				db:         tt.fields.db,
			}
			if got := n.Meta(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Meta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInternalNode_Capacity(t *testing.T) {
	type fields struct {
		id         string
		address    string
		rpcaddress string
		p          Power
		c          Capacity
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
				p:          Power{1.1},
				c:          Capacity{2.3},
				db:         nil,
			},
			want: Capacity{2.3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &InternalNode{
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
		p          Power
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
				p:          Power{1.1},
				c:          Capacity{2.3},
				db:         nil,
			},
			want: "test_addr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &InternalNode{
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
		p          Power
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
				p:          Power{1.1},
				c:          Capacity{2.3},
				db:         nil,
			},
			want: "test_id",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &InternalNode{
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
		p          Power
		c          Capacity
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
				p:          Power{1.1},
				c:          Capacity{2.3},
				db:         nil,
			},
			want: Power{1.1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &InternalNode{
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
		p          Power
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
				p:          Power{1.1},
				c:          Capacity{2.3},
				db:         nil,
			},
			want: "test_raddr",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &InternalNode{
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
		wantBeforeRemove        []kv
		wantExploreBeforeRemove []string
		wantAfterRemove         []kv
		wantExploreAfterRemove  []string
	}{
		{
			name: "test",
			args: args{
				kvs: []kv{
					{"t0", []byte("t0val")},
					{"t1", []byte("t0val")},
					{"t2", []byte("t0val")},
					{"t3", []byte("t0val")},
					{"t4", []byte("t0val")},
				},
				removeKeys: []string{"t0", "t2", "t3"},
			},
			wantErr: false,
			wantBeforeRemove: []kv{
				{"t0", []byte("t0val")},
				{"t1", []byte("t0val")},
				{"t2", []byte("t0val")},
				{"t3", []byte("t0val")},
				{"t4", []byte("t0val")},
			},
			wantExploreBeforeRemove: []string{"t0", "t1", "t2", "t3", "t4"},
			wantAfterRemove: []kv{
				{"t1", []byte("t0val")},
				{"t4", []byte("t0val")},
				{"t0", nil},
				{"t2", nil},
				{"t3", nil},
			},
			wantExploreAfterRemove: []string{"t1", "t4"},
		},
		{
			name: "remove not existing key",
			args: args{
				kvs: []kv{
					{"t0", []byte("t0val")},
					{"t1", []byte("t0val")},
					{"t2", []byte("t0val")},
					{"t3", []byte("t0val")},
					{"t4", []byte("t0val")},
				},
				removeKeys: []string{"k0", "k2", "k3"},
			},
			wantErr: false,
			wantBeforeRemove: []kv{
				{"t0", []byte("t0val")},
				{"t1", []byte("t0val")},
				{"t2", []byte("t0val")},
				{"t3", []byte("t0val")},
				{"t4", []byte("t0val")},
			},
			wantExploreBeforeRemove: []string{"t0", "t1", "t2", "t3", "t4"},
			wantAfterRemove: []kv{
				{"t0", []byte("t0val")},
				{"t1", []byte("t0val")},
				{"t2", []byte("t0val")},
				{"t3", []byte("t0val")},
				{"t4", []byte("t0val")},
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
			n := &InternalNode{
				db: db,
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

			for _, kv := range tt.wantBeforeRemove {
				val, err := n.Receive(kv.key)
				if (err != nil) != tt.wantErr {
					t.Errorf("Before remove: Receive(%s) = %v, want %v", kv.key, err, tt.wantErr)
				}
				if !reflect.DeepEqual(val, kv.val) {
					t.Errorf("Before remove: Receive(%s) = %v, want %v", kv.key, val, kv.val)
				}
			}

			for _, key := range tt.args.removeKeys {
				if err := n.Remove(key); (err != nil) != tt.wantErr {
					t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
				}
			}

			explore, err = n.Explore()
			if (err != nil) != tt.wantErr {
				t.Errorf("After remove: Explore() error = %v, wantErr %v", err, tt.wantErr)
			}

			sort.Strings(explore)
			sort.Strings(tt.wantExploreAfterRemove)
			if !reflect.DeepEqual(explore, tt.wantExploreAfterRemove) {
				t.Errorf("After remove: Explore() = %v, want %v", explore, tt.wantExploreAfterRemove)
			}

			for _, kv := range tt.wantAfterRemove {
				val, err := n.Receive(kv.key)
				if (err != nil) != tt.wantErr {
					t.Errorf("After remove: Receive(%s) = %v, want %v", kv.key, err, tt.wantErr)
				}
				if !reflect.DeepEqual(val, kv.val) {
					t.Errorf("After remove: Receive(%s) = %v, want %v", kv.key, val, kv.val)
				}
			}
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
