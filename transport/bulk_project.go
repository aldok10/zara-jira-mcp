package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerBulkProjectTools(s *server.MCPServer, h *tools.Handlers) {
	// Bulk operations
	s.AddTool(
		mcp.NewTool("jira_bulk_transition",
			mcp.WithDescription("Transition multiple issues at once. Loops over each key and applies the transition."),
			mcp.WithString("issue_keys", mcp.Required(), mcp.Description("Comma-separated issue keys (e.g. PROJ-1,PROJ-2,PROJ-3)")),
			mcp.WithString("transition_id", mcp.Required(), mcp.Description("Transition ID to apply (from jira_transitions)")),
		),
		h.BulkTransition,
	)

	s.AddTool(
		mcp.NewTool("jira_bulk_assign",
			mcp.WithDescription("Assign multiple issues to one person at once."),
			mcp.WithString("issue_keys", mcp.Required(), mcp.Description("Comma-separated issue keys")),
			mcp.WithString("assignee_id", mcp.Required(), mcp.Description("Assignee account ID")),
		),
		h.BulkAssign,
	)

	s.AddTool(
		mcp.NewTool("jira_bulk_label",
			mcp.WithDescription("Add a label to multiple issues at once."),
			mcp.WithString("issue_keys", mcp.Required(), mcp.Description("Comma-separated issue keys")),
			mcp.WithString("label", mcp.Required(), mcp.Description("Label to add")),
		),
		h.BulkLabel,
	)

	// Project tools
	s.AddTool(
		mcp.NewTool("jira_projects",
			mcp.WithDescription("List all accessible Jira projects with key, name, lead, and type."),
		),
		h.ListProjects,
	)

	s.AddTool(
		mcp.NewTool("jira_project_detail",
			mcp.WithDescription("Get full project details including components and versions."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Project key (e.g. PROJ)")),
		),
		h.ProjectDetail,
	)

	// Raw API access
	s.AddTool(
		mcp.NewTool("jira_raw_request",
			mcp.WithDescription("Make an arbitrary Jira REST API request. Escape hatch for unsupported operations."),
			mcp.WithString("method", mcp.Required(), mcp.Description("HTTP method: GET, POST, PUT, DELETE")),
			mcp.WithString("path", mcp.Required(), mcp.Description("API path (e.g. /rest/api/3/issue/PROJ-1)")),
			mcp.WithString("body", mcp.Description("Request body as JSON string (for POST/PUT)")),
		),
		h.RawRequest,
	)
}
