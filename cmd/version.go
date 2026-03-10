package cmd

import (
	"fmt"
	"go-samba4/internal/buildinfo"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Go-Samba4 version %s (Build Date: %s)\n", buildinfo.Version, buildinfo.BuildDate)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
