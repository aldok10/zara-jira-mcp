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
}
