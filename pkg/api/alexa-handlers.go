package api

import (
	"electricity-prices/pkg/model"
	"electricity-prices/pkg/service"
	"electricity-prices/pkg/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/text/language"
	"io"
	"net/http"
	"time"
)

// GetFullFeed @Summary Get full feed
// @Description Returns the full feed for an alexa flash briefing.
// @Tags Alexa
// @ID get-full-feed
// @Produce  json
// @Param lang query string false "Language in format es or en"
// @Success 200 {object} model.AlexaResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /alexa [get]
func GetFullFeed(c *gin.Context) {
	lang, err := language.Parse(c.DefaultQuery("lang", "es"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: "Failed to parse language. Ensure it is in the format es or en."})
		return
	}

	now := time.Now()

	title := service.GetTitle(lang)

	feed, err := service.GetFullFeed(now, lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Message: err.Error()})
		return
	}

	response := utils.WrapAlexaResponse(title, feed)
	c.IndentedJSON(http.StatusOK, response)
}

// ProcessSkillRequest @Summary Process request from the Alexa skill
// @Description Processes the request from the Alexa skill.
// @Tags Alexa
// @ID process-skill-request
// @Accept  json
// @Produce  json
// @Success 200 {object} model.AlexaResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /alexa-skill [post]
func ProcessSkillRequest(c *gin.Context) {
	// Get Raw JSON body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading body"})
		return
	}
	rawJSON := string(body)

	// Unmarshal JSON into AlexaRequest struct
	var request model.AlexaRequest
	err = json.Unmarshal(body, &request)

	// Validate the request
	if err := service.ValidateAlexaRequest(c.Request, rawJSON, request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Message: err.Error()})
		return
	}

	// Get Locale
	locale := request.Request.Locale
	lang, err := language.Parse(locale)
	if err != nil {
		lang = language.Spanish
	}

	// Parse the request
	response := service.ParseAlexaSkillRequest(request.Request.Intent, lang)
	c.IndentedJSON(http.StatusOK, response)
}
