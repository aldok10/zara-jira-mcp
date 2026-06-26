# zara-jira-mcp — Zed Setup

Zed uses `context_servers` in its settings, not `mcpServers`.

## MCP Configuration

Press `Cmd+,` (macOS) or `Ctrl+,` (Linux) > Open Settings (JSON).

Add to `~/.config/zed/settings.json`:

```json
{
  "context_servers": {
    "jira-pm": {
      "command": {
        "path": "zara-jira-mcp-wrapper",
        "args": []
      }
    }
  }
}
```

Or with full path:

```json
{
  "context_servers": {
    "jira-pm": {
      "command": {
        "path": "/Users/you/.local/bin/zara-jira-mcp-wrapper",
        "args": []
      }
    }
  }
}
```

## With Inline Env

```json
{
  "context_servers": {
    "jira-pm": {
      "command": {
        "path": "zara-jira-mcp",
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
}
```

## Verify

Restart Zed. Open the Agent Panel and confirm the jira-pm server is active (green indicator).

## Note

Zed gotcha: the key is `context_servers`, not `mcpServers`. Zed infers stdio transport from the presence of `command` (vs `url` for remote servers).
