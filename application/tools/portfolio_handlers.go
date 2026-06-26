package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// PortfolioOverview shows cross-project health: open issues per project.
func (h *Handlers) PortfolioOverview(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projects, err := h.Jira.GetProjects(ctx)
	if err != nil {
		return sanitizedError("failed to list portfolio projects", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Portfolio Overview (%d projects)\n\n", len(projects)))

	for _, p := range projects {
		jql := fmt.Sprintf("project = %s AND resolution = Unresolved", p.Key)
		result, _ := h.Jira.SearchIssues(ctx, jql, 1, 0)
		total := 0
		if result != nil {
			total = result.Total
		}
		sb.WriteString(fmt.Sprintf("- %s (%s): %d open issues | Lead: %s\n", p.Key, p.Name, total, p.Lead))
	}

	return textResult(sb.String()), nil
}

// PortfolioBlockers shows all active blockers and open dependencies across the portfolio.
func (h *Handlers) PortfolioBlockers(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	deps, _ := h.Memory.GetOpenDependencies(ctx)

	var sb strings.Builder
	if len(blockers) > 0 {
		sb.WriteString(fmt.Sprintf("Active Blockers (%d):\n", len(blockers)))
		for _, b := range blockers {
			days := int(time.Since(b.BlockedSince).Hours() / 24)
			sb.WriteString(fmt.Sprintf("  [%d days] %s (issue: %s, owner: %s)\n", days, b.Description, b.IssueKey, b.Owner))
		}
		sb.WriteString("\n")
	}

	if len(deps) > 0 {
		sb.WriteString(fmt.Sprintf("Open Dependencies (%d):\n", len(deps)))
		for _, d := range deps {
			days := int(time.Since(d.CreatedAt).Hours() / 24)
			sb.WriteString(fmt.Sprintf("  %s -> %s [%s] (%d days)\n", d.FromIssueKey, d.ToIssueKey, d.DependencyType, days))
		}
		sb.WriteString("\n")
	}

	if len(blockers) == 0 && len(deps) == 0 {
		return textResult("No cross-project blockers or dependencies. Clean!"), nil
	}

	return textResult(sb.String()), nil
}

// PortfolioWorkload shows per-person issue count across all projects.
func (h *Handlers) PortfolioWorkload(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	jql := "resolution = Unresolved AND assignee IS NOT EMPTY ORDER BY assignee ASC"
	result, err := h.Jira.SearchIssues(ctx, jql, 200, 0)
	if err != nil {
		return sanitizedError("portfolio operation failed", err), nil
	}

	type personLoad struct {
		total    int
		projects map[string]int
	}
	workload := map[string]*personLoad{}
	for _, issue := range result.Issues {
		if issue.Assignee == "" {
			continue
		}
		if workload[issue.Assignee] == nil {
			workload[issue.Assignee] = &personLoad{projects: map[string]int{}}
		}
		workload[issue.Assignee].total++
		parts := strings.SplitN(issue.Key, "-", 2)
		if len(parts) > 0 {
			workload[issue.Assignee].projects[parts[0]]++
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Portfolio Workload (%d people, %d issues)\n\n", len(workload), len(result.Issues)))
	for person, load := range workload {
		projects := []string{}
		for proj, count := range load.projects {
			projects = append(projects, fmt.Sprintf("%s:%d", proj, count))
		}
		sb.WriteString(fmt.Sprintf("- %s: %d total [%s]\n", person, load.total, strings.Join(projects, ", ")))
	}

	return textResult(sb.String()), nil
}

// PortfolioRisks shows aggregate risks across all projects sorted by severity.
func (h *Handlers) PortfolioRisks(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	risks, err := h.Memory.GetOpenRisks(ctx)
	if err != nil {
		return sanitizedError("portfolio operation failed", err), nil
	}
	if len(risks) == 0 {
		return textResult("No open risks across portfolio."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Portfolio Risk Radar (%d open)\n\n", len(risks)))
	for _, r := range risks {
		days := int(time.Since(r.IdentifiedAt).Hours() / 24)
		sb.WriteString(fmt.Sprintf("[%s] %s | owner: %s | %d days | sprint: %s\n", strings.ToUpper(r.Severity), r.Title, r.Owner, days, r.SprintName))
	}
	return textResult(sb.String()), nil
}

// PortfolioSummary generates an AI-powered executive portfolio summary.
func (h *Handlers) PortfolioSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var ctxData strings.Builder

	projects, _ := h.Jira.GetProjects(ctx)
	ctxData.WriteString(fmt.Sprintf("Projects: %d\n", len(projects)))

	risks, _ := h.Memory.GetOpenRisks(ctx)
	ctxData.WriteString(fmt.Sprintf("Open risks: %d\n", len(risks)))
	for _, r := range risks {
		ctxData.WriteString(fmt.Sprintf("  [%s] %s\n", r.Severity, r.Title))
	}

	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	ctxData.WriteString(fmt.Sprintf("Active blockers: %d\n", len(blockers)))

	deps, _ := h.Memory.GetOpenDependencies(ctx)
	ctxData.WriteString(fmt.Sprintf("Open dependencies: %d\n", len(deps)))

	actions, _ := h.Memory.GetPendingActionItems(ctx)
	ctxData.WriteString(fmt.Sprintf("Pending action items: %d\n", len(actions)))

	systemPrompt := `You are writing an executive portfolio status summary for a PM. Be concise (under 200 words).
Include: overall health assessment, top risks, blockers needing escalation, key numbers.
Format: short paragraph + bullet points for action items.`

	summary, err := h.aiComplete(ctx, systemPrompt, ctxData.String())
	if err != nil {
		return textResult("Portfolio data:\n" + ctxData.String()), nil
	}
	return textResult(summary), nil
}
