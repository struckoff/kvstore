// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	rpcapi "github.com/struckoff/kvstore/router/rpcapi"
)

// RPCNodeServer is an autogenerated mock type for the RPCNodeServer type
type RPCNodeServer struct {
	mock.Mock
}

// RPCExplore provides a mock function with given fields: _a0, _a1
func (_m *RPCNodeServer) RPCExplore(_a0 context.Context, _a1 *rpcapi.Empty) (*rpcapi.DataItems, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *rpcapi.DataItems
	if rf, ok := ret.Get(0).(func(context.Context, *rpcapi.Empty) *rpcapi.DataItems); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rpcapi.DataItems)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *rpcapi.Empty) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RPCMeta provides a mock function with given fields: _a0, _a1
func (_m *RPCNodeServer) RPCMeta(_a0 context.Context, _a1 *rpcapi.Empty) (*rpcapi.NodeMeta, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *rpcapi.NodeMeta
	if rf, ok := ret.Get(0).(func(context.Context, *rpcapi.Empty) *rpcapi.NodeMeta); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rpcapi.NodeMeta)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *rpcapi.Empty) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RPCMove provides a mock function with given fields: _a0, _a1
func (_m *RPCNodeServer) RPCMove(_a0 context.Context, _a1 *rpcapi.MoveReq) (*rpcapi.Empty, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *rpcapi.Empty
	if rf, ok := ret.Get(0).(func(context.Context, *rpcapi.MoveReq) *rpcapi.Empty); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rpcapi.Empty)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *rpcapi.MoveReq) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RPCReceive provides a mock function with given fields: _a0, _a1
func (_m *RPCNodeServer) RPCReceive(_a0 context.Context, _a1 *rpcapi.DataItems) (*rpcapi.KeyValues, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *rpcapi.KeyValues
	if rf, ok := ret.Get(0).(func(context.Context, *rpcapi.DataItems) *rpcapi.KeyValues); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rpcapi.KeyValues)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *rpcapi.DataItems) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RPCRemove provides a mock function with given fields: _a0, _a1
func (_m *RPCNodeServer) RPCRemove(_a0 context.Context, _a1 *rpcapi.DataItems) (*rpcapi.DataItems, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *rpcapi.DataItems
	if rf, ok := ret.Get(0).(func(context.Context, *rpcapi.DataItems) *rpcapi.DataItems); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rpcapi.DataItems)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *rpcapi.DataItems) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RPCStore provides a mock function with given fields: _a0, _a1
func (_m *RPCNodeServer) RPCStore(_a0 context.Context, _a1 *rpcapi.KeyValue) (*rpcapi.DataItem, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *rpcapi.DataItem
	if rf, ok := ret.Get(0).(func(context.Context, *rpcapi.KeyValue) *rpcapi.DataItem); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rpcapi.DataItem)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *rpcapi.KeyValue) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RPCStorePairs provides a mock function with given fields: _a0, _a1
func (_m *RPCNodeServer) RPCStorePairs(_a0 context.Context, _a1 *rpcapi.KeyValues) (*rpcapi.DataItems, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *rpcapi.DataItems
	if rf, ok := ret.Get(0).(func(context.Context, *rpcapi.KeyValues) *rpcapi.DataItems); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*rpcapi.DataItems)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *rpcapi.KeyValues) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
