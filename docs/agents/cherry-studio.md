# zara-jira-mcp — Cherry Studio Setup

Cherry Studio is an open-source multi-model desktop client with built-in MCP support.

## MCP Configuration

Open Cherry Studio > Settings > MCP Servers > Add.

Or configure via the config panel:

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

## Verify

After adding, the MCP server tools should appear in your chat. Ask:

> "List my Jira boards"

## Notes

- Cherry Studio supports stdio MCP servers natively
- Works with any model provider (OpenAI, Anthropic, Gemini, Ollama, etc.)
- See `docs/agents/README.md` for the wrapper script setup
