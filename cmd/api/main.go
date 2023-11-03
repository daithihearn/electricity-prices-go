// @title Electricity Prices API
// @version 2.0.0
// @description Returns PVPC electricity prices for a given range
// @BasePath /api/v1
package main

import (
	_ "electricity-prices/docs"
	"electricity-prices/pkg/api"
	"electricity-prices/pkg/db"
	"os"
	"os/signal"
	"syscall"

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
	// Catch SIGINT or SIGTERM and gracefully close the database connection.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		db.CloseMongoConnection()
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
}
