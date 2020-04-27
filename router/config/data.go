package config

import (
	"github.com/pkg/errors"
	"strings"
)

type DataModeType int8

const (
	KVData DataModeType = iota + 1
	GeoData
)

func (dm *DataModeType) UnmarshalJSON(cb []byte) error {
	c := strings.ToLower(string(cb))
	c = strings.Trim(c, "\"")
	switch c {
	case "kv":
		*dm = KVData
		return nil
	case "geo":
		*dm = GeoData
		return nil
	default:
		return errors.New("unknown data mode")
	}
}
