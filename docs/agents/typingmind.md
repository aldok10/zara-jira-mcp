# zara-jira-mcp — TypingMind Setup

TypingMind supports remote MCP servers natively and local MCP via the Private MCP Connector.

## Option 1: Private MCP Connector (Local)

1. Install the TypingMind Private MCP Connector
2. Add server config:

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

## Option 2: Remote MCP (if you deploy as HTTP server)

If you deploy zara-jira-mcp as a remote server with Streamable HTTP transport:

Go to TypingMind > Plugins > MCP > Add MCP URL:

```
https://your-server.com/mcp
```

## With Inline Env (Private Connector)

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

- TypingMind is web-based, so local MCP needs the Private Connector bridge
- Works with any model provider configured in TypingMind
- See `docs/agents/README.md` for the wrapper script setup
