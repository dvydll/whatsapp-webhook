## Why

El sistema actual tiene stubs en TTS y Delivery que no generan audio real. Se necesita implementar la integración completa con un proveedor TTS y el pipeline de audio.

## What Changes

- Implementar TTSProvider interface y StyleTTS
- Implementar audio processing (conversión a formato WhatsApp)
- Conectar stages al pipeline real
- Agregar métricas y health check

## Capabilities

### New Capabilities
- tts-provider: Integración con proveedor TTS real
- audio-processing: Conversión y validación de audio

### Modified Capabilities
- (ninguno)
