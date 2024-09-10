package seglog

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSegmentedLog(t *testing.T) {
	dir := t.TempDir()
	sl, err := NewSegmentedLog(dir, 100)
	require.NoError(t, err, "Failed to create SegmentedLog")
	defer sl.Close()

	assert.Equal(t, dir, sl.dir, "Unexpected directory")
	assert.Equal(t, 100, sl.segmentSize, "Unexpected segment size")
	assert.Len(t, sl.segments, 1, "Expected 1 segment")
}

func TestAppendAndRead(t *testing.T) {
	dir := t.TempDir()
	sl, err := NewSegmentedLog(dir, 100)
	require.NoError(t, err, "Failed to create SegmentedLog")
	defer sl.Close()

	testData := []string{
		"First entry",
		"Second entry",
		"Third entry",
	}

	offsets := make([]int64, len(testData))
	for i, data := range testData {
		offset, err := sl.Append([]byte(data))
		require.NoError(t, err, "Failed to append entry")
		offsets[i] = offset
	}

	for i, offset := range offsets {
		entry, err := sl.Read(offset)
		require.NoError(t, err, "Failed to read entry")
		assert.Equal(t, testData[i], string(entry), "Unexpected entry content")
	}
}

func TestSegmentRotation(t *testing.T) {
	dir := t.TempDir()
	sl, err := NewSegmentedLog(dir, 2) // Small segment size to force rotation
	require.NoError(t, err, "Failed to create SegmentedLog")
	defer sl.Close()

	for i := 0; i < 5; i++ {
		_, err := sl.Append([]byte("test entry"))
		require.NoError(t, err, "Failed to append entry")
	}

	segments := sl.GetAllSegmentPaths()
	assert.Len(t, segments, 3, "Expected 3 segments")
}

func TestGetActiveSegmentPath(t *testing.T) {
	dir := t.TempDir()
	sl, err := NewSegmentedLog(dir, 100)
	require.NoError(t, err, "Failed to create SegmentedLog")
	defer sl.Close()

	activePath := sl.GetActiveSegmentPath()
	expectedPath := filepath.Join(dir, "log-0.seg")
	assert.Equal(t, expectedPath, activePath, "Unexpected active segment path")
}

func TestGetAllSegmentPaths(t *testing.T) {
	dir := t.TempDir()
	sl, err := NewSegmentedLog(dir, 2) // Small segment size to force rotation
	require.NoError(t, err, "Failed to create SegmentedLog")
	defer sl.Close()

	for i := 0; i < 5; i++ {
		_, err := sl.Append([]byte("test entry"))
		require.NoError(t, err, "Failed to append entry")
	}

	paths := sl.GetAllSegmentPaths()
	assert.Len(t, paths, 3, "Expected 3 segment paths")

	for i, path := range paths {
		expectedPath := filepath.Join(dir, fmt.Sprintf("log-%d.seg", i*2))
		assert.Equal(t, expectedPath, path, "Unexpected segment path")
	}
}

func TestClose(t *testing.T) {
	dir := t.TempDir()
	sl, err := NewSegmentedLog(dir, 100)
	require.NoError(t, err, "Failed to create SegmentedLog")

	assert.NoError(t, sl.Close(), "Failed to close SegmentedLog")

	// Try to append after closing
	_, err = sl.Append([]byte("test"))
	assert.Error(t, err, "Expected error when appending to closed log")
}
