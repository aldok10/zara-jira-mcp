package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// PMBurnoutRisk calculates per-person burnout risk from Jira signals.
// Signals: WIP overload, carryover pattern, assignment concentration, blocker accumulation.
func (h *Handlers) PMBurnoutRisk(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil || len(sprints) == 0 {
		return errorResult("no active sprint"), nil
	}
	issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
	if len(issues) == 0 {
		return textResult("No sprint issues to analyze."), nil
	}

	// Collect per-person signals
	type personData struct {
		assigned int
		inProgress int
		blocked  int
		done     int
		carryover int
	}
	people := map[string]*personData{}
	for _, issue := range issues {
		if issue.Assignee == "" {
			continue
		}
		p, ok := people[issue.Assignee]
		if !ok {
			p = &personData{}
			people[issue.Assignee] = p
		}
		p.assigned++
		lower := strings.ToLower(issue.Status)
		if strings.Contains(lower, "progress") || strings.Contains(lower, "review") {
			p.inProgress++
		} else if strings.Contains(lower, "done") || strings.Contains(lower, "closed") {
			p.done++
		} else if strings.Contains(lower, "block") {
			p.blocked++
		}
	}

	// Check carryover from memory
	if h.Memory != nil {
		metrics, _ := h.Memory.GetTeamOverview(ctx, "")
		for _, m := range metrics {
			if p, ok := people[m.MemberName]; ok {
				p.carryover = m.CarryoverCount
			}
		}
	}

	// Calculate risk score per person (0-100)
	type risk struct {
		name    string
		score   int
		signals []string
	}
	var risks []risk
	teamSize := len(people)
	avgAssigned := float64(len(issues)) / float64(max(teamSize, 1))

	for name, p := range people {
		score := 0
		var signals []string

		// WIP overload (>2 items in progress)
		if p.inProgress > 2 {
			score += 25
			signals = append(signals, fmt.Sprintf("WIP=%d (context switching)", p.inProgress))
		}

		// Assignment concentration (>1.5x team avg)
		if float64(p.assigned) > avgAssigned*1.5 {
			score += 25
			signals = append(signals, fmt.Sprintf("%d items (team avg %.0f)", p.assigned, avgAssigned))
		}

		// Blocked items
		if p.blocked >= 2 {
			score += 20
			signals = append(signals, fmt.Sprintf("%d blocked items (frustration)", p.blocked))
		}

		// Carryover pattern
		if p.carryover >= 2 {
			score += 20
			signals = append(signals, fmt.Sprintf("%d carryover items (chronic overcommit)", p.carryover))
		}

		// Low completion ratio
		if p.assigned > 3 && p.done == 0 {
			score += 10
			signals = append(signals, "no items completed yet")
		}

		if score > 100 {
			score = 100
		}
		if score > 0 {
			risks = append(risks, risk{name, score, signals})
		}
	}

	if len(risks) == 0 {
		return textResult("No burnout risk signals detected. Workload looks balanced."), nil
	}

	// Sort by score desc
	for i := range risks {
		for j := i + 1; j < len(risks); j++ {
			if risks[j].score > risks[i].score {
				risks[i], risks[j] = risks[j], risks[i]
			}
		}
	}

	var sb strings.Builder
	sb.WriteString("Burnout Risk Assessment\n\n")
	for _, r := range risks {
		level := "LOW"
		if r.score >= 60 {
			level = "HIGH"
		} else if r.score >= 35 {
			level = "MEDIUM"
		}
		sb.WriteString(fmt.Sprintf("%s: %d/100 [%s]\n", r.name, r.score, level))
		for _, s := range r.signals {
			sb.WriteString(fmt.Sprintf("  - %s\n", s))
		}
	}

	high := 0
	for _, r := range risks {
		if r.score >= 60 {
			high++
		}
	}
	if high > 0 {
		sb.WriteString(fmt.Sprintf("\n%d person(s) at HIGH risk. Consider redistributing work or reducing scope.", high))
	}
	return textResult(sb.String()), nil
}

// PMRealityCheck compares perceived velocity (points/items done) vs actual delivery (releases, customer impact).
func (h *Handlers) PMRealityCheck(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}

	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
	if len(snaps) < 2 {
		return textResult("Need at least 2 sprint snapshots for comparison."), nil
	}

	// Calculate velocity trend (items done)
	var velocities []int
	for _, s := range snaps {
		velocities = append(velocities, s.Done)
	}

	// Check if velocity up but carryover also up (activity theater)
	latestVel := velocities[0]
	olderVel := velocities[len(velocities)-1]
	latestCarry := snaps[0].Carryover
	olderCarry := snaps[len(snaps)-1].Carryover

	var sb strings.Builder
	sb.WriteString("Reality Check: Activity vs Delivery\n\n")
	sb.WriteString(fmt.Sprintf("Velocity (items done): %d -> %d", olderVel, latestVel))
	if latestVel > olderVel {
		sb.WriteString(" (up)\n")
	} else if latestVel < olderVel {
		sb.WriteString(" (down)\n")
	} else {
		sb.WriteString(" (flat)\n")
	}

	sb.WriteString(fmt.Sprintf("Carryover: %d -> %d", olderCarry, latestCarry))
	if latestCarry > olderCarry {
		sb.WriteString(" (up — overcommitting)\n")
	} else {
		sb.WriteString(" (stable/down)\n")
	}

	// Predictability = done / (done + carryover)
	if snaps[0].TotalIssues > 0 {
		predictability := float64(snaps[0].Done) / float64(snaps[0].TotalIssues) * 100
		sb.WriteString(fmt.Sprintf("Predictability: %.0f%%", predictability))
		if predictability < 70 {
			sb.WriteString(" (below 70% — commitments unreliable)\n")
		} else {
			sb.WriteString(" (healthy)\n")
		}
	}

	// Detect "productivity theater": velocity up but carryover also up
	if latestVel > olderVel && latestCarry > olderCarry {
		sb.WriteString("\nWARNING: Velocity up but carryover also up. Team is doing more but finishing less. This is the 'AI productivity illusion' — activity without delivery.")
	} else if latestVel > olderVel && latestCarry <= olderCarry {
		sb.WriteString("\nHEALTHY: Velocity up AND carryover down. Real delivery improvement.")
	}

	// Goal hit rate from memory
	goals, _ := h.Memory.GetGoalHistory(ctx, boardID, 5)
	if len(goals) > 0 {
		hit := 0
		for _, g := range goals {
			if g.Status == "achieved" || g.Status == "hit" {
				hit++
			}
		}
		sb.WriteString(fmt.Sprintf("\n\nGoal hit rate: %d/%d (%.0f%%)", hit, len(goals), float64(hit)/float64(len(goals))*100))
		if hit < len(goals)/2 {
			sb.WriteString(" — below 50%. Volume isn't value.")
		}
	}

	return textResult(sb.String()), nil
}

// PMSafetySignals detects psychological safety from observable behavioral data.
// Based on 7-category behavioral framework (Edmondson/Nembhard 2020).
func (h *Handlers) PMSafetySignals(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}

	var signals []string
	score := 50 // Start neutral

	// Signal 1: Retro participation (silence = fear)
	retros, _ := h.Memory.GetRetrospectives(ctx, 3)
	emptyRetros := 0
	for _, r := range retros {
		if strings.TrimSpace(r.Improvements) == "" {
			emptyRetros++
		}
	}
	if emptyRetros >= 2 {
		score -= 15
		signals = append(signals, "Retros have empty 'improvements' — team may not feel safe to critique")
	} else if len(retros) >= 2 && emptyRetros == 0 {
		score += 10
		signals = append(signals, "Retros have rich content — team speaks up")
	}

	// Signal 2: Bug reporting distribution (concentrated = others afraid to report)
	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) > 0 {
		issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
		bugReporters := map[string]int{}
		for _, i := range issues {
			if strings.ToLower(i.Type) == "bug" && i.Reporter != "" {
				bugReporters[i.Reporter]++
			}
		}
		if len(bugReporters) == 1 && len(issues) > 5 {
			score -= 10
			signals = append(signals, "Only 1 person reports bugs — others may fear blame")
		} else if len(bugReporters) >= 3 {
			score += 5
			signals = append(signals, "Multiple bug reporters — healthy error culture")
		}
	}

	// Signal 3: Action item diversity (same person always = no shared accountability)
	pending, _ := h.Memory.GetPendingActionItems(ctx)
	owners := map[string]int{}
	for _, a := range pending {
		owners[a.Owner]++
	}
	if len(pending) > 3 && len(owners) == 1 {
		score -= 10
		signals = append(signals, "One person owns all action items — no shared accountability")
	}

	// Signal 4: Blocker reporting by new members
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) == 0 && len(pending) > 5 {
		score -= 5
		signals = append(signals, "Zero blockers reported with 10+ items — team may not ask for help")
	} else if len(blockers) > 0 {
		score += 5
		signals = append(signals, "Blockers being reported — team comfortable raising issues")
	}

	// Signal 5: Decision participation (no dissent = groupthink)
	decisions, _ := h.Memory.GetDecisions(ctx, 10)
	if len(decisions) >= 3 {
		hasAlternatives := 0
		for _, d := range decisions {
			if strings.Contains(d.Context, "alternative") || strings.Contains(d.Context, "option") {
				hasAlternatives++
			}
		}
		if hasAlternatives == 0 {
			score -= 10
			signals = append(signals, "No alternatives recorded in decisions — possible groupthink")
		}
	}

	// Clamp score
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}

	var sb strings.Builder
	level := "HEALTHY"
	if score < 35 {
		level = "AT RISK"
	} else if score < 55 {
		level = "WATCH"
	}
	sb.WriteString(fmt.Sprintf("Psychological Safety: %d/100 [%s]\n\n", score, level))

	sb.WriteString("Behavioral Signals:\n")
	for _, s := range signals {
		sb.WriteString(fmt.Sprintf("- %s\n", s))
	}

	if score < 50 && h.AI != nil {
		suggestion, err := h.AI.Complete(ctx, EmpathySystemPrompt+"\nGiven these low psychological safety signals, suggest ONE small action the SM can take in the next ceremony. Be specific and practical. Max 2 sentences.", strings.Join(signals, "\n"))
		if err == nil {
			sb.WriteString("\nSuggestion: " + strings.TrimSpace(suggestion) + "\n")
		}
	}

	return textResult(sb.String()), nil
}
