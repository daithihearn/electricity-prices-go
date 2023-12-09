// @title Electricity Prices API
// @version 2.1.5
// @description Returns PVPC electricity prices for a given range
// @BasePath /api/v1
package main

import (
	"context"
	_ "electricity-prices/docs"
	"electricity-prices/pkg/alexa"
	"electricity-prices/pkg/db"
	"electricity-prices/pkg/i18n"
	"electricity-prices/pkg/price"
	"log"
	"os"
	"os/signal"
	"strings"
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

	// Initialise Translations
	i18n.InitialiseTranslations()

	// Get the db name and collection name
	dbName := os.Getenv("MONGODB_DB")
	if dbName == "" {
		dbName = "electricity-prices"
	}

	colName := os.Getenv("MONGODB_COLLECTION")
	if colName == "" {
		colName = "prices"
	}

	// Configure services
	col, err := db.GetCollection(ctx, dbName, colName)
	if err != nil {
		cancel()
		log.Fatal("Failed to get collection: ", err)
	}
	priceCollection := price.PriceCollection{Col: col}
	priceService := price.Receiver{Collection: priceCollection}
	priceHandler := price.Handler{PriceService: &priceService}
	alexaService := alexa.Service{PriceService: &priceService}
	alexaHandler := alexa.Handler{AlexaService: alexaService}

	// Set up the API routes.
	router := gin.Default()

	// Configure CORS with custom settings
	// Get the environment variable
	origins := os.Getenv("CORS_ALLOWED_ORIGINS")

	// Check if the environment variable is empty and set a default value
	if origins == "" {
		origins = "http://localhost:888,http://localhost:3000" // Replace with your default list
	}

	config := cors.Config{
		AllowOrigins:  strings.Split(origins, ","),
		AllowMethods:  []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:  []string{"Origin", "Content-Length", "Content-Type"},
		ExposeHeaders: []string{"Content-Length"},
	}
	router.Use(cors.New(config))

	// Configure the routes
	router.GET("/api/v1/price", priceHandler.GetPrices)
	router.GET("/api/v1/price/averages", priceHandler.GetThirtyDayAverages)
	router.GET("/api/v1/price/dailyinfo", priceHandler.GetDailyInfo)
	router.GET("/api/v1/alexa", alexaHandler.GetFullFeed)
	router.POST("/api/v1/alexa-skill", alexaHandler.ProcessSkillRequest)

	// Use the generated docs in the docs package.
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))

	// Start the server.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	err = router.Run(":" + port)
	if err != nil {
		return
	}

	// Wait for the cancellation of the context (due to signal handling)
	<-ctx.Done()
}
