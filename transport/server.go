package transport

import (
	"os"
	"strings"

	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	s *server.MCPServer
}

type regFunc func(*server.MCPServer, *tools.Handlers)

func NewMCPServer(handlers *tools.Handlers) *MCPServer {
	s := server.NewMCPServer(
		"zara-jira-mcp",
		"0.5.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	enabled := enabledModules()

	modules := map[string][]regFunc{
		// Jira (split by depth)
		"jira":        {registerJiraTools},          // 10: search, get, boards, sprint, create, comment, transitions, update, health
		"jira-ops":    {registerIssueOpsTools},      // 5: assign, unassign, delete, subtask, find_user
		"jira-deep":   {registerEpicSprintTools, registerBulkProjectTools, registerLinkWorklogTools, registerVersionTools, registerTraceTools}, // 28: epics, bulk, worklog, versions, trace

		// PM (split by function)
		"pm-memory":   {registerPMTools, registerMemoryTools},       // 16: my_issues, overdue, workload, snapshot, risks, decisions, blockers, retros
		"pm-analysis": {registerPMIntelTools, registerForecastTools, registerFlowTools}, // 15: recommendations, velocity, standup, retro, forecast, anti-patterns, scope, flow
		"pm-planning": {registerAdvancedPMTools, registerDeepPMTools, registerProcessTools, registerRecipeTools, registerFacilitationTools}, // 30: health, goals, DoD, capacity, deps, burndown, process, recipes, retro format, meeting audit
		"pm-intel":    {registerOKRKPITools, registerCoachingTools, registerInsightTools, registerV6Tools, registerSafetyTools, registerStoryPointsTools, registerImprovementTools}, // 38: OKR, KPI, coaching, experiments, safety, metrics

		// AI
		"ai": {registerAITools}, // 2: analyze, sprint_report

		// Notifications (split by platform)
		"notify-lark":  {registerLarkTools},   // 1: lark notify
		"notify-slack": {registerSlackTools},  // 4: slack channel/msg/reactions
		"notify-all":   {registerPlatformTools, registerRoutingTools}, // 9: multi-platform routing

		// Stakeholder (split by depth)
		"stakeholder":      {registerStakeholderTools, registerTechDebtTools, registerLeverageTools, registerManagementTools, registerReportingTools, registerWhatNextTools, registerCommunicationTools}, // 25: exec, scorecard, kb, tech debt, meetings, leverage
		"stakeholder-deep": {registerCommsGapTools, registerSafetyTools, registerCareTools, registerOutcomeTools, registerTechSkillTools}, // 44: sentiment, feedback, safety, NVC, hard convos

		// Portfolio, GitHub, Integrations
		"portfolio":    {registerPortfolioTools},                   // 5
		"github":       {registerGitHubTools},                      // 3: link PR, branch, smart commit
		"github-deep":  {registerGitHubFullTools, registerGitIntegrationTools}, // 17: issues, repos, CI, full GitHub ops
		"integrations": {registerCalendarTools, registerNotionTools, registerLinearTools, registerPagerDutyTools, registerClockifyTools, registerSheetsTools, registerLarkOKRTools}, // 19

		// Smart routing (primary entry points)
		"smart-router": {registerSmartContextTools}, // 6: pm_smart, pm_do, pm_report + predictive
		"pm-quick":     {registerPMShortcuts},       // 5: pm, pm_create, pm_decide, pm_risk, pm_next
		"help":         {registerHelpTools},         // 3: pm_help, pm_quickstart, pm_workflow
	}

	for mod, fns := range modules {
		if enabled[mod] {
			for _, fn := range fns {
				fn(s, handlers)
			}
		}
	}

	return &MCPServer{s: s}
}

func enabledModules() map[string]bool {
	// Profile presets — research-backed tool counts (BCG 2026, Microsoft Research, Google Research)
	// See docs/research-blueprint.md for methodology.
	profile := os.Getenv("PM_PROFILE")
	switch profile {
	case "chatgpt":
		// ~11 tools: smart router + quick actions. Safe for ChatGPT Desktop (under 15).
		return map[string]bool{"smart-router": true, "pm-quick": true}
	case "lite":
		// ~24 tools: smart router + quick + help + Jira core.
		// Under Google's 16-tool threshold once factoring agent overhead.
		return map[string]bool{"smart-router": true, "pm-quick": true, "help": true, "jira": true}
	case "standard":
		// ~40 tools: lite + jira-ops + pm-memory + AI.
		// Under Microsoft Research's <40 guidance for PM teams.
		return map[string]bool{"smart-router": true, "pm-quick": true, "jira": true, "jira-ops": true, "pm-memory": true, "ai": true}
	case "full":
		// ~48 tools: standard + help + pm-analysis + basic stakeholder.
		// Acceptable with good tool descriptions (research-backed: 40-60).
		return map[string]bool{"smart-router": true, "pm-quick": true, "help": true, "jira": true, "jira-ops": true, "pm-memory": true, "pm-analysis": true, "ai": true, "notify-lark": true, "stakeholder": true}
	}

	// Custom module selection via PM_ENABLED_MODULES
	env := os.Getenv("PM_ENABLED_MODULES")
	if env == "" || env == "all" {
		return map[string]bool{
			"jira": true, "jira-ops": true, "jira-deep": true,
			"pm-memory": true, "pm-analysis": true, "pm-planning": true, "pm-intel": true,
			"ai": true,
			"notify-lark": true, "notify-slack": true, "notify-all": true,
			"stakeholder": true, "stakeholder-deep": true,
			"portfolio": true,
			"github": true, "github-deep": true,
			"integrations": true,
			"smart-router": true, "pm-quick": true, "help": true,
		}
	}
	m := map[string]bool{}
	for _, s := range strings.Split(env, ",") {
		m[strings.TrimSpace(s)] = true
	}
	return m
}

func (m *MCPServer) Server() *server.MCPServer {
	return m.s
}

func registerJiraTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_search",
			mcp.WithDescription("Search Jira issues via JQL. Returns key, summary, status, assignee."),
			mcp.WithString("jql", mcp.Required(), mcp.Description("JQL query string")),
			mcp.WithNumber("max_results", mcp.Description("Maximum results (default 20, max 50)")),
			mcp.WithNumber("start_at", mcp.Description("Pagination offset (default 0)")),
		),
		h.SearchIssues,
	)

	s.AddTool(
		mcp.NewTool("jira_get_issue",
			mcp.WithDescription("Get full details of a single Jira issue by key."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
		),
		h.GetIssue,
	)

	s.AddTool(
		mcp.NewTool("jira_boards",
			mcp.WithDescription("List all accessible Jira boards with their IDs and types."),
		),
		h.GetBoards,
	)

	s.AddTool(
		mcp.NewTool("jira_sprint_summary",
			mcp.WithDescription("Get active sprint status breakdown and issue list for a board."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID (get from jira_boards)")),
		),
		h.GetSprintSummary,
	)

	s.AddTool(
		mcp.NewTool("jira_create_issue",
			mcp.WithDescription("Create a new Jira issue. Returns the created issue key."),
			mcp.WithString("project", mcp.Required(), mcp.Description("Project key (e.g. SIT, VTPRO)")),
			mcp.WithString("summary", mcp.Required(), mcp.Description("Issue title/summary")),
			mcp.WithString("issue_type", mcp.Description("Issue type: Task, Bug, Story (default: Task)")),
			mcp.WithString("description", mcp.Description("Detailed description")),
			mcp.WithString("priority", mcp.Description("Priority: Highest, High, Medium, Low, Lowest")),
			mcp.WithString("assignee_id", mcp.Description("Assignee account ID")),
			mcp.WithString("labels", mcp.Description("Comma-separated labels")),
		),
		h.CreateIssue,
	)

	s.AddTool(
		mcp.NewTool("jira_add_comment",
			mcp.WithDescription("Add a comment to a Jira issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. SIT-123)")),
			mcp.WithString("body", mcp.Required(), mcp.Description("Comment text")),
		),
		h.AddComment,
	)

	s.AddTool(
		mcp.NewTool("jira_transitions",
			mcp.WithDescription("Available status transitions for an issue. Use IDs with jira_transition."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key")),
		),
		h.GetTransitions,
	)

	s.AddTool(
		mcp.NewTool("jira_transition",
			mcp.WithDescription("Transition an issue to a new status. Get transition IDs from jira_transitions."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key")),
			mcp.WithString("transition_id", mcp.Required(), mcp.Description("Transition ID (from jira_transitions)")),
		),
		h.TransitionIssue,
	)

	s.AddTool(
		mcp.NewTool("jira_update_issue",
			mcp.WithDescription("Update an existing Jira issue. Only provided fields are changed."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
			mcp.WithString("summary", mcp.Description("New summary/title")),
			mcp.WithString("description", mcp.Description("New description")),
			mcp.WithString("priority", mcp.Description("New priority: Highest, High, Medium, Low, Lowest")),
			mcp.WithString("assignee_id", mcp.Description("New assignee account ID")),
			mcp.WithString("labels", mcp.Description("Comma-separated labels (replaces existing)")),
		),
		h.UpdateIssue,
	)

	s.AddTool(
		mcp.NewTool("jira_health",
			mcp.WithDescription("Health check and version info for zara-jira-mcp server."),
		),
		h.Health,
	)
}

func registerIssueOpsTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_assign",
			mcp.WithDescription("Assign a Jira issue to a user by account ID."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
			mcp.WithString("assignee_id", mcp.Required(), mcp.Description("Assignee account ID (use jira_find_user to look up)")),
		),
		h.AssignIssue,
	)

	s.AddTool(
		mcp.NewTool("jira_unassign",
			mcp.WithDescription("Remove the assignee from a Jira issue."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
		),
		h.UnassignIssue,
	)

	s.AddTool(
		mcp.NewTool("jira_delete_issue",
			mcp.WithDescription("Delete a Jira issue permanently."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Issue key (e.g. PROJ-123)")),
		),
		h.DeleteIssue,
	)

	s.AddTool(
		mcp.NewTool("jira_create_subtask",
			mcp.WithDescription("Create a subtask under a parent issue."),
			mcp.WithString("parent_key", mcp.Required(), mcp.Description("Parent issue key (e.g. PROJ-123)")),
			mcp.WithString("summary", mcp.Required(), mcp.Description("Subtask summary")),
			mcp.WithString("description", mcp.Description("Subtask description")),
			mcp.WithString("assignee_id", mcp.Description("Assignee account ID")),
			mcp.WithString("priority", mcp.Description("Priority: Highest, High, Medium, Low, Lowest")),
		),
		h.CreateSubtask,
	)

	s.AddTool(
		mcp.NewTool("jira_find_user",
			mcp.WithDescription("Search for Jira users by name or email. Use to find account IDs for assignment."),
			mcp.WithString("query", mcp.Required(), mcp.Description("Search by name or email")),
		),
		h.FindUser,
	)
}

func registerPMTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_my_issues",
			mcp.WithDescription("Show issues assigned to the current user (you). Optionally filter by status."),
			mcp.WithString("status", mcp.Description("Filter by status name (e.g. 'In Progress', 'To Do')")),
		),
		h.MyIssues,
	)

	s.AddTool(
		mcp.NewTool("jira_overdue",
			mcp.WithDescription("Find stale/overdue issues with no updates in N days. Great for PM follow-ups."),
			mcp.WithNumber("days", mcp.Description("Days without update to consider stale (default: 14)")),
			mcp.WithString("project", mcp.Description("Filter by project key")),
		),
		h.Overdue,
	)

	s.AddTool(
		mcp.NewTool("jira_workload",
			mcp.WithDescription("Show workload distribution — how many open issues each team member has."),
			mcp.WithString("project", mcp.Description("Filter by project key")),
		),
		h.Workload,
	)
}

func registerAITools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_ai_analyze",
			mcp.WithDescription("AI analysis of Jira tickets. Ask about blockers, health, staleness."),
			mcp.WithString("query", mcp.Required(), mcp.Description("Your question about the project/tickets")),
			mcp.WithString("jql", mcp.Description("JQL to scope the analysis")),
			mcp.WithNumber("max_results", mcp.Description("Max tickets to analyze (default 30)")),
		),
		h.AIAnalyze,
	)

	s.AddTool(
		mcp.NewTool("jira_ai_sprint_report",
			mcp.WithDescription("AI sprint report with health assessment and recommendations. Sends to Lark."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithBoolean("send_to_lark", mcp.Description("Send report to Lark group")),
		),
		h.AISprintReport,
	)
}

func registerLarkTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_notify_lark",
			mcp.WithDescription("Send a markdown message to the configured Lark group."),
			mcp.WithString("title", mcp.Description("Card title (default: 'Jira Update')")),
			mcp.WithString("content", mcp.Required(), mcp.Description("Message content in markdown")),
		),
		h.NotifyLark,
	)
}

func registerMemoryTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_snapshot_sprint",
			mcp.WithDescription("Snapshot sprint state to memory. Auto-calculates from Jira. Call each sprint end."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithNumber("velocity", mcp.Description("Story points completed this sprint (manual input)")),
			mcp.WithNumber("carryover", mcp.Description("Issues carried over from previous sprint")),
			mcp.WithString("notes", mcp.Description("Any notes about this sprint snapshot")),
		),
		h.SnapshotSprint,
	)

	s.AddTool(
		mcp.NewTool("pm_record_risk",
			mcp.WithDescription("Record a project risk: severity, owner, mitigation. Stored in register."),
			mcp.WithString("title", mcp.Required(), mcp.Description("Risk title")),
			mcp.WithString("description", mcp.Description("Detailed description")),
			mcp.WithString("severity", mcp.Description("critical, high, medium, low (default: medium)")),
			mcp.WithString("owner", mcp.Description("Who owns mitigating this risk")),
			mcp.WithString("mitigation", mcp.Description("Mitigation plan")),
			mcp.WithString("sprint_name", mcp.Description("Sprint where risk was identified")),
		),
		h.RecordRisk,
	)

	s.AddTool(
		mcp.NewTool("pm_update_risk",
			mcp.WithDescription("Update a risk's status (open, mitigating, resolved, accepted)."),
			mcp.WithNumber("risk_id", mcp.Required(), mcp.Description("Risk ID")),
			mcp.WithString("status", mcp.Required(), mcp.Description("New status: open, mitigating, resolved, accepted")),
			mcp.WithString("mitigation", mcp.Description("Updated mitigation plan")),
			mcp.WithString("owner", mcp.Description("Updated owner")),
			mcp.WithString("severity", mcp.Description("Updated severity")),
		),
		h.UpdateRisk,
	)

	s.AddTool(
		mcp.NewTool("pm_risk_dashboard",
			mcp.WithDescription("Show all open risks sorted by severity. The PM's risk radar."),
		),
		h.GetRiskDashboard,
	)

	s.AddTool(
		mcp.NewTool("pm_record_decision",
			mcp.WithDescription("Record a project decision with context and rationale. Builds memory."),
			mcp.WithString("title", mcp.Required(), mcp.Description("Decision title (e.g. 'Use PostgreSQL over MongoDB')")),
			mcp.WithString("decision", mcp.Required(), mcp.Description("What was decided")),
			mcp.WithString("context", mcp.Description("What situation led to this decision")),
			mcp.WithString("rationale", mcp.Description("Why this over alternatives")),
			mcp.WithString("made_by", mcp.Description("Who made the decision")),
			mcp.WithString("tags", mcp.Description("Comma-separated tags (e.g. 'architecture,database')")),
		),
		h.RecordDecision,
	)

	s.AddTool(
		mcp.NewTool("pm_search_decisions",
			mcp.WithDescription("Search decision log by keyword or list recent decisions."),
			mcp.WithString("query", mcp.Description("Search keyword (searches title, decision, tags)")),
			mcp.WithNumber("limit", mcp.Description("Max results (default 10)")),
		),
		h.SearchDecisions,
	)

	s.AddTool(
		mcp.NewTool("pm_record_blocker",
			mcp.WithDescription("Record an impediment/blocker. Track how long things stay blocked."),
			mcp.WithString("description", mcp.Required(), mcp.Description("What is blocked and why")),
			mcp.WithString("issue_key", mcp.Description("Related Jira issue key")),
			mcp.WithString("owner", mcp.Description("Who is responsible for resolving")),
		),
		h.RecordBlocker,
	)

	s.AddTool(
		mcp.NewTool("pm_resolve_blocker",
			mcp.WithDescription("Mark a blocker as resolved with resolution details."),
			mcp.WithNumber("blocker_id", mcp.Required(), mcp.Description("Blocker ID")),
			mcp.WithString("resolution", mcp.Required(), mcp.Description("How was it resolved")),
		),
		h.ResolveBlocker,
	)

	s.AddTool(
		mcp.NewTool("pm_blockers",
			mcp.WithDescription("Show active blockers or blocker history."),
			mcp.WithBoolean("show_history", mcp.Description("Show resolved blockers too (default: false)")),
		),
		h.GetBlockers,
	)

	s.AddTool(
		mcp.NewTool("pm_record_team_metric",
			mcp.WithDescription("Record a team member's sprint performance. Builds workload patterns over time."),
			mcp.WithString("member_name", mcp.Required(), mcp.Description("Team member name")),
			mcp.WithString("sprint_name", mcp.Required(), mcp.Description("Sprint name")),
			mcp.WithNumber("issues_assigned", mcp.Description("Issues assigned")),
			mcp.WithNumber("issues_done", mcp.Description("Issues completed")),
			mcp.WithNumber("blocker_count", mcp.Description("Times blocked")),
			mcp.WithNumber("carryover_count", mcp.Description("Issues carried to next sprint")),
			mcp.WithString("notes", mcp.Description("Notes about this member's sprint")),
		),
		h.RecordTeamMetric,
	)

	s.AddTool(
		mcp.NewTool("pm_team_health",
			mcp.WithDescription("Show team workload overview for a sprint, or individual member history."),
			mcp.WithString("sprint_name", mcp.Description("Sprint name for team overview")),
			mcp.WithString("member_name", mcp.Description("Member name for individual history")),
		),
		h.GetTeamHealth,
	)

	s.AddTool(
		mcp.NewTool("pm_record_retro",
			mcp.WithDescription("Record a sprint retrospective: what went well, what to improve, action items."),
			mcp.WithString("sprint_name", mcp.Required(), mcp.Description("Sprint name")),
			mcp.WithString("went_well", mcp.Description("What went well")),
			mcp.WithString("improvements", mcp.Description("What needs improvement")),
			mcp.WithString("action_items", mcp.Description("Action items (newline-separated)")),
		),
		h.RecordRetrospective,
	)

	s.AddTool(
		mcp.NewTool("pm_action_items",
			mcp.WithDescription("Show pending action items from retrospectives. Never let retro actions die."),
		),
		h.GetActionItems,
	)
}

func registerPMIntelTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_recommendations",
			mcp.WithDescription("AI recommendations from all PM memory: trends, risks, blockers, metrics."),
			mcp.WithNumber("board_id", mcp.Description("Board ID for sprint context")),
			mcp.WithString("focus", mcp.Description("Focus area: general, velocity, risks, team, process (default: general)")),
		),
		h.PMRecommendations,
	)

	s.AddTool(
		mcp.NewTool("pm_velocity_trend",
			mcp.WithDescription("Velocity and completion trends. Detects improvement or decline."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.VelocityTrend,
	)

	s.AddTool(
		mcp.NewTool("pm_standup_prep",
			mcp.WithDescription("Standup prep: Jira data + blockers, risks, actions as talking points."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.StandupPrep,
	)

	s.AddTool(
		mcp.NewTool("pm_retro_analysis",
			mcp.WithDescription("AI retro analysis: recurring issues, trends, root causes, suggestions."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.SprintRetroAnalysis,
	)
}

func registerAdvancedPMTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_sprint_health",
			mcp.WithDescription("Sprint health (0-100): velocity, blockers, scope creep, balance. Saves history."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.SprintHealthScore,
	)

	s.AddTool(
		mcp.NewTool("pm_health_history",
			mcp.WithDescription("Health score trends over time. Is team improving or declining?"),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.HealthHistory,
	)

	s.AddTool(
		mcp.NewTool("pm_record_dependency",
			mcp.WithDescription("Record dependency between issues (blocks, blocked_by, external)."),
			mcp.WithString("from_issue", mcp.Required(), mcp.Description("Issue that is blocked/dependent")),
			mcp.WithString("to_issue", mcp.Required(), mcp.Description("Issue/team it depends on")),
			mcp.WithString("type", mcp.Description("blocks, blocked_by, relates_to, external (default: blocks)")),
			mcp.WithString("description", mcp.Description("Context about the dependency")),
		),
		h.RecordDependency,
	)

	s.AddTool(
		mcp.NewTool("pm_resolve_dependency",
			mcp.WithDescription("Mark a dependency as resolved."),
			mcp.WithNumber("dependency_id", mcp.Required(), mcp.Description("Dependency ID")),
		),
		h.ResolveDependency,
	)

	s.AddTool(
		mcp.NewTool("pm_dependencies",
			mcp.WithDescription("Show dependency map — all open dependencies or for a specific issue."),
			mcp.WithString("issue_key", mcp.Description("Filter by issue key (shows all if empty)")),
		),
		h.GetDependencies,
	)

	s.AddTool(
		mcp.NewTool("pm_record_meeting",
			mcp.WithDescription("Record meeting notes: decisions made, action items, key discussion points."),
			mcp.WithString("meeting_type", mcp.Required(), mcp.Description("standup, planning, retro, grooming, adhoc")),
			mcp.WithString("notes", mcp.Description("Key discussion points")),
			mcp.WithString("decisions", mcp.Description("Decisions made during meeting")),
			mcp.WithString("action_items", mcp.Description("Follow-up actions")),
			mcp.WithString("attendees", mcp.Description("Comma-separated attendee names")),
			mcp.WithString("sprint_name", mcp.Description("Sprint context")),
		),
		h.RecordMeeting,
	)

	s.AddTool(
		mcp.NewTool("pm_meetings",
			mcp.WithDescription("Show meeting notes history. Filter by type (standup, planning, retro, grooming)."),
			mcp.WithString("meeting_type", mcp.Description("Filter by type")),
			mcp.WithNumber("limit", mcp.Description("Max results (default 10)")),
		),
		h.GetMeetings,
	)

	s.AddTool(
		mcp.NewTool("pm_capacity_plan",
			mcp.WithDescription("Capacity plan from velocity. Recommends next sprint points by availability."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithNumber("team_size", mcp.Description("Number of team members")),
			mcp.WithNumber("sprint_days", mcp.Description("Sprint duration in days (default: 10)")),
			mcp.WithNumber("planned_leave_days", mcp.Description("Total planned leave days across team")),
		),
		h.CapacityPlan,
	)

	s.AddTool(
		mcp.NewTool("pm_auto_detect_risks",
			mcp.WithDescription("Scan for risk signals: stale tickets, overload, chronic blockers. Auto-records."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.AutoDetectRisks,
	)
}

func registerDeepPMTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_track_daily",
			mcp.WithDescription("Track today's sprint progress (burndown data). Call daily."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.TrackDailyProgress,
	)

	s.AddTool(
		mcp.NewTool("pm_burndown",
			mcp.WithDescription("Show sprint burndown data with daily progress tracking."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithString("sprint_name", mcp.Description("Sprint name (default: active sprint)")),
		),
		h.GetBurndown,
	)

	s.AddTool(
		mcp.NewTool("pm_set_sprint_goal",
			mcp.WithDescription("Define sprint goal with measurable key results. Tracks achievement."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithString("goal", mcp.Required(), mcp.Description("Sprint goal statement")),
			mcp.WithString("key_results", mcp.Description("Measurable key results (newline-separated)")),
			mcp.WithString("sprint_name", mcp.Description("Sprint name (default: active sprint)")),
		),
		h.SetSprintGoal,
	)

	s.AddTool(
		mcp.NewTool("pm_close_sprint_goal",
			mcp.WithDescription("Close a sprint goal with outcome assessment."),
			mcp.WithNumber("goal_id", mcp.Required(), mcp.Description("Goal ID")),
			mcp.WithString("status", mcp.Required(), mcp.Description("achieved, partially_achieved, missed")),
			mcp.WithString("outcome", mcp.Description("What actually happened")),
		),
		h.CloseSprintGoal,
	)

	s.AddTool(
		mcp.NewTool("pm_sprint_goals",
			mcp.WithDescription("Show active sprint goals or goal achievement history."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithBoolean("show_history", mcp.Description("Show past goals with outcomes (default: false)")),
		),
		h.GetSprintGoals,
	)

	s.AddTool(
		mcp.NewTool("pm_dod",
			mcp.WithDescription("Manage Definition of Done checklist. Actions: list (default), add, remove."),
			mcp.WithString("action", mcp.Description("list, add, remove (default: list)")),
			mcp.WithString("project", mcp.Description("Project key (default: * for global)")),
			mcp.WithString("item", mcp.Description("DoD item text (required for action=add)")),
			mcp.WithString("category", mcp.Description("code, testing, docs, review, deploy (default: general)")),
			mcp.WithNumber("item_id", mcp.Description("Item ID (required for action=remove)")),
		),
		h.ManageDoD,
	)

	s.AddTool(
		mcp.NewTool("pm_escalate",
			mcp.WithDescription("Auto-escalate to Lark: risks/blockers >3d or health <50."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.Escalate,
	)

	s.AddTool(
		mcp.NewTool("pm_escalations",
			mcp.WithDescription("Show escalation history — what was escalated, when, acknowledged status."),
		),
		h.GetEscalations,
	)

	s.AddTool(
		mcp.NewTool("pm_dashboard",
			mcp.WithDescription("Full PM dashboard: progress, health, risks, blockers, deps, goals, actions."),
			mcp.WithNumber("board_id", mcp.Description("Board ID for sprint-specific data")),
		),
		h.PMDashboard,
	)

	s.AddTool(
		mcp.NewTool("pm_release_notes",
			mcp.WithDescription("Release notes from completed sprint issues. Can send to Lark."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithBoolean("send_to_lark", mcp.Description("Send release notes to Lark group")),
		),
		h.GenerateReleaseNotes,
	)
}

func registerFlowTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_flow_metrics",
			mcp.WithDescription("Flow metrics: WIP, throughput, cycle/lead time. Detects flow problems."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.FlowMetrics,
	)

	s.AddTool(
		mcp.NewTool("pm_sprint_compare",
			mcp.WithDescription("Compare current vs previous sprint: velocity, completion, blockers, carryover."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.SprintComparison,
	)

	s.AddTool(
		mcp.NewTool("pm_facilitate",
			mcp.WithDescription("AI ceremony facilitator with contextual prompts for any Scrum ceremony."),
			mcp.WithString("ceremony", mcp.Required(), mcp.Description("standup, planning, retro, grooming, review")),
			mcp.WithNumber("board_id", mcp.Description("Board ID for context (optional)")),
		),
		h.CeremonyFacilitator,
	)

	s.AddTool(
		mcp.NewTool("pm_confidence",
			mcp.WithDescription("Record sprint confidence (1-5). Tracks pre-sprint vs actual outcome."),
			mcp.WithString("sprint_name", mcp.Required(), mcp.Description("Sprint name")),
			mcp.WithNumber("score", mcp.Required(), mcp.Description("1=very worried, 2=worried, 3=neutral, 4=confident, 5=very confident")),
			mcp.WithString("member", mcp.Description("Team member name (default: 'team')")),
			mcp.WithString("note", mcp.Description("Why this confidence level?")),
		),
		h.RecordConfidence,
	)

	s.AddTool(
		mcp.NewTool("pm_goal_check",
			mcp.WithDescription("AI sprint goal progress check. Evaluates if on track from current data."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.SprintGoalCheck,
	)
}

func registerForecastTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_forecast",
			mcp.WithDescription("Monte Carlo forecast: 10K simulations from throughput. 50/70/85/95% dates."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithNumber("remaining_items", mcp.Description("Items remaining (default: from active sprint)")),
			mcp.WithNumber("sprint_days", mcp.Description("Sprint length in days (default: 10)")),
		),
		h.MonteCarloForecast,
	)

	s.AddTool(
		mcp.NewTool("pm_anti_patterns",
			mcp.WithDescription("Detect Scrum anti-patterns: zombie sprints, hero culture, scope creep."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.DetectAntiPatterns,
	)

	s.AddTool(
		mcp.NewTool("pm_coaching",
			mcp.WithDescription("AI coaching for Scrum Masters. Data-driven improvement suggestions."),
			mcp.WithString("topic", mcp.Required(), mcp.Description("team_dynamics, velocity, blockers, morale, conflict, growth, predictability")),
			mcp.WithNumber("board_id", mcp.Description("Board ID for data context")),
			mcp.WithString("situation", mcp.Description("Describe the specific situation you need advice on")),
		),
		h.CoachingAdvice,
	)


	s.AddTool(
		mcp.NewTool("pm_nl_to_jql",
			mcp.WithDescription("Convert natural language to JQL query."),
			mcp.WithString("query", mcp.Required(), mcp.Description("Natural language query (e.g. 'my open bugs with high priority')")),
		),
		h.NLToJQL,
	)

	s.AddTool(
		mcp.NewTool("pm_scope_creep",
			mcp.WithDescription("Detect mid-sprint scope changes. Compares current items vs baseline snapshot."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.ScopeCreep,
	)

	s.AddTool(
		mcp.NewTool("pm_backlog_groom",
			mcp.WithDescription("Find stale backlog items needing grooming. Untouched N days, not in sprint."),
			mcp.WithString("project", mcp.Description("Project key filter")),
			mcp.WithNumber("days", mcp.Description("Days without update to consider stale (default: 90)")),
		),
		h.BacklogGroom,
	)
}

func registerRecipeTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_recipe_start_work",
		mcp.WithDescription("Start work: assign, transition to In Progress, suggest branch. One-click."),
		mcp.WithString("key", mcp.Required(), mcp.Description("Issue key")),
		mcp.WithString("assignee_id", mcp.Description("Your account ID (use jira_find_user to look up)")),
	), h.RecipeStartWork)

	s.AddTool(mcp.NewTool("pm_recipe_done",
		mcp.WithDescription("Mark done: transition, optionally log time and add comment."),
		mcp.WithString("key", mcp.Required(), mcp.Description("Issue key")),
		mcp.WithString("time_spent", mcp.Description("Time to log (e.g. '2h', '30m')")),
		mcp.WithString("comment", mcp.Description("Completion comment")),
	), h.RecipeDone)

	s.AddTool(mcp.NewTool("pm_recipe_block",
		mcp.WithDescription("Flag blocked: record in memory, comment on issue, create impediment trail."),
		mcp.WithString("key", mcp.Required(), mcp.Description("Issue key")),
		mcp.WithString("reason", mcp.Required(), mcp.Description("Why is it blocked?")),
		mcp.WithString("owner", mcp.Description("Who should resolve this?")),
	), h.RecipeBlock)
}

func registerStakeholderTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_exec_report",
			mcp.WithDescription("Executive report: outcomes, risks, health. No jargon. VP-ready."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithBoolean("send_to_lark", mcp.Description("Send to Lark group")),
		),
		h.ExecutiveReport,
	)

	s.AddTool(
		mcp.NewTool("pm_scorecard",
			mcp.WithDescription("Sprint scorecard (0-100): completion, goals, predictability, quality."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.SprintScorecard,
	)

	s.AddTool(
		mcp.NewTool("pm_team_kb",
			mcp.WithDescription("Team knowledge base: DoD, decisions, patterns, metrics. Ask or browse."),
			mcp.WithNumber("board_id", mcp.Description("Board ID for metrics context")),
			mcp.WithString("question", mcp.Description("Ask a question about how the team works (AI-powered answer)")),
		),
		h.TeamKnowledgeBase,
	)

	s.AddTool(
		mcp.NewTool("pm_record_learning",
			mcp.WithDescription("Record team learning or tribal knowledge. Builds searchable memory."),
			mcp.WithString("title", mcp.Required(), mcp.Description("Learning title")),
			mcp.WithString("learning", mcp.Required(), mcp.Description("What was learned")),
			mcp.WithString("context", mcp.Description("What situation triggered this learning")),
			mcp.WithString("tags", mcp.Description("Comma-separated tags")),
			mcp.WithString("author", mcp.Description("Who learned this")),
		),
		h.RecordLearning,
	)

	s.AddTool(
		mcp.NewTool("pm_weekly_digest",
			mcp.WithDescription("Weekly digest: decisions, risks, blockers, wins. AI-summarized."),
			mcp.WithNumber("board_id", mcp.Description("Board ID")),
			mcp.WithBoolean("send_to_lark", mcp.Description("Send digest to Lark")),
		),
		h.WeeklyDigest,
	)
}

func registerProcessTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_dor",
			mcp.WithDescription("Manage Definition of Ready (entry gate for sprint). Actions: list, add, remove."),
			mcp.WithString("action", mcp.Description("list, add, remove (default: list)")),
			mcp.WithString("item", mcp.Description("DoR item (for add)")),
			mcp.WithString("project", mcp.Description("Project key (default: *)")),
			mcp.WithString("category", mcp.Description("clarity, estimation, dependencies, design, size, support")),
			mcp.WithNumber("item_id", mcp.Description("Item ID (for remove)")),
		),
		h.ManageDoR,
	)

	s.AddTool(
		mcp.NewTool("pm_check_ready",
			mcp.WithDescription("AI evaluates if story meets DoR + INVEST. Returns READY/NOT READY."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Jira issue key to evaluate")),
		),
		h.CheckStoryReady,
	)

	s.AddTool(
		mcp.NewTool("pm_agreements",
			mcp.WithDescription("Team working agreements (rules team commits to). Actions: list, add."),
			mcp.WithString("action", mcp.Description("list, add (default: list)")),
			mcp.WithString("agreement", mcp.Description("Agreement text (for add)")),
			mcp.WithString("why", mcp.Description("Why this agreement exists")),
		),
		h.ManageAgreements,
	)

	s.AddTool(
		mcp.NewTool("pm_experiment",
			mcp.WithDescription("Record improvement experiment: hypothesis, action, measurement, duration."),
			mcp.WithString("hypothesis", mcp.Required(), mcp.Description("What we think will improve (e.g. 'reducing WIP will decrease cycle time')")),
			mcp.WithString("action", mcp.Required(), mcp.Description("What we will try (e.g. 'limit WIP to 2 per person')")),
			mcp.WithString("measure", mcp.Description("How we'll know it worked (default: observe)")),
			mcp.WithString("duration", mcp.Description("How long to run (default: 1 sprint)")),
			mcp.WithString("sprint_name", mcp.Description("Sprint context")),
			mcp.WithString("context", mcp.Description("What prompted this experiment")),
		),
		h.RecordExperiment,
	)

	s.AddTool(
		mcp.NewTool("pm_experiments",
			mcp.WithDescription("Show all improvement experiments and their status."),
		),
		h.ReviewExperiments,
	)

	s.AddTool(
		mcp.NewTool("pm_planning_prep",
			mcp.WithDescription("Sprint planning prep: outcome, capacity, carryover, risks, deps."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.SprintPlanningSummary,
	)
}

func registerTechDebtTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_tech_debt_add",
			mcp.WithDescription("Record tech debt: code shortcuts, architectural issues, testing gaps."),
			mcp.WithString("title", mcp.Required(), mcp.Description("Tech debt title")),
			mcp.WithString("description", mcp.Description("Detailed description")),
			mcp.WithString("impact", mcp.Description("high (blocks velocity), medium (slows), low (cosmetic). Default: medium")),
			mcp.WithString("category", mcp.Description("code, architecture, testing, infra, docs. Default: code")),
			mcp.WithString("owner", mcp.Description("Who should fix this")),
			mcp.WithString("fix_approach", mcp.Description("How to fix it")),
		),
		h.RecordTechDebt,
	)

	s.AddTool(
		mcp.NewTool("pm_tech_debt",
			mcp.WithDescription("Tech debt dashboard: all items by impact, resolution velocity, status."),
		),
		h.TechDebtDashboard,
	)

	s.AddTool(
		mcp.NewTool("pm_tech_debt_budget",
			mcp.WithDescription("Recommend sprint capacity for tech debt based on debt load and velocity."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.TechDebtBudget,
	)

	s.AddTool(
		mcp.NewTool("pm_review_prep",
			mcp.WithDescription("Sprint review prep: demo order, talking points, shipped vs not."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.SprintReviewPrep,
	)

	s.AddTool(
		mcp.NewTool("pm_mcp_stats",
			mcp.WithDescription("MCP server self-monitoring: memory contents, data freshness, storage stats."),
			mcp.WithNumber("board_id", mcp.Description("Board ID for sprint-specific stats")),
		),
		h.MCPStats,
	)
}
