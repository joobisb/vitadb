package store

import (
	"os"
	"testing"

	"github.com/joobisb/vitadb/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKVStore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "kvstore_test")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tempDir)

	cfg := &config.Config{
		WALDir: tempDir,
	}

	store, err := NewKVStore(cfg)
	require.NoError(t, err, "Failed to create KVStore")

	t.Run("Set and Get", func(t *testing.T) {
		err := store.Set("key1", "value1")
		assert.NoError(t, err, "Set failed")

		value, ok := store.Get("key1")
		assert.True(t, ok, "Get failed: key not found")
		assert.Equal(t, "value1", value, "Unexpected value")
	})

	t.Run("Delete", func(t *testing.T) {
		err := store.Set("key2", "value2")
		assert.NoError(t, err, "Set failed")

		err = store.Delete("key2")
		assert.NoError(t, err, "Delete failed")

		_, ok := store.Get("key2")
		assert.False(t, ok, "Key should have been deleted")
	})

	t.Run("RecoverFromWAL", func(t *testing.T) {
		err = store.Close()
		assert.NoError(t, err, "Failed to close store")

		newStore, err := NewKVStore(cfg)
		require.NoError(t, err, "Failed to create new KVStore")

		err = newStore.RecoverFromWAL()
		assert.NoError(t, err, "Failed to recover from WAL")

		value, ok := newStore.Get("key1")
		assert.True(t, ok, "Recovery failed: key1 not found")
		assert.Equal(t, "value1", value, "Recovery failed: unexpected value")

		_, ok = newStore.Get("key2")
		assert.False(t, ok, "Recovery failed: key2 should have been deleted")
	})
}
