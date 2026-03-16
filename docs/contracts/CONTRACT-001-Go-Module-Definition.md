# CONTRACT-001: Go Module Definition

## 1. Módulo Go

```go
// go.mod
module github.com/whatsapp-tts

go 1.21
```

## 2. Nombre del Módulo

**Nombre**: `github.com/whatsapp-tts`
**Versión Go**: `1.21` (mínimo)

## 3. Política de Dependencias

### Dependencias Mínimas

Solo usar dependencias externas cuando sea necesario. Preferir estándar library:

| Categoría | Preferir | Alternativa si es necesario |
|-----------|----------|----------------------------|
| HTTP Server | net/http | chi, gorilla/mux |
| Logging | log/slog | zerolog, zap |
| JSON | encoding/json | jsoniter |
| Testing | testing | testify |
| Validation | manual o go-playground/validator | - |
| Config | os.Getenv, flag | envconfig |

### Dependencias Esperadas

```go
// Dependencias que很可能 serán necesarias:
import (
    // Testing
    "testing"
    
    // HTTP
    "net/http"
    
    // JSON
    "encoding/json"
    
    // Context
    "context"
    
    // Time
    "time"
    
    // Errors
    "errors"
    "fmt"
)
```

## 4. Estructura de Módulos

El proyecto puede usar un único módulo para servicios pequeños/medios:

```
github.com/whatsapp-tts/
├── go.mod
├── cmd/
├── internal/
├── adapters/
└── pkg/
```

## 5. Notas

- No agregar dependencias sin justificación clara
- Mantener el módulo liviano para facilitar testing
- Usar Go 1.21+ para soporte de generics y mejor concurrency