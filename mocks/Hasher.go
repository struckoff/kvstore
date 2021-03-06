// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	rpcapi "github.com/struckoff/kvstore/router/rpcapi"
)

// Hasher is an autogenerated mock type for the Hasher type
type Hasher struct {
	mock.Mock
}

// Sum provides a mock function with given fields: _a0
func (_m *Hasher) Sum(_a0 *rpcapi.NodeMeta) (uint64, error) {
	ret := _m.Called(_a0)

	var r0 uint64
	if rf, ok := ret.Get(0).(func(*rpcapi.NodeMeta) uint64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*rpcapi.NodeMeta) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
