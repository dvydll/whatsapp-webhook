## ADDED Requirements

### Requirement: El sistema define UserMessage
El sistema DEBE representar mensajes entrantes del usuario.

#### Scenario: Crear UserMessage
- **WHEN** domain.NewUserMessage(messageID, externalMessageID, userID, channel, text, timestamp, metadata)
- **THEN** retorna *UserMessage con todos los campos

### Requirement: El sistema define ResponseMessage
El sistema DEBE representar respuestas generadas.

#### Scenario: Crear ResponseMessage
- **WHEN** domain.NewResponseMessage(targetUser, responseType, text, phoneNumberID)
- **THEN** retorna *ResponseMessage con ResponseID generado

### Requirement: El sistema define AudioAsset
El sistema DEBE representar archivos de audio.

#### Scenario: Crear AudioAsset
- **WHEN** domain.NewAudioAsset(format, codec, binaryData)
- **THEN** retorna *AudioAsset con AudioID generado

### Requirement: El sistema define Channel
El sistema DEBE definir tipos de canal soportados.

#### Scenario: Canales definidos
- **WHEN** se usa domain.ChannelWhatsApp
- **THEN** retorna "whatsapp"

### Requirement: El sistema define ResponseType
El sistema DEBE definir tipos de respuesta.

#### Scenario: Tipos de respuesta
- **WHEN** se usa domain.ResponseTypeText o domain.ResponseTypeAudio
- **THEN** retorna "text" o "audio"
