package lsm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemtable(t *testing.T) {
	m := NewMemtable()
	assert.NotNil(t, m)
	assert.NotNil(t, m.data)
	assert.Equal(t, 0, m.size)
}

func TestMemtableSetAndGet(t *testing.T) {
	m := NewMemtable()

	m.Set("key1", "value1")
	m.Set("key2", "value2")

	value, ok := m.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", value)

	value, ok = m.Get("key2")
	assert.True(t, ok)
	assert.Equal(t, "value2", value)

	value, ok = m.Get("nonexistent")
	assert.False(t, ok)
	assert.Empty(t, value)
}

func TestMemtableSize(t *testing.T) {
	m := NewMemtable()

	m.Set("key1", "value1")
	assert.Equal(t, len("key1")+len("value1"), m.Size())

	m.Set("key2", "value2")
	assert.Equal(t, len("key1")+len("value1")+len("key2")+len("value2"), m.Size())

	// Update existing key
	m.Set("key1", "newvalue1")
	assert.Equal(t, len("key2")+len("value2")+len("key1")+len("newvalue1"), m.Size())
}
