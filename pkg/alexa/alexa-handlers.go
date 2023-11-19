package alexa

import (
	"electricity-prices/pkg/api"
	"electricity-prices/pkg/i18n"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

type Handler struct {
	AlexaService Service
}

// GetFullFeed @Summary Get full feed
// @Description Returns the full feed for an alexa flash briefing.
// @Tags Alexa
// @ID get-full-feed
// @Produce  json
// @Param lang query string false "Language in format es or en"
// @Success 200 {object} alexa.AlexaResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /alexa [get]
func (s *Handler) GetFullFeed(c *gin.Context) {
	// Parse language from request
	lang := i18n.ParseLanguage(c.DefaultQuery("lang", "es"))

	now := time.Now()

	title := GetTitle(lang)

	// Get the context from the request
	ctx := c.Request.Context()

	feed, err := s.AlexaService.GetFullFeed(ctx, now, lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, api.ErrorResponse{Message: err.Error()})
		return
	}

	response := WrapAlexaResponse(title, feed)
	c.IndentedJSON(http.StatusOK, response)
}

// ProcessSkillRequest @Summary Process request from the Alexa skill
// @Description Processes the request from the Alexa skill.
// @Tags Alexa
// @ID process-skill-request
// @Accept  json
// @Produce  json
// @Success 200 {object} alexa.AlexaResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /alexa-skill [post]
func (s *Handler) ProcessSkillRequest(c *gin.Context) {
	// Get Raw JSON body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading body"})
		return
	}
	rawJSON := string(body)

	// Unmarshal JSON into AlexaRequest struct
	var request AlexaRequest
	err = json.Unmarshal(body, &request)

	// Validate the request
	if err := ValidateAlexaRequest(c.Request, rawJSON, request); err != nil {
		c.JSON(http.StatusBadRequest, api.ErrorResponse{Message: err.Error()})
		return
	}

	// Parse language from request
	locale := request.Request.Locale
	lang := i18n.ParseLanguage(locale)

	// Get the context from the request
	ctx := c.Request.Context()

	// Parse the request
	response := s.AlexaService.ProcessAlexaSkillRequest(ctx, request.Request.Intent, lang)
	c.IndentedJSON(http.StatusOK, response)
}
