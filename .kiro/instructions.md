# Kiro Instructions: zara-jira-mcp

AI-powered Scrum Master MCP server. 124 tools for Jira operations, sprint intelligence, risk management, forecasting, and multi-channel notifications.

## First Use

Call `jira_boards` to get board_id. Almost every PM tool needs it.

## Key Workflows

- Daily standup: `pm_standup_prep(board_id)`
- Sprint planning: `pm_planning_prep(board_id)`
- Health check: `pm_sprint_health(board_id)`
- Forecast: `pm_forecast(board_id, remaining_items:N)`
- Executive update: `pm_exec_report(board_id)` (NEVER pm_dashboard for execs)
- End of sprint: `pm_snapshot_sprint(board_id, velocity:N, carryover:N)`

## Record Immediately

- Risks: `pm_record_risk(title, severity, owner, mitigation)`
- Decisions: `pm_record_decision(title, decision, rationale)`
- Blockers: `pm_record_blocker(description, issue_key, owner)`

## Weekly

- `pm_auto_detect_risks(board_id)` — proactive risk scanning
- `pm_anti_patterns(board_id)` — detect team dysfunctions

## Full Reference

See SKILL.md for all 124 tools with parameters and usage patterns.
