package stages

import (
	"context"

	"github.com/whatsapp-tts/internal/domain"
	"github.com/whatsapp-tts/internal/pipeline"
)

// TTSGenerationStage convierte texto a audio.
type TTSGenerationStage struct{}

func NewTTSGenerationStage() *TTSGenerationStage {
	return &TTSGenerationStage{}
}

func (s *TTSGenerationStage) Name() string { return "tts_generation" }

func (s *TTSGenerationStage) CanProcess(input interface{}) bool {
	_, ok := input.(*domain.ResponseMessage)
	return ok
}

func (s *TTSGenerationStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
	resp, ok := input.(*domain.ResponseMessage)
	if !ok {
		return nil, pipeline.ErrInvalidInput
	}

	// Stub: generate placeholder audio
	audio := domain.NewAudioAsset(
		domain.FormatAAC,
		domain.CodecOpus,
		[]byte("stub-audio-data"),
	)

	// Store response in metadata for delivery stage
	resp.AudioID = audio.AudioID

	return audio, nil
}
