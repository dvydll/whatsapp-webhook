# SPEC-005: Response Contract

## 1. Propósito

El contrato de respuesta define una abstracción independiente del canal de entrega para las respuestas generadas por el sistema. Esto permite que el mismo pipeline pueda enviar respuestas a múltiples canales en el futuro (WhatsApp, Telegram, Web, etc.).

## 2. Estructura del Modelo de Respuesta

```go
package domain

import "time"

// Response representa una respuesta generada por el sistema.
type Response struct {
    // ID único de la respuesta
    ResponseID string
    
    // Usuario objetivo de la respuesta
    TargetUser string
    
    // Tipo de respuesta
    ResponseType ResponseType
    
    // Contenido textual de la respuesta
    ResponseText string
    
    // Referencia al audio (URL o ID interno)
    AudioReference string
    
    // Canal de entrega
    DeliveryChannel Channel
    
    // Timestamp de creación
    CreatedAt time.Time
    
    // Metadata adicional
    Metadata ResponseMetadata
}
```

## 3. Tipos de Respuesta

```go
// ResponseType define el tipo de respuesta.
type ResponseType string

const (
    ResponseTypeText   ResponseType = "text"
    ResponseTypeAudio  ResponseType = "audio"
    ResponseTypeImage  ResponseType = "image"
    ResponseTypeVideo  ResponseType = "video"
    ResponseTypeDoc    ResponseType = "document"
    ResponseTypeButtons ResponseType = "interactive_buttons"
    ResponseTypeList   ResponseType = "interactive_list"
)
```

## 4. Metadata de Respuesta

```go
// ResponseMetadata contiene metadata adicional de la respuesta.
type ResponseMetadata struct {
    // ID del mensaje al que responde (si es reply)
    ReplyToMessageID string
    
    // Mensaje original recibido
    OriginalMessageID string
    
    // Phone number ID para enviar a WhatsApp
    PhoneNumberID string
    
    // Language code para TTS
    Language string
    
    // Voz a usar para TTS
    VoiceID string
    
    // ID de correlación con la solicitud original
    CorrelationID string
    
    // Flags de entrega
    IsUrgent bool
}
```

## 5. Constructor

```go
// NewResponse crea una nueva respuesta.
func NewResponse(
    targetUser string,
    responseType ResponseType,
    responseText string,
    phoneNumberID string,
) *Response {
    return &Response{
        ResponseID:     generateUUID(),
        TargetUser:     targetUser,
        ResponseType:   responseType,
        ResponseText:   responseText,
        DeliveryChannel: ChannelWhatsApp,
        CreatedAt:       time.Now(),
        Metadata: ResponseMetadata{
            PhoneNumberID: phoneNumberID,
        },
    }
}

// NewAudioResponse crea una respuesta de audio con referencia.
func NewAudioResponse(
    targetUser string,
    responseText string,
    audioReference string,
    phoneNumberID string,
) *Response {
    resp := NewResponse(targetUser, ResponseTypeAudio, responseText, phoneNumberID)
    resp.AudioReference = audioReference
    return resp
}
```

## 6. Serialización para WhatsApp

```go
// ToWhatsAppPayload convierte la respuesta al formato de WhatsApp Cloud API.
func (r *Response) ToWhatsAppPayload() (map[string]interface{}, error) {
    switch r.ResponseType {
    case ResponseTypeAudio:
        return r.toWhatsAppAudioPayload()
    case ResponseTypeText:
        return r.toWhatsAppTextPayload()
    default:
        return nil, ErrUnsupportedResponseType
    }
}

func (r *Response) toWhatsAppAudioPayload() (map[string]interface{}, error) {
    return map[string]interface{}{
        "messaging_product": "whatsapp",
        "to":                r.TargetUser,
        "type":              "audio",
        "audio": map[string]interface{}{
            "id": r.AudioReference,
        },
    }, nil
}

func (r *Response) toWhatsAppTextPayload() (map[string]interface{}, error) {
    return map[string]interface{}{
        "messaging_product": "whatsapp",
        "to":                r.TargetUser,
        "type":              "text",
        "text": map[string]interface{}{
            "body": r.ResponseText,
        },
    }, nil
}
```

## 7. Interfaz del Adapter de Entrega

```go
package adapter

import (
    "context"
    "errors"
)

var ErrDeliveryFailed = errors.New("delivery failed")

// DeliveryAdapter define la interfaz para enviar respuestas.
type DeliveryAdapter interface {
    // Send envía una respuesta al canal configurado.
    Send(ctx context.Context, response *domain.Response) error
    
    // SendAudio envía un audio al usuario.
    SendAudio(ctx context.Context, targetUser string, audioID string) error
    
    // SendText envía un texto al usuario.
    SendText(ctx context.Context, targetUser string, text string) error
    
    // Name retorna el nombre del adapter.
    Name() string
}

// DeliveryResult representa el resultado de una entrega.
type DeliveryResult struct {
    Success        bool
    MessageID      string
    ExternalID     string
    Error          error
    Timestamp      time.Time
}
```

## 8. Ejemplo de Uso

```go
// En el stage de Response Generation
func GenerateResponse(msg *domain.Message) *domain.Response {
    // Placeholder response (v1)
    responseText := "Message received. Generating audio response."
    
    return domain.NewAudioResponse(
        msg.UserID,
        responseText,
        "", // Audio reference se completa después de TTS
        msg.Metadata.BusinessPhoneID,
    )
}

// En el stage de Delivery
func (s *DeliveryStage) Deliver(response *domain.Response) error {
    // El audio reference ya está populated desde TTS stage
    return s.adapter.SendAudio(
        context.Background(),
        response.TargetUser,
        response.AudioReference,
    )
}
```