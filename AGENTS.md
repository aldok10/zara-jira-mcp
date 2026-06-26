# AGENTS.md

Instructions for AI agents working with this project.

## What This Is

`zara-jira-mcp` is a ~285-tool MCP server that acts as an AI-powered Scrum Master with persistent memory, empathy, and learning capability. It connects to Jira Cloud, runs Monte Carlo forecasts, detects anti-patterns, reads team sentiment, manages feedback lifecycles, tracks OKR/KPI progress, and sends notifications across Lark/Slack/Discord/Telegram/Teams/Email.

Built with Go 1.26. SQLite for persistent memory (16+ tables). OpenAI-compatible API for AI intelligence. Lark OKR bi-directional sync.

**Design principle:** Tool-surface compression (BCG 2026, Microsoft Research). Use profiles to control tool count. For most use cases, start with `pm_smart`.

## First Things First

Before using any PM/sprint tools, get the board ID:

```
jira_boards -> returns board IDs
```

Store the board_id. Almost every PM tool needs it.

## Profiles (Tool-Surface Compression)

Choose based on your context (BCG research: productivity peaks at 3 tools, degrades at 4+):

| Profile | Tools | Best For |
|---------|-------|----------|
| `chatgpt` | ~14 | ChatGPT Desktop (token-limited) |
| `lite` | ~65 | Solo PM, daily workflow |
| `standard` | ~120 | Full PM team + Jira |
| `full` | ~200 | Power user + GitHub + portfolio |
| `all` | ~285 | Developer/debugging (all modules) |

Set via `PM_PROFILE=chatgpt` or `PM_ENABLED_MODULES=jira,pm,ai`.

## Tool Categories

| Category | Key Tools | Use When |
|----------|-----------|----------|
| Smart Router | `pm_smart`, `pm_do`, `pm_report`, `pm_team`, `pm_plan` | Don't know which tool — just ask |
| Jira Operations | `jira_search`, `jira_get_issue`, `jira_create_issue` | CRUD on issues, sprints, epics |
| PM Memory | `pm_record_decision`, `pm_record_risk`, `pm_record_blocker` | Record events immediately |
| Intelligence | `pm_sentiment`, `pm_coaching`, `pm_anti_patterns` | Understand team state |
| Communication | `pm_comms_nudge`, `pm_conversation_prep`, `pm_feedback_log` | Human-centered PM work |
| OKR/KPI | `pm_okr_health`, `pm_okr_suggest`, `pm_kpi_trend` | Goal alignment |
| Reporting | `pm_exec_report`, `pm_daily_digest`, `pm_forecast` | Stakeholder updates |
| Notifications | `notify_routed`, `lark_send`, `slack_send` | Multi-channel messaging |
| Lark OKR | `lark_okr_pull`, `lark_okr_sync` | Bi-directional OKR sync |

## Critical Rules

1. **Memory builds over time.** Intelligence tools need historical data. Encourage snapshot recording at sprint boundaries.
2. **`pm_snapshot_sprint` at end of EVERY sprint.** Without it, forecasting/velocity/capacity tools return nothing useful.
3. **Record immediately.** After any decision, risk, or blocker: record it. Don't wait.
4. **Never use `pm_dashboard` for executives.** Use `pm_exec_report` instead (no jargon, business outcomes only).
5. **Run `pm_auto_detect_risks` weekly.** Proactive scanning catches problems early.
6. **Use `pm_context_note` for human stories.** Record WHY someone is stuck, not just THAT they're stuck.
7. **Close the feedback loop.** `pm_feedback_log` -> `pm_feedback_due` -> `pm_feedback_close`.

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
pm_dependency_report -> who blocks whom, across teams
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
pm_sentiment(board_id) -> read team mood from data signals + AI coaching
pm_anti_patterns(board_id) -> zombie sprints, hero culture, scope creep
pm_lencioni(board_id) -> 5 Dysfunctions diagnosis with coaching
pm_coaching(topic:"team_dynamics", situation:"...")
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

### Difficult Conversation
```
pm_conversation_prep(type:"performance", context:"...", person:"...") -> framework-based prep
pm_hard_conversation(situation:"...", person:"...") -> STATE + SBI + SCARF
pm_nvc_reframe(message:"...") -> rewrite using Nonviolent Communication
```

### OKR Alignment
```
pm_okr_suggest(board_id) -> AI: which sprint items serve which OKRs?
pm_okr_health -> progress vs time elapsed, flags at-risk objectives
pm_kpi_trend(name:"cycle_time") -> single KPI trend over time
pm_kpi_to_okr -> AI: suggest Key Results from current metrics
```

### Feedback Loop
```
pm_feedback_log(person, topic, type) -> record feedback given
pm_feedback_due -> show overdue follow-ups
pm_feedback_close(id, outcome) -> close the loop
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
