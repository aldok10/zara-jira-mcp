package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerSafetyTools(s *server.MCPServer, h *tools.Handlers) {
	// Communication Templates (P1)
	s.AddTool(mcp.NewTool("pm_communicate",
		mcp.WithDescription("Minto Pyramid message for any audience. Rewrites or generates from data."),
		mcp.WithString("topic", mcp.Required(), mcp.Description("What to communicate")),
		mcp.WithString("audience", mcp.Required(), mcp.Description("exec, team, po, stakeholder")),
		mcp.WithString("message", mcp.Description("Optional: existing message to rewrite for audience")),
		mcp.WithNumber("board_id", mcp.Description("Board ID for sprint context")),
	), h.PMCommunicate)

	s.AddTool(mcp.NewTool("pm_feedback_prep",
		mcp.WithDescription("SBI feedback (Situation-Behavior-Impact). Data-backed when available."),
		mcp.WithString("member", mcp.Required(), mcp.Description("Team member name")),
		mcp.WithString("observation", mcp.Required(), mcp.Description("What you observed")),
		mcp.WithString("type", mcp.Description("positive or constructive (default: constructive)")),
	), h.PMFeedbackPrep)

	s.AddTool(mcp.NewTool("pm_escalation_draft",
		mcp.WithDescription("Escalation draft: 1-line ask, context, impact, next step, deadline."),
		mcp.WithString("issue", mcp.Required(), mcp.Description("What needs escalation")),
		mcp.WithString("severity", mcp.Description("critical, high, medium (default: high)")),
		mcp.WithString("deadline", mcp.Description("When resolution is needed")),
	), h.PMEscalationDraft)

	s.AddTool(mcp.NewTool("pm_decision_record",
		mcp.WithDescription("Create ADR (Architecture Decision Record). Stores to memory."),
		mcp.WithString("title", mcp.Required(), mcp.Description("Decision title")),
		mcp.WithString("decision", mcp.Required(), mcp.Description("What was decided")),
		mcp.WithString("context", mcp.Description("What situation led to this")),
		mcp.WithString("alternatives", mcp.Description("Alternatives considered")),
		mcp.WithString("consequences", mcp.Description("Known consequences")),
	), h.PMDecisionRecordEnhanced)

	// Psychological Safety (P2)
	s.AddTool(mcp.NewTool("pm_safety_survey",
		mcp.WithDescription("Record psych safety survey (7 Aristotle questions). Returns avg score."),
		mcp.WithString("sprint_name", mcp.Required(), mcp.Description("Sprint name")),
		mcp.WithString("responses", mcp.Required(), mcp.Description("JSON: {\"member\": {\"q1\": 4, \"q2\": 3, \"q3\": 5, \"q4\": 4, \"q5\": 3, \"q6\": 4, \"q7\": 5}}")),
	), h.PMSafetySurvey)

	s.AddTool(mcp.NewTool("pm_safety_trend",
		mcp.WithDescription("Show psychological safety scores over time, grouped by sprint."),
	), h.PMSafetyTrend)

	s.AddTool(mcp.NewTool("pm_team_aristotle",
		mcp.WithDescription("5-pillar Aristotle assessment: Safety, Dependability, Clarity, Meaning."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.PMTeamAristotle)
}
