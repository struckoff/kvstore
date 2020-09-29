package router

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/struckoff/kvstore/mocks"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/nodes"
	"github.com/struckoff/kvstore/router/rpcapi"
	balancer "github.com/struckoff/sfcframework"
	balancermocks "github.com/struckoff/sfcframework/mocks"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"

	"testing"
)

func TestRouter_HTTPHandler_API(t *testing.T) {
	type args struct {
		method string
		path   string
		body   io.Reader
	}
	type fields struct {
		nodes map[string][]string
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   *httptest.ResponseRecorder
	}{
		{
			name: "GET /get",
			args: args{
				method: "GET",
				path:   "/get/test-node-1-key-1",
				body:   nil,
			},
			fields: fields{
				nodes: map[string][]string{
					"test-node-0": {"test-node-0-key-0", "test-node-0-key-1", "test-node-0-key-2"},
					"test-node-1": {"test-node-1-key-0", "test-node-1-key-1", "test-node-1-key-2"},
				},
			},
			want: &httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBuffer([]byte("[{\"Key\":\"test-node-1-key-1\",\"Value\":\"test-node-1\",\"Found\":true}]\n")),
			},
		},
		{
			name: "OPTIONS /optimize",
			args: args{
				method: "OPTIONS",
				path:   "/optimize",
				body:   nil,
			},
			fields: fields{
				nodes: map[string][]string{
					"test-node-0": {"test-node-0-key-0", "test-node-0-key-1", "test-node-0-key-2"},
					"test-node-1": {"test-node-1-key-0", "test-node-1-key-1", "test-node-1-key-2"},
				},
			},
			want: &httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBuffer([]byte("Optimize complete")),
			},
		},
		{
			name: "GET /config",
			args: args{
				method: "GET",
				path:   "/config",
				body:   nil,
			},
			fields: fields{
			},
			want: &httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBuffer([]byte("{\"Mode\":\"SFC\",\"SFC\":{\"Dimensions\":2,\"Size\":8,\"Curve\":\"Hilbert\"},\"NodeHash\":\"GeoSFC\",\"DataMode\":\"Geo\",\"Latency\":\"0s\",\"State\":false}\n")),
			},
		},
		{
			name: "POST /put",
			args: args{
				method: "POST",
				path:   "/put/test-key-1",
				body:   bytes.NewBuffer([]byte("test-body")),
			},
			fields: fields{
			},
			want: &httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBuffer([]byte("OK")),
			},
		},
		{
			name: "GET /list",
			args: args{
				method: "GET",
				path:   "/list",
				body:   nil,
			},
			fields: fields{
				nodes: map[string][]string{
					"test-node-0": {"test-node-0-key-0", "test-node-0-key-1", "test-node-0-key-2"},
					"test-node-1": {"test-node-1-key-0", "test-node-1-key-1", "test-node-1-key-2"},
				},
			},
			want: &httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBuffer([]byte("{\"test-node-0\":[\"test-node-0-key-0\",\"test-node-0-key-1\",\"test-node-0-key-2\"],\"test-node-1\":[\"test-node-1-key-0\",\"test-node-1-key-1\",\"test-node-1-key-2\"]}")),
			},
		},
		{
			name: "GET /cid",
			args: args{
				method: "GET",
				path:   "/cid",
				body:   nil,
			},
			fields: fields{
				nodes: map[string][]string{
					"test-node-1": {"test-node-1-key-0", "test-node-1-key-1", "test-node-1-key-2"},
				},
			},
			want: &httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBuffer([]byte("{\"1\":[\"test-node-1-key-0\",\"test-node-1-key-1\",\"test-node-1-key-2\"]}")),
			},
		},
		{
			name: "GET /nodes",
			args: args{
				method: "GET",
				path:   "/nodes",
				body:   nil,
			},
			fields: fields{
				nodes: map[string][]string{
					"test-node-0": {"test-node-0-key-0", "test-node-0-key-1", "test-node-0-key-2"},
					"test-node-1": {"test-node-1-key-0", "test-node-1-key-1", "test-node-1-key-2"},
				},
			},
			want: &httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBuffer([]byte("[{\"ID\":\"test-node-0\"},{\"ID\":\"test-node-1\"}]\n")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ns []nodes.Node
			for name, keys := range tt.fields.nodes {
				c := &balancermocks.Capacity{}
				c.On("Get").Return(42.42, nil)


				mn := &mocks.Node{}
				mn.On("ID").Return(name)
				mn.On("Meta").Return(&rpcapi.NodeMeta{ID: name})
				mn.On("Explore").Return(keys, nil)
				mn.On("Store", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
				mn.On("Move", mock.Anything).Return(nil)
				mn.On("Capacity", mock.Anything).Return(c)
				ns = append(ns, mn)
			}
			sort.Slice(ns, func(i, j int) bool { return strings.Compare(ns[i].ID(), ns[j].ID()) < 1 })

			mn := &mocks.Node{}
			mn.On("Store", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
			kvs := &rpcapi.KeyValues{
				KVs: []*rpcapi.KeyValue{{Key: "test-node-1-key-1", Value: "test-node-1", Found: true}},
			}
			mn.On("Receive", []string{"test-node-1-key-1"}).Return(kvs, nil)
			mn.On("ID").Return("test-node-1")
			mn.On("Move", mock.Anything).Return(nil)

			mbal := &mocks.Balancer{}
			mbal.On("LocateData", mock.Anything).Return(mn, uint64(1), nil)
			mbal.On("Nodes").Return(ns, nil)
			mbal.On("Reset").Return(nil)
			mbal.On("Optimize").Return(nil)
			mbal.On("AddData", mock.Anything).Return(mn, uint64(1), nil)
			mbal.On("SetNodes", mock.Anything).Return(nil)
			mbal.On("AddNode", mock.Anything).Return(nil)
			h := &Router{
				bal: mbal,
				ndf: func(s string) (balancer.DataItem, error) {
					di := &balancermocks.DataItem{}
					di.On("ID").Return(s)
					return di, nil
				},
				conf: &config.BalancerConfig{
					Mode:     config.SFCMode,
					SFC:      &config.SFCConfig{
						Dimensions: 2,
						Size:       8,
					},
					NodeHash: 1,
					DataMode: config.GeoData,
				},
			}
			handler := h.HTTPHandler()

			req, err := http.NewRequest(tt.args.method, tt.args.path, tt.args.body)
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.want.Code, rr.Code)
			assert.Equal(t, tt.want.Body.String(), rr.Body.String())
		})
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateData(n int) [][]byte {
	res := make([][]byte, n)
	for i := range res {
		b := make([]rune, 256)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		res[i] = []byte(fmt.Sprintf("%d", rand.Int()))
	}

	return res
}

func Benchmark_byteSlice2String(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	data := generateData(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		log.Print(byteSlice2String(data[i]))
	}
}

func Benchmark_byteSlice2String_normal(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	data := generateData(b.N)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		log.Print(string(data[i]))
	}
}
