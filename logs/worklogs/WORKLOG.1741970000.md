---
title: "004-Minimal-Flow - Flujo Funcional WhatsApp Completo"
summary: "Implementación del flujo end-to-end: webhook handler, cliente WhatsApp real, parsing de mensajes y envío de respuestas verificado"
description: "Cuarta épica del proyecto. Se implementó el flujo funcional completo: cliente WhatsApp real (internal/adapters/whatsapp/client.go) para enviar mensajes, webhook handler con verificación y parsing, integración del pipeline. Se agregó endpoint /test-webhook para simular mensajes debido a limitación de Meta (apps en modo desarrollo no reciben webhooks de usuarios reales)."
createdAt: "2025-03-14T20:15:00Z"
tags:
  - whatsapp
  - tts
  - webhook
  - integration
  - ngrok
  - meta-limitations
  - working
metadata:
  epic: "004-minimal-flow"
  status: "completed"
  repository: "whatsapp-tts"
  language: "es-ES"
---

# Épica 004: Flujo Minimal Funcional

## Objetivo

Implementar el primer versión funcional real del sistema que pueda completar el flujo completo con WhatsApp.

## Componentes Implementados

### Cliente WhatsApp (internal/adapters/whatsapp/client.go)
- Envío de mensajes de texto (SendTextMessage)
- Envío de mensajes de audio (SendAudioMessage)
- Manejo de errores de API
- Tests unitarios pasando (TestSendTextMessage, TestSendAudioMessage, TestSendMessageAPIError)

### Handler de Webhook (internal/webhook/handler.go)
- Verificación de webhook (GET /webhook?hub.mode=subscribe&hub.verify_token=...)
- Parsing de payloads JSON de WhatsApp
- Procesamiento de mensajes entrantes
- Integración con cliente de WhatsApp
- Endpoint de prueba (/test-webhook) para simulación

## Problemas Encontrados y Soluciones

### Problema 1: Puerto en uso
**Descripción:** Al intentar iniciar el servidor, el puerto 8080 estaba ocupado por procesos anteriores.
**Solución:** Identificar y matar procesos con `lsof -ti :8080`

### Problema 2: Servidor sin URL pública
**Descripción:** WhatsApp Cloud API necesita URL pública para enviar webhooks.
**Solución:** Usar ngrok para exponer localhost públicamente.
**Comando:** `ngrok http 8080`

### Problema 3: Webhook no llega (App en modo desarrollo)
**Descripción:** Las apps de WhatsApp en modo desarrollo no reciben webhooks de usuarios reales.
**Mensaje de Meta:** "Las aplicaciones solo podrán recibir webhooks de prueba enviados desde el panel de la aplicación mientras esta no esté publicada."
**Solución temporal:** Crear endpoint /test-webhook que simula mensajes entrantes para verificar el flujo.

## Verificación Exitosa

```
test_webhook_triggered
message_received from=34685107027 type=text
user_message: "Test message from simulation"
sending_response: "Message received. Generating audio response."
response_sent whatsapp_message_id=wamid.HBgLMzQ2ODUxMDcwMjcVAgARGBI2RDZCN0FDQTZDM0ZCMEI2RTQA
```

El mensaje de respuesta fue recibido en WhatsApp ✓

## Configuración Utilizada

```
PHONE_NUMBER_ID=1052703781258489
META_ACCESS_TOKEN=EAAnZAnWGHZA5UBQZCrjztEZChk3csjZAZCYqE6il0MPkn9ovFrPLCNroEHUJHfcbaAA1GW7hKBEwfE7lwkBQaoGOugm9hfO57gkZCjAQ0fkDH358sJVjo5mZBQN0Av4jikZCJIfI2164yUA4yjqgt7Pwsh7Iw1gONMrjgvhycncsZB12ou6rkUY4nGpMa0pyV4RQ0GZAJgOOBu2qKZA3yTadAsCjb24kZCzO3E8rhZCkU6UbZANhHApH6V2ZA5SsYZCg3Y5GRFmWtL722jbaevYRY7o1hXOKXINUHdtNysPbiCwZDZD
WHATSAPP_VERIFY_TOKEN=my_verify_token
```

## Estado

- Sistema funcional para el flujo básico WhatsApp → respuesta
- Verificación completada exitosamente
- Listo para siguientes iteraciones (TTS, audio)

## Limitaciones Actuales

1. App en modo desarrollo → no recibe webhooks de usuarios reales
2. Solo mensajes de texto → no hay integración TTS aún
3. Necesita ngrok para desarrollo local

## Siguiente

Implementar integración con servicio TTS para generar respuestas de audio.