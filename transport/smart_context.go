package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerSmartContextTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_predict_blockers",
		mcp.WithDescription("Predict who is likely to get blocked based on historical patterns. Proactive impediment prevention."),
		mcp.WithNumber("board_id", mcp.Description("Board ID for current assignments")),
	), h.PredictiveBlockers)

	s.AddTool(mcp.NewTool("pm_sprint_similarity",
		mcp.WithDescription("Compare current sprint signals to historical sprints. Warns if current state matches a sprint that failed."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.SprintSimilarity)

	s.AddTool(mcp.NewTool("pm_early_warning",
		mcp.WithDescription("Multi-signal early warning system: chronic blockers, pace behind, unmitigated risks, declining velocity, ignored retro actions. Leading indicators before sprint failure."),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
	), h.EarlyWarningSystem)

	// Smart routing tools (NL interface)
	s.AddTool(mcp.NewTool("pm_smart",
		mcp.WithDescription("Natural language PM assistant. Ask anything: 'who is blocked?', 'sprint health?', 'when will we finish?'. Routes to the right tool automatically."),
		mcp.WithString("ask", mcp.Required(), mcp.Description("What you want to know (natural language)")),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
	), h.PMSmart)

	s.AddTool(mcp.NewTool("pm_do",
		mcp.WithDescription("Natural language action: 'create task X', 'record risk Y', 'record decision Z'. Routes to the right action."),
		mcp.WithString("what", mcp.Required(), mcp.Description("What to do (natural language)")),
	), h.PMDo)

	s.AddTool(mcp.NewTool("pm_report",
		mcp.WithDescription("Generate any report by type: status, executive, release_notes, weekly, health, velocity, scorecard."),
		mcp.WithString("type", mcp.Required(), mcp.Description("Report type")),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
	), h.PMReport)
}
