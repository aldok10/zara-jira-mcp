# Skill: zara-pm-brain

Use when managing IT projects, running sprints, tracking risks, making decisions, conducting retrospectives, preparing standups, or needing PM/Scrum Master recommendations.

## What This MCP Does

`zara-jira-mcp` is an AI-powered PM/Scrum Master brain with persistent memory. It connects to Jira for live data, maintains historical context in SQLite, and uses AI to provide proactive recommendations. **52 tools** across 13 categories.

## Tool Categories

### 1. Live Jira Data (10 tools)

| Tool | Purpose |
|------|---------|
| `jira_search` | Search issues with JQL |
| `jira_get_issue` | Get full issue details |
| `jira_boards` | List boards |
| `jira_sprint_summary` | Active sprint breakdown |
| `jira_create_issue` | Create new issue |
| `jira_add_comment` | Add comment to issue |
| `jira_transitions` | List available status transitions |
| `jira_transition` | Move issue to new status |
| `jira_my_issues` | Current user's issues |
| `jira_overdue` | Stale/overdue issues |
| `jira_workload` | Team workload distribution |

### 2. PM Memory — Sprint & Velocity (3 tools)

| Tool | When to Use |
|------|-------------|
| `pm_snapshot_sprint` | End of every sprint — captures state for velocity tracking |
| `pm_velocity_trend` | Sprint planning — see trends, detect decline |
| `pm_capacity_plan` | Sprint planning — recommended commitment based on history + availability |

### 3. PM Memory — Risk Management (4 tools)

| Tool | When to Use |
|------|-------------|
| `pm_record_risk` | When risk identified |
| `pm_update_risk` | When mitigated/resolved/accepted |
| `pm_risk_dashboard` | Sprint planning, standup, weekly review |
| `pm_auto_detect_risks` | Weekly — proactively scans Jira for risk signals |

### 4. PM Memory — Blockers & Dependencies (6 tools)

| Tool | When to Use |
|------|-------------|
| `pm_record_blocker` | When team is blocked |
| `pm_resolve_blocker` | When unblocked |
| `pm_blockers` | Daily standup — view active/history |
| `pm_record_dependency` | When cross-team/cross-issue dependency found |
| `pm_resolve_dependency` | When dependency satisfied |
| `pm_dependencies` | Sprint planning — dependency map |

### 5. PM Memory — Decisions & Knowledge (2 tools)

| Tool | When to Use |
|------|-------------|
| `pm_record_decision` | After any significant technical/process decision |
| `pm_search_decisions` | Before re-deciding — check institutional memory |

### 6. PM Memory — Team Health (2 tools)

| Tool | When to Use |
|------|-------------|
| `pm_record_team_metric` | End of sprint — per-member stats |
| `pm_team_health` | Sprint planning — workload overview, burnout signals |

### 7. PM Memory — Retrospectives & Actions (3 tools)

| Tool | When to Use |
|------|-------------|
| `pm_record_retro` | After retro ceremony |
| `pm_action_items` | Every standup — pending retro follow-ups |
| `pm_retro_analysis` | Before retro — AI pattern analysis across retros |

### 8. PM Memory — Meetings (2 tools)

| Tool | When to Use |
|------|-------------|
| `pm_record_meeting` | After any ceremony (standup, planning, grooming, adhoc) |
| `pm_meetings` | Review meeting history, find past decisions |

### 9. Sprint Health (2 tools)

| Tool | When to Use |
|------|-------------|
| `pm_sprint_health` | Anytime — computes 0-100 health score with breakdown |
| `pm_health_history` | Trend tracking — is team getting healthier? |

### 10. AI Intelligence (4 tools)

| Tool | When to Use |
|------|-------------|
| `pm_recommendations` | Weekly — AI recs from ALL historical memory |
| `pm_standup_prep` | Before daily standup — auto-generated talking points |
| `jira_ai_analyze` | Ad-hoc AI analysis of any tickets |
| `jira_ai_sprint_report` | End of sprint — full report, optional Lark send |

### 11. Notifications (1 tool)

| Tool | Purpose |
|------|---------|
| `jira_notify_lark` | Send markdown to Lark group |

### 12. Burndown & Daily Tracking (2 tools)

| Tool | When to Use |
|------|-------------|
| `pm_track_daily` | Daily — capture today's sprint progress for burndown |
| `pm_burndown` | View burndown with burn rate + days-to-complete estimate |

### 13. Sprint Goals (3 tools)

| Tool | When to Use |
|------|-------------|
| `pm_set_sprint_goal` | Sprint planning — define goal + key results |
| `pm_close_sprint_goal` | End of sprint — record achieved/missed + outcome |
| `pm_sprint_goals` | View active goals or achievement history |

### 14. Definition of Done (1 tool)

| Tool | When to Use |
|------|-------------|
| `pm_dod` | Manage DoD checklist (add/remove/list per project) |

### 15. Escalation System (3 tools)

| Tool | When to Use |
|------|-------------|
| `pm_escalate` | Auto-escalate critical risks/blockers to Lark |
| `pm_escalations` | View escalation history |
| `pm_dashboard` | One-shot full PM view (everything in one call) |

### 16. Release (1 tool)

| Tool | When to Use |
|------|-------------|
| `pm_release_notes` | Generate categorized release notes from done issues |

## PM Workflow Patterns
```
1. pm_standup_prep(board_id)           -> AI-generated talking points
2. pm_blockers()                       -> what's stuck
3. pm_action_items()                   -> retro follow-ups overdue
4. pm_dependencies()                   -> blocking/blocked chains
```

### Sprint Planning
```
1. pm_capacity_plan(board_id, team_size, planned_leave_days)  -> recommended points
2. pm_velocity_trend(board_id)         -> historical context
3. pm_risk_dashboard()                 -> risks for next sprint
4. pm_dependencies()                   -> unresolved dependencies
5. pm_team_health(sprint_name: prev)   -> who's overloaded?
6. pm_recommendations(focus:"velocity")-> AI capacity advice
```

### Mid-Sprint Health Check
```
1. pm_sprint_health(board_id)          -> health score 0-100
2. pm_auto_detect_risks(board_id)      -> proactive risk scan
3. pm_blockers()                       -> stuck items
```

### End of Sprint
```
1. pm_snapshot_sprint(board_id, velocity:X, carryover:Y)  -> capture state
2. pm_record_team_metric(...)          -> for each member
3. pm_sprint_health(board_id)          -> final health score
4. jira_ai_sprint_report(board_id, send_to_lark:true)     -> share report
```

### Retrospective
```
1. pm_retro_analysis(board_id)         -> AI patterns from history
2. pm_record_retro(sprint, well, improve, actions)
3. pm_action_items()                   -> verify old items closed
4. pm_record_meeting(type:"retro", decisions, action_items)
```

### Risk Management (Ongoing)
```
1. pm_auto_detect_risks(board_id)      -> automated scan
2. pm_record_risk(title, severity, owner, mitigation)
3. pm_risk_dashboard()                 -> prioritized view
4. pm_update_risk(id, status:"resolved")
5. pm_recommendations(focus:"risks")   -> AI risk advice
```

### Decision Logging
```
1. pm_record_decision(title, decision, context, rationale, tags)
2. pm_search_decisions("database")     -> before re-deciding
```

### After Any Meeting
```
1. pm_record_meeting(type, notes, decisions, action_items, attendees)
```

### Daily Burndown
```
1. pm_track_daily(board_id)            -> capture today's data point
2. pm_burndown(board_id)              -> view chart + burn rate
```

### Sprint Goal Lifecycle
```
1. pm_set_sprint_goal(board_id, goal, key_results)      -> at planning
2. pm_sprint_goals(board_id)                            -> check during sprint
3. pm_close_sprint_goal(goal_id, status, outcome)       -> at sprint end
4. pm_sprint_goals(board_id, show_history:true)         -> track over time
```

### Definition of Done Setup
```
1. pm_dod(action:"add", item:"Unit tests pass", category:"testing")
2. pm_dod(action:"add", item:"Code reviewed", category:"review")
3. pm_dod(action:"add", item:"No critical bugs", category:"testing")
4. pm_dod(project:"MYPROJ")           -> view project-specific DoD
```

### Escalation & Alerting
```
1. pm_escalate(board_id)              -> auto-check thresholds, send to Lark
2. pm_escalations()                   -> review what was escalated
```

### Release Day
```
1. pm_release_notes(board_id, send_to_lark:true) -> generate + share
```

### Full Status Check (One Command)
```
1. pm_dashboard(board_id)             -> everything in one view
```

## Health Score Breakdown (pm_sprint_health)

| Component | Max | Measures |
|-----------|-----|----------|
| Velocity | 25 | Completion rate (done / total) |
| Blockers | 25 | Blocked ratio (fewer = higher) |
| Scope | 25 | Scope change vs previous sprint |
| Team | 25 | Workload distribution balance |

Overall: 0-49 = AT RISK, 50-69 = WATCH, 70-100 = HEALTHY

## Auto Risk Detection (pm_auto_detect_risks)

Automatically scans for:
1. Stale tickets (7+ days no update in active sprint)
2. Workload imbalance (someone has 2x average load)
3. High blocked count (3+ issues blocked)
4. Chronic blockers (5+ days unresolved)
5. Overdue retro action items

All findings auto-recorded to risk register.

## Auto Escalation (pm_escalate)

Thresholds that trigger Lark alerts:
1. Critical/High risks open >3 days without resolution
2. Blockers unresolved >3 days
3. Sprint health score below 50

Escalation history tracked — prevents duplicate alerts.

## Principles

1. **Record everything that matters** — Decisions, risks, blockers. Future you will thank present you.
2. **Snapshot every sprint** — Without data, recommendations are guesses.
3. **Track blockers with resolution** — Patterns emerge over time.
4. **Never let retro actions die** — Check `pm_action_items` every standup.
5. **Use AI after building history** — Recommendations improve with more data.
6. **Team metrics are for support, not surveillance** — Detect overload early, not punish.
7. **Dependencies are risks in disguise** — Track them before they block.
8. **Decisions without rationale are forgotten** — Always record the "why".
9. **Health scores are signals, not grades** — Use for conversation, not blame.
10. **Automate detection, humanize response** — Let tools find problems, humans solve them.

## Data Lives In

SQLite at `~/.zara-jira-mcp/pm_memory.db` (configurable via `PM_MEMORY_DB_PATH`).
Back up this file to preserve all PM memory across machines.

Tables: `sprint_snapshots`, `risks`, `decisions`, `blockers`, `team_metrics`, `retrospectives`, `action_items`, `dependencies`, `meeting_notes`, `health_scores`, `daily_progress`, `sprint_goals`, `dod_items`, `escalations`
