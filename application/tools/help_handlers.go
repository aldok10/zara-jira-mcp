package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

var toolCategories = map[string][]struct {
	name string
	desc string
}{
	"Getting Started": {
		{"pm_help", "Discover tools by topic"},
		{"pm_quickstart", "First-time setup guide"},
		{"pm_workflow", "Pre-built workflow recipes"},
		{"pm_dashboard", "Full PM overview in one shot"},
		{"jira_health", "Server health check"},
	},
	"Daily Work": {
		{"jira_my_issues", "Your assigned issues"},
		{"pm_standup_prep", "Daily standup talking points"},
		{"pm_calendar_today", "Today's meetings"},
		{"pm_recipe_start_work", "Start working on an issue (assign + transition)"},
		{"pm_recipe_done", "Mark issue done (transition + log time)"},
		{"pm_github_prs", "Open PRs needing attention"},
	},
	"Sprint Management": {
		{"jira_sprint_summary", "Active sprint breakdown"},
		{"pm_sprint_health", "Sprint health score (0-100)"},
		{"pm_burndown", "Sprint burndown data"},
		{"pm_forecast", "Monte Carlo 'when will it be done?'"},
		{"pm_snapshot_sprint", "Save sprint state to memory"},
		{"pm_scope_creep", "Detect mid-sprint scope changes"},
	},
	"Risks & Blockers": {
		{"pm_record_risk", "Log a new risk"},
		{"pm_risk_dashboard", "All open risks by severity"},
		{"pm_record_blocker", "Record an impediment"},
		{"pm_blockers", "Active blockers"},
		{"pm_auto_detect_risks", "AI scan for risk signals"},
		{"pm_incidents", "PagerDuty incidents"},
	},
	"Team & People": {
		{"pm_team_health", "Team workload overview"},
		{"jira_workload", "Issue distribution per member"},
		{"pm_oncall", "Who's on call now"},
		{"pm_time_report", "Time tracked by team"},
		{"pm_confidence", "Team confidence voting"},
	},
	"Decisions & Learning": {
		{"pm_record_decision", "Record a decision with rationale"},
		{"pm_search_decisions", "Search decision log"},
		{"pm_record_learning", "Record team knowledge"},
		{"pm_team_kb", "Team knowledge base + Q&A"},
	},
	"Retrospectives": {
		{"pm_record_retro", "Record retro outcomes"},
		{"pm_action_items", "Pending retro actions"},
		{"pm_retro_analysis", "AI pattern detection across retros"},
		{"pm_experiment", "Improvement experiments"},
		{"pm_facilitate", "Ceremony facilitation prompts"},
	},
	"Planning": {
		{"pm_planning_prep", "Complete planning package"},
		{"pm_capacity_plan", "Data-driven capacity recommendation"},
		{"pm_check_ready", "Story readiness check (INVEST + DoR)"},
		{"pm_backlog_groom", "Find stale backlog items"},
		{"pm_linear_cycles", "Linear sprint cycles"},
	},
	"Reporting": {
		{"pm_exec_report", "Executive stakeholder report"},
		{"pm_weekly_digest", "Weekly team digest"},
		{"pm_release_notes", "Sprint release notes"},
		{"pm_scorecard", "Sprint grade (A-F)"},
		{"pm_sprint_compare", "This sprint vs last"},
	},
	"Notifications": {
		{"jira_notify_lark", "Send to Lark"},
		{"pm_notify_slack", "Send to Slack"},
		{"pm_notify_discord", "Send to Discord"},
		{"pm_notify_telegram", "Send to Telegram"},
		{"pm_notify_teams", "Send to MS Teams"},
		{"pm_notify_email", "Send email"},
	},
	"Integrations": {
		{"pm_github_prs", "GitHub PRs"},
		{"pm_github_activity", "GitHub repo activity"},
		{"pm_notion_search", "Search Notion"},
		{"pm_linear_issues", "Linear issues"},
		{"pm_calendar_events", "Calendar events"},
		{"pm_sheet_read", "Read Google Sheet"},
	},
	"AI Analysis": {
		{"jira_ai_analyze", "AI analysis of tickets"},
		{"pm_recommendations", "AI recommendations from history"},
		{"pm_coaching", "Data-driven coaching advice"},
		{"pm_anti_patterns", "Detect scrum anti-patterns"},
		{"pm_nl_to_jql", "Natural language to JQL"},
	},
}

var categoryOrder = []string{
	"Getting Started", "Daily Work", "Sprint Management",
	"Risks & Blockers", "Team & People", "Decisions & Learning",
	"Retrospectives", "Planning", "Reporting", "Notifications",
	"Integrations", "AI Analysis",
}

// PMHelp provides tool discovery by topic.
func (h *Handlers) PMHelp(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topic := req.GetString("topic", "")

	if topic == "" {
		var sb strings.Builder
		sb.WriteString("# zara-jira-mcp Tool Guide\n\n")
		for _, cat := range categoryOrder {
			tools := toolCategories[cat]
			sb.WriteString(fmt.Sprintf("## %s\n", cat))
			for _, t := range tools {
				sb.WriteString(fmt.Sprintf("  - `%s` — %s\n", t.name, t.desc))
			}
			sb.WriteString("\n")
		}
		sb.WriteString("---\nTip: Use `pm_help topic=\"sprint\"` to get detailed info on a specific area.\n")
		sb.WriteString("Tip: Use `pm_workflow workflow=\"standup\"` for step-by-step recipes.\n")
		return mcp.NewToolResultText(sb.String()), nil
	}

	// Search by topic
	topic = strings.ToLower(topic)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Tools for: %s\n\n", topic))
	found := false
	for cat, tools := range toolCategories {
		if strings.Contains(strings.ToLower(cat), topic) {
			found = true
			sb.WriteString(fmt.Sprintf("## %s\n", cat))
			for _, t := range tools {
				sb.WriteString(fmt.Sprintf("  - `%s` — %s\n", t.name, t.desc))
			}
			sb.WriteString("\n")
		}
	}
	if !found {
		// Search in tool descriptions
		sb.WriteString("Matching tools:\n")
		for _, cat := range categoryOrder {
			for _, t := range toolCategories[cat] {
				if strings.Contains(strings.ToLower(t.name+t.desc), topic) {
					found = true
					sb.WriteString(fmt.Sprintf("  - `%s` — %s (in %s)\n", t.name, t.desc, cat))
				}
			}
		}
	}
	if !found {
		sb.WriteString("No tools found for that topic. Try: sprint, risks, team, planning, reporting, notifications, integrations, daily, retro, decisions\n")
	}
	return mcp.NewToolResultText(sb.String()), nil
}

// PMQuickstart provides first-time getting started guidance.
func (h *Handlers) PMQuickstart(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var sb strings.Builder
	sb.WriteString("# Welcome to zara-jira-mcp!\n\n")
	sb.WriteString("## Your Integrations\n\n")

	checks := []struct {
		name       string
		configured bool
	}{
		{"Jira", h.Config.Jira.BaseURL != ""},
		{"AI Provider", h.Config.AI.BaseURL != "" && h.Config.AI.APIKey != ""},
		{"Slack", h.Config.Slack.BotToken != "" || h.Config.Slack.WebhookURL != ""},
		{"Lark", h.Config.Lark.WebhookURL != ""},
		{"Discord", h.Config.Discord.BotToken != ""},
		{"Telegram", h.Config.Telegram.BotToken != ""},
		{"Teams", h.Config.Teams.WebhookURL != ""},
		{"Email", h.Config.Email.SMTPHost != ""},
		{"Confluence", h.Config.Confluence.BaseURL != ""},
		{"GitHub", h.Config.GitHub.Token != ""},
		{"Notion", h.Config.Notion.APIKey != ""},
		{"Google Calendar", h.Config.GoogleCalendar.APIKey != ""},
		{"Linear", h.Config.Linear.APIKey != ""},
		{"PagerDuty", h.Config.PagerDuty.APIKey != ""},
		{"Clockify", h.Config.Clockify.APIKey != ""},
		{"Google Sheets", h.Config.GoogleSheets.APIKey != ""},
	}

	ready := 0
	for _, c := range checks {
		status := "[ ] Not configured"
		if c.configured {
			status = "[x] Ready"
			ready++
		}
		sb.WriteString(fmt.Sprintf("  %s %s\n", status, c.name))
	}
	sb.WriteString(fmt.Sprintf("\n%d/%d integrations configured.\n\n", ready, len(checks)))

	sb.WriteString("## Get Started (3 steps)\n\n")
	sb.WriteString("1. Run `pm_dashboard` — see your full project status at a glance\n")
	sb.WriteString("2. Run `pm_standup_prep board_id=YOUR_BOARD` — get daily talking points\n")
	sb.WriteString("3. Run `pm_workflow workflow=\"standup\"` — learn the daily workflow\n\n")

	sb.WriteString("## Most Powerful Workflows\n\n")
	sb.WriteString("**Daily Standup** — `pm_workflow workflow=\"standup\"`\n")
	sb.WriteString("  Calendar + blockers + PRs + talking points in 30 seconds.\n\n")
	sb.WriteString("**Sprint End** — `pm_workflow workflow=\"sprint_end\"`\n")
	sb.WriteString("  Health score + release notes + retro + exec report. Complete sprint closure.\n\n")
	sb.WriteString("**Risk Management** — `pm_workflow workflow=\"incident\"`\n")
	sb.WriteString("  Incidents + auto-detect risks + escalation. Proactive risk radar.\n\n")

	sb.WriteString("## Tips\n\n")
	sb.WriteString("- Use `pm_help` to discover tools by topic\n")
	sb.WriteString("- Use `pm_help topic=\"risks\"` for specific areas\n")
	sb.WriteString("- AI tools need JIRA_AI_BASE_URL + JIRA_AI_API_KEY configured\n")
	sb.WriteString("- All notification platforms are optional — configure only what you use\n")

	return mcp.NewToolResultText(sb.String()), nil
}

var workflows = map[string]string{
	"standup": `# Daily Standup Workflow

1. pm_calendar_today → See today's meetings
2. pm_standup_prep board_id=X → AI-generated talking points
3. pm_blockers → Check active impediments
4. pm_incidents → Any production issues?
5. pm_github_prs → PRs needing review

Tip: Replace X with your board ID (get it from jira_boards).`,

	"sprint_start": `# Sprint Start Workflow

1. pm_planning_prep board_id=X → Full planning package (velocity, capacity, carryover)
2. pm_set_sprint_goal board_id=X goal="..." → Define sprint goal
3. pm_check_ready key=ISSUE-123 → Verify top stories are ready (repeat per story)
4. pm_capacity_plan board_id=X team_size=N → Recommended story points
5. pm_confidence sprint_name="Sprint X" score=4 → Record team confidence`,

	"sprint_end": `# Sprint End Workflow

1. pm_sprint_health board_id=X → Health score
2. pm_scorecard board_id=X → Sprint grade (A-F)
3. pm_snapshot_sprint board_id=X → Save sprint to memory (velocity tracking)
4. pm_release_notes board_id=X → Generate categorized release notes
5. pm_record_retro sprint_name="Sprint X" → Capture retrospective
6. pm_close_sprint_goal goal_id=X status=achieved → Close sprint goal
7. pm_exec_report board_id=X → Executive stakeholder summary`,

	"planning": `# Sprint Planning Workflow

1. pm_planning_prep board_id=X → Velocity history + carryover + risks + dependencies
2. pm_backlog_groom project=PROJ → Find stale items to clean up
3. pm_capacity_plan board_id=X team_size=N sprint_days=10 → Recommended capacity
4. pm_check_ready key=ISSUE-123 → Evaluate story readiness (repeat per candidate)
5. pm_forecast board_id=X remaining_items=N → "When will it be done?"
6. pm_set_sprint_goal board_id=X goal="..." → Define the sprint goal`,

	"retro": `# Retrospective Workflow

1. pm_facilitate ceremony=retro → Fresh retro facilitation format
2. pm_sprint_compare board_id=X → This sprint vs last (data for discussion)
3. pm_anti_patterns board_id=X → Auto-detected process issues
4. pm_retro_analysis board_id=X → AI pattern detection across past retros
5. pm_record_retro sprint_name="Sprint X" went_well="..." improvements="..." → Save outcomes
6. pm_experiment hypothesis="..." action="..." → Track improvement experiment`,

	"incident": `# Incident Response Workflow

1. pm_incidents → Current open incidents
2. pm_incident_summary → Impact on sprint (duration, count)
3. pm_oncall → Who's handling it
4. pm_auto_detect_risks board_id=X → Scan for new risk signals
5. pm_record_risk title="..." severity=high → Log to risk register
6. pm_escalate board_id=X → Auto-escalate critical items`,

	"weekly_review": `# Weekly Review Workflow

1. pm_weekly_digest board_id=X → AI weekly summary
2. pm_github_activity days=7 → Repo activity this week
3. pm_time_report range=week → Hours tracked by team
4. pm_risk_dashboard → Open risks status
5. pm_action_items → Pending retro actions
6. pm_experiments → Active improvement experiments
7. pm_health_history board_id=X → Health trend over time`,
}

// PMWorkflow returns step-by-step tool sequences for common PM activities.
func (h *Handlers) PMWorkflow(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	wf, err := req.RequireString("workflow")
	if err != nil {
		return errorResult("workflow required. Options: standup, sprint_start, sprint_end, planning, retro, incident, weekly_review"), nil
	}

	recipe, ok := workflows[wf]
	if !ok {
		available := make([]string, 0, len(workflows))
		for k := range workflows {
			available = append(available, k)
		}
		return errorResult(fmt.Sprintf("Unknown workflow '%s'. Available: %s", wf, strings.Join(available, ", "))), nil
	}

	return mcp.NewToolResultText(recipe), nil
}
