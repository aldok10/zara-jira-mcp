# Research: Jira ↔ OKR/KPI Integration Problem & Solutions

## The Core Problem

**"OKRs live in a doc, a slide deck, or a spreadsheet. Leadership reviews them quarterly. Jira lives in Jira."** [8]

Teams struggle to connect:
- Daily work (Jira tickets, sprints, story points) → **Outputs**
- Business goals (OKRs, KPIs) → **Outcomes**

The gap: **52% of Key Results written by teams are actually tasks or KPIs in disguise** [5]. Teams connect goals to outputs, not outcomes. Yet teams that connect goals to outcomes are **30% more likely to hit them** [4].

---

## Why Jira → OKR Is Hard

### 1. Output ≠ Outcome
- **Output:** "Shipped feature X" (Jira ticket done)
- **Outcome:** "User activation increased 15%" (business result)
- Jira tracks outputs. OKRs should track outcomes.
- You can't derive an outcome from a ticket status alone.

### 2. Aggregation Problem
- 50 tickets done ≠ Key Result progress
- Story points ≠ business value
- Sprint velocity ≠ strategy execution
- How many tickets = 1% OKR progress? Nobody knows.

### 3. Measurement Gap
- Jira: binary (done/not done) or status (open/in progress/done)
- OKR Key Results: need percentage, numeric value, or measurable change
- No natural mapping between "PROJ-123 Done" and "Reduce churn by 5%"

### 4. Timing Mismatch
- OKRs: quarterly (90 days)
- Sprints: 2 weeks
- There's a 6:1 ratio that requires manual bridging

### 5. Lark OKR Specific
- Lark OKR has its own system with objectives, key results, progress
- No native Jira integration
- Manual copy-paste between systems
- Progress updates require human interpretation

---

## Solutions & Approaches

### Approach A: Epic = Key Result (Recommended)

Map Jira hierarchy to OKR hierarchy:
```
Lark OKR Objective     ←→  Jira Epic (or Initiative)
Lark OKR Key Result    ←→  Jira Epic completion %
Sprint Work            ←→  Stories/Tasks under that Epic
```

**Auto-calculate KR progress:**
```
KR Progress = (Done issues in epic / Total issues in epic) × 100
```

Or weighted by story points:
```
KR Progress = (Done points / Total points) × 100
```

### Approach B: Label-Based Mapping

Tag Jira issues with OKR labels:
```
Labels: okr:Q3-growth, kr:activation-15pct
```

Then aggregate all issues with that label → KR progress.

### Approach C: Sprint Goal = KR Contribution

Each sprint goal maps to a Key Result:
```
Sprint Goal: "Ship onboarding flow" → contributes to KR: "Increase activation to 40%"
Sprint Goal achieved? → +X% to KR progress (manually defined weight)
```

### Approach D: Metric-Driven (Best for Outcomes)

Key Results measured by external metrics, Jira tracks the WORK to get there:
```
KR: "Reduce page load to <2s" — measured by monitoring tool
Jira: Epic "Performance optimization" — tracks the work
Connection: "This epic CONTRIBUTES to this KR" (not IS the KR)
```

---

## Proposed Tool Implementation

### `pm_okr_link` — Connect Jira work to OKRs

```
pm_okr_link(
  objective: "Improve developer experience",
  key_result: "Reduce onboarding time from 2 weeks to 3 days",
  jira_epic: "PROJ-100",          // Epic that contributes
  weight: 40,                      // This epic accounts for 40% of KR
  measurement: "epic_completion"   // or "story_points", "custom"
)
```

### `pm_okr_progress` — Auto-calculate OKR progress from Jira

Reads linked epics → calculates weighted progress:
```
Objective: "Improve developer experience" — 45% complete
  KR1: "Reduce onboarding time" — 60% (Epic PROJ-100: 6/10 done)
  KR2: "Zero setup issues" — 30% (Epic PROJ-200: 3/10 done)
  KR3: "NPS >8" — measured externally, manual: 7.5/8 = 94%
```

### `pm_okr_report` — Generate OKR status for leadership

Pulls progress from Jira + stored OKR definitions:
- What's on track / at risk / behind
- Connection between sprint work and strategic goals
- Formatted for Lark OKR update (copy-paste ready)

### `pm_kpi_derive` — Auto-derive KPIs from Jira data

Computable KPIs (no external tool needed):
```
Lead Time         = avg days from created to done
Cycle Time        = avg days from in-progress to done
Throughput        = items completed per sprint
Sprint Goal Rate  = goals achieved / total sprints
Defect Escape Rate = bugs in prod / total items shipped
Deploy Frequency  = releases per sprint
Carryover Rate    = items carried / items planned
Team Utilization  = items done / items assigned
```

### `pm_okr_sync_lark` — Push progress to Lark OKR (future)

If Lark OKR opens API for progress updates:
- Auto-update KR progress from Jira data
- Sync weekly or on-demand
- PM never manually updates Lark OKR again

---

## KPI Framework: What Can Be Auto-Derived from Jira

### Delivery KPIs
| KPI | Formula | Source |
|-----|---------|--------|
| Sprint Goal Success Rate | goals_achieved / total_sprints | pm_snapshot_sprint |
| Velocity Trend | points_done per sprint (3-sprint avg) | pm_velocity_trend |
| Throughput | items_done / sprint | pm_flow_metrics |
| Cycle Time | avg(started_to_done) in days | pm_flow_metrics |
| Lead Time | avg(created_to_done) in days | Jira timestamps |

### Quality KPIs
| KPI | Formula | Source |
|-----|---------|--------|
| Defect Density | bugs / total_items × 100 | pm_qa_health |
| Escape Rate | prod_bugs / total_shipped | need tracking |
| Carryover Rate | carryover / planned × 100 | pm_snapshot_sprint |
| Tech Debt Ratio | debt_items / total × 100 | pm_tech_debt_ratio |

### Team Health KPIs
| KPI | Formula | Source |
|-----|---------|--------|
| WIP per Person | in_progress / team_size | pm_overload_check |
| Blocker Resolution Time | avg(days_blocked) | pm_blockers |
| Action Item Completion | done_actions / total_actions | pm_action_items |
| Team Confidence | avg(confidence_scores) | pm_confidence |

### Predictability KPIs
| KPI | Formula | Source |
|-----|---------|--------|
| Commitment Reliability | done / committed × 100 | pm_commitment_check |
| Velocity Variance | std_dev / mean × 100 | pm_forecast |
| Forecast Accuracy | actual vs predicted | need tracking |

---

## Lark OKR Integration Path

### Current State
- Lark OKR is a standalone feature within Lark Suite
- No public REST API for OKR CRUD (as of 2026)
- Progress updates are manual in Lark OKR UI
- Lark does have: messaging API, docs API, calendar API

### Possible Integrations Now
1. **Lark Message → OKR Update Reminder**: Send formatted progress to OKR owner via Lark bot
2. **Lark Doc → OKR Dashboard**: Auto-generate OKR progress doc in Lark Docs
3. **Manual Bridge**: `pm_okr_report` generates copy-paste text for Lark OKR updates

### Future (When Lark OKR API opens)
1. Auto-sync KR progress from Jira epic completion
2. Bi-directional: OKR changes → update Jira epic priority
3. Full dashboard: OKR tree with live Jira data

---

## TODO: OKR/KPI Tools to Build

### P0 (Must Have)
- [ ] `pm_okr_define` — Define OKR with Jira links (stored in SQLite)
- [ ] `pm_okr_progress` — Auto-calculate progress from linked Jira epics
- [ ] `pm_kpi_dashboard` — Auto-derived KPIs from Jira data (no manual input)
- [ ] `pm_okr_report` — Formatted OKR progress for leadership/Lark

### P1 (Should Have)
- [ ] `pm_okr_link` — Link additional epics/issues to KRs
- [ ] `pm_kpi_trend` — KPI trends over sprints (improving/declining)
- [ ] `pm_okr_risk` — Flag KRs that are falling behind schedule
- [ ] Lark message integration for OKR reminders

### P2 (Nice to Have)
- [ ] `pm_okr_suggest` — AI suggests KRs based on sprint patterns
- [ ] Lark Doc auto-generation with OKR progress
- [ ] KPI export to spreadsheet format
- [ ] Historical OKR cycle comparison

---

## Key Insight

**The solution is NOT to make Jira do OKRs.** The solution is:

1. **Define OKRs in the MCP** (simple: objective + key results + measurement method)
2. **Link to Jira** (epic = KR, issues = work contributing to KR)
3. **Auto-calculate** (progress derived from Jira status, no manual update)
4. **Format for audience** (Lark OKR update, executive report, team dashboard)

This way:
- PM defines OKRs once
- Jira tracks execution naturally (tickets, sprints, epics)
- MCP bridges the gap automatically
- Lark OKR gets formatted updates ready to paste
- No double-entry, no manual percentage calculation
