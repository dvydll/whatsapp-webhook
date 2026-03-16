package whatsapp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendTextMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			t.Error("Authorization header missing")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"messaging_product": "whatsapp",
			"contacts": [{"input": "34685107027", "wa_id": "34685107027"}],
			"messages": [{"id": "wamid.test123"}]
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{
		PhoneNumberID: "123456",
		AccessToken:   "test-token",
		BaseURL:       server.URL,
		Version:       "v25.0",
	})

	resp, err := client.SendTextMessage(context.Background(), "34685107027", "Hello")

	if err != nil {
		t.Errorf("SendTextMessage() error = %v", err)
		return
	}

	if resp == nil {
		t.Error("SendTextMessage() returned nil response")
		return
	}

	if len(resp.Messages) == 0 {
		t.Error("SendTextMessage() returned no messages")
	}

	if resp.Messages[0].ID == "" {
		t.Error("SendTextMessage() returned empty message ID")
	}
}

func TestSendAudioMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"messaging_product": "whatsapp",
			"messages": [{"id": "wamid.test456"}]
		}`))
	}))
	defer server.Close()

	client := NewClient(Config{
		PhoneNumberID: "123456",
		AccessToken:   "test-token",
		BaseURL:       server.URL,
		Version:       "v25.0",
	})

	resp, err := client.SendAudioMessage(context.Background(), "34685107027", "https://example.com/audio.ogg")

	if err != nil {
		t.Errorf("SendAudioMessage() error = %v", err)
		return
	}

	if resp == nil {
		t.Error("SendAudioMessage() returned nil response")
	}

	if len(resp.Messages) == 0 {
		t.Error("SendAudioMessage() returned no messages")
	}
}

func TestSendMessageAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": {"message": "Invalid token", "type": "OAuthException"}}`))
	}))
	defer server.Close()

	client := NewClient(Config{
		PhoneNumberID: "123456",
		AccessToken:   "invalid-token",
		BaseURL:       server.URL,
		Version:       "v25.0",
	})

	_, err := client.SendTextMessage(context.Background(), "34685107027", "Test")

	if err == nil {
		t.Error("SendTextMessage() expected error, got nil")
	}
}
