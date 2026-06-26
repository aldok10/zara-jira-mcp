package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handlers) SheetsRead(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Sheets == nil || !h.Sheets.Available() {
		return errorResult("Google Sheets not configured. Set GOOGLE_SHEETS_API_KEY."), nil
	}

	spreadsheetID, err := req.RequireString("spreadsheet_id")
	if err != nil {
		return errorResult("spreadsheet_id is required"), nil
	}
	rangeStr := req.GetString("range", "Sheet1")

	data, err := h.Sheets.ReadRange(ctx, spreadsheetID, rangeStr)
	if err != nil {
		return sanitizedError("Sheets: failed to read range", err), nil
	}
	if len(data) == 0 {
		return textResult("No data found in range."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sheet data (%d rows):\n\n", len(data)))
	for i, row := range data {
		if i == 0 {
			sb.WriteString("| " + strings.Join(row, " | ") + " |\n")
			sb.WriteString("|" + strings.Repeat(" --- |", len(row)) + "\n")
		} else {
			sb.WriteString("| " + strings.Join(row, " | ") + " |\n")
		}
	}
	return textResult(sb.String()), nil
}

func (h *Handlers) SheetsExport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	format := req.GetString("format", "csv")
	data := req.GetString("data", "")

	if data == "" {
		return errorResult("data parameter is required (provide sprint metrics as key:value pairs, one per line)"), nil
	}

	lines := strings.Split(strings.TrimSpace(data), "\n")

	var sb strings.Builder
	switch format {
	case "csv":
		for _, line := range lines {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				sb.WriteString(fmt.Sprintf("%s,%s\n", strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])))
			} else {
				sb.WriteString(line + "\n")
			}
		}
	default:
		for _, line := range lines {
			sb.WriteString(line + "\n")
		}
	}

	return textResult(fmt.Sprintf("Export (%s):\n\n```\n%s```\n\nCopy and paste into your spreadsheet.", format, sb.String())), nil
}
