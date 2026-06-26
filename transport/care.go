package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerCareTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_daily_delta",
		mcp.WithDescription("What changed since yesterday? Shows new completions, new blockers, burn rate, and who might need help. Your morning briefing."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.DailyDelta)

	s.AddTool(mcp.NewTool("pm_overload_check",
		mcp.WithDescription("Detect team members who may be overloaded. Shows WIP per person, blocked people, and sustainable pace signals. Protects the team."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.OverloadCheck)

	s.AddTool(mcp.NewTool("pm_commitment_check",
		mcp.WithDescription("Are we overcommitting? Compares sprint items against historical completion rate. Flags overcommitment before it causes burnout."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.CommitmentCheck)

	s.AddTool(mcp.NewTool("pm_team_care",
		mcp.WithDescription("Team care report: wellbeing signals, not productivity metrics. Checks workload, chronic carryover, blocked people, retro follow-through. Reminds you that sustainable pace > velocity."),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
	), h.TeamCareReport)
}
