package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	smcp "github.com/aldok10/zara-jira-mcp/modules/sprint/interfaces/mcp"
)

func RegisterSprintTools(s *server.MCPServer, h *smcp.Handlers) {
	s.AddTool(
		mcp.NewTool("pm",
			mcp.WithDescription("Quick project status. Shows sprint progress, blockers, risks, pending actions."),
			mcp.WithNumber("board_id", mcp.Description("Board ID")),
		),
		h.PMQuickStatus,
	)
	s.AddTool(
		mcp.NewTool("pm_create",
			mcp.WithDescription("Create work item anywhere: Jira, GitHub, or GitLab."),
			mcp.WithString("title", mcp.Required(), mcp.Description("What needs to be done")),
			mcp.WithString("description", mcp.Description("Details")),
			mcp.WithString("project", mcp.Description("Project key")),
			mcp.WithString("assignee", mcp.Description("Who should do this")),
			mcp.WithString("type", mcp.Description("Task, Bug, Story")),
			mcp.WithString("priority", mcp.Description("High, Medium, Low")),
			mcp.WithString("labels", mcp.Description("Comma-separated labels")),
			mcp.WithString("platform", mcp.Description("jira (default), github, gitlab")),
		),
		h.PMCreate,
	)
	s.AddTool(
		mcp.NewTool("pm_decide",
			mcp.WithDescription("Quick decision recording. Just say what was decided."),
			mcp.WithString("what", mcp.Required(), mcp.Description("What was decided")),
			mcp.WithString("who", mcp.Description("Who decided (default: team)")),
			mcp.WithString("why", mcp.Description("Why (optional)")),
		),
		h.PMDecide,
	)
	s.AddTool(
		mcp.NewTool("pm_risk",
			mcp.WithDescription("Quick risk recording. Just say what could go wrong."),
			mcp.WithString("what", mcp.Required(), mcp.Description("What could go wrong")),
			mcp.WithString("severity", mcp.Description("critical, high, medium, low")),
			mcp.WithString("owner", mcp.Description("Who should handle this")),
		),
		h.PMRisk,
	)
	s.AddTool(
		mcp.NewTool("pm_next",
			mcp.WithDescription("Suggest next high-priority PM action based on memory state."),
			mcp.WithNumber("board_id", mcp.Description("Board ID")),
		),
		h.PMNext,
	)
	s.AddTool(
		mcp.NewTool("pm_snapshot",
			mcp.WithDescription("Snapshot sprint state to memory. Auto-calculates from Jira."),
			mcp.WithNumber("board_id", mcp.Description("Board ID")),
			mcp.WithString("sprint_name", mcp.Description("Sprint name")),
			mcp.WithNumber("total_issues", mcp.Description("Total issues")),
			mcp.WithNumber("done", mcp.Description("Completed issues")),
			mcp.WithNumber("in_progress", mcp.Description("In progress")),
			mcp.WithNumber("todo", mcp.Description("To do")),
			mcp.WithNumber("blocked", mcp.Description("Blocked")),
			mcp.WithNumber("carryover", mcp.Description("Carried from previous sprint")),
			mcp.WithNumber("velocity", mcp.Description("Story points completed")),
			mcp.WithString("notes", mcp.Description("Any notes")),
		),
		h.PMSnapshotSprint,
	)
	s.AddTool(
		mcp.NewTool("pm_record_decision",
			mcp.WithDescription("Record a project decision with context and rationale."),
			mcp.WithString("title", mcp.Required(), mcp.Description("Decision title")),
			mcp.WithString("decision", mcp.Required(), mcp.Description("What was decided")),
			mcp.WithString("context", mcp.Description("What situation led to this")),
			mcp.WithString("rationale", mcp.Description("Why this over alternatives")),
			mcp.WithString("made_by", mcp.Description("Who decided")),
			mcp.WithString("tags", mcp.Description("Comma-separated tags")),
		),
		h.PMRecordDecision,
	)
	s.AddTool(
		mcp.NewTool("pm_record_risk",
			mcp.WithDescription("Record a project risk with mitigation."),
			mcp.WithString("title", mcp.Required(), mcp.Description("Risk title")),
			mcp.WithString("description", mcp.Description("Detailed description")),
			mcp.WithString("severity", mcp.Description("critical, high, medium, low")),
			mcp.WithString("owner", mcp.Description("Who owns mitigating this")),
			mcp.WithString("mitigation", mcp.Description("Mitigation plan")),
		),
		h.PMRecordRisk,
	)
	s.AddTool(
		mcp.NewTool("pm_record_blocker",
			mcp.WithDescription("Record an impediment/blocker."),
			mcp.WithString("description", mcp.Required(), mcp.Description("What is blocked and why")),
			mcp.WithString("issue_key", mcp.Description("Related Jira issue key")),
			mcp.WithString("owner", mcp.Description("Who resolves this")),
		),
		h.PMRecordBlocker,
	)
	s.AddTool(
		mcp.NewTool("pm_record_retro",
			mcp.WithDescription("Record a sprint retrospective."),
			mcp.WithString("sprint_name", mcp.Required(), mcp.Description("Sprint name")),
			mcp.WithString("went_well", mcp.Description("What went well")),
			mcp.WithString("improvements", mcp.Description("What needs improvement")),
			mcp.WithString("action_items", mcp.Description("Newline-separated action items")),
		),
		h.PMRecordRetro,
	)
	s.AddTool(
		mcp.NewTool("pm_record_meeting",
			mcp.WithDescription("Record meeting notes with decisions and actions."),
			mcp.WithString("meeting_type", mcp.Required(), mcp.Description("standup, planning, retro, grooming, adhoc")),
			mcp.WithString("notes", mcp.Description("Key discussion points")),
			mcp.WithString("attendees", mcp.Description("Comma-separated attendee names")),
			mcp.WithString("decisions", mcp.Description("Decisions made")),
			mcp.WithString("action_items", mcp.Description("Follow-up actions")),
			mcp.WithString("sprint_name", mcp.Description("Sprint context")),
		),
		h.PMRecordMeeting,
	)
	s.AddTool(
		mcp.NewTool("pm_risks",
			mcp.WithDescription("Show risk dashboard — all risks by severity."),
		),
		h.PMRiskDashboard,
	)
	s.AddTool(
		mcp.NewTool("pm_blockers",
			mcp.WithDescription("Show active blockers or blocker history."),
			mcp.WithBoolean("show_history", mcp.Description("Show resolved blockers too (default: false)")),
		),
		h.PMBlockers,
	)
	s.AddTool(
		mcp.NewTool("pm_decisions",
			mcp.WithDescription("Show recent decisions with context."),
			mcp.WithNumber("limit", mcp.Description("Max results (default 10)")),
		),
		h.PMDecisions,
	)
	s.AddTool(
		mcp.NewTool("pm_actions",
			mcp.WithDescription("Show pending action items from retrospectives."),
		),
		h.PMActionItems,
	)
	s.AddTool(
		mcp.NewTool("pm_dependencies",
			mcp.WithDescription("Show dependency map — all open dependencies or for a specific issue."),
			mcp.WithString("issue_key", mcp.Description("Filter by issue key")),
		),
		h.PMDependencies,
	)
	s.AddTool(
		mcp.NewTool("pm_health",
			mcp.WithDescription("Show sprint health history."),
			mcp.WithNumber("board_id", mcp.Description("Board ID")),
		),
		h.PMSprintHealth,
	)
	s.AddTool(
		mcp.NewTool("pm_forecast",
			mcp.WithDescription("Monte Carlo forecast: 10K simulations from historical throughput. Predicts completion sprints at 50/70/85/95% confidence."),
			mcp.WithNumber("board_id", mcp.Description("Board ID")),
			mcp.WithNumber("remaining_items", mcp.Description("Items remaining (default: from active sprint)")),
		),
		h.PMForecast,
	)
	s.AddTool(
		mcp.NewTool("pm_velocity_trend",
			mcp.WithDescription("Velocity and completion trends over recent sprints. Detects improvement or decline."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.PMVelocityTrend,
	)
	s.AddTool(
		mcp.NewTool("pm_anti_patterns",
			mcp.WithDescription("Detect Scrum anti-patterns: zombie sprints, scope creep, blocked issues, inconsistent delivery, declining velocity."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.PMAntiPatterns,
	)
	s.AddTool(
		mcp.NewTool("pm_flow_metrics",
			mcp.WithDescription("Flow metrics: WIP, throughput, cycle time, completion rate. Detects flow problems and bottlenecks."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.PMFlowMetrics,
	)
	s.AddTool(
		mcp.NewTool("pm_sprint_compare",
			mcp.WithDescription("Compare current vs previous sprint: issues, completion, carryover, velocity."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.PMSprintCompare,
	)
	s.AddTool(
		mcp.NewTool("pm_predictability",
			mcp.WithDescription("Sprint predictability score (0-100). How consistent is delivery across recent sprints?"),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.PMPredictability,
	)
	s.AddTool(
		mcp.NewTool("pm_scorecard",
			mcp.WithDescription("Sprint scorecard (0-100): completion, velocity, blocked ratio, carryover, predictability. Comprehensive sprint health."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.PMScorecard,
	)
	s.AddTool(
		mcp.NewTool("pm_calibration",
			mcp.WithDescription("Forecast accuracy: committed vs delivered over time. Shows on-target rate, over-commit, and under-deliver patterns."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.PMCalibration,
	)
	s.AddTool(
		mcp.NewTool("pm_set_sprint_goal",
			mcp.WithDescription("Define sprint goal with measurable key results. Tracks achievement across sprints."),
			mcp.WithString("goal", mcp.Required(), mcp.Description("Sprint goal statement")),
			mcp.WithString("key_results", mcp.Description("Measurable key results (newline-separated)")),
			mcp.WithNumber("board_id", mcp.Description("Board ID")),
			mcp.WithString("sprint_name", mcp.Description("Sprint name (default: current)")),
		),
		h.PMSetSprintGoal,
	)
	s.AddTool(
		mcp.NewTool("pm_sprint_goals",
			mcp.WithDescription("Show active sprint goals or goal achievement history."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
			mcp.WithBoolean("show_history", mcp.Description("Show past goals with outcomes (default: false)")),
		),
		h.PMSprintGoals,
	)
	s.AddTool(
		mcp.NewTool("pm_goal_check",
			mcp.WithDescription("AI sprint goal progress check. Evaluates if on track from current health and sprint data."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.PMGoalCheck,
	)
	s.AddTool(
		mcp.NewTool("pm_track_daily",
			mcp.WithDescription("Track today's sprint progress. Captures burndown data from Jira and saves to memory."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.PMTrackDaily,
	)
	s.AddTool(
		mcp.NewTool("pm_burndown",
			mcp.WithDescription("Show sprint burndown chart with daily progress tracking, ideal line, and on-track assessment."),
			mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
		),
		h.PMBurndown,
	)
}
