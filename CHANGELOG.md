# Changelog

## v0.4.0 (2026-06-26)

### Added
- **Performance Profiles** — `PM_PROFILE` env var: lite(30), pm(60), standard(100), full(150), all(224)
- **20 AI client pre-built configs** — Claude, ChatGPT, Cursor, Windsurf, Zed, Gemini CLI, Goose, Amazon Q, Cline, OpenCode, Kiro, Codex, Cherry Studio, Jan, Msty, LibreChat, TypingMind, Copilot Studio
- **Management reporting tools** — `pm_management_brief`, `pm_dependency_report`, `pm_escalation_report`, `pm_resource_utilization`, `pm_blocker_aging`, `pm_commitment_report`
- **SM leverage tools** — `pm_maturity_assessment`, `pm_sm_impact`, `pm_team_autonomy`, `pm_improvement_velocity`, `pm_outcome_map`
- **Stakeholder tools** — `pm_stakeholder_pulse`, `pm_stakeholder_trend`
- **Tech skill tools** — PM engineering literacy support
- **Database tools** — Postgres, MySQL, MongoDB read-only queries
- **GitHub Actions CI** — build + test + vet on push/PR
- **llms.txt** — AI crawler metadata for discoverability
- **CONTRIBUTING.md** — contributor guidelines
- **GitHub repo metadata** — description, 20 topics, homepage link

### Documentation
- `docs/reporting-guide.md` — scenario-based management reporting (10 scenarios)
- `docs/engineering-literacy.md` — engineering concepts for PM/SM + learning path
- `docs/communication-frameworks.md` — 10 communication frameworks with tool mapping
- `docs/okr-kpi-integration.md` — OKR/KPI research + Lark OKR sync design
- `docs/context.md` — research-backed project context (PwC, DORA, PMI, Gartner)
- `docs/implementation-roadmap.md` — 7-sprint prioritized backlog
- `docs/profiles.md` — performance profile guide per AI client
- `AGENTS.md` — universal AI agent integration guide
- Per-client setup in `docs/agents/` (21 files)

### Changed
- README rewritten — promotional, conversational, author info, LinkedIn/SociaBuzz links
- ROADMAP updated to v0.4.0 — 224 tools, profile system, security findings tracked
- `.env.example` expanded — all integrations + PM_PROFILE

### Infrastructure
- `.claude/` — Claude Code settings + CLAUDE.md
- `.opencode/` — OpenCode config + instructions
- `.cursor/` — Cursor rules
- `.github/copilot-instructions.md` — GitHub Copilot
- `.vscode/mcp.json` — VS Code Copilot MCP config
- `.zed/settings.json` — Zed context_servers
- `.kiro/` — Kiro config + instructions
- `.codex/` — Codex CLI instructions
- `.windsurfrules` — Windsurf project rules

---

## v0.3.0 (2026-06-20)

### Added
- SM leverage tools: early warning, individual signals, WIP guardian, question coach
- Portfolio tools: overview, risks, blockers, workload, AI summary
- Outcome tools: impediment aging, SM impact, stakeholder pulse, improvement velocity, team autonomy, outcome map
- Care tools: commitment check, team energy
- Flow tools: flow metrics, sprint comparison, ceremony facilitator, confidence, goal check
- Forecast tools: Monte Carlo (10k simulations), sprint forecast, capacity plan
- Process tools: DoD, DoR, agreements, experiments, sprint goals
- Recipe tools: start_work, done, block
- Coaching tools: AI coaching, anti-patterns, retro analysis, facilitate, check_ready
- Deep PM tools: planning prep, review prep, backlog groom, scope creep
- Help tools: pm_help, pm_quickstart, pm_workflow (7 built-in workflows)
- PM shortcuts: pm, pm_create, pm_decide, pm_risk, pm_next
- Multi-channel notifications: Lark, Slack, Discord, Telegram, Teams, Email, Confluence
- Smart routing + broadcast
- GitHub/GitLab: full CRUD, milestones, MRs, file reading
- Git integration: branch tracing, smart commit, link PR
- Linear, PagerDuty, Clockify, Notion, Google Sheets, Calendar integrations
- Module system: PM_ENABLED_MODULES for custom tool selection

---

## v0.2.0 (2026-06-10)

### Added
- PM memory layer (14 SQLite tables)
- Sprint snapshots, risks, decisions, blockers, team metrics, retros
- Dependencies, meeting notes, health scores, sprint goals, DoD, escalations
- AI intelligence: recommendations, standup prep, exec report, weekly digest
- Jira CRUD: search, create, update, transition, bulk ops, epics, sprints
- Lark webhook notifications
- Slack notifications

---

## v0.1.0 (2026-06-01)

### Added
- Initial MCP server (stdio transport)
- Basic Jira search and issue operations
- SQLite memory foundation
- OpenAI-compatible AI provider
- mcp-go integration
