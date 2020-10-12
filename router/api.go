package router

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/struckoff/kvstore/logger"
	"github.com/struckoff/kvstore/router/nodes"
	"github.com/struckoff/kvstore/router/rpcapi"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"unsafe"
)

func (h *Router) HTTPHandler() *httprouter.Router {
	r := httprouter.New()
	//h.POST("/node", h.HTTPRegister)
	r.GET("/nodes", h.Nodes)
	r.POST("/put/:key", h.Store)
	r.GET("/get/*key", h.Receive)
	r.GET("/list", h.Explore)
	r.GET("/cid", h.Cid)
	r.GET("/config", h.Config)
	//r.OPTIONS("/config/log/enable", h.EnableLog)
	//r.OPTIONS("/config/log/disable", h.DisableLog)
	r.OPTIONS("/optimize", h.CallOptimize)
	return r
}

//func (h *Router) EnableLog(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
//	msg := "logs enabled"
//	log.SetOutput(os.Stdout)
//	log.Println(msg)
//	if _, err := w.Write([]byte(msg)); err != nil {
//		log.Println(err)
//	}
//}
//
//func (h *Router) DisableLog(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
//	msg := "logs disabled"
//	log.Println(msg)
//	log.SetOutput(ioutil.Discard)
//	if _, err := w.Write([]byte(msg)); err != nil {
//		log.Println(err)
//	}
//}

func (h *Router) CallOptimize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	if err := h.Optimize(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Logger().Error("optimizer error", zap.Error(err))
	}

	if _, err := w.Write([]byte("Optimize complete")); err != nil {
		logger.Logger().Error("optimizer error", zap.Error(err))
	}
}

func (h *Router) Config(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	if err := json.NewEncoder(w).Encode(h.conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Store value for a given key on the remote node
func (h *Router) Store(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	if r.Body != nil {
		defer r.Body.Close()
	}
	key := ps.ByName("key")

	var uploaded = false
	defer func() {
		if !uploaded {
			if err := h.RemoveData(key); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		if _, err := fmt.Fprint(w, "OK"); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}()

	di, err := h.ndf(key, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	n, cID, err := h.bal.LocateData(di)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kv := &rpcapi.KeyValue{
		Key:   di.RPCApi(),
		Value: b,
	}
	rdi, err := n.Store(ctx, kv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	di = h.rpcndf(rdi)
	if err := h.AddData(cID, di); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uploaded = true
}

//Receive value for a given key from the remote node
func (h *Router) Receive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	if r.Body != nil {
		defer r.Body.Close()
	}
	k := ps.ByName("key")
	keys := strings.Split(k[1:], "/")
	nmk, err := h.keysOnNodes(ctx, keys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	kvsCh := make(chan *rpcapi.KeyValues, len(nmk))
	for n, dis := range nmk {
		go func(ctx context.Context, n nodes.Node, dis []*rpcapi.DataItem, kvsCh chan<- *rpcapi.KeyValues) {
			var kvs *rpcapi.KeyValues
			defer func() {
				kvsCh <- kvs
			}()
			kvs, err = n.Receive(ctx, dis)
			if err != nil {
				logger.Logger().Error("recieve error", zap.Error(err))
				return
			}
		}(ctx, n, dis, kvsCh)
	}
	type rec struct {
		Key   string
		Value string
	}

	//var resp rpcapi.KeyValues
	//resp.KVs = make([]*rpcapi.KeyValue, 0)
	resp := make([]rec, 0)
	for i := 0; i < len(nmk); i++ {
		select {
		case <-ctx.Done():
			return
		case kvs := <-kvsCh:
			if kvs == nil {
				continue
			}
			for _, kv := range kvs.KVs {
				if kv.Found {
					resp = append(resp, rec{string(kv.Key.ID), string(kv.Value)})
				}
			}
		}
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func byteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func (h *Router) keysOnNodes(ctx context.Context, keys []string) (map[nodes.Node][]*rpcapi.DataItem, error) {
	nmk := make(map[nodes.Node][]*rpcapi.DataItem)
	for i := range keys {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			di, err := h.ndf(keys[i], 0)
			if err != nil {
				return nil, err
			}
			n, _, err := h.bal.LocateData(di)
			if err != nil {
				return nil, err
			}
			nmk[n] = append(nmk[n], di.RPCApi())
		}

	}
	return nmk, nil
}

//Explore returns a list of keys on nodes
func (h *Router) Explore(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	//h.opLock.Lock()
	//defer h.opLock.Unlock()
	res, err := h.nodeKeys(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := res.JsonMarshal()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Router) cidToKeys(ctx context.Context) (*SyncMap, error) {
	sm := NewSyncMap()
	ns, err := h.bal.Nodes()
	if err != nil {
		return nil, err
	}

	eg, ectx := errgroup.WithContext(ctx)
	for _, n := range ns {
		eg.Go(h.cidToKeysNode(n, sm, ectx))
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}
	return sm, nil
}

func (h *Router) cidToKeysNode(n nodes.Node, sm *SyncMap, ctx context.Context) func() error {
	return func() (err error) {
		var rdis []*rpcapi.DataItem
		select {
		case <-ctx.Done():
			return nil
		default:
			rdis, err = n.Explore(ctx)
			if err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return nil
			default:
				for i := range rdis {
					select {
					case <-ctx.Done():
						return nil
					default:
						di := h.rpcndf(rdis[i])
						_, cid, err := h.bal.LocateData(di)
						if err != nil {
							return err
						}
						k := fmt.Sprintf("%d", cid)
						sm.Append(k, string(rdis[i].ID))
					}
				}
			}
		}
		return nil
	}
}

func (h *Router) Cid(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	sm, err := h.cidToKeys(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := sm.JsonMarshal()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// Nodes returns a list of nodes
func (h *Router) Nodes(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	metas, err := h.nodes(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(metas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Router) nodes(ctx context.Context) ([]*rpcapi.NodeMeta, error) {
	ns, err := h.bal.Nodes()
	if err != nil {
		return nil, err
	}
	metas := make([]*rpcapi.NodeMeta, len(ns))
	for i, n := range ns {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			metas[i] = n.Meta(ctx)
		}
	}
	return metas, nil
}

func (h *Router) nodeKeys(ctx context.Context) (*SyncMap, error) {
	var wg sync.WaitGroup
	res := NewSyncMap()
	ns, err := h.bal.Nodes()
	if err != nil {
		return nil, err
	}
	for _, n := range ns {
		wg.Add(1)
		go func(ctx context.Context, wg *sync.WaitGroup, n nodes.Node, sm *SyncMap) {
			defer wg.Done()
			dis, err := n.Explore(ctx)
			if err != nil {
				logger.Logger().Error("node keys error", zap.String("Node", n.ID()), zap.Error(err))
				return
			}
			select {
			case <-ctx.Done():
				return
			default:
				keys := make([]string, len(dis))
				for i := range keys {
					keys[i] = string(dis[i].ID)
				}
				sm.Put(n.ID(), keys)
			}
		}(ctx, &wg, n, res)
	}
	wg.Wait()
	return res, nil
}

func (h *Router) Optimize(ctx context.Context) error {
	//h.opLock.Lock()
	//defer h.opLock.Unlock()
	return h.optimize(ctx)
}

func (h *Router) optimize(ctx context.Context) error {
	logger.Logger().Info("Optimize started")
	if err := h.fillBalancer(ctx); err != nil {
		return err
	}
	if err := h.bal.Optimize(); err != nil {
		return err
	}
	if err := h.redistributeKeys(ctx); err != nil {
		return err
	}
	logger.Logger().Info("Optimize complete")
	return nil
}
