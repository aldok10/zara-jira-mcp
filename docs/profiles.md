# Performance Profiles

224 tools is powerful but heavy. Most AI clients load all tool definitions into their context window — that means slower responses and higher token usage. Profiles let you load only what you need.

## Quick Setup

Add one env var to your config:

```bash
PM_PROFILE=pm
```

Done. Your AI loads 60 tools instead of 224. Responses are faster, context is cleaner.

## Profiles

| Profile | Tools | Includes |
|---------|-------|----------|
| `lite` | ~30 | Jira CRUD + PM memory + shortcuts |
| `pm` | ~60 | Jira + PM memory + AI intelligence + stakeholder reports + portfolio |
| `standard` | ~100 | PM + notifications (Lark, Slack, Discord, Telegram, Teams, Email) |
| `full` | ~150 | Standard + GitHub/GitLab visibility |
| (default) | ~224 | Everything including Linear, PagerDuty, Clockify, Notion, Sheets |

## Which Profile Should I Use?

**Non-technical PM on ChatGPT Desktop:** `lite` or `pm`
- You get sprint management, forecasting, coaching, executive reports
- No developer tools cluttering your context

**Scrum Master daily work:** `pm`
- Everything you need: ceremonies, risks, blockers, decisions, velocity, health
- Stakeholder reporting and portfolio view included
- No notification channels (send reports manually or copy-paste)

**SM who sends reports via Slack/Lark:** `standard`
- PM profile + all notification channels
- Auto-send weekly digest, exec reports, escalations

**Technical PM / Engineering Manager:** `full`
- Standard + GitHub PR visibility, branch tracing, release tracking
- See developer activity without leaving your PM context

**Developer building/contributing to this project:** no profile (default)
- All 224 tools loaded

## Recommended by Client

| Client | Profile | Reason |
|--------|---------|--------|
| ChatGPT Desktop | `lite` or `pm` | Limited context, gets sluggish with 200+ tools |
| Claude Desktop | `pm` or `standard` | 200k context but cleaner with fewer tools |
| Cursor / Windsurf | `full` or default | IDE context is generous |
| Gemini CLI | `pm` or `standard` | Fast CLI, tools still eat context |
| Goose | `standard` | Good daily balance |
| OpenCode / Claude Code | `full` or default | Developer tool, wants everything |
| Kiro | `standard` or `full` | IDE with good context management |

## How to Set

### Option 1: In MCP config (per-client)

```json
{
  "mcpServers": {
    "jira-pm": {
      "command": "zara-jira-mcp",
      "env": {
        "PM_PROFILE": "pm",
        "JIRA_BASE_URL": "...",
        "JIRA_EMAIL": "...",
        "JIRA_API_TOKEN": "..."
      }
    }
  }
}
```

### Option 2: In wrapper script (all clients)

```bash
#!/usr/bin/env bash
set -a
[ -f ~/.config/zara-jira-mcp/.env ] && source ~/.config/zara-jira-mcp/.env
set +a
export PM_PROFILE=${PM_PROFILE:-pm}
exec zara-jira-mcp "$@"
```

### Option 3: In .env file

```bash
# ~/.config/zara-jira-mcp/.env
PM_PROFILE=pm
JIRA_BASE_URL=https://company.atlassian.net
JIRA_EMAIL=you@company.com
JIRA_API_TOKEN=your-token
```

## Custom Module Selection

For granular control, skip profiles and select modules directly:

```bash
PM_ENABLED_MODULES=jira,pm,ai,stakeholder
```

Available modules:
- `jira` — Jira CRUD, search, sprints, epics, bulk ops
- `pm` — Memory, AI intelligence, forecasting, coaching, process tools
- `ai` — AI analysis tools
- `notifications` — Lark, Slack, Discord, Telegram, Teams, Email
- `stakeholder` — Executive reports, management briefs, tech debt
- `portfolio` — Cross-project overview
- `github` — GitHub/GitLab PRs, activity, tracing
- `integrations` — Linear, PagerDuty, Clockify, Notion, Sheets
- `shortcuts` — Simplified PM commands (pm, pm_create, pm_decide, pm_risk, pm_next, pm_help)

## Performance Impact

| Profile | Context Tokens | ChatGPT Desktop | Claude Desktop |
|---------|---------------|-----------------|----------------|
| lite | ~3k | Fast | Fast |
| pm | ~6k | Fast | Fast |
| standard | ~10k | Normal | Fast |
| full | ~15k | Slightly slower | Fast |
| all | ~22k | Can lag noticeably | Normal |

Rule of thumb: if your AI takes > 3 seconds to start responding, reduce your profile.
