# zara-jira-mcp — OpenCode Setup

## MCP Configuration

Add to `opencode.json` or `.opencode/opencode.json`:

```json
{
  "mcp": {
    "jira-pm": {
      "type": "local",
      "command": ["zara-jira-mcp"],
      "env": {
        "JIRA_BASE_URL": "https://company.atlassian.net",
        "JIRA_EMAIL": "you@company.com",
        "JIRA_API_TOKEN": "{{env:JIRA_API_TOKEN}}",
        "JIRA_AI_BASE_URL": "https://api.openai.com",
        "JIRA_AI_API_KEY": "{{env:JIRA_AI_API_KEY}}",
        "JIRA_AI_MODEL": "gpt-4o-mini"
      },
      "timeout": 30000,
      "enabled": true
    }
  }
}
```

## Skill File

Copy `SKILL.md` to `.opencode/skills/zara-pm-brain/SKILL.md`:

```bash
mkdir -p .opencode/skills/zara-pm-brain
cp /path/to/zara-jira-mcp/SKILL.md .opencode/skills/zara-pm-brain/SKILL.md
```

Or for global access:

```bash
mkdir -p ~/.config/opencode/skills/zara-pm-brain
cp /path/to/zara-jira-mcp/SKILL.md ~/.config/opencode/skills/zara-pm-brain/SKILL.md
```

## Agent Instructions (`.opencode/instructions/pm.md`)

```markdown
## PM/Scrum Master Tools

When the user asks about sprints, Jira, team health, risks, blockers, decisions, 
retrospectives, or project management:

1. Load skill `zara-pm-brain` for full tool reference
2. Use `jira_boards` first to get board_id if not known
3. For quick status: `pm_dashboard(board_id)`
4. For standup prep: `pm_standup_prep(board_id)`
5. Record decisions, risks, blockers immediately when mentioned
6. Run `pm_snapshot_sprint` at end of every sprint

The MCP server has persistent memory — data accumulates over time and 
AI intelligence tools improve with more historical data.
```

## Wrapper Script (optional, for env loading)

`~/.local/bin/zara-jira-mcp-wrapper.sh`:

```bash
#!/bin/bash
set -a
source ~/.config/zara-jira-mcp/.env 2>/dev/null
set +a
exec zara-jira-mcp "$@"
```

Then in config: `"command": ["zara-jira-mcp-wrapper.sh"]`
