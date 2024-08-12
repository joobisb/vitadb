package main

import (
	"log"

	"github.com/joobisb/patterns/wal/internal/cli"
	"github.com/joobisb/patterns/wal/internal/store"
)

func main() {
	walFile := "kvstore.wal"
	kvStore, err := store.NewKVStore(walFile)
	if err != nil {
		log.Fatalf("Failed to create KVStore: %v", err)
	}
	defer func() {
		err := kvStore.Close()
		if err != nil {
			log.Printf("error closing store %v", err)
			return
		}
	}()

	if err := kvStore.RecoverFromWAL(walFile); err != nil {
		log.Printf("Failed to recover from WAL: %v", err)
	}

	cli := cli.NewCLI(kvStore)
	cli.Run()
}
