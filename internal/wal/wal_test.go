package wal

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/joobisb/vitadb/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWAL(t *testing.T) {
	t.Run("SingleFileWAL", func(t *testing.T) {
		testWAL(t, false)
	})

	t.Run("SegmentedLogWAL", func(t *testing.T) {
		testWAL(t, true)
	})
}

func testWAL(t *testing.T, useSegmentedLogs bool) {
	tempDir, err := os.MkdirTemp("", "wal_test")
	require.NoError(t, err, "Failed to create temp directory")
	defer os.RemoveAll(tempDir) // Clean up the directory after the test

	cfg := &config.Config{
		WALDir:           tempDir,
		UseSegmentedLogs: useSegmentedLogs,
		SegmentSize:      1024,
	}

	wal, err := NewWAL(cfg)
	require.NoError(t, err, "Failed to create WAL")

	t.Run("AppendSet", func(t *testing.T) {
		err := wal.AppendSet("key1", "value1")
		assert.NoError(t, err, "AppendSet failed")
	})

	t.Run("AppendDelete", func(t *testing.T) {
		err := wal.AppendDelete("key2")
		assert.NoError(t, err, "AppendDelete failed")
	})

	walFilePath := wal.GetWALFilePath()
	assert.NotEmpty(t, walFilePath, "WAL file path should not be empty")

	err = wal.Close()
	assert.NoError(t, err, "Failed to close WAL")

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

	// Test GetAllSegmentPaths
	paths := wal.GetAllSegmentPaths()
	assert.NotEmpty(t, paths, "GetAllSegmentPaths should return at least one path")
	if useSegmentedLogs {
		assert.GreaterOrEqual(t, len(paths), 1, "Segmented log should have at least one segment")
	} else {
		assert.Len(t, paths, 1, "Single file WAL should have exactly one path")
	}
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
