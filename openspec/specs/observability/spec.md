## ADDED Requirements

### Requirement: El sistema provee Logger interface
El sistema DEBE definir una interfaz de logging con Info, Warn, Error, Debug.

#### Scenario: Interface Logger
- **WHEN** se usa Logger interface
- **THEN** tiene métodos: Info, Warn, Error, Debug

### Requirement: El sistema provee StdLogger implementación
El sistema DEBE implementar un logger básico usando log estándar.

#### Scenario: Log con campos
- **WHEN** logger.Info("message", F("key", "value"))
- **THEN** output: "[INFO] message key=value"

### Requirement: El sistema provee helper F para campos
El sistema DEBE facilitar la creación de campos de logging.

#### Scenario: Crear campo
- **WHEN** F("key", "value")
- **THEN** retorna Field{Key: "key", Value: "value"}
