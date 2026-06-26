package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	s *server.MCPServer
}

func NewMCPServer(handlers *tools.Handlers) *MCPServer {
	s := server.NewMCPServer(
		"zara-jira-mcp",
		"0.1.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	registerJiraTools(s, handlers)
	registerAITools(s, handlers)
	registerLarkTools(s, handlers)

	return &MCPServer{s: s}
}

func (m *MCPServer) Server() *server.MCPServer {
	return m.s
}

func registerJiraTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_search",
			mcp.WithDescription("Search Jira issues using JQL. Returns key, summary, status, priority, assignee."),
			mcp.WithString("jql", mcp.Required(), mcp.Description("JQL query string (e.g. 'project = PROJ AND status != Done')")),
			mcp.WithNumber("max_results", mcp.Description("Maximum results to return (default 20, max 50)")),
		),
		h.SearchIssues,
	)

	s.AddTool(
		mcp.NewTool("jira_get_issue",
			mcp.WithDescription("Get full details of a single Jira issue by key."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
		),
		h.GetIssue,
	)

	s.AddTool(
		mcp.NewTool("jira_boards",
			mcp.WithDescription("List all accessible Jira boards with their IDs and types."),
		),
		h.GetBoards,
	)

	s.AddTool(
		mcp.NewTool("jira_sprint_summary",
			mcp.WithDescription("Get active sprint status breakdown and issue list for a board."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID (get from jira_boards)")),
		),
		h.GetSprintSummary,
	)
}

func registerAITools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_ai_analyze",
			mcp.WithDescription("AI-powered analysis of Jira tickets. Ask questions like: 'What are the blockers?', 'Which tickets are stale?', 'Sprint health?'. The AI reads your Jira data and provides PM-relevant insights."),
			mcp.WithString("query", mcp.Required(), mcp.Description("Your question about the project/tickets")),
			mcp.WithString("jql", mcp.Description("JQL to scope the analysis (default: all unresolved, ordered by updated)")),
			mcp.WithNumber("max_results", mcp.Description("Max tickets to analyze (default 30)")),
		),
		h.AIAnalyze,
	)

	s.AddTool(
		mcp.NewTool("jira_ai_sprint_report",
			mcp.WithDescription("Generate an AI-powered sprint report with health assessment, blockers, and recommendations. Optionally sends to Lark."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID for the sprint")),
			mcp.WithBoolean("send_to_lark", mcp.Description("If true, also sends the report to the configured Lark group")),
		),
		h.AISprintReport,
	)
}

func registerLarkTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_notify_lark",
			mcp.WithDescription("Send a message to the configured Lark group. Use for sharing summaries, alerts, or updates."),
			mcp.WithString("title", mcp.Description("Card title (default: 'Jira Update')")),
			mcp.WithString("content", mcp.Required(), mcp.Description("Message content in markdown format")),
		),
		h.NotifyLark,
	)
}
