# Communication Frameworks for PM/Scrum Master in the AI Era

Panduan komunikasi berbasis riset — framework, template, dan tools yang relevan untuk PM/SM yang bekerja dengan AI.

---

## Kenapa Ini Penting di Era AI

AI mengubah cara tim berkomunikasi:
- **Informasi bergerak lebih cepat** — AI bisa generate report dalam detik, tapi manusia tetap butuh context
- **Asynchronous makin dominan** — distributed teams butuh komunikasi tertulis yang lebih terstruktur
- **Decision fatigue naik** — lebih banyak data tersedia, tapi keputusan tetap harus diambil manusia
- **Trust menjadi currency** — Frontiers in Psychology (2025): informasi yang di-omit AI agent menurunkan trust dan performa tim secara signifikan

**Peran PM/SM bergeser:** dari "status reporter" ke "communication architect" — yang mendesain bagaimana informasi mengalir.

---

## 1. Pyramid Principle (Minto)

**Apa:** Komunikasi dimulai dari kesimpulan, baru supporting arguments, baru data. Kebalikan dari cara kebanyakan orang bicara.

**Kapan pakai:** Setiap kali komunikasi ke atas (management, PO, stakeholder).

**Struktur:**
```
1. Answer first (1 kalimat)
2. Key arguments (2-3 point)
3. Supporting data (detail kalau ditanya)
```

**Contoh:**

Bad:
> "Minggu ini kita selesaikan 12 dari 18 items. Ada 3 blocker dari Tim Platform. John sakit 2 hari. Sprint health 65..."

Good:
> "Sprint AT RISK — 3 cross-team blockers aging > 5 hari. Butuh escalation ke Head of Platform.
> - 67% completion (target 85%)
> - Blockers: API migration (5d), Auth service (4d), DB schema (3d)
> - Team health OK (John back tomorrow)"

**Tool:**
```
pm_exec_report(board_id:X)          → otomatis pakai Pyramid: status first
pm_management_brief(board_id:X, audience:"vp")  → tailored per level
```

---

## 2. SCARF Model (David Rock, NeuroLeadership Institute)

**Apa:** 5 domain sosial yang otak monitor untuk threat/reward: **S**tatus, **C**ertainty, **A**utonomy, **R**elatedness, **F**airness.

**Kenapa PM harus tahu:** Setiap komunikasi bisa memicu threat response di otak penerima. PM yang paham SCARF bisa frame pesan agar diterima, bukan ditolak.

| Domain | Threat (hindari) | Reward (kejar) |
|--------|-----------------|----------------|
| Status | "Kamu satu-satunya yang belum selesai" | "Tim butuh expertise kamu di area ini" |
| Certainty | "Kita nggak tahu kapan selesai" | "Forecast: 85% chance done by Thursday" |
| Autonomy | "Kamu harus pakai approach ini" | "Gimana menurutmu cara terbaik?" |
| Relatedness | Blame individu di depan tim | "Kita bareng-bareng stuck di sini" |
| Fairness | Workload nggak rata tanpa acknowledgment | "Gue lihat kamu carry lebih banyak sprint ini" |

**Tool:**
```
pm_resource_utilization(board_id:X)    → data buat address Fairness
pm_forecast(board_id:X, ...)           → provide Certainty
pm_team_autonomy(board_id:X)           → measure Autonomy level
pm_anti_patterns(board_id:X)           → detect hero culture (Status/Fairness threat)
```

---

## 3. SBI Feedback Model (Center for Creative Leadership)

**Apa:** Feedback terstruktur: **S**ituation (kapan/dimana), **B**ehavior (apa yang dilakukan, observable), **I**mpact (efek ke orang/tim/project).

**Kenapa penting:** Feedback yang nggak structured terasa personal attack. SBI menjaga feedback tetap factual.

**Template:**
```
"Di [standup kemarin], waktu kamu [bilang blocking issue sudah resolved tapi ternyata belum],
dampaknya [Tim QA nunggu 2 hari, sprint health turun 15 point]."
```

**Tool untuk gather data SBI:**
```
pm_blockers(show_history:true)         → Situation: kapan blocker muncul
pm_burndown(board_id:X)                → Impact: sprint progress terganggu
pm_sprint_health(board_id:X)           → Impact: quantified health drop
```

---

## 4. RACI / DACI Decision Framework

**RACI:** Responsible, Accountable, Consulted, Informed — untuk task assignment.
**DACI:** Driver, Approver, Contributor, Informed — untuk decision making.

**Kapan pakai DACI:** Sebelum meeting yang butuh keputusan. Tentukan dulu siapa yang drive, siapa yang approve.

**AI-era twist:** AI bisa jadi Contributor (provide data/analysis), tapi NEVER Approver. Manusia selalu yang memutuskan.

**Tool:**
```
pm_record_decision(title, decision, rationale, made_by)  → record siapa Approver
pm_decide(what, why, who)                                → quick version
pm_dependencies                                          → map RACI across teams
```

---

## 5. Radical Candor (Kim Scott)

**Apa:** Care Personally + Challenge Directly = Radical Candor. 4 kuadran:

| | Challenge Directly | Don't Challenge |
|---|---|---|
| **Care Personally** | Radical Candor | Ruinous Empathy |
| **Don't Care** | Obnoxious Aggression | Manipulative Insincerity |

**Buat PM/SM:** Kebanyakan SM jatuh ke Ruinous Empathy — terlalu baik, nggak berani bilang "sprint ini gagal karena kita overcommit" atau "team, kita punya hero culture problem."

**Order of Operations (4 step):**
1. **Get** — minta feedback ke diri sendiri dulu
2. **Give** — specific praise & kind-but-clear criticism
3. **Gauge** — cek apakah diterima atau defensive
4. **Encourage** — build culture where everyone does this

**Tool:**
```
pm_anti_patterns(board_id:X)           → data buat Challenge Directly (bukan opini)
pm_coaching(topic, situation)           → AI-assisted coaching script
pm_retro_analysis(board_id:X)          → patterns yang perlu di-address
pm_stakeholder_pulse(stakeholder, score, feedback)  → Gauge stakeholder reaction
```

---

## 6. 5W1H (Kipling Method)

**Apa:** Who, What, When, Where, Why, How — checklist completeness untuk setiap komunikasi penting.

**Kapan pakai:** Sebelum kirim announcement, decision record, atau escalation.

| W/H | Cek |
|-----|-----|
| What | Apa yang terjadi / apa keputusannya? |
| Why | Kenapa ini penting / kenapa sekarang? |
| Who | Siapa yang affected / responsible? |
| When | Kapan deadline / kapan terjadi? |
| Where | Di mana impact-nya (project, team, customer)? |
| How | Bagaimana next steps? |

**Tool:**
```
pm_record_decision(title, decision, context, rationale, made_by)  → built-in 5W1H
pm_record_risk(title, severity, owner, mitigation)                → Who, What, How
pm_escalate(board_id:X)                                           → auto-generates 5W1H
```

---

## 7. Asynchronous Communication Protocol

Di era remote/hybrid, PM/SM harus master async. Rules:

**Write-first culture:**
1. **Context up front** — jangan assume orang tahu background. Link to decision record.
2. **Action clear** — setiap message harus jelas: FYI, need input by X, or decision needed.
3. **Deadline explicit** — "need response by EOD Thursday" bukan "ASAP".
4. **One topic per thread** — jangan campur 3 hal dalam 1 message.

**AI-assisted async:**
```
pm_weekly_digest(board_id:X)           → replace status meeting with async digest
pm_standup_prep(board_id:X)            → async standup alternative
pm_release_notes(board_id:X)           → async sprint review for stakeholders
broadcast(message)                     → send to all channels simultaneously
```

---

## 8. Ceremony Communication Patterns

### Standup (2 min/person, synchronous)

**Anti-pattern:** monolog status update.
**Better:** focuskan on blockers, help needed, plan hari ini.

**AI-enhanced standup:**
```
pm_standup_prep(board_id:X)  → auto-detect talking points, skip manual reporting
```

### Retro (safe space, structured)

**Formats berdasarkan situasi:**
- **4Ls** (Liked, Learned, Lacked, Longed For) — reflective, safe
- **Start/Stop/Continue** — action-oriented
- **Mad/Sad/Glad** — emotion-first, good for tension release
- **DAKI** (Drop, Add, Keep, Improve) — concrete

**AI-enhanced retro:**
```
pm_facilitate(ceremony:"retro", board_id:X)  → fresh format each sprint
pm_retro_analysis(board_id:X)                → pattern detection across retros
pm_record_retro(sprint_name, went_well, improvements, action_items)
```

### Sprint Review (demo + stakeholder)

**Framework:** Show, don't tell. Demo > slides.

**Structure:**
1. Sprint goal recap (10 sec)
2. Demo what shipped (bulk of time)
3. What didn't ship + why (brief, honest)
4. Next sprint preview (30 sec)

**AI-enhanced review:**
```
pm_review_prep(board_id:X)             → demo order, talking points
pm_release_notes(board_id:X)           → categorized what shipped
pm_goal_check(board_id:X)              → goal achievement assessment
```

---

## 9. Escalation Communication

**Rule:** Escalation bukan failure. Late escalation is failure.

**TIRED framework buat escalation:**
- **T**imeframe — berapa lama udah stuck?
- **I**mpact — apa yang delayed/affected?
- **R**equested action — apa yang butuh dari penerima?
- **E**vidence — data yang support urgency
- **D**eadline — kapan perlu resolved?

**Tool:**
```
pm_impediment_aging                    → Timeframe + Evidence
pm_escalation_report(board_id:X)       → auto-generate TIRED format
pm_escalate(board_id:X)                → send to channels with context
pm_blocker_aging                       → SLA tracking per blocker
```

---

## 10. Stakeholder Mapping + Communication Cadence

| Stakeholder | Needs | Cadence | Tool |
|-------------|-------|---------|------|
| VP/C-Level | Business outcomes, risks | Bi-weekly | `pm_exec_report` |
| Product Owner | Goal progress, scope, timeline | Daily/as-needed | `pm_goal_check`, `pm_scope_creep` |
| Engineering Lead | Health, blockers, capacity | Weekly | `pm_sprint_health`, `pm_resource_utilization` |
| Cross-team leads | Dependencies, blockers | Weekly | `pm_dependency_report`, `portfolio_blockers` |
| Team | Ceremonies, coaching | Daily | `pm_standup_prep`, `pm_facilitate` |
| New members | Context, how-we-work | Onboarding | `pm_team_kb`, `pm_dod`, `pm_agreements` |

---

## Quick Reference: Framework per Situation

| Situation | Framework | Tool |
|-----------|-----------|------|
| Update ke management | Pyramid Principle | `pm_exec_report`, `pm_management_brief` |
| Giving feedback to team | SBI + SCARF awareness | `pm_coaching`, `pm_anti_patterns` |
| Making a decision | DACI + record | `pm_decide`, `pm_record_decision` |
| Escalating a blocker | TIRED | `pm_escalation_report`, `pm_escalate` |
| Running a retro | 4Ls / Start-Stop-Continue | `pm_facilitate`, `pm_retro_analysis` |
| Cross-team sync | RACI + dependency map | `pm_dependencies`, `portfolio_blockers` |
| Async weekly update | Pyramid + 5W1H | `pm_weekly_digest` |
| Sprint review with PO | Demo + honest status | `pm_review_prep`, `pm_goal_check` |
| Dealing with conflict | NVC (observe, feel, need, request) | `pm_coaching(topic:"conflict")` |
| Proving SM value | Data + narrative | `pm_sm_impact`, `pm_maturity_assessment` |

---

## AI sebagai Communication Partner

Di 2026, AI bukan pengganti komunikator — tapi amplifier. AI bantu PM/SM:

1. **Gather data** sebelum komunikasi (bukan asumsi, tapi fakta)
2. **Structure** pesan sesuai audience (exec vs team vs PO)
3. **Detect patterns** yang perlu dikomunikasikan (anti-patterns, risks aging)
4. **Automate routine** (weekly digest, standup prep, release notes)
5. **Track effectiveness** (stakeholder pulse, improvement velocity)

Yang AI TIDAK bisa gantikan:
- Empathy dalam percakapan 1-on-1
- Reading the room di retro
- Trust-building melalui consistency
- Difficult conversations yang butuh nuance

**PM/SM terbaik di era AI:** yang pakai AI buat eliminate busywork, dan pakai waktu yang di-save buat human connection.

---

## 11. Signal-over-Noise Protocol

**Apa:** Prinsip komunikasi dari LaaS/Litmus framework — move from "constant chatter" to highly structured, async-by-default communication yang meninggalkan "searchable decision trail" untuk future team members.

**Rules:**
1. **Context over pings** — tanya full question up front, no "hey, ada waktu?"
2. **Async friendly** — assume reader respond later, kasih enough info untuk act tanpa follow-up
3. **Signal review states** — make status visible tanpa perlu tanya (PR labels, Jira transitions, dashboards)
4. **One topic per thread** — no spaghetti threads
5. **Searchable trail** — every decision, every rationale, discoverable later

**Stats:**
- Teams yang pakai structured async: 25% fewer meetings, 3x faster decisions (Microsoft 2026 Work Trend)
- $75M at risk per $1B dari komunikasi yang buruk (PMI)
- Communication overload → 3x slower simple decisions (ITS Dart research)

**Tool:**
```
pm_record_decision(...)         → searchable decision trail
pm_search_decisions(query)      → retrieve past context
pm_weekly_digest(board_id:X)    → replace sync meetings with structured async
notify_routed(message, severity) → auto-route to right channel
pm_team_kb(question)            → self-service knowledge, no need to ask someone
```

---

## 12. Crucial Conversations (Patterson, Grenny, McMillan, Switzler)

**Apa:** Framework buat high-stakes, opposing opinions, strong emotions conversations — dimana hasilnya sangat berpengaruh.

**Core Steps:**
1. **Start with heart** — clarify what you really want (for yourself, for them, for the relationship)
2. **Learn to look** — detect when safety is at risk (silence or violence)
3. **Make it safe** — mutual purpose + mutual respect
4. **Master my stories** — separate facts from the narrative you tell yourself
5. **STATE my path** — Share facts, Tell your story, Ask for their path, Talk tentatively, Encourage testing
6. **Explore others' path** — AMPP (Ask, Mirror, Paraphrase, Prime)
7. **Move to action** — Who does What by When, follow up How

**Kapan pakai:** Performance issue, missed commitment, team conflict, scope dispute with PO, escalation to management.

**Tool:**
```
pm_coaching(topic:"crucial_conversation", situation:"...")  → AI-prep conversation steps
pm_hard_conversation(context, data_points)                  → generate talking points with data
pm_blockers(show_history:true)                             → facts for "Master my stories"
pm_stakeholder_pulse(...)                                  → track relationship post-conversation
```

---

## 13. Communication Anti-Patterns in Software Teams

**Research basis:** "Community Smells" (sociotechnical anti-patterns), Signal Retention research, Microsoft golazo patterns.

| Anti-Pattern | Signal | Consequence |
|---|---|---|
| **Information hoarding** | Key knowledge in 1 person's head | Bus factor risk, decisions delayed |
| **Ghost stakeholders** | Stakeholder listed but never engaged | Surprise objections at 11th hour |
| **Over-communication** | 15+ messages/day in channel, nobody reads | False sense of alignment, decision paralysis |
| **One-way reporting** | PM broadcasts, team never responds | Zero feedback loop, missed concerns |
| **Re-deciding** | Same topic decided 3x in 3 sprints | Poor decision record, no searchable trail |
| **Meeting addiction** | Every question → "let's schedule a call" | Deep work destroyed, async skills atrophy |
| **Status theater** | Elaborate status reports nobody reads | PM time wasted on reporting, not facilitating |
| **Blame language** | "You didn't deliver" instead of "We're stuck" | SCARF threat (Status + Relatedness), trust erosion |

**Detection (existing tools):**
```
pm_anti_patterns(board_id:X)    → detects hero culture, dead retros (related patterns)
pm_impediment_aging             → chronic blockers = possible information hoarding
pm_stakeholder_trend            → declining pulse = ghost stakeholder signal
pm_sm_impact(sprint_name)       → high meeting count + low decision count = meeting addiction
```

**Future tool idea:** `pm_comms_health` — scan for these patterns across decision records, blocker aging, stakeholder pulse trends.

---

## 14. Trust-Building in AI-Augmented Teams

**Context:** Trust in AI tools dropped 31% (2025), agentic AI trust dropped 89% (Axis Intelligence 2026). PM yang pakai AI tools perlu actively build trust.

**Trust Pyramid for AI-era PM:**
```
Level 4: AUTONOMY    → Team trusts AI recommendations enough to act on them
Level 3: RELIABILITY → AI consistently delivers accurate, useful outputs
Level 2: OPENNESS    → PM transparent about what AI does/doesn't know
Level 1: SAFETY      → Team feels safe to say "AI is wrong here"
```

**Practical trust-building:**
1. **Show your work** — when using AI forecast, show confidence intervals not just a date
2. **Acknowledge uncertainty** — "AI says 85% chance, but doesn't factor in holiday next week"
3. **Let team override** — AI recommendation is input, not mandate
4. **Track accuracy** — publish how accurate past forecasts were (calibration)
5. **Start small** — use AI for low-stakes first (standup prep), build to high-stakes (forecasting)

**Tool:**
```
pm_forecast(board_id:X, ...)    → shows 50/70/85/95% (transparent uncertainty)
pm_confidence(sprint, score)     → human confidence alongside AI
pm_scorecard(board_id:X)        → tracked accuracy over time
pm_coaching(topic:"trust")       → how to build trust with specific situation
```

---

## References

- Minto, B. (1987). The Pyramid Principle. Logic in Writing and Thinking.
- Rock, D. (2008). SCARF: A Brain-Based Model for Collaborating With and Influencing Others.
- Scott, K. (2017). Radical Candor: Be a Kick-Ass Boss Without Losing Your Humanity.
- Center for Creative Leadership. SBI Feedback Model.
- Rosenberg, M. (2003). Nonviolent Communication: A Language of Life.
- Patterson, K. et al. (2002, 4th ed 2021). Crucial Conversations: Tools for Talking When Stakes are High.
- Atlassian Team Playbook. DACI Decision Framework.
- Frontiers in Psychology (2025). Trust in Human-AI Team Communication.
- DORA 2025 Report. AI as amplifier of existing capabilities.
- Gartner (2026). 40% enterprise apps will integrate task-specific AI agents.
- PMI. Ineffective communication = $75M at risk per $1B spent.
- Microsoft 2026 Work Trend Index. 15.4 hrs/week in meetings, 12.1 hrs deep work.
- Axis Intelligence (2026). Trust in company AI declined 31%, agentic AI trust -89%.
- LaaS/Litmus Framework. Signal-over-Noise Protocol for scaling teams.
- Microsoft golazo/communication. Context Over Pings, Async Friendly, Signal Review States.
- ITS Dart Research. Communication overload → 3x slower decisions.
- EmergentMind/Community Smells. Sociotechnical anti-patterns in software teams.
- Scrum.org (2025). Smart Scrum Masters using AI to read between the lines.
- AgileVelocity. Monte Carlo simulations for sprint planning confidence intervals.
