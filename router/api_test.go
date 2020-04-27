package router

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/struckoff/kvstore/router/balanceradapter"
	"github.com/struckoff/kvstore/router/nodes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_HTTPHandler(t *testing.T) {
	type args struct {
		method string
		path   string
		body   io.Reader
	}
	type fields struct {
		bal balanceradapter.Balancer
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
				bal: mockBalancer{nodes: []nodes.Node{
					nodes.mockNode{id: "mnode-0"},
					nodes.mockNode{id: "mnode-1"},
					nodes.mockNode{id: "mnode-2"},
				}},
			},
			want: &httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBuffer([]byte("[{\"ID\":\"mnode-0\"},{\"ID\":\"mnode-1\"},{\"ID\":\"mnode-2\"}]\n")),
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
				bal: mockBalancer{
					nodes: []nodes.Node{
						nodes.mockNode{
							id: "mnode-0",
							kv: map[string][]byte{
								"mnode-0-key-0": nil,
								"mnode-0-key-1": nil,
								"mnode-0-key-2": nil,
							}},
						nodes.mockNode{
							id: "mnode-1",
							kv: map[string][]byte{
								"mnode-1-key-0": nil,
								"mnode-1-key-1": nil,
								"mnode-1-key-2": nil,
							}},
					},
					kv: nil,
				},
			},
			want: &httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBuffer([]byte("{\"mnode-0\":[\"mnode-0-key-0\",\"mnode-0-key-1\",\"mnode-0-key-2\"],\"mnode-1\":[\"mnode-1-key-0\",\"mnode-1-key-1\",\"mnode-1-key-2\"]}")),
			},
		},
		{
			name: "GET /get",
			args: args{
				method: "GET",
				path:   "/get/mkey-1",
				body:   nil,
			},
			fields: fields{
				bal: mockBalancer{
					kv: map[string][]byte{
						"mkey-0": []byte("mkey-0-val"),
						"mkey-1": []byte("mkey-1-val"),
						"mkey-2": []byte("mkey-2-val"),
					},
				},
			},
			want: &httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBuffer([]byte("mkey-1-val")),
			},
		},
		{
			name: "POST /put",
			args: args{
				method: "POST",
				path:   "/put/mkey-3",
				body:   bytes.NewBuffer([]byte("mkey-3-val")),
			},
			fields: fields{
				bal: mockBalancer{
					kv: map[string][]byte{
						"mkey-0": []byte("mkey-0-val"),
						"mkey-1": []byte("mkey-1-val"),
						"mkey-2": []byte("mkey-2-val"),
					},
				},
			},
			want: &httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBuffer([]byte("OK")),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Router{
				bal: tt.fields.bal,
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

type mockBalancer struct {
	nodes []nodes.Node
	kv    map[string][]byte
}

func (bal mockBalancer) AddNode(n nodes.Node) error {
	bal.nodes = append(bal.nodes, n)
	return nil
}

func (m mockBalancer) RemoveNode(id string) error {
	panic("implement me")
}

func (m mockBalancer) SetNodes(ns []nodes.Node) error {
	m.nodes = ns
	return nil
}

func (m mockBalancer) LocateData(key string) (nodes.Node, error) {
	return nodes.mockNode{id: key, kv: m.kv}, nil
}

func (m mockBalancer) Nodes() ([]nodes.Node, error) {
	return m.nodes, nil
}

func (m mockBalancer) GetNode(id string) (nodes.Node, error) {
	for _, n := range m.nodes {
		if n.ID() == id {
			return n, nil
		}
	}
	return nil, errors.New("Node not found")
}
