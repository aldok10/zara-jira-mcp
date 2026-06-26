# ROADMAP: zara-jira-mcp

> Transform zara-jira-mcp from a basic MCP wrapper into a go-jira-class Jira powerhouse
> with AI-native PM intelligence — the tool a great Scrum Master would build for themselves.

---

## Vision

**go-jira gave developers CLI speed. We give AI agents PM intelligence.**

go-jira solved "Jira is slow, let me use terminal." We solve "Jira is dumb, let me have
an AI PM partner that understands flow, predicts risks, and acts on my behalf."

Our competitive advantage:
1. **MCP-native** — works with any AI agent (Claude, OpenCode, Cursor, etc.)
2. **PM intelligence** — not just CRUD, but insights (flow metrics, risk detection, workload balance)
3. **AI-augmented** — natural language to JQL, sprint health analysis, anomaly detection
4. **Memory-persistent** — learns team patterns, remembers decisions, tracks velocity trends
5. **go-jira feature parity** — every useful thing go-jira does, we do too (and better)

---

## Current State (Baseline)

### What Exists (14 MCP tools)
- Jira CRUD: search, get_issue, create_issue, add_comment, transitions, transition
- Boards & Sprints: boards, sprint_summary
- PM Views: my_issues, overdue, workload
- AI: ai_analyze, ai_sprint_report
- Notifications: notify_lark

### Critical Issues
- Memory subsystem implemented but NOT wired (handlers exist, not registered)
- Zero tests
- No error wrapping / structured errors
- No pagination for large results
- No update/edit issue capability
- go-sqlite3 in go.mod but memory provider not connected to DI

---

## Research Foundation

### PM/Scrum Science (What the Research Says)

| Insight | Source | Implication for Tool |
|---------|--------|---------------------|
| Flow metrics > velocity for predictability | DORA 2024, Actionable Agile | Build cycle time, throughput, WIP tracking |
| WIP limits reduce interruptions 30% | Little's Law studies | Surface WIP violations, alert on overload |
| Context switching costs 40% of productive time | HBR, DevEx research | Minimize tool interactions, batch operations |
| AI helps struggling PMs, hurts top performers | Beehiiv PM research | Be smart about when to intervene vs stay silent |
| Stable priorities = less burnout | DORA 2024 | Detect priority churn, surface to leadership |
| Monte Carlo > single-point estimates | Probabilistic forecasting research | Build throughput-based forecasting |
| Stale backlogs are anti-patterns | Scrum.org patterns | Auto-detect and flag stale items |
| 68% of automation failures = overcomplicated rules | Jira automation research | Keep our automations simple, one-purpose |
| Developer flow state needs 4hr blocks | GitHub DevEx 2024 | Don't spam notifications, batch insights |
| SM as system coach, not ceremony cop | ThinkLouder 2024 | Tool should coach, not just report |

### go-jira DNA (What to Steal)

| Pattern | Why It Works | Our Implementation |
|---------|-------------|-------------------|
| Hierarchical config (.jira.d/) | Context-aware per-project | MCP resources for project context |
| Custom commands with templates | Extensibility without code | Composable tool chains |
| Editor integration ($EDITOR) | Rich text without leaving terminal | AI handles rich descriptions |
| Template system (input + output) | Customizable everything | Prompt templates for AI operations |
| JQL as first-class citizen | Power users love JQL | Natural language -> JQL + raw JQL |
| Single binary, zero deps | Easy install | Single Go binary via MCP |
| Per-command config overrides | Flexibility | Per-tool configuration |
| `request` command (raw API) | Escape hatch | `jira_raw_request` tool |

### What go-jira Couldn't Do (Our Advantages)

| Capability | go-jira | Us |
|-----------|---------|-----|
| Sprint management commands | Missing | First-class |
| Board operations | Missing | Full CRUD |
| AI analysis | Impossible | Native |
| Flow metrics | Impossible | Built-in |
| Risk prediction | Impossible | AI-powered |
| Memory/learning | Impossible | SQLite + embeddings |
| Natural language queries | Impossible | NL -> JQL |
| Bulk operations | Shell scripting hack | Atomic tool |
| Multi-project intelligence | Separate configs | Unified view |

---

## Architecture Principles

1. **Tools are atomic** — each tool does ONE thing well, composable by the AI agent
2. **Intelligence at the edge** — raw data tools + smart analysis tools, never mix
3. **Fail fast, fail clear** — structured errors with actionable messages
4. **Memory is optional** — works without memory, better with it
5. **Pagination is mandatory** — no tool returns unbounded results
6. **Idempotent where possible** — safe to retry
7. **Batch > loop** — bulk tools over N single calls

---

## Phases

### Phase 0: Foundation Fix (Week 1)
> Fix what's broken. Make it solid.

| Task | Priority | Effort |
|------|----------|--------|
| Wire memory subsystem into DI + register tools | P0 | 2hr |
| Add structured error types (domain/errors.go) | P0 | 1hr |
| Add pagination to jira_search (startAt, maxResults, hasMore) | P0 | 1hr |
| Add `jira_update_issue` tool (summary, description, priority, assignee, labels, custom fields) | P0 | 2hr |
| Add basic test infrastructure (handlers_test.go with mocked interfaces) | P0 | 3hr |
| Fix go.mod consistency (ensure sqlite3 is used or remove) | P0 | 30min |
| Add health/version introspection tool (`jira_health`) | P1 | 30min |
| Add Makefile targets: test, lint, build, coverage | P1 | 30min |

**Exit criteria:** All tools compile, memory works, `make test` passes, update works.

---

### Phase 1: go-jira Feature Parity (Week 2-3)
> Match everything useful from go-jira. No excuses.

#### Issue Operations (Core CRUD)
| Tool | Description | go-jira Equivalent |
|------|-------------|-------------------|
| `jira_update_issue` | Update any field (summary, desc, priority, assignee, labels, components, custom) | `edit` |
| `jira_delete_issue` | Delete issue (with confirmation flag) | N/A |
| `jira_assign` | Assign issue to user | `assign` |
| `jira_unassign` | Remove assignee | `unassign` |
| `jira_clone_issue` | Clone issue with optional modifications | N/A |

#### Sub-tasks & Hierarchy
| Tool | Description |
|------|-------------|
| `jira_create_subtask` | Create subtask under parent |
| `jira_list_subtasks` | List subtasks of an issue |
| `jira_reparent` | Move subtask to different parent |

#### Epic Management
| Tool | Description | go-jira Equivalent |
|------|-------------|-------------------|
| `jira_epic_list` | List issues in an epic | `epic list` |
| `jira_epic_add` | Add issues to an epic | `epic add` |
| `jira_epic_remove` | Remove issues from an epic | `epic remove` |
| `jira_epic_create` | Create a new epic | `epic create` |

#### Links & Relations
| Tool | Description | go-jira Equivalent |
|------|-------------|-------------------|
| `jira_link_issues` | Create link (blocks, relates, duplicates, etc.) | `issuelink` |
| `jira_unlink_issues` | Remove a link | N/A |
| `jira_list_links` | Show all links for an issue | via `view` |
| `jira_link_types` | List available link types | `issuelinktypes` |

#### Attachments
| Tool | Description |
|------|-------------|
| `jira_attach_list` | List attachments on issue |
| `jira_attach_upload` | Upload file to issue |
| `jira_attach_download` | Download attachment (return path) |
| `jira_attach_delete` | Remove attachment |

#### Worklogs & Time
| Tool | Description |
|------|-------------|
| `jira_worklog_add` | Log time on issue |
| `jira_worklog_list` | List worklogs for issue |
| `jira_worklog_delete` | Remove worklog entry |

#### Labels & Components
| Tool | Description |
|------|-------------|
| `jira_labels_add` | Add labels to issue |
| `jira_labels_remove` | Remove labels from issue |
| `jira_components_list` | List project components |
| `jira_component_add` | Add component to project |

#### Watchers & Votes
| Tool | Description |
|------|-------------|
| `jira_watch` | Add watcher to issue |
| `jira_unwatch` | Remove watcher |
| `jira_watchers` | List watchers |
| `jira_vote` | Vote for issue |

#### Projects & Users
| Tool | Description |
|------|-------------|
| `jira_projects` | List accessible projects |
| `jira_project_detail` | Get project info (lead, components, versions) |
| `jira_user_search` | Find users by name/email (for assignee lookup) |
| `jira_myself` | Get current authenticated user info |

#### Versions & Releases
| Tool | Description |
|------|-------------|
| `jira_versions` | List project versions |
| `jira_version_create` | Create a version |
| `jira_version_release` | Mark version as released |

#### Bulk Operations
| Tool | Description |
|------|-------------|
| `jira_bulk_transition` | Transition multiple issues at once |
| `jira_bulk_assign` | Assign multiple issues to one user |
| `jira_bulk_update` | Update field on multiple issues |
| `jira_bulk_move` | Move issues to different project/epic |

#### Sprint Management (go-jira's biggest gap!)
| Tool | Description |
|------|-------------|
| `jira_sprints` | List sprints for a board (active, future, closed) |
| `jira_sprint_issues` | List all issues in a sprint |
| `jira_sprint_create` | Create a new sprint |
| `jira_sprint_start` | Start a sprint (set dates) |
| `jira_sprint_close` | Complete a sprint |
| `jira_sprint_move_issues` | Move issues to a sprint |

#### Raw Access & Escape Hatch
| Tool | Description | go-jira Equivalent |
|------|-------------|-------------------|
| `jira_raw_request` | Arbitrary REST API call (method, path, body) | `request` |
| `jira_fields` | List all available fields (including custom) | `fields` |
| `jira_createmeta` | Get creation metadata for project/issuetype | `createmeta` |

**Phase 1 Exit criteria:** Feature parity with go-jira. Every go-jira command has an equivalent MCP tool (or better). Full test coverage on handlers.

---

### Phase 2: PM Intelligence Layer (Week 4-5)
> This is where we surpass go-jira. AI-powered PM insights.

#### Flow Metrics Engine
| Tool | Description | Research Basis |
|------|-------------|---------------|
| `pm_cycle_time` | Average cycle time per issue type, last N sprints | Little's Law, DORA |
| `pm_throughput` | Issues completed per sprint, trend analysis | Flow metrics research |
| `pm_wip_status` | Current WIP per person/status, vs limits | WIP limits research |
| `pm_lead_time` | Time from creation to done, breakdown by stage | Kanban flow metrics |
| `pm_flow_efficiency` | Active time vs wait time ratio | Lean/flow research |

#### Predictive Analytics
| Tool | Description | Research Basis |
|------|-------------|---------------|
| `pm_forecast_sprint` | Monte Carlo simulation: what can we finish this sprint? | Probabilistic forecasting |
| `pm_forecast_epic` | When will this epic be done? (confidence intervals) | ActionableAgile methodology |
| `pm_velocity_trend` | Velocity trend with anomaly detection | Statistical process control |
| `pm_capacity_plan` | Team capacity vs planned work next sprint | Sprint planning best practices |

#### Risk & Health Detection
| Tool | Description | Research Basis |
|------|-------------|---------------|
| `pm_blockers` | Issues blocked or stale, with duration | Flow metrics, WIP |
| `pm_scope_creep` | Items added to sprint after start | Sprint commitment research |
| `pm_priority_churn` | Priority changes frequency (stability indicator) | DORA 2024 unstable priorities |
| `pm_team_health` | Composite health score (WIP, blockers, velocity, churn) | Multiple frameworks |
| `pm_sprint_risk` | Sprint goal at risk? Based on burndown + historical | Burn-up forecasting |
| `pm_tech_debt_ratio` | Bug/debt items vs feature items trend | TD management research |

#### Retrospective Support
| Tool | Description |
|------|-------------|
| `pm_retro_data` | Auto-generate retro data: what went well (early deliveries), what didn't (overdue, blockers), patterns |
| `pm_improvement_track` | Track retro action items across sprints |

#### Natural Language Interface
| Tool | Description |
|------|-------------|
| `pm_ask` | Natural language question -> appropriate tool call(s) + synthesized answer |
| `pm_nl_to_jql` | Convert natural language to JQL query |

**Phase 2 Exit criteria:** Flow metrics calculated from real data. Forecast tools produce confidence intervals. Risk detection runs against live sprint data.

---

### Phase 3: Memory & Learning (Week 6-7)
> The tool gets smarter over time. Learns team patterns.

#### Memory Infrastructure
| Component | Description |
|-----------|-------------|
| SQLite store (already coded) | Wire it. Snapshots, decisions, risks, patterns |
| Embedding index | Vector similarity for pattern matching |
| Temporal decay | Recent patterns weighted higher |

#### Memory Tools (Already Coded, Need Wiring)
| Tool | Description |
|------|-------------|
| `memory_snapshot` | Save sprint/project state snapshot |
| `memory_risk` | Record identified risk with context |
| `memory_decision` | Record decision with rationale |
| `memory_pattern` | Record observed team pattern |
| `memory_recall` | Semantic search over memories |
| `memory_timeline` | Chronological event history |
| `memory_forget` | GDPR-style deletion |

#### Learning Capabilities
| Tool | Description |
|------|-------------|
| `learn_velocity_model` | Build team velocity model from historical sprints |
| `learn_estimation_accuracy` | Track estimate vs actual, per person/type |
| `learn_blocker_patterns` | Which types of work get blocked? Why? |
| `learn_team_rhythm` | When does the team ship? Early/late sprint patterns |

**Phase 3 Exit criteria:** Memory persists across sessions. Pattern detection improves recommendations. Estimation accuracy tracked.

---

### Phase 4: Automation & Proactive (Week 8-9)
> Don't wait to be asked. Surface insights proactively.

#### Proactive Alerts (MCP Notifications)
| Notification | Trigger |
|--------------|---------|
| Sprint goal at risk | Burndown deviates >20% from ideal |
| WIP limit exceeded | Person has >N items in progress |
| Blocker aging | Item blocked >48 hours |
| Scope creep alert | >3 items added mid-sprint |
| Stale review | PR/item in review >24 hours |
| Deadline approaching | Items due within 2 days, not started |

#### Automation Actions
| Tool | Description |
|------|-------------|
| `auto_standup_report` | Generate daily standup from yesterday's activity |
| `auto_sprint_close_report` | End-of-sprint summary with metrics |
| `auto_backlog_groom` | Flag stale items (>90 days untouched), suggest archive |
| `auto_release_notes` | Generate release notes from completed issues |
| `auto_onboard_context` | Generate project context for new team member |

#### Workflow Recipes (Composable Sequences)
| Recipe | Steps |
|--------|-------|
| `recipe_start_work` | Assign to me -> transition to In Progress -> create branch name |
| `recipe_done` | Transition to Done -> log time -> update sprint burndown |
| `recipe_block` | Flag as blocked -> add blocker comment -> notify team |
| `recipe_split_story` | Create N subtasks from story -> distribute story points |
| `recipe_sprint_plan` | Show capacity -> suggest items from backlog -> move to sprint |

**Phase 4 Exit criteria:** Notifications fire on real conditions. Recipes execute multi-step workflows atomically.

---

### Phase 5: Multi-Project & Portfolio (Week 10-11)
> Scale from one team to multiple.

| Tool | Description |
|------|-------------|
| `portfolio_overview` | Cross-project health dashboard |
| `portfolio_dependencies` | Cross-project blocking relationships |
| `portfolio_resource_conflicts` | People overallocated across projects |
| `portfolio_roadmap_status` | Epic/initiative progress across projects |
| `portfolio_risk_radar` | Aggregate risks across all projects |

---

### Phase 6: Integration & Ecosystem (Week 12+)
> Connect to the world.

| Integration | Purpose |
|-------------|---------|
| GitHub/GitLab | PR -> issue linking, auto-transitions, commit activity |
| Confluence | Auto-generate/update project pages |
| Slack/Lark/Teams | Multi-channel notifications (Lark already started) |
| Calendar | Sprint ceremonies, deadline awareness |
| CI/CD (GitHub Actions) | Deploy events -> version release automation |
| Tempo/Clockify | Time tracking integration |
| OKR tools | Align sprint work to objectives |

---

## Tool Count Summary

| Phase | Tools Added | Running Total |
|-------|-------------|---------------|
| Current | 14 | 14 |
| Phase 0 (Fix) | +3 | 17 |
| Phase 1 (Parity) | +45 | 62 |
| Phase 2 (Intelligence) | +16 | 78 |
| Phase 3 (Memory) | +11 | 89 |
| Phase 4 (Automation) | +10 | 99 |
| Phase 5 (Portfolio) | +5 | 104 |
| Phase 6 (Integration) | +7 | 111 |

---

## Technical Architecture (Target State)

```
cmd/server/main.go                    # fx bootstrap

config/
  config.go                           # Env-based config

domain/
  errors.go                           # Structured error types
  jira/
    domain.go                         # Issue, Sprint, Board, Epic models
    client.go                         # Jira Client interface (expanded)
  ai/
    provider.go                       # AI Provider interface
  memory/
    models.go                         # Memory domain models
    store.go                          # Memory Store interface
  metrics/
    models.go                         # Flow metric models (CycleTime, Throughput, etc.)
    calculator.go                     # Calculator interface

internal/
  jira/
    client.go                         # Jira REST v3 + Agile implementation
    bulk.go                           # Bulk operation helpers
    pagination.go                     # Pagination utilities
  ai/
    client.go                         # OpenAI-compatible client
    jql.go                            # NL -> JQL conversion
    forecast.go                       # Monte Carlo simulation
  memory/
    sqlite.go                         # SQLite persistence
    embeddings.go                     # Vector similarity (optional)
  metrics/
    calculator.go                     # Flow metrics calculation
    historical.go                     # Historical data aggregation
  lark/
    webhook.go                        # Lark notifications

application/
  tools/
    jira_core.go                      # CRUD tools
    jira_search.go                    # Search & JQL tools
    jira_sprint.go                    # Sprint management tools
    jira_epic.go                      # Epic tools
    jira_bulk.go                      # Bulk operation tools
    jira_link.go                      # Link tools
    jira_attach.go                    # Attachment tools
    jira_worklog.go                   # Time tracking tools
    jira_meta.go                      # Fields, createmeta, raw request
    pm_metrics.go                     # Flow metrics tools
    pm_forecast.go                    # Prediction tools
    pm_risk.go                        # Risk detection tools
    pm_automation.go                  # Proactive tools
    memory.go                         # Memory tools
    recipes.go                        # Workflow recipes

transport/
  server.go                           # MCP server + tool registration
  notifications.go                    # Proactive notification dispatch
```

---

## Quality Gates (Per Phase)

| Gate | Requirement |
|------|-------------|
| Compilation | `make build` succeeds |
| Linting | `make lint` (golangci-lint) passes |
| Tests | `make test` passes, >70% coverage on handlers |
| Integration | Manual test against real Jira instance |
| Documentation | Tool descriptions are clear, parameters documented |
| Performance | No tool takes >5s for typical operations |
| Error handling | All errors are structured, actionable, non-panicking |

---

## Non-Goals (Explicitly Out of Scope)

- **Web UI** — this is an MCP server, not a web app
- **Multi-tenant** — single user/team focus
- **Real-time sync** — pull-based, not push (no webhooks server... yet)
- **Replace Jira** — augment it, not compete with it
- **Board layout management** — that's a visual concern
- **Permission management** — too dangerous for automation

---

## Success Metrics

| Metric | Target | How to Measure |
|--------|--------|----------------|
| Tool coverage vs go-jira | 100% parity | Feature matrix comparison |
| Time saved per sprint | >2 hours | User self-report |
| Risk detection accuracy | >70% of actual blockers caught early | Retrospective validation |
| Forecast reliability | Within 20% of actual | Sprint actual vs predicted |
| Adoption | Daily use by owner | Memory call frequency |

---

## References

- DORA 2024 Report: https://dora.dev/report/2024
- GitHub DevEx Research 2024: https://github.blog/news-insights/research/good-devex-increases-productivity/
- go-jira: https://github.com/go-jira/jira (unmaintained since 2020)
- ankitpokhrel/jira-cli: https://github.com/ankitpokhrel/jira-cli (modern alternative)
- Atlassian Rovo Dev CLI: https://community.atlassian.com/forums/Jira-articles/Introducing-Rovo-Dev-CLI-Manage-Jira-Issues-From-Your-Terminal/ba-p/3046980
- 18th State of Agile (2025): digital.ai
- DX Core 4 Framework: https://getdx.com/research/measuring-developer-productivity-with-the-dx-core-4/
- Probabilistic Forecasting: https://agiledigest.com/probabilistic-forecasting-in-agile/
- WIP Limits Research: Little's Law applied to software development
- Scrum.org Jira Anti-Patterns: https://scrum.org/resources/blog/jira-anti-patterns-and-how-overcome-them
