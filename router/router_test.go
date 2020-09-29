package router

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/struckoff/kvstore/mocks"
	"github.com/struckoff/kvstore/router/nodes"
	balancer "github.com/struckoff/sfcframework"
	balancermocs "github.com/struckoff/sfcframework/mocks"
	"testing"
)

func TestRouter_AddNode(t *testing.T) {
	mn := &mocks.Node{}
	mn.On("Explore").Return([]string{}, nil)
	mn.On("Move", mock.Anything).Return(nil)

	mbal := &mocks.Balancer{}
	mbal.On("AddNode", mn).Return(nil)
	mbal.On("Reset").Return(nil)
	mbal.On("Nodes").Return([]nodes.Node{mn}, nil)
	mbal.On("Optimize").Return(nil)
	mbal.On("AddNode", mock.Anything).Return(nil)
	h := &Router{
		bal: mbal,
	}
	if err := h.AddNode(mn); err != nil {
		t.Errorf("AddNode() error = %v", err)
	}

	mbal.AssertCalled(t, "AddNode", mn)
}

func TestRouter_RemoveNode(t *testing.T) {
	name := "test-node"

	mbal := &mocks.Balancer{}
	mbal.On("RemoveNode", name).Return(nil)
	h := &Router{
		bal: mbal,
	}
	if err := h.RemoveNode(name); err != nil {
		t.Errorf("RemoveNode() error = %v", err)
	}

	mbal.AssertCalled(t, "RemoveNode", name)
}

func TestRouter_LocateKey(t *testing.T) {
	key := "test-key"

	name := "test-node"
	mn := &mocks.Node{}
	mn.On("ID").Return(name)

	mbal := &mocks.Balancer{}
	mbal.On("LocateData", mock.AnythingOfType("*mocks.DataItem")).Return(mn, uint64(1), nil)

	h := &Router{
		bal: mbal,
		ndf: func(s string) (balancer.DataItem, error) {
			di := &balancermocs.DataItem{}
			di.On("ID").Return(s)
			return di, nil
		},
	}
	n, err := h.LocateKey(key)
	if err != nil {
		t.Errorf("LocateKey() error = %v", err)
	}
	mbal.AssertCalled(t, "LocateData", mock.AnythingOfType("*mocks.DataItem"))
	assert.Equal(t, mn.ID(), n.ID())
}

func TestRouter_SetNodes(t *testing.T) {
	ns := []nodes.Node{
		&mocks.Node{},
		&mocks.Node{},
		&mocks.Node{},
	}

	mbal := &mocks.Balancer{}
	mbal.On("SetNodes", ns).Return(nil)

	h := &Router{
		bal: mbal,
		ndf: func(s string) (balancer.DataItem, error) {
			di := &balancermocs.DataItem{}
			di.On("ID").Return(s)
			return di, nil
		},
	}

	if err := h.SetNodes(ns); err != nil {
		t.Errorf("SetNodes() error = %v", err)
	}
	mbal.AssertCalled(t, "SetNodes", ns)
}

func TestRouter_GetNode(t *testing.T) {
	name := "test-node"
	mn := &mocks.Node{}
	mn.On("ID").Return(name)

	mbal := &mocks.Balancer{}
	mbal.On("GetNode", name).Return(mn, nil)

	h := &Router{
		bal: mbal,
		ndf: func(s string) (balancer.DataItem, error) {
			di := &balancermocs.DataItem{}
			di.On("ID").Return(s)
			return di, nil
		},
	}
	n, err := h.GetNode(name)
	if err != nil {
		t.Errorf("GetNode() error = %v", err)
	}
	mbal.AssertCalled(t, "GetNode", name)
	assert.Equal(t, mn.ID(), n.ID())
}
