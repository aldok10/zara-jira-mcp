package tools

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handlers) TeamMaturityAssess(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 10)
	blockers, _ := h.Memory.GetBlockerHistory(ctx, 30)
	risks, _ := h.Memory.GetOpenRisks(ctx)

	var signals []string
	var maturityScore int

	if len(snaps) >= 4 {
		sum := 0
		for _, s := range snaps {
			sum += s.Done
		}
		avg := sum / len(snaps)
		variance := 0
		for _, s := range snaps {
			diff := s.Done - avg
			variance += diff * diff
		}
		variance /= len(snaps)
		if variance < avg {
			signals = append(signals, "STABLE VELOCITY: Team delivers predictably (Performing indicator)")
			maturityScore += 25
		} else {
			signals = append(signals, "VOLATILE VELOCITY: Delivery unpredictable (Storming/Norming indicator)")
			maturityScore += 10
		}
	} else {
		signals = append(signals, "INSUFFICIENT DATA: Need 4+ sprint snapshots for velocity analysis")
		maturityScore += 5
	}

	if len(blockers) > 0 {
		totalDays := 0
		resolved := 0
		for _, b := range blockers {
			if b.ResolvedAt != nil {
				totalDays += b.DaysBlocked
				resolved++
			}
		}
		if resolved > 0 {
			avgDays := totalDays / resolved
			if avgDays <= 2 {
				signals = append(signals, fmt.Sprintf("FAST BLOCKER RESOLUTION: avg %d days (Performing)", avgDays))
				maturityScore += 25
			} else if avgDays <= 5 {
				signals = append(signals, fmt.Sprintf("MODERATE BLOCKER RESOLUTION: avg %d days (Norming)", avgDays))
				maturityScore += 15
			} else {
				signals = append(signals, fmt.Sprintf("SLOW BLOCKER RESOLUTION: avg %d days (Storming)", avgDays))
				maturityScore += 5
			}
		}
	} else {
		maturityScore += 15
	}

	if len(risks) > 0 {
		signals = append(signals, fmt.Sprintf("RISK AWARENESS: %d active risks tracked (proactive team)", len(risks)))
		maturityScore += 20
	} else if len(snaps) > 3 {
		signals = append(signals, "NO RISKS TRACKED: Either no risks exist or team isn't identifying them")
		maturityScore += 5
	}

	if len(snaps) >= 3 {
		highCompletion := 0
		for _, s := range snaps {
			if s.CompletionRate >= 80 {
				highCompletion++
			}
		}
		ratio := float64(highCompletion) / float64(len(snaps))
		if ratio >= 0.7 {
			signals = append(signals, "HIGH COMPLETION CONSISTENCY: Team regularly hits >80% (Performing)")
			maturityScore += 30
		} else if ratio >= 0.4 {
			signals = append(signals, "MODERATE COMPLETION: Team sometimes misses commitments (Norming)")
			maturityScore += 15
		} else {
			signals = append(signals, "LOW COMPLETION: Team frequently overcommits (Forming/Storming)")
			maturityScore += 5
		}
	}

	stage := "FORMING"
	coaching := "Direct: Provide clear structure, define roles, set explicit expectations."
	if maturityScore >= 80 {
		stage = "PERFORMING"
		coaching = "Delegate: Challenge with stretch goals, reduce ceremony overhead, focus on continuous improvement."
	} else if maturityScore >= 55 {
		stage = "NORMING"
		coaching = "Support: Reinforce good patterns, address remaining friction points, build autonomy."
	} else if maturityScore >= 30 {
		stage = "STORMING"
		coaching = "Coach: Facilitate conflict resolution, clarify shared goals, protect psychological safety."
	}

	var sb strings.Builder
	sb.WriteString("Team Maturity Assessment\n\n")
	sb.WriteString(fmt.Sprintf("Stage: %s (score: %d/100)\n", stage, maturityScore))
	sb.WriteString(fmt.Sprintf("SM Stance: %s\n\n", coaching))
	sb.WriteString("Evidence:\n")
	for _, s := range signals {
		sb.WriteString(fmt.Sprintf("  - %s\n", s))
	}
	sb.WriteString(fmt.Sprintf("\nRecommendation: Adjust your leadership style to %s level.\n", stage))

	return textResult(sb.String()), nil
}

func (h *Handlers) ImprovementVelocity(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	retros, _ := h.Memory.GetRetrospectives(ctx, 10)
	actions, _ := h.Memory.GetPendingActionItems(ctx)

	if len(retros) == 0 {
		return textResult("No retrospectives recorded. Use pm_record_retro to start tracking improvement."), nil
	}

	pending := len(actions)
	totalRetros := len(retros)

	var sb strings.Builder
	sb.WriteString("Improvement Velocity Report\n\n")
	sb.WriteString(fmt.Sprintf("Retrospectives recorded: %d\n", totalRetros))
	sb.WriteString(fmt.Sprintf("Pending action items: %d\n", pending))

	if pending > totalRetros {
		sb.WriteString("\nWARNING: Action items accumulating faster than being resolved.\n")
		sb.WriteString("This is a 'Dead Retro' anti-pattern. Actions are identified but never executed.\n")
		sb.WriteString("Recommendation: Limit to 1 action item per retro. Make it a sprint backlog item.\n")
	} else if pending == 0 && totalRetros > 0 {
		sb.WriteString("\nGood: All action items resolved (or none created).\n")
		sb.WriteString("Check: Are retros generating genuine improvement actions?\n")
	}

	themes := map[string]int{}
	for _, r := range retros {
		if r.Improvements != "" {
			words := strings.Fields(strings.ToLower(r.Improvements))
			for _, w := range words {
				if len(w) > 5 {
					themes[w]++
				}
			}
		}
	}
	var repeating []string
	for word, count := range themes {
		if count >= 3 {
			repeating = append(repeating, fmt.Sprintf("'%s' (%dx)", word, count))
		}
	}
	if len(repeating) > 0 {
		sb.WriteString(fmt.Sprintf("\nRecurring themes (same issue keeps appearing): %s\n", strings.Join(repeating, ", ")))
		sb.WriteString("These indicate systemic issues that retro actions aren't fixing.\n")
	}

	return textResult(sb.String()), nil
}

func (h *Handlers) MeetingROI(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)
	teamSize := req.GetInt("team_size", 5)
	sprintDays := req.GetInt("sprint_days", 10)
	hourlyRate := req.GetInt("hourly_rate", 50)

	standupHours := float64(sprintDays) * 0.25 * float64(teamSize)
	planningHours := 2.0 * float64(teamSize)
	reviewHours := 1.0 * float64(teamSize)
	retroHours := 1.5 * float64(teamSize)
	groomingHours := 1.0 * float64(teamSize)
	totalCeremonyHours := standupHours + planningHours + reviewHours + retroHours + groomingHours

	totalSprintHours := float64(sprintDays) * 8 * float64(teamSize)
	ceremonyPercent := (totalCeremonyHours / totalSprintHours) * 100
	ceremonyCost := totalCeremonyHours * float64(hourlyRate)

	var throughputNote string
	if boardID > 0 {
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 3)
		if len(snaps) > 0 {
			avgCompletion := 0.0
			for _, s := range snaps {
				avgCompletion += s.CompletionRate
			}
			avgCompletion /= float64(len(snaps))
			if avgCompletion >= 80 {
				throughputNote = fmt.Sprintf("Sprint completion: %.0f%% - ceremonies appear effective.", avgCompletion)
			} else {
				throughputNote = fmt.Sprintf("Sprint completion: %.0f%% - ceremonies may not be translating to delivery.", avgCompletion)
			}
		}
	}

	var sb strings.Builder
	sb.WriteString("Meeting ROI Analysis\n\n")
	sb.WriteString(fmt.Sprintf("Team size: %d | Sprint: %d days | Rate: $%d/hr\n\n", teamSize, sprintDays, hourlyRate))
	sb.WriteString("Ceremony Breakdown:\n")
	sb.WriteString(fmt.Sprintf("  Standups:   %.1f team-hours ($%.0f)\n", standupHours, standupHours*float64(hourlyRate)))
	sb.WriteString(fmt.Sprintf("  Planning:   %.1f team-hours ($%.0f)\n", planningHours, planningHours*float64(hourlyRate)))
	sb.WriteString(fmt.Sprintf("  Review:     %.1f team-hours ($%.0f)\n", reviewHours, reviewHours*float64(hourlyRate)))
	sb.WriteString(fmt.Sprintf("  Retro:      %.1f team-hours ($%.0f)\n", retroHours, retroHours*float64(hourlyRate)))
	sb.WriteString(fmt.Sprintf("  Grooming:   %.1f team-hours ($%.0f)\n", groomingHours, groomingHours*float64(hourlyRate)))
	sb.WriteString(fmt.Sprintf("\n  TOTAL: %.1f team-hours (%.1f%% of sprint capacity)\n", totalCeremonyHours, ceremonyPercent))
	sb.WriteString(fmt.Sprintf("  Cost per sprint: $%.0f\n", ceremonyCost))

	if ceremonyPercent > 25 {
		sb.WriteString("\nWARNING: >25% of sprint in ceremonies. Consider:\n")
		sb.WriteString(fmt.Sprintf("  - Async standups (saves %.0f hours)\n", standupHours*0.7))
		sb.WriteString("  - Shorter planning (split into pre-planning + planning)\n")
		sb.WriteString("  - Skip retro if no actions from last one\n")
	} else if ceremonyPercent < 10 {
		sb.WriteString("\nNOTE: <10% in ceremonies. Team might be under-communicating.\n")
	} else {
		sb.WriteString("\nHealthy range (10-25%). Ceremonies are proportionate.\n")
	}

	if throughputNote != "" {
		sb.WriteString("\n" + throughputNote + "\n")
	}

	return textResult(sb.String()), nil
}

func (h *Handlers) SprintCommitmentAdvisor(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 8)
	if len(snaps) < 3 {
		return textResult("Need 3+ sprint snapshots. Use pm_snapshot_sprint at end of each sprint."), nil
	}

	doneValues := make([]int, len(snaps))
	sum := 0
	for i, s := range snaps {
		doneValues[i] = s.Done
		sum += s.Done
	}
	avg := sum / len(snaps)

	sorted := make([]int, len(doneValues))
	copy(sorted, doneValues)
	sort.Ints(sorted)

	pessimistic := sorted[0]
	optimistic := sorted[len(sorted)-1]
	p25 := sorted[len(sorted)/4]

	avgCarryover := 0
	for _, s := range snaps {
		avgCarryover += s.Carryover
	}
	avgCarryover /= len(snaps)

	recommended := avg - avgCarryover
	if recommended < p25 {
		recommended = p25
	}

	teamSize := req.GetInt("team_size", 0)
	leaveDays := req.GetInt("leave_days", 0)
	if teamSize > 0 && leaveDays > 0 {
		sprintDays := 10
		capacity := float64(teamSize*sprintDays-leaveDays) / float64(teamSize*sprintDays)
		recommended = int(float64(recommended) * capacity)
	}

	var sb strings.Builder
	sb.WriteString("Sprint Commitment Advisor\n\n")
	sb.WriteString(fmt.Sprintf("Based on %d sprints of data:\n", len(snaps)))
	sb.WriteString(fmt.Sprintf("  Average completed: %d items/sprint\n", avg))
	sb.WriteString(fmt.Sprintf("  Best sprint: %d items\n", optimistic))
	sb.WriteString(fmt.Sprintf("  Worst sprint: %d items\n", pessimistic))
	sb.WriteString(fmt.Sprintf("  Average carryover: %d items\n\n", avgCarryover))
	sb.WriteString(fmt.Sprintf("RECOMMENDATION: Commit to %d items\n", recommended))
	sb.WriteString(fmt.Sprintf("  Conservative (90%% confidence): %d items\n", p25))
	sb.WriteString(fmt.Sprintf("  Stretch (50%% confidence): %d items\n\n", optimistic))

	if avgCarryover > 2 {
		sb.WriteString(fmt.Sprintf("NOTE: Team carries over %d items on average. This suggests over-commitment.\n", avgCarryover))
		sb.WriteString("Reduce commitment by 2-3 items below average to build trust with stakeholders.\n")
	}

	return textResult(sb.String()), nil
}

func (h *Handlers) DysfunctionDetector(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	var dysfunctions []string

	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 6)
	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)

	if len(snaps) >= 3 {
		allHigh := true
		for _, s := range snaps[:3] {
			if s.CompletionRate < 70 {
				allHigh = false
				break
			}
		}
		if allHigh {
			goals, _ := h.Memory.GetGoalHistory(ctx, boardID, 3)
			if len(goals) == 0 {
				dysfunctions = append(dysfunctions, "ZOMBIE SCRUM: Team completes tasks but has no sprint goals. Work may not connect to value.")
			}
		}
	}

	if len(sprints) > 0 {
		issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
		assigneeCounts := map[string]int{}
		for _, issue := range issues {
			lower := strings.ToLower(issue.Status)
			if strings.Contains(lower, "done") || strings.Contains(lower, "closed") {
				if issue.Assignee != "" {
					assigneeCounts[issue.Assignee]++
				}
			}
		}
		if len(assigneeCounts) > 1 {
			total := 0
			maxPerson := ""
			maxCount := 0
			for p, c := range assigneeCounts {
				total += c
				if c > maxCount {
					maxCount = c
					maxPerson = p
				}
			}
			if total > 0 && float64(maxCount)/float64(total) > 0.5 {
				dysfunctions = append(dysfunctions, fmt.Sprintf("HERO CULTURE: %s completed %.0f%% of done items. Bus factor risk. Distribute work.", maxPerson, float64(maxCount)/float64(total)*100))
			}
		}
	}

	if len(snaps) >= 2 {
		creepCount := 0
		for i := 0; i < len(snaps)-1; i++ {
			if snaps[i].TotalIssues > snaps[i+1].TotalIssues+3 {
				creepCount++
			}
		}
		if creepCount >= 2 {
			dysfunctions = append(dysfunctions, "CHRONIC SCOPE CREEP: Sprint scope expanded in 2+ recent sprints. PO may be adding items mid-sprint.")
		}
	}

	if len(snaps) >= 3 {
		highCarryover := 0
		for _, s := range snaps {
			if s.Carryover > 3 {
				highCarryover++
			}
		}
		if highCarryover >= 3 {
			dysfunctions = append(dysfunctions, "CARRYOVER ADDICTION: 3+ sprints with high carryover. Team consistently overcommits. Need to reduce sprint scope 20-30%.")
		}
	}

	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	longBlockers := 0
	for _, b := range blockers {
		if time.Since(b.BlockedSince).Hours() > 120 {
			longBlockers++
		}
	}
	if longBlockers >= 2 {
		dysfunctions = append(dysfunctions, fmt.Sprintf("BLOCKER PARALYSIS: %d items blocked >5 days. Team may lack escalation culture or external dependencies are unmanaged.", longBlockers))
	}

	if len(dysfunctions) == 0 {
		return textResult("No team dysfunctions detected. Team appears healthy based on available data.\nKeep collecting sprint snapshots for better signal."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Team Dysfunction Detection: %d patterns found\n\n", len(dysfunctions)))
	for i, d := range dysfunctions {
		sb.WriteString(fmt.Sprintf("%d. %s\n\n", i+1, d))
	}
	sb.WriteString("Each pattern above has a specific coaching intervention. Use pm_coaching for advice.\n")

	return textResult(sb.String()), nil
}

func (h *Handlers) StakeholderReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}
	sendToLark := req.GetBool("send_to_lark", false)

	var contextData strings.Builder

	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) > 0 {
		sprint := sprints[0]
		issues, _ := h.Jira.GetSprintIssues(ctx, sprint.ID)
		var done, total int
		total = len(issues)
		for _, i := range issues {
			lower := strings.ToLower(i.Status)
			if strings.Contains(lower, "done") || strings.Contains(lower, "closed") {
				done++
			}
		}
		if total > 0 {
			contextData.WriteString(fmt.Sprintf("Sprint: %s (Goal: %s)\n", sprint.Name, sprint.Goal))
			contextData.WriteString(fmt.Sprintf("Progress: %d/%d items completed (%.0f%%)\n", done, total, float64(done)/float64(total)*100))
		}
	}

	risks, _ := h.Memory.GetOpenRisks(ctx)
	if len(risks) > 0 {
		contextData.WriteString(fmt.Sprintf("\nRisks: %d open\n", len(risks)))
		for _, r := range risks {
			if r.Severity == "critical" || r.Severity == "high" {
				contextData.WriteString(fmt.Sprintf("  [%s] %s\n", r.Severity, r.Title))
			}
		}
	}

	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) > 0 {
		contextData.WriteString(fmt.Sprintf("\nBlockers: %d active\n", len(blockers)))
	}

	systemPrompt := `Write a 100-word executive status update for a VP stakeholder.
Rules:
- No technical jargon, no story points, no sprint velocity
- Lead with: on track / at risk / blocked
- Mention: what shipped, what's coming, what needs attention
- Business language only
- Format: 3-4 short sentences max`

	report, err := h.aiComplete(ctx, systemPrompt, contextData.String())
	if err != nil {
		return textResult("Stakeholder Update:\n\n" + contextData.String()), nil
	}

	if sendToLark {
		_ = h.Lark.SendMarkdown(ctx, "Stakeholder Update", report)
		return textResult(report + "\n\n(Sent to Lark)"), nil
	}
	return textResult(report), nil
}
