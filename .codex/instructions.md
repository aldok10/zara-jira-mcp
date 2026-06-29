# Codex Instructions

## zara-jira-mcp

AI Scrum Master MCP server. 124 tools. Persistent SQLite memory. Monte Carlo forecasting. Anti-pattern detection. Multi-channel notifications (Lark, Slack, Discord, Telegram, Teams, Email).

## Stack

Go 1.26, uber-go/fx, mcp-go, go-sqlite3 (WAL), OpenAI-compatible API.

## Architecture

- `cmd/server/` — entry point
- `domain/` — interfaces (jira, memory, ai, lark)
- `internal/` — implementations
- `application/tools/` — MCP tool handlers
- `transport/` — MCP registration

## MCP Config

```toml
[mcp_servers.jira-pm]
command = "zara-jira-mcp"

[mcp_servers.jira-pm.env]
JIRA_BASE_URL = "https://company.atlassian.net"
JIRA_EMAIL = "you@company.com"
JIRA_API_TOKEN = "your-token"
JIRA_AI_BASE_URL = "https://api.openai.com"
JIRA_AI_API_KEY = "sk-..."
JIRA_AI_MODEL = "gpt-4o-mini"
```

## Key Tool Patterns

First call: `jira_boards` -> get board_id.

- Daily: `pm_standup_prep(board_id)`
- Planning: `pm_planning_prep(board_id)`
- Health: `pm_sprint_health(board_id)`
- Forecast: `pm_forecast(board_id, remaining_items:N)`
- Exec report: `pm_exec_report(board_id)` (NEVER pm_dashboard for execs)
- End sprint: `pm_snapshot_sprint(board_id, velocity:N)`

## Full Reference

See SKILL.md.
