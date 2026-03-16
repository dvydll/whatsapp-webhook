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
│   │   └── error.go
│   │
│   ├── pipeline/               # Motor de pipeline
│   │   ├── pipeline.go
│   │   ├── stage.go
│   │   └── context.go
│   │
│   ├── handler/                # Handlers HTTP
│   │   ├── webhook.go
│   │   └── health.go
│   │
│   ├── config/                 # Configuración
│   │   └── config.go
│   │
│   └── logging/                # Logging
│       └── logger.go
│
├── adapters/                   # Implementaciones de contratos externos
│   ├── tts/                    # Proveedores TTS
│   │   ├── provider.go         # Interfaz TTSProvider
│   │   ├── styletts/
│   │   └── factory.go
│   │
│   ├── delivery/               # Adapters de entrega
│   │   ├── adapter.go         # Interfaz DeliveryAdapter
│   │   └── whatsapp/
│   │
│   └── audio/                 # Procesamiento de audio
│       └── processor.go
│
├── pkg/                        # Paquetes públicos reutilizables
│   └── utils/
│
├── contracts/                  # Contratos técnicos (esta carpeta)
│
├── specs/                      # Especificaciones del sistema
│
├── worklogs/                   # Decisiones de trabajo
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
├── error.go         # Domain errors
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
└── errors.go        # Pipeline errors
```

### internal/handler/
- Manejo de HTTP requests
- Separación de transporte/dominio

**Contenido esperado**:
```
internal/handler/
├── webhook.go       # Webhook handler
├── health.go        # Health check
└── middleware.go    # HTTP middleware
```

### internal/config/
- Carga de configuración
- Validación de config

### internal/logging/
- Logging estructurado

### adapters/tts/
- Implementaciones de proveedores TTS
- Cada proveedor en su propio paquete

### adapters/delivery/
- Implementaciones de adapters de entrega

### adapters/audio/
- Procesamiento de audio

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