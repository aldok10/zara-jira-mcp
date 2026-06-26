package lark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aldok10/zara-jira-mcp/config"
)

// WebhookClient sends messages to Lark via custom bot webhook.
type WebhookClient struct {
	webhookURL string
	httpClient *http.Client
}

func NewWebhookClient(cfg *config.Config) *WebhookClient {
	return &WebhookClient{
		webhookURL: cfg.Lark.WebhookURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *WebhookClient) SendText(ctx context.Context, text string) error {
	payload := map[string]any{
		"msg_type": "text",
		"content":  map[string]string{"text": text},
	}
	return c.post(ctx, payload)
}

func (c *WebhookClient) SendMarkdown(ctx context.Context, title, content string) error {
	// Lark custom bot uses interactive card for rich content
	payload := map[string]any{
		"msg_type": "interactive",
		"card": map[string]any{
			"header": map[string]any{
				"title": map[string]string{
					"tag":     "plain_text",
					"content": title,
				},
				"template": "blue",
			},
			"elements": []map[string]any{
				{
					"tag":     "markdown",
					"content": content,
				},
			},
		},
	}
	return c.post(ctx, payload)
}

func (c *WebhookClient) post(ctx context.Context, payload any) error {
	if c.webhookURL == "" {
		return fmt.Errorf("lark webhook URL not configured")
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("lark webhook returned %d", resp.StatusCode)
	}

	return nil
}
