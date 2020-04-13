package router

import (
	"encoding/json"
	"sync"
)

type SyncMap struct {
	mu sync.RWMutex
	s  map[string][]string
}

func (sm *SyncMap) Put(key string, value []string) {
	sm.mu.Lock()
	sm.s[key] = value
	sm.mu.Unlock()
}

func (sm *SyncMap) Get(key string) ([]string, bool) {
	sm.mu.RLock()
	v, ok := sm.s[key]
	sm.mu.RUnlock()
	return v, ok
}

func (sm *SyncMap) Delete(key string) {
	sm.mu.Lock()
	delete(sm.s, key)
	sm.mu.Unlock()
}

func (sm *SyncMap) JsonMarshal() ([]byte, error) {
	sm.mu.RLock()
	b, err := json.Marshal(sm.s)
	sm.mu.RUnlock()
	return b, err
}

func (sm *SyncMap) Copy() map[string][]string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	res := make(map[string][]string, len(sm.s))
	for k, v := range sm.s {
		res[k] = v
	}
	return res
}

func NewSyncMap() *SyncMap {
	return &SyncMap{s: make(map[string][]string)}
}
