package router

import (
	"log"
	"net/http"
	"time"
)

func LatencyMiddleware(h http.Handler, wait time.Duration) *Middleware {
	prefunc := func(w http.ResponseWriter, r *http.Request) error {
		time.Sleep(wait)
		return nil
	}
	return &Middleware{h: h, prefunc: prefunc, postfunc: nil}
}

type Middleware struct {
	h        http.Handler
	prefunc  func(w http.ResponseWriter, r *http.Request) error
	postfunc func(w http.ResponseWriter, r *http.Request) error
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if m.prefunc != nil {
		if err := m.prefunc(w, r); err != nil {
			log.Print(err.Error())
			return
		}
	}

	if m.h != nil {
		m.h.ServeHTTP(w, r)
	}

	if m.postfunc != nil {
		if err := m.postfunc(w, r); err != nil {
			log.Print(err.Error())
			return
		}
	}
}
