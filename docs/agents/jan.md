# zara-jira-mcp — Jan AI Setup

Jan is an open-source desktop AI client with MCP host support.

## MCP Configuration

Open Jan > Settings > MCP > Add Server.

Or edit the MCP config directly:

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

- Jan supports downloading different MCP clients and servers
- Works with local models (Ollama, llama.cpp) and cloud providers
- See `docs/agents/README.md` for the wrapper script setup
