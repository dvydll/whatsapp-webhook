You are a senior Go systems architect evolving an existing messaging service into a **modular agent-based architecture**.

The repository already contains:

* a working WhatsApp → TTS → audio reply flow
* a clean architecture structure
* a processing pipeline
* service interfaces for TTS and delivery

Your task is to redesign the processing layer so that the system supports **agent skills and modular processing stages**.

The objective is to make the system capable of executing **pluggable capabilities** during message processing.

Do not break the current working flow.

---

# Core Architectural Goal

Transform the message processing pipeline into an **agent execution engine**.

The system must allow multiple skills to participate in generating the final response.

Example future capabilities:

* database queries
* API calls
* LLM reasoning
* workflow execution
* knowledge retrieval
* multi-step tasks

---

# Message Lifecycle

The lifecycle must remain consistent:

User message
→ WhatsApp webhook
→ message normalization
→ agent processing engine
→ response generation
→ TTS
→ audio delivery.

The **agent processing engine** replaces the current simple message processor.

---

# Agent Execution Model

Define a model where an agent processes a message by executing a sequence of skills.

Example conceptual pipeline:

NormalizeMessage
↓
AgentEngine
↓
SkillRouter
↓
SkillExecution
↓
ResponseComposer

Skills may:

* modify the response
* add context
* trigger asynchronous jobs
* request additional processing

---

# Skill Interface

Define a generic skill interface.

Each skill should receive a processing context and return a result.

Example conceptual structure:

Skill

Inputs:

* message context
* agent state
* conversation metadata

Outputs:

* skill result
* updated context
* optional response fragment

Skills must be stateless and easily replaceable.

---

# Skill Types

Define at least three conceptual categories of skills.

Input skills
Responsible for enriching the message context.

Examples:

* entity extraction
* language detection
* metadata enrichment

Processing skills
Perform core reasoning or operations.

Examples:

* LLM call
* database lookup
* external API call

Output skills
Responsible for shaping the final response.

Examples:

* text generation
* formatting
* response summarization

---

# Agent Context Model

Define a shared context object that flows through the skill pipeline.

It should contain:

* user message
* conversation ID
* working memory
* accumulated results
* response builder

This context object will allow skills to collaborate.

---

# Skill Registry

Create a skill registry responsible for:

* registering available skills
* resolving skills by name
* constructing the execution pipeline

Skills must be dynamically configurable.

---

# Agent Engine

Implement an agent engine responsible for:

* executing skills sequentially
* managing context propagation
* handling errors
* determining when execution ends

The engine should support:

* deterministic pipelines
* dynamic skill routing in the future

---

# Response Composition

Multiple skills may contribute to the final response.

Define a ResponseBuilder object responsible for:

* collecting response fragments
* prioritizing responses
* selecting final output

For now the system should still produce a single response text.

---

# Compatibility Requirement

The existing minimal behavior must remain intact.

If no additional skills are registered, the system must behave exactly as before:

User message → fixed response text → TTS → audio reply.

---

# Future Capability Requirements

The new architecture must support future capabilities such as:

* LLM reasoning agents
* tool usage
* multi-step workflows
* background job scheduling
* streaming responses
* conversational memory

The architecture should anticipate these without implementing them yet.

---

# Repository Changes

Introduce new directories such as:

internal/agent/
internal/skills/
internal/context/
internal/registry/

Update the message processing service to use the agent engine instead of the previous processor.

---

# Deliverables

Produce:

1. skill interface definition
2. agent context model
3. agent engine design
4. skill registry implementation
5. example skill (simple response generator)
6. integration with existing pipeline

The final result must keep the WhatsApp → TTS → audio response flow operational while introducing the new modular architecture.

Do not implement LLMs or external tools yet.

Focus on building the **agent execution framework** that will support them later.
