package kvstore

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/struckoff/kvstore/rpcapi"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func (h *Host) RunServer(addr string) error {
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

func (h *Host) RPCStore(ctx context.Context, in *rpcapi.StoreReq) (*rpcapi.StoreRes, error) {
	if err := h.n.Store(in.Key, in.Value); err != nil {
		return nil, err
	}
	return &rpcapi.StoreRes{}, nil
}
func (h *Host) RPCReceive(ctx context.Context, in *rpcapi.ReceiveReq) (*rpcapi.ReceiveRes, error) {
	b, err := h.n.Receive(in.Key)
	if err != nil {
		return nil, err
	}
	return &rpcapi.ReceiveRes{Key: in.Key, Value: b}, nil
}
func (h *Host) RPCExplore(ctx context.Context, in *rpcapi.ExploreReq) (*rpcapi.ExploreRes, error) {
	keys, err := h.n.Explore()
	if err != nil {
		return nil, err
	}
	return &rpcapi.ExploreRes{Keys: keys}, nil
}
func (h *Host) RPCMeta(ctx context.Context, in *rpcapi.NodeMetaReq) (*rpcapi.NodeMeta, error) {
	meta := &rpcapi.NodeMeta{
		ID:         h.n.id,
		Address:    h.n.address,
		RPCAddress: h.n.rpcaddress,
		Power:      h.n.Power().Get(),
		Capacity:   h.n.Capacity().Get(),
	}
	return meta, nil
}
