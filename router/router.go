package router

import (
	"context"
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/balanceradapter"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/dataitem"
	"github.com/struckoff/kvstore/router/nodehasher"
	"github.com/struckoff/kvstore/router/nodes"
	"github.com/struckoff/kvstore/router/rpcapi"
	balancer "github.com/struckoff/sfcframework"
	"golang.org/x/sync/errgroup"
	"log"
	"sync"
)

// Router represents bounding of network api with kvrouter lib and local node
type Router struct {
	bal    balanceradapter.Balancer
	hasher nodehasher.Hasher
	ndf    dataitem.NewDataItemFunc
	rpcndf dataitem.DataItemFromRpc
	conf   *config.BalancerConfig
	opLock sync.Mutex
}

func NewRouter(bal balanceradapter.Balancer, h nodehasher.Hasher, ndf dataitem.NewDataItemFunc, rpcndf dataitem.DataItemFromRpc, conf *config.BalancerConfig) (*Router, error) {
	r := &Router{
		bal:    bal,
		hasher: h,
		ndf:    ndf,
		rpcndf: rpcndf,
		conf:   conf,
	}
	return r, nil
}

// AddNode adds node to kvrouter
func (h *Router) AddNode(n nodes.Node) error {
	h.opLock.Lock()
	defer h.opLock.Unlock()
	if err := h.bal.AddNode(n); err != nil {
		if err := h.removeNode(n.ID()); err != nil {
			log.Println(errors.Wrap(err, "failed to RemoveNode"))
		}
		return err
	}
	if err := h.optimize(); err != nil {
		if err := h.removeNode(n.ID()); err != nil {
			log.Println(errors.Wrap(err, "failed to RemoveNode"))
		}
		if err := h.optimize(); err != nil {
			return errors.Wrap(err, "failed to optimize")
		}
		return errors.Wrap(err, "failed to optimize")
	}
	return nil
}

// RemoveNode removes node from kvrouter
func (h *Router) RemoveNode(id string) error {
	h.opLock.Lock()
	defer h.opLock.Unlock()
	return h.removeNode(id)
}
func (h *Router) removeNode(id string) error {
	log.Printf("removing node(%s)", id)
	return h.bal.RemoveNode(id)
}

// Returns node from kvrouter by given key.
func (h *Router) LocateKey(rdi *rpcapi.DataItem) (nodes.Node, uint64, error) {
	//h.opLock.Lock()
	//defer h.opLock.Unlock()
	//di, err := h.ndf(key, 0)
	di := h.rpcndf(rdi)
	nb, cid, err := h.bal.LocateData(di)
	if err != nil {
		return nil, 0, err
	}
	n, ok := nb.(nodes.Node)
	if !ok {
		return nil, 0, errors.New("wrong node type")
	}
	return n, cid, nil
}

func (h *Router) AddData(cID uint64, di balancer.DataItem) error {
	h.opLock.Lock()
	defer h.opLock.Unlock()
	return h.addData(cID, di)
}

func (h *Router) addData(cID uint64, di balancer.DataItem) error {
	return h.bal.AddData(cID, di)
}

func (h *Router) RemoveData(key string) error {
	di, err := h.ndf(key, 0)
	if err != nil {
		return err
	}
	h.opLock.Lock()
	defer h.opLock.Unlock()
	return h.bal.RemoveData(di)
}

// GetNodes - returns a list of nodes in the balancer
func (h *Router) GetNodes() ([]nodes.Node, error) {
	//h.opLock.Lock()
	//defer h.opLock.Unlock()
	return h.bal.Nodes()
}

// SetNodes - removes all nodes from the balancer and set a new ones
func (h *Router) SetNodes(ns []nodes.Node) error {
	//h.opLock.Lock()
	//defer h.opLock.Unlock()
	return h.bal.SetNodes(ns)
}

// GetNode returns the node with the given ID
func (h *Router) GetNode(id string) (nodes.Node, error) {
	//h.opLock.Lock()
	//defer h.opLock.Unlock()
	return h.bal.GetNode(id)
}

func (h *Router) Hasher() nodehasher.Hasher {
	return h.hasher
}

func (h *Router) fillBalancer() error {
	if err := h.bal.Reset(); err != nil {
		return err
	}
	ns, err := h.GetNodes()
	if err != nil {
		return err
	}

	diCh := make(chan rpcapi.DataItem)
	defer close(diCh)
	//
	if len(ns) == 0 {
		return nil
	}
	nseg, nsctx := errgroup.WithContext(context.Background())
	addeg, addctx := errgroup.WithContext(nsctx)
	addeg.Go(h.addKeysGoroutine(diCh, addctx))

	for _, n := range ns {
		nseg.Go(h.fillBalancerNode(n, diCh, addctx))
	}

	err = nseg.Wait()
	if err := addeg.Wait(); err != nil {
		return err
	}
	return err
}

func (h *Router) fillBalancerNode(n nodes.Node, diCh chan<- rpcapi.DataItem, ctx context.Context) func() error {
	return func() (err error) {
		var dis []*rpcapi.DataItem
		select {
		case <-ctx.Done():
			return nil
		default:
			dis, err = n.Explore()
			if err != nil {
				return err
			}
		}
		select {
		case <-ctx.Done():
			return nil
		default:
			for i := range dis {
				select {
				case <-ctx.Done():
					return nil
				case diCh <- *dis[i]:
				}

			}
		}
		return nil
	}
}

func (h *Router) addKeysGoroutine(diCh <-chan rpcapi.DataItem, ctx context.Context) func() error {
	return func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case rdi, ok := <-diCh:
				if !ok {
					return nil
				}
				di := h.rpcndf(&rdi)
				_, cID, err := h.bal.LocateData(di)
				if err != nil {
					return err
				}
				if err := h.addData(cID, di); err != nil {
					return err
				}
			}
		}
	}
}

func (h *Router) redistributeKeys() error {
	//var wg sync.WaitGroup

	ns, err := h.GetNodes()
	if err != nil {
		return err
	}
	eg, ectx := errgroup.WithContext(context.Background())
	for _, n := range ns {
		eg.Go(h.redistributeKeysNode(n, ectx))
	}
	return eg.Wait()
}

func (h *Router) redistributeKeysNode(n nodes.Node, ctx context.Context) func() error {
	return func() (err error) {
		res := make(map[nodes.Node][]*rpcapi.DataItem)
		var dis []*rpcapi.DataItem
		select {
		case <-ctx.Done():
			return nil
		default:
			dis, err = n.Explore()
			if err != nil {
				return err
			}
		}

		for i := range dis {
			select {
			case <-ctx.Done():
				return nil
			default:
				di := h.rpcndf(dis[i])
				en, _, err := h.bal.LocateData(di)
				if err != nil {
					return err
				}
				if en.ID() != n.ID() {
					res[en] = append(res[en], di.RPCApi())
				}
			}
		}
		select {
		case <-ctx.Done():
			return nil
		default:
			if err := n.Move(res); err != nil {
				return err
			}
		}
		return nil
	}
}
