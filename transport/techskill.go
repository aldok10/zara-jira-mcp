package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTechSkillTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_tech_glossary",
		mcp.WithDescription("Explain technical concepts in PM-friendly language. Covers: API, CI/CD, PR, tech debt, WIP, cycle time, deployment, microservices, test coverage, and 20+ terms. Ask any term for AI explanation."),
		mcp.WithString("term", mcp.Description("Technical term to explain (e.g. 'api', 'ci/cd', 'code review', 'regression')")),
	), h.TechGlossary)

	s.AddTool(mcp.NewTool("pm_qa_health",
		mcp.WithDescription("QA health check: bug ratio, priority breakdown, stale bugs, quality signals. Helps PM understand code quality state."),
		mcp.WithString("project", mcp.Description("Project key (all if empty)")),
	), h.QAHealthCheck)

	s.AddTool(mcp.NewTool("pm_dev_workflow",
		mcp.WithDescription("Understand developer workflow: what they do daily, what they need from PM, how testing/deployment works. Builds empathy + reduces friction."),
		mcp.WithString("phase", mcp.Description("overview, planning, testing, deployment, code_review")),
	), h.DevWorkflowExplainer)

	s.AddTool(mcp.NewTool("pm_metrics_guide",
		mcp.WithDescription("Engineering metrics explained for PMs: what to track, what to avoid, toxic metrics vs healthy metrics. Prevents measuring the wrong things."),
	), h.EngineeringMetricsExplainer)

	s.AddTool(mcp.NewTool("pm_quality_gate",
		mcp.WithDescription("Sprint quality gate: checks if sprint is release-ready based on open bugs, critical issues, completion rate. Objective release decision."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.SprintQualityGate)
}
