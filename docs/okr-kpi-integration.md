# OKR & KPI: Menjembatani Gap Antara Jira dan Business Objectives

## Masalah Utama

### "Ticket selesai, tapi OKR nggak bergerak"

Ini masalah klasik yang dialami hampir semua PM:

1. **Jira tracks output** (tickets done, story points burned, velocity). **OKR tracks outcome** (user behavior change, revenue impact, quality improvement).
2. **Tidak ada mapping otomatis** antara "PROJ-123 Done" dan "Key Result: Reduce churn by 20% progress +5%".
3. **PM jadi translator manual** — tiap minggu copy data dari Jira board, lalu interpretasikan ke format OKR. Atlassian Community (2025): "I sat in a sprint planning session where nobody could answer: does this ticket move any of our key results?"

Research (OKR Tool, 2026): **Teams yang connect goals ke outcomes (bukan outputs) 30% lebih likely achieve goals.** Tapi kebanyakan agile teams measure activity, bukan impact.

### Kenapa Ini Terlalu Sulit

| Problem | Detail |
|---------|--------|
| **Abstraction mismatch** | Jira: task-level (implement API endpoint). OKR: outcome-level (reduce API latency to <200ms). Beda level abstraksi. |
| **Temporal mismatch** | Sprint: 2 minggu. OKR: quarterly. Satu KR bisa span 6 sprint. Progress per sprint = fractional, hard to quantify. |
| **Attribution problem** | 10 tickets dikerjakan sprint ini. Mana yang contribute ke KR mana? Kadang 1 ticket contribute ke multiple KR. |
| **Metric lag** | Ticket done hari ini, tapi impact ke KR (e.g., user retention) baru terlihat 30 hari kemudian. |
| **Manual update fatigue** | "When each bi-weekly check-in requires hours of progress tracking and data input, the benefits of OKRs are simply lost." (Atlassian, 2025) |

### Spesifik: Lark OKR Users

Pengguna Lark OKR punya challenge tambahan:
- OKR di Lark, execution tracking di Jira — dua system yang nggak natively connected
- Lark OKR punya cycle management, tapi nggak bisa auto-pull progress dari Jira sprint data
- PM harus login ke 2 platform, cross-reference, manual update progress

---

## Solusi: Automated OKR Intelligence dari Sprint Data

### Arsitektur

```
Jira Sprint Data ──→ zara-jira-mcp ──→ OKR Progress Calculation
     (output)           (bridge)            (outcome)
                           │
                           ├──→ Lark OKR (update via API)
                           ├──→ Internal OKR tracking (SQLite)
                           └──→ Reports (exec, PO, team)
```

### Pendekatan: Outcome Mapping

Bukan tracking "berapa ticket selesai" tapi "berapa KR yang ter-advance oleh sprint work."

**Step 1:** PM define mapping OKR ↔ Jira work sekali di awal quarter:
```
pm_okr_map(
  objective: "Improve platform reliability",
  key_result: "Reduce P1 incidents from 5/month to 1/month",
  jira_signal: "jql:project=INFRA AND type=Bug AND priority=P1 AND resolved>=-30d",
  metric_type: "count_decrease",
  baseline: 5,
  target: 1
)
```

**Step 2:** System otomatis calculate KR progress dari Jira data:
```
pm_okr_progress(board_id:X)
→ "KR: Reduce P1 incidents 5→1"
→ "Current: 2 incidents in last 30 days (60% progress toward target)"
→ "Contributing sprint work: INFRA-234, INFRA-256, INFRA-289 (stability fixes)"
```

**Step 3:** Auto-update Lark OKR (kalau Lark API configured):
```
pm_okr_sync_lark
→ Push progress ke Lark OKR system
→ PM nggak perlu login ke Lark buat update
```

---

## Tool Design

### OKR Management Tools (Planned)

| Tool | Purpose | Params |
|------|---------|--------|
| `pm_okr_map` | Link KR to Jira signal (JQL, label, epic, or custom metric) | `objective`, `key_result`, `jira_signal`, `metric_type`, `baseline`, `target` |
| `pm_okr_progress` | Calculate all KR progress from live Jira data | `board_id`, `period` (current quarter default) |
| `pm_okr_report` | Generate OKR status report (per audience) | `audience` (team/po/exec), `format` |
| `pm_okr_sync_lark` | Push progress to Lark OKR via API | (auto from mapped OKRs) |
| `pm_okr_suggest` | AI suggest which sprint items contribute to which KR | `board_id` |
| `pm_okr_gaps` | Identify KRs with no linked sprint work (strategy-execution gap) | `board_id` |
| `pm_okr_health` | Overall OKR health: on-track, at-risk, off-track per KR | (auto) |
| `pm_kpi_track` | Track arbitrary KPI with time-series (manual or auto input) | `name`, `value`, `date` |
| `pm_kpi_dashboard` | Show all KPIs with trend | `period` |
| `pm_kpi_alert` | Alert when KPI crosses threshold | `name`, `threshold`, `direction` |

### Metric Types (Signal → Progress Calculation)

| Type | How It Works | Example |
|------|-------------|---------|
| `count_increase` | Count matching issues, more = better | "Ship 5 features" → count done stories with label "feature" |
| `count_decrease` | Count matching issues, fewer = better | "Reduce P1 bugs to 1/month" → count resolved P1 |
| `completion_rate` | % of linked epic/sprint items done | "Complete platform migration" → epic completion % |
| `velocity_target` | Average velocity over period | "Sustain 25pts/sprint velocity" → snapshot history |
| `time_metric` | Cycle time, lead time, MTTR | "Reduce deploy lead time to <1 day" → flow metrics |
| `custom_jql` | Raw JQL count or percentage | Any custom signal |
| `manual` | PM manually updates (for non-Jira metrics) | "NPS score", "Customer interviews conducted" |

### Lark OKR Integration

Lark Open Platform menyediakan OKR API (via `lark-openapi-mcp`):
- List OKR periods/cycles
- Get objectives and key results per user/team
- Update KR progress programmatically

**Flow:**
1. `pm_okr_map` — PM creates mapping once
2. Sprint end → `pm_okr_progress` auto-calculates
3. `pm_okr_sync_lark` — pushes calculated progress to Lark OKR
4. Lark OKR reflects reality tanpa PM manual update

---

## Contoh Real-World

### Engineering Team OKR

**Objective:** Improve platform reliability

| Key Result | Jira Signal | Metric | Baseline | Target | Current |
|-----------|-------------|--------|----------|--------|---------|
| Reduce P1 incidents to 1/month | `project=INFRA AND type=Bug AND priority=P1 AND resolved>=-30d` | count_decrease | 5 | 1 | 2 |
| Achieve 99.9% uptime | manual (from monitoring) | manual | 99.2% | 99.9% | 99.7% |
| Deploy time < 30 minutes | `pm_flow_metrics` cycle time for deploy-tagged items | time_metric | 120min | 30min | 45min |

**Output dari `pm_okr_progress`:**
```
Q3 OKR Progress: "Improve platform reliability"

KR1: Reduce P1 incidents (5→1)
  Status: ON TRACK (60% progress)
  Current: 2 incidents in last 30d
  Contributing work: INFRA-234, INFRA-256 (resolved this sprint)
  
KR2: 99.9% uptime
  Status: AT RISK (manual metric, last update 5 days ago)
  Current: 99.7%
  Gap: 0.2% remaining

KR3: Deploy time < 30min
  Status: IMPROVING (62% progress)
  Current: 45min avg (from 120min baseline)
  Contributing: DEVOPS-89, DEVOPS-91 (CI pipeline optimization)
```

### Product Team OKR

**Objective:** Increase user activation

| Key Result | Jira Signal | Metric |
|-----------|-------------|--------|
| Ship onboarding redesign by end of Q3 | Epic PROD-50 completion % | completion_rate |
| Reduce time-to-first-value from 15min to 5min | manual (analytics) | manual |
| 3 user interviews per week | label "user-research" issues created/week | count_increase |

---

## Output vs Outcome: Bridging the Gap

### The Translation Table

PM perlu translate sprint output ke business outcome:

| Sprint Output (Jira) | Business Outcome (OKR) | Bridge |
|----------------------|------------------------|--------|
| 12 stories done | → contributes to KR: "Ship onboarding v2" | epic completion % |
| 3 P1 bugs fixed | → contributes to KR: "Reduce incidents" | count decrease |
| API latency PR merged | → contributes to KR: "< 200ms response time" | metric will lag 1-2 weeks |
| 5 user research tickets | → contributes to KR: "3 interviews/week" | direct count |

### `pm_okr_suggest` — AI-Powered Attribution

"Sprint ini ada 15 items done. Mana yang contribute ke KR mana?"

AI reads:
- Sprint items (title, description, labels, epic)
- Mapped OKRs and their signals
- Produces attribution suggestion

This eliminates the manual "Friday afternoon slide-making" problem.

---

## KPI Tracking (Simpler Than OKR)

Buat yang nggak pakai full OKR framework tapi butuh track metrics:

```
pm_kpi_track(name:"Sprint Velocity", value:23)
pm_kpi_track(name:"Bug Escape Rate", value:2)
pm_kpi_track(name:"Customer NPS", value:42)
pm_kpi_dashboard
```

Output:
```
KPI Dashboard (last 5 data points):

Sprint Velocity:     18 → 20 → 22 → 23 → 23  (stable, +28% from start)
Bug Escape Rate:      5 →  4 →  3 →  2 →  2  (improving, -60%)
Customer NPS:        35 → 38 → 40 → 42 → 42  (improving, +20%)
Team Satisfaction:    3 →  3 →  4 →  4 → N/A  (due for update)
```

---

## Implementation Priority

| Phase | What | Effort |
|-------|------|--------|
| **Phase 1** | `pm_kpi_track` + `pm_kpi_dashboard` (simple time-series) | 4h |
| **Phase 2** | `pm_okr_map` + `pm_okr_progress` (Jira signal → KR progress) | 8h |
| **Phase 3** | `pm_okr_report` + `pm_okr_gaps` + `pm_okr_health` | 6h |
| **Phase 4** | `pm_okr_suggest` (AI attribution) | 4h |
| **Phase 5** | `pm_okr_sync_lark` (Lark OKR API push) | 6h |
| **Phase 6** | `pm_kpi_alert` (threshold alerts) | 3h |

Total: ~31 hours for complete OKR/KPI intelligence.

---

## References

- Atlassian Community (2025). "OKR-Jira Alignment Gap: Strategy Disconnected from Execution"
- OKR Tool (2026). "Teams connecting goals to outcomes are 30% more likely to hit them"
- Mooncamp (2026). "Output vs Outcome: differences and OKR examples"
- ResearchGate (2022). "How agile teams make OKRs work" — enabling/limiting situations
- Lark Open Platform. OKR API documentation
- larksuite/lark-openapi-mcp. Official Lark MCP server with OKR capabilities
