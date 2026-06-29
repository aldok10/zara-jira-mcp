# Skill: zara-jira-mcp

**89-tool MCP server** — AI-powered Scrum Master with persistent memory, Jira Cloud integration, multi-channel notifications, and 12 platform integrations. Written in Go, modular hexagonal architecture.

**Version:** v0.4.0 (modular) | **Transport:** MCP stdio | **Memory:** SQLite WAL

## Connection

Stdio MCP server. Configuration via environment variables only (loaded from `.env` or shell).

**Binary:** `~/.local/bin/zara-jira-mcp` | **Wrapper:** `~/.local/bin/zara-jira-mcp-wrapper.sh` (loads `.env` then exec)

**Key env vars:**
- `JIRA_BASE_URL`, `JIRA_EMAIL`, `JIRA_API_TOKEN` — Jira Cloud auth
- `PM_MEMORY_DB_PATH` — SQLite DB path (default: `~/.zara-jira-mcp/pm_memory.db`)
- `JIRA_AI_BASE_URL`, `JIRA_AI_API_KEY` — AI provider for PM intelligence

## Critical Usage Rules

1. **Get board_id first** — call `jira_boards` before any PM/sprint tool.
2. **First call per session** — `jira_boards` to discover boards. Store the board_id.
3. **Auto-memory is ON** — `jira_search`, `jira_get_issue`, `jira_sprint_summary`, `jira_epic_issues` automatically record blockers, stale risks, and sprint snapshots to PM memory.
4. **Auto-reconciliation is ON** — every read checks stored blockers/risks against current Jira state. Resolved items auto-close.
5. **Board-aware classification** — `jira_sprint_summary` uses board configuration for accurate status mapping (catches custom statuses like "Stalled", "On Hold").
6. **Record immediately** — after any decision, risk, or blocker, record it with `pm_record_decision`, `pm_record_risk`, `pm_record_blocker`.
7. **Use `pm` as primary entry point** — quick project status with blockers, risks, sprint progress.
8. **Forecast needs 3+ snapshots** — Monte Carlo (`pm_forecast`) requires historical sprint data.
9. **Full-sweep reconciliation** — run `pm_reconcile` periodically to sync stored blockers/risks with current Jira state.

---

## TOOL REFERENCE (89 tools)

### Jira Module (28 tools)

#### Search & Read
| Tool | Params | Returns |
|------|--------|---------|
| `jira_search` | `jql` (req), `max_results` | Issues: key, summary, status, assignee |
| `jira_get_issue` | `key` (req) | Full issue details |
| `jira_boards` | *(none)* | Board IDs, names, types |
| `jira_sprints` | `board_id` (req), `state` | Sprints for a board |
| `jira_sprint_summary` | `board_id` (req) | Active sprint breakdown |
| `jira_projects` | *(none)* | All accessible projects |
| `jira_versions` | `project` (req) | Project versions/releases |
| `jira_components` | `project` (req) | Components with leads |
| `jira_board_config` | `board_id` (req) | Column layout + status mappings |
| `jira_transitions` | `key` (req) | Available status transitions |
| `jira_link_types` | *(none)* | Issue link types |
| `jira_worklogs` | `key` (req) | Time entries on issue |
| `jira_attachments` | `key` (req) | Issue attachments |
| `jira_find_user` | `query` (req) | Search users (returns account IDs) |
| `jira_epic_issues` | `epic_key` (req), `max_results` | Issues in an epic |
| `pm_reconcile` | *(none)* | Full-sweep: sync all stored blockers/risks with Jira |

#### Create & Modify
| Tool | Params | Action |
|------|--------|--------|
| `jira_create` | `project` (req), `summary` (req), `issue_type`, `description`, `priority`, `assignee_id`, `labels` | Create issue |
| `jira_transition` | `key` (req), `transition_id` (req) | Change status |
| `jira_assign` | `key` (req), `account_id` (req) | Assign user |
| `jira_add_comment` | `key` (req), `body` (req) | Add comment |
| `jira_add_worklog` | `key` (req), `time_spent` (req), `comment` | Log time |
| `jira_link_issues` | `inward_key` (req), `outward_key` (req), `link_type` (req) | Create issue link |
| `jira_epic_add` | `epic_key` (req), `issue_keys` (req) | Add to epic |
| `jira_epic_remove` | `issue_keys` (req) | Remove from epic |
| `jira_version_create` | `project` (req), `name` (req), `description` | Create version |
| `jira_version_release` | `version_id` (req) | Release version |

#### Sprint Management
| Tool | Params | Action |
|------|--------|--------|
| `jira_start_sprint` | `sprint_id` (req), `start_date` (req), `end_date` (req) | Start sprint |
| `jira_move_to_sprint` | `sprint_id` (req), `issue_keys` (req) | Move issues to sprint |

---

### Sprint/PM Module (18 tools)

| Tool | Params | Description |
|------|--------|-------------|
| `pm` | `board_id` | Quick project status: sprint, blockers, risks, actions |
| `pm_create` | `title` (req), `description`, `project`, `assignee`, `type`, `priority`, `labels`, `platform` | Create work item (Jira/GitHub/GitLab) |
| `pm_decide` | `what` (req), `who`, `why` | Quick decision record |
| `pm_risk` | `what` (req), `severity`, `owner` | Quick risk record |
| `pm_next` | `board_id` | Suggest next PM action from memory state |
| `pm_snapshot` | `board_id`, `sprint_name`, `total_issues`, `done`, `in_progress`, `todo`, `blocked`, `carryover`, `velocity`, `notes` | Sprint snapshot to memory |
| `pm_record_decision` | `title` (req), `decision` (req), `context`, `rationale`, `made_by`, `tags` | Decision with rationale |
| `pm_record_risk` | `title` (req), `description`, `severity`, `owner`, `mitigation` | Risk with mitigation |
| `pm_record_blocker` | `description` (req), `issue_key`, `owner` | Impediment record |
| `pm_record_retro` | `sprint_name` (req), `went_well`, `improvements`, `action_items` | Retrospective |
| `pm_record_meeting` | `meeting_type` (req), `notes`, `attendees`, `decisions`, `action_items`, `sprint_name` | Meeting notes |
| `pm_risks` | *(none)* | Risk dashboard |
| `pm_blockers` | `show_history` | Active blocker list |
| `pm_decisions` | `limit` | Recent decisions |
| `pm_actions` | *(none)* | Pending action items |
| `pm_dependencies` | `issue_key` | Dependency map |
| `pm_health` | `board_id` | Sprint health history |
| `pm_forecast` | `board_id`, `remaining_items` | Monte Carlo 50/70/85/95% completion forecast |

---

### Notification Module (5 tools)

| Tool | Params | Channel |
|------|--------|---------|
| `jira_notify_lark` | `content` (req), `title` | Lark group |
| `jira_notify_slack` | `content` (req), `title`, `channel` | Slack |
| `jira_notify_discord` | `content` (req), `title`, `channel` | Discord |
| `jira_notify_telegram` | `content` (req), `title`, `chat_id` | Telegram |
| `notify_routed` | `content` (req), `severity`, `audience`, `title` | Auto-picks channel by severity |

---

### GitHub Module (10 tools)

| Tool | Params | Description |
|------|--------|-------------|
| `pm_github_prs` | `state`, `limit` | Open PRs with age, review status |
| `pm_github_pr_metrics` | `stale_days` | Avg PR age, stale count |
| `pm_github_activity` | `days` | Commits + PRs merged + issues closed |
| `pm_github_releases` | `limit` | Recent releases/tags |
| `pm_github_search_branches` | `pattern` (req) | Search branches by pattern |
| `pm_github_search_prs_by_branch` | `branch` (req) | Find PRs for a branch |
| `pm_github_issues` | `state`, `labels`, `limit` | List GitHub issues |
| `pm_github_create_issue` | `title` (req), `body`, `labels`, `assignees`, `milestone` | Create GitHub issue |
| `pm_github_milestones` | `state` | Milestones with progress |
| `pm_github_create_milestone` | `title` (req), `description`, `due_date` | Create milestone |

---

### GitLab Module (9 tools)

| Tool | Params | Description |
|------|--------|-------------|
| `pm_gitlab_issues` | `state`, `labels`, `limit` | List GitLab issues |
| `pm_gitlab_create_issue` | `title` (req), `description`, `labels`, `assignee_id`, `milestone_id` | Create GitLab issue |
| `pm_gitlab_merge_requests` | `state`, `limit` | List merge requests |
| `pm_gitlab_milestones` | `state` | List milestones |
| `pm_gitlab_create_milestone` | `title` (req), `description`, `due_date` | Create milestone |
| `pm_gitlab_search_branches` | `pattern` (req) | Search branches |
| `pm_gitlab_search_mrs_by_branch` | `branch` (req) | Find MRs by branch |
| `pm_gitlab_read_file` | `path` (req), `ref` | Read file from repo |
| `pm_gitlab_list_files` | `path`, `ref` | List directory contents |

---

### Calendar (Lark) — 3 tools

| Tool | Params | Description |
|------|--------|-------------|
| `pm_calendar_create` | `summary` (req), `start` (req), `end`, `description` | Create Lark event |
| `pm_calendar_events` | `days` | List upcoming events |
| `pm_calendar_schedule_meeting` | `summary` (req), `start` (req), `description`, `duration_minutes`, `attendees` | Schedule with VC |

### Clockify — 2 tools

| Tool | Params | Description |
|------|--------|-------------|
| `pm_time_entries` | `days` | Recent time entries |
| `pm_time_report` | `days` | Time tracked per person/project |

### Confluence — 3 tools

| Tool | Params | Description |
|------|--------|-------------|
| `pm_confluence_search` | `query` (req), `limit` | Search pages by CQL |
| `pm_confluence_get_page` | `page_id` (req) | Get page content |
| `pm_confluence_create_page` | `space_key` (req), `title` (req), `body`, `parent_id` | Create page |

### Linear — 3 tools

| Tool | Params | Description |
|------|--------|-------------|
| `pm_linear_issues` | `team`, `state` | List issues |
| `pm_linear_cycles` | *(none)* | Current/recent cycles |
| `pm_linear_activity` | *(none)* | Recent changes |

### Notion — 3 tools

| Tool | Params | Description |
|------|--------|-------------|
| `pm_notion_search` | `query` (req), `limit` | Search pages/databases |
| `pm_notion_create_page` | `title` (req), `content`, `parent_id` | Create page |
| `pm_notion_query_db` | `database_id`, `filter`, `limit` | Query database |

### PagerDuty — 2 tools

| Tool | Params | Description |
|------|--------|-------------|
| `pm_incidents` | `status` | Incidents with severity |
| `pm_oncall` | *(none)* | Who's on call now |

### Google Sheets — 1 tool

| Tool | Params | Description |
|------|--------|-------------|
| `pm_sheet_read` | `spreadsheet_id` (req), `range` | Read sheet data |

### Backup & Onboard — 2 tools

| Tool | Params | Description |
|------|--------|-------------|
| `pm_backup` | *(none)* | Export PM memory to JSON |
| `pm_onboard` | *(none)* | First-run config wizard |

---

## AUTO-MEMORY SYSTEM

The server records and reconciles PM memory automatically on data access:

### Auto-Record (on read)
| Trigger | What's Recorded |
|---------|-----------------|
| `jira_search` | Blockers (status=blocked), stale risks (Highest/Critical, no update >7d) |
| `jira_get_issue` | Blockers + stale risks for single issue |
| `jira_sprint_summary` | Blockers + stale risks + sprint snapshot (done/progress/blocked/todo counts) |
| `jira_epic_issues` | Blockers + stale risks for epic |

### Auto-Reconciliation (on every read)
| Trigger | What's Resolved |
|---------|-----------------|
| Any read | Blockers whose issues are no longer blocked → auto-resolve |
| Any read | Risks whose issues are done → auto-resolve |
| Any read | Risks whose issues recently updated → auto-mitigate |

### Full-Sweep
`pm_reconcile` fetches every stored issue key and reconciles all at once.

### Deduplication
- Blockers: deduped by `issue_key` (only one active blocker per issue)
- Risks: deduped by issue key prefix in title (e.g., "Stale: PROJ-123 — ...")

---

## BOARD-AWARE CLASSIFICATION

The server uses board configuration (`jira_board_config`) for accurate status classification:

| Category | Description |
|----------|-------------|
| `todo` | Column names like "To Do", "Backlog", "Open", or unmapped |
| `progress` | Column names containing "Progress", "Review", "Testing", "Dev" |
| `blocked` | Column names containing "Blocked", "Stalled", "Waiting", "Impediment" |
| `done` | Column names containing "Done", "Closed", "Resolved", "Complete" |

**Fallback:** When board config isn't available (e.g., search without board context), uses heuristic string matching on status names.

**Caching:** Board configurations are cached per board ID in memory for the server lifetime.

---

## DATA MODEL

16+ SQLite tables, auto-created on first run at `~/.zara-jira-mcp/pm_memory.db`:

| Table | Purpose |
|-------|---------|
| `sprint_snapshots` | Sprint history (velocity, completion, carryover) |
| `daily_progress` | Burndown data |
| `risks` | Risk register |
| `decisions` | Decision log + agreements + experiments + learnings |
| `blockers` | Impediment tracker |
| `team_metrics` | Individual sprint metrics |
| `retrospectives` | Retro outcomes |
| `action_items` | Retro follow-ups |
| `dependencies` | Dependency map |
| `meeting_notes` | Ceremony outcomes |
| `health_scores` | Health over time |
| `sprint_goals` | Goal tracking |
| `dod_items` | DoD + DoR checklists |
| `escalations` | Escalation audit trail |
| `coaching` | Coaching records |
| `okrs` | OKR tracking |

Additional advanced tables: `kb_articles`, `feedback_log`, `experiments`, `learnings`, `pulse_surveys`, `radar_dimensions`, `safety_surveys`, `stakeholder_pulse`, `kpi_definitions`, `kpi_measurements`, `key_results`, `kr_issues`, `okr_link`.

---

## WORKFLOW PATTERNS

### Daily Standup
```
1. pm board_id=X → quick sprint status
2. pm_blockers → active impediments
3. pm_risks → risk dashboard
4. pm_github_prs → PRs needing review
```

### Sprint Planning
```
1. pm_health board_id=X → previous sprint history
2. pm_forecast board_id=X → Monte Carlo forecast
3. pm_next board_id=X → AI-suggested action
```

### Sprint End
```
1. pm_snapshot board_id=X → save to memory
2. pm_record_retro sprint_name="Sprint X" → capture retro outcomes
3. pm_reconcile → sync blockers/risks
```

### Decision Record
```
1. pm_record_decision title="..." decision="..." context="..."
2. pm_decisions → browse decision log
```

### Weekly Review
```
1. pm_github_activity days=7 → dev activity
2. pm_time_report days=7 → hours tracked
3. pm_risks → open risks
4. pm_actions → pending retro actions
```

---

## IMPORTANT NOTES FOR AI AGENTS

1. **First call**: `jira_boards` to get board_id. Store it.
2. **Board config**: After getting board_id, optionally call `jira_board_config board_id=X` — this populates the board-aware classifier cache for more accurate sprint reporting.
3. **Auto-memory reduces manual work**: No need to manually `pm_record_blocker` for every blocked issue — auto-memory catches them. Still record blockers for non-Jira reasons.
4. **Reconciliation is NOT background**: It runs synchronously on every read. For large datasets, this adds a small latency (individual issue fetches).
5. **Run `pm_reconcile` periodically** — especially if you've changed statuses/triaged issues outside the server.
6. **Forecast needs 3+ historical snapshots** — call `pm_snapshot` at sprint end.
7. **Board-aware classification > heuristic** — `jira_sprint_summary` is more accurate than `jira_search` for status categorization.
8. **Tools NOT yet wired** (from monolithic): bulk operations, update/delete issue, portfolio, AI analysis, stakeholder reporting, most communication intelligence tools. The modular codebase has 89 tools; the old monolithic had ~279.
9. **Migration status**: Modular codebase at `apps/api/` with 89 tools. Monolithic binary still installed but being replaced.

---

## SETUP PER AGENT PLATFORM

### OpenCode
```json
{
  "mcp": {
    "jira-pm": {
      "type": "stdio",
      "command": "zara-jira-mcp",
      "env": {
        "JIRA_BASE_URL": "...",
        "JIRA_EMAIL": "...",
        "JIRA_API_TOKEN": "...",
        "PM_MEMORY_DB_PATH": "~/.zara-jira-mcp/pm_memory.db"
      }
    }
  }
}
```

### Claude Code / Claude Desktop
```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "~/.local/bin/zara-jira-mcp",
      "env": { "JIRA_BASE_URL": "...", "JIRA_EMAIL": "...", "JIRA_API_TOKEN": "...", "PM_MEMORY_DB_PATH": "~/.zara-jira-mcp/pm_memory.db" }
    }
  }
}
```

### Cursor / VS Code / Kiro / Zed
Use the same format as above with the platform's MCP configuration file.
