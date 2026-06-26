# zara-jira-mcp — Claude Code Setup

## MCP Configuration

Add to `.claude/settings.json` (project-level) or `~/.claude/settings.json` (global):

```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "zara-jira-mcp",
      "env": {
        "JIRA_BASE_URL": "https://company.atlassian.net",
        "JIRA_EMAIL": "you@company.com",
        "JIRA_API_TOKEN": "your-token",
        "JIRA_AI_BASE_URL": "https://api.openai.com",
        "JIRA_AI_API_KEY": "sk-...",
        "JIRA_AI_MODEL": "gpt-4o-mini"
      }
    }
  }
}
```

## CLAUDE.md Instructions

Add to your project root `CLAUDE.md`:

```markdown
## MCP: Jira PM Brain (zara-jira-mcp)

131-tool MCP server for PM/Scrum Master workflows with persistent memory.

### Quick Commands
- Status check: `pm_dashboard(board_id:X)`
- Standup prep: `pm_standup_prep(board_id:X)`
- Forecast: `pm_forecast(board_id:X, remaining_items:N)`
- Record risk: `pm_record_risk(title, severity, owner, mitigation)`
- Record decision: `pm_record_decision(title, decision, rationale, tags)`

### Rules
- Get board_id with `jira_boards` first
- Call `pm_snapshot_sprint` at end of every sprint
- Use `pm_exec_report` for stakeholders (not pm_dashboard)
- Record decisions, risks, blockers immediately
- Run `pm_auto_detect_risks` weekly

### Tool Categories
- `jira_*` — Jira CRUD and operations
- `pm_*` — PM intelligence, memory, forecasting
- `slack_*` — Slack notifications
- `portfolio_*` — Cross-project views

### Memory
SQLite at ~/.zara-jira-mcp/pm_memory.db. Persists across sessions.
AI tools get better with more recorded data (snapshots, decisions, metrics).
```

## Slash Command (optional)

`.claude/commands/pm-status.md`:

```markdown
---
description: Get full PM dashboard
---

Call pm_dashboard with board_id from jira_boards. Show sprint progress, health, risks, blockers.
```

`.claude/commands/standup.md`:

```markdown
---
description: Prepare for daily standup
---

Call pm_standup_prep. Show talking points, blockers, action items.
```
