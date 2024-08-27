package wal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/joobisb/vitadb/internal/config"
	"github.com/joobisb/vitadb/internal/seglog"
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
	mu              sync.Mutex
	useSegmentedLog bool
	singleLog       *os.File
	segmentedLog    *seglog.SegmentedLog
}

func NewWAL(cfg *config.Config) (*WAL, error) {
	// Ensure the WAL directory exists
	if err := os.MkdirAll(cfg.WALDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create WAL directory: %v", err)
	}
	if cfg.UseSegmentedLogs {
		log, err := seglog.NewSegmentedLog(cfg.WALDir, cfg.SegmentSize)
		if err != nil {
			return nil, fmt.Errorf("failed to create segmented log: %v", err)
		}
		return &WAL{
			useSegmentedLog: true,
			segmentedLog:    log,
		}, nil
	}
	// Existing single file implementation
	file, err := os.OpenFile(filepath.Join(cfg.WALDir, "wal.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open WAL file: %v", err)
	}
	return &WAL{
		useSegmentedLog: false,
		singleLog:       file,
	}, nil
}

func (w *WAL) GetWALFilePath() string {
	if w.useSegmentedLog {
		return w.segmentedLog.GetActiveSegmentPath()
	}
	return w.singleLog.Name()
}

func (w *WAL) GetAllSegmentPaths() []string {
	if w.useSegmentedLog {
		return w.segmentedLog.GetAllSegmentPaths()
	}
	return []string{w.singleLog.Name()}
}

func (w *WAL) AppendSet(key, value string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	entry := LogEntry{Operation: OperationSet, Key: key, Value: value}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %v", err)
	}

	if w.useSegmentedLog {
		_, err = w.segmentedLog.Append(data)
	} else {
		_, err = fmt.Fprintf(w.singleLog, "%s\n", data)
	}

	return err
}

func (w *WAL) AppendDelete(key string) error {
	// Similar to AppendSet, but for delete operation
	w.mu.Lock()
	defer w.mu.Unlock()

	entry := LogEntry{Operation: OperationDel, Key: key}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal log entry: %v", err)
	}

	if w.useSegmentedLog {
		_, err = w.segmentedLog.Append(data)
	} else {
		_, err = fmt.Fprintf(w.singleLog, "%s\n", data)
	}

	return err
}

func (w *WAL) Close() error {
	if w.useSegmentedLog {
		return w.segmentedLog.Close()
	}
	return w.singleLog.Close()
}
