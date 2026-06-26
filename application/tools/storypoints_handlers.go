package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// StoryPointsSummary calculates total story points from a JQL query, sprint, or epic.
func (h *Handlers) StoryPointsSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)
	jql := req.GetString("jql", "")
	epicKey := req.GetString("epic_key", "")
	groupBy := req.GetString("group_by", "") // status, assignee, type

	// Determine what to query
	if jql == "" {
		if epicKey != "" {
			jql = fmt.Sprintf("\"Epic Link\" = %s OR parent = %s ORDER BY status ASC", epicKey, epicKey)
		} else if boardID > 0 {
			// Use active sprint
			sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
			if err != nil || len(sprints) == 0 {
				return textResult("No active sprint found. Provide jql or epic_key."), nil
			}
			jql = fmt.Sprintf("sprint = %d ORDER BY status ASC", sprints[0].ID)
		} else {
			return errorResult("Provide board_id (for sprint), epic_key, or jql"), nil
		}
	}

	result, err := h.Jira.SearchIssues(ctx, jql, 100, 0)
	if err != nil {
		return errorResult("Search failed: " + err.Error()), nil
	}

	// Calculate totals
	var totalPoints float64
	var estimated, unestimated int
	byStatus := map[string]float64{}
	byAssignee := map[string]float64{}
	byType := map[string]float64{}

	for _, issue := range result.Issues {
		if issue.StoryPoints > 0 {
			totalPoints += issue.StoryPoints
			estimated++
			byStatus[issue.Status] += issue.StoryPoints
			byAssignee[issue.Assignee] += issue.StoryPoints
			byType[issue.Type] += issue.StoryPoints
		} else {
			unestimated++
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Story Points Summary (%d issues)\n\n", len(result.Issues)))
	sb.WriteString(fmt.Sprintf("Total Points: %.0f\n", totalPoints))
	sb.WriteString(fmt.Sprintf("Estimated: %d | Unestimated: %d\n\n", estimated, unestimated))

	// Group by requested dimension
	switch groupBy {
	case "status":
		sb.WriteString("By Status:\n")
		for status, pts := range byStatus {
			sb.WriteString(fmt.Sprintf("  %s: %.0f pts\n", status, pts))
		}
	case "assignee":
		sb.WriteString("By Assignee:\n")
		for person, pts := range byAssignee {
			name := person
			if name == "" {
				name = "(unassigned)"
			}
			sb.WriteString(fmt.Sprintf("  %s: %.0f pts\n", name, pts))
		}
	case "type":
		sb.WriteString("By Type:\n")
		for issueType, pts := range byType {
			sb.WriteString(fmt.Sprintf("  %s: %.0f pts\n", issueType, pts))
		}
	default:
		// Show all groupings
		if len(byStatus) > 0 {
			sb.WriteString("By Status:\n")
			for status, pts := range byStatus {
				pct := pts / totalPoints * 100
				sb.WriteString(fmt.Sprintf("  %s: %.0f pts (%.0f%%)\n", status, pts, pct))
			}
			sb.WriteString("\n")
		}
		if len(byAssignee) > 0 {
			sb.WriteString("By Assignee:\n")
			for person, pts := range byAssignee {
				name := person
				if name == "" {
					name = "(unassigned)"
				}
				sb.WriteString(fmt.Sprintf("  %s: %.0f pts\n", name, pts))
			}
		}
	}

	if unestimated > 0 {
		sb.WriteString(fmt.Sprintf("\nWARNING: %d items have no story points. Estimation incomplete.\n", unestimated))
	}

	return textResult(sb.String()), nil
}

// SprintPointsBurndown shows story points progress for active sprint.
func (h *Handlers) SprintPointsBurndown(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil || len(sprints) == 0 {
		return textResult("No active sprint."), nil
	}

	sprint := sprints[0]
	issues, err := h.Jira.GetSprintIssues(ctx, sprint.ID)
	if err != nil {
		return errorResult("Failed to get issues: " + err.Error()), nil
	}

	var totalPts, donePts, inProgressPts, todoPts, blockedPts float64
	for _, issue := range issues {
		pts := issue.StoryPoints
		totalPts += pts
		lower := strings.ToLower(issue.Status)
		switch {
		case strings.Contains(lower, "done") || strings.Contains(lower, "closed") || strings.Contains(lower, "resolved"):
			donePts += pts
		case strings.Contains(lower, "progress") || strings.Contains(lower, "review"):
			inProgressPts += pts
		case strings.Contains(lower, "block"):
			blockedPts += pts
		default:
			todoPts += pts
		}
	}

	remainingPts := totalPts - donePts
	burnPct := 0.0
	if totalPts > 0 {
		burnPct = donePts / totalPts * 100
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sprint Points Burndown: %s\n\n", sprint.Name))
	sb.WriteString(fmt.Sprintf("Total:       %.0f pts\n", totalPts))
	sb.WriteString(fmt.Sprintf("Done:        %.0f pts (%.0f%%)\n", donePts, burnPct))
	sb.WriteString(fmt.Sprintf("In Progress: %.0f pts\n", inProgressPts))
	sb.WriteString(fmt.Sprintf("Todo:        %.0f pts\n", todoPts))
	sb.WriteString(fmt.Sprintf("Blocked:     %.0f pts\n", blockedPts))
	sb.WriteString(fmt.Sprintf("Remaining:   %.0f pts\n", remainingPts))

	if totalPts == 0 {
		sb.WriteString("\nNOTE: No story points found. Either issues are unestimated or the story points field ID doesn't match.\n")
		sb.WriteString("Check: jira_fields to find your story points custom field ID.\n")
	}

	return textResult(sb.String()), nil
}
