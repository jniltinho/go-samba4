package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	Version   = "v1.0.0"
	BuildDate = "Unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Go-Samba4 version %s (Build Date: %s)\n", Version, BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
