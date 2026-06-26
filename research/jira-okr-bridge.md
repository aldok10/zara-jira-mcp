# Jira → OKR/KPI Bridge: The Unsolved Problem

Research date: 2026-06-26

---

## The Core Problem

Jira tracks TASKS (tickets, sprints, story points).
OKRs track OUTCOMES (business impact, customer value).

**The gap:** There's no natural bridge between "we closed 45 tickets" and "customer retention improved 15%."

This is THE most common complaint from PM/SM using Jira:
> "OKRs live in a doc, a slide deck, or a spreadsheet. Jira lives in Jira. They never talk to each other."
> — Atlassian Community, 2026

> "The trickiest part comes when you need to regularly track and measure the progress of your OKRs, and how your day-to-day work impacts them."
> — Oboard, 2026

---

## Why It's Hard (Root Causes)

### 1. Level Mismatch
- OKR level: "Increase user activation by 20%"
- Jira level: "Implement onboarding tooltip for settings page"
- One KR may span 3 epics across 5 sprints. One ticket may serve 2 KRs.

### 2. Measurement Gap
- Jira measures: done/not done, story points, time
- KR measures: business metrics (revenue, NPS, churn, conversion)
- Business metrics live OUTSIDE Jira (analytics, CRM, finance)

### 3. Attribution Problem
- "Did closing those 45 tickets actually move the KR?"
- Multiple factors affect business metrics — isolating sprint contribution is near-impossible

### 4. Cadence Mismatch
- OKR cycle: quarterly (90 days)
- Sprint cycle: 2 weeks
- A KR won't show progress in a single sprint. Teams lose motivation.

### 5. Cultural Problem
- Engineering thinks in tickets
- Product/leadership thinks in outcomes
- Neither speaks the other's language fluently

---

## Existing Solutions (and why they fail)

| Solution | Problem |
|----------|---------|
| Jira epics linked to KRs manually | Nobody maintains the links. Stale in 2 weeks. |
| OKR Jira plugins (Oboard, etc.) | Progress = % tickets done. That's OUTPUT not OUTCOME. |
| Separate OKR tool (Lattice, 15Five) | Disconnected from daily work. Gets updated quarterly. |
| Confluence OKR pages | Static docs nobody reads after writing. |
| Lark OKR | Good UI but no AUTO connection to Jira data. Manual updates. |

---

## Lark OKR: What's Available

### Lark OKR Features:
- Set objectives + key results per period (quarterly)
- Align OKRs across org hierarchy
- Progress tracking (manual or metric-based)
- Comments and check-ins
- Dashboard with progress overview
- Integration with Lark chat (notifications)

### Lark OKR API (via oapi-sdk-go v3):
```
service/okr/v1/
├── ListPeriods      — Get OKR periods (quarters)
├── ListUserOkrs     — Get user's OKRs
├── CreateOkr        — Create new objective
├── UpdateOkr        — Update objective/KR
├── GetOkrDetail     — Full OKR with KRs and progress
├── CreateProgress   — Add progress record
└── ListProgress     — Get progress history
```

Already in our go.mod: `github.com/larksuite/oapi-sdk-go/v3`

---

## Our Solution: The Jira→OKR Bridge

### Design Principles:
1. **Don't force OKR structure on Jira** — Jira stays a task tracker
2. **Don't force ticket thinking on OKRs** — OKRs stay outcome-focused
3. **Bridge = AI-powered translation layer** — Interpret Jira data THROUGH OKR lens
4. **Auto-update** — No manual linking. System infers alignment.
5. **Lark OKR as source of truth** — We READ objectives from Lark, CALCULATE progress from Jira

### Architecture:
```
Lark OKR (objectives + KRs)
        ↓ read via API
zara-jira-mcp (bridge layer)
        ↓ pull sprint data
Jira (tickets, sprints, epics)
        ↓ AI interprets
"This sprint's work contributed ~30% to KR2"
        ↓ push progress
Lark OKR (auto-updated progress)
```

---

## Tools to Build

### Tier 1: Read + Display (no AI needed)
| Tool | What |
|------|------|
| `pm_okr_list` | List current period OKRs from Lark |
| `pm_okr_detail` | Show specific objective + KRs + current progress |
| `pm_okr_my` | My personal OKRs |

### Tier 2: Bridge (AI-powered)
| Tool | What |
|------|------|
| `pm_okr_sprint_alignment` | AI analyzes: which tickets this sprint align to which KRs? |
| `pm_okr_progress_calc` | Calculate KR progress from Jira data (tickets done that relate to KR) |
| `pm_okr_contribution` | "This sprint contributed X% toward objective Y" |
| `pm_okr_gap` | Find work NOT connected to any OKR (output without outcome) |

### Tier 3: Auto-Update (Lark write-back)
| Tool | What |
|------|------|
| `pm_okr_sync` | Push calculated progress back to Lark OKR |
| `pm_okr_checkin` | Auto-generate weekly OKR check-in from Jira activity |
| `pm_okr_report` | Quarterly OKR report: what moved, what didn't, why |

### Tier 4: KPI Dashboard
| Tool | What |
|------|------|
| `pm_kpi_define` | Define KPI with formula (e.g., "velocity / committed = predictability") |
| `pm_kpi_calculate` | Auto-calculate KPIs from Jira + memory data |
| `pm_kpi_trend` | KPI trend over time (sprints/months) |
| `pm_kpi_alert` | Alert when KPI drops below threshold |

---

## How AI Solves the Attribution Problem

The key insight: **AI can infer alignment even without explicit links.**

Given:
- Objective: "Improve onboarding completion rate"
- Sprint work: [PROJ-101 "Add tooltip to settings", PROJ-102 "Fix onboarding email", PROJ-103 "Refactor auth module"]

AI can:
1. Score relevance: PROJ-101 (high), PROJ-102 (high), PROJ-103 (low/none)
2. Estimate contribution: "2 of 3 done tickets directly serve this objective"
3. Generate narrative: "Sprint work moderately aligned to onboarding OKR. 67% of effort was outcome-aligned."

This is NOT precise measurement — it's **directional signal**. Better than nothing, worse than real analytics. But it's available TODAY from data we already have.

---

## Implementation Strategy

### Phase 1 (1 day): Lark OKR Read
- Add `internal/lark/okr.go` using existing SDK
- Tools: list periods, list OKRs, get detail
- Pure read-only. No AI. Just surface OKRs alongside Jira data.

### Phase 2 (1 day): AI Bridge
- `pm_okr_sprint_alignment`: Send sprint issues + OKR objectives to AI
- AI returns: alignment score per ticket per KR
- Store in memory for trend analysis

### Phase 3 (half day): Auto-update
- `pm_okr_sync`: Use Lark OKR CreateProgress API
- Push calculated progress after each sprint snapshot

### Phase 4 (half day): KPI Engine
- Define KPIs as formulas over existing data
- Auto-calculate from sprint snapshots + memory
- Alert via notification channels

---

## KPI Formulas We Can Calculate Today

From existing Jira + memory data, NO external analytics needed:

| KPI | Formula | Source |
|-----|---------|--------|
| Sprint Predictability | (done / committed) * 100 | sprint_snapshots |
| Velocity Trend | avg(velocity) over N sprints | sprint_snapshots |
| Blocker Resolution Time | avg(resolved_at - created_at) | blockers table |
| Risk Mitigation Rate | resolved_risks / total_risks | risks table |
| Action Completion Rate | completed_actions / total_actions | action_items |
| Team Happiness | avg(pulse_score) | team_pulse |
| Flow Efficiency | active_time / total_time | Jira transitions |
| Sprint Goal Hit Rate | achieved_goals / total_goals | sprint_goals |
| Bug Escape Rate | bugs_found_in_prod / total_items | Jira search |
| Cycle Time | avg(done_date - created_date) | Jira issues |
| WIP Compliance | avg_wip / wip_limit | flow_metrics |
| Meeting ROI | (value_score * attendees) / (duration * cost) | meeting_effectiveness |

---

## Summary

The Jira→OKR gap is real. Nobody has solved it well. Our advantage:
1. We already have persistent memory (sprint history, decisions, risks)
2. We already have AI integration (can infer alignment)
3. We already have Lark SDK (can read/write OKRs)
4. We already have the PM context (team metrics, health, velocity)

**The bridge isn't a new product — it's a natural extension of what we already built.**
