package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// DailyDelta shows what changed since yesterday — the PM's morning briefing.
func (h *Handlers) DailyDelta(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) == 0 {
		return textResult("No active sprint."), nil
	}

	sprint := sprints[0]
	issues, _ := h.Jira.GetSprintIssues(ctx, sprint.ID)

	// Compare with yesterday's snapshot
	progress, _ := h.Memory.GetDailyProgress(ctx, boardID, sprint.Name)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Daily Delta — %s (%s)\n\n", sprint.Name, time.Now().Format("Mon Jan 2")))

	// Current state
	var done, inProgress, blocked, todo int
	for _, i := range issues {
		l := strings.ToLower(i.Status)
		switch {
		case strings.Contains(l, "done") || strings.Contains(l, "closed"):
			done++
		case strings.Contains(l, "progress") || strings.Contains(l, "review"):
			inProgress++
		case strings.Contains(l, "block"):
			blocked++
		default:
			todo++
		}
	}

	if len(progress) > 0 {
		yesterday := progress[len(progress)-1]
		doneChange := done - yesterday.Done
		blockedChange := blocked - yesterday.Blocked

		sb.WriteString("CHANGES SINCE LAST CHECK:\n")
		if doneChange > 0 {
			sb.WriteString(fmt.Sprintf("  +%d items completed\n", doneChange))
		}
		if blockedChange > 0 {
			sb.WriteString(fmt.Sprintf("  +%d NEW blockers\n", blockedChange))
		} else if blockedChange < 0 {
			sb.WriteString(fmt.Sprintf("  %d blockers resolved\n", -blockedChange))
		}
		if doneChange == 0 && blockedChange == 0 {
			sb.WriteString("  No significant changes\n")
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("CURRENT: Done %d | In Progress %d | Todo %d | Blocked %d | Total %d\n", done, inProgress, todo, blocked, len(issues)))

	// Items updated in last 24h
	recentlyMoved := 0
	for _, i := range issues {
		if time.Since(i.Updated).Hours() < 24 {
			recentlyMoved++
		}
	}
	sb.WriteString(fmt.Sprintf("Active today: %d items updated in last 24h\n\n", recentlyMoved))

	// Burn rate
	remaining := len(issues) - done
	if len(progress) >= 2 {
		first := progress[0]
		daysElapsed := time.Since(first.Date).Hours() / 24
		itemsDone := done - first.Done
		if daysElapsed > 0 && itemsDone > 0 {
			rate := float64(itemsDone) / daysElapsed
			daysLeft := float64(remaining) / rate
			sb.WriteString(fmt.Sprintf("PACE: %.1f items/day → %d remaining → ~%d days to finish\n", rate, remaining, int(daysLeft)))
		}
	}

	// Who might need help (high WIP)
	personWIP := map[string]int{}
	for _, i := range issues {
		l := strings.ToLower(i.Status)
		if strings.Contains(l, "progress") || strings.Contains(l, "review") || strings.Contains(l, "block") {
			personWIP[i.Assignee]++
		}
	}
	var overloaded []string
	for person, wip := range personWIP {
		if wip >= 3 && person != "" {
			overloaded = append(overloaded, fmt.Sprintf("%s (%d items)", person, wip))
		}
	}
	if len(overloaded) > 0 {
		sb.WriteString(fmt.Sprintf("\nNEEDS ATTENTION: %s — high WIP, may need help\n", strings.Join(overloaded, ", ")))
	}

	return textResult(sb.String()), nil
}

// OverloadCheck detects team members who may be overloaded.
func (h *Handlers) OverloadCheck(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) == 0 {
		return textResult("No active sprint."), nil
	}

	issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)

	// Per-person load analysis
	type personLoad struct {
		assigned   int
		inProgress int
		blocked    int
		done       int
	}
	loads := map[string]*personLoad{}

	for _, i := range issues {
		if i.Assignee == "" {
			continue
		}
		if loads[i.Assignee] == nil {
			loads[i.Assignee] = &personLoad{}
		}
		p := loads[i.Assignee]
		p.assigned++
		l := strings.ToLower(i.Status)
		switch {
		case strings.Contains(l, "done") || strings.Contains(l, "closed"):
			p.done++
		case strings.Contains(l, "block"):
			p.blocked++
		case strings.Contains(l, "progress") || strings.Contains(l, "review"):
			p.inProgress++
		}
	}

	if len(loads) == 0 {
		return textResult("No assignees in sprint."), nil
	}

	// Calculate average
	totalAssigned := 0
	for _, p := range loads {
		totalAssigned += p.assigned
	}
	avgLoad := totalAssigned / len(loads)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Team Load Analysis (%s)\n\n", sprints[0].Name))
	sb.WriteString(fmt.Sprintf("Average load: %d items/person\n", avgLoad))
	sb.WriteString(fmt.Sprintf("Team size: %d\n\n", len(loads)))
	sb.WriteString("Person          | Assigned | WIP | Blocked | Done | Signal\n")
	sb.WriteString("----------------|----------|-----|---------|------|--------\n")

	for person, p := range loads {
		signal := "OK"
		wip := p.inProgress + p.blocked
		if p.assigned > avgLoad*2 {
			signal = "OVERLOADED"
		} else if wip >= 3 {
			signal = "HIGH WIP"
		} else if p.blocked > 0 && p.inProgress == 0 {
			signal = "STUCK"
		} else if p.assigned == 0 {
			signal = "IDLE?"
		}

		name := person
		if len(name) > 15 {
			name = name[:15]
		}
		sb.WriteString(fmt.Sprintf("%-15s | %d | %d | %d | %d | %s\n",
			name, p.assigned, wip, p.blocked, p.done, signal))
	}

	// Summary signals
	sb.WriteString("\n")
	overloaded := 0
	stuck := 0
	for _, p := range loads {
		if p.assigned > avgLoad*2 {
			overloaded++
		}
		if p.blocked > 0 && p.inProgress == 0 {
			stuck++
		}
	}

	if overloaded > 0 {
		sb.WriteString(fmt.Sprintf("WARNING: %d team member(s) have 2x+ average load.\n", overloaded))
		sb.WriteString("  Suggestion: Redistribute work or reduce sprint scope.\n")
	}
	if stuck > 0 {
		sb.WriteString(fmt.Sprintf("WARNING: %d team member(s) are fully blocked.\n", stuck))
		sb.WriteString("  Suggestion: Prioritize unblocking them in standup.\n")
	}

	// Sustainable pace check
	sb.WriteString("\nSUSTAINABLE PACE:\n")
	sb.WriteString("  Recommended max WIP per person: 2 items\n")
	sb.WriteString(fmt.Sprintf("  Recommended max load: %d items (80%% of avg capacity)\n", int(float64(avgLoad)*1.5)))
	highWIP := 0
	for _, p := range loads {
		if p.inProgress+p.blocked > 2 {
			highWIP++
		}
	}
	if highWIP > 0 {
		sb.WriteString(fmt.Sprintf("  %d person(s) above WIP limit — encourage finishing over starting.\n", highWIP))
	}

	return textResult(sb.String()), nil
}

// CommitmentCheck validates if sprint commitment is realistic.
func (h *Handlers) CommitmentCheck(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	// Get current sprint load
	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) == 0 {
		return textResult("No active sprint."), nil
	}

	issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
	currentLoad := len(issues)

	// Get historical velocity (items done per sprint)
	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Commitment Check: %s\n\n", sprints[0].Name))
	sb.WriteString(fmt.Sprintf("Items committed: %d\n", currentLoad))

	if len(snaps) >= 2 {
		var totalDone int
		for _, s := range snaps {
			totalDone += s.Done
		}
		avgDone := totalDone / len(snaps)
		ratio := float64(currentLoad) / float64(avgDone) * 100

		sb.WriteString(fmt.Sprintf("Historical avg completion: %d items/sprint (%d sprints)\n", avgDone, len(snaps)))
		sb.WriteString(fmt.Sprintf("Commitment ratio: %.0f%%\n\n", ratio))

		switch {
		case ratio > 120:
			sb.WriteString("VERDICT: OVERCOMMITTED\n")
			sb.WriteString(fmt.Sprintf("  You committed %d but historically finish ~%d.\n", currentLoad, avgDone))
			sb.WriteString("  Risk: carryover, burnout, missed sprint goal.\n")
			sb.WriteString(fmt.Sprintf("  Suggestion: Remove %d items or mark as stretch goals.\n", currentLoad-avgDone))
		case ratio > 100:
			sb.WriteString("VERDICT: STRETCHING\n")
			sb.WriteString("  Slightly above historical average. Achievable if no surprises.\n")
		case ratio >= 80:
			sb.WriteString("VERDICT: SUSTAINABLE\n")
			sb.WriteString("  Good balance between ambition and realism.\n")
		default:
			sb.WriteString("VERDICT: CONSERVATIVE\n")
			sb.WriteString("  Below historical average. Team has capacity for more OR buffer for quality.\n")
		}
	} else {
		sb.WriteString("\nNot enough history for comparison. Capture 3+ sprint snapshots.\n")
	}

	return textResult(sb.String()), nil
}

// TeamCareReport generates a "how is the team doing?" report focused on wellbeing signals.
func (h *Handlers) TeamCareReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)

	var sb strings.Builder
	sb.WriteString("=== TEAM CARE REPORT ===\n")
	sb.WriteString("(Focus: wellbeing, not productivity)\n\n")

	// 1. Overload signals
	if boardID > 0 {
		sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
		if len(sprints) > 0 {
			issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
			personWIP := map[string]int{}
			personBlocked := map[string]int{}
			for _, i := range issues {
				if i.Assignee == "" {
					continue
				}
				l := strings.ToLower(i.Status)
				if strings.Contains(l, "progress") || strings.Contains(l, "review") {
					personWIP[i.Assignee]++
				}
				if strings.Contains(l, "block") {
					personBlocked[i.Assignee]++
				}
			}

			sb.WriteString("WORKLOAD:\n")
			overloaded := 0
			for person, wip := range personWIP {
				if wip >= 3 {
					sb.WriteString(fmt.Sprintf("  %s has %d items in progress (recommended: 1-2)\n", person, wip))
					overloaded++
				}
			}
			if overloaded == 0 {
				sb.WriteString("  Everyone within healthy WIP limits.\n")
			}
			sb.WriteString("\n")

			// Blocked people
			if len(personBlocked) > 0 {
				sb.WriteString("BLOCKED TEAM MEMBERS:\n")
				for person, count := range personBlocked {
					sb.WriteString(fmt.Sprintf("  %s — %d items blocked (they may feel frustrated/stuck)\n", person, count))
				}
				sb.WriteString("  Action: Check in personally. Ask 'what do you need?'\n\n")
			}
		}
	}

	// 2. Carryover pattern (chronic = demotivating)
	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 3)
	if len(snaps) >= 2 {
		chronicCarryover := 0
		for _, s := range snaps {
			if s.Carryover > 2 {
				chronicCarryover++
			}
		}
		if chronicCarryover >= 2 {
			sb.WriteString("RECURRING CARRYOVER:\n")
			sb.WriteString("  Work keeps rolling over. This can feel demoralizing.\n")
			sb.WriteString("  Ask team: 'Are we overcommitting? Should we reduce scope?'\n")
			sb.WriteString("  Research: Chronic overcommitment is an ethical issue (PMI 2024)\n\n")
		}
	}

	// 3. Pending retro actions (team feels unheard if actions don't happen)
	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 3 {
		sb.WriteString("RETRO FOLLOW-THROUGH:\n")
		sb.WriteString(fmt.Sprintf("  %d action items still pending.\n", len(actions)))
		sb.WriteString("  When retro actions don't happen, team stops raising issues.\n")
		sb.WriteString("  This erodes psychological safety over time.\n")
		sb.WriteString("  Action: Close 2 items this sprint or explicitly cancel with explanation.\n\n")
	}

	// 4. Blocker duration (frustration signal)
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	longBlockers := 0
	for _, b := range blockers {
		if time.Since(b.BlockedSince).Hours() > 72 {
			longBlockers++
		}
	}
	if longBlockers > 0 {
		sb.WriteString("CHRONIC BLOCKERS:\n")
		sb.WriteString(fmt.Sprintf("  %d blockers open >3 days. Blocked developers feel helpless.\n", longBlockers))
		sb.WriteString("  Action: Escalate today. Remove impediments — that's YOUR job as SM.\n\n")
	}

	// Summary
	sb.WriteString("---\n")
	sb.WriteString("Remember: Sustainable pace means the team can maintain this pace indefinitely.\n")
	sb.WriteString("If anyone is working overtime to hit sprint goals, the sprint goal is wrong.\n")
	sb.WriteString("The goal of Scrum is not velocity — it's sustainable delivery of value.\n")

	return textResult(sb.String()), nil
}
