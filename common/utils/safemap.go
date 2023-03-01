package utils

import "sync"

type SafeMap struct {
	lock   *sync.RWMutex
	status map[interface{}]interface{}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		lock:   new(sync.RWMutex),
		status: make(map[interface{}]interface{}),
	}
}

func (m *SafeMap) Get(k interface{}) interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if val, ok := m.status[k]; ok {
		return val
	}
	return nil

}

func (m *SafeMap) Set(k, v interface{}) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	if val, ok := m.status[k]; !ok {
		m.status[k] = v
	} else if val != nil {
		m.status[k] = v
	} else {
		return false
	}
	return true
}

// Check Returns true if k is exist in the map.

func (m *SafeMap) Check(k interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	_, ok := m.status[k]
	return ok
}

// Delete the given key and value.

func (m *SafeMap) Delete(k interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.status, k)

}

// Items returns all items in safemap.

func (m *SafeMap) Items() map[interface{}]interface{} {
	m.lock.RLock()
	defer m.lock.RUnlock()
	r := make(map[interface{}]interface{})
	for k, v := range m.status {
		r[k] = v
	}
	return r

}

// Count returns the number of items within the map.

func (m *SafeMap) Count() int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.status)

}
