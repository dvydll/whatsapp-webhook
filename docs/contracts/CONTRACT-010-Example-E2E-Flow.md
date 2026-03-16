# CONTRACT-010: Example End-to-End Flow

## 1. Flujo Completo

```
┌─────────────┐     ┌──────────────┐     ┌───────────────┐
│   WhatsApp  │────>│   Ingestion  │────>│ Normalization │
│    Webhook  │     │    Stage     │     │    Stage      │
└─────────────┘     └──────────────┘     └───────────────┘
                                                 |
                                                 v
┌─────────────┐     ┌──────────────┐     ┌───────────────┐
│   WhatsApp  │<────│   Delivery   │<────│  TTS Engine   │
│     API      │     │    Stage     │     │    Stage      │
└─────────────┘     └──────────────┘     └───────────────┘
                                                  |
                                                  v
                                           ┌───────────────┐
                                           │   Response    │
                                           │   Generation  │
                                           │    Stage      │
                                           └───────────────┘
```

## 2. Pseudo-code del Flujo

```go
package main

import (
    "context"
    "net/http"
)

func main() {
    // Inicialización
    config := loadConfig()
    
    // Crear componentes
    ttsProvider := tts.NewProvider(config.TTS)
    whatsappAdapter := whatsapp.NewAdapter(config.WhatsApp)
    audioProcessor := audio.NewProcessor(config.Audio)
    
    // Crear stages
    stages := []pipeline.Stage{
        pipeline.NewIngestionStage(config.VerifyToken),
        pipeline.NewNormalizationStage(),
        pipeline.NewResponseGenerationStage(),
        pipeline.NewTTSGenerationStage(ttsProvider),
        pipeline.NewAudioProcessingStage(audioProcessor),
        pipeline.NewDeliveryStage(whatsappAdapter),
    }
    
    // Crear pipeline
    p := pipeline.NewPipeline(stages...)
    
    // Handler HTTP
    http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
        handleWebhook(p, w, r)
    })
    
    http.ListenAndServe(":8080", nil)
}

func handleWebhook(p pipeline.Pipeline, w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // 1. Ingestion Stage recibe el request
    rawEvent, err := p.Stages()[0].Process(ctx, r)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // 2. Normalization Stage convierte a modelo interno
    message, err := p.Stages()[1].Process(ctx, rawEvent)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 3. Response Generation crea respuesta de texto
    response, err := p.Stages()[2].Process(ctx, message)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 4. TTS Generation convierte texto a audio
    audio, err := p.Stages()[3].Process(ctx, response)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 5. Audio Processing prepara para WhatsApp
    processedAudio, err := p.Stages()[4].Process(ctx, audio)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // 6. Delivery Stage envía el audio
    err = p.Stages()[5].Process(ctx, struct{
        Response *domain.ResponseMessage
        Audio    *domain.AudioAsset
    }{response, processedAudio})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    w.WriteHeader(http.StatusOK)
}
```

## 3. Flujo Detallado por Stage

### 3.1 Ingestion Stage

```go
// Input: HTTP Request
// Output: RawEvent

func (s *IngestionStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
    req := input.(*http.Request)
    
    // Verificar GET request (webhook verification)
    if req.Method == "GET" {
        challenge := req.URL.Query().Get("hub.challenge")
        token := req.URL.Query().Get("hub.verify_token")
        
        if token != s.verifyToken {
            return nil, ErrVerificationFailed
        }
        
        // Return challenge para verificación
        return &VerificationResult{Challenge: challenge}, nil
    }
    
    // Parsear POST request (webhook events)
    body, _ := io.ReadAll(req.Body)
    
    return &RawEvent{
        Payload: body,
        Headers: mapHeader(req.Header),
        Method:  req.Method,
    }, nil
}
```

### 3.2 Normalization Stage

```go
// Input: RawEvent
// Output: *domain.UserMessage

func (s *NormalizationStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
    event := input.(*RawEvent)
    
    // Parsear payload de WhatsApp
    var webhook WhatsAppWebhook
    json.Unmarshal(event.Payload, &webhook)
    
    // Extraer mensaje
    msg := webhook.Entry[0].Changes[0].Value.Messages[0]
    
    // Convertir a modelo interno
    return domain.NewUserMessage(
        generateUUID(),
        msg.ID,
        msg.From,
        domain.ChannelWhatsApp,
        msg.Text.Body,
        parseTimestamp(msg.Timestamp),
        domain.MessageMetadata{
            UserPhone:       msg.From,
            BusinessPhoneID: webhook.Entry[0].Changes[0].Value.Metadata.PhoneNumberID,
        },
    ), nil
}
```

### 3.3 Response Generation Stage

```go
// Input: *domain.UserMessage
// Output: *domain.ResponseMessage

func (s *ResponseGenerationStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
    msg := input.(*domain.UserMessage)
    
    // Generar respuesta (placeholder v1)
    responseText := "Message received. Generating audio response."
    
    return domain.NewResponseMessage(
        msg.UserID,
        domain.ResponseTypeAudio,
        responseText,
        msg.Metadata.BusinessPhoneID,
    ), nil
}
```

### 3.4 TTS Generation Stage

```go
// Input: *domain.ResponseMessage
// Output: *domain.AudioAsset

func (s *TTSGenerationStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
    response := input.(*domain.ResponseMessage)
    
    // Generar audio
    audio, err := s.ttsProvider.Generate(ctx, response.Text)
    if err != nil {
        return nil, ErrTTSFailed
    }
    
    return audio, nil
}
```

### 3.5 Audio Processing Stage

```go
// Input: *domain.AudioAsset
// Output: *domain.AudioAsset (procesado)

func (s *AudioProcessingStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
    audio := input.(*domain.AudioAsset)
    
    // Convertir a formato WhatsApp
    processed, err := s.processor.ConvertToWhatsAppFormat(ctx, audio)
    if err != nil {
        return nil, ErrAudioProcessingFailed
    }
    
    return processed, nil
}
```

### 3.6 Delivery Stage

```go
// Input: {Response, Audio}
// Output: error (nil = éxito)

func (s *DeliveryStage) Process(ctx context.Context, input interface{}) error {
    data := input.(struct {
        Response *domain.ResponseMessage
        Audio    *domain.AudioAsset
    })
    
    // Subir audio a WhatsApp
    mediaID, err := s.adapter.UploadMedia(ctx, data.Audio.BinaryData, "audio/aac")
    if err != nil {
        return ErrUploadFailed
    }
    
    // Enviar mensaje de audio
    _, err = s.adapter.SendAudioMessage(ctx, data.Response.TargetUser, mediaID)
    if err != nil {
        return ErrDeliveryFailed
    }
    
    return nil
}
```

## 4. Manejo de Errores

```go
// El pipeline maneja errores en cada stage
func (p *Pipeline) Execute(ctx context.Context, input interface{}) (*PipelineContext, error) {
    pc := NewPipelineContext()
    current := input
    
    for _, stage := range p.Stages() {
        output, err := stage.Process(ctx, current)
        
        if err != nil {
            pc.Errors = append(pc.Errors, PipelineError{
                Stage: stage.Name(),
                Err:   err,
            })
            
            // Decidir si continuar o detener
            if isFatal(err) {
                return pc, err
            }
        }
        
        current = output
    }
    
    return pc, nil
}
```

## 5. Notas

- El flujo muestra cómo las interfaces interactúan
- Cada stage tiene responsabilidad clara
- Errores se propagan con contexto
- El flujo es síncrono (v1), pero diseñado para async futuro