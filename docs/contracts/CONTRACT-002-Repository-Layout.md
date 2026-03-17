# CONTRACT-002: Repository Layout

## 1. Estructura de Directorios

```
github.com/whatsapp-tts/
├── cmd/
│   └── whatsapp-tts/           # Entry point, main.go
│
├── internal/                   # Paquetes privados (no exportables)
│   ├── domain/                 # Entidades de dominio puras
│   │   ├── message.go
│   │   ├── response.go
│   │   ├── audio.go
│   │   └── types.go
│   │
│   ├── pipeline/               # Motor de pipeline
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
│   ├── webhook/                # Handlers HTTP
│   │   └── handler.go
│   │
│   ├── config/                 # Configuración
│   │   └── config.go
│   │
│   ├── logging/                # Logging
│   │   └── logger.go
│   │
│   ├── observability/          # Observabilidad
│   │   └── logger.go
│   │
│   └── adapters/               # Implementaciones de contratos externos
│       ├── tts/                # Proveedores TTS
│       │   ├── provider.go     # Interfaz TTSProvider
│       │   ├── styletts/
│       │   └── factory.go
│       │
│       ├── delivery/           # Adapters de entrega
│       │   ├── adapter.go      # Interfaz DeliveryAdapter
│       │   └── whatsapp/
│       │
│       └── whatsapp/            # Cliente WhatsApp Cloud API
│           ├── client.go
│           └── client_test.go
│
├── contracts/                  # Contratos técnicos (esta carpeta)
│
├── specs/                      # Especificaciones del sistema
│
├── logs/worklogs/              # Decisiones de trabajo
│
├── scripts/                    # Scripts utilitarios
│
├── docs/                       # Documentación
│
├── go.mod
└── README.md
```

## 2. Responsabilidad de Cada Directorio

### cmd/whatsapp-tts/
- Punto de entrada de la aplicación
- Inicialización de dependencias
- Arranque del servidor HTTP

**Contenido esperado**:
```
cmd/whatsapp-tts/
├── main.go
└── wire.go  # si se usa wire para DI
```

### internal/domain/
- Entidades de dominio puras
- Sin dependencias externas
- Types, interfaces de dominio

**Contenido esperado**:
```
internal/domain/
├── message.go       # UserMessage, MessageMetadata
├── response.go      # ResponseMessage, ResponseType
├── audio.go         # AudioAsset, AudioFormat
└── types.go         # Channel, ContentType, etc.
```

### internal/pipeline/
- Orquestación de stages
- Contexto de pipeline
- Integración de componentes

**Contenido esperado**:
```
internal/pipeline/
├── pipeline.go      # Pipeline execution
├── stage.go         # Stage interface
├── context.go       # PipelineContext
├── stages/          # Implementaciones de stages
│   ├── ingestion.go
│   ├── normalization.go
│   ├── response_generation.go
│   ├── tts_generation.go
│   └── delivery.go
└── errors.go        # Pipeline errors
```

### internal/webhook/
- Manejo de HTTP requests de webhook
- Separación de transporte/dominio

**Contenido esperado**:
```
internal/webhook/
├── handler.go       # Webhook handler
└── handler_test.go  # Tests con httptest
```

### internal/config/
- Carga de configuración
- Validación de config

### internal/logging/ y internal/observability/
- Logging estructurado

### internal/adapters/tts/
- Implementaciones de proveedores TTS
- Cada proveedor en su propio paquete

### internal/adapters/delivery/
- Implementaciones de adapters de entrega

### internal/adapters/whatsapp/
- Cliente HTTP directo para WhatsApp Cloud API
- Envío de mensajes de texto y audio

### contracts/
- Contratos técnicos (este directorio)
- Define interfaces que deben implementar los adapters

## 3. Reglas de visibility

- `internal/*` - No importable desde fuera del módulo
- `adapters/*` - Implementaciones de contratos
- `pkg/*` - Reutilizable externamente si se necesita
- `domain/*` - Solo tipos y constantes, sin lógica

## 4. Testing

Tests junto con código (mismo paquete) o en `*_test.go`:

```
internal/domain/
├── message.go
├── message_test.go
```

## 5. Convenciones

- Un paquete por directorio
- Nombres en lowercase (snake_case para archivos)
- Test files: `*_test.go`
- Interfaces: nombre descriptivo, no I prefix