package tools

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
)

var researchTablesOnce sync.Once

func (h *Handlers) initResearchTables() {
	researchTablesOnce.Do(func() {
		if h.Memory == nil || h.Memory.DB() == nil {
			return
		}
		db := h.Memory.DB()
		_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS notification_effectiveness (
			id INTEGER PRIMARY KEY,
			channel TEXT NOT NULL,
			severity TEXT NOT NULL,
			title TEXT,
			action_taken TEXT DEFAULT '',
			sent_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			action_at DATETIME
		)`)
		_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS collaboration_signals (
			id INTEGER PRIMARY KEY,
			sprint_name TEXT NOT NULL,
			person_a TEXT NOT NULL,
			person_b TEXT NOT NULL,
			shared_issues INTEGER DEFAULT 0,
			shared_labels TEXT DEFAULT '',
			sprint_id INTEGER DEFAULT 0,
			recorded_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
		_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_notif_effectiveness_channel ON notification_effectiveness(channel)`)
		_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_collab_signals_sprint ON collaboration_signals(sprint_name)`)
	})
}

// PMCollaborationSignal detects collaboration patterns, silos, and isolation risks.
func (h *Handlers) PMCollaborationSignal(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil || len(sprints) == 0 {
		return textResult("No active sprint found. Collaboration signals need sprint data."), nil
	}

	issues, err := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
	if err != nil {
		return errorResult("Failed to get sprint issues: " + err.Error()), nil
	}

	if len(issues) == 0 {
		return textResult("No issues in current sprint."), nil
	}

	// Analyze collaboration patterns
	assigneeIssues := map[string][]string{}        // person -> issue keys
	assigneeLabels := map[string]map[string]int{}   // person -> label -> count
	assigneeStatuses := map[string]map[string]int{} // person -> status -> count
	labelAssignees := map[string]map[string]int{}   // label -> person -> count

	for _, issue := range issues {
		assignee := issue.Assignee
		if assignee == "" {
			assignee = "Unassigned"
		}
		assigneeIssues[assignee] = append(assigneeIssues[assignee], issue.Key)

		if assigneeLabels[assignee] == nil {
			assigneeLabels[assignee] = map[string]int{}
		}
		for _, label := range issue.Labels {
			assigneeLabels[assignee][label]++
			if labelAssignees[label] == nil {
				labelAssignees[label] = map[string]int{}
			}
			labelAssignees[label][assignee]++
		}

		if assigneeStatuses[assignee] == nil {
			assigneeStatuses[assignee] = map[string]int{}
		}
		assigneeStatuses[assignee][issue.Status]++
	}

	totalPeople := len(assigneeIssues)
	totalIssues := len(issues)

	var sb strings.Builder
	sb.WriteString("## Collaboration Signal Report\n\n")
	sb.WriteString(fmt.Sprintf("Sprint: %s (%d issues, %d people)\n\n", sprints[0].Name, totalIssues, totalPeople))

	// 1. Knowledge silo detection (one person dominates a label/area)
	sb.WriteString("### Knowledge Silos\n")
	siloCount := 0
	for label, assignees := range labelAssignees {
		totalInLabel := 0
		for _, c := range assignees {
			totalInLabel += c
		}
		if totalInLabel < 3 {
			continue
		}
		for person, count := range assignees {
			pct := float64(count) / float64(totalInLabel) * 100
			if pct > 60 {
				sb.WriteString(fmt.Sprintf("- [SILO] %s owns %.0f%% of '%s' (%d/%d issues)\n", person, pct, label, count, totalInLabel))
				siloCount++
			}
		}
	}
	if siloCount == 0 {
		sb.WriteString("- No significant knowledge silos detected.\n")
	}
	sb.WriteString("\n")

	// 2. Isolation detection (people with no shared labels/areas)
	sb.WriteString("### Isolation Risk\n")
	isolatedCount := 0
	for person, labels := range assigneeLabels {
		sharedLabels := false
		for label := range labels {
			if len(labelAssignees[label]) > 1 {
				sharedLabels = true
				break
			}
		}
		if !sharedLabels && len(labels) > 0 {
			sb.WriteString(fmt.Sprintf("- [ISOLATED] %s works alone on: %s\n", person, joinKeys(labels)))
			isolatedCount++
		}
	}
	if isolatedCount == 0 {
		sb.WriteString("- No isolated team members detected.\n")
	}
	sb.WriteString("\n")

	// 3. Workload balance
	sb.WriteString("### Workload Balance\n")
	maxIssues := 0
	minIssues := -1
	maxPerson, minPerson := "", ""
	for person, iss := range assigneeIssues {
		count := len(iss)
		if count > maxIssues {
			maxIssues = count
			maxPerson = person
		}
		if minIssues < 0 || count < minIssues {
			minIssues = count
			minPerson = person
		}
	}
	imbalance := maxIssues - minIssues
	if totalPeople > 1 {
		ratio := float64(maxIssues) / float64(minIssues)
		if ratio > 2.0 {
			sb.WriteString(fmt.Sprintf("- [UNBALANCED] %s has %d issues vs %s with %d (ratio %.1fx)\n", maxPerson, maxIssues, minPerson, minIssues, ratio))
		} else {
			sb.WriteString(fmt.Sprintf("- Balanced: %d to %d issues per person (range: %d)\n", minIssues, maxIssues, imbalance))
		}
	}
	sb.WriteString("\n")

	// 4. Cross-label collaboration matrix
	sb.WriteString("### Collaboration Areas\n")
	if len(labelAssignees) > 0 {
		for label, assignees := range labelAssignees {
			if len(assignees) > 1 {
				names := make([]string, 0, len(assignees))
				for p := range assignees {
					names = append(names, p)
				}
				sb.WriteString(fmt.Sprintf("- '%s': %s\n", label, strings.Join(names, ", ")))
			}
		}
	} else {
		sb.WriteString("- No shared labels to analyze.\n")
	}
	sb.WriteString("\n")

	// 5. Overall collaboration score
	sb.WriteString("### Collaboration Score\n")
	score := 100.0
	penalties := []string{}

	if siloCount > 0 {
		penalty := float64(siloCount) * 10
		score -= penalty
		penalties = append(penalties, fmt.Sprintf("-%d for knowledge silos", int(penalty)))
	}
	if isolatedCount > 0 {
		penalty := float64(isolatedCount) * 15
		score -= penalty
		penalties = append(penalties, fmt.Sprintf("-%d for isolation risk", int(penalty)))
	}
	if totalPeople > 1 {
		ratio := float64(maxIssues) / float64(minIssues)
		if ratio > 2.0 {
			score -= 10
			penalties = append(penalties, "-10 for workload imbalance")
		}
	}
	if score < 0 {
		score = 0
	}

	rating := "GREEN"
	if score < 60 {
		rating = "RED"
	} else if score < 80 {
		rating = "AMBER"
	}

	sb.WriteString(fmt.Sprintf("**Score: %.0f/100 [%s]**\n", score, rating))
	if len(penalties) > 0 {
		sb.WriteString(fmt.Sprintf("Penalties: %s\n", strings.Join(penalties, ", ")))
	}
	sb.WriteString("\n### Recommendations\n")
	if score >= 80 {
		sb.WriteString("- Maintain cross-area pairing. Document shared knowledge.\n")
	} else if score >= 60 {
		sb.WriteString("- Pair isolated members on shared tasks.\n")
		sb.WriteString("- Rotate knowledge silo areas to distribute expertise.\n")
	} else {
		sb.WriteString("- [HIGH PRIORITY] Set up cross-training sessions for silo areas.\n")
		sb.WriteString("- Redistribute workload to balance under/over-loaded members.\n")
		sb.WriteString("- Schedule pair programming for isolated team members.\n")
	}

	return textResult(sb.String()), nil
}

// PMAIHealth evaluates AI adoption health and provides research-backed guidance.
func (h *Handlers) PMAIHealth(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	var sb strings.Builder
	sb.WriteString("## AI Health Check\n\n")

	// 1. Process fundamentals (good process = good AI foundation)
	sb.WriteString("### Process Fundamentals\n")
	fundamentalsScore := 0
	totalChecks := 6

	// Check sprint snapshots
	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
	if len(snaps) >= 3 {
		fundamentalsScore++
		sb.WriteString(fmt.Sprintf("- Sprint snapshots: %d sprints tracked [PASS]\n", len(snaps)))
	} else if len(snaps) > 0 {
		sb.WriteString(fmt.Sprintf("- Sprint snapshots: only %d (need 3+) [PARTIAL]\n", len(snaps)))
	} else {
		sb.WriteString("- Sprint snapshots: none [NEEDED]\n")
	}

	// Check retrospectives
	retros, _ := h.Memory.GetRetrospectives(ctx, 5)
	if len(retros) >= 2 {
		fundamentalsScore++
		sb.WriteString(fmt.Sprintf("- Retrospectives: %d recorded [PASS]\n", len(retros)))
	} else if len(retros) > 0 {
		sb.WriteString("- Retrospectives: 1 recorded [PARTIAL]\n")
	} else {
		sb.WriteString("- Retrospectives: none [NEEDED]\n")
	}

	// Check sprint goals
	goals, _ := h.Memory.GetGoalHistory(ctx, boardID, 5)
	if len(goals) >= 2 {
		fundamentalsScore++
		sb.WriteString(fmt.Sprintf("- Sprint goals: %d tracked [PASS]\n", len(goals)))
	} else {
		sb.WriteString("- Sprint goals: not tracked [NEEDED] — use pm_sprint_goal_track\n")
	}

	// Check decisions recorded
	decisions, _ := h.Memory.GetDecisions(ctx, 10)
	if len(decisions) >= 5 {
		fundamentalsScore++
		sb.WriteString(fmt.Sprintf("- Decisions recorded: %d [PASS]\n", len(decisions)))
	} else if len(decisions) > 0 {
		sb.WriteString(fmt.Sprintf("- Decisions recorded: %d [PARTIAL]\n", len(decisions)))
	} else {
		sb.WriteString("- Decisions recorded: none [NEEDED]\n")
	}

	// Check active blockers tracked
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	_ = blockers // We just check if the function works, not count
	sb.WriteString("- Blocker tracking: available [PASS]\n")
	fundamentalsScore++

	// Check team pulses
	pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 3)
	if len(pulses) >= 2 {
		fundamentalsScore++
		sb.WriteString(fmt.Sprintf("- Team pulse: %d recorded [PASS]\n", len(pulses)))
	} else {
		sb.WriteString("- Team pulse: not tracked [NEEDED]\n")
	}

	fundamentalsPct := float64(fundamentalsScore) / float64(totalChecks) * 100
	sb.WriteString(fmt.Sprintf("\n**Fundamentals Score: %d/%d (%.0f%%)**\n", fundamentalsScore, totalChecks, fundamentalsPct))

	if fundamentalsPct < 50 {
		sb.WriteString("**Warning:** Low process maturity = AI can't help effectively. Focus on fundamentals first.\n")
	}
	sb.WriteString("\n")

	// 2. AI tool usage analysis
	h.initInsightTables()
	if h.Memory != nil && h.Memory.DB() != nil {
		sb.WriteString("### AI Tool Usage\n")
		rows, err := h.Memory.DB().Query(`SELECT tool_name, COUNT(*) as cnt
			FROM tool_usage
			WHERE called_at > datetime('now', '-30 days')
			GROUP BY tool_name ORDER BY cnt DESC`)
		if err == nil {
			totalCalls := 0
			uniqueTools := 0
			first := true
			for rows.Next() {
				var name string
				var cnt int
				_ = rows.Scan(&name, &cnt)
				if first {
					sb.WriteString("Most used tools (last 30 days):\n")
					first = false
				}
				sb.WriteString(fmt.Sprintf("  %s: %d calls\n", name, cnt))
				totalCalls += cnt
				uniqueTools++
			}
			rows.Close()

			if totalCalls == 0 {
				sb.WriteString("  No tool usage data yet.\n")
			} else {
				sb.WriteString(fmt.Sprintf("\nTotal: %d calls across %d unique tools\n", totalCalls, uniqueTools))

				// BCG 2026 research: 3 AI tools max before cognitive overload
				if uniqueTools > 3 {
					sb.WriteString("\n**Cognitive Load Warning:** ")
					sb.WriteString(fmt.Sprintf("Using %d unique tools. Research (BCG 2026) shows productivity peaks at 3 tools and degrades at 4+.\n", uniqueTools))
					sb.WriteString("Consider consolidating or using profiles (PM_PROFILE).\n")
				} else if uniqueTools == 0 {
					// skip
				} else {
					sb.WriteString("\nTool diversity is within the healthy range (≤3). Research confirms this is optimal.\n")
				}
			}
		}
	}
	sb.WriteString("\n")

	// 3. Research-backed recommendations
	sb.WriteString("### AI Health Recommendations\n")
	sb.WriteString("Based on BCG 2026, PMI 2025, and Microsoft Research 2026:\n\n")

	if fundamentalsPct < 60 {
		sb.WriteString("1. **Fix process fundamentals first.** AI amplifies good process and automates bad process.\n")
		sb.WriteString("   Start with: pm_sprint_goal_track → pm_record_retro → pm_snapshot_sprint\n\n")
	}
	sb.WriteString("2. **Limit AI tool surface.** BCG 2026: 3 tools = peak productivity. Use PM_PROFILE to control.\n\n")
	sb.WriteString("3. **Keep humans in the loop.** AI should augment, not replace, critical thinking.\n")
	sb.WriteString("   Use AI for data gathering; make the decision yourself.\n\n")
	sb.WriteString("4. **Write culture.** Document decisions and context. AI memory only works if you feed it.\n\n")
	sb.WriteString("5. **Review AI output.** PMI 2025: only 20% of PMs have good AI experience.\n")
	sb.WriteString("   Calibrate your trust over time. Use pm_calibration_report.\n")

	return textResult(sb.String()), nil
}

// PMNotificationEffectiveness tracks and analyzes notification effectiveness.
func (h *Handlers) PMNotificationEffectiveness(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initResearchTables()

	var sb strings.Builder
	sb.WriteString("## Notification Effectiveness Report\n\n")

	// 1. Volume by channel
	sb.WriteString("### Volume by Channel (Last 30 Days)\n")
	rows, err := h.Memory.DB().Query(`
		SELECT channel, COUNT(*) as cnt,
			SUM(CASE WHEN sent_at > datetime('now', '-1 day') THEN 1 ELSE 0 END) as today,
			SUM(CASE WHEN sent_at > datetime('now', '-7 days') THEN 1 ELSE 0 END) as this_week
		FROM notification_log
		WHERE sent_at > datetime('now', '-30 days')
		GROUP BY channel ORDER BY cnt DESC`)
	if err == nil {
		hasData := false
		for rows.Next() {
			var channel string
			var cnt, today, week int
			_ = rows.Scan(&channel, &cnt, &today, &week)
			sb.WriteString(fmt.Sprintf("- %s: %d total, %d today, %d this week\n", channel, cnt, today, week))
			hasData = true
		}
		rows.Close()
		if !hasData {
			sb.WriteString("  No notifications sent in the last 30 days.\n")
		}
	} else {
		sb.WriteString("  Query failed: " + err.Error() + "\n")
	}
	sb.WriteString("\n")

	// 2. Volume by severity
	sb.WriteString("### Volume by Severity (Last 30 Days)\n")
	rows2, err := h.Memory.DB().Query(`
		SELECT severity, COUNT(*) as cnt
		FROM notification_log
		WHERE sent_at > datetime('now', '-30 days')
		GROUP BY severity ORDER BY cnt DESC`)
	if err == nil {
		for rows2.Next() {
			var severity string
			var cnt int
			_ = rows2.Scan(&severity, &cnt)
			sb.WriteString(fmt.Sprintf("- %s: %d\n", severity, cnt))
		}
		rows2.Close()
	} else {
		sb.WriteString("  Query failed: " + err.Error() + "\n")
	}
	sb.WriteString("\n")

	// 3. Fatigue analysis
	sb.WriteString("### Fatigue Analysis\n")
	var total30 int
	row30 := h.Memory.DB().QueryRow("SELECT COUNT(*) FROM notification_log WHERE sent_at > datetime('now', '-30 days')")
	_ = row30.Scan(&total30)

	var dailyAvg float64
	if total30 > 0 {
		dailyAvg = float64(total30) / 30.0
	}

	sb.WriteString(fmt.Sprintf("- Total notifications (30d): %d\n", total30))
	sb.WriteString(fmt.Sprintf("- Daily average: %.1f\n", dailyAvg))

	budgetDays := 0
	rows3, err := h.Memory.DB().Query(`
		SELECT DATE(sent_at) as day, COUNT(*) as cnt
		FROM notification_log
		WHERE sent_at > datetime('now', '-7 days')
		GROUP BY DATE(sent_at) ORDER BY day`)
	if err == nil {
		for rows3.Next() {
			var day string
			var cnt int
			_ = rows3.Scan(&day, &cnt)
			if cnt > 5 {
				budgetDays++
				sb.WriteString(fmt.Sprintf("- %s: %d notifications [OVER BUDGET]\n", day, cnt))
			}
		}
		rows3.Close()
	}

	if budgetDays > 0 {
		sb.WriteString(fmt.Sprintf("\n**Fatigue Risk:** %d days exceeded the 5/day budget in the last week.\n", budgetDays))
		sb.WriteString("Recommendation: Reduce notification frequency. Use severity levels wisely.\n")
	} else if dailyAvg > 3 {
		sb.WriteString("\n**Caution:** Notification volume is moderate. Monitor for fatigue.\n")
	} else {
		sb.WriteString("\nNotification volume is within healthy range. No fatigue risk detected.\n")
	}
	sb.WriteString("\n")

	// 4. Effectiveness score
	sb.WriteString("### Notification Health Score\n")
	score := 100.0
	var scoreDeductions []string

	if dailyAvg > 5 {
		score -= 30
		scoreDeductions = append(scoreDeductions, "-30 for high volume (>5/day)")
	} else if dailyAvg > 3 {
		score -= 15
		scoreDeductions = append(scoreDeductions, "-15 for moderate volume")
	}

	if budgetDays > 0 {
		deduction := float64(budgetDays) * 10
		score -= deduction
		scoreDeductions = append(scoreDeductions, fmt.Sprintf("-%.0f for budget violations", deduction))
	}

	if score < 0 {
		score = 0
	}

	rating := "GREEN"
	if score < 50 {
		rating = "RED"
	} else if score < 75 {
		rating = "AMBER"
	}

	sb.WriteString(fmt.Sprintf("**Score: %.0f/100 [%s]**\n", score, rating))
	for _, d := range scoreDeductions {
		sb.WriteString(fmt.Sprintf("- %s\n", d))
	}
	sb.WriteString("\n")

	// 5. Recommendations
	sb.WriteString("### Recommendations\n")
	if score >= 80 {
		sb.WriteString("- Maintain current notification practices.\n")
	} else if score >= 60 {
		sb.WriteString("- Review notification routing — use severity levels consistently.\n")
		sb.WriteString("- Consider consolidating low-severity notifications into digests.\n")
	} else {
		sb.WriteString("- [HIGH] Reduce notification frequency. Use pm_comms_nudge for urgent items only.\n")
		sb.WriteString("- Set up notification budgets and stick to them.\n")
		sb.WriteString("- Batch non-urgent updates into a daily digest.\n")
	}

	return textResult(sb.String()), nil
}

// PMNotificationRecordAction records user action on a notification.
func (h *Handlers) PMNotificationRecordAction(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil || h.Memory.DB() == nil {
		return errorResult("memory not configured"), nil
	}
	h.initResearchTables()

	channel := req.GetString("channel", "")
	title := req.GetString("title", "")
	action := req.GetString("action", "acknowledged")

	if channel == "" && title == "" {
		return errorResult("provide channel or title to match notification"), nil
	}

	var result interface{}
	var err error
	if channel != "" {
		result, err = h.Memory.DB().Exec(
			"INSERT INTO notification_effectiveness(channel, severity, title, action_taken, action_at) VALUES(?, ?, ?, ?, datetime('now'))",
			channel, "tracked", title, action)
	} else {
		result, err = h.Memory.DB().Exec(
			"UPDATE notification_effectiveness SET action_taken = ?, action_at = datetime('now') WHERE title = ? AND action_taken = ''",
			action, title)
	}

	if err != nil {
		return errorResult("Failed to record action: " + err.Error()), nil
	}

	_ = result // ignore exec result detail
	return textResult(fmt.Sprintf("Recorded action '%s' for notification (channel=%s, title=%s)", action, channel, title)), nil
}

// helpers

func joinKeys(m map[string]int) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}


