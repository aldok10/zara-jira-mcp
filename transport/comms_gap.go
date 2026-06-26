package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerCommsGapTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_comms_health",
		mcp.WithDescription("Comms health (0-100): decision speed, blocker resolution, follow-through."),
		mcp.WithNumber("board_id", mcp.Description("Board ID for board-scoped metrics")),
	), h.CommsHealth)

	s.AddTool(mcp.NewTool("pm_silence_detector",
		mcp.WithDescription("Find silent stakeholders. Detects ghosts who may cause surprise issues."),
		mcp.WithNumber("days_threshold", mcp.Description("Days of silence to flag (default: 30)")),
	), h.SilenceDetector)

	s.AddTool(mcp.NewTool("pm_comms_anti_patterns",
		mcp.WithDescription("Detect comms dysfunctions: re-deciding, dead actions, escalation hoarding."),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
	), h.CommsAntiPatterns)

	s.AddTool(mcp.NewTool("pm_nvc_reframe",
		mcp.WithDescription("Rewrite messages using NVC (Observation-Feeling-Need-Request)."),
		mcp.WithString("message", mcp.Required(), mcp.Description("The message to reframe")),
	), h.NVCReframe)

	s.AddTool(mcp.NewTool("pm_hard_conversation",
		mcp.WithDescription("Prep difficult conversation: STATE + SBI + SCARF. Facts + opening lines."),
		mcp.WithString("situation", mcp.Required(), mcp.Description("Describe the situation that needs addressing")),
		mcp.WithNumber("board_id", mcp.Description("Board ID for data context")),
		mcp.WithString("person", mcp.Description("Who the conversation is with")),
	), h.HardConversation)

	s.AddTool(mcp.NewTool("pm_trust_signals",
		mcp.WithDescription("Trust indicators: forecast accuracy, escalation response, transparency."),
		mcp.WithNumber("board_id", mcp.Description("Board ID for board-scoped metrics")),
	), h.TrustSignals)

	s.AddTool(mcp.NewTool("pm_lencioni",
		mcp.WithDescription("Lencioni 5 Dysfunctions diagnosis from Jira data with coaching advice."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.LencioniDysfunction)

	// Communication lifecycle
	s.AddTool(mcp.NewTool("pm_cadence_check",
		mcp.WithDescription("Check communication cadences. Flags overdue items."),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
	), h.CadenceCheck)

	s.AddTool(mcp.NewTool("pm_comms_nudge",
		mcp.WithDescription("Proactive communication suggestions from team signals."),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
	), h.CommsNudge)

	s.AddTool(mcp.NewTool("pm_comms_effectiveness",
		mcp.WithDescription("Communication effectiveness score (0-100)."),
		mcp.WithNumber("sprints", mcp.Description("Sprint window (default: 5)")),
	), h.CommsEffectiveness)

	s.AddTool(mcp.NewTool("pm_conversation_prep",
		mcp.WithDescription("Framework-based conversation prep (SBI, NVC, STATE, SCARF)."),
		mcp.WithString("type", mcp.Required(), mcp.Description("performance, conflict, scope, bad_news, recognition")),
		mcp.WithString("context", mcp.Required(), mcp.Description("Describe the situation")),
		mcp.WithString("person", mcp.Description("Who the conversation is with")),
		mcp.WithNumber("board_id", mcp.Description("Board ID for data context")),
	), h.ConversationPrep)

	// Feedback lifecycle
	s.AddTool(mcp.NewTool("pm_feedback_log",
		mcp.WithDescription("Record feedback given. Auto-schedules follow-up."),
		mcp.WithString("person", mcp.Required(), mcp.Description("Who received feedback")),
		mcp.WithString("topic", mcp.Required(), mcp.Description("What the feedback was about")),
		mcp.WithString("type", mcp.Description("constructive, positive, coaching")),
		mcp.WithNumber("follow_up_days", mcp.Description("Days until follow-up (default: 7)")),
	), h.FeedbackLog)

	s.AddTool(mcp.NewTool("pm_feedback_due",
		mcp.WithDescription("Show overdue feedback follow-ups."),
	), h.FeedbackDue)

	s.AddTool(mcp.NewTool("pm_feedback_close",
		mcp.WithDescription("Mark feedback as followed up."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Feedback log ID")),
		mcp.WithString("outcome", mcp.Required(), mcp.Description("improved, no_change, escalated, acknowledged")),
	), h.FeedbackClose)

	// Sentiment & empathy
	s.AddTool(mcp.NewTool("pm_sentiment",
		mcp.WithDescription("Team sentiment from multiple signals + coaching suggestion."),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
		mcp.WithString("message", mcp.Description("Your concern or observation")),
	), h.PMSentiment)

	s.AddTool(mcp.NewTool("pm_context_note",
		mcp.WithDescription("Record the human story behind data (why something is stuck)."),
		mcp.WithString("subject", mcp.Required(), mcp.Description("Person, blocker, or situation")),
		mcp.WithString("note", mcp.Required(), mcp.Description("What's really going on")),
		mcp.WithString("sentiment", mcp.Description("positive, neutral, frustrated, anxious, overwhelmed")),
	), h.PMContextNote)

	s.AddTool(mcp.NewTool("pm_context_recall",
		mcp.WithDescription("Recall context notes about a subject."),
		mcp.WithString("subject", mcp.Description("Filter by subject")),
	), h.PMContextRecall)
}
