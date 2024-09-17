package store

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/joobisb/vitadb/internal/config"
	"github.com/joobisb/vitadb/internal/lsm"
	"github.com/joobisb/vitadb/internal/wal"
)

type KVStore struct {
	mu   sync.RWMutex
	data map[string]string
	wal  *wal.WAL
	lsm  *lsm.LSM
}

func NewKVStore(cfg *config.Config) (*KVStore, error) {
	w, err := wal.NewWAL(cfg)
	if err != nil {
		return nil, err
	}

	l, err := lsm.NewLSM(cfg)
	if err != nil {
		return nil, err
	}

	return &KVStore{
		data: make(map[string]string),
		wal:  w,
		lsm:  l,
	}, nil
}

func (s *KVStore) Set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.wal.AppendSet(key, value); err != nil {
		return err
	}

	if err := s.lsm.Set(key, value); err != nil {
		return err
	}

	//TODO remove this once we have a proper LSM implementation
	s.data[key] = value
	return nil
}

func (s *KVStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]

	//TODO: Implement once we have LSM and SST
	//TODO: s.lsm.Get(key)

	return value, ok
}

func (s *KVStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.wal.AppendDelete(key); err != nil {
		return err
	}

	//TODO: Implement once we have LSM and SST
	//TODO: s.lsm.Get(key)

	delete(s.data, key)
	return nil
}

func (s *KVStore) Close() error {
	return s.wal.Close()
}

func (s *KVStore) RecoverFromWAL() error {
	paths := s.wal.GetAllSegmentPaths()

	for _, path := range paths {
		if err := s.recoverFromFile(path); err != nil {
			return err
		}
	}

	return nil
}

func (s *KVStore) recoverFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open WAL file %s: %v", filePath, err)
	}
	defer file.Close()

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
		return fmt.Errorf("error reading WAL file %s: %v", filePath, err)
	}

	return nil
}
