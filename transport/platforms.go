package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerPlatformTools(s *server.MCPServer, h *tools.Handlers) {
	// Discord
	s.AddTool(
		mcp.NewTool("discord_send",
			mcp.WithDescription("Send a message to a Discord channel. Supports plain text or embed with title."),
			mcp.WithString("content", mcp.Required(), mcp.Description("Message content")),
			mcp.WithString("channel", mcp.Description("Channel ID (uses default if empty)")),
			mcp.WithString("title", mcp.Description("Embed title (sends as rich embed when set)")),
		),
		h.DiscordSend,
	)

	// Telegram
	s.AddTool(
		mcp.NewTool("telegram_send",
			mcp.WithDescription("Send a message to a Telegram chat/group. Supports Markdown formatting."),
			mcp.WithString("text", mcp.Required(), mcp.Description("Message text (Markdown supported)")),
			mcp.WithString("chat_id", mcp.Description("Chat ID (uses default if empty)")),
		),
		h.TelegramSend,
	)

	// Microsoft Teams
	s.AddTool(
		mcp.NewTool("teams_send",
			mcp.WithDescription("Send a message to Microsoft Teams via incoming webhook. Supports Adaptive Cards."),
			mcp.WithString("content", mcp.Required(), mcp.Description("Message content")),
			mcp.WithString("title", mcp.Description("Card title (sends as Adaptive Card when set)")),
		),
		h.TeamsSend,
	)

	// Email
	s.AddTool(
		mcp.NewTool("email_send",
			mcp.WithDescription("Send an email via SMTP. Use for formal notifications, stakeholder updates, or escalations."),
			mcp.WithString("to", mcp.Required(), mcp.Description("Recipient email address")),
			mcp.WithString("subject", mcp.Required(), mcp.Description("Email subject")),
			mcp.WithString("body", mcp.Required(), mcp.Description("Email body (plain text)")),
		),
		h.EmailSend,
	)

	// Confluence
	s.AddTool(
		mcp.NewTool("confluence_search",
			mcp.WithDescription("Search Confluence pages by CQL query. Find documentation, specs, decisions."),
			mcp.WithString("query", mcp.Required(), mcp.Description("CQL query (e.g. 'type=page AND space=DEV AND text~sprint')")),
			mcp.WithNumber("limit", mcp.Description("Max results (default 10)")),
		),
		h.ConfluenceSearch,
	)

	s.AddTool(
		mcp.NewTool("confluence_get_page",
			mcp.WithDescription("Get a Confluence page content by ID."),
			mcp.WithString("page_id", mcp.Required(), mcp.Description("Page ID")),
		),
		h.ConfluenceGetPage,
	)

	s.AddTool(
		mcp.NewTool("confluence_create_page",
			mcp.WithDescription("Create a new Confluence page. Use for sprint reports, decision records, meeting notes."),
			mcp.WithString("space_key", mcp.Required(), mcp.Description("Space key (e.g. DEV, TEAM)")),
			mcp.WithString("title", mcp.Required(), mcp.Description("Page title")),
			mcp.WithString("body", mcp.Required(), mcp.Description("Page content in XHTML storage format")),
			mcp.WithString("parent_id", mcp.Description("Parent page ID (creates as child page)")),
		),
		h.ConfluenceCreatePage,
	)

	// Broadcast (all channels at once)
	s.AddTool(
		mcp.NewTool("broadcast",
			mcp.WithDescription("Send a message to ALL configured notification channels at once (Lark, Slack, Discord, Telegram, Teams). Use for critical announcements."),
			mcp.WithString("content", mcp.Required(), mcp.Description("Message content")),
			mcp.WithString("title", mcp.Description("Message title (default: PM Update)")),
		),
		h.Broadcast,
	)
}
