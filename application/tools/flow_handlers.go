package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/mark3labs/mcp-go/mcp"
)

// FlowMetrics calculates cycle time, throughput, WIP, and lead time from Jira.
func (h *Handlers) FlowMetrics(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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
		return sanitizedError("failed to get flow issues", err), nil
	}

	now := time.Now()
	var totalCycleTime float64
	var doneCount, wipCount, todoCount int

	for _, issue := range issues {
		lower := strings.ToLower(issue.Status)
		switch {
		case strings.Contains(lower, "done") || strings.Contains(lower, "closed") || strings.Contains(lower, "resolved"):
			doneCount++
			// Cycle time = created → done (approximation using updated as done date)
			cycleTime := issue.Updated.Sub(issue.Created).Hours() / 24
			totalCycleTime += cycleTime
		case strings.Contains(lower, "progress") || strings.Contains(lower, "review") || strings.Contains(lower, "dev"):
			wipCount++
		default:
			if !strings.Contains(lower, "block") {
				todoCount++
			} else {
				wipCount++ // blocked counts as WIP
			}
		}
	}

	// Calculate metrics
	avgCycleTime := 0.0
	if doneCount > 0 {
		avgCycleTime = totalCycleTime / float64(doneCount)
	}

	// Throughput: items done per day since sprint started
	sprintDays := now.Sub(issues[0].Created).Hours() / 24
	if sprintDays < 1 {
		sprintDays = 1
	}
	throughput := float64(doneCount) / sprintDays

	// Lead time approximation using Little's Law: Lead Time = WIP / Throughput
	leadTime := 0.0
	if throughput > 0 {
		leadTime = float64(wipCount) / throughput
	}

	// Flow efficiency: active work time vs total time (approximation)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Flow Metrics — Sprint: %s\n\n", sprint.Name))
	sb.WriteString(fmt.Sprintf("WIP (Work In Progress): %d items\n", wipCount))
	sb.WriteString(fmt.Sprintf("Throughput: %.1f items/day\n", throughput))
	sb.WriteString(fmt.Sprintf("Avg Cycle Time: %.1f days\n", avgCycleTime))
	sb.WriteString(fmt.Sprintf("Estimated Lead Time: %.1f days (Little's Law)\n", leadTime))
	sb.WriteString(fmt.Sprintf("\nDone: %d | WIP: %d | Todo: %d | Total: %d\n", doneCount, wipCount, todoCount, len(issues)))

	// Recommendations based on metrics
	sb.WriteString("\nSignals:\n")
	if wipCount > doneCount && wipCount > 5 {
		sb.WriteString("  - HIGH WIP: Too much in progress. Focus on finishing over starting.\n")
	}
	if avgCycleTime > 10 {
		sb.WriteString("  - LONG CYCLE TIME: Stories may be too large. Consider slicing.\n")
	}
	if throughput < 0.5 {
		sb.WriteString("  - LOW THROUGHPUT: Team is delivering slowly. Check for blockers.\n")
	}
	if wipCount <= 3 && throughput >= 1 {
		sb.WriteString("  - GOOD FLOW: Low WIP + good throughput. Keep it up.\n")
	}

	return textResult(sb.String()), nil
}

// SprintComparison compares current sprint against the previous one.
func (h *Handlers) SprintComparison(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	// Get last 2 snapshots
	snaps, err := h.Memory.GetSprintSnapshots(ctx, boardID, 2)
	if err != nil || len(snaps) < 2 {
		return textResult("Need at least 2 sprint snapshots for comparison. Capture more with pm_snapshot_sprint."), nil
	}

	current := snaps[0]
	previous := snaps[1]

	var sb strings.Builder
	sb.WriteString("Sprint Comparison\n\n")
	sb.WriteString(fmt.Sprintf("%-20s | %-15s | %-15s | Change\n", "Metric", current.SprintName, previous.SprintName))
	sb.WriteString(strings.Repeat("-", 70) + "\n")

	sb.WriteString(formatComparison("Total Issues", current.TotalIssues, previous.TotalIssues))
	sb.WriteString(formatComparison("Done", current.Done, previous.Done))
	sb.WriteString(formatComparison("Blocked", current.Blocked, previous.Blocked))
	sb.WriteString(formatComparison("Carryover", current.Carryover, previous.Carryover))
	sb.WriteString(formatComparison("Velocity", current.Velocity, previous.Velocity))
	sb.WriteString(fmt.Sprintf("%-20s | %-15.0f | %-15.0f | %+.0f%%\n", "Completion %",
		current.CompletionRate, previous.CompletionRate, current.CompletionRate-previous.CompletionRate))

	// Verdict
	sb.WriteString("\nVerdict: ")
	improvements := 0
	if current.Velocity > previous.Velocity {
		improvements++
	}
	if current.CompletionRate > previous.CompletionRate {
		improvements++
	}
	if current.Blocked < previous.Blocked {
		improvements++
	}
	if current.Carryover < previous.Carryover {
		improvements++
	}

	switch {
	case improvements >= 3:
		sb.WriteString("IMPROVING — team is getting better across multiple metrics\n")
	case improvements >= 2:
		sb.WriteString("STABLE — some improvements, some regressions\n")
	default:
		sb.WriteString("DECLINING — multiple metrics worse than previous sprint\n")
	}

	return textResult(sb.String()), nil
}

func formatComparison(metric string, current, previous int) string {
	diff := current - previous
	sign := ""
	if diff > 0 {
		sign = "+"
	}
	return fmt.Sprintf("%-20s | %-15d | %-15d | %s%d\n", metric, current, previous, sign, diff)
}

// CeremonyFacilitator provides AI-powered facilitation prompts for Scrum ceremonies.
func (h *Handlers) CeremonyFacilitator(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ceremony, err := req.RequireString("ceremony")
	if err != nil {
		return errorResult("ceremony required (standup, planning, retro, grooming, review)"), nil
	}

	boardID := req.GetInt("board_id", 0)
	var contextData strings.Builder

	// Gather context
	if boardID > 0 {
		sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
		if len(sprints) > 0 {
			issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
			contextData.WriteString(fmt.Sprintf("Sprint: %s, Goal: %s, Issues: %d\n", sprints[0].Name, sprints[0].Goal, len(issues)))
		}
	}
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) > 0 {
		contextData.WriteString(fmt.Sprintf("Active blockers: %d\n", len(blockers)))
	}
	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 0 {
		contextData.WriteString(fmt.Sprintf("Pending retro actions: %d\n", len(actions)))
	}

	prompts := map[string]string{
		"standup": `Generate a focused daily standup facilitation guide. Include:
1. Opening (1 min): Energy check, sprint goal reminder
2. Round-robin prompts: Instead of boring "what did you do", suggest focused questions like:
   - "What's the ONE thing that will move the sprint goal forward today?"
   - "Who needs help? What's blocking you RIGHT NOW?"
   - "Any dependency that might surprise us this week?"
3. Parking lot: topics to take offline
4. Closing (30 sec): confirm actions

Keep it conversational, not robotic. Under 150 words.`,

		"planning": `Generate a sprint planning facilitation guide. Include:
1. Review previous sprint outcomes
2. Capacity check questions
3. Goal-setting prompts (make it SMART)
4. Backlog grooming questions to ask before committing
5. Commitment confidence vote prompt
6. Risk identification round
7. Definition of Done reminder

Be practical, not textbook. Under 200 words.`,

		"retro": `Generate a fresh retrospective facilitation approach. Do NOT use boring "what went well / what didn't".
Instead, pick ONE of these formats randomly:
- "Start/Stop/Continue" 
- "Sailboat" (wind, anchor, rocks, island)
- "4Ls" (Liked, Learned, Lacked, Longed for)
- "Hot Air Balloon" (what lifts us, what weighs us down)
- "Timeline" (plot the sprint as a story)

Provide the full facilitation script with:
1. Icebreaker (2 min)
2. Data gathering (10 min)
3. Generate insights (10 min)
4. Decide actions (5 min)
5. Close with appreciation

Include specific prompts/questions for each step. Under 250 words.`,

		"grooming": `Generate backlog grooming facilitation prompts:
1. Priority validation questions
2. Story readiness checklist (INVEST criteria check)
3. Estimation facilitation (planning poker tips)
4. Splitting large stories prompts
5. Acceptance criteria review questions
6. Dependency identification round

Practical, conversational. Under 150 words.`,

		"review": `Generate sprint review facilitation guide:
1. Demo order suggestions
2. Stakeholder engagement questions
3. Feedback capture framework
4. "So what?" connection (how does this delivered work connect to business goals)
5. Next sprint preview teaser

Keep it engaging for stakeholders, not just devs. Under 150 words.`,
	}

	systemPrompt, ok := prompts[ceremony]
	if !ok {
		return errorResult("Unknown ceremony. Use: standup, planning, retro, grooming, review"), nil
	}

	userPrompt := "Team context:\n" + contextData.String()
	if contextData.Len() == 0 {
		userPrompt = "No specific context available. Generate generic but useful facilitation prompts."
	}

	result, err := h.aiComplete(ctx, systemPrompt, userPrompt)
	if err != nil {
		return sanitizedError("ai analysis failed in flow", err), nil
	}

	return textResult(result), nil
}

// RecordConfidence tracks team confidence before/during sprint.
func (h *Handlers) RecordConfidence(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintName, err := req.RequireString("sprint_name")
	if err != nil {
		return errorResult("sprint_name required"), nil
	}
	score := req.GetInt("score", 0)
	if score < 1 || score > 5 {
		return errorResult("score required (1-5): 1=very worried, 3=neutral, 5=very confident"), nil
	}

	note := req.GetString("note", "")
	member := req.GetString("member", "team")

	m := &memdom.TeamMetric{
		MemberName:     member,
		SprintName:     sprintName,
		RecordedAt:     time.Now(),
		IssuesAssigned: score, // repurpose for confidence score
		Notes:          fmt.Sprintf("confidence:%d %s", score, note),
	}

	if err := h.Memory.SaveTeamMetric(ctx, m); err != nil {
		return sanitizedError("failed to save flow data", err), nil
	}

	labels := []string{"", "very worried", "worried", "neutral", "confident", "very confident"}
	return textResult(fmt.Sprintf("Confidence recorded: %s gives %d/5 (%s) for sprint %s\nNote: %s",
		member, score, labels[score], sprintName, note)), nil
}

// SprintGoalCheck evaluates whether sprint goal is being met based on data.
func (h *Handlers) SprintGoalCheck(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	// Get active goals
	goals, err := h.Memory.GetActiveGoals(ctx, boardID)
	if err != nil || len(goals) == 0 {
		return textResult("No active sprint goals. Set one with pm_set_sprint_goal."), nil
	}

	// Get sprint data
	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	var sprintInfo string
	if len(sprints) > 0 {
		issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
		var done, total int
		total = len(issues)
		for _, issue := range issues {
			lower := strings.ToLower(issue.Status)
			if strings.Contains(lower, "done") || strings.Contains(lower, "closed") {
				done++
			}
		}
		sprintInfo = fmt.Sprintf("Sprint: %s | Progress: %d/%d done (%.0f%%)", sprints[0].Name, done, total, float64(done)/float64(total)*100)
	}

	var sb strings.Builder
	sb.WriteString("Sprint Goal Status:\n\n")
	sb.WriteString(sprintInfo + "\n\n")

	for _, g := range goals {
		sb.WriteString(fmt.Sprintf("Goal: %s\n", g.Goal))
		if g.KeyResults != "" {
			sb.WriteString(fmt.Sprintf("Key Results:\n%s\n", g.KeyResults))
		}
		sb.WriteString("\n")
	}

	// AI assessment
	systemPrompt := `You are a Scrum Master evaluating sprint goal progress.
Based on the sprint data and goal definition, assess:
1. On track / At risk / Behind
2. What evidence supports your assessment
3. What needs to happen to achieve the goal
Keep it to 3-4 bullet points. Be honest.`

	analysis, err := h.aiComplete(ctx, systemPrompt, sb.String())
	if err != nil {
		return textResult(sb.String()), nil // fallback to raw data
	}

	return textResult(sb.String() + "\nAssessment:\n" + analysis), nil
}
