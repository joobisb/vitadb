package lsm

import (
	"fmt"

	"github.com/joobisb/vitadb/internal/config"
)

type LSM struct {
	memtable *Memtable
	config   *config.Config

	// keep track of SSTables
	sstables   []*SSTable
	sstCounter int
}

func NewLSM(cfg *config.Config) (*LSM, error) {
	return &LSM{
		memtable:   NewMemtable(),
		config:     cfg,
		sstables:   make([]*SSTable, 0),
		sstCounter: 0,
	}, nil
}

func (l *LSM) Set(key, value string) error {

	// insert into Memtable
	l.memtable.Set(key, value)

	// Check if Memtable needs to be flushed (we'll implement this later)
	//TODO implement the concept of 2 active memtables
	//when flushMemtable is called, we switch the active memtable and write the previous one to disk
	if l.memtable.Size() >= l.config.MemtableSize {
		if err := l.flushMemtable(); err != nil {
			return err
		}
	}

	return nil
}

func (l *LSM) flushMemtable() error {
	ssTable, err := NewSSTable(l.config, l.sstCounter)
	if err != nil {
		return fmt.Errorf("failed to create new SSTable: %v", err)
	}

	// Iterate through the memtable and write entries to the SST file
	err = l.memtable.Iterate(func(key, value string) error {
		if err := ssTable.writeEntry(key, value); err != nil {
			return fmt.Errorf("failed to write entry to SST: %v", err)
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Add the new SST to the list of SSTables
	l.sstables = append(l.sstables, ssTable)
	l.sstCounter++

	// Create a new memtable
	l.memtable = NewMemtable()

	return nil
}

// Add this new method to iterate over the Memtable
func (m *Memtable) Iterate(fn func(key, value string) error) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for e := m.data.Front(); e != nil; e = e.Next() {
		if err := fn(e.Key().(string), e.Value.(string)); err != nil {
			return err
		}
	}
	return nil
}
