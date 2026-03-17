# SPEC-006: Audio Generation Contract

## 1. Propósito

El contrato de generación de audio define la interfaz para proveedores Text-to-Speech (TTS). Permite intercambiar proveedores sin modificar el código del pipeline.

## 2. Interfaz TTSProvider

```go
package tts

import (
    "context"
    "errors"
    "github.com/whatsapp-tts/internal/domain"
)

var ErrTTSGenerationFailed = errors.New("tts generation failed")

type TTSProvider interface {
    Generate(ctx context.Context, text string) (*domain.AudioAsset, error)
    GenerateWithOptions(ctx context.Context, text string, opts Options) (*domain.AudioAsset, error)
    Name() string
}

type Options struct {
    Voice    string
    Language string
    Speed    float64
    Format   domain.AudioFormat
    Quality  string
}
```

## 3. Factory

```go
package tts

import "errors"

var (
    ErrUnknownTTSProvider = errors.New("unknown tts provider")
)

type ProviderType string

const (
    ProviderStyleTTS ProviderType = "styletts"
    ProviderOpenAI   ProviderType = "openai"
    ProviderGoogle   ProviderType = "google"
    ProviderLocal    ProviderType = "local"
)

type Config struct {
    Type            ProviderType
    BaseURL         string
    APIKey          string
    DefaultVoice    string
    DefaultLanguage string
    DefaultFormat   domain.AudioFormat
}

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
        return nil, ErrUnknownTTSProvider
    }
}
```

## 4. Interfaz Extendida para Features Futuros

```go
package tts

import "io"

// StreamingTTSProvider permite generación de audio en streaming.
type StreamingTTSProvider interface {
    TTSProvider
    
    // GenerateStreaming genera audio en streaming.
    GenerateStreaming(ctx context.Context, text string, opts Options, cb func([]byte)) error
}

// BatchTTSProvider permite procesamiento por lotes.
type BatchTTSProvider interface {
    TTSProvider
    
    // GenerateBatch genera audio para múltiples textos.
    GenerateBatch(ctx context.Context, texts []string, opts Options) ([]Result, error)
}

type Result struct {
    Audio          *domain.AudioAsset
    Duration       float64
    CharactersUsed int
    Cost           float64
}
```

## 5. Notas de Diseño

- La interfaz permite swap de proveedores en tiempo de ejecución
- El dominio no depende de implementaciones específicas
- Las opciones se injectan via constructor
- Para implementación real ver CONTRACT-006-TTS-Provider-Contract.md