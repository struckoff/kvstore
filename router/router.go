package router

import (
	"context"
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/router/balanceradapter"
	"github.com/struckoff/kvstore/router/config"
	"github.com/struckoff/kvstore/router/dataitem"
	"github.com/struckoff/kvstore/router/nodehasher"
	"github.com/struckoff/kvstore/router/nodes"
	"golang.org/x/sync/errgroup"
	"log"
	"sync"
)

// Router represents bounding of network api with kvrouter lib and local node
type Router struct {
	bal    balanceradapter.Balancer
	hasher nodehasher.Hasher
	ndf    dataitem.NewDataItemFunc
	conf   *config.BalancerConfig
	opLock sync.Mutex
}

func NewRouter(bal balanceradapter.Balancer, h nodehasher.Hasher, ndf dataitem.NewDataItemFunc, conf *config.BalancerConfig) (*Router, error) {
	r := &Router{
		bal:    bal,
		hasher: h,
		ndf:    ndf,
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
func (h *Router) LocateKey(key string) (nodes.Node, error) {
	//h.opLock.Lock()
	//defer h.opLock.Unlock()
	di, err := h.ndf(key)
	if err != nil {
		return nil, err
	}
	nb, _, err := h.bal.LocateData(di)
	if err != nil {
		return nil, err
	}
	n, ok := nb.(nodes.Node)
	if !ok {
		return nil, errors.New("wrong node type")
	}
	return n, nil
}

func (h *Router) AddData(key string) (nodes.Node, error) {
	h.opLock.Lock()
	defer h.opLock.Unlock()
	return h.addData(key)
}

func (h *Router) addData(key string) (nodes.Node, error) {
	di, err := h.ndf(key)
	if err != nil {
		return nil, err
	}
	nb, _, err := h.bal.AddData(di)
	if err != nil {
		return nil, err
	}
	n, ok := nb.(nodes.Node)
	if !ok {
		return nil, errors.New("wrong node type")
	}
	return n, nil
}

func (h *Router) RemoveData(key string) error {
	di, err := h.ndf(key)
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

	keyCh := make(chan string)
	defer close(keyCh)
	//
	if len(ns) == 0 {
		return nil
	}
	nseg, nsctx := errgroup.WithContext(context.Background())
	addeg, addctx := errgroup.WithContext(nsctx)
	addeg.Go(h.addKeysGoroutine(keyCh, addctx))

	for _, n := range ns {
		nseg.Go(h.fillBalancerNode(n, keyCh, addctx))
	}

	err = nseg.Wait()
	if err := addeg.Wait(); err != nil {
		return err
	}
	return err
}

func (h *Router) fillBalancerNode(n nodes.Node, keyCh chan<- string, ctx context.Context) func() error {
	return func() (err error) {
		var keys []string
		select {
		case <-ctx.Done():
			return nil
		default:
			keys, err = n.Explore()
			if err != nil {
				return err
			}
		}
		select {
		case <-ctx.Done():
			return nil
		default:
			for iter := range keys {
				select {
				case <-ctx.Done():
					return nil
				case keyCh <- keys[iter]:
				}
			}
		}
		return nil
	}
}

func (h *Router) addKeysGoroutine(keysCh <-chan string, ctx context.Context) func() error {
	return func() error {
		for {
			select {
			case <-ctx.Done():
				return nil
			case key, ok := <-keysCh:
				if !ok {
					return nil
				}
				if _, err := h.addData(key); err != nil {
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
		res := make(map[nodes.Node][]string)
		var keys []string
		select {
		case <-ctx.Done():
			return nil
		default:
			keys, err = n.Explore()
			if err != nil {
				return err
			}
		}

		for iter := range keys {
			select {
			case <-ctx.Done():
				return nil
			default:
				en, err := h.LocateKey(keys[iter])
				if err != nil {
					return err
				}
				if en.ID() != n.ID() {
					res[en] = append(res[en], keys[iter])
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

func (h *Router) copies() error {
	ks := make(map[string]map[string]string)
	ns, err := h.GetNodes()
	if err != nil {
		return err
	}
	for _, n := range ns {
		keys, err := n.Explore()
		if err != nil {
			return err
		}
		for _, key := range keys {
			if ks[key] == nil {
				ks[key] = make(map[string]string)
			}
			en, err := h.LocateKey(key)
			if err != nil {
				return err
			}
			ks[key][n.ID()] = en.ID()
		}
	}
	for key := range ks {
		if len(ks[key]) > 1 {
			log.Println("COPY ", len(ns), len(ks[key]), key, ks[key])
		}
	}
	return nil
}
