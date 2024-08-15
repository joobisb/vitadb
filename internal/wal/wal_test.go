package wal

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWAL(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "wal_test")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tempDir) // Clean up the directory after the test

	wal, err := NewWAL(tempDir)
	require.NoError(t, err, "Failed to create WAL")

	t.Run("AppendSet", func(t *testing.T) {
		err := wal.AppendSet("key1", "value1")
		assert.NoError(t, err, "AppendSet failed")
	})

	t.Run("AppendDelete", func(t *testing.T) {
		err := wal.AppendDelete("key2")
		assert.NoError(t, err, "AppendDelete failed")
	})

	err = wal.Close()
	assert.NoError(t, err, "Failed to close WAL")

	walFilePath := filepath.Join(tempDir, walName)
	content, err := os.ReadFile(walFilePath)
	require.NoError(t, err, "Failed to read WAL file")

	lines := splitLines(content)
	assert.Len(t, lines, 2, "Expected 2 lines in WAL")

	var entry LogEntry
	err = json.Unmarshal(lines[0], &entry)
	assert.NoError(t, err, "Failed to unmarshal first entry")
	assert.Equal(t, OperationSet, entry.Operation, "Unexpected operation in first entry")
	assert.Equal(t, "key1", entry.Key, "Unexpected key in first entry")
	assert.Equal(t, "value1", entry.Value, "Unexpected value in first entry")

	err = json.Unmarshal(lines[1], &entry)
	assert.NoError(t, err, "Failed to unmarshal second entry")
	assert.Equal(t, OperationDel, entry.Operation, "Unexpected operation in second entry")
	assert.Equal(t, "key2", entry.Key, "Unexpected key in second entry")
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	for len(data) > 0 {
		i := 0
		for i < len(data) && data[i] != '\n' {
			i++
		}
		lines = append(lines, data[:i])
		data = data[i+1:]
	}
	return lines
}
