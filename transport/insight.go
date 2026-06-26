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
}
