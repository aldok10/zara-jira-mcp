package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerLinkWorklogTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_link_issues",
			mcp.WithDescription("Create a link between two Jira issues (e.g. Blocks, Relates, Duplicates)."),
			mcp.WithString("inward_key", mcp.Required(), mcp.Description("Inward issue key (e.g. PROJ-1)")),
			mcp.WithString("outward_key", mcp.Required(), mcp.Description("Outward issue key (e.g. PROJ-2)")),
			mcp.WithString("link_type", mcp.Required(), mcp.Description("Link type name: Blocks, Relates, Duplicates, etc.")),
		),
		h.LinkIssues,
	)

	s.AddTool(
		mcp.NewTool("jira_link_types",
			mcp.WithDescription("List available issue link types (Blocks, Relates, Duplicates, etc.)."),
		),
		h.LinkTypes,
	)

	s.AddTool(
		mcp.NewTool("jira_worklog_add",
			mcp.WithDescription("Log time spent on a Jira issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
			mcp.WithString("time_spent", mcp.Required(), mcp.Description("Time spent (e.g. 2h, 30m, 1d)")),
			mcp.WithString("comment", mcp.Description("Work description")),
		),
		h.WorklogAdd,
	)

	s.AddTool(
		mcp.NewTool("jira_worklog_list",
			mcp.WithDescription("List worklogs (time entries) for a Jira issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
		),
		h.WorklogList,
	)

	s.AddTool(
		mcp.NewTool("jira_watch",
			mcp.WithDescription("Add a watcher to a Jira issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
			mcp.WithString("account_id", mcp.Required(), mcp.Description("Watcher's account ID (use jira_find_user to look up)")),
		),
		h.Watch,
	)

	s.AddTool(
		mcp.NewTool("jira_watchers",
			mcp.WithDescription("List watchers on a Jira issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
		),
		h.Watchers,
	)

	s.AddTool(
		mcp.NewTool("jira_labels_set",
			mcp.WithDescription("Set labels on a Jira issue (replaces all existing labels)."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
			mcp.WithString("labels", mcp.Required(), mcp.Description("Comma-separated labels")),
		),
		h.LabelsSet,
	)
}
