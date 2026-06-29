# zara-jira-mcp

**Your AI Scrum Master that actually remembers.**  
89 tools. Persistent memory. Jira + GitHub + GitLab + 9 more platforms. Works in any AI editor.

[![CI](https://github.com/aldok10/zara-jira-mcp/actions/workflows/ci.yml/badge.svg)](https://github.com/aldok10/zara-jira-mcp/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![MCP Compatible](https://img.shields.io/badge/MCP-compatible-blue)](https://modelcontextprotocol.io)
[![Go](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go)](https://golang.org)

---

## You're a Scrum Master, Not a Secretary

Every sprint, the same cycle:

| You Spend Time On | Instead Of |
|---|---|
| 20 min writing status updates | Coaching your team |
| 15 min chasing blockers manually | Removing impediments |
| Guessing "when will it be done?" | Running data-backed forecasts |
| Copy-pasting Jira data into slides | Actually improving the process |
| Forgetting retro actions after 3 days | Building continuous improvement |

**zara-jira-mcp breaks that cycle.** It connects your AI editor directly to Jira, GitHub, GitLab, and your PM memory — so every status update, blocker report, and sprint forecast happens in seconds, not hours.

---

## What Changes

### Before

```
You: "Let me open Jira, check what's in progress,
     see who's blocked, check GitHub PRs,
     manually count story points, write status,
     copy to Slack..."

→ 30 minutes later, you have a status update
  that's already outdated.
```

### After

```
You (in your AI editor): "pm board_id=84"

→ Instant: sprint status, blockers, risks, actions.
  Memory knows last sprint's decisions.
  Board config knows your custom statuses.

→ 10 seconds. Always up to date. Zero copy-paste.
```

**That's 4+ hours per week** back to things that actually matter.

---

## What 89 Tools Do For You

```
# 👀 Sprint Status — 10 seconds
pm board_id=84           → 12/20 done, 3 blocked, 2 risks, 4 pending actions

# 🚫 Blockers — detected automatically
jira_search jql=...       → auto-records blockers to memory
pm_blockers                → see all active + aging

# 🔮 Forecast — data, not guesses
pm_forecast board_id=84   → 50% Sprint 13, 85% Sprint 14 (Monte Carlo)

# 🧠 Memory — remembers everything
pm_record_risk title="API perf" severity=high
pm_record_decision title="Use PostgreSQL" decision="..."
pm_record_retro sprint_name="Sprint 12"

# 🔗 Dev Workflow — across platforms
pm_github_prs              → open PRs with ages
pm_github_activity days=7  → commits, merges, issues
pm_gitlab_merge_requests   → GitLab MRs

# 📢 Notify — multi-channel
notify_routed content="..." severity=high audience=team
jira_notify_lark content="Sprint done"

# 📊 Board-aware — understands your custom statuses
jira_board_config board_id=84  → shows your column layout
```

---

## Quick Start

```bash
git clone https://github.com/aldok10/zara-jira-mcp.git && cd zara-jira-mcp
cp .env.example .env
make build && make install
```

Add to your MCP client:

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

**Works with:** OpenCode, Claude Desktop, Cursor, VS Code Copilot, ChatGPT, Zed, Kiro, Gemini CLI, Goose — any MCP-compatible client.

---

## Why This Exists

Most Scrum tools are either:

| This | vs | Jira AI | vs | Monday/Asana | vs | Custom scripts |
|------|----|---------|----|-------------|----|----------------|
| **Remembers** across sprints | | No memory at all | | Some memory | | None |
| **Works in your AI editor** | | Jira UI only | | Their UI only | | Terminal only |
| **Forecasts with Monte Carlo** | | Basic trends | | No forecasting | | Nothing |
| **0$ — MIT license** | | Paid Premium | | Paid per seat | | Your time |
| **Self-hosted, data stays local** | | Atlassian cloud | | Their cloud | | Your infra |

Built for the AI-native workflow: **talk to your AI editor, it talks to Jira.** No dashboards to check, no tabs to switch.

---

## Architecture — Simple on Purpose

```bash
# One binary. Zero runtime deps. Starts in <1s.
~/.local/bin/zara-jira-mcp
```

| Module | Tools | What |
|--------|-------|------|
| Jira | 28 | Full CRUD + board config + sprint management |
| Sprint/PM | 18 | Memory, forecast, decisions, risks, retros |
| GitHub | 10 | PRs, issues, milestones, activity, branches |
| GitLab | 9 | MRs, issues, files, milestones, branches |
| Notifications | 5 | Lark, Slack, Discord, Telegram + smart routing |
| Calendar | 3 | Lark calendar events + meetings |
| Time | 2 | Clockify entries + reports |
| Wiki | 6 | Confluence + Notion read/create/search |
| Incidents | 2 | PagerDuty + on-call |
| Linear | 3 | Issues, cycles, activity |
| Sheets | 1 | Google Sheets read |
| Tools | 2 | Backup (JSON) + Onboard wizard |

**Total: 89 tools.** All optional. Works with just Jira + SQLite.

See [SKILL.md](SKILL.md) for the full reference with all parameters and workflow patterns.

---

## How Auto-Memory Works

Every read (`jira_search`, `jira_get_issue`, `jira_sprint_summary`) automatically:

1. **Records blockers** — detects blocked statuses, dedup by issue key
2. **Flags stale risks** — Highest/Critical untouched >7 days → auto-recorded as risks
3. **Snapshots sprints** — saves done/in-progress/blocked/todo counts
4. **Reconciles** — compares stored items against Jira state, auto-resolves

All persists in **SQLite WAL** at `~/.zara-jira-mcp/pm_memory.db`.

---

## Research

Built on validated research, not vibes:

- **DORA 2024**: AI without process fundamentals hurts stability
- **Monte Carlo Forecasting**: Historical throughput beats gut-feel estimates
- **Little's Law**: WIP limits + cycle time = only levers for faster delivery
- **Tuckman's Stages**: Team maturity determines the right SM stance
- **McKinsey 2026**: AI-augmented PM roles are the fastest-growing function

---

## For Everyone

| Role | What You Get |
|------|-------------|
| **Scrum Master** | Stop being a meeting scheduler. Start being a system coach. Data for every conversation. |
| **Engineering Manager** | Defensible answers to "when will it be done?" and "where are the risks?" |
| **Product Owner** | Sprint progress, blocked items needing your decision, value delivery — in one command. |
| **Developer** | Talk to Jira through your AI editor. No slow UI. No context switching. |

---

## Build

```bash
make build      # → bin/zara-jira-mcp
make test       # All tests
make lint       # golangci-lint
make install    # → ~/.local/bin/
```

---

**MIT License** — use it, fork it, break it, make it yours.

**[Aldo Karendra](https://github.com/aldok10)** — Making Scrum Masters more effective, one tool at a time.
