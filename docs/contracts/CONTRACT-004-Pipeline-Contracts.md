# CONTRACT-004: Pipeline Contracts

## 1. Stage Interface

```go
package pipeline

import "context"

// Stage representa una etapa del pipeline.
type Stage interface {
    // Name retorna el nombre de la etapa.
    Name() string
    
    // Process ejecuta la lógica de la etapa.
    // Retorna el contexto modificado y un error si falla.
    Process(ctx context.Context, input interface{}) (interface{}, error)
    
    // CanProcess determina si esta etapa puede procesar el input.
    CanProcess(input interface{}) bool
}
```

## 2. Pipeline Context

```go
package pipeline

import "github.com/whatsapp-tts/internal/domain"

// PipelineContext contiene el contexto compartido entre etapas.
type PipelineContext struct {
    // RequestID es el ID de correlación de la solicitud.
    RequestID string
    
    // TraceID para trazabilidad.
    TraceID string
    
    // Message es el mensaje normalizado.
    Message *domain.UserMessage
    
    // Response es la respuesta generada.
    Response *domain.ResponseMessage
    
    // Audio es el audio generado.
    Audio *domain.AudioAsset
    
    // Errors acumula errores durante el pipeline.
    Errors []PipelineError
    
    // Metadata adicional del pipeline.
    Metadata map[string]interface{}
}

// PipelineError representa un error en el pipeline.
type PipelineError struct {
    Stage   string
    Err     error
    IsFatal bool
}
```

## 3. Pipeline Interface

```go
package pipeline

import "context"

// Pipeline ejecuta el flujo completo de procesamiento.
type Pipeline interface {
    // Execute ejecuta el pipeline con el input dado.
    Execute(ctx context.Context, input interface{}) (*PipelineContext, error)
    
    // RegisterStage registra una etapa en el pipeline.
    RegisterStage(stage Stage)
    
    // Stages retorna las etapas registradas.
    Stages() []Stage
}
```

## 4. Stage Contracts Específicos

### 4.1 EventReceiver

```go
package pipeline

// EventReceiver recibe eventos HTTP del webhook.
type EventReceiver interface {
    // ReceiveExtrae el evento raw del request HTTP.
    Receive(ctx context.Context, request interface{}) (RawEvent, error)
}

type RawEvent struct {
    Payload   []byte
    Headers   map[string]string
    Method    string
}
```

### 4.2 MessageNormalizer

```go
package pipeline

import "github.com/whatsapp-tts/internal/domain"

// MessageNormalizer normaliza eventos raw a UserMessage.
type MessageNormalizer interface {
    // Normalize convierte un evento raw a UserMessage.
    Normalize(ctx context.Context, event RawEvent) (*domain.UserMessage, error)
}
```

### 4.3 MessageProcessor

```go
package pipeline

import "github.com/whatsapp-tts/internal/domain"

// MessageProcessor procesa el mensaje y genera una respuesta.
type MessageProcessor interface {
    // Process genera una respuesta para el mensaje.
    Process(ctx context.Context, message *domain.UserMessage) (*domain.ResponseMessage, error)
}
```

### 4.4 TTSEngine

```go
package pipeline

import "github.com/whatsapp-tts/internal/domain"

// TTSEngine convierte texto a audio.
type TTSEngine interface {
    // GenerateAudio genera audio a partir de texto.
    GenerateAudio(ctx context.Context, text string) (*domain.AudioAsset, error)
}
```

### 4.5 DeliveryAdapter

```go
package pipeline

import "github.com/whatsapp-tts/internal/domain"

// DeliveryAdapter entrega la respuesta al usuario.
type DeliveryAdapter interface {
    // Deliver envía la respuesta al usuario.
    Deliver(ctx context.Context, response *domain.ResponseMessage, audio *domain.AudioAsset) error
}
```

## 5. Interfaz Unificada (Wrapper)

```go
package pipeline

// StageWrapper envolvuelve las interfaces específicas en una interfaz Stage.
type StageWrapper struct {
    name       string
    receiver   interface{}  // Una de las interfaces específicas
}

func NewStageWrapper(name string, receiver interface{}) *StageWrapper {
    return &StageWrapper{
        name:     name,
        receiver: receiver,
    }
}

func (s *StageWrapper) Name() string { return s.name }

func (s *StageWrapper) Process(ctx context.Context, input interface{}) (interface{}, error) {
    // Delegate al receiver apropiado basado en tipo
    switch r := s.receiver.(type) {
    case EventReceiver:
        return r.Receive(ctx, input)
    case MessageNormalizer:
        return r.Normalize(ctx, input.(RawEvent))
    case MessageProcessor:
        return r.Process(ctx, input.(*domain.UserMessage))
    case TTSEngine:
        return r.GenerateAudio(ctx, input.(string))
    case DeliveryAdapter:
        resp := input.(struct {
            Response *domain.ResponseMessage
            Audio    *domain.AudioAsset
        })
        return nil, r.Deliver(ctx, resp.Response, resp.Audio)
    }
    return nil, ErrUnknownStageType
}
```

## 6. Errores del Pipeline

```go
package pipeline

import "errors"

var (
    ErrInvalidInput        = errors.New("invalid input for stage")
    ErrUnknownStageType    = errors.New("unknown stage type")
    ErrPipelineFailed      = errors.New("pipeline execution failed")
    ErrStageNotRegistered  = errors.New("stage not registered")
)
```

## 7. Notas

- Las interfaces específicas permiten implementación más precisa
- El Stage wrapper permite usar cualquier interfaz como Stage
- El PipelineContext es el medio de comunicación entre stages
- Cada stage puede agregar errores al contexto sin detener el pipeline (si no es fatal)