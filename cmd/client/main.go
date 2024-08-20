package main

import (
	"fmt"
	"os"

	"github.com/joobisb/vitadb/cmd/client/command"
)

func main() {
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
