package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerImprovementTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_improvement_dashboard",
		mcp.WithDescription("Meta-metrics: is the team getting better at getting better? Shows velocity trend, action item completion rate, sprint goal hit rate, retro cadence, blocker resolution speed."),
		mcp.WithNumber("board_id", mcp.Description("Board ID for velocity data")),
	), h.PMImprovementDashboard)

	s.AddTool(mcp.NewTool("pm_bus_factor",
		mcp.WithDescription("Detect single points of failure. Shows workload concentration per person and per area/label. Flags when one person owns >50% of items in any area."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.PMBusFactor)

	s.AddTool(mcp.NewTool("pm_async_standup",
		mcp.WithDescription("Generate async standup from Jira data. Per member: what moved in last 24h, what's in progress, active blockers. Copy-paste ready for Slack/Lark."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.PMAsyncStandup)

	s.AddTool(mcp.NewTool("pm_sprint_goal_track",
		mcp.WithDescription("Record sprint goal outcome (hit or miss). Shows cumulative hit rate. Research shows this single metric is the best predictor of team effectiveness."),
		mcp.WithString("sprint_name", mcp.Required(), mcp.Description("Sprint name")),
		mcp.WithBoolean("hit", mcp.Description("Did the team hit the sprint goal? (default: false)")),
		mcp.WithString("notes", mcp.Description("Optional context about the outcome")),
	), h.PMSprintGoalTrack)
}
