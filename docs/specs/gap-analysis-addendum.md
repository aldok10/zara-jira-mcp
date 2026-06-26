# SPEC ADDENDUM: Gap Analysis & Corrections

Status: Active
Date: 2026-06-26
Applies to: `communication-enhancement-spec.md` + `okr-kpi-bridge-spec.md`

---

## Critical Corrections

Both specs were written without full awareness of existing implementations. This addendum corrects factual errors and identifies what actually needs building vs what already exists.

---

## 1. OKR/KPI Spec — What ALREADY EXISTS

### Existing OKR System (2 generations)

**Generation 1** (`okr_handlers.go` + `transport/okr_kpi.go`):
| Tool | Description | Status |
|------|-------------|--------|
| `pm_okr_define` | OKR with linked epics, auto quarter detection | Registered, working |
| `pm_okr_progress` | Auto-calculate from Jira epic completion | Registered, working |
| `pm_okr_report` | Formatted report for leadership | Registered, working |
| `pm_kpi_dashboard` | Goal rate, velocity, carryover, blocker resolution, tech debt, WIP | Registered, working |
| `pm_bus_factor` | Single point of failure detection | Registered, working |
| `pm_improvement_dashboard` | Meta-metrics: are we improving? | Registered, working |

**Generation 2** (`okr_v2_handlers.go` + `transport/okr.go`):
| Tool | Description | Status |
|------|-------------|--------|
| `pm_okr_tree` | Full hierarchy: objective → KR → initiative (Jira link) | Registered, working |
| `pm_okr_sync` | Auto-calculate KR progress from JQL + initiatives | Registered, working |
| `pm_okr_kpi` | Sprint KPIs + OKR alignment | Registered, working |
| `pm_okr_report` | AI executive summary with OKR progress | Registered, working (duplicate name!) |

**Existing DB tables (V2):**
```sql
okr_objectives (id, title, description, period, owner, status, board_id, created_at)
okr_key_results (id, objective_id, title, metric_type, target_value, current_value, unit, jql_query, status, created_at)
okr_initiatives (id, key_result_id, issue_key, label, epic_key, description, created_at)
okr_kpi_snapshots (id, key_result_id, sprint_name, value, delta, recorded_at)
```

### Existing Lark OKR Client (`internal/lark/okr.go`):
| Method | Description | Status |
|--------|-------------|--------|
| `ListPeriods()` | Get all OKR periods | Implemented |
| `ListUserOKRs(userID, periodID)` | Get user's OKRs by period | Implemented |
| `BatchGetOKRs(okrIDs)` | Get OKR details by ID | Implemented |
| Token management | Auto-refresh tenant_access_token | Implemented |

**Missing from Lark client:** `CreateProgressRecord()` (write operation)

### What the OKR Spec ACTUALLY Needs to Build

Given existing code, the spec should be REVISED to:

| Proposed in Spec | Actual Status | What to Do |
|------------------|---------------|------------|
| `pm_kpi_dashboard` | EXISTS (both gen1 + gen2) | SKIP — enhance existing if needed |
| `pm_kpi_trend` | NOT exists | BUILD — add trend view per individual KPI |
| `pm_okr_link` | EXISTS as `pm_okr_tree` (add_initiative) | SKIP — already does epic/issue/label → KR linking |
| `pm_okr_progress` | EXISTS (both gen1 + gen2) | SKIP — `pm_okr_sync` already does this |
| `pm_okr_suggest` | NOT exists | BUILD — AI suggestion of OKR alignment |
| `lark_okr_pull` | PARTIAL — client exists, no MCP tool wrapping it | BUILD — wrap existing `OKRClient.ListUserOKRs` as MCP tool |
| `lark_okr_sync_progress` | NOT exists (no write method in client) | BUILD — add `CreateProgressRecord` to client + MCP tool |
| `lark_okr_periods` | PARTIAL — `ListPeriods()` exists, no MCP tool | BUILD — wrap as MCP tool |
| `pm_okr_health` | NOT exists | BUILD — time-elapsed vs progress analysis |
| `pm_okr_report` | EXISTS (name collision!) | RENAME to `pm_okr_quarterly_report` or enhance existing |
| `pm_kpi_to_okr` | NOT exists | BUILD — AI suggestion engine |

**Revised scope: 7 tools to BUILD (not 11)**

---

## 2. Communication Spec — What ALREADY EXISTS

### Existing Communication Tools (registered)

**In `transport/communication.go`:**
| Tool | Framework | Status |
|------|-----------|--------|
| `pm_compose` | Pyramid + BLUF, multi-audience | Registered |
| `pm_status_draft` | Auto-pull Jira + format per audience | Registered |
| `pm_feedback_coach` | SBI + Radical Candor | Registered |
| `pm_escalate_message` | SCQA structured escalation | Registered |
| `pm_announce_decision` | DACI framework | Registered |
| `pm_comms_plan` | Stakeholder communication planning | Registered |

**In `transport/safety.go`:**
| Tool | Framework | Status |
|------|-----------|--------|
| `pm_communicate` | Minto Pyramid (data-enhanced) | Registered |
| `pm_feedback_prep` | SBI + team data context | Registered |
| `pm_escalation_draft` | Pyramid escalation (severity-based) | Registered |
| `pm_decision_record` | ADR format + memory storage | Registered |
| `pm_safety_survey` | Project Aristotle 7Q | Registered |
| `pm_safety_trend` | Safety score over time | Registered |
| `pm_team_aristotle` | 5-pillar team assessment | Registered |

**NOTE:** There are overlapping tools:
- `pm_compose` vs `pm_communicate` — both do audience-adaptive messaging
- `pm_feedback_coach` vs `pm_feedback_prep` — both do SBI feedback
- `pm_escalate_message` vs `pm_escalation_draft` — both do escalation formatting

### What the Communication Spec ACTUALLY Needs to Build

| Proposed in Spec | Actual Status | What to Do |
|------------------|---------------|------------|
| `pm_comms_health` | NOT exists | BUILD — communication anti-pattern scanner |
| `pm_cadence_check` | NOT exists | BUILD — cadence compliance checker |
| `pm_conversation_prep` | PARTIAL — `pm_feedback_prep` covers feedback only | BUILD — expand to cover all conversation types |
| `pm_hard_conversation` | NOT exists (but similar to `pm_escalation_draft`) | BUILD — focused on internal team conversations |
| `pm_feedback_log` | NOT exists | BUILD — track feedback given |
| `pm_feedback_due` | NOT exists | BUILD — follow-up reminders |
| `pm_feedback_close` | NOT exists | BUILD — close the loop |
| `pm_comms_effectiveness` | NOT exists | BUILD — aggregate effectiveness metric |
| `pm_comms_nudge` | NOT exists | BUILD — proactive suggestions |

**Communication spec scope: 9 tools to build (all genuinely new)**

---

## 3. Cross-Cutting Gaps Neither Spec Addresses

### A. Tool Duplication Problem

The project has accumulated duplicate tools across transport registrations:
- 2 OKR define tools (gen1 + gen2)
- 2 OKR progress tools (gen1 + gen2)
- 2 OKR report tools (same name!)
- 2 message composition tools
- 2 feedback tools
- 2 escalation tools

**Recommendation:** Before adding MORE tools, consolidate:
1. Deprecate gen1 OKR tools (keep gen2 as canonical)
2. Deprecate `pm_compose` in favor of `pm_communicate` (data-enhanced version)
3. Choose one feedback tool as canonical

### B. Missing: Lark OKR WRITE Operations

The existing `internal/lark/okr.go` only has READ methods. To sync progress back:

```go
// Missing — needs implementation
func (c *OKRClient) CreateProgressRecord(ctx context.Context, targetID string, content string) error
func (c *OKRClient) UpdateProgressRecord(ctx context.Context, progressID string, content string) error
```

API endpoint: `POST /open-apis/okr/v1/progress_records`
Required scope: `okr:okr.content:writeonly`

### C. Missing: `pm_smart` Router Awareness

The project has a `pm_smart` tool (in `transport/pm_shortcuts.go`) that routes natural language to the right tool. New tools from both specs need to be registered in the router's action map:

```go
// In smart_handlers.go, the action router needs to know about:
case "comms_health": return h.CommsHealth(ctx, req)
case "cadence": return h.CadenceCheck(ctx, req)
case "okr_health": return h.OKRHealth(ctx, req)
// etc.
```

### D. Missing: Module Assignment for New Tools

Both specs must specify which module their tools belong to. Current module map (from `transport/server.go`):

| Module | Profiles Available In |
|--------|----------------------|
| `jira` | standard, full, all |
| `pm` | lite, standard, full, all |
| `ai` | standard, full, all |
| `notifications` | standard, full, all |
| `stakeholder` | full, all |
| `portfolio` | full, all |
| `github` | full, all |
| `integrations` | full, all |
| `shortcuts` | chatgpt, lite, standard, full, all |

**Assignment:**
- Communication tools → module `stakeholder` (full + all profiles)
- OKR/KPI tools → module `pm` (available from lite up)
- Lark OKR sync tools → module `integrations` (full + all only)

### E. Missing: Integration with `pm_daily_digest`

Both specs mention embedding into daily digest but don't specify how. The existing `PMDailyDigestCoaching` in `coaching_handlers.go` already aggregates:
- Sprint status
- Blockers
- Risks
- Pending actions
- Dependencies

**Integration point:** Add optional sections:
- "Communication" section (from `pm_comms_nudge`)
- "OKR alignment" section (from `pm_okr_health` summary)

### F. Missing: FlowMetrics Integration

`pm_flow_metrics` (registered in `transport/server.go` line 637) already calculates:
- Throughput
- Cycle time
- WIP

The OKR spec's Phase 1 KPI dashboard should CALL this handler rather than recalculating.

### G. Missing: Testing Strategy

Neither spec mentions:
- Unit tests for new handlers (follow pattern in `*_test.go` files)
- Integration test with mock Jira/Memory/AI
- Test for Lark OKR client (mock HTTP)

Existing test pattern: `application/tools/handlers_test.go` defines `mockMemory`, `mockJira`, `mockAI` — new handlers should use these.

---

## 4. Revised Implementation Priority

Given existing code, the TRUE gaps ranked by value:

### Immediate Value (builds on existing, no new external deps)
1. `pm_comms_health` — scans existing data for communication anti-patterns
2. `pm_cadence_check` — uses existing table timestamps
3. `pm_okr_health` — time-elapsed vs progress (uses existing V2 OKR tables)
4. `pm_kpi_trend` — individual KPI over time (extends existing dashboard)
5. `pm_kpi_to_okr` — AI suggestion from current metrics to OKR format

### Medium Value (new features, some integration)
6. `pm_comms_nudge` — rule-based proactive suggestions
7. `pm_conversation_prep` — AI conversation framework selector
8. `pm_feedback_log/due/close` — feedback tracking lifecycle (1 new table)
9. `pm_okr_suggest` — AI alignment detection

### High Value, Higher Effort (Lark OKR write integration)
10. `lark_okr_pull` (MCP tool wrapping existing client)
11. `lark_okr_periods` (MCP tool wrapping existing client)
12. `lark_okr_sync_progress` (new client method + MCP tool)

### Housekeeping (should happen in parallel)
- Consolidate duplicate OKR tools (deprecate gen1, keep gen2)
- Fix `pm_okr_report` name collision between gen1 and gen2
- Register new tools in `pm_smart` router

---

## 5. Existing Files Reference (COMPLETE)

| Purpose | File |
|---------|------|
| OKR v1 handlers (gen1) | `application/tools/okr_handlers.go` |
| OKR v2 handlers (gen2, canonical) | `application/tools/okr_v2_handlers.go` |
| OKR v1 registration | `transport/okr_kpi.go` |
| OKR v2 registration | `transport/okr.go` |
| Communication handlers (original) | `application/tools/communication_handlers.go` |
| Communication handlers (enhanced) | `application/tools/communication_enhanced_handlers.go` |
| Communication registration (original) | `transport/communication.go` |
| Communication registration (enhanced) | `transport/safety.go` |
| Lark OKR client (read-only) | `internal/lark/okr.go` |
| Lark webhook client | `internal/lark/webhook.go` |
| Lark calendar client | `internal/calendar/client.go` |
| Flow metrics handler | referenced in `transport/server.go:637` |
| Coaching handlers (daily digest) | `application/tools/coaching_handlers.go` |
| SM leverage (maturity, dysfunction) | `application/tools/sm_leverage_handlers.go` |
| Outcomes (stakeholder pulse, OKR map) | `application/tools/outcomes_handlers.go` |
| Reporting (escalation brief) | `application/tools/reporting_handlers.go` |
| Smart router | `application/tools/smart_handlers.go` |
| Bootstrap/DI | `internal/bootstrap/bootstrap.go` |
| Config | `config/config.go` |
| Mock memory (tests) | `application/tools/handlers_test.go` |

---

## 6. Summary

| Spec | Originally Proposed | Actually Needed | Reason |
|------|--------------------:|----------------:|--------|
| OKR/KPI | 11 tools | 7 tools | 4 already exist (pm_kpi_dashboard, pm_okr_link, pm_okr_progress, pm_okr_report) |
| Communication | 9 tools | 9 tools | All genuinely new (no existing equivalents) |
| **Total** | **20 tools** | **16 tools** | Plus 1 Lark client method + housekeeping |

The existing codebase is more capable than the specs assumed. The real value-add is:
1. **Communication health/cadence/nudge** — completely new capability
2. **Lark OKR write** — the one truly missing bridge (read exists, write doesn't)
3. **OKR health analysis** — time-aware progress assessment
4. **Feedback lifecycle** — tracking whether feedback lands
