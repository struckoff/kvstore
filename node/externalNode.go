package node

import (
	"context"
	balancer "github.com/struckoff/SFCFramework"
	"github.com/struckoff/kvstore/proto"
	"google.golang.org/grpc"
	"log"
	"sync"
)

// ExternalNode represents compunction API with cluster unit
// It also contains meta information
type ExternalNode struct {
	mu         sync.RWMutex
	id         string
	address    string
	rpcaddress string
	p          Power
	c          Capacity
	rpcclient  proto.RPCListenerClient
}

func (n *ExternalNode) ID() string                  { return n.id }
func (n *ExternalNode) Power() balancer.Power       { return n.p }
func (n *ExternalNode) Capacity() balancer.Capacity { return n.c }

//Save value for a given key on the remote node
func (n *ExternalNode) Store(key string, body []byte) error {
	log.Printf("Store key(%s) on %s", key, n.id)
	req := proto.StoreReq{Key: key, Value: body}
	if _, err := n.rpcclient.RPCStore(context.TODO(), &req); err != nil {
		return err
	}

	//p := strings.Join([]string{"http:/", n.address, "kv", key}, "/")
	//b := bytes.NewBuffer(body)
	//r, err := http.Post(p, "application/text", b)
	//if err != nil {
	//	return err
	//}
	//if r.StatusCode >= 400 {
	//	return errors.New(r.Status)
	//}
	return nil
}

//Receive value for a given key from the remote node
func (n *ExternalNode) Receive(key string) ([]byte, error) {
	log.Printf("Receive key(%s) from %s", key, n.id)
	req := proto.ReceiveReq{Key: key}
	res, err := n.rpcclient.RPCReceive(context.TODO(), &req)
	if err != nil {
		return nil, err
	}
	return res.Value, nil
	//
	//p := strings.Join([]string{"http:/", n.address, "kv", key}, "/")
	//r, err := http.Get(p)
	//if err != nil {
	//	return nil, err
	//}
	//if r.StatusCode >= 400 {
	//	return nil, errors.New(r.Status)
	//}
	//defer r.Body.Close()
	//b, err := ioutil.ReadAll(r.Body)
	//return b, err
}

//TODO: implement
func (n *ExternalNode) Explore() ([]string, error) {
	log.Printf("Exploring %s", n.id)
	req := proto.ExploreReq{}
	res, err := n.rpcclient.RPCExplore(context.TODO(), &req)
	if err != nil {
		return nil, err
	}
	return res.Keys, nil
}

// Return meta information about the node
func (n *ExternalNode) Meta() NodeMeta {
	n.mu.RLock()
	defer n.mu.RUnlock()
	return NodeMeta{
		ID:       n.id,
		Address:  n.address,
		Power:    n.p.Get(),
		Capacity: n.p.Get(),
	}
}

func NewExternalNode(meta NodeMeta) (*ExternalNode, error) {
	conn, err := grpc.Dial(meta.RPCAddress, grpc.WithInsecure()) // TODO Make it secure
	if err != nil {
		return nil, err
	}
	c := proto.NewRPCListenerClient(conn)
	return &ExternalNode{
		id:         meta.ID,
		address:    meta.Address,
		rpcaddress: meta.RPCAddress,
		p:          NewPower(meta.Power),
		c:          NewCapacity(meta.Capacity),
		rpcclient:  c,
	}, nil
}
