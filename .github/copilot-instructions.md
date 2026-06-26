# Copilot Instructions: zara-jira-mcp

## Project

AI-powered Scrum Master MCP server with 124 tools. Persistent SQLite memory, Monte Carlo forecasting, anti-pattern detection, multi-channel notifications.

Built with Go 1.26, uber-go/fx, mcp-go, SQLite WAL mode.

## Architecture

```
cmd/server/          Entry point (DI via uber-go/fx)
config/              Env-based configuration
domain/              Interfaces + models
internal/            Implementations (jira, ai, lark, memory, slack)
application/tools/   MCP tool handlers
transport/           MCP server + tool registration
```

## Using The PM Tools

First: `jira_boards` to get board_id.

Key tools:
- `pm_standup_prep(board_id)` — daily talking points
- `pm_planning_prep(board_id)` — sprint planning prep
- `pm_sprint_health(board_id)` — 0-100 health score
- `pm_forecast(board_id, remaining_items:N)` — Monte Carlo dates
- `pm_exec_report(board_id)` — executive summary (no jargon)
- `pm_snapshot_sprint(board_id, velocity:N)` — end-of-sprint capture
- `pm_record_risk/decision/blocker` — record immediately
- `pm_auto_detect_risks(board_id)` — weekly proactive scan

## Rules

- Never `pm_dashboard` for executives
- Sprint snapshots compound forecasting accuracy
- See SKILL.md for full tool reference
