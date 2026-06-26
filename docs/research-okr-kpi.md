# Research: Jira → OKR/KPI Integration Problem & Solutions

> The #1 pain point: Sprint work lives in Jira. Strategy lives in OKR tools. They never connect.

---

## The Problem (Deeply Understood)

### Why Jira Data Doesn't Map Cleanly to OKRs

1. **Abstraction mismatch**: Jira tracks TASKS (output). OKRs track OUTCOMES (impact).
   - Completing PROJ-123 "Add login page" ≠ KR "Increase user onboarding rate by 20%"
   - The gap between "done ticket" and "achieved result" is where all OKR tools fail

2. **Granularity mismatch**: One Key Result may span 3 epics across 2 teams and 6 sprints.
   - Jira hierarchy: Issue → Epic → Sprint (work decomposition)
   - OKR hierarchy: Company Objective → Team Objective → Key Result (outcome decomposition)
   - These are DIFFERENT hierarchies. Mapping is many-to-many, not one-to-one.

3. **Measurement mismatch**: 
   - Jira measures: story points, velocity, cycle time (activity metrics)
   - OKRs measure: revenue, NPS, conversion rate, uptime (business metrics)
   - Completing 100% of sprint issues ≠ 100% OKR progress

4. **Temporal mismatch**:
   - Sprints: 2 weeks
   - OKR cycles: quarterly
   - The cadence doesn't align. Progress is not linear.

5. **The "People Problem" (Atlassian Community)**:
   - "The strategy-execution gap isn't a Jira problem. It's a cadence and accountability problem." [1]
   - OKRs live in docs/slides. Jira lives in Jira. They drift apart within days.
   - Teams review OKRs quarterly but update Jira daily. Gap grows exponentially.

### Lark OKR Specific Challenges

1. **Lark OKR is separate from Lark's other tools** — no native Jira connection
2. **Manual progress updates** — PMs must manually calculate % and update Lark OKR
3. **No formula for "Jira done items → KR progress"** — it's always a judgment call
4. **API exists but underdocumented** — Feishu OKR API supports CRUD but few use it programmatically
5. **Language barrier** — most Lark OKR API docs are in Chinese

---

## The Solution Framework

### Principle: Don't Map Tasks to OKRs. Map SIGNALS to Key Results.

Instead of: "If PROJ-123 is done, update KR progress by 5%"
Do this: "KR: Reduce cycle time to <3 days. SIGNAL: Calculate actual cycle time from Jira. Auto-update."

### Three Types of Key Results and How to Auto-Calculate

| KR Type | Example | Jira Signal | Formula |
|---------|---------|-------------|---------|
| **Completion-based** | "Ship auth feature" | Epic % done | `done_issues / total_issues * 100` |
| **Metric-based** | "Reduce cycle time to 3d" | Flow metrics | `avg(resolved_date - created_date)` from Jira |
| **Count-based** | "Close 50 customer bugs" | JQL count | `count(type=Bug AND resolved >= startOfQuarter())` |

### The Auto-Progress Engine

```
┌─────────────┐     ┌──────────────────┐     ┌─────────────┐
│  Lark OKR   │ ←── │  zara-jira-mcp   │ ──→ │  Jira Cloud │
│  (target)   │     │  (bridge/calc)   │     │  (source)   │
└─────────────┘     └──────────────────┘     └─────────────┘
      ↕                      ↕                      ↕
  Objectives          JQL → calculate         Issues, Sprints
  Key Results         signal → progress       Epics, Velocity
  Progress %          push → Lark OKR API     Cycle Time
```

---

## Lark OKR API (What's Available)

From Feishu/Lark Open Platform:

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/open-apis/okr/v1/periods` | GET | List OKR cycles (Q1, Q2, etc.) |
| `/open-apis/okr/v1/users/{user_id}/okrs` | GET | Get user's OKRs |
| `/open-apis/okr/v1/okrs/{okr_id}` | GET | Get specific OKR details |
| `/open-apis/okr/v1/progress_records` | POST | Create progress record |
| `/open-apis/okr/v1/progress_records/{id}` | PUT | Update progress record |
| `/open-apis/okr/v1/images/upload` | POST | Upload image for progress |

**Key capability**: We can READ OKRs and WRITE progress updates programmatically.

**SDK**: `github.com/larksuite/oapi-sdk-go/v3` already in our go.mod — same SDK we use for messaging.

---

## KPI Calculation from Jira (DORA + Custom)

### DORA Metrics (Calculable from Jira + Git)

| Metric | How to Calculate from Jira | Elite Benchmark |
|--------|---------------------------|-----------------|
| **Deployment Frequency** | Count issues transitioned to "Done" per day/week | On-demand |
| **Lead Time for Changes** | `avg(done_date - created_date)` for completed issues | <1 day |
| **Change Failure Rate** | `bugs_in_sprint / total_done_in_sprint * 100` | <5% |
| **Recovery Time** | `avg(resolved_date - created_date)` for P1 bugs only | <1 hour |

### Engineering KPIs from Jira Data

| KPI | JQL / Calculation |
|-----|-------------------|
| Sprint Predictability | `done_items / committed_items * 100` |
| Carryover Rate | `carryover_items / total_sprint_items * 100` |
| Bug Escape Rate | `bugs_found_in_prod / total_stories_shipped` |
| Team Utilization | `assigned_items / available_capacity` |
| Backlog Health | `items_without_estimate / total_backlog * 100` |
| Story Cycle Time | `avg(done_date - in_progress_date)` |
| Blocked Time Ratio | `time_in_blocked / total_cycle_time` |
| Defect Density | `bugs / story_points_shipped` |

---

## Proposed Tools for zara-jira-mcp

### Phase A: OKR Bridge (Jira → OKR Progress)

| Tool | Purpose |
|------|---------|
| `okr_define_signal` | Define how a KR maps to Jira data (JQL + formula) |
| `okr_calculate_progress` | Run the signal, calculate KR progress percentage |
| `okr_sync_to_lark` | Push calculated progress to Lark OKR API |
| `okr_status_report` | Show all KRs with auto-calculated progress from Jira |
| `okr_alignment_check` | What % of sprint work connects to an OKR? |

### Phase B: KPI Dashboard

| Tool | Purpose |
|------|---------|
| `kpi_dora_metrics` | Calculate DORA 4 from Jira data |
| `kpi_sprint_predictability` | Committed vs delivered ratio |
| `kpi_team_kpis` | Composite team health KPIs |
| `kpi_trend` | KPI trends over time (improving/declining) |
| `kpi_report` | Executive KPI report (formatted for management) |

### Phase C: Lark OKR Full Integration

| Tool | Purpose |
|------|---------|
| `lark_okr_list` | List current OKR cycle objectives |
| `lark_okr_progress` | View KR progress |
| `lark_okr_update` | Push progress update with comment |
| `lark_okr_sync_all` | Batch sync all Jira-linked KRs |

---

## Implementation Plan

### Step 1: Signal Definition (in-memory, no Lark API needed)

Store OKR-to-Jira mappings in SQLite:
```sql
CREATE TABLE okr_signals (
  id INTEGER PRIMARY KEY,
  objective TEXT NOT NULL,
  key_result TEXT NOT NULL,
  signal_type TEXT NOT NULL, -- completion, metric, count
  jql TEXT NOT NULL,         -- JQL to run against Jira
  formula TEXT NOT NULL,     -- how to calculate: count, avg_cycle_time, pct_done
  target_value REAL,         -- what 100% looks like
  current_value REAL,        -- last calculated value
  lark_kr_id TEXT,           -- Lark OKR key result ID (for sync)
  updated_at DATETIME
);
```

### Step 2: Auto-Calculation Engine

```go
// Signal formulas:
// "count" → count issues matching JQL
// "pct_done" → done/total from JQL results
// "avg_cycle_time" → avg days from created to resolved
// "sum_points" → total story points from JQL
// "custom" → raw value from a specific field
```

### Step 3: Lark OKR API Integration

Using existing `larksuite/oapi-sdk-go/v3`:
```go
// Already in go.mod! Just need to call OKR endpoints.
// POST /open-apis/okr/v1/progress_records
// Body: { okr_id, kr_id, content: "Auto-synced from Jira: 75% (18/24 stories done)" }
```

---

## Key Insight: The Formula Library

The breakthrough is giving PM/SM a **pre-built formula library** so they don't have to figure out HOW to map:

| I want to track... | Use formula | Signal JQL |
|--------------------|-------------|------------|
| Feature delivery | `pct_done` | `"Epic Link" = EPIC-123 AND resolution = Done` |
| Bug reduction | `count_inverse` | `type = Bug AND created >= -30d` |
| Cycle time improvement | `avg_cycle_time` | `resolved >= -14d AND type = Story` |
| Sprint predictability | `ratio` | `sprint = X` (done/committed) |
| Customer issues resolved | `count` | `type = Bug AND labels = customer AND resolved >= startOfQuarter()` |
| Deployment frequency | `count_per_day` | `status changed to Done AFTER -7d` |
| Quality | `inverse_ratio` | bugs/stories shipped |

---

## What Makes This HARD (and Our Solution)

| Challenge | Why It's Hard | Our Solution |
|-----------|--------------|--------------|
| Jira tasks ≠ outcomes | Completing tasks doesn't prove value | Signal-based: measure the METRIC, not the task |
| Many-to-many mapping | 1 KR = multiple epics/teams | JQL is flexible enough to query across boundaries |
| Manual updates kill adoption | People forget to update OKR tool | Auto-calculate + auto-sync (push to Lark OKR) |
| Different cadences | OKR quarterly, Sprint biweekly | Calculate on-demand, store history |
| Business metrics not in Jira | Revenue, NPS live elsewhere | Hybrid: auto for engineering KRs, manual input for business KRs |
| Subjectivity | "Is this KR really 60% done?" | Formula removes subjectivity. Data = data. |

---

## Priority Order for Implementation

### Tomorrow (P0):
1. `okr_define_signal` + `okr_calculate_progress` — the core value prop
2. `kpi_sprint_predictability` — simplest meaningful KPI
3. SQLite schema for signals

### This Week (P1):
4. `okr_status_report` — show all KRs with auto-progress
5. `kpi_dora_metrics` — DORA 4 from Jira
6. `okr_alignment_check` — what % of work connects to OKR
7. `lark_okr_update` — push progress to Lark

### Next Sprint (P2):
8. `lark_okr_list` / `lark_okr_progress` — full Lark OKR read
9. `lark_okr_sync_all` — batch auto-sync
10. `kpi_trend` + `kpi_report` — historical tracking + exec report

---

## Sources

[1] https://community.atlassian.com/forums/App-Central-articles/The-OKR-Jira-Gap-Is-a-People-Problem-Not-a-Tooling-Problem/ba-p/3220406
[2] https://oboard.io/blog/how-to-implement-okrs-in-jira-a-practical-guide-for-beginners
[3] https://medium.com/@margosakova/jira-vs-jira-vs-okr-board-for-jira-8cbb01ad9454
[4] https://www.larksuite.com/hc/en-US/articles/854393465133
[5] https://o-mega.ai/skills/lark-okr (Lark OKR API reference)
[6] https://appfire.atlassian.net/wiki/spaces/OFJ/pages/2540699673/Auto+KRs (auto progress calculation)
[7] https://www.cortex.io/post/what-is-slowing-progress-towards-engineering-okrs-survey
[8] https://www.atlassian.com/blog/confluence/okr-jira-confluence
[9] https://www.taskade.com/blog/dora-metrics-explained (DORA 2024 benchmarks)
[10] https://deepwiki.com/larksuite/oapi-sdk-go (Lark Go SDK)
