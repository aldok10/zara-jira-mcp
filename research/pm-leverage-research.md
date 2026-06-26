# PM/SM Leverage Research: What Creates 10x Impact with AI Tooling

Compiled: 2026-06-26 | Sources: 15+ web searches, 508 academic papers (existing), DORA research, industry reports

---

## 1. What Makes the Top 1% Scrum Masters Different

### The Invisible Activities That Create 10x Outcomes

**Elite SMs operate as system sensors, not ceremony facilitators.**

| Average SM | Top 1% SM |
|------------|-----------|
| Facilitates meetings | Reads the system between meetings |
| Waits for problems to emerge | Detects signals 3-5 days before failure |
| Removes blockers when reported | Preempts blockers before team notices |
| Runs retrospectives | Ensures retro actions actually complete (only 35% of teams do this) |
| Tracks velocity | Tracks velocity *trends* and *anomalies* |
| Updates boards | Uses boards as information radiators that tell stories |
| Asks "any blockers?" | Notices who stopped asking questions |

**Key behavioral differences (research-backed):**

1. **Pattern Recognition Over Process** — They correlate sprint velocity drops with PR review delays, team composition changes, or external dependency patterns. They don't just see numbers; they see *why*.

2. **Predictive Scrum Mastery** — By correlating historical team velocity with current complexity of open PRs, they predict sprint failure by Day 3. This gives time to negotiate scope *before* burnout. [Source: AgileVisa "The Predictive Scrum Master"]

3. **Psychological Safety as Primary Metric** — Google's Project Aristotle: #1 factor separating great teams was psychological safety, not talent or resources. Elite SMs measure this through behavioral proxies, not surveys.

4. **"Silence Reading"** — They detect who has gone quiet. Pull requests getting shorter/sloppier, once-proactive members stopping questions, defensive comment tones — all burnout precursors visible in work artifacts 2-3 weeks before resignation.

5. **Coaching Stance** — Ask questions before giving answers. The research shows agile coaches affect effort, strategies, knowledge, and skills of teams. Most essential traits: emphatic, people-oriented, able to listen, diplomatic, persistent.

6. **Transfer Leadership** — Mature SM grows themselves out of the job. They transfer 9 leadership roles to the team over time. D4 teams need delegation, not direction.

### How They Read Team Dynamics Without Surveys

Observable behavioral proxies for psychological safety:
- **Voice frequency** in ceremonies (who speaks, who doesn't)
- **Defensive voice patterns** (qualifying every statement, hedging)
- **Silence behaviors** (particular individuals going quiet)
- **Supportive vs unsupportive responses** in discussions
- **Learning-oriented behaviors** (admitting mistakes, asking for help)
- **Patterns of constant agreement** = dissent feels risky

**From work artifacts:**
- PR comment tone shifts (collaborative → transactional → defensive)
- Standup update length decay (rich → perfunctory → "same as yesterday")
- Question frequency drop (engaged → executing tickets → disappeared)
- Review turnaround time divergence (engaged members respond faster)

---

## 2. PM Pain Points That AI Can Actually Solve (2025-2026)

### The "60% Tax" — Work About Work

Research: PMs spend up to 60% of their time on "work about work" — status chasing, writing updates, switching between tools. One PM reported AI cut monthly management hours by 70%. [Sources: Zemith 2026, Medium/Karen 2024]

### Specific Cognitive Drains

| Pain Point | Time Cost | AI Solvability |
|------------|-----------|----------------|
| Status report generation | 3-5h/week | HIGH — summarize from Jira data |
| Stakeholder update emails | 2-3h/week | HIGH — generate from sprint state |
| Chasing updates from team | 2-4h/week | HIGH — detect stale tickets, auto-nudge |
| Sprint planning prep | 2-3h/sprint | HIGH — capacity + carryover + risk calc |
| Retro facilitation prep | 1-2h/sprint | HIGH — surface data-driven talking points |
| Blocker escalation tracking | 1-2h/week | HIGH — auto-age, auto-escalate |
| Risk register maintenance | 1h/week | MEDIUM — auto-detect from signals |
| Cross-team dependency tracking | 2-3h/week | MEDIUM — detect from Jira links |
| Board hygiene | 1-2h/week | HIGH — stale item detection |
| Meeting notes → actions | 1h/meeting | HIGH — record + track completion |

### Decisions PMs Make Daily That Should Be Data-Informed

1. **"Is this sprint on track?"** — Usually gut feel. Should be: burndown rate vs historical, WIP analysis, blocker age, scope change delta.

2. **"Who should I check on?"** — Usually whoever's loudest. Should be: who has unusual PR patterns, high WIP, stale items, or went quiet.

3. **"Should we cut scope?"** — Usually last day panic. Should be: Monte Carlo simulation showing 50%/85% completion probability by Day 3.

4. **"Is this team healthy?"** — Usually "no one complained." Should be: cycle time trends, review turnaround, WIP distribution, engagement proxies.

5. **"Which risk needs attention NOW?"** — Usually recency bias. Should be: risk severity × age × blast radius ranking.

### Where PMs Fail Due to Information Asymmetry

- **Developers know** code is fragile but PM doesn't see tech debt accumulating
- **One person** knows they're blocked but doesn't escalate for 3 days
- **External team** deprioritized your dependency but no one told you
- **Sprint scope** grew by 15% through "quick tweaks" that nobody tracked
- **Team member** is burning out (visible in activity data) but says "I'm fine"
- **Recurring retro themes** show the same problem for 4 sprints but nobody connects the dots

---

## 3. Leverage Multipliers

### If a PM Could Only Have 5 AI-Powered Actions

**Ranked by impact-per-effort:**

#### 1. Sprint Failure Early Warning (Day 3 Alert)
- Correlate: current burndown rate + open PRs complexity + blocker age + WIP distribution
- Output: "Sprint is at 35% completion probability. Recommend: cut [X], [Y], or negotiate timeline."
- **Impact**: Prevents the #1 PM failure — delivering the bad news on sprint review day

#### 2. Team Health Anomaly Detection (Continuous)
- Monitor: PR turnaround, comment tone, standup engagement, WIP per person, blocker self-reporting delay
- Alert: "Dev X hasn't submitted a PR in 5 days (usually 1.5/day). Check in."
- **Impact**: Catches burnout/struggle/blocker 1-2 weeks before it manifests as missed commitments

#### 3. Scope Creep Guardian (Continuous)
- Baseline: Sprint start snapshot (items + story points)
- Monitor: Any additions post-planning without removals
- Alert: "3 items added mid-sprint (+18 SP). No items removed. Scope grew 22%."
- **Impact**: Makes invisible scope creep visible immediately — the "dozen quick tweaks" pattern

#### 4. Retro Action Item Enforcer (Weekly)
- Track: Every retro action item with owner + deadline
- Persist: In next sprint planning prep, standup prep, next retro
- Alert: "4 of 6 retro actions from last 2 sprints are still incomplete."
- **Impact**: Only 35% of teams complete retro actions. This tool alone differentiates from 65%

#### 5. Right-Context-Right-Time Surfacing (Event-Triggered)
- Before standup: "Dev A has 4 items in progress (WIP limit: 2). Dev B finished all items."
- Before planning: "Last sprint velocity was 34 SP. Team capacity is 80% (2 on PTO). Historical similar sprints completed 28 SP."
- Before 1:1: "This team member had 3 blockers last sprint, carryover ratio 60%, PR review time 3x team average."
- **Impact**: Eliminates "I didn't know" information asymmetry

### What "Predictive Project Management" Looks Like in Practice

**The shift: From "what happened" → "what's about to happen" → "what should we do"**

| Level | Traditional | Predictive |
|-------|-------------|------------|
| Reporting | "We completed 28 SP last sprint" | "Based on 10,000 Monte Carlo simulations, we have 85% probability of completing 26-32 SP" |
| Risk | "We have a risk register" | "Risk X has been open 14 days, correlates with Blocker Y, and historically this pattern delays delivery by 5 days" |
| People | "Everyone seems fine" | "Dev X's PR volume dropped 60% this week, review response time doubled — early burnout signal" |
| Scope | "We planned 35 SP" | "Scope has grown 18% since planning. At current velocity, we'll complete 85% of original scope or 70% of current scope" |
| Dependencies | "Team B said they'd deliver" | "Team B's dependency item hasn't moved in 4 days. Their sprint has 12 items in progress (WIP limit: 6). HIGH risk of delay" |

---

## 4. Team Health Indicators That Predict Problems

### Leading Indicators (Early Warning — Action Window: Days)

| Indicator | Signal | Tool Implementation |
|-----------|--------|---------------------|
| WIP spike per person | >2 items in progress | Real-time WIP monitoring per assignee |
| PR review latency increase | Time-to-first-review growing | Track review request → first comment delta |
| Standup engagement decay | Updates getting shorter/same | NLP on standup notes (if recorded) |
| Blocker self-report delay | Time between block and report >24h | Compare activity gap vs blocker creation |
| Scope additions without removals | Items added post-sprint-start | Snapshot comparison (already built) |
| Velocity trend break | >15% deviation from 5-sprint rolling avg | Statistical anomaly detection |
| Single-person concentration | >40% of sprint work on one person | Workload distribution monitoring |
| Backlog volatility | High add/remove/reorder rate | Track backlog churn rate |
| After-hours commits | Increasing frequency of late commits | Git timestamp analysis |
| PR size inflation | PRs getting larger over time | Track lines changed per PR |

### Lagging Indicators (Too Late — Damage Done)

| Indicator | What It Means |
|-----------|---------------|
| Sprint commitment miss (>25%) | Team was already overloaded |
| Resignation | Burnout was visible 4-6 weeks earlier |
| Quality cliff (bugs spike) | Corners were cut weeks ago |
| Customer escalation | Internal signals were ignored |
| Velocity collapse | Accumulated debt finally broke something |

### Code Review Patterns That Signal Team Dysfunction

From research on 8M+ PRs:
- **Elite teams**: Full PR cycle <26 hours
- **Lagging teams**: PRs sit for a week before first look
- **Median**: 2-5 days cycle time

**Dysfunction signals:**
1. **Review queue growing** — Nobody's reviewing, everyone's producing → bottleneck ahead
2. **Rubber-stamp reviews** (approve in <5 min on 500+ line PRs) → quality decay
3. **One reviewer doing all reviews** → hero culture, bus factor = 1
4. **Nitpick-only reviews** (formatting but no logic discussion) → disengagement
5. **Review comment tone shift** (helpful → terse → adversarial) → interpersonal friction
6. **PR size inflation** (small focused PRs → large risky PRs) → deadline pressure, cutting corners
7. **Abandoned PRs** (opened but never merged/closed) → scope confusion or ownership vacuum

### Communication Pattern Changes That Predict Blockers

Research shows consistent communication patterns predict success:
- **Volume spike** in one person → overloaded, seeking help
- **Volume drop** in previously active person → disengagement or burnout
- **Frequency decrease** between specific people → broken collaboration
- **After-hours message increase** → unsustainable pace
- **Question frequency drop** → either mastery (good) or learned helplessness (bad, check other signals)

### Velocity Anomalies That Indicate Hidden Problems

| Pattern | Likely Cause |
|---------|--------------|
| Steadily declining velocity | Unidentified process problem, tech debt accumulation |
| Velocity spike then crash | Team was "borrowing" from quality, debt caught up |
| High volatility (>30% swing) | Inconsistent estimation, scope churn, or team instability |
| Velocity stable but customer satisfaction dropping | Velocity measuring wrong things, feature factory pattern |
| Velocity increasing but cycle time increasing too | Items getting smaller (split gaming) without real throughput gain |
| Individual velocity stable, team velocity dropping | Collaboration breaking down, silo formation |

---

## 5. The "PM as Coach" Model (Research-Backed)

### How Coaching Interventions Improve Sprint Outcomes

From research (Papers 17-23 in existing research + new findings):
- Agile coaches **affect effort, strategies, knowledge, and skills** of teams
- Coaching is distinct from directing — it's about **unlocking team's own capacity**
- Teams with dedicated coaching show **stronger interpersonal trust (r=0.93)** and **organizational performance (r=0.90)**
- The coaching stance: "Questions are the primary tool for facilitation"

### Powerful Questions Framework (Timing-Based)

**Before Sprint Planning:**
- "What was our biggest learning from last sprint that should change how we plan this one?"
- "What dependency makes you most nervous about this sprint?"
- "If we could only deliver ONE thing this sprint, what would make the biggest difference?"

**During Sprint (Day 3-5 Check):**
- "What's your confidence level that we'll hit the sprint goal? (1-5)"
- "What's the one thing that could derail us that we haven't talked about?"
- "Who on the team could use help right now?"

**Before Retrospective:**
- "What pattern have you noticed repeating across sprints?"
- "What did we say we'd do last retro that we actually did / didn't do?"
- "If a new team member joined, what would surprise them about how we work?"

**1:1 Coaching:**
- "What's taking more energy than it should?"
- "What challenge do you see in collaborating with the team?" (not "Are you having difficulties?")
- "What would make next sprint different from this one for you?"

### Detecting When a Team Member Is Struggling (Before It's a Blocker)

**Data-observable signals (no surveillance needed — just Jira/Git artifacts):**

1. **PR pattern change** — Shorter, sloppier, defensive tone, or complete stop
2. **WIP accumulation** — Items started but not finished, growing over time
3. **Blocker reporting delay** — Longer time between "stuck" and "asked for help"
4. **Estimation drift** — Same-sized stories taking 2-3x longer than historical
5. **After-hours commit increase** — Compensating for lost daytime productivity
6. **Review participation drop** — Stopped reviewing others' code (withdrawal)
7. **Carryover rate increase** — Personal sprint completion rate declining

**The key insight:** These signals are visible in existing data. No surveillance tools needed. Just correlation and trend detection on Jira + Git data that already exists.

### Psychological Safety Measurement Proxies From Work Artifacts

Based on research (observational measure with 31 behaviors in 7 categories):

| Category | Positive Proxy | Negative Proxy |
|----------|---------------|----------------|
| Voice | Questions in reviews, dissenting comments, suggestions | Silence in discussions, only agreements |
| Learning | "I made a mistake" in retros, asking for help openly | Never admitting uncertainty |
| Support | Collaborative PR comments, pair programming offers | "Not my problem" patterns |
| Familiarity | Cross-pollination (reviewing outside your area) | Strict ownership silos |
| Safety | Experiments proposed, new approaches tried | Only proven patterns, risk avoidance |

---

## 6. Practical MCP Tool Ideas (Highest Impact) — NEW

Based on all research above, mapped against existing 139 tools in zara-jira-mcp:

### Tier 1: PREVENT Problems (Highest Leverage)

#### Tool: `pm_individual_health_signal`
**What it does:** Per-team-member health scoring from Jira + Git data.
- Track: PR frequency trend, WIP count, carryover rate, blocker delay, review participation
- Score: 1-5 health per person per sprint
- Alert: When any individual's score drops >1.5 points from their baseline
**Why:** Catches burnout/struggle 2-3 weeks before it becomes a missed commitment or resignation.
**Differs from existing:** `pm_team_health` is aggregate. This is individual signal detection.

#### Tool: `pm_sprint_day3_forecast`
**What it does:** Automatic sprint health check on Day 3 (of 10-day sprint).
- Runs Monte Carlo with current burndown rate
- Factors in: open PR complexity, WIP distribution, blocker count
- Output: "Sprint completion probability: 72%. Risk factors: [X]"
- Suggests: scope cut options ranked by business value / effort
**Why:** Day 3 is the intervention window. By Day 7, it's too late to recover gracefully.
**Differs from existing:** `pm_forecast` exists but isn't triggered at the RIGHT MOMENT with decision options.

#### Tool: `pm_dependency_tracker_live`
**What it does:** Monitors cross-team dependency health in real-time.
- Watches linked issues in other projects
- Detects: stale dependencies (no movement in N days)
- Detects: dependency team's sprint overload (their WIP vs capacity)
- Alert: "Dependency PROJ-456 hasn't moved in 4 days. Owner team has 18 items WIP (limit: 10). HIGH risk."
**Why:** Cross-team dependencies are where PMs consistently fail due to information asymmetry.
**Differs from existing:** `pm_dependencies` records them. This one monitors and alerts on health.

#### Tool: `pm_question_coach`
**What it does:** Context-aware coaching question generator.
- Input: Context (ceremony type, team state, individual situation)
- Uses: Current sprint data + historical patterns + coaching frameworks
- Output: 3-5 specific questions tailored to the situation
- Example: Before 1:1 with struggling dev → "Dev X has 60% carryover rate and 3x WIP. Consider asking: 'What's taking more energy than it should this sprint?'"
**Why:** The right question at the right time is the SM's highest-leverage tool. Most SMs default to the same 5 questions.
**Differs from existing:** `pm_coaching` gives general advice. This generates specific questions with data context.

### Tier 2: SURFACE Right Context at Right Moment

#### Tool: `pm_ceremony_brief`
**What it does:** Auto-generated briefing document before each ceremony.
- **Before Standup:** WIP violations, stale items, blocker ages, who finished what since last standup
- **Before Planning:** Velocity trend, capacity calc, carryover items, dependency status, retro actions pending
- **Before Retro:** Sprint data summary, recurring themes from last 3 retros, action item completion rate, anti-pattern signals
- **Before Review:** Demo order, what shipped vs didn't, stakeholder-facing summary, risks to mention
**Why:** Eliminates 2-3h prep per ceremony while making every ceremony data-driven.
**Differs from existing:** Parts exist (`pm_standup_prep`, `pm_planning_prep`, `pm_review_prep`) but not unified with the "right context at right moment" timing engine.

#### Tool: `pm_notification_budget`
**What it does:** Intelligent notification throttling and prioritization.
- Tracks: How many notifications sent per day/week
- Scores: Each potential alert by: severity × novelty × action-ability
- Throttles: Low-value alerts get batched into daily digest
- Ensures: Critical alerts never buried under noise
- Respects: Notification fatigue research (>10/week → user mutes)
**Why:** Research shows proactive AI that sends 10+ notifications/week gets muted. The notification budget is the architectural blind spot of proactive-AI.
**Differs from existing:** `notify_routed` routes but doesn't budget or throttle.

#### Tool: `pm_pattern_connector`
**What it does:** Cross-signal correlation engine.
- Connects: Velocity drop + PR review delay + specific person overloaded → root cause story
- Connects: Recurring retro themes + unchanged metrics → "dead retro" pattern
- Connects: Blocker age + dependency team WIP + escalation delay → predicted delivery miss
- Output: Narrative explanation, not just numbers
**Why:** Individual metrics are noise. Correlated metrics are signal. This is what elite SMs do in their heads.
**Differs from existing:** `pm_anti_patterns` detects known patterns. This discovers NEW correlations.

### Tier 3: REDUCE Cognitive Load While Maintaining Awareness

#### Tool: `pm_decision_recommender`
**What it does:** When decisions need to be made, presents options with data backing.
- Trigger: Detected situation requiring decision (scope cut, escalation, resource realloc)
- Output: 2-3 options with: projected impact, historical precedent, risks of each, recommended action
- Example: "Sprint at 40% probability. Options: (A) Cut items X,Y — recover to 75% probability. (B) Negotiate 2-day extension. (C) Move Z to next sprint — lowest value item per stakeholder ranking."
**Why:** Reduces decision fatigue by doing the analysis, while leaving the decision to the human.
**Differs from existing:** Current tools present data. This presents decisions.

#### Tool: `pm_retro_effectiveness_score`
**What it does:** Measures whether retrospectives are creating actual improvement.
- Tracks: Action item completion rate across sprints
- Tracks: Whether the SAME themes recur sprint after sprint
- Tracks: Whether metrics changed after retro actions were implemented
- Scores: Retro effectiveness 0-100
- Alerts: "Same issue raised in 4 consecutive retros. Action items completed: 1/7. Retro is dead."
**Why:** Dead retros (going through motions without change) are the #1 anti-pattern of average teams. Only 35% of teams complete retro actions.
**Differs from existing:** `pm_retro_analysis` analyzes content. This measures EFFECTIVENESS (did anything change?).

#### Tool: `pm_wip_guardian`
**What it does:** Real-time WIP limit enforcement with smart nudging.
- Monitors: Per-person WIP continuously
- Detects: WIP violations before they become cycle time problems
- Suggests: "Dev X has 4 items in progress. Optimal: 2. Suggest: pair with Dev Y on item Z to finish before starting new work."
- Tracks: WIP compliance over time as team health metric
**Why:** WIP limits reduce cycle time by 30-40% (research). Teams with numerical WIP boundaries have 30% fewer interruptions and measurably better throughput.
**Differs from existing:** `pm_flow_metrics` calculates WIP. This enforces and coaches on it.

#### Tool: `pm_stakeholder_update_generator`
**What it does:** Auto-generates stakeholder-appropriate updates at different levels.
- Input: Current sprint state + audience level (team, manager, VP, customer)
- **For VP (30 seconds):** "Sprint on track. 2 features shipping Friday. 1 risk: dependency on Team B."
- **For manager (2 min):** Sprint progress, blockers, risks, team health, what's shipping.
- **For team (detail):** Full sprint data, individual progress, code review queue, dependency status.
**Why:** PMs spend 2-3h/week writing the SAME information at different abstraction levels.
**Differs from existing:** `pm_exec_report` exists but is one-shot. This is continuous, multi-audience.

### Tier 4: COACH Without Being Patronizing

#### Tool: `pm_team_growth_tracker`
**What it does:** Tracks team maturity signals over time.
- Measures: Self-organization level (how many decisions team makes without SM)
- Measures: Cross-pollination (are people reviewing outside their area?)
- Measures: Ownership breadth (bus factor per code area)
- Measures: Quality of self-reported impediments (do they surface problems early?)
- Suggests: When to step back (team is ready for more autonomy)
**Why:** The goal of an SM is to transfer leadership. This measures progress toward that goal.
**Differs from existing:** No current tool measures maturity progression.

#### Tool: `pm_experiment_evaluator`
**What it does:** After an improvement experiment runs, evaluates if it worked.
- Takes: Experiment hypothesis + measurement criteria
- Pulls: Actual metrics from before/after period
- Determines: Did the metric actually improve? Statistically significant?
- Output: "Experiment 'limit WIP to 2' ran for 2 sprints. Cycle time: 4.2 days → 2.8 days (-33%). Verdict: EFFECTIVE. Recommend: make permanent."
**Why:** Most experiments die because nobody measures the outcome rigorously.
**Differs from existing:** `pm_experiments` records them. This evaluates them.

---

## Summary: The PM Leverage Stack

```
┌─────────────────────────────────────────────────────┐
│  TIER 4: COACH (Long-term team growth)              │
│  • Team growth tracker • Experiment evaluator       │
├─────────────────────────────────────────────────────┤
│  TIER 3: REDUCE LOAD (Free PM cognitive space)      │
│  • Decision recommender • Retro effectiveness       │
│  • WIP guardian • Stakeholder update generator      │
├─────────────────────────────────────────────────────┤
│  TIER 2: SURFACE (Right context, right moment)      │
│  • Ceremony brief • Notification budget             │
│  • Pattern connector                                │
├─────────────────────────────────────────────────────┤
│  TIER 1: PREVENT (Highest leverage)                 │
│  • Individual health signal • Day 3 forecast        │
│  • Live dependency tracker • Question coach         │
├─────────────────────────────────────────────────────┤
│  EXISTING: 139 tools (foundation layer)             │
│  • Jira CRUD • Sprint management • Memory           │
│  • Flow metrics • Anti-patterns • Notifications     │
└─────────────────────────────────────────────────────┘
```

### Implementation Priority (by ROI)

| # | Tool | Effort | Impact | Build On |
|---|------|--------|--------|----------|
| 1 | pm_sprint_day3_forecast | Low | Critical | pm_forecast + pm_sprint_health |
| 2 | pm_retro_effectiveness_score | Low | High | pm_retro_analysis + pm_action_items |
| 3 | pm_wip_guardian | Low | High | pm_flow_metrics |
| 4 | pm_individual_health_signal | Medium | Critical | pm_team_health + jira_workload |
| 5 | pm_question_coach | Medium | High | pm_coaching + sprint context |
| 6 | pm_dependency_tracker_live | Medium | High | pm_dependencies |
| 7 | pm_notification_budget | Medium | Medium | notify_routed |
| 8 | pm_pattern_connector | High | Critical | All existing metrics |
| 9 | pm_decision_recommender | High | High | Multiple data sources |
| 10 | pm_stakeholder_update_generator | Low | Medium | pm_exec_report |
| 11 | pm_team_growth_tracker | Medium | Medium | pm_team_health |
| 12 | pm_experiment_evaluator | Low | Medium | pm_experiments |

---

## Key Insight

**The gap between the existing 139 tools and "exceptional PM tooling" is NOT more data — it's TIMING and CORRELATION.**

The tools already collect data. What's missing:
1. **Triggering at the RIGHT MOMENT** (Day 3, not Day 9)
2. **Correlating across signals** (individual + sprint + dependency = root cause)
3. **Presenting DECISIONS not just dashboards** (options with projected outcomes)
4. **Measuring OUTCOMES of actions** (did the retro actually change anything?)
5. **Respecting attention** (notification budget, not notification flood)

The mental model shift: **From "tools that answer questions" → "tools that tell you what question to ask and when to ask it."**
