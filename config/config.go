package config

import (
	"github.com/pkg/errors"
	"github.com/struckoff/SFCFramework/curve"
	"strings"
)

type Config struct {
	Name        string
	Address     string
	RPCAddress  string
	Power       float64
	Capacity    float64
	DBpath      string
	Entrypoints []string  //If not empty node tries to connect to each entrypoint, send its meta and receive cluster info
	Dimensions  uint64    //Amount of space filling curve dimensions
	Size        uint64    //Size of space filling curve
	Curve       CurveType //Space filling curve type
}

type CurveType struct {
	curve.CurveType
}

func (ct *CurveType) UnmarshalJSON(cb []byte) error {
	c := strings.ToLower(string(cb))
	c = strings.Trim(c, "\"")
	switch c {
	case "morton":
		ct.CurveType = curve.Morton
		return nil
	case "hilbert":
		ct.CurveType = curve.Hilbert
		return nil
	default:
		return errors.New("unknown curve type")
	}
}
