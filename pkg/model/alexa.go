package model

type AlexaApplication struct {
	ApplicationId string `json:"applicationId"`
}
type AlexaSlot struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
type AlexaIntent struct {
	Name  string      `json:"name"`
	Slots []AlexaSlot `json:"slots"`
}

type AlexaUser struct {
	UserId      string `json:"userId"`
	AccessToken string `json:"accessToken"`
}

type AlexaSession struct {
	SessionId   string           `json:"sessionId"`
	Application AlexaApplication `json:"application"`
}

type Request struct {
	Type      string      `json:"type"`
	Locale    string      `json:"locale"`
	RequestId string      `json:"requestId"`
	Timestamp string      `json:"timestamp"`
	Intent    AlexaIntent `json:"intent"`
}

type AlexaRequest struct {
	Version string       `json:"version"`
	Session AlexaSession `json:"session"`
	Request Request      `json:"request"`
}

type OutputSpeech struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type AlexaResponseBody struct {
	OutputSpeech     OutputSpeech `json:"outputSpeech"`
	ShouldEndSession bool         `json:"shouldEndSession"`
}
type AlexaSkillResponse struct {
	Version  string            `json:"version"`
	Response AlexaResponseBody `json:"response"`
}

type AlexaResponse struct {
	Uid            string `json:"uid"`
	UpdateDate     string `json:"updateDate"`
	TitleText      string `json:"titleText"`
	MainText       string `json:"mainText"`
	RedirectionUrl string `json:"redirectionUrl"`
}
