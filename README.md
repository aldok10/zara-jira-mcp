# zara-jira-mcp

**89-tool MCP server** — AI-powered Scrum Master with persistent memory, Jira Cloud integration, multi-channel notifications, and 12 platform integrations.

[![CI](https://github.com/aldok10/zara-jira-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/aldok10/zara-jira-mcp/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![MCP Compatible](https://img.shields.io/badge/MCP-compatible-blue)](https://modelcontextprotocol.io)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://golang.org)

Stop copy-pasting status updates. Start coaching your team with data.

---

## Quick Start

```bash
git clone https://github.com/aldok10/zara-jira-mcp.git && cd zara-jira-mcp
cp .env.example .env   # Add JIRA_BASE_URL, JIRA_EMAIL, JIRA_API_TOKEN
make build && make install
```

Add to your MCP client (OpenCode, Claude, Cursor, etc.):

```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "~/.local/bin/zara-jira-mcp",
      "env": {
        "JIRA_BASE_URL": "https://company.atlassian.net",
        "JIRA_EMAIL": "you@company.com",
        "JIRA_API_TOKEN": "your-token",
        "PM_MEMORY_DB_PATH": "~/.zara-jira-mcp/pm_memory.db",
        "JIRA_AI_BASE_URL": "https://api.openai.com/v1",
        "JIRA_AI_API_KEY": "your-key"
      }
    }
  }
}
```

**First commands:**
```
jira_boards                  → discover your boards
jira_board_config board_id=X → see column layout & status mappings
pm board_id=X                → sprint status + blockers + risks
```

---

## What You Can Do With 89 Tools

| Category | Tools | What |
|----------|-------|------|
| **Jira** | 28 | Search, create, transition, assign, link, comment, worklog, epics, versions |
| **Sprint/PM** | 18 | Snapshots, decisions, risks, blockers, retros, meetings, forecast |
| **Notifications** | 5 | Lark, Slack, Discord, Telegram + smart routing |
| **GitHub** | 10 | PRs, issues, milestones, releases, activity, branches |
| **GitLab** | 9 | MRs, issues, milestones, branches, files |
| **Calendar** | 3 | Create/list events, schedule meetings (Lark) |
| **Time** | 2 | Clockify entries & reports |
| **Wiki** | 6 | Confluence + Notion search, read, create |
| **Incidents** | 2 | PagerDuty incidents & on-call |
| **Linear** | 3 | Issues, cycles, activity |
| **Sheets** | 1 | Google Sheets read |
| **Tools** | 2 | `pm_backup` (JSON export), `pm_onboard` (config wizard) |

### Quick Examples

```
# What's happening right now?
pm board_id=84               → sprint status: 12/20 done, 3 blocked
pm_blockers                  → active impediments
pm_risks                     → risk dashboard

# Record things (memory keeps it across sessions)
pm_decide what="Use PostgreSQL" who="team" why="Operational experience"
pm_record_risk title="API performance" severity=high owner="Febrian"
pm_record_blocker description="Waiting for MixPanel access" issue_key="PROJ-123"

# Sprint ceremonies
pm_record_retro sprint_name="Sprint 12" went_well="Sisa waktu untuk improvement"
pm_record_meeting meeting_type="planning" notes="Sprint goal: stabilize July"

# Forecast & planning
pm_forecast board_id=84      → Monte Carlo: 50% in Sprint 13, 85% in Sprint 14
pm_snapshot board_id=84      → save sprint state to memory
pm_health board_id=84        → health trend across sprints

# Dev workflow bridge
pm_github_prs                → open PRs with age & reviewers
pm_github_activity days=7   → commits, merges, issues closed
pm_github_create_issue title="Add login" labels=backend

# Notifications
jira_notify_lark content="Sprint 12 done. 12/20 completed. 3 blockers."
notify_routed content="Prod incident" severity=critical audience=team

# Board-aware classification
jira_board_config board_id=84 → see which statuses map to done/blocked/progress

# Memory management
pm_backup                    → full JSON export of all PM memory
pm_reconcile                 → sync stored blockers/risks with current Jira state
```

---

## How It Works

The server runs as an MCP stdio process alongside your AI editor. Every read (`jira_search`, `jira_get_issue`, `jira_sprint_summary`) automatically:

1. **Records blockers** → detects blocked statuses, stores in SQLite
2. **Records stale risks** → flags Highest/Critical issues untouched >7 days
3. **Records sprint snapshots** → saves done/in-progress/blocked/todo counts
4. **Reconciles** → checks stored items against current state, auto-resolves

Board-aware classification uses each board's column configuration for accurate status mapping — no more hardcoded string matching.

All memory persists in SQLite WAL at `~/.zara-jira-mcp/pm_memory.db`.

---

## Architecture

Single Go binary. No runtime deps. Starts in <1 second.

```
apps/api/                  # Entry & DI wiring
├── cmd/server/main.go     # Entry point
└── internal/
    ├── bootstrap/         # Manual DI (all 89 tools wired here)
    └── mcp/               # 13 registration files

modules/                   # Hexagonal domain modules
├── jira/                  # 28 tools — full Jira Cloud CRUD
├── sprint/                # 18 tools — PM memory + analysis
└── notification/          # 5 tools — multi-channel send

shared/                    # Shared kernel
├── domain/                # Cross-module entities
├── infrastructure/        # 16+ clients (AI, GitHub, GitLab, Notion...)
└── usecase/               # Shared business logic

agents/                    # Agent architecture (event-driven)
```

**Stack:** Go 1.26.4 | `mark3labs/mcp-go` v0.55.1 | SQLite WAL | Manual DI

**Config:** Environment variables only (`.env` or shell). See `.env.example`.

---

## Integrations (All Optional)

| Category | Services | Tools |
|----------|----------|-------|
| Project Tracking | Jira Cloud | 28 |
| Sprint/PM | SQLite memory | 18 |
| Code | GitHub, GitLab | 10 + 9 |
| Communication | Lark, Slack, Discord, Telegram, Teams, Email | 5 |
| Wiki | Confluence, Notion | 3 + 3 |
| Incidents | PagerDuty | 2 |
| Time | Clockify | 2 |
| Calendar | Lark Calendar | 3 |
| Spreadsheets | Google Sheets | 1 |
| AI | OpenAI, Anthropic, Gemini, Groq, Ollama, OpenRouter, DeepSeek | (internal) |

Works with just Jira + SQLite. Everything else is optional.

---

## Full Tool Reference

See [SKILL.md](SKILL.md) for the complete 89-tool reference with all parameters, descriptions, and workflow patterns.

---

## Research Foundation

Built on validated research, not hype:

- **DORA 2024**: AI increases individual productivity but hurts stability without fundamentals
- **State of Agile (18th Ed)**: AI is the "Fourth Wave" of software delivery
- **Little's Law**: WIP limits + cycle time = faster delivery
- **Monte Carlo Forecasting**: Historical throughput beats gut-feel estimates
- **Tuckman's Stages**: Team maturity determines SM stance
- **McKinsey 2026**: AI-augmented SM/PM roles are the fastest-growing function

---

## Build & Develop

```bash
make build      # Build modular → bin/zara-jira-mcp
make test       # Run all tests
make lint       # golangci-lint
make install    # Copy to ~/.local/bin/

# Test individual tools
echo '{"jsonrpc":"2.0","id":1,"method":"tools/list","params":{}}' | ./bin/zara-jira-mcp
```

---

## Author

**[Aldo Karendra](https://github.com/aldok10)** — Building AI tools that make engineers and Scrum Masters more effective.

## License

MIT — use it, fork it, make it yours.
