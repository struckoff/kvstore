package nodehasher

import (
	"github.com/cespare/xxhash"
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/rpcapi"
)

type XXHash struct{}

func (xh XXHash) Sum(meta *rpcapi.NodeMeta) (uint64, error) {
	if meta == nil {
		return 0, errors.New("meta data not found")
	}
	return xxhash.Sum64String(meta.ID), nil
}

func NewXXHash() XXHash {
	return XXHash{}
}
