package ttl

import (
	"github.com/struckoff/kvstore/router/rpcapi"
	"sync"
	"time"
)

type Check struct {
	mu          sync.RWMutex
	deadman     time.Duration // How long to wait unit node will be declared dead
	removeAfter time.Duration // How long to wait unit node will be removed

	timerDead   *time.Timer // Timer that set node dead and start Check.timerRemove
	timerRemove *time.Timer // Timer that remove node
}

func NewTTLCheck(hc *rpcapi.HealthCheck, onDead, onRemove func()) (*Check, error) {
	var err error
	t := &Check{
		deadman:     0,
		removeAfter: 0,
	}
	t.deadman, err = time.ParseDuration(hc.Timeout)
	if err != nil {
		return nil, err
	}
	if len(hc.DeregisterCriticalServiceAfter) > 0 {
		t.removeAfter, err = time.ParseDuration(hc.DeregisterCriticalServiceAfter)
		if err != nil {
			return nil, err
		}
	}
	t.timerDead = time.AfterFunc(t.deadman, t.deadHandler(onDead, onRemove))
	return t, nil
}

func (t *Check) deadHandler(onDead, onRemove func()) func() {
	return func() {
		if t.removeAfter > 0 {
			t.timerRemove = time.AfterFunc(t.removeAfter, onRemove)
		}
		onDead()
	}
}

func (t *Check) Update() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.timerDead.Reset(t.deadman)
	if t.timerRemove != nil {
		t.timerRemove.Stop()
	}
}
