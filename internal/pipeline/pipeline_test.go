package pipeline

import (
	"context"
	"testing"
)

type mockStage struct {
	name       string
	canProcess bool
	processFn  func(ctx context.Context, input interface{}) (interface{}, error)
}

func (s *mockStage) Name() string { return s.name }

func (s *mockStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
	return s.processFn(ctx, input)
}

func (s *mockStage) CanProcess(input interface{}) bool { return s.canProcess }

func TestPipelineExecution(t *testing.T) {
	p := NewPipeline(
		&mockStage{
			name:       "stage1",
			canProcess: true,
			processFn: func(ctx context.Context, input interface{}) (interface{}, error) {
				return "output1", nil
			},
		},
		&mockStage{
			name:       "stage2",
			canProcess: true,
			processFn: func(ctx context.Context, input interface{}) (interface{}, error) {
				return "output2", nil
			},
		},
	)

	ctx := context.Background()
	result, err := p.Execute(ctx, "input")

	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	if result == nil {
		t.Error("Execute() returned nil result")
	}

	if len(result.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(result.Errors))
	}
}

func TestPipelineStageError(t *testing.T) {
	p := NewPipeline(
		&mockStage{
			name:       "stage1",
			canProcess: true,
			processFn: func(ctx context.Context, input interface{}) (interface{}, error) {
				return nil, ErrInvalidInput
			},
		},
	)

	ctx := context.Background()
	result, err := p.Execute(ctx, "input")

	if err == nil {
		t.Error("Expected error, got nil")
	}

	if result == nil {
		t.Error("Execute() returned nil result")
	}

	if len(result.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(result.Errors))
	}
}

func TestPipelineStages(t *testing.T) {
	p := NewPipeline(
		&mockStage{name: "stage1", canProcess: true, processFn: nil},
		&mockStage{name: "stage2", canProcess: true, processFn: nil},
	)

	stages := p.Stages()
	if len(stages) != 2 {
		t.Errorf("Expected 2 stages, got %d", len(stages))
	}

	if stages[0].Name() != "stage1" {
		t.Errorf("Expected stage1, got %s", stages[0].Name())
	}

	if stages[1].Name() != "stage2" {
		t.Errorf("Expected stage2, got %s", stages[1].Name())
	}
}
