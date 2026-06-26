package tools

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// PMRACI generates a RACI matrix from Jira sprint assignments.
func (h *Handlers) PMRACI(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}
	accountable := req.GetString("accountable", "")

	if h.Jira == nil {
		return errorResult("Jira not configured"), nil
	}

	// Get active sprint issues
	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil || len(sprints) == 0 {
		return errorResult("No active sprint found"), nil
	}

	sprintIssues, err := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
	if err != nil {
		return errorResult("Failed to fetch sprint issues"), nil
	}

	type matrixEntry struct {
		Member   string
		IssueKey string
		Summary  string
		Role     string // R, A, C, I
		Epic     string
	}

	var entries []matrixEntry
	memberSet := map[string]bool{}
	issueMap := map[string]string{} // key -> summary

	for _, iss := range sprintIssues {
		issueMap[iss.Key] = iss.Summary

		// R = Assignee (who does the work)
		if iss.Assignee != "" {
			epic := extractEpic(iss.Labels)
			entries = append(entries, matrixEntry{
				Member: iss.Assignee, IssueKey: iss.Key, Summary: truncateStr(iss.Summary, 50),
				Role: "R", Epic: epic,
			})
			memberSet[iss.Assignee] = true
		}

		// A = Accountable (reporter or specified)
		aPerson := accountable
		if aPerson == "" {
			aPerson = iss.Reporter
		}
		if aPerson != "" && aPerson != iss.Assignee {
			epic := extractEpic(iss.Labels)
			entries = append(entries, matrixEntry{
				Member: aPerson, IssueKey: iss.Key, Summary: truncateStr(iss.Summary, 50),
				Role: "A", Epic: epic,
			})
			memberSet[aPerson] = true
		}

		// I from labels that look like team names or stakeholders
		for _, label := range iss.Labels {
			if strings.HasPrefix(strings.ToLower(label), "team-") || strings.HasPrefix(strings.ToLower(label), "stake-") || strings.HasPrefix(strings.ToLower(label), "ui-") {
				memberName := strings.TrimPrefix(strings.TrimPrefix(strings.TrimPrefix(label, "team-"), "stake-"), "ui-")
				if memberName != "" && !memberSet[memberName] {
					entries = append(entries, matrixEntry{
						Member: memberName, IssueKey: iss.Key, Summary: truncateStr(iss.Summary, 50),
						Role: "I", Epic: extractEpic(iss.Labels),
					})
					memberSet[memberName] = true
				}
			}
		}
	}

	if len(entries) == 0 {
		return textResult("No issues with assignees found in current sprint."), nil
	}

	// Group by epic then member
	epicGroups := map[string][]matrixEntry{}
	epicOrder := []string{}
	for _, e := range entries {
		ep := e.Epic
		if ep == "" {
			ep = "(no epic)"
		}
		if _, ok := epicGroups[ep]; !ok {
			epicOrder = append(epicOrder, ep)
		}
		epicGroups[ep] = append(epicGroups[ep], e)
	}

	// Deduplicate: for each (member, issue, role), keep first
	type dedupKey struct {
		member, issue, role string
	}
	seen := map[dedupKey]bool{}
	deduped := []matrixEntry{}
	for _, ep := range epicOrder {
		for _, e := range epicGroups[ep] {
			k := dedupKey{e.Member, e.IssueKey, e.Role}
			if seen[k] {
				continue
			}
			seen[k] = true
			deduped = append(deduped, e)
		}
	}

	// Build RACI matrix
	// Collect all members and issues
	memberList := []string{}
	membersSeen := map[string]bool{}
	issueList := []string{}
	issuesSeen := map[string]bool{}

	for _, e := range deduped {
		if !membersSeen[e.Member] {
			memberList = append(memberList, e.Member)
			membersSeen[e.Member] = true
		}
		if !issuesSeen[e.IssueKey] {
			issueList = append(issueList, e.IssueKey)
			issuesSeen[e.IssueKey] = true
		}
	}

	sort.Strings(memberList)
	sort.Strings(issueList)

	// Build matrix: member x issue -> role
	matrix := map[string]map[string]string{}
	for _, member := range memberList {
		matrix[member] = map[string]string{}
	}
	for _, e := range deduped {
		if _, ok := matrix[e.Member][e.IssueKey]; !ok {
			matrix[e.Member][e.IssueKey] = e.Role
		}
	}

	// Render
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("RACI Matrix — Sprint: %s\n", sprints[0].Name))
	if sprints[0].Goal != "" {
		sb.WriteString(fmt.Sprintf("Goal: %s\n", sprints[0].Goal))
	}
	sb.WriteString("\nLegend: R=Responsible A=Accountable C=Consulted I=Informed\n\n")

	// Check if we have too many issues for a readable table
	if len(issueList) > 25 {
		// Compact view: aggregate by member
		sb.WriteString("Aggregate View (by member):\n\n")
		for _, member := range memberList {
			var roles []string
			for _, iss := range issueList {
				if r, ok := matrix[member][iss]; ok {
					roles = append(roles, r)
				}
			}
			if len(roles) > 0 {
				rCount := countRole(roles, "R")
				aCount := countRole(roles, "A")
				iCount := countRole(roles, "I")
				sb.WriteString(fmt.Sprintf("  %s: R=%d A=%d I=%d\n", member, rCount, aCount, iCount))
			}
		}

		sb.WriteString("\nDetailed View (sample):\n")
		// Show matrix for first 10 issues
		display := issueList
		if len(display) > 10 {
			display = display[:10]
		}
		for _, member := range memberList {
			var parts []string
			for _, iss := range display {
				if r, ok := matrix[member][iss]; ok {
					parts = append(parts, fmt.Sprintf("%s:%s", iss, r))
				}
			}
			if len(parts) > 0 {
				sb.WriteString(fmt.Sprintf("  %s: %s\n", member, strings.Join(parts, ", ")))
			}
		}
		if len(issueList) > 10 {
			sb.WriteString(fmt.Sprintf("  ... and %d more issues\n", len(issueList)-10))
		}
	} else {
		// Full matrix
		// Header
		sb.WriteString(fmt.Sprintf("%-20s", "Member"))
		for _, iss := range issueList {
			sb.WriteString(fmt.Sprintf(" %-12s", iss))
		}
		sb.WriteString("\n" + strings.Repeat("-", 20+13*len(issueList)) + "\n")

		for _, member := range memberList {
			sb.WriteString(fmt.Sprintf("%-20s", member))
			for _, iss := range issueList {
				if r, ok := matrix[member][iss]; ok {
					sb.WriteString(fmt.Sprintf(" %-12s", r))
				} else {
					sb.WriteString(fmt.Sprintf(" %-12s", "-"))
				}
			}
			sb.WriteString("\n")
		}
	}

	// Summary
	sb.WriteString(fmt.Sprintf("\nSummary: %d members, %d issues\n", len(memberList), len(issueList)))
	return textResult(sb.String()), nil
}

func extractEpic(labels []string) string {
	for _, l := range labels {
		if strings.HasPrefix(strings.ToLower(l), "epic-") || strings.HasPrefix(strings.ToLower(l), "epic:") {
			return strings.TrimPrefix(strings.TrimPrefix(l, "epic-"), "epic:")
		}
	}
	return ""
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func countRole(roles []string, role string) int {
	n := 0
	for _, r := range roles {
		if r == role {
			n++
		}
	}
	return n
}
