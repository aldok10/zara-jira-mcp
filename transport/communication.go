package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerCommunicationTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_compose",
		mcp.WithDescription("Write effective PM messages. Adapts content for target audience (executive, developer, stakeholder). Uses Pyramid Principle + BLUF."),
		mcp.WithString("message", mcp.Required(), mcp.Description("What you want to communicate")),
		mcp.WithString("audience", mcp.Description("executive, developer, stakeholder, team (default: team)")),
		mcp.WithString("tone", mcp.Description("formal, casual, urgent (default: professional)")),
	), h.AdaptMessage)

	s.AddTool(mcp.NewTool("pm_status_draft",
		mcp.WithDescription("Generate ready-to-send status update. Auto-pulls project data, formats for audience. BLUF format."),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
		mcp.WithString("audience", mcp.Description("team, executive, stakeholder, po")),
	), h.WriteUpdate)

	s.AddTool(mcp.NewTool("pm_feedback_coach",
		mcp.WithDescription("Help give effective feedback using SBI (Situation-Behavior-Impact) + Radical Candor. Get exact words to say."),
		mcp.WithString("situation", mcp.Required(), mcp.Description("Describe what happened")),
		mcp.WithString("person", mcp.Description("Who the feedback is for")),
	), h.GiveFeedback)

	s.AddTool(mcp.NewTool("pm_escalate_message",
		mcp.WithDescription("Write a structured escalation using SCQA (Situation-Complication-Question-Answer). Clear, concise, actionable."),
		mcp.WithString("issue", mcp.Required(), mcp.Description("What needs escalation")),
		mcp.WithString("impact", mcp.Description("Business impact")),
		mcp.WithString("ask", mcp.Description("What you need from leadership")),
	), h.EscalateWithSCQA)

	s.AddTool(mcp.NewTool("pm_announce_decision",
		mcp.WithDescription("Communicate a decision to the team. Structured: what was decided, why, what changes, what stays the same."),
		mcp.WithString("decision", mcp.Required(), mcp.Description("What was decided")),
		mcp.WithString("rationale", mcp.Description("Why this decision")),
		mcp.WithString("audience", mcp.Description("Who needs to know")),
	), h.CommunicateDecision)

	s.AddTool(mcp.NewTool("pm_comms_plan",
		mcp.WithDescription("Generate a communication plan: who needs what info, when, via which channel. Based on stakeholder analysis."),
		mcp.WithString("context", mcp.Required(), mcp.Description("What's the situation (new project, major change, incident, release)")),
	), h.CommunicationPlan)
}
