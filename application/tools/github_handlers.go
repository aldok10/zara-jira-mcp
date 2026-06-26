package tools

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// LinkPR adds a comment to a Jira issue linking a PR/commit URL.
func (h *Handlers) LinkPR(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key required"), nil
	}
	prURL, err := req.RequireString("pr_url")
	if err != nil {
		return errorResult("pr_url required"), nil
	}

	title := req.GetString("title", "Pull Request")
	comment := fmt.Sprintf("PR linked: [%s](%s)", title, prURL)

	if err := h.Jira.AddComment(ctx, key, comment); err != nil {
		return errorResult("Failed to add PR link: " + err.Error()), nil
	}
	return textResult(fmt.Sprintf("Linked PR to %s: %s", key, prURL)), nil
}

// IssueFromBranch extracts a Jira issue key from a git branch name and fetches issue details.
func (h *Handlers) IssueFromBranch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	branch, err := req.RequireString("branch")
	if err != nil {
		return errorResult("branch required"), nil
	}

	re := regexp.MustCompile(`[A-Z][A-Z0-9]+-\d+`)
	match := re.FindString(strings.ToUpper(branch))
	if match == "" {
		return textResult("No Jira issue key found in branch: " + branch), nil
	}

	issue, err := h.Jira.GetIssue(ctx, match)
	if err != nil {
		return textResult(fmt.Sprintf("Found key %s in branch but failed to fetch: %s", match, err.Error())), nil
	}

	return textResult(fmt.Sprintf("Issue: %s\nSummary: %s\nStatus: %s\nAssignee: %s",
		issue.Key, issue.Summary, issue.Status, issue.Assignee)), nil
}

// SmartCommit parses a commit message and applies Jira smart commit actions.
func (h *Handlers) SmartCommit(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	message, err := req.RequireString("message")
	if err != nil {
		return errorResult("message required"), nil
	}

	re := regexp.MustCompile(`[A-Z][A-Z0-9]+-\d+`)
	key := re.FindString(message)
	if key == "" {
		return textResult("No Jira issue key found in commit message."), nil
	}

	var actions []string

	// Check for #done / #close / #resolve
	if strings.Contains(message, "#done") || strings.Contains(message, "#close") || strings.Contains(message, "#resolve") {
		transitions, _ := h.Jira.GetTransitions(ctx, key)
		for _, t := range transitions {
			lower := strings.ToLower(t.Name)
			if strings.Contains(lower, "done") || strings.Contains(lower, "close") || strings.Contains(lower, "resolve") {
				_ = h.Jira.TransitionIssue(ctx, key, t.ID)
				actions = append(actions, "Transitioned to "+t.Name)
				break
			}
		}
	}

	// Check for #time Xh
	timeRe := regexp.MustCompile(`#time\s+(\d+[hmd])`)
	if timeMatch := timeRe.FindStringSubmatch(message); len(timeMatch) > 1 {
		_ = h.Jira.AddWorklog(ctx, key, timeMatch[1], "")
		actions = append(actions, "Logged "+timeMatch[1])
	}

	// Check for #comment text
	commentRe := regexp.MustCompile(`#comment\s+(.+?)(?:\s+#|$)`)
	if commentMatch := commentRe.FindStringSubmatch(message); len(commentMatch) > 1 {
		_ = h.Jira.AddComment(ctx, key, commentMatch[1])
		actions = append(actions, "Comment added")
	}

	if len(actions) == 0 {
		return textResult(fmt.Sprintf("Found issue %s but no actions (#done, #time, #comment) in message.", key)), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Smart commit applied to %s:\n", key))
	for _, a := range actions {
		sb.WriteString(fmt.Sprintf("  - %s\n", a))
	}
	return textResult(sb.String()), nil
}
