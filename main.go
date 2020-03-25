package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

const addr = "0.0.0.0:9191"

func main() {
	r := httprouter.New()
	r.POST("/:key", Store)
	r.GET("/:key", Receive)
	r.GET("/", Explore)

	log.Printf("Run server [%s]", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		panic(err)
	}
}

func Store(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName("key")
	if _, err := fmt.Fprintf(w, "Store: %s", key); err != nil {
		panic(err)
	}
}

func Receive(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	key := ps.ByName("key")
	if _, err := fmt.Fprintf(w, "Receive: %s", key); err != nil {
		panic(err)
	}
}

func Explore(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if _, err := fmt.Fprint(w, "Explore"); err != nil {
		panic(err)
	}

}
