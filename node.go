package main

import (
	"io"
)

type Node interface {
	Power() Power
	Capacity() Capacity
	Store(key string, body io.Reader) error
	Receive(key string) ([]byte, error)
	Explore() ([]byte, error)
}
