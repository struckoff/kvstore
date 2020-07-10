package ttl

import "sync"

type ChecksMap struct {
	mu sync.RWMutex
	st map[string]*Check
}

func (cm *ChecksMap) Store(name string, check *Check) {
	cm.mu.Lock()
	cm.st[name] = check
	cm.mu.Unlock()
}

func (cm *ChecksMap) Get(name string) (*Check, bool) {
	cm.mu.RLock()
	check, ok := cm.st[name]
	cm.mu.RUnlock()
	return check, ok
}

func (cm *ChecksMap) Delete(name string) {
	cm.mu.Lock()
	delete(cm.st, name)
	cm.mu.Unlock()
}

func (cm *ChecksMap) Update(name string) bool {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if _, ok := cm.st[name]; ok {
		cm.st[name].Update()
		return true
	}
	return false
}

func NewChecksMap() *ChecksMap {
	return &ChecksMap{st: make(map[string]*Check)}
}
