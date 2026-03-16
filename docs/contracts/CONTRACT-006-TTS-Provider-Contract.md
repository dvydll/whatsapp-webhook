# CONTRACT-006: TTS Provider Contract

## 1. TTSProvider Interface

```go
package tts

import (
    "context"
    "github.com/whatsapp-tts/internal/domain"
)

// TTSProvider define la interfaz para proveedores TTS.
type TTSProvider interface {
    // Generate genera audio a partir de texto.
    Generate(ctx context.Context, text string) (*domain.AudioAsset, error)
    
    // GenerateWithOptions genera audio con opciones específicas.
    GenerateWithOptions(ctx context.Context, text string, opts Options) (*domain.AudioAsset, error)
    
    // Name retorna el nombre del proveedor.
    Name() string
}
```

## 2. Options

```go
package tts

// Options contiene opciones para la generación de audio.
type Options struct {
    // Voice es el identificador de voz.
    Voice string
    
    // Language es el código de idioma (e.g., "en-US", "es-ES").
    Language string
    
    // Speed es la velocidad de habla (0.5 - 2.0).
    Speed float64
    
    // Format es el formato de audio deseado.
    Format domain.AudioFormat
    
    // Quality es la calidad de audio (low, medium, high).
    Quality string
}
```

## 3. TTS Result

```go
package tts

import "github.com/whatsapp-tts/internal/domain"

// Result representa el resultado de la generación TTS.
type Result struct {
    // Audio es el asset de audio generado.
    Audio *domain.AudioAsset
    
    // Duration es la duración en segundos.
    Duration float64
    
    // CharactersUsed es el número de caracteres convertidos.
    CharactersUsed int
    
    // Cost es el costo de la operación (si aplica).
    Cost float64
}
```

## 4. Factory

```go
package tts

import "errors"

var (
    ErrUnknownProvider = errors.New("unknown tts provider")
    ErrGenerationFailed = errors.New("tts generation failed")
)

// ProviderType define los tipos de proveedor disponibles.
type ProviderType string

const (
    ProviderStyleTTS ProviderType = "styletts"
    ProviderOpenAI   ProviderType = "openai"
    ProviderGoogle   ProviderType = "google"
    ProviderLocal    ProviderType = "local"
)

// Config contiene la configuración para crear un provider.
type Config struct {
    Type            ProviderType
    BaseURL         string
    APIKey          string
    DefaultVoice    string
    DefaultLanguage string
    DefaultFormat   domain.AudioFormat
}

// NewProvider crea un provider según la configuración.
func NewProvider(config Config) (TTSProvider, error) {
    switch config.Type {
    case ProviderStyleTTS:
        return newStyleTTSProvider(config)
    case ProviderOpenAI:
        return newOpenAIProvider(config)
    case ProviderGoogle:
        return newGoogleProvider(config)
    case ProviderLocal:
        return newLocalProvider(config)
    default:
        return nil, ErrUnknownProvider
    }
}
```

## 5. Interfaz Extendida para Features Futuros

```go
package tts

import "io"

// StreamingTTSProvider permite generación de audio en streaming.
type StreamingTTSProvider interface {
    TTSProvider
    
    // GenerateStreaming genera audio en streaming.
    // El callback es llamado con chunks de audio.
    GenerateStreaming(ctx context.Context, text string, opts Options, cb func([]byte)) error
}

// BatchTTSProvider permite procesamiento por lotes.
type BatchTTSProvider interface {
    TTSProvider
    
    // GenerateBatch genera audio para múltiples textos.
    GenerateBatch(ctx context.Context, texts []string, opts Options) ([]Result, error)
}
```

## 6. Notas

- La interfaz permite swap de proveedores en tiempo de ejecución
- Options permite configuración flexible por llamada
- Las interfaces extendidas (Streaming, Batch) son opcionales
- El dominio no depende de implementaciones específicas