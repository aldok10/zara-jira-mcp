# zara-jira-mcp — Goose (Block) Setup

## Interactive Setup

```bash
goose configure
# Select "Add Extension"
# Name: jira-pm
# Type: stdio
# Command: zara-jira-mcp-wrapper
```

## Manual Config

Edit `~/.config/goose/config.yaml`:

```yaml
extensions:
  jira-pm:
    type: stdio
    command: zara-jira-mcp-wrapper
    enabled: true
```

Or with env vars:

```yaml
extensions:
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
    enabled: true
```

## Usage

```bash
$ goose

> What boards do I have access to?
> Prep my standup for board 5
> What's the forecast for the remaining 8 items?
```

## Wrapper Script

See `docs/agents/README.md` for the recommended wrapper approach that keeps credentials out of config.
