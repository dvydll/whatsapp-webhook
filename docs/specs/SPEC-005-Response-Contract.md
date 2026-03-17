# SPEC-005: Response Contract

## 1. Propósito

El contrato de respuesta define una abstracción independiente del canal de entrega para las respuestas generadas por el sistema. Esto permite que el mismo pipeline pueda enviar respuestas a múltiples canales en el futuro (WhatsApp, Telegram, Web, etc.).

## 2. Estructura del Modelo de Respuesta

```go
package domain

import "time"

// ResponseMessage representa una respuesta generada por el sistema.
type ResponseMessage struct {
    // ResponseID es el ID único de la respuesta.
    ResponseID string
    
    // TargetUser es el usuario objetivo de la respuesta.
    TargetUser string
    
    // ResponseType es el tipo de respuesta.
    ResponseType ResponseType
    
    // Text es el contenido textual de la respuesta.
    Text string
    
    // AudioURL es la URL del audio generado.
    AudioURL string
    
    // AudioID es el ID del audio en WhatsApp (después de upload).
    AudioID string
    
    // CreatedAt es el timestamp de creación.
    CreatedAt time.Time
    
    // Metadata contiene información adicional.
    Metadata ResponseMetadata
}

// ResponseMetadata contiene metadata adicional de la respuesta.
type ResponseMetadata struct {
    // OriginalMessageID es el ID del mensaje original.
    OriginalMessageID string
    
    // PhoneNumberID es el ID del número para enviar.
    PhoneNumberID string
    
    // CorrelationID es el ID de correlación con la solicitud.
    CorrelationID string
}

// ResponseType representa el tipo de respuesta.
type ResponseType string

const (
    ResponseTypeText  ResponseType = "text"
    ResponseTypeAudio ResponseType = "audio"
)
```

## 3. Constructor

```go
// NewResponseMessage crea una nueva ResponseMessage.
func NewResponseMessage(
    targetUser string,
    responseType ResponseType,
    text string,
    phoneNumberID string,
) *ResponseMessage {
    return &ResponseMessage{
        ResponseID:   generateID(),
        TargetUser:   targetUser,
        ResponseType: responseType,
        Text:         text,
        CreatedAt:    time.Now(),
        Metadata: ResponseMetadata{
            PhoneNumberID: phoneNumberID,
        },
    }
}
```

## 4. Notas de Diseño

- El modelo es independiente del canal de entrega
- AudioURL se usa para URLs externas (como TTS que retorna URL)
- AudioID se usa para IDs de media en WhatsApp
- Para integración real con WhatsApp ver CONTRACT-005-WhatsApp-Adapter-Contract.md