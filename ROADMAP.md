# ROADMAP: zara-jira-mcp

> AI-powered Scrum Master MCP — 224 tools. Persistent memory. 20 AI clients. Profile system.
> The PM/SM's unfair advantage.

---

## Current State (v0.4.0)

**224 MCP tools** | 14 SQLite tables | 8 notification platforms | 20 AI client configs | 5 performance profiles

### What's Built

| Domain | Tools | Coverage |
|--------|-------|----------|
| Jira Core | 55 | Full CRUD + epics + sprints + bulk ops + versions |
| PM Intelligence | 30 | Flow metrics, Monte Carlo, coaching, anti-patterns, forecasting |
| PM Memory | 22 | Sprints, risks, decisions, blockers, team, retros, deps, goals, DoD/DoR |
| Notifications | 15 | Lark, Slack, Discord, Telegram, Teams, Email, Confluence, routing |
| GitHub/GitLab | 13 | Issues, milestones, MRs, file reading, branch tracing |
| Management Reporting | 10 | Exec report, management brief, escalation, dependency, resource util |
| SM Leverage | 10 | Maturity, dysfunction, meeting ROI, commitment, impact, autonomy |
| Portfolio | 5 | Cross-project overview, risks, workload, blockers, AI summary |
| Shortcuts & Help | 8 | pm, pm_create, pm_decide, pm_risk, pm_next, pm_help, pm_quickstart, pm_workflow |
| Stakeholder | 8 | Pulse, trend, outcome map, improvement velocity, scorecard, KB |
| Integrations | 17 | Linear, PagerDuty, Clockify, Notion, Google Sheets, Calendar |
| Database | 5 | Postgres, MySQL, MongoDB read queries |
| Tech Skill | 5 | PM technical literacy tools |
| Meta | 3 | Health check, MCP stats, NL-to-JQL |

### What's New in v0.4.0

- Performance profiles (PM_PROFILE: lite/pm/standard/full/all)
- 20 AI client pre-built configs
- Management reporting tools (brief, dependency, escalation, resource, commitment)
- SM leverage tools (maturity, dysfunction, meeting ROI)
- Engineering literacy documentation for PM/SM
- CI pipeline (GitHub Actions)
- llms.txt for AI discoverability

---

## Phase 7: Production Hardening (In Progress)

**Goal:** Make this reliable for daily use by real PM teams.

| Item | Effort | Impact | Status |
|------|--------|--------|--------|
| Module-level enable/disable (PM_PROFILE) | Low | High | Done |
| Connection health checks (Jira, AI, Slack) | Low | Medium | Partial |
| Rate limiting for Jira API calls | Low | Medium | Pending |
| Graceful degradation when AI unavailable | Low | High | Done |
| SQLite backup command (pm_backup) | Low | Medium | Pending |
| Error messages that suggest fix actions | Low | High | Pending |
| Fix security findings (raw_request, db query, file perms) | Medium | Critical | Pending |
| Input validation (issue keys, JQL params) | Low | High | Pending |
| io.LimitReader on all HTTP responses | Low | Medium | Pending |
| Test coverage to 60% | Medium | High | Pending |

---

## Phase 8: Smart Context (Planned)

**Goal:** The MCP learns team patterns and proactively surfaces insights without being asked.

| Item | Effort | Impact |
|------|--------|--------|
| Auto-snapshot sprint end (detect sprint close event) | Medium | High |
| Pattern recognition: "this sprint looks like Sprint 7 which failed" | High | High |
| Predictive blockers: "Alice usually gets blocked on external API tasks" | High | High |
| Auto-generate retro data points from sprint history | Medium | Medium |
| Meeting effectiveness scoring (decisions/actions ratio) | Low | Medium |
| Confidence calibration (track prediction accuracy over time) | Medium | High |

---

## Phase 9: Developer Integration (Planned)

**Goal:** Bridge the PM-Developer gap. Bidirectional visibility.

| Item | Effort | Impact |
|------|--------|--------|
| GitHub Actions webhook → auto-update Jira status | Medium | High |
| GitLab pipeline status → sprint health factor | Medium | High |
| PR review time tracking → flow metrics | Medium | Medium |
| Branch → Jira auto-link (on branch create) | Medium | Medium |
| Deploy frequency tracking (DORA metric) | Medium | High |
| Escaped defects detection (prod bugs from recent releases) | High | High |

---

## Phase 10: Team Autonomy (Vision)

**Goal:** The team gradually needs the SM less. MCP coaches the team directly.

| Item | Effort | Impact |
|------|--------|--------|
| Individual developer dashboards (my flow, my debt, my blockers) | Medium | Medium |
| Self-service sprint health (team can run pm without SM) | Low | High |
| Automated working agreement enforcement | High | Medium |
| Sprint auto-scoring at close (no manual snapshot needed) | Medium | High |
| Maturity model tracking (team progress toward self-organization) | High | Medium |
| Onboarding guide generation from team KB | Medium | Medium |

---

## Phase 11: Multi-Team / Enterprise (Future)

**Goal:** Scale from single team to program/portfolio level.

| Item | Effort | Impact |
|------|--------|--------|
| Multi-board aggregation (Scrum of Scrums) | High | High |
| Cross-team dependency visualization | High | High |
| Program-level forecasting | High | Medium |
| Normalized velocity across teams | Medium | Medium |
| Enterprise risk rollup | Medium | High |
| Multi-tenant SQLite (per team DB) | Medium | Medium |

---

## Phase 12: Ecosystem (Future)

**Goal:** Become the standard PM MCP that works with any AI agent.

| Item | Effort | Impact | Status |
|------|--------|--------|--------|
| Pre-built configs for 20 AI clients | Low | High | Done |
| Docker Hub image (one-command deploy) | Low | Medium | Pending |
| Helm chart for K8s | Medium | Low | Pending |
| SSE/HTTP transport (remote MCP) | High | High | Partial |
| OAuth2 for multi-user access | High | Medium | Pending |
| Plugin system (custom tools per team) | High | Medium | Pending |
| Open-source community (docs, contributing guide) | Medium | High | Done |
| MCP marketplace / registry listing | Low | High | Pending |

---

## Key Metrics to Track

| Metric | Current | Target v1.0 |
|--------|---------|-------------|
| Tools | 224 | 230+ (stable) |
| Test coverage | ~27% | 70%+ |
| AI client configs | 20 | 25+ |
| Daily active use | 0 | 1 team |
| Sprint snapshots captured | 0 | 10+ per board |
| Decisions recorded | 0 | 50+ |
| Forecast accuracy | unmeasured | within 20% at 85% confidence |
| PM time saved | unmeasured | 5+ hours/week |

---

## Principles

1. **Ship to learn** — Every phase ships. No "big bang" releases.
2. **Data before features** — More sprint snapshots = better AI. Prioritize data capture.
3. **PM friction = bug** — If a PM has to think about which tool to use, we failed.
4. **Automate the boring** — Snapshots, digests, escalations should be automatic.
5. **Measure value** — Track "PM time saved" not "tools added".
6. **One team first** — Perfect for one team before scaling to many.
7. **Security is non-negotiable** — Fix critical findings before adding features.

---

## Research Foundation

- `research/scrum-master-papers.md` — 508 academic papers on SM effectiveness
- `research/pm-integration-platforms.md` — Notification routing, escalation patterns
- `research/pm-leverage-research.md` — DORA metrics, priority churn, tech debt frameworks
- DORA 2025: PRs merged +98%, incidents +242% — velocity metrics lie without quality signals
- Industry standard: 15-20% sprint capacity for tech debt (confirmed across 6 sources)
- Sprint goal success rate: only 52% of teams achieve goals (Scrum Alliance)
- Flow metrics > velocity for predicting delivery (cycle time, throughput, WIP)
