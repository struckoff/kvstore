package router

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/struckoff/kvstore/router/nodes"
	"github.com/struckoff/kvstore/router/rpcapi"
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
	r.GET("/config", h.Config)
	r.OPTIONS("/config/log/enable", h.EnableLog)
	r.OPTIONS("/config/log/disable", h.DisableLog)
	r.OPTIONS("/optimize", h.Optimize)
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

func (h *Router) Optimize(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log.Println("optimize started")
	if err := h.bal.Optimize(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}
	if err := h.redistributeKeys(); err != nil {
		log.Printf("Error redistributing keys: %s", err.Error())
		return
	}
	if _, err := w.Write([]byte("optimize complete")); err != nil {
		log.Println(err)
	}
	log.Println("optimize complete")
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
	n, err := h.LocateKey(key)
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
	if _, err := fmt.Fprint(w, "OK"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
		func(n nodes.Node, keys []string, kvsCh chan<- *rpcapi.KeyValues) {
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
