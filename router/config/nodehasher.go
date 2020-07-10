package config

import (
	"github.com/pkg/errors"
	"strings"
)

// TODO: consider using plugins
type NodeHashType int

const (
	GeoSfc NodeHashType = iota + 1
	XXHash
)

func (nh *NodeHashType) MarshalJSON() ([]byte, error) {
	var s string
	switch *nh {
	case GeoSfc:
		s = "GeoSFC"
	case XXHash:
		s = "xxhash"
	default:
		return nil, errors.New("wrong node hash type")
	}
	return []byte("\"" + s + "\""), nil
}

func (nh *NodeHashType) Decode(c string) error {
	return nh.UnmarshalJSON([]byte(c))
}

func (nh *NodeHashType) UnmarshalJSON(cb []byte) error {
	c := strings.ToLower(string(cb))
	c = strings.Trim(c, "\"")
	switch c {
	case "geosfc":
		*nh = GeoSfc
	case "xxhash":
		*nh = XXHash
	default:
		return errors.New("wrong node hash type")
	}
	return nil
}
