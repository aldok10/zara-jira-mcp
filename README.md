# zara-jira-mcp

Your AI-powered Scrum Master that actually remembers what happened last sprint.

**124 tools.** Jira operations, sprint intelligence, risk management, team health tracking, Monte Carlo forecasting, coaching, and multi-channel notifications. All in one MCP server.

Built as part of the [Zara Agent OPC](https://github.com/aldok10/zara-agent-opc) ecosystem by [Aldo Karendra](https://www.linkedin.com/in/aldok10/).

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![MCP Compatible](https://img.shields.io/badge/MCP-compatible-blue)](https://modelcontextprotocol.io)
[![Donate](https://img.shields.io/badge/Support-SociaBuzz-orange)](https://sociabuzz.com/aldok10)

---

## The Problem

PM tools are dumb. Jira gives you data, but it doesn't think. It doesn't notice that your team has been carrying the same 3 tickets for 4 sprints. It won't tell you that one person is doing 60% of the work. It doesn't remember the decision you made 3 sprints ago about why you chose microservices.

Scrum Masters spend hours every week pulling reports, writing summaries, preparing ceremonies, chasing blockers. Manual. Repetitive. Time that should go to actually helping the team.

## What This Solves

`zara-jira-mcp` is not another Jira wrapper. It's a PM brain with memory.

- **It remembers.** Decisions, risks, retro outcomes, blocker patterns. Across sprints, across months.
- **It thinks.** Monte Carlo forecasting tells you *when* you'll ship, not just *what's done*. Anti-pattern detection catches problems before they become crises.
- **It acts.** Auto-escalate stale blockers. Generate standup briefs. Prep your sprint planning. Write executive reports that don't sound like a dashboard export.

One tool call replaces 30 minutes of clicking through Jira boards.

## Who Is This For?

**Scrum Masters** who are tired of being glorified note-takers. You should be coaching teams and removing impediments, not copying ticket statuses into Slack.

**Engineering Managers** who need visibility without micromanaging. Sprint health scores, velocity trends, and risk dashboards without asking "how's the sprint going?" in standup.

**Product Managers** who want honest delivery forecasts. Not "we'll try to finish by Friday" but "there's an 85% probability this ships by next Thursday based on historical throughput."

**Solo devs and small teams** who don't have a dedicated PM but still want structured sprint hygiene.

## Why Now?

AI coding assistants are everywhere. But the PM side of software delivery is still stuck in 2015. You have Copilot writing your code, but your Scrum Master is still manually copying ticket updates into a Google Doc.

MCP (Model Context Protocol) changes this. Your AI assistant can now talk directly to Jira, run forecasts, detect risks, and manage ceremonies, all through a natural conversation. No context switching, no dashboard fatigue.

This is what AI-augmented project management looks like. Not replacing humans, but giving them superpowers.

## Works With Every Major AI Tool

| Client | Type | Status |
|--------|------|--------|
| Claude Desktop | Desktop | Ready |
| Claude Code | CLI | Ready |
| ChatGPT Desktop | Desktop | Ready |
| Cursor | IDE | Ready |
| Windsurf | IDE | Ready |
| VS Code + Copilot | IDE | Ready |
| Cline / Roo Code | VS Code ext | Ready |
| Zed | IDE | Ready |
| Gemini CLI | CLI | Ready |
| Goose (Block) | CLI/Desktop | Ready |
| Amazon Q Developer | CLI/IDE | Ready |
| OpenCode | CLI | Ready |
| Kiro | IDE | Ready |
| Codex CLI | CLI | Ready |
| Cherry Studio | Desktop | Ready |
| Jan | Desktop | Ready |
| Msty | Desktop | Ready |
| LibreChat | Self-hosted | Ready |
| TypingMind | Web | Ready |
| Copilot Studio | Enterprise | Ready |

Works with **any** MCP-compatible client. Stdio transport, zero external dependencies.

Pre-built config files included for all platforms. See [docs/agents/](docs/agents/) for copy-paste setup per client.

## Quick Start

```bash
cp .env.example .env
# Edit .env with your Jira, AI, and notification credentials

make build
make install  # copies to ~/.local/bin/
```

### MCP Configuration (Claude / ChatGPT / Cursor / Windsurf / Gemini CLI / Amazon Q)

```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "zara-jira-mcp",
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

### Zed (uses `context_servers`)

```json
{
  "context_servers": {
    "jira-pm": {
      "command": {
        "path": "zara-jira-mcp",
        "args": []
      }
    }
  }
}
```

### Goose

```yaml
extensions:
  jira-pm:
    type: stdio
    command: zara-jira-mcp
    enabled: true
```

### LibreChat

```yaml
mcpServers:
  jira-pm:
    type: stdio
    command: zara-jira-mcp
```

Full setup guides for all 20 clients: [docs/agents/](docs/agents/)

## What You Get (124 Tools)

### Jira Operations (45 tools)

Full Jira Cloud CRUD. Search, create, update, transition, bulk operations, worklogs, issue links, watchers, sprints, epics, subtasks, labels, raw API access.

No more switching between your AI assistant and the Jira UI.

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

Your sprint history doesn't disappear when you close the terminal.

### AI Intelligence (15 tools)

Powered by your historical data + live Jira state:

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

### Management & Stakeholder Reporting

| Tool | Audience | Purpose |
|------|----------|---------|
| `pm_exec_report` | VP / C-Level | Business outcomes, no jargon, 30-second read |
| `pm_weekly_digest` | All stakeholders | AI weekly summary: wins, concerns, next focus |
| `pm_release_notes` | Stakeholders | What shipped this sprint (features, fixes) |
| `pm_stakeholder_pulse` | PM internal | Track stakeholder satisfaction over time |
| `pm_stakeholder_trend` | PM internal | Is the relationship improving or degrading? |
| `pm_sm_impact` | SM's manager | Prove SM value: blockers resolved, risks mitigated |
| `pm_outcome_map` | PO / leadership | Connect sprints to OKR/business objectives |
| `pm_escalate` | Auto | Alert management when blockers/risks go chronic |
| `portfolio_summary` | Steering committee | AI executive summary across all projects |
| `pm_maturity_assessment` | Eng leadership | Team stage: Forming/Storming/Norming/Performing |

> Full guide: [docs/reporting-guide.md](docs/reporting-guide.md) — scenario-based guide for every type of stakeholder communication.

### Workflow Recipes (3 tools)

One-click workflows:
- `pm_recipe_start_work` — Assign + transition + suggest branch name
- `pm_recipe_done` — Transition + log time + comment
- `pm_recipe_block` — Record blocker + comment on issue

### Notifications (9 tools)

Multi-channel: Lark, Slack, Discord, Telegram, Teams, Email, Confluence, broadcast.

## Real-World Examples

**Monday morning standup prep:**
> "Prep my standup" -> gives you blockers, what moved yesterday, what's at risk, who might need help. 30 seconds instead of 10 minutes of clicking.

**Sprint planning:**
> "Prep planning for next sprint" -> last sprint outcome, carryover items, team capacity based on velocity history, risks to discuss, experiments to review.

**Mid-sprint health check:**
> "How's the sprint?" -> health score, scope creep detection, blocker age, forecast of completion probability.

**Executive update:**
> "Write exec report" -> business outcomes, delivery risks, team health. No story points, no Jira jargon. Ready to paste into an email.

**PO asks "will we hit the goal?":**
> "Check sprint goal progress" -> AI evaluates current data vs key results. Gives honest On Track / At Risk / Off Track verdict.

**Escalation to management:**
> "Show impediment aging" -> all blockers with age in days, which ones are chronic. Data to back up your escalation request.

**Cross-team dependency tracking:**
> "Show all open dependencies" -> who is waiting on whom, across teams. Bring this to your cross-team sync.

**Monthly steering committee:**
> "Portfolio summary" -> AI-generated status across all projects. Health, risks, delivery confidence. One-pager for the boardroom.

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

## Stack

- Go 1.26
- [mcp-go](https://github.com/mark3labs/mcp-go) — MCP protocol implementation
- [jirasdk](https://github.com/felixgeelhaar/jirasdk) — Jira Cloud REST SDK
- [oapi-sdk-go/v3](https://github.com/larksuite/oapi-sdk-go) — Lark SDK
- [uber-go/fx](https://github.com/uber-go/fx) — Dependency injection
- [go-sqlite3](https://github.com/mattn/go-sqlite3) — SQLite with WAL mode
- OpenAI-compatible API — AI analysis (any provider)

## What Makes This Different

1. **Memory that persists.** Your Jira dashboard resets every time you open it. This doesn't. It knows what happened 5 sprints ago and can spot patterns you'd miss.

2. **Forecasting that works.** Monte Carlo simulation with 10,000 runs. Not "we think maybe next week" but "there's a 70% chance by Thursday, 95% by the following Monday."

3. **Catches what humans miss.** Hero culture (one person doing everything). Zombie tickets (alive but not moving). Scope creep mid-sprint. Dead retro actions nobody follows up on.

4. **Speaks multiple languages.** Executive report for your VP. Sprint data for the team. Blocker alerts for engineering. Different audience, different format, same source of truth.

5. **Gets smarter over time.** The more sprints you track, the better the forecasts. The more retros you record, the better the pattern detection. It compounds.

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

## Part of the Zara Ecosystem

This project is built as part of [Zara Agent OPC](https://github.com/aldok10/zara-agent-opc), an empathetic AI engineering partner with cognitive memory, multi-agent orchestration, and self-improving capabilities. Zara is designed to be the AI companion that grows with you, not just another stateless tool.

If you find this useful, check out the full Zara ecosystem for AI-augmented development workflows.

## Author

**[Aldo Karendra](https://www.linkedin.com/in/aldok10/)** — Lead Backend Developer & AI Systems Architect based in Jakarta, Indonesia. Building AI tools that actually help engineers ship better software, faster.

- GitHub: [@aldok10](https://github.com/aldok10)
- LinkedIn: [linkedin.com/in/aldok10](https://www.linkedin.com/in/aldok10/)
- Support: [sociabuzz.com/aldok10](https://sociabuzz.com/aldok10)

## License

MIT
