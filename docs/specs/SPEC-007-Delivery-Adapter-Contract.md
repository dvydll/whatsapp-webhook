# SPEC-007: Delivery Adapter Contract

## 1. Propósito

El adapter de entrega es responsable de enviar respuestas al canal de mensajería externo. La interfaz permite soportar múltiples canales (WhatsApp, Telegram, etc.) sin modificar la lógica del pipeline.

## 2. WhatsApp Client Interface

El cliente de WhatsApp se encuentra en `internal/adapters/whatsapp/client.go`.

```go
package whatsapp

import (
    "context"
)

type Client struct {
    httpClient  *http.Client
    baseURL     string
    version     string
    phoneID     string
    accessToken string
}

type Config struct {
    PhoneNumberID string
    AccessToken   string
    BaseURL       string
    Version       string
}

func NewClient(config Config) *Client
```

## 3. Métodos del Cliente

```go
func (c *Client) SendTextMessage(ctx context.Context, to, text string) (*SendMessageResponse, error)

func (c *Client) SendAudioMessage(ctx context.Context, to, audioURL string) (*SendMessageResponse, error)
```

## 4. DTOs de Request/Response

```go
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
    MessagingProduct string       `json:"messaging_product"`
    To               string       `json:"to"`
    Type             string       `json:"type"`
    Audio            AudioContent `json:"audio"`
}

type AudioContent struct {
    Link string `json:"link,omitempty"`
    ID   string `json:"id,omitempty"`
}

// SendMessageResponse representa la respuesta de enviar mensaje.
type SendMessageResponse struct {
    MessagingProduct string    `json:"messaging_product"`
    Contacts         []Contact `json:"contacts"`
    Messages         []Message `json:"messages"`
}

type Contact struct {
    Input string `json:"input"`
    WaID  string `json:"wa_id"`
}

type Message struct {
    ID string `json:"id"`
}
```

## 5. Integración con Pipeline

```go
type DeliveryStage struct {
    whatsappClient *whatsapp.Client
}

func (s *DeliveryStage) Process(ctx context.Context, input interface{}) error {
    ttsOut := input.(*TTSOutput)
    _, err := s.whatsappClient.SendAudioMessage(ctx, ttsOut.Response.TargetUser, ttsOut.Audio.URL)
    return err
}
```

## 6. Notas de Diseño

- El cliente aísla todo el formato específico de WhatsApp
- Los DTOs mapean directamente al JSON de la API
- Permite cambio de implementación (mock, otro adapter) sin cambiar código cliente
- Para la interfaz completa ver CONTRACT-005-WhatsApp-Adapter-Contract.md