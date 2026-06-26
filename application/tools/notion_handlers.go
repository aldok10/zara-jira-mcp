package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// NotionSearch searches Notion pages and databases.
func (h *Handlers) NotionSearch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Notion == nil || !h.Notion.Available() {
		return errorResult("Notion not configured. Set NOTION_API_KEY."), nil
	}
	query, err := req.RequireString("query")
	if err != nil {
		return errorResult("query parameter is required"), nil
	}
	limit := req.GetInt("limit", 10)

	results, err := h.Notion.Search(ctx, query, limit)
	if err != nil {
		return sanitizedError("Notion: search failed", err), nil
	}
	if len(results) == 0 {
		return textResult("No results found for: " + query), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Notion search results for '%s' (%d):\n\n", query, len(results)))
	for _, r := range results {
		title := r.Title
		if title == "" {
			title = "(untitled)"
		}
		sb.WriteString(fmt.Sprintf("- [%s] %s | %s\n", r.Type, title, r.URL))
	}
	return textResult(sb.String()), nil
}

// NotionCreatePage creates a new Notion page.
func (h *Handlers) NotionCreatePage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Notion == nil || !h.Notion.Available() {
		return errorResult("Notion not configured. Set NOTION_API_KEY."), nil
	}
	title, err := req.RequireString("title")
	if err != nil {
		return errorResult("title parameter is required"), nil
	}
	parentID := req.GetString("parent_id", h.Notion.DefaultDatabaseID())
	if parentID == "" {
		return errorResult("parent_id is required (or set NOTION_DATABASE_ID as default)"), nil
	}
	content := req.GetString("content", "")

	page, err := h.Notion.CreatePage(ctx, parentID, title, content)
	if err != nil {
		return sanitizedError("Notion: failed to create page", err), nil
	}
	return textResult(fmt.Sprintf("Created page: %s\nURL: %s", page.Title, page.URL)), nil
}

// NotionQueryDB queries a Notion database.
func (h *Handlers) NotionQueryDB(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Notion == nil || !h.Notion.Available() {
		return errorResult("Notion not configured. Set NOTION_API_KEY."), nil
	}
	dbID := req.GetString("database_id", h.Notion.DefaultDatabaseID())
	if dbID == "" {
		return errorResult("database_id is required (or set NOTION_DATABASE_ID as default)"), nil
	}
	filter := req.GetString("filter", "")
	limit := req.GetInt("limit", 20)

	results, err := h.Notion.QueryDatabase(ctx, dbID, filter, limit)
	if err != nil {
		return sanitizedError("Notion: database query failed", err), nil
	}
	if len(results) == 0 {
		return textResult("No results from database query."), nil
	}

	// Format results showing properties
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Database query results (%d items):\n\n", len(results)))
	for i, item := range results {
		props, ok := item["properties"]
		if !ok {
			continue
		}
		data, _ := json.Marshal(props)
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, string(data)))
		if i >= 19 {
			sb.WriteString("... (truncated)\n")
			break
		}
	}
	return textResult(sb.String()), nil
}
