# Skill: zara-jira-mcp

239-tool MCP server: AI-powered Scrum Master with persistent memory, Jira Cloud integration, multi-channel notifications, and 16 platform integrations.

## Connection

Stdio MCP server. SQLite persistent memory at `~/.zara-jira-mcp/pm_memory.db`. Config file at `~/.zara-jira-mcp/config.json`.

## Profile System

The server loads tools based on `PM_PROFILE` env var:

| Profile | Tools | Modules |
|---------|-------|---------|
| `chatgpt` | ~14 | shortcuts only (smart routing) |
| `lite` | ~30 | shortcuts + pm |
| `standard` | ~80 | jira + pm + ai + notifications + shortcuts |
| `full` | ~150 | standard + stakeholder + portfolio + github |
| (none/all) | ~239 | everything |

Custom: `PM_ENABLED_MODULES=jira,pm,ai,notifications` (comma-separated).

Module names: `jira`, `pm`, `ai`, `notifications`, `stakeholder`, `portfolio`, `github`, `integrations`, `shortcuts`

## Critical Rules

1. **Get board_id first** — call `jira_boards` before any PM tool.
2. **Memory builds over time** — intelligence tools need historical data. Record snapshots, decisions, risks.
3. **`pm_snapshot_sprint` at end of EVERY sprint** — forecasting/velocity require it.
4. **Never `pm_dashboard` for executives** — use `pm_exec_report` instead.
5. **Record immediately** — after any decision, risk, or blocker.
6. **All AI tools gracefully degrade** — if AI provider is down, they return raw data.

---

## TOOL REFERENCE (239 tools)

### Jira Operations (~48 tools)

#### Search & Read
| Tool | Params | Returns |
|------|--------|---------|
| `jira_search` | `jql` (required), `max_results`, `start_at` | Issues: key, summary, status, priority, assignee |
| `jira_get_issue` | `key` | Full issue details |
| `jira_boards` | none | Board IDs (needed for PM tools) |
| `jira_sprint_summary` | `board_id` | Active sprint status breakdown |
| `jira_sprints` | `board_id`, `state` | All sprints for a board |
| `jira_my_issues` | `status` | Current user's unresolved issues |
| `jira_overdue` | `days`, `project` | Issues with no update in N days |
| `jira_workload` | `project` | Issue count per assignee |
| `jira_projects` | none | All accessible projects |
| `jira_project_detail` | `key` | Project metadata |
| `jira_transitions` | `key` | Available status transitions |
| `jira_link_types` | none | Available issue link types |
| `jira_watchers` | `key` | Issue watchers |
| `jira_worklog_list` | `key` | Time entries |
| `jira_find_user` | `query` | Search users (returns account IDs) |
| `jira_epic_issues` | `epic_key` | Issues in an epic |
| `jira_health` | none | Server health check |
| `jira_versions` | `project` | Project versions/releases |

#### Create & Modify
| Tool | Params | Action |
|------|--------|--------|
| `jira_create_issue` | `project`, `summary`, `issue_type`, `description`, `priority`, `assignee_id`, `labels` | Create issue |
| `jira_update_issue` | `key`, fields | Update fields |
| `jira_add_comment` | `key`, `body` | Add comment |
| `jira_transition` | `key`, `transition_id` | Change status |
| `jira_assign` | `key`, `account_id` | Assign user |
| `jira_unassign` | `key` | Remove assignee |
| `jira_delete_issue` | `key` | Delete permanently |
| `jira_create_subtask` | `parent_key`, `project`, `summary` | Create subtask |
| `jira_labels_set` | `key`, `labels` | Set labels |
| `jira_link_issues` | `from_key`, `to_key`, `link_type` | Create link |
| `jira_watch` | `key`, `account_id` | Add watcher |
| `jira_worklog_add` | `key`, `time_spent`, `comment` | Log time |
| `jira_epic_add` | `issue_key`, `epic_key` | Add to epic |
| `jira_epic_remove` | `issue_key` | Remove from epic |
| `jira_version_create` | `project`, `name` | Create version |
| `jira_version_release` | `version_id` | Release version |

#### Sprint Management
| Tool | Params | Action |
|------|--------|--------|
| `jira_sprint_create` | `board_id`, `name`, `goal` | Create sprint |
| `jira_sprint_start` | `sprint_id`, `start_date`, `end_date` | Start sprint |
| `jira_sprint_close` | `sprint_id` | Close sprint |
| `jira_sprint_move_issues` | `sprint_id`, `issue_keys` | Move issues to sprint |

#### Bulk Operations
| Tool | Params | Action |
|------|--------|--------|
| `jira_bulk_transition` | `keys`, `transition_id` | Transition multiple |
| `jira_bulk_assign` | `keys`, `account_id` | Assign multiple |
| `jira_bulk_label` | `keys`, `label` | Label multiple |

#### Advanced
| Tool | Params | Action |
|------|--------|--------|
| `jira_raw_request` | `method`, `path`, `body` | Raw REST API call |
| `jira_from_branch` | `branch` | Extract issue key from branch |
| `jira_smart_commit` | `key`, `message` | Parse commit commands |
| `jira_link_pr` | `key`, `url`, `title` | Link PR to issue |
| `jira_trace_branch` | `key` | Trace issue through git/deploy pipeline |

---

### PM Memory (~22 tools)

#### Sprint Data
| Tool | When | Params |
|------|------|--------|
| `pm_snapshot_sprint` | End of sprint | `board_id`, `velocity`, `carryover`, `notes` |
| `pm_track_daily` | Daily | `board_id` |
| `pm_burndown` | View chart | `board_id`, `sprint_name` |
| `pm_velocity_trend` | Planning | `board_id` |

#### Risks
| Tool | When | Params |
|------|------|--------|
| `pm_record_risk` | Risk found | `title`, `severity`, `owner`, `mitigation` |
| `pm_update_risk` | Status change | `risk_id`, `status` |
| `pm_risk_dashboard` | Standup/planning | none |
| `pm_auto_detect_risks` | Weekly | `board_id` |

#### Decisions & Knowledge
| Tool | When | Params |
|------|------|--------|
| `pm_record_decision` | After decision | `title`, `decision`, `context`, `rationale`, `made_by`, `tags` |
| `pm_search_decisions` | Before re-deciding | `query`, `limit` |
| `pm_record_learning` | Team learned something | `title`, `learning`, `context`, `tags` |
| `pm_team_kb` | Onboarding/Q&A | `question`, `board_id` |

#### Blockers
| Tool | When | Params |
|------|------|--------|
| `pm_record_blocker` | Team blocked | `description`, `issue_key`, `owner` |
| `pm_resolve_blocker` | Unblocked | `blocker_id`, `resolution` |
| `pm_blockers` | Standup | `show_history` |

#### Dependencies
| Tool | When | Params |
|------|------|--------|
| `pm_record_dependency` | Cross-team dep | `from_issue`, `to_issue`, `type`, `description` |
| `pm_resolve_dependency` | Satisfied | `dependency_id` |
| `pm_dependencies` | Planning | `issue_key` |

#### Team & Retros
| Tool | When | Params |
|------|------|--------|
| `pm_record_team_metric` | End sprint | `member_name`, `sprint_name`, `issues_assigned`, `issues_done` |
| `pm_team_health` | Planning | `sprint_name` or `member_name` |
| `pm_confidence` | Pre-sprint | `sprint_name`, `score`, `member`, `note` |
| `pm_record_retro` | After retro | `sprint_name`, `went_well`, `improvements`, `action_items` |
| `pm_action_items` | Standup | none |
| `pm_record_meeting` | After ceremony | `meeting_type`, `notes`, `decisions`, `action_items` |

#### Goals & Process
| Tool | When | Params |
|------|------|--------|
| `pm_set_sprint_goal` | Sprint planning | `board_id`, `goal`, `key_results` |
| `pm_close_sprint_goal` | End sprint | `goal_id`, `status`, `outcome` |
| `pm_sprint_goals` | Review | `board_id`, `show_history` |
| `pm_dod` | Manage DoD | `action`, `item`, `category`, `project` |
| `pm_dor` | Manage DoR | `action`, `item`, `category`, `project` |
| `pm_agreements` | Team rules | `action`, `agreement`, `why` |
| `pm_experiment` | Improvement | `hypothesis`, `action`, `measure`, `duration` |

---

### AI Intelligence (~15 tools)

| Tool | Purpose | Key Params |
|------|---------|------------|
| `pm_recommendations` | AI recs from ALL memory | `board_id`, `focus` |
| `pm_standup_prep` | Daily talking points | `board_id` |
| `pm_forecast` | Monte Carlo "when done?" | `board_id`, `remaining_items` |
| `pm_forecast_sprint` | Sprint completion probability | `board_id`, `items_remaining` |
| `pm_anti_patterns` | Detect dysfunctions | `board_id` |
| `pm_coaching` | Coaching advice | `topic`, `situation`, `board_id` |
| `pm_facilitate` | Ceremony prompts | `ceremony`, `board_id` |
| `pm_retro_analysis` | Pattern detection | `board_id` |
| `pm_check_ready` | Story readiness | `key` |
| `pm_exec_report` | Executive report | `board_id`, `send_to_lark` |
| `pm_flow_metrics` | WIP/throughput/cycle time | `board_id` |
| `pm_sprint_compare` | This vs last sprint | `board_id` |
| `pm_scope_creep` | Scope change detection | `board_id` |
| `jira_ai_analyze` | Ad-hoc analysis | `query`, `jql` |
| `jira_ai_sprint_report` | Full sprint report | `board_id`, `send_to_lark` |

---

### Stakeholder & Reporting (~18 tools)

| Tool | Audience | Purpose |
|------|----------|---------|
| `pm_dashboard` | SM/PM | One-shot sprint view |
| `pm_sprint_health` | SM/PM | 0-100 health score |
| `pm_health_history` | SM/PM | Health trend |
| `pm_scorecard` | SM/PM | A-F sprint grade |
| `pm_exec_report` | Executives | Business outcomes, no jargon |
| `pm_weekly_digest` | Stakeholders | AI weekly summary |
| `pm_release_notes` | Stakeholders | What shipped |
| `pm_review_prep` | SM | Demo prep |
| `pm_planning_prep` | SM | Full planning package |
| `pm_stakeholder_pulse` | PM | Track satisfaction |
| `pm_stakeholder_trend` | PM | Relationship trajectory |
| `pm_sm_impact` | SM's manager | Prove SM value |
| `pm_escalate` | Auto | Alert on chronic issues |
| `pm_maturity_assessment` | Eng leadership | Team stage |
| `pm_outcome_map` | PO/leadership | Sprint-to-OKR mapping |
| `pm_mcp_stats` | Admin | Memory contents |
| `pm_tech_debt` | Team | Debt inventory |
| `pm_tech_debt_budget` | Planning | Allocation % recommendation |

---

### Communication Intelligence (~14 tools)

| Tool | Framework | Purpose |
|------|-----------|---------
| `pm_comms_health` | Signal-over-Noise | Communication health score 0-100 (decision velocity, blocker resolution, action follow-through, engagement) |
| `pm_comms_anti_patterns` | Community Smells | Detect: re-deciding, escalation hoarding, ghost stakeholders, blocker silence |
| `pm_silence_detector` | Ghost Stakeholders | Find stakeholders with no recent activity |
| `pm_lencioni` | 5 Dysfunctions | Map team data to Lencioni pyramid levels with coaching per level |
| `pm_trust_signals` | Trust Pyramid | Dashboard: forecast accuracy, escalation responsiveness, consistency, transparency |
| `pm_hard_conversation` | Crucial Conversations | Prep difficult conversations: facts + stories + SCARF risks + opening lines |
| `pm_nvc_reframe` | NVC (Rosenberg) | Rewrite blaming language into Observation/Feeling/Need/Request |
| `pm_communicate` | Minto Pyramid | Generate or rewrite audience-specific messages (also absorbs pm_compose) |
| `pm_feedback_prep` | SBI + data | AI-generate structured feedback from team metrics |
| `pm_escalation_draft` | BLUF + TIRED | Draft escalation: ask + context + impact + next step + deadline |
| `pm_decision_record` | ADR/MADR | Enhanced decision record with alternatives + consequences |
| `pm_status_draft` | Pyramid + data | Auto-pull sprint data, format for audience |
| `pm_announce_decision` | DACI | Communicate decisions with roles and rationale |
| `pm_comms_plan` | RACI + Timing | Generate communication plan: who, what, when, channel |

---

### Notifications (~15 tools)

| Tool | Channel |
|------|---------|
| `jira_notify_lark` | Lark |
| `slack_send` | Slack |
| `slack_notify_team` | Slack |
| `slack_channels` | Slack (list) |
| `slack_history` | Slack (read) |
| `discord_send` | Discord |
| `telegram_send` | Telegram |
| `teams_send` | Teams |
| `email_send` | Email |
| `confluence_create_page` | Confluence |
| `confluence_get_page` | Confluence |
| `confluence_search` | Confluence |
| `notify_routed` | Smart routing (auto-picks channel by severity) |
| `broadcast` | All channels |
| `daily_digest` | Generates + sends digest |

---

### GitHub/GitLab (~12 tools)

| Tool | Purpose |
|------|---------|
| `pm_github_prs` | Open PRs with age |
| `pm_github_pr_metrics` | PR review health |
| `pm_github_activity` | Repo activity |
| `pm_github_releases` | Release history |
| `github_create_issue` | Create GitHub issue |
| `github_create_pr` | Create pull request |
| `github_repo_info` | Repository metadata |
| `github_actions` | Workflow status |
| `jira_from_branch` | Extract issue from branch |
| `jira_link_pr` | Link PR to Jira issue |
| `jira_trace_branch` | Trace through pipeline |
| `jira_smart_commit` | Parse commit commands |

---

### External Integrations (~12 tools)

| Tool | Service |
|------|---------|
| `pm_calendar_today` | Google Calendar — today's events |
| `pm_calendar_week` | Google Calendar — week view |
| `pm_calendar_create` | Google Calendar — create event |
| `notion_query` | Notion — query database |
| `notion_create_page` | Notion — create page |
| `linear_issues` | Linear — query issues |
| `linear_create_issue` | Linear — create issue |
| `pagerduty_incidents` | PagerDuty — active incidents |
| `pagerduty_oncall` | PagerDuty — on-call schedule |
| `clockify_time_entries` | Clockify — time entries |
| `clockify_log_time` | Clockify — log time |
| `sheets_read` | Google Sheets — read data |

---

### Shortcuts & Smart Context (~14 tools)

| Tool | Purpose |
|------|---------|
| `pm` | One-shot sprint status (alias for dashboard) |
| `pm_next` | AI suggests next action |
| `pm_help` | Topic-based tool discovery |
| `pm_quickstart` | First-time setup guide |
| `pm_workflow` | Step-by-step workflow sequences |
| `pm_decide` | Quick decision record (simplified) |
| `pm_risk` | Quick risk record (simplified) |
| `pm_create` | Quick issue create (simplified) |
| `pm_context` | Smart context for current sprint |
| `pm_nl_to_jql` | Natural language to JQL |

---

### Workflow Recipes (3 tools)

| Tool | Action |
|------|--------|
| `pm_recipe_start_work` | Assign + In Progress + branch name |
| `pm_recipe_done` | Transition Done + log time + comment |
| `pm_recipe_block` | Record blocker + comment on issue |

---

### Portfolio (5 tools)

| Tool | Purpose |
|------|---------|
| `portfolio_overview` | Cross-project status |
| `portfolio_risks` | Risks across all projects |
| `portfolio_blockers` | Blockers across all projects |
| `portfolio_workload` | Team load across projects |
| `portfolio_summary` | AI executive portfolio summary |

---

## WORKFLOW PATTERNS

### Daily Standup
```
1. pm_calendar_today → today's meetings
2. pm_standup_prep board_id=X → AI talking points
3. pm_blockers → active impediments
4. pm_incidents → production issues
5. pm_github_prs → PRs needing review
```

### Sprint Start
```
1. pm_planning_prep board_id=X → velocity, capacity, carryover
2. pm_set_sprint_goal board_id=X goal="..." → define goal
3. pm_check_ready key=ISSUE-123 → verify readiness (per story)
4. pm_capacity_plan board_id=X team_size=N → recommended points
5. pm_confidence sprint_name="Sprint X" score=4 → team confidence
```

### Sprint End
```
1. pm_sprint_health board_id=X → health score
2. pm_scorecard board_id=X → sprint grade
3. pm_snapshot_sprint board_id=X → save to memory
4. pm_release_notes board_id=X → categorized release notes
5. pm_record_retro sprint_name="Sprint X" → capture retro
6. pm_close_sprint_goal goal_id=X status=achieved → close goal
7. pm_exec_report board_id=X → executive summary
```

### Sprint Planning
```
1. pm_planning_prep board_id=X → history + carryover + risks + deps
2. pm_backlog_groom project=PROJ → stale items
3. pm_capacity_plan board_id=X team_size=N sprint_days=10 → capacity
4. pm_check_ready key=ISSUE-123 → story readiness (per candidate)
5. pm_forecast board_id=X remaining_items=N → "when done?"
6. pm_set_sprint_goal board_id=X goal="..." → define goal
```

### Retrospective
```
1. pm_facilitate ceremony=retro → fresh format
2. pm_sprint_compare board_id=X → this vs last
3. pm_anti_patterns board_id=X → detected issues
4. pm_retro_analysis board_id=X → AI patterns
5. pm_record_retro sprint_name="Sprint X" → save outcomes
6. pm_experiment hypothesis="..." action="..." → track improvement
```

### Incident Response
```
1. pm_incidents → open incidents
2. pm_incident_summary → sprint impact
3. pm_oncall → who's handling
4. pm_auto_detect_risks board_id=X → scan for risks
5. pm_record_risk title="..." severity=high → log risk
6. pm_escalate board_id=X → auto-escalate critical items
```

### Weekly Review
```
1. pm_weekly_digest board_id=X → AI weekly summary
2. pm_github_activity days=7 → repo activity
3. pm_time_report range=week → hours tracked
4. pm_risk_dashboard → open risks
5. pm_action_items → pending retro actions
6. pm_experiments → active experiments
7. pm_health_history board_id=X → health trend
```

### "When Will It Be Done?"
```
pm_forecast board_id=X remaining_items=30
→ 50%/70%/85%/95% confidence dates from 10,000 Monte Carlo simulations
```

### Executive Update
```
pm_exec_report board_id=X
→ Business outcomes, no jargon, 30-second read. NOT pm_dashboard.
```

---

## DATA MODEL

14 SQLite tables, auto-created on first run:

| Table | Purpose |
|-------|---------|
| sprint_snapshots | Sprint history (velocity, completion, carryover) |
| daily_progress | Burndown data |
| risks | Risk register + tech debt |
| decisions | Decision log + agreements + experiments + learnings |
| blockers | Impediment tracker |
| team_metrics | Individual sprint metrics |
| retrospectives | Retro outcomes |
| action_items | Retro follow-ups |
| dependencies | Dependency map |
| meeting_notes | Ceremony outcomes |
| health_scores | Health over time |
| sprint_goals | Goal tracking |
| dod_items | DoD + DoR checklists |
| escalations | Escalation audit trail |

---

## SETUP PER AGENT PLATFORM

### Claude Code / Claude Desktop
```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "zara-jira-mcp",
      "env": { "JIRA_BASE_URL": "...", "JIRA_EMAIL": "...", "JIRA_API_TOKEN": "...", "PM_PROFILE": "standard" }
    }
  }
}
```

### OpenCode
In `opencode.json`:
```json
{
  "mcp": {
    "jira-pm": {
      "type": "stdio",
      "command": "zara-jira-mcp",
      "env": { "JIRA_BASE_URL": "...", "JIRA_EMAIL": "...", "JIRA_API_TOKEN": "...", "PM_PROFILE": "standard" }
    }
  }
}
```

### Cursor
In `.cursor/mcp.json`:
```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "zara-jira-mcp",
      "env": { "JIRA_BASE_URL": "...", "JIRA_EMAIL": "...", "JIRA_API_TOKEN": "...", "PM_PROFILE": "standard" }
    }
  }
}
```

### VS Code + Copilot
In `.vscode/mcp.json`:
```json
{
  "servers": {
    "jira-pm": {
      "type": "stdio",
      "command": "zara-jira-mcp",
      "env": { "JIRA_BASE_URL": "...", "JIRA_EMAIL": "...", "JIRA_API_TOKEN": "...", "PM_PROFILE": "standard" }
    }
  }
}
```

### Kiro
In `.kiro/mcp.json`:
```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "zara-jira-mcp",
      "env": { "JIRA_BASE_URL": "...", "JIRA_EMAIL": "...", "JIRA_API_TOKEN": "...", "PM_PROFILE": "standard" }
    }
  }
}
```

### ChatGPT Desktop
```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "zara-jira-mcp",
      "env": { "JIRA_BASE_URL": "...", "JIRA_EMAIL": "...", "JIRA_API_TOKEN": "...", "PM_PROFILE": "chatgpt" }
    }
  }
}
```
Use `PM_PROFILE=chatgpt` for ChatGPT (14 tools, smart routing handles the rest).

### Zed
In `.zed/settings.json`:
```json
{
  "context_servers": {
    "jira-pm": {
      "command": { "path": "zara-jira-mcp", "args": [] }
    }
  }
}
```

### Goose
```yaml
extensions:
  jira-pm:
    type: stdio
    command: zara-jira-mcp
    enabled: true
```

### Gemini CLI
```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "zara-jira-mcp",
      "env": { "JIRA_BASE_URL": "...", "JIRA_EMAIL": "...", "JIRA_API_TOKEN": "...", "PM_PROFILE": "standard" }
    }
  }
}
```

---

## FIRST-TIME USAGE

1. Call `jira_boards` → get your board_id
2. Call `pm_quickstart` → guided setup walkthrough
3. Call `pm_standup_prep board_id=X` → verify everything works
4. Start recording: snapshots, decisions, risks
5. After 3+ sprints, forecasting and pattern detection become accurate

---

## IMPORTANT NOTES FOR AI AGENTS

1. **First interaction**: `jira_boards` to get board_id. Store it.
2. **Build memory early**: Encourage `pm_snapshot_sprint` + `pm_record_decision` + `pm_record_risk`.
3. **Don't overwhelm**: Use `pm` for quick checks. Drill into specific tools when needed.
4. **Proactive scanning**: Run `pm_auto_detect_risks` and `pm_anti_patterns` weekly.
5. **Escalation threshold**: `pm_escalate` auto-sends when: critical/high risk >3 days, blocker >3 days, health <50.
6. **Tech debt is stored as risks** with `sprint_name` prefix "tech_debt:category".
7. **Agreements, experiments, learnings** are stored as decisions with specific tags.
8. **DoR vs DoD**: DoR items have category prefix "dor:" in dod_items table.
9. **Confidence votes** are stored as team_metrics with notes prefix "confidence:N".
10. **Monte Carlo requires 3+ sprint snapshots** before producing useful forecasts.
11. **Health score 0-100**: velocity(25) + blockers(25) + scope(25) + team(25). Thresholds: 0-49 AT RISK | 50-69 WATCH | 70-100 HEALTHY.
12. **Anti-patterns detected**: zombie sprint, hero culture, scope creep, unpredictable velocity, dead retros, rubber-stamp DoD, no sprint goals.
13. **Smart routing (`notify_routed`)**: auto-picks notification channel based on severity and configured platforms.
14. **Profile affects tool availability** — if a tool isn't found, suggest changing PM_PROFILE.
