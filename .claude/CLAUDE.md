# Claude Code Instructions

## Project: zara-jira-mcp

AI-powered Scrum Master MCP server. 124 tools for Jira, sprint intelligence, risk management, forecasting, and notifications.

## When Working On This Codebase

- Go 1.26, uber-go/fx for DI, SQLite for persistence
- Domain-driven architecture: domain/ has interfaces, internal/ has implementations
- All MCP tools live in application/tools/ (12 handler files)
- Transport layer in transport/ registers tools with mcp-go

## Using The MCP Tools (as a PM/SM)

First call: `jira_boards` to get board_id.

### Daily
- `pm_standup_prep(board_id)` — talking points, blockers, risks

### Sprint Lifecycle
- Planning: `pm_planning_prep(board_id)`
- Mid-sprint: `pm_sprint_health(board_id)`, `pm_scope_creep(board_id)`
- End: `pm_snapshot_sprint(board_id, velocity:N)`, `pm_scorecard(board_id)`

### Record Immediately
- Decisions: `pm_record_decision(title, decision, rationale)`
- Risks: `pm_record_risk(title, severity, owner, mitigation)`
- Blockers: `pm_record_blocker(description, issue_key, owner)`

### Forecasting
- `pm_forecast(board_id, remaining_items:N)` — Monte Carlo probabilities

### For Executives (NEVER pm_dashboard)
- `pm_exec_report(board_id)` — business language, no jargon

## Rules
- Memory compounds. More sprint snapshots = better forecasts.
- Run `pm_auto_detect_risks(board_id)` weekly.
- See SKILL.md for complete tool reference.
