package domain

import "time"

// ResponseMessage representa una respuesta generada por el sistema.
type ResponseMessage struct {
	ResponseID   string
	TargetUser   string
	ResponseType ResponseType
	Text         string
	AudioURL     string
	AudioID      string
	CreatedAt    time.Time
	Metadata     ResponseMetadata
}

// ResponseMetadata contiene metadata adicional de la respuesta.
type ResponseMetadata struct {
	OriginalMessageID string
	PhoneNumberID     string
	CorrelationID     string
}

// ResponseType representa el tipo de respuesta.
type ResponseType string

const (
	ResponseTypeText  ResponseType = "text"
	ResponseTypeAudio ResponseType = "audio"
)

// NewResponseMessage crea un nuevo ResponseMessage.
func NewResponseMessage(
	targetUser string,
	responseType ResponseType,
	text string,
	phoneNumberID string,
) *ResponseMessage {
	return &ResponseMessage{
		ResponseID:   generateID(),
		TargetUser:   targetUser,
		ResponseType: responseType,
		Text:         text,
		CreatedAt:    time.Now(),
		Metadata: ResponseMetadata{
			PhoneNumberID: phoneNumberID,
		},
	}
}
