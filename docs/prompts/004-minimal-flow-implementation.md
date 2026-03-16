You are a senior Go backend engineer continuing a Specification-Driven Development workflow.

The repository already contains:

* architecture specifications
* domain models
* interface contracts
* a compilable project skeleton
* stub pipeline stages

Your task now is to implement the **first real functional version of the system**.

The objective is to make the service capable of completing a **real end-to-end flow** with WhatsApp.

---

# Implementation Goal

Implement the minimal functional flow:

User sends message on WhatsApp
→ WhatsApp webhook triggers service
→ message is parsed and normalized
→ a response text is generated
→ text is converted to speech via TTS
→ audio is converted to WhatsApp-compatible format
→ audio is sent back to the user via WhatsApp Cloud API.

This version does not need advanced features yet.

---

# Functional Requirements

The system must:

1. Receive webhook events from WhatsApp Cloud API.
2. Extract incoming user text messages.
3. Generate a simple response text.
4. Generate audio using the configured TTS service.
5. Convert audio to WhatsApp compatible format (OGG Opus).
6. Send the audio response to the user.

---

# WhatsApp Integration

Integrate with the WhatsApp Cloud API.

The implementation must support:

POST /{PHONE_NUMBER_ID}/messages

Audio message format:

{
"messaging_product": "whatsapp",
"to": "<user_phone>",
"type": "audio",
"audio": {
"link": "<public_audio_url>"
}
}

Implementation must include:

* HTTP client
* request serialization
* error handling
* logging of delivery responses

Webhook handler must parse incoming payloads from WhatsApp message events.

---

# Webhook Parsing

Implement parsing of the WhatsApp webhook payload.

Extract:

* sender phone number
* message text
* message ID

Ignore other message types for now.

If message type is not text, log and ignore.

---

# Response Text Generation

Implement a simple response generator.

Example response:

"Message received. Generating audio response."

Later versions will replace this with real logic.

---

# TTS Integration

Connect to an external TTS service.

Assume the TTS service exposes an HTTP endpoint:

POST /tts

Input:

{
"text": "text to synthesize"
}

Output:

binary WAV file.

Your implementation must:

* send request
* receive WAV audio
* store it temporarily in memory or temp file

---

# Audio Conversion

Convert WAV output to WhatsApp compatible format.

Use ffmpeg.

Required output format:

* container: OGG
* codec: Opus
* sample rate: 16000
* mono audio

Example command:

ffmpeg -i input.wav -c:a libopus -b:a 32k -ar 16000 -ac 1 output.ogg

Encapsulate this inside the AudioProcessor service.

---

# Audio Delivery Strategy

For the first version, audio must be accessible via public URL.

Implementation approach:

1. Save generated OGG file in local directory such as:

/tmp/audio/

2. Serve this directory via HTTP:

/audio/{filename}

3. Use this URL when sending the WhatsApp message.

Example:

https://yourserver/audio/response123.ogg

---

# HTTP Server Additions

The server must expose:

POST /webhook
GET /audio/{filename}

The audio endpoint should serve generated audio files.

---

# Logging Requirements

Log the following events:

* webhook received
* message parsed
* response generated
* TTS request started
* audio conversion completed
* WhatsApp delivery request
* delivery response

Logs should include message IDs when available.

---

# Error Handling

Failures should not crash the server.

Handle errors for:

* invalid webhook payload
* TTS request failure
* ffmpeg failure
* WhatsApp API errors

Log all errors clearly.

---

# Environment Configuration

The service must read the following environment variables:

WHATSAPP_TOKEN
WHATSAPP_PHONE_NUMBER_ID
TTS_ENDPOINT
SERVER_PORT
PUBLIC_BASE_URL

Example:

PUBLIC_BASE_URL=https://myserver.com

This will be used to generate audio URLs.

---

# Expected Behavior

When a user sends a text message on WhatsApp:

1. webhook is triggered
2. service logs event
3. response text is generated
4. audio is synthesized
5. audio converted to OGG Opus
6. audio stored and served
7. WhatsApp API called
8. user receives audio reply

---

# Validation Criteria

The system must:

* build successfully with `go build`
* start with `go run cmd/server/main.go`
* accept webhook requests
* generate audio files
* expose them via /audio
* send audio messages via WhatsApp API

---

# Do Not Implement Yet

* authentication middleware
* message queues
* persistent storage
* LLM integration
* async background workers

Those will be added in later iterations.

Focus only on making the **minimal WhatsApp → TTS → audio reply flow work**.

---

# Deliverables

Update the project with:

* real webhook handler
* real WhatsApp client
* real TTS integration
* audio conversion service
* audio file server
* pipeline wiring

Ensure the service can complete the full end-to-end interaction.
