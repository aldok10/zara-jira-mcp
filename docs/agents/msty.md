# zara-jira-mcp — Msty Studio Setup

Msty supports MCP servers via its Toolbox feature.

## Configuration

Open Msty > Toolbox > MCP Servers > Add Server.

Configure with:
- **Name:** jira-pm
- **Transport:** stdio
- **Command:** `zara-jira-mcp-wrapper`

Or manually edit the config:

```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "zara-jira-mcp-wrapper",
      "args": []
    }
  }
}
```

## With Inline Env

```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "zara-jira-mcp",
      "args": [],
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

## Notes

- Msty supports local models (Ollama) + cloud providers
- MCP tools are available in any chat session after setup
- See `docs/agents/README.md` for the wrapper script setup
