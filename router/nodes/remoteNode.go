package nodes

import (
	"context"
	"github.com/struckoff/kvstore/logger"
	"github.com/struckoff/kvstore/router/nodehasher"
	"github.com/struckoff/kvstore/router/rpcapi"
	"github.com/struckoff/sfcframework/node"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"sync"
	"time"
)

const MaxCallSendMsgSize = 100 * 1024 * 1024
const MaxRecvSendMsgSize = 100 * 1024 * 1024

// NewExternalNode - create a new instance of an external by given meta information.
func NewExternalNode(meta *rpcapi.NodeMeta, hasher nodehasher.Hasher) (*RemoteNode, error) {
	nodeC, capC, err := enClient(meta.RPCAddress)
	if err != nil {
		return nil, err
	}
	var hashsum uint64
	if hasher != nil {
		hashsum, err = hasher.Sum(meta)
		if err != nil {
			return nil, err
		}
	}
	en := newExternalNode(meta, nodeC, capC, hashsum)
	return en, nil
}

// NewExternalNodeByAddr - create a new instance of an external node.
// Function asks remote node for it meta information by RPC
func NewExternalNodeByAddr(rpcaddr string, hasher nodehasher.Hasher) (*RemoteNode, error) {
	nodeC, capC, err := enClient(rpcaddr)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	meta, err := nodeC.RPCMeta(ctx, &rpcapi.Empty{}, grpc.WaitForReady(true))
	if err != nil {
		return nil, err
	}
	hashsum, err := hasher.Sum(meta)
	if err != nil {
		return nil, err
	}
	en := newExternalNode(meta, nodeC, capC, hashsum)
	return en, nil
}

func newExternalNode(meta *rpcapi.NodeMeta, nodeC rpcapi.RPCNodeClient, c rpcapi.RPCCapacityClient, h uint64) *RemoteNode {
	return &RemoteNode{
		id:         meta.ID,
		address:    meta.Address,
		rpcaddress: meta.RPCAddress,
		p:          NewPower(meta.Power),
		c:          NewCapacity(c),
		rpcclient:  nodeC,
		h:          h,
		geo:        meta.Geo,
	}
}

func enClient(addr string) (rpcapi.RPCNodeClient, rpcapi.RPCCapacityClient, error) {
	// TODO Make it secure
	conn, err := grpc.Dial(addr, grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallSendMsgSize(MaxCallSendMsgSize),
			grpc.MaxCallRecvMsgSize(MaxRecvSendMsgSize),
		))

	if err != nil {
		return nil, nil, err
	}
	nodeC := rpcapi.NewRPCNodeClient(conn)
	capC := rpcapi.NewRPCCapacityClient(conn)
	return nodeC, capC, nil
}

// RemoteNode represents compunction API with cluster unit
// It also contains meta information
type RemoteNode struct {
	mu         sync.RWMutex
	id         string // uniq node ID
	address    string // node HTTP address
	rpcaddress string // node RPC address
	p          Power
	c          RemoteCapacity
	rpcclient  rpcapi.RPCNodeClient
	geo        *rpcapi.GeoData
	h          uint64
}

// ID  returns the node ID
func (n *RemoteNode) ID() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.id
}

func (n *RemoteNode) Power() node.Power {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.p
}
func (n *RemoteNode) Capacity() Capacity {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return &n.c
}

func (n *RemoteNode) Hash() uint64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.h
}

//Save value for a given key on the remote node
func (n *RemoteNode) Store(kv *rpcapi.KeyValue) (*rpcapi.DataItem, error) {
	logger.Logger().Debug("store key", zap.String("Key", string(kv.Key.ID)), zap.String("Node", n.id))
	di, err := n.rpcclient.RPCStore(context.TODO(), kv)
	if err != nil {
		return nil, err
	}
	return di, nil
}

// Save key/value pairs on remote node
func (n *RemoteNode) StorePairs(pairs []*rpcapi.KeyValue) ([]*rpcapi.DataItem, error) {
	logger.Logger().Debug("store pairs", zap.String("Node", n.id))
	req := rpcapi.KeyValues{KVs: pairs}
	dis, err := n.rpcclient.RPCStorePairs(context.TODO(), &req)
	if err != nil {
		return nil, err
	}
	return dis.DIs, nil
}

//Receive value for a given key from the remote node
func (n *RemoteNode) Receive(dis []*rpcapi.DataItem) (*rpcapi.KeyValues, error) {
	logger.Logger().Debug("receive keys", zap.String("Node", n.id))
	req := rpcapi.DataItems{DIs: dis}
	res, err := n.rpcclient.RPCReceive(context.TODO(), &req)
	return res, err
}

// Explore returns the list of keys on remote node
func (n *RemoteNode) Explore() ([]*rpcapi.DataItem, error) {
	logger.Logger().Debug("exploring", zap.String("Node", n.id))
	req := rpcapi.Empty{}
	dis, err := n.rpcclient.RPCExplore(context.TODO(), &req)
	if err != nil {
		return nil, err
	}
	return dis.DIs, nil
}

// Remove value for a given key
func (n *RemoteNode) Remove(dis []*rpcapi.DataItem) ([]*rpcapi.DataItem, error) {
	logger.Logger().Debug("remove keys", zap.String("Node", n.id))
	req := rpcapi.DataItems{DIs: dis}
	ds, err := n.rpcclient.RPCRemove(context.TODO(), &req)
	if err != nil {
		return nil, err
	}
	return ds.DIs, nil
}

// Return meta information about the node
func (n *RemoteNode) Meta() *rpcapi.NodeMeta {
	n.mu.RLock()
	defer n.mu.RUnlock()
	cp, err := n.c.Get()
	if err != nil {
		return nil
	}
	return &rpcapi.NodeMeta{
		ID:         n.id,
		Address:    n.address,
		RPCAddress: n.rpcaddress,
		Power:      n.p.Get(),
		Capacity:   cp,
		Geo:        n.geo,
	}
}

// Move keys from remote node to another remote node.
func (n *RemoteNode) Move(nk map[Node][]*rpcapi.DataItem) error {
	mr := &rpcapi.MoveReq{}
	for en, dis := range nk {
		meta := en.Meta()
		kl := &rpcapi.KeyList{
			Node: meta,
			Keys: &rpcapi.DataItems{DIs: dis},
		}
		mr.KLs = append(mr.KLs, kl)
	}
	_, err := n.rpcclient.RPCMove(context.TODO(), mr)
	return err
}
