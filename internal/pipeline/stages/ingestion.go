package stages

import (
	"context"
	"io"
	"net/http"

	"github.com/whatsapp-tts/internal/pipeline"
)

// IngestionStage recibe eventos HTTP del webhook.
type IngestionStage struct {
	verifyToken string
}

// NewIngestionStage crea un nuevo stage de ingestión.
func NewIngestionStage(verifyToken string) *IngestionStage {
	return &IngestionStage{verifyToken: verifyToken}
}

func (s *IngestionStage) Name() string { return "ingestion" }

func (s *IngestionStage) CanProcess(input interface{}) bool {
	_, ok := input.(*http.Request)
	return ok
}

func (s *IngestionStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
	req, ok := input.(*http.Request)
	if !ok {
		return nil, pipeline.ErrInvalidInput
	}

	// Handle GET for webhook verification (TODO: implement properly)
	if req.Method == "GET" {
		return nil, nil // Placeholder for verification
	}

	// Handle POST for webhook events
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	return &pipeline.RawEvent{
		Payload: body,
		Method:  req.Method,
		Headers: map[string]string{
			"content-type": req.Header.Get("content-type"),
		},
	}, nil
}
