package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"

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

	if err := kvStore.RecoverFromWAL(); err != nil {
		log.Printf("Failed to recover from WAL: %v", err)
	}

	listener, err := net.Listen("tcp", ":6370")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	fmt.Println("VitaDB server listening on :6370")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go handleConnection(conn, kvStore)
	}
}

func handleConnection(conn net.Conn, kvStore *store.KVStore) {
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		cmd := strings.Fields(scanner.Text())
		if len(cmd) == 0 {
			continue
		}
		switch strings.ToUpper(cmd[0]) {
		case "SET":
			if len(cmd) != 3 {
				fmt.Fprintf(conn, "Usage: set <key> <value>\n")
				continue
			}
			err := kvStore.Set(cmd[1], cmd[2])
			if err != nil {
				fmt.Fprintf(conn, "ERR %v\n", err)
			} else {
				fmt.Fprintf(conn, "OK\n")
			}
		case "GET":
			if len(cmd) != 2 {
				fmt.Fprintf(conn, "Usage: get <key>")
				continue
			}
			value, ok := kvStore.Get(cmd[1])
			if !ok {
				fmt.Fprintf(conn, "(nil)\n")
			} else {
				fmt.Fprintf(conn, "%s\n", value)
			}
		case "DEL":
			if len(cmd) != 2 {
				fmt.Fprintf(conn, "ERR wrong number of arguments for 'del' command\n")
				continue
			}
			err := kvStore.Delete(cmd[1])
			if err != nil {
				fmt.Fprintf(conn, "ERR %v\n", err)
			} else {
				fmt.Fprintf(conn, "OK\n")
			}
		default:
			fmt.Fprintf(conn, "ERR unknown command '%s'\n", cmd[0])
		}
	}
}
