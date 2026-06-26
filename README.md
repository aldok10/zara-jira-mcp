# zara-jira-mcp

MCP server for AI-powered Jira intelligence. Reads tickets, analyzes sprint health, and pushes summaries to Lark.

## Tools

| Tool | Description |
|------|-------------|
| `jira_search` | Search issues with JQL |
| `jira_get_issue` | Get full issue details |
| `jira_boards` | List accessible boards |
| `jira_sprint_summary` | Active sprint breakdown |
| `jira_ai_analyze` | AI-powered ticket analysis (ask questions about your board) |
| `jira_ai_sprint_report` | AI sprint health report, optionally sent to Lark |
| `jira_notify_lark` | Send markdown message to Lark group |

## Setup

```bash
cp .env.example .env
# Edit .env with your credentials

make build
# or: go build -o bin/zara-jira-mcp ./cmd/server
```

## Usage (MCP stdio)

```json
{
  "mcpServers": {
    "jira": {
      "command": "/path/to/zara-jira-mcp",
      "env": {
        "JIRA_BASE_URL": "https://company.atlassian.net",
        "JIRA_EMAIL": "you@company.com",
        "JIRA_API_TOKEN": "..."
      }
    }
  }
}
```

## Architecture

```
cmd/server/         - Entry point (fx wiring)
config/             - Env-based configuration
domain/             - Interfaces (jira.Client, ai.Provider, lark.Notifier)
internal/jira/      - Jira REST API v3 + Agile API client
internal/ai/        - OpenAI-compatible completion client
internal/lark/      - Lark webhook client
internal/bootstrap/ - DI module (uber-go/fx)
application/tools/  - MCP tool handlers
transport/          - MCP server + tool registration
```
