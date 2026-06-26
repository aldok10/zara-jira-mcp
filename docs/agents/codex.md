# zara-jira-mcp — OpenAI Codex CLI Setup

## MCP Configuration

Add to `~/.codex/config.toml` or project `.codex/config.toml`:

```toml
[mcp_servers.jira-pm]
command = "zara-jira-mcp"
args = []

[mcp_servers.jira-pm.env]
JIRA_BASE_URL = "https://company.atlassian.net"
JIRA_EMAIL = "you@company.com"
JIRA_API_TOKEN = "your-token"
JIRA_AI_BASE_URL = "https://api.openai.com"
JIRA_AI_API_KEY = "sk-..."
JIRA_AI_MODEL = "gpt-4o-mini"
```

## Codex Instructions

Add to `codex.md` or `AGENTS.md` in project root:

```markdown
## MCP: Jira PM Brain

A 131-tool MCP for PM/Scrum Master work with persistent memory.

### Available Tools (key ones)

**Jira**: jira_search, jira_get_issue, jira_create_issue, jira_transition, jira_boards, jira_sprint_summary
**PM Intel**: pm_dashboard, pm_standup_prep, pm_forecast, pm_recommendations, pm_flow_metrics
**Memory**: pm_record_risk, pm_record_decision, pm_record_blocker, pm_snapshot_sprint
**Reports**: pm_exec_report, pm_release_notes, pm_weekly_digest, pm_scorecard

### Usage Pattern

1. `jira_boards` to get board_id (do this once)
2. `pm_dashboard(board_id:X)` for full overview
3. `pm_standup_prep(board_id:X)` before standup
4. `pm_planning_prep(board_id:X)` before planning
5. Record decisions/risks/blockers as they come up
6. `pm_snapshot_sprint(board_id:X, velocity:N)` at sprint end

### Important
- Codex requires local STDIO transport (no remote servers)
- Memory persists in ~/.zara-jira-mcp/pm_memory.db
- AI tools need JIRA_AI_* env vars configured
- Run `pm_auto_detect_risks(board_id)` for proactive scanning
```

## Notes for Codex

- Codex only supports STDIO transport — this MCP is already STDIO-native
- Binary must be in PATH or use absolute path in config
- All 131 tools are available via tools/list
- No remote/SSE needed
