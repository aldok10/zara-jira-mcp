# zara-jira-mcp — Kiro / AWS Q Developer Setup

## MCP Configuration

Add to `kiro.json` or `.kiro/settings.json`:

```json
{
  "mcpServers": {
    "jira-pm": {
      "command": ["zara-jira-mcp"],
      "env": {
        "JIRA_BASE_URL": "https://company.atlassian.net",
        "JIRA_EMAIL": "you@company.com",
        "JIRA_API_TOKEN": "your-token",
        "JIRA_AI_BASE_URL": "https://api.openai.com",
        "JIRA_AI_API_KEY": "sk-...",
        "JIRA_AI_MODEL": "gpt-4o-mini"
      },
      "timeout": 30000
    }
  }
}
```

## Instructions

Add PM context to your project instructions:

```markdown
## MCP: Jira PM Brain (131 tools)

Persistent Scrum Master brain connected via MCP.

### Core Workflow
1. `jira_boards` — get board_id (once)
2. `pm_dashboard(board_id)` — full status
3. `pm_standup_prep(board_id)` — standup prep
4. `pm_forecast(board_id)` — Monte Carlo forecast
5. Record decisions/risks/blockers to memory

### Key Tools
- `pm_recommendations(board_id, focus)` — AI advice
- `pm_anti_patterns(board_id)` — detect dysfunctions
- `pm_coaching(topic, situation)` — coaching tips
- `pm_exec_report(board_id)` — stakeholder report
- `pm_planning_prep(board_id)` — sprint planning package
- `pm_snapshot_sprint(board_id, velocity)` — end-of-sprint capture
```
