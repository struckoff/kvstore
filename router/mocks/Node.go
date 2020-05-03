// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import balancer "github.com/struckoff/SFCFramework"
import mock "github.com/stretchr/testify/mock"
import nodes "github.com/struckoff/kvstore/router/nodes"
import rpcapi "github.com/struckoff/kvstore/router/rpcapi"

// Node is an autogenerated mock type for the Node type
type Node struct {
	mock.Mock
}

// Capacity provides a mock function with given fields:
func (_m *Node) Capacity() balancer.Capacity {
	ret := _m.Called()

	var r0 balancer.Capacity
	if rf, ok := ret.Get(0).(func() balancer.Capacity); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(balancer.Capacity)
		}
	}

	return r0
}

// Explore provides a mock function with given fields:
func (_m *Node) Explore() ([]string, error) {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Hash provides a mock function with given fields:
func (_m *Node) Hash() uint64 {
	ret := _m.Called()

	var r0 uint64
	if rf, ok := ret.Get(0).(func() uint64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint64)
	}

	return r0
}

// ID provides a mock function with given fields:
func (_m *Node) ID() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Meta provides a mock function with given fields:
func (_m *Node) Meta() *rpcapi.NodeMeta {
	ret := _m.Called()

	var r0 *rpcapi.NodeMeta
	if rf, ok := ret.Get(0).(func() *rpcapi.NodeMeta); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rpcapi.NodeMeta)
		}
	}

	return r0
}

// Move provides a mock function with given fields: _a0
func (_m *Node) Move(_a0 map[nodes.Node][]string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(map[nodes.Node][]string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Power provides a mock function with given fields:
func (_m *Node) Power() balancer.Power {
	ret := _m.Called()

	var r0 balancer.Power
	if rf, ok := ret.Get(0).(func() balancer.Power); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(balancer.Power)
		}
	}

	return r0
}

// Receive provides a mock function with given fields: _a0
func (_m *Node) Receive(_a0 []string) (*rpcapi.KeyValues, error) {
	ret := _m.Called(_a0)

	var r0 *rpcapi.KeyValues
	if rf, ok := ret.Get(0).(func([]string) *rpcapi.KeyValues); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rpcapi.KeyValues)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func([]string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Remove provides a mock function with given fields: _a0
func (_m *Node) Remove(_a0 []string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func([]string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Store provides a mock function with given fields: _a0, _a1
func (_m *Node) Store(_a0 string, _a1 []byte) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []byte) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StorePairs provides a mock function with given fields: _a0
func (_m *Node) StorePairs(_a0 []*rpcapi.KeyValue) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*rpcapi.KeyValue) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}