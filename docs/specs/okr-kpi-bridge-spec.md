# SPEC: OKR/KPI Bridge — Jira Data to Business Outcomes

Status: Draft (see gap-analysis-addendum.md for corrections)
Author: Zara (AI) + Aldo
Date: 2026-06-26
Priority: High
Depends on: Lark SDK v3.9.7 (already in go.mod), existing pm_outcome_map

> **IMPORTANT:** This spec was written before discovering existing OKR v1/v2 systems.
> See `docs/specs/gap-analysis-addendum.md` for what already exists and revised scope.
> Actual build scope: 7 tools (not 11). Existing: pm_okr_tree, pm_okr_sync, pm_okr_kpi,
> pm_kpi_dashboard. Lark OKR read client exists at internal/lark/okr.go.

---

## The Core Problem

### Why Jira Data Doesn't Convert to OKR/KPI Naturally

**The fundamental mismatch:**

| Jira tracks... | OKR/KPI needs... |
|---|---|
| Tasks (output) | Impact (outcome) |
| Story points | Business metrics |
| Sprint velocity | Progress toward objectives |
| Issue status | Key Result achievement % |
| Who did what | Whether it moved the needle |

Research confirms this is universal:

1. **Sprint cadence trains output thinking** — teams write Key Results that track work completed, not impact created. "Ship retention dashboard" is an output. "Reduce churn by 20%" is an outcome. [5](https://www.okrstool.com/blog/okrs-agile)

2. **The connection gap** — "Nobody could confidently answer: does this ticket move any of our key results?" — Atlassian Community [4](https://community.atlassian.com/forums/Advanced-Planning-in-Jira/OKR-Jira-Alignment-Gap-Strategy-Disconnected-from-Execution/td-p/3224657)

3. **30% hit rate gap** — teams that connect goals to outcomes (not outputs) are 30% more likely to achieve them [1](https://www.okrstool.com/blog/okrs-agile)

4. **Measurement fragmentation** — OKR lives in Lark OKR / Weekdone / spreadsheet. Work lives in Jira. Nobody maintains the bridge manually.

### The Lark OKR Specific Pain

For teams using Lark OKR (our primary usecase):
- OKRs are set quarterly in Lark OKR
- Sprint work happens in Jira
- **Nobody updates Lark OKR progress** because it requires manual translation
- End of quarter: "How did we do?" → scramble to reconstruct from Jira
- Progress records in Lark OKR are empty or outdated
- Leaders can't see if sprints actually serve the stated objectives

---

## Current State in zara-jira-mcp

**What exists:**
- `pm_outcome_map` — manually map sprint → objective (stored locally in SQLite)
- `pm_outcome_history` — show past mappings
- `pm_set_sprint_goal` / `pm_goal_check` — sprint-level goals (not OKR-level)
- Lark SDK v3.9.7 with full OKR module (`service/okr/v1`) — **unused**

**What's missing:**
1. No Lark OKR integration (SDK ready but not connected)
2. No automatic progress calculation from Jira data
3. No KPI derivation from sprint metrics
4. No "contribution analysis" (which tickets actually moved a KR)
5. No bi-directional sync (Lark OKR ↔ local memory)

---

## Research: How to Bridge Output → Outcome

### The Translation Layer Model

The key insight from research: you need an **explicit translation layer** between task completion and outcome measurement.

```
JIRA (output)              TRANSLATION LAYER              OKR (outcome)
─────────────────────    ──────────────────────────    ─────────────────────
Tickets completed    →   Contribution mapping       →   KR progress %
Sprint velocity      →   Velocity-to-capacity       →   Delivery confidence
Blocked items        →   Risk-to-objective impact   →   Objective at-risk flag
Epic completion %    →   Feature-to-KR mapping      →   KR progress update
Cycle time           →   Efficiency KPI             →   Team health metric
```

### Three Types of Key Results (and how Jira data feeds each)

| KR Type | Example | Jira Signal |
|---------|---------|-------------|
| **Metric KR** | "Reduce p95 latency to < 200ms" | External metric — Jira can't measure this directly. Track whether enabling work was completed. |
| **Milestone KR** | "Ship payment v2 to production" | Epic/version completion % from Jira. Binary: shipped or not. |
| **Activity KR** | "Complete 3 customer interviews" | Specific issue/task completion. Direct 1:1 mapping. |

**Practical approach:**
- Milestone KRs: auto-calculate from Jira epics/versions (highest value, easiest)
- Activity KRs: direct issue linkage (1:1 mapping)
- Metric KRs: manual input + Jira tracks enabling work contribution

### DORA/SPACE as Engineering KPIs

For engineering teams, these are the proven KPIs derivable from Jira data:

| KPI | Source | Formula from Jira |
|-----|--------|------------------|
| **Throughput** | Issues done / time | `count(status=Done) / sprint_days` |
| **Cycle Time** | Start → Done duration | `avg(done_date - in_progress_date)` |
| **Predictability** | Planned vs Delivered | `delivered / committed * 100` |
| **WIP Ratio** | Items in flight | `count(status=InProgress) / team_size` |
| **Blocker Rate** | % time blocked | `blocked_days / total_sprint_days` |
| **Carryover Rate** | Incomplete → next sprint | `carryover / committed * 100` |
| **Quality Signal** | Bug escape rate | `bugs_in_prod / items_shipped` |

These map naturally to OKR Key Results:
- "Improve delivery predictability to > 80%" → measured by sprint completion rate
- "Reduce cycle time by 30%" → measured by avg done-to-start duration
- "Zero critical production incidents" → measured by bug escape rate

---

## Proposed Architecture

### Design Principle: Output → Outcome Translation Engine

```
┌────────────────────────────────────────────────────────┐
│                    USER (PM/SM)                          │
│   "How are we doing on our OKRs this quarter?"          │
└─────────────────────────┬──────────────────────────────┘
                          │
┌─────────────────────────▼──────────────────────────────┐
│              TRANSLATION ENGINE                          │
│                                                         │
│  ┌──────────┐  ┌──────────────┐  ┌─────────────────┐  │
│  │  Jira    │  │  KR Mapping  │  │  KPI Calculator  │  │
│  │  Reader  │  │  (epic→KR)   │  │  (DORA/flow)    │  │
│  └────┬─────┘  └──────┬───────┘  └────────┬────────┘  │
│       │               │                    │            │
│  ┌────▼───────────────▼────────────────────▼────────┐  │
│  │           Progress Synthesizer (AI)               │  │
│  │  "Epic 70% done → KR1 at 70% → update Lark OKR" │  │
│  └──────────────────────┬───────────────────────────┘  │
│                         │                               │
└─────────────────────────┼───────────────────────────────┘
                          │
              ┌───────────▼───────────┐
              │                       │
     ┌────────▼──────┐    ┌─────────▼────────┐
     │  Local SQLite  │    │   Lark OKR API   │
     │  (memory)      │    │   (sync target)  │
     └────────────────┘    └──────────────────┘
```

---

## Implementation Spec (4 Phases)

### Phase 1: KPI Engine — Derive Engineering KPIs from Jira (Week 1-2)

No Lark dependency. Pure computation from existing data.

#### Tool: `pm_kpi_dashboard`

Auto-calculate engineering KPIs from sprint snapshots + Jira data.

```
pm_kpi_dashboard(board_id: int, sprints: int) -> KPIDashboard
```

**Logic:**
1. Pull last N sprint snapshots from memory
2. Calculate each KPI from stored data
3. Show trend (improving/stable/declining)
4. Compare to benchmarks

**Output:**
```
Engineering KPI Dashboard (last 5 sprints)
Board: SIT Board

DELIVERY
  Throughput:      4.2 items/day (↑ from 3.8)         [Good]
  Predictability:  72% (target: 80%)                   [Watch]
  Carryover Rate:  18% (↓ from 25%)                    [Improving]

FLOW
  Cycle Time:      6.2 days avg (target: < 5)          [Over]
  WIP Ratio:       2.8 per person (target: < 2)        [Over]
  Blocker Rate:    12% of sprint time blocked           [Watch]

QUALITY
  Bug Escape:      0.3 per release (↓ from 0.8)        [Good]

TREND: 4/7 KPIs improving. Focus area: WIP management and cycle time.
RECOMMENDATION: Reduce WIP limits. Current 2.8 items/person causes context switching.
```

**Implementation:**
- File: `application/tools/kpi_handlers.go`
- Data source: existing `sprint_snapshots` table + Jira API (for cycle time)
- No AI required (pure math)
- No new tables

---

#### Tool: `pm_kpi_trend`

Show single KPI trend over time with sparkline.

```
pm_kpi_trend(board_id: int, kpi: string, sprints: int) -> Trend
```

`kpi` values: `throughput`, `predictability`, `cycle_time`, `wip`, `blocker_rate`, `carryover`, `bug_escape`

---

### Phase 2: OKR Mapping Engine — Connect Work to Objectives (Week 3-4)

#### Tool: `pm_okr_link`

Link a Jira epic/version to a specific Key Result. The fundamental bridge.

```
pm_okr_link(
  jira_ref: string,       // epic key (e.g. "PROJ-100") or version name
  kr_id: string,          // Key Result identifier (local or Lark OKR ID)
  kr_description: string, // "Reduce p95 latency to <200ms"
  contribution: string,   // "enabling" | "direct" | "partial"
  weight: float           // 0.0-1.0 how much this work contributes to the KR
) -> Confirmation
```

**DB Schema (new table):**
```sql
CREATE TABLE IF NOT EXISTS okr_links (
  id INTEGER PRIMARY KEY,
  jira_ref TEXT NOT NULL,          -- epic key or version
  jira_ref_type TEXT DEFAULT 'epic', -- epic, version, label, jql
  kr_id TEXT NOT NULL,             -- local ID or Lark OKR KR ID
  kr_description TEXT,
  objective TEXT,                   -- parent objective text
  contribution_type TEXT DEFAULT 'direct', -- enabling, direct, partial
  weight REAL DEFAULT 1.0,
  lark_okr_id TEXT,                -- Lark OKR objective ID (for sync)
  lark_kr_id TEXT,                 -- Lark OKR key result ID (for sync)
  cycle TEXT,                      -- "Q3 2026"
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

#### Tool: `pm_okr_progress`

Calculate OKR progress from linked Jira work. The auto-translation.

```
pm_okr_progress(
  cycle: string,          // optional: "Q3 2026" (default: current)
  board_id: int           // optional: filter by board
) -> OKRProgressReport
```

**Logic:**
1. Get all `okr_links` for the cycle
2. For each link, query Jira for epic/version completion %
3. Calculate weighted KR progress: `sum(epic_completion * weight) / sum(weight)`
4. Group by Objective
5. AI generates narrative summary

**Output:**
```
OKR Progress Report — Q3 2026

OBJECTIVE 1: Improve platform reliability
  Overall: 62% (on track for 85% confidence)

  KR1: "Reduce p95 latency to <200ms" [45%]
    - PROJ-100 (Gateway optimization): 80% complete [direct, weight 0.6]
    - PROJ-112 (Redis cache layer): 30% complete [enabling, weight 0.4]
    - Weighted progress: 80*0.6 + 30*0.4 = 60% → KR at 45% (needs external metric validation)

  KR2: "Zero critical incidents in production" [100%]
    - Current: 0 critical incidents this quarter
    - Tracked via: PROJ bug issues with priority=Critical, fixVersion=Q3

  KR3: "Ship monitoring dashboard v2" [70%]
    - PROJ-150 (Monitoring epic): 70% complete [direct, weight 1.0]

OBJECTIVE 2: Accelerate delivery speed
  Overall: 55%

  KR1: "Cycle time < 5 days" [40%]
    - Current: 6.2 days (was 8.1 at quarter start)
    - Progress: (8.1-6.2)/(8.1-5.0) = 61% toward target

  KR2: "Predictability > 80%" [68%]
    - Last 3 sprints: 72%, 75%, 68%
    - Trend: volatile, not yet at 80% consistently

AT RISK:
  - KR1 Obj1 needs acceleration. Gateway epic must complete by sprint 4.
  - KR1 Obj2 backsliding. Cycle time increased last sprint due to WIP overload.
```

---

#### Tool: `pm_okr_suggest`

AI suggests which OKRs current sprint work might serve (for teams that haven't linked yet).

```
pm_okr_suggest(board_id: int) -> Suggestions
```

**Logic:**
1. Get current sprint issues (summaries, epics, labels)
2. Get known objectives from `okr_links` or `outcome_map`
3. AI matches work to objectives
4. Suggest linkages

**Output:**
```
OKR Alignment Suggestions:

Sprint items without OKR linkage:
  PROJ-200 "Implement rate limiting" → likely serves "Improve platform reliability"
  PROJ-205 "Add Prometheus metrics" → likely serves "Ship monitoring dashboard v2"
  PROJ-210 "Fix login timeout" → likely serves "Zero critical incidents" (bug fix)

Unlinked items: 3/15 (80% alignment score)

Action: Run pm_okr_link(jira_ref:"PROJ-200", kr_id:"KR-reliability-1", contribution:"enabling")
```

---

### Phase 3: Lark OKR Integration — Bi-directional Sync (Week 5-6)

Leverages `github.com/larksuite/oapi-sdk-go/v3/service/okr/v1` (already in go.mod).

#### Tool: `lark_okr_pull`

Fetch OKRs from Lark OKR and store locally for mapping.

```
lark_okr_pull(
  user_id: string,        // Lark user ID (or "me")
  cycle: string           // optional: filter by period
) -> OKRList
```

**Logic:**
1. Call `userOkr.List()` to get user's OKRs
2. For each OKR, get objectives + key results
3. Store in local `okr_links` table with `lark_okr_id` and `lark_kr_id`
4. Display structured list

**Required config:**
```
LARK_APP_ID=xxx          # already configured
LARK_APP_SECRET=xxx      # already configured
# New: app needs okr:okr.content:readonly scope
```

---

#### Tool: `lark_okr_sync_progress`

Push calculated progress back to Lark OKR as progress records.

```
lark_okr_sync_progress(
  cycle: string,          // optional
  dry_run: bool           // default true — show what would be updated
) -> SyncResult
```

**Logic:**
1. Run `pm_okr_progress` logic internally
2. For each KR with `lark_kr_id` set:
   - Call `progressRecord.Create()` with markdown summary
   - Include: Jira-derived progress %, contributing epics, trend
3. In dry_run mode: show what would be posted, don't post

**Output (dry_run=true):**
```
Lark OKR Sync Preview:

Would update 4 Key Results:
  KR "Reduce latency" → progress note: "60% (Gateway 80% + Cache 30%, weighted)"
  KR "Zero incidents" → progress note: "100% (0 critical incidents Q3)"
  KR "Ship monitoring v2" → progress note: "70% (Epic PROJ-150 at 70%)"
  KR "Cycle time <5d" → progress note: "40% toward target (6.2d, was 8.1d)"

Run with dry_run:false to push to Lark OKR.
```

**Safety:**
- Default dry_run=true (never auto-push without user confirmation)
- Creates progress RECORDS (additive), never overwrites OKR structure
- Preserves human-written progress notes

---

#### Tool: `lark_okr_periods`

List available OKR periods/cycles from Lark.

```
lark_okr_periods() -> PeriodList
```

Simple wrapper around `period.List()`.

---

### Phase 4: Smart OKR Assistant (Week 7-8)

#### Tool: `pm_okr_health`

Unified health check: are OKRs on track based on Jira data?

```
pm_okr_health(cycle: string) -> OKRHealthReport
```

**Logic:**
1. Calculate time elapsed in quarter (e.g. 60% through Q3)
2. Calculate weighted OKR progress (e.g. 45%)
3. If progress < time_elapsed * 0.8: flag as AT RISK
4. For at-risk KRs: identify which Jira work is lagging
5. AI generates recommended actions

**Output:**
```
OKR Health — Q3 2026 (60% of quarter elapsed)

OBJECTIVE 1: Improve platform reliability [62%] ✓ On Track
OBJECTIVE 2: Accelerate delivery speed [45%] ⚠ At Risk

AT RISK DETAILS:
  KR "Cycle time < 5d" — Progress 40% vs expected 48%
    Cause: WIP overload last 2 sprints (2.8 items/person)
    Fix: Reduce sprint commitment by 2 items next sprint

  KR "Predictability > 80%" — Progress 68% but volatile
    Cause: Scope additions mid-sprint (detected 3 times)
    Fix: Enforce scope freeze after planning, use pm_scope_creep alerts

ACTIONS:
1. Reduce next sprint scope to improve focus (KR: cycle time)
2. Activate pm_scope_creep alerts for PO (KR: predictability)
3. Complete Gateway epic (PROJ-100) this sprint — critical path for KR: latency
```

---

#### Tool: `pm_okr_report`

Generate quarterly OKR report for leadership. Bridges the gap between Jira activity and business communication.

```
pm_okr_report(
  cycle: string,
  audience: string    // "leadership", "team", "po"
) -> FormattedReport
```

**For leadership:** Business outcomes language. No Jira keys. "We reduced system latency by 25% toward our 50% target."
**For team:** Specific tickets and what moved. "PROJ-100 was the biggest contributor to latency KR."
**For PO:** Goal progress + what tradeoffs were made. "We deprioritized feature X to focus on reliability KR."

---

#### Tool: `pm_kpi_to_okr`

Suggest how to convert current KPI data into well-formed Key Results.

```
pm_kpi_to_okr(board_id: int) -> Suggestions
```

**Logic (AI-powered):**
1. Calculate current KPIs (throughput, cycle time, predictability, etc.)
2. AI suggests OKR-formatted Key Results based on areas needing improvement
3. Include baseline (current) and realistic target

**Output:**
```
Suggested Key Results (based on current metrics):

From CYCLE TIME (current: 6.2 days, industry benchmark: 3-5 days):
  → "Reduce average cycle time from 6.2 to 4.5 days by end of Q4"

From PREDICTABILITY (current: 72%, target: 80%):
  → "Achieve sprint completion rate of 80%+ for 4 consecutive sprints"

From WIP (current: 2.8/person, healthy: < 2):
  → "Maintain WIP limit of 2 items per developer throughout Q4"

From CARRYOVER (current: 18%, was 25%):
  → "Reduce sprint carryover to below 10% by end of Q4"

These are OUTCOME Key Results (not output). They measure team capability improvement,
not "number of tickets completed."
```

---

## Architecture Notes

### File Structure
```
application/tools/
  kpi_handlers.go           # Phase 1: pm_kpi_dashboard, pm_kpi_trend
  okr_bridge_handlers.go    # Phase 2: pm_okr_link, pm_okr_progress, pm_okr_suggest
  lark_okr_handlers.go      # Phase 3: lark_okr_pull, lark_okr_sync_progress, lark_okr_periods
  okr_intelligence_handlers.go  # Phase 4: pm_okr_health, pm_okr_report, pm_kpi_to_okr

internal/lark/
  okr_client.go             # Lark OKR API wrapper (thin layer over SDK)

transport/
  okr.go                    # Tool registration (new file)
```

### DB Changes
```sql
-- 1 new table (Phase 2)
CREATE TABLE IF NOT EXISTS okr_links (
  id INTEGER PRIMARY KEY,
  jira_ref TEXT NOT NULL,
  jira_ref_type TEXT DEFAULT 'epic',
  kr_id TEXT NOT NULL,
  kr_description TEXT,
  objective TEXT,
  contribution_type TEXT DEFAULT 'direct',
  weight REAL DEFAULT 1.0,
  lark_okr_id TEXT,
  lark_kr_id TEXT,
  cycle TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Index for cycle-based queries
CREATE INDEX IF NOT EXISTS idx_okr_links_cycle ON okr_links(cycle);
CREATE INDEX IF NOT EXISTS idx_okr_links_jira_ref ON okr_links(jira_ref);
```

### Config Changes
```
# Existing (no change needed for Phase 1-2):
LARK_APP_ID=xxx
LARK_APP_SECRET=xxx

# For Phase 3 (Lark OKR sync):
# App needs additional scope: okr:okr.content:readonly, okr:okr.content:writeonly
# No new env vars — uses existing Lark credentials
```

### Module Assignment
- Phase 1-2: module `pm` (available from `lite` profile up)
- Phase 3-4: module `integrations` (available in `full` and `all` profiles)

---

## Key Design Decisions

### 1. Translation, Not Replacement

We do NOT replace Lark OKR. We bridge:
- OKR definition stays in Lark OKR (quarterly cadence, leadership-driven)
- Work execution stays in Jira (sprint cadence, team-driven)
- zara-jira-mcp **translates** between them automatically

### 2. Outcome > Output Focus

Every KPI and OKR tool emphasizes *outcome metrics* over *output metrics*:
- Bad KR: "Complete 15 stories per sprint" (output)
- Good KR: "Reduce cycle time to < 5 days" (outcome = team capability)
- The `pm_kpi_to_okr` tool explicitly helps teams make this shift

### 3. Conservative Sync (Safety First)

Lark OKR sync is:
- **Read-heavy, write-light** — mostly pull OKRs, occasionally push progress records
- **Additive only** — creates progress records, never modifies OKR structure
- **Dry-run default** — never pushes without explicit user confirmation
- **Idempotent** — running sync twice doesn't duplicate progress records

### 4. Works Without Lark OKR

Phase 1-2 work independently of Lark OKR:
- Local `okr_links` table stores mappings
- `pm_okr_progress` calculates from Jira data alone
- Teams without Lark OKR still get full value

Phase 3-4 add Lark OKR integration as optional enhancement.

---

## New Tools Summary (10 total)

| Phase | Tool | Complexity | AI | Lark OKR |
|-------|------|-----------|-----|----------|
| 1 | `pm_kpi_dashboard` | Medium | No | No |
| 1 | `pm_kpi_trend` | Low | No | No |
| 2 | `pm_okr_link` | Low | No | No |
| 2 | `pm_okr_progress` | High | Yes (narrative) | No |
| 2 | `pm_okr_suggest` | Medium | Yes | No |
| 3 | `lark_okr_pull` | Medium | No | Yes (read) |
| 3 | `lark_okr_sync_progress` | High | No | Yes (write) |
| 3 | `lark_okr_periods` | Low | No | Yes (read) |
| 4 | `pm_okr_health` | Medium | Yes | No |
| 4 | `pm_okr_report` | Medium | Yes | No |
| 4 | `pm_kpi_to_okr` | Medium | Yes | No |

---

## Success Criteria

1. PM can answer "how are our OKRs doing?" in < 5 seconds (vs manual reconstruction)
2. Lark OKR progress records update automatically from Jira data (dry_run → confirm → push)
3. KPI dashboard shows meaningful engineering metrics without manual calculation
4. `pm_okr_suggest` correctly identifies OKR alignment for > 70% of sprint items
5. `pm_kpi_to_okr` produces well-formed, outcome-focused Key Results (not output KRs)
6. Zero breaking changes to existing `pm_outcome_map` / `pm_outcome_history` tools

---

## Migration from pm_outcome_map

Existing `pm_outcome_map` is preserved but enhanced:
- Data in `outcome_map` table remains valid
- `pm_okr_link` is the more structured successor (explicit KR linking vs freetext objective)
- Future: `pm_okr_progress` can optionally read from both `okr_links` AND `outcome_map`

---

## Dependencies

**Existing (no changes):**
- `github.com/larksuite/oapi-sdk-go/v3` (v3.9.7, includes `service/okr/v1`)
- `h.Jira.*` interface (GetActiveSprints, GetSprintIssues, GetEpicIssues)
- `h.Memory.*` interface (GetSprintSnapshots, DB())
- `h.AI.Complete()` for narrative generation

**New:**
- Lark app permission: `okr:okr.content:readonly` (Phase 3 read)
- Lark app permission: `okr:okr.content:writeonly` (Phase 3 write)
- New `internal/lark/okr_client.go` for OKR API calls

---

## Relevant Files (for implementor)

| Purpose | File |
|---------|------|
| Existing outcome handlers | `application/tools/outcomes_handlers.go` |
| Existing outcome registration | `transport/outcomes.go` |
| Existing Lark webhook client | `internal/lark/webhook.go` |
| Existing Lark calendar client | `internal/calendar/client.go` |
| Lark SDK OKR module | `service/okr/v1/` in go module cache |
| Sprint snapshots (KPI source) | `internal/memory/sqlite_deep.go` |
| Jira epic/sprint methods | `internal/jirasdk/client.go` |
| Config struct | `config/config.go` |
| Bootstrap/DI | `internal/bootstrap/bootstrap.go` |
| Existing pm_forecast (Monte Carlo) | `application/tools/forecast_handlers.go` |

---

## Implementation Order

```
Phase 1 (standalone, highest immediate value):
  pm_kpi_dashboard -> pm_kpi_trend
  Delivers: engineering KPI visibility without any new integration

Phase 2 (core bridge, medium effort):
  pm_okr_link -> pm_okr_progress -> pm_okr_suggest
  Delivers: the translation layer between Jira work and OKR outcomes

Phase 3 (Lark integration, requires app permission setup):
  lark_okr_pull -> lark_okr_periods -> lark_okr_sync_progress
  Delivers: bi-directional Lark OKR connection

Phase 4 (intelligence layer, depends on Phase 1+2):
  pm_okr_health -> pm_okr_report -> pm_kpi_to_okr
  Delivers: proactive OKR coaching and leadership reporting
```

Each phase is independently shippable. Phase 1 has zero external dependencies.

---

## References

- Atlassian Community. "OKR-Jira Alignment Gap: Strategy Disconnected from Execution." 2026.
- OKRs Tool. "Why Sprint Teams Write Bad KRs." 2026.
- Mooncamp. "Output vs Outcome: Differences & OKR Examples." 2026.
- DORA/SPACE Framework. Engineering productivity metrics. 2025.
- Lark OKR API. Open Platform SDK `service/okr/v1`. github.com/larksuite/oapi-sdk-go.
- Larksuite CLI Skills. "lark-okr" agent skill documentation. o-mega.ai.
- Plandek. "DORA vs SPACE Metrics in AI-Enabled Engineering." 2026.
- KindaTechnical. "Measuring What Matters: Outcome Over Output." 2026.
