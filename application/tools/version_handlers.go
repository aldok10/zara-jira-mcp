package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// GetAttachments lists attachments on an issue.
func (h *Handlers) GetAttachments(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key parameter required"), nil
	}
	attachments, err := h.Jira.GetAttachments(ctx, key)
	if err != nil {
		return sanitizedError("failed to get attachments", err), nil
	}
	if len(attachments) == 0 {
		return textResult(fmt.Sprintf("No attachments on %s.", key)), nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Attachments on %s (%d):\n\n", key, len(attachments)))
	for _, a := range attachments {
		sb.WriteString(fmt.Sprintf("- [%s] %s (%s, %d KB) by %s — %s\n", a.ID, a.Filename, a.MimeType, a.Size/1024, a.Author, a.URL))
	}
	return textResult(sb.String()), nil
}

// GetVersions lists project versions.
func (h *Handlers) GetVersions(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	project, err := req.RequireString("project")
	if err != nil {
		return errorResult("project parameter required"), nil
	}
	versions, err := h.Jira.GetVersions(ctx, project)
	if err != nil {
		return sanitizedError("failed to get versions", err), nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Versions for %s (%d):\n\n", project, len(versions)))
	for _, v := range versions {
		status := "unreleased"
		if v.Released {
			status = "released"
		}
		sb.WriteString(fmt.Sprintf("- [%s] %s (%s) %s\n", v.ID, v.Name, status, v.ReleaseDate))
	}
	return textResult(sb.String()), nil
}

// CreateVersion creates a new project version.
func (h *Handlers) CreateVersion(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	project, err := req.RequireString("project")
	if err != nil {
		return errorResult("project parameter required"), nil
	}
	name, err := req.RequireString("name")
	if err != nil {
		return errorResult("name parameter required"), nil
	}
	desc := req.GetString("description", "")

	version, err := h.Jira.CreateVersion(ctx, project, name, desc)
	if err != nil {
		return sanitizedError("failed to create version", err), nil
	}
	return textResult(fmt.Sprintf("Created version: %s (ID: %s)", version.Name, version.ID)), nil
}

// ReleaseVersion marks a version as released.
func (h *Handlers) ReleaseVersion(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	versionID, err := req.RequireString("version_id")
	if err != nil {
		return errorResult("version_id parameter required"), nil
	}
	if err := h.Jira.ReleaseVersion(ctx, versionID); err != nil {
		return sanitizedError("failed to release version", err), nil
	}
	return textResult(fmt.Sprintf("Version %s marked as released.", versionID)), nil
}

// GetComponents lists project components.
func (h *Handlers) GetComponents(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	project, err := req.RequireString("project")
	if err != nil {
		return errorResult("project parameter required"), nil
	}
	components, err := h.Jira.GetComponents(ctx, project)
	if err != nil {
		return sanitizedError("failed to get components", err), nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Components for %s (%d):\n\n", project, len(components)))
	for _, c := range components {
		lead := c.Lead
		if lead == "" {
			lead = "unassigned"
		}
		sb.WriteString(fmt.Sprintf("- [%s] %s (lead: %s)\n", c.ID, c.Name, lead))
	}
	return textResult(sb.String()), nil
}

// GetFields lists all Jira fields (system + custom).
func (h *Handlers) GetFields(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	onlyCustom := req.GetBool("custom_only", false)

	fields, err := h.Jira.GetFields(ctx)
	if err != nil {
		return sanitizedError("failed to get fields", err), nil
	}

	var sb strings.Builder
	count := 0
	for _, f := range fields {
		if onlyCustom && !f.Custom {
			continue
		}
		count++
		custom := ""
		if f.Custom {
			custom = " [custom]"
		}
		sb.WriteString(fmt.Sprintf("- %s: %s (%s)%s\n", f.ID, f.Name, f.Type, custom))
	}
	return textResult(fmt.Sprintf("Fields (%d):\n\n%s", count, sb.String())), nil
}

// TechDebtRatio calculates bug/debt vs feature ratio.
func (h *Handlers) TechDebtRatio(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	project := req.GetString("project", "")
	sprintCount := req.GetInt("sprints", 3)

	// Count bugs/tech-debt
	debtJQL := "resolution = Unresolved AND (type = Bug OR labels = tech-debt)"
	featureJQL := "resolution = Unresolved AND type in (Story, Task) AND labels != tech-debt"
	if project != "" {
		debtJQL = fmt.Sprintf("project = %s AND %s", project, debtJQL)
		featureJQL = fmt.Sprintf("project = %s AND %s", project, featureJQL)
	}

	debtResult, err := h.Jira.SearchIssues(ctx, debtJQL, 1, 0)
	if err != nil {
		return sanitizedError("failed to query tech debt ratio", err), nil
	}
	featureResult, err := h.Jira.SearchIssues(ctx, featureJQL, 1, 0)
	if err != nil {
		return sanitizedError("failed to query features count", err), nil
	}

	total := debtResult.Total + featureResult.Total
	ratio := 0.0
	if total > 0 {
		ratio = float64(debtResult.Total) / float64(total) * 100
	}

	health := "healthy"
	if ratio > 30 {
		health = "critical"
	} else if ratio > 20 {
		health = "concerning"
	} else if ratio > 10 {
		health = "moderate"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Tech Debt Ratio: %.1f%% (%s)\n\n", ratio, health))
	sb.WriteString(fmt.Sprintf("- Bugs/Debt items: %d\n", debtResult.Total))
	sb.WriteString(fmt.Sprintf("- Feature/Task items: %d\n", featureResult.Total))
	sb.WriteString(fmt.Sprintf("- Total open: %d\n", total))

	if ratio > 20 {
		sb.WriteString("\nRecommendation: Dedicate 20-30% of sprint capacity to debt reduction.")
	}

	_ = sprintCount // TODO: historical trend across sprints
	return textResult(sb.String()), nil
}

// PriorityChurn detects priority instability.
func (h *Handlers) PriorityChurn(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	project := req.GetString("project", "")
	days := req.GetInt("days", 14)

	jql := fmt.Sprintf("resolution = Unresolved AND priority changed DURING (-%dd, now()) ORDER BY updated DESC", days)
	if project != "" {
		jql = fmt.Sprintf("project = %s AND %s", project, jql)
	}

	result, err := h.Jira.SearchIssues(ctx, jql, 50, 0)
	if err != nil {
		return sanitizedError("failed to detect priority churn", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Priority Churn (last %d days): %d issues had priority changes\n\n", days, len(result.Issues)))

	if len(result.Issues) > 10 {
		sb.WriteString("WARNING: High churn indicates unstable priorities. This correlates with team burnout (DORA 2024).\n\n")
	}

	for _, issue := range result.Issues {
		sb.WriteString(fmt.Sprintf("- %s [%s] %s (Priority: %s)\n", issue.Key, issue.Status, issue.Summary, issue.Priority))
	}
	return textResult(sb.String()), nil
}
