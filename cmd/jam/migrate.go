package main

import (
	"context"
	"time"

	"github.com/bhbdev/jam/config"
	"github.com/bhbdev/jam/lib/database"
	"github.com/bhbdev/jam/lib/logger"
	"github.com/bhbdev/jam/models"

	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(models.MODELS...)
	if err != nil {
		logger.Error("error automigrating models", "error", err)
		return
	}
}

func migrate(ctx context.Context, cfg *config.Config) {
	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database", "error", err)
		return
	}
	Migrate(db.Instance())
}

func init() {
	var migrateCmd = &cobra.Command{
		Use:     "migrate",
		Short:   "JAM migration command",
		Example: "jam migrate",
		Run: func(cmd *cobra.Command, args []string) {
			cfg, err := config.Load()
			if err != nil {
				logger.Fatal("Failed to get config", "error", err)
			}

			start := time.Now()
			logger.Info("Starting migration")
			migrate(context.Background(), cfg)
			logger.Info("Migration completed", "time", time.Since(start))
		},
	}

	rootCmd.AddCommand(migrateCmd)
}
