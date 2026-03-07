package main

import (
	"log/slog"
	"os"

	"samba4-admin/cmd"
)

func main() {
	if err := cmd.Execute(TemplatesFS, StaticFS); err != nil {
		slog.Error("Startup failure", "error", err)
		os.Exit(1)
	}
}
