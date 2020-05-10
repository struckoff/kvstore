package config

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"strings"
	"time"
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
	Latency  Duration     `envconfig:"HTTP_LATENCY"`
}

//func (bc *BalancerConfig) MarshalJSON() ([]byte, error) {
//
//}

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

func (bm *BalancerModeType) MarshalJSON() ([]byte, error) {
	var s string
	switch *bm {
	case SFCMode:
		s = "SFC"
	case ConsistentMode:
		s = "Consistent"
	default:
		return nil, errors.New("unknown balancer mode")
	}
	return []byte("\"" + s + "\""), nil
}

func (bm *BalancerModeType) Decode(c string) error {
	return bm.UnmarshalJSON([]byte(c))
}

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

type Duration struct {
	time.Duration
}

func (d *Duration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + d.String() + "\""), nil
}

func (d *Duration) Decode(c string) error {
	return d.UnmarshalJSON([]byte(c))
}

func (d *Duration) UnmarshalJSON(cb []byte) error {
	dur, err := time.ParseDuration(string(cb))
	if err != nil {
		return err
	}
	d.Duration = dur
	log.Printf("HTTP latency: %s", d.Duration.String())
	return nil
}
