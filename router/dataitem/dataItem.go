package dataitem

import (
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/rpcapi"
	balancer "github.com/struckoff/sfcframework"
)

type DataItem interface {
	balancer.DataItem
	RPCApi() *rpcapi.DataItem
}

type NewDataItemFunc func(string, uint64) (DataItem, error)

func GetDataItemFunc(dmt config.DataModeType) (NewDataItemFunc, error) {
	switch dmt {
	case config.KVData:
		return NewKVDataItem, nil
	case config.GeoData:
		return NewSpaceDataItem, nil
	default:
		return nil, errors.New("wrong data mode")
	}
}

type DataItemFromRpc func(*rpcapi.DataItem) DataItem

func GetDataItemFromRpcFunc(dmt config.DataModeType) (DataItemFromRpc, error) {
	switch dmt {
	case config.KVData:
		return KVDataItemFromRPC, nil
	case config.GeoData:
		return SpaceDataItemFromRPC, nil
	default:
		return nil, errors.New("wrong data mode")
	}
}
