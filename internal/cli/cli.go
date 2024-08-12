package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joobisb/patterns/wal/internal/store"
)

type CLI struct {
	store *store.KVStore
}

func NewCLI(store *store.KVStore) *CLI {
	return &CLI{store: store}
}

func (c *CLI) Run() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		parts := strings.Fields(input)
		if len(parts) == 0 {
			continue
		}

		switch parts[0] {
		case "set":
			if len(parts) != 3 {
				fmt.Println("Usage: set <key> <value>")
				continue
			}
			c.store.Set(parts[1], parts[2])
			fmt.Println("OK")
		case "get":
			if len(parts) != 2 {
				fmt.Println("Usage: get <key>")
				continue
			}
			value, ok := c.store.Get(parts[1])
			if !ok {
				fmt.Println("Key not found")
			} else {
				fmt.Println(value)
			}
		case "exit":
			return
		default:
			fmt.Println("Unknown command. Available commands: set, get, exit")
		}
	}
}
