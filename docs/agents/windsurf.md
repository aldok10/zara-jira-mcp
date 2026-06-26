# zara-jira-mcp — Windsurf Setup

## MCP Configuration

Open Windsurf Settings (Cmd+,) > search "MCP" > View raw config, or edit `~/.codeium/windsurf/mcp_config.json`:

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

## Windsurf Rules

Add to `.windsurfrules` or Windsurf's AI Rules settings:

```markdown
# PM/Scrum Master MCP

When user asks about project management, Jira, sprints, team health, or planning:

1. Use zara-jira-mcp tools (131 available)
2. Get board_id with `jira_boards` first
3. For status: `pm_dashboard(board_id)`
4. For standup: `pm_standup_prep(board_id)`
5. For forecasting: `pm_forecast(board_id, remaining_items)`
6. For stakeholders: `pm_exec_report(board_id)`
7. Record risks/decisions/blockers to persistent memory immediately
8. Run `pm_snapshot_sprint` at end of every sprint

Tool prefixes: jira_* (Jira ops), pm_* (PM intelligence), slack_* (notifications)
```

## Cascade Integration

Windsurf's Cascade agent automatically discovers MCP tools. Key behaviors:

- Tools show as available actions in Cascade's tool panel
- Cascade will auto-select relevant tools based on conversation
- For explicit tool use, mention the tool name in your prompt
- Memory persists at `~/.zara-jira-mcp/pm_memory.db`
