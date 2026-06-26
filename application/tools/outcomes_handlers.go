package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handlers) initOutcomeTables() {
	db := h.Memory.DB()
	db.Exec(`CREATE TABLE IF NOT EXISTS stakeholder_pulse (
		id INTEGER PRIMARY KEY,
		stakeholder TEXT,
		score INTEGER,
		sprint_name TEXT,
		feedback TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	db.Exec(`CREATE TABLE IF NOT EXISTS outcome_map (
		id INTEGER PRIMARY KEY,
		sprint_name TEXT,
		objective TEXT,
		key_results TEXT,
		progress_note TEXT,
		board_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
}

// PMImpedimentAging tracks how long blockers stay alive.
func (h *Handlers) PMImpedimentAging(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	active, _ := h.Memory.GetActiveBlockers(ctx)
	history, _ := h.Memory.GetBlockerHistory(ctx, 50)

	var sb strings.Builder
	sb.WriteString("Impediment Aging Report\n\n")

	// Active blockers with age
	now := time.Now()
	chronic := 0
	if len(active) > 0 {
		sb.WriteString(fmt.Sprintf("ACTIVE BLOCKERS (%d):\n", len(active)))
		for _, b := range active {
			days := int(now.Sub(b.BlockedSince).Hours() / 24)
			flag := ""
			if days > 3 {
				flag = " [CHRONIC]"
				chronic++
			}
			sb.WriteString(fmt.Sprintf("  - %s (%d days)%s", b.Description, days, flag))
			if b.IssueKey != "" {
				sb.WriteString(fmt.Sprintf(" [%s]", b.IssueKey))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("\n")
	} else {
		sb.WriteString("No active blockers.\n\n")
	}

	// Avg resolution time from history
	resolved := 0
	totalDays := 0
	for _, b := range history {
		if b.ResolvedAt != nil {
			resolved++
			totalDays += int(b.ResolvedAt.Sub(b.BlockedSince).Hours() / 24)
		}
	}

	if resolved > 0 {
		avg := float64(totalDays) / float64(resolved)
		sb.WriteString(fmt.Sprintf("RESOLUTION STATS:\n  Avg resolution time: %.1f days\n  Resolved count: %d\n  Chronic (>3 days active): %d\n", avg, resolved, chronic))
	} else {
		sb.WriteString(fmt.Sprintf("RESOLUTION STATS:\n  No resolved blockers yet.\n  Chronic (>3 days active): %d\n", chronic))
	}

	return textResult(sb.String()), nil
}

// PMSMImpact tracks Scrum Master's measurable impact for a sprint.
func (h *Handlers) PMSMImpact(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintName := req.GetString("sprint_name", "")

	// Blockers resolved
	history, _ := h.Memory.GetBlockerHistory(ctx, 100)
	blockersResolved := 0
	totalResDays := 0
	for _, b := range history {
		if b.ResolvedAt == nil {
			continue
		}
		if sprintName != "" && b.IssueKey == "" {
			// Can't filter by sprint without better data, include all
		}
		blockersResolved++
		totalResDays += int(b.ResolvedAt.Sub(b.BlockedSince).Hours() / 24)
	}
	avgRes := 0.0
	if blockersResolved > 0 {
		avgRes = float64(totalResDays) / float64(blockersResolved)
	}

	// Risks mitigated
	allRisks, _ := h.Memory.GetAllRisks(ctx, 100)
	risksMitigated := 0
	for _, r := range allRisks {
		if r.Status == "resolved" || r.Status == "mitigating" {
			risksMitigated++
		}
	}

	// Pending action items (fewer = better SM follow-through)
	pending, _ := h.Memory.GetPendingActionItems(ctx)

	// Retros (count as facilitation impact)
	retros, _ := h.Memory.GetRetrospectives(ctx, 10)

	var sb strings.Builder
	sb.WriteString("SM Impact Report")
	if sprintName != "" {
		sb.WriteString(fmt.Sprintf(" (%s)", sprintName))
	}
	sb.WriteString("\n\n")
	sb.WriteString(fmt.Sprintf("Blockers Resolved: %d\n", blockersResolved))
	sb.WriteString(fmt.Sprintf("Avg Resolution Time: %.1f days\n", avgRes))
	sb.WriteString(fmt.Sprintf("Risks Mitigated/Mitigating: %d\n", risksMitigated))
	sb.WriteString(fmt.Sprintf("Pending Action Items: %d (lower is better)\n", len(pending)))
	sb.WriteString(fmt.Sprintf("Retros Facilitated: %d\n", len(retros)))
	sb.WriteString("\n")

	// Narrative score
	score := blockersResolved*10 + risksMitigated*5 - len(pending)*2
	if score < 0 {
		score = 0
	}
	sb.WriteString(fmt.Sprintf("Impact Score: %d points\n", score))
	if score > 50 {
		sb.WriteString("High impact. Blockers cleared fast, risks managed.")
	} else if score > 20 {
		sb.WriteString("Moderate impact. Room to improve follow-through on action items.")
	} else {
		sb.WriteString("Low measured impact. Consider: are blockers being recorded? Are retro actions tracked?")
	}

	return textResult(sb.String()), nil
}

// PMStakeholderPulse records stakeholder satisfaction.
func (h *Handlers) PMStakeholderPulse(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOutcomeTables()

	stakeholder, err := req.RequireString("stakeholder")
	if err != nil {
		return errorResult("stakeholder required"), nil
	}
	score, err := req.RequireInt("score")
	if err != nil {
		return errorResult("score required (1-5)"), nil
	}
	if score < 1 || score > 5 {
		return errorResult("score must be 1-5"), nil
	}
	sprintName := req.GetString("sprint_name", "")
	feedback := req.GetString("feedback", "")

	db := h.Memory.DB()
	_, execErr := db.Exec("INSERT INTO stakeholder_pulse (stakeholder, score, sprint_name, feedback) VALUES (?, ?, ?, ?)",
		stakeholder, score, sprintName, feedback)
	if execErr != nil {
		return errorResult("Failed to save: " + execErr.Error()), nil
	}

	return textResult(fmt.Sprintf("Recorded: %s satisfaction %d/5 for %s.", stakeholder, score, sprintName)), nil
}

// PMStakeholderTrend shows stakeholder satisfaction over time.
func (h *Handlers) PMStakeholderTrend(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOutcomeTables()

	db := h.Memory.DB()
	rows, err := db.Query("SELECT stakeholder, score, sprint_name, created_at FROM stakeholder_pulse ORDER BY created_at DESC LIMIT 100")
	if err != nil {
		return errorResult("Query failed: " + err.Error()), nil
	}
	defer rows.Close()

	type entry struct {
		score      int
		sprint     string
		createdAt  string
	}
	byStakeholder := map[string][]entry{}
	var order []string
	seen := map[string]bool{}

	for rows.Next() {
		var stakeholder, sprint, createdAt string
		var score int
		if err := rows.Scan(&stakeholder, &score, &sprint, &createdAt); err != nil {
			continue
		}
		if !seen[stakeholder] {
			seen[stakeholder] = true
			order = append(order, stakeholder)
		}
		byStakeholder[stakeholder] = append(byStakeholder[stakeholder], entry{score, sprint, createdAt})
	}

	if len(order) == 0 {
		return textResult("No stakeholder pulse data. Use pm_stakeholder_pulse to record."), nil
	}

	var sb strings.Builder
	sb.WriteString("Stakeholder Satisfaction Trends:\n\n")
	for _, name := range order {
		entries := byStakeholder[name]
		total := 0
		for _, e := range entries {
			total += e.score
		}
		avg := float64(total) / float64(len(entries))
		sb.WriteString(fmt.Sprintf("  %s: avg %.1f/5 (%d entries)\n", name, avg, len(entries)))
		for i, e := range entries {
			if i >= 5 {
				sb.WriteString(fmt.Sprintf("    ... and %d more\n", len(entries)-5))
				break
			}
			label := e.sprint
			if label == "" {
				label = e.createdAt
			}
			sb.WriteString(fmt.Sprintf("    %s: %d/5\n", label, e.score))
		}
	}
	return textResult(sb.String()), nil
}

// PMImprovementVelocity tracks if retro actions are getting done faster.
func (h *Handlers) PMImprovementVelocity(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	retros, _ := h.Memory.GetRetrospectives(ctx, 10)
	pending, _ := h.Memory.GetPendingActionItems(ctx)

	if len(retros) < 2 {
		return textResult("Need at least 2 retrospectives to calculate improvement velocity."), nil
	}

	// Count action items per retro by parsing ActionItems field
	var sb strings.Builder
	sb.WriteString("Improvement Velocity:\n\n")

	totalActions := 0
	for _, r := range retros {
		items := strings.Split(r.ActionItems, "\n")
		count := 0
		for _, item := range items {
			if strings.TrimSpace(item) != "" {
				count++
			}
		}
		totalActions += count
		sb.WriteString(fmt.Sprintf("  %s: %d actions created\n", r.SprintName, count))
	}

	sb.WriteString(fmt.Sprintf("\nTotal actions created: %d\n", totalActions))
	sb.WriteString(fmt.Sprintf("Currently pending: %d\n", len(pending)))

	completed := totalActions - len(pending)
	if completed < 0 {
		completed = 0
	}
	if totalActions > 0 {
		rate := float64(completed) / float64(totalActions) * 100
		sb.WriteString(fmt.Sprintf("Completion rate: %.0f%%\n", rate))
		if rate >= 80 {
			sb.WriteString("\nStrong follow-through. Retro actions are being executed.")
		} else if rate >= 50 {
			sb.WriteString("\nModerate. Some actions slip through. Consider fewer, more specific actions.")
		} else {
			sb.WriteString("\nLow completion. Retro actions are dying. Reduce quantity, increase accountability.")
		}
	}

	return textResult(sb.String()), nil
}

// PMTeamAutonomy assesses team self-organization level.
func (h *Handlers) PMTeamAutonomy(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	_, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	if h.AI == nil {
		return errorResult("AI provider not configured"), nil
	}

	// Gather data
	var contextData strings.Builder

	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	history, _ := h.Memory.GetBlockerHistory(ctx, 30)
	contextData.WriteString(fmt.Sprintf("Active blockers: %d\n", len(blockers)))
	contextData.WriteString(fmt.Sprintf("Historical blockers resolved: %d\n", len(history)))

	// Owner distribution in blockers
	ownerCount := map[string]int{}
	for _, b := range history {
		if b.Owner != "" {
			ownerCount[b.Owner]++
		}
	}
	if len(ownerCount) > 0 {
		contextData.WriteString("Blocker resolution by owner:\n")
		for owner, count := range ownerCount {
			contextData.WriteString(fmt.Sprintf("  %s: %d\n", owner, count))
		}
	}

	// Action items ownership spread
	pending, _ := h.Memory.GetPendingActionItems(ctx)
	actionOwners := map[string]int{}
	for _, a := range pending {
		if a.Owner != "" {
			actionOwners[a.Owner]++
		}
	}
	contextData.WriteString(fmt.Sprintf("\nPending action items: %d\n", len(pending)))
	if len(actionOwners) > 0 {
		contextData.WriteString("Action item owners:\n")
		for owner, count := range actionOwners {
			contextData.WriteString(fmt.Sprintf("  %s: %d\n", owner, count))
		}
	}

	// Team metrics for self-assignment signals
	metrics, _ := h.Memory.GetTeamMetrics(ctx, "", 50)
	if len(metrics) > 0 {
		contextData.WriteString(fmt.Sprintf("\nTeam metric entries: %d\n", len(metrics)))
	}

	systemPrompt := `You are assessing team self-organization and autonomy (1-5 scale):
1=SM does everything, 2=SM drives most actions, 3=Shared ownership, 4=Team drives most, 5=Fully autonomous.

Based on the data, assess:
- Blocker resolution: Is it concentrated in one person (likely SM) or spread?
- Action item ownership: diverse or single-owner?
- Overall autonomy score with 2-3 specific recommendations.

Be concise. No fluff. Data-driven verdict.`

	result, aiErr := h.AI.Complete(ctx, systemPrompt, contextData.String())
	if aiErr != nil {
		return errorResult("AI analysis failed: " + aiErr.Error()), nil
	}

	return textResult(result), nil
}

// PMOutcomeMap records OKR-sprint mapping.
func (h *Handlers) PMOutcomeMap(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOutcomeTables()

	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}
	objective, err := req.RequireString("objective")
	if err != nil {
		return errorResult("objective required"), nil
	}
	keyResults := req.GetString("key_results", "")

	// Get active sprint name
	sprintName := ""
	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) > 0 {
		sprintName = sprints[0].Name
	}

	db := h.Memory.DB()
	_, execErr := db.Exec("INSERT INTO outcome_map (sprint_name, objective, key_results, board_id) VALUES (?, ?, ?, ?)",
		sprintName, objective, keyResults, boardID)
	if execErr != nil {
		return errorResult("Failed to save: " + execErr.Error()), nil
	}

	msg := fmt.Sprintf("Mapped: Sprint '%s' serves objective '%s'.", sprintName, objective)
	if keyResults != "" {
		msg += fmt.Sprintf("\nKey Results: %s", keyResults)
	}
	return textResult(msg), nil
}

// PMOutcomeHistory shows OKR alignment over time.
func (h *Handlers) PMOutcomeHistory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOutcomeTables()

	db := h.Memory.DB()
	rows, err := db.Query("SELECT sprint_name, objective, key_results, created_at FROM outcome_map ORDER BY created_at DESC LIMIT 50")
	if err != nil {
		return errorResult("Query failed: " + err.Error()), nil
	}
	defer rows.Close()

	var sb strings.Builder
	sb.WriteString("Outcome Map History:\n\n")
	count := 0
	for rows.Next() {
		var sprint, objective, kr, createdAt string
		if err := rows.Scan(&sprint, &objective, &kr, &createdAt); err != nil {
			continue
		}
		count++
		sb.WriteString(fmt.Sprintf("  %s: %s\n", sprint, objective))
		if kr != "" {
			sb.WriteString(fmt.Sprintf("    KRs: %s\n", kr))
		}
	}

	if count == 0 {
		return textResult("No outcome mappings yet. Use pm_outcome_map to connect sprints to business objectives."), nil
	}
	return textResult(sb.String()), nil
}
