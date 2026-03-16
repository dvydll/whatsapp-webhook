# SPEC-010: Future Capability Space

## 1. Propósito

Esta especificación reserva espacio para capacidades futuras que deben ser anticipadas en la arquitectura actual sin ser implementadas todavía.

## 2. Procesamiento Asíncrono

### 2.1 Job Queue

**Descripción**: Permitir procesamiento no-bloqueante de mensajes.

**Arquitectura sugerida**:
- Cola de mensajes (Redis, RabbitMQ, o similar)
- Workers dedicados para procesamiento
- Estado del jobtrackeable

**Interfaces a preparar**:
```go
type JobQueue interface {
    Enqueue(ctx context.Context, job Job) (string, error)
    Dequeue(ctx context.Context) (Job, error)
    UpdateStatus(ctx context.Context, jobID string, status JobStatus) error
}

type JobStatus string

const (
    JobStatusPending   JobStatus = "pending"
    JobStatusProcessing JobStatus = "processing"
    JobStatusCompleted JobStatus = "completed"
    JobStatusFailed    JobStatus = "failed"
)
```

## 3. Background Workflows

### 3.1 Workflow Engine

**Descripción**: Orquestar flujos de múltiples pasos.

**Capacidades**:
- Pasos condicionales
- Ramificaciones
- Timeouts por paso
- Rollback

**Interfaces a preparar**:
```go
type WorkflowEngine interface {
    Start(ctx context.Context, wf Workflow) (string, error)
    GetStatus(ctx context.Context, workflowID string) (WorkflowStatus, error)
}

type Workflow struct {
    ID      string
    Name    string
    Steps   []Step
    Context map[string]interface{}
}

type Step struct {
    Name      string
    Type      StepType
    Handler   StepHandler
    Timeout   time.Duration
    OnError   string  // siguiente step en caso de error
}
```

## 4. Multi-Message Responses

### 4.1 Respuestas Múltiples

**Descripción**: Enviar múltiples mensajes como respuesta.

**Arquitectura**:
```go
type MultiResponse struct {
    Responses []Response
    Strategy  DeliveryStrategy  // parallel, sequential, batch
}

type DeliveryStrategy string

const (
    StrategyParallel   DeliveryStrategy = "parallel"
    StrategySequential DeliveryStrategy = "sequential"
    StrategyBatch      DeliveryStrategy = "batch"
)
```

## 5. Streaming Audio

### 5.1 Audio Streaming

**Descripción**: Enviar audio en chunks en lugar de archivo completo.

**Interfaces a preparar**:
```go
type AudioStreamer interface {
    Stream(ctx context.Context, userID string, audio io.Reader) error
}
```

## 6. Conversational Context Memory

### 6.1 Context Storage

**Descripción**: Mantener historial de conversación por usuario.

**Arquitectura**:
```go
type ConversationStore interface {
    AddMessage(ctx context.Context, userID string, msg *Message) error
    GetHistory(ctx context.Context, userID string, limit int) ([]Message, error)
    ClearHistory(ctx context.Context, userID string) error
}

// Implementaciones posibles:
// - Redis (cache)
// - PostgreSQL (persistencia)
// - In-memory (desarrollo)
```

### 6.2 Session Management

```go
type Session struct {
    UserID      string
    CreatedAt   time.Time
    LastActive  time.Time
    Context     map[string]interface{}
}

type SessionManager interface {
    GetOrCreate(ctx context.Context, userID string) (*Session, error)
    Update(ctx context.Context, session *Session) error
    Expire(ctx context.Context, userID string) error
}
```

## 7. Integración con LLM Systems

### 7.1 LLM Provider Interface

**Descripción**: Conectar con sistemas LLM para respuestas generativas.

**Interfaces a preparar**:
```go
type LLMProvider interface {
    Generate(ctx context.Context, prompt string, opts LLMOptions) (LLMResponse, error)
    StreamGenerate(ctx context.Context, prompt string, opts LLMOptions, cb func(chunk string)) error
}

type LLMOptions struct {
    Model       string
    Temperature float64
    MaxTokens   int
    SystemPrompt string
}

type LLMResponse struct {
    Content   string
    Usage     Usage
    FinishReason string
}

type Usage struct {
    PromptTokens     int
    CompletionTokens int
    TotalTokens     int
}
```

### 7.2 Prompt Templates

```go
type PromptTemplate struct {
    System string
    User   string
}

var DefaultPromptTemplate = PromptTemplate{
    System: "Eres un asistente conversacional. Responde de manera útil y concisa.",
    User:   "Mensaje del usuario: {{.Message}}",
}
```

## 8. Lista de Capacidades Reservadas

| Capability | Prioridad | Dependencias |
|------------|-----------|--------------|
| Job Queue | Alta | Pipeline actual |
| Background Workflows | Media | Job Queue |
| Multi-Message Responses | Media | Ninguna |
| Streaming Audio | Baja | WhatsApp API support |
| Context Memory | Alta | Storage backend |
| LLM Integration | Alta | Context Memory |

## 9. Notas de Diseño

- No implementar todavía, pero diseñar la arquitectura para facilitar estas adiciones
- Usar interfaces para permitir reemplazo de implementaciones
- Mantener el pipeline v1 simple y añadir complejidad después