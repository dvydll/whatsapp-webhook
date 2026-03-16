You are a senior Go backend engineer continuing a Specification-Driven Development workflow.

The repository already contains:

* architecture specifications
* domain models
* interface contracts
* pipeline definitions

Your task now is to **generate the full project skeleton** for the Go application that satisfies these contracts.

Do NOT implement real business logic yet.

Instead:

* create the project structure
* define interfaces
* create stub implementations
* wire the application entry point
* make the project compile successfully

The goal is that a developer can immediately begin implementing features inside a well-structured codebase.

---

# System Overview

The application is a Go service that integrates with the WhatsApp Cloud API.

The basic workflow is:

1. receive webhook event from WhatsApp
2. normalize incoming message
3. process the message through a pipeline
4. generate response text
5. generate audio using TTS
6. send the audio response back to the user

Future features must be supported such as:

* additional processing stages
* asynchronous jobs
* database integrations
* LLM integrations
* delayed responses

The skeleton must support these expansions without structural changes.

---

# Goal

Produce a **compilable Go project skeleton** with:

* clear module structure
* dependency inversion
* replaceable adapters
* clean architecture separation

---

# Go Module Setup

Create a go.mod with:

* Go version ≥ 1.22
* minimal external dependencies
* standard library preferred

---

# Project Structure

Generate the following directory structure:

cmd/server/
main.go

internal/
app/
app.go

```
pipeline/
    pipeline.go
    stage.go

webhook/
    handler.go

domain/
    message.go
    response.go
    audio.go

services/
    message_processor.go
    tts_service.go
    audio_service.go

adapters/
    whatsapp/
        client.go
        webhook_parser.go

    tts/
        provider.go

delivery/
    dispatcher.go

config/
    config.go

observability/
    logger.go
```

tests/
contracts/

Explain the purpose of each directory.

---

# Interfaces to Implement

Create Go interfaces for:

EventReceiver
MessageNormalizer
MessageProcessor
TTSEngine
AudioProcessor
DeliveryAdapter

These must match the previously defined contracts.

---

# Stub Implementations

Create placeholder implementations that return predictable results.

Example:

TestMessageProcessor

Input:
UserMessage

Output:
ResponseMessage with text:

"Message received. Generating audio response."

These stubs allow the application to run end-to-end.

---

# HTTP Server

Create a basic HTTP server that:

* exposes /webhook endpoint
* receives POST requests from WhatsApp
* forwards them into the pipeline

Do not implement webhook verification yet; just leave a TODO marker.

---

# Pipeline Wiring

Create a simple pipeline runner that executes:

Webhook → Normalizer → Processor → TTS → Delivery

Each stage should be injected via interfaces.

The pipeline should support adding new stages later.

---

# Configuration

Add a minimal configuration loader that reads:

* WhatsApp API token
* WhatsApp phone number ID
* TTS endpoint
* server port

Configuration can be read from environment variables.

---

# Logging

Create a minimal structured logger interface.

Use standard library logging for now.

The logger must support:

* info logs
* error logs
* request IDs

---

# Example Run Flow

When the server receives a webhook message:

1. parse incoming payload
2. convert to UserMessage
3. run through pipeline
4. generate placeholder text
5. call stub TTS
6. call stub delivery adapter

The system should log each step.

---

# Output Requirements

Generate:

1. all directories
2. all Go files
3. interface definitions
4. stub implementations
5. pipeline runner
6. server entrypoint

The project must:

* compile successfully with `go build`
* run with `go run cmd/server/main.go`
* expose `/webhook`
* log the pipeline steps when a request arrives

Do not implement real WhatsApp API calls or real TTS calls yet.

Those will be implemented later.
