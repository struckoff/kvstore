package nodehasher

import (
	"github.com/OneOfOne/xxhash"
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/rpcapi"
)

type XXHash struct{}

func (xh XXHash) Sum(meta *rpcapi.NodeMeta) (uint64, error) {
	if meta == nil {
		return 0, errors.New("meta data not found")
	}
	hasher := xxhash.New64()
	_, err := hasher.WriteString(meta.ID)
	if err != nil {
		return 0, err
	}
	return hasher.Sum64(), nil
}

func NewXXHash() XXHash {
	return XXHash{}
}
