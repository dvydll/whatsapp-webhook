package observability

import (
	"log"
	"os"
	"strings"
)

// Logger es una interfaz para logging estructurado.
type Logger interface {
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
}

// Field representa un campo de logging.
type Field struct {
	Key   string
	Value any
}

// F crea un campo de logging.
func F(key string, value any) Field {
	return Field{Key: key, Value: value}
}

// StdLogger es una implementación simple de Logger usando log estándar.
type StdLogger struct {
	logger *log.Logger
}

// NewStdLogger crea un nuevo logger estándar.
func NewStdLogger() *StdLogger {
	return &StdLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *StdLogger) format(msg string, fields ...Field) string {
	var line strings.Builder
	line.WriteString(msg)

	if len(fields) > 0 {
		for _, f := range fields {
			line.WriteString(" " + f.Key + "=" + toString(f.Value))
		}
	}

	return line.String()
}

func toString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	default:
		return "unknown"
	}
}

// Info logs a message at INFO level.
func (l *StdLogger) Info(msg string, fields ...Field) {
	l.logger.Println("[INFO] " + l.format(msg, fields...))
}

// Warn logs a message at WARN level.
func (l *StdLogger) Warn(msg string, fields ...Field) {
	l.logger.Println("[WARN] " + l.format(msg, fields...))
}

// Error logs a message at ERROR level.
func (l *StdLogger) Error(msg string, fields ...Field) {
	l.logger.Println("[ERROR] " + l.format(msg, fields...))
}

// Debug logs a message at DEBUG level.
func (l *StdLogger) Debug(msg string, fields ...Field) {
	l.logger.Println("[DEBUG] " + l.format(msg, fields...))
}
