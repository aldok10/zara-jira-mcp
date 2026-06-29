package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	shmcp "github.com/aldok10/zara-jira-mcp/shared/infrastructure/sheets/mcp"
)

func RegisterSheetsTools(s *server.MCPServer, h *shmcp.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_sheet_read",
			mcp.WithDescription("Read data from a Google Sheet (must be publicly accessible or shared)."),
			mcp.WithString("spreadsheet_id", mcp.Required(), mcp.Description("Google Sheets spreadsheet ID (from URL)")),
			mcp.WithString("range", mcp.Description("Cell range (e.g. Sheet1!A1:D10). Default: first sheet")),
		),
		h.ReadRange,
	)
}
