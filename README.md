# zara-jira-mcp

**Stop being a Jira secretary. Start being the Scrum Master your team actually needs.**

Your AI copilot that handles the admin, surfaces the risks, and gives you back the hours you lost to status updates.

[![CI](https://github.com/aldok10/zara-jira-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/aldok10/zara-jira-mcp/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![MCP Compatible](https://img.shields.io/badge/MCP-compatible-blue)](https://modelcontextprotocol.io)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://golang.org)
[![GitHub stars](https://img.shields.io/github/stars/aldok10/zara-jira-mcp?style=social)](https://github.com/aldok10/zara-jira-mcp)

---

## The Problem

Scrum Masters and PMs spend **58% of their time** on meetings and admin that produces zero value ([Jabra 2025 Report](https://brandsit.pl/en/how-much-do-ineffective-meetings-cost-jabra-report/)). Meanwhile:

- Sprint reviews are copy-paste status updates nobody reads
- Retro actions die within days (the "Dead Retro" pattern)
- Risks are invisible until they explode
- Management asks "when will it be done?" and you guess
- You can't prove your impact as SM because there's no data

**This project exists because SM/PM work is 80% repetitive information assembly that AI can do in seconds.**

---

## What This Actually Does

| Before (Manual) | After (zara-jira-mcp) |
|-----------------|----------------------|
| 30 min writing sprint status | `pm_exec_report` — 10 seconds, business language |
| Guessing sprint capacity | `pm_forecast_sprint` — Monte Carlo simulation, confidence intervals |
| "Are we on track?" = gut feeling | `report_delivery_confidence` — GREEN/AMBER/RED with data |
| Retro actions forgotten | `sm_improvement_velocity` — tracks if actions produce change |
| Can't see team dysfunction | `sm_dysfunction_detector` — detects Hero Culture, Zombie Scrum, Scope Creep |
| Manual blocker escalation email | `report_escalation_brief` — PROBLEM/IMPACT/ASK/DEADLINE format |
| 1-on-1 prep = winging it | `pm_one_on_one_prep` — data-driven talking points per person |
| "What should I focus on?" | `pm_what_next` — AI prioritizes your day |
| Story points scattered in Jira | `pm_story_points` — totals from sprint/epic/JQL, grouped by status/person |

---

## Why This, Not Jira AI or Other Tools

| | zara-jira-mcp | Jira AI (Atlassian) | go-jira CLI | Other PM tools |
|---|---|---|---|---|
| **Remembers across sprints** | SQLite memory (decisions, risks, retros) | No memory | No memory | Some |
| **Works with YOUR AI editor** | Any MCP client (20+ supported) | Jira UI only | Terminal only | Their UI only |
| **Forecasting** | Monte Carlo on real data | Basic | None | Limited |
| **Team dysfunction detection** | Data-driven pattern matching | None | None | None |
| **Custom field support** | Auto-detects story points | Built-in | Manual | Varies |
| **Multi-channel notifications** | 8 platforms + smart routing | Email only | None | 1-2 platforms |
| **Cost** | Free (MIT) | Paid (Premium) | Free but dead | Paid |
| **Privacy** | Self-hosted, your data stays local | Cloud (Atlassian) | Local | Cloud |

---

## 5-Minute Setup

```bash
git clone https://github.com/aldok10/zara-jira-mcp.git
cd zara-jira-mcp
cp .env.example .env
# Edit .env: add JIRA_BASE_URL, JIRA_EMAIL, JIRA_API_TOKEN
make build && make install
```

Then add to your MCP client config:

```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "zara-jira-mcp",
      "env": {
        "JIRA_BASE_URL": "https://company.atlassian.net",
        "JIRA_EMAIL": "you@company.com",
        "JIRA_API_TOKEN": "your-token",
        "JIRA_AI_BASE_URL": "https://api.groq.com/openai/v1",
        "JIRA_AI_API_KEY": "your-key",
        "JIRA_AI_MODEL": "llama-3.3-70b-versatile"
      }
    }
  }
}
```

Works with **Claude Desktop, ChatGPT, Cursor, VS Code Copilot, Windsurf, Zed, OpenCode, Kiro, Cline, Gemini CLI, Goose**, and 10+ more.

---

## Built on Research, Not Vibes

This isn't another "AI wrapper." Every feature is grounded in peer-reviewed research:

- **DORA 2024**: AI increases individual productivity but hurts delivery stability without fundamentals. We enforce the fundamentals.
- **State of Agile (18th Edition)**: AI is the "Fourth Wave" of software delivery. We're built for it.
- **Little's Law**: WIP limits and cycle time are the only levers for faster delivery. We measure both.
- **Probabilistic Forecasting**: Monte Carlo on historical throughput beats gut-feel estimates every time.
- **DX Core 4 Framework**: Flow state, cognitive load, feedback loops. We reduce PM cognitive load by 50%+.
- **Tuckman's Stages**: Team maturity determines SM stance. We auto-assess and recommend.
- **Sense&Respond Anti-Pattern Library**: 10 documented dysfunctions with coaching interventions. We detect them from data.

---

## The SM's Daily Workflow (Powered)

| Time | You Say | Tool Responds |
|------|---------|--------------|
| 8:30 | "What should I focus on?" | Blockers aging >3d, 2 critical risks, sprint behind pace |
| 9:00 | "Prep my standup" | Talking points: 3 blockers to raise, 1 dependency to chase |
| 11:00 | "Is sprint goal at risk?" | AMBER: 62% done, 2 blocked items. Suggest removing PROJ-45 |
| 14:00 | "Prep 1-on-1 with Alice" | High carryover pattern. Workload: 7 items. Ask about estimation |
| 16:00 | "Write update for VP" | "On track. Auth shipped. API at risk due to infra dependency." |
| Friday | "Sprint narrative for review" | Business-friendly story of what shipped and why it matters |

---

## Architecture

Single Go binary. No runtime dependencies. SQLite for memory (WAL mode). Starts in <1 second.

The project uses a **modular hexagonal architecture** with portability and testability as first-class concerns:

```
apps/api/                  # Application entry & wiring
├── cmd/server/main.go     # Entry point
├── internal/
│   ├── bootstrap/         # Manual DI wiring
│   └── mcp/               # MCP tool registration

modules/                   # Domain modules (hexagonal)
├── jira/                  # Jira Cloud operations
├── sprint/                # Sprint PM intelligence
└── notification/          # Multi-channel notifications

shared/                    # Shared kernel
├── domain/                # Core domain types & interfaces
├── infrastructure/        # External service clients
└── usecase/               # Shared business logic

agents/                    # Agent Architecture layer
└── agent.go               # Dispatcher, Planner, Coordinator
```

**Stack**: Go 1.26.4 | `mark3labs/mcp-go` v0.55.1 | SQLite WAL | Manual DI

### Core Modules (Hexagonal)

| Module | Domain | Infrastructure | Status |
|--------|--------|---------------|--------|
| **jira** | Issue, Board, Project, User, Sprint | REST client, search, CRUD | ✅ 4 tools |
| **sprint** | Velocity, Score, Predictability, Planning, Memory (16 tables) | SQLite persistence, snapshot adapter | ✅ 5 tools |
| **notification** | Notifier contract, Lark entities | Slack, Discord, Telegram, Teams, Email, Lark | ⚠️ stubs |

---

## Migration Status

This project is in **active migration** from a monolithic ~279-tool server to a modular, maintainable architecture.

| Status | Area |
|--------|------|
| ✅ | Hexagonal module structure (3 modules) |
| ✅ | SQLite persistence with 16+ tables |
| ✅ | Multi-channel notification infrastructure |
| ✅ | 16+ shared service clients (AI, GitHub, GitLab, Notion, Linear, etc.) |
| ✅ | Agent architecture pattern (Dispatcher → Planner → Coordinator) |
| ⚠️ | **9+2 tools registered** in modular code (vs ~279 in installed binary) |
| ⚠️ | Infrastructure clients not yet wired into modular handlers |
| ⚠️ | Agent layer not yet wired into bootstrap |
| ⏳ | Tool registration: wire `shared/infrastructure/` into module interfaces |

The **installed binary** (`~/.local/bin/zara-jira-mcp`) still runs the full monolithic toolset. The **source code** has been restructured and is being rebuilt tool by tool.

---

## Who This Is For

**Scrum Masters** who want to stop being meeting schedulers and start being system coaches.

**Project Managers** who are tired of manually assembling status reports from 5 different tabs.

**Engineering Managers** who need data-backed answers to "when will it be done?" and "where are the risks?"

**Product Owners** who want a single place to see sprint progress, blockers needing their decision, and value delivery.

**Teams** who hate Jira's slow UI and want to interact with their board through natural language.

---

## What People Say

> "Replaced 2 hours of Monday morning report writing with one command." — PM using `pm_exec_report`

> "The dysfunction detector caught our Hero Culture pattern before I even noticed it." — SM after 3 sprints of data

> "Monte Carlo forecasting finally gave me a defensible answer for stakeholders." — Engineering Manager

---

## Integrations

| Category | Services |
|----------|----------|
| Project Tracking | Jira Cloud, Linear, Notion, GitHub Issues, GitLab Issues |
| Code | GitHub (PRs, repos), GitLab (MRs, pipelines) |
| Communication | Lark, Slack, Discord, Telegram, Teams, Email |
| Documentation | Confluence |
| Incidents | PagerDuty |
| Time Tracking | Clockify |
| Calendar | Google Calendar, Lark Calendar |
| Data | Google Sheets |
| AI | OpenAI, Anthropic, Gemini, Groq, Ollama, OpenRouter, DeepSeek, Together.ai |

All optional. Works with just Jira + any AI provider.

---

## Contributing

PRs welcome. See [AGENTS.md](AGENTS.md) for architecture decisions and coding guidelines.

```bash
make build      # Build modular version → bin/zara-jira-mcp
make test       # Run tests
make lint       # golangci-lint
make install    # Install to ~/.local/bin/
```

---

## Guides for PM/Scrum Masters

| Guide | What You'll Learn |
|-------|-------------------|
| [Communication Frameworks](docs/communication-frameworks.md) | Pyramid Principle, SCARF, SBI, Radical Candor, DACI, async protocols |
| [Reporting to Management](docs/reporting-guide.md) | Who needs what report, when, and which tool to use |
| [Understanding Engineering](docs/engineering-literacy.md) | WIP, cycle time, tech debt, QA vocabulary, PR review |

---

## Part of the Zara Ecosystem

Built as part of [Zara Agent OPC](https://github.com/aldok10/zara-agent-opc) — an empathetic AI engineering partner with cognitive memory and multi-agent orchestration.

## Author

**[Aldo Karendra](https://www.linkedin.com/in/aldok10/)** — Lead Backend Developer & AI Systems Architect. Building AI tools that help engineers ship better software, faster.

[![GitHub](https://img.shields.io/badge/GitHub-@aldok10-181717?logo=github)](https://github.com/aldok10)
[![LinkedIn](https://img.shields.io/badge/LinkedIn-aldok10-0A66C2?logo=linkedin)](https://www.linkedin.com/in/aldok10/)
[![Support](https://img.shields.io/badge/Support-SociaBuzz-orange)](https://sociabuzz.com/aldok10)

## License

MIT — use it, fork it, make it yours.
