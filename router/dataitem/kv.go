package dataitem

import balancer "github.com/struckoff/sfcframework"

type KVDataItem string

func NewKVDataItem(key string) (balancer.DataItem, error) {
	return KVDataItem(key), nil
}

func (di KVDataItem) ID() string {
	return string(di)
}

func (di KVDataItem) Size() uint64 {
	return 1
}

func (di KVDataItem) Values() []interface{} {
	return []interface{}{string(di)}
}
