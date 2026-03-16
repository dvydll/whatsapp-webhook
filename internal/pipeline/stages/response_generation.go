package stages

import (
	"context"

	"github.com/whatsapp-tts/internal/domain"
	"github.com/whatsapp-tts/internal/pipeline"
)

// ResponseGenerationStage genera una respuesta de texto.
type ResponseGenerationStage struct{}

func NewResponseGenerationStage() *ResponseGenerationStage {
	return &ResponseGenerationStage{}
}

func (s *ResponseGenerationStage) Name() string { return "response_generation" }

func (s *ResponseGenerationStage) CanProcess(input interface{}) bool {
	_, ok := input.(*domain.UserMessage)
	return ok
}

func (s *ResponseGenerationStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
	msg, ok := input.(*domain.UserMessage)
	if !ok {
		return nil, pipeline.ErrInvalidInput
	}

	// Placeholder response (v1)
	responseText := "Message received. Generating audio response."

	resp := domain.NewResponseMessage(
		msg.UserID,
		domain.ResponseTypeAudio,
		responseText,
		msg.Metadata.BusinessPhoneID,
	)

	return resp, nil
}
