# SPEC-003: Internal Message Model

## 1. Propósito

El modelo de mensaje interno normaliza los datos provenientes del webhook de WhatsApp para decoupling del formato raw externo. Esto permite que el sistema procese mensajes independientemente de la fuente (WhatsApp u otros canales futuros).

## 2. Estructura del Modelo

```go
package domain

import "time"

// Message representa el modelo interno normalizado de un mensaje.
type Message struct {
    // ID único del mensaje en el sistema
    MessageID string
    
    // ID del mensaje en la plataforma externa (WhatsApp)
    ExternalMessageID string
    
    // Identificador único del usuario emisor
    UserID string
    
    // Canal por el cual llegó el mensaje
    Channel Channel
    
    // Tipo de contenido del mensaje
    ContentType ContentType
    
    // Contenido textual del mensaje
    TextContent string
    
    // Timestamp del mensaje (cuando fue enviado por el usuario)
    Timestamp time.Time
    
    // Metadata adicional del mensaje
    Metadata MessageMetadata
    
    // Estado actual del procesamiento
    ProcessingStatus ProcessingStatus
}

// Channel define los canales de entrada/salida soportados.
type Channel string

const (
    ChannelWhatsApp Channel = "whatsapp"
    ChannelTelegram Channel = "telegram"  // Futuro
    ChannelWebChat  Channel = "webchat"    // Futuro
    ChannelSMS      Channel = "sms"        // Futuro
)

// ContentType define el tipo de contenido del mensaje.
type ContentType string

const (
    ContentTypeText     ContentType = "text"
    ContentTypeImage    ContentType = "image"
    ContentTypeAudio    ContentType = "audio"
    ContentTypeVideo    ContentType = "video"
    ContentTypeDocument ContentType = "document"
    ContentTypeSticker  ContentType = "sticker"
    ContentTypeUnknown  ContentType = "unknown"
)

// ProcessingStatus representa el estado del procesamiento.
type ProcessingStatus string

const (
    StatusReceived   ProcessingStatus = "received"
    StatusNormalized ProcessingStatus = "normalized"
    StatusProcessing ProcessingStatus = "processing"
    StatusCompleted   ProcessingStatus = "completed"
    StatusFailed      ProcessingStatus = "failed"
)

// MessageMetadata contiene información adicional del mensaje.
type MessageMetadata struct {
    // Número de teléfono del usuario en formato internacional
    UserPhone string
    
    // Nombre de perfil del usuario (si está disponible)
    UserProfileName string
    
    // ID del número de teléfono del negocio (para respuestas)
    BusinessPhoneID string
    
    // Número de teléfono del negocio (display)
    BusinessDisplayPhone string
    
    // ID de la cuenta de WhatsApp Business
    BusinessAccountID string
    
    // Hash de identidad del usuario (si existe)
    IdentityHash string
    
    // Flags de mensajes reenviados
    IsForwarded bool
    
    // Indica si el mensaje ha sido reenviado frecuentemente
    IsFrequentlyForwarded bool
    
    // Mensaje al que responde (si es reply)
    ReplyToMessageID string
    
    // Datos crudos originales del webhook
    RawPayload []byte
}
```

## 3. Constructor y Factory

```go
// NewMessage crea un nuevo mensaje normalizado a partir del payload de WhatsApp.
func NewMessage(
    externalMsgID string,
    userID string,
    channel Channel,
    contentType ContentType,
    textContent string,
    timestamp time.Time,
    metadata MessageMetadata,
) *Message {
    return &Message{
        MessageID:           generateUUID(),
        ExternalMessageID:   externalMsgID,
        UserID:              userID,
        Channel:             channel,
        ContentType:         contentType,
        TextContent:         textContent,
        Timestamp:           timestamp,
        Metadata:            metadata,
        ProcessingStatus:    StatusReceived,
    }
}

// FromWhatsAppPayload crea un mensaje normalizado desde el payload del webhook de WhatsApp.
func FromWhatsAppPayload(payload *WhatsAppWebhookPayload) (*Message, error) {
    if len(payload.Entry) == 0 || len(payload.Entry[0].Changes) == 0 {
        return nil, ErrInvalidWebhookPayload
    }
    
    change := payload.Entry[0].Changes[0]
    value := change.Value
    
    if len(value.Messages) == 0 {
        return nil, ErrNoMessagesInPayload
    }
    
    msg := value.Messages[0]
    
    // Extraer metadata
    metadata := MessageMetadata{
        UserPhone:            msg.From,
        UserProfileName:      extractProfileName(value.Contacts, msg.From),
        BusinessPhoneID:      value.Metadata.PhoneNumberID,
        BusinessDisplayPhone: value.Metadata.DisplayPhoneNumber,
        BusinessAccountID:    payload.Entry[0].ID,
    }
    
    // Parsear timestamp
    ts, err := strconv.ParseInt(msg.Timestamp, 10, 64)
    if err != nil {
        ts = time.Now().Unix()
    }
    
    return NewMessage(
        msg.ID,
        msg.From,
        ChannelWhatsApp,
        mapMessageType(msg.Type),
        msg.Text.Body,
        time.Unix(ts, 0),
        metadata,
    ), nil
}

// mapMessageType mapea el tipo de mensaje de WhatsApp al tipo interno.
func mapMessageType(watype string) ContentType {
    switch watype {
    case "text":
        return ContentTypeText
    case "image":
        return ContentTypeImage
    case "audio":
        return ContentTypeAudio
    case "video":
        return ContentTypeVideo
    case "document":
        return ContentTypeDocument
    case "sticker":
        return ContentTypeSticker
    default:
        return ContentTypeUnknown
    }
}

// extractProfileName busca el nombre del perfil en los contacts.
func extractProfileName(contacts []Contact, waID string) string {
    for _, c := range contacts {
        if c.WaID == waID {
            return c.Profile.Name
        }
    }
    return ""
}
```

## 4. Validaciones

```go
// Validate verifica que el mensaje tenga los campos obligatorios.
func (m *Message) Validate() error {
    if m.MessageID == "" {
        return ErrMessageIDRequired
    }
    if m.UserID == "" {
        return ErrUserIDRequired
    }
    if m.Channel == "" {
        return ErrChannelRequired
    }
    if m.ProcessingStatus == "" {
        return ErrStatusRequired
    }
    return nil
}
```

## 5. Estado y Transiciones

```
[StatusReceived] --> [StatusNormalized] --> [StatusProcessing] --> [StatusCompleted]
                                        |
                                        v
                                    [StatusFailed]
```

### 5.1 Descripción de Estados

| Estado | Descripción |
|--------|-------------|
| `received` | Mensaje recibido del webhook |
| `normalized` | Convertido al modelo interno |
| `processing` | En proceso dentro del pipeline |
| `completed` | Procesamiento exitoso completado |
| `failed` | Fallo en el procesamiento |

## 6. Notas de Diseño

- El `MessageID` interno es un UUID generado por el sistema
- El `ExternalMessageID` preserva el ID de WhatsApp para trazabilidad
- El `Channel` permite extensibilidad futura para otros canales
- El `ContentType` permite manejar diferentes tipos de contenido
- La `Metadata` contiene información específica del canal que no debe filtrar al modelo de dominio