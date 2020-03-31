package kvstore

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

func (h *Host) RunServer(addr string) error {
	r := httprouter.New()
	r.POST("/node", h.Register)
	r.GET("/nodes", h.Nodes)
	r.POST("/put/:key", h.Store)
	r.GET("/get/:key", h.Receive)
	r.GET("/list", h.Explore)

	log.Printf("Run server [%s]", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		return err
	}
	return nil
}

// Store save value for a given key
func (h *Host) Store(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body != nil {
		defer r.Body.Close()
	}
	key := ps.ByName("key")
	n, err := h.GetNode(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := n.Store(key, r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := fmt.Fprint(w, "OK"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

//func (h *Host) StoreOutside(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	if r.Body != nil {
//		defer r.Body.Close()
//	}
//	key := ps.ByName("key")
//	if err := h.n.Store(key, r.Body); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	if _, err := fmt.Fprint(w, "OK"); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//}
func (h *Host) Receive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func (h *Host) Explore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, err := fmt.Fprint(w, "Explore"); err != nil {
		panic(err)
	}
}

func (h *Host) Register(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if r.Body == nil {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var meta NodeMeta
	if err := json.NewDecoder(r.Body).Decode(&meta); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	en := NewExternalNode(meta)
	if err := h.AddNode(en); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	metas, err := h.nodes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(metas); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Host) Nodes(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
	nbs := h.bal.Nodes()
	metas := make([]NodeMeta, len(nbs))
	for iter, nb := range nbs {
		n, ok := nb.(Node)
		if !ok {
			return nil, errors.New("Wrong node type")
		}
		metas[iter] = n.Meta()
	}
	return metas, nil
}
