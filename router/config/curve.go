package config

import (
	"github.com/pkg/errors"
	"github.com/struckoff/sfcframework/curve"
	"strings"
)

type CurveType struct {
	curve.CurveType
}

func (ct *CurveType) MarshalJSON() ([]byte, error) {
	return []byte("\"" + ct.CurveType.String() + "\""), nil
}

func (ct *CurveType) Decode(c string) error {
	return ct.UnmarshalJSON([]byte(c))
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
