package transform

import (
	"errors"
	"fmt"
	"github.com/struckoff/sfcframework/curve"
)

const latStep = 90.0
const lonStep = 180.0

func SpaceTransform(values []interface{}, sfc curve.Curve) ([]uint64, error) {
	dimSize := sfc.DimensionSize()
	if len(values) != 2 || sfc.Dimensions() != 2 {
		return nil, errors.New("number of dimensions must be 2")
	}
	res := make([]uint64, 2)
	lat, ok := values[0].(float64)
	if !ok {
		return nil, errors.New("first value must be float64 latitude")
	}
	if -latStep >= lat || lat >= latStep {
		return nil, fmt.Errorf("%f latitude exceeds limit [%f, %f]", lat, -latStep, latStep)
	}
	res[0] = uint64((lat + latStep) / (latStep * 2) * float64(dimSize))
	lon, ok := values[1].(float64)
	if !ok {
		return nil, errors.New("second value must be float64 longitude")
	}
	if -lonStep >= lon || lon >= lonStep {
		return nil, fmt.Errorf("%f longitude exceeds limit [%f, %f]", lon, -lonStep, lonStep)
	}
	res[1] = uint64((lon + lonStep) / (lonStep * 2) * float64(dimSize))
	return res, nil
}
