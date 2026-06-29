package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/confluence"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)

type Handlers struct {
	client *confluence.Client
	errH   *mcputil.ErrorHandler
}

func NewHandlers(client *confluence.Client) *Handlers {
	return &Handlers{
		client: client,
		errH:   mcputil.NewErrorHandler(nil),
	}
}

func (h *Handlers) SearchPages(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Confluence not configured. Set CONFLUENCE_BASE_URL and CONFLUENCE_API_TOKEN."), nil
	}
	query, err := req.RequireString("query")
	if err != nil {
		return mcputil.ErrInvalid("query parameter is required"), nil
	}
	limit := req.GetInt("limit", 10)

	pages, err := h.client.SearchPages(ctx, query, limit)
	if err != nil {
		return h.errH.Wrap("search pages", err), nil
	}
	if len(pages) == 0 {
		return mcputil.TextResult("No pages found."), nil
	}
	var b strings.Builder
	for _, p := range pages {
		b.WriteString(fmt.Sprintf("[%s] %s — %s (%s)\n", p.SpaceKey, p.Title, p.WebURL, p.Type))
	}
	return mcputil.TextResult(b.String()), nil
}

func (h *Handlers) GetPage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Confluence not configured."), nil
	}
	pageID, err := req.RequireString("page_id")
	if err != nil {
		return mcputil.ErrInvalid("page_id parameter is required"), nil
	}

	page, err := h.client.GetPage(ctx, pageID)
	if err != nil {
		return h.errH.Wrap("get page", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Title: %s (v%d, %s)\n\n%s", page.Title, page.Version, page.SpaceKey, page.Body)), nil
}

func (h *Handlers) CreatePage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Confluence not configured."), nil
	}
	spaceKey, err := req.RequireString("space_key")
	if err != nil {
		return mcputil.ErrInvalid("space_key parameter is required"), nil
	}
	title, err := req.RequireString("title")
	if err != nil {
		return mcputil.ErrInvalid("title parameter is required"), nil
	}
	body := req.GetString("body", "")
	parentID := req.GetString("parent_id", "")

	page, err := h.client.CreatePage(ctx, spaceKey, title, body, parentID)
	if err != nil {
		return h.errH.Wrap("create page", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Created page: %s (%s)", page.Title, page.WebURL)), nil
}
