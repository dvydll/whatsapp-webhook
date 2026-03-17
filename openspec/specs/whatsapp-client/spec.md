## ADDED Requirements

### Requirement: El sistema define WhatsApp Client
El sistema DEBE proporcionar cliente HTTP para WhatsApp Cloud API.

#### Scenario: Crear cliente
- **WHEN** se llama whatsapp.NewClient(config)
- **THEN** retorna *Client con httpClient configurado

### Requirement: El sistema envía mensajes de texto
El sistema DEBE poder enviar mensajes de texto a WhatsApp.

#### Scenario: SendTextMessage
- **WHEN** Client.SendTextMessage(ctx, to, text)
- **THEN** retorna (*SendMessageResponse, error) con message ID

### Requirement: El sistema envía mensajes de audio
El sistema DEBE poder enviar mensajes de audio a WhatsApp.

#### Scenario: SendAudioMessage
- **WHEN** Client.SendAudioMessage(ctx, to, audioURL)
- **THEN** retorna (*SendMessageResponse, error) con message ID

### Requirement: El sistema maneja errores de API
El sistema DEBE retornar error cuando WhatsApp API retorna status >= 400.

#### Scenario: Error de API
- **WHEN** WhatsApp retorna status 400+
- **THEN** retorna error con mensaje descriptivo
