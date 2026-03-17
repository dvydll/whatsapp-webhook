## Context

El sistema actual procesa mensajes pero con stubs que no generan audio real. El pipeline está definido pero los stages TTS y Delivery no funcionan.

## Goals / Non-Goals

**Goals:**
- Implementar TTS real
- Hacer que el pipeline procese audio de verdad
- Enviar audio a WhatsApp

**Non-Goals:**
- Cambiar arquitectura
- Agregar base de datos

## Decisions

- Usar interfaz TTSProvider para poder swap de proveedores
- FFmpeg para procesamiento de audio
- Upload a WhatsApp o usar URL directa

## Risks

- TTS API puede no estar disponible
- Audio puede exceder límites de WhatsApp
