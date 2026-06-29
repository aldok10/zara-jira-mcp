# zara-jira-mcp — System Guidelines

## Project Overview

**AI-powered Scrum Master MCP server** with **~150 focused tools** across 17 categories. Operates as AI assistant for Jira project management with Monte Carlo forecasting, anti-pattern detection, persistent memory, and multi-channel notifications.

**Architecture**: Modular Ports & Adapters (DDD), uber-go/fx DI, SQLite memory, MCP stdio transport.

---

## Core Patterns

### 1. Module Structure (Primary)

Each bounded context follows this exact layout:

```
modules/<name>/
  domain/           # Entities, value objects, interfaces (NO external deps)
    structs.go      # Data types
    methods.go      # Domain methods
  application/
    port/           # Inbound interface (consumer-defined)
      port.go       # struct Inbound interface
    service/        # Implement port
      service.go    # Business logic
  infrastructure/   # Adapters: DB, HTTP, cache
    client/        # HTTP clients
    store/         # Database access
  interfaces/       # Delivery: MCP, REST, gRPC
    mcp/           # MCP handler registration
  test/             # Integration tests
```

**Rule**: Domain layer imports ONLY stdlib (`context`, `time`, `fmt`, `errors`, `math`)

### 2. Interface Definition Convention

```go
// modules/jira/application/port/port.go
package port

// Inbound defines the Jira use cases exposed by this module.
type Inbound interface {
	SearchIssues(ctx context.Context, jql string, maxResults int) (*domain.SearchResult, error)
	GetIssue(ctx context.Context, key string) (*domain.Issue, error)
	// ... other methods
}
```

**Rule**: Define interfaces at CONSUMER site, not where implemented

### 3. Implementation Pattern

```go
// modules/jira/application/service/jira_service.go
package service

import (
	"context"
	"github.com/aldok10/zara-jira-mcp/modules/jira/application/port"
	"github.com/aldok10/zara-jira-mcp/modules/jira/domain"
)

// Compile-time safety
var _ port.Inbound = (*JiraService)(nil)

// Concrete implementation
type JiraService struct {
	client domain.Client
	cache  port.Cache
}

// Constructor
type JiraService struct {
	client domain.Client
	cache  port.Cache
}

func NewJiraService(client domain.Client, cache port.Cache) *JiraService {
	return &JiraService{client: client, cache: cache}
}
```

---

## Go Code Rules

### Compile-Time Interface Checks (MANDATORY)

Every port implementation MUST have:

```go
var _ port.Inbound = (*Type)(nil)
```

**Why**: Catches interface breakage at compile time, not runtime.

**Location**: Every application/service directory has this line.

### Error Handling Convention

```go
// GOOD - Operation context
return fmt.Errorf("search issues: %w", err)

// BAD - Redundant "failed to"
return fmt.Errorf("failed to search issues: %w", err)
```

**Pattern**: Errors include what was attempted, not what failed.

### Error Types Naming

```go
// Shared errors package
package domain

type ErrJiraAPI struct {
	StatusCode int
	Message    string
	Endpoint   string
}

func (e *ErrJiraAPI) Error() string {
	return fmt.Sprintf("jira api error %d on %s: %s", e.StatusCode, e.Endpoint, e.Message)
}
```

**Rule**: - Exported: `ErrXxx` - Unexported: `errXxx`

### Return Convention

```go
// Return concrete types, accept interfaces
type JiraService struct {
	client domain.Client
	cache  port.Cache
}

func (j *JiraService) SearchIssues(ctx context.Context, jql string, maxResults int) (*domain.SearchResult, error) {
	return j.client.SearchIssues(ctx, jql, maxResults, 0)
}
```

### Mutex Pattern

```go
// GOOD - Private field, exported methods
type simpleCache struct {
	mu   sync.RWMutex        // private, exported operations only
	data map[string][]byte
}

func (c *simpleCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.data[key], true
}
```

**Rule**: Mutex always a private unexported field in struct.

---

## File Organization

### Files Structure

```
apps/api/
  internal/
    bootstrap/           # Entry point
    mcp/                 # Tool registration
  cmd/
    server/
      main.go          # Binary entry

modules/<name>/         # Each bounded context
  domain/               # ENTITIES ONLY
  application/          # USE CASES
  infrastructure/       # DRIVERS
  interfaces/           # GATEWAYS
  test/                 # INTEGRATION TESTS

shared/                 # Cross-cutting concerns
  kernel/              # Base types, events
  infrastructure/      # Config, validation
```

### Package Naming

- Lowercase, descriptive, no underscores (except separators)
- No naming conflicts with stdlib (`context`, `http`, `errors`)
- Each module has distinct, clear name

---

## Specific Module Examples

### Jira Module (`modules/jira/`)

**Domain Files** (`modules/jira/domain/`):
- `types.go` - Issue, Sprint, Board, SearchResult
- `client.go` - Jira API client interface  
- `errors.go` - Jira-specific error types

**Application Files** (`modules/jira/application/`):
- `port/port.go` - Inbound interface definition
- `service/jira_service.go` - Jira operations implementation

**Infrastructure Files** (`modules/jira/infrastructure/`):
- `client/rest.go` - HTTP REST implementation

**Interface Files** (`modules/jira/interfaces/mcp/`):
- `handlers.go` - MCP handler implementation

### Sprint Module (`modules/sprint/`)

**Domain Files** (`modules/sprint/domain/`):
- `sprint/velocity.go` - Velocity calculations
- `sprint/domain/predictability.go` - Predictability analysis
- `memory/` - Retrospective, risks, decisions storage

**Application Files** (`modules/sprint/application/`):
- `sprint/application/port/` - Sprint interfaces
- `sprint/application/service/` - Sprint services

**Infrastructure Files** (`modules/sprint/infrastructure/`):
- `sprintstore/` - SQLite persistence

---

## Testing Patterns

### Table-Driven Tests (DEFAULT)

```go
func TestScore_Grade(t *testing.T) {
	tests := []struct {
		name string
		score Score
		want string
	}{
		{name: "healthy", score: 85, want: "A"},
		{name: "fair", score: 65, want: "B"},
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

### Test Location Rule

- **Unit tests**: Same package as code
- **Integration tests**: `/test/` subdirectory (requires real infra)
- **Shared tests**: `shared/infrastructure/<name>/*_test.go`

---

## Import Rules

### Import Groups

```go
import (
	// 1. Standard library
	"context"
	"fmt"
	"time"

	// 2. External packages
	"github.com/mark3labs/mcp-go/mcp"
	"go.uber.org/zap"

	// 3. Internal packages
	"github.com/aldok10/zara-jira-mcp/modules/jira/domain"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)
```

### Import Aliases

```go
// ONLY for conflicts
import (
	domain "github.com/aldok10/zara-jira-mcp/modules/jira/domain"
)

// NO aliases for convenience
import (
	mcputil "github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil" // WRONG
)
```

---

## Naming Conventions

### Function Names

- **Getters**: `Name()` not `GetName()`
- **Setters**: `SetName(n)` (setter IS prefixed)
- **Predicates**: `IsDone()`, `HasChildren()`, `CanRetry()`
- **Constructors**: `New()` or `NewTypeName()`

### Constant Definitions

```go
const (
	ScoreHealthy Score = 80
	ScoreFair    Score = 60
	ScoreAtRisk  Score = 40
	ScoreMin     Score = 0
	ScoreMax     Score = 100
)
```

**Rule**: Enums start at 1 (`iota + 1`), zero means "unset".

### Variable Scope

```go
// GOOD - Scope limited to where needed
if err := validate(input); err != nil {
	return err
}
```

---

## MCP Handler Rules

### Handler Pattern

```go
// modules/<name>/interfaces/mcp/handlers.go
package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/aldok10/zara-jira-mcp/modules/<name>/application/port"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)

// Handler struct
type Handlers struct {
	<name> port.Inbound
}

// Handler implementation
func (h *<Name>Handlers) <ToolName>(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Parse parameters
	param, err := req.RequireString("param")
	if err != nil {
		return mcputil.ErrInvalid("param is required"), nil
	}

	// Delegate to service
	result, err := h.<Name>.DoSomething(ctx, param)
	if err != nil {
		return mcputil.ErrJira("do something", err), nil
	}

	// Format result
	return mcputitextResult(formatResult(result)), nil
}
```

### Tool Registration

```go
// apps/api/internal/mcp/<name>.go
package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

jmcp "github.com/aldok10/zara-jira-mcp/modules/<name>/interfaces/mcp"
)

func Register<name>Tools(s *server.MCPServer, h *jmcp.Handlers) {
	s.AddTool(
		mcp.NewTool("<tool_name>",
			mcp.WithDescription("<tool description>"),
			mcp.WithString("param", mcp.Required(), mcp.Description("<param description>")),
		),
		h.<ToolName>,
	)
}
```

---

## Production Rules

### No init() Functions

```go
// BAD - Hidden initialization, can't test, no error handling
func init() {
	raw, _ := os.ReadFile("config.yaml")
	yaml.Unmarshal(raw, &config)
}
```

**Rule**: Use explicit functions called from main() or bootstrap().

### Exit Only in main()

```go
// main.go
func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatalf("server error: %v", err) // ONLY here
	}
}
```

### No Panic in Production

```go
// BAD - Crashes entire process
func parseConfig(path string) Config {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err) // CRASHES SERVER
	}
	// ...
}
```

```go
// GOOD - Return error
func parseConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	// ...
}
```

### No Fire-and-Forget Goroutines

```go
// BAD - No lifecycle, resource leak
func StartBackground() {
	go func() {
		for {
			flush()
			time.Sleep(time.Second)
		}
	}()
}
```

```go
// GOOD - Has stop signal, caller can wait
func StartBackground(ctx context.Context) {
	stop := make(chan struct{})

go func() {
	defer close(stop)
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
			flush()
		}
	}
}()

	// Caller can wait
	<-stop
}
```

---

## Performance Guidelines

### Number Conversion

```go
// GOOD - Fast
import "strconv"
s := strconv.Itoa(n)

// BAD - Slow, reflection
s := fmt.Sprintf("%d", n)
```

### String Building

```go
var sb strings.Builder
sb.Grow(estimatedSize) // Pre-allocate
for _, item := range items {
	sb.WriteString(item.Name)
}
result := sb.String()
```

### Slice Initialization

```go
// Pre-allocate when size known
items := make([]Item, 0, len(input))

// Return copy to prevent mutation
func (s *Store) Items() []Item {
	result := make([]Item, len(s.items))
	copy(result, s.items)
	return result
}
```

---

## Configuration Rules

### Central Config Package

**Single source**: `shared/infrastructure/config/config.go`

```go
// ALL imports use this config package
import (
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
)
```

### Config Loading

```go
// NEVER config.go at project root for new code
func Load() (*Config, error) {
	// Environment variable loading
	// Validation
	// Defaults
}
```

---

## Documentation Rules

### Function Documentation

```go
// CalculateHealth computes a 0-100 health score for the active sprint.
// Returns an error if no active sprint or insufficient data.
func (s *sprintService) CalculateHealth(ctx context.Context, boardID int) (*port.HealthResult, error) {
```

### Package Documentation

```go
// Package sprint implements sprint planning and execution intelligence.
// It provides velocity tracking, forecasting, health scoring, and anti-pattern detection.
package sprint
```

---

## Common Antipatterns to Avoid

### Zombie Sprint Pattern

When carryover exceeds 30% (detected by `SprintSnapshot.IsZombie()`)

### Hero Culture

When one person handles >40% of sprint items (detect via `TeamMetric`)

### Dead Retro

Retro without action items, or items pending >2 sprints

### Rubber-Stamp DoD  

DoD exists but completion rate <70%

### No Sprint Goal

Sprint without goal set via `pm_set_sprint_goal`

---

## Build & CI Commands

```bash
make build      # Build binary to bin/
make test       # Run all tests
make lint       # Run golangci-lint
make tidy       # Tidy all go.mod files
```

### Go Mod Updates

```bash
# After import changes
make tidy

# Updates these go.mod files:
# - go.mod (root)
# - shared/go.mod
# - modules/jira/go.mod
# - modules/sprint/go.mod
# - apps/api/go.mod
```

---

## Checklist for New Code

Before committing, verify:

**Architecture**
- [ ] Follow module structure exactly
- [ ] Domain layer has NO external dependencies
- [ ] Interfaces defined in application/port
- [ ] Interface compliance check present

**Code Quality**
- [ ] Compile-time interface check: `var _ Interface = (*Type)(nil)`
- [ ] Error messages use operation context
- [ ] Errors follow `ErrXxx` / `errXxx` naming
- [ ] Mutex is private field, never embedded
- [ ] Slices pre-allocated when size known
- [ ] Use `strconv` instead of `fmt.Sprintf` for numbers
- [ ] Structured logging with `slog`

**Testing**
- [ ] Table-driven tests with `t.Run()`
- [ ] Pass `go test -race` in CI
- [ ] No `time.Sleep()` in tests

**Build**
- [ ] Run `golangci-lint` (`make lint`)
- [ ] Run `make tidy` if imports changed
- [ ] Write unit test when adding logic

---

**Summary**: Follow these rules exactly for consistency and maintainability.
Every deviation requires documented justification in the PR.