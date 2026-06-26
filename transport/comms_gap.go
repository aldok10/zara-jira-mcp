package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerCommsGapTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_comms_health",
		mcp.WithDescription("Communication health score (0-100). Measures: decision velocity, blocker resolution speed, action follow-through, stakeholder engagement."),
		mcp.WithNumber("board_id", mcp.Description("Board ID for board-scoped metrics")),
	), h.CommsHealth)

	s.AddTool(mcp.NewTool("pm_silence_detector",
		mcp.WithDescription("Find stakeholders with no recent pulse/activity. Detects ghost stakeholders who may cause surprise objections."),
		mcp.WithNumber("days_threshold", mcp.Description("Days of silence to flag (default: 30)")),
	), h.SilenceDetector)

	s.AddTool(mcp.NewTool("pm_comms_anti_patterns",
		mcp.WithDescription("Detect communication dysfunctions: re-deciding, dead actions, escalation hoarding, ghost stakeholders, blocker silence."),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
	), h.CommsAntiPatterns)

	s.AddTool(mcp.NewTool("pm_nvc_reframe",
		mcp.WithDescription("Rewrite blaming/judgmental messages using Nonviolent Communication (Observation-Feeling-Need-Request). Makes feedback receivable."),
		mcp.WithString("message", mcp.Required(), mcp.Description("The message to reframe")),
	), h.NVCReframe)

	s.AddTool(mcp.NewTool("pm_hard_conversation",
		mcp.WithDescription("Prepare for a difficult conversation. Uses Crucial Conversations (STATE path) + SBI + SCARF. Gets you facts, opening lines, and safety restoration strategies."),
		mcp.WithString("situation", mcp.Required(), mcp.Description("Describe the situation that needs addressing")),
		mcp.WithNumber("board_id", mcp.Description("Board ID for data context")),
		mcp.WithString("person", mcp.Description("Who the conversation is with")),
	), h.HardConversation)

	s.AddTool(mcp.NewTool("pm_trust_signals",
		mcp.WithDescription("Trust indicator dashboard: forecast accuracy, escalation responsiveness, health consistency, decision transparency."),
		mcp.WithNumber("board_id", mcp.Description("Board ID for board-scoped metrics")),
	), h.TrustSignals)

	s.AddTool(mcp.NewTool("pm_lencioni",
		mcp.WithDescription("Diagnose team dysfunction using Lencioni's 5 Dysfunctions pyramid (Trust > Conflict > Commitment > Accountability > Results). Maps Jira data to dysfunction levels with coaching."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.LencioniDysfunction)
}
