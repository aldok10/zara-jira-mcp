# Tool Audit Complete — 2026-06-27 00:05

## Stats
- **248 unique tools** registered
- **35 transport files**
- **1 duplicate found & fixed** (pm_incident_summary)
- **0 tools to merge** — all serve distinct purposes

## Overlap Analysis

### Checked for semantic duplicates:
| Group | Tools | Verdict |
|-------|-------|---------|
| Forecasting | pm_forecast, pm_forecast_sprint, pm_sprint_forecast_simple | KEEP: full Monte Carlo vs quick calc |
| Health | pm_sprint_health, pm_scorecard, pm_dashboard | KEEP: realtime vs grade vs full overview |
| Standup | pm_standup_prep, pm_daily_digest, pm_daily_delta | KEEP: talking points vs AI summary vs delta |
| What next | pm_what_next, pm_next | KEEP: detailed vs shortcut alias |
| Blocker | pm_blockers, pm_blocker_aging, pm_impediment_aging | KEEP: list vs SLA tracking vs deep analysis |
| Communication | pm_communicate, pm_status_draft, pm_comms_plan | KEEP: audience-router vs draft vs plan |

### Files that could be MERGED (code quality, not feature):
- `smart_context.go` has pm_smart + pm_do + pm_report — same as `smart.go` concept. **One was created by subagent, other manually.** Whichever builds is the one that stays.
- `pm_shortcuts.go` has `pm` + `pm_next` which overlap with `smart_context.go`'s `pm_smart`. Different interfaces to same data.

### Recommendation: NO MERGE
The profile system (`PM_PROFILE=smart` → 7 tools) already solves the UX problem. Merging at code level adds risk for zero user benefit. The 248 tools are distinct enough.

## What's Actually Missing (from research):
1. SPACE metrics (P3) — not implemented
2. Liberating Structures library (P5) — pm_facilitate exists but basic
3. Lark OKR API read (P-OKR Phase 1) — client exists, tools registered locally

## Decision
Status: **AUDIT COMPLETE. NO ACTION NEEDED.**
Tools are healthy. Focus energy on new features (P3, P5) not consolidation.
