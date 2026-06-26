# Communication Framework Implementation — Task Breakdown

> Created: 2026-06-26
> Status: Planning
> Goal: Make zara-jira-mcp communication-literate for AI-era PMs

---

## Priority Matrix

### P0: Quick Wins (Apply to existing tools, no new code)

- [ ] **Audit all AI prompts for BLUF compliance** — Every tool output should lead with the answer
  - Files: `application/tools/reporting_handlers.go`, `coaching_handlers.go`, `smart_handlers.go`
  - Check: Does `pm_exec_report` output BLUF first? Does `pm_coaching` lead with actionable insight?
  - Effort: 2-3h (prompt rewording only)

- [ ] **Add confidence signaling to AI outputs** — "Based on X sprints..." / "High/Medium/Low confidence"
  - Files: `application/tools/forecast_handlers.go`, `smart_handlers.go`
  - Pattern: `[INSIGHT] ... [CONFIDENCE] ... [DATA SOURCE] ... [RECOMMENDATION]`
  - Effort: 2-4h

- [ ] **Audit language calibration** — Ensure `pm_exec_report` never outputs "story points", "velocity", "WIP" to executives
  - Files: `application/tools/reporting_handlers.go`, `management_handlers.go`
  - Replace: "velocity" → "delivery speed", "WIP" → "items in progress", "story points" → "work items"
  - Effort: 1-2h

### P1: New Communication Tools (Medium effort, high value)

- [ ] **`pm_communication_health`** — Score team communication patterns
  - Metrics: decision documentation rate, meeting ROI, escalation speed, retro action follow-through
  - Input: board_id
  - Output: score 0-100 + breakdown + recommendations
  - Dependencies: needs data from `pm_record_decision`, `pm_meeting_roi`, `pm_impediment_aging`
  - Effort: 1 day

- [ ] **`pm_audience_translate`** — Reformat any report for different audience level
  - Input: original_text + target_audience (exec/po/team/external)
  - Output: reformatted text with appropriate vocabulary and depth
  - Uses AI to translate jargon levels
  - Effort: half day

- [ ] **`pm_notification_audit`** — Analyze notification volume and response rates
  - Input: board_id, days (default 30)
  - Output: notifications/day per member, response rate, fatigue indicators
  - Needs: notification log table in SQLite
  - Effort: 1-2 days (includes schema change)

- [ ] **`pm_decision_debt`** — Surface undocumented decisions, orphaned ADRs
  - Input: board_id
  - Output: decisions without DACI, meetings without action items, stale knowledge base pages
  - Effort: half day

### P2: Communication Infrastructure (Higher effort, foundational)

- [ ] **Notification batching engine** — Aggregate non-urgent notifications into digest
  - Rule: Max 5-10 notifications/day per user. Batch P3-P5 into morning digest
  - Requires: notification queue + user timezone + preference store
  - Files: `internal/` new module
  - Effort: 2-3 days

- [ ] **Escalation acknowledgment tracking** — Closed-loop communication
  - Pattern: sent → acknowledged → acted-on. Auto-re-escalate on silence
  - Requires: schema addition for escalation state machine
  - Effort: 1-2 days

- [ ] **Communication plan generator** — `pm_communication_plan`
  - Auto-generate stakeholder x channel x frequency matrix from team/board structure
  - Uses Mendelow grid logic: map Jira roles → power/interest → engagement strategy
  - Effort: 1 day

- [ ] **Tone/culture configuration** — Allow teams to set communication style
  - Options: direct (Western) / supportive (Asian) / mixed
  - Affects: coaching prompts, escalation language, notification tone
  - Store: team_settings table
  - Effort: 1 day

### P3: Research & Validation (Needs investigation before implementation)

- [ ] **Measure notification effectiveness** — Do people read and act on AI-generated reports?
  - Design: Add delivery tracking (sent/opened/acted-on) to notification system
  - Validate: Is our notification volume within 5-10/day sweet spot?
  - Effort: research 2h, implementation 1 day

- [ ] **A/B test report formats** — Narrative storytelling vs bullet list vs dashboard
  - Hypothesis: Narrative format improves exec engagement
  - Method: Generate same data in 2 formats, track which gets responses
  - Effort: research + design 4h, implementation TBD

- [ ] **Cross-cultural communication study** — How do Indonesian teams prefer AI communication?
  - Questions: Direct vs indirect tone? Bahasa vs English for reports? Individual vs group notifications?
  - Method: Interview PM users, iterate on tone
  - Effort: ongoing

- [ ] **AI de-skilling risk assessment** — Does the tool make PMs lazier or more capable?
  - Design: Track PM engagement patterns over time (are they overriding less? Asking fewer questions?)
  - Metric: Override rate, custom input frequency, coaching topic diversity
  - Effort: research 2h, monitoring implementation 1 day

---

## Implementation Order (Recommended)

```
Week 1 (Quick Wins):
  Day 1: BLUF audit + language calibration on existing tools
  Day 2: Add confidence signaling to forecast + coaching outputs
  Day 3: pm_audience_translate (simple AI-powered rewriter)

Week 2 (New Tools):
  Day 1: pm_communication_health (scoring)
  Day 2: pm_decision_debt + pm_communication_plan
  Day 3: Notification batching design + schema

Week 3 (Infrastructure):
  Day 1-2: Notification batching engine
  Day 3: Escalation acknowledgment tracking

Week 4 (Polish):
  Day 1: Tone/culture configuration
  Day 2: Notification audit tool
  Day 3: Testing + documentation
```

---

## Spec Notes for Key Tools

### pm_communication_health

```
Input:  board_id (required)
Output: {
  score: 72,
  breakdown: {
    decision_documentation_rate: 80,  // % decisions with DACI + record
    meeting_action_followthrough: 65, // % action items that became tickets
    escalation_speed: 90,             // avg time from blocker → escalation
    retro_action_completion: 55,      // % retro items completed
    stakeholder_pulse_trend: 78,      // satisfaction direction
  },
  recommendations: [
    "Decision documentation rate dropped 15% vs last sprint. Record decisions immediately after meetings.",
    "3 retro action items from Sprint 22 still unresolved. Consider adding to Sprint 24 commitment."
  ]
}
```

### pm_audience_translate

```
Input:  text (required), target_audience: "exec" | "po" | "team" | "external" (required)
Output: reformatted text with:
  - Vocabulary adapted (no jargon for exec)
  - Depth adjusted (3 sentences for exec, full detail for team)
  - Tone matched (formal for external, casual for team)
  - Action items highlighted
```

### pm_notification_audit

```
Input:  board_id, days: 30
Output: {
  total_sent: 245,
  avg_per_user_per_day: 8.2,  // WARNING: above 5-10 threshold
  by_severity: { P1: 3, P2: 12, P3: 80, P4: 150 },
  response_rate: 0.34,        // Only 34% acknowledged
  fatigue_indicators: [
    "User Alice has 0% response rate on P4 notifications — consider suppressing",
    "Channel #sprint-updates has <5% engagement — consider digest format"
  ],
  recommendation: "Reduce P4 notifications by 60%. Batch into weekly digest."
}
```

---

## Research Backlog (For Deeper Dives Later)

| Topic | Why | Priority |
|-------|-----|----------|
| JTBD for notifications | "What job is this notification hired to do?" — filter useless notifs | Medium |
| Narrative arc in AI reports | Situation → Complication → Resolution structure | Medium |
| Closed-loop communication protocol | Aviation CRM adapted for PM escalations | High |
| Media Richness mapping | Which PM ceremony needs which channel richness level | Low |
| Nudge design for PM behavior | Default to recording decisions, nudge documentation | Medium |
| Multi-language AI output | Seamless Indo/English mixing in generated reports | Low |
| PM de-skilling longitudinal study | Long-term effect of AI tool dependency | Low (future) |
| Notification cost calculator | "This notification costs 23 min of deep work. Is it worth it?" | High |

---

## Success Criteria

| Metric | Current | Target |
|--------|---------|--------|
| All AI outputs BLUF-first | ~50% | 100% |
| Confidence signaling in forecasts | 0% | 100% |
| Max notifications/user/day | uncapped | < 10 |
| Executive reports without jargon | ~70% | 100% |
| Decision documentation rate | unmeasured | > 80% |
| Notification response rate | unmeasured | > 50% |
| Escalation acknowledgment | not tracked | tracked |
| Communication health score | not implemented | available |
