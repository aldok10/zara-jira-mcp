package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerVersionTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_attachments",
			mcp.WithDescription("List attachments on a Jira issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key")),
		),
		h.GetAttachments,
	)

	s.AddTool(
		mcp.NewTool("jira_versions",
			mcp.WithDescription("List project versions/releases."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project key")),
		),
		h.GetVersions,
	)

	s.AddTool(
		mcp.NewTool("jira_version_create",
			mcp.WithDescription("Create a new project version for release tracking."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project key")),
			mcp.WithString("name", mcp.Required(), mcp.Description("Version name (e.g. v1.2.0)")),
			mcp.WithString("description", mcp.Description("Version description")),
		),
		h.CreateVersion,
	)

	s.AddTool(
		mcp.NewTool("jira_version_release",
			mcp.WithDescription("Mark a version as released."),
			mcp.WithString("version_id", mcp.Required(), mcp.Description("Version ID (from jira_versions)")),
		),
		h.ReleaseVersion,
	)

	s.AddTool(
		mcp.NewTool("jira_components",
			mcp.WithDescription("List project components with leads."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project key")),
		),
		h.GetComponents,
	)

	s.AddTool(
		mcp.NewTool("jira_fields",
			mcp.WithDescription("List all available Jira fields (system + custom). Useful for finding custom field IDs."),
			mcp.WithBoolean("custom_only", mcp.Description("Show only custom fields (default: false)")),
		),
		h.GetFields,
	)
}
