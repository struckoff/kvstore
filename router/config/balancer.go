package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"strings"
)

// If config implies use of consul, this options will be taken from consul KV.
// Otherwise it will be taken from config file.
type BalancerConfig struct {
	// type of balancer to use. Possible: SFC, Consistent.
	Mode BalancerModeType `envconfig:"MODE"`
	// SFC configuration.
	SFC *SFCConfig `json:",omitempty", envconfig:"SFC"`
	// Consistent ring configuration
	//Ring *consistent.Config `json:",omitempty", envconfig:"RING"`
	// Which way to use for node hashing and sorting
	// Possible: geosfc, xxhash.
	NodeHash NodeHashType `envconfig:"NODE_HASH"`
	// Which data to store in the database
	// Possible: kv, geo
	DataMode DataModeType `envconfig:"DATA_MODE"`
}

func (bc *BalancerConfig) UnmarshalJSON(cb []byte) error {
	type clone BalancerConfig
	if err := json.Unmarshal(cb, (*clone)(bc)); err != nil {
		return err
	}
	switch bc.Mode {
	case SFCMode:
		if bc.SFC == nil {
			return errors.New("unable to find SFC config")
		}
	case ConsistentMode:
		//if bc.Ring == nil {
		//	return errors.New("unable to find consistent ring config")
		//}
		if bc.NodeHash == GeoSfc {
			return errors.New("SFC node hasher should be used with SFC balancer")
		}
	default:
		return errors.New("wrong balancer mode")
	}
	return nil
}

type BalancerModeType int8

const (
	SFCMode BalancerModeType = iota + 1
	ConsistentMode
)

func (bm *BalancerModeType) UnmarshalJSON(cb []byte) error {
	c := strings.ToLower(string(cb))
	c = strings.Trim(c, "\"")
	switch c {
	case "sfc":
		*bm = SFCMode
		return nil
	case "consistent":
		*bm = ConsistentMode
		return nil
	default:
		return errors.New("unknown balancer mode")
	}
}