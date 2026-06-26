package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// PredictiveBlockers identifies team members likely to get blocked based on historical patterns.
func (h *Handlers) PredictiveBlockers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)

	// Analyze blocker history per person
	history, _ := h.Memory.GetBlockerHistory(ctx, 50)
	if len(history) < 5 {
		return textResult("Need more blocker history (5+ resolved) for prediction. Keep recording blockers with pm_record_blocker."), nil
	}

	// Count blockers per owner and per issue pattern
	ownerBlockCount := map[string]int{}
	for _, b := range history {
		if b.Owner != "" {
			ownerBlockCount[b.Owner]++
		}
	}

	// Current sprint: who is working on what
	var currentAssignments map[string][]string
	if boardID > 0 {
		sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
		if len(sprints) > 0 {
			issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
			currentAssignments = map[string][]string{}
			for _, i := range issues {
				l := strings.ToLower(i.Status)
				if !strings.Contains(l, "done") && !strings.Contains(l, "closed") && i.Assignee != "" {
					currentAssignments[i.Assignee] = append(currentAssignments[i.Assignee], i.Key)
				}
			}
		}
	}

	var sb strings.Builder
	sb.WriteString("Predictive Blocker Analysis:\n\n")

	// Find people with high block frequency
	var atRisk []string
	for person, count := range ownerBlockCount {
		if count >= 3 {
			atRisk = append(atRisk, fmt.Sprintf("  %s: blocked %d times historically", person, count))
		}
	}

	if len(atRisk) > 0 {
		sb.WriteString("HIGH RISK (frequently blocked):\n")
		for _, r := range atRisk {
			sb.WriteString(r + "\n")
		}

		// Cross-reference with current assignments
		if currentAssignments != nil {
			sb.WriteString("\nCurrently active (check on them):\n")
			for person := range ownerBlockCount {
				if issues, ok := currentAssignments[person]; ok && ownerBlockCount[person] >= 3 {
					sb.WriteString(fmt.Sprintf("  %s: working on %s\n", person, strings.Join(issues, ", ")))
				}
			}
		}
		sb.WriteString("\nAction: Proactively check with these people in standup.\n")
	} else {
		sb.WriteString("No high-risk patterns found. Team has low historical block rate.\n")
	}

	// Avg resolution time
	var totalDays int
	resolved := 0
	for _, b := range history {
		if b.DaysBlocked > 0 {
			totalDays += b.DaysBlocked
			resolved++
		}
	}
	if resolved > 0 {
		sb.WriteString(fmt.Sprintf("\nAvg blocker resolution: %d days (%d blockers)\n", totalDays/resolved, resolved))
		if totalDays/resolved > 3 {
			sb.WriteString("Warning: avg resolution >3 days. Escalation process may need improvement.\n")
		}
	}

	return textResult(sb.String()), nil
}

// SprintSimilarity compares current sprint signals to historical ones to find similar (potentially failing) sprints.
func (h *Handlers) SprintSimilarity(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	// Get current sprint state
	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) == 0 {
		return textResult("No active sprint."), nil
	}

	issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
	var nowDone, nowBlocked, nowTotal int
	nowTotal = len(issues)
	for _, i := range issues {
		l := strings.ToLower(i.Status)
		if strings.Contains(l, "done") || strings.Contains(l, "closed") {
			nowDone++
		} else if strings.Contains(l, "block") {
			nowBlocked++
		}
	}

	nowCompletion := 0.0
	if nowTotal > 0 {
		nowCompletion = float64(nowDone) / float64(nowTotal) * 100
	}

	// Compare against historical sprints
	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 10)
	if len(snaps) < 3 {
		return textResult("Need 3+ sprint snapshots for pattern matching."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sprint Similarity: %s\n", sprints[0].Name))
	sb.WriteString(fmt.Sprintf("Current: %d/%d done (%.0f%%), %d blocked\n\n", nowDone, nowTotal, nowCompletion, nowBlocked))

	// Find most similar historical sprint by completion rate + blocked ratio
	type match struct {
		name       string
		similarity float64
		completion float64
		outcome    string
	}
	var matches []match

	for _, snap := range snaps {
		if snap.TotalIssues == 0 {
			continue
		}
		snapCompletion := snap.CompletionRate
		blockedRatio := float64(snap.Blocked) / float64(snap.TotalIssues) * 100
		nowBlockedRatio := float64(nowBlocked) / float64(nowTotal) * 100

		// Simple similarity: how close is completion + blocked ratio
		completionDiff := abs(nowCompletion - snapCompletion)
		blockedDiff := abs(nowBlockedRatio - blockedRatio)
		similarity := 100 - (completionDiff + blockedDiff) / 2

		outcome := "OK"
		if snap.CompletionRate < 70 {
			outcome = "MISSED GOAL"
		}
		if snap.Carryover > 3 {
			outcome = "HIGH CARRYOVER"
		}

		matches = append(matches, match{
			name: snap.SprintName, similarity: similarity,
			completion: snap.CompletionRate, outcome: outcome,
		})
	}

	sb.WriteString("Most similar historical sprints:\n")
	for _, m := range matches {
		if m.similarity > 60 {
			warning := ""
			if m.outcome != "OK" {
				warning = " *** WARNING: that sprint " + m.outcome + " ***"
			}
			sb.WriteString(fmt.Sprintf("  %s (%.0f%% similar, ended %.0f%% complete)%s\n",
				m.name, m.similarity, m.completion, warning))
		}
	}

	// Overall prediction
	sb.WriteString("\n")
	if nowBlocked > 2 && nowCompletion < 50 {
		sb.WriteString("PREDICTION: AT RISK — similar to sprints that missed goals.\n")
		sb.WriteString("Action: Reduce scope or escalate blockers NOW.\n")
	} else if nowCompletion > 60 {
		sb.WriteString("PREDICTION: ON TRACK — similar to successful sprints.\n")
	} else {
		sb.WriteString("PREDICTION: WATCH — could go either way. Focus on finishing over starting.\n")
	}

	return textResult(sb.String()), nil
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// EarlyWarningSystem scans multiple signals and raises alerts proactively.
func (h *Handlers) EarlyWarningSystem(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)

	var warnings []string

	// Signal 1: Active blockers > 3 days
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	chronic := 0
	for _, b := range blockers {
		if time.Since(b.BlockedSince).Hours() > 72 {
			chronic++
		}
	}
	if chronic > 0 {
		warnings = append(warnings, fmt.Sprintf("BLOCKER: %d items stuck >3 days", chronic))
	}

	// Signal 2: Sprint completion too low for time elapsed
	if boardID > 0 {
		sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
		if len(sprints) > 0 {
			issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
			var done, total int
			total = len(issues)
			for _, i := range issues {
				l := strings.ToLower(i.Status)
				if strings.Contains(l, "done") || strings.Contains(l, "closed") {
					done++
				}
			}
			if total > 5 && float64(done)/float64(total) < 0.3 {
				warnings = append(warnings, fmt.Sprintf("PACE: Only %d/%d done (%.0f%%) — behind", done, total, float64(done)/float64(total)*100))
			}
		}
	}

	// Signal 3: High-priority risks without mitigation
	risks, _ := h.Memory.GetOpenRisks(ctx)
	unmitRisks := 0
	for _, r := range risks {
		if (r.Severity == "critical" || r.Severity == "high") && r.Mitigation == "" {
			unmitRisks++
		}
	}
	if unmitRisks > 0 {
		warnings = append(warnings, fmt.Sprintf("RISK: %d high/critical risks without mitigation plan", unmitRisks))
	}

	// Signal 4: Retro actions piling up
	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 5 {
		warnings = append(warnings, fmt.Sprintf("PROCESS: %d retro actions ignored (team trust eroding)", len(actions)))
	}

	// Signal 5: Velocity declining
	if boardID > 0 {
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 3)
		if len(snaps) >= 3 && snaps[0].Velocity < snaps[2].Velocity {
			warnings = append(warnings, fmt.Sprintf("VELOCITY: Declining (%d → %d over 3 sprints)", snaps[2].Velocity, snaps[0].Velocity))
		}
	}

	if len(warnings) == 0 {
		return textResult("All clear. No early warning signals detected."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("EARLY WARNING: %d signal(s) detected\n\n", len(warnings)))
	for i, w := range warnings {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, w))
	}
	sb.WriteString("\nThese are leading indicators — act now before they become sprint failures.")

	return textResult(sb.String()), nil
}
