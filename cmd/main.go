package main

import (
	"github.com/joobisb/patterns/wal/internal/cli"
	"github.com/joobisb/patterns/wal/internal/store"
)

func main() {
	kvStore := store.NewKVStore()
	cli := cli.NewCLI(kvStore)
	cli.Run()
}
