package main

import (
	"log/slog"
	"os"

	"go-samba4/cmd"
)

func main() {
	if err := cmd.Execute(TemplatesFS, StaticFS); err != nil {
		slog.Error("Startup failure", "error", err)
		os.Exit(1)
	}
}
