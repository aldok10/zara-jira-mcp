# zara-jira-mcp — Universal Agent Setup

## Install Binary

```bash
# From source
cd zara-jira-mcp && make build && make install

# Binary location after install
~/.local/bin/zara-jira-mcp
```

## Configuration Matrix

| Agent | Config File | Format |
|-------|------------|--------|
| **OpenCode** | `opencode.json` or `.opencode/opencode.json` | JSON `mcp.servers` |
| **Claude Code** | `.claude/settings.json` | JSON `mcpServers` |
| **Codex CLI** | `~/.codex/config.toml` | TOML `[mcp_servers]` |
| **Cursor** | `~/.cursor/mcp.json` | JSON `mcpServers` |
| **Copilot (VS Code)** | `.vscode/mcp.json` or `settings.json` | JSON `servers` / `github.copilot.chat.mcp.servers` |
| **Windsurf** | `~/.codeium/windsurf/mcp_config.json` | JSON `mcpServers` |
| **ChatGPT Desktop** | Settings > MCP | JSON `mcpServers` |
| **Kiro** | `kiro.json` | JSON `mcpServers` |

## Minimal Config (all agents use same env vars)

```
JIRA_BASE_URL=https://company.atlassian.net
JIRA_EMAIL=you@company.com
JIRA_API_TOKEN=your-jira-api-token
JIRA_AI_BASE_URL=https://api.openai.com
JIRA_AI_API_KEY=sk-your-key
JIRA_AI_MODEL=gpt-4o-mini
```

## Instructions Template (universal, paste into any agent)

```markdown
## PM Brain (zara-jira-mcp, 131 tools)

AI Scrum Master with persistent memory. Connected via MCP.

### First Use
1. `jira_boards` — get board_id (store it, needed everywhere)

### Daily
- `pm_standup_prep(board_id)` — talking points
- `pm_dashboard(board_id)` — full view

### Sprint Lifecycle
- Planning: `pm_planning_prep(board_id)`
- Mid-sprint: `pm_sprint_health(board_id)`, `pm_flow_metrics(board_id)`
- End: `pm_snapshot_sprint(board_id, velocity:N)`, `pm_scorecard(board_id)`
- Retro: `pm_facilitate(ceremony:"retro")`, `pm_record_retro(...)`

### When Asked "When Will It Be Done?"
- `pm_forecast(board_id, remaining_items:N)` — Monte Carlo probabilities

### Record Immediately
- Decisions: `pm_record_decision(title, decision, rationale)`
- Risks: `pm_record_risk(title, severity, owner, mitigation)`
- Blockers: `pm_record_blocker(description, issue_key, owner)`

### For Executives (NEVER use pm_dashboard for them)
- `pm_exec_report(board_id)` — business language, no jargon

### Rules
- Memory builds over time. More snapshots = better forecasts.
- `pm_auto_detect_risks(board_id)` weekly for proactive scanning.
- Tech debt: `pm_tech_debt_add`, `pm_tech_debt_budget(board_id)`.
```

## Wrapper Script (recommended for all agents)

`~/.local/bin/zara-jira-mcp-wrapper.sh`:

```bash
#!/usr/bin/env bash
# Load env from secure location
set -a
[ -f ~/.config/zara-jira-mcp/.env ] && source ~/.config/zara-jira-mcp/.env
set +a
export PM_PROFILE=${PM_PROFILE:-pm}  # lite, pm, standard, full, or remove for all 224 tools
exec zara-jira-mcp "$@"
```

```bash
chmod +x ~/.local/bin/zara-jira-mcp-wrapper
mkdir -p ~/.config/zara-jira-mcp
# Put your .env there (not in project repo)
```

Then use `zara-jira-mcp-wrapper` as command in all configs.

## Performance Profiles

224 tools can make ChatGPT Desktop or Claude Desktop slow. Use profiles:

| Profile | Tools | Recommended For |
|---------|-------|-----------------|
| `lite` | ~30 | Slow machines, mobile, minimal |
| `pm` | ~60 | **PM/Scrum Master daily work** |
| `standard` | ~100 | PM + all notification channels |
| `full` | ~150 | PM + GitHub/developer visibility |
| (none) | ~224 | Developers who want everything |

Set `PM_PROFILE=pm` in your env or wrapper script. See [docs/profiles.md](profiles.md) for details.

## Verification

Test MCP is working:

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}}}' | zara-jira-mcp
```

Expected: JSON response with `serverInfo.name: "zara-jira-mcp"`
