# zara-jira-mcp — AI Agent Rules

**Non-negotiable coding standards for all AI agents working on this codebase.**

Every rule includes WHY it exists. Breaking a rule requires documented justification in the PR.

---

## Project Identity

- **Language**: Go 1.26 (uses modern features: `slices.*`, `maps.*`, `for range n`, `cmp.Or`, `omitzero`)
- **Architecture**: Ports & Adapters (Hexagonal Architecture), modular DDD
- **Transport**: MCP stdio (primary), SSE/HTTP secondary. Uses `mcp-go` library
- **Persistence**: SQLite WAL for PM memory. No ORM — raw SQL with minimal helpers
- **DI**: Manual constructor injection. No uber-go/fx in production (per current go.mod)
- **Module system**: `apps/api/` (entry), `modules/` (bounded contexts), `shared/` (kernel + infra)

---

## 1. Architecture Rules (Non-Negotiable)

### 1.1 Module Structure — Every New Module Must Follow This Layout

```
modules/<name>/
  domain/           # Entities, value objects, interfaces (ZERO external deps)
  application/
    port/           # Inbound interface (use case contract)
    service/        # Application service implementing port
  infrastructure/   # Adapters: DB, HTTP, cache implementations
  interfaces/       # Delivery mechanisms: mcp/, http/, grpc/
  test/             # Integration tests (requires real infra)
```

**WHY**: Ports/Adapters separates "what to do" from "how to do it." Domain layer has ZERO external dependencies. This enables testing and technology swaps without touching business logic.

**VIOLATION**: If `domain/` imports anything outside `context`, `time`, `fmt`, `errors`, `math`, or stdlib — it's wrong.

### 1.2 Domain Layer Is Sacred

- Domain packages (`modules/*/domain/`) must NOT import: infrastructure, interfaces, application, config, or any third-party library
- Domain only imports: stdlib (`context`, `time`, `fmt`, `errors`, `math`)
- Cross-module domain imports are allowed only via `shared/kernel/` types

### 1.3 Interface Location — Consumer-Side Only

Define interfaces WHERE THEY ARE CONSUMED, not where implemented.

```go
// GOOD — defined in application/port, consumed by application/service
// modules/jira/application/port/port.go
type Inbound interface {
    SearchIssues(ctx context.Context, jql string, maxResults int) (*domain.SearchResult, error)
}

// BAD — interface in domain or infrastructure
```

**WHY**: Consumer defines only what it needs. Provider doesn't know about callers. Enables mocking.

### 1.4 No Circular Dependencies Between Modules

Modules may communicate only via:
1. `shared/kernel/` types (domain events, shared entities)
2. Port interfaces defined by the consumer module
3. Domain events published through `event.Bus`

### 1.5 Config Package — Single Source

`shared/infrastructure/config/config.go` is the canonical config package.

- Do NOT create or import `config/` at project root for new code
- All config structs have explicit field tags and env var comments
- Load function returns `(*Config, error)` — never panics

---

## 2. Go Code Rules (Uber Style Guide Enforced)

### 2.1 Compile-Time Interface Checks — MANDATORY

Every concrete type implementing a port interface MUST have this line:

```go
var _ port.Inbound = (*JiraService)(nil)  // for pointer receivers
var _ port.Inbound = (*Score)(0)          // for value receivers (rare)
```

**WHY**: Catch interface breakage at compile time, not runtime. Zero cost.

**EXISTING PATTERN**: `modules/jira/application/service/jira_service.go` already does this.

### 2.2 Error Handling — Operation Context, Not "Failed To"

```go
// GOOD
return fmt.Errorf("search issues: %w", err)

// BAD
return fmt.Errorf("failed to search issues: %w", err)
```

**Stack reads cleanly**: `"search issues: connection refused"` vs `"failed to search issues: failed to connect: connection refused"`

### 2.3 Error Types — Domain-Layer Convention

- Exported error types: suffix `Error` (e.g., `ErrJiraAPI`, `ErrNotFound`)
- Unexported error variables: prefix `err` (e.g., `errNotFound`)
- Use `errors.As[T]` for type matching, `errors.Is` for value matching
- Wrap errors with `%w` when callers need to inspect
- Never `log AND return` — pick one layer to handle

**EXISTING PATTERN**: `modules/jira/domain/errors.go` and `shared/kernel/errors/errors.go`

### 2.4 Return Concrete Types, Accept Interfaces

```go
// GOOD
func NewJiraService(client domain.Client, cache port.Cache) *JiraService {
    return &JiraService{client: client, cache: cache}
}

// BAD — hides the concrete type, prevents callers from using specific methods
func NewJiraService(client domain.Client, cache port.Cache) port.Inbound {
```

### 2.5 Mutex Must Be Private Field — Never Embedded

```go
// GOOD
type simpleCache struct {
    mu   sync.RWMutex
    data map[string][]byte
}

// BAD — Lock/Unlock become public API
type simpleCache struct {
    sync.RWMutex
    data map[string][]byte
}
```

### 2.6 No init() Functions

Use explicit constructors called from `main()` or `bootstrap.go`. `init()` runs uncontrollably, can't return errors, and hides dependencies.

**EXCEPTION**: `database/sql` driver registration and `encoding` type registration only.

### 2.7 Exit Only in main()

```go
// main.go
func main() {
    if err := bootstrap.Run(); err != nil {
        log.Fatalf("server error: %v", err) // only here
    }
}
```

All other functions return `error`. `log.Fatal` and `os.Exit` forbidden outside `main()`.

### 2.8 Never Panic — Except in main/Tests

Return errors for all recoverable situations. Panic kills the process. In MCP server serving multiple clients, one bad request kills everything.

**Acceptable panics**: `template.Must()` at init, `regexp.MustCompile()` for package-level regex.

### 2.9 No Fire-and-Forget Goroutines

Every goroutine must have:
1. A way to stop (context cancellation, done channel)
2. A way for the caller to wait (WaitGroup, done channel)

```go
// GOOD
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
go func() {
    select {
    case <-ctx.Done():
        return
    case <-ticker.C:
        flush()
    }
}()
```

### 2.10 Pre-allocate Slices and Maps

```go
// When size is known
items := make([]Item, 0, len(input))

// Map with expected size
m := make(map[string]int, len(keys))
```

**WHY**: Avoids rehashing and re-allocations on append.

### 2.11 Use time.Duration, Not int

```go
// GOOD
func poll(delay time.Duration) { time.Sleep(delay) }

// BAD — what unit is 10?
func poll(delay int) { time.Sleep(time.Duration(delay) * time.Millisecond) }
```

For JSON/config fields where `time.Duration` can't be used directly, include unit in field name: `TimeoutMillis`, `IntervalSec`.

---

## 3. MCP Handler Rules

### 3.1 One Tool, One Handler Function

Each MCP tool maps to exactly ONE handler function. Handler functions:
- Live in `modules/<name>/interfaces/mcp/handlers.go`
- Have signature: `func (h *Handlers) ToolName(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)`
- Return ONLY `*mcp.CallToolResult` — never raw errors to MCP client

### 3.2 Register Tools in apps/api/internal/mcp/

Tool registration lives in `apps/api/internal/mcp/<module>.go`:
```go
func RegisterJiraTools(s *server.MCPServer, h *jmcp.Handlers) {
    s.AddTool(
        mcp.NewTool("jira_search",
            mcp.WithDescription("..."),
            mcp.WithString("jql", mcp.Required(), mcp.Description("...")),
        ),
        h.SearchIssues,
    )
}
```

### 3.3 Use mcputil for Error Handling

All MCP tool errors go through `shared/infrastructure/mcputil/helpers.go`:

```go
// Validation error
return mcputil.ErrInvalid("jql parameter is required"), nil

// Jira API error (auto-classifies auth/rate-limit/not-found/network)
return mcputil.ErrJira("search issues", err), nil

// Internal error (no details leaked to user)
return mcputil.ErrInternal("get board", err), nil
```

### 3.4 Never Return Raw Errors to MCP Client

```go
// GOOD
result, err := h.Jira.SearchIssues(ctx, jql, int(maxResults))
if err != nil {
    return mcputil.ErrJira("search issues", err), nil
}
return mcputil.TextResult(formatResults(result)), nil

// BAD — exposes raw error to AI client
if err != nil {
    return nil, err
}
```

### 3.5 Tool Description — Be Specific

Description should answer: "What does this do AND when should I use it?"

```go
mcp.WithDescription("Get full details of a Jira issue including description, comments, and linked issues. Use when you need complete issue context.")
```

---

## 4. Domain Model Rules

### 4.1 Entities vs Value Objects

- **Entities**: Have identity (`ID` field), mutable, represent things that change over time (e.g., `SprintSnapshot`, `Risk`, `Decision`)
- **Value Objects**: No identity, immutable, defined by attributes (e.g., `Score`, `Velocity`, `Trend`)

### 4.2 Domain Methods — Behavior on Entities

```go
// GOOD — behavior belongs on the entity
func (s *SprintSnapshot) IsZombie() bool {
    return s.CarryoverRate() > 30
}

// BAD — behavior in service that doesn't need the entity
func isZombie(carryover, total int) bool {
    return float64(carryover)/float64(total) > 0.30
}
```

### 4.3 Domain Events — Cross-Module Communication

Define events in `shared/kernel/event/events.go`. Module-specific events live in module's domain layer.

```go
// Event naming: domain.entity.verb_past_tense
type HealthScoreComputed struct {
    BoardID    int
    SprintName string
    Score      int
}
func (e HealthScoreComputed) EventName() string { return "health_score.computed" }
```

### 4.4 No God-Types

Single responsibility: if a struct has 10+ fields covering 3+ unrelated concerns, split it.

**EXISTING PATTERN**: `shared/domain/jira/domain.go` — Issue struct is acceptable because all fields relate to a single Jira issue.

---

## 5. Testing Rules

### 5.1 Table-Driven Tests — Default Pattern

```go
func TestScore_Grade(t *testing.T) {
    tests := []struct {
        name  string
        score Score
        want  string
    }{
        {name: "healthy", score: 85, want: "A"},
        {name: "fair", score: 65, want: "B"},
        {name: "at_risk", score: 45, want: "C"},
        {name: "critical", score: 20, want: "D"},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := tt.score.Grade(); got != tt.want {
                t.Errorf("Grade() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 5.2 Race Detector — Always

```bash
go test -race ./...
```

No exceptions. If code fails `-race`, it's broken.

### 5.3 Test Location

- Unit tests: same package as code (e.g., `modules/sprint/domain/score_test.go`)
- Integration tests: `modules/<name>/test/` (requires real infrastructure)
- Shared tests: `shared/infrastructure/<name>/*_test.go`

### 5.4 No Sleep in Tests

```go
// BAD
time.Sleep(100 * time.Millisecond)
assert(t, condition)

// GOOD — poll with timeout
deadline := time.Now().Add(5 * time.Second)
for time.Now().Before(deadline) {
    if condition {
        break
    }
    time.Sleep(10 * time.Millisecond)
}
assert(t, condition)
```

---

## 6. Naming Rules

### 6.1 Functions

- Getters: `Name()` not `GetName()` (Go convention)
- Setters: `SetName(n)` (setter IS prefixed with Set)
- Predicates: `IsDone()`, `HasChildren()`, `CanRetry()`
- Constructors: `New()` or `NewTypeName()`

### 6.2 Interfaces

- One-method interfaces preferred (e.g., `io.Reader`, `io.Closer`)
- Multi-method interfaces acceptable for bounded contexts (e.g., `domain.Client`)
- Suffix with `-er` for single-method: `Notifier`, `Reader`, `Writer`

### 6.3 Package Names

- Short, lowercase, no underscores: `jira`, `sprint`, `mcputil`
- Avoid: `util`, `common`, `shared` (use descriptive names instead)
- Avoid: package name same as common stdlib (`context`, `http`, `errors`)

### 6.4 Constants

```go
// Group related constants
const (
    ScoreHealthy Score = 80
    ScoreFair    Score = 60
    ScoreAtRisk  Score = 40
    ScoreMin     Score = 0
    ScoreMax     Score = 100
)

// Enum starting at 1 (zero = unset)
type Severity int
const (
    SeverityLow Severity = iota + 1
    SeverityMedium
    SeverityHigh
    SeverityCritical
)
```

---

## 7. Import Rules

### 7.1 Import Groups — Strict

```go
import (
    // 1. Standard library
    "context"
    "fmt"
    "time"

    // 2. External packages
    "github.com/mark3labs/mcp-go/mcp"
    "go.uber.org/zap"

    // 3. Internal packages (this project)
    "github.com/aldok10/zara-jira-mcp/modules/jira/domain"
    "github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)
```

### 7.2 Import Aliases — Only on Conflict

```go
// GOOD — alias needed for domain conflict
import (
    domain "github.com/aldok10/zara-jira-mcp/modules/jira/domain"
)

// BAD — alias for convenience
import (
    mcputil "github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)
```

---

## 8. Performance Rules

### 8.1 strconv Over fmt for Number Conversion

```go
// GOOD
s := strconv.Itoa(n)

// BAD — 5-10x slower, allocates
s := fmt.Sprintf("%d", n)
```

### 8.2 strings.Builder with Grow()

```go
var sb strings.Builder
sb.Grow(estimatedSize) // if you know approximate size
for _, item := range items {
    sb.WriteString(item.Name)
}
result := sb.String()
```

### 8.3 Copy Slices/Maps at Boundaries

```go
// When receiving from caller
func (s *Store) SetItems(items []Item) {
    s.items = make([]Item, len(items))
    copy(s.items, items)
}

// When returning to caller
func (s *Store) Items() []Item {
    result := make([]Item, len(s.items))
    copy(result, s.items)
    return result
}
```

### 8.4 Pre-size Maps When Size Known

```go
m := make(map[string]int, len(keys))
```

---

## 9. Anti-Patterns to Avoid

### 9.1 Zombie Sprint Pattern

When carryover exceeds 30% of sprint items — this is a **zombie sprint**. The system should detect and flag this automatically.

**Detection**: `SprintSnapshot.IsZombie()` — already implemented in `modules/sprint/domain/memory/models.go`

### 9.2 Hero Culture

When one person handles >40% of sprint items — unhealthy dependency.

**Detection**: Track via `TeamMetric` and flag when `IssuesAssigned / TotalSprintIssues > 0.4`

### 9.3 Dead Retro

Retrospective without action items, or action items with status "pending" for >2 sprints.

### 9.4 Rubber-Stamp DoD

DoD checklist exists but completion rate <70% of defined items.

### 9.5 No Sprint Goal

Sprint started without explicit goal set via `pm_set_sprint_goal`.

### 9.6 Config Duplication

Do NOT create duplicate config packages. Use `shared/infrastructure/config` everywhere.

### 9.7 Interface in Wrong Location

Do NOT define interfaces in `domain/` for infrastructure implementations. Define in consumer (application/port/).

---

## 10. Security Rules

### 10.1 Secrets Never in Code

- All secrets via environment variables only
- Never log tokens, API keys, or passwords
- Use `crypto/subtle` for constant-time comparison of secrets

### 10.2 Input Validation at Handler Boundary

```go
func (h *Handlers) SearchIssues(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    jql, err := req.RequireString("jql")
    if err != nil {
        return mcputil.ErrInvalid("jql parameter is required"), nil
    }
    // Use mcputil.IsSafeJQLValue() for JQL injection prevention
}
```

### 10.3 SQL Injection Prevention

Use parameterized queries. Never string-interpolate user input into SQL.

### 10.4 JQL Injection Prevention

Use `mcputil.IsSafeJQLValue()` to validate JQL values before interpolation.

---

## 11. Logging Rules

### 11.1 Use Structured Logging

```go
import "log/slog"

// GOOD
slog.Info("server started", "port", cfg.Server.Port, "version", "0.4.0")
slog.Error("search failed", "jql", jql, "error", err)

// BAD
log.Printf("Server started on port %s", cfg.Server.Port)
```

### 11.2 Log at Appropriate Level

- `DEBUG`: Values, traces (disabled in production)
- `INFO`: Lifecycle events (started, stopped, connected)
- `WARN`: Recoverable issues (validation failed, retrying)
- `ERROR`: Unrecoverable in this handler (API failed, DB error)

### 11.3 Never Log Secrets

```go
// GOOD
slog.Info("Jira connected", "base_url", cfg.Jira.BaseURL)

// BAD — leaks token
slog.Info("Jira connected", "token", cfg.Jira.Token)
```

---

## 12. Documentation Rules

### 12.1 Exported Functions — Always Document

```go
// CalculateHealth computes a 0-100 health score for the active sprint.
// Returns an error if no active sprint or insufficient data.
func (s *sprintService) CalculateHealth(ctx context.Context, boardID int) (*port.HealthResult, error) {
```

### 12.2 Package Documentation

Every package must have a `doc.go` or comment at top of first file:

```go
// Package sprint implements sprint planning and execution intelligence.
// It provides velocity tracking, forecasting, health scoring, and anti-pattern detection.
package sprint
```

### 12.3 TODO Convention

```go
// TODO(aldo): Implement caching for Jira search results
// TODO(#42): Add retry logic for transient failures
```

---

## 13. Build & CI Rules

### 13.1 Makefile Commands

```bash
make build       # Build binary to bin/
make test        # Run all tests
make lint        # Run golangci-lint
make tidy        # Tidy all go.mod files
```

### 13.2 Linter Compliance

Project uses `.golangci.yml` with:
- `go vet`, `staticcheck`, `gosec`, `bodyclose`, `noctx`, `sqlclosecheck`
- `revive`, `gofumpt`, `misspell`
- `prealloc`, `gocritic`
- `errcheck`, `ineffassign`, `unparam`, `unused`

**Run `make lint` before committing.**

### 13.3 go.mod Tidy

After changing imports, run `make tidy` to update all 5 `go.mod` files:
- Root `go.mod`
- `shared/go.mod`
- `modules/jira/go.mod`
- `modules/sprint/go.mod`
- `apps/api/go.mod`

---

## 14. MCP Transport Rules

### 14.1 Stdio Is Default

The server starts with `server.ServeStdio(s)`. This is the primary transport for AI clients.

### 14.2 Tool Registration Pattern

```go
// apps/api/internal/mcp/jira.go
func RegisterJiraTools(s *server.MCPServer, h *jmcp.Handlers) {
    s.AddTool(
        mcp.NewTool("jira_search", ...),
        h.SearchIssues,
    )
}
```

### 14.3 Tool Naming Convention

Format: `<module>_<action>` or `pm_<action>`

- `jira_search`, `jira_get_issue`, `jira_create_issue`
- `pm_health`, `pm_forecast`, `pm_record_risk`

---

## 15. Performance Budget

### 15.1 MCP Tool Response Time

- P50 < 200ms
- P99 < 2s
- Cache Jira responses: 60s TTL for read operations

### 15.2 Binary Size

- Target: <15MB (Go static binary)
- Build flags: `-ldflags="-s -w"`

### 15.3 Memory

- SQLite WAL mode for concurrent reads
- No unbounded goroutines
- No unbounded caches

---

## Checklist for New Code

Before committing, verify:

**Architecture & Structure**
- [ ] Follow Ports & Adapters architecture
- [ ] Domain layer has ZERO external dependencies
- [ ] Interfaces defined in application/port, NOT in domain
- [ ] No circular dependencies between modules
- [ ] Use `shared/infrastructure/config/config.go` for all config

**Go Code Quality**
- [ ] No `init()` functions
- [ ] No `panic()` in non-test code
- [ ] Error messages use operation context, NOT "failed to"
- [ ] Error types follow `ErrXxx` (exported) or `errXxx` (unexported) naming
- [ ] Compile-time interface checks: `var _ port.Inbound = (*Type)(nil)`
- [ ] Mutex is private field (`mu sync.Mutex`), never embedded
- [ ] Slices pre-allocated with `make([]T, 0, n)` when size known
- [ ] Use `strconv` instead of `fmt.Sprintf` for number conversion
- [ ] Structured logging with `slog`
- [ ] No secrets in code or logs
- [ ] Use `mcputil` helpers for MCP error responses
- [ ] `defer` used for resource cleanup (locks, close, etc.)

**Import Rules**
- [ ] Import groups: stdlib / external / internal
- [ ] No import aliasing unless required for conflicts
- [ ] Use `_prefix` for unexported package-level variables

**Function & Method Naming**
- [ ] Getters: `Name()` not `GetName()`
- [ ] Predicates: `IsDone()`, `HasChildren()`, `CanRetry()`
- [ ] Constructors: `New()` or `NewTypeName()`

**Testing**
- [ ] Table-driven tests with `t.Run()`
- [ ] Pass `-race` flag in CI
- [ ] No `time.Sleep()` in tests
- [ ] Proper cleanup in tests (`t.Cleanup()`)

**Build & CI**
- [ ] Run `golangci-lint` (`make lint`)
- [ ] Run `make tidy` if imports changed
- [ ] Write unit test when adding logic

**MCP Specific**
- [ ] One handler function per tool
- [ ] Tool registration in `apps/api/internal/mcp/`
- [ ] Use `mcputil.ErrInvalid`, `mcputil.ErrJira`, `mcputil.ErrAI` for errors

**Performance**
- [ ] Channel buffers size 0 or 1
- [ ] Pre-allocate maps when size known: `make(map[K]V, n)`
- [ ] Be aware of copy-on-write with slices

**Enum Rules**
- [ ] Enums start at 1 (`iota + 1`)
- [ ] Define constants in groups in parentheses

**String Handling**
- [ ] Use `strings.Builder` with `Grow()` for concatenation
- [ ] Use `strconv.AppendInt()` for efficient number-to-string conversion

**Debuggability**
- [ ] Sensitive fields have JSON tags (`json:"omitempty"`)
- [ ] Field names descriptive, avoid single-letter variables

**Documentation**
- [ ] Document exported functions
- [ ] Package docstrings at top of first file

---

**Anti-Patterns to Check For**

1. **Zombie Sprint Pattern** — Carryover >30% (already implemented in models)
2. **Hero Culture** — One person >40% of sprint items
3. **Dead Retro** — No action items or "pending" status >2 sprints
4. **Rubber-Stamp DoD** — DoD defined but completion rate <70%
5. **God-Types** — Structs with 10+ fields covering unrelated concerns
6. **Fire-and-Forget Goroutines** — Missing context cancellation or WaitGroup
7. **No Sprint Goal** — Sprint without goal set via `pm_set_sprint_goal`
8. **Config Duplication** — Using config outside `shared/infrastructure/config/`
9. **Interface in Wrong Location** — Defined in domain instead of application/port

---

*Last updated: 2026-06-28*
*Based on: Uber Go Style Guide (40 rules), 100 Go Mistakes, Hexagonal Architecture, MCP best practices*

---

**Project-Specific Notes**

The project currently **deviates** from the strict architectural recommendations:

1. Module structure uses **uber-go/fx DI** (not manual constructors)
2. Interface checks are **NOT** present in existing code - must be added for new code
3. Some MCP tools are incomplete (has TODO: "not implemented")
4. Error handling partially follows mcputil pattern
5. Tests exist only in `shared/infrastructure/validate/`
6. String building uses `fmt.Sprintf` in handlers (could use `strings.Builder`)

**Recommended Changes for Compliance:**

1. Add interface compliance checks to all service methods
2. Add missing implementation for TODO: "not implemented" services
3. Update existing MCP handlers to follow `mcputil` error patterns
4. Implement interface compliance for all existing service implementations
5. Add the missing `AddWatcher`, `GetWatchers`, `AddLabel`, `GetProjects` methods to sprint service
