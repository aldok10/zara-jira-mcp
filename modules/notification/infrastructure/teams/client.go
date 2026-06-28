package teams

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
)

// Client sends messages to Microsoft Teams via incoming webhook.
type Client struct {
	webhookURL string
	httpClient *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		webhookURL: cfg.Teams.WebhookURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
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
	data, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return sanitizeReqErr(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("teams webhook returned %d", resp.StatusCode)
	}
	return nil
}

// sanitizeReqErr strips the URL from HTTP client errors to prevent
// credential leakage (webhook URLs contain tokens in the path).
func sanitizeReqErr(err error) error {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return urlErr.Err
	}
	return err
}
