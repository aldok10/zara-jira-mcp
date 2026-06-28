// Package mcp provides MCP tool handlers for the jira module.
package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/modules/jira/application/port"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)

// Handlers holds dependencies for jira MCP tool handlers.
type Handlers struct {
	Jira port.Inbound
}

// NewHandlers creates a new jira MCP handlers instance.
func NewHandlers(jiraService port.Inbound) *Handlers {
	return &Handlers{Jira: jiraService}
}

// Health returns server version and status.
func (h *Handlers) Health(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("zara-jira-mcp v0.4.0 | status: ok | modular handlers"), nil
}

// SearchIssues searches Jira issues using JQL.
func (h *Handlers) SearchIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	jql, err := req.RequireString("jql")
	if err != nil {
		return mcputil.ErrInvalid("jql parameter is required"), nil
	}
	maxResults := req.GetInt("max_results", 20)

	results, err := h.Jira.SearchIssues(ctx, jql, int(maxResults))
	if err != nil {
		return mcputil.ErrJira("Jira search", err), nil
	}
	if results == nil {
		return mcputil.TextResult("No results found."), nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d issues\n", len(results.Issues)))
	for _, issue := range results.Issues {
		sb.WriteString(fmt.Sprintf("  %s - %s [%s]\n", issue.Key, issue.Summary, issue.Status))
	}
	return mcputil.TextResult(sb.String()), nil
}

// GetIssue returns full details of a single Jira issue.
func (h *Handlers) GetIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}

	issue, err := h.Jira.GetIssue(ctx, key)
	if err != nil {
		return mcputil.ErrJira("get issue", err), nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**%s** - %s\n", issue.Key, issue.Summary))
	sb.WriteString(fmt.Sprintf("Type: %s | Status: %s | Priority: %s\n", issue.Type, issue.Status, issue.Priority))
	sb.WriteString(fmt.Sprintf("Assignee: %s\n", issue.Assignee))
	sb.WriteString(fmt.Sprintf("Description: %s\n", issue.Description))
	return mcputil.TextResult(sb.String()), nil
}

// GetBoards lists all accessible Jira boards.
func (h *Handlers) GetBoards(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boards, err := h.Jira.GetBoards(ctx)
	if err != nil {
		return mcputil.ErrJira("get boards", err), nil
	}

	var sb strings.Builder
	for _, b := range boards {
		sb.WriteString(fmt.Sprintf("%d: %s (%s)\n", b.ID, b.Name, b.Type))
	}
	return mcputil.TextResult(sb.String()), nil
}

// GetSprintSummary returns the active sprint status for a board.
func (h *Handlers) GetSprintSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return mcputil.ErrInvalid("board_id parameter is required"), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, int(boardID))
	if err != nil {
		return mcputil.ErrJira("get active sprints", err), nil
	}

	if len(sprints) == 0 {
		return mcputil.TextResult("No active sprints found for this board."), nil
	}

	sprint := sprints[0]
	issues, err := h.Jira.GetSprintIssues(ctx, sprint.ID)
	if err != nil {
		return mcputil.ErrJira("get sprint issues", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sprint: %s (ID: %d)\n", sprint.Name, sprint.ID))
	sb.WriteString(fmt.Sprintf("Goals: %s\n", sprint.Goal))
	sb.WriteString(fmt.Sprintf("Start: %s | End: %s\n", sprint.StartDate, sprint.EndDate))
	sb.WriteString(fmt.Sprintf("Issues: %d\n", len(issues)))

	for _, i := range issues {
		sb.WriteString(fmt.Sprintf("  %s - %s [%s]\n", i.Key, i.Summary, i.Status))
	}
	return mcputil.TextResult(sb.String()), nil
}

// CreateIssue creates a new Jira issue.
func (h *Handlers) CreateIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Create issue - NOT YET IMPLEMENTED"), nil
}

// CreateSubtask creates a subtask.
func (h *Handlers) CreateSubtask(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Create subtask - NOT YET IMPLEMENTED"), nil
}

// UpdateIssue updates an existing Jira issue.
func (h *Handlers) UpdateIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Update issue - NOT YET IMPLEMENTED"), nil
}

// DeleteIssue deletes a Jira issue.
func (h *Handlers) DeleteIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	return mcputil.TextResult("Deleted: " + key), nil
}

// TransitionIssue transitions a Jira issue through workflow.
func (h *Handlers) TransitionIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Transition issue - NOT YET IMPLEMENTED"), nil
}

// GetTransitions lists available transitions for a Jira issue.
func (h *Handlers) GetTransitions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Get transitions - NOT YET IMPLEMENTED"), nil
}

// AddComment adds a comment to a Jira issue.
func (h *Handlers) AddComment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Add comment - NOT YET IMPLEMENTED"), nil
}

// AssignIssue assigns a Jira issue to a user.
func (h *Handlers) AssignIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Assign issue - NOT YET IMPLEMENTED"), nil
}

// UnassignIssue unassigns a Jira issue.
func (h *Handlers) UnassignIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Unassign issue - NOT YET IMPLEMENTED"), nil
}

// FindUser searches for a Jira user by query.
func (h *Handlers) FindUser(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Find user - NOT YET IMPLEMENTED"), nil
}

// ListProjects lists all accessible Jira projects.
func (h *Handlers) ListProjects(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("List projects - NOT YET IMPLEMENTED"), nil
}

// MyIssues lists issues assigned to the current user.
func (h *Handlers) MyIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("My issues - NOT YET IMPLEMENTED"), nil
}

// Workload shows workload distribution across assignees.
func (h *Handlers) Workload(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Workload - NOT YET IMPLEMENTED"), nil
}

// Overdue lists overdue issues for a board.
func (h *Handlers) Overdue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Overdue - NOT YET IMPLEMENTED"), nil
}

// EpicIssues lists issues belonging to an epic.
func (h *Handlers) EpicIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Epic issues - NOT YET IMPLEMENTED"), nil
}

// EpicAdd adds issues to an epic.
func (h *Handlers) EpicAdd(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Epic add - NOT YET IMPLEMENTED"), nil
}

// EpicRemove removes issues from an epic.
func (h *Handlers) EpicRemove(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Epic remove - NOT YET IMPLEMENTED"), nil
}

// ListSprints lists sprints for a board.
func (h *Handlers) ListSprints(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("List sprints - NOT YET IMPLEMENTED"), nil
}

// CreateSprint creates a new sprint on a board.
func (h *Handlers) CreateSprint(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Create sprint - NOT YET IMPLEMENTED"), nil
}

// StartSprint starts (activates) a sprint.
func (h *Handlers) StartSprint(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Start sprint - NOT YET IMPLEMENTED"), nil
}

// CloseSprint closes (completes) a sprint.
func (h *Handlers) CloseSprint(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Close sprint - NOT YET IMPLEMENTED"), nil
}

// MoveIssuesToSprint moves issues to a sprint.
func (h *Handlers) MoveIssuesToSprint(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Move to sprint - NOT YET IMPLEMENTED"), nil
}

// StoryPointsSummary shows story points summary for a board.
func (h *Handlers) StoryPointsSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Story points - NOT YET IMPLEMENTED"), nil
}

// SprintPointsBurndown shows sprint points burndown for a board's active sprint.
func (h *Handlers) SprintPointsBurndown(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Burndown - NOT YET IMPLEMENTED"), nil
}

// RightSize suggests right-sized estimates for issues on a board.
func (h *Handlers) RightSize(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Right size - NOT YET IMPLEMENTED"), nil
}

// GetAttachments lists attachments on a Jira issue.
func (h *Handlers) GetAttachments(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Get attachments - NOT YET IMPLEMENTED"), nil
}

// GetVersions lists versions (releases) for a project.
func (h *Handlers) GetVersions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Get versions - NOT YET IMPLEMENTED"), nil
}

// CreateVersion creates a version (release) in a project.
func (h *Handlers) CreateVersion(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Create version - NOT YET IMPLEMENTED"), nil
}

// ReleaseVersion releases a version (marks as released) in a project.
func (h *Handlers) ReleaseVersion(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Release version - NOT YET IMPLEMENTED"), nil
}

// GetComponents lists components in a project.
func (h *Handlers) GetComponents(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Get components - NOT YET IMPLEMENTED"), nil
}

// GetFields lists all custom fields available in Jira.
func (h *Handlers) GetFields(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Get fields - NOT YET IMPLEMENTED"), nil
}

// LinkIssues links two Jira issues together.
func (h *Handlers) LinkIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Link issues - NOT YET IMPLEMENTED"), nil
}

// LinkTypes lists all available issue link types.
func (h *Handlers) LinkTypes(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Link types - NOT YET IMPLEMENTED"), nil
}

// WorklogAdd adds a worklog entry to a Jira issue.
func (h *Handlers) WorklogAdd(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Worklog add - NOT YET IMPLEMENTED"), nil
}

// WorklogList lists worklog entries on a Jira issue.
func (h *Handlers) WorklogList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Worklog list - NOT YET IMPLEMENTED"), nil
}

// Watch starts watching a Jira issue.
func (h *Handlers) Watch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Watch - NOT YET IMPLEMENTED"), nil
}

// Watchers gets watchers of a Jira issue.
func (h *Handlers) Watchers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Watchers - NOT YET IMPLEMENTED"), nil
}

// LabelsSet sets labels on a Jira issue (replaces all existing labels).
func (h *Handlers) LabelsSet(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Labels set - NOT YET IMPLEMENTED"), nil
}

// BulkTransition bulk transitions multiple issues to a new status.
func (h *Handlers) BulkTransition(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Bulk transition - NOT YET IMPLEMENTED"), nil
}

// BulkAssign bulk assigns multiple issues to a user.
func (h *Handlers) BulkAssign(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Bulk assign - NOT YET IMPLEMENTED"), nil
}

// BulkLabel bulk adds labels to multiple issues.
func (h *Handlers) BulkLabel(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Bulk label - NOT YET IMPLEMENTED"), nil
}

// ProjectDetail gets full project details by key.
func (h *Handlers) ProjectDetail(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Project detail - NOT YET IMPLEMENTED"), nil
}

// TraceTicketBranch traces a Jira issue back to its development branch.
func (h *Handlers) TraceTicketBranch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return mcputil.TextResult("Trace branch - NOT YET IMPLEMENTED"), nil
}

