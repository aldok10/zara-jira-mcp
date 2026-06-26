# Extended Communication Research for zara-jira-mcp

> Compiled: 2026-06-26
> Purpose: Deep research summaries for communication layer implementation. Each section = one research area with summary, key findings, and how it maps to zara-jira-mcp.
> Action: Use this as reference when implementing Comms Phase 1-4 tools.

---

## INDEX

1. [Psychological Safety & Project Aristotle](#1-psychological-safety--project-aristotle)
2. [Five Dysfunctions of a Team (Lencioni)](#2-five-dysfunctions-of-a-team-lencioni)
3. [Westrum Organizational Culture](#3-westrum-organizational-culture)
4. [SPACE Framework (Developer Productivity)](#4-space-framework-developer-productivity)
5. [Evidence-Based Management (EBM)](#5-evidence-based-management-ebm)
6. [Notification Fatigue & Attention Budget](#6-notification-fatigue--attention-budget)
7. [Team Topologies & Cognitive Load](#7-team-topologies--cognitive-load)
8. [Nudge Theory & Choice Architecture](#8-nudge-theory--choice-architecture)
9. [Liberating Structures](#9-liberating-structures)
10. [Motivational Interviewing for Coaching](#10-motivational-interviewing-for-coaching)
11. [Crucial Conversations](#11-crucial-conversations)
12. [Thomas-Kilmann Conflict Modes](#12-thomas-kilmann-conflict-modes)
13. [Kotter + ADKAR Change Communication](#13-kotter--adkar-change-communication)
14. [Johari Window](#14-johari-window)
15. [Shannon-Weaver & Signal Retention](#15-shannon-weaver--signal-retention)
16. [Information Asymmetry / Principal-Agent](#16-information-asymmetry--principal-agent)
17. [Tuckman Stages + AI Assessment](#17-tuckman-stages--ai-assessment)
18. [Communication Metrics & Measurement](#18-communication-metrics--measurement)

---

## 1. Psychological Safety & Project Aristotle

**Source:** Google re:Work (2016), Amy Edmondson (Harvard), Frontiers in Psychology (2020)

**Summary:** Google studied 180+ teams over 2 years. #1 predictor of team performance: psychological safety — "shared belief that the team is safe for interpersonal risk taking." Correlated with 43% of variance in performance. Teams with high psych safety: less turnover, more diverse ideas harnessed, 2x rated effective by executives.

**5 Dynamics (ranked):**
1. Psychological Safety — can I take risks without punishment?
2. Dependability — can I count on teammates?
3. Structure & Clarity — are goals/roles/plans clear?
4. Meaning — is work personally meaningful?
5. Impact — do we believe our work matters?

**Key Insight:** Without safety, agile ceremonies become hollow rituals. Standups become status theater. Retros become performative. Anti-patterns flourish.

**Mapping to zara-jira-mcp:**
- `pm_anti_patterns` — detects symptoms of low safety (dead retros, hero culture = people afraid to ask for help)
- `pm_confidence` — anonymous confidence votes = proxy for safety
- `pm_coaching(topic:"psychological_safety")` — actionable coaching
- **Gap:** No direct measurement of safety. Could add `pm_safety_check` — 7-question Edmondson survey tracked over time.

**Further Research Needed:**
- Edmondson's 7-item safety questionnaire (validated instrument)
- How to detect safety decline from Jira data patterns (increasing re-assignment, blocker age, carryover)

---

## 2. Five Dysfunctions of a Team (Lencioni)

**Source:** Patrick Lencioni (2002), Stratrix framework analysis

**Summary:** Pyramid of dysfunctions — each builds on the one below:

```
5. INATTENTION TO RESULTS  (ego > team goals)
4. AVOIDANCE OF ACCOUNTABILITY  (no peer pressure)
3. LACK OF COMMITMENT  (ambiguity, no buy-in)
2. FEAR OF CONFLICT  (artificial harmony)
1. ABSENCE OF TRUST  (vulnerability-based trust missing)
```

**Key Insight:** Most SM interventions address symptoms (level 4-5) without fixing root cause (level 1-2). If trust is absent, accountability tools are useless.

**Mapping to zara-jira-mcp:**
- `pm_anti_patterns` — can detect some level 4-5 (hero culture = inattention to collective results)
- `pm_maturity_assessment` — related to trust/stage
- **Gap:** No dysfunction-level diagnosis. Could map anti-patterns → Lencioni levels. "Dead retros" = fear of conflict. "Hero culture" = avoidance of accountability. "No sprint goals" = lack of commitment.

**Further Research Needed:**
- Lencioni's Team Assessment questionnaire (public domain summary)
- Mapping specific Jira patterns → dysfunction level

---

## 3. Westrum Organizational Culture

**Source:** Ron Westrum (2004), DORA State of DevOps (2019), dora.dev

**Summary:** Three culture types based on how organizations process information:

| Aspect | Pathological | Bureaucratic | Generative |
|--------|-------------|--------------|-----------|
| Power | Fear-based | Rule-based | Mission-based |
| Messengers | Shot | Neglected | Trained |
| Responsibilities | Shirked | Narrow | Shared |
| Bridging | Discouraged | Tolerated | Encouraged |
| Failure | Punished | Justice | Inquiry |
| Novelty | Crushed | Problems | Implemented |
| Information | Hidden | Ignored | Exploited |

**Key Insight:** DORA 2019 proved generative culture correlates with superior delivery performance, job satisfaction, and organizational goal attainment. Psychological safety IS the marker of generative culture.

**Mapping to zara-jira-mcp:**
- `pm_maturity_assessment` — could include Westrum culture indicator
- Blocker patterns reveal culture: chronic blockers + no escalation = pathological. Blockers resolved quickly + lessons recorded = generative.
- **Gap:** No explicit culture assessment. Could infer from: blocker resolution time, decision record quality, retro action follow-through, escalation patterns.

**Further Research Needed:**
- Westrum culture survey (short form, validated)
- Correlation between culture score and sprint health over time

---

## 4. SPACE Framework (Developer Productivity)

**Source:** Nicole Forsgren et al. (2021), Microsoft Research, ACM Queue

**Summary:** Developer productivity has 5 dimensions:
- **S**atisfaction and Well-being
- **P**erformance (quality of outcomes)
- **A**ctivity (volume of output)
- **C**ommunication and Collaboration
- **E**fficiency and Flow

**Key Insight for Communication:** The "C" dimension explicitly measures: quality of handoffs, ease of integration, documentation discoverability, PR review quality. GitHub Copilot study (2024) found AI impacts "Performance" and "Communication" dimensions differently than perceived.

**Mapping to zara-jira-mcp:**
- `pm_flow_metrics` — covers Efficiency (cycle time, WIP)
- `pm_github_pr_metrics` — covers Communication (PR review time)
- `pm_team_health` — partial Satisfaction
- **Gap:** No explicit SPACE dashboard. Could add `pm_space_metrics` — integrate all 5 dimensions from existing data.

**Further Research Needed:**
- SPACE survey items for Communication dimension
- How to infer Satisfaction from behavioral signals (not just surveys)

---

## 5. Evidence-Based Management (EBM)

**Source:** Scrum.org (2020+), Ken Schwaber / Patricia Kong

**Summary:** EBM framework measures value delivery through 4 Key Value Areas (KVAs):
1. **Current Value (CV)** — value delivered NOW to stakeholders
2. **Unrealized Value (UV)** — potential future value (market gap)
3. **Ability to Innovate (A2I)** — capacity for new value (tech debt, skills)
4. **Time to Market (T2M)** — speed of delivering new capability

**Key Insight:** Traditional metrics (velocity, burndown) measure Activity, not Value. EBM connects sprint work → business outcomes. Exec reports should use EBM language.

**Mapping to zara-jira-mcp:**
- `pm_exec_report` — reports outcomes, not activity (already EBM-aligned)
- `pm_outcome_map` — connects sprints to objectives
- `pm_tech_debt_budget` — A2I dimension (tech debt erodes innovation)
- **Gap:** No explicit EBM KVA tracking. Could add `pm_ebm_dashboard` — map existing data to 4 KVAs.

**Further Research Needed:**
- EBM guide metrics list (which are measurable from Jira+GitHub data)
- UV measurement (requires product analytics integration — out of scope?)

---

## 6. Notification Fatigue & Attention Budget

**Source:** TianPan.co (2026), ACM CACM, Coralogix, Courier

**Summary:** Users have a daily ceiling of 3-5 unsolicited AI updates across ALL sources. The agent that sends 10th notification/week gets muted by Friday, uninstalled by next month.

**Key findings:**
- Proactive AI hits hard ceiling at user attention (TianPan.co 2026)
- AI agents should: receive event → evaluate → take action if possible → ONLY alert if human judgment required
- Alert fatigue desensitizes users, causing missed critical alerts
- Best alerting systems "earn trust by interrupting only when it counts"
- Notification is a product problem, not a marketing problem — fix in infrastructure

**Mapping to zara-jira-mcp:**
- `notify_routed` — already routes by severity (good!)
- `pm_escalate` — only fires on thresholds (3-day blocker, health<50)
- **Gap:** No notification budget tracking. No "did user act on notification?" feedback loop. Could add `pm_notification_budget` — track sends/week, click-through, dismiss rate. Auto-throttle when budget exceeded.

**Further Research Needed:**
- Optimal notification frequency for PM context (daily? weekly?)
- "Notification budget" implementation in existing MCP tools
- A/B testing frameworks for notification effectiveness

---

## 7. Team Topologies & Cognitive Load

**Source:** Skelton & Pais (2019), Martin Fowler, DevOps community

**Summary:** 4 team types optimized for flow:
1. **Stream-aligned** — delivers end-to-end value for one stream
2. **Platform** — provides self-service capability to reduce cognitive load
3. **Enabling** — helps stream-aligned teams overcome obstacles
4. **Complicated subsystem** — handles deep specialist knowledge

3 interaction modes: Collaboration, X-as-a-Service, Facilitating

**Key Insight:** When teams exceed 15 people, trust relationships break → cognitive load increases. Communication patterns MUST change based on team type and interaction mode. PM tool = enabling team.

**Mapping to zara-jira-mcp:**
- zara-jira-mcp IS the platform/enabling pattern for PM capability
- `pm_dependencies` — maps cross-team interactions
- `portfolio_workload` — cognitive load proxy (too many domains per person)
- **Gap:** No team topology awareness. Could use Jira component/project assignments to detect: is this team stream-aligned? Who's the platform team? Are interaction modes healthy?

**Further Research Needed:**
- Detecting team type from Jira data patterns
- Cognitive load scoring from issue assignments (breadth of domains)

---

## 8. Nudge Theory & Choice Architecture

**Source:** Thaler & Sunstein (2008), PNAS meta-analysis, SUE Behavioural Design

**Summary:** Nudging = redesigning choice environment so desired behavior is easiest/most natural option. Not telling people what to do, but adjusting: order of options, default settings, social cues, message framing.

**Key Insight for PM tools:** Every AI recommendation is a nudge. The ORDER in which you present options matters. Defaults matter. Social proof matters ("80% of teams that did X saw Y improvement").

**Application ideas:**
- `pm_recommendations` already nudges — but could be more explicit about defaults
- Sprint planning: present "recommended commitment" as default (anchoring bias)
- Risk recording: make it frictionless (one-click = nudge toward recording)
- Retro actions: show "teams that follow up on actions improve velocity 15%" (social proof)

**Mapping to zara-jira-mcp:**
- `pm_next` — already a nudge ("here's what you should do next")
- `pm_quickstart` — onboarding nudges
- **Gap:** Not designed with deliberate choice architecture. Could audit all tool outputs for nudge opportunities.

**Further Research Needed:**
- PNAS meta-analysis results: which nudge types most effective?
- Ethical boundaries of nudging in PM tools (libertarian paternalism)
- Default effect in sprint planning recommendations

---

## 9. Liberating Structures

**Source:** Lipmanowicz & McCandless, liberatingstructures.com, Scrum.org

**Summary:** 33 microstructures (interaction patterns) that include every voice in group work. Replace "big five" conventional structures (presentations, managed discussions, open discussions, status reports, brainstorms).

**Key structures for PM:**
- **1-2-4-All** — silent reflection → pairs → quads → full group. Engages introverts.
- **TRIZ** — "What must we stop doing?" Identifies counterproductive activities.
- **Wise Crowds** — structured peer coaching in 15 min.
- **15% Solutions** — "What can you do RIGHT NOW with what you have?"
- **Troika Consulting** — 3-person coaching rounds.

**Key Insight:** Retros and standups default to "open discussion" which lets loudest voices dominate. LS provides equal-voice alternatives.

**Mapping to zara-jira-mcp:**
- `pm_facilitate(ceremony:"retro")` — already suggests different formats
- **Gap:** Could explicitly reference LS structures. "Try 1-2-4-All format for this retro" with step-by-step. Could add timing suggestions and facilitation scripts.

**Further Research Needed:**
- Which LS structures map to which agile ceremonies?
- Can AI generate facilitation scripts for each LS?
- LS combinations (LS strings) for sprint planning, retros, demos

---

## 10. Motivational Interviewing for Coaching

**Source:** Miller & Rollnick (2013), Positive Psychology, coaching psychology research

**Summary:** Evidence-based approach for working with ambivalence about change. NOT pushing — evoking internal motivation. Core technique: OARS (Open questions, Affirmations, Reflective listening, Summary reflections).

**Core sequence:** ambivalence → evoked change talk → commitment → action

**Key types of talk:**
- **Change talk** — client expressing reasons/desire/ability to change
- **Sustain talk** — client expressing reasons to stay same
- Therapist role: amplify change talk, don't argue with sustain talk

**Key Insight for PM coaching:** When team is resistant to process change (e.g., refusing to do retros, not recording decisions), pushing harder = more resistance. MI approach: "What would recording decisions give you?" (evoke) vs "You need to record decisions" (direct).

**Mapping to zara-jira-mcp:**
- `pm_coaching(topic, situation)` — could use MI principles in coaching output
- **Gap:** Current coaching is directive. Could add MI framing: open questions, affirmations of existing good behavior, reflective restatements.

**Further Research Needed:**
- MI techniques adapted for text-based (async) coaching
- How to detect "sustain talk" from team behavior patterns
- OARS examples for common SM coaching scenarios

---

## 11. Crucial Conversations

**Source:** Patterson, Grenny, McMillan, Switzler (2002, 4th ed 2021)

**Summary:** Framework for high-stakes conversations where opinions vary and emotions run strong.

**7 Steps:**
1. Start with heart — clarify what you really want
2. Learn to look — detect when safety is at risk (silence or violence)
3. Make it safe — establish mutual purpose + mutual respect
4. Master my stories — separate facts from narratives
5. STATE my path — Share facts, Tell story, Ask their path, Talk tentatively, Encourage testing
6. Explore others' path — AMPP (Ask, Mirror, Paraphrase, Prime)
7. Move to action — Who does What by When, follow up How

**Key Insight:** "People who are skilled at dialogue do their best to make it safe for everyone to add their meaning to the shared pool." The bigger the shared pool of meaning, the better decisions.

**Mapping to zara-jira-mcp:**
- `pm_coaching(topic:"crucial_conversation")` — can generate prep
- `pm_hard_conversation` (planned) — data-backed conversation prep
- Data from blockers, health scores, anti-patterns = "facts" for step 4
- **Gap:** No explicit conversation prep template. Should generate: facts (from data), possible stories (interpretations), opening line options, SCARF risks to watch for.

**Further Research Needed:**
- STATE path examples for common PM situations (missed deadline, scope creep dispute, performance concern)
- How to detect "silence" and "violence" from async communication patterns

---

## 12. Thomas-Kilmann Conflict Modes

**Source:** Thomas & Kilmann (1974), TKI instrument (50+ years validated)

**Summary:** 5 conflict-handling modes based on 2 dimensions (assertiveness x cooperativeness):

```
High Assertive  │ COMPETING        │ COLLABORATING
                │ (win-lose)       │ (win-win)
                │                  │
                ├──── COMPROMISING ─┤ (split-the-difference)
                │                  │
Low Assertive   │ AVOIDING         │ ACCOMMODATING
                │ (lose-lose)      │ (lose-win)
                ├──────────────────┤
                  Low Cooperative    High Cooperative
```

**When each is appropriate:**
- Competing: emergency, unpopular but necessary decision
- Collaborating: both concerns too important to compromise
- Compromising: under time pressure, temporary settlement
- Avoiding: trivial issue, need cooling off
- Accommodating: you're wrong, or relationship preservation critical

**Key Insight:** Agile teams default to Avoiding (ignore conflict) or Accommodating (SM agrees with everyone). Healthy teams need Collaborating for important issues, Competing for safety/security decisions.

**Mapping to zara-jira-mcp:**
- `pm_coaching(topic:"conflict")` — could recommend mode based on situation
- `pm_conflict_mediation` (planned) — diagnose conflict mode + suggest resolution
- **Gap:** No conflict detection. Could infer from: blocked issues that ping-pong between assignees, decisions that get reversed, escalations that come back down.

**Further Research Needed:**
- TKI scoring from behavioral patterns (not questionnaire)
- Which mode maps to which PM situation (decision matrix)

---

## 13. Kotter + ADKAR Change Communication

**Source:** John Kotter (1996), Prosci/Jeff Hiatt (2003), Change Management Hub

**Summary:**

**Kotter 8-step (organizational change):**
1. Create urgency
2. Build guiding coalition
3. Form strategic vision
4. Enlist volunteer army
5. Enable action by removing barriers
6. Generate short-term wins
7. Sustain acceleration
8. Institute change

**ADKAR (individual change):**
- **A**wareness — why change is needed
- **D**esire — personal motivation to support
- **K**nowledge — how to change
- **A**bility — skills/behavior to implement
- **R**einforcement — sustain the change

**Integration:** Map Kotter steps → ADKAR elements so individual + organizational change reinforce each other.

**Key Insight for PM:** Every process change (new retro format, new DoD, new tool adoption) needs change communication. Most PMs announce change (Kotter step 1-3) but skip step 5-8 and never address individual ADKAR.

**Mapping to zara-jira-mcp:**
- `pm_experiment` — maps to "generate short-term wins"
- `pm_agreements` — maps to "institute change"
- `pm_change_communication` (planned) — full change comms plan
- **Gap:** No change communication template. No ADKAR tracking per change initiative.

**Further Research Needed:**
- Change communication template for PM context
- How to detect ADKAR blockers (person stuck at which stage?)
- Communication cadence per Kotter stage

---

## 14. Johari Window

**Source:** Luft & Ingham (1955), Harvard Business Review, Agile Laws

**Summary:** 4 quadrants of self-awareness:

```
              KNOWN TO SELF     UNKNOWN TO SELF
KNOWN TO    ┌─────────────────┬─────────────────┐
OTHERS      │   OPEN AREA     │   BLIND SPOT    │
            │  (shared)       │  (feedback)     │
            ├─────────────────┼─────────────────┤
UNKNOWN TO  │   HIDDEN AREA   │   UNKNOWN AREA  │
OTHERS      │  (self-disclose)│  (discovery)    │
            └─────────────────┴─────────────────┘
```

Goal: expand OPEN area through feedback + self-disclosure.

**Key Insight:** Only 10-15% of people are genuinely self-aware (HBR). In leadership roles, blind spots become expensive. SM's job includes expanding team's open area through structured feedback (retros, 1-on-1s).

**Mapping to zara-jira-mcp:**
- `pm_record_retro` — expands open area (team shares)
- `pm_anti_patterns` — reveals blind spots (data shows what team can't see)
- `pm_feedback_prep` (planned) — structured feedback using SBI
- **Gap:** No explicit blind spot tracking. Could track: what anti-patterns persist across sprints? (team sees data but doesn't change = organizational blind spot)

**Further Research Needed:**
- Johari-based team exercise templates
- How to detect organizational blind spots from persistent anti-patterns

---

## 15. Shannon-Weaver & Signal Retention

**Source:** Shannon & Weaver (1948/1949), hirefraction.com (2025), LaaS Litmus framework

**Summary:** Communication = Sender → Encoder → Channel → Noise → Decoder → Receiver. Key addition: Feedback loop (added later). Three problem levels: Technical, Semantic, Effectiveness.

**Signal Retention (modern interpretation):**
- Communication in software teams is "inherently lossy"
- Requirements degrade through handoffs
- Feedback gets diluted across layers
- Misaligned expectations → rework
- Signal retention = discipline of keeping critical info intact sender→receiver

**Key Insight:** Every handoff (PM→dev, dev→QA, team→stakeholder) is a potential signal loss point. AI can: reduce handoffs, preserve context, detect when signal has degraded.

**Mapping to zara-jira-mcp:**
- `pm_record_decision` — preserves signal (context + rationale don't degrade)
- `pm_team_kb` — reduces handoff loss (anyone can query context)
- `pm_search_decisions` — prevents re-decisions from signal loss
- **Gap:** No explicit "signal loss detection." Could track: stories returned to dev (requirement misunderstood = signal loss), decisions re-decided (context lost = signal loss).

**Further Research Needed:**
- Measuring signal retention: return-to-dev rate, re-decision rate
- Which channels lose most signal (Slack chat vs recorded decision vs story AC)?

---

## 16. Information Asymmetry / Principal-Agent

**Source:** Agency Theory (economics), arxiv 2601.23211 (multi-agent principal-agent), Wikipedia

**Summary:** Principal-agent problem: when one party (agent) acts on behalf of another (principal), information asymmetry creates misaligned incentives. Agent knows more about their actions than principal can observe.

**Application to PM:**
- PM = agent, Management = principal (PM knows sprint reality, management sees reports)
- Dev team = agent, PM = principal (team knows blocker severity, PM sees status)
- AI tool = agent, PM = principal (AI knows confidence level, PM sees recommendation)

**Key Insight:** PM's communication job is REDUCING information asymmetry without overwhelming the principal. The `pm_exec_report` vs `pm_dashboard` split IS addressing this — right level of info for right audience.

**Also relevant for AI agents:** arxiv paper argues multi-agent systems should be treated as principal-agent problems. Information asymmetry between agents not problematic when incentives aligned (agents report truthfully). Relevant to zara-jira-mcp's multi-agent dispatch.

**Mapping to zara-jira-mcp:**
- Entire reporting layer = reducing asymmetry
- `pm_forecast` with confidence intervals = honest uncertainty reporting
- `pm_escalate` = principal should know when agent is stuck
- **Gap:** No explicit "transparency score." Could measure: how much does dashboard match reality? (committed items vs delivered = honesty metric)

**Further Research Needed:**
- How AI can detect "hiding" behavior (issues moved to done without verification, inflated story points, scope quietly reduced)
- Incentive alignment in AI recommendation systems

---

## 17. Tuckman Stages + AI Assessment

**Source:** Bruce Tuckman (1965), Upscale Tech (2025 AI adaptation)

**Summary:**
1. **Forming** — politeness, independence, silos, anxiety
2. **Storming** — conflict, testing boundaries, power struggles
3. **Norming** — resolving differences, building cohesion, developing trust
4. **Performing** — high efficiency, autonomous, interdependent
5. **Adjourning** — disbanding, knowledge transfer

**AI-era twist (Upscale Tech 2025):** Teams adopting AI tools go through same stages:
- Forming: "What can AI do?"
- Storming: "AI is wrong / AI will replace us / AI slows us down"
- Norming: "We know when to use AI and when not to"
- Performing: "AI is seamless part of our workflow"

**Mapping to zara-jira-mcp:**
- `pm_maturity_assessment` — already maps to Tuckman stages based on data
- Communication needs differ per stage:
  - Forming → more structure, more check-ins
  - Storming → conflict facilitation, safe space
  - Norming → agreements, shared practices
  - Performing → delegate, get out of the way
- **Gap:** Communication recommendations not stage-aware. Could adjust all coaching/facilitation suggestions based on detected stage.

**Further Research Needed:**
- Behavioral markers per stage (measurable from Jira/GitHub data)
- Stage-specific communication frequency recommendations

---

## 18. Communication Metrics & Measurement

**Source:** Superhuman Blog (2025), Ragan Research (2024), Springer PCE scale

**Summary:** Communication effectiveness can be measured across 5 categories:
1. **Engagement** — message open/read rates, response rates
2. **Leadership** — alignment understanding, strategy clarity
3. **Collaboration** — handoff quality, feedback loops
4. **Productivity** — time spent gathering context, decision speed
5. **Sentiment** — satisfaction, trust, psychological safety

**Key metrics for PM context:**
- Decision speed: time from "need decision" → "decision recorded"
- Escalation speed: time from "blocker identified" → "escalated to right person"
- Re-decision rate: same topic decided multiple times = signal loss
- Stakeholder response time: how fast stakeholders respond to requests
- Meeting effectiveness: decisions per meeting hour (not meetings per week)
- Information freshness: are people using stale info? (checking old dashboards?)

**Mapping to zara-jira-mcp:**
- `pm_sm_impact` — measures resolution speed
- `pm_impediment_aging` — measures escalation speed
- `pm_meeting_roi` (existing?) — decisions per meeting
- `pm_stakeholder_pulse` — sentiment tracking
- **Gap:** No composite "communication health score." Could aggregate: decision speed + escalation speed + re-decision rate + stakeholder responsiveness into one score.

**Further Research Needed:**
- PCE (Perceived Communication Effectiveness) scale — validated instrument
- Benchmarks: what's "good" decision speed for software teams?
- How to measure without adding survey burden (infer from existing data)

---

## SYNTHESIS: Prioritized Implementation Impact

Based on all 18 research areas, ranked by: (1) measurability from existing data, (2) impact on PM communication, (3) implementation effort:

### Tier 1 — High Impact, Low Effort (do first)

| Research | Tool/Feature | Why |
|----------|-------------|-----|
| Signal Retention (#15) | `pm_comms_health` — re-decision rate, return-to-dev rate | Measurable from existing decision log + Jira data |
| Notification Fatigue (#6) | Notification budget tracking in `notify_routed` | Prevents tool abandonment |
| Crucial Conversations (#11) | `pm_hard_conversation` — data-backed prep | Uses existing data, high human value |
| Lencioni (#2) | Map anti-patterns → dysfunction levels | No new tools, enhance existing `pm_anti_patterns` |

### Tier 2 — High Impact, Medium Effort

| Research | Tool/Feature | Why |
|----------|-------------|-----|
| SPACE (#4) | `pm_space_metrics` — Communication dimension | Aggregates existing data into new view |
| EBM (#5) | `pm_ebm_dashboard` — 4 KVAs from existing data | Exec-speak, connects activity→value |
| Liberating Structures (#9) | Enhanced `pm_facilitate` with LS scripts | Better retros/planning facilitation |
| Nudge Theory (#8) | Audit + redesign recommendation outputs | Behavioral design of existing tools |

### Tier 3 — Medium Impact, Higher Effort (stretch)

| Research | Tool/Feature | Why |
|----------|-------------|-----|
| Psychological Safety (#1) | `pm_safety_check` + trend tracking | Requires survey integration or inference model |
| Westrum Culture (#3) | Culture indicator from behavioral patterns | Requires cross-referencing multiple signals |
| MI Coaching (#10) | Evocative coaching style in `pm_coaching` | Prompt engineering, not code |
| TKI Conflict (#12) | `pm_conflict_mediation` with mode recommendation | Needs situation classification |

---

## NEXT STEPS

1. Update `docs/communication-frameworks.md` with cross-references to this research
2. Create GitHub issues for Tier 1 items
3. Prototype `pm_comms_health` (aggregate existing data)
4. Enhance `pm_anti_patterns` with Lencioni dysfunction mapping
5. Add notification budget to `notify_routed` (counter + throttle)
6. Design `pm_hard_conversation` prompt (Crucial Conversations + SBI + data)
