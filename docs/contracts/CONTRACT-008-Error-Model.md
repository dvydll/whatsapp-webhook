# CONTRACT-008: Error Model

## 1. Error Structure

```go
package domain

import "fmt"

// AppError representa un error estructurado de la aplicación.
type AppError struct {
    // Code es el código de error único.
    Code ErrorCode
    
    // Message es la descripción legible del error.
    Message string
    
    // Stage indica en qué etapa del pipeline ocurrió.
    Stage PipelineStage
    
    // Retryable indica si el error puede reintentarse.
    Retryable bool
    
    // IsInternal indica si es un error interno o externo.
    IsInternal bool
    
    // Details contiene información adicional.
    Details map[string]interface{}
    
    // Cause es el error original que causó este error.
    Cause error
}
```

## 2. Error Codes

```go
package domain

// ErrorCode define códigos de error únicos.
type ErrorCode string

const (
    // Errores de entrada
    ErrCodeInvalidPayload     ErrorCode = "INVALID_PAYLOAD"
    ErrCodeVerificationFailed ErrorCode = "VERIFICATION_FAILED"
    ErrCodeUnsupportedType   ErrorCode = "UNSUPPORTED_TYPE"
    
    // Errores de normalización
    ErrCodeNormalizationFailed ErrorCode = "NORMALIZATION_FAILED"
    ErrCodeMissingField        ErrorCode = "MISSING_FIELD"
    
    // Errores de procesamiento
    ErrCodeProcessingFailed    ErrorCode = "PROCESSING_FAILED"
    ErrCodeResponseGeneration ErrorCode = "RESPONSE_GENERATION_FAILED"
    
    // Errores de TTS
    ErrCodeTTSFailed         ErrorCode = "TTS_FAILED"
    ErrCodeTTSUnavailable    ErrorCode = "TTS_UNAVAILABLE"
    
    // Errores de audio processing
    ErrCodeAudioProcessing    ErrorCode = "AUDIO_PROCESSING_FAILED"
    ErrCodeInvalidAudio      ErrorCode = "INVALID_AUDIO"
    
    // Errores de entrega
    ErrCodeDeliveryFailed    ErrorCode = "DELIVERY_FAILED"
    ErrCodeWhatsAppAPIError  ErrorCode = "WHATSAPP_API_ERROR"
    ErrCodeRateLimited       ErrorCode = "RATE_LIMITED"
    
    // Errores internos
    ErrCodeInternalError    ErrorCode = "INTERNAL_ERROR"
    ErrCodeConfigError      ErrorCode = "CONFIG_ERROR"
)

func (e ErrorCode) String() string { return string(e) }
```

## 3. Pipeline Stage Enum

```go
package domain

// PipelineStage indica la etapa del pipeline.
type PipelineStage string

const (
    StageIngestion        PipelineStage = "ingestion"
    StageNormalization    PipelineStage = "normalization"
    StageProcessing       PipelineStage = "processing"
    StageResponseGen      PipelineStage = "response_generation"
    StageTTS              PipelineStage = "tts_generation"
    StageAudioProcessing  PipelineStage = "audio_processing"
    StageDelivery         PipelineStage = "delivery"
)
```

## 4. Error Construction

```go
package domain

import "fmt"

// NewError crea un nuevo AppError.
func NewError(code ErrorCode, message string, stage PipelineStage, retryable bool) *AppError {
    return &AppError{
        Code:      code,
        Message:   message,
        Stage:     stage,
        Retryable: retryable,
        IsInternal: true,
    }
}

// WrapError envuelve un error existente.
func WrapError(err error, code ErrorCode, message string, stage PipelineStage, retryable bool) *AppError {
    return &AppError{
        Code:      code,
        Message:   message,
        Stage:     stage,
        Retryable: retryable,
        IsInternal: false,
        Cause:     err,
    }
}

func (e *AppError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Cause
}
```

## 5. Error List

```go
package domain

// ErrorList contiene múltiples errores.
type ErrorList struct {
    Errors []*AppError
}

func (el *ErrorList) Add(err *AppError) {
    el.Errors = append(el.Errors, err)
}

func (el *ErrorList) HasErrors() bool {
    return len(el.Errors) > 0
}

func (el *ErrorList) HasFatal() bool {
    for _, e := range el.Errors {
        if !e.Retryable {
            return true
        }
    }
    return false
}

func (el *ErrorList) Error() string {
    if len(el.Errors) == 0 {
        return ""
    }
    msg := "multiple errors: "
    for i, e := range el.Errors {
        if i > 0 {
            msg += ", "
        }
        msg += e.Error()
    }
    return msg
}
```

## 6. Predefined Errors

```go
package domain

var (
    ErrInvalidPayload = NewError(
        ErrCodeInvalidPayload,
        "invalid webhook payload",
        StageIngestion,
        false,
    )
    
    ErrTTSUnavailable = NewError(
        ErrCodeTTSUnavailable,
        "TTS service unavailable",
        StageTTS,
        true,
    )
    
    ErrDeliveryFailed = NewError(
        ErrCodeDeliveryFailed,
        "failed to deliver message",
        StageDelivery,
        true,
    )
)
```

## 7. Error Handler Interface

```go
package handler

import "github.com/whatsapp-tts/internal/domain"

// ErrorHandler maneja errores del pipeline.
type ErrorHandler interface {
    // HandleError maneja un error y decide la respuesta.
    HandleError(ctx context.Context, err *domain.AppError) (response interface{}, shouldStop bool)
    
    // LogError registra el error.
    LogError(err *domain.AppError)
}
```

## 8. Notas

- Estructura consistente para todos los errores
- Clasificación clara de retryable vs permanent
- Trazabilidad de etapa donde ocurrió el error
- Wrapping de errores originales para debugging