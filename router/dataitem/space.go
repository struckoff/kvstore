package dataitem

import (
	"encoding/json"
	balancer "github.com/struckoff/SFCFramework"
)

// SpaceDataItem represents geospatial key as balancer item
type SpaceDataItem struct {
	Key string
	Lat float64
	Lon float64
}

func NewSpaceDataItem(key string) (balancer.DataItem, error) {
	var item SpaceDataItem
	err := json.Unmarshal([]byte(key), &item)
	item.Key = key
	return item, err
}

func (di SpaceDataItem) ID() string {
	return di.Key
}

func (di SpaceDataItem) Size() uint64 {
	return 1
}

func (di SpaceDataItem) Values() []interface{} {
	return []interface{}{di.Lat, di.Lon}
}
