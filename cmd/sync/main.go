package main

import (
	"electricity-prices/pkg/db"
	"electricity-prices/pkg/service"
	"github.com/joho/godotenv"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	// Load .env file if it exists
	_ = godotenv.Load()
}

func main() {
	// Catch SIGINT or SIGTERM and gracefully close the database connection.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		db.CloseMongoConnection()
		os.Exit(0)
	}()

	// Sync with the API.
	service.SyncWithAPI()
}
