package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joobisb/vitadb/internal/config"
	"github.com/joobisb/vitadb/internal/wal"
)

type KVStore struct {
	mu   sync.RWMutex
	data map[string]string
	wal  *wal.WAL
}

func NewKVStore(cfg *config.Config) (*KVStore, error) {
	w, err := wal.NewWAL(cfg.WALDir)
	if err != nil {
		return nil, err
	}

	return &KVStore{
		data: make(map[string]string),
		wal:  w,
	}, nil
}

func (s *KVStore) Set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.wal.AppendSet(key, value); err != nil {
		return err
	}

	s.data[key] = value
	return nil
}

func (s *KVStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

func (s *KVStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.wal.AppendDelete(key); err != nil {
		return err
	}

	delete(s.data, key)
	return nil
}

func (s *KVStore) Close() error {
	return s.wal.Close()
}

func (s *KVStore) RecoverFromWAL() error {
	file, err := os.Open(s.wal.GetWALFilePath())
	if err != nil {
		return fmt.Errorf("failed to open WAL file: %v", err)
	}

	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("error closing file %v", err)
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry wal.LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return fmt.Errorf("failed to unmarshal log entry: %v", err)
		}

		switch entry.Operation {
		case wal.OperationSet:
			s.data[entry.Key] = entry.Value
		case wal.OperationDel:
			delete(s.data, entry.Key)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading WAL file: %v", err)
	}

	return nil
}
