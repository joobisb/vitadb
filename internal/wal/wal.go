package wal

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type OperationType string

const (
	OperationSet OperationType = "SET"
	OperationDel OperationType = "DEL"

	walName string = "kvstore.wal"
)

type LogEntry struct {
	Operation OperationType `json:"op"`
	Key       string        `json:"key"`
	Value     string        `json:"value,omitempty"`
}

type WAL struct {
	filePath string
	file     *os.File
	mu       sync.Mutex
}

func NewWAL(walDir string) (*WAL, error) {
	walFilePath := walDir + "/" + walName

	fmt.Println("walFilePath", walFilePath)
	file, err := os.OpenFile(walFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open WAL file: %v", err)
	}

	return &WAL{
		file:     file,
		filePath: walFilePath,
	}, nil
}

func (w *WAL) GetWALFilePath() string {
	return w.filePath
}

func (w *WAL) AppendSet(key, value string) error {
	return w.append(LogEntry{Operation: OperationSet, Key: key, Value: value})
}

func (w *WAL) AppendDelete(key string) error {
	return w.append(LogEntry{Operation: OperationDel, Key: key})
}

func (w *WAL) append(entry LogEntry) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %v", err)
	}

	if _, err := fmt.Fprintf(w.file, "%s\n", data); err != nil {
		return fmt.Errorf("failed to write log entry: %v", err)
	}

	return nil
}

func (w *WAL) Close() error {
	return w.file.Close()
}
