package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// ReportToPO generates a Product Owner briefing focused on value delivery, blocked items, and scope decisions.
func (h *Handlers) ReportToPO(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	var data strings.Builder

	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) > 0 {
		sprint := sprints[0]
		issues, _ := h.Jira.GetSprintIssues(ctx, sprint.ID)
		var done, inProgress, blocked, total int
		total = len(issues)
		var doneItems, blockedItems []string
		for _, i := range issues {
			lower := strings.ToLower(i.Status)
			switch {
			case strings.Contains(lower, "done") || strings.Contains(lower, "closed"):
				done++
				doneItems = append(doneItems, fmt.Sprintf("%s: %s", i.Key, i.Summary))
			case strings.Contains(lower, "block"):
				blocked++
				blockedItems = append(blockedItems, fmt.Sprintf("%s: %s (assignee: %s)", i.Key, i.Summary, i.Assignee))
			case strings.Contains(lower, "progress") || strings.Contains(lower, "review"):
				inProgress++
			}
		}
		pct := 0.0
		if total > 0 {
			pct = float64(done) / float64(total) * 100
		}
		data.WriteString(fmt.Sprintf("Sprint: %s | Goal: %s\n", sprint.Name, sprint.Goal))
		data.WriteString(fmt.Sprintf("Progress: %d/%d done (%.0f%%)\n", done, total, pct))
		data.WriteString(fmt.Sprintf("In flight: %d | Blocked: %d\n\n", inProgress, blocked))
		if len(doneItems) > 0 {
			data.WriteString("Delivered this sprint:\n")
			for _, item := range doneItems {
				data.WriteString("  - " + item + "\n")
			}
			data.WriteString("\n")
		}
		if len(blockedItems) > 0 {
			data.WriteString("Blocked (needs PO attention):\n")
			for _, item := range blockedItems {
				data.WriteString("  - " + item + "\n")
			}
			data.WriteString("\n")
		}
	}

	risks, _ := h.Memory.GetOpenRisks(ctx)
	scopeRisks := 0
	for _, r := range risks {
		if strings.Contains(strings.ToLower(r.Title), "scope") || strings.Contains(strings.ToLower(r.Mitigation), "remove") {
			scopeRisks++
		}
	}
	if scopeRisks > 0 {
		data.WriteString(fmt.Sprintf("Scope decisions pending: %d items may need removal from sprint\n\n", scopeRisks))
	}

	systemPrompt := `Format this data as a brief for the Product Owner. Focus on:
1. What value was delivered (use feature names, not ticket IDs)
2. What's blocked and what decision the PO needs to make
3. Sprint goal status (on track / at risk)
4. Any scope changes recommended
Keep it under 200 words. Use clear, non-technical language.`

	report, err := h.aiComplete(ctx, systemPrompt, data.String())
	if err != nil {
		return textResult("PO Briefing:\n\n" + data.String()), nil
	}
	return textResult(report), nil
}

// EscalationBrief generates a structured impediment escalation for management with PROBLEM, IMPACT, ASK, and DEADLINE.
func (h *Handlers) EscalationBrief(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)

	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) == 0 {
		return textResult("No active impediments to escalate. All clear."), nil
	}

	var data strings.Builder
	data.WriteString(fmt.Sprintf("IMPEDIMENT ESCALATION BRIEF\nDate: %s\n\n", time.Now().Format("2006-01-02")))
	data.WriteString(fmt.Sprintf("Active Impediments: %d\n\n", len(blockers)))

	for i, b := range blockers {
		days := int(time.Since(b.BlockedSince).Hours() / 24)
		severity := "MEDIUM"
		if days > 5 {
			severity = "HIGH"
		}
		if days > 10 {
			severity = "CRITICAL"
		}
		data.WriteString(fmt.Sprintf("--- Impediment #%d [%s] ---\n", i+1, severity))
		data.WriteString(fmt.Sprintf("Issue: %s\n", b.IssueKey))
		data.WriteString(fmt.Sprintf("Description: %s\n", b.Description))
		data.WriteString(fmt.Sprintf("Blocked since: %s (%d days)\n", b.BlockedSince.Format("2006-01-02"), days))
		data.WriteString(fmt.Sprintf("Owner: %s\n\n", b.Owner))
	}

	if boardID > 0 {
		sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
		if len(sprints) > 0 {
			sprint := sprints[0]
			issues, _ := h.Jira.GetSprintIssues(ctx, sprint.ID)
			total := len(issues)
			if total > 0 {
				blockedPct := float64(len(blockers)) / float64(total) * 100
				data.WriteString(fmt.Sprintf("SPRINT IMPACT: %.0f%% of sprint items are blocked\n", blockedPct))
				if blockedPct > 20 {
					data.WriteString("RISK: Sprint goal is at risk without management intervention\n")
				}
			}
		}
	}

	systemPrompt := `Format as a management escalation brief. For each impediment, specify:
1. PROBLEM: one sentence
2. IMPACT: what happens if unresolved (sprint goal at risk, delivery date slips, etc.)
3. ASK: specific decision or action needed from management
4. DEADLINE: by when the decision is needed

Keep professional, concise, action-oriented. No jargon.`

	report, err := h.aiComplete(ctx, systemPrompt, data.String())
	if err != nil {
		return textResult(data.String()), nil
	}
	return textResult(report), nil
}

// CrossTeamDependencyReport shows what we're waiting on from other teams and what they wait on from us.
func (h *Handlers) CrossTeamDependencyReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	deps, _ := h.Memory.GetOpenDependencies(ctx)
	if len(deps) == 0 {
		return textResult("No cross-team dependencies tracked. Use pm_record_dependency to add them."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CROSS-TEAM DEPENDENCY STATUS\nDate: %s\n\n", time.Now().Format("2006-01-02")))

	var blocking, blockedBy, external []string
	for _, d := range deps {
		days := int(time.Since(d.CreatedAt).Hours() / 24)
		entry := fmt.Sprintf("%s -> %s [%s] | %d days | %s", d.FromIssueKey, d.ToIssueKey, d.DependencyType, days, d.Description)
		switch d.DependencyType {
		case "blocks":
			blocking = append(blocking, entry)
		case "blocked_by":
			blockedBy = append(blockedBy, entry)
		case "external":
			external = append(external, entry)
		default:
			blockedBy = append(blockedBy, entry)
		}
	}

	if len(blockedBy) > 0 {
		sb.WriteString(fmt.Sprintf("WE ARE WAITING ON (%d):\n", len(blockedBy)))
		for _, e := range blockedBy {
			sb.WriteString("  - " + e + "\n")
		}
		sb.WriteString("\n")
	}
	if len(blocking) > 0 {
		sb.WriteString(fmt.Sprintf("OTHERS WAITING ON US (%d):\n", len(blocking)))
		for _, e := range blocking {
			sb.WriteString("  - " + e + "\n")
		}
		sb.WriteString("\n")
	}
	if len(external) > 0 {
		sb.WriteString(fmt.Sprintf("EXTERNAL DEPENDENCIES (%d):\n", len(external)))
		for _, e := range external {
			sb.WriteString("  - " + e + "\n")
		}
		sb.WriteString("\n")
	}

	critical := 0
	for _, d := range deps {
		if time.Since(d.CreatedAt).Hours() > 120 {
			critical++
		}
	}
	if critical > 0 {
		sb.WriteString(fmt.Sprintf("ALERT: %d dependencies aged >5 days - escalation recommended\n", critical))
	}

	return textResult(sb.String()), nil
}

// DeliveryConfidenceReport provides a GREEN/AMBER/RED confidence assessment for management.
func (h *Handlers) DeliveryConfidenceReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	var done, inProgress, todo, blocked, total int
	total = len(issues)
	for _, i := range issues {
		lower := strings.ToLower(i.Status)
		switch {
		case strings.Contains(lower, "done") || strings.Contains(lower, "closed"):
			done++
		case strings.Contains(lower, "progress") || strings.Contains(lower, "review"):
			inProgress++
		case strings.Contains(lower, "block"):
			blocked++
		default:
			todo++
		}
	}

	completion := 0.0
	if total > 0 {
		completion = float64(done) / float64(total) * 100
	}

	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
	avgCompletion := 0.0
	if len(snaps) > 0 {
		for _, s := range snaps {
			avgCompletion += s.CompletionRate
		}
		avgCompletion /= float64(len(snaps))
	}

	confidence := "HIGH"
	color := "GREEN"
	if blocked > 2 || (total > 0 && float64(blocked)/float64(total) > 0.2) {
		confidence = "LOW"
		color = "RED"
	} else if (total-done) > done && todo > inProgress {
		confidence = "MEDIUM"
		color = "AMBER"
	}

	var sb strings.Builder
	sb.WriteString("DELIVERY CONFIDENCE REPORT\n\n")
	sb.WriteString(fmt.Sprintf("Sprint: %s\n", sprint.Name))
	sb.WriteString(fmt.Sprintf("Goal: %s\n\n", sprint.Goal))
	sb.WriteString(fmt.Sprintf("Status: %s (%s)\n\n", color, confidence))
	sb.WriteString(fmt.Sprintf("Progress: %d/%d (%.0f%%)\n", done, total, completion))
	sb.WriteString(fmt.Sprintf("  Done: %d | In Progress: %d | Todo: %d | Blocked: %d\n\n", done, inProgress, todo, blocked))

	if len(snaps) > 0 {
		sb.WriteString(fmt.Sprintf("Historical avg completion: %.0f%%\n", avgCompletion))
		if completion < avgCompletion-10 {
			sb.WriteString("BELOW AVERAGE - team is behind typical pace\n")
		} else if completion > avgCompletion {
			sb.WriteString("ABOVE AVERAGE - team tracking well\n")
		}
		sb.WriteString("\n")
	}

	if blocked > 0 {
		sb.WriteString(fmt.Sprintf("Blockers impacting delivery: %d items (%.0f%% of sprint)\n", blocked, float64(blocked)/float64(total)*100))
		sb.WriteString("Action: See impediment escalation brief for details\n\n")
	}

	sb.WriteString("Recommendation: ")
	switch confidence {
	case "HIGH":
		sb.WriteString("Sprint goal achievable. No intervention needed.\n")
	case "MEDIUM":
		sb.WriteString("Sprint goal at risk. Monitor daily. Consider scope reduction if no progress in 2 days.\n")
	case "LOW":
		sb.WriteString("Sprint goal unlikely without intervention. Recommend: remove blocked items, reduce scope, or extend timeline.\n")
	}

	return textResult(sb.String()), nil
}

// ResourcePlanningReport shows team capacity and throughput trends for management resource planning.
func (h *Handlers) ResourcePlanningReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 6)

	var sb strings.Builder
	sb.WriteString("RESOURCE & CAPACITY REPORT\n\n")

	if len(snaps) >= 3 {
		totalDone := 0
		totalCarryover := 0
		for _, s := range snaps {
			totalDone += s.Done
			totalCarryover += s.Carryover
		}
		avgDone := totalDone / len(snaps)
		avgCarryover := totalCarryover / len(snaps)
		utilizationPct := 0.0
		if avgDone > 0 {
			utilizationPct = float64(avgDone-avgCarryover) / float64(avgDone) * 100
		}

		sb.WriteString(fmt.Sprintf("Team Throughput (last %d sprints):\n", len(snaps)))
		sb.WriteString(fmt.Sprintf("  Average items delivered: %d/sprint\n", avgDone))
		sb.WriteString(fmt.Sprintf("  Average carryover: %d/sprint\n", avgCarryover))
		sb.WriteString(fmt.Sprintf("  Effective delivery rate: %.0f%%\n\n", utilizationPct))

		if len(snaps) >= 4 {
			recent := (snaps[0].Done + snaps[1].Done) / 2
			older := (snaps[len(snaps)-2].Done + snaps[len(snaps)-1].Done) / 2
			if recent > older {
				sb.WriteString("Trend: INCREASING capacity (team ramping up)\n")
			} else if recent < older {
				sb.WriteString("Trend: DECREASING capacity (investigate: attrition? complexity? blockers?)\n")
			} else {
				sb.WriteString("Trend: STABLE\n")
			}
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("Insufficient historical data (need 3+ sprint snapshots)\n\n")
	}

	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) > 0 {
		issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
		assigneeLoad := map[string]int{}
		for _, i := range issues {
			if i.Assignee != "" {
				assigneeLoad[i.Assignee]++
			}
		}
		sb.WriteString(fmt.Sprintf("Current Sprint Workload (%d items, %d people):\n", len(issues), len(assigneeLoad)))
		for person, count := range assigneeLoad {
			indicator := ""
			if count > 7 {
				indicator = " [OVERLOADED]"
			}
			sb.WriteString(fmt.Sprintf("  %s: %d items%s\n", person, count, indicator))
		}
	}

	return textResult(sb.String()), nil
}
