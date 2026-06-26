# Tool Consolidation Plan

Status: Ready for execution
Date: 2026-06-26
Author: Zara (multi-agent audit)

---

## Executive Summary

Audit found **30 tool registrations** that should be eliminated through deduplication and merging.
- 4 CRITICAL: same tool name registered twice (runtime collision risk)
- 2 CRITICAL: module double-registration bugs
- 7 HIGH: functional duplicates confusing users
- 6 MEDIUM: overlapping but mergeable

---

## CRITICAL FIXES (must do first — potential runtime bugs)

### 1. Exact Name Duplicates (4 tools)

| Tool | Keep | Remove From |
|------|------|-------------|
| `jira_trace_branch` | `transport/trace.go` | `transport/leverage.go` |
| `pm_tech_debt_ratio` | `transport/leverage.go` | `transport/versions.go` |
| `pm_priority_churn` | `transport/leverage.go` | `transport/versions.go` |
| `pm_incident_summary` | `transport/pagerduty.go` (rename to `pm_pagerduty_incident`) | `transport/leverage.go` (rename to `pm_incident_analysis`) |

### 2. Double Module Registration (2 bugs in server.go)

| Function | Registered In | Fix |
|----------|--------------|-----|
| `registerCommsGapTools` | `pm` module AND `stakeholder` module | Remove from `pm` module (belongs in `stakeholder`) |
| `registerCommunicationTools` | `stakeholder` module AND `shortcuts` module | Remove from `shortcuts` module (belongs in `stakeholder`) |

**Impact:** Fixes 13 tools being registered twice under "all" profile.

---

## HIGH PRIORITY MERGES (7 tools deprecated)

### 3. Forecasting: `pm_forecast_sprint` → deprecated

- **Keep:** `pm_forecast` (full Monte Carlo, takes `remaining_items`)
- **Keep:** `pm_sprint_forecast_simple` (different algorithm: burn rate)
- **Deprecate:** `pm_forecast_sprint` (subset of `pm_forecast`)
- **Action:** Remove from `transport/server.go` registerForecastTools

### 4. What Next: `pm_what_next` → deprecated

- **Keep:** `pm_next` (shorter name, same utility, in shortcuts)
- **Deprecate:** `pm_what_next`
- **Action:** Remove from `transport/whatnext.go`, redirect to `pm_next`

### 5. Message Composition: `pm_communicate` → deprecated

- **Keep:** `pm_compose` (canonical "adapt message for audience")
- **Deprecate:** `pm_communicate` (same purpose, in safety.go)
- **Action:** Remove from `transport/safety.go`

### 6. Feedback: merge into `pm_feedback`

- **Keep:** `pm_feedback_prep` approach (data-backed, takes member + observation)
- **Enhance:** Add Radical Candor framing from `pm_feedback_coach`
- **Deprecate:** `pm_feedback_coach`
- **New name:** `pm_feedback` (rename `pm_feedback_prep`)
- **Action:** Merge handler logic, update `transport/safety.go` + remove from `transport/communication.go`

### 7. Blocker Aging: `pm_impediment_aging` → deprecated

- **Keep:** `pm_blocker_aging` (better name, has SLA tracking + ownership)
- **Enhance:** Add avg resolution stats from `pm_impediment_aging`
- **Deprecate:** `pm_impediment_aging`
- **Action:** Merge logic into `management_handlers.go`, remove from `transport/outcomes.go`

### 8. Daily Digest: merge into one tool

- **Keep:** `pm_daily_digest` (generates content)
- **Enhance:** Add optional `send` param to push to notification channel (from `daily_digest`)
- **Deprecate:** `daily_digest` (in routing.go)
- **Action:** Add send capability to coaching_handlers.go, remove from routing.go

### 9. Dependency Report: `report_cross_team_deps` → deprecated

- **Keep:** `pm_dependency_report` (in management.go, has aging + overdue)
- **Deprecate:** `report_cross_team_deps` (same output, different name)
- **Action:** Remove from `transport/reporting.go`

---

## MEDIUM PRIORITY MERGES (6 tools consolidated)

### 10. Outcome Map: `pm_outcome_map` → deprecated

- **Keep:** `pm_okr_tree` (V2, full hierarchy: objective → KR → initiative)
- **Keep:** `pm_okr_define` (V1, simpler epic-linked OKRs)
- **Deprecate:** `pm_outcome_map` (subset of `pm_okr_define`)
- **Migration:** Existing `outcome_map` table data can be read by `pm_outcome_history` (keep as read-only)

### 11. Hypothesis: `pm_hypothesis` → merged into `pm_experiment`

- **Keep:** `pm_experiment` (richer model: action + duration + sprint context)
- **Enhance:** Add `deadline` field from `pm_hypothesis` + validate capability
- **Deprecate:** `pm_hypothesis` and `pm_hypothesis_validate`
- **Action:** Merge into smart_handlers.go RecordExperiment

### 12. Decision Recording: merge `pm_record_decision` + `pm_decision_record`

- **Keep:** `pm_record_decision` (canonical, in memory tools)
- **Enhance:** Add optional ADR fields (alternatives, consequences) from `pm_decision_record`
- **Deprecate:** `pm_decision_record` (in safety.go)
- **Keep:** `pm_decide` (quick shortcut, delegates to above)
- **Keep:** `pm_announce_decision` (different purpose: communication)

### 13. Commitment: `pm_commitment_report` → deprecated

- **Keep:** `pm_commitment_check` (in care.go)
- **Enhance:** Add `audience` param ("team" vs "management" format)
- **Deprecate:** `pm_commitment_report`

### 14. Management Brief: `pm_management_brief` → deprecated

- **Keep:** `pm_status_draft` (auto-data-enriched, takes `audience` param)
- **Deprecate:** `pm_management_brief` (subset of `pm_status_draft(audience:"management")`)
- **Action:** Remove from `transport/management.go`

### 15. Escalation Writing: merge into one tool

- **Keep:** `pm_escalation_draft` (rename to `pm_escalation_write`)
- **Enhance:** Add `format` param: "pyramid" (default) | "scqa"
- **Deprecate:** `pm_escalate_message`
- **Keep (separate purpose):** `pm_escalate` (auto-action), `pm_escalation_report` (data view)

---

## ENHANCEMENT RECOMMENDATIONS (from research)

Based on `research/` findings, these improvements to EXISTING tools:

### Communication Tools (from communication-research-extended.md + communication-ai-era.md)

| Tool | Enhancement | Source |
|------|-------------|--------|
| `pm_compose` | Add SCARF awareness (flag potential threat triggers) | SCARF Model research |
| `pm_status_draft` | Add confidence signaling (how sure are we?) | communication-implementation-tasks.md |
| `pm_feedback` (merged) | Add SCARF check before delivery | communication-research-extended.md |
| `pm_escalate` | Add acknowledgment tracking (was it received?) | P2 infrastructure research |

### OKR/KPI Tools (from okr-kpi-integration.md + jira-okr-bridge.md)

| Tool | Enhancement | Source |
|------|-------------|--------|
| `pm_okr_sync` (V2) | Add Lark OKR write-back (CreateProgressRecord) | jira-okr-bridge.md Tier 3 |
| `pm_kpi_dashboard` | Add DORA metrics (deploy freq, lead time, MTTR, change failure) | pm-ai-era-research.md |
| `pm_okr_tree` | Add AI inference: auto-suggest ticket → KR alignment | jira-okr-bridge.md innovation |
| New: `pm_kpi_trend` | Individual KPI over time with sparkline | okr-kpi-integration.md |
| New: `pm_okr_health` | Time-elapsed vs progress risk detection | gap-analysis-addendum.md |

### Team Care Tools (from team-care-research.md)

| Tool | Enhancement | Source |
|------|-------------|--------|
| `pm_overload_check` | Add per-person overload risk score formula | team-care-research.md |
| `pm_sprint_health` | Add sustainable pace signals | team-care-research.md |
| `pm_team_pulse` | Add burnout early warning (Maslach signals) | team-care-research.md |

### SM Leverage Tools (from pm-leverage-research.md)

| Tool | Enhancement | Source |
|------|-------------|--------|
| `pm_daily_digest` | Add communication nudges section | communication-implementation-tasks.md |
| `pm_standup_prep` | Add "day 3 forecast" (early sprint risk) | pm-leverage-research.md |
| `pm_anti_patterns` | Add Lencioni dysfunction level mapping | communication-research-extended.md |

---

## NEW TOOLS TO BUILD (genuinely missing, not duplicates)

From research + gap analysis, these are confirmed NOT existing:

| Tool | Purpose | Priority | Source |
|------|---------|----------|--------|
| `pm_comms_health` | Communication anti-pattern scanner | P0 | communication-implementation-tasks.md |
| `pm_cadence_check` | Is PM meeting communication commitments? | P0 | gap-analysis-addendum.md |
| `pm_kpi_trend` | Single KPI trend over time | P0 | okr-kpi-integration.md |
| `pm_okr_health` | Time vs progress risk analysis | P0 | gap-analysis-addendum.md |
| `pm_comms_nudge` | Proactive communication suggestions | P1 | communication-implementation-tasks.md |
| `pm_conversation_prep` | Framework-based conversation preparation | P1 | communication-enhancement-spec.md |
| `pm_feedback_log/due/close` | Feedback lifecycle tracking | P1 | communication-enhancement-spec.md |
| `pm_kpi_to_okr` | AI: suggest OKR KRs from current metrics | P1 | okr-kpi-bridge-spec.md |
| `pm_okr_suggest` | AI: which tickets serve which OKRs? | P1 | jira-okr-bridge.md |
| `lark_okr_pull` | MCP wrapper for existing Lark OKR client reads | P2 | gap-analysis-addendum.md |
| `lark_okr_sync_progress` | Write progress to Lark OKR | P2 | jira-okr-bridge.md |
| `pm_notification_budget` | Track notification volume/fatigue | P2 | pm-leverage-research.md |
| `pm_cognitive_load` | Detect context-switching overload | P2 | pm-ai-era-research.md |

---

## Implementation Order

```
Week 1: CRITICAL fixes (runtime safety)
  - Remove 4 exact-name duplicates
  - Fix 2 module double-registrations
  - Verify: go build && go test ./...

Week 2: HIGH merges (user confusion)
  - Deprecate 7 tools (mark with "[DEPRECATED]" in description first)
  - Merge handler logic where needed
  - Update SKILL.md tool reference

Week 3: MEDIUM merges + enhancements
  - 6 consolidations
  - Add enhancement params to existing tools
  - Update docs/reporting-guide.md references

Week 4+: New tools from specs
  - Build genuinely new tools (pm_comms_health, pm_kpi_trend, etc.)
  - Wire Lark OKR write operations
```

---

## Deprecation Strategy

Do NOT delete tools immediately. Use this 3-step process:

1. **Mark:** Change description to start with `[DEPRECATED: use pm_X instead]`
2. **Redirect:** Handler returns `textResult("This tool is deprecated. Use pm_X instead.")`
3. **Remove:** After 2 releases, delete registration + handler

This prevents breaking existing users/prompts that reference the old tool name.

---

## Files to Modify

| File | Changes |
|------|---------|
| `transport/leverage.go` | Remove `jira_trace_branch`, `pm_incident_summary` (or rename) |
| `transport/versions.go` | Remove `pm_tech_debt_ratio`, `pm_priority_churn` |
| `transport/server.go` | Fix double module registration, remove `pm_forecast_sprint` |
| `transport/safety.go` | Remove `pm_communicate`, `pm_decision_record` |
| `transport/communication.go` | Remove `pm_feedback_coach`, `pm_escalate_message` |
| `transport/whatnext.go` | Remove `pm_what_next` |
| `transport/routing.go` | Remove `daily_digest` |
| `transport/outcomes.go` | Remove `pm_impediment_aging`, mark `pm_outcome_map` deprecated |
| `transport/management.go` | Remove `pm_management_brief`, `pm_commitment_report` |
| `transport/reporting.go` | Remove `report_cross_team_deps`, `report_escalation_brief` |
| `transport/okr_kpi.go` | Remove `pm_hypothesis`, `pm_hypothesis_validate` |
| `SKILL.md` | Update tool reference (remove deprecated, update descriptions) |
| `docs/reporting-guide.md` | Update tool references |
