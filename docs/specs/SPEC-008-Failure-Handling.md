# SPEC-008: Failure Handling

## 1. Propósito

Esta especificación define cómo los errores y fallos se propagan a través del pipeline, las estrategias de retry, y las expectativas de logging.

## 2. Categorías de Fallos

### 2.1 Fallos de Entrada
- Payload de webhook malformado
- Token de verificación inválido
- Mensaje de tipo no soportado
- Metadata faltante

### 2.2 Fallos de Procesamiento
- Error en normalización
- Error en generación de respuesta
- Timeout en procesamiento

### 2.3 Fallos de Integración Externa
- TTS service unavailable
- WhatsApp API error
- Timeout de red

### 2.4 Fallos de Entrega
- Usuario no encontrado
- Rate limiting
- Audio inválido

## 3. Tipos de Error

```go
package errors

import "fmt"

// Classification define la clasificación del error.
type Classification string

const (
    ClassificationTransient Classification = "transient"  // Retry puede ayudar
    ClassificationPermanent Classification = "permanent" // No retry
    ClassificationUnknown   Classification = "unknown"    // Determinar automáticamente
)

// ErrorWithClassification envuelve errores con clasificación.
type ErrorWithClassification struct {
    Err           error
    Classification Classification
    Retryable     bool
    LogLevel      string  // "error", "warn", "info"
}

// NewTransientError crea un error transitorio.
func NewTransientError(err error) *ErrorWithClassification {
    return &ErrorWithClassification{
        Err:            err,
        Classification: ClassificationTransient,
        Retryable:      true,
        LogLevel:       "warn",
    }
}

// NewPermanentError crea un error permanente.
func NewPermanentError(err error) *ErrorWithClassification {
    return &ErrorWithClassification{
        Err:            err,
        Classification: ClassificationPermanent,
        Retryable:      false,
        LogLevel:       "error",
    }
}

// Unwrap retorna el error subyacente.
func (e *ErrorWithClassification) Unwrap() error {
    return e.Err
}

func (e *ErrorWithClassification) Error() string {
    return fmt.Sprintf("[%s] %s", e.Classification, e.Err.Error())
}
```

## 4. Errores Predefinidos

```go
package errors

import "errors"

var (
    // ErrInvalidWebhookPayload = payload inválido
    ErrInvalidWebhookPayload = errors.New("invalid webhook payload")
    
    // ErrNoMessages = no hay mensajes en el payload
    ErrNoMessages = errors.New("no messages in payload")
    
    // ErrUnsupportedMessageType = tipo de mensaje no soportado
    ErrUnsupportedMessageType = errors.New("unsupported message type")
    
    // ErrTTSGenerationFailed = generación TTS fallida
    ErrTTSGenerationFailed = errors.New("tts generation failed")
    
    // ErrDeliveryFailed = entrega fallida
    ErrDeliveryFailed = errors.New("delivery failed")
    
    // ErrWhatsAppAPIError = error de API de WhatsApp
    ErrWhatsAppAPIError = errors.New("whatsapp api error")
)
```

## 5. Estrategias de Retry

### 5.1 Retry Config

```go
type RetryConfig struct {
    // MaxRetries número máximo de reintentos
    MaxRetries int
    
    // InitialBackoff tiempo inicial de backoff
    InitialBackoff time.Duration
    
    // MaxBackoff tiempo máximo de backoff
    MaxBackoff time.Duration
    
    // BackoffMultiplier multiplicador de backoff
    BackoffMultiplier float64
    
    // RetryableErrors lista de errores que permiten retry
    RetryableErrors []error
}

var DefaultRetryConfig = RetryConfig{
    MaxRetries:        3,
    InitialBackoff:    1 * time.Second,
    MaxBackoff:        30 * time.Second,
    BackoffMultiplier: 2.0,
}
```

### 5.2 Retry Logic

```go
package retry

import (
    "context"
    "time"
)

func Do(ctx context.Context, config RetryConfig, fn func() error) error {
    var lastErr error
    backoff := config.InitialBackoff
    
    for attempt := 0; attempt <= config.MaxRetries; attempt++ {
        if attempt > 0 {
            select {
            case <-ctx.Done():
                return ctx.Err()
            case <-time.After(backoff):
            }
            backoff = time.Duration(float64(backoff) * config.BackoffMultiplier)
            if backoff > config.MaxBackoff {
                backoff = config.MaxBackoff
            }
        }
        
        if err := fn(); err != nil {
            lastErr = err
            
            // Verificar si es retryable
            if !isRetryable(err, config.RetryableErrors) {
                return err
            }
            
            continue
        }
        
        return nil
    }
    
    return lastErr
}

func isRetryable(err error, retryableErrors []error) bool {
    for _, r := range retryableErrors {
        if errors.Is(err, r) {
            return true
        }
    }
    return false
}
```

### 5.3 Aplicación en Stages

```go
// Retry config para TTS (más reintentos, mayor backoff)
var TTSRetryConfig = errors.RetryConfig{
    MaxRetries:        5,
    InitialBackoff:    2 * time.Second,
    MaxBackoff:        60 * time.Second,
    RetryableErrors:   []error{errors.ErrTTSGenerationFailed},
}

// Retry config para Delivery (menos reintentos)
var DeliveryRetryConfig = errors.RetryConfig{
    MaxRetries:        3,
    InitialBackoff:    1 * time.Second,
    MaxBackoff:        10 * time.Second,
    RetryableErrors:   []error{errors.ErrDeliveryFailed},
}
```

## 6. Manejo en Pipeline

```go
type Pipeline struct {
    stages   []Stage
    retryCfg map[string]RetryConfig  // Config por stage
}

func (p *Pipeline) Execute(ctx context.Context, input interface{}) (*PipelineContext, error) {
    pc := newPipelineContext()
    currentInput := input
    
    for _, stage := range p.stages {
        // Obtener config de retry para este stage
        retryCfg := p.retryCfg[stage.Name()]
        
        err := retry.Do(ctx, retryCfg, func() error {
            out, err := stage.Process(ctx, currentInput)
            if err != nil {
                return err
            }
            currentInput = out
            return nil
        })
        
        if err != nil {
            pc.Errors = append(pc.Errors, PipelineError{
                Stage:   stage.Name(),
                Err:     err,
                IsFatal: !isRetryable(err),
            })
            
            // Si es error fatal, detener pipeline
            if !isRetryable(err) {
                return pc, err
            }
        }
    }
    
    return pc, nil
}
```

## 7. Logging de Errores

```go
package logging

func LogPipelineError(err *PipelineError, logger Logger) {
    level := err.Err.LogLevel()
    
    logger.Log(level, "pipeline_error",
        "stage", err.Stage,
        "error", err.Err.Error(),
        "classification", err.Err.Classification(),
        "retryable", err.Err.Retryable(),
    )
}
```

## 8. Circuit Breaker (Futuro)

Para sistemas de alto tráfico, considerar implementar circuit breaker:

- After N consecutive failures, open circuit
- After timeout, allow one request through (half-open)
- If success, close circuit; if fail, reopen