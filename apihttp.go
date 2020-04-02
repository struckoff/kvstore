package kvstore

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func (h *Host) RunHTTPServer(addr string) error {
	r := httprouter.New()
	//r.POST("/node", h.HTTPRegister)
	r.GET("/nodes", h.HTTPNodes)
	r.POST("/put/:key", h.HTTPStore)
	r.GET("/get/:key", h.HTTPReceive)
	r.GET("/list", h.HTTPExplore)

	log.Printf("Run server [%s]", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		return err
	}
	return nil
}

func (h *Host) HTTPStore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body != nil {
		defer r.Body.Close()
	}
	key := ps.ByName("key")
	n, err := h.GetNode(key)
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
func (h *Host) HTTPReceive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body != nil {
		defer r.Body.Close()
	}
	key := ps.ByName("key")
	n, err := h.GetNode(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var body []byte
	body, err = n.Receive(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (h *Host) HTTPExplore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var wg sync.WaitGroup
	res := NewSyncMap()
	ns, err := h.bal.Nodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, n := range ns {
		wg.Add(1)
		go func(wg *sync.WaitGroup, n Node, sm *SyncMap) {
			defer wg.Done()
			keys, err := n.Explore()
			if err != nil {
				log.Printf("%s: %w", n.ID(), err)
				return
			}
			sm.Put(n.ID(), keys)
		}(&wg, n, res)
	}
	wg.Wait()
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

//func (h *Host) HTTPRegister(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	if r.Body == nil {
//		http.Error(w, "empty body", http.StatusBadRequest)
//		return
//	}
//	defer r.Body.Close()
//	var meta NodeMeta
//	if err := json.NewDecoder(r.Body).Decode(&meta); err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	en, err := NewExternalNode(meta)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	if err := h.AddNode(en); err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	metas, err := h.nodes()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	if err := json.NewEncoder(w).Encode(metas); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//	}
//}
func (h *Host) HTTPNodes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	metas, err := h.nodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(metas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Host) nodes() ([]NodeMeta, error) {
	ns, err := h.bal.Nodes()
	if err != nil {
		return nil, err
	}
	metas := make([]NodeMeta, len(ns))
	for iter, n := range ns {
		metas[iter] = n.Meta()
	}
	return metas, nil
}
