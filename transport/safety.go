package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerSafetyTools(s *server.MCPServer, h *tools.Handlers) {
	// Communication Templates (P1)
	s.AddTool(mcp.NewTool("pm_communicate",
		mcp.WithDescription("Generate Minto Pyramid-structured message for any audience. Conclusion first, then arguments, then data. Adapts tone per audience."),
		mcp.WithString("topic", mcp.Required(), mcp.Description("What to communicate")),
		mcp.WithString("audience", mcp.Required(), mcp.Description("exec, team, po, stakeholder")),
		mcp.WithNumber("board_id", mcp.Description("Board ID for sprint context")),
	), h.PMCommunicate)

	s.AddTool(mcp.NewTool("pm_feedback_prep",
		mcp.WithDescription("Generate SBI feedback (Situation-Behavior-Impact). Data-backed when available. Supports positive and constructive feedback."),
		mcp.WithString("member", mcp.Required(), mcp.Description("Team member name")),
		mcp.WithString("observation", mcp.Required(), mcp.Description("What you observed")),
		mcp.WithString("type", mcp.Description("positive or constructive (default: constructive)")),
	), h.PMFeedbackPrep)

	s.AddTool(mcp.NewTool("pm_escalation_draft",
		mcp.WithDescription("Generate pyramid-structured escalation: 1-line ask, context, impact, next step, deadline."),
		mcp.WithString("issue", mcp.Required(), mcp.Description("What needs escalation")),
		mcp.WithString("severity", mcp.Description("critical, high, medium (default: high)")),
		mcp.WithString("deadline", mcp.Description("When resolution is needed")),
	), h.PMEscalationDraft)

	s.AddTool(mcp.NewTool("pm_decision_record",
		mcp.WithDescription("Create an Architecture Decision Record (ADR). Stores to memory. Formats: Status, Context, Decision, Alternatives, Consequences."),
		mcp.WithString("title", mcp.Required(), mcp.Description("Decision title")),
		mcp.WithString("decision", mcp.Required(), mcp.Description("What was decided")),
		mcp.WithString("context", mcp.Description("What situation led to this")),
		mcp.WithString("alternatives", mcp.Description("Alternatives considered")),
		mcp.WithString("consequences", mcp.Description("Known consequences")),
	), h.PMDecisionRecordEnhanced)

	// Psychological Safety (P2)
	s.AddTool(mcp.NewTool("pm_safety_survey",
		mcp.WithDescription("Record psychological safety survey (7 questions from Project Aristotle). Returns average score."),
		mcp.WithString("sprint_name", mcp.Required(), mcp.Description("Sprint name")),
		mcp.WithString("responses", mcp.Required(), mcp.Description("JSON: {\"member\": {\"q1\": 4, \"q2\": 3, \"q3\": 5, \"q4\": 4, \"q5\": 3, \"q6\": 4, \"q7\": 5}}")),
	), h.PMSafetySurvey)

	s.AddTool(mcp.NewTool("pm_safety_trend",
		mcp.WithDescription("Show psychological safety scores over time, grouped by sprint."),
	), h.PMSafetyTrend)

	s.AddTool(mcp.NewTool("pm_team_aristotle",
		mcp.WithDescription("Full 5-pillar team assessment (Project Aristotle): Safety, Dependability, Clarity, Meaning, Impact. AI-powered with data from all PM sources."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.PMTeamAristotle)
}
