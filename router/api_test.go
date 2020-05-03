package router

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	balancer "github.com/struckoff/SFCFramework"
	balancermocs "github.com/struckoff/SFCFramework/mocks"
	"github.com/struckoff/kvstore/router/mocks"
	"github.com/struckoff/kvstore/router/nodes"
	"github.com/struckoff/kvstore/router/rpcapi"
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

func TestRouter_HTTPHandler_GET(t *testing.T) {
	type args struct {
		method string
		path   string
		body   io.Reader
	}
	type fields struct {
		//bal balanceradapter.Balancer
		nodes map[string][]string
	}
	tests := []struct {
		name   string
		args   args
		fields fields
		want   *httptest.ResponseRecorder
	}{
		{
			name: "GET /nodes",
			args: args{
				method: "GET",
				path:   "/nodes",
				body:   nil,
			},
			fields: fields{
				nodes: map[string][]string{
					"test-node-0": nil,
					"test-node-1": nil,
					"test-node-2": nil,
				},
			},
			want: &httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBuffer([]byte("[{\"ID\":\"test-node-0\"},{\"ID\":\"test-node-1\"},{\"ID\":\"test-node-2\"}]\n")),
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
				Body: bytes.NewBuffer([]byte("{\"KVs\":[{\"Key\":\"test-node-1-key-1\",\"Value\":\"test-node-1\"}]}\n")),
			},
		},
		{
			name: "POST /put/test-key-3",
			args: args{
				method: "POST",
				body:   bytes.NewBuffer([]byte("test-key-3-val")),
				path:   "/put/test-key-3",
			},
			want: &httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBuffer([]byte("OK")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ns []nodes.Node
			for name, keys := range tt.fields.nodes {
				mn := &mocks.Node{}
				mn.On("ID").Return(name)
				mn.On("Meta").Return(&rpcapi.NodeMeta{ID: name})
				mn.On("Explore").Return(keys, nil)
				mn.On("Store", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
				ns = append(ns, mn)
			}
			sort.Slice(ns, func(i, j int) bool { return strings.Compare(ns[i].ID(), ns[j].ID()) < 1 })

			mn := &mocks.Node{}
			mn.On("Store", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
			kvs := &rpcapi.KeyValues{
				KVs: []*rpcapi.KeyValue{{Key: "test-node-1-key-1", Value: "test-node-1"}},
			}
			mn.On("Receive", []string{"test-node-1-key-1"}).Return(kvs, nil)

			mbal := &mocks.Balancer{}
			mbal.On("LocateData", mock.AnythingOfType("*mocks.DataItem")).Return(mn, nil)
			mbal.On("Nodes").Return(ns, nil)
			h := &Router{
				bal: mbal,
				ndf: func(s string) (balancer.DataItem, error) {
					di := &balancermocs.DataItem{}
					di.On("ID").Return(s)
					return di, nil
				},
			}
			handler := h.HTTPHandler()

			req, err := http.NewRequest(tt.args.method, tt.args.path, tt.args.body)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.want.Code, rr.Code)
			assert.Equal(t, tt.want.Body, rr.Body)
		})
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func generateData(n int) [][]byte {
	res := make([][]byte, n)
	for iter := range res {
		b := make([]rune, 256)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		res[iter] = []byte(fmt.Sprintf("%d", rand.Int()))
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
