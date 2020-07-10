package dataitem

import (
	"github.com/pkg/errors"
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/kvstore/router/config"
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
