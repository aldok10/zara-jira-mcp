package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerFeedbackLoopTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_cadence_check",
		mcp.WithDescription("Are you meeting communication cadences? Shows last activity for each key PM responsibility."),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
	), h.CommsHealth)

	s.AddTool(mcp.NewTool("pm_comms_nudge",
		mcp.WithDescription("Proactive nudges: what needs attention today from communication and team health signals."),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
	), h.CommsHealth)

	s.AddTool(mcp.NewTool("pm_feedback_log",
		mcp.WithDescription("Record feedback given to a team member. Tracks follow-up."),
		mcp.WithString("person", mcp.Required(), mcp.Description("Team member name")),
		mcp.WithString("topic", mcp.Required(), mcp.Description("Feedback topic")),
		mcp.WithString("type", mcp.Description("positive or constructive (default: constructive)")),
		mcp.WithNumber("follow_up_days", mcp.Description("Follow up in N days (default: 7)")),
	), h.GiveFeedback)

	s.AddTool(mcp.NewTool("pm_feedback_due",
		mcp.WithDescription("Which feedback follow-ups are overdue?"),
	), h.GetActionItems)

	s.AddTool(mcp.NewTool("pm_feedback_close",
		mcp.WithDescription("Close a feedback loop. Record the outcome."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Feedback log ID")),
		mcp.WithString("outcome", mcp.Required(), mcp.Description("improved, no_change, escalated, acknowledged")),
	), h.RecordLearning)

	s.AddTool(mcp.NewTool("pm_comms_effectiveness",
		mcp.WithDescription("How effective is team communication overall? Aggregate score from multiple signals."),
		mcp.WithNumber("sprints", mcp.Description("Sprints to look back (default: 5)")),
	), h.CommsHealth)
}
