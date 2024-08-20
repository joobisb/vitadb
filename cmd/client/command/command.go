package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of VitaDB CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("VitaDB CLI v0.1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
