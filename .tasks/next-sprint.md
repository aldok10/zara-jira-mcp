# Next Sprint Tasks

Priority order. Each item independently shippable.

## P0: Context Engineering (unblocks everything)
- [ ] Shorten tool descriptions to <80 chars (saves ~30% tokens)
- [ ] Add `pm_load_module` for dynamic tool loading mid-session
- [ ] Tool usage tracking to SQLite (which tools called, frequency)
- [ ] Auto-profile suggestion based on usage patterns

## P1: Communication Templates
- [ ] `pm_communicate` — Minto-structured message for any audience
- [ ] `pm_feedback_prep` — SBI feedback draft from team data
- [ ] `pm_escalation_draft` — Pyramid-structured escalation
- [ ] `pm_decision_record` — ADR format (context/options/consequences)

## P2: Psychological Safety (Project Aristotle)
- [ ] `pm_safety_survey` — 7-question survey (1-5 scale per member)
- [ ] `pm_safety_trend` — Safety score over time
- [ ] `pm_team_aristotle` — Full 5-pillar assessment

## P3: SPACE Metrics + DevEx
- [ ] `pm_space_metrics` — Aggregate S/P/A/C/E from Jira+GitHub
- [ ] `pm_flow_disruption` — Detect broken flow signals
- [ ] `pm_maker_time` — Calendar analysis for deep work blocks
- [ ] `pm_right_size` — Find oversized tickets hurting predictability

## P4: Hypothesis-Driven Development
- [ ] `pm_hypothesis` — Record belief + expected outcome + measure
- [ ] `pm_hypothesis_review` — Post-sprint validation check
- [ ] `pm_estimation_accuracy` — Estimates vs actuals feedback loop

## P5: Enhanced Facilitation
- [ ] Expand `pm_facilitate` with 10+ Liberating Structures
- [ ] `pm_retro_format` — AI-suggest format from team context
- [ ] `pm_meeting_audit` — "Could this be async?"

## P6: Smart Communication Routing
- [ ] `pm_audience_router` — Same data, 3 versions (exec/PO/team)
- [ ] `pm_silence_detector` — Flag disengaged stakeholders
- [ ] `pm_raci` — Auto-generate from Jira assignments
- [ ] `pm_communication_plan` — Who/what/when/how

## P7: EBM Dashboard
- [ ] `pm_ebm_dashboard` — 4 KVAs tracked over time
- [ ] `pm_value_check` — "Is this output or outcome?"

## P8: Team Topology
- [ ] `pm_cognitive_load` — Assess per-team mental overhead
- [ ] `pm_team_dependencies_map` — Cross-team interaction modes

---

## Start Here Tomorrow
1. Pick P0 (30 min) — highest leverage, makes everything else lighter
2. Then P1 (1 hr) — most visible PM value
3. Then P2 (30 min) — differentiator no competitor has

---

## P-OKR: Jira→OKR/KPI Bridge (HIGH PRIORITY)

Based on `research/jira-okr-bridge.md`. The #1 unsolved PM pain point.

### OKR-1: Lark OKR Read (Phase 1, 1 day)
- [ ] `internal/lark/okr.go` — Lark OKR API client (ListPeriods, ListUserOkrs, GetOkrDetail)
- [ ] `pm_okr_list` — Show current period objectives
- [ ] `pm_okr_detail` — Full objective + KRs + progress
- [ ] `pm_okr_my` — My personal OKRs

### OKR-2: AI Bridge (Phase 2, 1 day)
- [ ] `pm_okr_sprint_alignment` — AI scores: which tickets align to which KRs
- [ ] `pm_okr_contribution` — "This sprint contributed X% to objective Y"
- [ ] `pm_okr_gap` — Find work not connected to any OKR

### OKR-3: Auto-Update (Phase 3, half day)
- [ ] `pm_okr_sync` — Push progress to Lark OKR via CreateProgress API
- [ ] `pm_okr_checkin` — Auto-generate weekly OKR check-in

### OKR-4: KPI Engine (Phase 4, half day)
- [ ] `pm_kpi_define` — Define KPI with formula
- [ ] `pm_kpi_calculate` — Auto-calculate from existing data
- [ ] `pm_kpi_trend` — KPI over time
- [ ] `pm_kpi_alert` — Alert on threshold breach

### Pre-built KPIs (zero config needed):
- Sprint Predictability = done/committed * 100
- Blocker Resolution Time = avg(resolved - created)
- Action Completion Rate = completed/total actions
- Sprint Goal Hit Rate = achieved/total goals
- Team Happiness = avg(pulse score)
- Cycle Time = avg(done - created) from Jira
