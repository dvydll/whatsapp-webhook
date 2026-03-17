# SPEC-011: Repository Layout

## 1. Propósito

Define la estructura del repositorio para el servicio Go, siguiendo principios de clean architecture y las mejores prácticas del ecosistema Go.

## 2. Estructura Propuesta

```
whatsapp-tts/
├── cmd/                          # Punto de entrada de la aplicación
│   └── whatsapp-tts/            # Main application
│       └── main.go
│
├── internal/                     # Paquetes privados (no importables externamente)
│   ├── domain/                  # Modelos de dominio y lógica de negocio
│   │   ├── message.go
│   │   ├── response.go
│   │   ├── audio.go
│   │   ├── types.go
│   │   └── errors.go
│   │
│   ├── pipeline/                # Motor de pipeline de procesamiento
│   │   ├── pipeline.go
│   │   ├── stage.go
│   │   ├── stages/
│   │   │   ├── ingestion.go
│   │   │   ├── normalization.go
│   │   │   ├── response_generation.go
│   │   │   ├── tts_generation.go
│   │   │   └── delivery.go
│   │   └── context.go
│   │
│   ├── handler/                 # Handlers HTTP
│   │   └── webhook.go
│   │
│   ├── config/                 # Configuración de la aplicación
│   │   └── config.go
│   │
│   ├── logging/                # Logging estructurado
│   │   └── logger.go
│   │
│   └── observability/          # Observabilidad (logging, metrics)
│       └── logger.go
│
├── internal/adapters/          # Implementaciones de interfaces externas
│   ├── tts/                    # Proveedores TTS
│   │   ├── provider.go        # Interfaz TTSProvider
│   │   └── styletts/          # Implementación StyleTTS
│   │
│   ├── delivery/              # Adapters de entrega
│   │   ├── adapter.go         # Interfaz DeliveryAdapter
│   │   └── whatsapp/          # Implementación WhatsApp
│   │       └── adapter.go
│   │
│   └── whatsapp/              # Cliente WhatsApp Cloud API
│       └── client.go
│
├── pkg/                        # Paquetes públicos reutilizables
│   └── utils/                  # Utilidades generales
│
├── specs/                      # Especificaciones del sistema
│   ├── SPEC-001-System-Overview.md
│   ├── SPEC-002-Event-Contracts.md
│   └── ...
│
├── worklogs/                   # Decisiones de trabajo
│   └── WORKLOG.*.md
│
├── scripts/                    # Scripts utilitarios
│   └── whatsapp_send_message.sh
│
├── docs/                       # Documentación
│   └── WhatsApp Cloud API.postman_collection.json
│
├── .env                        # Variables de entorno (no versionar)
│
├── go.mod                      # Módulo Go
├── go.sum
│
├── README.md
│
└── Makefile                    # Comandos de build y desarrollo
```

## 3. Descripción de Paquetes

### 3.1 cmd/whatsapp-tts

Punto de entrada de la aplicación. Inicializa dependencias y arranca el servidor.

```go
// cmd/whatsapp-tts/main.go
func main() {
    cfg := config.Load()
    
    logger := logging.New(cfg.LogLevel)
    
    ttsProvider := styletts.NewStyleTTSProvider(cfg.TTSEndpoint, cfg.TTSAPIKey)
    deliveryAdapter := whatsapp.NewAdapter(cfg.WhatsAppEndpoint, cfg.PhoneNumberID, cfg.AccessToken)
    
    pipeline := pipeline.NewPipeline(
        stages.NewIngestionStage(cfg.VerifyToken),
        stages.NewNormalizationStage(),
        stages.NewResponseGenerationStage(),
        stages.NewTTSGenerationStage(ttsProvider),
        stages.NewDeliveryStage(deliveryAdapter),
    )
    
    handler := handler.NewWebhookHandler(pipeline, logger)
    
    http.ListenAndServe(":8080", handler)
}
```

### 3.2 internal/domain

Contiene los modelos de dominio puros. Sin dependencias externas.

```go
// internal/domain/message.go
package domain

type Message struct {
    MessageID         string
    ExternalMessageID string
    UserID           string
    Channel          Channel
    ContentType      ContentType
    TextContent      string
    Timestamp        time.Time
    Metadata         MessageMetadata
    ProcessingStatus ProcessingStatus
}
```

### 3.3 internal/pipeline

Orquestador del procesamiento. Define la interfaz Stage y ejecuta el flujo.

```go
// internal/pipeline/pipeline.go
package pipeline

type Pipeline struct {
    stages []Stage
}

func (p *Pipeline) Execute(ctx context.Context, input interface{}) (*PipelineContext, error)
```

### 3.4 internal/handler

Manejo de requests HTTP. Separa la capa de transporte del dominio.

```go
// internal/handler/webhook.go
package handler

type WebhookHandler struct {
    pipeline *pipeline.Pipeline
    logger   *logging.Logger
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

### 3.5 internal/adapters/tts

Implementaciones de proveedores TTS. Cumple con la interfaz TTSProvider.

```
internal/adapters/tts/
├── provider.go      # Interfaz que deben cumplir los providers
├── styletts/        # StyleTTS2 implementation
└── factory.go       # Factory para crear providers
```

### 3.6 internal/adapters/delivery

Implementaciones de adapters de entrega.

```
internal/adapters/delivery/
├── adapter.go       # Interfaz que deben cumplir los adapters
├── whatsapp/        # WhatsApp Cloud API implementation
└── registry.go      # Registry para múltiples adapters
```

### 3.7 internal/adapters/whatsapp

Cliente directo para WhatsApp Cloud API.

```
internal/adapters/whatsapp/
├── client.go        # Cliente HTTP para WhatsApp API
└── client_test.go   # Tests del cliente
```

## 4. Convenciones de Nomenclatura

### 4.1 Nombres de Archivos

- Lower snake_case: `message.go`, `webhook_handler.go`
- Test files: `message_test.go`, `handler_test.go`
- Interfaces: `nouns.go` (no I prefix en Go)

### 4.2 Nombres de Paquetes

- Lowercase: `pipeline`, `handler`, `domain`
- Evitar nombres genéricos: `util`, `common`, `misc`
- Un paquete por directorio

### 4.3 Nombres de Variables

- camelCase: `messageID`, `processingStatus`
- Acrónimos en uppercase si > 2 letras: `userID` no `userId`
- Evitar nombres cortos excepto en loops: `i`, `j`, `k`

## 5. Dependencias Externas

### 5.1 Frameworks y Librerías

| Categoría | Librería sugerida |
|-----------|-------------------|
| HTTP Router | chi o gorilla/mux |
| Logging | zerolog o zap |
| Config | envconfig o standard JSON env |
| Metrics | prometheus client |
| Validation | go-playground/validator |
| HTTP Client | standard net/http o resty |

### 5.2 Gestión de Dependencias

```makefile
# Makefile
.PHONY: deps
deps:
	go mod download
	go mod verify

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test -race -cover ./...
```

## 6. Testing

```
internal/
├── handler/
│   ├── handler.go
│   └── handler_test.go    # Tests con httptest
│
├── pipeline/
│   ├── pipeline_test.go   # Tests de integración del pipeline
│   └── stages/
│       └── stages_test.go
│
└── domain/
    ├── message_test.go    # Tests unitarios del modelo
    └── response_test.go
```

## 7. Docker (Opcional)

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /whatsapp-tts ./cmd/whatsapp-tts

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /whatsapp-tts /whatsapp-tts
ENV PORT=8080
EXPOSE 8080
CMD ["/whatsapp-tts"]
```

## 8. Notas de Implementación

- El paquete `internal` no puede ser importado por otros módulos
- Los adapters implementan interfaces definidas en `internal`
- Mantener el dominio libre de dependencias de infraestructura
- Usar inyección de dependencias para facilitar testing