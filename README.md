# zara-jira-mcp

AI-powered Scrum Master MCP server. 86 tools for Jira intelligence, sprint management, team health tracking, and Lark notifications. Designed as Zara's PM brain.

## Quick Start

```bash
cp .env.example .env
# Edit .env with your credentials

make build
make install  # copies to ~/.local/bin/
```

## Tools (86)

### Jira Core

| Tool | Description |
|------|-------------|
| `jira_search` | Search issues with JQL |
| `jira_get_issue` | Full issue details |
| `jira_boards` | List boards |
| `jira_sprint_summary` | Active sprint breakdown |
| `jira_create_issue` | Create new issue |
| `jira_update_issue` | Update issue fields |
| `jira_add_comment` | Add comment |
| `jira_transitions` | List available transitions |
| `jira_transition` | Change issue status |
| `jira_assign` / `jira_unassign` | Manage assignees |
| `jira_find_user` | Search users by name/email |
| `jira_my_issues` | Current user's issues |
| `jira_overdue` | Stale issues (N days no update) |
| `jira_workload` | Team workload distribution |
| `jira_link_issues` | Create issue relationships |
| `jira_worklog_add` / `jira_worklog_list` | Time tracking |
| `jira_bulk_transition` / `jira_bulk_assign` | Batch operations |

### PM Intelligence

| Tool | Description |
|------|-------------|
| `pm_recommendations` | AI recommendations from historical memory |
| `pm_standup_prep` | Daily standup brief |
| `pm_velocity_trend` | Velocity over sprints |
| `pm_burndown` | Sprint burndown status |
| `pm_flow_metrics` | Cycle time, throughput, WIP |
| `pm_sprint_comparison` | Current vs previous sprint |
| `pm_ceremony_facilitator` | AI facilitation guide for ceremonies |
| `pm_forecast` | Monte Carlo sprint forecast |
| `pm_anti_patterns` | Detect Scrum anti-patterns |
| `pm_coaching` | Context-aware coaching advice |
| `pm_retro_analysis` | AI pattern analysis across retros |
| `pm_dashboard` | Full PM overview |

### PM Memory (Persistent)

| Tool | Description |
|------|-------------|
| `pm_snapshot_sprint` | Save sprint state for trend tracking |
| `pm_record_risk` / `pm_update_risk` | Risk register |
| `pm_risk_dashboard` | Open risks by severity |
| `pm_record_decision` / `pm_search_decisions` | Decision log |
| `pm_record_blocker` / `pm_resolve_blocker` | Blocker tracker |
| `pm_record_team_metric` / `pm_team_health` | Team performance |
| `pm_record_retro` / `pm_action_items` | Retrospective outcomes |
| `pm_confidence_vote` | Team confidence tracking |
| `pm_set_sprint_goal` / `pm_goal_check` | Sprint goal management |

### Lark / AI

| Tool | Description |
|------|-------------|
| `jira_notify_lark` | Send card message to Lark group |
| `jira_ai_analyze` | AI-powered ticket analysis |
| `jira_ai_sprint_report` | AI sprint report (optional Lark send) |

## Architecture

```
cmd/server/         - Entry point (uber-go/fx DI)
config/             - Env-based configuration
domain/jira/        - Jira domain interfaces
domain/memory/      - PM memory domain models
internal/jira/      - felixgeelhaar/jirasdk adapter
internal/ai/        - OpenAI-compatible client
internal/lark/      - larksuite/oapi-sdk-go + webhook fallback
internal/memory/    - SQLite store (WAL mode)
internal/bootstrap/ - DI wiring
application/tools/  - MCP tool handlers
transport/          - MCP server + tool registration
research/           - 508 academic papers on Scrum effectiveness
```

## Stack

- Go 1.26
- [jirasdk](https://github.com/felixgeelhaar/jirasdk) — Jira Cloud SDK
- [oapi-sdk-go/v3](https://github.com/larksuite/oapi-sdk-go) — Lark SDK
- [mcp-go](https://github.com/mark3labs/mcp-go) — MCP protocol
- [uber-go/fx](https://github.com/uber-go/fx) — Dependency injection
- SQLite — PM memory persistence

## MCP Config (OpenCode/Zara)

```json
{
  "Jira PM": {
    "type": "local",
    "command": ["/path/to/zara-jira-mcp-wrapper.sh"],
    "timeout": 30000,
    "enabled": true
  }
}
```

## Research

See `research/scrum-master-papers.md` — 508 papers on Scrum Master effectiveness, organized by 44 categories with key findings distilled.
