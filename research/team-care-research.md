# PM/Scrum Master Deep Research: Metrics, Psychology, and Team Care

> Compiled from academic papers, industry research (DORA, SPACE, Google Project Aristotle),
> and practitioner evidence. Focus: making zara-jira-mcp a tool that CARES about the team.

---

## 1. METRICS THAT ACTUALLY MATTER

### The Three Pillars of Team Health (SPACE Framework, Microsoft/GitHub 2021)

| Pillar | What It Means | How We Measure |
|--------|--------------|----------------|
| **Satisfaction** | Developers happy with work | Confidence votes, retro sentiment |
| **Performance** | Quality of output | Escaped defects, sprint goal success rate |
| **Activity** | Volume (but NOT a target) | Throughput, deployment frequency |
| **Communication** | Collaboration quality | PR review time, knowledge sharing |
| **Efficiency** | Flow state, low friction | Cycle time, WIP, context switching |

**Key Insight:** Measuring only Activity (velocity, story points) without Satisfaction and Efficiency leads to burnout. [Source: Cortex 2024 State of Developer Productivity]

### Evidence-Based Management (EBM, Scrum.org)

Four Key Value Areas:
1. **Current Value** — How happy are users NOW? (Customer satisfaction, revenue per employee)
2. **Unrealized Value** — What value COULD we deliver? (Market share gap)
3. **Ability to Innovate** — Can we deliver new things? (Tech debt ratio, time on features vs bugs)
4. **Time to Market** — How fast from idea to production? (Lead time, deployment frequency)

### The Numbers That Matter

| Metric | Healthy Range | Danger Signal | Source |
|--------|-------------|---------------|--------|
| Sprint Goal Success Rate | >70% | <52% (industry avg!) | Scrum Alliance |
| Carryover Ratio | <15% | >30% | Industry standard |
| WIP per person | 1-2 items | >3 items | Flow theory |
| Cycle Time (stories) | 1-5 days | >10 days | Agile seekers |
| Focus Factor | 60-70% | <50% | Capacity planning |
| Tech Debt Ratio | <20% | >30% | Multiple sources |
| Context Switching Cost | 23 min per switch | >10 switches/day = 4hr lost | Research (Cortex 2024) |
| Velocity Variance (CV) | <30% | >40% = unpredictable | Statistical standard |
| PR Review Time | <24 hours | >3 days = bottleneck | DORA |
| Developer Time Lost | baseline | 5-15 hrs/week avg | Cortex 2024 |

---

## 2. PSYCHOLOGY: THE INVISIBLE METRICS

### Google Project Aristotle (180 teams studied)

Top 5 factors for team effectiveness (ranked):
1. **Psychological Safety** (43% of variance in performance!)
2. **Dependability** — Can I count on teammates?
3. **Structure & Clarity** — Are roles/goals clear?
4. **Meaning** — Does work matter personally?
5. **Impact** — Do I believe our work matters?

**Critical finding:** Individual talent matters less than team dynamics.
The best teams report MORE mistakes — not because they make more, but because they detect and surface them earlier.

### Burnout Signals (Maslach Burnout Inventory, adapted for agile)

| Signal | Observable in Data | Early Warning |
|--------|-------------------|---------------|
| Emotional exhaustion | Declining velocity, more sick days | Consistent overcommitment |
| Depersonalization | Less retro participation, terse comments | Rising carryover |
| Reduced efficacy | More bugs, longer cycle time | Story points inflation |

### Sustainable Pace (Agile Manifesto Principle 8)

"Agile processes promote sustainable development. The sponsors, developers, and users should be able to maintain a constant pace indefinitely."

**What this means in numbers:**
- Focus factor should be 60-70% (NOT 100%)
- Reserve 15-20% for debt/learning/unexpected
- Leave 10-15% buffer for interruptions
- NEVER plan >80% capacity utilization

### Overcommitment Research

"The Ethics of Over-Allocation in Sprints" (PMI, 2024):
- Teams routinely pushed beyond sustainable velocity
- Creates ethical breach of "respect for human capital"
- Leads to: hidden overtime, quality shortcuts, technical debt accumulation
- Solution: Capacity planning based on ACTUAL availability, not idealized

---

## 3. STORY POINTS & OVERLOAD PREVENTION

### Why Story Points Lie

"Your story points are a mathematical lie" — Scrum Alliance, 2024

Problems:
- Different teams use different scales (13 ≠ 13)
- Points inflate over time (velocity looks good, output doesn't change)
- Used as TARGETS instead of SIGNALS
- Individual productivity tracking destroys collaboration

### What To Do Instead

1. **Count items** — Throughput (items/sprint) is more reliable than points
2. **Track cycle time** — How long from start to done (reveals bottlenecks)
3. **Measure outcomes** — Sprint goal hit rate, not velocity
4. **Use points ONLY for relative sizing** — "Is this bigger than that?"

### Overload Detection Formula

```
Overload Risk = (assigned_points / avg_velocity) * 100

> 100% = OVERCOMMITTED (red flag)
> 85% = pushing limits
60-80% = sustainable
< 50% = underutilized (might add more)
```

Per-person variant:
```
Individual Load = person_assigned_points / (team_avg_velocity / team_size)

> 150% = THIS PERSON IS OVERLOADED
> 120% = watch closely
80-120% = balanced
< 60% = might have blockers or be underutilized
```

---

## 4. SCHEDULING & DAILY SUMMARIES

### What a PM Needs to Know Every Day

1. **What changed since yesterday?** (new blockers, completed items, status changes)
2. **Who might need help?** (stuck items, high WIP, overdue items)
3. **Are we on track for sprint goal?** (burn rate projection)
4. **Any external risks materialized?** (new high-priority bugs, failed deployments)

### Optimal Cadence (Research-backed)

| Report | Frequency | Audience | Content |
|--------|-----------|----------|---------|
| Daily digest | Every morning | PM/SM | Delta from yesterday, blockers, burn rate |
| Sprint health | Twice/sprint | Team | Health score, flow metrics, goal progress |
| Weekly digest | Friday | Stakeholders | Wins, concerns, decisions, next week |
| Sprint report | End of sprint | All | Scorecard, release notes, retro themes |
| Executive summary | Bi-weekly | Leadership | Business outcomes, risks, forecast |

---

## 5. THE "CARING TOOLS" — WHAT MAKES THIS DIFFERENT

Traditional PM tools: **track**, **measure**, **report**, **demand**
This MCP should: **detect**, **protect**, **support**, **grow**

| Traditional | Caring Alternative |
|-------------|-------------------|
| "Velocity is dropping" | "Team might be overloaded — check WIP per person" |
| "Sprint goal missed" | "What blocked us? Was the goal realistic given capacity?" |
| "Alice has low throughput" | "Alice has been blocked 3 times this sprint — she needs help" |
| "We need more points" | "At sustainable pace, we can deliver X. Want to reduce scope?" |
| "Why is this late?" | "This item has high cycle time — what's the bottleneck?" |

### Signals That Trigger "Care Mode"

1. **Overcommitment detected** → Suggest scope reduction, NOT pressure
2. **Same person blocked repeatedly** → Suggest pair programming or ownership change
3. **Carryover >3 sprints for same item** → Suggest splitting or abandoning
4. **Velocity decline + no new members** → Ask about team wellbeing, not performance
5. **High WIP + long cycle time** → Suggest WIP limits, finish before starting

---

## 6. KEY RESEARCH NUMBERS

- **23 minutes 15 seconds** — cost of single context switch to regain focus
- **50%** of developers lose 10+ hours weekly to workflow disruptions
- **43%** of team performance variance explained by psychological safety
- **52%** of teams consistently meet sprint goals (industry average)
- **15-20%** recommended sprint allocation for tech debt
- **60-70%** optimal focus factor (not 100%)
- **5-15 hours/week** per developer lost to unproductive work
- **40%** higher turnover in teams that push CI without DevEx investment
- **98% increase** in PRs merged with AI tools, but **242% increase** in incidents (DORA 2025)
- **80%** maximum capacity utilization for sustainable pace
- **1-2 items** optimal WIP per person (flow theory)
- **<5 days** healthy cycle time for user stories
- **<24 hours** healthy PR review time

---

## 7. IMPLEMENTATION PRIORITIES FOR zara-jira-mcp

Based on all research, highest-impact additions:

1. **Daily Delta Report** — "What changed since yesterday" auto-generated
2. **Overload Detection** — Per-person WIP and assignment load monitoring
3. **Sprint Commitment Validation** — "Are we overcommitting?" at planning time
4. **Psychological Safety Signals** — Track retro participation, blame-free metrics
5. **Sustainable Pace Monitor** — Is the team burning out? (velocity + carryover + overtime signals)
6. **Sprint Goal Focus** — Are we working on goal-related items or getting pulled into other work?
7. **Context Switching Cost** — How many items is each person juggling?
8. **Cycle Time Alerts** — Items stuck too long get flagged automatically

These tools don't DEMAND — they PROTECT.
They don't BLAME — they DETECT early so PM can HELP.
