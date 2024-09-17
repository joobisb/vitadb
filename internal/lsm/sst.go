package lsm

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joobisb/vitadb/internal/config"
)

type SSTable struct {
	path string
}

func NewSSTable(cfg *config.Config, id int) (*SSTable, error) {
	// Ensure the SST directory exists
	if err := os.MkdirAll(cfg.SSTDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create WAL directory: %v", err)
	}

	path := filepath.Join(cfg.SSTDir, fmt.Sprintf("sst_%d.db", id))
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create SST file: %v", err)
	}
	defer file.Close()

	return &SSTable{path: path}, nil
}

// SST file format:
// [key_size (4 bytes)][key][value_size (4 bytes)][value]...
func (sst *SSTable) writeEntry(key, value string) error {
	file, err := os.OpenFile(sst.path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open SST file: %v", err)
	}
	defer file.Close()

	// Write key size
	//This approach allows for keys and values of varying lengths, up to a maximum of 2^32 - 1 bytes.
	if err := binary.Write(file, binary.LittleEndian, uint32(len(key))); err != nil {
		return fmt.Errorf("failed to write key size: %v", err)
	}

	// Write key
	if _, err := file.Write([]byte(key)); err != nil {
		return fmt.Errorf("failed to write key: %v", err)
	}

	// Write value size
	if err := binary.Write(file, binary.LittleEndian, uint32(len(value))); err != nil {
		return fmt.Errorf("failed to write value size: %v", err)
	}

	// Write value
	if _, err := file.Write([]byte(value)); err != nil {
		return fmt.Errorf("failed to write value: %v", err)
	}

	return nil
}
