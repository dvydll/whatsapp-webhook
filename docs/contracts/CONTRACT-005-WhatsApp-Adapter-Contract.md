# CONTRACT-005: WhatsApp Adapter Contract

## 1. WhatsAppAdapter Interface

```go
package whatsapp

import (
    "context"
)

// WhatsAppAdapter define la interfaz para integrar con WhatsApp Cloud API.
type WhatsAppAdapter interface {
    // VerifyWebhook verifica el token del webhook.
    // Retorna el challenge si es válido, error si no.
    VerifyWebhook(ctx context.Context, mode, token, challenge string) (string, error)
    
    // ParseMessage extrae mensajes del payload del webhook.
    // Retorna array de mensajes parseados.
    ParseMessage(ctx context.Context, payload []byte) ([]WhatsAppMessage, error)
    
    // SendTextMessage envía un mensaje de texto.
    SendTextMessage(ctx context.Context, to, text string) (string, error)
    
    // SendAudioMessage envía un mensaje de audio.
    // El audioID es el ID del media después de subido a WhatsApp.
    SendAudioMessage(ctx context.Context, to, audioID string) (string, error)
    
    // UploadMedia sube un archivo de audio y retorna el ID.
    UploadMedia(ctx context.Context, audioData []byte, mimeType string) (string, error)
    
    // Name retorna el nombre del adapter.
    Name() string
}
```

## 2. WhatsAppMessage DTO

```go
package whatsapp

import "time"

// WhatsAppMessage representa un mensaje parseado del webhook.
type WhatsAppMessage struct {
    // MessageID es el ID del mensaje de WhatsApp.
    MessageID string
    
    // From es el número del usuario emisor.
    From string
    
    // Timestamp es el timestamp del mensaje.
    Timestamp time.Time
    
    // Type es el tipo de mensaje.
    Type string
    
    // Text es el contenido si es texto.
    Text string
    
    // Image es la info de imagen si es imagen.
    Image *ImageInfo
    
    // Audio es la info de audio si es audio.
    Audio *AudioInfo
    
    // Metadata contiene metadata adicional.
    Metadata MessageMetadata
}

type ImageInfo struct {
    ID       string
    MIMEType string
    Caption  string
    SHA256   string
}

type AudioInfo struct {
    ID       string
    MIMEType string
    SHA256   string
}

type MessageMetadata struct {
    PhoneNumberID      string
    BusinessAccountID  string
}
```

## 3. Webhook Payload DTOs

```go
package whatsapp

// WebhookRequest representa el payload completo del webhook.
type WebhookRequest struct {
    Object string    `json:"object"`
    Entry  []Entry   `json:"entry"`
}

type Entry struct {
    ID      string   `json:"id"`
    Changes []Change `json:"changes"`
}

type Change struct {
    Value  WebhookValue `json:"value"`
    Field  string       `json:"field"`
}

type WebhookValue struct {
    MessagingProduct string    `json:"messaging_product"`
    Metadata         Metadata `json:"metadata"`
    Contacts         []Contact `json:"contacts"`
    Messages         []Message `json:"messages"`
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

type Message struct {
    From     string     `json:"from"`
    ID       string     `json:"id"`
    Timestamp string   `json:"timestamp"`
    Type     string     `json:"type"`
    Text     TextBody  `json:"text"`
    Image    ImageBody `json:"image"`
    Audio    AudioBody `json:"audio"`
}

type TextBody struct {
    Body string `json:"body"`
}

type ImageBody struct {
    ID       string `json:"id"`
    MIMEType string `json:"mime_type"`
    SHA256   string `json:"sha256"`
    Caption  string `json:"caption"`
}

type AudioBody struct {
    ID       string `json:"id"`
    MIMEType string `json:"mime_type"`
    SHA256   string `json:"sha256"`
}
```

## 4. API Response DTOs

```go
package whatsapp

// MessageResponse representa la respuesta de enviar mensaje.
type MessageResponse struct {
    MessagingProduct string   `json:"messaging_product"`
    Contacts         []ContactResult `json:"contacts"`
    Messages         []MessageResult `json:"messages"`
}

type ContactResult struct {
    Input string `json:"input"`
    WaID  string `json:"wa_id"`
}

type MessageResult struct {
    ID string `json:"id"`
}

// MediaUploadResponse representa la respuesta de subir media.
type MediaUploadResponse struct {
    ID string `json:"id"`
}

// ErrorResponse representa un error de la API.
type ErrorResponse struct {
    Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
    Message   string `json:"message"`
    Type      string `json:"type"`
    Code      int    `json:"code"`
    ErrorData ErrorData `json:"error_data"`
}

type ErrorData struct {
    Details string `json:"details"`
}
```

## 5. Send Message DTOs

```go
package whatsapp

// SendTextRequest representa el request para enviar texto.
type SendTextRequest struct {
    MessagingProduct string      `json:"messaging_product"`
    To               string      `json:"to"`
    Type             string      `json:"type"`
    Text             TextContent `json:"text"`
}

type TextContent struct {
    Body string `json:"body"`
}

// SendAudioRequest representa el request para enviar audio.
type SendAudioRequest struct {
    MessagingProduct string      `json:"messaging_product"`
    To               string      `json:"to"`
    Type             string      `json:"type"`
    Audio            AudioContent `json:"audio"`
}

type AudioContent struct {
    ID string `json:"id"`
}
```

## 6. Factory

```go
package whatsapp

import "context"

// Config contiene la configuración del adapter.
type Config struct {
    BaseURL     string
    PhoneNumberID string
    AccessToken string
    VerifyToken string
}

// NewWhatsAppAdapter crea una nueva instancia del adapter.
func NewWhatsAppAdapter(config Config) WhatsAppAdapter {
    return &whatsAppAdapterImpl{
        config: config,
    }
}
```

## 7. Notas

- El adapter aísla todo el formato específico de WhatsApp
- Los DTOs mapean directamente al JSON de la API
- Permite cambio de implementación (mock, otro adapter) sin cambiar código cliente
- Ver CONTRACT-004 para cómo se integra con el pipeline