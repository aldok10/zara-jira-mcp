package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// ManagementBrief generates a concise management-level status brief.
// Written for someone with 30 seconds: outcomes, risks, asks.
func (h *Handlers) ManagementBrief(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}
	audience := req.GetString("audience", "manager") // manager, director, po

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil || len(sprints) == 0 {
		return errorResult("No active sprint"), nil
	}

	sprint := sprints[0]
	issues, _ := h.Jira.GetSprintIssues(ctx, sprint.ID)

	done, inProgress, blocked, total := 0, 0, 0, len(issues)
	for _, issue := range issues {
		switch {
		case issue.Status == "Done" || issue.Status == "Closed":
			done++
		case issue.Status == "Blocked":
			blocked++
		case issue.Status != "To Do" && issue.Status != "Backlog":
			inProgress++
		}
	}

	// Blockers from memory
	var blockerLines []string
	if h.Memory != nil {
		blockers, _ := h.Memory.GetActiveBlockers(ctx)
		for _, b := range blockers {
			days := int(time.Since(b.BlockedSince).Hours() / 24)
			blockerLines = append(blockerLines, fmt.Sprintf("- %s: %s (%dd)", b.IssueKey, b.Description, days))
		}
	}

	// Risks from memory
	var riskLines []string
	if h.Memory != nil {
		risks, _ := h.Memory.GetOpenRisks(ctx)
		for _, r := range risks {
			if r.Severity == "critical" || r.Severity == "high" {
				riskLines = append(riskLines, fmt.Sprintf("- [%s] %s (owner: %s)", strings.ToUpper(r.Severity), r.Title, r.Owner))
			}
		}
	}

	completion := 0.0
	if total > 0 {
		completion = float64(done) / float64(total) * 100
	}

	status := "On Track"
	if completion < 30 && blocked > 0 {
		status = "At Risk"
	} else if blocked > 2 || completion < 20 {
		status = "Behind"
	}

	var sb strings.Builder

	switch audience {
	case "director", "vp", "executive":
		// Ultra-brief: 3 lines max
		sb.WriteString(fmt.Sprintf("**Sprint: %s — %s**\n", sprint.Name, status))
		sb.WriteString(fmt.Sprintf("Progress: %d/%d (%.0f%%) | Blocked: %d\n", done, total, completion, blocked))
		if len(riskLines) > 0 {
			sb.WriteString(fmt.Sprintf("Top risk: %s\n", riskLines[0]))
		}

	case "po", "product_owner":
		// PO needs: what's shipping, what's not, why
		sb.WriteString(fmt.Sprintf("## Sprint Status for PO: %s\n\n", sprint.Name))
		sb.WriteString(fmt.Sprintf("**Status: %s** | %.0f%% complete\n\n", status, completion))
		sb.WriteString(fmt.Sprintf("**Will ship (%d items):**\n", done))
		for _, issue := range issues {
			if issue.Status == "Done" || issue.Status == "Closed" {
				sb.WriteString(fmt.Sprintf("- %s: %s\n", issue.Key, issue.Summary))
			}
		}
		if blocked > 0 {
			sb.WriteString(fmt.Sprintf("\n**Blocked (%d items) — needs your input:**\n", blocked))
			for _, issue := range issues {
				if issue.Status == "Blocked" {
					sb.WriteString(fmt.Sprintf("- %s: %s\n", issue.Key, issue.Summary))
				}
			}
		}
		if len(riskLines) > 0 {
			sb.WriteString("\n**Risks:**\n" + strings.Join(riskLines, "\n") + "\n")
		}

	default: // manager
		sb.WriteString(fmt.Sprintf("## Sprint Brief: %s\n\n", sprint.Name))
		sb.WriteString(fmt.Sprintf("**Status: %s**\n", status))
		sb.WriteString(fmt.Sprintf("- Goal: %s\n", sprint.Goal))
		sb.WriteString(fmt.Sprintf("- Progress: %d/%d done (%.0f%%)\n", done, total, completion))
		sb.WriteString(fmt.Sprintf("- In Progress: %d | Blocked: %d\n\n", inProgress, blocked))

		if len(blockerLines) > 0 {
			sb.WriteString("**Active Blockers:**\n")
			sb.WriteString(strings.Join(blockerLines, "\n") + "\n\n")
		}
		if len(riskLines) > 0 {
			sb.WriteString("**Risks (High/Critical):**\n")
			sb.WriteString(strings.Join(riskLines, "\n") + "\n\n")
		}
	}

	return textResult(sb.String()), nil
}

// DependencyReport shows cross-team/external dependencies and their status.
func (h *Handlers) DependencyReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)

	var sb strings.Builder
	sb.WriteString("## Cross-Team Dependency Report\n\n")

	// From memory
	if h.Memory != nil {
		deps, _ := h.Memory.GetOpenDependencies(ctx)
		if len(deps) > 0 {
			overdue := 0
			for _, d := range deps {
				age := int(time.Since(d.CreatedAt).Hours() / 24)
				status := "OPEN"
				if d.ResolvedAt != nil {
					status = "RESOLVED"
				} else if age > 5 {
					status = "OVERDUE"
					overdue++
				}
				sb.WriteString(fmt.Sprintf("- [%s] %s → %s: %s (%dd)\n", status, d.FromIssueKey, d.ToIssueKey, d.Description, age))
			}
			sb.WriteString(fmt.Sprintf("\nTotal: %d | Overdue (>5d): %d\n", len(deps), overdue))
		} else {
			sb.WriteString("No dependencies recorded. Use pm_record_dependency to track.\n")
		}
	}

	// Also check Jira for blocked items pointing to other projects
	if boardID > 0 {
		sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
		if len(sprints) > 0 {
			issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
			var externalBlocks []string
			for _, issue := range issues {
				if issue.Status == "Blocked" {
					externalBlocks = append(externalBlocks, fmt.Sprintf("- %s: %s (assignee: %s)", issue.Key, issue.Summary, issue.Assignee))
				}
			}
			if len(externalBlocks) > 0 {
				sb.WriteString(fmt.Sprintf("\n**Blocked Sprint Items (%d):**\n", len(externalBlocks)))
				sb.WriteString(strings.Join(externalBlocks, "\n") + "\n")
			}
		}
	}

	return textResult(sb.String()), nil
}

// EscalationReport generates a report of items that need management attention.
func (h *Handlers) EscalationReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)

	var sb strings.Builder
	sb.WriteString("## Escalation Report\n\n")
	sb.WriteString("Items requiring management attention:\n\n")

	escalations := 0

	// 1. Long-standing blockers (>3 days)
	if h.Memory != nil {
		blockers, _ := h.Memory.GetActiveBlockers(ctx)
		var aged []string
		for _, b := range blockers {
			days := int(time.Since(b.BlockedSince).Hours() / 24)
			if days >= 3 {
				aged = append(aged, fmt.Sprintf("- %s: %s (%dd blocked, owner: %s)", b.IssueKey, b.Description, days, b.Owner))
			}
		}
		if len(aged) > 0 {
			escalations += len(aged)
			sb.WriteString(fmt.Sprintf("**Blockers (>3 days): %d**\n", len(aged)))
			sb.WriteString(strings.Join(aged, "\n") + "\n\n")
		}
	}

	// 2. Critical/High risks unmitigated
	if h.Memory != nil {
		risks, _ := h.Memory.GetOpenRisks(ctx)
		var critRisks []string
		for _, r := range risks {
			if r.Severity == "critical" || r.Severity == "high" {
				days := int(time.Since(r.IdentifiedAt).Hours() / 24)
				critRisks = append(critRisks, fmt.Sprintf("- [%s] %s (%dd, owner: %s)", strings.ToUpper(r.Severity), r.Title, days, r.Owner))
			}
		}
		if len(critRisks) > 0 {
			escalations += len(critRisks)
			sb.WriteString(fmt.Sprintf("**Unresolved Critical/High Risks: %d**\n", len(critRisks)))
			sb.WriteString(strings.Join(critRisks, "\n") + "\n\n")
		}
	}

	// 3. Sprint health at risk
	if boardID > 0 {
		sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
		if len(sprints) > 0 {
			issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
			done, total, blocked := 0, len(issues), 0
			for _, i := range issues {
				if i.Status == "Done" || i.Status == "Closed" {
					done++
				}
				if i.Status == "Blocked" {
					blocked++
				}
			}
			if total > 0 {
				rate := float64(done) / float64(total) * 100
				if rate < 30 || blocked > 2 {
					escalations++
					sb.WriteString(fmt.Sprintf("**Sprint At Risk:** %.0f%% done, %d blocked\n\n", rate, blocked))
				}
			}
		}
	}

	// 4. Overdue items
	stale, err := h.Jira.SearchIssues(ctx, "resolution = Unresolved AND updated <= -7d AND priority in (High, Highest) ORDER BY priority DESC", 10, 0)
	if err == nil && len(stale.Issues) > 0 {
		escalations += len(stale.Issues)
		sb.WriteString(fmt.Sprintf("**High-Priority Stale Items (>7d no update): %d**\n", len(stale.Issues)))
		for _, i := range stale.Issues {
			sb.WriteString(fmt.Sprintf("- %s [%s]: %s (assignee: %s)\n", i.Key, i.Priority, i.Summary, i.Assignee))
		}
		sb.WriteString("\n")
	}

	if escalations == 0 {
		sb.WriteString("No items currently require escalation. All clear.")
	} else {
		sb.WriteString(fmt.Sprintf("---\n**Total escalation items: %d**\n", escalations))
		sb.WriteString("Action needed: Review and resolve or reassign above items.")
	}

	return textResult(sb.String()), nil
}

// ResourceUtilization shows team workload distribution for management.
func (h *Handlers) ResourceUtilization(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
	if err != nil || len(sprints) == 0 {
		return errorResult("No active sprint"), nil
	}

	issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)

	type memberStats struct {
		Total   int
		Done    int
		WIP     int
		Blocked int
	}
	stats := map[string]*memberStats{}

	for _, issue := range issues {
		name := issue.Assignee
		if name == "" {
			name = "(Unassigned)"
		}
		if stats[name] == nil {
			stats[name] = &memberStats{}
		}
		s := stats[name]
		s.Total++
		switch {
		case issue.Status == "Done" || issue.Status == "Closed":
			s.Done++
		case issue.Status == "Blocked":
			s.Blocked++
			s.WIP++
		case issue.Status != "To Do" && issue.Status != "Backlog":
			s.WIP++
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## Resource Utilization: %s\n\n", sprints[0].Name))
	sb.WriteString("| Member | Assigned | Done | WIP | Blocked | Load |\n")
	sb.WriteString("|--------|----------|------|-----|---------|------|\n")

	for name, s := range stats {
		load := "Normal"
		if s.WIP > 4 {
			load = "OVERLOADED"
		} else if s.WIP > 2 {
			load = "High"
		} else if s.Total == s.Done {
			load = "Available"
		}
		sb.WriteString(fmt.Sprintf("| %s | %d | %d | %d | %d | %s |\n", name, s.Total, s.Done, s.WIP, s.Blocked, load))
	}

	// Flag issues
	overloaded := 0
	available := 0
	for _, s := range stats {
		if s.WIP > 4 {
			overloaded++
		}
		if s.Total == s.Done && s.Total > 0 {
			available++
		}
	}
	if overloaded > 0 || available > 0 {
		sb.WriteString(fmt.Sprintf("\n**Flags:** %d overloaded, %d available for more work\n", overloaded, available))
	}

	return textResult(sb.String()), nil
}

// BlockerAgingReport shows blockers with aging info for management review.
func (h *Handlers) BlockerAgingReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("Memory not configured"), nil
	}

	blockers, err := h.Memory.GetActiveBlockers(ctx)
	if err != nil {
		return errorResult("Failed to get blockers: " + err.Error()), nil
	}

	if len(blockers) == 0 {
		return textResult("No active blockers. All clear."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("## Blocker Aging Report (%d active)\n\n", len(blockers)))
	sb.WriteString("| Ticket | Description | Owner | Days Blocked | SLA |\n")
	sb.WriteString("|--------|-------------|-------|-------------|-----|\n")

	critical := 0
	for _, b := range blockers {
		days := int(time.Since(b.BlockedSince).Hours() / 24)
		sla := "OK"
		if days >= 5 {
			sla = "BREACH"
			critical++
		} else if days >= 3 {
			sla = "WARNING"
		}
		desc := b.Description
		if len(desc) > 40 {
			desc = desc[:40] + "..."
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %s | %d | %s |\n", b.IssueKey, desc, b.Owner, days, sla))
	}

	if critical > 0 {
		sb.WriteString(fmt.Sprintf("\n**%d blocker(s) breaching SLA (>5 days).** Escalation recommended.\n", critical))
	}

	return textResult(sb.String()), nil
}

// SprintCommitmentReport shows commitment vs delivery for management accountability.
func (h *Handlers) SprintCommitmentReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	if h.Memory == nil {
		return errorResult("Memory required for historical comparison"), nil
	}

	// Current sprint
	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) == 0 {
		return errorResult("No active sprint"), nil
	}
	issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
	done, total := 0, len(issues)
	for _, i := range issues {
		if i.Status == "Done" || i.Status == "Closed" {
			done++
		}
	}

	// Historical from snapshots
	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)

	var sb strings.Builder
	sb.WriteString("## Sprint Commitment vs Delivery\n\n")
	sb.WriteString(fmt.Sprintf("**Current: %s** — %d/%d delivered (%.0f%%)\n\n", sprints[0].Name, done, total, float64(done)/float64(total)*100))

	if len(snaps) > 0 {
		sb.WriteString("**Historical (last 5 sprints):**\n")
		sb.WriteString("| Sprint | Committed | Delivered | Rate |\n")
		sb.WriteString("|--------|-----------|-----------|------|\n")
		for _, s := range snaps {
			rate := 0.0
			if s.TotalIssues > 0 {
				rate = float64(s.Done) / float64(s.TotalIssues) * 100
			}
			sb.WriteString(fmt.Sprintf("| %s | %d | %d | %.0f%% |\n", s.SprintName, s.TotalIssues, s.Done, rate))
		}

		// Calculate average
		var totalRate float64
		for _, s := range snaps {
			if s.TotalIssues > 0 {
				totalRate += float64(s.Done) / float64(s.TotalIssues) * 100
			}
		}
		avg := totalRate / float64(len(snaps))
		sb.WriteString(fmt.Sprintf("\n**Average delivery rate: %.0f%%**\n", avg))

		if avg < 70 {
			sb.WriteString("\nRecommendation: Team consistently over-committing. Consider reducing sprint scope by 20-30%.\n")
		} else if avg > 95 {
			sb.WriteString("\nNote: Very high delivery rate may indicate under-commitment. Consider stretching capacity.\n")
		}
	}

	return textResult(sb.String()), nil
}
