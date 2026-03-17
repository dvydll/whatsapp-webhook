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

## 3. Modelo de Errores

El modelo de errores está definido en CONTRACT-008-Error-Model.md. Incluye:

- **AppError**: Estructura principal con código, mensaje, stage, retryable, details
- **ErrorCode**: Códigos únicos para cada tipo de error
- **PipelineStage**: Etapa donde ocurrió el error
- **ErrorList**: Para acumular múltiples errores

## 4. Estrategias de Retry

### 4.1 Retry Config

```go
type RetryConfig struct {
    MaxRetries        int
    InitialBackoff    time.Duration
    MaxBackoff        time.Duration
    BackoffMultiplier float64
    RetryableErrors   []error
}

var DefaultRetryConfig = RetryConfig{
    MaxRetries:        3,
    InitialBackoff:    1 * time.Second,
    MaxBackoff:        30 * time.Second,
    BackoffMultiplier: 2.0,
}
```

### 4.2 Retry Logic

```go
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
            
            if !isRetryable(err, config.RetryableErrors) {
                return err
            }
            
            continue
        }
        
        return nil
    }
    
    return lastErr
}
```

### 4.3 Aplicación en Stages

```go
var TTSRetryConfig = RetryConfig{
    MaxRetries:        5,
    InitialBackoff:    2 * time.Second,
    MaxBackoff:        60 * time.Second,
}

var DeliveryRetryConfig = RetryConfig{
    MaxRetries:        3,
    InitialBackoff:    1 * time.Second,
    MaxBackoff:        10 * time.Second,
}
```

## 5. Manejo en Pipeline

```go
func (p *Pipeline) Execute(ctx context.Context, input interface{}) (*PipelineContext, error) {
    pc := &PipelineContext{
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
            pc.Errors = append(pc.Errors, PipelineError{
                Stage:   stage.Name(),
                Err:     err,
                IsFatal: true,
            })
            return pc, err
        }
        
        currentInput = output
    }
    
    return pc, nil
}
```

## 6. Logging de Errores

```go
func LogPipelineError(err PipelineError, logger Logger) {
    logger.Error("pipeline_error",
        "stage", err.Stage,
        "error", err.Err.Error(),
        "is_fatal", err.IsFatal,
    )
}
```

## 7. Circuit Breaker (Futuro)

Para sistemas de alto tráfico, considerar implementar circuit breaker:

- After N consecutive failures, open circuit
- After timeout, allow one request through (half-open)
- If success, close circuit; if fail, reopen

## 8. Notas

- El modelo de errores detallado está en CONTRACT-008-Error-Model.md
- Los códigos de error definidos en CONTRACT-008 deben usarse consistentemente
- La clasificación retryable vs permanent guía la estrategia de retry
