# Research: Communication Frameworks for PM in the AI Era

> Compiled: 2026-06-26
> Focus: How PMs/Scrum Masters should communicate in a world where AI agents handle 60% of operational tasks.
> Goal: Make zara-jira-mcp the most communication-literate PM tool available.

---

## 1. FOUNDATIONAL COMMUNICATION MODELS

### The Pyramid Principle (Barbara Minto / McKinsey)

**Core rule:** Lead with the conclusion. Then support with grouped arguments. Then data.

```
BLUF (Bottom Line Up Front):
  → Conclusion / Recommendation
    → Key supporting argument 1
      → Data / Evidence
    → Key supporting argument 2
      → Data / Evidence
    → Key supporting argument 3
      → Data / Evidence
```

**Why this matters for PM tools:**
- Executives read the first 2 sentences. If your report buries the conclusion, it won't be read.
- AI-generated reports MUST follow this structure. `pm_exec_report` already does this.
- Every notification should answer: What happened? What should I do? How urgent?

[Source: managementconsulted.com/pyramid-principle, untools.co/minto-pyramid]

### BLUF (Bottom Line Up Front)

Military-origin framework adopted by business:
1. State the conclusion or request FIRST
2. Then provide context
3. Then supporting detail

**Application in zara-jira-mcp:**
- `pm_exec_report` → status line FIRST, then breakdown
- `pm_weekly_digest` → 3 bullet summary FIRST, then detail sections
- Notifications → action required FIRST, then context

[Source: en.wikipedia.org/wiki/BLUF_(communication)]

### Radical Candor (Kim Scott)

Two axes: **Care Personally** x **Challenge Directly**

| | Challenge Directly | Don't Challenge |
|---|---|---|
| **Care Personally** | Radical Candor (ideal) | Ruinous Empathy |
| **Don't Care** | Obnoxious Aggression | Manipulative Insincerity |

**The 4-step feedback order:** Get → Give → Gauge → Encourage

**Application for AI-powered PM tools:**
- Sprint health reports should be Radical Candor: honest data + empathy
- `pm_coaching` should challenge assumptions while acknowledging effort
- Anti-pattern detection should name the problem directly, not soften it into uselessness
- Reports to management should be candid about risks, not optimistic spin

[Source: radicalcandor.com/our-approach, em-tools.io/frameworks/radical-candor]

### Nonviolent Communication (NVC, Marshall Rosenberg)

Four components: **Observation → Feeling → Need → Request**

| Step | Example in PM context |
|------|----------------------|
| Observation | "Sprint velocity dropped 30% from Sprint 22 to 23" |
| Feeling | "The team seems stressed and disengaged" |
| Need | "We need sustainable pace to maintain quality" |
| Request | "Can we reduce commitment by 20% next sprint?" |

**Why NVC matters for PM AI tools:**
- Coaching outputs should never blame ("you failed to...")
- Instead: observation (data) → impact → need → actionable suggestion
- `pm_coaching` and `pm_anti_patterns` should follow this pattern
- Reports that trigger defensiveness are reports that get ignored

[Source: dave-bailey.com/blog/nonviolent-communication]

---

## 2. STRUCTURED COMMUNICATION FRAMEWORKS FOR PM

### RACI / DACI Matrix

**RACI:** Responsible, Accountable, Consulted, Informed
**DACI (Atlassian/tech variant):** Driver, Approver, Contributor, Informed

| Role | Who | Communication Need |
|------|-----|-------------------|
| Driver (D) | PM/SM | Full context, real-time updates |
| Approver (A) | PO/Stakeholder | Decision-relevant summary only |
| Contributor (C) | Dev team | Task-level detail + blockers |
| Informed (I) | Management | Outcomes + risks, no detail |

**Application:** Every `pm_record_decision` should capture DACI roles. Every notification should check: "Is this person R/A/C/I for this?" and route accordingly.

[Source: atlassian.com/templates/decision, ideaplan.io/glossary/raci-matrix]

### Situational Leadership (Hersey-Blanchard) Applied to Communication

| Team Readiness | Communication Style | Report Depth |
|---------------|--------------------|--------------| 
| R1 (Forming) | Directing | Detailed, prescriptive, frequent check-ins |
| R2 (Storming) | Coaching | Explain WHY, provide context, encourage questions |
| R3 (Norming) | Supporting | High-level, let team self-organize, available on request |
| R4 (Performing) | Delegating | Minimal, trust the team, exception-based only |

**Application:** `pm_maturity_assessment` output should influence notification frequency and report verbosity. A Performing team doesn't need daily nudges.

[Source: scribd.com/document/368863536/Situational-Leadership]

### Information Radiator Principles

**Core idea:** Make information visible without requiring people to ask for it.

Principles:
1. **Minimal interpretation** — the visual speaks for itself
2. **Shared ownership** — all teams contribute and rely on it
3. **Focus on flow** — highlight bottlenecks, not vanity metrics
4. **Always current** — stale data is worse than no data
5. **Ambient awareness** — glanceable, not demanding attention

**Application:** Sprint health dashboards, burndown, WIP limits should be always-available, not push notifications. Reserve push for exceptions only.

[Source: ituonline.com/tech-definitions/what-is-an-information-radiator, skills.visual-paradigm.com]

### Empiricism Pillars (Scrum)

**Transparency → Inspection → Adaptation**

- **Transparency**: Everyone sees the same data. No hidden information. Tools must surface truth, not comfortable summaries.
- **Inspection**: Regular checkpoints (standups, reviews, retros). AI should detect when inspection is being skipped.
- **Adaptation**: Change based on inspection. AI should detect when inspection happens but adaptation doesn't (zombie retros).

[Source: agileambition.com/empiricism-in-scrum, scrum.org]

---

## 3. COMMUNICATION IN THE AI ERA

### The Shift: PM as Communication Orchestrator (not Information Gatherer)

**Before AI (2020):**
- PM spends 5-10h/week collecting status from team members
- PM spends 3-5h/week formatting reports for different audiences
- PM spends 2-3h/week routing information between stakeholders
- Total communication overhead: 10-18h/week (25-45% of PM capacity)

**With AI-powered tools (2026):**
- AI collects and synthesizes status automatically from Jira/Git/Chat
- AI formats for audience (exec vs PO vs team) automatically
- AI routes via appropriate channels based on severity/audience
- PM role shifts to: **curating**, **interpreting**, **deciding**, **relating**

[Source: techademy.com/ai-stakeholder-communication, get-alfred.ai/blog/best-ai-tools-for-project-managers]

### Gartner Predictions (2025-2026)

- 40% of enterprise applications will feature task-specific AI agents by 2026
- 64% of communications leaders already use GenAI for executive messaging
- 68% of US IT firms adopted AI-enabled PM software (Gartner 2024)
- AI on its own does NOT solve the "unread" problem — but it enables experimentation with formats until one lands

[Source: gartner.com/en/communications/research, axis-intelligence.com/ai-stakeholder-management-communication-2026]

### Microsoft Work Trend Index 2025

- Employees interrupted 275 times/day during core work hours
- 3x more meetings/calls per week since 2020
- Information workers spend 20-30% of workweek searching for information
- AI meeting assistants now commodity (transcription accuracy commoditized by 2026)

[Source: questworks.games/blog/async-communication-guide-remote-teams, questionbase.com]

### 4 Patterns of Human-Agent Work (Microsoft WorkLab 2025)

| Pattern | What Changes | Human Responsibility |
|---------|-------------|---------------------|
| 1. AI as Tool | Agent executes discrete tasks | Human decides what to do |
| 2. AI as Assistant | Agent drafts, suggests, summarizes | Human reviews and approves |
| 3. AI as Collaborator | Agent takes initiative, proposes actions | Human sets direction and boundaries |
| 4. AI as Delegate | Agent owns outcomes within constraints | Human monitors exceptions |

**Each pattern only becomes possible because the previous one exists.**

For zara-jira-mcp: currently at Pattern 2-3. Moving toward Pattern 4 for routine communications (digests, escalations, status reports) while keeping Pattern 2 for decisions and coaching.

[Source: microsoft.com/en-us/worklab/ai-at-work-one-function-wrote-the-ai-playbook]

### Trust in Human-AI Communication

**Key research findings:**
- Only 29% of developers trust AI outputs to be accurate (Stack Overflow 2025)
- Information omissions by AI agents significantly reduce team trust
- Transparency and accountability are foundational for establishing trust
- "Meaningful human oversight" must preserve both AI agency and human agency

**Implications for zara-jira-mcp:**
- Always show data source and confidence level
- Never hide bad news in AI-generated reports
- Make AI reasoning visible ("Based on last 5 sprints...")
- Allow override/correction at every point
- Track accuracy over time (forecast calibration)

[Source: frontiersin.org/fpsyg.2025.1637339, arxiv.org/abs/2601.06223, augmentcode.com/guides/agent-handoff-patterns]

---

## 4. COGNITIVE LOAD & INFORMATION DESIGN

### Cognitive Load Theory (Sweller)

**Core insight:** Working memory can only process 2-4 items simultaneously.

| Load Type | Definition | PM Communication Impact |
|-----------|-----------|------------------------|
| Intrinsic | Complexity of the content itself | Sprint scope, technical decisions |
| Extraneous | Poorly designed presentation | Bad report format, unclear notifications |
| Germane | Productive learning/processing | Understanding patterns, making connections |

**Goal:** Minimize extraneous load. Keep notifications to ONE actionable insight. Structure reports for scanning (headers, bullets, bold key numbers).

[Source: get-alfred.ai/blog/cognitive-load-theory, frontiersin.org/forgp.2026.1812361]

### The "5-10 Notifications" Rule

- Max 5-10 notifications per active user per day across ALL channels
- Only 3% of alerts require immediate action
- DND elimination removes 92% of non-urgent interruptions
- Teams receive 2,000+ alerts/week, only 3% need action
- 57% of consumers actively avoid brands that flood with messages

**Application:** `notify_routed` should enforce daily caps. Batch aggressively. Send "nothing to report" = don't send.

[Source: suprsend.com, incident.io, knock.app]

### Async-First Communication

**Principle:** Default to async. Sync only when fidelity demands it.

| When Async | When Sync |
|-----------|-----------|
| Status updates | Conflict resolution |
| Decision documentation | Complex brainstorming |
| Progress reports | Emotional/sensitive topics |
| Code reviews | Real-time debugging |
| Sprint metrics | Retrospective discussions |

**Evidence:**
- Teams with async-first achieve 54% faster issue resolution
- Documented Slack etiquette → 22% faster sprint planning, 31% lower cognitive fatigue (N=312)
- Slack median first-response: 2-5 minutes vs email 30-90 minutes

**Application:** zara-jira-mcp should default all non-critical communications to async channels (digests, reports, summaries). Reserve sync triggers for P1 escalations only.

[Source: questworks.games, lifetips.alibaba.com, clearfeed.ai, moldstud.com]

---

## 5. AUDIENCE-ADAPTIVE COMMUNICATION

### The "Same Information, Different Format" Principle

One sprint produces many audiences. Each needs different framing:

| Audience | Needs | Format | Cadence | Tool |
|----------|-------|--------|---------|------|
| VP/C-Level | Business outcomes, risks, decisions needed | 3 sentences + 1 number | Weekly | `pm_exec_report` |
| Product Owner | Goal progress, scope changes, forecasts | Summary + detail drill-down | Daily-ish | `pm_goal_check` |
| Engineering Manager | Health, anti-patterns, resource issues | Dashboard + exceptions | Weekly | `pm_sprint_health` |
| Dev Team | Their tasks, blockers, PR reviews | Actionable list | Daily digest | `pm_standup_prep` |
| External Stakeholder | Deliverables, timelines, risks | Formal, jargon-free | Bi-weekly | Email via `email_send` |
| Steering Committee | Portfolio health, cross-team deps | Executive summary | Monthly | `portfolio_summary` |

### Language Calibration

| Audience | Avoid | Use Instead |
|----------|-------|-------------|
| Executive | "Story points", "velocity", "WIP" | "Items completed", "delivery speed", "parallel work" |
| Product Owner | "Refactoring", "tech debt" | "Stability work", "reducing future risk" |
| Developer | "Deliverables", "stakeholder alignment" | "What we're shipping", "what product wants" |
| External | Any internal jargon | Plain business language |

**AI must translate.** `pm_exec_report` should NEVER contain "story points" or "velocity". These are internal team metrics that mean nothing to leadership.

---

## 6. AI-SPECIFIC COMMUNICATION PATTERNS

### AI as Communication Mediator

The PM AI tool serves as a translation layer between:
- **Raw data** (Jira, Git, time tracking) → **Insight** (pattern, trend, anomaly)
- **Technical language** → **Business language** (audience-appropriate)
- **Many small signals** → **One actionable message** (cognitive load reduction)
- **Historical context** → **Prediction** (forecasting, pattern recognition)

### AI Communication Anti-Patterns

| Anti-Pattern | Why It Fails | Better |
|-------------|--------------|--------|
| AI dumps raw data | No insight, high cognitive load | Summarize + highlight exceptions |
| AI uses generic advice | Feels disconnected from reality | Ground in team's actual data |
| AI hides uncertainty | Erodes trust when predictions fail | Show confidence intervals |
| AI over-personalizes | Feels creepy or presumptuous | Be helpful, not intimate |
| AI generates without context | Hallucinated recommendations | Always cite the data source |
| AI replaces human judgment | PM loses authority and trust | AI proposes, human decides |

### The "Confidence + Source" Pattern

Every AI-generated communication should include:
```
[INSIGHT]: Sprint 24 is at risk of missing goal.
[CONFIDENCE]: High (based on 5 similar sprints, 4 missed with this burndown pattern)
[DATA]: 60% items still In Progress at mid-sprint. Historical average at this point: 35%.
[RECOMMENDATION]: Consider scope reduction. Remove 2-3 items to protect sprint goal.
[OVERRIDE]: If team has capacity not reflected in board, ignore this signal.
```

### Proactive vs Reactive Communication

| Type | When | Example |
|------|------|---------|
| Proactive | AI detects pattern before human asks | "Sprint 24 looks like Sprint 19 which failed" |
| Reactive | Human asks, AI responds | "What's sprint health?" → `pm_sprint_health` |
| Escalation | AI detects threshold breach | "Blocker > 3 days, escalating to manager" |
| Coaching | AI identifies learning opportunity | "This sprint had 40% carryover. Want to explore why?" |

**Rule:** Proactive communication must be HIGH signal. If proactive alerts have >50% false positive rate, they will be ignored. Better to miss one signal than cry wolf daily.

---

## 7. MEETING COMMUNICATION OPTIMIZATION

### Meeting ROI Framework

**The formula:** Meeting ROI = (Decisions Made + Actions Assigned + Alignment Created) / (Time Spent x Attendees)

**Indicators of low ROI:**
- Information-sharing meetings (could be async digest)
- Meetings without decisions or action items
- Status updates in sync format
- >8 attendees (beyond Scrum team size)

### AI-Augmented Meeting Flow

```
PRE-MEETING:
  → AI generates agenda from sprint data + open blockers + pending decisions
  → AI identifies: "These 3 items need sync discussion. These 5 can be resolved async."

DURING MEETING:
  → AI transcribes + extracts: decisions, action items, owners, deadlines
  → AI detects: unresolved items, missing assignments, time overruns

POST-MEETING:
  → AI generates: meeting notes → Confluence
  → AI extracts: action items → Jira tickets (auto-created)
  → AI distributes: summary to I (Informed) stakeholders
  → AI schedules: follow-up for unresolved items
```

**Application:** `pm_meeting_roi` already exists. Extend with pre/post meeting automation.

[Source: fellow.app, techademy.com/ai-meeting-summaries-for-pms, grain.com/blog/ai-meeting-summaries]

---

## 8. ESCALATION COMMUNICATION

### Progressive Escalation Model

```
Level 0: Self-service (information radiator, dashboards)
Level 1: Automated nudge (Slack DM to assignee)
Level 2: Team visibility (team channel mention)
Level 3: Manager alert (direct message + email)
Level 4: Executive escalation (formal report)
Level 5: Cross-org escalation (steering committee)
```

**Each level adds:**
- More people
- More formal channel
- More context
- Higher urgency language

**Rule:** Never jump levels. Escalation without prior levels = alarm fatigue.

### Escalation Communication Template

```
SUBJECT: [LEVEL X ESCALATION] [Brief description]

SITUATION: What is blocked and since when
IMPACT: What happens if unresolved (timeline, deliverable, revenue)
ATTEMPTED: What has already been tried
ASK: What specifically is needed from the recipient
DEADLINE: By when a response is needed
```

---

## 9. CROSS-CULTURAL COMMUNICATION CONSIDERATIONS

### High-Context vs Low-Context (Hall)

| Dimension | High-Context (Asian, LATAM) | Low-Context (US, Northern EU) |
|-----------|---------------------------|-------------------------------|
| Communication | Implicit, read between lines | Explicit, say exactly what you mean |
| Disagreement | Indirect ("that might be difficult") | Direct ("I disagree because...") |
| Decision style | Consensus-building | Individual authority |
| Notification tone | Softer, more context | Blunt, actionable |
| Meeting format | Relationship-first | Agenda-first |

**Application for zara-jira-mcp:**
- Coaching tone should be configurable (direct vs supportive)
- Escalation language should adapt to team culture
- Asian teams (Lark users) may prefer group notifications over individual DMs
- Reports for Japanese/Korean stakeholders should include more context

### Indonesian Communication Style (Project Context)

- Mix of formal (baku) and informal (non-baku) depending on audience
- Hierarchy-aware: different language for "atasan" vs "tim"
- Indirect refusal common: "nanti kita lihat" = probably no
- Value: gotong royong (mutual cooperation), musyawarah (consensus)
- PM tools in Indo context should avoid confrontational language in coaching

---

## 10. FRAMEWORKS SPECIFICALLY FOR AI-ERA PM

### The CLEAR Framework (AI-Adapted Communication)

**C** — Context: Ground every message in specific, verifiable data
**L** — Level: Match communication depth to audience readiness level
**E** — Evidence: Show source, confidence, and limitations
**A** — Action: Every communication ends with a clear next step
**R** — Review: Feedback loop to improve future communications

### The "3-Layer Report" Pattern

Every AI-generated report should have 3 accessible layers:

```
Layer 1: HEADLINE (1 line, BLUF, answer the question)
  "Sprint 24 is ON TRACK. 78% complete, 3 days remaining."

Layer 2: KEY INSIGHTS (3-5 bullets, exceptions only)
  - 2 items at risk: AUTH-45 blocked 2 days, FEAT-12 no reviewer
  - Velocity tracking 10% above last 5-sprint average
  - No critical risks. 1 medium risk (dependency on Team B)

Layer 3: FULL DETAIL (for those who want to drill down)
  - Complete item list with status
  - Burndown chart data
  - Individual contributor progress
  - Historical comparison
```

**Rule:** Layer 1 must be readable in 5 seconds. Layer 2 in 30 seconds. Layer 3 is optional.

### The Communication Debt Concept

Just as teams accumulate tech debt, they accumulate **communication debt:**

| Communication Debt | Symptom | Cost |
|-------------------|---------|------|
| Undocumented decisions | "Why did we do this?" meetings | Repeated discussions |
| No stakeholder pulse | Surprise negative feedback | Relationship damage |
| Skipped retros | Same mistakes repeated | Efficiency loss |
| No risk log | Surprised by known unknowns | Late escalation |
| No knowledge base | 40% repeat questions | 5-15h/week wasted |

**Application:** `pm_record_decision`, `pm_stakeholder_pulse`, `pm_record_risk` are communication debt reduction tools. Frame them this way to PM users.

---

## 11. PRACTICAL RECOMMENDATIONS FOR ZARA-JIRA-MCP

### What Makes This Tool Communication-Literate

1. **BLUF-first in all outputs** — Every tool output leads with the answer, not the journey
2. **Audience-adaptive language** — `pm_exec_report` vs `pm_standup_prep` use different vocabulary
3. **Cognitive load respect** — Max 5-7 items in any list. Batch aggressively. "Nothing to report" = silence
4. **Confidence signaling** — AI recommendations show basis and certainty level
5. **Escalation intelligence** — Progressive levels, never skip, auto-detect threshold breach
6. **Cultural adaptability** — Tone/language adapts to team culture settings
7. **Async-by-default** — Push only for exceptions. Digest for everything else
8. **Feedback loops** — Track: generated → delivered → read → acted-on
9. **Communication debt tracking** — Surface when decisions/risks/retros are skipped
10. **Human override everywhere** — AI proposes, human disposes. Never autonomous for irreversible actions

### New Tool Opportunities (Based on Research)

| Tool | What It Does | Framework Applied |
|------|-------------|-------------------|
| `pm_communication_health` | Score team's communication patterns (decision doc rate, meeting ROI, escalation speed) | SPACE + EBM |
| `pm_audience_translate` | Take any report and reformat for different audience level | Pyramid Principle + Situational Leadership |
| `pm_notification_audit` | Analyze notification volume/response-rate per team member, detect fatigue | Cognitive Load Theory |
| `pm_async_standup` | Collect individual updates at their timezone, AI-cross-reference with board, post summary | Async-First + NVC |
| `pm_decision_debt` | List undocumented decisions, orphaned ADRs, stale knowledge | Communication Debt |
| `pm_confidence_calibration` | Track AI prediction accuracy over time, show trust score | Trust Research |
| `pm_escalation_history` | Show escalation patterns: which blockers escalated, resolution time, was it effective | Progressive Escalation |

### Communication Principles to Embed in Every Tool

```
1. LEAD WITH THE ANSWER (BLUF)
2. SHOW YOUR WORK (Confidence + Source)
3. RESPECT ATTENTION (Cognitive Load < 7 items)
4. ADAPT TO AUDIENCE (Language Calibration)
5. DEFAULT TO ASYNC (Push only for exceptions)
6. ESCALATE PROGRESSIVELY (Never skip levels)
7. ENABLE OVERRIDE (Human always has final say)
8. TRACK EFFECTIVENESS (Feedback loop on delivery)
9. BE HONEST (Radical Candor > Ruinous Empathy)
10. REDUCE DEBT (Document decisions, risks, outcomes)
```

---

## 12. SOURCES & CITATIONS

### Communication Frameworks
1. Barbara Minto — The Pyramid Principle — managementconsulted.com/pyramid-principle
2. Kim Scott — Radical Candor — radicalcandor.com/our-approach
3. Marshall Rosenberg — Nonviolent Communication — dave-bailey.com/blog/nonviolent-communication
4. BLUF — Bottom Line Up Front — en.wikipedia.org/wiki/BLUF_(communication)
5. DACI Framework — atlassian.com/templates/decision
6. Information Radiators — ituonline.com/tech-definitions/what-is-an-information-radiator

### AI & Communication Research
7. Microsoft Work Trend Index 2025 — questworks.games/blog/async-communication-guide-remote-teams
8. Microsoft WorkLab — 4 Patterns of Human-Agent Work — microsoft.com/en-us/worklab
9. Gartner — 40% enterprise apps with AI agents by 2026 — axis-intelligence.com/ai-stakeholder-management-communication-2026
10. Gartner — 64% comms leaders use GenAI — gartner.com/en/communications/research
11. Frontiers in Psychology — Trust in Human-AI Communication — frontiersin.org/fpsyg.2025.1637339
12. arxiv.org/abs/2601.06223 — Transparency, Accountability, Trustworthiness in AI
13. Stack Overflow 2025 — 29% trust AI accuracy — augmentcode.com/guides/agent-handoff-patterns
14. Microsoft Research — Agentic AI: Reimagining Human-Agent Communication — microsoft.com/en-us/research

### Cognitive Load & Notification Science
15. John Sweller — Cognitive Load Theory — get-alfred.ai/blog/cognitive-load-theory
16. Frontiers in Org Psychology 2026 — Cognitive Overload as Ergonomic Risk — frontiersin.org/forgp.2026.1812361
17. incident.io — 2000+ alerts/week, 3% need action
18. suprsend.com — 5-10 notifications/day max
19. knock.app — 57% consumers avoid message-flooding brands
20. Slack Engineering — Notification rebuild research

### Async & Team Communication
21. Ben Balter (GitHub) — Tools of the Trade — ben.balter.com/2020/08/14/tools-of-the-trade/
22. N=312 engineers study — Slack etiquette → 22% faster planning — lifetips.alibaba.com
23. clearfeed.ai — Slack median response 2-5min vs email 30-90min
24. moldstud.com — 54% faster issue resolution with dedicated chat
25. questionbase.com — 20-30% of week searching for information, 40% questions are repeats

### PM & Agile Specific
26. Scrum.org — Evidence-Based Management (EBM)
27. Google Project Aristotle — Psychological safety explains 43% of performance variance
28. DORA/SPACE Framework — Microsoft/GitHub 2021
29. PMI 2024 — Ethics of Over-Allocation in Sprints
30. Scrum Alliance — Only 52% teams achieve sprint goals
31. Cortex 2024 — State of Developer Productivity

### Meeting & Stakeholder
32. fellow.app — AI meeting summary tools 2026
33. techademy.com — PMs attend 25h/week meetings
34. grain.com — AI meeting summaries: decisions, owners, due dates
35. projectmanagement.com — Using AI for team communication without losing trust

---

## 13. ADDITIONAL RESEARCH SUMMARIES (To Explore Further)

Below are additional communication models, theories, and research areas that are relevant but need deeper investigation to determine specific applicability to zara-jira-mcp tooling.

### Classic Communication Models

| Model | Core Idea | PM Relevance | Riset Lanjutan |
|-------|-----------|-------------|----------------|
| **Lasswell's 5W** (1948) | Who says what, in which channel, to whom, with what effect | Maps directly to notification routing: sender, message, channel, audience, intended action | Explore how to make every tool output answer all 5W |
| **Shannon-Weaver** (1948) | Sender → Encoder → Channel → Noise → Decoder → Receiver | "Noise" = notification fatigue, information overload, context loss. Encoding = format for audience | Identify all noise sources in PM communication pipeline |
| **Schramm's Interactive Model** (1954) | Communication is circular, both parties encode/decode, shared "field of experience" required | Team needs shared mental model. Reports fail when sender/receiver have different context | How to ensure AI-generated comms match receiver's field of experience |
| **Osgood-Schramm** (1954) | Continuous feedback loop, no fixed sender/receiver | Every communication creates response. AI must listen for feedback on its own outputs | Build feedback mechanisms into every AI-generated report |

**Key insight:** Semua model klasik mengarah ke hal yang sama — komunikasi bukan one-way broadcast. AI tool yang cuma generate report tanpa feedback loop = half the model.

### Behavioral & Psychological Models

| Model | Core Idea | PM Relevance | Riset Lanjutan |
|-------|-----------|-------------|----------------|
| **Media Richness Theory** (Daft & Lengel 1986) | Rich media (video/face-to-face) for ambiguous tasks, lean media (text) for routine | Match channel to message complexity. Sprint planning = rich. Status update = lean | Map each PM ceremony to ideal media richness level |
| **Nudge Theory** (Thaler & Sunstein) | Small environmental changes steer behavior without restricting choice | Notifications as "nudges" — default actions, smart ordering, framing | Design notifications as choice architecture |
| **Cognitive Load Theory** (Sweller) | Working memory handles 2-4 items max. Reduce extraneous load | Every notification/report competes for limited attention. Less = more | Audit all tool outputs for cognitive load score |
| **Communication Accommodation Theory** | People adjust speech to reduce/increase social distance | AI should adapt tone to team maturity, culture, individual preferences | Build tone configuration into coaching tools |
| **Context Switching Cost** (23min recovery per interruption) | Every notification costs 23+ minutes of deep work if it interrupts flow | Push notifications are expensive. Batch unless P1 | Calculate "interruption cost" of notification strategy |

**Key insight:** Setiap notifikasi punya "hidden cost" 23 menit. Kalau 10 notif/hari = 3.8 jam deep work hilang. Tool harus default ke batching.

### Organizational & Leadership Communication

| Model | Core Idea | PM Relevance | Riset Lanjutan |
|-------|-----------|-------------|----------------|
| **Mendelow's Power/Interest Grid** | Map stakeholders by power x interest → different engagement strategy per quadrant | High-power/high-interest = manage closely. Low/low = monitor only. Route notifications accordingly | Auto-map stakeholders from Jira roles → notification profile |
| **PMBOK Communication Plan** | Audience x Message x Channel x Frequency x Owner matrix | Standard PM artifact. AI should auto-generate and maintain this | `pm_communication_plan` tool that generates from team structure |
| **Psychological Safety** (Edmondson/Google) | 43% of team performance variance explained by feeling safe to speak up | AI coaching tone must never shame. Reports should normalize failure as learning | Audit coaching outputs for psychological safety compliance |
| **Transparent Communication** (TAEO model) | Transparency + Authenticity + Empathy + Optimism → trust → engagement | AI reports must be transparent about bad news. Never spin. Show empathy for impact | Ensure `pm_exec_report` and `pm_coaching` follow TAEO |
| **Servant Leadership Communication** | Lead by removing barriers, not commanding. Listen first, speak second | PM tool should ask "what do you need?" not "here's what you should do" | Redesign coaching prompts to servant-leader style |

**Key insight:** AI PM tool harus jadi servant leader — enable team, bukan micromanage. Proactive detection, tapi reactive prescription (offer help, don't force).

### AI-Era Specific Research

| Topic | Key Finding | PM Relevance | Riset Lanjutan |
|-------|-------------|-------------|----------------|
| **Human-in-the-Loop (HITL)** | "HITL alone is not a governance strategy" (IBM). Must be meaningful oversight, not rubber stamp | Every AI escalation/recommendation must allow easy override with clear "why" | Design override UX for every proactive AI action |
| **AI Trust Gap** | 29% trust AI accuracy. Trust drops when AI omits information | Always show what data AI used and what it couldn't access | Add "data sources" section to all AI-generated outputs |
| **Beyond Psychological Safety** (AI Journal 2025) | AI tools reduce moments where assumptions are tested. Learning becomes private | PM tool should create spaces for collective learning, not just individual dashboards | Design retro/coaching tools that prompt team discussion |
| **Narrative Intelligence / Data Storytelling** | Stories with data boost retention 65% vs raw data presentation | AI reports should follow narrative arc: situation → complication → resolution | Rewrite all report templates with narrative structure |
| **Closed-Loop Communication** (Aviation CRM) | Send → Confirm → Verify. If no confirmation, assume not received | Every escalation should track: sent → acknowledged → acted-on. Escalate on silence | Add "acknowledgment required" flag to critical notifications |
| **Ben Balter: AI-First PM** (2026) | "AI-augmented PM is natural evolution of async-first. Amplifies judgment, doesn't replace it" | Zara-jira-mcp positioned perfectly here. AI + async + human judgment | Use as positioning statement for project |
| **AI SM Augmentation** (Multiple 2025-26) | SM administrative load reduced 30-50% by AI. Freed time goes to coaching/relationships | Our value prop: not "replace SM" but "give SM 10-15h/week back for human work" | Frame all marketing/docs around time-saved for human work |
| **Modular Prompting** (PMI 2026) | Reusable prompt structures improve decision-making and team alignment | Internal: how we structure AI prompts for report generation | Audit and standardize all LLM prompts in codebase |

**Key insight:** AI PM tool yang sukses = yang membuat human PM LEBIH human, bukan less human. Free dari admin → more coaching, more empathy, more relationship.

### Communication Anti-Patterns in AI-Era PM (Need Research)

| Anti-Pattern | Hypothesis | What to Validate |
|-------------|-----------|-----------------|
| "AI said so" authority | Teams stop questioning AI recommendations | Does our coaching output get accepted without challenge? |
| Report blindness | Too many AI-generated reports → nobody reads any | Track read rates and action rates on generated reports |
| Learned helplessness | PM delegates all thinking to AI, loses situational awareness | Does PM engagement decrease over time with tool usage? |
| Echo chamber | AI trained on team's own data reinforces existing biases | Are anti-patterns detected that the team is blind to? |
| Trust erosion via hallucination | One wrong forecast = team ignores all future AI signals | Track forecast accuracy, show calibration scores |
| Notification arms race | Multiple AI tools competing for attention | Monitor total notification volume across all sources |

---

## 14. RESEARCH GAPS & OPEN QUESTIONS

Things we don't yet know and should investigate:

1. **What's the optimal notification:action ratio?** (Current guess: 1 action per 3 notifications max)
2. **Does AI-generated coaching change behavior?** (Need longitudinal study on teams using `pm_coaching`)
3. **Cross-cultural notification preferences?** (Do Asian teams want more or fewer notifications than Western teams?)
4. **Ideal AI report length by audience?** (Executive: <100 words? PO: <300 words? Team: <500 words?)
5. **When does proactive AI become annoying?** (What's the threshold between helpful and nagging?)
6. **Does data storytelling in AI reports improve decision quality?** (A/B test narrative vs bullet format)
7. **HITL friction: how much override is too much?** (If PM overrides >50% of AI suggestions, is the AI useful?)
8. **PM de-skilling risk:** Does reliance on AI PM tools atrophy PM judgment over time?
9. **Feedback loop design:** What's the minimum viable feedback mechanism for AI-generated comms?
10. **Multi-language communication:** How should AI handle mixed-language teams (e.g., Indo/English)?

---

## 15. ADDITIONAL SOURCES (Batch 2)

36. Lasswell's 5W Model (1948) — en.wikipedia.org/wiki/Lasswell's_model_of_communication
37. Shannon-Weaver Model — en.wikipedia.org/wiki/Shannon-Weaver_model
38. Schramm Interactive Model — en.wikipedia.org/wiki/Schramm's_model_of_communication
39. Media Richness Theory — thecommspot.com/communication-theories/media-richness-theory
40. Nudge Theory / Choice Architecture — pnas.org/doi/10.1073/pnas.2107346118 (meta-analysis)
41. Communication Accommodation Theory — emergentmind.com/topics/communication-accommodation-theory-cat
42. Mendelow Power/Interest Grid — mutomorro.com/tools/mendelow-power-interest-matrix
43. PMBOK Communication Plan — riskpublishing.com/project-communication-plan-example
44. Psychological Safety — universidadisep.com/en/2026/06/psychological-safety-high-performing-teams-en
45. Beyond Psychological Safety in AI era — aijourn.com/beyond-psychological-safety
46. Transparent Communication & Trust — instituteforpr.org/how-does-leadership-communication-impact-employee-trust-during-crisis
47. Context Switching Cost 223% slower PRs — usehaystack.io/blog/the-true-cost-of-context-switching
48. Context Switching 40% productivity loss — APA via taskade.com/wiki/productivity/deep-work
49. Data Storytelling 65% retention boost — moldstud.com/articles/p-mastering-data-storytelling
50. Closed-Loop Communication (Aviation CRM) — ncbi.nlm.nih.gov/sites/books/NBK551708
51. Servant Leadership in Agile — toptal.com/project-managers/agile/agile-servant-leadership
52. Human-in-the-Loop not governance — ibm.com/think/insights/liability-laundering-problem
53. AI-First Program Management — ben.balter.com/2026/05/31/ai-first-program-management
54. AI SM reduces admin 30-50% — techademy.com/best-ai-tools-for-scrum-masters
55. AI won't replace good SM — scrum.org/resources/blog/why-ai-wont-replace-good-scrum-master
56. Modular Prompting for PM — pmi.org/blog/modular-prompting-practical-guide
57. Prompt Engineering for PMs — ideaplan.io/guides/prompt-engineering-for-pms
58. JTBD for Notifications — mrx.sivoinsights.com/blog/jtbd-for-mobile-ux
59. Forbes: Why Open Door Policies Fail — forbes.com/sites/benjaminlaker/2026
60. Deloitte: Transparency in Workplace — deloitte.com/us/en/insights/topics/talent/human-capital-trends/2024/transparency-in-the-workplace
