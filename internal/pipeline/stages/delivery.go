package stages

import (
	"context"

	"github.com/whatsapp-tts/internal/domain"
	"github.com/whatsapp-tts/internal/pipeline"
)

// DeliveryStage entrega la respuesta al usuario.
type DeliveryStage struct{}

func NewDeliveryStage() *DeliveryStage {
	return &DeliveryStage{}
}

func (s *DeliveryStage) Name() string { return "delivery" }

func (s *DeliveryStage) CanProcess(input interface{}) bool {
	_, ok := input.(*domain.AudioAsset)
	return ok
}

func (s *DeliveryStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
	audio, ok := input.(*domain.AudioAsset)
	if !ok {
		return nil, pipeline.ErrInvalidInput
	}

	// Stub: log that we would send audio
	// In real implementation, this would call WhatsApp API
	_ = audio

	return nil, nil
}
