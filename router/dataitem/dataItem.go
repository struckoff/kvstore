package dataitem

import (
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/config"
	balancer "github.com/struckoff/sfcframework"
)

type NewDataItemFunc func(string) (balancer.DataItem, error)

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
