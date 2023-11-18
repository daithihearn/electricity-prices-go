package alexa

import (
	"github.com/google/uuid"
	"time"
)

func WrapAlexaResponse(title, message string) AlexaResponse {
	return AlexaResponse{
		Uid:            uuid.New().String(),
		UpdateDate:     time.Now().Format("2006-01-02T15:04:05.000Z"),
		TitleText:      title,
		MainText:       message,
		RedirectionUrl: "https://elec.daithiapp.com/",
	}
}

func WrapAlexaSkillResponse(message string, endSession bool) AlexaSkillResponse {
	return AlexaSkillResponse{
		Version: "1.0",
		Response: AlexaResponseBody{
			OutputSpeech: OutputSpeech{
				Type: "PlainText",
				Text: message,
			},
			ShouldEndSession: endSession,
		},
	}
}
