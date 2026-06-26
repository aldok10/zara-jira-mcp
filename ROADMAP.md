# ROADMAP: zara-jira-mcp

> AI-powered Scrum Master MCP server — 139 tools for Jira intelligence,
> sprint management, team health tracking, and multi-platform notifications.

---

## Current State (v0.2.0)

**139 MCP tools** across 10 categories. All roadmap phases 0-6 substantially complete.

| Category | Tools | Status |
|----------|-------|--------|
| Jira Core CRUD | ~55 | Complete (exceeds go-jira parity) |
| PM Intelligence | ~25 | Complete (flow, forecast, coaching, anti-patterns) |
| PM Memory | ~20 | Complete (SQLite: sprints, risks, decisions, blockers, metrics, retros) |
| Notifications | ~15 | Complete (8 platforms + routing + broadcast + digest) |
| Portfolio | ~5 | Complete |
| Automation/Recipes | ~8 | Complete (start_work, done, block, planning prep) |
| Confluence | 3 | Complete |
| Version/Release | 4 | Complete |
| Meta/Health | ~4 | Complete |

### Platform Integrations

| Platform | Status | Auth |
|----------|--------|------|
| Jira Cloud | Full | API Token via jirasdk |
| Lark/Feishu | Full | oapi-sdk-go + webhook |
| Slack | Full | Bot token + webhook |
| Discord | Basic | Bot token |
| Telegram | Basic | Bot token |
| MS Teams | Basic | Incoming webhook |
| Email | Basic | SMTP |
| Confluence | Full | API Token |

### Quality

- Build: Clean (Go 1.26, golangci-lint)
- Tests: 27.4% coverage (handler tests with mocks)
- Lint: Zero violations
- Docker: Multi-stage build ready

---

## Completed Phases

### Phase 0: Foundation
Memory wired, structured errors, pagination, update tool, test infra, health tool.

### Phase 1: go-jira Parity (98%)
All core operations: CRUD, subtasks, epics, links, worklogs, labels, watchers,
projects, users, sprints (create/start/close/move), bulk ops, raw request,
versions, components, fields, attachments.

Missing only: clone_issue, createmeta (low value).

### Phase 2: PM Intelligence
Flow metrics, Monte Carlo forecast, velocity trend, capacity planning,
burndown, scope creep, anti-patterns, coaching advice, ceremony facilitator,
tech debt ratio, priority churn, NL-to-JQL, sprint comparison, confidence voting.

### Phase 3: Memory & Learning
SQLite persistence: sprint snapshots, risk register, decision log, blocker tracker,
team metrics, retrospectives, action items, learning records, dependencies.

### Phase 4: Automation
Standup prep, release notes, backlog groom, daily digest, planning prep,
sprint review prep, recipes (start_work, done, block), escalation.

### Phase 5: Portfolio
Overview, summary, blockers, risks, workload across projects.

### Phase 6: Integration (Exceeded)
8 platforms. Smart routing engine. Broadcast. Research-backed notification cadence.

---

## Remaining Gaps (Nice-to-Have)

| Item | Effort | Value |
|------|--------|-------|
| Vector embeddings for pattern matching | High | Medium |
| Calendar integration (Google/Outlook) | Medium | Medium |
| CI/CD event ingestion | Medium | Low |
| Webhook server (push notifications) | High | Medium |
| jira_clone_issue | Low | Low |
| jira_createmeta | Low | Low |
| Multi-tenant support | High | Low |

---

## Research Foundation

- `research/scrum-master-papers.md` — 508 academic papers on SM effectiveness
- `research/pm-integration-platforms.md` — Deep dive on notification routing, escalation patterns, anti-patterns

---

## Architecture

```
cmd/server/              Entry point (uber-go/fx DI)
config/                  Env-based config (all platforms optional)
domain/                  Interfaces + domain models
internal/jira/           felixgeelhaar/jirasdk adapter
internal/ai/             OpenAI-compatible client
internal/lark/           larksuite/oapi-sdk-go + webhook
internal/slack/          slack-go/slack
internal/discord/        bwmarrin/discordgo
internal/telegram/       go-telegram-bot-api
internal/teams/          MS Teams webhook
internal/email/          net/smtp
internal/confluence/     Confluence REST API
internal/memory/         SQLite (WAL mode)
internal/bootstrap/      DI wiring
application/tools/       MCP tool handlers (~15 files)
transport/               MCP server + tool registration (~10 files)
research/                Academic papers + integration research
```

## Success Metrics

| Metric | Target | Actual |
|--------|--------|--------|
| Tool count | 111 | **139** (125% of target) |
| go-jira parity | 100% | 98% |
| Platform integrations | 4 | **8** (200% of target) |
| Research papers | — | 508 |
| Test coverage | >70% | 27.4% (in progress) |
| Lint violations | 0 | 0 |
