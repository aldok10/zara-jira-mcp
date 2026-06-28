package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	smcp "github.com/aldok10/zara-jira-mcp/modules/sprint/interfaces/mcp"
)

func RegisterSprintTools(s *server.MCPServer, h *smcp.Handlers) {
	s.AddTool(
		mcp.NewTool("pm",
			mcp.WithDescription("Quick project status. THE main tool for PMs."),
			mcp.WithNumber("board_id", mcp.Description("Board ID")),
		),
		h.PMQuickStatus,
	)
	s.AddTool(
		mcp.NewTool("pm_create",
			mcp.WithDescription("Create work item anywhere."),
			mcp.WithString("title", mcp.Required(), mcp.Description("What needs to be done")),
		),
		h.PMCreate,
	)
	s.AddTool(
		mcp.NewTool("pm_decide",
			mcp.WithDescription("Record a decision."),
			mcp.WithString("what", mcp.Required(), mcp.Description("What was decided")),
		),
		h.PMDecide,
	)
	s.AddTool(
		mcp.NewTool("pm_risk",
			mcp.WithDescription("Record a risk."),
			mcp.WithString("what", mcp.Required(), mcp.Description("What could go wrong")),
		),
		h.PMRisk,
	)
	s.AddTool(
		mcp.NewTool("pm_next",
			mcp.WithDescription("Suggest next PM action."),
			mcp.WithNumber("board_id", mcp.Description("Board ID")),
		),
		h.PMNext,
	)
}
