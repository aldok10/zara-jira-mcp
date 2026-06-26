# PM/SM in the AI Era — Comprehensive Research Summary

Research date: 2026-06-26. Compiled for roadmap implementation.

---

## 1. Evidence-Based Management (EBM) — Scrum.org

**What:** 4 Key Value Areas to measure organizational agility:
- **Current Value** (CV): Customer satisfaction, revenue per employee
- **Unrealized Value** (UV): Market gap, customer wishes unfulfilled
- **Ability to Innovate** (A2I): Innovation rate, cycle time, defect trends
- **Time to Market** (T2M): Release frequency, lead time, integration frequency

**Why it matters:** Most PMs measure OUTPUT (velocity, story points). EBM measures OUTCOME (business value delivered).

**Tool opportunity:**
- `pm_ebm_dashboard` — Track all 4 KVAs over time
- `pm_value_check` — "Is this sprint creating value or just output?"

---

## 2. Google Project Aristotle — Psychological Safety

**What:** 5 pillars of high-performing teams:
1. **Psychological Safety** (#1 predictor — 43% variance in performance)
2. **Dependability** (team delivers on commitments)
3. **Structure & Clarity** (clear roles, plans, goals)
4. **Meaning** (work matters personally)
5. **Impact** (work matters to the world)

**7 statements to measure psychological safety:**
1. "If I make a mistake, it won't be held against me"
2. "Team members can bring up problems and tough issues"
3. "People never reject others for being different"
4. "It's safe to take a risk"
5. "It's easy to ask for help"
6. "No one would sabotage my efforts"
7. "My unique skills are valued"

**Tool opportunity:**
- `pm_safety_survey` — Record team answers to 7 questions (1-5 scale)
- `pm_safety_trend` — Track psychological safety over time
- `pm_team_aristotle` — Full 5-pillar assessment

---

## 3. SPACE Framework — Developer Productivity

**What:** 5 dimensions (Microsoft/GitHub research):
- **S**atisfaction and well-being
- **P**erformance (outcomes, quality)
- **A**ctivity (actions taken, measurable)
- **C**ommunication and collaboration
- **E**fficiency and flow

**Key insight:** Never measure just ONE dimension. Activity without Satisfaction = burnout. Performance without Communication = silos.

**Tool opportunity:**
- `pm_space_metrics` — Track all 5 dimensions from Jira + survey data
- `pm_flow_disruption` — Detect context-switching signals (many small tasks, rapid priority changes)

---

## 4. Team Topologies — Interaction Modes

**What:** 4 team types × 3 interaction modes:
- Types: Stream-aligned, Platform, Enabling, Complicated-subsystem
- Modes: Collaboration, X-as-a-Service, Facilitating

**Key insight:** "Permanent collaboration = hidden dependency = slower delivery." Interactions should be intentional and time-bounded.

**Tool opportunity:**
- `pm_team_dependencies` — Map cross-team interaction modes
- `pm_cognitive_load` — Assess team cognitive load (too many services, domains, responsibilities)

---

## 5. Liberating Structures — Facilitation

**What:** 33 micro-structures for group engagement. Most relevant for SM:
- **1-2-4-All** — Progressive disclosure (individual → pairs → groups → whole)
- **Troika Consulting** — 3-person peer coaching rounds
- **15% Solutions** — What can you do RIGHT NOW without permission?
- **TRIZ** — "What must we stop doing?"
- **Conversation Cafe** — Structured dialogue with rounds
- **Discovery & Action Dialogue** — Find positive deviance

**Tool opportunity:**
- `pm_facilitate` (already exists) — Expand with Liberating Structures library
- `pm_retro_format` — Suggest format based on team situation + history

---

## 6. #NoEstimates + Probabilistic Forecasting

**What:** Replace deterministic estimates with:
- Monte Carlo simulation (already implemented!)
- Throughput-based forecasting
- Right-sizing (break to similar-sized items)
- "How many items can we do?" not "How long will these take?"

**Key insight:** "Story points are the WORST metric except for all the others. But Monte Carlo + throughput is better than both."

**Tool opportunity:**
- Already have `pm_forecast` — Monte Carlo ✓
- Add: `pm_right_size` — Detect oversized tickets that reduce predictability
- Add: `pm_estimation_accuracy` — How accurate were past estimates?

---

## 7. OKR-Sprint Alignment

**What:** Connect every sprint to a business objective:
- Objective: Qualitative, inspiring goal
- Key Results: Quantitative measures of success
- Hypothesis: "We believe [doing X] will result in [outcome Y] as measured by [metric Z]"

**Key insight:** Teams that connect sprints to OKRs are 30% more likely to hit goals (OKR Institute data).

**Tool opportunity:**
- `pm_outcome_map` (already exists) — enhance with hypothesis tracking
- `pm_hypothesis` — Record hypothesis per feature, track if validated
- `pm_outcome_review` — "Did last sprint's work actually move the OKR needle?"

---

## 8. Servant Leadership Measurement

**What:** SM value is INVISIBLE if measured wrong. Wrong metrics: velocity, burndown. Right metrics:
- Time to impediment removal (faster = more SM value)
- Team autonomy growth (less SM intervention = success)
- Psychological safety trend (up = SM doing coaching right)
- Stakeholder satisfaction trajectory
- Experiment completion rate (retro actions actually done)

**Key insight:** "SM's best work makes the SM unnecessary. Measure team independence, not SM busyness."

**Tool opportunity:**
- `pm_sm_impact` (already exists) — enhance with autonomy metrics
- `pm_invisible_work` — Track the untracked: conversations, unblockings, coaching moments

---

## 9. DevEx & Flow State

**What:** Developer Experience is PM's responsibility too:
- **Flow state** — uninterrupted deep work time
- **Cognitive load** — how many things compete for attention
- **Feedback loops** — how fast do you know if your code works

**Key PM actions:** Protect maker time, reduce WIP, shorten meetings, batch interruptions.

**Tool opportunity:**
- `pm_maker_time` — Analyze calendar for meeting-free blocks
- `pm_wip_guardian` (already exists concept) — enhance with cognitive load signals

---

## 10. Async-First Communication

**What:** The future of PM communication is:
- Written > spoken (searchable, inclusive, timezone-friendly)
- Decisions documented > decisions in meetings
- Status auto-generated > status asked in standups
- Context pushed > context pulled

**Key insight (Discourse CEO 2026):** "Companies with written communication develop verifiable shared memory. Everything else is amnesia."

**Tool opportunity:**
- `pm_async_update` — Generate team status that replaces a meeting
- `pm_decision_trail` — All decisions searchable, linked to outcomes
- `pm_meeting_audit` — "Could this meeting be an async update instead?"

---

## Summary: What We Already Have vs What's Missing

### Already Implemented (strong):
- Monte Carlo forecasting ✓
- Sprint health + predictability ✓
- Team pulse + radar ✓
- Anti-pattern detection ✓
- SM impact tracking ✓
- Async standup ✓
- 1-on-1 prep ✓
- Commitment tracking ✓

### Partially Done (needs enhancement):
- OKR alignment (basic `pm_outcome_map`, needs hypothesis tracking)
- Facilitation (`pm_facilitate` exists, needs Liberating Structures library)
- Skills matrix (exists, needs cognitive load assessment)

### Not Yet Implemented (high value):
- Psychological safety survey (Project Aristotle)
- EBM 4 KVAs dashboard
- SPACE metrics aggregation
- Hypothesis-driven development tracking
- Estimation accuracy feedback loop
- Meeting audit ("should this be async?")
- Cognitive load assessment
- Right-sizing ticket detection
- Team topology dependency mapping

