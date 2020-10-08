package dataitem

import (
	"encoding/json"
	"github.com/struckoff/kvstore/router/rpcapi"
)

// SpaceDataItem represents geospatial key as balancer item
type SpaceDataItem struct {
	Key  string
	size uint64
	Lat  float64
	Lon  float64
}

func SpaceDataItemFromRPC(rdi *rpcapi.DataItem) DataItem {
	return SpaceDataItem{string(rdi.ID), rdi.Size, rdi.Geo.Latitude, rdi.Geo.Longitude}
}

func NewSpaceDataItem(key string, size uint64) (DataItem, error) {
	var item SpaceDataItem
	err := json.Unmarshal([]byte(key), &item)
	item.Key = key
	item.size = size
	return item, err
}

func (di SpaceDataItem) ID() string {
	return di.Key
}

func (di SpaceDataItem) Size() uint64 {
	return di.size
}

func (di SpaceDataItem) Values() []interface{} {
	return []interface{}{di.Lat, di.Lon}
}

func (di SpaceDataItem) RPCApi() *rpcapi.DataItem {
	return &rpcapi.DataItem{
		ID:   []byte(di.Key),
		Size: di.size,
		Geo: &rpcapi.GeoData{
			Longitude: di.Lon,
			Latitude:  di.Lat,
		},
	}
}
