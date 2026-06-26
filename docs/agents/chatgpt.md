# zara-jira-mcp — ChatGPT Desktop Setup

## MCP Configuration

ChatGPT Desktop (macOS/Windows) supports MCP servers. Configure in Settings > MCP:

```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "/path/to/zara-jira-mcp",
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

## Custom Instructions

Add to ChatGPT system instructions or custom GPT:

```
You have access to a Jira PM MCP server with 131 tools for project management.

Key tools:
- pm_dashboard(board_id) — full sprint status
- pm_standup_prep(board_id) — standup talking points
- pm_forecast(board_id, remaining_items) — "when will it be done?" with probabilities
- pm_exec_report(board_id) — executive summary
- pm_record_risk/decision/blocker — persistent memory
- jira_search(jql) — search Jira
- jira_create_issue(project, summary) — create tickets

Always get board_id first with jira_boards.
Record decisions and risks immediately to memory.
Use pm_exec_report for stakeholders, pm_dashboard for the team.
```

## Notes

- ChatGPT Desktop MCP launched 2025, now standard in 2026
- Use absolute path for command (ChatGPT doesn't use shell PATH)
- Tools appear as actions in the conversation
- Memory persists in ~/.zara-jira-mcp/pm_memory.db
