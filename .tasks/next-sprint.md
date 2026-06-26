# Next Sprint Tasks

Priority order. Each item independently shippable.

## P0: Context Engineering (DONE)
- [x] Shorten tool descriptions to <80 chars (saves ~30% tokens)
- [x] Add tool usage tracking to SQLite (`pm_tool_usage`)
- [ ] Add `pm_load_module` for dynamic tool loading mid-session
- [ ] Auto-profile suggestion based on usage patterns

## P1: Communication Templates (DONE)
- [x] `pm_communicate` — Minto-structured message for any audience
- [x] `pm_feedback_prep` — SBI feedback draft from team data
- [x] `pm_escalation_draft` — Pyramid-structured escalation
- [x] `pm_decision_record` — ADR format (context/options/consequences)

## P2: Psychological Safety (Partial)
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

## P5: Enhanced Facilitation (DONE)
- [x] Expand `pm_facilitate` with Liberating Structures (1-2-4-All, TRIZ, 15% Solutions)
- [ ] `pm_retro_format` — AI-suggest format from team context
- [ ] `pm_meeting_audit` — "Could this be async?"

## P6: Smart Communication Routing (Partial)
- [x] `pm_communicate` — audience-aware routing
- [x] `pm_silence_detector` — Flag disengaged stakeholders
- [ ] `pm_raci` — Auto-generate from Jira assignments
- [x] `pm_comms_anti_patterns` — Detect communication dysfunctions

## P7: EBM Dashboard
- [ ] `pm_ebm_dashboard` — 4 KVAs tracked over time
- [ ] `pm_value_check` — "Is this output or outcome?"

## P8: Team Topology
- [ ] `pm_cognitive_load` — Assess per-team mental overhead
- [ ] `pm_team_dependencies_map` — Cross-team interaction modes

---

## Completed This Session (v0.5.0)

1. Tool descriptions shortened to <80 chars (all 262 tools)
2. Tool usage tracking (`pm_tool_usage`)
3. Team sentiment analysis (`pm_sentiment`)
4. Notification budget/throttle (5/day, auto-enforced)
5. Forecast calibration (`pm_calibration`)
6. Meeting ROI (`pm_meeting_roi`)
7. Liberating Structures in `pm_facilitate`
8. Security: input validation, path restriction on `jira_raw_request`
9. Removed deprecated tools (`pm_status_draft`, `pm_forecast_sprint`)
10. Fixed duplicate tool registration (`registerCommunicationTools`)
11. Updated README, AGENTS.md, version to v0.5.0

---

## P-OKR: Jira→OKR/KPI Bridge (DONE)

Based on `research/jira-okr-bridge.md`.

### OKR-1: Local OKR System (DONE)
- [x] `pm_okr_define` — Create OKR with key results
- [x] `pm_okr_list` — Show current objectives
- [x] `pm_kr_link` — Link Jira issues to KRs
- [x] `pm_kr_progress` — Auto-calculate from Jira

### OKR-2: AI Bridge (DONE)
- [x] `pm_outcome_review` — AI: did sprint work move OKR?
- [x] `pm_okr_health` — At-risk objectives
- [x] `pm_goal_hit_rate` — Sprint goal success rate

### OKR-3: KPI Engine (DONE)
- [x] `pm_kpi_define` — Define KPI with thresholds
- [x] `pm_kpi_snapshot` — Record measurement
- [x] `pm_kpi_dashboard` — All KPIs with trends
