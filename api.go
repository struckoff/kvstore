package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func RunServer(addr string) error {
	r := httprouter.New()
	r.POST("/:key", Store)
	r.GET("/:key", Receive)
	r.GET("/", Explore)

	log.Printf("Run server [%s]", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		return err
	}
	return nil
}

func Store(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName("key")
	n := GetNode(key)

	defer r.Body.Close()
	if err := n.Store(key, r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := fmt.Fprint(w, "OK"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Receive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName("key")
	n := GetNode(key)

	defer r.Body.Close()
	body, err := n.Receive(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Explore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, err := fmt.Fprint(w, "Explore"); err != nil {
		panic(err)
	}
}
