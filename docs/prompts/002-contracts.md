You are a senior Go backend engineer continuing a **Specification-Driven Development workflow**.

The previous step produced the **high-level system specifications** for a messaging service that integrates with the WhatsApp Cloud API and responds with audio generated via a TTS service.

Your task now is to convert those architectural specifications into **concrete technical contracts for Go implementation**.

Do NOT implement business logic yet.

Your job is to produce:

* interfaces
* domain models
* DTOs
* pipeline contracts
* adapter contracts
* testable specifications

The goal is that a future implementation agent can implement the entire system purely by satisfying these contracts.

---

# System Context

The application is a Go service that:

1. receives WhatsApp webhook events
2. extracts incoming user messages
3. normalizes them into an internal message format
4. processes them through a modular pipeline
5. generates a text response
6. converts the response text to audio via a TTS provider
7. sends the audio back via the WhatsApp API

The architecture must allow future features such as:

* multi-step workflows
* asynchronous processing
* background jobs
* intermediate services (LLM calls, database queries, etc.)

---

# Primary Design Constraints

The system must be designed around:

* **clean architecture**
* **dependency inversion**
* **testability**
* **replaceable adapters**

The core domain must not depend on external APIs.

External systems must be accessed through adapters.

---

# Deliverable 1 — Go Module Definition

Define the base Go module structure.

Specify:

* module name
* Go version
* dependency policy (minimal external dependencies)

Prefer standard library where possible.

---

# Deliverable 2 — Repository Layout

Define a concrete project layout suitable for a production Go service.

Example direction:

cmd/
internal/
domain/
adapters/
services/
pipelines/
contracts/
tests/

Explain the responsibility of each directory.

---

# Deliverable 3 — Domain Models

Define the core domain entities.

At minimum include:

UserMessage

Fields such as:

* MessageID
* UserID
* Channel
* Text
* Timestamp
* Metadata

ResponseMessage

Fields such as:

* ResponseID
* TargetUser
* ResponseType
* Text
* AudioURL
* CreatedAt

AudioAsset

Fields such as:

* Format
* Codec
* BinaryData or URL
* Duration

These must be pure Go structs with no external dependencies.

---

# Deliverable 4 — Pipeline Contracts

Define interfaces for each stage of the message pipeline.

Examples:

EventReceiver

Input: HTTP request
Output: raw event

MessageNormalizer

Input: raw event
Output: UserMessage

MessageProcessor

Input: UserMessage
Output: ResponseMessage

TTSEngine

Input: text
Output: AudioAsset

DeliveryAdapter

Input: ResponseMessage
Output: delivery result

Each stage must be defined as a Go interface.

---

# Deliverable 5 — WhatsApp Adapter Contract

Define the interface that integrates with the WhatsApp Cloud API.

Responsibilities:

* verify webhook
* parse message events
* send message responses
* send audio messages

Define request and response DTOs for this adapter.

The adapter must isolate all WhatsApp-specific payload formats.

---

# Deliverable 6 — TTS Provider Contract

Define a generic interface for TTS providers.

Example interface conceptually:

GenerateAudio(text, voice, language) → AudioAsset

The interface must allow future implementations such as:

* StyleTTS
* OpenAI TTS
* Local models

---

# Deliverable 7 — Audio Processing Contract

Define a service responsible for:

* audio format conversion
* codec handling
* preparing audio compatible with WhatsApp

Example operations:

ConvertToWhatsAppFormat
NormalizeSampleRate
EncodeOpus

This service must abstract the use of tools such as ffmpeg.

---

# Deliverable 8 — Error Model

Define a consistent error structure.

Include:

* error code
* error message
* stage of failure
* retryable flag

Errors must propagate across pipeline stages in a structured way.

---

# Deliverable 9 — Contract Tests

Define test scaffolding that verifies that implementations satisfy the interfaces.

Examples:

* pipeline stage tests
* adapter compliance tests
* TTS provider contract tests

Do not implement full tests yet, only the structure and expectations.

---

# Deliverable 10 — Example End-to-End Flow

Define the canonical flow as pseudo-code:

Webhook → Normalization → Processing → TTS → Delivery

Show how interfaces interact without implementing concrete logic.

---

# Critical Requirement

All definitions must allow the following future extensions without breaking contracts:

* background processing queues
* long-running tasks
* streaming audio generation
* multi-channel messaging
* conversational memory

---

# Output Format

Produce:

1. Go interface definitions
2. domain structs
3. repository layout
4. contract test structure
5. pipeline interaction diagram

Do not implement production logic yet.

The objective is to make the system **fully spec-driven and implementation-ready**.
