package stages

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/whatsapp-tts/internal/domain"
	"github.com/whatsapp-tts/internal/pipeline"
)

// NormalizationStage normaliza eventos raw a UserMessage.
type NormalizationStage struct{}

func NewNormalizationStage() *NormalizationStage {
	return &NormalizationStage{}
}

func (s *NormalizationStage) Name() string { return "normalization" }

func (s *NormalizationStage) CanProcess(input interface{}) bool {
	_, ok := input.(*pipeline.RawEvent)
	return ok
}

func (s *NormalizationStage) Process(ctx context.Context, input interface{}) (interface{}, error) {
	event, ok := input.(*pipeline.RawEvent)
	if !ok {
		return nil, pipeline.ErrInvalidInput
	}

	var webhook WhatsAppWebhook
	if err := json.Unmarshal(event.Payload, &webhook); err != nil {
		return nil, err
	}

	if len(webhook.Entry) == 0 || len(webhook.Entry[0].Changes) == 0 {
		return nil, pipeline.ErrInvalidInput
	}

	changes := webhook.Entry[0].Changes[0]
	value := changes.Value

	if len(value.Messages) == 0 {
		return nil, pipeline.ErrInvalidInput
	}

	msg := value.Messages[0]
	ts, _ := strconv.ParseInt(msg.Timestamp, 10, 64)

	userMsg := domain.NewUserMessage(
		generateID(),
		msg.ID,
		msg.From,
		domain.ChannelWhatsApp,
		msg.Text.Body,
		time.Unix(ts, 0),
		domain.MessageMetadata{
			UserPhone:         msg.From,
			BusinessPhoneID:   value.Metadata.PhoneNumberID,
			BusinessAccountID: webhook.Entry[0].ID,
		},
	)

	return userMsg, nil
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

func generateID() string {
	return "msg-" + time.Now().Format("20060102150405")
}
