# SPEC-009: Observability

## 1. Propósito

Esta especificación define los requisitos de observabilidad del sistema: logging estructurado, tracing, y métricas.

## 2. Estructura de Logging

### 2.1 Formato JSON estructurado

```go
package logging

type LogEntry struct {
    Timestamp   string                 `json:"timestamp"`
    Level       string                 `json:"level"`
    Message     string                 `json:"message"`
    RequestID   string                 `json:"request_id,omitempty"`
    TraceID     string                 `json:"trace_id,omitempty"`
    SpanID      string                 `json:"span_id,omitempty"`
    Service     string                 `json:"service"`
    Environment string                 `json:"environment"`
    Fields      map[string]interface{} `json:"fields,omitempty"`
}

func (e *LogEntry) String() string {
    b, _ := json.Marshal(e)
    return string(b)
}
```

### 2.2 Niveles de Log

| Nivel | Uso |
|-------|-----|
| DEBUG | Información detallada de debug |
| INFO | Eventos normales del sistema |
| WARN | Situaciones anómalas pero manejables |
| ERROR | Errores que requieren atención |
| FATAL | Errores críticos que detenienen el servicio |

### 2.3 Logging en Stages

```go
func (s *IngestionStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
    logger := logging.FromContext(ctx)
    
    logger.Info("starting_ingestion",
        "stage", "ingestion",
        "request_id", getRequestID(ctx),
    )
    
    // Processing...
    
    logger.Info("ingestion_completed",
        "stage", "ingestion",
        "message_count", len(messages),
    )
    
    return output, nil
}
```

## 3. Request Tracing

### 3.1 IDs de Trazabilidad

```go
type TraceContext struct {
    RequestID string  // ID único por solicitud
    TraceID   string  // ID de trace completo (padre-hijo)
    SpanID    string  // ID del span actual
}
```

### 3.2 Propagation

```go
func ExtractTrace(ctx context.Context, headers http.Header) context.Context {
    traceID := headers.Get("X-Trace-ID")
    if traceID == "" {
        traceID = generateTraceID()
    }
    
    return context.WithValue(ctx, TraceKey, &TraceContext{
        RequestID: generateRequestID(),
        TraceID:   traceID,
        SpanID:    generateSpanID(),
    })
}

func InjectTrace(ctx context.Context, headers http.Header) {
    if tc, ok := ctx.Value(TraceKey).(*TraceContext); ok {
        headers.Set("X-Trace-ID", tc.TraceID)
        headers.Set("X-Span-ID", tc.SpanID)
    }
}
```

### 3.3 Logging con Trace

```go
func (s *Pipeline) Execute(ctx context.Context, input interface{}) (*PipelineContext, error) {
    tc := getTraceContext(ctx)
    
    logger := logging.FromContext(ctx)
    logger.Info("pipeline_start",
        "trace_id", tc.TraceID,
        "request_id", tc.RequestID,
    )
    
    // Execute stages...
    
    logger.Info("pipeline_end",
        "trace_id", tc.TraceID,
        "success", err == nil,
    )
}
```

## 4. Métricas

### 4.1 Métricas de Negocio

| Métrica | Tipo | Descripción |
|---------|------|-------------|
| messages_received | Counter | Mensajes recibidos |
| messages_processed | Counter | Mensajes procesados exitosamente |
| messages_failed | Counter | Mensajes que fallaron |
| response_time_ms | Histogram | Tiempo de respuesta |

### 4.2 Métricas Técnicas

| Métrica | Tipo | Descripción |
|---------|------|-------------|
| http_requests_total | Counter | Total de requests HTTP |
| http_request_duration_ms | Histogram | Duración de requests |
| tts_generation_duration_ms | Histogram | Tiempo de generación TTS |
| whatsapp_api_calls_total | Counter | Llamadas a API de WhatsApp |
| whatsapp_api_errors_total | Counter | Errores de API de WhatsApp |

### 4.3 Implementación

```go
package metrics

var (
    MessagesReceived = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "whatsapp_tts_messages_received_total",
            Help: "Total number of messages received",
        },
    )
    
    MessagesProcessed = prometheus.NewCounter(
        prometheus.CounterOpts{
            Name: "whatsapp_tts_messages_processed_total",
            Help: "Total number of messages processed successfully",
        },
    )
    
    ResponseTime = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "whatsapp_tts_response_time_ms",
            Help:    "Response time in milliseconds",
            Buckets: []float64{100, 250, 500, 1000, 2500, 5000},
        },
    )
)
```

## 5. Health Checks

### 5.1 Endpoint de Health

```
GET /health
```

Response:
```json
{
  "status": "healthy",
  "timestamp": "2025-03-14T10:00:00Z",
  "checks": {
    "webhook": "ok",
    "tts": "ok",
    "whatsapp_api": "ok"
  }
}
```

### 5.2 Readiness vs Liveness

- **Liveness**: El servicio está vivo (/)
- **Readiness**: El servicio puede recibir tráfico (/ready)

## 6. Notas de Implementación

- Usar logger estructurado (zap, zerolog, o similar)
- Incluir siempre RequestID y TraceID en logs
- Configurar muestreo de traces para alto volumen
- Exponer métricas en formato Prometheus