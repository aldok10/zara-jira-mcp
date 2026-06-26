package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerTraceTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_trace_branch",
			mcp.WithDescription("Trace a Jira ticket to its branch in GitHub/GitLab. Shows if branch exists, PRs/MRs status, and whether it's been merged to target branch. Use to verify ticket implementation status."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Jira issue key (e.g. SIT-3658)")),
		),
		h.TraceTicketBranch,
	)
}
