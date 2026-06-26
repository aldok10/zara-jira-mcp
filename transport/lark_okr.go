package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerLarkOKRTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("lark_okr_periods",
		mcp.WithDescription("List Lark OKR periods."),
	), h.LarkOKRPeriods)

	s.AddTool(mcp.NewTool("lark_okr_pull",
		mcp.WithDescription("Pull OKRs from Lark for a user and period."),
		mcp.WithString("user_id", mcp.Required(), mcp.Description("Lark user open_id")),
		mcp.WithString("period_id", mcp.Required(), mcp.Description("Period ID")),
	), h.LarkOKRPull)

	s.AddTool(mcp.NewTool("lark_okr_sync_progress",
		mcp.WithDescription("Push progress to Lark OKR (dry_run default)."),
		mcp.WithString("okr_id", mcp.Required(), mcp.Description("Lark OKR target ID")),
		mcp.WithString("content", mcp.Required(), mcp.Description("Progress update")),
		mcp.WithBoolean("dry_run", mcp.Description("Preview only (default: true)")),
	), h.LarkOKRSyncProgress)
}
