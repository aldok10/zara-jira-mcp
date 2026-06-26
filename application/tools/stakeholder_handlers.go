package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/mark3labs/mcp-go/mcp"
)

// ExecutiveReport generates a stakeholder-friendly report (not a sprint report — a BUSINESS report).
func (h *Handlers) ExecutiveReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	// Gather data
	var contextData strings.Builder

	// Sprint state
	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) > 0 {
		sprint := sprints[0]
		issues, _ := h.Jira.GetSprintIssues(ctx, sprint.ID)
		var done, total int
		total = len(issues)
		for _, issue := range issues {
			lower := strings.ToLower(issue.Status)
			if strings.Contains(lower, "done") || strings.Contains(lower, "closed") {
				done++
			}
		}
		contextData.WriteString(fmt.Sprintf("Current Sprint: %s (Goal: %s)\nProgress: %d/%d done (%.0f%%)\n\n",
			sprint.Name, sprint.Goal, done, total, float64(done)/float64(total)*100))
	}

	// Velocity trend
	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
	if len(snaps) > 0 {
		contextData.WriteString("Recent Sprint Results:\n")
		for _, s := range snaps {
			contextData.WriteString(fmt.Sprintf("  %s: velocity=%d, completion=%.0f%%\n", s.SprintName, s.Velocity, s.CompletionRate))
		}
		contextData.WriteString("\n")
	}

	// Health
	scores, _ := h.Memory.GetHealthScores(ctx, boardID, 1)
	if len(scores) > 0 {
		contextData.WriteString(fmt.Sprintf("Team Health Score: %d/100\n\n", scores[0].OverallScore))
	}

	// Risks
	risks, _ := h.Memory.GetOpenRisks(ctx)
	if len(risks) > 0 {
		contextData.WriteString(fmt.Sprintf("Open Risks: %d\n", len(risks)))
		for _, r := range risks {
			if r.Severity == "critical" || r.Severity == "high" {
				contextData.WriteString(fmt.Sprintf("  [%s] %s\n", r.Severity, r.Title))
			}
		}
		contextData.WriteString("\n")
	}

	// Goals
	goals, _ := h.Memory.GetActiveGoals(ctx, boardID)
	goalHistory, _ := h.Memory.GetGoalHistory(ctx, boardID, 3)
	if len(goalHistory) > 0 {
		contextData.WriteString("Goal Achievement:\n")
		for _, g := range goalHistory {
			contextData.WriteString(fmt.Sprintf("  %s: %s (%s)\n", g.SprintName, g.Goal, g.Status))
		}
		contextData.WriteString("\n")
	}
	if len(goals) > 0 {
		contextData.WriteString(fmt.Sprintf("Current Goal: %s\n\n", goals[0].Goal))
	}

	// Blockers
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) > 0 {
		contextData.WriteString(fmt.Sprintf("Active Blockers: %d\n", len(blockers)))
		for _, b := range blockers {
			days := int(time.Since(b.BlockedSince).Hours() / 24)
			contextData.WriteString(fmt.Sprintf("  [%d days] %s\n", days, b.Description))
		}
	}

	systemPrompt := `You are writing an executive stakeholder report for a software team.
Executives DON'T want: burndowns, story points, velocity numbers, ceremony details.
Executives WANT: business outcomes, risks to timeline, team health signal, what shipped, what's blocked.

Write a concise report (under 250 words) with these sections:
1. **Status** (one line: On Track / Watch / At Risk + why)
2. **Delivered This Sprint** (business value, not ticket numbers)
3. **Coming Next** (what stakeholders should expect)
4. **Risks & Blockers** (what could delay, what needs executive action)
5. **Team Health** (one sentence signal)

Tone: confident, concise, no jargon. Write for a VP who has 30 seconds.`

	report, err := h.AI.Complete(ctx, systemPrompt, contextData.String())
	if err != nil {
		return errorResult("AI failed: " + err.Error()), nil
	}

	sendToLark := req.GetBool("send_to_lark", false)
	if sendToLark {
		if err := h.Lark.SendMarkdown(ctx, "Executive Update", report); err != nil {
			return textResult(report + "\n\n(Lark send failed: " + err.Error() + ")"), nil
		}
		return textResult(report + "\n\n(Sent to Lark)"), nil
	}

	return textResult(report), nil
}

// SprintScorecard generates end-of-sprint quantified outcome card.
func (h *Handlers) SprintScorecard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	var done, total, blocked int
	total = len(issues)
	for _, issue := range issues {
		lower := strings.ToLower(issue.Status)
		switch {
		case strings.Contains(lower, "done") || strings.Contains(lower, "closed"):
			done++
		case strings.Contains(lower, "block"):
			blocked++
		}
	}

	// Scores (each 0-20, total 0-100)
	completionScore := 0
	if total > 0 {
		completionScore = int(float64(done) / float64(total) * 20)
	}

	// Goal achievement
	goalScore := 10 // neutral if no goal
	goals, _ := h.Memory.GetActiveGoals(ctx, boardID)
	if len(goals) > 0 {
		goalScore = 15 // has a goal = better
	}

	// Predictability (low carryover from previous)
	predictScore := 15
	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 2)
	if len(snaps) >= 2 {
		prev := snaps[1]
		if prev.Carryover > 0 && prev.TotalIssues > 0 {
			carryRatio := float64(prev.Carryover) / float64(prev.TotalIssues)
			if carryRatio > 0.3 {
				predictScore = 5
			} else if carryRatio > 0.15 {
				predictScore = 10
			} else {
				predictScore = 20
			}
		}
	}

	// Quality (no blocked items = good)
	qualityScore := 20
	if total > 0 {
		blockRatio := float64(blocked) / float64(total)
		qualityScore = 20 - int(blockRatio*40)
		if qualityScore < 0 {
			qualityScore = 0
		}
	}

	// Team balance
	teamScore := 15
	assignees := map[string]int{}
	for _, issue := range issues {
		if issue.Assignee != "" {
			assignees[issue.Assignee]++
		}
	}
	if len(assignees) > 1 {
		avg := total / len(assignees)
		maxLoad := 0
		for _, c := range assignees {
			if c > maxLoad {
				maxLoad = c
			}
		}
		if maxLoad > avg*2 {
			teamScore = 8
		} else {
			teamScore = 20
		}
	}

	totalScore := completionScore + goalScore + predictScore + qualityScore + teamScore

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Sprint Scorecard: %s\n\n", sprint.Name))
	sb.WriteString(fmt.Sprintf("TOTAL: %d/100\n\n", totalScore))
	sb.WriteString(fmt.Sprintf("  Completion:    %d/20 (%d/%d done)\n", completionScore, done, total))
	sb.WriteString(fmt.Sprintf("  Goal Focus:    %d/20\n", goalScore))
	sb.WriteString(fmt.Sprintf("  Predictability:%d/20\n", predictScore))
	sb.WriteString(fmt.Sprintf("  Quality:       %d/20 (%d blocked)\n", qualityScore, blocked))
	sb.WriteString(fmt.Sprintf("  Team Balance:  %d/20 (%d contributors)\n", teamScore, len(assignees)))

	grade := "F"
	switch {
	case totalScore >= 85:
		grade = "A"
	case totalScore >= 70:
		grade = "B"
	case totalScore >= 55:
		grade = "C"
	case totalScore >= 40:
		grade = "D"
	}
	sb.WriteString(fmt.Sprintf("\nGrade: %s\n", grade))

	return textResult(sb.String()), nil
}

// TeamKnowledgeBase provides onboarding context — how this team works.
func (h *Handlers) TeamKnowledgeBase(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)
	question := req.GetString("question", "")

	// Build context from all available memory
	var kb strings.Builder

	// DoD
	dod, _ := h.Memory.GetDoD(ctx, "*")
	if len(dod) > 0 {
		kb.WriteString("Definition of Done:\n")
		for _, d := range dod {
			kb.WriteString(fmt.Sprintf("  [%s] %s\n", d.Category, d.Item))
		}
		kb.WriteString("\n")
	}

	// Recent decisions
	decisions, _ := h.Memory.GetDecisions(ctx, 10)
	if len(decisions) > 0 {
		kb.WriteString("Key Decisions:\n")
		for _, d := range decisions {
			kb.WriteString(fmt.Sprintf("  [%s] %s: %s (why: %s)\n", d.MadeAt.Format("2006-01-02"), d.Title, d.Decision, d.Rationale))
		}
		kb.WriteString("\n")
	}

	// Sprint patterns
	if boardID > 0 {
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
		if len(snaps) > 0 {
			var totalVel int
			for _, s := range snaps {
				totalVel += s.Velocity
			}
			avgVel := totalVel / len(snaps)
			kb.WriteString(fmt.Sprintf("Team Metrics:\n  Avg Velocity: %d/sprint\n  Sprints tracked: %d\n\n", avgVel, len(snaps)))
		}

		goals, _ := h.Memory.GetGoalHistory(ctx, boardID, 5)
		achieved := 0
		for _, g := range goals {
			if g.Status == "achieved" {
				achieved++
			}
		}
		if len(goals) > 0 {
			kb.WriteString(fmt.Sprintf("Goal Success Rate: %d/%d (%.0f%%)\n\n", achieved, len(goals), float64(achieved)/float64(len(goals))*100))
		}
	}

	// Retro patterns
	retros, _ := h.Memory.GetRetrospectives(ctx, 3)
	if len(retros) > 0 {
		kb.WriteString("Recent Retro Themes:\n")
		for _, r := range retros {
			kb.WriteString(fmt.Sprintf("  %s: Well=%s | Improve=%s\n", r.SprintName, r.WentWell, r.Improvements))
		}
		kb.WriteString("\n")
	}

	// Recurring risks
	allRisks, _ := h.Memory.GetAllRisks(ctx, 20)
	if len(allRisks) > 0 {
		kb.WriteString(fmt.Sprintf("Risk History: %d risks tracked\n", len(allRisks)))
		riskTypes := map[string]int{}
		for _, r := range allRisks {
			riskTypes[r.Severity]++
		}
		for sev, count := range riskTypes {
			kb.WriteString(fmt.Sprintf("  %s: %d\n", sev, count))
		}
		kb.WriteString("\n")
	}

	if question != "" {
		// AI-powered Q&A about the team
		systemPrompt := `You are a team knowledge assistant. A new team member or stakeholder is asking about how this team works.
Answer based ONLY on the provided data. If you don't have data to answer, say so honestly.
Be concise and helpful. Format as bullet points.`

		userPrompt := fmt.Sprintf("Question: %s\n\nTeam Knowledge Base:\n%s", question, kb.String())
		answer, err := h.AI.Complete(ctx, systemPrompt, userPrompt)
		if err != nil {
			return textResult(kb.String()), nil
		}
		return textResult(answer), nil
	}

	if kb.Len() == 0 {
		return textResult("Team knowledge base is empty. Start recording decisions, DoD, retros, and sprint snapshots to build it."), nil
	}

	return textResult("Team Knowledge Base:\n\n" + kb.String()), nil
}

// RecordLearning captures team learnings/tribal knowledge.
func (h *Handlers) RecordLearning(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := req.RequireString("title")
	if err != nil {
		return errorResult("title required"), nil
	}
	learning, err := req.RequireString("learning")
	if err != nil {
		return errorResult("learning required (what did you learn)"), nil
	}

	d := &memdom.Decision{
		Title:     title,
		Context:   req.GetString("context", ""),
		Decision:  learning,
		Rationale: "learning",
		MadeBy:    req.GetString("author", "team"),
		MadeAt:    time.Now(),
		Tags:      "learning," + req.GetString("tags", ""),
	}

	if err := h.Memory.SaveDecision(ctx, d); err != nil {
		return errorResult("Failed to save: " + err.Error()), nil
	}

	return textResult(fmt.Sprintf("Learning recorded: %s\n%s", title, learning)), nil
}

// WeeklyDigest generates a weekly summary of all team activity.
func (h *Handlers) WeeklyDigest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)

	var contextData strings.Builder
	contextData.WriteString("Weekly Team Digest:\n\n")

	// This week's meetings
	meetings, _ := h.Memory.GetMeetingNotes(ctx, "", 10)
	weekAgo := time.Now().AddDate(0, 0, -7)
	weekMeetings := 0
	for _, m := range meetings {
		if m.Date.After(weekAgo) {
			weekMeetings++
			if m.Decisions != "" {
				contextData.WriteString(fmt.Sprintf("Decision [%s]: %s\n", m.MeetingType, m.Decisions))
			}
		}
	}
	contextData.WriteString(fmt.Sprintf("\nMeetings this week: %d\n", weekMeetings))

	// Risks opened/resolved this week
	allRisks, _ := h.Memory.GetAllRisks(ctx, 20)
	newRisks, resolvedRisks := 0, 0
	for _, r := range allRisks {
		if r.IdentifiedAt.After(weekAgo) {
			newRisks++
		}
		if r.ResolvedAt != nil && r.ResolvedAt.After(weekAgo) {
			resolvedRisks++
		}
	}
	contextData.WriteString(fmt.Sprintf("Risks: +%d new, -%d resolved\n", newRisks, resolvedRisks))

	// Blockers
	blockerHistory, _ := h.Memory.GetBlockerHistory(ctx, 20)
	newBlockers, resolvedBlockers := 0, 0
	for _, b := range blockerHistory {
		if b.BlockedSince.After(weekAgo) {
			newBlockers++
		}
		if b.ResolvedAt != nil && b.ResolvedAt.After(weekAgo) {
			resolvedBlockers++
		}
	}
	contextData.WriteString(fmt.Sprintf("Blockers: +%d new, -%d resolved\n", newBlockers, resolvedBlockers))

	// Sprint progress
	if boardID > 0 {
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 1)
		if len(snaps) > 0 {
			contextData.WriteString(fmt.Sprintf("\nSprint: %s (velocity: %d, completion: %.0f%%)\n",
				snaps[0].SprintName, snaps[0].Velocity, snaps[0].CompletionRate))
		}

		scores, _ := h.Memory.GetHealthScores(ctx, boardID, 1)
		if len(scores) > 0 {
			contextData.WriteString(fmt.Sprintf("Health: %d/100\n", scores[0].OverallScore))
		}
	}

	// Decisions made this week
	decisions, _ := h.Memory.GetDecisions(ctx, 10)
	weekDecisions := 0
	for _, d := range decisions {
		if d.MadeAt.After(weekAgo) {
			weekDecisions++
		}
	}
	contextData.WriteString(fmt.Sprintf("Decisions recorded: %d\n", weekDecisions))

	systemPrompt := `Summarize this weekly team data into a concise digest.
Format:
1. Headline (one sentence: what defined this week)
2. Wins (what went well)
3. Concerns (what needs attention)
4. Next week focus (what to prioritize)

Keep under 150 words. Be specific, reference the data.`

	digest, err := h.AI.Complete(ctx, systemPrompt, contextData.String())
	if err != nil {
		return textResult(contextData.String()), nil
	}

	sendToLark := req.GetBool("send_to_lark", false)
	if sendToLark {
		h.Lark.SendMarkdown(ctx, "Weekly Digest", digest)
		return textResult(digest + "\n\n(Sent to Lark)"), nil
	}

	return textResult(digest), nil
}
