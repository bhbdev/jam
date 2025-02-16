package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bhbdev/jam/app"
	"github.com/bhbdev/jam/config"
	"github.com/bhbdev/jam/lib/logger"

	"github.com/spf13/cobra"
)

func init() {
	runAll := &cobra.Command{
		Use:   "run",
		Short: "run server",
		Run:   run,
	}
	runAll.Flags().BoolP("migrate", "m", false, "Migrate database")

	// add to root command
	rootCmd.AddCommand(runAll)
}

func run(cmd *cobra.Command, args []string) {

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to receive termination signals
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Load config
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to get config", "error", err)
	}

	// Set debug level if development mode is enabled
	logger.SetDebugLevel(cfg.Development)

	// if migrate flag is set, run migration for all enabled services
	if cmd.Flag("migrate").Value.String() == "true" {
		migrate(ctx, cfg)
	}

	// Wait group to wait for all components to shut down
	var wg sync.WaitGroup

	// Start app server
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.ServerApp(cfg)
	}()

	// Wait for termination signal or cancellation
	select {
	case <-ctx.Done():
		// Shutdown triggered, wait for all components to shut down
		wg.Wait()
	case <-sigs:
		// Termination signal received, initiate shutdown
		cancel()
	}
}
