package nodes

import (
	"context"
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/kvstore/router/nodehasher"
	"github.com/struckoff/kvstore/router/rpcapi"
	"google.golang.org/grpc"
	"log"
	"sync"
)

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

func (n *RemoteNode) Power() balancer.Power {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.p
}
func (n *RemoteNode) Capacity() balancer.Capacity {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.c
}

func (n *RemoteNode) Hash() uint64 {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return n.h
}

//Save value for a given key on the remote node
func (n *RemoteNode) Store(key string, body []byte) error {
	log.Printf("Store key(%s) on %s", key, n.id)
	req := rpcapi.KeyValue{Key: key, Value: string(body)}
	if _, err := n.rpcclient.RPCStore(context.TODO(), &req); err != nil {
		return err
	}
	return nil
}

// Save key/value pairs on remote node
func (n *RemoteNode) StorePairs(pairs []*rpcapi.KeyValue) error {
	log.Printf("Store pairs on %s", n.id)
	req := rpcapi.KeyValues{KVs: pairs}
	if _, err := n.rpcclient.RPCStorePairs(context.TODO(), &req); err != nil {
		return err
	}
	return nil
}

//Receive value for a given key from the remote node
func (n *RemoteNode) Receive(keys []string) (*rpcapi.KeyValues, error) {
	log.Printf("Receive keys(%v) from %s", keys, n.id)
	req := rpcapi.KeyReq{Keys: keys}
	res, err := n.rpcclient.RPCReceive(context.TODO(), &req)
	return res, err
	//if err != nil {
	//	return nil, err
	//}
	//
	//b, err := json.Marshal(res)
	//if err != nil {
	//	return nil, errors.Wrap(err, "failed to marshal response from external node")
	//}
	//return b, nil
}

// Explore returns the list of keys on remote node
func (n *RemoteNode) Explore() ([]string, error) {
	log.Printf("Exploring %s", n.id)
	req := rpcapi.Empty{}
	res, err := n.rpcclient.RPCExplore(context.TODO(), &req)
	if err != nil {
		return nil, err
	}
	return res.Keys, nil
}

// Remove value for a given key
func (n *RemoteNode) Remove(key []string) error {
	log.Printf("Remove key(%s) from %s", key, n.id)
	req := rpcapi.KeyReq{Keys: key}
	_, err := n.rpcclient.RPCRemove(context.TODO(), &req)
	return err
}

// Return meta information about the node
func (n *RemoteNode) Meta() *rpcapi.NodeMeta {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return &rpcapi.NodeMeta{
		ID:         n.id,
		Address:    n.address,
		RPCAddress: n.rpcaddress,
		Power:      n.p.Get(),
		Capacity:   n.p.Get(),
	}
}

// Move keys  from remote node to another remote node.
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

// NewExternalNode - create a new instance of an external by given meta information.
func NewExternalNode(meta *rpcapi.NodeMeta, hasher nodehasher.Hasher) (*RemoteNode, error) {
	conn, err := grpc.Dial(meta.RPCAddress, grpc.WithInsecure()) // TODO Make it secure
	if err != nil {
		return nil, err
	}
	c := rpcapi.NewRPCNodeClient(conn)
	var hashsum uint64
	if hasher != nil {
		hashsum, err = hasher.Sum(meta)
		if err != nil {
			return nil, err
		}
	}
	return &RemoteNode{
		id:         meta.ID,
		address:    meta.Address,
		rpcaddress: meta.RPCAddress,
		p:          NewPower(meta.Power),
		c:          NewCapacity(meta.Capacity),
		rpcclient:  c,
		h:          hashsum,
	}, nil
}

// NewExternalNodeByAddr - create a new instance of an external node.
// Function asks remote node for it meta information by RPC
func NewExternalNodeByAddr(rpcaddr string, hasher nodehasher.Hasher) (*RemoteNode, error) {
	c, err := enClient(rpcaddr)
	if err != nil {
		return nil, err
	}
	meta, err := c.RPCMeta(context.TODO(), &rpcapi.Empty{})
	if err != nil {
		return nil, err
	}
	hashsum, err := hasher.Sum(meta)
	if err != nil {
		return nil, err
	}
	return &RemoteNode{
		id:         meta.ID,
		address:    meta.Address,
		rpcaddress: meta.RPCAddress,
		p:          NewPower(meta.Power),
		c:          NewCapacity(meta.Capacity),
		rpcclient:  c,
		h:          hashsum,
	}, nil
}

func enClient(addr string) (rpcapi.RPCNodeClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure()) // TODO Make it secure
	if err != nil {
		return nil, err
	}
	c := rpcapi.NewRPCNodeClient(conn)
	return c, nil
}
