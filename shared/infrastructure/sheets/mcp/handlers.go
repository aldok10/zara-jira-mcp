// Package mcp provides MCP tool handlers for Google Sheets.
package mcp

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/sheets"
)

// Handlers exposes Google Sheets operations as MCP tool handlers.
type Handlers struct {
	client *sheets.Client
	errH   *mcputil.ErrorHandler
}

func NewHandlers(client *sheets.Client) *Handlers {
	return &Handlers{
		client: client,
		errH:   mcputil.NewErrorHandler(nil),
	}
}

// ReadRange reads data from a Google Sheet.
func (h *Handlers) ReadRange(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if !h.client.Available() {
		return mcputil.ErrorResult("Google Sheets not configured. Set GOOGLE_SHEETS_API_KEY."), nil
	}

	spreadsheetID, err := req.RequireString("spreadsheet_id")
	if err != nil {
		return mcputil.ErrorResult("Missing required parameter: spreadsheet_id"), nil
	}
	rangeStr := req.GetString("range", "")

	values, err := h.client.ReadRange(ctx, spreadsheetID, rangeStr)
	if err != nil {
		return h.errH.Wrap("read sheet", err), nil
	}

	if len(values) == 0 {
		return mcputil.TextResult("Sheet is empty or range not found."), nil
	}

	var b strings.Builder
	for _, row := range values {
		b.WriteString(strings.Join(row, " | "))
		b.WriteString("\n")
	}
	return mcputil.TextResult(b.String()), nil
}
