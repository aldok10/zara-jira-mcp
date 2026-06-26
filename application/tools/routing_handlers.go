package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// NotifyRouted sends a notification to the optimal channel based on severity, audience, and time.
func (h *Handlers) NotifyRouted(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	content, err := req.RequireString("content")
	if err != nil {
		return errorResult("content parameter required"), nil
	}
	severity := req.GetString("severity", "medium") // critical, high, medium, low, info
	audience := req.GetString("audience", "team")   // individual, team, stakeholder, executive
	title := req.GetString("title", "PM Notification")

	channels := routeNotification(severity, audience)

	var results []string
	for _, ch := range channels {
		var sendErr error
		switch ch {
		case "slack_dm", "slack_channel":
			if h.Slack != nil && h.Slack.Available() {
				sendErr = h.Slack.SendRichMessage(ctx, "", title, content)
				if sendErr == nil {
					results = append(results, "Slack: sent")
				}
			}
		case "telegram":
			if h.Telegram != nil && h.Telegram.Available() {
				sendErr = h.Telegram.SendMessage(ctx, 0, fmt.Sprintf("*%s*\n\n%s", title, content))
				if sendErr == nil {
					results = append(results, "Telegram: sent")
				}
			}
		case "lark":
			if h.Lark != nil {
				sendErr = h.Lark.SendMarkdown(ctx, title, content)
				if sendErr == nil {
					results = append(results, "Lark: sent")
				}
			}
		case "teams":
			if h.Teams != nil && h.Teams.Available() {
				sendErr = h.Teams.SendCard(ctx, title, content)
				if sendErr == nil {
					results = append(results, "Teams: sent")
				}
			}
		case "discord":
			if h.Discord != nil && h.Discord.Available() {
				sendErr = h.Discord.SendEmbed(ctx, "", title, content, severityColor(severity))
				if sendErr == nil {
					results = append(results, "Discord: sent")
				}
			}
		case "email":
			if h.Email != nil && h.Email.Available() {
				// Email requires explicit recipient — skip in auto-routing unless stakeholder/executive
				results = append(results, "Email: skipped (no recipient specified)")
			}
		}
		if sendErr != nil {
			results = append(results, fmt.Sprintf("%s: %s", ch, sendErr.Error()))
		}
	}

	if len(results) == 0 {
		return textResult("No channels configured for this route. Configure at least one notification platform."), nil
	}
	return textResult(fmt.Sprintf("Routed [%s/%s]:\n%s", severity, audience, strings.Join(results, "\n"))), nil
}

// DailyDigest generates and sends a daily digest of overnight changes.
func (h *Handlers) DailyDigest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)

	var sections []string

	// Blockers
	if h.Memory != nil {
		blockers, _ := h.Memory.GetActiveBlockers(ctx)
		if len(blockers) > 0 {
			var bb strings.Builder
			bb.WriteString(fmt.Sprintf("*Active Blockers (%d):*\n", len(blockers)))
			for _, b := range blockers {
				days := int(time.Since(b.BlockedSince).Hours() / 24)
				bb.WriteString(fmt.Sprintf("- %s: %s (%d days)\n", b.IssueKey, b.Description, days))
			}
			sections = append(sections, bb.String())
		}

		// Pending action items
		actions, _ := h.Memory.GetPendingActionItems(ctx)
		if len(actions) > 0 {
			var ab strings.Builder
			ab.WriteString(fmt.Sprintf("*Pending Actions (%d):*\n", len(actions)))
			for _, a := range actions {
				ab.WriteString(fmt.Sprintf("- %s (owner: %s)\n", a.Description, a.Owner))
			}
			sections = append(sections, ab.String())
		}

		// Open risks
		risks, _ := h.Memory.GetOpenRisks(ctx)
		if len(risks) > 0 {
			critical := 0
			for _, r := range risks {
				if r.Severity == "critical" || r.Severity == "high" {
					critical++
				}
			}
			if critical > 0 {
				sections = append(sections, fmt.Sprintf("*Risks:* %d open (%d critical/high)", len(risks), critical))
			}
		}
	}

	// Overdue from Jira
	if h.Jira != nil && boardID > 0 {
		result, err := h.Jira.SearchIssues(ctx, "resolution = Unresolved AND assignee = currentUser() AND updated <= -3d ORDER BY updated ASC", 5, 0)
		if err == nil && len(result.Issues) > 0 {
			var ob strings.Builder
			ob.WriteString(fmt.Sprintf("*Your Stale Items (%d):*\n", len(result.Issues)))
			for _, i := range result.Issues {
				ob.WriteString(fmt.Sprintf("- %s: %s\n", i.Key, i.Summary))
			}
			sections = append(sections, ob.String())
		}
	}

	if len(sections) == 0 {
		return textResult("Nothing to report today. All clear."), nil
	}

	digest := strings.Join(sections, "\n\n")
	title := fmt.Sprintf("Daily Digest — %s", time.Now().Format("Mon Jan 2"))

	// Send to primary channel
	if h.Slack != nil && h.Slack.Available() {
		_ = h.Slack.SendRichMessage(ctx, "", title, digest)
	} else if h.Lark != nil {
		_ = h.Lark.SendMarkdown(ctx, title, digest)
	} else if h.Teams != nil && h.Teams.Available() {
		_ = h.Teams.SendCard(ctx, title, digest)
	}

	return textResult(fmt.Sprintf("%s\n\n%s", title, digest)), nil
}

// routeNotification determines channels based on severity and audience.
func routeNotification(severity, audience string) []string {
	switch severity {
	case "critical":
		return []string{"telegram", "slack_dm", "lark", "teams", "discord"}
	case "high":
		switch audience {
		case "stakeholder", "executive":
			return []string{"slack_channel", "lark", "teams"}
		default:
			return []string{"slack_dm", "lark", "teams"}
		}
	case "medium":
		return []string{"slack_channel", "lark", "teams"}
	case "low":
		return []string{"slack_channel", "lark"}
	default: // info
		return []string{"slack_channel"}
	}
}

func severityColor(severity string) int {
	switch severity {
	case "critical":
		return 0xFF0000
	case "high":
		return 0xFF8C00
	case "medium":
		return 0xFFD700
	case "low":
		return 0x3498DB
	default:
		return 0x95A5A6
	}
}
