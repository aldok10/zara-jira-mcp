package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	jmcp "github.com/aldok10/zara-jira-mcp/modules/jira/interfaces/mcp"
)

func RegisterJiraTools(s *server.MCPServer, h *jmcp.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_search",
			mcp.WithDescription("Search Jira issues using JQL."),
			mcp.WithString("jql", mcp.Required(), mcp.Description("JQL query string")),
			mcp.WithNumber("max_results", mcp.Description("Maximum results (default 20)")),
		),
		h.SearchIssues,
	)
	s.AddTool(
		mcp.NewTool("jira_get_issue",
			mcp.WithDescription("Get full details of a Jira issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
		),
		h.GetIssue,
	)
	s.AddTool(
		mcp.NewTool("jira_boards",
			mcp.WithDescription("List all accessible Jira boards."),
		),
		h.GetBoards,
	)
	s.AddTool(
		mcp.NewTool("jira_create",
			mcp.WithDescription("Create a new Jira issue."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project key (e.g. PROJ)")),
			mcp.WithString("summary", mcp.Required(), mcp.Description("Issue title/summary")),
			mcp.WithString("issue_type", mcp.Description("Issue type: Task, Bug, Story (default Task)")),
			mcp.WithString("description", mcp.Description("Detailed description")),
			mcp.WithString("priority", mcp.Description("Priority: Highest, High, Medium, Low, Lowest")),
			mcp.WithString("assignee_id", mcp.Description("Assignee account ID")),
			mcp.WithString("labels", mcp.Description("Comma-separated labels")),
		),
		h.CreateIssue,
	)
	s.AddTool(
		mcp.NewTool("jira_transition",
			mcp.WithDescription("Transition a Jira issue to a new status."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
			mcp.WithString("transition_id", mcp.Required(), mcp.Description("Transition ID (from jira_transitions)")),
		),
		h.TransitionIssue,
	)
	s.AddTool(
		mcp.NewTool("jira_transitions",
			mcp.WithDescription("Get available status transitions for an issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
		),
		h.GetTransitions,
	)
	s.AddTool(
		mcp.NewTool("jira_assign",
			mcp.WithDescription("Assign a Jira issue to a user by account ID."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
			mcp.WithString("account_id", mcp.Required(), mcp.Description("Assignee account ID (from jira_find_user)")),
		),
		h.AssignIssue,
	)
	s.AddTool(
		mcp.NewTool("jira_find_user",
			mcp.WithDescription("Search Jira users by name or email."),
			mcp.WithString("query", mcp.Required(), mcp.Description("Search by name or email")),
		),
		h.FindUser,
	)
	s.AddTool(
		mcp.NewTool("jira_add_comment",
			mcp.WithDescription("Add a comment to a Jira issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
			mcp.WithString("body", mcp.Required(), mcp.Description("Comment text")),
		),
		h.AddComment,
	)
	s.AddTool(
		mcp.NewTool("jira_sprints",
			mcp.WithDescription("List sprints for a board. Filter by state: active, future, closed."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithString("state", mcp.Description("Filter: active, future, closed (default: all)")),
		),
		h.GetSprints,
	)
	s.AddTool(
		mcp.NewTool("jira_start_sprint",
			mcp.WithDescription("Start a sprint with start and end dates."),
			mcp.WithNumber("sprint_id", mcp.Required(), mcp.Description("Sprint ID")),
			mcp.WithString("start_date", mcp.Required(), mcp.Description("Start date (YYYY-MM-DD)")),
			mcp.WithString("end_date", mcp.Required(), mcp.Description("End date (YYYY-MM-DD)")),
		),
		h.StartSprint,
	)
	s.AddTool(
		mcp.NewTool("jira_move_to_sprint",
			mcp.WithDescription("Move issues into a sprint."),
			mcp.WithNumber("sprint_id", mcp.Required(), mcp.Description("Sprint ID")),
			mcp.WithString("issue_keys", mcp.Required(), mcp.Description("Comma-separated issue keys")),
		),
		h.MoveIssuesToSprint,
	)
	s.AddTool(
		mcp.NewTool("jira_link_issues",
			mcp.WithDescription("Create a link between two issues."),
			mcp.WithString("inward_key", mcp.Required(), mcp.Description("Inward issue key")),
			mcp.WithString("outward_key", mcp.Required(), mcp.Description("Outward issue key")),
			mcp.WithString("link_type", mcp.Required(), mcp.Description("Link type: Blocks, Relates, Duplicates, etc.")),
		),
		h.LinkIssues,
	)
	s.AddTool(
		mcp.NewTool("jira_link_types",
			mcp.WithDescription("List available issue link types."),
		),
		h.GetLinkTypes,
	)
	s.AddTool(
		mcp.NewTool("jira_add_worklog",
			mcp.WithDescription("Log time spent on an issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
			mcp.WithString("time_spent", mcp.Required(), mcp.Description("Time spent (e.g. 2h, 30m, 1d)")),
			mcp.WithString("comment", mcp.Description("Work description")),
		),
		h.AddWorklog,
	)
	s.AddTool(
		mcp.NewTool("jira_worklogs",
			mcp.WithDescription("List worklogs on a Jira issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
		),
		h.GetWorklogs,
	)
	s.AddTool(
		mcp.NewTool("jira_projects",
			mcp.WithDescription("List all accessible Jira projects."),
		),
		h.GetProjects,
	)
	s.AddTool(
		mcp.NewTool("jira_sprint_summary",
			mcp.WithDescription("Get active sprint status and issues for a board."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.GetSprintSummary,
	)

	// Epic tools
	s.AddTool(
		mcp.NewTool("jira_epic_add",
			mcp.WithDescription("Add issues to an epic by setting the parent/epic link."),
			mcp.WithString("epic_key", mcp.Required(), mcp.Description("Epic key to add issues to")),
			mcp.WithString("issue_keys", mcp.Required(), mcp.Description("Comma-separated issue keys to add")),
		),
		h.SetEpicLink,
	)
	s.AddTool(
		mcp.NewTool("jira_epic_remove",
			mcp.WithDescription("Remove issues from their epic by clearing the parent link."),
			mcp.WithString("issue_keys", mcp.Required(), mcp.Description("Comma-separated issue keys to remove")),
		),
		h.RemoveEpicLink,
	)
	s.AddTool(
		mcp.NewTool("jira_epic_issues",
			mcp.WithDescription("List all issues in an epic."),
			mcp.WithString("epic_key", mcp.Required(), mcp.Description("Epic issue key (e.g. PROJ-100)")),
			mcp.WithNumber("max_results", mcp.Description("Maximum results (default 50)")),
		),
		h.GetEpicIssues,
	)

	// Version tools
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

	// Component tools
	s.AddTool(
		mcp.NewTool("jira_components",
			mcp.WithDescription("List project components with leads."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project key")),
		),
		h.GetComponents,
	)

	// Attachment tools
	s.AddTool(
		mcp.NewTool("jira_attachments",
			mcp.WithDescription("List attachments on a Jira issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key")),
		),
		h.GetAttachments,
	)
}
