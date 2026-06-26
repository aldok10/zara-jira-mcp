# Skill: zara-pm-brain

Use when managing IT projects, running sprints, tracking risks, making decisions, conducting retrospectives, preparing standups, forecasting delivery, coaching teams, or needing PM/Scrum Master recommendations.

## What This MCP Does

`zara-jira-mcp` is an AI-powered PM/Scrum Master brain with persistent memory. 124 tools. Connects to Jira for live data, maintains historical context in SQLite, uses AI for proactive recommendations, sends alerts to Lark/Slack.

## Quick Reference: When to Use What

### Before Standup
```
pm_standup_prep(board_id)
```

### Sprint Planning
```
pm_planning_prep(board_id)       <- everything in one shot
pm_forecast(board_id)            <- "when will it be done?"
pm_capacity_plan(board_id, team_size, planned_leave_days)
pm_check_ready(key:"PROJ-123")   <- is this story ready?
pm_confidence(sprint_name, score:4)
```

### During Sprint
```
pm_dashboard(board_id)           <- one-shot full view
pm_flow_metrics(board_id)        <- WIP, cycle time, throughput
pm_track_daily(board_id)         <- burndown data point
pm_goal_check(board_id)          <- on track?
pm_auto_detect_risks(board_id)   <- proactive scan
pm_scope_creep(board_id)         <- scope change?
```

### End of Sprint
```
pm_snapshot_sprint(board_id, velocity:X)
pm_scorecard(board_id)           <- A-F grade
pm_sprint_compare(board_id)      <- vs last sprint
pm_release_notes(board_id, send_to_lark:true)
pm_close_sprint_goal(goal_id, status:"achieved")
```

### Retrospective
```
pm_facilitate(ceremony:"retro")  <- fresh format each time
pm_retro_analysis(board_id)      <- patterns from history
pm_record_retro(sprint_name, went_well, improvements, action_items)
pm_experiment(hypothesis, action, measure)
```

### Stakeholder Reporting
```
pm_exec_report(board_id)         <- VP-friendly, no jargon
pm_weekly_digest(board_id, send_to_lark:true)
```

### Risk & Blocker Management
```
pm_record_risk(title, severity, owner, mitigation)
pm_risk_dashboard()
pm_record_blocker(description, issue_key, owner)
pm_escalate(board_id)            <- auto-alert critical items
```

### Team Health
```
pm_sprint_health(board_id)       <- 0-100 score
pm_anti_patterns(board_id)       <- detect dysfunctions
pm_coaching(topic:"team_dynamics", situation:"...")
pm_team_kb(question:"how does this team work?")
```

### Process Setup (One-Time)
```
pm_dod(action:"add", item:"Unit tests pass", category:"testing")
pm_dor(action:"add", item:"Acceptance criteria defined", category:"clarity")
pm_agreements(action:"add", agreement:"PRs reviewed within 24h")
```

## Key Tools by Category

| Category | Count | Key Tools |
|----------|-------|-----------|
| Jira Operations | 45 | search, create, update, transition, bulk ops, epics, sprints |
| PM Memory | 20 | snapshot, risks, decisions, blockers, team, retros, deps |
| AI Intelligence | 15 | forecast, coaching, facilitate, anti-patterns, recommendations |
| Process & Health | 18 | health score, velocity, capacity, goals, DoD/DoR, experiments |
| Reporting | 8 | exec report, weekly digest, scorecard, release notes |
| Notifications | 9 | Lark, Slack, Discord, Telegram, Teams, Email, Confluence |
| Recipes | 3 | start_work, done, block (one-click) |

## Health Score Breakdown

| Component | Max | Measures |
|-----------|-----|----------|
| Velocity | 25 | Completion rate |
| Blockers | 25 | Blocked ratio |
| Scope | 25 | Change vs baseline |
| Team | 25 | Workload balance |

0-49 = AT RISK, 50-69 = WATCH, 70-100 = HEALTHY

## Anti-Pattern Detection

Detects from real data:
- Zombie Sprint (>30% carryover consistently)
- Hero Culture (one person does >50% of work)
- Scope Creep (>20% growth mid-sprint)
- Unpredictable (>40% velocity variance)
- Dead Retros (>5 pending action items never done)
- Rubber-Stamp DoD (>95% completion with zero blockers)
- No Sprint Goals (no recorded goals)

## Monte Carlo Forecasting

`pm_forecast(board_id, remaining_items:30)`

Returns probability-based dates:
- 50% confidence (coin flip)
- 70% confidence
- 85% confidence (recommended for stakeholder commitments)
- 95% confidence (almost certain)

Based on 10,000 simulations using historical throughput.

## Data Persistence

SQLite at `~/.zara-jira-mcp/pm_memory.db` (14 tables).
Back up this file to preserve all PM memory.

## Principles

1. Record everything that matters — future you will thank present you
2. Snapshot every sprint — without data, recommendations are guesses
3. Track blockers with resolution — patterns emerge over time
4. Never let retro actions die — check every standup
5. Flow over velocity — WIP and cycle time predict delivery better
6. Team metrics are for support, not surveillance
7. Automate detection, humanize response
8. Decisions without rationale are forgotten
9. Forecast with probability, not promises
10. Process gates (DoR/DoD) prevent rework
