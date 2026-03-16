# SPEC-007: Delivery Adapter Contract

## 1. Propósito

El adapter de entrega es responsable de enviar respuestas al canal de mensajería externo. La interfaz permite soportar múltiples canales (WhatsApp, Telegram, etc.) sin modificar la lógica del pipeline.

## 2. Interfaz DeliveryAdapter

```go
package adapter

import (
    "context"
    "errors"
)

var (
    ErrDeliveryFailed      = errors.New("delivery failed")
    ErrInvalidPayload      = errors.New("invalid payload")
    ErrChannelNotSupported = errors.New("channel not supported")
)

type DeliveryAdapter interface {
    SendAudio(ctx context.Context, targetUser string, audioID string) error
    SendText(ctx context.Context, targetUser string, text string) error
    SendImage(ctx context.Context, targetUser string, imageURL, caption string) error
    SendVideo(ctx context.Context, targetUser string, videoURL, caption string) error
    SendDocument(ctx context.Context, targetUser string, docURL, filename string) error
    Name() string
    SupportedChannels() []string
}
```

## 3. WhatsApp Adapter

```go
package whatsapp

type WhatsAppAdapter struct {
    httpClient  *http.Client
    baseURL     string
    phoneNumber string
    accessToken string
}

func (a *WhatsAppAdapter) Name() string { return "whatsapp" }

func (a *WhatsAppAdapter) SupportedChannels() []string {
    return []string{"whatsapp"}
}

func (a *WhatsAppAdapter) SendAudio(ctx context.Context, targetUser, audioID string) error {
    payload := map[string]interface{}{
        "messaging_product": "whatsapp",
        "to":                targetUser,
        "type":              "audio",
        "audio":            map[string]interface{}{"id": audioID},
    }
    return a.sendMessage(ctx, targetUser, payload)
}

func (a *WhatsAppAdapter) SendText(ctx context.Context, targetUser, text string) error {
    payload := map[string]interface{}{
        "messaging_product": "whatsapp",
        "to":                targetUser,
        "type":              "text",
        "text":             map[string]interface{}{"body": text},
    }
    return a.sendMessage(ctx, targetUser, payload)
}

func (a *WhatsAppAdapter) sendMessage(ctx context.Context, to string, payload map[string]interface{}) error {
    endpoint := fmt.Sprintf("%s/%s/messages", a.baseURL, a.phoneNumber)
    
    reqBody, _ := json.Marshal(payload)
    req, _ := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(reqBody))
    req.Header.Set("Authorization", "Bearer "+a.accessToken)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := a.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("http request failed: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode >= 400 {
        return fmt.Errorf("whatsapp api error (status=%d)", resp.StatusCode)
    }
    return nil
}
```

## 4. Media Uploader

Para enviar audio, primero debe subirse a WhatsApp:

```go
type MediaUploader interface {
    UploadAudio(ctx context.Context, audioData []byte, mimeType string) (string, error)
}

func (a *WhatsAppAdapter) UploadAudio(ctx context.Context, audioData []byte, mimeType string) (string, error) {
    endpoint := fmt.Sprintf("%s/%s/media", a.baseURL, a.phoneNumber)
    // WhatsApp requiere form-data upload
    // Retorna el media ID de WhatsApp
}
```

## 5. Adapter Registry

```go
package adapter

type AdapterRegistry struct {
    adapters map[string]DeliveryAdapter
}

func (r *AdapterRegistry) Register(adapter DeliveryAdapter) {
    for _, ch := range adapter.SupportedChannels() {
        r.adapters[ch] = adapter
    }
}

func (r *AdapterRegistry) GetAdapter(channel string) (DeliveryAdapter, error) {
    adapter, ok := r.adapters[channel]
    if !ok {
        return nil, ErrChannelNotSupported
    }
    return adapter, nil
}
```

## 6. Integración con Pipeline

```go
type DeliveryStage struct {
    registry *AdapterRegistry
}

func (s *DeliveryStage) Process(ctx context.Context, input interface{}) error {
    ttsOut := input.(*TTSOutput)
    adapter, _ := s.registry.GetAdapter(ttsOut.Response.DeliveryChannel)
    return adapter.SendAudio(ctx, ttsOut.Response.TargetUser, ttsOut.Audio.Reference)
}
```