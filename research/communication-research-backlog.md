# Communication Research Backlog & Roadmap

> Compiled: 2026-06-26
> Purpose: Summary of all relevant communication research areas for zara-jira-mcp, with actionable next steps.

---

## A. FRAMEWORKS ALREADY DOCUMENTED (in communication-frameworks.md)

| # | Framework | Status | Relevance |
|---|-----------|--------|-----------|
| 1 | Minto Pyramid Principle | Documented | `pm_exec_report`, `pm_communicate` |
| 2 | RACI Matrix | Documented | `pm_raci` |
| 3 | SBI/SBII Feedback Model | Documented | `pm_feedback_prep`, `pm_coaching` |
| 4 | Radical Candor | Documented | `pm_hard_conversation` |
| 5 | SCARF Model | Documented | `pm_scarf_check` |
| 6 | NVC (Nonviolent Communication) | Documented | `pm_conflict_mediation` |
| 7 | BLUF (Bottom Line Up Front) | Documented | All notification outputs |
| 8 | SCQA | Documented | Escalation templates |
| 9 | Pyramid of Communication Needs | Documented | L1-L4 automation strategy |
| 10 | 4 C's Framework | Documented | Human/AI boundary design |
| 11 | ADR/RFC | Documented | `pm_decision_record` |
| 12 | Stakeholder Sentiment | Documented | `pm_silence_detector` |

---

## B. FRAMEWORKS TO RESEARCH DEEPER (Summaries + Next Steps)

### B1. ADKAR Change Management Model (Prosci)
**Summary**: Individual change framework: Awareness → Desire → Knowledge → Ability → Reinforcement. Each stage needs specific communication.
- Awareness: "Why change?" — urgency communication
- Desire: "What's in it for me?" — motivation/vision
- Knowledge: "How do I change?" — training/docs
- Ability: "Can I do it?" — practice/support
- Reinforcement: "Will it stick?" — celebration/accountability

**Relevance**: When PM introduces new processes, tools, or team changes. AI could detect which ADKAR stage a team member is stuck at based on their behavior patterns (resistance = Desire issue, mistakes = Knowledge/Ability issue).

**Tool opportunity**: `pm_change_readiness` — assess team ADKAR state for a proposed change; `pm_change_communication` — generate stage-appropriate comms.

**Next step**: Research ADKAR + agile transformation communication patterns. How to detect ADKAR stage from Jira data (e.g., low adoption of new workflow = stuck at Desire).

---

### B2. Kotter's 8-Step Change Model
**Summary**: Leadership-driven organizational change in 8 phases:
1. Create urgency
2. Build guiding coalition
3. Form strategic vision
4. Communicate the vision (enlist support)
5. Remove barriers (empower action)
6. Generate short-term wins
7. Sustain acceleration
8. Anchor in culture

**Relevance**: Large-scale agile transformations, introducing AI tools to teams, shifting from Scrum to Kanban, etc. Integration with Agile and Lean methodologies shown to improve outcomes.

**Tool opportunity**: `pm_change_plan` — generate Kotter-aligned communication plan for any team change (sprint length change, new tool adoption, team restructuring).

**Next step**: Map Kotter steps to PM ceremonies. Which sprint events naturally support which step? (Retro = Step 7-8, Sprint Review = Step 6).

---

### B3. Thomas-Kilmann Conflict Instrument (TKI)
**Summary**: 5 conflict-handling modes based on 2 dimensions (assertiveness × cooperativeness):
- **Competing** (high assertive, low cooperative) — win/lose
- **Collaborating** (high/high) — win/win, but time-intensive
- **Compromising** (moderate/moderate) — split the difference
- **Avoiding** (low/low) — withdraw, delay
- **Accommodating** (low assertive, high cooperative) — yield

**Relevance**: Retro discussions, inter-team conflicts, stakeholder disagreements. AI could diagnose which mode is being used and suggest alternatives. Most team conflicts default to Avoiding — SM's job is to surface them.

**Tool opportunity**: `pm_conflict_diagnosis` — analyze a described conflict and suggest handling mode + communication approach; integrate with `pm_coaching`.

**Next step**: Research TKI × agile team dynamics. Which modes correlate with healthy teams? (Collaborating for technical decisions, Compromising for scope negotiations). How to detect conflict avoidance from retro patterns.

---

### B4. Crucial Conversations / Difficult Conversations
**Summary**: Framework for high-stakes discussions where emotions run strong and opinions differ:
- Start with heart (clarify what you want)
- Learn to look (notice when safety is at risk)
- Make it safe (restore mutual purpose/respect)
- Master your stories (separate facts from interpretations)
- STATE your path (Share facts, Tell your story, Ask for their path, Talk tentatively, Encourage testing)

Key stat: 60% of employees encounter "unwritten rules" causing miscommunication (Workhuman 2024). Difficult conversations are a SKILL, not personality trait — improve through deliberate practice.

**Relevance**: 1-on-1s, performance conversations, escalations. AI can prep the PM: gather facts, separate from stories, suggest safe opening.

**Tool opportunity**: `pm_hard_conversation` (already planned) — should use STATE path structure. Add prep template: Facts gathered → My story → Their possible story → Mutual purpose → Opening statement.

**Next step**: Research STATE path integration with sprint data. Can we auto-populate the "facts" section from Jira metrics?

---

### B5. ORID Focused Conversation Method
**Summary**: Developed by Institute of Cultural Affairs. Structured dialogue in 4 levels:
- **O**bjective — What happened? (facts, data, observations)
- **R**eflective — How do we feel about it? (reactions, emotions)
- **I**nterpretive — What does it mean? (insights, patterns, significance)
- **D**ecisional — What do we do? (actions, next steps)

**Relevance**: PERFECT for retrospectives. Most retros jump from O straight to D (what happened → what to do). Missing R and I leads to superficial actions. AI could facilitate ORID-structured retros.

**Tool opportunity**: `pm_retro_facilitate` — guide discussion through ORID phases; generate prompts for each stage using sprint data.

**Next step**: Research ORID × sprint retrospective effectiveness. Design retro templates that force I stage (interpretation) before D (decisions).

---

### B6. Liberating Structures
**Summary**: 33 simple, powerful interaction patterns that include every voice. Developed over decades, tested globally. Key structures for PM:
- **1-2-4-All** — individual → pairs → groups → whole group (prevents groupthink)
- **Troika Consulting** — 3 people, rotating roles of client/consultants
- **15% Solutions** — what can you do RIGHT NOW without permission?
- **TRIZ** — what could make this absolutely fail? (pre-mortem variant)
- **Wicked Questions** — surface paradoxes the team is navigating

**Relevance**: Retro facilitation, sprint planning engagement, breaking meeting fatigue. AI could suggest which structure to use based on team situation.

**Tool opportunity**: `pm_facilitation_suggest` — given context (retro, planning, conflict, brainstorm), recommend a Liberating Structure with instructions.

**Next step**: Map Liberating Structures to PM contexts. Which structure for which ceremony?

---

### B7. Lean Coffee
**Summary**: Democratic meeting format (2009, Jim Benson):
1. Generate topics (everyone writes)
2. Vote (dot voting)
3. Discuss in order of votes (timeboxed per topic)
4. Move topics: To Discuss → Discussing → Discussed

**Relevance**: Prevents agenda hijacking. Perfect for retros, community of practice meetings, ad-hoc problem-solving. AI could auto-generate topic suggestions from sprint data.

**Tool opportunity**: `pm_lean_coffee_prep` — suggest topics based on sprint events, vote results, past unresolved items.

**Next step**: Research how to auto-generate Lean Coffee topics from Jira data (stale items, blockers, velocity anomalies).

---

### B8. Nudge Theory (Thaler & Sunstein)
**Summary**: Subtle changes to choice architecture that influence behavior without restricting autonomy:
- Default settings (opt-in vs opt-out)
- Reminders at decision points
- Social proof ("80% of teams do X")
- Simplification of desired behavior
- Feedback loops (immediate consequences visible)

Workplace applications: reminders, default nudges, implementation intentions, priming.

**Relevance**: PM tools shouldn't DEMAND behavior — they should NUDGE. Example: instead of "you should do standup", surface "3 items haven't been updated in 2 days" at the right moment. Choice architecture for Sprint planning (default capacity at 80%, not 100%).

**Tool opportunity**: Every tool output is a nudge, not a command. `pm_auto_detect_risks` is already a nudge. Expand: default suggestions that are easy to accept (e.g., "Shall I create a blocker ticket for this?").

**Next step**: Audit all tool outputs for nudge vs command framing. Rewrite any that demand action into choice-preserving nudges.

---

### B9. Cognitive Load Theory (Sweller) + Notification Design
**Summary**: Working memory has strict capacity limits. Overloading reduces comprehension, increases errors.
- Fragmented workflows + persistent interruptions = sustained cognitive strain
- 79% of workers get distracted within 1 hour; 59% can't maintain 30 min focus
- Structured roundups reduce task-switching latency 47%, daily interruptions from 22 to 3.8
- Progressive disclosure: reveal info incrementally
- Chunking: group related data into digestible sections

**Relevance**: Tool outputs must respect cognitive limits. A notification with 20 data points is worse than one with 3 actionable items. Profile system (14 tools for ChatGPT vs 124 for experts) already addresses this. Go further: output length/density should adapt to context.

**Tool opportunity**: `PM_OUTPUT_DENSITY` config — brief/normal/detailed. All tools respect this. Notifications: max 3 items, expandable. Dashboard: progressive disclosure (headline → detail on request).

**Next step**: Research optimal notification density. How many items per alert? How many alerts per day? Design attention budget system.

---

### B10. Proactive AI Communication (arxiv 2025-2026)
**Summary**: Cutting-edge research on AI agents that anticipate needs before explicit prompts:
- Proactive agents reduce cognitive load and streamline workflows
- Communication policy evolution — learning WHEN to communicate, not just what
- Context-aware perceptions enable better user intent understanding
- Proactive communication policies learned via reward models trained on human acceptance/rejection
- Key challenge: proactiveness must be USEFUL, not annoying (notification fatigue)

**Relevance**: Our tools should be PROACTIVE, not just reactive. `pm_auto_detect_risks` is already proactive. Expand: predict blockers before they happen, surface stale items before PM asks, suggest stakeholder updates before they're overdue.

**Tool opportunity**: `pm_proactive_scan` (scheduled) — daily scan that generates a prioritized "what you should know/do today" without being asked. Like a smart daily briefing.

**Next step**: Research proactive communication policies. When is proactive helpful vs annoying? Design threshold system (only surface if confidence > X%, impact > Y).

---

### B11. AI Burnout Detection via Communication Patterns
**Summary**: AI can detect burnout BEFORE employees recognize symptoms:
- Monitor: calendar density, meeting fragmentation, work-hour patterns
- Slack/communication sentiment shifts weeks before self-recognition
- Context switching frequency as burnout proxy
- Not about surveillance — about early intervention
- Forbes 2025: AI flags patterns from calendar data + work fragmentation

**Relevance**: Our existing tools (WIP per person, velocity decline, carryover) are already indirect burnout signals. Add communication pattern analysis: response time changes, message brevity shifts, silence periods.

**Tool opportunity**: `pm_wellbeing_signals` — aggregate available signals (from Jira: WIP, cycle time, overtime patterns; from chat: response patterns) into wellbeing risk score. Private to PM/SM only.

**Next step**: Research ethical boundaries of wellbeing monitoring. What's helpful vs invasive? Design opt-in consent model. Map Jira-only signals (no chat monitoring) that still provide value.

---

### B12. AI Coaching Directiveness Research (Frontiers 2026)
**Summary**: Surprising finding — directive AI coaching was rated HIGHER than non-directive for:
- Technology performance expectancy
- Working alliance quality
- Goal attainment

People with high extraversion, conscientiousness, and openness preferred directive approach. Contradicts traditional coaching dogma (always ask, never tell).

**Relevance**: Our `pm_coaching` tool could adapt style. New SM (D1-D2) → directive ("Here's what to do"). Experienced SM (D3-D4) → non-directive ("What do you think about..."). Match coaching style to user's developmental level.

**Tool opportunity**: `pm_coaching` style parameter: directive/supportive/delegating (matches situational leadership D1-D4).

**Next step**: Research adaptive coaching style selection. Can we infer user's developmental stage from their query patterns?

---

### B13. Meeting Statistics & Meeting Replacement
**Summary** (2026 data):
- 21.7 meetings/week per knowledge worker (avg)
- 15.4 hours/week in meetings vs 12.1 hours uninterrupted focus
- 71% of meetings considered unproductive (HBR)
- 76% of decisions forgotten within 24h if no follow-up notes
- Meeting time increased 252% since 2020 (Microsoft)
- $400B/year wasted on unnecessary meetings (US)
- 47% of workers feel drained on heavy meeting days
- 40% feel dread when meetings appear on calendar
- 16.85% of productive capacity consumed by meetings

**Relevance**: MASSIVE opportunity. Every meeting our tools can replace = direct productivity gain. `pm_standup_prep` replaces standup asking. `pm_daily_delta` replaces "what's the status" sync. `pm_async_update` replaces weekly status meeting.

**Tool opportunity**: `pm_meeting_cost` — calculate actual cost of a recurring meeting ($salary × attendees × duration × frequency). Show what async alternative saves. Nudge: "This standup costs $2,400/month. Here's the async alternative."

**Next step**: Quantify meeting replacement ROI for each tool. Create comparison: "pm_standup_prep (5 min AI-generated) vs 15-min standup with 8 people = $X saved/sprint".

---

### B14. Stakeholder Power/Interest Grid (Mendelow Matrix)
**Summary**: 2×2 matrix plotting stakeholders by:
- Power/Influence (ability to affect project) — vertical
- Interest/Involvement (how much they care) — horizontal

Four quadrants → four strategies:
- High Power + High Interest → **Manage Closely** (frequent, detailed comms)
- High Power + Low Interest → **Keep Satisfied** (periodic, summary comms)
- Low Power + High Interest → **Keep Informed** (regular, transparent comms)
- Low Power + Low Interest → **Monitor** (minimal comms)

**Relevance**: `pm_audience_router` should USE this grid to decide communication frequency and detail level. Auto-detect quadrant from stakeholder interaction patterns.

**Tool opportunity**: `pm_stakeholder_map` — build and maintain power/interest grid; auto-suggest communication strategy per stakeholder; detect quadrant shifts.

**Next step**: Research how to auto-detect power/interest from Jira/Lark/Slack data (frequency of comments, approval authority, escalation patterns).

---

### B15. Communication Cadence Design
**Summary**: Research-backed cadence patterns:
- **Daily**: Operational metrics, blockers, real-time unblocking (PM/SM only)
- **Weekly**: Wins, concerns, decisions, next week (stakeholders)
- **Bi-weekly**: Sprint health, goal progress (team)
- **Monthly**: Executive summary, business outcomes (leadership)
- **Quarterly**: Strategic review, portfolio health (C-suite)

Key insight: "A fixed cadence creates natural touchpoints for planning, asset requests, and campaign changes" — reduces anxiety, builds trust. Single weekly digest can prevent frantic Slack messages and cut status calls in half.

**Relevance**: Our notification system should support cadence-based delivery, not just on-demand queries. PM configures cadence once, tools deliver automatically.

**Tool opportunity**: `pm_cadence_config` — set up automated delivery schedule; `pm_digest` — generate cadence-appropriate report on schedule. Integration with Lark/Slack/Email scheduled messages.

**Next step**: Design cadence configuration schema. Which tools map to which cadence? (daily: `pm_daily_delta`; weekly: `pm_weekly_digest`; end-of-sprint: `pm_release_notes`).

---

### B16. Sprint Review as Communication Event (Not Demo)
**Summary**: Sprint Review is NOT a demo. Key insights:
- It's inspection + adaptation, not presentation
- Two-way conversation, not one-directional show
- Start with sprint goal — did we deliver? If not, why?
- Most reviews are "status theater" — no useful feedback, misalignment compounds
- 76% of decisions forgotten in 24h without notes

**Relevance**: `pm_release_notes` + `pm_goal_check` already support better reviews. Add: pre-review prep that frames discussion around goals, not ticket lists. Post-review action capture.

**Tool opportunity**: `pm_sprint_review_prep` — generate review agenda focused on goal achievement, stakeholder questions, and feedback prompts (not a ticket list).

**Next step**: Research effective sprint review formats. What makes stakeholders ACTUALLY give feedback? Design facilitation guide.

---

## C. ADDITIONAL RESEARCH AREAS (Not Yet Explored)

| # | Topic | Why It Matters | Priority |
|---|-------|---------------|----------|
| C1 | Motivational Interviewing (MI) | Coaching resistant team members | Medium |
| C2 | Appreciative Inquiry (AI/4D) | Strength-based retros, positive framing | Medium |
| C3 | DISC/MBTI Communication Styles | Adapting communication per personality | Low |
| C4 | Information Radiators (Agile) | Dashboard design, glanceability | High |
| C5 | Shannon-Weaver Communication Model | Channel noise, signal degradation in async | Low |
| C6 | Johari Window | Team trust building, feedback culture | Medium |
| C7 | Active Listening in Digital Context | How to "listen" in async/written comms | Medium |
| C8 | Storytelling for Data (Cole Nussbaumer) | Making metrics compelling to executives | High |
| C9 | The Elephant and the Rider (Heath) | Emotional + rational persuasion for change | Medium |
| C10 | Working Out Loud (WOL) | Transparency culture, showing work in progress | Medium |
| C11 | OKR Communication | Connecting sprint work to company objectives | High |
| C12 | Psychological Ownership Theory | Why people resist change to "their" process | Medium |
| C13 | Feedback Loops in Systems Thinking | Reinforcing/balancing loops in team dynamics | Medium |
| C14 | Trust Equation (Maister) | Credibility × Reliability × Intimacy / Self-orientation | High |
| C15 | Five Dysfunctions of a Team (Lencioni) | Absence of trust → fear of conflict → lack of commitment → avoidance of accountability → inattention to results | High |

---

## D. IMPLEMENTATION ROADMAP (Prioritized)

### Week 1: Foundation — Communication Framework Engine

- [ ] Design `CommunicationStyle` struct: audience, framework, density, language
- [ ] Implement Pyramid Principle output formatter (conclusion-first for all tools)
- [ ] Add BLUF template to all notification outputs
- [ ] Design `pm_communicate` tool — audience-aware message generation
- [ ] Audit existing tool outputs for nudge-vs-command framing

### Week 2: Feedback & Coaching Tools

- [ ] Implement `pm_feedback_prep` — SBI-structured feedback from sprint data
- [ ] Enhance `pm_coaching` — add directiveness parameter (D1-D4)
- [ ] Implement `pm_hard_conversation` — STATE path + data gathering
- [ ] Add SCARF-check to all change/announcement notifications

### Week 3: Proactive & Async Communication

- [ ] Design `pm_daily_delta` improvements — cognitive load limits (max 5 items)
- [ ] Implement `pm_async_update` — generate async status replacing meetings
- [ ] Design cadence configuration schema (daily/weekly/sprint/monthly)
- [ ] Implement `pm_silence_detector` — flag disengaged stakeholders

### Week 4: Stakeholder Intelligence

- [ ] Implement `pm_stakeholder_map` — power/interest grid
- [ ] Implement `pm_audience_router` — reframe same data per audience
- [ ] Implement `pm_communication_plan` — who/what/when/how for any event
- [ ] Add meeting cost calculator to sprint planning prep

### Week 5: Facilitation & Retro Enhancement

- [ ] Implement `pm_retro_facilitate` — ORID-structured retro guide
- [ ] Add Liberating Structures suggestions to ceremony prep tools
- [ ] Implement Lean Coffee topic generation from sprint data
- [ ] Enhance `pm_conflict_mediation` — add TKI diagnosis + NVC framing

### Week 6: Wellbeing & Trust

- [ ] Implement `pm_wellbeing_signals` — aggregate burnout risk indicators
- [ ] Implement `pm_trust_signals` — track psychological safety indicators
- [ ] Design ethical boundaries: what's helpful vs invasive
- [ ] Add ADKAR stage detection for change initiatives

### Stretch / Future

- [ ] `pm_meeting_cost` — ROI of meeting replacement
- [ ] `pm_communication_debt` — audit undocumented decisions
- [ ] `pm_change_readiness` — ADKAR + Kotter assessment
- [ ] `pm_proactive_scan` — daily AI briefing (scheduled)
- [ ] Sentiment analysis integration (opt-in, from Lark/Slack)
- [ ] Multi-language output (Indonesian/English toggle)

---

## E. KEY STATISTICS (Quick Reference)

| Stat | Source |
|------|--------|
| 21.7 meetings/week per worker | Microsoft Work Trend Index 2026 |
| 15.4 hrs/week in meetings | Microsoft 2026 |
| 71% meetings unproductive | HBR |
| 76% decisions forgotten in 24h (no notes) | Laxis 2026 |
| 252% meeting time increase since 2020 | Microsoft |
| $400B/year wasted on meetings (US) | Flowtrace/Atlassian |
| $1.2T/year miscommunication cost | Talaera |
| $3.2M/year async savings (60-person team) | JetThoughts |
| 89% trust drop in agentic AI | Axis Intelligence 2025 |
| 31% trust drop in company GenAI | Axis Intelligence 2025 |
| 13% orgs see real performance from AI | Glean 2026 |
| 43% performance variance = psychological safety | Google Aristotle |
| 23 min context switch recovery | Multiple sources |
| 79% workers distracted within 1 hour | Makerstations 2026 |
| 47% reduction in task-switching with structured roundups | Carnegie Mellon |
| 60% encounter "unwritten rules" causing miscommunication | Workhuman 2024 |
| AI coaching market: $6.25B | Osmo/PRNewswire 2026 |
| AI PM market: $3.08B → $7.4B (2024-2029) | Medium/Kanerika |

---

## F. RESEARCH SOURCES TO DEEP-DIVE

1. Frontiers in Psychology 2026 — AI coaching directiveness study
2. Frontiers in Psychology 2025 — Trust in human-AI team communication
3. Frontiers in Organizational Psychology 2026 — Cognitive overload as ergonomic risk
4. arxiv 2505.14668 — Context-aware proactive LLM agents
5. arxiv 2602.04482 — Evaluating LLM agents for proactive assistance
6. arxiv 2606.14314 — Communication policy evolution for proactive agents
7. Springer 2024 — Proactive AI implications for competence & satisfaction
8. ResearchGate — ADKAR comprehensive guide
9. ResearchGate — Thomas-Kilmann conflict modes & leadership styles
10. PMI 2026 Report — soft skills > tools for high-performing orgs
11. Laxis State of Meetings 2026 — benchmark report
12. Prosci — ADKAR vs Kotter comparison
13. LeadDev — Async-first blueprint
14. Forbes 2025 — AI burnout detection from calendar/work patterns
15. Scrum.org — Liberating Structures for Scrum events
