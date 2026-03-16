package webhook

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/whatsapp-tts/internal/adapters/whatsapp"
	"github.com/whatsapp-tts/internal/config"
	"github.com/whatsapp-tts/internal/observability"
	"github.com/whatsapp-tts/internal/pipeline"
	"github.com/whatsapp-tts/internal/pipeline/stages"
)

type Handler struct {
	pipeline       *pipeline.Pipeline
	whatsappClient *whatsapp.Client
	logger         observability.Logger
	verifyToken    string
}

func NewHandler(cfg *config.Config, logger observability.Logger) *Handler {
	p := pipeline.NewPipeline(
		stages.NewIngestionStage(cfg.WhatsAppVerifyToken),
		stages.NewNormalizationStage(),
		stages.NewResponseGenerationStage(),
		stages.NewTTSGenerationStage(),
		stages.NewDeliveryStage(),
	)

	whatsappClient := whatsapp.NewClient(whatsapp.Config{
		PhoneNumberID: cfg.WhatsAppPhoneNumberID,
		AccessToken:   cfg.WhatsAppAccessToken,
		BaseURL:       cfg.WhatsAppAPIURL,
		Version:       "v25.0",
	})

	return &Handler{
		pipeline:       p,
		whatsappClient: whatsappClient,
		logger:         logger,
		verifyToken:    cfg.WhatsAppVerifyToken,
	}
}

type WhatsAppWebhook struct {
	Object string         `json:"object"`
	Entry  []WebhookEntry `json:"entry"`
}

type WebhookEntry struct {
	ID      string   `json:"id"`
	Changes []Change `json:"changes"`
}

type Change struct {
	Value WebhookValue `json:"value"`
	Field string       `json:"field"`
}

type WebhookValue struct {
	MessagingProduct string    `json:"messaging_product"`
	Metadata         Metadata  `json:"metadata"`
	Messages         []Message `json:"messages"`
}

type Metadata struct {
	DisplayPhoneNumber string `json:"display_phone_number"`
	PhoneNumberID      string `json:"phone_number_id"`
}

type Message struct {
	From      string   `json:"from"`
	ID        string   `json:"id"`
	Timestamp string   `json:"timestamp"`
	Type      string   `json:"type"`
	Text      TextBody `json:"text"`
}

type TextBody struct {
	Body string `json:"body"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("webhook_request",
		observability.F("method", r.Method),
		observability.F("path", r.URL.Path),
	)

	// Test endpoint for simulating webhooks (must be before verification check)
	if r.URL.Path == "/test-webhook" && r.Method == "GET" {
		h.testWebhook(w, r)
		return
	}

	// Handle GET for webhook verification
	if r.Method == "GET" {
		h.verifyWebhook(w, r)
		return
	}

	// Handle POST for webhook events
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse webhook payload
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed_to_read_body", observability.F("error", err.Error()))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var webhook WhatsAppWebhook
	if err := json.Unmarshal(body, &webhook); err != nil {
		h.logger.Error("failed_to_parse_webhook", observability.F("error", err.Error()))
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Process each message
	for _, entry := range webhook.Entry {
		for _, change := range entry.Changes {
			value := change.Value
			h.logger.Info("processing_messages",
				observability.F("message_count", len(value.Messages)),
			)

			for _, msg := range value.Messages {
				h.processMessage(msg, value.Metadata)
			}
		}
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) verifyWebhook(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	h.logger.Info("webhook_verification",
		observability.F("mode", mode),
		observability.F("token_matches", token == h.verifyToken),
	)

	if mode == "subscribe" && token == h.verifyToken {
		h.logger.Info("webhook_verified")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		return
	}

	http.Error(w, "Verification failed", http.StatusForbidden)
}

func (h *Handler) processMessage(msg Message, metadata Metadata) {
	h.logger.Info("message_received",
		observability.F("message_id", msg.ID),
		observability.F("from", msg.From),
		observability.F("type", msg.Type),
	)

	// Ignore non-text messages for now
	if msg.Type != "text" {
		h.logger.Info("ignoring_non_text_message", observability.F("type", msg.Type))
		return
	}

	userMessage := msg.Text.Body
	h.logger.Info("user_message", observability.F("text", userMessage))

	// Generate response (placeholder)
	responseText := "Message received. Generating audio response."

	// TODO: Generate audio with TTS and convert to OGG
	// For now, send a text response to test the flow

	h.logger.Info("sending_response", observability.F("text", responseText))

	ctx := context.Background()
	resp, err := h.whatsappClient.SendTextMessage(ctx, msg.From, responseText)
	if err != nil {
		h.logger.Error("failed_to_send_response", observability.F("error", err.Error()))
		return
	}

	h.logger.Info("response_sent", observability.F("whatsapp_message_id", resp.Messages[0].ID))
}

// TestWebhook simulates a WhatsApp webhook for testing
func (h *Handler) testWebhook(w http.ResponseWriter, r *http.Request) {
	// Simular un mensaje entrante de WhatsApp
	testMsg := Message{
		From:      "34685107027", // Tu número
		ID:        "test_msg_123",
		Timestamp: "1741977600",
		Type:      "text",
		Text:      TextBody{Body: "Test message from simulation"},
	}

	testMetadata := Metadata{
		DisplayPhoneNumber: "+1234567890",
		PhoneNumberID:      "1052703781258489",
	}

	h.logger.Info("test_webhook_triggered")
	h.processMessage(testMsg, testMetadata)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status": "test_completed"}`))
}

// PipelineBuilder creates a pipeline with stub stages (for backward compatibility)
func PipelineBuilder() *pipeline.Pipeline {
	return pipeline.NewPipeline(
		stages.NewIngestionStage(""),
		stages.NewNormalizationStage(),
		stages.NewResponseGenerationStage(),
		stages.NewTTSGenerationStage(),
		stages.NewDeliveryStage(),
	)
}
