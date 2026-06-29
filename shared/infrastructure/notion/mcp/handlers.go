package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/notion"
)

type Handlers struct {
	client *notion.Client
	errH   *mcputil.ErrorHandler
}

func NewHandlers(client *notion.Client) *Handlers {
	return &Handlers{
		client: client,
		errH:   mcputil.NewErrorHandler(nil),
	}
}

func (h *Handlers) SearchPages(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Notion not configured. Set NOTION_API_KEY."), nil
	}
	query, err := req.RequireString("query")
	if err != nil {
		return mcputil.ErrInvalid("query parameter is required"), nil
	}
	limit := req.GetInt("limit", 10)

	results, err := h.client.Search(ctx, query, limit)
	if err != nil {
		return h.errH.Wrap("search", err), nil
	}
	if len(results) == 0 {
		return mcputil.TextResult(fmt.Sprintf("No Notion results for %q", query)), nil
	}
	var b strings.Builder
	for _, r := range results {
		b.WriteString(fmt.Sprintf("[%s] %s — %s\n", r.Type, r.Title, r.URL))
	}
	return mcputil.TextResult(b.String()), nil
}

func (h *Handlers) CreatePage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Notion not configured."), nil
	}
	title, err := req.RequireString("title")
	if err != nil {
		return mcputil.ErrInvalid("title parameter is required"), nil
	}
	parentID := req.GetString("parent_id", h.client.DefaultDatabaseID())
	content := req.GetString("content", "")

	page, err := h.client.CreatePage(ctx, parentID, title, content)
	if err != nil {
		return h.errH.Wrap("create page", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Created Notion page: %s (%s)", page.Title, page.URL)), nil
}

func (h *Handlers) QueryDatabase(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Notion not configured."), nil
	}
	databaseID := req.GetString("database_id", h.client.DefaultDatabaseID())
	filterJSON := req.GetString("filter", "")
	limit := req.GetInt("limit", 20)

	results, err := h.client.QueryDatabase(ctx, databaseID, filterJSON, limit)
	if err != nil {
		return h.errH.Wrap("query database", err), nil
	}
	if len(results) == 0 {
		return mcputil.TextResult("No results found."), nil
	}
	out, _ := json.MarshalIndent(results, "", "  ")
	return mcputil.TextResult(string(out)), nil
}
