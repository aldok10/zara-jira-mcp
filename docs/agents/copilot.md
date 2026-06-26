# zara-jira-mcp — GitHub Copilot (VS Code) Setup

## MCP Configuration

Add to VS Code settings (`settings.json`) or use the MCP panel:

```json
{
  "github.copilot.chat.mcp.servers": {
    "jira-pm": {
      "type": "stdio",
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

Alternative: `.vscode/mcp.json` (project-level):

```json
{
  "servers": {
    "jira-pm": {
      "type": "stdio",
      "command": "zara-jira-mcp",
      "env": {
        "JIRA_BASE_URL": "https://company.atlassian.net",
        "JIRA_EMAIL": "you@company.com",
        "JIRA_API_TOKEN": "${env:JIRA_API_TOKEN}",
        "JIRA_AI_BASE_URL": "${env:JIRA_AI_BASE_URL}",
        "JIRA_AI_API_KEY": "${env:JIRA_AI_API_KEY}",
        "JIRA_AI_MODEL": "gpt-4o-mini"
      }
    }
  }
}
```

## Copilot Instructions

Add to `.github/copilot-instructions.md`:

```markdown
## PM/Scrum Master MCP (zara-jira-mcp)

This project has a Jira PM MCP server connected. Use it for:
- Sprint management and status
- Risk, blocker, decision tracking
- Team health and velocity
- Forecasting and capacity planning
- Stakeholder reporting

### Key Tools
- `pm_dashboard(board_id)` — full PM overview
- `pm_standup_prep(board_id)` — standup talking points
- `pm_forecast(board_id, remaining_items)` — delivery forecast
- `pm_record_risk/decision/blocker` — persistent memory
- `pm_exec_report(board_id)` — stakeholder-friendly report
- `jira_search(jql)` — search Jira issues
- `jira_create_issue(project, summary)` — create tickets

### Rules
- Get board_id with `jira_boards` first
- Always record decisions and risks to memory
- Use `pm_exec_report` for executives (no technical details)
- `pm_snapshot_sprint` at end of every sprint
```

## Notes

- VS Code Copilot MCP support requires GitHub Copilot Chat extension
- Tools appear in Copilot Chat when you mention project management topics
- Use `#tool:jira-pm` in chat to explicitly invoke MCP tools
- All 131 tools available but Copilot will auto-select relevant ones
