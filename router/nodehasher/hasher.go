package nodehasher

import "github.com/struckoff/kvstore/router/rpcapi"

type Hasher interface {
	Sum(*rpcapi.NodeMeta) (uint64, error)
}
