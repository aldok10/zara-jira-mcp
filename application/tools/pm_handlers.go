package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// PMRecommendations generates AI-powered recommendations based on historical memory.
func (h *Handlers) PMRecommendations(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)

	// Gather all context from memory
	var memoryContext strings.Builder

	// Recent sprint snapshots
	if boardID > 0 {
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
		if len(snaps) > 0 {
			memoryContext.WriteString("Sprint History (last 5):\n")
			for _, s := range snaps {
				memoryContext.WriteString(fmt.Sprintf("  %s [%s]: Total=%d Done=%d InProgress=%d Todo=%d Blocked=%d Velocity=%d Completion=%.0f%%\n",
					s.SprintName, s.SnapshotDate.Format("2006-01-02"), s.TotalIssues, s.Done, s.InProgress, s.Todo, s.Blocked, s.Velocity, s.CompletionRate))
			}
			memoryContext.WriteString("\n")
		}
	}

	// Open risks
	risks, _ := h.Memory.GetOpenRisks(ctx)
	if len(risks) > 0 {
		memoryContext.WriteString(fmt.Sprintf("Open Risks (%d):\n", len(risks)))
		for _, r := range risks {
			days := int(time.Since(r.IdentifiedAt).Hours() / 24)
			memoryContext.WriteString(fmt.Sprintf("  [%s] %s (owner: %s, days open: %d)\n", r.Severity, r.Title, r.Owner, days))
		}
		memoryContext.WriteString("\n")
	}

	// Active blockers
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) > 0 {
		memoryContext.WriteString(fmt.Sprintf("Active Blockers (%d):\n", len(blockers)))
		for _, b := range blockers {
			days := int(time.Since(b.BlockedSince).Hours() / 24)
			memoryContext.WriteString(fmt.Sprintf("  %s (owner: %s, days: %d, issue: %s)\n", b.Description, b.Owner, days, b.IssueKey))
		}
		memoryContext.WriteString("\n")
	}

	// Pending action items
	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 0 {
		memoryContext.WriteString(fmt.Sprintf("Pending Retro Actions (%d):\n", len(actions)))
		for _, a := range actions {
			memoryContext.WriteString(fmt.Sprintf("  - %s (owner: %s)\n", a.Description, a.Owner))
		}
		memoryContext.WriteString("\n")
	}

	// Recent decisions
	decisions, _ := h.Memory.GetDecisions(ctx, 5)
	if len(decisions) > 0 {
		memoryContext.WriteString("Recent Decisions:\n")
		for _, d := range decisions {
			memoryContext.WriteString(fmt.Sprintf("  [%s] %s: %s\n", d.MadeAt.Format("2006-01-02"), d.Title, d.Decision))
		}
		memoryContext.WriteString("\n")
	}

	if memoryContext.Len() == 0 {
		return textResult("No historical data yet. Use snapshot_sprint, record_risk, record_decision, and record_blocker to build up PM memory first."), nil
	}

	focus := req.GetString("focus", "general")

	systemPrompt := `You are an expert Scrum Master / PM advisor with full historical context of this team.
Based on the historical data provided, give actionable recommendations.

Focus areas based on request:
- "general": Overall health, top 3 priorities for PM this week
- "velocity": Velocity trends, capacity planning suggestions  
- "risks": Risk mitigation priorities, escalation suggestions
- "team": Workload balance, burnout signals, skill gaps
- "process": Process improvements, ceremony effectiveness

Rules:
- Be specific, reference data points
- Prioritize: what's most urgent vs important
- Suggest concrete next actions, not vague advice
- Flag patterns (recurring blockers, declining velocity, etc.)
- Keep recommendations to 5 max, ranked by impact`

	userPrompt := fmt.Sprintf("Focus: %s\n\nHistorical PM Data:\n%s\nGive me your top recommendations.",
		focus, memoryContext.String())

	analysis, err := h.AI.Complete(ctx, systemPrompt, userPrompt)
	if err != nil {
		return errorResult("AI analysis failed: " + err.Error()), nil
	}

	return textResult(analysis), nil
}

// VelocityTrend shows velocity over recent sprints with trend analysis.
func (h *Handlers) VelocityTrend(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	snaps, err := h.Memory.GetSprintSnapshots(ctx, boardID, 10)
	if err != nil {
		return errorResult("Failed to get snapshots: " + err.Error()), nil
	}

	if len(snaps) == 0 {
		return textResult("No sprint snapshots yet. Use pm_snapshot_sprint to capture sprint data."), nil
	}

	var sb strings.Builder
	sb.WriteString("Velocity Trend:\n\n")
	sb.WriteString("Sprint | Velocity | Completion | Done/Total | Blocked | Carryover\n")
	sb.WriteString("-------|----------|------------|-----------|---------|----------\n")

	var totalVelocity int
	for _, s := range snaps {
		sb.WriteString(fmt.Sprintf("%s | %d | %.0f%% | %d/%d | %d | %d\n",
			s.SprintName, s.Velocity, s.CompletionRate, s.Done, s.TotalIssues, s.Blocked, s.Carryover))
		totalVelocity += s.Velocity
	}

	avgVelocity := 0
	if len(snaps) > 0 {
		avgVelocity = totalVelocity / len(snaps)
	}

	sb.WriteString(fmt.Sprintf("\nAverage Velocity: %d\n", avgVelocity))
	sb.WriteString(fmt.Sprintf("Sprints Tracked: %d\n", len(snaps)))

	// Simple trend detection
	if len(snaps) >= 3 {
		recent := snaps[0].Velocity
		older := snaps[len(snaps)-1].Velocity
		if recent > older {
			sb.WriteString("Trend: IMPROVING (velocity increasing)\n")
		} else if recent < older {
			sb.WriteString("Trend: DECLINING (velocity decreasing - investigate)\n")
		} else {
			sb.WriteString("Trend: STABLE\n")
		}
	}

	return textResult(sb.String()), nil
}

// StandupPrep generates a daily standup preparation brief.
func (h *Handlers) StandupPrep(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	var briefContext strings.Builder

	// Current sprint state from Jira
	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err == nil && len(sprints) > 0 {
		sprint := sprints[0]
		issues, _ := h.Jira.GetSprintIssues(ctx, sprint.ID)
		if len(issues) > 0 {
			briefContext.WriteString(fmt.Sprintf("Current Sprint: %s (Goal: %s)\n", sprint.Name, sprint.Goal))
			var done, inProgress, blocked int
			var blockedIssues []string
			for _, issue := range issues {
				switch {
				case strings.Contains(strings.ToLower(issue.Status), "done") || strings.Contains(strings.ToLower(issue.Status), "closed"):
					done++
				case strings.Contains(strings.ToLower(issue.Status), "progress") || strings.Contains(strings.ToLower(issue.Status), "review"):
					inProgress++
				case strings.Contains(strings.ToLower(issue.Status), "block"):
					blocked++
					blockedIssues = append(blockedIssues, fmt.Sprintf("  - %s: %s (%s)", issue.Key, issue.Summary, issue.Assignee))
				}
			}
			briefContext.WriteString(fmt.Sprintf("Progress: Done=%d, In Progress=%d, Blocked=%d, Total=%d\n", done, inProgress, blocked, len(issues)))
			if len(blockedIssues) > 0 {
				briefContext.WriteString("Blocked Issues:\n" + strings.Join(blockedIssues, "\n") + "\n")
			}
			briefContext.WriteString("\n")
		}
	}

	// Active blockers from memory
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) > 0 {
		briefContext.WriteString(fmt.Sprintf("Tracked Blockers (%d):\n", len(blockers)))
		for _, b := range blockers {
			days := int(time.Since(b.BlockedSince).Hours() / 24)
			briefContext.WriteString(fmt.Sprintf("  - [%d days] %s (owner: %s)\n", days, b.Description, b.Owner))
		}
		briefContext.WriteString("\n")
	}

	// Open risks
	risks, _ := h.Memory.GetOpenRisks(ctx)
	criticalRisks := 0
	for _, r := range risks {
		if r.Severity == "critical" || r.Severity == "high" {
			criticalRisks++
		}
	}
	if criticalRisks > 0 {
		briefContext.WriteString(fmt.Sprintf("High/Critical Risks: %d (review in standup)\n\n", criticalRisks))
	}

	// Pending action items
	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 0 {
		briefContext.WriteString(fmt.Sprintf("Pending Action Items (%d):\n", len(actions)))
		for _, a := range actions {
			briefContext.WriteString(fmt.Sprintf("  - %s (owner: %s)\n", a.Description, a.Owner))
		}
		briefContext.WriteString("\n")
	}

	systemPrompt := `You are a Scrum Master preparing for daily standup.
Generate a brief standup prep note for the PM/Scrum Master with:
1. Key talking points (what to bring up)
2. Blockers to address (prioritized)
3. Follow-ups needed (who to check with)
4. Sprint health signal (on track / watch / at risk)

Keep it concise - this is a quick prep, not a report. Bullet points. Under 200 words.`

	analysis, err := h.AI.Complete(ctx, systemPrompt, briefContext.String())
	if err != nil {
		// Fallback: return raw context if AI fails
		return textResult("Standup Prep (raw):\n\n" + briefContext.String()), nil
	}

	return textResult(analysis), nil
}

// SprintRetroAnalysis generates AI analysis comparing current sprint to historical patterns.
func (h *Handlers) SprintRetroAnalysis(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	var analysisContext strings.Builder

	// Sprint history
	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 10)
	if len(snaps) > 0 {
		analysisContext.WriteString("Sprint History:\n")
		for _, s := range snaps {
			analysisContext.WriteString(fmt.Sprintf("  %s: velocity=%d, completion=%.0f%%, blocked=%d, carryover=%d\n",
				s.SprintName, s.Velocity, s.CompletionRate, s.Blocked, s.Carryover))
		}
		analysisContext.WriteString("\n")
	}

	// Past retros
	retros, _ := h.Memory.GetRetrospectives(ctx, 5)
	if len(retros) > 0 {
		analysisContext.WriteString("Past Retrospectives:\n")
		for _, r := range retros {
			analysisContext.WriteString(fmt.Sprintf("  %s: Well=%s | Improve=%s | Actions=%s\n",
				r.SprintName, r.WentWell, r.Improvements, r.ActionItems))
		}
		analysisContext.WriteString("\n")
	}

	// Blocker patterns
	blockerHistory, _ := h.Memory.GetBlockerHistory(ctx, 20)
	if len(blockerHistory) > 0 {
		analysisContext.WriteString(fmt.Sprintf("Blocker History (%d total):\n", len(blockerHistory)))
		avgDays := 0
		for _, b := range blockerHistory {
			avgDays += b.DaysBlocked
		}
		if len(blockerHistory) > 0 {
			avgDays /= len(blockerHistory)
		}
		analysisContext.WriteString(fmt.Sprintf("  Avg resolution time: %d days\n", avgDays))
	}

	if analysisContext.Len() == 0 {
		return textResult("Not enough historical data for retrospective analysis. Record more sprint snapshots and retros first."), nil
	}

	systemPrompt := `You are an expert Scrum Master performing a sprint retrospective analysis.
Based on historical sprint data, identify:
1. Patterns: What keeps repeating (good and bad)?
2. Trends: Is the team improving, stable, or declining?
3. Root Causes: For recurring problems, what's the underlying issue?
4. Specific Actions: What concrete changes would have the highest impact?
5. Celebrations: What's genuinely going well (be specific)?

Be data-driven. Reference specific sprints and metrics. No generic advice.
If previous retro action items keep appearing without resolution, call it out.`

	analysis, err := h.AI.Complete(ctx, systemPrompt, analysisContext.String())
	if err != nil {
		return errorResult("AI analysis failed: " + err.Error()), nil
	}

	return textResult(analysis), nil
}
