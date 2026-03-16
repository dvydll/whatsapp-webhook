---
title: "003-Boilerplate - Skeleton y Estructura del Proyecto"
summary: "Creación del esqueleto completo del proyecto Go con estructura, interfaces stub, pipeline y servidor HTTP funcional"
description: "Tercera épica del proyecto. Se creó el boilerplate completo: go.mod, estructura de directorios, interfaces de pipeline, stages stub (ingestion, normalization, response_generation, tts_generation, delivery), handler HTTP, configuración, logger y aplicación principal. Todo compila y ejecuta correctamente."
createdAt: "2025-03-14T15:45:00Z"
tags:
  - whatsapp
  - tts
  - boilerplate
  - golang
  - skeleton
  - pipeline
metadata:
  epic: "003-boilerplate"
  status: "completed"
  repository: "whatsapp-tts"
  language: "es-ES"
---

# Épica 003: Boilerplate y Skeleton

## Objetivo

Crear un esqueleto de proyecto Go compilable y ejecutable que satisfaga los contratos definidos.

## Estructura Creada

```
github.com/whatsapp-tts/
├── go.mod                          # Módulo Go 1.22
├── cmd/server/
│   └── main.go                    # Entry point
├── internal/
│   ├── app/
│   │   └── app.go                 # Aplicación principal
│   ├── pipeline/
│   │   ├── pipeline.go            # Implementación del pipeline
│   │   ├── stage.go              # Interfaz Stage
│   │   ├── pipeline_test.go       # Tests del pipeline
│   │   └── stages/
│   │       ├── ingestion.go      # Stage de ingestión (stub)
│   │       ├── normalization.go   # Stage de normalización (stub)
│   │       ├── response_generation.go  # Stage de respuesta (stub)
│   │       ├── tts_generation.go  # Stage TTS (stub)
│   │       └── delivery.go        # Stage de entrega (stub)
│   ├── webhook/
│   │   └── handler.go            # Handler HTTP
│   ├── domain/
│   │   ├── message.go            # Modelo UserMessage
│   │   ├── response.go            # Modelo ResponseMessage
│   │   ├── audio.go               # Modelo AudioAsset
│   │   ├── types.go              # Tipos auxiliares
│   │   └── message_test.go        # Tests de dominio
│   ├── config/
│   │   └── config.go              # Configuración
│   └── observability/
│       └── logger.go              # Logger estructurado
```

## Tests TDD Creados

- TestNewUserMessage
- TestNewResponseMessage  
- TestNewAudioAsset
- TestChannelValues
- TestAudioFormatValues
- TestPipelineExecution
- TestPipelineStageError
- TestPipelineStages

Todos los tests pasan: `go test ./...` ✓

## Verificaciones

- ✓ Compila: `go build ./cmd/server/`
- ✓ Genera ejecutable binario
- ✓ Expone endpoint `/webhook`
- ✓ Logs los pasos del pipeline

## Decisiones de Diseño

- Interface de Logger en paquete observability para permitir reemplazo
- Pipeline con stages intercambiables usando interfaz Stage
- Stages stub que retornan datos predecibles para testing
- Configuración desde variables de entorno

## Estado

El proyecto está listo para la épica 004 donde se implementó el flujo real con WhatsApp.

## Siguiente

Épica 004: Implementación del flujo real de WhatsApp.