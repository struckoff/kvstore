package balanceradapter

import (
	"github.com/lafikl/consistent"
	"github.com/stretchr/testify/assert"
	"github.com/struckoff/kvstore/router/nodes"
	"github.com/struckoff/kvstore/router/nodes/mocks"
	"sort"
	"sync"
	"testing"
)

func TestConsistent_AddNode(t *testing.T) {
	type fields struct {
		ring  *consistent.Consistent
		nodes sync.Map
	}
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantCalls map[string]int
		wantErr   bool
	}{
		{
			name: "test",
			fields: fields{
				ring:  consistent.New(),
				nodes: sync.Map{},
			},
			args: args{
				name: "test-node",
			},
			wantCalls: map[string]int{
				"ID": 2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//var n nodes.Node
			mn := mocks.Node{}
			mn.On("ID").Return(tt.args.name)
			c := &Consistent{
				ring:  tt.fields.ring,
				nodes: tt.fields.nodes,
			}
			if err := c.AddNode(&mn); (err != nil) != tt.wantErr {
				t.Errorf("AddNode() error = %v, wantErr %v", err, tt.wantErr)
			}
			got, _ := c.nodes.Load(tt.args.name)
			names := c.ring.Hosts()

			mn.AssertExpectations(t)
			assert.Equal(t, &mn, got)
			assert.Equal(t, []string{tt.args.name}, names)
			for method, count := range tt.wantCalls {
				mn.AssertNumberOfCalls(t, method, count)
			}
		})
	}
}

func TestConsistent_RemoveNode(t *testing.T) {
	type fields struct {
		ring  *consistent.Consistent
		nodes sync.Map
	}
	type args struct {
		names  []string
		remove string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantNodes []string
		wantErr   bool
	}{
		{
			name: "remove 1 of 1",
			fields: fields{
				ring:  consistent.New(),
				nodes: sync.Map{},
			},
			args: args{
				names: []string{
					"test-node-1",
				},
				remove: "test-node-1",
			},
			wantNodes: nil,
			wantErr:   false,
		},
		{
			name: "remove 1 of 3",
			fields: fields{
				ring:  consistent.New(),
				nodes: sync.Map{},
			},
			args: args{
				names: []string{
					"test-node-0",
					"test-node-1",
					"test-node-2",
				},
				remove: "test-node-1",
			},
			wantNodes: []string{
				"test-node-0",
				"test-node-2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//var n nodes.Node
			//Fill it ip
			for _, name := range tt.args.names {
				tt.fields.ring.Add(name)
				mn := &mocks.Node{}
				//mn.On("ID").Return(name)
				tt.fields.nodes.Store(name, mn)
			}

			c := &Consistent{
				ring:  tt.fields.ring,
				nodes: tt.fields.nodes,
			}
			if err := c.RemoveNode(tt.args.remove); (err != nil) != tt.wantErr {
				t.Errorf("RemoveNode() error = %v, wantErr %v", err, tt.wantErr)
			}

			_, ok := c.nodes.Load(tt.args.remove)
			assert.True(t, !ok)
			nodes := c.ring.Hosts()

			var mapNodes []string
			c.nodes.Range(func(key, value interface{}) bool {
				mapNodes = append(mapNodes, key.(string))
				return true
			})

			sort.Strings(tt.wantNodes)
			sort.Strings(nodes)
			sort.Strings(mapNodes)

			assert.Equal(t, tt.wantNodes, nodes)
			assert.Equal(t, tt.wantNodes, mapNodes)

		})
	}
}

func TestConsistent_SetNodes(t *testing.T) {
	type fields struct {
		ring  *consistent.Consistent
		nodes sync.Map
	}
	type args struct {
		wasNodes []string
		newNodes []string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantNodes []string
		wantErr   bool
	}{
		{
			name: "replace multiple with multiple",
			fields: fields{
				ring:  consistent.New(),
				nodes: sync.Map{},
			},
			args: args{
				wasNodes: []string{
					"test-node-1",
					"test-node-2",
					"test-node-3",
				},
				newNodes: []string{
					"test-node-1-1",
					"test-node-2-2",
					"test-node-3-3",
					"test-node-4-4",
					"test-node-5-5",
					"test-node-3",
				},
			},
			wantNodes: []string{
				"test-node-1-1",
				"test-node-2-2",
				"test-node-3-3",
				"test-node-4-4",
				"test-node-5-5",
				"test-node-3",
			},
			wantErr: false,
		},
		{
			name: "replace empty with multiple",
			fields: fields{
				ring:  consistent.New(),
				nodes: sync.Map{},
			},
			args: args{
				wasNodes: []string{},
				newNodes: []string{
					"test-node-1-1",
					"test-node-2-2",
					"test-node-3-3",
					"test-node-4-4",
					"test-node-5-5",
					"test-node-3",
				},
			},
			wantNodes: []string{
				"test-node-1-1",
				"test-node-2-2",
				"test-node-3-3",
				"test-node-4-4",
				"test-node-5-5",
				"test-node-3",
			},
			wantErr: false,
		},
		{
			name: "replace multiple with empty",
			fields: fields{
				ring:  consistent.New(),
				nodes: sync.Map{},
			},
			args: args{
				wasNodes: []string{
					"test-node-1",
					"test-node-2",
					"test-node-3",
				},
				newNodes: []string{},
			},
			wantNodes: nil,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			//var n nodes.Node
			//Fill it ip
			for _, name := range tt.args.wasNodes {
				tt.fields.ring.Add(name)
				mn := &mocks.Node{}
				mn.On("ID").Return(name)
				tt.fields.nodes.Store(name, mn)
			}

			var newNodes []nodes.Node
			for _, name := range tt.args.newNodes {
				mn := &mocks.Node{}
				mn.On("ID").Return(name)
				newNodes = append(newNodes, mn)
			}

			c := &Consistent{
				ring:  tt.fields.ring,
				nodes: tt.fields.nodes,
			}
			if err := c.SetNodes(newNodes); (err != nil) != tt.wantErr {
				t.Errorf("SetNodes() error = %v, wantErr %v", err, tt.wantErr)
			}

			nodes := c.ring.Hosts()
			var mapNodes []string
			c.nodes.Range(func(key, value interface{}) bool {
				mapNodes = append(mapNodes, key.(string))
				return true
			})

			sort.Strings(tt.wantNodes)
			sort.Strings(nodes)
			sort.Strings(mapNodes)

			assert.Equal(t, tt.wantNodes, nodes)
			assert.Equal(t, tt.wantNodes, mapNodes)
		})
	}
}
