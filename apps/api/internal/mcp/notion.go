package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	nmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/notion/mcp"
)

func RegisterNotionTools(s *server.MCPServer, h *nmcp.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_notion_search",
			mcp.WithDescription("Search Notion pages and databases by keyword."),
			mcp.WithString("query", mcp.Required(), mcp.Description("Search keyword")),
			mcp.WithNumber("limit", mcp.Description("Max results (default: 10)")),
		),
		h.SearchPages,
	)
	s.AddTool(
		mcp.NewTool("pm_notion_create_page",
			mcp.WithDescription("Create a Notion page (for meeting notes, decisions, etc)."),
			mcp.WithString("title", mcp.Required(), mcp.Description("Page title")),
			mcp.WithString("content", mcp.Description("Page body content (plain text)")),
			mcp.WithString("parent_id", mcp.Description("Parent database ID (uses NOTION_DATABASE_ID if empty)")),
		),
		h.CreatePage,
	)
	s.AddTool(
		mcp.NewTool("pm_notion_query_db",
			mcp.WithDescription("Query a Notion database for tracking items."),
			mcp.WithString("database_id", mcp.Description("Database ID (uses NOTION_DATABASE_ID if empty)")),
			mcp.WithString("filter", mcp.Description("Notion filter as JSON (optional)")),
			mcp.WithNumber("limit", mcp.Description("Max results (default: 20)")),
		),
		h.QueryDatabase,
	)
}
