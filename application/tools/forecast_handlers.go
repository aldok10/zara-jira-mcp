package tools

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// MonteCarloForecast runs Monte Carlo simulation to answer "when will it be done?"
func (h *Handlers) MonteCarloForecast(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}
	remainingItems := req.GetInt("remaining_items", 0)

	// Get historical throughput data
	snaps, err := h.Memory.GetSprintSnapshots(ctx, boardID, 10)
	if err != nil || len(snaps) < 3 {
		return textResult("Need at least 3 sprint snapshots for Monte Carlo forecasting."), nil
	}

	// If remaining not provided, get from current sprint
	if remainingItems == 0 {
		sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
		if len(sprints) > 0 {
			issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
			for _, issue := range issues {
				lower := strings.ToLower(issue.Status)
				if !strings.Contains(lower, "done") && !strings.Contains(lower, "closed") {
					remainingItems++
				}
			}
		}
	}
	if remainingItems == 0 {
		return textResult("No remaining items to forecast. Provide remaining_items or have an active sprint."), nil
	}

	// Historical throughput per sprint (items done)
	throughputs := make([]int, len(snaps))
	for i, s := range snaps {
		throughputs[i] = s.Done
	}

	// Run 10,000 Monte Carlo simulations
	simulations := 10000
	sprintDays := req.GetInt("sprint_days", 10)
	results := make([]int, simulations)

	for i := 0; i < simulations; i++ {
		remaining := remainingItems
		sprints := 0
		for remaining > 0 {
			// Random sample from historical throughput
			throughput := throughputs[rand.Intn(len(throughputs))]
			if throughput <= 0 {
				throughput = 1
			}
			remaining -= throughput
			sprints++
			if sprints > 50 { // safety cap
				break
			}
		}
		results[i] = sprints
	}

	sort.Ints(results)

	// Calculate percentiles
	p50 := results[simulations*50/100]
	p70 := results[simulations*70/100]
	p85 := results[simulations*85/100]
	p95 := results[simulations*95/100]

	// Convert to dates
	now := time.Now()
	dateP50 := now.AddDate(0, 0, p50*sprintDays)
	dateP70 := now.AddDate(0, 0, p70*sprintDays)
	dateP85 := now.AddDate(0, 0, p85*sprintDays)
	dateP95 := now.AddDate(0, 0, p95*sprintDays)

	var sb strings.Builder
	sb.WriteString("Monte Carlo Forecast\n\n")
	sb.WriteString(fmt.Sprintf("Remaining items: %d\n", remainingItems))
	sb.WriteString(fmt.Sprintf("Historical throughput: %v items/sprint (%d sprints sampled)\n", throughputs, len(throughputs)))
	sb.WriteString(fmt.Sprintf("Sprint length: %d days\n", sprintDays))
	sb.WriteString(fmt.Sprintf("Simulations: %d\n\n", simulations))

	sb.WriteString("When will it be done?\n")
	sb.WriteString(fmt.Sprintf("  50%% confidence: %d sprints (%s)\n", p50, dateP50.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("  70%% confidence: %d sprints (%s)\n", p70, dateP70.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("  85%% confidence: %d sprints (%s)  <- recommended commitment\n", p85, dateP85.Format("2006-01-02")))
	sb.WriteString(fmt.Sprintf("  95%% confidence: %d sprints (%s)\n", p95, dateP95.Format("2006-01-02")))

	sb.WriteString("\nInterpretation:\n")
	sb.WriteString(fmt.Sprintf("  Tell stakeholders: '%s' (85%% confidence)\n", dateP85.Format("Jan 2")))
	sb.WriteString(fmt.Sprintf("  Best case: '%s' (50%% confidence — coin flip)\n", dateP50.Format("Jan 2")))

	return textResult(sb.String()), nil
}

// DetectAntiPatterns scans data for Scrum anti-patterns.
func (h *Handlers) DetectAntiPatterns(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	var findings []string

	// Get sprint history
	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 6)

	// 1. Zombie Sprint: high carryover consistently (>30% carryover)
	if len(snaps) >= 3 {
		highCarryover := 0
		for _, s := range snaps[:3] {
			if s.TotalIssues > 0 && float64(s.Carryover)/float64(s.TotalIssues) > 0.3 {
				highCarryover++
			}
		}
		if highCarryover >= 2 {
			findings = append(findings, "ZOMBIE SPRINT: Team consistently carries over >30% of work. Sprint commitment is unrealistic or scope is growing mid-sprint. Fix: reduce commitment, protect sprint from interruptions.")
		}
	}

	// 2. Velocity Rollercoaster: high variance means unpredictability
	if len(snaps) >= 4 {
		var velocities []float64
		for _, s := range snaps {
			velocities = append(velocities, float64(s.Velocity))
		}
		mean := avg(velocities)
		stdDev := stddev(velocities, mean)
		cv := stdDev / mean * 100 // coefficient of variation
		if cv > 40 {
			findings = append(findings, fmt.Sprintf("UNPREDICTABLE: Velocity variation is %.0f%% (std dev: %.1f, mean: %.1f). Team cannot plan reliably. Fix: stabilize WIP, consistent sprint length, reduce external interruptions.", cv, stdDev, mean))
		}
	}

	// 3. Hero Culture: one person does most of the work
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
			var maxPerson string
			var maxDone int
			for person, count := range assigneeDone {
				if count > maxDone {
					maxDone = count
					maxPerson = person
				}
			}
			if float64(maxDone)/float64(totalDone) > 0.5 {
				findings = append(findings, fmt.Sprintf("HERO CULTURE: %s completed %d/%d items (%.0f%%). Team is over-dependent on one person. Bus factor = 1. Fix: pair programming, knowledge sharing, distribute ownership.", maxPerson, maxDone, totalDone, float64(maxDone)/float64(totalDone)*100))
			}
		}
	}

	// 4. Scope Creep: total issues growing mid-sprint
	if len(snaps) >= 2 {
		progress, _ := h.Memory.GetDailyProgress(ctx, boardID, snaps[0].SprintName)
		if len(progress) >= 3 {
			first := progress[0].TotalIssues
			last := progress[len(progress)-1].TotalIssues
			if last > first && float64(last-first)/float64(first) > 0.2 {
				findings = append(findings, fmt.Sprintf("SCOPE CREEP: Sprint started with %d items, now has %d (+%.0f%%). Sprint backlog is not being protected. Fix: strict change control, PM as gatekeeper, make scope changes visible.", first, last, float64(last-first)/float64(first)*100))
			}
		}
	}

	// 5. Rubber Stamping DoD: completion rate suspiciously high with no blockers
	if len(snaps) >= 3 {
		perfectSprints := 0
		for _, s := range snaps[:3] {
			if s.CompletionRate > 95 && s.Blocked == 0 {
				perfectSprints++
			}
		}
		if perfectSprints >= 2 {
			findings = append(findings, "TOO PERFECT: Multiple sprints with >95% completion and zero blockers. Either DoD is too low, stories are trivially small, or problems aren't being surfaced. Fix: challenge DoD, increase story complexity, create psychological safety for raising issues.")
		}
	}

	// 6. Retro Actions Never Done
	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 5 {
		findings = append(findings, fmt.Sprintf("DEAD RETROS: %d retro action items pending. Team discusses improvements but never implements them. Retros are becoming ceremony without change. Fix: limit to 2 actions per retro, track in sprint as first-class work, celebrate closed actions.", len(actions)))
	}

	// 7. No Sprint Goal or vague goals
	goals, _ := h.Memory.GetActiveGoals(ctx, boardID)
	goalHistory, _ := h.Memory.GetGoalHistory(ctx, boardID, 5)
	if len(goals) == 0 && len(goalHistory) == 0 {
		findings = append(findings, "NO SPRINT GOALS: Team has no recorded sprint goals. Without a goal, the sprint is just a time-box with a ticket list. Fix: set a clear, testable goal every sprint that connects tickets to business value.")
	}

	if len(findings) == 0 {
		return textResult("No anti-patterns detected. Team practices look healthy based on available data."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Scrum Anti-Patterns Detected: %d\n\n", len(findings)))
	for i, f := range findings {
		sb.WriteString(fmt.Sprintf("%d. %s\n\n", i+1, f))
	}

	return textResult(sb.String()), nil
}

// CoachingAdvice provides AI coaching suggestions based on team data.
func (h *Handlers) CoachingAdvice(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topic, err := req.RequireString("topic")
	if err != nil {
		return errorResult("topic required (team_dynamics, velocity, blockers, morale, conflict, growth)"), nil
	}

	boardID := req.GetInt("board_id", 0)

	// Gather context
	var contextData strings.Builder

	if boardID > 0 {
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
		if len(snaps) > 0 {
			contextData.WriteString("Sprint History:\n")
			for _, s := range snaps {
				contextData.WriteString(fmt.Sprintf("  %s: velocity=%d, completion=%.0f%%, blocked=%d, carryover=%d\n",
					s.SprintName, s.Velocity, s.CompletionRate, s.Blocked, s.Carryover))
			}
		}
	}

	risks, _ := h.Memory.GetOpenRisks(ctx)
	if len(risks) > 0 {
		contextData.WriteString(fmt.Sprintf("\nOpen risks: %d\n", len(risks)))
	}

	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) > 0 {
		contextData.WriteString(fmt.Sprintf("Active blockers: %d\n", len(blockers)))
	}

	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 0 {
		contextData.WriteString(fmt.Sprintf("Pending retro actions: %d\n", len(actions)))
	}

	situation := req.GetString("situation", "")
	if situation != "" {
		contextData.WriteString(fmt.Sprintf("\nSpecific situation: %s\n", situation))
	}

	systemPrompt := fmt.Sprintf(`You are an experienced Agile Coach / Scrum Master giving coaching advice.
Topic: %s

Provide advice that is:
1. Specific and actionable (not generic textbook answers)
2. Based on the team data provided
3. Focused on root causes, not symptoms
4. Includes one concrete experiment the SM can try THIS SPRINT
5. Acknowledges what's going well (if anything)

Format: 
- Observation (what you see in the data)
- Root cause hypothesis  
- Coaching approach (question to ask the team, not directive)
- Experiment to try
- What success looks like

Keep it under 200 words. Be direct.`, topic)

	result, err := h.aiComplete(ctx, systemPrompt, contextData.String())
	if err != nil {
		return errorResult("AI failed: " + err.Error()), nil
	}

	return textResult(result), nil
}

// ForecastSprint predicts what can be delivered this sprint using Monte Carlo on historical velocity.
func (h *Handlers) ForecastSprint(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 10)
	if len(snaps) < 3 {
		return textResult("Need at least 3 sprint snapshots for forecasting. Capture more with pm_snapshot_sprint."), nil
	}

	velocities := make([]int, len(snaps))
	for i, s := range snaps {
		velocities[i] = s.Done
	}

	itemsRemaining := req.GetInt("items_remaining", 0)
	if itemsRemaining == 0 {
		sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
		if len(sprints) > 0 {
			issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
			for _, issue := range issues {
				lower := strings.ToLower(issue.Status)
				if !strings.Contains(lower, "done") && !strings.Contains(lower, "closed") {
					itemsRemaining++
				}
			}
		}
	}
	if itemsRemaining == 0 {
		return textResult("No remaining items detected. Provide items_remaining manually."), nil
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	simulations := 1000
	completions := make([]int, simulations)
	for i := 0; i < simulations; i++ {
		completions[i] = velocities[rng.Intn(len(velocities))]
	}
	sort.Ints(completions)

	p15 := completions[int(float64(simulations)*0.15)]
	p50 := completions[simulations/2]
	p85 := completions[int(float64(simulations)*0.85)]

	canFinish50 := p50 >= itemsRemaining
	canFinish85 := p85 >= itemsRemaining

	var sb strings.Builder
	sb.WriteString("Sprint Forecast (Monte Carlo, 1000 simulations)\n\n")
	sb.WriteString(fmt.Sprintf("Items remaining: %d\n", itemsRemaining))
	sb.WriteString(fmt.Sprintf("Historical velocity: %v (last %d sprints)\n\n", velocities, len(velocities)))
	sb.WriteString(fmt.Sprintf("50%% confidence: team will complete %d items\n", p50))
	sb.WriteString(fmt.Sprintf("85%% confidence: team will complete at least %d items\n", p15))
	sb.WriteString(fmt.Sprintf("15%% upside: team could complete up to %d items\n\n", p85))

	if canFinish50 {
		sb.WriteString("Verdict: LIKELY ON TRACK (50%+ chance of completing all remaining work)\n")
	} else if canFinish85 {
		sb.WriteString("Verdict: STRETCH (need above-average sprint to finish all)\n")
	} else {
		sb.WriteString("Verdict: AT RISK (unlikely to finish all remaining items)\n")
		deficit := itemsRemaining - p50
		sb.WriteString(fmt.Sprintf("Consider removing %d items from sprint scope.\n", deficit))
	}

	return textResult(sb.String()), nil
}

// NLToJQL converts natural language to JQL using AI.
func (h *Handlers) NLToJQL(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	query, err := req.RequireString("query")
	if err != nil {
		return errorResult("query required"), nil
	}

	systemPrompt := `Convert natural language to Jira JQL. Return ONLY the JQL query on the first line, then a brief explanation on the second line.

JQL reference:
- assignee = currentUser() / assignee = "name"
- status = "In Progress" / status in ("To Do", "In Progress")
- priority in (High, Highest)
- issuetype = Bug / Story / Task / Epic
- project = KEY
- sprint in openSprints() / sprint in futureSprints()
- created >= -7d / updated >= startOfWeek()
- resolution = Unresolved
- labels = "label-name"
- ORDER BY priority DESC, created DESC

Examples:
- "my open bugs" -> assignee = currentUser() AND issuetype = Bug AND resolution = Unresolved
- "high priority tasks in current sprint" -> priority in (High, Highest) AND sprint in openSprints() AND resolution = Unresolved`

	result, err := h.aiComplete(ctx, systemPrompt, query)
	if err != nil {
		return errorResult("AI conversion failed: " + err.Error()), nil
	}
	return textResult(result), nil
}

// ScopeCreep detects mid-sprint scope additions by comparing current state to baseline snapshot.
func (h *Handlers) ScopeCreep(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil || len(sprints) == 0 {
		return textResult("No active sprint."), nil
	}

	issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
	currentTotal := len(issues)

	snap, _ := h.Memory.GetLatestSnapshot(ctx, boardID)
	if snap == nil {
		return textResult(fmt.Sprintf("Current sprint has %d issues. No baseline snapshot to compare. Use pm_snapshot_sprint at sprint start.", currentTotal)), nil
	}

	delta := currentTotal - snap.TotalIssues
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Scope Analysis: %s\n\n", sprints[0].Name))
	sb.WriteString(fmt.Sprintf("Baseline (snapshot): %d issues\n", snap.TotalIssues))
	sb.WriteString(fmt.Sprintf("Current: %d issues\n", currentTotal))
	sb.WriteString(fmt.Sprintf("Delta: %+d\n\n", delta))

	if delta > 0 {
		pct := float64(delta) / float64(snap.TotalIssues) * 100
		sb.WriteString(fmt.Sprintf("SCOPE CREEP DETECTED: +%.0f%% increase since sprint start\n", pct))
		if pct > 20 {
			sb.WriteString("WARNING: >20% increase - sprint goal likely at risk. Consider removing items.\n")
		}
	} else if delta < 0 {
		sb.WriteString("Scope REDUCED (items removed from sprint). Good discipline.\n")
	} else {
		sb.WriteString("Scope STABLE. No items added or removed.\n")
	}

	return textResult(sb.String()), nil
}

// BacklogGroom finds stale backlog items that need grooming or archiving.
func (h *Handlers) BacklogGroom(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	days := req.GetInt("days", 90)
	project := req.GetString("project", "")

	jql := fmt.Sprintf("resolution = Unresolved AND updated <= -%dd AND sprint not in openSprints() ORDER BY updated ASC", days)
	if project != "" {
		jql = fmt.Sprintf("project = %s AND resolution = Unresolved AND updated <= -%dd AND sprint not in openSprints() ORDER BY updated ASC", project, days)
	}

	result, err := h.Jira.SearchIssues(ctx, jql, 50, 0)
	if err != nil {
		return errorResult("Search failed: " + err.Error()), nil
	}

	if len(result.Issues) == 0 {
		return textResult(fmt.Sprintf("Backlog is clean! No items untouched for %d+ days.", days)), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Backlog Grooming: %d stale items (no update in %d+ days)\n\n", len(result.Issues), days))
	sb.WriteString("Consider archiving or re-prioritizing:\n\n")
	for _, issue := range result.Issues {
		daysSince := int(time.Since(issue.Updated).Hours() / 24)
		sb.WriteString(fmt.Sprintf("- %s [%s] %s (%d days stale, assignee: %s)\n",
			issue.Key, issue.Type, issue.Summary, daysSince, issue.Assignee))
	}
	sb.WriteString(fmt.Sprintf("\nTotal: %d items. If >100, your backlog needs a purge.\n", result.Total))

	return textResult(sb.String()), nil
}

// helpers

func avg(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

func stddev(vals []float64, mean float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += (v - mean) * (v - mean)
	}
	return math.Sqrt(sum / float64(len(vals)))
}
