# CONTRACT-003: Domain Models

## 1. UserMessage

```go
package domain

import "time"

// UserMessage representa un mensaje entrante del usuario.
type UserMessage struct {
    // MessageID es el ID único del mensaje en el sistema.
    MessageID string
    
    // ExternalMessageID es el ID del mensaje en la plataforma externa.
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
```

## 2. MessageMetadata

```go
package domain

// MessageMetadata contiene información adicional del mensaje.
type MessageMetadata struct {
    // UserPhone es el número de teléfono del usuario.
    UserPhone string
    
    // UserProfileName es el nombre del perfil del usuario.
    UserProfileName string
    
    // BusinessPhoneID es el ID del número de teléfono del negocio.
    BusinessPhoneID string
    
    // BusinessAccountID es el ID de la cuenta de WhatsApp Business.
    BusinessAccountID string
    
    // MessageType es el tipo de mensaje de WhatsApp.
    MessageType string
    
    // RawPayload contiene los datos crudos originales.
    RawPayload []byte
}
```

## 3. ResponseMessage

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
```

## 4. ResponseMetadata

```go
package domain

// ResponseMetadata contiene metadata adicional de la respuesta.
type ResponseMetadata struct {
    // OriginalMessageID es el ID del mensaje original.
    OriginalMessageID string
    
    // PhoneNumberID es el ID del número para enviar.
    PhoneNumberID string
    
    // CorrelationID es el ID de correlación con la solicitud.
    CorrelationID string
}
```

## 5. AudioAsset

```go
package domain

import "time"

// AudioAsset representa un archivo de audio generado por TTS.
type AudioAsset struct {
    // AudioID es el ID único del audio.
    AudioID string
    
    // Format es el formato del audio.
    Format AudioFormat
    
    // Codec es el códec de audio.
    Codec AudioCodec
    
    // BinaryData contiene los datos binarios del audio.
    BinaryData []byte
    
    // URL es la URL de descarga (si está disponible).
    URL string
    
    // Duration es la duración en segundos.
    Duration float64
    
    // Size es el tamaño en bytes.
    Size int64
    
    // SampleRate es la tasa de muestreo.
    SampleRate int
    
    // CreatedAt es el timestamp de generación.
    CreatedAt time.Time
}
```

## 6. Tipos Enum

```go
package domain

// Channel representa el canal de mensajería.
type Channel string

const (
    ChannelWhatsApp Channel = "whatsapp"
    ChannelTelegram Channel = "telegram"
    ChannelWebChat  Channel = "webchat"
)

// ContentType representa el tipo de contenido.
type ContentType string

const (
    ContentTypeText     ContentType = "text"
    ContentTypeImage   ContentType = "image"
    ContentTypeAudio   ContentType = "audio"
    ContentTypeVideo   ContentType = "video"
    ContentTypeDoc     ContentType = "document"
)

// ResponseType representa el tipo de respuesta.
type ResponseType string

const (
    ResponseTypeText  ResponseType = "text"
    ResponseTypeAudio ResponseType = "audio"
)

// AudioFormat representa el formato de audio.
type AudioFormat string

const (
    FormatAAC  AudioFormat = "aac"
    FormatMP3  AudioFormat = "mp3"
    FormatWAV  AudioFormat = "wav"
    FormatOGG  AudioFormat = "ogg"
)

// AudioCodec representa el códec de audio.
type AudioCodec string

const (
    CodecOpus  AudioCodec = "opus"
    CodecMP3   AudioCodec = "mp3"
    CodecPCM   AudioCodec = "pcm"
)
```

## 7. Constructores

```go
package domain

// NewUserMessage crea un nuevo UserMessage.
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

// NewResponseMessage crea un nuevo ResponseMessage.
func NewResponseMessage(
    targetUser string,
    responseType ResponseType,
    text string,
    phoneNumberID string,
) *ResponseMessage {
    return &ResponseMessage{
        ResponseID:   generateUUID(),
        TargetUser:   targetUser,
        ResponseType: responseType,
        Text:         text,
        CreatedAt:    time.Now(),
        Metadata: ResponseMetadata{
            PhoneNumberID: phoneNumberID,
        },
    }
}

// NewAudioAsset crea un nuevo AudioAsset.
func NewAudioAsset(
    format AudioFormat,
    codec AudioCodec,
    binaryData []byte,
) *AudioAsset {
    return &AudioAsset{
        AudioID:    generateUUID(),
        Format:     format,
        Codec:      codec,
        BinaryData: binaryData,
        CreatedAt:  time.Now(),
    }
}
```

## 8. Notas

- Todos los structs son puro Go, sin dependencias externas
- Los tipos enum son aliases de string para serialización JSON
- Los constructores facilitan la creación con defaults
- Los campos sonexported (mayúscula) para acceso desde otros paquetes