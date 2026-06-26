package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// WhatNext tells the SM what to focus on right now based on all signals.
func (h *Handlers) WhatNext(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)

	var signals []string
	var priorities []string

	// Check blockers
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	for _, b := range blockers {
		days := int(time.Since(b.BlockedSince).Hours() / 24)
		if days >= 3 {
			priorities = append(priorities, fmt.Sprintf("ESCALATE blocker: %s (%d days, owner: %s)", b.Description, days, b.Owner))
		} else if days >= 1 {
			signals = append(signals, fmt.Sprintf("Monitor blocker: %s (%d days)", b.Description, days))
		}
	}

	// Check overdue actions
	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 3 {
		priorities = append(priorities, fmt.Sprintf("Follow up on %d pending retro actions (they're rotting)", len(actions)))
	}

	// Check risks
	risks, _ := h.Memory.GetOpenRisks(ctx)
	critical := 0
	for _, r := range risks {
		if r.Severity == "critical" { critical++ }
	}
	if critical > 0 {
		priorities = append(priorities, fmt.Sprintf("Address %d CRITICAL risks immediately", critical))
	}

	// Sprint progress check
	if boardID > 0 {
		sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
		if len(sprints) > 0 {
			issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
			var done, blocked, total int
			total = len(issues)
			for _, i := range issues {
				lower := strings.ToLower(i.Status)
				if strings.Contains(lower, "done") || strings.Contains(lower, "closed") {
					done++
				} else if strings.Contains(lower, "block") {
					blocked++
				}
			}
			if total > 0 {
				completion := float64(done) / float64(total) * 100
				if completion < 30 && blocked > 0 {
					priorities = append(priorities, "Sprint behind + blockers. Urgent: unblock the team or reduce scope")
				}
				if blocked > 2 {
					priorities = append(priorities, fmt.Sprintf("Swarming needed: %d items blocked. Facilitate impediment removal", blocked))
				}
			}
		}
	}

	// Dependencies
	deps, _ := h.Memory.GetOpenDependencies(ctx)
	agingDeps := 0
	for _, d := range deps {
		if time.Since(d.CreatedAt).Hours() > 72 { agingDeps++ }
	}
	if agingDeps > 0 {
		signals = append(signals, fmt.Sprintf("Chase %d aging dependencies (>3 days)", agingDeps))
	}

	// Build response
	var sb strings.Builder
	sb.WriteString("What to Focus on Now\n\n")

	if len(priorities) > 0 {
		sb.WriteString("URGENT:\n")
		for i, p := range priorities {
			sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, p))
		}
		sb.WriteString("\n")
	}

	if len(signals) > 0 {
		sb.WriteString("MONITOR:\n")
		for _, s := range signals {
			sb.WriteString(fmt.Sprintf("  - %s\n", s))
		}
		sb.WriteString("\n")
	}

	if len(priorities) == 0 && len(signals) == 0 {
		sb.WriteString("All clear. Good day to:\n")
		sb.WriteString("  - Invest in team development (1-on-1s, coaching)\n")
		sb.WriteString("  - Review and groom the backlog\n")
		sb.WriteString("  - Work on process improvements from last retro\n")
	}

	return textResult(sb.String()), nil
}

// OneOnOnePrep generates prep notes for a 1-on-1 with a team member.
func (h *Handlers) OneOnOnePrep(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	memberName, err := req.RequireString("member")
	if err != nil {
		return errorResult("member name required"), nil
	}

	var data strings.Builder
	data.WriteString(fmt.Sprintf("1-on-1 Prep: %s\n\n", memberName))

	// Get their metrics
	metrics, _ := h.Memory.GetTeamMetrics(ctx, memberName, 5)
	if len(metrics) > 0 {
		data.WriteString("Recent Performance:\n")
		for _, m := range metrics {
			completion := 0.0
			if m.IssuesAssigned > 0 {
				completion = float64(m.IssuesDone) / float64(m.IssuesAssigned) * 100
			}
			data.WriteString(fmt.Sprintf("  %s: %d/%d done (%.0f%%), blockers: %d, carryover: %d\n",
				m.SprintName, m.IssuesDone, m.IssuesAssigned, completion, m.BlockerCount, m.CarryoverCount))
		}
		data.WriteString("\n")

		// Patterns
		totalBlockers := 0
		totalCarryover := 0
		for _, m := range metrics {
			totalBlockers += m.BlockerCount
			totalCarryover += m.CarryoverCount
		}
		if totalBlockers > len(metrics)*2 {
			data.WriteString("Pattern: Frequently blocked. Discuss: what's causing the blocks?\n")
		}
		if totalCarryover > len(metrics)*2 {
			data.WriteString("Pattern: High carryover. Discuss: estimation accuracy, scope of stories\n")
		}
	}

	// Current workload from Jira
	jql := fmt.Sprintf("assignee = \"%s\" AND resolution = Unresolved ORDER BY priority DESC", memberName)
	result, _ := h.Jira.SearchIssues(ctx, jql, 10, 0)
	if result != nil && len(result.Issues) > 0 {
		data.WriteString(fmt.Sprintf("Current Open Items (%d):\n", len(result.Issues)))
		for _, i := range result.Issues {
			data.WriteString(fmt.Sprintf("  [%s] %s - %s\n", i.Priority, i.Key, i.Summary))
		}
		data.WriteString("\n")
		if len(result.Issues) > 7 {
			data.WriteString("NOTE: Heavy workload (>7 items). Check if overloaded.\n\n")
		}
	}

	// AI suggestions for talking points
	systemPrompt := `Based on this team member's data, suggest 3-4 talking points for a 1-on-1.
Focus on: growth opportunities, blockers to discuss, workload concerns, recognition.
Be specific and actionable. No generic questions.`

	suggestions, err := h.aiComplete(ctx, systemPrompt, data.String())
	if err != nil {
		data.WriteString("Suggested topics: workload balance, blockers, career growth, feedback\n")
		return textResult(data.String()), nil
	}

	return textResult(data.String() + "Talking Points:\n" + suggestions), nil
}

// SprintNarrative generates a "story" of the sprint for review/demo.
func (h *Handlers) SprintNarrative(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil || len(sprints) == 0 {
		return textResult("No active sprint."), nil
	}

	sprint := sprints[0]
	issues, _ := h.Jira.GetSprintIssues(ctx, sprint.ID)

	// Categorize work
	var features, bugs, tech []string
	var done, total int
	total = len(issues)
	for _, i := range issues {
		lower := strings.ToLower(i.Status)
		if !strings.Contains(lower, "done") && !strings.Contains(lower, "closed") {
			continue
		}
		done++
		entry := fmt.Sprintf("%s: %s", i.Key, i.Summary)
		switch strings.ToLower(i.Type) {
		case "story", "feature":
			features = append(features, entry)
		case "bug":
			bugs = append(bugs, entry)
		default:
			tech = append(tech, entry)
		}
	}

	// Get blockers that happened
	blockers, _ := h.Memory.GetBlockerHistory(ctx, 10)
	var sprintBlockers []string
	for _, b := range blockers {
		if b.ResolvedAt != nil {
			sprintBlockers = append(sprintBlockers, fmt.Sprintf("%s (resolved in %d days)", b.Description, b.DaysBlocked))
		}
	}

	var data strings.Builder
	data.WriteString(fmt.Sprintf("Sprint: %s | Goal: %s\n", sprint.Name, sprint.Goal))
	data.WriteString(fmt.Sprintf("Completed: %d/%d items\n\n", done, total))
	if len(features) > 0 {
		data.WriteString("Features delivered:\n")
		for _, f := range features { data.WriteString("  - " + f + "\n") }
	}
	if len(bugs) > 0 {
		data.WriteString("Bugs fixed:\n")
		for _, b := range bugs { data.WriteString("  - " + b + "\n") }
	}
	if len(tech) > 0 {
		data.WriteString("Technical work:\n")
		for _, t := range tech { data.WriteString("  - " + t + "\n") }
	}
	if len(sprintBlockers) > 0 {
		data.WriteString("Challenges overcome:\n")
		for _, b := range sprintBlockers { data.WriteString("  - " + b + "\n") }
	}

	systemPrompt := `Write a sprint narrative for a Sprint Review demo. This is the "story" of the sprint.
Structure:
1. Sprint Goal and whether we achieved it (1 sentence)
2. Key highlights (what users/stakeholders will notice)
3. Challenges we overcame
4. What we learned

Tone: confident, transparent, stakeholder-friendly. Under 150 words.
Do NOT list ticket IDs. Describe outcomes in business language.`

	narrative, err := h.aiComplete(ctx, systemPrompt, data.String())
	if err != nil {
		return textResult("Sprint Narrative:\n\n" + data.String()), nil
	}

	return textResult(narrative), nil
}
