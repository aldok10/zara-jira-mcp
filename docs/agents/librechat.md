# zara-jira-mcp — LibreChat Setup

LibreChat configures MCP servers in `librechat.yaml`.

## Configuration

Add to your `librechat.yaml`:

```yaml
mcpServers:
  jira-pm:
    type: stdio
    command: zara-jira-mcp-wrapper
```

Or with env vars:

```yaml
mcpServers:
  jira-pm:
    type: stdio
    command: zara-jira-mcp
    env:
      JIRA_BASE_URL: https://company.atlassian.net
      JIRA_EMAIL: you@company.com
      JIRA_API_TOKEN: your-token
      JIRA_AI_BASE_URL: https://api.openai.com
      JIRA_AI_API_KEY: sk-...
      JIRA_AI_MODEL: gpt-4o-mini
```

## Restart Required

After editing `librechat.yaml`, restart LibreChat to initialize the MCP connection.

## Verify

In chat, the MCP tools should be available. Ask:

> "List my Jira boards"

## Notes

- LibreChat supports both stdio and SSE MCP transports
- Works with any configured model provider
- See `docs/agents/README.md` for the wrapper script setup
