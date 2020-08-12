package router

//
//import (
//	"encoding/json"
//	"sync"
//)
//
//type SyncMapUint2String struct {
//	mu sync.RWMutex
//	s  map[uint64]string
//}
//
//func (sm *SyncMapUint2String) Put(key uint64, value string) {
//	sm.mu.Lock()
//	sm.s[key] = value
//	sm.mu.Unlock()
//}
//
//func (sm *SyncMapUint2String) Get(key uint64) (string, bool) {
//	sm.mu.RLock()
//	v, ok := sm.s[key]
//	sm.mu.RUnlock()
//	return v, ok
//}
//
//func (sm *SyncMapUint2String) Delete(key uint64) {
//	sm.mu.Lock()
//	delete(sm.s, key)
//	sm.mu.Unlock()
//}
//
//func (sm *SyncMapUint2String) JsonMarshal() ([]byte, error) {
//	sm.mu.RLock()
//	b, err := json.Marshal(sm.s)
//	sm.mu.RUnlock()
//	return b, err
//}
//
//func (sm *SyncMapUint2String) Copy() map[uint64]string {
//	sm.mu.RLock()
//	defer sm.mu.RUnlock()
//	res := make(map[uint64]string, len(sm.s))
//	for k, v := range sm.s {
//		res[k] = v
//	}
//	return res
//}
//
//func NewSyncMapUint2String() *SyncMapUint2String {
//	return &SyncMapUint2String{s: make(map[uint64]string)}
//}
