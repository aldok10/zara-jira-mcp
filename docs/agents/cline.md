# zara-jira-mcp — Cline / Roo Code Setup

## MCP Configuration

Open VS Code > Cline panel > MCP Servers icon > Configure tab > Edit Config.

Or edit directly at the path shown in the Cline settings panel:

```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "zara-jira-mcp-wrapper",
      "disabled": false
    }
  }
}
```

For Roo Code, the config format is identical. Open Roo Code settings > MCP > Add Server.

## With Inline Env

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
      },
      "disabled": false
    }
  }
}
```

## Verify

After adding, the MCP server should appear with a green indicator in the Cline/Roo Code MCP panel. Ask:

> "List my Jira boards"

## Wrapper Script

See `docs/agents/README.md` for the recommended wrapper approach.
