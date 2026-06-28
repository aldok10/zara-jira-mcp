# PM Brain — zara-jira-mcp

## About

124-tool MCP server. AI Scrum Master with persistent memory, Monte Carlo forecasting, anti-pattern detection, and multi-channel notifications.

## First Use

Call `jira_boards` to get board_id. Store it — needed for almost every PM tool.

## Quick Reference

| Task | Tool |
|------|------|
| Standup prep | `pm_standup_prep(board_id)` |
| Sprint planning | `pm_planning_prep(board_id)` |
| Health check | `pm_sprint_health(board_id)` |
| Forecasting | `pm_forecast(board_id, remaining_items:N)` |
| Record risk | `pm_record_risk(title, severity, owner, mitigation)` |
| Record decision | `pm_record_decision(title, decision, rationale)` |
| Record blocker | `pm_record_blocker(description, issue_key, owner)` |
| Executive report | `pm_exec_report(board_id)` |
| End of sprint | `pm_snapshot_sprint(board_id, velocity:N)` |
| Anti-patterns | `pm_anti_patterns(board_id)` |
| Scope creep | `pm_scope_creep(board_id)` |

## Rules

1. Never use `pm_dashboard` for executives — use `pm_exec_report`
2. Record decisions/risks/blockers immediately, don't wait
3. `pm_snapshot_sprint` at end of every sprint — forecasting needs history
4. Run `pm_auto_detect_risks` weekly for proactive scanning
5. See SKILL.md for full 124-tool documentation
