package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	httpClient  *http.Client
	baseURL     string
	version     string
	phoneID     string
	accessToken string
}

type Config struct {
	PhoneNumberID string
	AccessToken   string
	BaseURL       string
	Version       string
}

func NewClient(config Config) *Client {
	return &Client{
		httpClient:  &http.Client{Timeout: 30 * time.Second},
		baseURL:     config.BaseURL,
		version:     config.Version,
		phoneID:     config.PhoneNumberID,
		accessToken: config.AccessToken,
	}
}

type SendAudioRequest struct {
	MessagingProduct string       `json:"messaging_product"`
	To               string       `json:"to"`
	Type             string       `json:"type"`
	Audio            AudioContent `json:"audio"`
}

type AudioContent struct {
	Link string `json:"link,omitempty"`
	ID   string `json:"id,omitempty"`
}

type SendTextRequest struct {
	MessagingProduct string      `json:"messaging_product"`
	To               string      `json:"to"`
	Type             string      `json:"type"`
	Text             TextContent `json:"text"`
}

type TextContent struct {
	Body string `json:"body"`
}

type SendMessageResponse struct {
	MessagingProduct string    `json:"messaging_product"`
	Contacts         []Contact `json:"contacts"`
	Messages         []Message `json:"messages"`
}

type Contact struct {
	Input string `json:"input"`
	WaID  string `json:"wa_id"`
}

type Message struct {
	ID string `json:"id"`
}

func (c *Client) SendAudioMessage(ctx context.Context, to, audioURL string) (*SendMessageResponse, error) {
	req := SendAudioRequest{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "audio",
		Audio:            AudioContent{Link: audioURL},
	}

	return c.sendMessage(ctx, req)
}

func (c *Client) SendTextMessage(ctx context.Context, to, text string) (*SendMessageResponse, error) {
	req := SendTextRequest{
		MessagingProduct: "whatsapp",
		To:               to,
		Type:             "text",
		Text:             TextContent{Body: text},
	}

	return c.sendMessage(ctx, req)
}

func (c *Client) sendMessage(ctx context.Context, payload interface{}) (*SendMessageResponse, error) {
	url := fmt.Sprintf("%s/%s/%s/messages", c.baseURL, c.version, c.phoneID)

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("WhatsApp API error (status=%d): %s", resp.StatusCode, string(respBody))
	}

	var result SendMessageResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
