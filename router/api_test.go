package router

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	balancer "github.com/struckoff/SFCFramework"
	balancermocs "github.com/struckoff/SFCFramework/mocks"
	balanceradaptermock "github.com/struckoff/kvstore/router/balanceradapter/mocks"
	"github.com/struckoff/kvstore/router/nodes"
	nodesmock "github.com/struckoff/kvstore/router/nodes/mocks"
	"github.com/struckoff/kvstore/router/rpcapi"
	"sort"
	"strings"

	//"github.com/struckoff/kvstore/router/nodes/mocks"
	"io"
	"net/http"
	"net/http/httptest"
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
				Body: bytes.NewBuffer([]byte("test-node-1")),
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
				mn := &nodesmock.Node{}
				mn.On("ID").Return(name)
				mn.On("Meta").Return(&rpcapi.NodeMeta{ID: name})
				mn.On("Explore").Return(keys, nil)
				mn.On("Store", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
				ns = append(ns, mn)
			}
			sort.Slice(ns, func(i, j int) bool { return strings.Compare(ns[i].ID(), ns[j].ID()) < 1 })

			mn := &nodesmock.Node{}
			mn.On("Store", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
			mn.On("Receive", "test-node-1-key-1").Return([]byte("test-node-1"), nil)

			mbal := &balanceradaptermock.Balancer{}
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
