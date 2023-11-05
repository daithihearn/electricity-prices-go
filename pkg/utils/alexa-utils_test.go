package utils

import (
	"electricity-prices/pkg/model"
	"testing"
)

func TestWrapAlexaResponse(t *testing.T) {
	testCases := []struct {
		name     string
		title    string
		message  string
		expected model.AlexaResponse
	}{
		{
			name:    "Test WrapAlexaResponse",
			title:   "Title",
			message: "Message",
		},
		{
			name:    "Test WrapAlexaResponse 2",
			title:   "Title 2",
			message: "Message 2",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := WrapAlexaResponse(tc.title, tc.message)
			if actual.TitleText != tc.title {
				t.Errorf("expected %s, got %s", tc.title, actual.TitleText)
			}
			if actual.MainText != tc.message {
				t.Errorf("expected %s, got %s", tc.message, actual.MainText)
			}
			if actual.RedirectionUrl != "https://elec.daithiapp.com/" {
				t.Errorf("expected https://elec.daithiapp.com/, got %s", actual.RedirectionUrl)
			}
			if actual.UpdateDate == "" {
				t.Errorf("expected non-empty UpdateDate, got %s", actual.UpdateDate)
			}
			if actual.Uid == "" {
				t.Errorf("expected non-empty Uid, got %s", actual.Uid)
			}
		})
	}
}

func TestWrapAlexaSkillResponse(t *testing.T) {
	testCases := []struct {
		name     string
		message  string
		end      bool
		expected model.AlexaSkillResponse
	}{
		{
			name:    "End session",
			message: "Message",
			end:     true,
		},
		{
			name:    "Keep session alive",
			message: "Message 2",
			end:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := WrapAlexaSkillResponse(tc.message, tc.end)
			if actual.Version != "1.0" {
				t.Errorf("expected 1.0, got %s", actual.Version)
			}
			if actual.Response.OutputSpeech.Type != "PlainText" {
				t.Errorf("expected PlainText, got %s", actual.Response.OutputSpeech.Type)
			}
			if actual.Response.OutputSpeech.Text != tc.message {
				t.Errorf("expected %s, got %s", tc.message, actual.Response.OutputSpeech.Text)
			}
			if actual.Response.ShouldEndSession != tc.end {
				t.Errorf("expected %t, got %t", tc.end, actual.Response.ShouldEndSession)
			}
		})
	}
}
