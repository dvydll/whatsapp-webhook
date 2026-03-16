# SPEC-004: Processing Pipeline Specification

## 1. Visión General

El pipeline de procesamiento es el núcleo del sistema que orquesta el flujo de un mensaje desde su recepción hasta la respuesta final. Está diseñado como una cadena de responsabilidad (Chain of Responsibility) que permite插入 nuevos pasos sin modificar los existentes.

## 2. Arquitectura del Pipeline

```
┌─────────────────────────────────────────────────────────────────────┐
│                         Pipeline Engine                              │
├─────────────────────────────────────────────────────────────────────┤
│                                                                      │
│   ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐    │
│   │   Stage  │--> │   Stage  │--> │   Stage  │--> │   Stage  │    │
│   │    1     │    │    2     │    │    3     │    │    N     │    │
│   │ Ingestion│    │Normalize │    │ Response │    │ Deliver  │    │
│   └──────────┘    └──────────┘    │  Generate │    │  Adapter │    │
│                                   └──────────┘    └──────────┘    │
│                                                                      │
└─────────────────────────────────────────────────────────────────────┘
```

## 3. Stages del Pipeline (v1)

### Stage 1: Event Ingestion
- **Input**: HTTP Request (Webhook)
- **Output**: Raw message payload
- **Responsabilidad**: Validar webhook, extraer payload

### Stage 2: Message Normalization
- **Input**: Raw payload
- **Output**: Internal Message model
- **Responsabilidad**: Convertir payload externo a modelo interno

### Stage 3: Processing Pipeline
- **Input**: Normalized message
- **Output**: Processed message with context
- **Responsabilidad**: Ejecutar lógica de negocio (placeholder v1)

### Stage 4: Response Generation
- **Input**: Processed message
- **Output**: Response object with text
- **Responsabilidad**: Generar texto de respuesta

### Stage 5: TTS Generation
- **Input**: Response text
- **Output**: Audio binary/audio URL
- **Responsabilidad**: Convertir texto a audio

### Stage 6: Response Delivery
- **Input**: Audio + metadata
- **Output**: Confirmation from external API
- **Responsabilidad**: Enviar audio al usuario

## 4. Contratos de Stage

### 4.1 Interfaz Base de Stage

```go
package pipeline

import (
    "context"
    "errors"
)

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

// PipelineContext contiene el contexto compartido entre etapas.
type PipelineContext struct {
    // Message es el mensaje interno normalizado
    Message *domain.Message
    
    // RequestID es el ID de correlación de la solicitud
    RequestID string
    
    // TraceID para trazabilidad
    TraceID string
    
    // Response es la respuesta generada
    Response *domain.Response
    
    // Audio es el audio generado
    Audio *domain.Audio
    
    // Errors acumula errores durante el pipeline
    Errors []PipelineError
    
    // Metadata adicional del pipeline
    Metadata map[string]interface{}
}

// PipelineError representa un error en el pipeline.
type PipelineError struct {
    Stage   string
    Err     error
    IsFatal bool
}
```

### 4.2 Stage de Ingestión (Event Ingestion)

```go
type IngestionStage struct {
    verifyToken string
}

func (s *IngestionStage) Name() string { return "ingestion" }

func (s *IngestionStage) CanProcess(input interface{}) bool {
    _, ok := input.(*WebhookRequest)
    return ok
}

func (s *IngestionStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
    req, ok := input.(*WebhookRequest)
    if !ok {
        return nil, ErrInvalidInput
    }
    
    // Validar verification token para GET requests
    if isVerificationRequest(req) {
        return nil, ErrVerificationRequest
    }
    
    // Validar payload
    if err := validateWebhook(req); err != nil {
        return nil, err
    }
    
    // Extraer mensajes
    messages := extractMessages(req)
    if len(messages) == 0 {
        return nil, ErrNoMessages
    }
    
    return &IngestionOutput{
        RawPayload:  req,
        Messages:    messages,
        PhoneNumberID: getPhoneNumberID(req),
    }, nil
}
```

### 4.3 Stage de Normalización (Message Normalization)

```go
type NormalizationStage struct{}

func (s *NormalizationStage) Name() string { return "normalization" }

func (s *NormalizationStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
    ingestionOut, ok := input.(*IngestionOutput)
    if !ok {
        return nil, ErrInvalidInput
    }
    
    // Tomar el primer mensaje (v1)
    waMsg := ingestionOut.Messages[0]
    
    // Convertir a modelo interno
    msg := domain.NewMessage(
        waMsg.ID,
        waMsg.From,
        domain.ChannelWhatsApp,
        mapContentType(waMsg.Type),
        waMsg.Text.Body,
        parseTimestamp(waMsg.Timestamp),
        domain.MessageMetadata{
            UserPhone:       waMsg.From,
            BusinessPhoneID: ingestionOut.PhoneNumberID,
        },
    )
    
    msg.ProcessingStatus = domain.StatusNormalized
    
    return &NormalizationOutput{
        Message: msg,
    }, nil
}
```

### 4.4 Stage de Generación de Respuesta (Response Generation)

```go
type ResponseGenerationStage struct{}

func (s *ResponseGenerationStage) Name() string { return "response_generation" }

func (s *ResponseGenerationStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
    normOut, ok := input.(*NormalizationOutput)
    if !ok {
        return nil, ErrInvalidInput
    }
    
    // Generar respuesta (placeholder v1)
    responseText := "Message received. Generating audio response."
    
    response := domain.NewResponse(
        normOut.Message.UserID,
        domain.ResponseTypeAudio,
        responseText,
        normOut.Message.Metadata.BusinessPhoneID,
    )
    
    return &ResponseGenerationOutput{
        Response: response,
    }, nil
}
```

### 4.5 Stage de Generación TTS (TTS Generation)

```go
type TTSGenerationStage struct {
    ttsProvider TTSProvider
}

func (s *TTSGenerationStage) Name() string { return "tts_generation" }

func (s *TTSGenerationStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
    respGenOut, ok := input.(*ResponseGenerationOutput)
    if !ok {
        return nil, ErrInvalidInput
    }
    
    // Generar audio
    audio, err := s.ttsProvider.Generate(ctx, respGenOut.Response.Text)
    if err != nil {
        return nil, fmt.Errorf("tts generation failed: %w", err)
    }
    
    return &TTSGenerationOutput{
        Audio: audio,
    }, nil
}
```

### 4.6 Stage de Entrega (Delivery)

```go
type DeliveryStage struct {
    deliveryAdapter DeliveryAdapter
}

func (s *DeliveryStage) Name() string { return "delivery" }

func (s *DeliveryStage) Process(ctx context.Context, input interface{}) error {
    ttsOut, ok := input.(*TTSGenerationOutput)
    if !ok {
        return ErrInvalidInput
    }
    
    // Enviar mensaje
    err := s.deliveryAdapter.Send(ctx, ttsOut.Audio)
    if err != nil {
        return fmt.Errorf("delivery failed: %w", err)
    }
    
    return nil
}
```

## 5. Motor del Pipeline

```go
type Pipeline struct {
    stages []Stage
}

func NewPipeline(stages ...Stage) *Pipeline {
    return &Pipeline{stages: stages}
}

func (p *Pipeline) Execute(ctx context.Context, input interface{}) (*PipelineContext, error) {
    pipelineCtx := &PipelineContext{
        RequestID: generateRequestID(),
        TraceID:   generateTraceID(),
        Metadata:  make(map[string]interface{}),
    }
    
    currentInput := input
    
    for _, stage := range p.stages {
        if !stage.CanProcess(currentInput) {
            return nil, fmt.Errorf("stage %s cannot process input", stage.Name())
        }
        
        output, err := stage.Process(ctx, currentInput)
        if err != nil {
            pipelineCtx.Errors = append(pipelineCtx.Errors, PipelineError{
                Stage:   stage.Name(),
                Err:     err,
                IsFatal: true,
            })
            return pipelineCtx, err
        }
        
        currentInput = output
    }
    
    return pipelineCtx, nil
}
```

## 6. Flujo de Ejecución

```
[Webhook Request]
       │
       v
┌──────────────────┐
│  IngestionStage  │ ──> Valida webhook, extrae mensajes
└────────┬─────────┘
         │
         v
┌──────────────────┐
│ NormalizationStage│ ──> Convierte a modelo interno
└────────┬─────────┘
         │
         v
┌──────────────────────┐
│ ResponseGenerationStage│ ──> Genera texto de respuesta
└────────┬─────────────┘
         │
         v
┌─────────────────┐
│ TTSGenerationStage│ ──> Genera audio
└────────┬────────┘
         │
         v
┌───────────────┐
│ DeliveryStage │ ──> Envía audio a WhatsApp
└───────────────┘
```

## 7. Manejo de Errores

- Cada stage puede retornar un error
- Errores fatales detienen el pipeline inmediatamente
- Errores no-fatales se acumulan en `PipelineContext.Errors`
- El pipeline retorna el contexto con errores para debugging