package main

import (
	"context"
	"electricity-prices/pkg/db"
	"electricity-prices/pkg/esios"
	"electricity-prices/pkg/price"
	"electricity-prices/pkg/ree"
	"electricity-prices/pkg/sync"
	"github.com/joho/godotenv"
	"log"
	"net/http"
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

	// Configure services
	col, err := db.GetCollection(ctx)
	if err != nil {
		cancel()
		log.Fatal("Failed to get collection: ", err)
	}
	priceCollection := price.PriceCollection{Col: col}
	priceService := price.Service{Collection: priceCollection}
	reeClient := ree.Client{
		Http: &http.Client{Timeout: time.Second * 30},
	}
	esiosClient := esios.Client{
		Http: &http.Client{Timeout: time.Second * 30},
	}
	syncService := sync.Syncer{PriceService: &priceService, PrimaryClient: &reeClient, SecondaryClient: &esiosClient}

	// Sync with the API.
	synced, err := syncService.Sync(ctx, time.Now().AddDate(0, 0, 1))
	if err != nil {
		cancel()
		log.Fatal("Failed to sync with API: ", err)
	}
	if synced {
		log.Println("Synced successfully")
	} else {
		cancel()
		log.Fatal("Failed to sync fully...")
	}
	cancel()

	// Wait for the cancellation of the context (due to signal handling)
	<-ctx.Done()
}
