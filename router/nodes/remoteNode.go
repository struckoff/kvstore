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
	c          Capacity
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
func (n *RemoteNode) Capacity() node.Capacity {
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
func (n *RemoteNode) Store(key string, body []byte) error {
	logger.Logger().Debug("store key", zap.String("Key", key), zap.String("Node", n.id))
	req := rpcapi.KeyValue{Key: key, Value: string(body)}
	if _, err := n.rpcclient.RPCStore(context.TODO(), &req); err != nil {
		return err
	}
	return nil
}

// Save key/value pairs on remote node
func (n *RemoteNode) StorePairs(pairs []*rpcapi.KeyValue) error {
	logger.Logger().Debug("store pairs", zap.String("Node", n.id))
	req := rpcapi.KeyValues{KVs: pairs}
	if _, err := n.rpcclient.RPCStorePairs(context.TODO(), &req); err != nil {
		return err
	}
	return nil
}

//Receive value for a given key from the remote node
func (n *RemoteNode) Receive(keys []string) (*rpcapi.KeyValues, error) {
	logger.Logger().Debug("receive keys", zap.Strings("Keys", keys), zap.String("Node", n.id))
	req := rpcapi.KeyReq{Keys: keys}
	res, err := n.rpcclient.RPCReceive(context.TODO(), &req)
	return res, err
}

// Explore returns the list of keys on remote node
func (n *RemoteNode) Explore() ([]string, error) {
	logger.Logger().Debug("exploring", zap.String("Node", n.id))
	req := rpcapi.Empty{}
	res, err := n.rpcclient.RPCExplore(context.TODO(), &req)
	if err != nil {
		return nil, err
	}
	return res.Keys, nil
}

// Remove value for a given key
func (n *RemoteNode) Remove(keys []string) error {
	logger.Logger().Debug("remove keys", zap.Strings("Keys", keys), zap.String("Node", n.id))
	req := rpcapi.KeyReq{Keys: keys}
	_, err := n.rpcclient.RPCRemove(context.TODO(), &req)
	return err
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
func (n *RemoteNode) Move(nk map[Node][]string) error {
	mr := &rpcapi.MoveReq{}
	for en, keys := range nk {
		meta := en.Meta()
		kl := &rpcapi.KeyList{
			Node: meta,
			Keys: keys,
		}
		mr.KL = append(mr.KL, kl)
	}
	_, err := n.rpcclient.RPCMove(context.TODO(), mr)
	return err
}
