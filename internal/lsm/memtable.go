package lsm

import (
	"sync"

	"github.com/huandu/skiplist"
)

type Memtable struct {
	data *skiplist.SkipList
	size int
	mu   sync.RWMutex
}

func NewMemtable() *Memtable {

	return &Memtable{
		data: skiplist.New(skiplist.String),
		size: 0,
	}
}

func (m *Memtable) Set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	oldValue := m.data.Get(key)
	if oldValue != nil {
		m.size -= len(oldValue.Value.(string))
		m.size += len(value)
	} else {
		m.size += len(key) + len(value)
	}
	m.data.Set(key, value)

}

func (m *Memtable) Get(key string) (string, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	element := m.data.Get(key)
	if element == nil {
		return "", false
	}

	value, ok := element.Value.(string)
	if !ok {
		return "", false
	}
	return value, true
}

func (m *Memtable) Size() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.size
}
