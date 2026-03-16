package pipeline

import "context"

// Stage representa una etapa del pipeline.
type Stage interface {
	Name() string
	Process(ctx context.Context, input interface{}) (interface{}, error)
	CanProcess(input interface{}) bool
}

// RawEvent representa un evento raw del webhook.
type RawEvent struct {
	Payload []byte
	Headers map[string]string
	Method  string
}
