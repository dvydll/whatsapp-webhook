package domain

import (
	"testing"
	"time"
)

func TestNewUserMessage(t *testing.T) {
	metadata := MessageMetadata{
		UserPhone:       "+1234567890",
		BusinessPhoneID: "123456789",
	}

	msg := NewUserMessage(
		"msg-123",
		"wamid.abc",
		"user-456",
		ChannelWhatsApp,
		"Hello world",
		time.Now(),
		metadata,
	)

	if msg == nil {
		t.Fatal("NewUserMessage returned nil")
	}

	if msg.MessageID != "msg-123" {
		t.Errorf("MessageID = %v, want msg-123", msg.MessageID)
	}

	if msg.ExternalMessageID != "wamid.abc" {
		t.Errorf("ExternalMessageID = %v, want wamid.abc", msg.ExternalMessageID)
	}

	if msg.UserID != "user-456" {
		t.Errorf("UserID = %v, want user-456", msg.UserID)
	}

	if msg.Channel != ChannelWhatsApp {
		t.Errorf("Channel = %v, want whatsapp", msg.Channel)
	}

	if msg.Text != "Hello world" {
		t.Errorf("Text = %v, want Hello world", msg.Text)
	}
}

func TestNewResponseMessage(t *testing.T) {
	resp := NewResponseMessage(
		"+1234567890",
		ResponseTypeAudio,
		"Test response",
		"123456789",
	)

	if resp == nil {
		t.Fatal("NewResponseMessage returned nil")
	}

	if resp.ResponseID == "" {
		t.Error("ResponseID should not be empty")
	}

	if resp.TargetUser != "+1234567890" {
		t.Errorf("TargetUser = %v, want +1234567890", resp.TargetUser)
	}

	if resp.ResponseType != ResponseTypeAudio {
		t.Errorf("ResponseType = %v, want audio", resp.ResponseType)
	}

	if resp.Text != "Test response" {
		t.Errorf("Text = %v, want Test response", resp.Text)
	}

	if resp.Metadata.PhoneNumberID != "123456789" {
		t.Errorf("PhoneNumberID = %v, want 123456789", resp.Metadata.PhoneNumberID)
	}
}

func TestNewAudioAsset(t *testing.T) {
	audio := NewAudioAsset(
		FormatAAC,
		CodecOpus,
		[]byte{0x00, 0x01, 0x02},
	)

	if audio == nil {
		t.Fatal("NewAudioAsset returned nil")
	}

	if audio.AudioID == "" {
		t.Error("AudioID should not be empty")
	}

	if audio.Format != FormatAAC {
		t.Errorf("Format = %v, want aac", audio.Format)
	}

	if audio.Codec != CodecOpus {
		t.Errorf("Codec = %v, want opus", audio.Codec)
	}

	if len(audio.BinaryData) != 3 {
		t.Errorf("BinaryData length = %v, want 3", len(audio.BinaryData))
	}
}

func TestChannelValues(t *testing.T) {
	tests := []struct {
		name     string
		got      Channel
		expected Channel
	}{
		{"whatsapp", ChannelWhatsApp, "whatsapp"},
		{"telegram", ChannelTelegram, "telegram"},
		{"webchat", ChannelWebChat, "webchat"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("got %v, want %v", tt.got, tt.expected)
			}
		})
	}
}

func TestAudioFormatValues(t *testing.T) {
	tests := []struct {
		name     string
		got      AudioFormat
		expected AudioFormat
	}{
		{"aac", FormatAAC, "aac"},
		{"mp3", FormatMP3, "mp3"},
		{"wav", FormatWAV, "wav"},
		{"ogg", FormatOGG, "ogg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.expected {
				t.Errorf("got %v, want %v", tt.got, tt.expected)
			}
		})
	}
}
