# SPEC-006: Audio Generation Contract

## 1. Propósito

El contrato de generación de audio define la interfaz para proveedores Text-to-Speech (TTS). Permite intercambiar proveedores sin modificar el código del pipeline.

## 2. Interfaz TTSProvider

```go
package tts

import (
    "context"
    "errors"
)

var ErrTTSGenerationFailed = errors.New("tts generation failed")

type TTSProvider interface {
    Generate(ctx context.Context, text string) (*Audio, error)
    GenerateWithOptions(ctx context.Context, text string, opts Options) (*Audio, error)
    Name() string
}

type Options struct {
    Voice    string
    Language string
    Speed    float64
    Format   AudioFormat
}

type AudioFormat string

const (
    FormatAAC AudioFormat = "aac"
    FormatMP3 AudioFormat = "mp3"
    FormatWAV AudioFormat = "wav"
    FormatOGG AudioFormat = "ogg"
)
```

## 3. Modelo Audio

```go
package domain

type Audio struct {
    AudioID    string
    Format     AudioFormat
    Reference  string  // ID externo o URL
    Data       []byte  // Audio binary
    URL        string
    Duration   float64
    GeneratedAt time.Time
}
```

## 4. Implementación StyleTTS

```go
package styletts

type StyleTTSProvider struct {
    baseURL      string
    apiKey       string
    defaultVoice string
    defaultLang  string
}

func (p *StyleTTSProvider) Generate(ctx context.Context, text string) (*Audio, error) {
    return p.GenerateWithOptions(ctx, text, Options{
        Voice:   p.defaultVoice,
        Language: p.defaultLang,
        Speed:    1.0,
        Format:  FormatAAC,
    })
}

func (p *StyleTTSProvider) GenerateWithOptions(ctx context.Context, text string, opts Options) (*Audio, error) {
    // HTTP call a StyleTTS API
    resp, err := p.callAPI(ctx, text, opts)
    if err != nil {
        return nil, fmt.Errorf("styletts api call failed: %w", err)
    }
    
    return &Audio{
        AudioID:   resp.AudioID,
        Format:    opts.Format,
        Reference: resp.AudioID,
        URL:       resp.DownloadURL,
        Duration:  resp.Duration,
    }, nil
}
```

## 5. Factory

```go
package tts

func NewProvider(config Config) (TTSProvider, error) {
    switch config.Type {
    case "styletts":
        return styletts.NewStyleTTSProvider(config.BaseURL, config.APIKey, config.DefaultVoice, config.DefaultLanguage), nil
    case "google":
        return google.NewGoogleTTSProvider(config.APIKey), nil
    default:
        return nil, ErrUnknownTTSProvider
    }
}

type Config struct {
    Type            string
    BaseURL         string
    APIKey          string
    DefaultVoice    string
    DefaultLanguage string
}
```

## 6. Notas de Diseño

- La interfaz permite swap de proveedores en tiempo de ejecución
- Reference puede ser un ID externo de WhatsApp (después de upload)
- Las opciones se injectan via constructor