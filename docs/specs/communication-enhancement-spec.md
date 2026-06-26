# SPEC: Communication Framework Enhancement for AI-Era PM

Status: Draft
Author: Zara (AI) + Aldo
Date: 2026-06-26
Priority: Medium-High

---

## Context

zara-jira-mcp is a 239-tool MCP server acting as AI-powered Scrum Master. It already has strong communication foundations:

**Already Implemented:**
- `pm_compose` — audience-adaptive message writing (Pyramid Principle + BLUF)
- `pm_feedback_coach` — SBI feedback structuring + Radical Candor tone
- `pm_escalate_message` — SCQA escalation format
- `pm_announce_decision` — DACI decision communication
- `pm_comms_plan` — stakeholder communication planning
- `pm_status_draft` — auto-pull Jira data + format per audience
- `pm_coaching` — AI coaching advice with team data context
- `pm_maturity_assessment` — Tuckman stages + leadership stance
- `pm_team_pulse` / `pm_team_radar` — multi-dimension health tracking
- `pm_meeting_effectiveness` — ceremony value measurement
- `pm_stakeholder_pulse` — stakeholder satisfaction over time
- Monte Carlo forecasting (50/70/85/95% confidence)
- 6-channel notifications (Lark, Slack, Discord, Telegram, Teams, Email)
- Auto-escalation when health < 50 or blockers > 3 days

**Research Basis (see docs/communication-frameworks.md):**
- Pyramid Principle, SCARF Model, SBI, RACI/DACI, Radical Candor, 5W1H
- Crucial Conversations (STATE + AMPP), NVC
- Async-First 3-Tier, Information Radiators
- Psychological Safety (Edmondson), Situational Leadership
- Trust-Building in AI-Augmented Teams
- Communication Anti-Patterns (13 identified patterns)

---

## Problem Statement

Despite strong individual tools, the system lacks:

1. **Communication intelligence** — no tracking of communication patterns, effectiveness, or gaps
2. **Proactive nudges** — no reminder when cadence is missed or relationship decaying
3. **Conversation preparation** — no structured prep for difficult conversations (NVC/Crucial)
4. **Communication health scoring** — no unified metric for "how well is this team communicating?"
5. **Learning loop** — feedback given via `pm_feedback_coach` is fire-and-forget, no follow-up tracking

---

## Proposed Enhancements (4 Phases)

### Phase 1: Communication Cadence & Health (Week 1-2)

**Goal:** Make communication gaps visible before they become problems.

#### Tool: `pm_comms_health`

Scans existing data to detect communication anti-patterns.

```
pm_comms_health(board_id: int) -> CommsHealthReport
```

**Logic (no AI needed, pure data):**
- Check: days since last `pm_exec_report` sent (cadence breach if > 14 days)
- Check: stakeholder_pulse trend (declining 2+ sprints = "ghost stakeholder" risk)
- Check: decision_records with no `informed` field populated (DACI gap)
- Check: blockers aging > 3 days with no escalation recorded (escalation avoidance)
- Check: retro action items pending > 2 sprints (follow-through failure)
- Check: meeting_effectiveness scores declining (meeting addiction signal)

**Output:**
```
Communication Health: 72/100

ISSUES:
- [CADENCE] No exec report in 18 days (target: bi-weekly)
- [STAKEHOLDER] PO pulse declining: 4.2 -> 3.8 -> 3.1
- [FOLLOW-THROUGH] 5 retro actions pending > 2 sprints
- [ESCALATION] 2 blockers > 5 days, no escalation recorded

STRENGTHS:
- Decision records complete (8/8 have DACI roles)
- Meeting effectiveness stable (avg 3.8/5)

RECOMMENDATION: Schedule PO sync this week. Stakeholder trust at risk.
```

**Implementation:**
- File: `application/tools/comms_health_handlers.go`
- Registration: `transport/communication.go`
- Module: `stakeholder` (existing)
- DB: uses existing tables (stakeholder_pulse, decisions, blockers, action_items, meeting_effectiveness)
- No new tables needed

---

#### Tool: `pm_cadence_check`

Quick view: is PM meeting communication commitments?

```
pm_cadence_check(board_id: int) -> CadenceReport
```

**Logic:**
- Define default cadences: exec_report=14d, weekly_digest=7d, retro=per-sprint, stakeholder_pulse=per-sprint
- Query last occurrence of each from DB
- Flag overdue items

**Output:**
```
Communication Cadence Status:

  Exec Report:       Last 18 days ago [OVERDUE - target 14d]
  Weekly Digest:     Last 3 days ago [OK]
  Retro Recorded:    Last sprint [OK]
  Stakeholder Pulse: Last 2 sprints ago [OVERDUE]
  Risk Scan:         Last 12 days ago [OVERDUE - target 7d]

Action: 3 cadence items overdue. Run pm_exec_report, pm_stakeholder_pulse, pm_auto_detect_risks.
```

**Implementation:**
- File: same `comms_health_handlers.go`
- DB: query timestamps from existing tables (no new schema)

---

### Phase 2: Conversation Preparation (Week 3-4)

**Goal:** Help PM prepare for difficult conversations using proven frameworks.

#### Tool: `pm_conversation_prep`

Prepares structured talking points for high-stakes conversations using Crucial Conversations + NVC + SCARF awareness.

```
pm_conversation_prep(
  type: string,       // "performance", "conflict", "scope_negotiation", "bad_news", "recognition"
  context: string,    // what's the situation
  person: string,     // optional: who
  board_id: int       // optional: pull relevant data
) -> ConversationPrep
```

**Logic:**
1. Pull relevant data from memory (blockers, sprint health, pulse scores)
2. AI generates structured prep using framework appropriate to `type`:
   - `performance` -> SBI + Radical Candor
   - `conflict` -> NVC (Observation, Feeling, Need, Request)
   - `scope_negotiation` -> Crucial Conversations STATE path
   - `bad_news` -> Pyramid (answer first) + SCARF (reduce threat)
   - `recognition` -> Specific praise + growth connection

**Output format:**
```
CONVERSATION PREP: Performance Feedback
Target: [person]
Framework: SBI + Radical Candor

BEFORE THE CONVERSATION:
- Your intent: [what you want for them, for the team, for the relationship]
- SCARF check: which domains might feel threatened? [Status, Certainty]
- Safety signal: how will you show you care? [mention specific good work first]

OPENING (Care Personally):
"[exact opening sentence suggestion]"

SBI CORE:
- Situation: [from data]
- Behavior: [observable, objective]
- Impact: [quantified where possible]

BRIDGE TO ACTION:
"What's your perspective on this?"
"What would help you with X?"

IF THEY GET DEFENSIVE:
- Return to safety: "[mirror their feeling] + [restate shared goal]"
- AMPP: Ask what they see, Mirror emotion, Paraphrase, Prime if stuck

CLOSE:
- Agreement on next step
- Timeline for check-in
- Reaffirm relationship

DATA CONTEXT (from Jira/memory):
- Sprint completion: 45% (team avg 78%)
- Blocked items assigned: 3 (team avg 0.8)
- Trend: declining 3 sprints
```

**Implementation:**
- File: `application/tools/conversation_prep_handlers.go`
- Registration: `transport/communication.go`
- Uses: `h.AI.Complete` + existing memory queries
- No new tables

---

#### Tool: `pm_hard_conversation`

Shorter version of conversation_prep — just the talking points for a specific difficult message.

```
pm_hard_conversation(
  what: string,       // "need to tell team we're cutting scope"
  why: string,        // "sprint overcommitted, health at 42"
  audience: string    // "team", "individual", "stakeholder"
) -> TalkingPoints
```

**Output:** 5-7 bullet talking points using appropriate framework. Brief, actionable.

**Implementation:** Lightweight AI completion, same file as conversation_prep.

---

### Phase 3: Communication Learning Loop (Week 5-6)

**Goal:** Track whether communication is landing and improving over time.

#### Tool: `pm_feedback_log`

Record that feedback was given, and track follow-up.

```
pm_feedback_log(
  person: string,
  type: string,          // "constructive", "positive", "coaching"
  topic: string,         // brief description
  follow_up_date: string // when to check in
) -> Confirmation
```

**DB Schema (new table):**
```sql
CREATE TABLE IF NOT EXISTS feedback_log (
  id INTEGER PRIMARY KEY,
  person TEXT NOT NULL,
  type TEXT DEFAULT 'constructive',
  topic TEXT NOT NULL,
  follow_up_date TEXT,
  followed_up INTEGER DEFAULT 0,
  outcome TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

---

#### Tool: `pm_feedback_due`

Show overdue feedback follow-ups.

```
pm_feedback_due() -> list of pending follow-ups
```

**Logic:** Query feedback_log where follow_up_date <= today AND followed_up = 0.

---

#### Tool: `pm_feedback_close`

Mark feedback as followed up, record outcome.

```
pm_feedback_close(
  id: int,
  outcome: string    // "improved", "no_change", "escalated", "acknowledged"
)
```

---

#### Tool: `pm_comms_effectiveness`

Aggregate view: are communications landing?

```
pm_comms_effectiveness(sprints: int) -> EffectivenessReport
```

**Logic:**
- Stakeholder pulse trend (going up = communications working)
- Feedback follow-through rate (followed_up / total)
- Meeting effectiveness trend
- Escalation response time (from escalation to resolution)
- Repeat-decision rate (same topic decided >1x = poor original communication)

---

### Phase 4: Proactive Communication Intelligence (Week 7-8)

**Goal:** AI-driven nudges that prevent communication failures.

#### Tool: `pm_comms_nudge`

Run as part of daily digest or standalone. Proactive suggestions.

```
pm_comms_nudge(board_id: int) -> list of nudges
```

**Logic (rule-based, no AI needed):**

| Condition | Nudge |
|-----------|-------|
| Stakeholder pulse < 3.0 | "Schedule sync with [stakeholder]. Satisfaction declining." |
| No exec report > 12 days | "Exec report due soon. Run pm_exec_report." |
| Blocker > 3 days, no escalation | "Blocker [X] aging. Consider pm_escalate." |
| New team member joined (Jira assignee new) | "Onboard [person]: share team agreements, DoD." |
| Sprint goal at risk + no PO update | "PO doesn't know sprint goal is at risk. Run pm_goal_check." |
| Feedback follow-up overdue | "Follow up with [person] on [topic] (due 3 days ago)." |
| Retro actions stale > 2 sprints | "5 retro actions stale. Address in next retro or close." |

**Output:**
```
Communication Nudges (3 items):

1. [STAKEHOLDER] PO satisfaction dropped to 3.1. Schedule a sync to understand concerns.
   -> pm_stakeholder_pulse or pm_conversation_prep(type:"scope_negotiation")

2. [CADENCE] Exec report overdue (16 days). Management may be in the dark.
   -> pm_exec_report(board_id:X)

3. [FEEDBACK] Follow-up with Harry on code review bottleneck (due 5 days ago).
   -> pm_feedback_close(id:7, outcome:"...")
```

**Integration point:** Can be embedded into `pm_daily_digest` output (existing tool) as a "Communication" section.

---

## Architecture Notes

### File Structure
```
application/tools/
  comms_health_handlers.go      # Phase 1: pm_comms_health, pm_cadence_check
  conversation_prep_handlers.go  # Phase 2: pm_conversation_prep, pm_hard_conversation
  feedback_loop_handlers.go      # Phase 3: pm_feedback_log, pm_feedback_due, pm_feedback_close, pm_comms_effectiveness
  comms_nudge_handlers.go        # Phase 4: pm_comms_nudge

transport/
  communication.go               # Add new registrations to existing file
```

### DB Changes
- 1 new table: `feedback_log` (Phase 3 only)
- All other phases use existing tables

### Module Assignment
- All tools belong to module: `stakeholder`
- Available in profiles: `full`, `all`

### Pattern
Follow existing pattern in project:
1. Handler in `application/tools/` with `func (h *Handlers) MethodName(ctx, req) (*mcp.CallToolResult, error)`
2. Registration in `transport/` using `s.AddTool(mcp.NewTool(...))`
3. AI calls via `h.AI.Complete(ctx, systemPrompt, data)` with graceful fallback to raw data
4. Memory queries via `h.Memory.*` interfaces
5. Module gating via existing profile system

---

## New Tools Summary (9 total)

| Phase | Tool | Complexity | AI Required |
|-------|------|-----------|-------------|
| 1 | `pm_comms_health` | Medium | No (pure data) |
| 1 | `pm_cadence_check` | Low | No |
| 2 | `pm_conversation_prep` | Medium | Yes |
| 2 | `pm_hard_conversation` | Low | Yes |
| 3 | `pm_feedback_log` | Low | No |
| 3 | `pm_feedback_due` | Low | No |
| 3 | `pm_feedback_close` | Low | No |
| 3 | `pm_comms_effectiveness` | Medium | No |
| 4 | `pm_comms_nudge` | Medium | No (rule-based) |

---

## Success Criteria

1. `pm_comms_health` accurately detects communication gaps from existing data
2. `pm_conversation_prep` produces actionable, framework-correct prep in < 3 seconds
3. `pm_comms_nudge` produces 0 false positives (every nudge is actionable)
4. Feedback follow-through rate trackable across sprints
5. No performance regression (all new tools are lightweight queries + optional AI)

---

## Dependencies

- Existing: `h.Memory.*` interface methods (GetSprintSnapshots, GetOpenRisks, GetActiveBlockers, GetBlockerHistory, GetPendingActionItems, GetRetrospectives)
- Existing: `h.AI.Complete(ctx, systemPrompt, data)` for AI-powered tools
- Existing: `stakeholder_pulse`, `decisions`, `meeting_effectiveness`, `action_items` tables
- New: `feedback_log` table (Phase 3, one CREATE TABLE statement)

---

## Non-Goals

- No UI/dashboard (MCP server only)
- No real-time chat integration (async tool calls only)
- No replacement of existing tools (additive only)
- No breaking changes to existing API

## Related Specs

- `docs/specs/okr-kpi-bridge-spec.md` — OKR/KPI integration (Jira → Lark OKR translation engine)

---

## Relevant Files (for implementor)

| Purpose | File |
|---------|------|
| Existing communication handlers | `application/tools/communication_handlers.go` |
| Existing communication registration | `transport/communication.go` |
| Existing coaching handlers | `transport/coaching.go` |
| SM leverage (maturity) | `application/tools/sm_leverage_handlers.go` |
| Outcomes (stakeholder pulse) | `application/tools/outcomes_handlers.go` |
| Reporting (escalation) | `application/tools/reporting_handlers.go` |
| Forecast (coaching) | `application/tools/forecast_handlers.go` |
| Communication frameworks doc | `docs/communication-frameworks.md` |
| Reporting guide | `docs/reporting-guide.md` |
| Full tool reference | `SKILL.md` |
| Architecture overview | `AGENTS.md` |

---

## Implementation Order

```
Phase 1 (highest value, lowest effort):
  pm_comms_health -> pm_cadence_check

Phase 2 (high value, medium effort):
  pm_conversation_prep -> pm_hard_conversation

Phase 3 (medium value, low effort):
  pm_feedback_log -> pm_feedback_due -> pm_feedback_close -> pm_comms_effectiveness

Phase 4 (high value, integrates everything):
  pm_comms_nudge (depends on Phase 1 + 3 data)
```

Each phase is independently shippable and testable.
