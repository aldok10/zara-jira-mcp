# Skill: zara-pm-brain

Use when managing IT projects, running sprints, tracking risks, making decisions, conducting retrospectives, preparing standups, or needing PM/Scrum Master recommendations.

## What This MCP Does

`zara-jira-mcp` is an AI-powered PM/Scrum Master brain with persistent memory. It connects to Jira for live data, maintains historical context in SQLite, and uses AI to provide proactive recommendations.

## Tool Categories

### 1. Live Jira Data

| Tool | Purpose |
|------|---------|
| `jira_search` | Search issues with JQL |
| `jira_get_issue` | Get full issue details |
| `jira_boards` | List boards |
| `jira_sprint_summary` | Active sprint breakdown |
| `jira_my_issues` | Current user's issues |
| `jira_overdue` | Stale/overdue issues |
| `jira_workload` | Team workload distribution |

### 2. PM Memory (Persistent State)

| Tool | Purpose | When to Use |
|------|---------|-------------|
| `pm_snapshot_sprint` | Save sprint state | End of every sprint |
| `pm_record_risk` | Add to risk register | When risk identified |
| `pm_update_risk` | Change risk status | When mitigated/resolved |
| `pm_risk_dashboard` | View open risks | Sprint planning, standup |
| `pm_record_decision` | Log a decision + rationale | After any significant decision |
| `pm_search_decisions` | Find past decisions | Before re-deciding something |
| `pm_record_blocker` | Track impediment | When blocked |
| `pm_resolve_blocker` | Mark resolved | When unblocked |
| `pm_blockers` | View active/history | Daily standup |
| `pm_record_team_metric` | Member sprint stats | End of sprint |
| `pm_team_health` | Workload overview | Sprint planning |
| `pm_record_retro` | Save retrospective | After retro ceremony |
| `pm_action_items` | Pending retro actions | Every standup, sprint planning |

### 3. AI Intelligence (Memory-Powered)

| Tool | Purpose | When to Use |
|------|---------|-------------|
| `pm_recommendations` | AI recommendations from ALL history | Weekly, sprint planning |
| `pm_velocity_trend` | Velocity over time + trend detection | Capacity planning |
| `pm_standup_prep` | Standup talking points | Before daily standup |
| `pm_retro_analysis` | Pattern analysis across retros | Before retro, quarterly review |
| `jira_ai_analyze` | Ad-hoc AI analysis of tickets | Any time |
| `jira_ai_sprint_report` | Full sprint report | End of sprint |

### 4. Notifications

| Tool | Purpose |
|------|---------|
| `jira_notify_lark` | Send updates to Lark |

## PM Workflow Patterns

### Daily Standup Prep
```
1. pm_standup_prep(board_id) -> talking points
2. pm_blockers() -> what's stuck
3. pm_action_items() -> retro follow-ups
```

### Sprint Planning
```
1. pm_velocity_trend(board_id) -> how much can we commit?
2. pm_risk_dashboard() -> what risks affect next sprint?
3. pm_recommendations(board_id, focus:"velocity") -> AI capacity advice
4. pm_team_health(sprint_name: last_sprint) -> who's overloaded?
```

### End of Sprint
```
1. pm_snapshot_sprint(board_id, velocity:X, carryover:Y) -> capture state
2. pm_record_team_metric(...) -> for each member
3. jira_ai_sprint_report(board_id, send_to_lark:true) -> share report
```

### Retrospective
```
1. pm_retro_analysis(board_id) -> patterns from history
2. pm_record_retro(sprint_name, went_well, improvements, action_items)
3. pm_action_items() -> verify old items resolved
```

### Risk Management
```
1. pm_record_risk(title, severity, owner, mitigation)
2. pm_risk_dashboard() -> review in planning
3. pm_update_risk(id, status:"resolved") -> close when done
4. pm_recommendations(focus:"risks") -> AI risk advice
```

### Decision Logging
```
1. pm_record_decision(title, decision, context, rationale, tags)
2. pm_search_decisions("database") -> before re-deciding
```

## Principles

1. **Record everything that matters** - Decisions, risks, blockers. Future you will thank present you.
2. **Snapshot every sprint** - Without data, recommendations are guesses.
3. **Track blockers with resolution** - Patterns emerge over time.
4. **Never let retro actions die** - Check `pm_action_items` every standup.
5. **Use AI after building history** - Recommendations improve with more data.
6. **Team metrics are for support, not surveillance** - Detect overload early, not punish.

## Data Lives In

SQLite at `~/.zara-jira-mcp/pm_memory.db` (configurable via `PM_MEMORY_DB_PATH`).
Back up this file to preserve all PM memory across machines.
