package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerSheetsTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_sheet_read",
		mcp.WithDescription("Read data from a Google Sheet (must be publicly accessible or shared)."),
		mcp.WithString("spreadsheet_id", mcp.Required(), mcp.Description("Google Sheets spreadsheet ID (from URL)")),
		mcp.WithString("range", mcp.Description("Cell range (e.g. Sheet1!A1:D10). Default: Sheet1")),
	), h.SheetsRead)

	s.AddTool(mcp.NewTool("pm_sheet_export",
		mcp.WithDescription("Export sprint data to CSV format for pasting into spreadsheets."),
		mcp.WithString("data", mcp.Required(), mcp.Description("Data as key:value pairs, one per line")),
		mcp.WithString("format", mcp.Description("Export format: csv (default)")),
	), h.SheetsExport)
}
