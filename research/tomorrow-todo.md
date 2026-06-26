# Research Summary + TODO for Tomorrow

> Compiled research across: team psychology, AI era challenges, communication,
> servant leadership, and continuous improvement. Each section ends with TODO.

---

## 1. TEAM MATURITY (Tuckman + Scrum.org)

**Finding:** Teams go through Forming → Storming → Norming → Performing. SM role changes at each stage:
- Forming: Direct (teach process)
- Storming: Coach (resolve conflict)
- Norming: Support (enable decisions)
- Performing: Delegate (step back)

**Finding (Springer 2021):** Scrum Masters initially play 9 leadership roles. They transfer roles to team as it matures. The "leadership gap" (intentionally stepping back) is what enables team autonomy.

**Finding:** Servant leadership → 42.5% higher on-time delivery, 37% less stakeholder dissatisfaction vs command-and-control.

**TODO:**
- [ ] Add `pm_maturity_assess` — assess team stage based on data signals (conflict in retros = storming, stable velocity = norming, self-correcting = performing)
- [ ] Add `pm_sm_role_suggest` — suggest SM behavior based on team stage ("your team is norming, step back from facilitating")
- [ ] Track "leadership gap" metric: ratio of SM-initiated actions vs team-initiated

---

## 2. RETROSPECTIVE EFFECTIVENESS

**Finding:** Most retros fail because:
- Action items never completed (>5 pending = "dead retros")
- Same format every time → engagement drops
- No measurement of improvement over time
- Actions too vague ("improve communication" vs "pair on API tasks Tues/Thu")

**Finding:** SMART action items have 3x higher completion rate.

**TODO:**
- [ ] Add `pm_retro_score` — rate retro quality: were actions SMART? were previous actions completed? did participation improve?
- [ ] Add `pm_improvement_velocity` — track how fast team improves (improvements implemented per sprint)
- [ ] Auto-detect vague actions and suggest SMART rewording

---

## 3. AI BRAIN FRY & TOOL OVERLOAD (BCG 2026)

**Finding:** 3 AI tools = productive. 4+ = cognitive collapse.
**Finding:** High performers most affected (they try all tools).
**Finding:** "Cognitive debt" — AI erodes critical thinking over time.

**Key quote:** "Organizations have a work design problem, not a skills problem."

**TODO:**
- [x] Profile system (done — chatgpt=20 tools)
- [x] pm_smart NL router (done — one entry point)
- [ ] Add "tool usage stats" — track which tools PM actually uses (80/20 rule)
- [ ] Auto-suggest profile upgrade/downgrade based on usage pattern

---

## 4. STAKEHOLDER TRUST & SATISFACTION

**Finding:** Only 44% perception match between stakeholders and project team.
**Finding:** Transparency builds trust, but too much information destroys it.
**Finding:** Clients and end users are most important stakeholders; they also cause most problems.

**Finding (Gallup):** 70% of team engagement variance is attributed to the manager.

**TODO:**
- [ ] Add `pm_stakeholder_register` — track who needs what info, how often
- [ ] Add `pm_stakeholder_pulse` — quick satisfaction check (1-5 after every sprint review)
- [ ] Auto-suggest communication frequency based on stakeholder power/interest matrix
- [ ] Add "expectation gap" detection: are stakeholders expecting more than team can deliver?

---

## 5. REMOTE/ASYNC TEAM DYNAMICS

**Finding (Gallup):** Remote workers are most engaged BUT loneliest AND most likely to leave.
**Finding:** One in-person day/month boosts performance AND retention.
**Finding:** Async standups work IF: written, brief, include "needs help" signal.
**Finding:** Engagement must be designed, not assumed.

**Finding (Frontiers in Psychology 2024):** Interactional monitoring (check-ins) improves engagement. Electronic monitoring (surveillance) destroys it.

**TODO:**
- [ ] Add `pm_async_standup` — structured async standup format (progress, plan, blockers, mood 1-5)
- [ ] Add `pm_team_connectedness` — track signals of isolation (no PR reviews, no retro comments, no collaboration)
- [ ] Distinguish "check-in" (caring) from "monitoring" (surveillance) in tool descriptions
- [ ] Add "remote team health" signals to pm_team_care

---

## 6. KNOWLEDGE SHARING & BUS FACTOR

**Finding:** AI pair programming reduces knowledge sharing (devs work with AI, not humans).
**Finding:** Teams with high code review participation have lower bus factor.
**Finding:** Bus factor = 1 is the #1 risk for team continuity.

**TODO:**
- [ ] Add `pm_bus_factor` — detect single points of failure (one person owns all of an area)
- [ ] Track PR review distribution (same reviewer always = knowledge silo)
- [ ] Add "knowledge map" — who knows what (based on issue assignment history)

---

## 7. CONTINUOUS IMPROVEMENT MEASUREMENT

**Finding:** Teams that track improvement velocity (changes implemented per sprint) improve 2x faster.
**Finding:** The best teams have a "meta-process" for improving their process.
**Finding:** Sprint Goal Achievement Rate is the single best predictor of team effectiveness (Scrum Alliance).

**Key metrics to track:**
- Sprint Goal Success Rate (binary: hit/miss per sprint)
- Action Item Completion Rate (retro actions completed / total)
- Cycle Time Trend (improving = getting better)
- Escaped Defect Rate (decreasing = quality improving)
- Team Satisfaction Score (pulse check 1-5)

**TODO:**
- [ ] Add `pm_improvement_dashboard` — meta-metrics: are we getting better at getting better?
- [ ] Track sprint goal hit rate as binary yes/no (most important single metric)
- [ ] Add "time to value" metric: from idea to user hands

---

## 8. COMMUNICATION IN PRACTICE

**Finding:** Best orgs have "writing culture" — writing IS the decision medium, not recording medium.
**Finding:** RFCs before ADRs — the "maybe" layer that saves bad decisions.
**Finding:** Specs are infrastructure in the age of agents.
**Finding:** In-person: tone, body language. Text: all that disappears → misunderstandings 3x more likely.

**TODO:**
- [x] Communication frameworks implemented (BLUF, SBI, SCQA, Pyramid)
- [ ] Add `pm_rfc_template` — generate RFC structure for technical decisions
- [ ] Add `pm_decision_announce` — template for communicating decisions to different audiences
- [ ] Track "decision quality" — decisions that got reversed (signal of poor initial communication)

---

## 9. EVIDENCE-BASED MANAGEMENT (EBM, Scrum.org)

Four Key Value Areas:
1. **Current Value** — customer satisfaction, employee satisfaction
2. **Unrealized Value** — customer satisfaction gap, market share gap
3. **Ability to Innovate** — % time on new vs maintenance, tech debt ratio
4. **Time to Market** — lead time, deploy frequency, cycle time

**TODO:**
- [ ] Add `pm_ebm_report` — generate EBM scorecard from available data
- [ ] Map existing tools to EBM areas:
  - Current Value: sprint goal hit rate, stakeholder pulse
  - Unrealized Value: backlog size growth
  - Ability to Innovate: tech debt ratio, time on features vs bugs
  - Time to Market: cycle time, deploy frequency

---

## 10. NUMBERS TO REMEMBER

| Stat | Source |
|------|--------|
| 42.5% higher delivery with servant leadership | IJSRA 2024 |
| 37% less stakeholder dissatisfaction | Same |
| 43% performance variance = psychological safety | Google Aristotle |
| 52% teams hit sprint goals | Scrum Alliance |
| 70% engagement variance = the manager | Gallup |
| 44% perception match stakeholders vs team | ResearchGate |
| 3 AI tools max before brain fry | BCG 2026 |
| 20% PMs have good AI experience | PMI |
| 23 min per context switch | Multiple |
| 1 in-person day/month → +performance | CEPR 2025 |
| 9 leadership roles of SM (transfer over time) | Springer |
| 3x better action completion with SMART format | Multiple |

---

## PRIORITY ORDER FOR TOMORROW

### P0 (High Impact, Low Effort)
1. `pm_improvement_dashboard` — are we getting better?
2. `pm_bus_factor` — single point of failure detection
3. `pm_async_standup` — structured async format
4. Sprint Goal binary tracking (hit/miss per sprint)

### P1 (High Impact, Medium Effort)
5. `pm_maturity_assess` — team stage detection
6. `pm_stakeholder_register` — who needs what info
7. `pm_ebm_report` — Evidence-Based Management scorecard
8. Tool usage stats (which tools PM actually uses)

### P2 (Medium Impact, Nice to Have)
9. `pm_rfc_template` — RFC generator
10. `pm_team_connectedness` — isolation detection
11. Knowledge map from issue history
12. Decision reversal tracking
