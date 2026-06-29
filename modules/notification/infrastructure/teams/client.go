package teams

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/httpclient"
)

// Client sends messages to Microsoft Teams via incoming webhook.
type Client struct {
	webhookURL string
	httpClient *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		webhookURL: cfg.Teams.WebhookURL,
		httpClient: httpclient.NewWithTimeout(10 * time.Second),
	}
}

func (c *Client) Available() bool {
	return c.webhookURL != ""
}

// SendMessage sends a simple text message.
func (c *Client) SendMessage(ctx context.Context, text string) error {
	payload := map[string]string{"text": text}
	return c.post(ctx, payload)
}

// SendCard sends an Adaptive Card with title and body.
func (c *Client) SendCard(ctx context.Context, title, body string) error {
	card := map[string]any{
		"type": "message",
		"attachments": []map[string]any{
			{
				"contentType": "application/vnd.microsoft.card.adaptive",
				"content": map[string]any{
					"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
					"type":    "AdaptiveCard",
					"version": "1.4",
					"body": []map[string]any{
						{"type": "TextBlock", "text": title, "weight": "Bolder", "size": "Medium"},
						{"type": "TextBlock", "text": body, "wrap": true},
					},
				},
			},
		},
	}
	return c.post(ctx, card)
}

func (c *Client) post(ctx context.Context, payload any) error {
	if c.webhookURL == "" {
		return fmt.Errorf("teams webhook not configured: set TEAMS_WEBHOOK_URL")
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return httpclient.SanitizeError(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		// Don't leak full error body to prevent information leakage
		_, _ = io.Copy(io.Discard, resp.Body)
		return fmt.Errorf("teams webhook returned %d", resp.StatusCode)
	}
	return nil
}
