# Skill: zara-jira-mcp

PM/Scrum Master MCP server with persistent memory, AI intelligence, and multi-channel notifications.
Use when the user asks about Jira, sprints, team management, risks, blockers, decisions, retrospectives, forecasting, tech debt, or any PM/Scrum Master activity.

## Connection

This MCP runs as a stdio server. All tools are available via the configured MCP connection. The server maintains a SQLite database at `~/.zara-jira-mcp/pm_memory.db` that persists across sessions.

## Core Concept

This is NOT just a Jira wrapper. It is a **Scrum Master brain** with:
- **Persistent memory** (14 SQLite tables) — remembers everything across sessions
- **AI analysis** — uses historical + live data for recommendations
- **Proactive detection** — auto-finds risks, anti-patterns, scope creep
- **Multi-audience** — different reports for engineers vs executives

## Critical Rules

1. **board_id is required** for most PM tools. Get it first with `jira_boards`.
2. **Memory builds over time** — PM intelligence tools improve with more data. Encourage recording snapshots, decisions, blockers.
3. **Never expose raw data to executives** — use `pm_exec_report` for stakeholders, not `pm_dashboard`.
4. **Record before you forget** — After any decision, risk, blocker: record it immediately.
5. **Sprint snapshots are essential** — Call `pm_snapshot_sprint` at end of EVERY sprint. Without it, forecasting/velocity/capacity tools return nothing.

---

## TOOL REFERENCE (131 tools)

### Jira Operations

#### Search & Read
| Tool | Params | Returns |
|------|--------|---------|
| `jira_search` | `jql` (required), `max_results`, `start_at` | Issues: key, summary, status, priority, assignee |
| `jira_get_issue` | `key` (required) | Full issue: description, labels, dates, sprint, comments |
| `jira_boards` | none | Board ID, name, type (needed for sprint tools) |
| `jira_sprint_summary` | `board_id` | Active sprint: status breakdown + issue list |
| `jira_sprints` | `board_id`, `state` (active/future/closed) | All sprints for a board |
| `jira_my_issues` | `status` (optional filter) | Current user's unresolved issues |
| `jira_overdue` | `days` (default 14), `project` | Issues with no update in N days |
| `jira_workload` | `project` | Issue count per assignee |
| `jira_projects` | none | All accessible projects |
| `jira_project_detail` | `key` | Project metadata |
| `jira_transitions` | `key` | Available status transitions for an issue |
| `jira_link_types` | none | Available issue link types |
| `jira_watchers` | `key` | Issue watchers |
| `jira_worklog_list` | `key` | Time entries for an issue |
| `jira_find_user` | `query` | Search users (returns account IDs for assignment) |
| `jira_epic_issues` | `epic_key` | Issues in an epic |
| `jira_health` | none | Server health check |
| `jira_nl_to_jql` (via `pm_nl_to_jql`) | `query` (natural language) | Converts to JQL |

#### Create & Modify
| Tool | Params | Action |
|------|--------|--------|
| `jira_create_issue` | `project`, `summary` (required), `issue_type`, `description`, `priority`, `assignee_id`, `labels` | Create issue |
| `jira_update_issue` | `key` (required), any field to update | Update fields |
| `jira_add_comment` | `key`, `body` | Add comment |
| `jira_transition` | `key`, `transition_id` (from `jira_transitions`) | Change status |
| `jira_assign` | `key`, `account_id` | Assign to user |
| `jira_unassign` | `key` | Remove assignee |
| `jira_delete_issue` | `key` | Delete permanently |
| `jira_create_subtask` | `parent_key`, `project`, `summary` | Create subtask |
| `jira_labels_set` | `key`, `labels` | Set labels |
| `jira_link_issues` | `from_key`, `to_key`, `link_type` | Create relationship |
| `jira_watch` | `key`, `account_id` | Add watcher |
| `jira_worklog_add` | `key`, `time_spent`, `comment` | Log time |
| `jira_epic_add` | `issue_key`, `epic_key` | Add to epic |
| `jira_epic_remove` | `issue_key` | Remove from epic |

#### Sprint Management
| Tool | Params | Action |
|------|--------|--------|
| `jira_sprint_create` | `board_id`, `name`, `goal` | Create future sprint |
| `jira_sprint_start` | `sprint_id`, `start_date`, `end_date` | Start sprint |
| `jira_sprint_close` | `sprint_id` | Close sprint |
| `jira_sprint_move_issues` | `sprint_id`, `issue_keys` | Move issues into sprint |

#### Bulk Operations
| Tool | Params | Action |
|------|--------|--------|
| `jira_bulk_transition` | `keys`, `transition_id` | Transition multiple issues |
| `jira_bulk_assign` | `keys`, `account_id` | Assign multiple issues |
| `jira_bulk_label` | `keys`, `label` | Add label to multiple |

#### Advanced
| Tool | Params | Action |
|------|--------|--------|
| `jira_raw_request` | `method`, `path`, `body` | Raw Jira REST API call |
| `jira_from_branch` | `branch` | Extract issue key from git branch name |
| `jira_smart_commit` | `key`, `message` | Parse commit-style commands |
| `jira_link_pr` | `key`, `url`, `title` | Link PR/MR to issue |

---

### PM Memory Tools (Record & Retrieve)

#### Sprint Data
| Tool | When to Call | Params |
|------|-------------|--------|
| `pm_snapshot_sprint` | **End of every sprint** | `board_id`, `velocity` (points), `carryover`, `notes` |
| `pm_track_daily` | **Daily** (for burndown) | `board_id` |
| `pm_burndown` | View burndown chart | `board_id`, `sprint_name` |
| `pm_velocity_trend` | Capacity planning | `board_id` |

#### Risks
| Tool | When to Call | Params |
|------|-------------|--------|
| `pm_record_risk` | Risk identified | `title`, `severity` (critical/high/medium/low), `owner`, `mitigation`, `sprint_name` |
| `pm_update_risk` | Status changes | `risk_id`, `status` (open/mitigating/resolved/accepted) |
| `pm_risk_dashboard` | Sprint planning, standup | none |
| `pm_auto_detect_risks` | Weekly scan | `board_id` — auto-records findings |

#### Decisions & Knowledge
| Tool | When to Call | Params |
|------|-------------|--------|
| `pm_record_decision` | After significant decision | `title`, `decision`, `context`, `rationale`, `made_by`, `tags` |
| `pm_search_decisions` | Before re-deciding | `query`, `limit` |
| `pm_record_learning` | Team learned something | `title`, `learning`, `context`, `tags` |
| `pm_team_kb` | Onboarding, Q&A | `question` (AI answers), `board_id` |

#### Blockers
| Tool | When to Call | Params |
|------|-------------|--------|
| `pm_record_blocker` | Team is blocked | `description`, `issue_key`, `owner` |
| `pm_resolve_blocker` | Unblocked | `blocker_id`, `resolution` |
| `pm_blockers` | Standup | `show_history` (bool) |

#### Dependencies
| Tool | When to Call | Params |
|------|-------------|--------|
| `pm_record_dependency` | Cross-team dep found | `from_issue`, `to_issue`, `type` (blocks/blocked_by/external), `description` |
| `pm_resolve_dependency` | Satisfied | `dependency_id` |
| `pm_dependencies` | Planning | `issue_key` (optional filter) |

#### Team
| Tool | When to Call | Params |
|------|-------------|--------|
| `pm_record_team_metric` | End of sprint | `member_name`, `sprint_name`, `issues_assigned`, `issues_done`, `blocker_count`, `carryover_count` |
| `pm_team_health` | Planning | `sprint_name` OR `member_name` |
| `pm_confidence` | Pre-sprint | `sprint_name`, `score` (1-5), `member`, `note` |

#### Retrospectives
| Tool | When to Call | Params |
|------|-------------|--------|
| `pm_record_retro` | After retro | `sprint_name`, `went_well`, `improvements`, `action_items` |
| `pm_action_items` | Every standup | none — shows pending items |

#### Meetings
| Tool | When to Call | Params |
|------|-------------|--------|
| `pm_record_meeting` | After any ceremony | `meeting_type` (standup/planning/retro/grooming/adhoc), `notes`, `decisions`, `action_items`, `attendees` |
| `pm_meetings` | Review history | `meeting_type`, `limit` |

#### Goals
| Tool | When to Call | Params |
|------|-------------|--------|
| `pm_set_sprint_goal` | Sprint planning | `board_id`, `goal`, `key_results`, `sprint_name` |
| `pm_close_sprint_goal` | End of sprint | `goal_id`, `status` (achieved/partially_achieved/missed), `outcome` |
| `pm_sprint_goals` | Review | `board_id`, `show_history` |
| `pm_goal_check` | Mid-sprint | `board_id` — AI evaluates progress |

#### Process Gates
| Tool | When to Call | Params |
|------|-------------|--------|
| `pm_dod` | Manage DoD | `action` (list/add/remove), `item`, `category`, `project` |
| `pm_dor` | Manage DoR | `action` (list/add/remove), `item`, `category`, `project` |
| `pm_agreements` | Team rules | `action` (list/add), `agreement`, `why` |
| `pm_experiment` | Retro improvement | `hypothesis`, `action`, `measure`, `duration` |
| `pm_experiments` | Review experiments | none |

#### Tech Debt
| Tool | When to Call | Params |
|------|-------------|--------|
| `pm_tech_debt_add` | Debt identified | `title`, `description`, `impact` (high/medium/low), `category` (code/architecture/testing/infra/docs), `owner`, `fix_approach` |
| `pm_tech_debt` | Review debt | none |
| `pm_tech_debt_budget` | Sprint planning | `board_id` — recommends allocation % |

---

### AI Intelligence Tools

| Tool | Purpose | When |
|------|---------|------|
| `pm_recommendations` | AI recs from ALL memory | Weekly, planning. Params: `board_id`, `focus` (general/velocity/risks/team/process) |
| `pm_standup_prep` | Talking points | Before standup. Params: `board_id` |
| `pm_forecast` | Monte Carlo "when done?" | Backlog planning. Params: `board_id`, `remaining_items`, `sprint_days` |
| `pm_forecast_sprint` | Sprint completion prob | Planning. Params: `board_id`, `items_remaining` |
| `pm_anti_patterns` | Detect dysfunctions | Monthly. Params: `board_id` |
| `pm_coaching` | Coaching advice | When stuck. Params: `topic`, `board_id`, `situation` |
| `pm_facilitate` | Ceremony prompts | Before ceremony. Params: `ceremony` (standup/planning/retro/grooming/review), `board_id` |
| `pm_retro_analysis` | Pattern detection | Before retro. Params: `board_id` |
| `pm_check_ready` | Story readiness | Grooming. Params: `key` (issue key) |
| `pm_exec_report` | Executive report | Weekly/bi-weekly. Params: `board_id`, `send_to_lark` |
| `pm_flow_metrics` | WIP/throughput/cycle | Mid-sprint. Params: `board_id` |
| `pm_sprint_compare` | This vs last | End of sprint. Params: `board_id` |
| `pm_scope_creep` | Scope change detect | Mid-sprint. Params: `board_id` |
| `jira_ai_analyze` | Ad-hoc analysis | Any time. Params: `query`, `jql`, `max_results` |
| `jira_ai_sprint_report` | Full report | End sprint. Params: `board_id`, `send_to_lark` |

---

### Dashboards & Reports

| Tool | Audience | Content |
|------|----------|---------|
| `pm_dashboard` | SM/PM | Sprint progress, health, risks, blockers, deps, goals, actions |
| `pm_sprint_health` | SM/PM | 0-100 score: velocity(25) + blockers(25) + scope(25) + team(25) |
| `pm_health_history` | SM/PM | Health trend over time |
| `pm_scorecard` | SM/PM | A-F grade: completion, goals, predictability, quality, balance |
| `pm_exec_report` | Executives | Business outcomes, no jargon, 30-second read |
| `pm_weekly_digest` | Team/Stakeholders | AI summary of week's activity |
| `pm_release_notes` | Stakeholders | Categorized: features, bugs, tasks |
| `pm_review_prep` | SM | Demo order, talking points, completion % |
| `pm_planning_prep` | SM | Full prep: capacity, carryover, risks, deps, experiments |
| `pm_mcp_stats` | Admin | Memory contents, data freshness |

---

### Workflow Recipes (One-Click)

| Tool | Action |
|------|--------|
| `pm_recipe_start_work` | Assign + In Progress + branch name. Params: `key`, `assignee_id` |
| `pm_recipe_done` | Transition Done + log time + comment. Params: `key`, `time_spent`, `comment` |
| `pm_recipe_block` | Record blocker + comment on issue. Params: `key`, `reason`, `owner` |

---

### Notifications (Multi-Channel)

| Tool | Channel | Params |
|------|---------|--------|
| `jira_notify_lark` | Lark | `title`, `content` |
| `slack_send` | Slack | `channel`, `message` |
| `slack_notify_team` | Slack | `message`, `channel` |
| `slack_channels` | Slack | none (list channels) |
| `slack_history` | Slack | `channel`, `limit` |
| `discord_send` | Discord | `message` |
| `telegram_send` | Telegram | `message` |
| `teams_send` | Teams | `message` |
| `email_send` | Email | `to`, `subject`, `body` |
| `confluence_create_page` | Confluence | `space`, `title`, `content` |
| `confluence_get_page` | Confluence | `page_id` |
| `confluence_search` | Confluence | `query` |
| `notify_routed` | Smart routing | `message`, `severity` — auto-picks channel |
| `broadcast` | All channels | `message` |
| `daily_digest` | Scheduled | `board_id` — generates + sends digest |

---

### Portfolio Tools

| Tool | Purpose |
|------|---------|
| `portfolio_overview` | Cross-project status summary |
| `portfolio_risks` | Risks across all projects |
| `portfolio_blockers` | Blockers across all projects |
| `portfolio_workload` | Team load across projects |
| `portfolio_summary` | AI executive portfolio summary |

---

## WORKFLOW PATTERNS

### Daily Standup (2 min prep)
```
pm_standup_prep(board_id:X)
```
Returns: talking points, blockers, action items, sprint health signal.

### Sprint Planning (full prep)
```
pm_planning_prep(board_id:X)
```
Returns: last sprint outcome, capacity recommendation, carryover, open risks, dependencies, active experiments, checklist.

### "When Will It Be Done?"
```
pm_forecast(board_id:X, remaining_items:30)
```
Returns: 50%/70%/85%/95% confidence dates from 10,000 Monte Carlo simulations.

### End of Sprint (capture everything)
```
pm_snapshot_sprint(board_id:X, velocity:21, carryover:3)
pm_scorecard(board_id:X)
pm_release_notes(board_id:X, send_to_lark:true)
pm_close_sprint_goal(goal_id:Y, status:"achieved")
```

### Risk Emerges
```
pm_record_risk(title:"API vendor may shutdown", severity:"high", owner:"alice", mitigation:"Build fallback")
```

### Team Seems Off
```
pm_anti_patterns(board_id:X)
pm_coaching(topic:"team_dynamics", situation:"team seems disengaged in retros")
```

### New Member Joins
```
pm_team_kb(question:"how does this team work?")
pm_agreements()
pm_dod()
pm_dor()
```

### Stakeholder Asks for Update
```
pm_exec_report(board_id:X)
```
NOT `pm_dashboard` (too technical for execs).

---

## DATA MODEL

14 SQLite tables, auto-created on first run:

| Table | Key Fields | Purpose |
|-------|-----------|---------|
| sprint_snapshots | sprint_name, board_id, velocity, completion_rate, blocked, carryover | Sprint history |
| daily_progress | sprint_name, date, done, in_progress, todo, blocked | Burndown |
| risks | title, severity, status, owner, mitigation | Risk register (also stores tech debt) |
| decisions | title, decision, rationale, tags | Decision log (also stores agreements, experiments, learnings) |
| blockers | description, blocked_since, resolved_at, resolution | Impediment tracker |
| team_metrics | member_name, sprint_name, issues_assigned, issues_done | Individual tracking |
| retrospectives | sprint_name, went_well, improvements, action_items | Retro outcomes |
| action_items | description, owner, status, due_date | Retro follow-ups |
| dependencies | from_issue_key, to_issue_key, dependency_type, status | Dependency map |
| meeting_notes | meeting_type, date, decisions, action_items | Ceremony outcomes |
| health_scores | overall_score, velocity_score, blocker_score, scope_score, team_score | Health over time |
| sprint_goals | goal, key_results, status, outcome | Goal tracking |
| dod_items | project, item, category, active | DoD + DoR checklists |
| escalations | type, title, severity, channel, acknowledged | Escalation audit |

---

## HEALTH SCORE (0-100)

Computed by `pm_sprint_health`:
- **Velocity (0-25)**: completion rate (done/total)
- **Blockers (0-25)**: fewer blocked = higher score
- **Scope (0-25)**: stable scope vs previous = higher
- **Team (0-25)**: balanced workload = higher

Thresholds: 0-49 AT RISK | 50-69 WATCH | 70-100 HEALTHY

---

## ANTI-PATTERNS DETECTED

`pm_anti_patterns` checks for:
1. **Zombie Sprint** — >30% carryover for 2+ sprints
2. **Hero Culture** — one person completes >50% of items
3. **Scope Creep** — >20% items added mid-sprint
4. **Unpredictable** — >40% velocity coefficient of variation
5. **Dead Retros** — >5 pending action items never done
6. **Rubber-Stamp DoD** — >95% completion with zero blockers (suspiciously perfect)
7. **No Sprint Goals** — no recorded goals exist

---

## MONTE CARLO FORECAST

`pm_forecast` runs 10,000 simulations:
- Samples randomly from historical throughput (items done per sprint)
- Produces probability distribution of completion dates
- **85% confidence** = recommended stakeholder commitment
- Requires minimum 3 sprint snapshots

---

## IMPORTANT NOTES FOR AI AGENTS

1. **First interaction**: Call `jira_boards` to get board_id. Store it — almost every PM tool needs it.
2. **Build memory early**: Encourage `pm_snapshot_sprint` + `pm_record_decision` + `pm_record_risk`. Intelligence tools need data.
3. **Don't overwhelm**: Use `pm_dashboard` for quick checks. Only drill into specific tools when needed.
4. **Proactive scanning**: Run `pm_auto_detect_risks` and `pm_anti_patterns` periodically (weekly).
5. **Escalation threshold**: `pm_escalate` auto-sends to Lark when: critical/high risk >3 days, blocker >3 days, health <50.
6. **Tech debt is stored as risks** with `sprint_name` prefix "tech_debt:category". The `pm_tech_debt*` tools filter for this.
7. **Agreements, experiments, learnings** are stored as decisions with specific tags. The dedicated tools handle filtering.
8. **DoR vs DoD**: DoR items have category prefix "dor:" in the dod_items table.
9. **Confidence votes** are stored as team_metrics with notes prefix "confidence:N".
10. **All AI tools gracefully degrade** — if AI provider is down, they return raw data instead of failing.
