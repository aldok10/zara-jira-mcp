package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerCoachingTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_team_pulse",
		mcp.WithDescription("Record team health pulse survey. Track morale trends across sprints."),
		mcp.WithString("sprint_name", mcp.Required(), mcp.Description("Sprint name")),
		mcp.WithString("ratings", mcp.Required(), mcp.Description("JSON: {\"member_name\": score(1-5), ...}")),
		mcp.WithString("notes", mcp.Description("Optional context")),
	), h.PMTeamPulse)

	s.AddTool(mcp.NewTool("pm_team_pulse_history",
		mcp.WithDescription("Show team pulse trends over time. See morale trajectory."),
	), h.PMTeamPulseHistory)

	s.AddTool(mcp.NewTool("pm_predictability",
		mcp.WithDescription("Sprint predictability score (0-100). How consistent is your team's delivery?"),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.PMPredictability)

	s.AddTool(mcp.NewTool("pm_meeting_effectiveness",
		mcp.WithDescription("Rate ceremony effectiveness. Track which meetings add value and which waste time."),
		mcp.WithString("ceremony", mcp.Required(), mcp.Description("standup, planning, retro, review, grooming")),
		mcp.WithNumber("duration_minutes", mcp.Required(), mcp.Description("How long the meeting lasted")),
		mcp.WithNumber("score", mcp.Required(), mcp.Description("Effectiveness 1-5 (1=waste, 5=excellent)")),
		mcp.WithString("notes", mcp.Description("What made it good or bad")),
		mcp.WithString("sprint_name", mcp.Description("Sprint context")),
	), h.PMMeetingEffectiveness)

	s.AddTool(mcp.NewTool("pm_meeting_trends",
		mcp.WithDescription("Show meeting effectiveness trends. Which ceremonies are improving or degrading?"),
		mcp.WithString("ceremony", mcp.Description("Filter by ceremony type (all if empty)")),
	), h.PMMeetingTrends)

	s.AddTool(mcp.NewTool("pm_team_radar",
		mcp.WithDescription("Multi-dimension team health assessment (like Spotify Health Check). Track delivery, quality, fun, learning, teamwork, speed, mission."),
		mcp.WithString("sprint_name", mcp.Required(), mcp.Description("Sprint name")),
		mcp.WithString("dimensions", mcp.Required(), mcp.Description("JSON: {\"delivery\": 4, \"quality\": 3, \"fun\": 5, \"learning\": 3, \"teamwork\": 4, \"speed\": 3, \"mission\": 4}")),
	), h.PMTeamRadar)

	s.AddTool(mcp.NewTool("pm_team_radar_history",
		mcp.WithDescription("Show team radar dimension trends across sprints."),
	), h.PMTeamRadarHistory)

	s.AddTool(mcp.NewTool("pm_maturity_assessment",
		mcp.WithDescription("AI-powered agile maturity assessment (1-5 scale). Combines predictability, process discipline, continuous improvement, autonomy, stakeholder trust."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.PMMaturityAssessment)

	s.AddTool(mcp.NewTool("pm_daily_digest",
		mcp.WithDescription("Auto-generated morning brief: blockers, sprint status, risks, pending actions, dependencies. Everything needing attention in one view."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.PMDailyDigestCoaching)
}
