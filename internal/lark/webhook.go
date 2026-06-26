package lark

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"

	"github.com/aldok10/zara-jira-mcp/config"
)

// WebhookClient sends messages to Lark via SDK or webhook fallback.
type WebhookClient struct {
	webhookURL string
	chatID     string
	sdk        *lark.Client
	httpClient *http.Client
}

func NewWebhookClient(cfg *config.Config) *WebhookClient {
	c := &WebhookClient{
		webhookURL: cfg.Lark.WebhookURL,
		chatID:     cfg.Lark.ChatID,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}

	if cfg.Lark.AppID != "" && cfg.Lark.AppSecret != "" {
		c.sdk = lark.NewClient(cfg.Lark.AppID, cfg.Lark.AppSecret)
	}

	return c
}

func (c *WebhookClient) SendText(ctx context.Context, text string) error {
	if c.sdk != nil && c.chatID != "" {
		return c.sdkSendText(ctx, text)
	}
	return c.webhookPost(ctx, map[string]any{
		"msg_type": "text",
		"content":  map[string]string{"text": text},
	})
}

func (c *WebhookClient) SendMarkdown(ctx context.Context, title, content string) error {
	if c.sdk != nil && c.chatID != "" {
		return c.sdkSendCard(ctx, title, content)
	}
	return c.webhookPost(ctx, map[string]any{
		"msg_type": "interactive",
		"card": map[string]any{
			"header": map[string]any{
				"title":    map[string]string{"tag": "plain_text", "content": title},
				"template": "blue",
			},
			"elements": []map[string]any{
				{"tag": "markdown", "content": content},
			},
		},
	})
}

// SDK methods

func (c *WebhookClient) sdkSendText(ctx context.Context, text string) error {
	content := larkim.NewTextMsgBuilder().Text(text).Build()
	resp, err := c.sdk.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(c.chatID).
			Content(content).
			Build()).
		Build())
	if err != nil {
		return err
	}
	if !resp.Success() {
		return fmt.Errorf("lark SDK error %d: %s", resp.Code, resp.Msg)
	}
	return nil
}

func (c *WebhookClient) sdkSendCard(ctx context.Context, title, markdown string) error {
	header := larkcard.NewMessageCardHeader().
		Template(larkcard.TemplateBlue).
		Title(larkcard.NewMessageCardPlainText().Content(title).Build()).
		Build()

	mdElement := larkcard.NewMessageCardMarkdown().Content(markdown).Build()

	cardContent, err := larkcard.NewMessageCard().
		Header(header).
		Elements([]larkcard.MessageCardElement{mdElement}).
		String()
	if err != nil {
		return fmt.Errorf("build card: %w", err)
	}

	resp, err := c.sdk.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(c.chatID).
			Content(cardContent).
			Build()).
		Build())
	if err != nil {
		return err
	}
	if !resp.Success() {
		return fmt.Errorf("lark SDK error %d: %s", resp.Code, resp.Msg)
	}
	return nil
}

// Webhook fallback

func (c *WebhookClient) webhookPost(ctx context.Context, payload any) error {
	if c.webhookURL == "" {
		return fmt.Errorf("lark not configured: set JIRA_LARK_WEBHOOK_URL or LARK_APP_ID+LARK_APP_SECRET+LARK_CHAT_ID")
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
