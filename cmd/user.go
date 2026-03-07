package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage local administrator bypass users or perform skeleton tasks",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("User management CLI tool.")
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
}
