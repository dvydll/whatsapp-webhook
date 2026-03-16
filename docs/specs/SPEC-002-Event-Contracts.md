# SPEC-002: Event Contracts

## 1. Webhook Event Schema (Incoming)

El sistema recibe eventos vía HTTP POST desde WhatsApp Cloud API.

### 1.1 Estructura General del Webhook

```json
{
  "object": "whatsapp_business_account",
  "entry": [
    {
      "id": "<WABA_ID>",
      "changes": [
        {
          "value": { ... },
          "field": "messages"
        }
      ]
    }
  ]
}
```

### 1.2 Payload de Mensaje de Texto (v1)

```json
{
  "object": "whatsapp_business_account",
  "entry": [
    {
      "id": "8856996819413533",
      "changes": [
        {
          "value": {
            "messaging_product": "whatsapp",
            "metadata": {
              "display_phone_number": "+1234567890",
              "phone_number_id": "27681414235104944"
            },
            "contacts": [
              {
                "profile": {
                  "name": "John Doe"
                },
                "wa_id": "1234567890"
              }
            ],
            "messages": [
              {
                "from": "16315551234",
                "id": "wamid.ABGGFlCGg0cvAgo-sJQh43L5Pe4W",
                "timestamp": "1603059201",
                "text": {
                  "body": "Hello, this is my message"
                },
                "type": "text"
              }
            ]
          },
          "field": "messages"
        }
      ]
    }
  ]
}
```

## 2. Esquemas JSON (Go Type Definitions)

### 2.1 WebhookRequest

```go
type WebhookRequest struct {
    Object string       `json:"object"`
    Entry  []Entry     `json:"entry"`
}

type Entry struct {
    ID      string    `json:"id"`
    Changes []Change  `json:"changes"`
}

type Change struct {
    Value WebhookValue `json:"value"`
    Field string       `json:"field"`
}

type WebhookValue struct {
    MessagingProduct string     `json:"messaging_product"`
    Metadata         Metadata   `json:"metadata"`
    Contacts         []Contact  `json:"contacts"`
    Messages         []Message  `json:"messages"`
}

type Metadata struct {
    DisplayPhoneNumber string `json:"display_phone_number"`
    PhoneNumberID      string `json:"phone_number_id"`
}

type Contact struct {
    Profile Profile `json:"profile"`
    WaID    string  `json:"wa_id"`
}

type Profile struct {
    Name string `json:"name"`
}
```

### 2.2 Message (para mensajes de texto)

```go
type Message struct {
    From    string     `json:"from"`
    ID      string     `json:"id"`
    Timestamp string   `json:"timestamp"`
    Type    string     `json:"type"`  // "text", "image", "audio", etc.
    Text    TextBody   `json:"text"`
}

type TextBody struct {
    Body string `json:"body"`
}
```

## 3. Campos Relevantes para el Sistema

### 3.1 Identificación del Remitente

| Campo | Ubicación | Descripción |
|-------|-----------|-------------|
| `from` | `entry[].changes[].value.messages[].from` | Número de teléfono del usuario emisor (formato internacional) |
| `wa_id` | `entry[].changes[].value.contacts[].wa_id` | WhatsApp ID del contacto |

### 3.2 Contenido del Mensaje

| Campo | Ubripción | Descripción |
|-------|-----------|-------------|
| `type` | `entry[].changes[].value.messages[].type` | Tipo de mensaje ("text", "image", etc.) |
| `text.body` | `entry[].changes[].value.messages[].text.body` | Texto del mensaje |
| `timestamp` | `entry[].changes[].value.messages[].timestamp` | Timestamp UNIX del mensaje |

### 3.3 metadata para Respuestas

| Campo | Ubicación | Descripción |
|-------|-----------|-------------|
| `phone_number_id` | `entry[].changes[].value.metadata.phone_number_id` | ID del número de teléfono para enviar respuestas |

## 4. Verificación del Webhook

### 4.1 Request de Verificación (GET)

WhatsApp envía una solicitud GET para verificar el endpoint:

```
GET /webhook?hub.mode=subscribe&hub.verify_token=<TOKEN>&hub.challenge=<CHALLENGE>
```

### 4.2 Respuesta de Verificación

El servidor debe:
1. Verificar que `hub.verify_token` coincida con el token configurado
2. Retornar el valor de `hub.challenge` como plain text (HTTP 200)

## 5. Tipos de Mensaje Soportados (v1)

| Tipo | Descripción | Soportado |
|------|-------------|-----------|
| `text` | Mensaje de texto | **Sí** |
| `image` | Imagen con caption | No (v1) |
| `audio` | Audio/Voice message | No (v1) |
| `video` | Video | No (v1) |
| `document` | Documento | No (v1) |
| `sticker` | Sticker | No (v1) |
| `reaction` | Reacción a mensaje | No (v1) |
| `interactive` | Mensaje interactivo (botones, listas) | No (v1) |

## 6. Contractos de Eventos

### 6.1 Evento de Mensaje Entrante

```go
type IncomingMessageEvent struct {
    EventType    string    `json:"event_type"`  // "message"
    MessageID    string    `json:"message_id"`
    From         string    `json:"from"`
    Timestamp    int64     `json:"timestamp"`
    MessageType  string    `json:"message_type"`
    Text         string    `json:"text"`
    RawPayload   json.RawMessage `json:"raw_payload"`
}
```

### 6.2 Evento de Verificación

```go
type WebhookVerificationRequest struct {
    Mode      string `query:"hub.mode"`
    Token     string `query:"hub.verify_token"`
    Challenge string `query:"hub.challenge"`
}
```