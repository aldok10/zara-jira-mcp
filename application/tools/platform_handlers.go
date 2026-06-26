package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// Discord handlers

func (h *Handlers) DiscordSend(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Discord == nil || !h.Discord.Available() {
		return errorResult("Discord not configured. Set DISCORD_BOT_TOKEN."), nil
	}
	content, err := req.RequireString("content")
	if err != nil {
		return errorResult("content parameter required"), nil
	}
	channel := req.GetString("channel", "")
	title := req.GetString("title", "")

	if title != "" {
		if err := h.Discord.SendEmbed(ctx, channel, title, content, 0x3498DB); err != nil {
			return errorResult("Discord send failed: " + err.Error()), nil
		}
	} else {
		if err := h.Discord.SendMessage(ctx, channel, content); err != nil {
			return errorResult("Discord send failed: " + err.Error()), nil
		}
	}
	return textResult("Message sent to Discord."), nil
}

// Telegram handlers

func (h *Handlers) TelegramSend(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Telegram == nil || !h.Telegram.Available() {
		return errorResult("Telegram not configured. Set TELEGRAM_BOT_TOKEN."), nil
	}
	text, err := req.RequireString("text")
	if err != nil {
		return errorResult("text parameter required"), nil
	}
	chatIDStr := req.GetString("chat_id", "")
	var chatID int64
	if chatIDStr != "" {
		chatID, _ = strconv.ParseInt(chatIDStr, 10, 64)
	}

	if err := h.Telegram.SendMessage(ctx, chatID, text); err != nil {
		return errorResult("Telegram send failed: " + err.Error()), nil
	}
	return textResult("Message sent to Telegram."), nil
}

// Teams handlers

func (h *Handlers) TeamsSend(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Teams == nil || !h.Teams.Available() {
		return errorResult("Teams not configured. Set TEAMS_WEBHOOK_URL."), nil
	}
	content, err := req.RequireString("content")
	if err != nil {
		return errorResult("content parameter required"), nil
	}
	title := req.GetString("title", "")

	if title != "" {
		if err := h.Teams.SendCard(ctx, title, content); err != nil {
			return errorResult("Teams send failed: " + err.Error()), nil
		}
	} else {
		if err := h.Teams.SendMessage(ctx, content); err != nil {
			return errorResult("Teams send failed: " + err.Error()), nil
		}
	}
	return textResult("Message sent to Teams."), nil
}

// Email handlers

func (h *Handlers) EmailSend(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Email == nil || !h.Email.Available() {
		return errorResult("Email not configured. Set EMAIL_SMTP_HOST and EMAIL_FROM."), nil
	}
	to, err := req.RequireString("to")
	if err != nil {
		return errorResult("to parameter required"), nil
	}
	subject, err := req.RequireString("subject")
	if err != nil {
		return errorResult("subject parameter required"), nil
	}
	body, err := req.RequireString("body")
	if err != nil {
		return errorResult("body parameter required"), nil
	}

	if err := h.Email.Send(ctx, to, subject, body); err != nil {
		return errorResult("Email send failed: " + err.Error()), nil
	}
	return textResult(fmt.Sprintf("Email sent to %s.", to)), nil
}

// Confluence handlers

func (h *Handlers) ConfluenceSearch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Confluence == nil || !h.Confluence.Available() {
		return errorResult("Confluence not configured. Set CONFLUENCE_BASE_URL and CONFLUENCE_API_TOKEN."), nil
	}
	query, err := req.RequireString("query")
	if err != nil {
		return errorResult("query parameter required"), nil
	}
	limit := req.GetInt("limit", 10)

	pages, err := h.Confluence.SearchPages(ctx, query, limit)
	if err != nil {
		return errorResult("Confluence search failed: " + err.Error()), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d pages:\n\n", len(pages)))
	for _, p := range pages {
		sb.WriteString(fmt.Sprintf("- [%s] %s (%s) %s\n", p.SpaceKey, p.Title, p.Type, p.WebURL))
	}
	return textResult(sb.String()), nil
}

func (h *Handlers) ConfluenceGetPage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Confluence == nil || !h.Confluence.Available() {
		return errorResult("Confluence not configured."), nil
	}
	pageID, err := req.RequireString("page_id")
	if err != nil {
		return errorResult("page_id parameter required"), nil
	}

	page, err := h.Confluence.GetPage(ctx, pageID)
	if err != nil {
		return errorResult("Failed to get page: " + err.Error()), nil
	}

	body := page.Body
	if len(body) > 3000 {
		body = body[:3000] + "\n\n... (truncated)"
	}
	return textResult(fmt.Sprintf("Title: %s\nSpace: %s\nVersion: %d\n\n%s", page.Title, page.SpaceKey, page.Version, body)), nil
}

func (h *Handlers) ConfluenceCreatePage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Confluence == nil || !h.Confluence.Available() {
		return errorResult("Confluence not configured."), nil
	}
	spaceKey, err := req.RequireString("space_key")
	if err != nil {
		return errorResult("space_key parameter required"), nil
	}
	title, err := req.RequireString("title")
	if err != nil {
		return errorResult("title parameter required"), nil
	}
	body, err := req.RequireString("body")
	if err != nil {
		return errorResult("body parameter required (XHTML storage format)"), nil
	}
	parentID := req.GetString("parent_id", "")

	page, err := h.Confluence.CreatePage(ctx, spaceKey, title, body, parentID)
	if err != nil {
		return errorResult("Failed to create page: " + err.Error()), nil
	}
	return textResult(fmt.Sprintf("Created: %s — %s", page.Title, page.WebURL)), nil
}

// Broadcast sends to ALL configured channels at once.
func (h *Handlers) Broadcast(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	content, err := req.RequireString("content")
	if err != nil {
		return errorResult("content parameter required"), nil
	}
	title := req.GetString("title", "PM Update")

	var results []string

	if h.Lark != nil {
		if err := h.Lark.SendMarkdown(ctx, title, content); err == nil {
			results = append(results, "Lark: sent")
		} else {
			results = append(results, "Lark: "+err.Error())
		}
	}
	if h.Slack != nil && h.Slack.Available() {
		if err := h.Slack.SendRichMessage(ctx, "", title, content); err == nil {
			results = append(results, "Slack: sent")
		} else {
			results = append(results, "Slack: "+err.Error())
		}
	}
	if h.Discord != nil && h.Discord.Available() {
		if err := h.Discord.SendEmbed(ctx, "", title, content, 0x3498DB); err == nil {
			results = append(results, "Discord: sent")
		} else {
			results = append(results, "Discord: "+err.Error())
		}
	}
	if h.Telegram != nil && h.Telegram.Available() {
		msg := fmt.Sprintf("*%s*\n\n%s", title, content)
		if err := h.Telegram.SendMessage(ctx, 0, msg); err == nil {
			results = append(results, "Telegram: sent")
		} else {
			results = append(results, "Telegram: "+err.Error())
		}
	}
	if h.Teams != nil && h.Teams.Available() {
		if err := h.Teams.SendCard(ctx, title, content); err == nil {
			results = append(results, "Teams: sent")
		} else {
			results = append(results, "Teams: "+err.Error())
		}
	}

	if len(results) == 0 {
		return errorResult("No notification channels configured."), nil
	}
	return textResult("Broadcast results:\n" + strings.Join(results, "\n")), nil
}
