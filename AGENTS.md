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

> **Migration status**: The source code has been restructured to a modular architecture. The installed binary at `~/.local/bin/zara-jira-mcp` still runs the older monolithic tools (~279). The new modular code in `apps/api/` currently registers **9 core tools** with all domain logic ready for expansion.

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
| **`jira`** | Issue, Board, Project, User, Sprint, Commands, Types | REST client (search, issue, sprint, admin, misc) | `jira_search`, `jira_get_issue`, `jira_boards`, `jira_sprint_summary` |
| **`sprint`** | Velocity, Score, Predictability, Planning, Memory (16 tables), Team/Workload | SQLite persistence, Snapshot adapter | `pm`, `pm_create`, `pm_decide`, `pm_risk`, `pm_next` |
| **`notification`** | Notifier contract, Lark entities | Slack, Discord, Telegram, Teams, Email, Lark (bot + webhook) | `jira_notify_lark`, `notify_routed` (stubs) |

### Shared Infrastructure (ready but not wired)

The `shared/infrastructure/` directory contains client implementations for **16+ external services** that were part of the monolithic version. These need to be wired into modular handlers:

- **AI**: OpenAI client (`shared/infrastructure/ai/`)
- **Calendar**: Lark Calendar (`shared/infrastructure/calendar/`)
- **Time tracking**: Clockify (`shared/infrastructure/clockify/`)
- **Docs**: Confluence (`shared/infrastructure/confluence/`)
- **Source control**: GitHub, GitLab (`shared/infrastructure/github/`, `gitlab/`)
- **Lark**: Bot, Webhook, OKR, Commands (`shared/infrastructure/lark/`)
- **Project tracking**: Linear (`shared/infrastructure/linear/`)
- **Logger**: Structured logging (`shared/infrastructure/logger/`)
- **MCP**: Utility helpers (`shared/infrastructure/mcputil/`)
- **Wiki**: Notion (`shared/infrastructure/notion/`)
- **Observability**: Metrics/tracing (`shared/infrastructure/observability/`)
- **Incidents**: PagerDuty (`shared/infrastructure/pagerduty/`)
- **Sheets**: Google Sheets (`shared/infrastructure/sheets/`)
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

### Jira Module (4 tools)
| Tool | Params | Returns |
|------|--------|---------|
| `jira_search` | `jql` (required), `max_results` | Issues: key, summary, status, priority, assignee |
| `jira_get_issue` | `key` | Full issue details |
| `jira_boards` | none | Board IDs |
| `jira_sprint_summary` | `board_id` | Active sprint status breakdown |

### Sprint/PM Module (5 tools)
| Tool | Params | Returns |
|------|--------|---------|
| `pm` | `board_id` | Quick project status |
| `pm_create` | `title`, `description`, `project`, `type`, `assignee`, `priority`, `labels`, `platform` | Create work item |
| `pm_decide` | `what`, `who`, `why` | Record a decision |
| `pm_risk` | `what`, `severity`, `owner` | Record a risk |
| `pm_next` | `board_id` | AI-suggested next action |

### Notification Module (2 stubs)
| Tool | Params |
|------|--------|
| `jira_notify_lark` | `content`, `title` |
| `notify_routed` | `content`, `severity`, `audience`, `title` |

> **Tool gap**: The monolithic installed binary has ~279 tools across 17 categories. The modular codebase currently registers 9+2 tools. To restore full functionality in the modular version, register tools from the old `tools/handlers.go` pattern into the new module interfaces.

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
4. **Tool limitation**: The modular codebase has only 9+2 tools registered. If you need tools not listed above, they're still available in the installed monolithic binary.
5. **Memory persistence**: SQLite WAL at `~/.zara-jira-mcp/pm.db` — shared between both builds.
6. **To build modular version**: `make build` (produces `bin/zara-jira-mcp`).
7. **Migration priority**: Wire more tools from `shared/infrastructure/` into module interfaces.
8. **Agent layer**: New `agents/` directory implements event-driven agent architecture — not yet wired into bootstrap.

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
- ⚠️ **9+2 tools registered** in modular code (vs ~279 in monolithic binary)
- ⚠️ **Infrastructure clients exist** in `shared/infrastructure/` but not wired into modular handlers
- ⚠️ **Agent layer** not yet wired into bootstrap
- ⏳ **Migration in progress**: wire tools, register modules, replace monolithic binary

The project has evolved from a monolithic 279-tool server (still in installed binary) to a modular, maintainable hexagonal architecture with clear separation of concerns and improved developer experience. The shared infrastructure and domain logic are largely intact — the remaining work is wiring them into the new module interfaces.
