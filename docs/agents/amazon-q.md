# zara-jira-mcp — Amazon Q Developer Setup

Amazon Q CLI uses `~/.aws/amazonq/mcp.json` for global MCP config, or `.amazonq/mcp.json` for per-project.

## Global Config

`~/.aws/amazonq/mcp.json`:

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

## Per-Project Config

`.amazonq/mcp.json` in your project root:

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

## Usage

```bash
$ q chat

> What are my active sprints?
> Prep my standup for board 5
> What's the sprint health score?
```

## Notes

- Amazon Q supports both local (stdio) and remote (HTTP) MCP servers
- The config format is identical to Claude Desktop
- See `docs/agents/README.md` for the wrapper script setup
