package store

import (
	"context"
	"github.com/struckoff/kvstore/router/rpcapi"
	"sync"
)

func NewCapacity(c float64) Capacity {
	return Capacity{c: c}
}

type Capacity struct {
	l sync.RWMutex
	c float64
}

func (c *Capacity) Get() (float64, error) {
	c.l.RLock()
	defer c.l.RUnlock()
	return c.c, nil
}

func (c *Capacity) Add(arg float64) {
	//c.l.Lock()
	//c.c += arg
	//c.l.Unlock()
}

func (c *Capacity) RPCGet(_ context.Context, _ *rpcapi.Empty) (*rpcapi.Capacity, error) {
	c.l.RLock()
	defer c.l.RUnlock()
	return &rpcapi.Capacity{Capacity: c.c}, nil
}
