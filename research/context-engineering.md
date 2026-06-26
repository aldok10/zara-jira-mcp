# Context Engineering for MCP Tool Servers

Research date: 2026-06-26

## The Problem We Already Solved (partially)

> "With 50+ tools, definitions alone consume 70,000-130,000 tokens per request."
> — tianpan.co, 2026

> "Every MCP server you connect injects hundreds of tool definitions into the context window. Once those definitions consume enough tokens, the agent starts hallucinating."
> — albato.com, 2026

Our 231 tools = ~50,000+ tokens of schema overhead. This is why ChatGPT Desktop lags.

### Current Solution: PM_PROFILE
- `unified` (1 tool) = ~200 tokens
- `smart` (7 tools) = ~1,400 tokens
- `lite` (30 tools) = ~6,000 tokens
- `standard` (80 tools) = ~16,000 tokens
- `full` (150 tools) = ~30,000 tokens
- `all` (231 tools) = ~50,000 tokens

This is GOOD — but static. The ideal: **dynamic tool selection at runtime**.

---

## Industry Patterns (2026)

### 1. Intent-Based Dynamic Tool Selection (Lunar.dev)
- Agent first classifies intent
- Only relevant tools loaded based on intent
- Result: 5-10 tools per request instead of 200+

### 2. Tool Groups / Namespaces
- Group tools into domains
- Agent loads one group at a time
- Our module system already does this

### 3. Code Mode (LeanIX)
- Instead of 50 tool schemas, give agent a TYPE DEFINITION
- Agent writes code that calls the API
- Massively reduces tokens

### 4. Hierarchical Tool Discovery
- Level 1: "I have these capabilities: jira, pm, notifications, reports"
- Level 2: "Within PM, I can: track sprint, manage risks, forecast"
- Level 3: "For sprint tracking: snapshot, health, burndown, goals"
- Agent drills down as needed

### 5. Context Offloading (Karpathy's approach)
- Don't put everything in context window
- Store in external memory, retrieve on demand
- Our SQLite memory already does this for DATA
- But TOOL SCHEMAS are still all in context

---

## What This Means for Our Project

### Already Done:
- Profile system (static tool reduction)
- Module toggles (PM_ENABLED_MODULES)
- Unified/Smart tools (router pattern)

### Not Yet Done (Roadmap):

#### Dynamic Tool Discovery
Agent asks "what can you do about X?" → server returns only relevant tools.
MCP protocol doesn't natively support this, but we can:
1. Use `pm_help topic=X` as discovery tool (already exists!)
2. Add `pm_load_module module=risks` that dynamically registers tools mid-session
3. Server-side: track which tools are actually called, auto-suggest optimal profile

#### Contextual Tool Descriptions
Instead of full descriptions for all tools, use ULTRA SHORT descriptions when in large profile, detailed descriptions only in small profiles.

#### Usage Analytics
Track which tools each PM actually uses → auto-recommend their ideal profile.

---

## Implementation Roadmap

### Phase A: Optimize Current (this week)
- [ ] Shorten all tool descriptions to <50 chars in `all` profile
- [ ] Add `pm_load_module` tool for dynamic tool loading
- [ ] Track tool usage to SQLite (which tools called, frequency)
- [ ] Auto-suggest profile based on usage patterns

### Phase B: Dynamic Discovery (next sprint)
- [ ] Implement hierarchical tool listing (categories → tools)
- [ ] `pm_discover action=X` returns only relevant tools
- [ ] Lazy tool registration (tools registered on first category access)

### Phase C: Context-Aware Routing (future)
- [ ] AI-powered intent classification before tool selection
- [ ] Auto-compress tool descriptions based on conversation history
- [ ] Cross-session tool preference learning

---

## Key Takeaway

> "From Collector to Curator: The Strategic Shift" — stork.ai

The project has 231 tools. That's the COLLECTOR phase (done).
Now we need the CURATOR phase: **serve the right 5-10 tools per conversation, not all 231.**

Profile system is step 1. Dynamic discovery is step 2. Intent-based selection is step 3.
