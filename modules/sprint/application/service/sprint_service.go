// Package service implements sprint application service.
package service

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/aldok10/zara-jira-mcp/modules/sprint/application/port"
	memory "github.com/aldok10/zara-jira-mcp/modules/sprint/domain/memory"
	"github.com/aldok10/zara-jira-mcp/modules/sprint/domain/planning"
)

// Ensure service implements the interface at compile time.
var _ port.Inbound = (*sprintService)(nil)

// sprintService implements port.Inbound.
type sprintService struct {
	snapshots port.SnapshotRepository
	health    port.HealthRepository
	risks     port.RiskRepository
	blockers  port.BlockerRepository
	goals     port.GoalRepository
	daily     port.DailyProgressRepository
	workflow  port.WorkflowRepository
	jira      port.JiraClient
	ai        port.AIProvider
	events    port.EventBus
}

// CalculateHealth computes a 0-100 health score for the active sprint.
// Scores 4 dimensions (each 0-25): Velocity, Blocker, Scope, Team.
func (s *sprintService) CalculateHealth(ctx context.Context, boardID int) (*port.HealthResult, error) {
	if boardID <= 0 {
		return nil, fmt.Errorf("board_id is required")
	}

	// Get active sprint from Jira.
	sprints, err := s.jira.GetActiveSprints(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("get sprints: %w", err)
	}
	if len(sprints) == 0 {
		return nil, fmt.Errorf("no active sprint found for board %d", boardID)
	}
	activeSprint := sprints[0]

	// Get sprint issues for status breakdown.
	issues, err := s.jira.GetSprintIssues(ctx, activeSprint.ID)
	if err != nil {
		return nil, fmt.Errorf("get sprint issues: %w", err)
	}

	// Count by status (workflow-aware classification).
	var done, inProgress, todo, blocked int
	for _, iss := range issues {
		cls := s.classifyWithOverrides(ctx, boardID, iss.Status)
		switch cls {
		case "done":
			done++
		case "progress":
			inProgress++
		case "blocked":
			blocked++
		default:
			todo++
		}
	}
	total := len(issues)

	// VelocityScore (0-25): based on completion rate.
	var velocityScore int
	if total > 0 {
		rate := float64(done) / float64(total)
		velocityScore = int(math.Round(rate * 25))
	} else {
		velocityScore = 25 // empty sprint is trivially healthy
	}

	// BlockerScore (0-25): inverse of blocked ratio.
	var blockerScore int
	if total > 0 {
		blockedRatio := float64(blocked) / float64(total)
		blockerScore = int(math.Round((1 - blockedRatio) * 25))
	} else {
		blockerScore = 25
	}

	// ScopeScore (0-25): measure scope stability.
	// Compare current total vs recent snapshot total.
	scopeScore := 25 // default healthy
	if latest, err := s.snapshots.FindLatest(ctx, boardID); err == nil && latest != nil && latest.TotalIssues > 0 {
		// If current sprint has significantly more/fewer issues, scope is unstable.
		ratio := float64(total) / float64(latest.TotalIssues)
		penalty := math.Abs(1-ratio) * 100
		// Allow up to 50% drift before max penalty.
		scopeScore = 25 - int(math.Min(25, penalty))
		if scopeScore < 0 {
			scopeScore = 0
		}
	}

	// TeamScore (0-25): based on predictability and carryover from historical snapshots.
	teamScore := 20 // default slightly above midpoint
	if snaps, err := s.snapshots.FindByBoard(ctx, boardID, 5); err == nil && len(snaps) > 0 {
		var carryoverRates []float64
		for _, sn := range snaps {
			if sn.TotalIssues > 0 {
				carryoverRates = append(carryoverRates, float64(sn.Carryover)/float64(sn.TotalIssues))
			}
		}
		if len(carryoverRates) > 0 {
			avgCarryover := 0.0
			for _, c := range carryoverRates {
				avgCarryover += c
			}
			avgCarryover /= float64(len(carryoverRates))
			// Lower carryover = better score.
			teamScore = int(math.Round((1 - avgCarryover) * 25))
		}
	}

	overall := velocityScore + blockerScore + scopeScore + teamScore

	// Determine rating.
	rating := "Critical"
	switch {
	case overall >= 80:
		rating = "Healthy"
	case overall >= 60:
		rating = "Fair"
	case overall >= 40:
		rating = "At Risk"
	}

	// Find weakest dimension.
	scores := map[string]int{
		"Velocity": velocityScore,
		"Blocker":  blockerScore,
		"Scope":    scopeScore,
		"Team":     teamScore,
	}
	weakest := "Velocity"
	weakestVal := velocityScore
	for dim, val := range scores {
		if val < weakestVal {
			weakestVal = val
			weakest = dim
		}
	}

	result := &port.HealthResult{
		Score:         overall,
		Rating:        rating,
		WeakestDim:    weakest,
		SprintName:    activeSprint.Name,
		VelocityScore: velocityScore,
		BlockerScore:  blockerScore,
		ScopeScore:    scopeScore,
		TeamScore:     teamScore,
	}

	return result, nil
}

// Forecast predicts completion sprints using Monte Carlo simulation.
func (s *sprintService) Forecast(ctx context.Context, boardID int, remaining int) (*port.ForecastResult, error) {
	snaps, err := s.snapshots.FindByBoard(ctx, boardID, 20)
	if err != nil {
		return nil, fmt.Errorf("fetch snapshots: %w", err)
	}
	if len(snaps) < 3 {
		return nil, fmt.Errorf("need at least 3 sprint snapshots, got %d", len(snaps))
	}

	throughput := make([]float64, len(snaps))
	for i := range snaps {
		if snaps[i].Done <= 0 {
			throughput[i] = 1
		} else {
			throughput[i] = float64(snaps[i].Done)
		}
	}

	res := planning.Forecast(throughput, remaining, 0)

	return &port.ForecastResult{
		MeanSprints: res.MeanSprints,
		MinSprints:  res.MinSprints,
		MaxSprints:  res.MaxSprints,
		Percentiles: res.Percentiles,
		Remaining:   res.Remaining,
		Simulations: res.Simulations,
	}, nil
}

// FlowMetrics computes WIP, throughput, and inferred cycle time.
func (s *sprintService) FlowMetrics(ctx context.Context, boardID int) (*port.FlowMetricsResult, error) {
	if boardID <= 0 {
		return nil, fmt.Errorf("board_id is required")
	}

	// Get active sprint.
	sprints, err := s.jira.GetActiveSprints(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("get sprints: %w", err)
	}

	// Get sprint issues to count statuses.
	issues, err := s.jira.GetSprintIssues(ctx, sprints[0].ID)
	if err != nil {
		return nil, fmt.Errorf("get sprint issues: %w", err)
	}

	var wip, done, blocked int
	for _, iss := range issues {
		cls := s.classifyWithOverrides(ctx, boardID, iss.Status)
		switch cls {
		case "progress":
			wip++
		case "done":
			done++
		case "blocked":
			blocked++
		}
	}
	total := len(issues)

	// Throughput from historical snapshots.
	snaps, err := s.snapshots.FindByBoard(ctx, boardID, 10)
	if err != nil {
		return nil, fmt.Errorf("fetch snapshots: %w", err)
	}

	avgThroughput := 0.0
	if len(snaps) > 0 {
		var sum int
		for _, sn := range snaps {
			sum += sn.Done
		}
		avgThroughput = float64(sum) / float64(len(snaps))
	}

	// Inferred cycle time: WIP / (throughput / sprintDays)
	sprintDays := 10.0
	avgCycleTime := 0.0
	if avgThroughput > 0 && wip > 0 {
		avgCycleTime = float64(wip) / (avgThroughput / sprintDays)
	}

	completionPct := 0.0
	if total > 0 {
		completionPct = float64(done) / float64(total) * 100
	}

	// Trend: compare recent vs older throughput.
	trend := "stable"
	if len(snaps) >= 6 {
		mid := len(snaps) / 2
		recent := avgDone(snaps[:mid])
		older := avgDone(snaps[mid:])
		if recent > older*1.15 {
			trend = "up ↑"
		} else if recent < older*0.85 {
			trend = "down ↓"
		}
	}

	return &port.FlowMetricsResult{
		BoardID:       boardID,
		CurrentWIP:    wip,
		AvgThroughput: math.Round(avgThroughput*100) / 100,
		AvgCycleTime:  math.Round(avgCycleTime*100) / 100,
		CompletionPct: math.Round(completionPct*100) / 100,
		TotalIssues:   total,
		DoneIssues:    done,
		BlockedIssues: blocked,
		Trend:         trend,
	}, nil
}

// SprintCompare compares current vs previous sprint metrics.
func (s *sprintService) SprintCompare(ctx context.Context, boardID int) (string, error) {
	snaps, err := s.snapshots.FindByBoard(ctx, boardID, 3)
	if err != nil {
		return "", fmt.Errorf("fetch snapshots: %w", err)
	}
	if len(snaps) < 2 {
		return "Need at least 2 sprint snapshots to compare. Snapshot each sprint end to build history.", nil
	}

	curr := snaps[0]
	prev := snaps[1]

	pct := func(val, total int) string {
		if total <= 0 {
			return "N/A"
		}
		return fmt.Sprintf("%.0f%%", float64(val)/float64(total)*100)
	}

	delta := func(curr, prev int) string {
		diff := curr - prev
		if diff > 0 {
			return fmt.Sprintf("+%d ↑", diff)
		} else if diff < 0 {
			return fmt.Sprintf("%d ↓", diff)
		}
		return "0 →"
	}

	var sb strings.Builder
	sb.WriteString("# Sprint Comparison\n\n")
	sb.WriteString(fmt.Sprintf("| Metric | %s | %s | Δ |\n", curr.SprintName, prev.SprintName))
	sb.WriteString("|--------|------|------|----|\n")
	sb.WriteString(fmt.Sprintf("| Total Issues | %d | %d | %s |\n", curr.TotalIssues, prev.TotalIssues, delta(curr.TotalIssues, prev.TotalIssues)))
	sb.WriteString(fmt.Sprintf("| Done | %d (%s) | %d (%s) | %s |\n", curr.Done, pct(curr.Done, curr.TotalIssues), prev.Done, pct(prev.Done, prev.TotalIssues), delta(curr.Done, prev.Done)))
	sb.WriteString(fmt.Sprintf("| In Progress | %d (%s) | %d (%s) | %s |\n", curr.InProgress, pct(curr.InProgress, curr.TotalIssues), prev.InProgress, pct(prev.InProgress, prev.TotalIssues), delta(curr.InProgress, prev.InProgress)))
	sb.WriteString(fmt.Sprintf("| Blocked | %d (%s) | %d (%s) | %s |\n", curr.Blocked, pct(curr.Blocked, curr.TotalIssues), prev.Blocked, pct(prev.Blocked, prev.TotalIssues), delta(curr.Blocked, prev.Blocked)))
	sb.WriteString(fmt.Sprintf("| Carryover | %d (%s) | %d (%s) | %s |\n", curr.Carryover, pct(curr.Carryover, curr.TotalIssues), prev.Carryover, pct(prev.Carryover, prev.TotalIssues), delta(curr.Carryover, prev.Carryover)))
	sb.WriteString(fmt.Sprintf("| Velocity | %d pts | %d pts | %s |\n", curr.Velocity, prev.Velocity, delta(curr.Velocity, prev.Velocity)))

	return sb.String(), nil
}

// Predictability computes consistency of sprint completion over recent sprints.
func (s *sprintService) Predictability(ctx context.Context, boardID int) (string, error) {
	snaps, err := s.snapshots.FindByBoard(ctx, boardID, 10)
	if err != nil {
		return "", fmt.Errorf("fetch snapshots: %w", err)
	}
	if len(snaps) < 3 {
		return "Need at least 3 sprint snapshots. Snapshot each sprint end to build predictability history.", nil
	}

	// Calculate completion rates.
	var rates []float64
	var velocities []int
	for _, sn := range snaps {
		if sn.TotalIssues > 0 {
			rates = append(rates, float64(sn.Done)/float64(sn.TotalIssues)*100)
		}
		if sn.Velocity > 0 {
			velocities = append(velocities, sn.Velocity)
		}
	}

	if len(rates) < 2 {
		return "Not enough completion data to calculate predictability.", nil
	}

	// Average completion rate.
	avgRate := avgFloat64(rates)

	// Standard deviation (consistency).
	_, rateStdDev := stddev(rates)

	// Score: high completion + low variance = predictable.
	score := math.Min(100, avgRate*0.6+(100-rateStdDev)*0.4)
	if score < 0 {
		score = 0
	}

	rating := "Low"
	switch {
	case score >= 80:
		rating = "High"
	case score >= 60:
		rating = "Medium"
	}

	var sb strings.Builder
	sb.WriteString("# Sprint Predictability\n\n")
	sb.WriteString(fmt.Sprintf("**Score**: %.0f/100 — %s\n\n", score, rating))

	if len(velocities) >= 2 {
		_, velStdDev := stddev(float64Slice(velocities))
		avgVel := avgInts(velocities)
		consistencyPct := math.Max(0, 100-velStdDev/float64(avgVel)*100)
		sb.WriteString(fmt.Sprintf("Average velocity: **%.0f** pts/sprint\n", float64(avgVel)))
		sb.WriteString(fmt.Sprintf("Velocity consistency: **%.0f%%**\n\n", consistencyPct))
	}

	sb.WriteString("| Sprint | Done/Total | Rate |\n")
	sb.WriteString("|--------|------------|------|\n")
	for _, sn := range snaps {
		pct := 0.0
		if sn.TotalIssues > 0 {
			pct = float64(sn.Done) / float64(sn.TotalIssues) * 100
		}
		sb.WriteString(fmt.Sprintf("| %s | %d/%d | %.0f%% |\n", sn.SprintName, sn.Done, sn.TotalIssues, pct))
	}

	return sb.String(), nil
}

// Scorecard computes a 0-100 sprint scorecard from multiple dimensions.
func (s *sprintService) Scorecard(ctx context.Context, boardID int) (string, error) {
	snaps, err := s.snapshots.FindByBoard(ctx, boardID, 5)
	if err != nil {
		return "", fmt.Errorf("fetch snapshots: %w", err)
	}
	if len(snaps) == 0 {
		return "No sprint snapshots yet. Use pm_snapshot at sprint end to build history.", nil
	}

	latest := snaps[0]

	// 1. Completion (0-30): ratio of done to total
	completionScore := 0.0
	if latest.TotalIssues > 0 {
		rate := float64(latest.Done) / float64(latest.TotalIssues)
		completionScore = math.Min(30, rate*30)
	}

	// 2. Velocity vs average (0-25): how current velocity compares to historical avg
	var velocities []int
	for _, sn := range snaps {
		if sn.Velocity > 0 {
			velocities = append(velocities, sn.Velocity)
		}
	}
	velocityScore := 20.0 // default
	if len(velocities) >= 2 {
		avgVel := avgInts(velocities)
		if avgVel > 0 {
			ratio := float64(latest.Velocity) / avgVel
			velocityScore = math.Min(25, ratio*25)
		}
	}

	// 3. Blocked ratio (0-20): fewer blocked = better
	blockedScore := 15.0
	if latest.TotalIssues > 0 {
		blockedRatio := float64(latest.Blocked) / float64(latest.TotalIssues)
		blockedScore = math.Max(0, 20-(blockedRatio*20))
	}

	// 4. Carryover (0-15): lower carryover = better
	carryoverScore := 12.0
	if latest.TotalIssues > 0 {
		carryoverRatio := float64(latest.Carryover) / float64(latest.TotalIssues)
		carryoverScore = math.Max(0, 15-(carryoverRatio*15))
	}

	// 5. Predictability bonus (0-10): consistency across recent sprints
	predictabilityBonus := 5.0
	if len(snaps) >= 3 {
		var rates []float64
		for _, sn := range snaps {
			if sn.TotalIssues > 0 {
				rates = append(rates, float64(sn.Done)/float64(sn.TotalIssues))
			}
		}
		if len(rates) >= 2 {
			_, std := stddev(rates)
			// Lower std = more predictable
			predictabilityBonus = math.Max(0, 10-(std*10))
		}
	}

	total := completionScore + velocityScore + blockedScore + carryoverScore + predictabilityBonus

	rating := "Critical"
	switch {
	case total >= 80:
		rating = "Excellent"
	case total >= 65:
		rating = "Good"
	case total >= 45:
		rating = "Fair"
	}

	trend := "stable"
	if total >= 80 {
		trend = "green"
	} else if total >= 45 {
		trend = "amber"
	} else {
		trend = "red"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Scorecard: **%.0f/100** — %s [%s]\n\n", total, rating, trend))
	sb.WriteString("| Dimension | Score | Max |\n")
	sb.WriteString("|-----------|-------|-----|\n")
	sb.WriteString(fmt.Sprintf("| Completion | %.0f | 30 |\n", completionScore))
	sb.WriteString(fmt.Sprintf("| Velocity | %.0f | 25 |\n", velocityScore))
	sb.WriteString(fmt.Sprintf("| Blocked Ratio | %.0f | 20 |\n", blockedScore))
	sb.WriteString(fmt.Sprintf("| Carryover | %.0f | 15 |\n", carryoverScore))
	sb.WriteString(fmt.Sprintf("| Predictability | %.0f | 10 |\n", predictabilityBonus))
	sb.WriteString(fmt.Sprintf("\n**Based on sprint**: %s\n", latest.SprintName))
	sb.WriteString(fmt.Sprintf("**Done**: %d/%d (%.0f%%)\n", latest.Done, latest.TotalIssues, float64(latest.Done)/float64(latest.TotalIssues)*100))

	// Weakest area
	weakest := "Completion"
	weakestVal := completionScore
	for s, v := range map[string]float64{"Velocity": velocityScore, "Blocked Ratio": blockedScore, "Carryover": carryoverScore, "Predictability": predictabilityBonus} {
		if v < weakestVal {
			weakestVal = v
			weakest = s
		}
	}
	sb.WriteString(fmt.Sprintf("\n**Focus**: %s (lowest score)\n", weakest))

	return sb.String(), nil
}

// Calibration shows forecast accuracy — committed vs delivered over time.
func (s *sprintService) Calibration(ctx context.Context, boardID int) (string, error) {
	snaps, err := s.snapshots.FindByBoard(ctx, boardID, 20)
	if err != nil {
		return "", fmt.Errorf("fetch snapshots: %w", err)
	}
	if len(snaps) < 2 {
		return "Need at least 2 sprint snapshots for calibration. Use pm_snapshot after each sprint.", nil
	}

	var overcommits, underdelivers, onTarget int
	var totalDiffPct float64

	var sb strings.Builder
	sb.WriteString("# Sprint Calibration\n\n")
	sb.WriteString("| Sprint | Committed | Delivered | Δ | Accuracy |\n")
	sb.WriteString("|--------|-----------|----------|----|----------|\n")

	for _, sn := range snaps {
		if sn.TotalIssues <= 0 {
			continue
		}
		diff := sn.Done - sn.TotalIssues
		pct := float64(sn.Done) / float64(sn.TotalIssues) * 100
		totalDiffPct += pct

		delta := fmt.Sprintf("%+d", diff)
		if diff > 0 {
			delta += " ↑"
		} else if diff < 0 {
			delta += " ↓"
		} else {
			delta = "0 →"
		}

		accuracy := "✓"
		if pct < 80 {
			underdelivers++
			accuracy = "✗ under"
		} else if pct > 120 {
			overcommits++
			accuracy = "✗ over"
		} else {
			onTarget++
			accuracy = "✓ on target"
		}

		sb.WriteString(fmt.Sprintf("| %s | %d | %d | %s | %.0f%% %s |\n",
			sn.SprintName, sn.TotalIssues, sn.Done, delta, pct, accuracy))
	}

	n := float64(len(snaps))
	avgAccuracy := totalDiffPct / n

	sb.WriteString(fmt.Sprintf("\n**Average accuracy**: %.0f%%\n", avgAccuracy))
	sb.WriteString(fmt.Sprintf("**On target**: %d/%d (%.0f%%)\n", onTarget, len(snaps), float64(onTarget)/n*100))
	sb.WriteString(fmt.Sprintf("**Under-delivered**: %d\n", underdelivers))
	sb.WriteString(fmt.Sprintf("**Over-committed**: %d\n", overcommits))

	rating := "Good"
	if avgAccuracy < 80 || float64(onTarget)/n < 0.5 {
		rating = "Needs Improvement"
	} else if avgAccuracy < 90 {
		rating = "Fair"
	}

	sb.WriteString(fmt.Sprintf("\n**Overall calibration**: %s\n", rating))

	return sb.String(), nil
}

// TrackDaily captures today's sprint progress from Jira and saves it as DailyProgress.
func (s *sprintService) TrackDaily(ctx context.Context, boardID int) (string, error) {
	if boardID <= 0 {
		return "", fmt.Errorf("board_id is required")
	}

	sprints, err := s.jira.GetActiveSprints(ctx, boardID)
	if err != nil {
		return "", fmt.Errorf("get active sprints: %w", err)
	}
	if len(sprints) == 0 {
		return "No active sprint found for this board.", nil
	}
	active := sprints[0]

	issues, err := s.jira.GetSprintIssues(ctx, active.ID)
	if err != nil {
		return "", fmt.Errorf("get sprint issues: %w", err)
	}

	var done, inProgress, todo, blocked int
	for _, iss := range issues {
		cls := s.classifyWithOverrides(ctx, boardID, iss.Status)
		switch cls {
		case "done":
			done++
		case "progress":
			inProgress++
		case "blocked":
			blocked++
		default:
			todo++
		}
	}
	total := len(issues)

	progress := &memory.DailyProgress{
		SprintName:  active.Name,
		BoardID:     boardID,
		Date:        time.Now(),
		TotalIssues: total,
		Done:        done,
		InProgress:  inProgress,
		Todo:        todo,
		Blocked:     blocked,
	}

	if err := s.daily.Save(ctx, progress); err != nil {
		return "", fmt.Errorf("save daily progress: %w", err)
	}

	pct := 0.0
	if total > 0 {
		pct = float64(done) / float64(total) * 100
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Daily Sprint Tracking — %s\n\n", active.Name))
	sb.WriteString(fmt.Sprintf("**Date**: %s\n\n", progress.Date.Format("2006-01-02")))
	sb.WriteString("| Status | Count |\n")
	sb.WriteString("|--------|-------|\n")
	sb.WriteString(fmt.Sprintf("| Total | %d |\n", total))
	sb.WriteString(fmt.Sprintf("| Done | %d (%.0f%%) |\n", done, pct))
	sb.WriteString(fmt.Sprintf("| In Progress | %d |\n", inProgress))
	sb.WriteString(fmt.Sprintf("| To Do | %d |\n", todo))
	sb.WriteString(fmt.Sprintf("| Blocked | %d |\n", blocked))

	return sb.String(), nil
}

// Burndown shows the sprint burndown chart from daily progress data.
func (s *sprintService) Burndown(ctx context.Context, boardID int) (string, error) {
	if boardID <= 0 {
		return "", fmt.Errorf("board_id is required")
	}

	// Get active sprint name.
	sprints, err := s.jira.GetActiveSprints(ctx, boardID)
	if err != nil {
		return "", fmt.Errorf("get active sprints: %w", err)
	}
	sprintName := "current"
	if len(sprints) > 0 {
		sprintName = sprints[0].Name
	}

	entries, err := s.daily.FindByBoardAndSprint(ctx, boardID, sprintName)
	if err != nil {
		return "", fmt.Errorf("get daily progress: %w", err)
	}
	if len(entries) == 0 {
		return fmt.Sprintf("No daily tracking data for sprint '%s'. Use pm_track_daily to record progress each day.", sprintName), nil
	}

	// Calculate ideal burndown (linear from total to 0).
	first := entries[0]
	last := entries[len(entries)-1]
	sprintDays := len(entries)
	idealPerDay := float64(first.TotalIssues) / float64(sprintDays)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Burndown — %s\n\n", sprintName))
	sb.WriteString(fmt.Sprintf("**Sprint length**: %d days tracked\n", sprintDays))
	sb.WriteString(fmt.Sprintf("**Starting total**: %d issues\n", first.TotalIssues))
	sb.WriteString(fmt.Sprintf("**Remaining**: %d issues\n", last.TotalIssues-last.Done))
	sb.WriteString(fmt.Sprintf("**Completed**: %d issues\n\n", last.Done))

	sb.WriteString("```\n")
	sb.WriteString("Day  Start  Ideal  Actual  Remaining\n")
	sb.WriteString("───  ─────  ─────  ────── ─────────\n")

	for i, e := range entries {
		remaining := e.TotalIssues - e.Done
		idealRemaining := first.TotalIssues - int(idealPerDay*float64(i+1))
		if idealRemaining < 0 {
			idealRemaining = 0
		}

		day := i + 1
		bar := ""
		if remaining > 0 {
			barLen := remaining * 20 / first.TotalIssues
			if barLen > 20 {
				barLen = 20
			}
			for j := 0; j < barLen; j++ {
				bar += "█"
			}
		}

		sb.WriteString(fmt.Sprintf("Day %-2d  %-5d  %-5d  %-6d %s\n",
			day, first.TotalIssues, idealRemaining, remaining, bar))
	}
	sb.WriteString("```\n")

	// Assessment
	onTrack := last.Done >= first.TotalIssues-int(idealPerDay*float64(sprintDays))
	if onTrack {
		sb.WriteString("\n✅ **On track** to complete sprint scope.\n")
	} else {
		sb.WriteString("\n⚠️ **Behind pace**. Consider reducing scope or addressing blockers.\n")
	}

	return sb.String(), nil
}

// LearnWorkflow scans all unique statuses on a board, classifies them, and stores patterns.
func (s *sprintService) LearnWorkflow(ctx context.Context, boardID int) (string, error) {
	if boardID <= 0 {
		return "", fmt.Errorf("board_id is required")
	}

	// Get active sprint issues to discover statuses.
	sprints, err := s.jira.GetActiveSprints(ctx, boardID)
	if err != nil {
		return "", fmt.Errorf("get active sprints: %w", err)
	}
	if len(sprints) == 0 {
		return "No active sprint found. Try with a board that has an active sprint.", nil
	}

	issues, err := s.jira.GetSprintIssues(ctx, sprints[0].ID)
	if err != nil {
		return "", fmt.Errorf("get sprint issues: %w", err)
	}

	// Collect unique status names.
	seen := make(map[string]bool)
	var statuses []string
	for _, iss := range issues {
		status := iss.Status
		if status != "" && !seen[status] {
			seen[status] = true
			statuses = append(statuses, status)
		}
	}

	if len(statuses) == 0 {
		return "No issues found in active sprint to learn from.", nil
	}

	// Delete old patterns for this board.
	if err := s.workflow.DeleteByBoard(ctx, boardID); err != nil {
		return "", fmt.Errorf("clear old patterns: %w", err)
	}

	// Classify each status and save.
	var saved int
	var sb strings.Builder
	sb.WriteString("# Board Workflow Learned\n\n")
	sb.WriteString(fmt.Sprintf("Board ID: **%d**\n", boardID))
	sb.WriteString(fmt.Sprintf("Sprint: **%s**\n", sprints[0].Name))
	sb.WriteString(fmt.Sprintf("Statuses found: **%d**\n\n", len(statuses)))
	sb.WriteString("| Status | Classification | Matched By |\n")
	sb.WriteString("|--------|---------------|------------|\n")

	for _, status := range statuses {
		classification, pattern := classifyStatus(status)
		p := &memory.WorkflowPattern{
			BoardID:        boardID,
			StatusName:     status,
			Classification: classification,
			Pattern:        pattern,
			IsAuto:         true,
		}
		if err := s.workflow.Upsert(ctx, p); err != nil {
			return "", fmt.Errorf("save pattern for %q: %w", status, err)
		}
		saved++
		sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", status, classification, pattern))
	}

	sb.WriteString(fmt.Sprintf("\n✅ **%d** patterns saved to database. Use `pm_workflow_map` to view.\n", saved))

	// Show summary.
	var doneCt, blockedCt, progressCt, todoCt int
	for _, c := range statuses {
		cls, _ := classifyStatus(c)
		switch cls {
		case "done":
			doneCt++
		case "blocked":
			blockedCt++
		case "progress":
			progressCt++
		default:
			todoCt++
		}
	}
	sb.WriteString(fmt.Sprintf("\n**Summary**: %d done · %d blocked · %d in progress · %d todo\n", doneCt, blockedCt, progressCt, todoCt))

	return sb.String(), nil
}

// classifyStatus applies regex keyword matching to classify a Jira status.
// Returns (classification, matchedPattern).
func classifyStatus(status string) (string, string) {
	lower := strings.ToLower(status)

	// Check done keywords
	doneKeywords := []string{"done", "closed", "resolved", "complete", "finish", "merged", "deployed", "released"}
	for _, kw := range doneKeywords {
		if strings.Contains(lower, kw) {
			return "done", kw
		}
	}

	// Check blocked keywords
	blockedKeywords := []string{"blocked", "block", "impediment", "waiting", "stuck"}
	for _, kw := range blockedKeywords {
		if strings.Contains(lower, kw) {
			return "blocked", kw
		}
	}

	// Check in-progress keywords
	progressKeywords := []string{"progress", "review", "testing", "dev", "development",
		"implement", "working", "in-progress", "in progress", "wip", "code"}
	for _, kw := range progressKeywords {
		if strings.Contains(lower, kw) {
			return "progress", kw
		}
	}

	// Default to todo
	return "todo", "default"
}

// classifyWithOverrides classifies a status using DB patterns first, then keyword fallback.
func (s *sprintService) classifyWithOverrides(ctx context.Context, boardID int, status string) string {
	// If we have DB patterns, try them first.
	if s.workflow != nil {
		patterns, err := s.workflow.FindByBoard(ctx, boardID)
		if err == nil {
			for _, p := range patterns {
				if strings.EqualFold(p.StatusName, status) {
					return p.Classification
				}
			}
		}
	}
	// Fallback to keyword matching.
	cls, _ := classifyStatus(status)
	return cls
}

// VelocityTrend returns the velocity direction over recent sprints.
func (s *sprintService) VelocityTrend(ctx context.Context, boardID int) (string, error) {
	snaps, err := s.snapshots.FindByBoard(ctx, boardID, 10)
	if err != nil {
		return "", fmt.Errorf("fetch snapshots: %w", err)
	}
	if len(snaps) < 2 {
		return "Not enough data (need at least 2 sprint snapshots)", nil
	}

	// Extract completed items per sprint (newest first from FindByBoard).
	completed := make([]int, len(snaps))
	for i, sn := range snaps {
		completed[i] = sn.Done
	}

	// Calculate simple trend: compare last 3 vs previous 3 (or fewer).
	n := len(completed)
	recent := completed[:int(math.Min(3, float64(n)))]
	older := completed[int(math.Min(3, float64(n))):]

	recentAvg := avgInts(recent)
	olderAvg := avgInts(older)

	trend := "stable"
	direction := "→"
	change := recentAvg - olderAvg
	if change > 1.5 {
		trend = "improving"
		direction = "↑"
	} else if change < -1.5 {
		trend = "declining"
		direction = "↓"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Velocity is **%s** %s\n", trend, direction))
	sb.WriteString(fmt.Sprintf("Recent avg: **%.1f** items/sprint\n", recentAvg))
	sb.WriteString(fmt.Sprintf("Previous avg: **%.1f** items/sprint\n", olderAvg))
	sb.WriteString(fmt.Sprintf("Change: **%+.1f** items\n", change))
	sb.WriteString(fmt.Sprintf("Based on %d sprint snapshots\n\n", n))
	sb.WriteString("| Sprint | Completed |\n")
	sb.WriteString("|--------|----------|\n")
	for i := len(completed) - 1; i >= 0; i-- {
		sb.WriteString(fmt.Sprintf("| Sprint %d | %d |\n", i+1, completed[i]))
	}

	return sb.String(), nil
}

// DetectAntiPatterns scans for known anti-patterns in sprint execution.
func (s *sprintService) DetectAntiPatterns(ctx context.Context, boardID int) ([]port.AntiPattern, error) {
	var patterns []port.AntiPattern

	snaps, err := s.snapshots.FindByBoard(ctx, boardID, 10)
	if err != nil {
		return nil, fmt.Errorf("fetch snapshots: %w", err)
	}
	blockers, _ := s.blockers.FindActive(ctx)

	// 1. Zombie Sprint: carryover > 30%.
	if len(snaps) > 0 {
		latest := snaps[0]
		if latest.TotalIssues > 0 {
			carryoverRate := float64(latest.Carryover) / float64(latest.TotalIssues) * 100
			if carryoverRate > 30 {
				patterns = append(patterns, port.AntiPattern{
					Name:        "Zombie Sprint",
					Description: fmt.Sprintf("Carryover is %.0f%% — items not completing sprint-to-sprint. Sustainable pace is broken.", carryoverRate),
					Severity:    "High",
					Suggestion:  "Reduce WIP, split large items, and ensure every sprint delivers value. Consider a 'cleanup sprint' to zero out carryover.",
				})
			}
		}
	}

	// 2. Scope Creep: significant change in total issues across sprints.
	if len(snaps) >= 2 {
		latest := snaps[0]
		prev := snaps[1]
		if latest.TotalIssues > 0 && prev.TotalIssues > 0 {
			change := math.Abs(float64(latest.TotalIssues-prev.TotalIssues)) / float64(prev.TotalIssues) * 100
			if change > 50 {
				patterns = append(patterns, port.AntiPattern{
					Name:        "Scope Creep",
					Description: fmt.Sprintf("Issue count changed by %.0f%% between recent sprints — inconsistent scope.", change),
					Severity:    "Medium",
					Suggestion:  "Stabilize sprint scope during planning. Use Definition of Ready before committing items.",
				})
			}
		}
	}

	// 3. Blocked Items
	if len(blockers) > 0 {
		severity := "Medium"
		if len(blockers) >= 5 {
			severity = "High"
		}
		patterns = append(patterns, port.AntiPattern{
			Name:        "Blocked Items",
			Description: fmt.Sprintf("%d active blockers in the sprint. Blocked items degrade predictability.", len(blockers)),
			Severity:    severity,
			Suggestion:  "Review each blocker during daily standup. Assign clear owners. Escalate if unresolved >48h.",
		})
	}

	// 4. Completion inconsistency (feast-or-famine).
	if len(snaps) >= 3 {
		completions := make([]float64, len(snaps))
		for i, sn := range snaps {
			if sn.TotalIssues > 0 {
				completions[i] = float64(sn.Done) / float64(sn.TotalIssues) * 100
			}
		}
		mean, stdDev := stddev(completions)
		if stdDev > 20 && mean > 0 {
			patterns = append(patterns, port.AntiPattern{
				Name:        "Inconsistent Delivery",
				Description: fmt.Sprintf("Completion rate std deviation is %.0f%% (mean %.0f%%) — feast-or-famine pattern.", stdDev, mean),
				Severity:    "Medium",
				Suggestion:  "Standardize sprint length. Break work into smaller, more uniform items. Use story points for better estimation.",
			})
		}
	}

	// 5. Slipping velocity (declining trend).
	if len(snaps) >= 4 {
		half := len(snaps) / 2
		recent := snaps[:half]
		older := snaps[half:]
		recentDone := avgFloat64(doneSlice(recent))
		olderDone := avgFloat64(doneSlice(older))
		if recentDone < olderDone*0.7 {
			patterns = append(patterns, port.AntiPattern{
				Name:        "Declining Velocity",
				Description: fmt.Sprintf("Recent sprints average %.1f done vs %.1f previously — velocity dropped >30%%.", recentDone, olderDone),
				Severity:    "High",
				Suggestion:  "Investigate root cause: team changes, technical debt, scope inflation, or morale. Run a retrospective focused on delivery.",
			})
		}
	}

	if len(patterns) == 0 {
		return nil, nil
	}

	// Sort by severity.
	sort.Slice(patterns, func(i, j int) bool {
		order := map[string]int{"High": 0, "Medium": 1, "Low": 2}
		return order[patterns[i].Severity] < order[patterns[j].Severity]
	})

	return patterns, nil
}

// --- helpers ---

func avgInts(nums []int) float64 {
	if len(nums) == 0 {
		return 0
	}
	sum := 0
	for _, n := range nums {
		sum += n
	}
	return float64(sum) / float64(len(nums))
}

func avgFloat64(nums []float64) float64 {
	if len(nums) == 0 {
		return 0
	}
	sum := 0.0
	for _, n := range nums {
		sum += n
	}
	return sum / float64(len(nums))
}

func avgDone(snaps []memory.SprintSnapshot) float64 {
	if len(snaps) == 0 {
		return 0
	}
	var sum int
	for _, sn := range snaps {
		sum += sn.Done
	}
	return float64(sum) / float64(len(snaps))
}

func float64Slice(nums []int) []float64 {
	out := make([]float64, len(nums))
	for i, n := range nums {
		out[i] = float64(n)
	}
	return out
}

func doneSlice(snaps []memory.SprintSnapshot) []float64 {
	out := make([]float64, len(snaps))
	for i, sn := range snaps {
		out[i] = float64(sn.Done)
	}
	return out
}

func stddev(nums []float64) (mean, std float64) {
	if len(nums) == 0 {
		return 0, 0
	}
	for _, n := range nums {
		mean += n
	}
	mean /= float64(len(nums))

	var variance float64
	for _, n := range nums {
		diff := n - mean
		variance += diff * diff
	}
	variance /= float64(len(nums))
	return mean, math.Sqrt(variance)
}

// NewSprintService creates a new port.Inbound with its dependencies.
func NewSprintService(
	snapshots port.SnapshotRepository,
	health port.HealthRepository,
	risks port.RiskRepository,
	blockers port.BlockerRepository,
	goals port.GoalRepository,
	daily port.DailyProgressRepository,
	workflow port.WorkflowRepository,
	jira port.JiraClient,
	ai port.AIProvider,
	events port.EventBus,
) port.Inbound {
	return &sprintService{
		snapshots: snapshots,
		health:    health,
		risks:     risks,
		blockers:  blockers,
		goals:     goals,
		daily:     daily,
		workflow:  workflow,
		jira:      jira,
		ai:        ai,
		events:    events,
	}
}
