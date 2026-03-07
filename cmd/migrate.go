package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go-samba4/internal/models"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("Starting Database Migrations...")

		var dialector gorm.Dialector
		if globalCfg.Database.Driver == "mysql" {
			dialector = mysql.Open(globalCfg.Database.DSN)
		} else {
			dialector = sqlite.Open(globalCfg.Database.Path)
		}

		db, err := gorm.Open(dialector, &gorm.Config{})
		if err != nil {
			slog.Error("Failed to connect to database for migration", "err", err)
			os.Exit(1)
		}

		if err := db.AutoMigrate(&models.Session{}, &models.AuditLog{}, &models.Setting{}); err != nil {
			slog.Error("Migration failed", "err", err)
			os.Exit(1)
		}

		slog.Info("Database migrations completed successfully.")
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
