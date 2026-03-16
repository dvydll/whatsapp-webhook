# CONTRACT-009: Contract Tests

## 1. Propósito

Los contract tests verifican que las implementaciones satisfacen las interfaces definidas. Son tests de cumplimiento, no de funcionalidad completa.

## 2. Test Structure

```
internal/
├── domain/
│   ├── message_test.go
│   ├── response_test.go
│   └── audio_test.go
│
├── pipeline/
│   ├── pipeline_test.go
│   └── stage_test.go
│
├── adapters/
│   ├── tts/
│   │   └── provider_test.go
│   ├── delivery/
│   │   └── adapter_test.go
│   └── audio/
│       └── processor_test.go
```

## 3. Interface Compliance Tests

### 3.1 TTS Provider Contract Test

```go
package tts

import (
    "context"
    "testing"
    "github.com/whatsapp-tts/internal/domain"
)

// TestTTSProviderContract verifica que un provider cumple la interfaz TTSProvider.
func TestTTSProviderContract(t *testing.T) {
    provider := NewProvider(Config{Type: ProviderStyleTTS})
    
    // Verify interface compliance
    var _ TTSProvider = provider
    
    // Test Generate returns valid AudioAsset
    ctx := context.Background()
    audio, err := provider.Generate(ctx, "Hello world")
    
    if err != nil {
        t.Errorf("Generate() error = %v", err)
        return
    }
    
    if audio == nil {
        t.Error("Generate() returned nil audio")
        return
    }
    
    if audio.AudioID == "" {
        t.Error("Generate() returned audio with empty AudioID")
    }
    
    if audio.Format == "" {
        t.Error("Generate() returned audio with empty Format")
    }
    
    // Test Name returns non-empty string
    if provider.Name() == "" {
        t.Error("Name() returned empty string")
    }
}
```

### 3.2 Delivery Adapter Contract Test

```go
package adapter

import (
    "context"
    "testing"
    "github.com/whatsapp-tts/internal/domain"
)

// TestDeliveryAdapterContract verifica que un adapter cumple la interfaz.
func TestDeliveryAdapterContract(t *testing.T) {
    adapter := NewWhatsAppAdapter(Config{...})
    
    // Verify interface compliance
    var _ DeliveryAdapter = adapter
    
    // Test SendTextMessage
    ctx := context.Background()
    msgID, err := adapter.SendTextMessage(ctx, "+1234567890", "Test message")
    
    if err != nil {
        t.Errorf("SendTextMessage() error = %v", err)
        return
    }
    
    if msgID == "" {
        t.Error("SendTextMessage() returned empty message ID")
    }
    
    // Test Name returns non-empty string
    if adapter.Name() == "" {
        t.Error("Name() returned empty string")
    }
}
```

## 3. Pipeline Stage Tests

```go
package pipeline

import (
    "context"
    "testing"
)

// TestStageInterfaceCompliance verifica cumplimiento de Stage.
func TestStageInterfaceCompliance(t *testing.T) {
    stages := []Stage{
        NewIngestionStage("token"),
        NewNormalizationStage(),
        NewResponseGenerationStage(),
        NewTTSGenerationStage(mockTTSProvider{}),
        NewDeliveryStage(mockDeliveryAdapter{}),
    }
    
    for _, stage := range stages {
        // Verify Name returns non-empty
        if stage.Name() == "" {
            t.Errorf("%T: Name() returned empty string", stage)
        }
        
        // Verify CanProcess works
        if !stage.CanProcess(nil) {
            t.Errorf("%T: CanProcess(nil) returned false", stage)
        }
    }
}

// TestPipelineExecution verifica la ejecución del pipeline.
func TestPipelineExecution(t *testing.T) {
    pipeline := NewPipeline(
        NewIngestionStage("token"),
        NewNormalizationStage(),
    )
    
    ctx := context.Background()
    result, err := pipeline.Execute(ctx, testWebhookPayload)
    
    if err != nil {
        t.Errorf("Execute() error = %v", err)
        return
    }
    
    if result == nil {
        t.Error("Execute() returned nil result")
    }
}
```

## 4. Domain Model Tests

```go
package domain

import (
    "testing"
    "time"
)

func TestUserMessageConstruction(t *testing.T) {
    metadata := MessageMetadata{
        UserPhone:       "+1234567890",
        BusinessPhoneID: "123456789",
    }
    
    msg := NewUserMessage(
        "msg-123",
        "wamid.abc",
        "user-456",
        ChannelWhatsApp,
        "Hello",
        time.Now(),
        metadata,
    )
    
    if msg.MessageID != "msg-123" {
        t.Errorf("MessageID = %v, want msg-123", msg.MessageID)
    }
    
    if msg.Text != "Hello" {
        t.Errorf("Text = %v, want Hello", msg.Text)
    }
}
```

## 5. Test Helpers

```go
package testutil

import "github.com/whatsapp-tts/internal/domain"

// MockTTSProvider es un mock para testing.
type MockTTSProvider struct {
    GenerateFunc func(ctx context.Context, text string) (*domain.AudioAsset, error)
}

func (m *MockTTSProvider) Generate(ctx context.Context, text string) (*domain.AudioAsset, error) {
    if m.GenerateFunc != nil {
        return m.GenerateFunc(ctx, text)
    }
    return &domain.AudioAsset{AudioID: "mock-id"}, nil
}

func (m *MockTTSProvider) GenerateWithOptions(ctx context.Context, text string, opts Options) (*domain.AudioAsset, error) {
    return m.Generate(ctx, text)
}

func (m *MockTTSProvider) Name() string { return "mock" }

// MockDeliveryAdapter es un mock para testing.
type MockDeliveryAdapter struct {
    DeliverFunc func(ctx context.Context, response *domain.ResponseMessage, audio *domain.AudioAsset) error
}

func (m *MockDeliveryAdapter) Deliver(ctx context.Context, response *domain.ResponseMessage, audio *domain.AudioAsset) error {
    if m.DeliverFunc != nil {
        return m.DeliverFunc(ctx, response, audio)
    }
    return nil
}

func (m *MockDeliveryAdapter) SendTextMessage(ctx context.Context, to, text string) (string, error) {
    return "msg-id", nil
}

func (m *MockDeliveryAdapter) SendAudioMessage(ctx context.Context, to, audioID string) (string, error) {
    return "msg-id", nil
}

func (m *MockDeliveryAdapter) Name() string { return "mock" }
```

## 6. Test Data

```go
package testutil

var (
    TestWebhookPayload = []byte(`{
        "object": "whatsapp_business_account",
        "entry": [{
            "id": "123456789",
            "changes": [{
                "value": {
                    "messaging_product": "whatsapp",
                    "metadata": {
                        "display_phone_number": "+1234567890",
                        "phone_number_id": "123456789"
                    },
                    "messages": [{
                        "from": "9876543210",
                        "id": "wamid.abc123",
                        "timestamp": "1603059201",
                        "type": "text",
                        "text": {
                            "body": "Hello world"
                        }
                    }]
                },
                "field": "messages"
            }]
        }]
    }`)
)
```

## 7. Running Tests

```makefile
test:
	go test -race -v ./...

test-contract:
	go test -v -run "Contract" ./...

test-unit:
	go test -v -run "TestUser|TestMessage|TestAudio" ./...
```

## 8. Notas

- Tests verifican interface compliance, no implementación completa
- Mocks permiten testing de componentes dependientes
- Test data facilita reproducibilidad
- Usar -race flag para detectar race conditions