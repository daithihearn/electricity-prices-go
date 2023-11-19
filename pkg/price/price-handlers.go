package price

import (
	"electricity-prices/pkg/api"
	"electricity-prices/pkg/date"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	PriceService Service
}

// GetPrices @Summary Get price info
// @Description Returns price info for the date provided. If no date is provided it defaults to today. The day should be given in a string form yyyy-MM-dd
// @Tags Price
// @ID get-prices
// @Produce  json
// @Param date query string false "Date in format yyyy-MM-dd"
// @Success 200 {object} []price.Price
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /price [get]
func (h *Handler) GetPrices(c *gin.Context) {

	// Get the date string from the request
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02")) // Default to today if not provided

	// Parse the date string
	d, err := date.ParseDate(dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: "Failed to parse date. Ensure it is in the format yyyy-MM-dd."})
		return
	}

	// Get the context from the request
	ctx := c.Request.Context()

	// Get the prices from the database
	prices, err := h.PriceService.GetDailyPrices(ctx, d)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
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
// @Success 200 {object} []price.DailyAverage
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /price/averages [get]
func (h *Handler) GetThirtyDayAverages(c *gin.Context) {

	// Get the date string from the request
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02")) // Default to today if not provided

	// Parse the date string
	d, err := date.ParseDate(dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: "Failed to parse date. Ensure it is in the format yyyy-MM-dd."})
		return
	}

	// Get the context from the request
	ctx := c.Request.Context()

	averages, err := h.PriceService.GetDailyAverages(ctx, d, 30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
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
// @Success 200 {object} price.DailyPriceInfo
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /price/dailyinfo [get]
func (h *Handler) GetDailyInfo(c *gin.Context) {

	// Get the date string from the request
	dateStr := c.DefaultQuery("date", time.Now().Format("2006-01-02")) // Default to today if not provided

	// Parse the date string
	d, err := date.ParseDate(dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: "Failed to parse date. Ensure it is in the format yyyy-MM-dd."})
		return
	}

	// Get the context from the request
	ctx := c.Request.Context()

	dailyInfo, err := h.PriceService.GetDailyInfo(ctx, d)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}
	if len(dailyInfo.Prices) == 0 {
		c.JSON(http.StatusNotFound, api.ErrorResponse{Message: "No data found for the given date."})
		return
	}

	c.IndentedJSON(http.StatusOK, dailyInfo)
}
