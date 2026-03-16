package pipeline

import (
	"context"
	"errors"
)

var (
	ErrInvalidInput   = errors.New("invalid input for stage")
	ErrPipelineFailed = errors.New("pipeline execution failed")
)

// PipelineContext contiene el contexto compartido entre etapas.
type PipelineContext struct {
	RequestID string
	TraceID   string
	Errors    []PipelineError
	Metadata  map[string]interface{}
}

// PipelineError representa un error en el pipeline.
type PipelineError struct {
	Stage   string
	Err     error
	IsFatal bool
}

// Pipeline ejecuta el flujo completo de procesamiento.
type Pipeline struct {
	stages []Stage
}

// NewPipeline crea un nuevo pipeline.
func NewPipeline(stages ...Stage) *Pipeline {
	return &Pipeline{stages: stages}
}

// Execute ejecuta el pipeline con el input dado.
func (p *Pipeline) Execute(ctx context.Context, input interface{}) (*PipelineContext, error) {
	pc := &PipelineContext{
		Metadata: make(map[string]interface{}),
	}

	current := input

	for _, stage := range p.stages {
		output, err := stage.Process(ctx, current)
		if err != nil {
			pc.Errors = append(pc.Errors, PipelineError{
				Stage:   stage.Name(),
				Err:     err,
				IsFatal: true,
			})
			return pc, err
		}
		current = output
	}

	return pc, nil
}

// RegisterStage registra una etapa en el pipeline.
func (p *Pipeline) RegisterStage(stage Stage) {
	p.stages = append(p.stages, stage)
}

// Stages retorna las etapas registradas.
func (p *Pipeline) Stages() []Stage {
	return p.stages
}
