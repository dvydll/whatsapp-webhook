# SPEC-003: Internal Message Model

## 1. Propósito

El modelo de mensaje interno normaliza los datos provenientes del webhook de WhatsApp para decoupling del formato raw externo. Esto permite que el sistema procese mensajes independientemente de la fuente (WhatsApp u otros canales futuros).

## 2. Estructura del Modelo

```go
package domain

import "time"

// UserMessage representa un mensaje entrante del usuario.
type UserMessage struct {
    // MessageID es el ID único del mensaje en el sistema.
    MessageID string
    
    // ExternalMessageID es el ID del mensaje en la plataforma externa (WhatsApp).
    ExternalMessageID string
    
    // UserID es el identificador único del usuario emisor.
    UserID string
    
    // Channel es el canal por el cual llegó el mensaje.
    Channel Channel
    
    // Text es el contenido textual del mensaje.
    Text string
    
    // Timestamp es cuando fue enviado el mensaje.
    Timestamp time.Time
    
    // Metadata contiene información adicional del mensaje.
    Metadata MessageMetadata
}

// MessageMetadata contiene información adicional del mensaje.
type MessageMetadata struct {
    // UserPhone es el número de teléfono del usuario en formato internacional.
    UserPhone string
    
    // UserProfileName es el nombre de perfil del usuario (si está disponible).
    UserProfileName string
    
    // BusinessPhoneID es el ID del número de teléfono del negocio (para respuestas).
    BusinessPhoneID string
    
    // BusinessAccountID es el ID de la cuenta de WhatsApp Business.
    BusinessAccountID string
    
    // MessageType es el tipo de mensaje de WhatsApp.
    MessageType string
    
    // RawPayload contiene los datos crudos originales del webhook.
    RawPayload []byte
}

// Channel define los canales de entrada/salida soportados.
type Channel string

const (
    ChannelWhatsApp Channel = "whatsapp"
    ChannelTelegram Channel = "telegram"  // Futuro
    ChannelWebChat  Channel = "webchat"    // Futuro
)
```

## 3. Constructor y Factory

```go
// NewUserMessage crea un nuevo mensaje normalizado a partir del payload de WhatsApp.
func NewUserMessage(
    messageID string,
    externalMessageID string,
    userID string,
    channel Channel,
    text string,
    timestamp time.Time,
    metadata MessageMetadata,
) *UserMessage {
    return &UserMessage{
        MessageID:         messageID,
        ExternalMessageID: externalMessageID,
        UserID:            userID,
        Channel:           channel,
        Text:              text,
        Timestamp:         timestamp,
        Metadata:          metadata,
    }
}
```

## 4. Validaciones

```go
// Validate verifica que el mensaje tenga los campos obligatorios.
func (m *UserMessage) Validate() error {
    if m.MessageID == "" {
        return ErrMessageIDRequired
    }
    if m.UserID == "" {
        return ErrUserIDRequired
    }
    if m.Channel == "" {
        return ErrChannelRequired
    }
    return nil
}
```

## 5. Notas de Diseño

- El `MessageID` interno es un UUID generado por el sistema
- El `ExternalMessageID` preserva el ID de WhatsApp para trazabilidad
- El `Channel` permite extensibilidad futura para otros canales
- La `Metadata` contiene información específica del canal que no debe filtrar al modelo de dominio
- Para el modelo de errores ver CONTRACT-008-Error-Model.md