package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MCPServer struct {
	s *server.MCPServer
}

func NewMCPServer(handlers *tools.Handlers) *MCPServer {
	s := server.NewMCPServer(
		"zara-jira-mcp",
		"0.3.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	registerJiraTools(s, handlers)
	registerIssueOpsTools(s, handlers)
	registerPMTools(s, handlers)
	registerAITools(s, handlers)
	registerLarkTools(s, handlers)
	registerSlackTools(s, handlers)
	registerMemoryTools(s, handlers)
	registerPMIntelTools(s, handlers)
	registerAdvancedPMTools(s, handlers)
	registerDeepPMTools(s, handlers)
	registerEpicSprintTools(s, handlers)
	registerBulkProjectTools(s, handlers)
	registerLinkWorklogTools(s, handlers)
	registerFlowTools(s, handlers)
	registerForecastTools(s, handlers)
	registerStakeholderTools(s, handlers)
	registerProcessTools(s, handlers)
	registerRecipeTools(s, handlers)
	registerPortfolioTools(s, handlers)
	registerGitHubTools(s, handlers)

	return &MCPServer{s: s}
}

func (m *MCPServer) Server() *server.MCPServer {
	return m.s
}

func registerJiraTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("jira_search",
			mcp.WithDescription("Search Jira issues using JQL. Returns key, summary, status, priority, assignee."),
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
			mcp.WithDescription("List available status transitions for an issue. Use transition IDs with jira_transition."),
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
			mcp.WithDescription("AI-powered analysis of Jira tickets. Ask questions like: 'What are the blockers?', 'Which tickets are stale?', 'Sprint health?'"),
			mcp.WithString("query", mcp.Required(), mcp.Description("Your question about the project/tickets")),
			mcp.WithString("jql", mcp.Description("JQL to scope the analysis")),
			mcp.WithNumber("max_results", mcp.Description("Max tickets to analyze (default 30)")),
		),
		h.AIAnalyze,
	)

	s.AddTool(
		mcp.NewTool("jira_ai_sprint_report",
			mcp.WithDescription("Generate AI-powered sprint report with health assessment and recommendations. Optionally sends to Lark."),
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
			mcp.WithDescription("Capture current sprint state into PM memory. Auto-calculates done/in-progress/todo/blocked from Jira. Call at end of each sprint for velocity tracking."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithNumber("velocity", mcp.Description("Story points completed this sprint (manual input)")),
			mcp.WithNumber("carryover", mcp.Description("Issues carried over from previous sprint")),
			mcp.WithString("notes", mcp.Description("Any notes about this sprint snapshot")),
		),
		h.SnapshotSprint,
	)

	s.AddTool(
		mcp.NewTool("pm_record_risk",
			mcp.WithDescription("Record a project risk to the risk register. Track risks, their severity, owners, and mitigation plans."),
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
			mcp.WithDescription("Record a project decision with context and rationale. Build institutional memory."),
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
			mcp.WithDescription("AI-powered PM recommendations based on ALL historical memory: sprint trends, risks, blockers, team metrics, decisions. Your AI Scrum Master brain."),
			mcp.WithNumber("board_id", mcp.Description("Board ID for sprint context")),
			mcp.WithString("focus", mcp.Description("Focus area: general, velocity, risks, team, process (default: general)")),
		),
		h.PMRecommendations,
	)

	s.AddTool(
		mcp.NewTool("pm_velocity_trend",
			mcp.WithDescription("Show velocity and completion trends over recent sprints. Detect if team is improving, stable, or declining."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.VelocityTrend,
	)

	s.AddTool(
		mcp.NewTool("pm_standup_prep",
			mcp.WithDescription("Generate daily standup preparation brief. Combines live Jira data with historical blockers, risks, and action items into talking points."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.StandupPrep,
	)

	s.AddTool(
		mcp.NewTool("pm_retro_analysis",
			mcp.WithDescription("AI analysis of sprint patterns across retrospectives: recurring issues, trend detection, root cause patterns, improvement suggestions."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.SprintRetroAnalysis,
	)
}

func registerAdvancedPMTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_sprint_health",
			mcp.WithDescription("Compute sprint health score (0-100) with breakdown: velocity, blockers, scope creep, team balance. Saves to history for trend tracking."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.SprintHealthScore,
	)

	s.AddTool(
		mcp.NewTool("pm_health_history",
			mcp.WithDescription("Show health score trends over time. See if the team is getting healthier or declining."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.HealthHistory,
	)

	s.AddTool(
		mcp.NewTool("pm_record_dependency",
			mcp.WithDescription("Record a dependency between issues (blocks, blocked_by, external). Track cross-team dependencies."),
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
			mcp.WithDescription("Capacity planning based on velocity history. Calculates recommended story points for next sprint based on team availability."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithNumber("team_size", mcp.Description("Number of team members")),
			mcp.WithNumber("sprint_days", mcp.Description("Sprint duration in days (default: 10)")),
			mcp.WithNumber("planned_leave_days", mcp.Description("Total planned leave days across team")),
		),
		h.CapacityPlan,
	)

	s.AddTool(
		mcp.NewTool("pm_auto_detect_risks",
			mcp.WithDescription("Proactively scan for risk signals: stale tickets, overloaded members, chronic blockers, overdue actions. Auto-records findings to risk register."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.AutoDetectRisks,
	)
}

func registerDeepPMTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_track_daily",
			mcp.WithDescription("Track today's sprint progress (burndown data point). Call daily for burndown chart data."),
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
			mcp.WithDescription("Define sprint goal with measurable key results. Track if the team achieves what matters."),
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
			mcp.WithDescription("Auto-escalate to Lark: critical risks >3 days, blockers >3 days, sprint health <50. Sends alert to configured Lark group."),
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
			mcp.WithDescription("One-shot PM dashboard: sprint progress, health score, risks, blockers, dependencies, goals, actions, escalations. Everything in one view."),
			mcp.WithNumber("board_id", mcp.Description("Board ID for sprint-specific data")),
		),
		h.PMDashboard,
	)

	s.AddTool(
		mcp.NewTool("pm_release_notes",
			mcp.WithDescription("Generate categorized release notes from completed sprint issues (features, bugs, tasks). Optionally sends to Lark."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithBoolean("send_to_lark", mcp.Description("Send release notes to Lark group")),
		),
		h.GenerateReleaseNotes,
	)
}

func registerFlowTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_flow_metrics",
			mcp.WithDescription("Calculate flow metrics from live Jira data: WIP, throughput, cycle time, lead time. Detects flow problems (high WIP, long cycle time, low throughput)."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.FlowMetrics,
	)

	s.AddTool(
		mcp.NewTool("pm_sprint_compare",
			mcp.WithDescription("Compare current sprint vs previous: velocity, completion, blockers, carryover. Shows if team is improving or declining."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.SprintComparison,
	)

	s.AddTool(
		mcp.NewTool("pm_facilitate",
			mcp.WithDescription("AI ceremony facilitator — generates fresh, contextual facilitation prompts for any Scrum ceremony."),
			mcp.WithString("ceremony", mcp.Required(), mcp.Description("standup, planning, retro, grooming, review")),
			mcp.WithNumber("board_id", mcp.Description("Board ID for context (optional)")),
		),
		h.CeremonyFacilitator,
	)

	s.AddTool(
		mcp.NewTool("pm_confidence",
			mcp.WithDescription("Record team/individual confidence level for current sprint (1-5 scale). Track pre-sprint confidence vs actual outcome."),
			mcp.WithString("sprint_name", mcp.Required(), mcp.Description("Sprint name")),
			mcp.WithNumber("score", mcp.Required(), mcp.Description("1=very worried, 2=worried, 3=neutral, 4=confident, 5=very confident")),
			mcp.WithString("member", mcp.Description("Team member name (default: 'team')")),
			mcp.WithString("note", mcp.Description("Why this confidence level?")),
		),
		h.RecordConfidence,
	)

	s.AddTool(
		mcp.NewTool("pm_goal_check",
			mcp.WithDescription("AI-powered sprint goal progress check. Evaluates if goal is on track based on current data."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.SprintGoalCheck,
	)
}

func registerForecastTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_forecast",
			mcp.WithDescription("Monte Carlo forecasting: 'When will it be done?' Runs 10,000 simulations using historical throughput. Returns probability-based dates (50%, 70%, 85%, 95% confidence)."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithNumber("remaining_items", mcp.Description("Items remaining (default: from active sprint)")),
			mcp.WithNumber("sprint_days", mcp.Description("Sprint length in days (default: 10)")),
		),
		h.MonteCarloForecast,
	)

	s.AddTool(
		mcp.NewTool("pm_anti_patterns",
			mcp.WithDescription("Detect Scrum anti-patterns from data: zombie sprints, hero culture, scope creep, unpredictable velocity, dead retros, rubber-stamp DoD, no sprint goals."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.DetectAntiPatterns,
	)

	s.AddTool(
		mcp.NewTool("pm_coaching",
			mcp.WithDescription("AI coaching advice for Scrum Masters. Get data-driven suggestions for team improvement."),
			mcp.WithString("topic", mcp.Required(), mcp.Description("team_dynamics, velocity, blockers, morale, conflict, growth, predictability")),
			mcp.WithNumber("board_id", mcp.Description("Board ID for data context")),
			mcp.WithString("situation", mcp.Description("Describe the specific situation you need advice on")),
		),
		h.CoachingAdvice,
	)

	s.AddTool(
		mcp.NewTool("pm_forecast_sprint",
			mcp.WithDescription("Monte Carlo sprint forecast. Predicts completion probability based on historical throughput. Shows 50%/85% confidence intervals."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithNumber("items_remaining", mcp.Description("Items left to do (auto-detected from sprint if empty)")),
		),
		h.ForecastSprint,
	)

	s.AddTool(
		mcp.NewTool("pm_nl_to_jql",
			mcp.WithDescription("Convert natural language to JQL. Ask in plain English, get JQL query back."),
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
			mcp.WithDescription("Find stale backlog items that need grooming or archiving. Items untouched for N days outside active sprints."),
			mcp.WithString("project", mcp.Description("Project key filter")),
			mcp.WithNumber("days", mcp.Description("Days without update to consider stale (default: 90)")),
		),
		h.BacklogGroom,
	)
}

func registerRecipeTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_recipe_start_work",
		mcp.WithDescription("Start working on an issue: assigns to you, transitions to In Progress, suggests git branch name. One-click workflow."),
		mcp.WithString("key", mcp.Required(), mcp.Description("Issue key")),
		mcp.WithString("assignee_id", mcp.Description("Your account ID (use jira_find_user to look up)")),
	), h.RecipeStartWork)

	s.AddTool(mcp.NewTool("pm_recipe_done",
		mcp.WithDescription("Mark issue as done: transitions to Done, optionally logs time and adds completion comment."),
		mcp.WithString("key", mcp.Required(), mcp.Description("Issue key")),
		mcp.WithString("time_spent", mcp.Description("Time to log (e.g. '2h', '30m')")),
		mcp.WithString("comment", mcp.Description("Completion comment")),
	), h.RecipeDone)

	s.AddTool(mcp.NewTool("pm_recipe_block",
		mcp.WithDescription("Flag an issue as blocked: records blocker in memory, adds comment to issue, creates impediment trail."),
		mcp.WithString("key", mcp.Required(), mcp.Description("Issue key")),
		mcp.WithString("reason", mcp.Required(), mcp.Description("Why is it blocked?")),
		mcp.WithString("owner", mcp.Description("Who should resolve this?")),
	), h.RecipeBlock)
}

func registerStakeholderTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(
		mcp.NewTool("pm_exec_report",
			mcp.WithDescription("Generate executive stakeholder report. Business outcomes, risks, team health — no jargon, no story points. Written for a VP with 30 seconds."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithBoolean("send_to_lark", mcp.Description("Send to Lark group")),
		),
		h.ExecutiveReport,
	)

	s.AddTool(
		mcp.NewTool("pm_scorecard",
			mcp.WithDescription("Sprint scorecard: quantified 0-100 grade across completion, goal focus, predictability, quality, team balance."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.SprintScorecard,
	)

	s.AddTool(
		mcp.NewTool("pm_team_kb",
			mcp.WithDescription("Team knowledge base: DoD, decisions, patterns, metrics, retro themes. Ask questions or browse. Great for onboarding."),
			mcp.WithNumber("board_id", mcp.Description("Board ID for metrics context")),
			mcp.WithString("question", mcp.Description("Ask a question about how the team works (AI-powered answer)")),
		),
		h.TeamKnowledgeBase,
	)

	s.AddTool(
		mcp.NewTool("pm_record_learning",
			mcp.WithDescription("Record a team learning or tribal knowledge. Builds searchable institutional memory."),
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
			mcp.WithDescription("Generate weekly team digest: decisions, risks, blockers, wins, concerns. AI-summarized from all week's activity."),
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
			mcp.WithDescription("AI evaluates if a Jira story meets Definition of Ready + INVEST criteria. Gives READY/NOT READY verdict."),
			mcp.WithString("key", mcp.Required(), mcp.Description("Jira issue key to evaluate")),
		),
		h.CheckStoryReady,
	)

	s.AddTool(
		mcp.NewTool("pm_agreements",
			mcp.WithDescription("Manage team working agreements. The rules the team commits to. Actions: list, add."),
			mcp.WithString("action", mcp.Description("list, add (default: list)")),
			mcp.WithString("agreement", mcp.Description("Agreement text (for add)")),
			mcp.WithString("why", mcp.Description("Why this agreement exists")),
		),
		h.ManageAgreements,
	)

	s.AddTool(
		mcp.NewTool("pm_experiment",
			mcp.WithDescription("Record an improvement experiment from retro. Track hypothesis, action, measurement, duration."),
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
			mcp.WithDescription("Complete sprint planning preparation package: last sprint outcome, capacity, carryover, risks, dependencies, experiments, checklist."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.SprintPlanningSummary,
	)
}
