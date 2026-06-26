package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handlers) PagerDutyIncidents(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.PagerDuty == nil || !h.PagerDuty.Available() {
		return errorResult("PagerDuty not configured. Set PAGERDUTY_API_KEY."), nil
	}
	status := req.GetString("status", "")

	incidents, err := h.PagerDuty.ListIncidents(ctx, status)
	if err != nil {
		return errorResult("PagerDuty API error: " + err.Error()), nil
	}
	if len(incidents) == 0 {
		return textResult("No incidents found."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Incidents (%d):\n\n", len(incidents)))
	for _, i := range incidents {
		sb.WriteString(fmt.Sprintf("- [%s] %s | %s | %s | Assigned: %s | Service: %s\n",
			i.Urgency, i.Title, i.Status, i.CreatedAt[:10], i.Assignee, i.ServiceName))
	}
	return textResult(sb.String()), nil
}

func (h *Handlers) PagerDutyIncidentSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.PagerDuty == nil || !h.PagerDuty.Available() {
		return errorResult("PagerDuty not configured. Set PAGERDUTY_API_KEY."), nil
	}

	incidents, err := h.PagerDuty.ListIncidents(ctx, "")
	if err != nil {
		return errorResult("PagerDuty API error: " + err.Error()), nil
	}

	triggered := 0
	acknowledged := 0
	resolved := 0
	highUrgency := 0
	for _, i := range incidents {
		switch i.Status {
		case "triggered":
			triggered++
		case "acknowledged":
			acknowledged++
		case "resolved":
			resolved++
		}
		if i.Urgency == "high" {
			highUrgency++
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Incident Summary (last %d incidents):\n\n", len(incidents)))
	sb.WriteString(fmt.Sprintf("- Triggered: %d\n", triggered))
	sb.WriteString(fmt.Sprintf("- Acknowledged: %d\n", acknowledged))
	sb.WriteString(fmt.Sprintf("- Resolved: %d\n", resolved))
	sb.WriteString(fmt.Sprintf("- High urgency: %d\n", highUrgency))
	sb.WriteString(fmt.Sprintf("\nImpact: %d active incidents requiring attention.\n", triggered+acknowledged))

	// Calculate average resolution time for resolved incidents
	if resolved > 0 {
		var totalDuration time.Duration
		count := 0
		for _, i := range incidents {
			if i.Status == "resolved" && i.CreatedAt != "" {
				created, err := time.Parse(time.RFC3339, i.CreatedAt)
				if err == nil {
					totalDuration += time.Since(created)
					count++
				}
			}
		}
		if count > 0 {
			avg := totalDuration / time.Duration(count)
			sb.WriteString(fmt.Sprintf("Average age of resolved incidents: %s\n", avg.Round(time.Hour)))
		}
	}
	return textResult(sb.String()), nil
}

func (h *Handlers) PagerDutyOnCall(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.PagerDuty == nil || !h.PagerDuty.Available() {
		return errorResult("PagerDuty not configured. Set PAGERDUTY_API_KEY."), nil
	}

	oncalls, err := h.PagerDuty.GetOnCalls(ctx)
	if err != nil {
		return errorResult("PagerDuty API error: " + err.Error()), nil
	}
	if len(oncalls) == 0 {
		return textResult("No on-call schedules found."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("On-call now (%d):\n\n", len(oncalls)))
	for _, o := range oncalls {
		sb.WriteString(fmt.Sprintf("- %s | Schedule: %s | Until: %s\n",
			o.UserName, o.Schedule, o.End))
	}
	return textResult(sb.String()), nil
}
