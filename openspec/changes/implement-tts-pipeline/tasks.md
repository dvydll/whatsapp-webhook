## 1. TTS Provider - NO IMPLEMENTADO

### 1.1 Interfaz TTSProvider
- [ ] 1.1.1 Definir interfaz TTSProvider en internal/adapters/tts/provider.go
- [ ] 1.1.2 Definir método Generate(ctx, text) -> (*domain.AudioAsset, error)
- [ ] 1.1.3 Definir método GenerateWithOptions(ctx, text, opts) -> (*domain.AudioAsset, error)

### 1.2 StyleTTS Implementation
- [ ] 1.2.1 Crear internal/adapters/tts/styletts/provider.go
- [ ] 1.2.2 Implementar HTTP client para StyleTTS API
- [ ] 1.2.3 Manejar respuesta y convertir a AudioAsset

### 1.3 Factory
- [ ] 1.3.1 Crear factory en internal/adapters/tts/factory.go
- [ ] 1.3.2 Soportar múltiples providers (StyleTTS, OpenAI, Google)

## 2. Audio Processing - NO IMPLEMENTADO

### 2.1 Conversión de Formato
- [ ] 2.1.1 Crear internal/adapters/audio/processor.go
- [ ] 2.1.2 Implementar conversión a formato WhatsApp (AAC, 48kHz, mono)
- [ ] 2.1.3 Usar FFmpeg o librería nativa

### 2.2 Validación
- [ ] 2.2.1 Validar duración máxima (16 segundos)
- [ ] 2.2.2 Validar tamaño máximo (16MB)
- [ ] 2.2.3 Validar códec soportado (AAC)

## 3. Pipeline Integration

### 3.1 Conectar TTSGenerationStage
- [ ] 3.1.1 Modificar TTSGenerationStage para usar TTSProvider real
- [ ] 3.1.2 Pasar TTSProvider al stage via constructor

### 3.2 Conectar DeliveryStage
- [ ] 3.2.1 Modificar DeliveryStage para usar WhatsApp Client
- [ ] 3.2.2 Enviar audio real vía SendAudioMessage
- [ ] 3.2.3 Necesita先 upload a WhatsApp o usar URL

## 4. Handler Pipeline Integration

### 4.1 Usar Pipeline desde Handler
- [ ] 4.1.1 Modificar webhook handler para usar pipeline.Execute()
- [ ] 4.1.2 Pasar http.Request como input al pipeline
- [ ] 4.1.3 Manejar respuesta del pipeline

## 5. Error Handling

### 5.1 AppError Estructurado
- [ ] 5.1.1 Definir AppError en internal/domain/errors.go
- [ ] 5.1.2 Definir ErrorCode enum
- [ ] 5.1.3 Definir PipelineStage enum
- [ ] 5.1.4 Implementar constructors NewError, WrapError

## 6. Observability

### 6.1 Métricas Prometheus
- [ ] 6.1.1 Agregar dependencia prometheus/client_golang
- [ ] 6.1.2 Definir métricas: messages_received, messages_processed, response_time

### 6.2 Health Check
- [ ] 6.2.1 Crear endpoint GET /health
- [ ] 6.2.2 Retornar status: healthy

## 7. Tests

- [ ] 7.1 Tests para TTSProvider interface
- [ ] 7.2 Tests para StyleTTS implementation
- [ ] 7.3 Tests para pipeline stages
- [ ] 7.4 Tests de integración handler -> pipeline
