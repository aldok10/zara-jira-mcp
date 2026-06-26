package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerPortfolioTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("portfolio_overview",
		mcp.WithDescription("Cross-project portfolio overview: issue counts, leads, health per project."),
	), h.PortfolioOverview)

	s.AddTool(mcp.NewTool("portfolio_blockers",
		mcp.WithDescription("All active blockers and open dependencies across the portfolio."),
	), h.PortfolioBlockers)

	s.AddTool(mcp.NewTool("portfolio_workload",
		mcp.WithDescription("Team workload distribution across ALL projects. Finds overloaded people."),
	), h.PortfolioWorkload)

	s.AddTool(mcp.NewTool("portfolio_risks",
		mcp.WithDescription("Aggregate risk radar across all tracked projects."),
	), h.PortfolioRisks)

	s.AddTool(mcp.NewTool("portfolio_summary",
		mcp.WithDescription("AI-generated executive portfolio status summary: health, risks, actions."),
	), h.PortfolioSummary)
}
