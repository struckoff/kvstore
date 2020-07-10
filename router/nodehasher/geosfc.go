package nodehasher

import (
	"github.com/pkg/errors"
	"github.com/struckoff/SFCFramework/curve"
	"github.com/struckoff/SFCFramework/transform"
	"github.com/struckoff/kvstore/router/rpcapi"
)

type GeoSfc struct {
	sfc curve.Curve
}

func (gs GeoSfc) Sum(meta *rpcapi.NodeMeta) (uint64, error) {
	if meta == nil {
		return 0, errors.New("meta data not found")
	}
	if meta.Geo == nil {
		return 0, errors.New("geo data not found")
	}

	coords, err := gs.transform(meta.Geo.Latitude, meta.Geo.Latitude)
	if err != nil {
		return 0, errors.Wrap(err, "unable to transform coordinates into layers")
	}
	h, err := gs.sfc.Encode(coords)
	if err != nil {
		return 0, errors.Wrap(err, "unable to encode coordinates with SFC")
	}
	return h, nil
}

func (gs GeoSfc) transform(lat, lon float64) ([]uint64, error) {
	vals := []interface{}{lat, lon}
	return transform.SpaceTransform(vals, gs.sfc)
}

func NewGeoSfc(sfc curve.Curve) GeoSfc {
	return GeoSfc{sfc}
}
