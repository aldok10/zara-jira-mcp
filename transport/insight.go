package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerInsightTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_tool_usage",
		mcp.WithDescription("Show which PM tools are used most/least."),
		mcp.WithNumber("days", mcp.Description("Lookback period in days (default: 30)")),
	), h.PMToolUsage)

	s.AddTool(mcp.NewTool("pm_calibration",
		mcp.WithDescription("Forecast accuracy: committed vs delivered over time."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.PMCalibrationReport)

	s.AddTool(mcp.NewTool("pm_meeting_roi",
		mcp.WithDescription("Meeting effectiveness: decisions+actions per meeting."),
	), h.PMMeetingROI)

	s.AddTool(mcp.NewTool("pm_notification_budget",
		mcp.WithDescription("Check notification budget: sent today vs daily limit."),
	), h.PMNotificationBudgetCheck)

	s.AddTool(mcp.NewTool("pm_collaboration_signal",
		mcp.WithDescription("Detect collaboration patterns, knowledge silos, and isolation risks from sprint data. Shows cross-label collaboration and workload balance."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.PMCollaborationSignal)

	s.AddTool(mcp.NewTool("pm_ai_health",
		mcp.WithDescription("Evaluate AI adoption health. Checks process fundamentals, tool usage diversity, and provides research-backed AI recommendations (BCG 2026, PMI 2025)."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.PMAIHealth)

	s.AddTool(mcp.NewTool("pm_notification_effectiveness",
		mcp.WithDescription("Track notification effectiveness: volume by channel/severity, fatigue analysis, budget compliance, and recommendations."),
	), h.PMNotificationEffectiveness)

	s.AddTool(mcp.NewTool("pm_notification_record_action",
		mcp.WithDescription("Record user action on a notification for effectiveness tracking."),
		mcp.WithString("channel", mcp.Description("Notification channel (e.g., slack, lark)")),
		mcp.WithString("title", mcp.Description("Notification title")),
		mcp.WithString("action", mcp.Description("Action taken (acknowledged, clicked, dismissed)")),
	), h.PMNotificationRecordAction)
}
