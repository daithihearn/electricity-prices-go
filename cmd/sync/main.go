package main

import (
	"context"
	"electricity-prices/pkg/db"
	"electricity-prices/pkg/sync"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	// Load .env file if it exists
	_ = godotenv.Load()
}

func main() {
	// Create a cancellable context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel() // Cancel the context upon receiving the signal

		// Create a new context for the graceful shutdown procedure
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()

		// Gracefully close the database connection
		if err := db.CloseMongoConnection(shutdownCtx); err != nil {
			// Handle error (e.g., log it)
			os.Exit(1)
		}
		os.Exit(0)
	}()

	// Sync with the API.
	sync.SyncWithAPI(ctx)
	cancel()

	// Wait for the cancellation of the context (due to signal handling)
	<-ctx.Done()
}
