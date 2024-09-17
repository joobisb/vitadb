package lsm

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/joobisb/vitadb/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNewSSTable(t *testing.T) {
	tempDir := t.TempDir()
	cfg := &config.Config{SSTDir: tempDir}

	sst, err := NewSSTable(cfg, 1)
	assert.NoError(t, err)
	assert.NotNil(t, sst)
	assert.Equal(t, filepath.Join(tempDir, "sst_1.db"), sst.path)

	// Check if the file was created
	_, err = os.Stat(sst.path)
	assert.NoError(t, err)
}

func TestSSTableWriteEntry(t *testing.T) {
	tempDir := t.TempDir()
	cfg := &config.Config{SSTDir: tempDir}

	sst, _ := NewSSTable(cfg, 1)

	err := sst.writeEntry("key1", "value1")
	assert.NoError(t, err)

	err = sst.writeEntry("key2", "value2")
	assert.NoError(t, err)

	// Read the file contents and verify
	data, err := os.ReadFile(sst.path)
	assert.NoError(t, err)

	expected := []byte{
		4, 0, 0, 0, // key1 length
		'k', 'e', 'y', '1',
		6, 0, 0, 0, // value1 length
		'v', 'a', 'l', 'u', 'e', '1',
		4, 0, 0, 0, // key2 length
		'k', 'e', 'y', '2',
		6, 0, 0, 0, // value2 length
		'v', 'a', 'l', 'u', 'e', '2',
	}

	assert.Equal(t, expected, data)
}
