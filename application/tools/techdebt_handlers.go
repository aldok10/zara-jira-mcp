package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/mark3labs/mcp-go/mcp"
)

// RecordTechDebt adds a tech debt item to the backlog with severity and impact.
func (h *Handlers) RecordTechDebt(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := req.RequireString("title")
	if err != nil {
		return errorResult("title required"), nil
	}

	// Store as risk with "tech_debt" in sprint_name field for filtering
	r := &memdom.Risk{
		Title:        title,
		Description:  req.GetString("description", ""),
		Severity:     req.GetString("impact", "medium"), // high = blocks velocity, medium = slows, low = cosmetic
		Status:       "open",
		Owner:        req.GetString("owner", ""),
		Mitigation:   req.GetString("fix_approach", ""),
		IdentifiedAt: time.Now(),
		SprintName:   "tech_debt:" + req.GetString("category", "code"), // code, architecture, testing, infra, docs
	}

	if err := h.Memory.SaveRisk(ctx, r); err != nil {
		return sanitizedError("failed to record tech debt", err), nil
	}

	return textResult(fmt.Sprintf("Tech debt recorded: [%s] %s\nCategory: %s\nFix: %s",
		r.Severity, title, strings.TrimPrefix(r.SprintName, "tech_debt:"), r.Mitigation)), nil
}

// TechDebtDashboard shows all tech debt items prioritized by impact.
func (h *Handlers) TechDebtDashboard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	allRisks, err := h.Memory.GetAllRisks(ctx, 100)
	if err != nil {
		return sanitizedError("failed to get risks", err), nil
	}

	var debts []memdom.Risk
	for _, r := range allRisks {
		if strings.HasPrefix(r.SprintName, "tech_debt:") {
			debts = append(debts, r)
		}
	}

	if len(debts) == 0 {
		return textResult("No tech debt tracked. Use pm_tech_debt_add to record items."), nil
	}

	// Group by status
	var open, resolved []memdom.Risk
	for _, d := range debts {
		if d.Status == "open" || d.Status == "mitigating" {
			open = append(open, d)
		} else {
			resolved = append(resolved, d)
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Tech Debt Dashboard (%d open, %d resolved)\n\n", len(open), len(resolved)))

	if len(open) > 0 {
		sb.WriteString("OPEN (by impact):\n")
		for _, d := range open {
			cat := strings.TrimPrefix(d.SprintName, "tech_debt:")
			days := int(time.Since(d.IdentifiedAt).Hours() / 24)
			sb.WriteString(fmt.Sprintf("  #%d [%s] [%s] %s (%d days old)\n", d.ID, strings.ToUpper(d.Severity), cat, d.Title, days))
			if d.Mitigation != "" {
				sb.WriteString(fmt.Sprintf("      Fix: %s\n", d.Mitigation))
			}
		}
	}

	// Debt velocity (how much resolved per month)
	thisMonth := 0
	monthAgo := time.Now().AddDate(0, -1, 0)
	for _, d := range resolved {
		if d.ResolvedAt != nil && d.ResolvedAt.After(monthAgo) {
			thisMonth++
		}
	}
	sb.WriteString(fmt.Sprintf("\nDebt paid this month: %d items\n", thisMonth))
	sb.WriteString(fmt.Sprintf("Recommended sprint allocation: 15-20%% capacity\n"))

	return textResult(sb.String()), nil
}

// TechDebtBudget recommends how much sprint capacity to allocate to tech debt.
func (h *Handlers) TechDebtBudget(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	// Count open debt
	allRisks, _ := h.Memory.GetAllRisks(ctx, 100)
	var openDebt int
	var highImpact int
	for _, r := range allRisks {
		if strings.HasPrefix(r.SprintName, "tech_debt:") && (r.Status == "open" || r.Status == "mitigating") {
			openDebt++
			if r.Severity == "high" || r.Severity == "critical" {
				highImpact++
			}
		}
	}

	// Get velocity for capacity calculation
	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
	avgVelocity := 0
	if len(snaps) > 0 {
		total := 0
		for _, s := range snaps {
			total += s.Velocity
		}
		avgVelocity = total / len(snaps)
	}

	// Calculate recommendation
	baseAllocation := 15 // 15% default
	if highImpact > 3 {
		baseAllocation = 25 // more if high-impact debt is piling up
	} else if openDebt > 10 {
		baseAllocation = 20
	}

	pointsBudget := 0
	if avgVelocity > 0 {
		pointsBudget = avgVelocity * baseAllocation / 100
	}

	var sb strings.Builder
	sb.WriteString("Tech Debt Budget Recommendation:\n\n")
	sb.WriteString(fmt.Sprintf("Open debt items: %d (high impact: %d)\n", openDebt, highImpact))
	sb.WriteString(fmt.Sprintf("Avg velocity: %d points/sprint\n\n", avgVelocity))
	sb.WriteString(fmt.Sprintf("Recommended allocation: %d%% = %d points/sprint\n", baseAllocation, pointsBudget))
	sb.WriteString(fmt.Sprintf("Priority: Fix high-impact items first (%d items)\n\n", highImpact))

	if highImpact == 0 && openDebt < 5 {
		sb.WriteString("Status: HEALTHY — debt is under control\n")
	} else if highImpact > 3 {
		sb.WriteString("Status: URGENT — high-impact debt accumulating, increase allocation\n")
	} else {
		sb.WriteString("Status: MANAGEABLE — maintain consistent allocation\n")
	}

	return textResult(sb.String()), nil
}

// SprintReviewPrep generates sprint review preparation with demo order and talking points.
func (h *Handlers) SprintReviewPrep(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	// Categorize done items for demo
	var features, bugs, tasks []string
	var notDone []string
	for _, issue := range issues {
		lower := strings.ToLower(issue.Status)
		if strings.Contains(lower, "done") || strings.Contains(lower, "closed") || strings.Contains(lower, "resolved") {
			entry := fmt.Sprintf("%s: %s (%s)", issue.Key, issue.Summary, issue.Assignee)
			switch strings.ToLower(issue.Type) {
			case "story", "feature":
				features = append(features, entry)
			case "bug":
				bugs = append(bugs, entry)
			default:
				tasks = append(tasks, entry)
			}
		} else {
			notDone = append(notDone, fmt.Sprintf("%s: %s [%s]", issue.Key, issue.Summary, issue.Status))
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== SPRINT REVIEW PREP: %s ===\n\n", sprint.Name))
	sb.WriteString(fmt.Sprintf("Goal: %s\n\n", sprint.Goal))

	sb.WriteString("DEMO ORDER (features first, then fixes):\n")
	for i, f := range features {
		sb.WriteString(fmt.Sprintf("  %d. %s\n", i+1, f))
	}
	for i, b := range bugs {
		sb.WriteString(fmt.Sprintf("  %d. [FIX] %s\n", len(features)+i+1, b))
	}
	sb.WriteString("\n")

	if len(tasks) > 0 {
		sb.WriteString(fmt.Sprintf("MENTION (don't demo): %d tasks/improvements\n", len(tasks)))
		for _, t := range tasks {
			sb.WriteString(fmt.Sprintf("  - %s\n", t))
		}
		sb.WriteString("\n")
	}

	if len(notDone) > 0 {
		sb.WriteString(fmt.Sprintf("NOT COMPLETED (%d — explain briefly):\n", len(notDone)))
		for _, n := range notDone {
			sb.WriteString(fmt.Sprintf("  - %s\n", n))
		}
		sb.WriteString("\n")
	}

	done := len(features) + len(bugs) + len(tasks)
	total := done + len(notDone)
	sb.WriteString(fmt.Sprintf("SUMMARY: %d/%d done (%.0f%%)\n\n", done, total, float64(done)/float64(total)*100))

	sb.WriteString("TALKING POINTS:\n")
	sb.WriteString("  1. Did we meet the sprint goal? (Yes/Partially/No + why)\n")
	sb.WriteString("  2. Key demo highlights (features that deliver user value)\n")
	sb.WriteString("  3. What didn't make it and why (brief, no excuses)\n")
	sb.WriteString("  4. Risks/concerns for next sprint\n")
	sb.WriteString("  5. Questions for stakeholders\n")

	return textResult(sb.String()), nil
}

// MCPStats shows server self-monitoring stats.
func (h *Handlers) MCPStats(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)

	var sb strings.Builder
	sb.WriteString("=== MCP Server Stats ===\n\n")

	// Memory stats
	allRisks, _ := h.Memory.GetAllRisks(ctx, 1000)
	decisions, _ := h.Memory.GetDecisions(ctx, 1000)
	blockerHistory, _ := h.Memory.GetBlockerHistory(ctx, 1000)
	retros, _ := h.Memory.GetRetrospectives(ctx, 100)
	meetings, _ := h.Memory.GetMeetingNotes(ctx, "", 1000)

	sb.WriteString("MEMORY CONTENTS:\n")
	sb.WriteString(fmt.Sprintf("  Risks: %d\n", len(allRisks)))
	sb.WriteString(fmt.Sprintf("  Decisions: %d\n", len(decisions)))
	sb.WriteString(fmt.Sprintf("  Blockers: %d\n", len(blockerHistory)))
	sb.WriteString(fmt.Sprintf("  Retrospectives: %d\n", len(retros)))
	sb.WriteString(fmt.Sprintf("  Meetings: %d\n", len(meetings)))

	if boardID > 0 {
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 100)
		scores, _ := h.Memory.GetHealthScores(ctx, boardID, 100)
		goals, _ := h.Memory.GetGoalHistory(ctx, boardID, 100)
		sb.WriteString(fmt.Sprintf("  Sprint Snapshots: %d\n", len(snaps)))
		sb.WriteString(fmt.Sprintf("  Health Scores: %d\n", len(scores)))
		sb.WriteString(fmt.Sprintf("  Sprint Goals: %d\n", len(goals)))
	}

	// Data freshness
	sb.WriteString("\nDATA FRESHNESS:\n")
	if len(allRisks) > 0 {
		sb.WriteString(fmt.Sprintf("  Latest risk: %s\n", allRisks[0].IdentifiedAt.Format("2006-01-02")))
	}
	if len(decisions) > 0 {
		sb.WriteString(fmt.Sprintf("  Latest decision: %s\n", decisions[0].MadeAt.Format("2006-01-02")))
	}
	if len(meetings) > 0 {
		sb.WriteString(fmt.Sprintf("  Latest meeting: %s\n", meetings[0].Date.Format("2006-01-02")))
	}

	sb.WriteString("\nSERVER:\n")
	sb.WriteString("  Version: 0.3.0\n")
	sb.WriteString("  Tools: 124+\n")
	sb.WriteString("  Storage: SQLite (WAL mode)\n")

	return textResult(sb.String()), nil
}
