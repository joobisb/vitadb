package main

import (
	"log"

	"github.com/joobisb/vitadb/internal/cli"
	"github.com/joobisb/vitadb/internal/config"
	"github.com/joobisb/vitadb/internal/store"
)

func main() {

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	kvStore, err := store.NewKVStore(cfg)
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

	if err := kvStore.RecoverFromWAL(); err != nil {
		log.Printf("Failed to recover from WAL: %v", err)
	}

	cli := cli.NewCLI(kvStore)
	cli.Run()
}
