package lark

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"

	domain "github.com/aldok10/zara-jira-mcp/shared/domain/agent"
)

// Helper functions for Lark message construction

func buildTextContent(text string) string {
	return larkim.NewTextMsgBuilder().Text(text).Build()
}

func buildCreateMsgReq(chatID, content string) *larkim.CreateMessageReq {
	return larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(chatID).
			Content(content).
			Build()).
		Build()
}

func buildCardContent(title, markdown string) (string, error) {
	header := larkcard.NewMessageCardHeader().
		Template(larkcard.TemplateBlue).
		Title(larkcard.NewMessageCardPlainText().Content(title).Build()).
		Build()

	mdElement := larkcard.NewMessageCardMarkdown().Content(markdown).Build()

	card, err := larkcard.NewMessageCard().
		Header(header).
		Elements([]larkcard.MessageCardElement{mdElement}).
		String()
	if err != nil {
		return "", fmt.Errorf("build card: %w", err)
	}
	return card, nil
}

func buildCreateMsgCardReq(chatID, cardContent string) *larkim.CreateMessageReq {
	return larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeInteractive).
			ReceiveId(chatID).
			Content(cardContent).
			Build()).
		Build()
}

// Card templates for different response types
const (
	templateBlue   = "blue"
	templateRed    = "red"
	templateOrange = "orange"
	templateGreen  = "green"
	templatePurple = "purple"
)

// Ensure GatewayAdapter satisfies domain.Gateway at compile time.
var _ domain.Gateway = (*GatewayAdapter)(nil)

// GatewayAdapter for Lark.
// Features smart card detection: structured responses get rich cards,
// simple text gets plain messages.
type GatewayAdapter struct {
	client *WebhookClient
	logger *slog.Logger
}

// NewGatewayAdapter creates a Lark gateway adapter.
func NewGatewayAdapter(client *WebhookClient, logger *slog.Logger) *GatewayAdapter {
	return &GatewayAdapter{client: client, logger: logger}
}

// SendText sends a text or card response depending on content structure.
// Smart detection:
//   - Sprint status → blue card with progress
//   - Blockers → red card
//   - Risks → orange card
//   - Decisions → green card
//   - Short text → plain text
func (g *GatewayAdapter) SendText(ctx context.Context, channelID, text string) error {
	if text == "" {
		return nil
	}

	// Try to detect response type and send as card
	cardType, title := detectResponseType(text)
	if cardType != "" && g.client.sdk != nil && channelID != "" {
		md := formatAsMarkdown(text)
		if md != text {
			// Has structure — send as card
			g.logger.Debug("sending as card", "type", cardType, "title", title)
			return g.client.sdkSendCardToChat(ctx, channelID, title, md)
		}
	}

	// Fallback: plain text
	if g.client.sdk != nil && channelID != "" {
		return g.client.sdkSendTextToChat(ctx, channelID, text)
	}
	return g.client.SendText(ctx, text)
}

func (g *GatewayAdapter) SendMarkdown(ctx context.Context, channelID, title, content string) error {
	if g.client.sdk != nil && channelID != "" {
		return g.client.sdkSendCardToChat(ctx, channelID, title, content)
	}
	return g.client.SendMarkdown(ctx, title, content)
}

func (g *GatewayAdapter) Channel() string { return "lark" }

// detectResponseType analyzes text and returns (cardTemplate, title).
// Empty template = plain text.
func detectResponseType(text string) (template, title string) {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) == 0 {
		return "", ""
	}

	firstLine := lines[0]

	// Sprint status
	if strings.Contains(firstLine, "Sprint:") || strings.Contains(firstLine, "Sprint ") {
		return templateBlue, "📋 Sprint Status"
	}

	// Blockers
	if strings.Contains(firstLine, "Blocker") || strings.Contains(firstLine, "blocker") {
		return templateRed, "🚫 Active Blockers"
	}

	// Risks
	if strings.Contains(firstLine, "Risk") || strings.Contains(firstLine, "risk") {
		return templateOrange, "⚠️ Open Risks"
	}

	// Action items
	if strings.Contains(firstLine, "Action") || strings.Contains(firstLine, "Pending") {
		return templatePurple, "📌 Action Items"
	}

	// Decisions
	if strings.Contains(firstLine, "Decision") || strings.Contains(firstLine, "decision") {
		return templateGreen, "✅ Decisions"
	}

	// Team overview
	if strings.Contains(firstLine, "Team Overview") || strings.Contains(firstLine, "Workload") {
		return templateBlue, "👥 Team Overview"
	}

	// Long structured messages (3+ lines with bullets or pipes)
	if len(lines) >= 3 {
		bulletCount := 0
		pipeCount := 0
		for _, l := range lines {
			trimmed := strings.TrimSpace(l)
			if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "├") || strings.HasPrefix(trimmed, "└") {
				bulletCount++
			}
			if strings.Contains(trimmed, "│") || strings.Contains(trimmed, "|") {
				pipeCount++
			}
		}
		if bulletCount >= 2 || pipeCount >= 2 {
			return templateBlue, "📋 Summary"
		}
	}

	return "", ""
}

// formatAsMarkdown converts structured text to Lark markdown.
// Preserves emoji, converts tree chars to plain markdown.
func formatAsMarkdown(text string) string {
	if text == "" {
		return text
	}

	var sb strings.Builder
	for _, line := range strings.Split(text, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "├") || strings.HasPrefix(trimmed, "└"):
			// Tree formatting → indent with dash
			content := strings.TrimLeft(trimmed, "├└─ ")
			sb.WriteString("  - " + content + "\n")
		case strings.HasPrefix(trimmed, "│"):
			// Continuation pipe → ignore
			continue
		default:
			sb.WriteString(line + "\n")
		}
	}

	result := strings.TrimSpace(sb.String())
	// If the result is the same as input, no transformation was needed
	if result == strings.TrimSpace(text) {
		return text
	}
	return result
}

// sdkSendTextToChat sends text to a specific chat ID.
func (c *WebhookClient) sdkSendTextToChat(ctx context.Context, chatID, text string) error {
	if c.sdk == nil {
		return fmt.Errorf("lark SDK not configured")
	}

	content := buildTextContent(text)
	resp, err := c.sdk.Im.Message.Create(ctx, buildCreateMsgReq(chatID, content))
	if err != nil {
		return fmt.Errorf("lark send text: %w", err)
	}
	if !resp.Success() {
		return fmt.Errorf("lark send error %d: %s", resp.Code, resp.Msg)
	}
	return nil
}

// sdkSendCardToChat sends an interactive card to a specific chat ID.
func (c *WebhookClient) sdkSendCardToChat(ctx context.Context, chatID, title, markdown string) error {
	if c.sdk == nil {
		return fmt.Errorf("lark SDK not configured")
	}

	cardContent, err := buildCardContent(title, markdown)
	if err != nil {
		return err
	}

	resp, err := c.sdk.Im.Message.Create(ctx, buildCreateMsgCardReq(chatID, cardContent))
	if err != nil {
		return fmt.Errorf("lark send card: %w", err)
	}
	if !resp.Success() {
		return fmt.Errorf("lark send error %d: %s", resp.Code, resp.Msg)
	}
	return nil
}
