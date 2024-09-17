package lsm

import (
	"fmt"
	"testing"

	"github.com/joobisb/vitadb/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewLSM(t *testing.T) {
	cfg := &config.Config{MemtableSize: 1024}
	lsm, err := NewLSM(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, lsm)
	assert.Equal(t, cfg, lsm.config)
	assert.NotNil(t, lsm.memtable)
	assert.Empty(t, lsm.sstables)
	assert.Equal(t, 0, lsm.sstCounter)
}

func TestLSMSet(t *testing.T) {
	cfg := &config.Config{MemtableSize: 100, SSTDir: t.TempDir()}
	lsm, _ := NewLSM(cfg)

	err := lsm.Set("key1", "value1")
	assert.NoError(t, err)
	assert.Equal(t, 1, lsm.memtable.data.Len())

	// Test flushing memtable
	for i := 0; i < 10; i++ {
		err := lsm.Set(fmt.Sprintf("key%d", i), "long_value_to_trigger_flush")
		assert.NoError(t, err)
	}

	assert.Equal(t, 2, len(lsm.sstables))
	assert.Equal(t, 2, lsm.sstCounter)
}

func TestLSMIterate(t *testing.T) {
	m := NewMemtable()
	m.Set("key1", "value1")
	m.Set("key2", "value2")

	count := 0
	err := m.Iterate(func(key, value string) error {
		count++
		assert.Contains(t, []string{"key1", "key2"}, key)
		assert.Contains(t, []string{"value1", "value2"}, value)
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, 2, count)
}
