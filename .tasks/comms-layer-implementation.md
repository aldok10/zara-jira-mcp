# Communication Layer Implementation Tasks

> Created: 2026-06-26
> Goal: Make zara-jira-mcp the most communication-literate PM MCP tool available.
> Research basis: `research/communication-research-extended.md` (18 research areas)
> Roadmap: `ROADMAP.md` → "Next Phase: AI Communication Layer"

---

## Phase 1: Signal-over-Noise (Next Sprint)

Priority: Ship tools that reduce PM communication noise and make decisions searchable.

### Task 1.1: `pm_comms_health` — Communication Health Score
- **What:** Composite score from existing data: re-decision rate, blocker escalation speed, stakeholder responsiveness, action item follow-through rate
- **Research:** Signal Retention (#15), Communication Metrics (#18)
- **Data sources:** decisions table (re-decisions = same tags, different outcomes), blockers table (time to first escalation), stakeholder_pulse (response trend), action_items (completion rate)
- **Output:** Score 0-100 + breakdown per dimension + trend over sprints
- **Effort:** Medium (aggregation of existing data, no new tables)
- **Files:** `application/tools/pm_comms_health.go`, `internal/sqlite/comms_health.go`

### Task 1.2: Enhance `pm_anti_patterns` with Lencioni Dysfunction Mapping
- **What:** Map existing anti-pattern detections → Lencioni pyramid level. Add new detections for Level 1-2 (trust, conflict avoidance).
- **Research:** Five Dysfunctions (#2)
- **New patterns:**
  - Dead retros + no improvements adopted = Fear of Conflict (Level 2)
  - Hero culture + no peer code review = Avoidance of Accountability (Level 4)
  - No sprint goals = Lack of Commitment (Level 3)
  - Carryover >40% repeatedly = Inattention to Results (Level 5)
  - Re-assignment ping-pong = Absence of Trust (Level 1)
- **Output:** Existing anti_patterns output + "Possible Dysfunction Level: X — [explanation]"
- **Effort:** Low (enhance existing AI prompt, add pattern rules)
- **Files:** `application/tools/pm_anti_patterns.go` (enhance prompt)

### Task 1.3: Notification Budget in `notify_routed`
- **What:** Track notifications sent per user/channel per day/week. Auto-throttle when budget exceeded. Surface stats in `pm_mcp_stats`.
- **Research:** Notification Fatigue (#6)
- **Rules:** Max 3-5 unsolicited notifications/day (configurable). Critical severity bypasses budget. Log all sends with outcome (acted-on? dismissed?).
- **Schema addition:** `notification_log` table (timestamp, channel, severity, message_hash, acknowledged)
- **Effort:** Low-Medium (new table, counter logic in routing)
- **Files:** `internal/notifications/budget.go`, `internal/sqlite/notification_log.go`

### Task 1.4: `pm_decision_record` (enhanced ADR/MADR format)
- **What:** Upgrade `pm_record_decision` with structured format: Context → Problem → Options Considered → Decision → Consequences → Review Date
- **Research:** Signal Retention (#15), Crucial Conversations (#11)
- **Backward compat:** Existing `pm_record_decision` still works. New tool adds structure.
- **Effort:** Low (new tool handler, same decisions table with richer JSON in context field)
- **Files:** `application/tools/pm_decision_record.go`

### Task 1.5: `pm_escalation_draft` with TIRED framework
- **What:** Generate escalation message from blocker/risk data using TIRED format (Timeframe, Impact, Requested action, Evidence, Deadline)
- **Research:** Already in `docs/communication-frameworks.md` section 9
- **Input:** `blocker_id` or `risk_id` or free-text situation
- **Output:** Ready-to-send message + suggested channel + severity classification
- **Effort:** Low (AI prompt + data gathering from existing tables)
- **Files:** `application/tools/pm_escalation_draft.go`

---

## Phase 2: Audience-Aware Routing (Sprint +1)

### Task 2.1: `pm_audience_router`
- **What:** Same data, auto-reframe per audience. Input: data + target audience. Output: reframed message.
- **Research:** SCARF (#docs section 2), Pyramid Principle, Information Asymmetry (#16)
- **Audiences:** exec (BLUF, no jargon, outcomes), tech-lead (metrics, blockers, capacity), PO (goal progress, scope, timeline), team (action items, recognition)
- **SCARF awareness:** exec needs Certainty (confidence intervals), dev needs Autonomy (options not mandates), PO needs Status updates (progress toward THEIR goals)
- **Effort:** Medium (AI prompt engineering per audience, wrapper over existing data tools)

### Task 2.2: `pm_silence_detector`
- **What:** Flag stakeholders with no interaction in N sprints. "Ghost stakeholder" alert.
- **Research:** Communication Anti-Patterns (#docs section 13)
- **Data:** stakeholder_pulse (no entries for stakeholder X in 3+ sprints), decision records (stakeholder not consulted), meeting_notes (stakeholder absent repeatedly)
- **Effort:** Low (query existing tables, threshold alert)

### Task 2.3: `pm_comms_anti_patterns`
- **What:** Dedicated communication anti-pattern detection (separate from sprint anti-patterns).
- **Research:** Communication Anti-Patterns (#docs section 13), Community Smells
- **Patterns to detect:**
  - Information hoarding (1 person on all blockers, never delegates)
  - Over-communication (too many notifications, no action taken)
  - Re-deciding (same topic in decisions table, 3+ times)
  - Meeting addiction (high meeting count, low decision-per-meeting)
  - Status theater (reports generated but health declining = nobody reads them)
  - One-way reporting (digests sent, no stakeholder pulse recorded)
- **Effort:** Medium

### Task 2.4: `pm_communication_plan`
- **What:** Generate stakeholder communication plan for any initiative. Who gets what, when, via which channel.
- **Research:** Kotter+ADKAR (#13), Stakeholder Mapping (#docs section 10)
- **Effort:** Medium (AI-generated, uses stakeholder data + severity rules)

---

## Phase 3: AI Coaching Communication (Sprint +2)

### Task 3.1: `pm_hard_conversation`
- **What:** Prep PM for difficult conversation with data + framework.
- **Research:** Crucial Conversations (#11), SBI, SCARF, Radical Candor
- **Output structure:**
  1. Facts from data (blocker age, health score, pattern)
  2. Possible interpretations (your story vs their story)
  3. SCARF risks (which domains threatened?)
  4. Opening line options (STATE path)
  5. Fallback if safety breaks (mutual purpose restoration)
- **Effort:** High (complex AI prompt, but high human value)

### Task 3.2: `pm_nvc_reframe`
- **What:** Reframe blaming/judgmental language → NVC format (Observation, Feeling, Need, Request)
- **Research:** NVC (#docs section), Motivational Interviewing (#10)
- **Input:** Raw message/feedback text
- **Output:** NVC-reframed version + explanation of what changed
- **Effort:** Low (pure AI prompt engineering)

### Task 3.3: `pm_trust_signals`
- **What:** Track trust indicators over time: forecast accuracy, override rate, team confidence trend, notification action rate
- **Research:** Trust Pyramid (#docs section 14), Notification Fatigue (#6)
- **Effort:** Medium (aggregate multiple data sources)

### Task 3.4: Enhanced `pm_coaching` with MI techniques
- **What:** When situation suggests resistance/ambivalence, coaching output uses MI approach (OARS: Open questions, Affirmations, Reflective listening, Summary)
- **Research:** Motivational Interviewing (#10)
- **Instead of:** "You should do X" → "What would recording decisions give your team?"
- **Effort:** Low (prompt engineering change in existing tool)

### Task 3.5: Enhanced `pm_facilitate` with Liberating Structures
- **What:** Add LS suggestions to ceremony facilitation. Include step-by-step scripts + timing.
- **Research:** Liberating Structures (#9)
- **Structures to include:** 1-2-4-All (retro), TRIZ (retro), 15% Solutions (planning), Wise Crowds (coaching), Troika (1-on-1 prep)
- **Effort:** Low (content addition to existing prompt)

---

## Phase 4: Organizational Communication (Sprint +3, stretch)

### Task 4.1: `pm_change_communication`
- **What:** Generate change communication plan using Kotter 8-step + ADKAR individual tracking
- **Research:** Kotter + ADKAR (#13)
- **Effort:** High

### Task 4.2: `pm_conflict_mediation`
- **What:** AI-assisted conflict diagnosis using Thomas-Kilmann modes + recommended resolution approach
- **Research:** TKI (#12)
- **Effort:** High

### Task 4.3: `pm_calibration_report`
- **What:** Historical accuracy of AI forecasts. "We said 85% chance by Thursday — did we make it?"
- **Research:** Trust Pyramid, EBM, Notification Fatigue (trust-building)
- **Effort:** Medium (compare forecast dates vs actual completion dates from snapshots)

### Task 4.4: `pm_ceremony_optimizer`
- **What:** Recommend which ceremonies to keep sync, go async, or kill entirely
- **Research:** Meeting ROI, Async Protocol, Signal-over-Noise
- **Input:** meeting_notes frequency + decisions-per-meeting + team size + stage
- **Effort:** Low (analysis of existing meeting data)

### Task 4.5: `pm_space_metrics`
- **What:** SPACE framework dashboard from existing data
- **Research:** SPACE (#4)
- **Mapping:** S=confidence+pulse, P=escaped defects+goal achievement, A=throughput, C=PR review time+decision speed, E=cycle time+WIP
- **Effort:** Medium

---

## Research Backlog (riset lanjutan sebelum implement)

| Topic | Why | Priority |
|-------|-----|----------|
| Edmondson's 7-item safety survey | Validated instrument for `pm_safety_check` | High (Phase 3) |
| PCE (Perceived Communication Effectiveness) scale | Validated comms measurement | Medium |
| PNAS nudge meta-analysis | Which nudge types work for PM context | Low |
| EBM metrics measurable from Jira+GitHub | Which KVAs can we auto-calculate | Medium |
| Westrum culture short-form survey | Culture indicator questions | Medium |
| LS + ceremony mapping | Which Liberating Structure for which ceremony | Low (before Task 3.5) |
| MI OARS examples for SM scenarios | Feed into coaching prompt | Low (before Task 3.4) |
| TKI behavioral markers from Jira data | Can we infer conflict mode without survey? | Low (Phase 4) |
| Notification frequency benchmarks | What's optimal for PM context | Medium (before Task 1.3) |
| Signal retention measurement | Return-to-dev rate, re-decision rate formulas | High (Phase 1) |

---

## Definition of Done (per task)

1. Tool implemented in `application/tools/`
2. AI prompt tested with 3+ scenarios
3. Registered in `transport/` tool registration
4. Added to `SKILL.md` tool reference
5. Added to relevant section in `docs/communication-frameworks.md`
6. Basic test coverage (happy path + error case)
7. Works with `PM_PROFILE=standard` or higher
