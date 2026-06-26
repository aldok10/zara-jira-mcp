package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerCommunicationTools(s *server.MCPServer, h *tools.Handlers) {

	s.AddTool(mcp.NewTool("pm_announce_decision",
		mcp.WithDescription("Announce decision using DACI structure."),
		mcp.WithString("decision", mcp.Required(), mcp.Description("What was decided")),
		mcp.WithString("rationale", mcp.Description("Why this decision")),
		mcp.WithString("audience", mcp.Description("Who needs to know")),
	), h.CommunicateDecision)

	s.AddTool(mcp.NewTool("pm_comms_plan",
		mcp.WithDescription("Communication plan: who, what, when, which channel. Includes DACI."),
		mcp.WithString("context", mcp.Required(), mcp.Description("What's the situation (new project, major change, incident, release)")),
		mcp.WithString("stakeholders", mcp.Description("Key stakeholders involved")),
		mcp.WithString("urgency", mcp.Description("normal, high, critical")),
	), h.CommunicationPlan)

	s.AddTool(mcp.NewTool("pm_raci",
		mcp.WithDescription("Generate RACI matrix from Jira sprint assignments."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		mcp.WithString("accountable", mcp.Description("Accountable person (default: reporter)")),
	), h.PMRACI)
}
