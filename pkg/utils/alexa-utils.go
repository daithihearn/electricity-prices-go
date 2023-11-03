package utils

import (
	"electricity-prices/pkg/model"
	"github.com/google/uuid"
	"time"
)

func WrapAlexaResponse(title, message string) model.AlexaResponse {
	return model.AlexaResponse{
		Uid:            uuid.New().String(),
		UpdateDate:     time.Now().Format("2006-01-02T15:04:05.000Z"),
		TitleText:      title,
		MainText:       message,
		RedirectionUrl: "https://elec.daithiapp.com/",
	}
}

func WrapAlexaSkillResponse(message string, endSession bool) model.AlexaSkillResponse {
	return model.AlexaSkillResponse{
		Version: "1.0",
		Response: model.AlexaResponseBody{
			OutputSpeech: model.OutputSpeech{
				Type: "PlainText",
				Text: message,
			},
			ShouldEndSession: endSession,
		},
	}
}
