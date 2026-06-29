package slack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	slackapi "github.com/slack-go/slack"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/httpclient"
)

// Client wraps slack-go for sending messages and managing channels.
type Client struct {
	api            *slackapi.Client
	defaultChannel string
	webhookURL     string
	httpClient     *http.Client
}

func NewClient(cfg *config.Config) *Client {
	c := &Client{
		defaultChannel: cfg.Slack.DefaultChannel,
		webhookURL:     cfg.Slack.WebhookURL,
		httpClient:     httpclient.NewWithTimeout(10 * time.Second),
	}
	if cfg.Slack.BotToken != "" {
		c.api = slackapi.New(cfg.Slack.BotToken)
	}
	return c
}

func (c *Client) Available() bool {
	return c.api != nil || c.webhookURL != ""
}

// SendMessage sends a text message to a channel.
func (c *Client) SendMessage(ctx context.Context, channel, text string) error {
	if channel == "" {
		channel = c.defaultChannel
	}
	if c.api != nil {
		_, _, err := c.api.PostMessageContext(ctx, channel, slackapi.MsgOptionText(text, false))
		return err
	}
	return c.webhookPost(ctx, text, "")
}

// SendRichMessage sends a message with blocks (markdown sections).
func (c *Client) SendRichMessage(ctx context.Context, channel, title, markdown string) error {
	if channel == "" {
		channel = c.defaultChannel
	}
	if c.api != nil {
		headerBlock := slackapi.NewHeaderBlock(slackapi.NewTextBlockObject("plain_text", title, false, false))
		sectionBlock := slackapi.NewSectionBlock(slackapi.NewTextBlockObject("mrkdwn", markdown, false, false), nil, nil)
		_, _, err := c.api.PostMessageContext(ctx, channel,
			slackapi.MsgOptionBlocks(headerBlock, sectionBlock),
		)
		return err
	}
	return c.webhookPost(ctx, fmt.Sprintf("*%s*\n\n%s", title, markdown), "")
}

// ListChannels returns accessible channels.
func (c *Client) ListChannels(ctx context.Context) ([]Channel, error) {
	if c.api == nil {
		return nil, fmt.Errorf("slack bot token not configured")
	}
	params := &slackapi.GetConversationsParameters{Limit: 100, Types: []string{"public_channel", "private_channel"}}
	channels, _, err := c.api.GetConversationsContext(ctx, params)
	if err != nil {
		return nil, err
	}
	out := make([]Channel, len(channels))
	for i := range channels {
		out[i] = Channel{ID: channels[i].ID, Name: channels[i].Name, Topic: channels[i].Topic.Value, MemberCount: channels[i].NumMembers}
	}
	return out, nil
}

// GetChannelHistory returns recent messages from a channel.
func (c *Client) GetChannelHistory(ctx context.Context, channel string, limit int) ([]Message, error) {
	if c.api == nil {
		return nil, fmt.Errorf("slack bot token not configured")
	}
	if limit <= 0 {
		limit = 20
	}
	params := &slackapi.GetConversationHistoryParameters{ChannelID: channel, Limit: limit}
	resp, err := c.api.GetConversationHistoryContext(ctx, params)
	if err != nil {
		return nil, err
	}
	out := make([]Message, 0, len(resp.Messages))
	for i := range resp.Messages {
		out = append(out, Message{User: resp.Messages[i].User, Text: resp.Messages[i].Text, Timestamp: resp.Messages[i].Timestamp})
	}
	return out, nil
}

// Channel is a simplified Slack channel.
type Channel struct {
	ID          string
	Name        string
	Topic       string
	MemberCount int
}

// Message is a simplified Slack message.
type Message struct {
	User      string
	Text      string
	Timestamp string
}

// webhookPost sends via incoming webhook (fallback).
func (c *Client) webhookPost(ctx context.Context, text, channel string) error {
	if c.webhookURL == "" {
		return fmt.Errorf("slack not configured: set SLACK_BOT_TOKEN or SLACK_WEBHOOK_URL")
	}
	payload := map[string]string{"text": text}
	if channel != "" {
		payload["channel"] = channel
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.webhookURL, bytes.NewReader(body))
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
		// Drain body to prevent connection leak
		_, _ = io.Copy(io.Discard, resp.Body)
		return fmt.Errorf("slack webhook returned %d", resp.StatusCode)
	}
	return nil
}
