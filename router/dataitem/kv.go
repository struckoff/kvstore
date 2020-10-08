package dataitem

import (
	"github.com/struckoff/kvstore/router/rpcapi"
)

type KVDataItem struct {
	string
	size uint64
}

func KVDataItemFromRPC(rdi *rpcapi.DataItem) DataItem {
	return KVDataItem{string(rdi.ID), rdi.Size}
}

func NewKVDataItem(k string, size uint64) (DataItem, error) {
	return KVDataItem{string: k, size: size}, nil
}

func (di KVDataItem) ID() string {
	return di.string
}

func (di KVDataItem) Size() uint64 {
	return di.size
}

func (di KVDataItem) Values() []interface{} {
	return []interface{}{di.string}
}

func (di KVDataItem) RPCApi() *rpcapi.DataItem {
	return &rpcapi.DataItem{
		ID:   []byte(di.string),
		Size: di.size,
	}
}
