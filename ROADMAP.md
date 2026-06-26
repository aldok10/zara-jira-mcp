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
- `docs/communication-frameworks.md` — 14 frameworks: Minto Pyramid, SCARF, SBI, RACI/DACI, Radical Candor, NVC, 5W1H, Async Protocol, Ceremony Patterns, Escalation/TIRED, Crucial Conversations, Signal-over-Noise, Communication Anti-Patterns, Trust Pyramid
- DORA 2025: PRs merged +98%, incidents +242% — velocity metrics lie without quality signals
- Industry standard: 15-20% sprint capacity for tech debt (confirmed across 6 sources)
- Sprint goal success rate: only 52% of teams achieve goals (Scrum Alliance)
- Flow metrics > velocity for predicting delivery (cycle time, throughput, WIP)

---

## Next Phase: AI Communication Layer

Based on `docs/communication-frameworks.md` + fresh research (Jun 2026). The thesis: **PM's #1 job in the AI era is communication — framing, routing, timing. AI handles data; PM handles meaning.**

**Research backing:**
- $75M at risk per $1B spent from bad communication (PMI)
- Teams with communication overload: 3x slower decisions (ITS Dart)
- Trust in AI tools -31%, agentic AI -89% (Axis Intelligence 2026)
- 15.4 hrs/week meetings vs 12.1 hrs deep work (Microsoft 2026)
- Structured async = 25% fewer meetings (multiple sources)
- Monte Carlo confidence intervals build stakeholder trust > single-point estimates

---

### Comms Phase 1: Signal-over-Noise (next — low effort, high impact)

**Thesis:** Reduce noise, amplify signal. Make decisions searchable, make status self-service.

| Tool | Framework | What | Effort |
|------|-----------|------|--------|
| `pm_communicate` | Minto Pyramid | Generate audience-specific updates (exec vs team vs PO) | Medium |
| `pm_feedback_prep` | SBI Model | AI-generate structured feedback from team data | Medium |
| `pm_escalation_draft` | TIRED + Pyramid | Draft escalation with context + ask + deadline | Low |
| `pm_decision_record` | ADR/MADR | Enhanced decisions: context → options → consequences | Low |
| `pm_comms_health` | Signal-over-Noise | Score: re-decision rate, blocker escalation speed, stakeholder responsiveness | Medium |

**Success metric:** Decision re-occurrence drops 50% (same topic not revisited). Blocker avg age drops 30%.

---

### Comms Phase 2: Audience-Aware Routing (medium effort, high impact)

**Thesis:** Same data, different framing per audience. SCARF-aware message construction.

| Tool | Framework | What | Effort |
|------|-----------|------|--------|
| `pm_audience_router` | Minto + SCARF | Same data, auto-reframe: remove jargon for exec, add certainty for anxious stakeholders | Medium |
| `pm_communication_plan` | RACI + Timing | Who needs what, when, via which channel — auto-enforced | Medium |
| `pm_raci` | RACI Matrix | Auto-generate from Jira assignments | Low |
| `pm_silence_detector` | Ghost Stakeholders | Flag disengaged stakeholders (no pulse, no decisions, no comments) | Low |
| `pm_comms_anti_patterns` | Community Smells | Detect: information hoarding, over-communication, meeting addiction, re-deciding | Medium |

**Success metric:** Stakeholder pulse scores improve. "I didn't know about that" incidents drop to zero.

---

### Comms Phase 3: AI Coaching Communication (medium effort, high impact)

**Thesis:** AI preps the PM for difficult human interactions — data + script, but human delivers.

| Tool | Framework | What | Effort |
|------|-----------|------|--------|
| `pm_hard_conversation` | Crucial Conversations + SBI + data | Prep: facts, your story, their likely perspective, SCARF risks, opening line | High |
| `pm_meeting_prep` | Pyramid + 5W1H | Agenda + talking points + pre-read links for any meeting | Low |
| `pm_async_update` | Signal-over-Noise | Structured status that replaces a sync meeting | Low |
| `pm_trust_signals` | Trust Pyramid | Track: forecast accuracy, override rate, team confidence trend | Medium |
| `pm_nvc_reframe` | NVC | Reframe blaming/judgmental language → observation + feeling + need + request | Medium |

**Success metric:** pm_stakeholder_pulse trend positive. Team confidence scores stable/rising.

---

### Comms Phase 4: Organizational Communication (stretch — high effort)

**Thesis:** Communication architecture at org level — change management, conflict resolution, influence mapping.

| Tool | Framework | What | Effort |
|------|-----------|------|--------|
| `pm_change_communication` | Kotter 8-step + ADKAR | Plan comms for org changes (reorgs, tool migrations, process shifts) | High |
| `pm_conflict_mediation` | Thomas-Kilmann + NVC | AI conflict diagnosis + suggest resolution approach | High |
| `pm_influence_map` | Power/Interest Matrix | Stakeholder influence strategy per initiative | Medium |
| `pm_calibration_report` | Superforecasting | How accurate were past AI forecasts? Build institutional trust | Medium |
| `pm_ceremony_optimizer` | Meeting ROI + async-first | Recommend: which ceremonies to keep sync, which to go async, which to kill | Low |

**Success metric:** PM time saved 5+ hrs/week. Forecast accuracy within 20% at 85% confidence. Zero surprise escalations.

---

## Context Engineering: Tool Delivery Optimization

Based on `research/context-engineering.md`. The problem: 231 tools = ~50K tokens overhead. Most clients choke.

### Phase 11: Optimize Current (low effort)

| Task | Impact |
|------|--------|
| Shorten all tool descriptions to <50 chars in `all` profile | -30% token overhead |
| Add `pm_load_module` for dynamic tool loading mid-session | Lazy registration |
| Track tool usage per PM (which tools actually called) | Data for auto-profile |
| Auto-suggest optimal profile based on usage history | PM friction = 0 |

### Phase 12: Dynamic Discovery (medium effort)

| Task | Impact |
|------|--------|
| Hierarchical tool listing (category → sub → tool) | Agent drills down |
| `pm_discover intent="X"` returns only relevant 5-10 tools | Smart routing |
| Lazy tool registration (tools loaded on first category access) | Minimal initial schema |

### Phase 13: Intent-Based Selection (high effort, transformative)

| Task | Impact |
|------|--------|
| AI intent classification before tool selection | Right tool every time |
| Auto-compress descriptions based on conversation context | Adaptive schema |
| Cross-session tool preference learning | Personalized per PM |

---

## Implementation Priority Matrix

| Priority | Phase | Effort | Impact | Ship target |
|----------|-------|--------|--------|-------------|
| 1 | 11 (Context optimize) | Low | High | This week |
| 2 | 7 (Communication templates) | Low | High | Next week |
| 3 | 8 (Smart routing) | Medium | High | Sprint +1 |
| 4 | 12 (Dynamic discovery) | Medium | High | Sprint +1 |
| 5 | 9 (AI coaching comms) | Medium | Medium | Sprint +2 |
| 6 | 13 (Intent-based) | High | Transformative | Sprint +3 |
| 7 | 10 (Org communication) | High | Medium | Backlog |

---

## North Star

> A PM opens ChatGPT/Claude/OpenCode, types "how's my sprint going?" and gets an answer in 2 seconds — backed by live Jira data, historical memory, team health signals, and stakeholder context. No manual tool selection. No profile configuration. No lag.

That's the end state. Everything else is a step toward it.
