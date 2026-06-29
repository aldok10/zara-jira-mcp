package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	cmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/confluence/mcp"
)

func RegisterConfluenceTools(s *server.MCPServer, h *cmcp.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_confluence_search",
			mcp.WithDescription("Search Confluence pages by CQL query. Find documentation, specs, decisions."),
			mcp.WithString("query", mcp.Required(), mcp.Description("CQL query (e.g. 'type=page AND space=DEV AND text~sprint')")),
			mcp.WithNumber("limit", mcp.Description("Max results (default 10)")),
		),
		h.SearchPages,
	)
	s.AddTool(
		mcp.NewTool("pm_confluence_get_page",
			mcp.WithDescription("Get a Confluence page content by ID."),
			mcp.WithString("page_id", mcp.Required(), mcp.Description("Page ID")),
		),
		h.GetPage,
	)
	s.AddTool(
		mcp.NewTool("pm_confluence_create_page",
			mcp.WithDescription("Create a new Confluence page. Use for sprint reports, decision records, meeting notes."),
			mcp.WithString("space_key", mcp.Required(), mcp.Description("Space key (e.g. DEV, TEAM)")),
			mcp.WithString("title", mcp.Required(), mcp.Description("Page title")),
			mcp.WithString("body", mcp.Description("Page content in XHTML storage format")),
			mcp.WithString("parent_id", mcp.Description("Parent page ID (creates as child page)")),
		),
		h.CreatePage,
	)
}
