## ADDED Requirements

### Requirement: El sistema expone endpoint GET /webhook para verificación
El sistema DEBE verificar el token de WhatsApp y retornar el challenge.

#### Scenario: Verificación exitosa
- **WHEN** GET /webhook?hub.mode=subscribe&hub.verify_token=<token>&hub.challenge=<valor>
- **THEN** retorna 200 con el valor de hub.challenge

#### Scenario: Verificación fallida
- **WHEN** GET /webhook con token inválido
- **THEN** retorna 403 Forbidden

### Requirement: El sistema procesa POST /webhook con mensajes
El sistema DEBE parsear el payload JSON y extraer mensajes.

#### Scenario: Procesar mensaje de texto
- **WHEN** POST /webhook con payload JSON conteniendo message.type="text"
- **THEN** extrae: from, id, timestamp, text.body

#### Scenario: Ignorar mensajes no-texto
- **WHEN** POST /webhook con message.type diferente a "text"
- **THEN** ignora el mensaje y continua

### Requirement: El sistema envía respuesta de texto
El sistema DEBE enviar un mensaje de texto al usuario (NO audio - v1).

#### Scenario: Enviar respuesta de texto
- **WHEN** se procesa un mensaje de texto entrante
- **THEN** envía mensaje de texto con contenido "Message received. Generating audio response."

### Requirement: El sistema provee endpoint de test
El sistema DEBE exponer GET /test-webhook para simular mensajes.

#### Scenario: Test webhook
- **WHEN** GET /test-webhook
- **THEN** simula un mensaje y ejecuta el flujo completo
