package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aldok10/zara-jira-mcp/internal/ai"
	"github.com/aldok10/zara-jira-mcp/internal/jira"
	"github.com/aldok10/zara-jira-mcp/internal/lark"
	"github.com/mark3labs/mcp-go/mcp"
)

// Handlers holds all MCP tool handler methods.
type Handlers struct {
	Jira *jira.RestClient
	AI   *ai.OpenAIClient
	Lark *lark.WebhookClient
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

// AIAnalyze uses AI to analyze Jira tickets and provide PM-relevant insights.
func (h *Handlers) AIAnalyze(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := req.RequireString("query")
	if err != nil {
		return errorResult("query parameter is required (describe what you want to know)"), nil
	}

	jql := req.GetString("jql", "resolution = Unresolved ORDER BY updated DESC")
	maxResults := req.GetInt("max_results", 30)

	result, err := h.Jira.SearchIssues(ctx, jql, maxResults)
	if err != nil {
		return errorResult("Jira search failed: " + err.Error()), nil
	}

	// Build context for AI
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

// NotifyLark sends a message to the configured Lark group.
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

// AISprintReport generates a full AI sprint report and optionally sends to Lark.
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

	systemPrompt := `You are a project management AI. Generate a concise sprint report for a PM.
Include:
1. Sprint health (on track / at risk / behind)
2. Key blockers
3. Tickets needing attention
4. Progress summary by status
5. Recommendations

Format in markdown suitable for a Lark message card. Keep it under 500 words.`

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
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: text},
		},
	}
}

func errorResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: msg},
		},
	}
}
