# AGENTS.md

Instructions for AI agents working with this project.

## What This Is

`zara-jira-mcp` is a 124-tool MCP server that acts as an AI-powered Scrum Master with persistent memory. It connects to Jira Cloud, runs Monte Carlo forecasts, detects anti-patterns, manages risks/blockers/decisions, and sends notifications across Lark/Slack/Discord/Telegram/Teams/Email.

Built with Go 1.26. SQLite for persistent memory (14 tables). OpenAI-compatible API for AI intelligence.

## First Things First

Before using any PM/sprint tools, get the board ID:

```
jira_boards -> returns board IDs
```

Store the board_id. Almost every PM tool needs it.

## Tool Categories

| Category | Count | Use When |
|----------|-------|----------|
| Jira Operations | 45 | CRUD on issues, sprints, epics, bulk ops |
| PM Memory | 20 | Recording decisions, risks, blockers, retros, team metrics |
| AI Intelligence | 15 | Forecasting, coaching, anti-patterns, recommendations |
| Process & Health | 18 | Sprint health, velocity, capacity, DoD/DoR, goals |
| Escalation & Reporting | 8 | Executive reports, release notes, weekly digests |
| Workflow Recipes | 3 | One-click start/done/block workflows |
| Notifications | 9 | Multi-channel messaging |
| Portfolio | 5 | Cross-project overview |
| Calendar | 4 | Lark calendar events and meetings |

## Critical Rules

1. **Memory builds over time.** Intelligence tools need historical data. Encourage snapshot recording at sprint boundaries.
2. **`pm_snapshot_sprint` at end of EVERY sprint.** Without it, forecasting/velocity/capacity tools return nothing useful.
3. **Record immediately.** After any decision, risk, or blocker: record it. Don't wait.
4. **Never use `pm_dashboard` for executives.** Use `pm_exec_report` instead (no jargon, business outcomes only).
5. **Run `pm_auto_detect_risks` weekly.** Proactive scanning catches problems early.

## Common Workflows

### Daily Standup
```
pm_standup_prep(board_id) -> talking points, blockers, risks, action items
```

### Sprint Planning
```
pm_planning_prep(board_id) -> capacity, carryover, risks, deps, experiments
```

### "When Will It Be Done?"
```
pm_forecast(board_id, remaining_items:N) -> 50%/70%/85%/95% confidence dates
```

### End of Sprint
```
pm_snapshot_sprint(board_id, velocity:N, carryover:N)
pm_scorecard(board_id)
pm_release_notes(board_id)
pm_close_sprint_goal(goal_id, status)
```

### Executive Update
```
pm_exec_report(board_id) -> business outcomes, no story points
```

### Report to PO / Product Owner
```
pm_goal_check(board_id) -> is sprint goal on track?
pm_scope_creep(board_id) -> what changed mid-sprint without approval?
pm_forecast(board_id, remaining_items:N) -> realistic delivery date
```

### Escalation to Management
```
pm_impediment_aging -> all blockers with age, chronic flags
pm_escalate(board_id) -> auto-alert if risk/blocker > 3 days or health < 50
pm_stakeholder_pulse(stakeholder, score:1-5, feedback) -> track satisfaction
```

### Cross-Team Dependencies
```
pm_dependencies -> who blocks whom, across teams
portfolio_blockers -> all blockers across all projects
portfolio_summary -> AI exec summary for steering committee
```

### Prove SM Value
```
pm_sm_impact(sprint_name) -> blockers resolved, resolution time, risks mitigated
pm_maturity_assessment(board_id) -> team stage with evidence
```

### Team Seems Off
```
pm_anti_patterns(board_id) -> zombie sprints, hero culture, scope creep
pm_coaching(topic:"team_dynamics", situation:"...")
```

## Reporting Guide

See `docs/reporting-guide.md` for scenario-based guide: who needs what report, when, and which tool to use.

## Architecture (for code contributions)

```
cmd/server/          Entry point (uber-go/fx DI)
config/              Env configuration
domain/              Interfaces + models (jira, memory, ai, lark)
internal/            Implementations (jirasdk, sqlite, openai, slack, lark)
application/tools/   MCP tool handlers
transport/           MCP server + tool registration
```

## Setup Guides

See `docs/agents/` for platform-specific setup:
- `docs/agents/opencode.md` — OpenCode
- `docs/agents/claude-code.md` — Claude Code
- `docs/agents/cursor.md` — Cursor
- `docs/agents/copilot.md` — GitHub Copilot (VS Code)
- `docs/agents/vscode-copilot.md` — VS Code + Copilot (mcp.json)
- `docs/agents/windsurf.md` — Windsurf
- `docs/agents/kiro.md` — Kiro
- `docs/agents/chatgpt.md` — ChatGPT Desktop
- `docs/agents/codex.md` — Codex CLI
- `docs/agents/gemini-cli.md` — Gemini CLI
- `docs/agents/zed.md` — Zed Editor
- `docs/agents/goose.md` — Goose (Block)
- `docs/agents/cline.md` — Cline / Roo Code

Pre-built config files included in the repo:
- `.claude/` — Claude Code settings + CLAUDE.md
- `.opencode/` — OpenCode config + instructions
- `.cursor/` — Cursor rules
- `.kiro/` — Kiro config + instructions
- `.vscode/` — VS Code Copilot MCP config
- `.zed/` — Zed context_servers config
- `.github/copilot-instructions.md` — GitHub Copilot instructions
- `.codex/` — Codex CLI instructions
- `.windsurfrules` — Windsurf project rules

## Full Tool Reference

See `SKILL.md` for complete tool documentation with parameters, return values, and usage patterns.
