package command

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	DefaultHost = "localhost"
	DefaultPort = "6370"
)

var (
	host string
	port string
)

var rootCmd = &cobra.Command{
	Use:   "vitadb-cli",
	Short: "VitaDB CLI - A command-line interface for VitaDB",
	Long:  `VitaDB CLI is a command-line interface for interacting with the VitaDB server.`,
	Run: func(cmd *cobra.Command, args []string) {
		if host == "" {
			host = DefaultHost
		}
		if port == "" {
			port = DefaultPort
		}

		serverAddress := net.JoinHostPort(host, port)
		conn, err := net.Dial("tcp", serverAddress)
		if err != nil {
			fmt.Printf("Error connecting to VitaDB server at %s: %v\n", serverAddress, err)
			return
		}
		defer conn.Close()

		fmt.Printf("Connected to VitaDB server at %s\n", serverAddress)
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			if !scanner.Scan() {
				break
			}
			command := scanner.Text()
			if command == "exit" {
				break
			}
			fmt.Fprintf(conn, "%s\n", command)
			response, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Println("Error reading response:", err)
				continue
			}
			fmt.Print(strings.TrimSpace(response), "\n")
		}
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "", "VitaDB server host")
	rootCmd.PersistentFlags().StringVarP(&port, "port", "p", "", "VitaDB server port")
	rootCmd.Flags().SortFlags = false
}

func Execute() error {
	return rootCmd.Execute()
}
