package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// TraceTicketBranch checks if a Jira ticket has corresponding branches in GitHub/GitLab
// and whether those branches have been merged.
func (h *Handlers) TraceTicketBranch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key parameter required (e.g. SIT-3658)"), nil
	}

	// Common branch naming patterns for a ticket
	patterns := generateBranchPatterns(key)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Branch trace for %s:\n\n", key))

	found := false

	// Search GitHub
	if h.GitHub != nil && h.GitHub.Available() {
		sb.WriteString("**GitHub:**\n")
		for _, pattern := range patterns {
			branches, err := h.GitHub.SearchBranches(ctx, pattern)
			if err != nil {
				sb.WriteString(fmt.Sprintf("  (error: %s)\n", err.Error()))
				break
			}
			for _, b := range branches {
				found = true
				sb.WriteString(fmt.Sprintf("  Branch: `%s`\n", b.Name))
				// Check PRs for this branch
				prs, err := h.GitHub.SearchPRsByBranch(ctx, b.Name)
				if err == nil && len(prs) > 0 {
					for _, pr := range prs {
						var icon string
						if pr.State == "merged" {
							icon = "✅"
						} else if pr.State == "closed" {
							icon = "⚫"
						} else {
							icon = "🟡"
						}
						mergedInfo := ""
						if pr.MergedTo != "" {
							mergedInfo = fmt.Sprintf(" → merged to `%s`", pr.MergedTo)
						}
						sb.WriteString(fmt.Sprintf("    %s PR #%d [%s]: %s%s\n", icon, pr.Number, pr.State, pr.Title, mergedInfo))
					}
				} else {
					sb.WriteString("    No PRs found for this branch\n")
				}
			}
		}
		if !found {
			sb.WriteString("  No matching branches found\n")
		}
		sb.WriteString("\n")
	}

	// Search GitLab
	gitlabFound := false
	if h.GitLab != nil && h.GitLab.Available() {
		sb.WriteString("**GitLab:**\n")
		for _, pattern := range patterns {
			branches, err := h.GitLab.SearchBranches(ctx, pattern)
			if err != nil {
				sb.WriteString(fmt.Sprintf("  (error: %s)\n", err.Error()))
				break
			}
			for _, b := range branches {
				gitlabFound = true
				found = true
				mergeStatus := ""
				if b.Merged {
					mergeStatus = " (merged into default branch)"
				}
				sb.WriteString(fmt.Sprintf("  Branch: `%s`%s\n", b.Name, mergeStatus))
				// Check MRs for this branch
				mrs, err := h.GitLab.SearchMRsByBranch(ctx, b.Name)
				if err == nil && len(mrs) > 0 {
					for _, mr := range mrs {
						var icon string
						if mr.State == "merged" {
							icon = "✅"
						} else if mr.State == "closed" {
							icon = "⚫"
						} else {
							icon = "🟡"
						}
						sb.WriteString(fmt.Sprintf("    %s MR !%d [%s]: %s → `%s`\n", icon, mr.IID, mr.State, mr.Title, mr.TargetBranch))
					}
				} else {
					sb.WriteString("    No MRs found for this branch\n")
				}
			}
		}
		if !gitlabFound {
			sb.WriteString("  No matching branches found\n")
		}
		sb.WriteString("\n")
	}

	if h.GitHub == nil && h.GitLab == nil {
		return errorResult("Neither GitHub nor GitLab configured. Set GITHUB_TOKEN or GITLAB_TOKEN."), nil
	}

	if !found {
		sb.WriteString(fmt.Sprintf("\nNo branches found matching ticket %s.\n", key))
		sb.WriteString("Expected patterns: " + strings.Join(patterns, ", "))
	}

	return textResult(sb.String()), nil
}

// generateBranchPatterns creates common branch name patterns from a Jira key.
func generateBranchPatterns(key string) []string {
	lower := strings.ToLower(key)
	return []string{
		lower,                                        // sit-3658
		key,                                          // SIT-3658
		strings.ReplaceAll(lower, "-", "/"),          // sit/3658
		"feature/" + lower,                           // feature/sit-3658
		"fix/" + lower,                               // fix/sit-3658
		"bugfix/" + lower,                            // bugfix/sit-3658
	}
}
