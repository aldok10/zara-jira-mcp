# AGENTS.md

Instructions for AI agents working with this project.

## What This Is

`zara-jira-mcp` is an **MCP server in active migration** from a monolithic ~279-tool server to a **modular hexagonal architecture** with clean domain separation.

| Aspect | Description |
|--------|-------------|
| **Purpose** | AI-powered Scrum Master with persistent memory, Jira Cloud integration, multi-channel notifications |
| **Stack** | Go 1.26.4, `github.com/mark3labs/mcp-go` v0.55.1, SQLite (WAL) via `mattn/go-sqlite3` |
| **Modules** | `jira`, `sprint`, `notification` — hexagonal (ports & adapters) |
| **Binary** | Single binary, starts in <1 second. Current installed binary (~18MB) still has the monolithic ~279 tools |
| **Version** | v0.4.0 (modular), v0.4.0 (monolithic binary) |
| **Transport** | MCP stdio |
| **Module path** | `github.com/aldok10/zara-jira-mcp` |

> **Migration status**: The source code has been restructured to a modular architecture. The installed binary at `~/.local/bin/zara-jira-mcp` still runs the older monolithic tools (~279). The new modular code in `apps/api/` currently registers **99 tools** (Jira: 25, Sprint PM: 31, Notification: 5, GitHub: 10, GitLab: 9, PagerDuty: 2, Confluence: 3, Linear: 3, Notion: 3, Calendar: 3, Clockify: 2, Sheets: 1, Backup: 1, Onboard: 1). Event bus (InMemoryBus with retry, depth tracking, metrics) and Agent dispatcher layer wired into bootstrap. All shared infrastructure clients are now wired.

## Architecture

```
apps/api/                    # Application entry points
├── cmd/server/main.go       # Modular server entry (new)
├── internal/
│   ├── bootstrap/           # DI wiring (manual, no framework)
│   └── mcp/                 # MCP tool registration (jira.go, sprint.go)

modules/                     # Domain modules (hexagonal)
├── jira/                    # Jira Cloud operations
│   ├── domain/              # Entities: Issue, Board, Project, User, Sprint
│   ├── application/         # Use cases + ports
│   ├── infrastructure/      # REST client (client/, rest_*.go)
│   └── interfaces/          # MCP handlers
├── sprint/                  # Sprint PM intelligence
│   ├── domain/              # Velocity, Score, Predictability, Memory, Planning
│   ├── application/         # Services: Sprint, Analysis, OKR, People, Jira
│   ├── infrastructure/      # SQLite persistence, sprint store
│   └── interfaces/          # MCP handlers
└── notification/            # Multi-channel notifications
    ├── domain/              # Notifier interface, Lark config
    ├── infrastructure/      # Slack, Discord, Telegram, Teams, Email, Lark
    └── interfaces/          # MCP handlers

shared/                      # Shared kernel across modules
├── domain/                  # Agent/*, AI, Event, Jira, Lark, Memory, Planning, Sprint, Team
├── infrastructure/          # Adapters, AI, Bus, Calendar, Clockify, Config, Confluence,
│                            # GitHub, GitLab, HTTP client, Lark, Linear, Logger, MCP util,
│                            # Notion, Observability, PagerDuty, Sheets, Validate
└── usecase/                 # Analysis, OKR, People, Sprint services

agents/                      # AI Agent Architecture layer
├── agent.go                 # Agent interface, Dispatcher, Planner, Coordinator, Executor
```

### Core Modules

| Module | Domain | Infrastructure | Tools Registered |
|--------|--------|----------------|------------------|
| **`jira`** | Issue, Board, Project, User, Sprint, Commands, Types | REST client (search, issue, sprint, admin, misc) | 25 tools — full Jira CRUD |
| **`sprint`** | Velocity, Score, Predictability, Planning, Memory (16 tables), Team/Workload | SQLite persistence, AI provider, Jira adapter | 17 tools — PM memory + analysis |
| **`notification`** | Notifier contract, Lark entities | Slack, Discord, Telegram, Teams, Email, Lark (bot + webhook) | 5 tools — multi-channel routing |
| **GitHub** (shared) | PRs, issues, milestones, releases, activity, branches | REST client (`shared/infrastructure/github/`) | 10 tools — dev workflow bridge |
| **GitLab** (shared) | MRs, issues, milestones, branches, files | REST client (`shared/infrastructure/gitlab/`) | 9 tools — dev workflow bridge |
| **PagerDuty** (shared) | Incidents, on-call schedules | REST client (`shared/infrastructure/pagerduty/`) | 2 tools — incident awareness |
| **Confluence** (shared) | Pages (search/get/create) | REST client (`shared/infrastructure/confluence/`) | 3 tools — wiki bridge |
| **Linear** (shared) | Issues, cycles, activity | GraphQL client (`shared/infrastructure/linear/`) | 3 tools — project tracking bridge |
| **Notion** (shared) | Pages (search/create), database queries | REST client (`shared/infrastructure/notion/`) | 3 tools — wiki bridge |

### Shared Infrastructure (all wired)

All clients in `shared/infrastructure/` are now wired into the MCP server:

- **AI**: OpenAI client (`shared/infrastructure/ai/`) — ✅ wired into sprint service
- **Calendar**: Lark Calendar (`shared/infrastructure/calendar/`) — ✅ 3 tools (`pm_calendar_create`, `pm_calendar_events`, `pm_calendar_schedule_meeting`)
- **Time tracking**: Clockify (`shared/infrastructure/clockify/`) — ✅ 2 tools (`pm_time_entries`, `pm_time_report`)
- **Docs**: Confluence (`shared/infrastructure/confluence/`) — ✅ 3 tools
- **Lark**: Bot, Webhook, OKR, Commands (`shared/infrastructure/lark/`)
- **Project tracking**: Linear (`shared/infrastructure/linear/`) — ✅ 3 tools
- **Logger**: Structured logging (`shared/infrastructure/logger/`)
- **MCP**: Utility helpers (`shared/infrastructure/mcputil/`) — ✅ used by all handlers
- **Wiki**: Notion (`shared/infrastructure/notion/`) — ✅ 3 tools
- **Observability**: Metrics/tracing (`shared/infrastructure/observability/`)
- **Sheets**: Google Sheets (`shared/infrastructure/sheets/`) — ✅ 1 tool (`pm_sheet_read`)
- **Validation**: Input validation (`shared/infrastructure/validate/`)

### Agent Architecture Layer

The `agents/` directory implements an **Agent Registry → Dispatcher → Planner → Coordinator → Executor** pattern:

- **Agent interface**: `Name()`, `Description()`, `EventTypes()`, `Execute()`
- **Dispatcher**: Routes system events to registered agents
- **Planner**: Decomposes goals into action sequences
- **Coordinator**: Orchestrates plan execution across tools
- **Executor**: Runs individual tool calls

Agents communicate via System Events — business domains never call agents directly.

## Data Model

16 SQLite tables, auto-created on first run via `modules/sprint/infrastructure/persistence/sqlite*.go`:

| Table | Purpose |
|-------|---------|
| `sprint_snapshots` | Sprint history (velocity, completion, carryover) |
| `daily_progress` | Burndown data |
| `risks` | Risk register + tech debt |
| `decisions` | Decision log + agreements + experiments + learnings |
| `blockers` | Impediment tracker |
| `team_metrics` | Individual sprint metrics |
| `retrospectives` | Retro outcomes |
| `action_items` | Retro follow-ups |
| `dependencies` | Dependency map |
| `meeting_notes` | Ceremony outcomes |
| `health_scores` | Health over time |
| `sprint_goals` | Goal tracking |
| `dod_items` | DoD + DoR checklists |
| `escalations` | Escalation audit trail |
| `coaching` | Coaching records |
| `okrs` | OKR tracking |

Additional advanced tables: `kb_articles`, `feedback_log`, `experiments`, `learnings`, `pulse_surveys`, `radar_dimensions`, `safety_surveys`, `stakeholder_pulse`, `kpi_definitions`, `kpi_measurements`, `key_results`, `kr_issues`, `okr_link`.

Persistence files: `sqlite.go`, `sqlite_deep.go`, `sqlite_advanced.go`, `sqlite_coaching.go`, `sqlite_okr.go`.

## Currently Registered Tools

These are the tools wired in the **new modular** codebase. The installed monolithic binary has ~279 tools.

### Jira Module (25 tools)
Full Jira CRUD: search, get issue, create/update/delete, transitions, assign, find user, add comment, sprints (list/start/move/close), create subtask, link issues, worklog, watchers, labels, projects, versions (create/release), components, attachments, epics (add/remove/list), boards, fields.

### Sprint/PM Module (17 tools)
PM memory + analysis: snapshot, sprint health, scorecard, predictability, velocity trend, flow metrics, forecast, calendar, risks, decisions, blockers, dependencies, meetings, retros, action items, next action, context notes.

### Notification Module (5 tools)
Multi-channel routing: Lark, Slack, Discord, Telegram, routed notification.

### GitHub Module (10 tools)
Dev workflow bridge: PRs (list), releases, activity summary, PR metrics, branch search, PRs by branch, create issue, list issues, milestones (create/list).

### GitLab Module (9 tools)
Dev workflow bridge: issues (list/create), merge requests, milestones (list/create), branch search, MRs by branch, read file, list files.

### PagerDuty Module (2 tools)
Incident awareness: list incidents, on-call schedule.
| Tool | Params |
|------|--------|
| `pm_incidents` | `status` |
| `pm_oncall` | none |

### Confluence Module (3 tools)
Wiki bridge: pages (search/get/create).
| Tool | Params |
|------|--------|
| `pm_confluence_search` | `query`, `limit` |
| `pm_confluence_get_page` | `page_id` |
| `pm_confluence_create_page` | `space_key`, `title`, `body`, `parent_id` |

### Linear Module (3 tools)
Project tracking bridge: issues, cycles, activity.
| Tool | Params |
|------|--------|
| `pm_linear_issues` | `team`, `state` |
| `pm_linear_cycles` | none |
| `pm_linear_activity` | none |

### Notion Module (3 tools)
Wiki bridge: pages (search/create), database queries.
| Tool | Params |
|------|--------|
| `pm_notion_search` | `query`, `limit` |
| `pm_notion_create_page` | `title`, `content`, `parent_id` |
| `pm_notion_query_db` | `database_id`, `filter`, `limit` |

> **Tool gap**: The monolithic installed binary has ~279 tools across 17 categories. The modular codebase currently registers 79 tools. To restore full functionality in the modular version, register remaining tools from the old `tools/handlers.go` pattern into new module interfaces.

## Key Design Decisions

### Why Hexagonal Architecture
- **Domain isolation**: `modules/jira`, `sprint`, `notification` each own their domain, ports, and infrastructure
- **Portability**: Swap Jira SDK, notification provider, or DB without touching domain
- **Testability**: Domain logic depends only on interfaces (ports)
- **Migration**: Old code in `internal/` moved to `shared/` — preserved but deprecated

### Module Structure Convention
```
module/
├── domain/           # Entities, value objects, aggregates, domain services
├── application/      # Use cases, ports (inbound/outbound interfaces)
├── infrastructure/   # Adapters, REST clients, persistence, external SDKs
├── interfaces/       # MCP handlers, CLI, or other entry points
└── test/             # Integration/fixture data
```

### No DI Framework
The modular bootstrap (`apps/api/internal/bootstrap/`) uses **manual dependency injection**. No uber-go/fx in the new code. Each handler is constructed explicitly:

```go
restClient := client.NewRestClient(cfg)
jiraSvc := service.NewJiraService(restClient)
jiraHandler := jira_mcp.NewHandlers(jiraSvc)
```

### Tool Registration Pattern
```go
// apps/api/internal/mcp/jira.go
func RegisterJiraTools(s *server.MCPServer, h *jmcp.Handlers) {
    s.AddTool(
        mcp.NewTool("jira_search",
            mcp.WithDescription("Search Jira issues using JQL."),
            mcp.WithString("jql", mcp.Required(), ...),
        ),
        h.SearchIssues,
    )
}
```

## Important Notes for AI Agents

1. **Two codebases exist**: Source code is modular (new), installed binary is monolithic (old with ~279 tools).
2. **Source code is truth**: When making changes, work in the modular structure (`apps/api/`, `modules/`, `shared/`).
3. **First interaction**: Call `jira_boards` to get board_id before using PM tools.
4. **Tool limitation**: The modular codebase has 85 tools registered. If you need tools not listed above, they're still available in the installed monolithic binary (~279 tools).
5. **Memory persistence**: SQLite WAL at `~/.zara-jira-mcp/pm.db` — shared between both builds.
6. **To build modular version**: `make build` (produces `bin/zara-jira-mcp`).
7. **Migration priority**: Wire remaining tools from monolithic binary (~279) into modular structure.
8. **Agent layer**: `agents/` directory implements event-driven agent architecture — now wired into bootstrap with BusBridge to event bus.

## Quick Build

```bash
make build        # Build modular version → bin/zara-jira-mcp
make install      # Copy to ~/.local/bin/
make test         # Run modular tests
make lint         # golangci-lint
```

## Project Status

Current state:
- ✅ **Modular architecture** with 3 hex modules + shared kernel + agent layer
- ✅ **SQLite persistence** with 16+ tables for full PM memory
- ✅ **Multi-channel notification infrastructure** (Slack, Discord, Telegram, Teams, Email, Lark)
- ✅ **Shared infrastructure** for 16+ external services (ready to wire)
- ✅ **Agent architecture** with dispatcher/planner/coordinator pattern
- ✅ **Jira client wired to sprint service** (duck-typing via domain.Client)
- ✅ **AI provider wired** (shared/infrastructure/ai.OpenAIClient → sprint port.AIProvider)
- ✅ **86 tools registered** in modular code (was 9+2)
- ✅ **All shared infrastructure clients wired** — Calendar, Clockify, Sheets now have MCP handlers
- ✅ **Event bus + Agent dispatcher** wired into bootstrap (InMemoryBus → BusBridge → Dispatcher)
- ✅ **pm_forecast** added (Monte Carlo simulation, already implemented in sprint service)
- ⏳ **Migration in progress**: wire remaining tools from monolithic binary (~279) into modular structure

The project has evolved from a monolithic 279-tool server (still in installed binary) to a modular, maintainable hexagonal architecture with clear separation of concerns and improved developer experience. The shared infrastructure and domain logic are largely intact — the remaining work is wiring them into the new module interfaces.
