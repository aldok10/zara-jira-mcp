# zara-jira-mcp — Cursor Setup

## MCP Configuration

Open Cursor Settings > MCP, or edit `~/.cursor/mcp.json`:

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

## Cursor Rules

Add to `.cursor/rules/pm-brain.mdc`:

```markdown
---
description: PM/Scrum Master MCP with persistent memory
globs: ["**/*"]
alwaysApply: false
---

# Jira PM Brain (zara-jira-mcp)

131-tool MCP for project management with AI intelligence and persistent memory.

## When to Use

Activate these tools when user asks about:
- Sprint status, Jira issues, team workload
- Risks, blockers, decisions, dependencies
- Retrospectives, planning, forecasting
- Team health, velocity, capacity
- Stakeholder reporting

## Quick Reference

| Need | Tool |
|------|------|
| Full status | `pm_dashboard(board_id)` |
| Standup | `pm_standup_prep(board_id)` |
| "When done?" | `pm_forecast(board_id, remaining_items)` |
| Record risk | `pm_record_risk(title, severity, owner, mitigation)` |
| Record decision | `pm_record_decision(title, decision, rationale)` |
| End of sprint | `pm_snapshot_sprint(board_id, velocity)` |
| For executives | `pm_exec_report(board_id)` |
| Team problems | `pm_anti_patterns(board_id)` |
| Coaching | `pm_coaching(topic, situation)` |

## Rules

1. Get board_id with `jira_boards` first (do once per session)
2. Memory accumulates — more data = better AI recommendations
3. Use `pm_exec_report` for stakeholders, `pm_dashboard` for team
4. Record risks/decisions/blockers immediately when mentioned
5. `pm_snapshot_sprint` at every sprint end (essential for forecasting)
```

## Project-Level Rules (`.cursorrules`)

If using legacy `.cursorrules` format:

```
When user mentions Jira, sprints, team management, or PM work:
- Use zara-jira-mcp MCP tools (prefixed jira_* and pm_*)
- Get board_id with jira_boards first
- pm_dashboard for status, pm_standup_prep for standups
- Record decisions/risks/blockers to PM memory immediately
- pm_exec_report for stakeholder updates (never raw data)
```
