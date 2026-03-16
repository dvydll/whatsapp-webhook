# CONTRACT-007: Audio Processing Contract

## 1. AudioProcessor Interface

```go
package audio

import (
    "context"
    "github.com/whatsapp-tts/internal/domain"
)

// AudioProcessor procesa audio para hacerlo compatible con WhatsApp.
type AudioProcessor interface {
    // ConvertToWhatsAppFormat convierte el audio al formato de WhatsApp.
    ConvertToWhatsAppFormat(ctx context.Context, audio *domain.AudioAsset) (*domain.AudioAsset, error)
    
    // NormalizeVolume normaliza el volumen del audio.
    NormalizeVolume(ctx context.Context, audio *domain.AudioAsset) (*domain.AudioAsset, error)
    
    // TrimSilence elimina silencio al inicio y final.
    TrimSilence(ctx context.Context, audio *domain.AudioAsset) (*domain.AudioAsset, error)
}
```

## 2. FFmpegProcessor (Implementación de referencia)

```go
package audio

import (
    "context"
    "os/exec"
)

// FFmpegProcessor implementa AudioProcessor usando FFmpeg.
type FFmpegProcessor struct {
    ffmpegPath string
}

// NewFFmpegProcessor crea un nuevo procesador FFmpeg.
func NewFFmpegProcessor(ffmpegPath string) *FFmpegProcessor {
    if ffmpegPath == "" {
        ffmpegPath = "ffmpeg"
    }
    return &FFmpegProcessor{
        ffmpegPath: ffmpegPath,
    }
}

func (p *FFmpegProcessor) ConvertToWhatsAppFormat(ctx context.Context, audio *domain.AudioAsset) (*domain.AudioAsset, error) {
    // Convertir a opus/aac, 48kHz, mono
    // Implementación usa exec.Command con ffmpeg
    return p.convert(ctx, audio, "libopus", 48000, 1)
}

func (p *FFmpegProcessor) NormalizeVolume(ctx context.Context, audio *domain.AudioAsset) (*domain.AudioAsset, error) {
    // Normalizar a -16 LUFS
    return audio, nil
}

func (p *FFmpegProcessor) TrimSilence(ctx context.Context, audio *domain.AudioAsset) (*domain.AudioAsset, error) {
    // Eliminar silencios menores a 0.5s
    return audio, nil
}
```

## 3. Audio Converter Interface

```go
package audio

import "github.com/whatsapp-tts/internal/domain"

// AudioConverter convierte entre formatos de audio.
type AudioConverter interface {
    // Convert convierte audio a un formato específico.
    Convert(ctx context.Context, input *domain.AudioAsset, outputFormat domain.AudioFormat) (*domain.AudioAsset, error)
    
    // SupportedFormats retorna los formatos soportados.
    SupportedFormats() []domain.AudioFormat
}
```

## 4. Audio Validator

```go
package audio

import "github.com/whatsapp-tts/internal/domain"

// AudioValidator valida que el audio sea aceptable.
type AudioValidator interface {
    // Validate valida el audio.
    Validate(ctx context.Context, audio *domain.AudioAsset) ValidationResult
    
    // IsWhatsAppCompatible verifica compatibilidad con WhatsApp.
    IsWhatsAppCompatible(audio *domain.AudioAsset) bool
}

type ValidationResult struct {
    Valid  bool
    Errors []string
    Warnings []string
}

// WhatsAppAudioRequirements define los requisitos de WhatsApp.
var WhatsAppAudioRequirements = struct {
    MaxDuration  float64  // 16 segundos max
    MaxSizeBytes int64   // 16MB max
    MinBitrate   int
    MaxBitrate   int
    AllowedCodecs []domain.AudioFormat
}{
    MaxDuration:  16.0,
    MaxSizeBytes: 16 * 1024 * 1024,
    MinBitrate:   32000,
    MaxBitrate:   128000,
    AllowedCodecs: []domain.AudioFormat{domain.FormatAAC},
}
```

## 5. Audio Pipeline

```go
package audio

import "context"

// ProcessorPipeline define un pipeline de procesamiento de audio.
type ProcessorPipeline struct {
    processors []AudioProcessor
}

// AddProcessor agrega un procesador al pipeline.
func (p *ProcessorPipeline) AddProcessor(processor AudioProcessor) {
    p.processors = append(p.processors, processor)
}

// Process ejecuta todos los procesadores en orden.
func (p *ProcessorPipeline) Process(ctx context.Context, audio *domain.AudioAsset) (*domain.AudioAsset, error) {
    current := audio
    for _, proc := range p.processors {
        result, err := proc.ConvertToWhatsAppFormat(ctx, current)
        if err != nil {
            return nil, err
        }
        current = result
    }
    return current, nil
}
```

## 6. Factory

```go
package audio

import "context"

// Config contiene configuración del procesador.
type Config struct {
    UseFFmpeg       bool
    FFmpegPath      string
    TargetFormat    domain.AudioFormat
    TargetSampleRate int
    TargetChannels  int
}

// NewProcessor crea un procesador según la configuración.
func NewProcessor(config Config) (AudioProcessor, error) {
    if config.UseFFmpeg {
        return NewFFmpegProcessor(config.FFmpegPath), nil
    }
    // Fallback a implementación nativa o error
    return nil, ErrFFmpegRequired
}

var ErrFFmpegRequired = &ConfigError{"ffmpeg is required for audio processing"}
```

## 7. Notas

- Abstrae el uso de herramientas como FFmpeg
- Permite reemplazo de implementación
- Incluye validación de requisitos de WhatsApp
- El pipeline permite encadenar procesadores