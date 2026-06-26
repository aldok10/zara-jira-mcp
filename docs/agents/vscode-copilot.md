# zara-jira-mcp — VS Code + GitHub Copilot Setup

## MCP Configuration

Create `.vscode/mcp.json` in your project root:

```json
{
  "servers": {
    "jira-pm": {
      "command": "zara-jira-mcp-wrapper",
      "type": "stdio"
    }
  }
}
```

Or add to VS Code user settings (`settings.json`):

```json
{
  "github.copilot.chat.mcp.servers": {
    "jira-pm": {
      "command": "zara-jira-mcp-wrapper",
      "type": "stdio"
    }
  }
}
```

## With Inline Env

```json
{
  "servers": {
    "jira-pm": {
      "command": "zara-jira-mcp",
      "type": "stdio",
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

## Copilot Instructions

Add `.github/copilot-instructions.md` to your project (already included in this repo) for Copilot to understand how to use the PM tools.

## Verify

Open Copilot Chat and ask:

> @workspace List my Jira boards

The MCP tools should be available in the Copilot agent mode.

## Wrapper Script

See `docs/agents/README.md` for the recommended wrapper approach.
