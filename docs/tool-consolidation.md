# Tool Consolidation Analysis

> 248 tools total. ~30 are overlapping/redundant. This doc maps canonical tools vs aliases.

## Recommended Consolidation (Don't Remove — Route)

### Principle: Keep all tool names (backward compat), but route duplicates to the canonical handler.

---

## 1. Forecasting (5 → 2 canonical)

| Canonical | Aliases (route to canonical) | Notes |
|-----------|------------------------------|-------|
| `pm_forecast` | `pm_forecast_sprint`, `pm_sprint_forecast_simple` | Monte Carlo, accepts board_id + items_remaining |
| `kpi_predictability` | `pm_predictability` | Committed vs delivered ratio |

**Action**: Make `pm_forecast_sprint` and `pm_sprint_forecast_simple` call the same handler as `pm_forecast`.

---

## 2. Daily/Weekly Digest (3 → 2 canonical)

| Canonical | Aliases | Notes |
|-----------|---------|-------|
| `pm_daily_digest` | `daily_digest` | Same tool, daily cadence |
| `pm_weekly_digest` | (none) | Keep, different cadence |

**Action**: Route `daily_digest` to `pm_daily_digest` handler.

---

## 3. Blockers & Escalation (8 → 4 canonical)

| Canonical | Aliases | Notes |
|-----------|---------|-------|
| `pm_blockers` | (none) | Show active/history |
| `pm_impediment_aging` | `pm_blocker_aging` | Same concept, keep both names |
| `report_escalation_brief` | `pm_escalation_report`, `pm_escalation_draft`, `pm_escalate` (partial) | SCQA-formatted escalation |
| `pm_escalations` | (none) | History log |

**Action**: Route `pm_blocker_aging` → `pm_impediment_aging`. Route `pm_escalation_report` + `pm_escalation_draft` → `report_escalation_brief`.

---

## 4. Executive Reporting (5 → 3 canonical)

| Canonical | Aliases | Notes |
|-----------|---------|-------|
| `pm_exec_report` | `pm_report` | Executive-level, no jargon |
| `report_to_po` | (none) | PO-specific (value + decisions) |
| `report_delivery_confidence` | `pm_status_draft` | GREEN/AMBER/RED |
| `report_cross_team_deps` | `pm_dependency_report` | Same thing |

**Action**: Route `pm_report` → `pm_exec_report`. Route `pm_dependency_report` → `report_cross_team_deps`.

---

## 5. Health Checks (7 → 3 canonical)

| Canonical | Aliases/Aspect | Notes |
|-----------|---------------|-------|
| `pm_sprint_health` | (none) | Composite score 0-100 |
| `pm_team_health` | `pm_comms_health` (aspect=comms), `pm_qa_health` (aspect=quality) | Could add aspect param |
| `pm_backlog_health` | (none) | Keep separate (different data source) |

**Action**: Consider merging qa/comms health into team_health with `aspect` parameter in future. Low priority.

---

## 6. OKR/KPI (keep all — newly created, no overlap)

All 7 `okr_*` and `kpi_*` tools are distinct:
- `okr_define_signal` — CRUD
- `okr_calculate` — compute
- `okr_status` — report
- `okr_alignment` — check
- `okr_delete_signal` — CRUD
- `kpi_dora` — DORA metrics
- `kpi_predictability` — predictability index

---

## Summary of Redundancy

| Category | Current | After Consolidation | Savings |
|----------|---------|--------------------:|---------|
| Forecast | 5 | 2 canonical (+3 aliases) | 3 handler merges |
| Digest | 3 | 2 canonical (+1 alias) | 1 handler merge |
| Blockers/Escalation | 8 | 4 canonical (+4 aliases) | 4 handler merges |
| Reporting | 5 | 3 canonical (+2 aliases) | 2 handler merges |
| Health | 7 | 3 canonical (keep rest) | 0 (low priority) |
| **Total** | **28** | **14 canonical** | **10 handler merges** |

## Recommendation

**Don't consolidate now.** The tool names are the API surface. Removing them breaks existing prompts/configs. Instead:

1. **Document canonical tools** in SKILL.md/README
2. **Add "See also" in descriptions** (e.g., `pm_forecast_sprint` description: "Alias for pm_forecast. Use pm_forecast instead.")
3. **In future v1.0**: deprecate aliases, route all to canonical handlers
4. **Focus energy on**: making the 14 canonical tools EXCELLENT rather than reducing count

## Tools That Should Be IMPROVED (not merged)

| Tool | Current Gap | Improvement |
|------|-------------|-------------|
| `pm_forecast` | Only sprint-level | Add epic_key param for epic-level forecast |
| `okr_calculate` | No Lark push | Add auto-push to Lark OKR when lark_kr_id set |
| `kpi_dora` | Jira-only | Needs GitHub/GitLab commit data for true lead time |
| `pm_sprint_health` | Static score | Add trend (improving/declining vs last sprint) |
| `report_escalation_brief` | No auto-send | Add send_to_lark/slack param |
| `pm_what_next` | No OKR context | Include OKR alignment in priority scoring |
