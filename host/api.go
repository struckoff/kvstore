package host

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"github.com/struckoff/kvstore/node"
	"github.com/struckoff/kvstore/proto"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func (h *Host) RunServer(addr string) error {
	r := httprouter.New()
	r.POST("/node", h.HTTPRegister)
	r.GET("/node", h.HTTPNodes)
	r.POST("/kv/:key", h.HTTPStore)
	r.GET("/kv/:key", h.HTTPReceive)
	r.GET("/kv", h.HTTPExplore)

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
	nbs := h.bal.Nodes()
	for _, nb := range nbs {
		n, ok := nb.(node.Node)
		if !ok {
			http.Error(w, "wrong node type", http.StatusInternalServerError)
			return
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, n node.Node, sm *SyncMap) {
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
func (h *Host) HTTPRegister(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body == nil {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var meta node.NodeMeta
	if err := json.NewDecoder(r.Body).Decode(&meta); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	en, err := node.NewExternalNode(meta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	metas, err := h.nodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := h.AddNode(en); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(metas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
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

func (h *Host) nodes() ([]node.NodeMeta, error) {
	nbs := h.bal.Nodes()
	metas := make([]node.NodeMeta, len(nbs))
	for iter, nb := range nbs {
		n, ok := nb.(node.Node)
		if !ok {
			return nil, errors.New("wrong node type")
		}
		metas[iter] = n.Meta()
	}
	return metas, nil
}

func (h *Host) RPCStore(ctx context.Context, in *proto.StoreReq) (*proto.StoreRes, error) {
	if err := h.n.Store(in.Key, in.Value); err != nil {
		return nil, err
	}
	return &proto.StoreRes{}, nil
}
func (h *Host) RPCReceive(ctx context.Context, in *proto.ReceiveReq) (*proto.ReceiveRes, error) {
	b, err := h.n.Receive(in.Key)
	if err != nil {
		return nil, err
	}
	return &proto.ReceiveRes{Key: in.Key, Value: b}, nil
}
func (h *Host) RPCExplore(ctx context.Context, in *proto.ExploreReq) (*proto.ExploreRes, error) {
	keys, err := h.n.Explore()
	if err != nil {
		return nil, err
	}
	return &proto.ExploreRes{Keys: keys}, nil
}
func (h *Host) RPCRegister(ctx context.Context, in *proto.NodeMeta) (*proto.NodeMetas, error) {
	return nil, nil
}
