package api

import (
	"electricity-prices/pkg/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"electricity-prices/pkg/model"
	"electricity-prices/pkg/utils"
)

// GetPrices @Summary Get price info
// @Description Returns price info for the date provided. If no date is provided it defaults to today. The day should be given in a string form yyyy-MM-dd
// @Tags Price
// @ID get-prices
// @Produce  json
// @Param date query string false "Date in format yyyy-MM-dd"
// @Success 200 {object} []model.Price
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /price [get]
func GetPrices(c *gin.Context) {

	// Get the date string from the request
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02")) // Default to today if not provided

	// Parse the date string
	date, err := utils.ParseDate(dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to parse date. Ensure it is in the format yyyy-MM-dd."})
		return
	}

	// Get the prices from the database
	prices, err := service.GetDailyPrices(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, prices)
}

// GetThirtyDayAverages @Summary Get daily averages
// @Description Returns daily averages for the date provided and the previous 30 days.
// @Tags Price
// @ID get-daily-averages
// @Produce  json
// @Param date query string false "Date in format yyyy-MM-dd"
// @Success 200 {object} []model.DailyAverage
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /price/averages [get]
func GetThirtyDayAverages(c *gin.Context) {

	// Get the date string from the request
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02")) // Default to today if not provided

	// Parse the date string
	date, err := utils.ParseDate(dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to parse date. Ensure it is in the format yyyy-MM-dd."})
		return
	}

	averages, err := service.GetDailyAverages(date, 30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, averages)
}

// GetDailyInfo @Summary Get daily info
// @Description Returns daily info for the date provided.
// @Tags Price
// @ID get-daily-info
// @Produce  json
// @Param date query string false "Date in format yyyy-MM-dd"
// @Success 200 {object} model.DailyPriceInfo
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /price/dailyinfo [get]
func GetDailyInfo(c *gin.Context) {

	// Get the date string from the request
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02")) // Default to today if not provided

	// Parse the date string
	date, err := utils.ParseDate(dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to parse date. Ensure it is in the format yyyy-MM-dd."})
		return
	}

	dailyInfo, err := service.GetDailyInfo(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: err.Error()})
		return
	}
	if len(dailyInfo.Prices) == 0 {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Message: "No data found for the given date."})
		return
	}

	c.IndentedJSON(http.StatusOK, dailyInfo)
}
