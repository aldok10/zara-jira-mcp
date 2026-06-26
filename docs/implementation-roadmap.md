# Implementation Roadmap

Prioritized backlog berdasarkan research findings. Setiap item linked ke data point yang justify-nya.

---

## Sprint Next: Security + Stability (1-2 minggu)

Rationale: nggak bisa kirim ke real PM team kalau ada security hole.

| # | Task | Research Basis | Effort |
|---|------|---------------|--------|
| 1 | Restrict `jira_raw_request` — path allowlist, block DELETE on sensitive paths | Security audit finding #2 | 2h |
| 2 | Restrict `DatabaseQuery` — read-only user enforcement, reject multi-statement | Security audit finding #1 | 2h |
| 3 | Fix file permissions (`0700` dir, `0600` config) | Security audit finding #3 | 30m |
| 4 | Add input validation: issue key regex, JQL param sanitization | Security audit finding #4 | 3h |
| 5 | Add `io.LimitReader(10MB)` on all HTTP response reads | Security audit finding #8 | 1h |
| 6 | Sanitize error messages — generic to client, full to log | Security audit finding #9 | 2h |
| 7 | Move Gemini API key from URL to header | Security audit finding #6 | 30m |
| 8 | Dashboard: bind localhost, add basic auth | Security audit finding #5 | 1h |

**Exit criteria:** Zero CRITICAL findings. All HIGH findings addressed or mitigated.

---

## Sprint +1: Daily Usability (1-2 minggu)

Rationale: "PM friction = bug" (principle #3). Research: 78% PM adopt faster when setup is simple.

| # | Task | Research Basis | Effort |
|---|------|---------------|--------|
| 1 | Auto-detect board_id (cache after first `jira_boards` call) | PM calls tools 50+ times/day, board_id every time is friction | 2h |
| 2 | `pm_morning` composite tool — standup + blockers + calendar + PRs in 1 call | PwC: save 7.4h/week starts with morning routine | 3h |
| 3 | Rate limiting: Jira (60 req/min), AI (10 req/min), Notifications (5/min) | Security audit + Jira Cloud rate limits | 3h |
| 4 | `pm_backup` — SQLite export to JSON | Data loss = trust loss | 1h |
| 5 | Better error messages: "Jira API rate limited, retry in 30s" not raw error | 35% PM churns tool after first confusing error | 2h |
| 6 | `pm_onboard` — first-run wizard: detect config, suggest profile, explain tools | Adoption barrier is setup complexity | 3h |

**Exit criteria:** PM can go from install to "prep my standup" in < 3 minutes.

---

## Sprint +2: Smart Context (2-3 minggu)

Rationale: DORA 2025 — AI works when teams learn where/when it's useful. Memory that accumulates = compounding value.

| # | Task | Research Basis | Effort |
|---|------|---------------|--------|
| 1 | Auto-snapshot sprint end (detect when sprint closes via Jira webhook or polling) | 52% teams miss goals — need data to improve. Manual snapshot = forgotten | 4h |
| 2 | Sprint pattern matching: "this sprint resembles Sprint 7 (which missed goal)" | Pattern recognition across history = predictive PM | 6h |
| 3 | Predictive blocker alerts: "this task type blocks 40% of time for this assignee" | Techademy: SM admin reduced 30-50% via proactive detection | 6h |
| 4 | Confidence calibration: track forecast accuracy, adjust over time | METR 2025: perception ≠ reality. Need honest calibration | 4h |
| 5 | Meeting effectiveness score: decisions+actions per meeting minute | "9.9% every dollar wasted" (PMI) — meetings are #1 time sink | 3h |

**Exit criteria:** After 5 sprints of data, system proactively surfaces 2-3 useful insights per week without being asked.

---

## Sprint +3: Communication Intelligence (2-3 minggu)

Rationale: PM role shift = "communication architect." Tools harus support structured communication, not just data.

| # | Task | Research Basis | Effort |
|---|------|---------------|--------|
| 1 | `pm_brief` — one tool, auto-detect audience from context, apply Pyramid Principle | Minto: executives process top-down | 4h |
| 2 | `pm_feedback` — generate SBI-structured feedback from sprint data | SBI model: reduces defensiveness | 3h |
| 3 | `pm_escalate_draft` — TIRED-formatted escalation message ready to send | Late escalation = project failure | 2h |
| 4 | `pm_retro_format` — suggest retro format based on team mood/history (4Ls, DAKI, Mad-Sad-Glad) | Techademy: fresh format = better participation | 2h |
| 5 | `pm_async_digest` — replace daily standup meeting with async written update | Async research: distributed teams perform better with write-first | 3h |
| 6 | Template library: common messages (sprint complete, risk escalation, goal missed) with framework applied | 5W1H completeness + Pyramid structure | 4h |

**Exit criteria:** PM can generate stakeholder-appropriate communication for any scenario in < 30 seconds.

---

## Sprint +4: Developer Bridge (2-3 minggu)

Rationale: DORA 2025 paradox — devs merge 98% more PRs, org delivery flat. Need system-level metrics.

| # | Task | Research Basis | Effort |
|---|------|---------------|--------|
| 1 | DORA metrics dashboard: deploy frequency, lead time, change failure rate, MTTR | DORA: these 4 predict delivery performance | 6h |
| 2 | PR review SLA tracking with alerts | PR age > 3 days = invisible bottleneck | 3h |
| 3 | Escaped defects correlation: prod bugs linked to sprint stories | Quality visibility for PM without reading code | 4h |
| 4 | `pm_flow_warning` — proactive alert when WIP > threshold or cycle time spikes | Flow metrics > velocity for prediction (research) | 3h |
| 5 | GitHub Actions/CI status → sprint health factor | Deploy failures = sprint health impact | 4h |

**Exit criteria:** PM has DORA-quality engineering visibility without asking developers for reports.

---

## Sprint +5: Team Autonomy (3-4 minggu)

Rationale: "The team gradually needs the SM less." SM role shift: 10% admin, 90% people.

| # | Task | Research Basis | Effort |
|---|------|---------------|--------|
| 1 | Sprint auto-scoring at close (no manual snapshot) | Automate the boring (principle #4) | 4h |
| 2 | Individual dev dashboard: "my flow, my debt, my blockers" | Self-service = autonomy (SCARF model) | 6h |
| 3 | Team self-health-check: team can run `pm` without SM involvement | Team ownership = maturity (Tuckman) | 2h |
| 4 | Working agreement enforcement alerts | Automated accountability | 4h |
| 5 | Onboarding guide auto-generation from team KB + decisions + DoD | New member productivity = team velocity | 4h |
| 6 | Maturity progress tracking: visualize team growth over time | SM proves value via team growth, not personal output | 3h |

**Exit criteria:** Team independently uses 5+ PM tools without SM prompting.

---

## Sprint +6: Scale + Ecosystem (4+ minggu)

| # | Task | Research Basis | Effort |
|---|------|---------------|--------|
| 1 | Docker image (one-command deploy) | Adoption barrier for non-Go teams | 4h |
| 2 | Remote MCP (Streamable HTTP) with auth | Enterprise teams need shared server | 8h |
| 3 | Multi-board aggregation (Scrum of Scrums view) | Enterprise PM needs | 8h |
| 4 | MCP marketplace/registry listing | Discoverability (ecosystem play) | 2h |
| 5 | Plugin system (custom tools per team) | Teams have unique workflows | 12h |

---

## Prioritization Logic

```
Priority = (Research Evidence × User Impact) / Effort

Critical path:
Security → Daily Usability → Smart Context → Communication → Developer Bridge → Autonomy → Scale
```

Setiap sprint delivers usable value. Nggak ada "big bang". PM bisa mulai pakai dari Sprint Next selesai.

---

## Sprint +7: OKR/KPI Intelligence (3-4 minggu)

Rationale: "Teams connecting goals to outcomes are 30% more likely to hit them" (OKR Tool, 2026). PM wastes hours/week manually translating Jira data to OKR progress. Lark OKR users doubly penalized — 2 systems, zero auto-sync.

| # | Task | Research Basis | Effort |
|---|------|---------------|--------|
| 1 | `pm_kpi_track` + `pm_kpi_dashboard` — simple time-series metric tracking | Foundation for all outcome measurement | 4h |
| 2 | `pm_okr_map` — link Key Results to Jira signals (JQL, epic, label) | Atlassian: "does this ticket move any KR?" | 6h |
| 3 | `pm_okr_progress` — auto-calculate KR progress from live Jira data | Eliminate "Friday afternoon slide-making" | 6h |
| 4 | `pm_okr_report` — generate OKR status per audience (team/PO/exec) | Different stakeholders need different granularity | 4h |
| 5 | `pm_okr_gaps` — identify KRs with zero linked sprint work | Strategy-execution gap detection | 2h |
| 6 | `pm_okr_suggest` — AI attribute sprint items to KRs | Output→Outcome bridge, 30% goal achievement improvement | 4h |
| 7 | `pm_okr_sync_lark` — push calculated progress to Lark OKR via API | Eliminate dual-system manual update | 6h |
| 8 | `pm_kpi_alert` — alert when KPI crosses threshold | Proactive vs reactive goal management | 3h |

**Exit criteria:** PM can answer "how does this sprint contribute to our quarterly OKRs?" in one tool call, and Lark OKR reflects reality without manual updates.

---

## Success Metrics

| Metric | Baseline | After 3 Sprints | After 6 Sprints |
|--------|----------|-----------------|-----------------|
| PM time on admin | ~15h/week | ~8h/week | ~4h/week |
| Time to generate exec report | 30-45 min | 30 seconds | 30 seconds |
| Sprint goal achievement | 52% (industry) | 65% | 75% |
| Blocker avg resolution time | unknown | tracked, improving | < 2 days |
| Forecast accuracy (at 85%) | unmeasured | within 30% | within 20% |
| Team autonomy score | unmeasured | tracked | improving trend |
| Stakeholder satisfaction | unmeasured | 3.5/5 | 4.2/5 |
