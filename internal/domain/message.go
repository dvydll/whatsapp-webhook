package domain

import (
	"time"
)

// UserMessage representa un mensaje entrante del usuario.
type UserMessage struct {
	MessageID         string
	ExternalMessageID string
	UserID            string
	Channel           Channel
	Text              string
	Timestamp         time.Time
	Metadata          MessageMetadata
}

// MessageMetadata contiene información adicional del mensaje.
type MessageMetadata struct {
	UserPhone         string
	UserProfileName   string
	BusinessPhoneID   string
	BusinessAccountID string
	MessageType       string
	RawPayload        []byte
}

// Channel representa el canal de mensajería.
type Channel string

const (
	ChannelWhatsApp Channel = "whatsapp"
	ChannelTelegram Channel = "telegram"
	ChannelWebChat  Channel = "webchat"
)

// NewUserMessage crea un nuevo UserMessage.
func NewUserMessage(
	messageID string,
	externalMessageID string,
	userID string,
	channel Channel,
	text string,
	timestamp time.Time,
	metadata MessageMetadata,
) *UserMessage {
	return &UserMessage{
		MessageID:         messageID,
		ExternalMessageID: externalMessageID,
		UserID:            userID,
		Channel:           channel,
		Text:              text,
		Timestamp:         timestamp,
		Metadata:          metadata,
	}
}
