---
title: "002-Contracts - Contratos Técnicos Go"
summary: "Creación de contratos técnicos para implementación en Go: interfaces, modelos de dominio, DTOs, contratos de pipeline y adapters"
description: "Segunda épica del proyecto. Se definieron contratos técnicos que permiten implementación determinística: Go module definition, repository layout, domain models (UserMessage, ResponseMessage, AudioAsset), pipeline contracts, WhatsApp adapter contract, TTS provider contract, audio processing contract, error model, contract tests y example E2E flow."
createdAt: "2025-03-14T15:30:00Z"
tags:
  - whatsapp
  - tts
  - contracts
  - golang
  - interfaces
  - domain-models
metadata:
  epic: "002-contracts"
  status: "completed"
  repository: "whatsapp-tts"
  language: "es-ES"
---

# Épica 002: Contratos Técnicos Go

## Objetivo

Convertir las especificaciones arquitecturales en contratos técnicos concretos para implementación Go.

## Entregables

| # | Archivo | Descripción |
|---|---------|-------------|
| 1 | CONTRACT-001-Go-Module-Definition.md | Módulo Go 1.21, política de dependencias mínimas |
| 2 | CONTRACT-002-Repository-Layout.md | Estructura de directorios con responsabilidades |
| 3 | CONTRACT-003-Domain-Models.md | Modelos puros: UserMessage, ResponseMessage, AudioAsset |
| 4 | CONTRACT-004-Pipeline-Contracts.md | Interfaces Stage, Pipeline, contratos específicos |
| 5 | CONTRACT-005-WhatsApp-Adapter-Contract.md | Interfaz WhatsAppAdapter, DTOs |
| 6 | CONTRACT-006-TTS-Provider-Contract.md | Interfaz TTSProvider, factory |
| 7 | CONTRACT-007-Audio-Processing-Contract.md | AudioProcessor, validación WhatsApp |
| 8 | CONTRACT-008-Error-Model.md | AppError estructurado, códigos, stages |
| 9 | CONTRACT-009-Contract-Tests.md | Estructura de tests de compliance |
| 10 | CONTRACT-010-Example-E2E-Flow.md | Pseudo-código del flujo completo |

## Notas de Diseño

- Todos los modelos de dominio son puros Go (sin dependencias externas)
- Las interfaces permiten reemplazo de implementaciones (DI)
- Los DTOs mapean directamente al formato externo (WhatsApp API)
- Los contracts tests verifican compliance de interfaces

## Estado

Contratos listos para implementación. El sistema está listo para la épica 003 (boilerplate).

## Siguiente

Proceder a épica 003 (boilerplate) donde se implementó código real.