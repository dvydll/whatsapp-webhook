## ADDED Requirements

### Requirement: El sistema define IngestionStage
El sistema DEBE definir un stage que reciba HTTP Request y extraiga el body.

#### Scenario: Procesar POST request
- **WHEN** IngestionStage.Process recibe *http.Request
- **THEN** retorna *pipeline.RawEvent con el body del request

#### Scenario: Ignorar GET request
- **WHEN** IngestionStage.Process recibe GET
- **THEN** retorna nil (placeholder para verificación)

### Requirement: El sistema define NormalizationStage
El sistema DEBE parsear el payload JSON y convertir a UserMessage.

#### Scenario: Normalizar mensaje de WhatsApp
- **WHEN** NormalizationStage.Process recibe *RawEvent con JSON válido
- **THEN** retorna *domain.UserMessage con los datos extraídos

#### Scenario: Payload inválido
- **WHEN** el JSON no tiene mensajes
- **THEN** retorna error

### Requirement: El sistema define ResponseGenerationStage (STUB)
El sistema DEBE generar una respuesta de texto.

#### Scenario: Generar respuesta placeholder
- **WHEN** ResponseGenerationStage.Process recibe *UserMessage
- **THEN** retorna *domain.ResponseMessage con texto "Message received. Generating audio response."

### Requirement: El sistema define TTSGenerationStage (STUB)
El sistema DEBE generar un AudioAsset dummy.

#### Scenario: Generar audio stub
- **WHEN** TTSGenerationStage.Process recibe *ResponseMessage
- **THEN** retorna *domain.AudioAsset con datos placeholder (no audio real)

### Requirement: El sistema define DeliveryStage (STUB)
El sistema DEBE tener un stage de entrega que no hace nada.

#### Scenario: Stage stub
- **WHEN** DeliveryStage.Process recibe *AudioAsset
- **THEN** retorna nil, nil (no entrega nada)
