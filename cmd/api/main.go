// @title Electricity Prices API
// @version 2.1.2
// @description Returns PVPC electricity prices for a given range
// @BasePath /api/v1
package main

import (
	"context"
	_ "electricity-prices/docs"
	"electricity-prices/pkg/api"
	"electricity-prices/pkg/db"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Set up the API routes.
	router := gin.Default()

	// Configure CORS with custom settings
	config := cors.Config{
		AllowOrigins:  []string{"https://elec.daithiapp.com", "http://localhost:888", "http://localhost:3000"},
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
	}
	router.Use(cors.New(config))

	// Configure the routes
	router.GET("/api/v1/price", api.GetPrices)
	router.GET("/api/v1/price/averages", api.GetThirtyDayAverages)
	router.GET("/api/v1/price/dailyinfo", api.GetDailyInfo)
	router.GET("/api/v1/alexa", api.GetFullFeed)
	router.POST("/api/v1/alexa-skill", api.ProcessSkillRequest)

	// Use the generated docs in the docs package.
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// Start the server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err := router.Run(":" + port)
	if err != nil {
		return
	}

	// Wait for the cancellation of the context (due to signal handling)
	<-ctx.Done()
}
