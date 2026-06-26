package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/aldok10/zara-jira-mcp/internal/ai"
	"github.com/aldok10/zara-jira-mcp/internal/jira"
	"github.com/aldok10/zara-jira-mcp/internal/lark"
	domain "github.com/aldok10/zara-jira-mcp/domain/jira"
	"github.com/mark3labs/mcp-go/mcp"
)

// Handlers holds all MCP tool handler methods.
type Handlers struct {
	Jira   *jira.RestClient
	AI     *ai.OpenAIClient
	Lark   *lark.WebhookClient
	Memory memory.Store
}

// SearchIssues searches Jira using JQL.
func (h *Handlers) SearchIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	jql, err := req.RequireString("jql")
	if err != nil {
		return errorResult("jql parameter is required"), nil
	}
	maxResults := req.GetInt("max_results", 20)

	result, err := h.Jira.SearchIssues(ctx, jql, maxResults)
	if err != nil {
		return errorResult("Jira search failed: " + err.Error()), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d issues (showing %d):\n\n", result.Total, len(result.Issues)))
	for _, issue := range result.Issues {
		sb.WriteString(fmt.Sprintf("**%s** [%s] %s\n  Status: %s | Priority: %s | Assignee: %s\n\n",
			issue.Key, issue.Type, issue.Summary, issue.Status, issue.Priority, issue.Assignee))
	}
	return textResult(sb.String()), nil
}

// GetIssue retrieves a single Jira issue by key.
func (h *Handlers) GetIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key parameter is required"), nil
	}
	issue, err := h.Jira.GetIssue(ctx, key)
	if err != nil {
		return errorResult("Failed to get issue: " + err.Error()), nil
	}
	data, _ := json.MarshalIndent(issue, "", "  ")
	return textResult(string(data)), nil
}

// GetBoards lists all accessible Jira boards.
func (h *Handlers) GetBoards(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boards, err := h.Jira.GetBoards(ctx)
	if err != nil {
		return errorResult("Failed to get boards: " + err.Error()), nil
	}
	var sb strings.Builder
	for _, b := range boards {
		sb.WriteString(fmt.Sprintf("- [%d] %s (%s)\n", b.ID, b.Name, b.Type))
	}
	return textResult(sb.String()), nil
}

// GetSprintSummary gets active sprint issues and generates a summary.
func (h *Handlers) GetSprintSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id parameter is required"), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil {
		return errorResult("Failed to get sprints: " + err.Error()), nil
	}
	if len(sprints) == 0 {
		return textResult("No active sprints found for this board."), nil
	}

	sprint := sprints[0]
	issues, err := h.Jira.GetSprintIssues(ctx, sprint.ID)
	if err != nil {
		return errorResult("Failed to get sprint issues: " + err.Error()), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sprint: %s (Goal: %s)\n\n", sprint.Name, sprint.Goal))

	statusCount := map[string]int{}
	for _, issue := range issues {
		statusCount[issue.Status]++
	}
	sb.WriteString("Status breakdown:\n")
	for status, count := range statusCount {
		sb.WriteString(fmt.Sprintf("  %s: %d\n", status, count))
	}
	sb.WriteString(fmt.Sprintf("\nTotal: %d issues\n\n", len(issues)))

	for _, issue := range issues {
		sb.WriteString(fmt.Sprintf("- %s [%s] %s (Assignee: %s)\n", issue.Key, issue.Status, issue.Summary, issue.Assignee))
	}
	return textResult(sb.String()), nil
}

// CreateIssue creates a new Jira issue.
func (h *Handlers) CreateIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	project, err := req.RequireString("project")
	if err != nil {
		return errorResult("project parameter is required"), nil
	}
	summary, err := req.RequireString("summary")
	if err != nil {
		return errorResult("summary parameter is required"), nil
	}
	issueType := req.GetString("issue_type", "Task")

	input := &domain.CreateIssueInput{
		Project:     project,
		Summary:     summary,
		IssueType:   issueType,
		Description: req.GetString("description", ""),
		Priority:    req.GetString("priority", ""),
		Assignee:    req.GetString("assignee_id", ""),
	}

	labelsRaw := req.GetString("labels", "")
	if labelsRaw != "" {
		input.Labels = strings.Split(labelsRaw, ",")
	}

	created, err := h.Jira.CreateIssue(ctx, input)
	if err != nil {
		return errorResult("Failed to create issue: " + err.Error()), nil
	}
	return textResult(fmt.Sprintf("Created: %s - %s", created.Key, created.Summary)), nil
}

// AddComment adds a comment to a Jira issue.
func (h *Handlers) AddComment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key parameter is required"), nil
	}
	body, err := req.RequireString("body")
	if err != nil {
		return errorResult("body parameter is required"), nil
	}

	if err := h.Jira.AddComment(ctx, key, body); err != nil {
		return errorResult("Failed to add comment: " + err.Error()), nil
	}
	return textResult(fmt.Sprintf("Comment added to %s", key)), nil
}

// TransitionIssue transitions an issue to a new status.
func (h *Handlers) TransitionIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key parameter is required"), nil
	}
	transitionID, err := req.RequireString("transition_id")
	if err != nil {
		return errorResult("transition_id parameter is required"), nil
	}

	if err := h.Jira.TransitionIssue(ctx, key, transitionID); err != nil {
		return errorResult("Failed to transition issue: " + err.Error()), nil
	}
	return textResult(fmt.Sprintf("Issue %s transitioned successfully", key)), nil
}

// GetTransitions lists available transitions for an issue.
func (h *Handlers) GetTransitions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key parameter is required"), nil
	}

	transitions, err := h.Jira.GetTransitions(ctx, key)
	if err != nil {
		return errorResult("Failed to get transitions: " + err.Error()), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Available transitions for %s:\n\n", key))
	for _, t := range transitions {
		sb.WriteString(fmt.Sprintf("- [%s] %s\n", t.ID, t.Name))
	}
	return textResult(sb.String()), nil
}

// MyIssues shows issues assigned to the current user.
func (h *Handlers) MyIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	status := req.GetString("status", "")
	jql := "assignee = currentUser() AND resolution = Unresolved ORDER BY updated DESC"
	if status != "" {
		jql = fmt.Sprintf("assignee = currentUser() AND resolution = Unresolved AND status = \"%s\" ORDER BY updated DESC", status)
	}

	result, err := h.Jira.SearchIssues(ctx, jql, 30)
	if err != nil {
		return errorResult("Failed to get my issues: " + err.Error()), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Your open issues (%d):\n\n", len(result.Issues)))
	for _, issue := range result.Issues {
		sb.WriteString(fmt.Sprintf("- **%s** [%s] %s | Priority: %s\n", issue.Key, issue.Status, issue.Summary, issue.Priority))
	}
	return textResult(sb.String()), nil
}

// Overdue shows issues that might be overdue or stale.
func (h *Handlers) Overdue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	days := req.GetInt("days", 14)
	project := req.GetString("project", "")

	jql := fmt.Sprintf("resolution = Unresolved AND updated <= -%dd ORDER BY updated ASC", days)
	if project != "" {
		jql = fmt.Sprintf("project = %s AND resolution = Unresolved AND updated <= -%dd ORDER BY updated ASC", project, days)
	}

	result, err := h.Jira.SearchIssues(ctx, jql, 30)
	if err != nil {
		return errorResult("Failed to get overdue issues: " + err.Error()), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Stale issues (no update in %d+ days): %d\n\n", days, len(result.Issues)))
	for _, issue := range result.Issues {
		sb.WriteString(fmt.Sprintf("- **%s** [%s] %s | Assignee: %s | Last updated: %s\n",
			issue.Key, issue.Status, issue.Summary, issue.Assignee, issue.Updated.Format("2006-01-02")))
	}
	return textResult(sb.String()), nil
}

// Workload shows workload distribution across team members.
func (h *Handlers) Workload(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	project := req.GetString("project", "")

	jql := "resolution = Unresolved AND assignee IS NOT EMPTY ORDER BY assignee ASC"
	if project != "" {
		jql = fmt.Sprintf("project = %s AND resolution = Unresolved AND assignee IS NOT EMPTY ORDER BY assignee ASC", project)
	}

	result, err := h.Jira.SearchIssues(ctx, jql, 200)
	if err != nil {
		return errorResult("Failed to get workload: " + err.Error()), nil
	}

	workload := map[string]int{}
	for _, issue := range result.Issues {
		if issue.Assignee != "" {
			workload[issue.Assignee]++
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Workload distribution (%d open issues):\n\n", len(result.Issues)))
	for person, count := range workload {
		sb.WriteString(fmt.Sprintf("- %s: %d issues\n", person, count))
	}
	return textResult(sb.String()), nil
}

// AIAnalyze uses AI to analyze Jira tickets.
func (h *Handlers) AIAnalyze(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := req.RequireString("query")
	if err != nil {
		return errorResult("query parameter is required"), nil
	}

	jql := req.GetString("jql", "resolution = Unresolved ORDER BY updated DESC")
	maxResults := req.GetInt("max_results", 30)

	result, err := h.Jira.SearchIssues(ctx, jql, maxResults)
	if err != nil {
		return errorResult("Jira search failed: " + err.Error()), nil
	}

	var issueContext strings.Builder
	for _, issue := range result.Issues {
		issueContext.WriteString(fmt.Sprintf("[%s] %s | Type: %s | Status: %s | Priority: %s | Assignee: %s | Labels: %s\n",
			issue.Key, issue.Summary, issue.Type, issue.Status, issue.Priority, issue.Assignee, strings.Join(issue.Labels, ",")))
		if issue.Description != "" {
			desc := issue.Description
			if len(desc) > 200 {
				desc = desc[:200] + "..."
			}
			issueContext.WriteString("  Description: " + desc + "\n")
		}
	}

	systemPrompt := `You are an AI assistant helping a Project Manager understand their Jira board.
Provide clear, actionable insights. Focus on:
- Blockers and risks
- Progress and velocity
- Patterns (bottlenecks, unassigned work, stale tickets)
- Recommendations for the PM

Be concise and data-driven. Reference specific ticket keys when relevant.`

	userPrompt := fmt.Sprintf("PM's question: %s\n\nJira tickets (%d total, showing %d):\n%s",
		query, result.Total, len(result.Issues), issueContext.String())

	analysis, err := h.AI.Complete(ctx, systemPrompt, userPrompt)
	if err != nil {
		return errorResult("AI analysis failed: " + err.Error()), nil
	}
	return textResult(analysis), nil
}

// NotifyLark sends a message to Lark.
func (h *Handlers) NotifyLark(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	content, err := req.RequireString("content")
	if err != nil {
		return errorResult("content parameter is required"), nil
	}
	title := req.GetString("title", "Jira Update")

	if err := h.Lark.SendMarkdown(ctx, title, content); err != nil {
		return errorResult("Failed to send to Lark: " + err.Error()), nil
	}
	return textResult("Message sent to Lark successfully."), nil
}

// AISprintReport generates an AI sprint report.
func (h *Handlers) AISprintReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id parameter is required"), nil
	}
	sendToLark := req.GetBool("send_to_lark", false)

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil {
		return errorResult("Failed to get sprints: " + err.Error()), nil
	}
	if len(sprints) == 0 {
		return textResult("No active sprints found."), nil
	}

	sprint := sprints[0]
	issues, err := h.Jira.GetSprintIssues(ctx, sprint.ID)
	if err != nil {
		return errorResult("Failed to get sprint issues: " + err.Error()), nil
	}

	var issueContext strings.Builder
	for _, issue := range issues {
		issueContext.WriteString(fmt.Sprintf("[%s] %s | Status: %s | Priority: %s | Assignee: %s\n",
			issue.Key, issue.Summary, issue.Status, issue.Priority, issue.Assignee))
	}

	systemPrompt := `You are a project management AI. Generate a concise sprint report.
Include:
1. Sprint health (on track / at risk / behind)
2. Key blockers
3. Tickets needing attention
4. Progress summary by status
5. Recommendations

Format in markdown. Keep it under 500 words.`

	userPrompt := fmt.Sprintf("Sprint: %s\nGoal: %s\nTotal issues: %d\n\n%s",
		sprint.Name, sprint.Goal, len(issues), issueContext.String())

	report, err := h.AI.Complete(ctx, systemPrompt, userPrompt)
	if err != nil {
		return errorResult("AI report generation failed: " + err.Error()), nil
	}

	if sendToLark {
		title := fmt.Sprintf("Sprint Report: %s", sprint.Name)
		if err := h.Lark.SendMarkdown(ctx, title, report); err != nil {
			return textResult(report + "\n\n(Warning: failed to send to Lark: " + err.Error() + ")"), nil
		}
		return textResult(report + "\n\n(Sent to Lark successfully)"), nil
	}
	return textResult(report), nil
}

func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.TextContent{Type: "text", Text: text}},
	}
}

func errorResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{mcp.TextContent{Type: "text", Text: msg}},
	}
}
