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

func (dn *NodeHashType) UnmarshalJSON(cb []byte) error {
	c := strings.ToLower(string(cb))
	c = strings.Trim(c, "\"")
	switch c {
	case "geosfc":
		*dn = GeoSfc
	case "xxhash":
		*dn = XXHash
	default:
		return errors.New("wrong node hash type")
	}
	return nil
}
