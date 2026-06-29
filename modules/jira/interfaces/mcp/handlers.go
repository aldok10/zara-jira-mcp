// Package mcp provides MCP tool handlers for the jira module.
package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/modules/jira/application/port"
	"github.com/aldok10/zara-jira-mcp/modules/jira/domain"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/validate"
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
	return mcputil.TextResult("zara-jira-mcp | status: ok | modular handlers"), nil
}

// SearchIssues searches Jira issues using JQL.
func (h *Handlers) SearchIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	jql, err := req.RequireString("jql")
	if err != nil {
		return mcputil.ErrInvalid("jql parameter is required"), nil
	}
	if err := validate.JQL(jql); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
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
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
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
// CreateIssue creates a new Jira issue.
func (h *Handlers) CreateIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	project, err := req.RequireString("project")
	if err != nil {
		return mcputil.ErrInvalid("project parameter is required"), nil
	}
	if err := validate.ProjectKey(project); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	summary, err := req.RequireString("summary")
	if err != nil {
		return mcputil.ErrInvalid("summary parameter is required"), nil
	}

	issueType := req.GetString("issue_type", "Task")
	priority := req.GetString("priority", "")
	desc := req.GetString("description", "")
	assigneeID := req.GetString("assignee_id", "")
	labelsStr := req.GetString("labels", "")

	var labels []string
	if labelsStr != "" {
		labels, err = validate.Labels(labelsStr)
		if err != nil {
			return mcputil.ErrInvalid(err.Error()), nil
		}
	}

	input := &domain.CreateIssueInput{
		Project:     project,
		Summary:     summary,
		IssueType:   issueType,
		Description: desc,
		Priority:    priority,
		Assignee:    assigneeID,
		Labels:      labels,
	}

	issue, err := h.Jira.CreateIssue(ctx, input)
	if err != nil {
		return mcputil.ErrJira("create issue", err), nil
	}

	return mcputil.TextResult(fmt.Sprintf("Created %s - %s", issue.Key, issue.Summary)), nil
}

// TransitionIssue transitions a Jira issue to a new status.
func (h *Handlers) TransitionIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	transitionID, err := req.RequireString("transition_id")
	if err != nil {
		return mcputil.ErrInvalid("transition_id parameter is required"), nil
	}
	if err := validate.TransitionID(transitionID); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	if err := h.Jira.TransitionIssue(ctx, key, transitionID); err != nil {
		return mcputil.ErrJira("transition issue", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Transitioned %s (transition: %s)", key, transitionID)), nil
}

// GetTransitions returns available transitions for an issue.
func (h *Handlers) GetTransitions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	transitions, err := h.Jira.GetTransitions(ctx, key)
	if err != nil {
		return mcputil.ErrJira("get transitions", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Available transitions for %s:\n", key))
	for _, t := range transitions {
		sb.WriteString(fmt.Sprintf("  %s: %s\n", t.ID, t.Name))
	}
	return mcputil.TextResult(sb.String()), nil
}

// AssignIssue assigns a Jira issue to a user.
func (h *Handlers) AssignIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	accountID, err := req.RequireString("account_id")
	if err != nil {
		return mcputil.ErrInvalid("account_id parameter is required"), nil
	}

	if err := h.Jira.AssignIssue(ctx, key, accountID); err != nil {
		return mcputil.ErrJira("assign issue", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Assigned %s to %s", key, accountID)), nil
}

// FindUser searches for Jira users.
func (h *Handlers) FindUser(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := req.RequireString("query")
	if err != nil {
		return mcputil.ErrInvalid("query parameter is required"), nil
	}

	users, err := h.Jira.FindUser(ctx, query)
	if err != nil {
		return mcputil.ErrJira("find user", err), nil
	}

	var sb strings.Builder
	if len(users) == 0 {
		return mcputil.TextResult("No users found."), nil
	}
	for _, u := range users {
		sb.WriteString(fmt.Sprintf("%s: %s (%s)\n", u.AccountID, u.DisplayName, u.Email))
	}
	return mcputil.TextResult(sb.String()), nil
}

// AddComment adds a comment to a Jira issue.
func (h *Handlers) AddComment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	body, err := req.RequireString("body")
	if err != nil {
		return mcputil.ErrInvalid("body parameter is required"), nil
	}

	if err := h.Jira.AddComment(ctx, key, body); err != nil {
		return mcputil.ErrJira("add comment", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Comment added to %s", key)), nil
}

// GetSprints lists sprints for a board with optional state filter.
func (h *Handlers) GetSprints(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return mcputil.ErrInvalid("board_id parameter is required"), nil
	}
	state := req.GetString("state", "")

	sprints, err := h.Jira.GetSprints(ctx, int(boardID), state)
	if err != nil {
		return mcputil.ErrJira("get sprints", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sprints for board %d (state: %s):\n", int(boardID), state))
	for _, s := range sprints {
		sb.WriteString(fmt.Sprintf("  %d: %s [%s]\n", s.ID, s.Name, s.State))
	}
	return mcputil.TextResult(sb.String()), nil
}

// StartSprint starts a sprint with start/end dates.
func (h *Handlers) StartSprint(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintID, err := req.RequireInt("sprint_id")
	if err != nil {
		return mcputil.ErrInvalid("sprint_id parameter is required"), nil
	}
	startDate := req.GetString("start_date", "")
	endDate := req.GetString("end_date", "")
	if startDate == "" || endDate == "" {
		return mcputil.ErrInvalid("start_date and end_date are required"), nil
	}

	if err := h.Jira.StartSprint(ctx, int(sprintID), startDate, endDate); err != nil {
		return mcputil.ErrJira("start sprint", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Sprint %d started (%s → %s)", int(sprintID), startDate, endDate)), nil
}

// MoveIssuesToSprint moves issues into a sprint.
func (h *Handlers) MoveIssuesToSprint(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintID, err := req.RequireInt("sprint_id")
	if err != nil {
		return mcputil.ErrInvalid("sprint_id parameter is required"), nil
	}
	issueKeysStr, err := req.RequireString("issue_keys")
	if err != nil {
		return mcputil.ErrInvalid("issue_keys parameter is required"), nil
	}
	issueKeys, err := validate.IssueKeys(issueKeysStr)
	if err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	if err := h.Jira.MoveIssuesToSprint(ctx, int(sprintID), issueKeys); err != nil {
		return mcputil.ErrJira("move to sprint", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Moved %d issues to sprint %d", len(issueKeys), int(sprintID))), nil
}

// LinkIssues creates a link between two issues.
func (h *Handlers) LinkIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	inwardKey, err := req.RequireString("inward_key")
	if err != nil {
		return mcputil.ErrInvalid("inward_key parameter is required"), nil
	}
	outwardKey, err := req.RequireString("outward_key")
	if err != nil {
		return mcputil.ErrInvalid("outward_key parameter is required"), nil
	}
	linkType, err := req.RequireString("link_type")
	if err != nil {
		return mcputil.ErrInvalid("link_type parameter is required"), nil
	}

	if err := h.Jira.LinkIssues(ctx, inwardKey, outwardKey, linkType); err != nil {
		return mcputil.ErrJira("link issues", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Linked %s → %s (%s)", inwardKey, outwardKey, linkType)), nil
}

// GetLinkTypes returns available link types.
func (h *Handlers) GetLinkTypes(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	types, err := h.Jira.GetLinkTypes(ctx)
	if err != nil {
		return mcputil.ErrJira("get link types", err), nil
	}

	var sb strings.Builder
	for _, t := range types {
		sb.WriteString(fmt.Sprintf("%s: %s / %s\n", t.Name, t.Inward, t.Outward))
	}
	return mcputil.TextResult(sb.String()), nil
}

// AddWorklog logs time on an issue.
func (h *Handlers) AddWorklog(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}
	timeSpent, err := req.RequireString("time_spent")
	if err != nil {
		return mcputil.ErrInvalid("time_spent parameter is required"), nil
	}
	comment := req.GetString("comment", "")

	if err := h.Jira.AddWorklog(ctx, key, timeSpent, comment); err != nil {
		return mcputil.ErrJira("add worklog", err), nil
	}
	return mcputil.TextResult(fmt.Sprintf("Logged %s on %s", timeSpent, key)), nil
}

// GetWorklogs returns worklogs for an issue.
func (h *Handlers) GetWorklogs(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return mcputil.ErrInvalid("key parameter is required"), nil
	}
	if err := validate.IssueKey(key); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
	}

	logs, err := h.Jira.GetWorklogs(ctx, key)
	if err != nil {
		return mcputil.ErrJira("get worklogs", err), nil
	}

	var sb strings.Builder
	if len(logs) == 0 {
		return mcputil.TextResult("No worklogs found."), nil
	}
	for _, l := range logs {
		sb.WriteString(fmt.Sprintf("%s: %s - %s\n", l.Author, l.TimeSpent, l.Comment))
	}
	return mcputil.TextResult(sb.String()), nil
}

// GetProjects lists all accessible projects.
func (h *Handlers) GetProjects(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projects, err := h.Jira.GetProjects(ctx)
	if err != nil {
		return mcputil.ErrJira("get projects", err), nil
	}

	var sb strings.Builder
	for _, p := range projects {
		sb.WriteString(fmt.Sprintf("%s: %s (%s)\n", p.Key, p.Name, p.Lead))
	}
	return mcputil.TextResult(sb.String()), nil
}

// GetSprintSummary returns the active sprint status for a board.
func (h *Handlers) GetSprintSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return mcputil.ErrInvalid("board_id parameter is required"), nil
	}
	if err := validate.BoardID(int(boardID)); err != nil {
		return mcputil.ErrInvalid(err.Error()), nil
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
