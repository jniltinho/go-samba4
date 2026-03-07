package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"go-samba4/internal/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("Starting Samba4-Admin Server")
		server.Serve(globalCfg, tplFS, statFS)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntP("port", "p", 8080, "Port to run the server on")
}
