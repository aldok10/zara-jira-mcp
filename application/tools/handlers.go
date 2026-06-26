package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	aiprovider "github.com/aldok10/zara-jira-mcp/domain/ai"
	domain "github.com/aldok10/zara-jira-mcp/domain/jira"
	larkdom "github.com/aldok10/zara-jira-mcp/domain/lark"
	"github.com/aldok10/zara-jira-mcp/config"
	"github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/aldok10/zara-jira-mcp/internal/cache"
	icalendar "github.com/aldok10/zara-jira-mcp/internal/calendar"
	"github.com/aldok10/zara-jira-mcp/internal/clockify"
	"github.com/aldok10/zara-jira-mcp/internal/confluence"
	idiscord "github.com/aldok10/zara-jira-mcp/internal/discord"
	iemail "github.com/aldok10/zara-jira-mcp/internal/email"
	igithub "github.com/aldok10/zara-jira-mcp/internal/github"
	igitlab "github.com/aldok10/zara-jira-mcp/internal/gitlab"
	"github.com/aldok10/zara-jira-mcp/internal/linear"
	inotion "github.com/aldok10/zara-jira-mcp/internal/notion"
	"github.com/aldok10/zara-jira-mcp/internal/pagerduty"
	"github.com/aldok10/zara-jira-mcp/internal/sheets"
	"github.com/aldok10/zara-jira-mcp/internal/database"
	islack "github.com/aldok10/zara-jira-mcp/internal/slack"
	iteams "github.com/aldok10/zara-jira-mcp/internal/teams"
	itelegram "github.com/aldok10/zara-jira-mcp/internal/telegram"
	ilark "github.com/aldok10/zara-jira-mcp/internal/lark"
	"github.com/mark3labs/mcp-go/mcp"
)

// Handlers holds all MCP tool handler methods.
type Handlers struct {
	Config     *config.Config
	Jira       domain.Client
	AI         aiprovider.Provider
	Lark       larkdom.Notifier
	Slack      *islack.Client
	Discord    *idiscord.Client
	Telegram   *itelegram.Client
	Teams      *iteams.Client
	Email      *iemail.Client
	Confluence *confluence.Client
	Memory     memory.Store
	Cache      *cache.Client
	Calendar   *icalendar.Client
	GitHub     *igithub.Client
	GitLab     *igitlab.Client
	OKR        *ilark.OKRClient
	Notion     *inotion.Client
	Linear     *linear.Client
	PagerDuty  *pagerduty.Client
	Clockify   *clockify.Client
	Sheets     *sheets.Client
	Database   *database.Client
}

// SearchIssues searches Jira using JQL.
func (h *Handlers) SearchIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	jql, err := req.RequireString("jql")
	if err != nil {
		return errorResult("jql parameter is required"), nil
	}
	maxResults := req.GetInt("max_results", 20)
	startAt := req.GetInt("start_at", 0)

	result, err := h.Jira.SearchIssues(ctx, jql, maxResults, startAt)
	if err != nil {
		return sanitizedError("Jira search failed", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d issues (showing %d, offset %d, hasMore: %v):\n\n", result.Total, len(result.Issues), result.StartAt, result.HasMore))
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
		return sanitizedError("Failed to get issue", err), nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("**%s** [%s] %s\n", issue.Key, issue.Type, issue.Summary))
	sb.WriteString(fmt.Sprintf("Status: %s | Priority: %s | Assignee: %s\n", issue.Status, issue.Priority, issue.Assignee))
	if issue.Description != "" {
		sb.WriteString(fmt.Sprintf("\nDescription:\n%s\n", issue.Description))
	}
	data, _ := json.MarshalIndent(issue, "", "  ")
	sb.WriteString(fmt.Sprintf("\n--- Raw Data ---\n%s\n", string(data)))
	return textResult(sb.String()), nil
}

// GetBoards lists all accessible Jira boards.
func (h *Handlers) GetBoards(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	const cacheKey = "jira:boards"
	if h.Cache.Available() {
		if cached, err := h.Cache.Get(ctx, cacheKey); err == nil {
			return textResult(cached), nil
		}
	}

	boards, err := h.Jira.GetBoards(ctx)
	if err != nil {
		return sanitizedError("Failed to get boards", err), nil
	}
	var sb strings.Builder
	for _, b := range boards {
		sb.WriteString(fmt.Sprintf("- [%d] %s (%s)\n", b.ID, b.Name, b.Type))
	}
	result := sb.String()
	_ = h.Cache.Set(ctx, cacheKey, result, 10*time.Minute)
	return textResult(result), nil
}

// GetSprintSummary gets active sprint issues and generates a summary.
func (h *Handlers) GetSprintSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id parameter is required"), nil
	}

	cacheKey := fmt.Sprintf("jira:sprint_summary:%d", boardID)
	if h.Cache.Available() {
		if cached, err := h.Cache.Get(ctx, cacheKey); err == nil {
			return textResult(cached), nil
		}
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil {
		return sanitizedError("Failed to get sprints", err), nil
	}
	if len(sprints) == 0 {
		return textResult("No active sprints found for this board."), nil
	}

	sprint := sprints[0]
	issues, err := h.Jira.GetSprintIssues(ctx, sprint.ID)
	if err != nil {
		return sanitizedError("Failed to get sprint issues", err), nil
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
	result := sb.String()
	_ = h.Cache.Set(ctx, cacheKey, result, 2*time.Minute)
	return textResult(result), nil
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
		return sanitizedError("Failed to create issue", err), nil
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
		return sanitizedError("Failed to add comment", err), nil
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
		return sanitizedError("Failed to transition issue", err), nil
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
		return sanitizedError("Failed to get transitions", err), nil
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

	cacheKey := "jira:my_issues:" + status
	if h.Cache.Available() {
		if cached, err := h.Cache.Get(ctx, cacheKey); err == nil {
			return textResult(cached), nil
		}
	}

	jql := "assignee = currentUser() AND resolution = Unresolved ORDER BY updated DESC"
	if status != "" {
		jql = fmt.Sprintf("assignee = currentUser() AND resolution = Unresolved AND status = \"%s\" ORDER BY updated DESC", status)
	}

	result, err := h.Jira.SearchIssues(ctx, jql, 30, 0)
	if err != nil {
		return sanitizedError("Failed to get my issues", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Your open issues (%d):\n\n", len(result.Issues)))
	for _, issue := range result.Issues {
		sb.WriteString(fmt.Sprintf("- **%s** [%s] %s | Priority: %s\n", issue.Key, issue.Status, issue.Summary, issue.Priority))
	}
	out := sb.String()
	_ = h.Cache.Set(ctx, cacheKey, out, 1*time.Minute)
	return textResult(out), nil
}

// Overdue shows issues that might be overdue or stale.
func (h *Handlers) Overdue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	days := req.GetInt("days", 14)
	project := req.GetString("project", "")

	jql := fmt.Sprintf("resolution = Unresolved AND updated <= -%dd ORDER BY updated ASC", days)
	if project != "" {
		jql = fmt.Sprintf("project = %s AND resolution = Unresolved AND updated <= -%dd ORDER BY updated ASC", project, days)
	}

	result, err := h.Jira.SearchIssues(ctx, jql, 30, 0)
	if err != nil {
		return sanitizedError("Failed to get overdue issues", err), nil
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

	cacheKey := "jira:workload:" + project
	if h.Cache.Available() {
		if cached, err := h.Cache.Get(ctx, cacheKey); err == nil {
			return textResult(cached), nil
		}
	}

	jql := "resolution = Unresolved AND assignee IS NOT EMPTY ORDER BY assignee ASC"
	if project != "" {
		jql = fmt.Sprintf("project = %s AND resolution = Unresolved AND assignee IS NOT EMPTY ORDER BY assignee ASC", project)
	}

	result, err := h.Jira.SearchIssues(ctx, jql, 200, 0)
	if err != nil {
		return sanitizedError("Failed to get workload", err), nil
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
	out := sb.String()
	_ = h.Cache.Set(ctx, cacheKey, out, 3*time.Minute)
	return textResult(out), nil
}

// AIAnalyze uses AI to analyze Jira tickets.
func (h *Handlers) AIAnalyze(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := req.RequireString("query")
	if err != nil {
		return errorResult("query parameter is required"), nil
	}

	jql := req.GetString("jql", "resolution = Unresolved ORDER BY updated DESC")
	maxResults := req.GetInt("max_results", 30)

	result, err := h.Jira.SearchIssues(ctx, jql, maxResults, 0)
	if err != nil {
		return sanitizedError("Jira search failed", err), nil
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

	analysis, err := h.aiComplete(ctx, systemPrompt, userPrompt)
	if err != nil {
		return sanitizedError("AI analysis failed", err), nil
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
		return sanitizedError("Failed to send to Lark", err), nil
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
		return sanitizedError("Failed to get sprints", err), nil
	}
	if len(sprints) == 0 {
		return textResult("No active sprints found."), nil
	}

	sprint := sprints[0]
	issues, err := h.Jira.GetSprintIssues(ctx, sprint.ID)
	if err != nil {
		return sanitizedError("Failed to get sprint issues", err), nil
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

	report, err := h.aiComplete(ctx, systemPrompt, userPrompt)
	if err != nil {
		return sanitizedError("AI report generation failed", err), nil
	}

	if sendToLark {
		title := fmt.Sprintf("Sprint Report: %s", sprint.Name)
		if err := h.Lark.SendMarkdown(ctx, title, report); err != nil {
			slog.Warn("Sprint report Lark send failed", "detail", err.Error())
			return textResult(report + "\n\n(Warning: failed to send to Lark — check server logs)"), nil
		}
		return textResult(report + "\n\n(Sent to Lark successfully)"), nil
	}
	return textResult(report), nil
}

// UpdateIssue updates an existing Jira issue.
func (h *Handlers) UpdateIssue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key parameter is required"), nil
	}

	input := &domain.UpdateIssueInput{Key: key}
	input.Summary = req.GetString("summary", "")
	input.Description = req.GetString("description", "")
	input.Priority = req.GetString("priority", "")
	input.Assignee = req.GetString("assignee_id", "")

	labelsRaw := req.GetString("labels", "")
	if labelsRaw != "" {
		input.Labels = strings.Split(labelsRaw, ",")
	}

	if err := h.Jira.UpdateIssue(ctx, input); err != nil {
		return sanitizedError("Failed to update issue", err), nil
	}
	return textResult(fmt.Sprintf("Issue %s updated successfully", key)), nil
}

// Health returns server version and status.
func (h *Handlers) Health(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return textResult("zara-jira-mcp v0.3.0 | status: ok | tools: 139"), nil
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

// sanitizedError logs the detailed error and returns a generic message to the client.
// Use for all external-facing calls where err.Error() might leak internals.
func sanitizedError(logMsg string, err error) *mcp.CallToolResult {
	slog.Error(logMsg, "detail", err.Error())
	return errorResult("Operation failed. Contact your administrator.")
}
