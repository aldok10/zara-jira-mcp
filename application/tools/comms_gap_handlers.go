package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// CommsHealth calculates a communication health score (0-100) across 4 dimensions.
func (h *Handlers) CommsHealth(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	boardID := req.GetInt("board_id", 0)

	var sb strings.Builder
	total := 0

	// Decision velocity: avg decisions per week in last 30 days
	decisions, _ := h.Memory.GetDecisions(ctx, 30)
	decisionScore := 5
	if len(decisions) > 0 {
		cutoff := time.Now().AddDate(0, 0, -30)
		recent := 0
		for _, d := range decisions {
			if d.MadeAt.After(cutoff) {
				recent++
			}
		}
		perWeek := float64(recent) / 4.0
		if perWeek >= 1.0 {
			decisionScore = 25
		} else if perWeek >= 0.5 {
			decisionScore = 15
		}
	}
	total += decisionScore
	sb.WriteString(fmt.Sprintf("Decision Velocity: %d/25", decisionScore))
	if decisionScore == 25 {
		sb.WriteString(" (healthy)\n")
	} else {
		sb.WriteString(" (slow)\n")
	}

	// Blocker resolution: avg days to resolve
	blockerScore := 5
	blockers, _ := h.Memory.GetBlockerHistory(ctx, 30)
	if len(blockers) > 0 {
		resolved := 0
		totalDays := 0.0
		for _, b := range blockers {
			if b.ResolvedAt != nil {
				resolved++
				totalDays += b.ResolvedAt.Sub(b.BlockedSince).Hours() / 24
			}
		}
		if resolved > 0 {
			avgDays := totalDays / float64(resolved)
			if avgDays < 2 {
				blockerScore = 25
			} else if avgDays < 5 {
				blockerScore = 15
			} else if avgDays < 10 {
				blockerScore = 10
			}
		}
	}
	total += blockerScore
	sb.WriteString(fmt.Sprintf("Blocker Resolution: %d/25", blockerScore))
	if blockerScore >= 15 {
		sb.WriteString(" (good)\n")
	} else {
		sb.WriteString(" (slow)\n")
	}

	// Action follow-through: completed vs total
	actionScore := 5
	pending, _ := h.Memory.GetPendingActionItems(ctx)
	// Query total action items via raw DB
	totalActions := len(pending)
	completedActions := 0
	if h.Memory.DB() != nil {
		row := h.Memory.DB().QueryRow("SELECT COUNT(*) FROM action_items WHERE status='done'")
		_ = row.Scan(&completedActions)
		totalActions += completedActions
	}
	if totalActions > 0 {
		ratio := float64(completedActions) / float64(totalActions)
		if ratio > 0.8 {
			actionScore = 25
		} else if ratio > 0.6 {
			actionScore = 15
		} else if ratio > 0.4 {
			actionScore = 10
		}
	}
	total += actionScore
	sb.WriteString(fmt.Sprintf("Action Follow-through: %d/25", actionScore))
	if totalActions > 0 {
		sb.WriteString(fmt.Sprintf(" (%d/%d done)\n", completedActions, totalActions))
	} else {
		sb.WriteString(" (no data)\n")
	}

	// Stakeholder engagement: pulse entries in last 60 days
	engagementScore := 0
	pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 100)
	cutoff60 := time.Now().AddDate(0, 0, -60)
	recentPulses := 0
	for _, p := range pulses {
		if p.CreatedAt.After(cutoff60) {
			recentPulses++
		}
	}
	if recentPulses > 5 {
		engagementScore = 25
	} else if recentPulses > 2 {
		engagementScore = 15
	} else if recentPulses > 0 {
		engagementScore = 10
	}
	total += engagementScore
	sb.WriteString(fmt.Sprintf("Stakeholder Engagement: %d/25 (%d entries)\n", engagementScore, recentPulses))

	// Summary
	_ = boardID // included for future board-scoped filtering
	summary := "Critical - communication infrastructure missing"
	if total >= 80 {
		summary = "Healthy communication patterns"
	} else if total >= 60 {
		summary = "Decent but gaps exist"
	} else if total >= 40 {
		summary = "Below average - address gaps"
	}

	return textResult(fmt.Sprintf("Communication Health: %d/100\n\n%s\nSummary: %s", total, sb.String(), summary)), nil
}

// SilenceDetector identifies stakeholders with no recent pulse activity.
func (h *Handlers) SilenceDetector(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	threshold := req.GetInt("days_threshold", 30)

	pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 100)
	cutoff := time.Now().AddDate(0, 0, -threshold)

	// Track last seen per member
	lastSeen := map[string]time.Time{}
	for _, p := range pulses {
		if t, ok := lastSeen[p.Member]; !ok || p.CreatedAt.After(t) {
			lastSeen[p.Member] = p.CreatedAt
		}
	}

	if len(lastSeen) == 0 {
		return textResult("No pulse data available. Record team pulse first."), nil
	}

	var silent []string
	for member, last := range lastSeen {
		if last.Before(cutoff) {
			days := int(time.Since(last).Hours() / 24)
			silent = append(silent, fmt.Sprintf("- %s (last seen: %s, %d days ago)", member, last.Format("2006-01-02"), days))
		}
	}

	if len(silent) == 0 {
		return textResult(fmt.Sprintf("All %d stakeholders active within last %d days.", len(lastSeen), threshold)), nil
	}

	return textResult(fmt.Sprintf("Silent Stakeholders (%d found, threshold: %d days):\n%s", len(silent), threshold, strings.Join(silent, "\n"))), nil
}

// CommsAntiPatterns detects communication anti-patterns from historical data.
func (h *Handlers) CommsAntiPatterns(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	_ = req.GetInt("board_id", 0)

	var patterns []string

	// 1. Re-deciding: same keywords in decisions appearing 3+ times
	decisions, _ := h.Memory.GetDecisions(ctx, 50)
	titleWords := map[string]int{}
	for _, d := range decisions {
		for _, w := range strings.Fields(strings.ToLower(d.Title)) {
			if len(w) > 4 {
				titleWords[w]++
			}
		}
	}
	for word, count := range titleWords {
		if count >= 3 {
			patterns = append(patterns, fmt.Sprintf("[HIGH] Re-deciding: '%s' appears in %d decisions. Team revisiting same topics = signal loss. Fix: record decisions with rationale and refer back.", word, count))
			break
		}
	}

	// 2. Escalation hoarding: blockers but no escalations
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	escalations, _ := h.Memory.GetRecentEscalations(ctx, 10)
	agingBlockers := 0
	for _, b := range blockers {
		if time.Since(b.BlockedSince).Hours()/24 > 5 {
			agingBlockers++
		}
	}
	if agingBlockers >= 3 && len(escalations) == 0 {
		patterns = append(patterns, fmt.Sprintf("[HIGH] Escalation Hoarding: %d blockers aging >5 days but 0 escalations. Team absorbing pain silently. Fix: escalate blockers >3 days automatically.", agingBlockers))
	}

	// 4. Ghost stakeholders: no pulse data at all
	pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 10)
	if len(pulses) == 0 {
		patterns = append(patterns, "[MEDIUM] Ghost Stakeholders: No pulse data recorded. No feedback loop with the team. Fix: run team pulse at end of each sprint.")
	}

	// 5. Blocker silence: active blockers aging with no resolution
	for _, b := range blockers {
		days := int(time.Since(b.BlockedSince).Hours() / 24)
		if days > 7 && b.Resolution == "" {
			patterns = append(patterns, fmt.Sprintf("[HIGH] Blocker Silence: '%s' blocked %d days with no resolution attempt. Fix: daily check on blockers, assign explicit owner + deadline.", b.Description, days))
			break
		}
	}

	if len(patterns) == 0 {
		return textResult("No communication anti-patterns detected. Team communication looks healthy."), nil
	}

	return textResult(fmt.Sprintf("Communication Anti-Patterns (%d detected):\n\n%s", len(patterns), strings.Join(patterns, "\n\n"))), nil
}

// NVCReframe rewrites a message using Nonviolent Communication (Observation, Feeling, Need, Request).
func (h *Handlers) NVCReframe(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	message, err := req.RequireString("message")
	if err != nil {
		return errorResult("message required"), nil
	}

	if h.AI == nil {
		// Graceful degradation: return template
		return textResult(fmt.Sprintf("Original: %s\n\nNVC Reframe (template - AI unavailable):\n- Observation: When I notice [specific behavior]...\n- Feeling: I feel [emotion]...\n- Need: Because I need [value/need]...\n- Request: Would you be willing to [specific action]?", message)), nil
	}

	systemPrompt := `Reframe the given message using Nonviolent Communication (NVC) by Marshall Rosenberg.
Output format:
ORIGINAL: [the input message]

NVC REFRAME:
- Observation: [factual, no judgment]
- Feeling: [emotion word, not "I feel that..."]
- Need: [universal need behind it]
- Request: [specific, doable, positive language]

ONE-LINER VERSION: [condensed NVC version in one sentence]

Rules: Keep the core meaning. Remove blame, judgment, demands. Be specific.`

	result, aiErr := h.AI.Complete(ctx, systemPrompt, message)
	if aiErr != nil {
		return textResult(fmt.Sprintf("Original: %s\n\nNVC Reframe (template - AI error):\n- Observation: When I notice [specific behavior]...\n- Feeling: I feel [emotion]...\n- Need: Because I need [value/need]...\n- Request: Would you be willing to [specific action]?", message)), nil
	}
	return textResult(result), nil
}

// HardConversation prepares for a crucial conversation using STATE path + SBI + SCARF frameworks.
func (h *Handlers) HardConversation(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	situation, err := req.RequireString("situation")
	if err != nil {
		return errorResult("situation required (describe what's happening)"), nil
	}
	boardID := req.GetInt("board_id", 0)
	person := req.GetString("person", "")

	var dataContext strings.Builder
	dataContext.WriteString(fmt.Sprintf("Situation: %s\n", situation))
	if person != "" {
		dataContext.WriteString(fmt.Sprintf("Person: %s\n", person))
	}

	if h.Memory != nil {
		if boardID > 0 {
			scores, _ := h.Memory.GetHealthScores(ctx, boardID, 3)
			if len(scores) > 0 {
				dataContext.WriteString(fmt.Sprintf("Recent health: %d/100 (%s)\n", scores[0].OverallScore, scores[0].SprintName))
			}
		}
		blockers, _ := h.Memory.GetActiveBlockers(ctx)
		if len(blockers) > 0 {
			dataContext.WriteString(fmt.Sprintf("Active blockers: %d\n", len(blockers)))
			for i, b := range blockers {
				if i >= 3 {
					break
				}
				dataContext.WriteString(fmt.Sprintf("  - %s (%d days)\n", b.Description, int(time.Since(b.BlockedSince).Hours()/24)))
			}
		}
	}

	if h.AI == nil {
		return textResult(fmt.Sprintf("Crucial Conversation Prep (AI unavailable)\n\nContext:\n%s\n\nFramework:\n1. STATE: Share facts, Tell your story, Ask for their path, Talk tentatively, Encourage testing\n2. SCARF risks: Status, Certainty, Autonomy, Relatedness, Fairness\n3. SBI: Situation, Behavior, Impact\n\nFill in manually.", dataContext.String())), nil
	}

	systemPrompt := `Prepare a crucial conversation using these frameworks:

1. FACTS (from data provided): List only verifiable facts, no interpretation
2. STORIES (interpretations): 2-3 possible explanations for the behavior
3. SCARF RISKS: Which of Status/Certainty/Autonomy/Relatedness/Fairness might be threatened?
4. OPENING LINES (STATE path): 3 options for opening the conversation - share facts tentatively
5. SAFETY RESTORATION: If they get defensive, how to restore safety

Rules:
- Start with facts, not conclusions
- Use tentative language ("I'm wondering if..." not "You always...")
- Separate behavior from intent
- Keep each section under 50 words
- Be direct and practical`

	result, aiErr := h.AI.Complete(ctx, systemPrompt, dataContext.String())
	if aiErr != nil {
		return errorResult("AI failed: " + aiErr.Error()), nil
	}
	return textResult(fmt.Sprintf("[Crucial Conversation Prep]\n\n%s", result)), nil
}

// TrustSignals calculates trust indicators from historical data.
func (h *Handlers) TrustSignals(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	boardID := req.GetInt("board_id", 0)

	var sb strings.Builder
	sb.WriteString("Trust Signal Dashboard\n\n")

	// Forecast accuracy: committed vs done in snapshots
	if boardID > 0 {
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
		if len(snaps) > 0 {
			totalCommitted := 0
			totalDone := 0
			for _, s := range snaps {
				totalCommitted += s.TotalIssues
				totalDone += s.Done
			}
			accuracy := 0.0
			if totalCommitted > 0 {
				accuracy = float64(totalDone) / float64(totalCommitted) * 100
			}
			rating := "low"
			if accuracy >= 80 {
				rating = "high"
			} else if accuracy >= 60 {
				rating = "medium"
			}
			sb.WriteString(fmt.Sprintf("Forecast Accuracy: %.0f%% [%s] (%d/%d items delivered over %d sprints)\n", accuracy, rating, totalDone, totalCommitted, len(snaps)))
		} else {
			sb.WriteString("Forecast Accuracy: no data\n")
		}
	} else {
		sb.WriteString("Forecast Accuracy: provide board_id for this metric\n")
	}

	// Escalation responsiveness: acknowledged vs total
	escalations, _ := h.Memory.GetRecentEscalations(ctx, 20)
	if len(escalations) > 0 {
		acked := 0
		for _, e := range escalations {
			if e.Acknowledged {
				acked++
			}
		}
		ratio := float64(acked) / float64(len(escalations)) * 100
		rating := "low"
		if ratio >= 80 {
			rating = "high"
		} else if ratio >= 50 {
			rating = "medium"
		}
		sb.WriteString(fmt.Sprintf("Escalation Responsiveness: %.0f%% acknowledged [%s] (%d/%d)\n", ratio, rating, acked, len(escalations)))
	} else {
		sb.WriteString("Escalation Responsiveness: no escalations recorded\n")
	}

	// Consistency: health score stability
	if boardID > 0 {
		scores, _ := h.Memory.GetHealthScores(ctx, boardID, 6)
		if len(scores) >= 3 {
			var vals []float64
			for _, s := range scores {
				vals = append(vals, float64(s.OverallScore))
			}
			mean := avg(vals)
			sd := stddev(vals, mean)
			rating := "low"
			if sd < 10 {
				rating = "high"
			} else if sd < 20 {
				rating = "medium"
			}
			sb.WriteString(fmt.Sprintf("Consistency: std dev %.1f [%s] (mean health: %.0f)\n", sd, rating, mean))
		} else {
			sb.WriteString("Consistency: insufficient health data\n")
		}
	} else {
		sb.WriteString("Consistency: provide board_id for this metric\n")
	}

	// Transparency: number of decisions recorded
	decisions, _ := h.Memory.GetDecisions(ctx, 100)
	rating := "low"
	if len(decisions) >= 20 {
		rating = "high"
	} else if len(decisions) >= 10 {
		rating = "medium"
	}
	sb.WriteString(fmt.Sprintf("Transparency: %d decisions recorded [%s]\n", len(decisions), rating))

	return textResult(sb.String()), nil
}

// LencioniDysfunction maps team data to Lencioni's 5 Dysfunctions pyramid.
func (h *Handlers) LencioniDysfunction(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}

	var findings []string

	// Level 1: Absence of Trust - no vulnerability, reassignment
	retros, _ := h.Memory.GetRetrospectives(ctx, 5)
	emptyRetros := 0
	for _, r := range retros {
		if strings.TrimSpace(r.Improvements) == "" && strings.TrimSpace(r.WentWell) == "" {
			emptyRetros++
		}
	}
	if emptyRetros >= 2 {
		findings = append(findings, "Level 1 - ABSENCE OF TRUST: Retros have empty improvements/went-well sections. Team not sharing vulnerabilities. Coaching: Start with personal histories exercise, leader goes first with vulnerability.")
	}

	// Level 2: Fear of Conflict - no improvements adopted, artificial harmony
	pending, _ := h.Memory.GetPendingActionItems(ctx)
	if len(retros) >= 3 && len(pending) == 0 {
		findings = append(findings, "Level 2 - FEAR OF CONFLICT: Multiple retros but zero pending actions. Team avoids productive disagreement. Coaching: Designate a 'miner of conflict' role, use real-time permission to disagree.")
	}

	// Level 3: Lack of Commitment - no sprint goals, high carryover
	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
	goals, _ := h.Memory.GetActiveGoals(ctx, boardID)
	if len(goals) == 0 && len(snaps) > 0 {
		findings = append(findings, "Level 3 - LACK OF COMMITMENT: No active sprint goals. Team working without clear shared commitment. Coaching: End every planning with explicit goal statement + commitment poll.")
	}
	if len(snaps) >= 3 {
		highCarryover := 0
		for _, s := range snaps[:3] {
			if s.TotalIssues > 0 && float64(s.Carryover)/float64(s.TotalIssues) > 0.3 {
				highCarryover++
			}
		}
		if highCarryover >= 2 {
			findings = append(findings, "Level 3 - LACK OF COMMITMENT: Consistent >30% carryover signals vague commitments. Coaching: Use cascading messaging - clarify decisions and deadlines at end of each meeting.")
		}
	}

	// Level 4: Avoidance of Accountability - hero culture
	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) > 0 {
		issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
		assigneeDone := map[string]int{}
		totalDone := 0
		for _, issue := range issues {
			lower := strings.ToLower(issue.Status)
			if strings.Contains(lower, "done") || strings.Contains(lower, "closed") {
				assigneeDone[issue.Assignee]++
				totalDone++
			}
		}
		if totalDone > 5 && len(assigneeDone) > 1 {
			var maxDone int
			for _, count := range assigneeDone {
				if count > maxDone {
					maxDone = count
				}
			}
			if float64(maxDone)/float64(totalDone) > 0.5 {
				findings = append(findings, "Level 4 - AVOIDANCE OF ACCOUNTABILITY: One person doing >50% of work. Peers not holding each other accountable. Coaching: Publish team goals + individual commitments, peer feedback rounds.")
			}
		}
	}

	// Level 5: Inattention to Results - missed goals
	goalHistory, _ := h.Memory.GetGoalHistory(ctx, boardID, 5)
	missed := 0
	for _, g := range goalHistory {
		if g.Status == "missed" {
			missed++
		}
	}
	if missed >= 2 {
		findings = append(findings, fmt.Sprintf("Level 5 - INATTENTION TO RESULTS: %d sprint goals missed. Team focused on activity not outcomes. Coaching: Public declaration of results, team-based rewards over individual.", missed))
	}

	if len(findings) == 0 {
		return textResult("No Lencioni dysfunction signals detected. Team appears healthy across all 5 levels."), nil
	}

	return textResult(fmt.Sprintf("Lencioni 5 Dysfunctions Analysis (%d signals):\n\n%s", len(findings), strings.Join(findings, "\n\n"))), nil
}
