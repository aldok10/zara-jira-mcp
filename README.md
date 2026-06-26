# zara-jira-mcp

AI-powered Scrum Master MCP server with persistent memory. **124 tools** for Jira operations, sprint intelligence, risk management, team health, forecasting, coaching, and multi-channel notifications.

Not just a Jira wrapper — a complete PM/Scrum Master brain that **remembers**, **learns**, and **recommends**.

## Quick Start

```bash
cp .env.example .env
# Edit .env with your Jira, AI, and notification credentials

make build
make install  # copies to ~/.local/bin/
```

### MCP Configuration

```json
{
  "mcpServers": {
    "jira-pm": {
      "command": ["zara-jira-mcp"],
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

## Tool Categories (124 tools)

### Jira Operations (45 tools)

Full Jira Cloud CRUD: search, issues, sprints, epics, bulk operations, worklogs, issue links, watchers, projects, transitions, subtasks, labels, raw API access.

### PM Memory (20 tools)

Persistent SQLite memory that survives sessions:

| Area | Tools |
|------|-------|
| Sprint Tracking | `pm_snapshot_sprint`, `pm_track_daily`, `pm_burndown` |
| Risks | `pm_record_risk`, `pm_update_risk`, `pm_risk_dashboard`, `pm_auto_detect_risks` |
| Decisions | `pm_record_decision`, `pm_search_decisions`, `pm_record_learning` |
| Blockers | `pm_record_blocker`, `pm_resolve_blocker`, `pm_blockers` |
| Team | `pm_record_team_metric`, `pm_team_health`, `pm_confidence` |
| Retros | `pm_record_retro`, `pm_action_items` |
| Dependencies | `pm_record_dependency`, `pm_resolve_dependency`, `pm_dependencies` |

### AI Intelligence (15 tools)

All powered by historical memory + live Jira data:

| Tool | What It Does |
|------|-------------|
| `pm_recommendations` | AI recommendations from ALL historical context |
| `pm_standup_prep` | Daily talking points (live + memory) |
| `pm_forecast` | Monte Carlo "when will it be done?" (10,000 simulations) |
| `pm_anti_patterns` | Detect: zombie sprints, hero culture, scope creep, dead retros |
| `pm_coaching` | Data-driven coaching advice by topic |
| `pm_facilitate` | Fresh ceremony facilitation (standup/retro/planning/grooming/review) |
| `pm_retro_analysis` | Pattern detection across retrospectives |
| `pm_goal_check` | AI evaluates sprint goal progress |
| `pm_check_ready` | Story readiness (INVEST + DoR) |
| `pm_exec_report` | Executive stakeholder report (no jargon) |
| `pm_flow_metrics` | WIP, throughput, cycle time, lead time |
| `pm_sprint_compare` | This sprint vs last |
| `pm_weekly_digest` | AI weekly activity summary |
| `pm_nl_to_jql` | Natural language to JQL conversion |
| `jira_ai_analyze` | Ad-hoc AI analysis of any tickets |

### Process & Health (18 tools)

| Tool | What It Does |
|------|-------------|
| `pm_sprint_health` | 0-100 health score (velocity + blockers + scope + team) |
| `pm_health_history` | Health trend over time |
| `pm_scorecard` | End-of-sprint grade (A-F) |
| `pm_velocity_trend` | Velocity over sprints + trend detection |
| `pm_capacity_plan` | Data-driven capacity recommendation |
| `pm_sprint_goals` | Set/track/close sprint goals |
| `pm_dod` / `pm_dor` | Definition of Done / Ready (entry + exit gates) |
| `pm_agreements` | Team working agreements |
| `pm_experiment` / `pm_experiments` | Improvement experiments from retros |
| `pm_planning_prep` | Complete sprint planning preparation package |
| `pm_dashboard` | One-shot full PM view |
| `pm_scope_creep` | Mid-sprint scope change detection |
| `pm_backlog_groom` | Find stale backlog items |

### Escalation & Reporting (8 tools)

| Tool | What It Does |
|------|-------------|
| `pm_escalate` | Auto-escalate critical items to notification channels |
| `pm_escalations` | Escalation history |
| `pm_release_notes` | Categorized release notes from done issues |
| `pm_exec_report` | Executive summary (business outcomes, not story points) |
| `pm_weekly_digest` | AI weekly team digest |
| `pm_team_kb` | Onboarding knowledge base + AI Q&A |
| `jira_ai_sprint_report` | AI sprint report |
| `jira_notify_lark` | Send to Lark |

### Workflow Recipes (3 tools)

One-click workflows:
- `pm_recipe_start_work` — Assign + transition + suggest branch name
- `pm_recipe_done` — Transition + log time + comment
- `pm_recipe_block` — Record blocker + comment on issue

### Notifications (9 tools)

Multi-channel: Lark, Slack, Discord, Telegram, Teams, Email, Confluence, broadcast.

## Architecture

```
cmd/server/          Entry point (uber-go/fx DI)
config/              Env-based configuration
domain/
  jira/              Jira domain interfaces + models
  memory/            PM memory domain (14 entities)
  ai/                AI provider interface
  lark/              Lark notifier interface
internal/
  jira/              jirasdk adapter (Jira Cloud REST)
  ai/                OpenAI-compatible client
  lark/              Lark SDK + webhook
  memory/            SQLite store (WAL mode, 14 tables)
  slack/             Slack API client
  bootstrap/         DI wiring
application/tools/   MCP tool handlers (12 files)
transport/           MCP server + tool registration
```

### SQLite Memory (14 tables)

```
sprint_snapshots    — velocity, completion, carryover per sprint
daily_progress      — burndown data points
risks               — risk register with severity + mitigation
decisions           — decision log with rationale + tags
blockers            — impediments with resolution time
team_metrics        — per-member per-sprint stats
retrospectives      — went well / improve / actions
action_items        — retro follow-ups
dependencies        — cross-issue/team dependency map
meeting_notes       — ceremony outcomes
health_scores       — computed health over time
sprint_goals        — goal + key results + outcome
dod_items           — DoD + DoR checklists
escalations         — escalation audit trail
```

## Stack

- Go 1.26
- [mcp-go](https://github.com/mark3labs/mcp-go) — MCP protocol implementation
- [jirasdk](https://github.com/felixgeelhaar/jirasdk) — Jira Cloud REST SDK
- [oapi-sdk-go/v3](https://github.com/larksuite/oapi-sdk-go) — Lark SDK
- [uber-go/fx](https://github.com/uber-go/fx) — Dependency injection
- [go-sqlite3](https://github.com/mattn/go-sqlite3) — SQLite with WAL mode
- OpenAI-compatible API — AI analysis (any provider)

## Key Differentiators

1. **Persistent Memory** — Knows what happened last sprint. Tracks decisions, risks, blockers across sessions.
2. **Monte Carlo Forecasting** — "When will it be done?" with probability ranges, not guesses.
3. **Anti-Pattern Detection** — Automatically detects zombie sprints, hero culture, scope creep, dead retros.
4. **AI Coaching** — Data-driven coaching suggestions, not generic textbook advice.
5. **Multi-Audience Reporting** — Executive reports (business outcomes) vs team reports (sprint data).
6. **Proactive Alerting** — Auto-escalate critical risks and chronic blockers.
7. **Ceremony Facilitation** — Fresh retro formats, planning checklists, standup prompts.
8. **Flow Over Velocity** — WIP, throughput, cycle time. Predicts bottlenecks better than velocity.
9. **Process Maturity** — DoR, DoD, working agreements, improvement experiments.

## Environment Variables

```bash
# Required
JIRA_BASE_URL=https://company.atlassian.net
JIRA_EMAIL=you@company.com
JIRA_API_TOKEN=your-api-token

# AI (required for intelligence tools)
JIRA_AI_BASE_URL=https://api.openai.com
JIRA_AI_API_KEY=sk-...
JIRA_AI_MODEL=gpt-4o-mini

# Lark (optional)
JIRA_LARK_WEBHOOK_URL=https://open.larksuite.com/open-apis/bot/v2/hook/xxx

# Slack (optional)
SLACK_BOT_TOKEN=xoxb-...
SLACK_DEFAULT_CHANNEL=general

# Memory (optional, default: ~/.zara-jira-mcp/pm_memory.db)
PM_MEMORY_DB_PATH=/custom/path/pm_memory.db
```

## Development

```bash
make build      # Build binary
make test       # Run tests
make lint       # Run linter
make install    # Install to ~/.local/bin/
```

## License

MIT
