package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerCommunicationTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_status_draft",
		mcp.WithDescription("Generate ready-to-send status update. Auto-pulls project data, formats for audience. BLUF format."),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
		mcp.WithString("audience", mcp.Description("team, executive, stakeholder, po")),
	), h.WriteUpdate)

	s.AddTool(mcp.NewTool("pm_announce_decision",
		mcp.WithDescription("Communicate a decision to the team using DACI. Structured: what was decided, why, what changes, who is D/A/C/I."),
		mcp.WithString("decision", mcp.Required(), mcp.Description("What was decided")),
		mcp.WithString("rationale", mcp.Description("Why this decision")),
		mcp.WithString("audience", mcp.Description("Who needs to know")),
	), h.CommunicateDecision)

	s.AddTool(mcp.NewTool("pm_comms_plan",
		mcp.WithDescription("Generate a communication plan: who needs what info, when, via which channel. Includes DACI role mapping."),
		mcp.WithString("context", mcp.Required(), mcp.Description("What's the situation (new project, major change, incident, release)")),
		mcp.WithString("stakeholders", mcp.Description("Key stakeholders involved")),
		mcp.WithString("urgency", mcp.Description("normal, high, critical")),
	), h.CommunicationPlan)
}
