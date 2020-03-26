package main

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type ExternalNode struct {
	id      string
	address string
	p       Power
	c       Capacity
}

func (n *ExternalNode) Power() Power       { return n.p }
func (n *ExternalNode) Capacity() Capacity { return n.c }
func (n *ExternalNode) Store(key string, body io.Reader) error {
	p := strings.Join([]string{addr, key}, "/")
	r, err := http.Post(p, "application/text", body)
	if err != nil {
		return err
	}
	if r.StatusCode >= 400 {
		return errors.New(r.Status)
	}
	return nil
}
func (n *ExternalNode) Receive(key string) ([]byte, error) {
	p := strings.Join([]string{addr, key}, "/")
	r, err := http.Get(p)
	if err != nil {
		return nil, err
	}
	if r.StatusCode >= 400 {
		return nil, errors.New(r.Status)
	}
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	return b, err
}
func (n *ExternalNode) Explore() error { return nil }

func NewExternalNode(id, address string, p float64, c float64) *ExternalNode {
	return &ExternalNode{
		id,
		address,
		NewPower(p),
		NewCapacity(c),
	}
}
