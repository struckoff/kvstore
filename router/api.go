package router

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/struckoff/kvstore/router/nodes"
	"github.com/struckoff/kvstore/router/rpcapi"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	r.OPTIONS("/config/log/enable", h.EnableLog)
	r.OPTIONS("/config/log/disable", h.DisableLog)
	r.OPTIONS("/optimize", h.CallOptimize)
	r.OPTIONS("/rebuild", h.CallRebuild)
	return r
}

func (h *Router) EnableLog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	msg := "logs enabled"
	log.SetOutput(os.Stdout)
	log.Println(msg)
	if _, err := w.Write([]byte(msg)); err != nil {
		log.Println(err)
	}
}

func (h *Router) DisableLog(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	msg := "logs disabled"
	log.Println(msg)
	log.SetOutput(ioutil.Discard)
	if _, err := w.Write([]byte(msg)); err != nil {
		log.Println(err)
	}
}

func (h *Router) CallOptimize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := h.Optimize(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
	}

	if _, err := w.Write([]byte("Optimize complete")); err != nil {
		log.Println(err)
	}
}

func (h *Router) Config(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if err := json.NewEncoder(w).Encode(h.conf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Store value for a given key on the remote node
func (h *Router) Store(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	n, err := h.AddData(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := n.Store(key, b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	uploaded = true
}

//Receive value for a given key from the remote node
func (h *Router) Receive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body != nil {
		defer r.Body.Close()
	}
	k := ps.ByName("key")
	keys := strings.Split(k[1:], "/")
	nmk, err := h.keysOnNodes(keys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	kvsCh := make(chan *rpcapi.KeyValues, len(nmk))
	for n, keys := range nmk {
		go func(n nodes.Node, keys []string, kvsCh chan<- *rpcapi.KeyValues) {
			var kvs *rpcapi.KeyValues
			defer func() {
				kvsCh <- kvs
			}()
			kvs, err = n.Receive(keys)
			if err != nil {
				log.Print(err)
				return
			}
		}(n, keys, kvsCh)
	}

	var resp rpcapi.KeyValues
	resp.KVs = make([]*rpcapi.KeyValue, 0)
	for iter := 0; iter < len(nmk); iter++ {
		kvs := <-kvsCh
		if kvs == nil {
			continue
		}
		for _, kv := range kvs.KVs {
			if kv.Found {
				resp.KVs = append(resp.KVs, kv)
			}
		}
	}

	if err := json.NewEncoder(w).Encode(resp.KVs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func byteSlice2String(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func (h *Router) keysOnNodes(keys []string) (map[nodes.Node][]string, error) {
	nmk := make(map[nodes.Node][]string)
	for iter := range keys {
		n, err := h.LocateKey(keys[iter])
		if err != nil {
			return nil, err
		}
		nmk[n] = append(nmk[n], keys[iter])
	}
	return nmk, nil
}

//Explore returns a list of keys on nodes
func (h *Router) Explore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	h.opLock.Lock()
	defer h.opLock.Unlock()
	res, err := h.nodeKeys()
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

func (h *Router) cidToKeys() (*SyncMap, error) {
	sm := NewSyncMap()
	ns, err := h.bal.Nodes()
	if err != nil {
		return nil, err
	}

	eg, ectx := errgroup.WithContext(context.Background())
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
		var keys []string
		select {
		case <-ctx.Done():
			return nil
		default:
			keys, err = n.Explore()
			if err != nil {
				return err
			}
			select {
			case <-ctx.Done():
				return nil
			default:
				for iter := range keys {
					select {
					case <-ctx.Done():
						return nil
					default:
						di, err := h.ndf(keys[iter])
						if err != nil {
							return err
						}
						_, cid, err := h.bal.LocateData(di)
						if err != nil {
							return err
						}
						k := fmt.Sprintf("%d", cid)
						sm.Append(k, keys[iter])
					}
				}
			}
		}
		return nil
	}
}

func (h *Router) Cid(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	sm, err := h.cidToKeys()
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
func (h *Router) Nodes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	metas, err := h.nodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(metas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Router) nodes() ([]*rpcapi.NodeMeta, error) {
	ns, err := h.bal.Nodes()
	if err != nil {
		return nil, err
	}
	metas := make([]*rpcapi.NodeMeta, len(ns))
	for iter, n := range ns {
		metas[iter] = n.Meta()
	}
	return metas, nil
}

func (h *Router) nodeKeys() (*SyncMap, error) {
	var wg sync.WaitGroup
	res := NewSyncMap()
	ns, err := h.bal.Nodes()
	if err != nil {
		return nil, err
	}
	for _, n := range ns {
		wg.Add(1)
		go func(wg *sync.WaitGroup, n nodes.Node, sm *SyncMap) {
			defer wg.Done()
			keys, err := n.Explore()
			if err != nil {
				log.Printf("%s: %s", n.ID(), err.Error())
				return
			}
			sm.Put(n.ID(), keys)
		}(&wg, n, res)
	}
	wg.Wait()
	return res, nil
}

func (h *Router) CallRebuild(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ns, err := h.bal.Nodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.bal.SetNodes(nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, n := range ns {
		if err := h.bal.AddNode(n); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	ns, err = h.bal.Nodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, n := range ns {
		c, _ := n.Capacity().Get()
		if _, err := fmt.Fprintf(w, "node: %s, cap: %d", n.ID(), c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *Router) Optimize() error {
	h.opLock.Lock()
	defer h.opLock.Unlock()
	return h.optimize()
}

func (h *Router) optimize() error {
	log.Println("Optimize started")
	if err := h.fillBalancer(); err != nil {
		return err
	}
	if err := h.bal.Optimize(); err != nil {
		return err
	}
	if err := h.redistributeKeys(); err != nil {
		return err
	}
	log.Println("Optimize complete")
	return nil
}
