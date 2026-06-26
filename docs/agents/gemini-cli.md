# zara-jira-mcp — Gemini CLI Setup

## MCP Configuration

Edit `~/.gemini/settings.json`:

```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "zara-jira-mcp-wrapper"
    }
  }
}
```

Or with inline env (less recommended):

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

## Usage

```
$ gemini

> What's the sprint health for board 5?
> Prep my standup
> When will the remaining 12 items be done?
```

## Wrapper Script

See `docs/agents/README.md` for the recommended wrapper approach that keeps credentials out of config.
