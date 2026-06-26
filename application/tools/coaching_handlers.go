package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strings"

	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/mark3labs/mcp-go/mcp"
)

// PMTeamPulse records a team health pulse survey.
func (h *Handlers) PMTeamPulse(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintName, err := req.RequireString("sprint_name")
	if err != nil {
		return errorResult("sprint_name required"), nil
	}
	ratingsStr, err := req.RequireString("ratings")
	if err != nil {
		return errorResult("ratings required (JSON: {\"member\": score, ...})"), nil
	}
	notes := req.GetString("notes", "")

	var ratings map[string]int
	if err := json.Unmarshal([]byte(ratingsStr), &ratings); err != nil {
		return errorResult("ratings must be valid JSON: {\"alice\": 4, \"bob\": 3}"), nil
	}

	for member, score := range ratings {
		p := &memdom.TeamPulse{SprintName: sprintName, Member: member, Score: score, Notes: notes}
		if err := h.Memory.SaveTeamPulse(ctx, p); err != nil {
			return errorResult("Failed to save pulse: " + err.Error()), nil
		}
	}

	// Calculate average
	total := 0
	for _, s := range ratings {
		total += s
	}
	avg := float64(total) / float64(len(ratings))

	// Get previous pulse for comparison
	history, _ := h.Memory.GetTeamPulseHistory(ctx, 50)
	prevAvg := 0.0
	prevCount := 0
	for _, p := range history {
		if p.SprintName != sprintName {
			prevAvg += float64(p.Score)
			prevCount++
			if prevCount >= len(ratings) {
				break
			}
		}
	}

	var trend string
	if prevCount > 0 {
		prevAvg /= float64(prevCount)
		diff := avg - prevAvg
		if diff > 0 {
			trend = fmt.Sprintf(" (up from %.1f last pulse)", prevAvg)
		} else if diff < 0 {
			trend = fmt.Sprintf(" (down from %.1f last pulse)", prevAvg)
		} else {
			trend = fmt.Sprintf(" (unchanged from %.1f)", prevAvg)
		}
	}

	return textResult(fmt.Sprintf("Team pulse: %.1f/5%s\nRecorded %d member ratings for %s.", avg, trend, len(ratings), sprintName)), nil
}

// PMTeamPulseHistory shows pulse trends.
func (h *Handlers) PMTeamPulseHistory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	history, err := h.Memory.GetTeamPulseHistory(ctx, 100)
	if err != nil {
		return errorResult("Failed: " + err.Error()), nil
	}
	if len(history) == 0 {
		return textResult("No pulse data recorded yet. Use pm_team_pulse to record."), nil
	}

	// Group by sprint
	type sprintData struct {
		total int
		count int
	}
	sprints := map[string]*sprintData{}
	var order []string
	for _, p := range history {
		if _, ok := sprints[p.SprintName]; !ok {
			sprints[p.SprintName] = &sprintData{}
			order = append(order, p.SprintName)
		}
		sprints[p.SprintName].total += p.Score
		sprints[p.SprintName].count++
	}

	var sb strings.Builder
	sb.WriteString("Team Pulse History:\n\n")
	for _, name := range order {
		d := sprints[name]
		avg := float64(d.total) / float64(d.count)
		bar := strings.Repeat("*", int(avg))
		sb.WriteString(fmt.Sprintf("  %s: %.1f/5 %s (%d ratings)\n", name, avg, bar, d.count))
	}
	return textResult(sb.String()), nil
}

// PMPredictability calculates sprint predictability score.
func (h *Handlers) PMPredictability(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	snaps, err := h.Memory.GetSprintSnapshots(ctx, boardID, 10)
	if err != nil || len(snaps) < 3 {
		return textResult("Need at least 3 sprint snapshots. Use pm_snapshot_sprint to record sprint data."), nil
	}

	var velocities []float64
	for _, s := range snaps {
		if s.Velocity > 0 {
			velocities = append(velocities, float64(s.Velocity))
		}
	}
	if len(velocities) < 3 {
		return textResult("Need at least 3 sprints with velocity data."), nil
	}

	mean := 0.0
	for _, v := range velocities {
		mean += v
	}
	mean /= float64(len(velocities))

	variance := 0.0
	for _, v := range velocities {
		variance += (v - mean) * (v - mean)
	}
	variance /= float64(len(velocities))
	sd := math.Sqrt(variance)

	score := int(100 - (sd/mean)*100)
	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Predictability: %d/100\n", score))
	sb.WriteString(fmt.Sprintf("Your team delivers within ±%d%% of commitment.\n\n", 100-score))
	sb.WriteString(fmt.Sprintf("Mean velocity: %.1f | Std dev: %.1f | Samples: %d\n", mean, sd, len(velocities)))

	if score >= 80 {
		sb.WriteString("\nVerdict: Highly predictable. Stakeholders can trust commitments.")
	} else if score >= 60 {
		sb.WriteString("\nVerdict: Moderately predictable. Some variance exists.")
	} else {
		sb.WriteString("\nVerdict: Low predictability. Consider stabilizing WIP, reducing interruptions, and improving estimation.")
	}

	return textResult(sb.String()), nil
}

// PMMeetingEffectiveness records ceremony effectiveness.
func (h *Handlers) PMMeetingEffectiveness(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ceremony, err := req.RequireString("ceremony")
	if err != nil {
		return errorResult("ceremony required (standup/planning/retro/review/grooming)"), nil
	}
	duration, err := req.RequireInt("duration_minutes")
	if err != nil {
		return errorResult("duration_minutes required"), nil
	}
	score, err := req.RequireInt("score")
	if err != nil {
		return errorResult("score required (1-5)"), nil
	}
	notes := req.GetString("notes", "")
	sprintName := req.GetString("sprint_name", "")

	m := &memdom.MeetingEffectiveness{
		Ceremony:        ceremony,
		DurationMinutes: duration,
		Score:           score,
		Notes:           notes,
		SprintName:      sprintName,
	}
	if err := h.Memory.SaveMeetingEffectiveness(ctx, m); err != nil {
		return errorResult("Failed: " + err.Error()), nil
	}

	// Get trend
	history, _ := h.Memory.GetMeetingEffectivenessHistory(ctx, ceremony, 5)
	if len(history) > 1 {
		prevAvg := 0.0
		for _, h := range history[1:] {
			prevAvg += float64(h.Score)
		}
		prevAvg /= float64(len(history) - 1)
		diff := float64(score) - prevAvg
		trend := "stable"
		if diff > 0.3 {
			trend = "improving"
		} else if diff < -0.3 {
			trend = "declining"
		}
		return textResult(fmt.Sprintf("Recorded: %s effectiveness %d/5 (%d min). Trend: %s (avg %.1f).", ceremony, score, duration, trend, prevAvg)), nil
	}

	return textResult(fmt.Sprintf("Recorded: %s effectiveness %d/5 (%d min). First entry - no trend yet.", ceremony, score, duration)), nil
}

// PMMeetingTrends shows meeting effectiveness over time.
func (h *Handlers) PMMeetingTrends(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ceremony := req.GetString("ceremony", "")

	history, err := h.Memory.GetMeetingEffectivenessHistory(ctx, ceremony, 50)
	if err != nil {
		return errorResult("Failed: " + err.Error()), nil
	}
	if len(history) == 0 {
		return textResult("No meeting effectiveness data. Use pm_meeting_effectiveness to record."), nil
	}

	// Group by ceremony
	type stats struct {
		totalScore    int
		totalDuration int
		count         int
	}
	byCeremony := map[string]*stats{}
	for _, m := range history {
		if _, ok := byCeremony[m.Ceremony]; !ok {
			byCeremony[m.Ceremony] = &stats{}
		}
		byCeremony[m.Ceremony].totalScore += m.Score
		byCeremony[m.Ceremony].totalDuration += m.DurationMinutes
		byCeremony[m.Ceremony].count++
	}

	var sb strings.Builder
	sb.WriteString("Meeting Effectiveness Trends:\n\n")
	for c, s := range byCeremony {
		avgScore := float64(s.totalScore) / float64(s.count)
		avgDuration := float64(s.totalDuration) / float64(s.count)
		sb.WriteString(fmt.Sprintf("  %s: avg %.1f/5, avg %.0f min (%d entries)\n", c, avgScore, avgDuration, s.count))
	}
	return textResult(sb.String()), nil
}

// PMTeamRadar records multi-dimension team assessment.
func (h *Handlers) PMTeamRadar(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintName, err := req.RequireString("sprint_name")
	if err != nil {
		return errorResult("sprint_name required"), nil
	}
	dimensionsStr, err := req.RequireString("dimensions")
	if err != nil {
		return errorResult("dimensions required (JSON: {\"delivery\": 4, \"quality\": 3, ...})"), nil
	}

	var dimensions map[string]int
	if err := json.Unmarshal([]byte(dimensionsStr), &dimensions); err != nil {
		return errorResult("dimensions must be valid JSON: {\"delivery\": 4, \"quality\": 3}"), nil
	}

	for dim, score := range dimensions {
		r := &memdom.TeamRadar{SprintName: sprintName, Dimension: dim, Score: score}
		if err := h.Memory.SaveTeamRadar(ctx, r); err != nil {
			return errorResult("Failed: " + err.Error()), nil
		}
	}

	// Get previous sprint data for comparison
	history, _ := h.Memory.GetTeamRadarHistory(ctx, 200)
	prevScores := map[string][]int{}
	for _, r := range history {
		if r.SprintName != sprintName {
			prevScores[r.Dimension] = append(prevScores[r.Dimension], r.Score)
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Team Radar - %s:\n\n", sprintName))
	for dim, score := range dimensions {
		line := fmt.Sprintf("  %s: %d/5", dim, score)
		if prev, ok := prevScores[dim]; ok && len(prev) > 0 {
			prevAvg := 0
			for _, s := range prev {
				prevAvg += s
			}
			diff := float64(score) - float64(prevAvg)/float64(len(prev))
			if diff > 0 {
				line += fmt.Sprintf(" (+%.1f)", diff)
			} else if diff < 0 {
				line += fmt.Sprintf(" (%.1f)", diff)
			}
		}
		sb.WriteString(line + "\n")
	}
	return textResult(sb.String()), nil
}

// PMTeamRadarHistory shows radar trends across sprints.
func (h *Handlers) PMTeamRadarHistory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	history, err := h.Memory.GetTeamRadarHistory(ctx, 200)
	if err != nil {
		return errorResult("Failed: " + err.Error()), nil
	}
	if len(history) == 0 {
		return textResult("No radar data. Use pm_team_radar to record."), nil
	}

	// Group by sprint then dimension
	type dimScore struct {
		total int
		count int
	}
	bySprint := map[string]map[string]*dimScore{}
	var sprintOrder []string
	seen := map[string]bool{}
	for _, r := range history {
		if !seen[r.SprintName] {
			seen[r.SprintName] = true
			sprintOrder = append(sprintOrder, r.SprintName)
		}
		if bySprint[r.SprintName] == nil {
			bySprint[r.SprintName] = map[string]*dimScore{}
		}
		if bySprint[r.SprintName][r.Dimension] == nil {
			bySprint[r.SprintName][r.Dimension] = &dimScore{}
		}
		bySprint[r.SprintName][r.Dimension].total += r.Score
		bySprint[r.SprintName][r.Dimension].count++
	}

	var sb strings.Builder
	sb.WriteString("Team Radar History:\n\n")
	for _, sprint := range sprintOrder {
		sb.WriteString(fmt.Sprintf("%s:\n", sprint))
		for dim, ds := range bySprint[sprint] {
			avg := float64(ds.total) / float64(ds.count)
			sb.WriteString(fmt.Sprintf("  %s: %.1f/5\n", dim, avg))
		}
		sb.WriteString("\n")
	}
	return textResult(sb.String()), nil
}

// PMMaturityAssessment provides AI-powered agile maturity assessment.
func (h *Handlers) PMMaturityAssessment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	if h.AI == nil {
		return errorResult("AI provider not configured"), nil
	}

	// Gather data points
	var contextData strings.Builder

	// Sprint predictability
	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 10)
	if len(snaps) > 0 {
		contextData.WriteString("Sprint History:\n")
		for _, s := range snaps {
			contextData.WriteString(fmt.Sprintf("  %s: done=%d, total=%d, velocity=%d, carryover=%d\n", s.SprintName, s.Done, s.TotalIssues, s.Velocity, s.Carryover))
		}
	}

	// Retro action completion
	actions, _ := h.Memory.GetPendingActionItems(ctx)
	contextData.WriteString(fmt.Sprintf("\nPending retro actions: %d\n", len(actions)))

	// Goals
	goals, _ := h.Memory.GetGoalHistory(ctx, boardID, 5)
	if len(goals) > 0 {
		contextData.WriteString("Sprint Goals:\n")
		for _, g := range goals {
			contextData.WriteString(fmt.Sprintf("  %s: %s - %s\n", g.SprintName, g.Goal, g.Status))
		}
	}

	// Blockers
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	contextData.WriteString(fmt.Sprintf("\nActive blockers: %d\n", len(blockers)))

	// Health scores
	healthScores, _ := h.Memory.GetHealthScores(ctx, boardID, 5)
	if len(healthScores) > 0 {
		contextData.WriteString("Health Scores:\n")
		for _, hs := range healthScores {
			contextData.WriteString(fmt.Sprintf("  %s: %d/100\n", hs.SprintName, hs.OverallScore))
		}
	}

	systemPrompt := `You are an Agile Coach assessing team maturity on a 1-5 scale:
1=Initial (chaotic, no process), 2=Managed (basic scrum, reactive), 3=Defined (consistent process, proactive), 4=Quantitatively Managed (data-driven, predictable), 5=Optimizing (continuous improvement culture).

Assess based on the data provided. Return:
- Overall maturity level (1-5)
- Score per dimension: Predictability, Process Discipline, Continuous Improvement, Team Autonomy, Stakeholder Trust
- Top 3 specific improvement actions

Be concise and data-driven. No fluff.`

	result, err := h.AI.Complete(ctx, systemPrompt, contextData.String())
	if err != nil {
		return errorResult("AI analysis failed: " + err.Error()), nil
	}

	return textResult(result), nil
}

// PMDailyDigestCoaching generates a morning brief combining all available data.
func (h *Handlers) PMDailyDigestCoaching(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	var sb strings.Builder
	sb.WriteString("Daily Digest\n\n")

	// Active blockers
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) > 0 {
		sb.WriteString(fmt.Sprintf("BLOCKERS (%d active):\n", len(blockers)))
		for _, b := range blockers {
			sb.WriteString(fmt.Sprintf("  - %s", b.Description))
			if b.IssueKey != "" {
				sb.WriteString(fmt.Sprintf(" [%s]", b.IssueKey))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Stale items from sprint
	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) > 0 {
		issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
		inProgress := 0
		todo := 0
		done := 0
		for _, issue := range issues {
			lower := strings.ToLower(issue.Status)
			if strings.Contains(lower, "done") || strings.Contains(lower, "closed") {
				done++
			} else if strings.Contains(lower, "progress") {
				inProgress++
			} else {
				todo++
			}
		}
		total := len(issues)
		sb.WriteString(fmt.Sprintf("SPRINT: %s\n", sprints[0].Name))
		sb.WriteString(fmt.Sprintf("  Done: %d/%d | In Progress: %d | To Do: %d\n\n", done, total, inProgress, todo))
	}

	// Open risks
	risks, _ := h.Memory.GetOpenRisks(ctx)
	if len(risks) > 0 {
		sb.WriteString(fmt.Sprintf("RISKS (%d open):\n", len(risks)))
		for _, r := range risks {
			sb.WriteString(fmt.Sprintf("  - [%s] %s", r.Severity, r.Title))
			if r.Owner != "" {
				sb.WriteString(fmt.Sprintf(" (owner: %s)", r.Owner))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Pending action items
	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 0 {
		sb.WriteString(fmt.Sprintf("PENDING ACTIONS (%d):\n", len(actions)))
		for i, a := range actions {
			if i >= 5 {
				sb.WriteString(fmt.Sprintf("  ... and %d more\n", len(actions)-5))
				break
			}
			sb.WriteString(fmt.Sprintf("  - %s", a.Description))
			if a.Owner != "" {
				sb.WriteString(fmt.Sprintf(" (@%s)", a.Owner))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	}

	// Open dependencies
	deps, _ := h.Memory.GetOpenDependencies(ctx)
	if len(deps) > 0 {
		sb.WriteString(fmt.Sprintf("DEPENDENCIES (%d open):\n", len(deps)))
		for _, d := range deps {
			sb.WriteString(fmt.Sprintf("  - %s -> %s: %s\n", d.FromIssueKey, d.ToIssueKey, d.Description))
		}
		sb.WriteString("\n")
	}

	if sb.Len() < 30 {
		sb.WriteString("All clear! No blockers, risks, or pending actions.")
	}

	return textResult(sb.String()), nil
}
