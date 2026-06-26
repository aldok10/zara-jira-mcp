# Research Summary: PM/SM in AI Era — What's Next

> Compiled from 50+ sources. Foundation for next iteration of zara-jira-mcp.

---

## 1. Developer Productivity Measurement (DX Core 4, SPACE, DORA)

### Key Findings
- **DX Core 4** (2024): Speed, Effectiveness, Quality, Business Impact — survives AI cycles [1]
- **SPACE** (GitHub/Microsoft): Satisfaction, Performance, Activity, Communication, Efficiency — holistic view [2]
- **DORA 2024**: AI increases individual output BUT decreases delivery stability by 7.2% [3]
- **AI ROI paradox**: 4-10x output boost for some, 19% SLOWDOWN for experienced devs on complex tasks [4]
- Traditional metrics (LOC, commits, PRs) are useless in AI era — measure OUTCOMES not output [5]
- Elite teams: cycle time <2.5 days, good teams: 4-7 days (LinearB 2025) [6]

### What This Means for Tool
- Need: `pm_engineering_metrics` — track DORA + SPACE from Jira data
- Need: `pm_ai_impact_check` — compare velocity pre/post AI adoption
- Need: outcome-tracking beyond story points

### Sources
[1] https://open.substack.com/pub/abinoda/p/revisiting-the-dx-core-4
[2] https://www.cortex.io/post/space-metrics
[3] DORA 2024 Report
[4] https://blog.exceeds.ai/frameworks-developer-productivity-ai-roi/
[5] https://www.cortex.io/post/metrics-for-measuring-developer-productivity
[6] https://www.worklytics.co/resources/2025-software-engineer-productivity-score-benchmarks

---

## 2. AI-Augmented Project Management (Systematic Reviews 2025-2026)

### Key Findings
- ML models (SVM, neural networks, ensemble) enhance estimation accuracy, resource utilization, risk reduction [1]
- LLMs can automate: reporting, requirements writing, scope definition [2]
- Transition: opinion-based PM → evidence-based PM → AI-augmented PM [3]
- AI as "Scrum Master": automates junior office work + complex creative tasks [4]
- PMI 2025: 68% organizations using/piloting AI in project workflows [5]
- Key risk: AI helps bad PMs, can hurt top performers (adds overhead) [6]

### What This Means for Tool
- Our PM intelligence tools are research-validated
- Need: `pm_estimation_calibration` — track estimate vs actual, improve over time
- Need: `pm_requirement_quality` — AI check story quality (INVEST criteria)
- Opportunity: auto-generate user stories from high-level objectives

### Sources
[1] https://link.springer.com/article/10.1007/s10515-025-00578-6
[2] https://link.springer.com/chapter/10.1007/978-3-031-72781-8_12
[3] https://pmworldjournal.com/article/from-opinion-based-to-ai-augmented-project-management
[4] https://www.mdpi.com/2079-8954/13/3/208
[5] PMI AI in PM Survey 2025
[6] Beehiiv PM Research

---

## 3. Team Topologies & Cognitive Load

### Key Findings
- Teams >15 people: trust breaks down, cognitive load spikes [1]
- **AI transforms cognitive load** — becomes "anticipation burden" + "decision throughput problem" [2]
- Platform teams reduce cognitive load on delivery teams by 40-60% [3]
- Fast flow requires: managed cognitive load + clear team interactions + bounded context [4]
- Interaction modes: Collaboration, X-as-a-Service, Facilitating [5]

### What This Means for Tool
- Need: `pm_team_cognitive_load` — assess if team is overloaded (scope breadth, tech stack diversity, domain complexity)
- Need: `pm_interaction_map` — track inter-team dependencies pattern
- SM role: reduce cognitive load, not add. Every tool output should be pre-digested.

### Sources
[1] https://teamtopologies.com/news-blogs-newsletters/when-teams-grow-too-large-solving-cognitive-load-issues
[2] https://blog.owulveryck.info/2026/06/24/who-does-what-team-topologies-for-the-agentic-platform.html
[3] https://teamtopologies.com/creditas
[4] https://teamtopologies.com/news-blogs-newsletters/moving-beyond-agile-rituals-designing-the-whole-organization-for-fast-flow
[5] Team Topologies (Skelton & Pais)

---

## 4. Remote/Async Work Effectiveness

### Key Findings
- Async-first: 23% faster project completion on distributed teams (GitLab 2025) [1]
- Microsoft 2025: 275 interruptions/day, 3x more meetings since 2020 [2]
- Buffer 2023: 71% say clear async processes improve productivity [3]
- Key metrics for remote: cycle time, deployment frequency, change failure rate [4]
- Deep work blocks (4hr) = most impactful intervention for engineer productivity [5]
- Meeting-free days: 35% increase in focused work output [6]

### What This Means for Tool
- Need: `pm_async_health` — measure team's async vs sync ratio
- Need: `pm_meeting_load` — quantify meeting overhead per person
- Existing `sm_meeting_roi` is validated by this research
- Auto-generated async updates replace meetings (already doing this)

### Sources
[1] https://stealthagents.com/research/asynchronous-work-statistics-2026
[2] Microsoft Work Trend Index 2025
[3] Buffer State of Remote Work 2023
[4] https://www.thirstysprout.com/post/manage-remote-teams
[5] GitHub DevEx Research 2024
[6] https://www.questworks.games/blog/async-communication-guide-remote-teams

---

## 5. Probabilistic Forecasting & Sprint Predictability

### Key Findings
- Monte Carlo on throughput beats velocity-based estimation consistently [1]
- Script tested: simulated hundreds of historical forecasts at 7/14/21/30/60/90 day horizons [2]
- Teams using probabilistic forecasting: 40% fewer missed deadlines [3]
- Key insight: count items completed, not story points — simpler, more reliable [4]
- "When will it be done?" is answerable with 85% confidence using 8+ sprints of data [5]

### What This Means for Tool
- Our `pm_forecast_sprint` is validated
- Need: `pm_forecast_epic` — when will this epic be done? (multi-sprint simulation)
- Need: `pm_forecast_release` — when will this set of features ship?
- Need: `pm_predictability_index` — how predictable is this team? (actual/committed ratio)

### Sources
[1] https://www.agilevelocity.com/blog/harnessing-monte-carlo-simulations-for-more-accurate-sprint-planning
[2] https://blog.leadingedje.com/post/physics-of-predictability.html
[3] https://www.techademy.com/ai-velocity-forecasting
[4] https://mariachec.substack.com/p/alternative-to-estimations-monte-carlo-simulation
[5] https://medium.com/expedia-group-tech/monte-carlo-forecasting-in-software-delivery-474bb49cb3f9

---

## 6. OKR Alignment & Outcome Thinking

### Key Findings
- Teams connecting goals to outcomes: 30% more likely to hit them [1]
- Output ≠ Outcome: shipping features means nothing if users don't benefit [2]
- Engineering OKRs: quality, reliability, performance — not just delivery [3]
- Sprint goals should map to OKRs — most teams fail here [4]
- Product management OKRs drive product strategy, not just task completion [5]

### What This Means for Tool
- Need: `pm_okr_alignment` — map sprint items to OKRs, show coverage gaps
- Our `pm_outcome_map` partially covers this
- Need: `pm_outcome_vs_output` — ratio of outcome-linked work vs busy-work
- Sprint goal → OKR key result connection should be explicit

### Sources
[1] https://www.okrstool.com/blog/okrs-agile
[2] https://mooncamp.com/blog/output-vs-outcome
[3] https://okrinstitute.org/it/okr-for-engineering/
[4] https://www.scrum.org/resources/blog/transforming-output-oriented-mindset-outcome-oriented-mindset-using-okrs
[5] https://www.opsmatters.com/posts/high-output-high-outcome-how-product-teams-can-shift-focus-okr-culture

---

## 7. Technical Debt Management

### Key Findings
- Developers spend 42% of time on tech debt and maintenance (Stripe 2023) [1]
- McKinsey: TD = 20-40% of technology estate value before depreciation [2]
- 80th percentile TDS companies: 20% higher revenue growth [3]
- Best practice: 15-20% of every sprint allocated to TD [4]
- Score debt with RICE so it competes fairly with features [5]
- Leading indicators: deployment frequency dropping, bug rate rising, onboarding time increasing [6]

### What This Means for Tool
- Need: `pm_tech_debt_budget` — enforce 15-20% sprint allocation tracking
- Our `pm_tech_debt_ratio` exists but could be smarter
- Need: `pm_debt_leading_indicators` — detect when debt is accumulating silently
- Need: `pm_debt_roi` — quantify debt paydown in business terms

### Sources
[1] https://wojciechowski.app/en/articles/technical-debt-quantification
[2] https://www.mckinsey.com/business-functions/mckinsey-digital/our-insights/demystifying-digital-dark-matter-a-new-standard-to-tame-technical-debt
[3] McKinsey Digital 2024
[4] https://www.ideaplan.io/guides/technical-debt-for-product-managers
[5] RICE scoring for debt prioritization
[6] https://sourcegraph.com/blog/technical-debt-management

---

## 8. Communication in AI Era

### Key Findings
- AI causes trust issues when communication is unclear in human-AI teams [1]
- Cognitive offloading to AI improves coping — only when well-structured [2]
- Minto Pyramid: lead with answer, then support — optimal for execs [3]
- SBI feedback model removes subjectivity from coaching [4]
- DACI: clear decision roles prevent decision unraveling [5]
- 58% meeting time = wasted (Jabra 2025) — structured async is the answer [6]
- Radical Candor: Care + Challenge = trust building [7]
- SCQA (McKinsey): Situation-Complication-Question-Answer for escalation [8]

### What This Means for Tool
- Our `comm_*` tools are research-validated
- All existing and future communication output should apply Minto by default
- Need: AI output that adapts AUTOMATICALLY based on detected audience
- Need: `comm_retro_prompt` — non-repetitive retrospective format generator

### Sources
[1] https://www.frontiersin.org/journals/psychology/articles/10.3389/fpsyg.2025.1637339/full
[2] https://www.frontiersin.org/journals/psychology/articles/10.3389/fpsyg.2025.1699320/full
[3] https://www.strategypunk.com/the-minto-pyramid-principle
[4] https://www.ccl.org/articles/leading-effectively-articles/closing-the-gap-between-intent-vs-impact-sbii/
[5] https://project-management.com/daci-model/
[6] Jabra Meeting Report 2025
[7] https://www.radicalcandor.com/blog/culture-of-feedback
[8] McKinsey Minto Pyramid / SCQA

---

## NEXT ITERATION: Prioritized Tasks

### P0 — High Value, Low Effort (do tomorrow)

| Task | Research Basis | Effort |
|------|---------------|--------|
| `pm_forecast_epic` — Monte Carlo for epic completion | Probabilistic forecasting | 2hr |
| `pm_predictability_index` — team predictability score | Sprint predictability research | 1hr |
| `pm_requirement_quality` — INVEST criteria check | AI-augmented PM papers | 2hr |
| `pm_outcome_vs_output` — work alignment check | OKR research, 30% improvement | 1hr |

### P1 — High Value, Medium Effort (this week)

| Task | Research Basis | Effort |
|------|---------------|--------|
| `pm_engineering_metrics` — DORA + SPACE from Jira | DX Core 4, DORA 2024 | 4hr |
| `pm_team_cognitive_load` — overload assessment | Team Topologies | 3hr |
| `pm_tech_debt_budget` — 15-20% allocation enforcement | McKinsey/Stripe research | 2hr |
| `pm_estimation_calibration` — track estimate accuracy | AI estimation papers | 3hr |
| `pm_async_health` — async vs sync ratio | Remote work research | 2hr |

### P2 — Medium Value, Needs More Research

| Task | Research Basis | Effort |
|------|---------------|--------|
| Auto-generate user stories from objectives | LLM automation paper | 6hr |
| `pm_ai_impact_check` — pre/post AI velocity | DORA AI paradox | 4hr |
| `pm_interaction_map` — Team Topologies viz | Team Topologies | 6hr |
| `pm_debt_roi` — quantify TD in business terms | McKinsey valuation | 4hr |
| Integration with actual DORA metrics source | DORA research | 8hr |

### P3 — Experimental / Future

| Task | Research Basis |
|------|---------------|
| Sentiment analysis on team communication | Frontiers 2025 empathy research |
| Auto-detect team stage changes over time | Tuckman longitudinal tracking |
| Predictive model for sprint failure risk | ML for PM paper |
| NLP on retro themes for pattern mining | LLM automation |
| Cross-team flow efficiency metrics | Team Topologies fast flow |

---

## Key Numbers to Remember

| Stat | Source |
|------|--------|
| 58% meeting time wasted | Jabra 2025 |
| 42% dev time on tech debt | Stripe 2023 |
| 20-40% of tech estate = debt | McKinsey |
| 275 interruptions/day | Microsoft 2025 |
| 20+ min to recover from interrupt | Context switching research |
| 23% faster delivery (async-first) | GitLab 2025 |
| 30% more goal hits (outcome-linked) | OKR research |
| 19% slower (AI on complex tasks) | AI ROI research |
| 7.2% delivery stability drop (AI) | DORA 2024 |
| 4-10x output boost (AI simple tasks) | AI ROI research |
| <2.5 day cycle time = elite | LinearB 2025 |
| 15-20% sprint = optimal TD budget | Industry consensus |
| 8+ sprints data = reliable forecast | Monte Carlo research |
| 68% orgs using AI in PM | PMI 2025 |
| 3x more meetings since 2020 | Microsoft 2025 |
