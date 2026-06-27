package tools

import (
	"context"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func containsAny(s string, terms ...string) bool {
	for _, t := range terms { if strings.Contains(s, t) { return true } }
	return false
}

func (h *Handlers) PMSmart(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ask, err := req.RequireString("ask")
	if err != nil { return errorResult("ask something"), nil }
	lower := strings.ToLower(ask)
	switch {
	case containsAny(lower, "blocker", "impediment", "blocked"): return h.GetBlockers(ctx, req)
	case containsAny(lower, "risk"): return h.GetRiskDashboard(ctx, req)
	case containsAny(lower, "health"): return h.SprintHealthScore(ctx, req)
	case containsAny(lower, "status", "dashboard", "overview"): return h.PMDashboard(ctx, req)
	case containsAny(lower, "standup"): return h.StandupPrep(ctx, req)
	case containsAny(lower, "forecast", "when will"): return h.MonteCarloForecast(ctx, req)
	case containsAny(lower, "velocity"): return h.VelocityTrend(ctx, req)
	case containsAny(lower, "my issue", "assigned"): return h.MyIssues(ctx, req)
	case containsAny(lower, "action item"): return h.GetActionItems(ctx, req)
	case containsAny(lower, "workload"): return h.Workload(ctx, req)
	case containsAny(lower, "sentiment", "mood", "morale", "team feel", "team happ"): return h.PMTeamPulseHistory(ctx, req)
	case containsAny(lower, "okr", "objective", "key result", "goal align"): return h.PMOKRSuggest(ctx, req)
	case containsAny(lower, "nudge", "follow.?up", "ping", "remind"): return h.CommsHealth(ctx, req)
	case containsAny(lower, "help"): return h.PMHelp(ctx, req)
	default:
		if h.AI != nil {
			prompt := "You are a PM assistant. The user asked: '" + ask + "'. " +
				"Available tools handle: blockers, risks, sprint health, dashboard, standup, forecast, velocity, my issues, " +
				"action items, workload, sentiment, OKRs, coaching, planning, capacity, retro, experiments, decisions, and more. " +
				"If their question matches one of these, briefly respond and suggest the right tool name. " +
				"Otherwise answer concisely in 2-3 sentences. Be helpful, not robotic."
			r, e := h.aiComplete(ctx, prompt, ask)
			if e == nil { return mcp.NewToolResultText(r), nil }
		}
		return h.PMHelp(ctx, req)
	}
}

func (h *Handlers) PMDo(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	what, err := req.RequireString("what")
	if err != nil { return errorResult("what to do?"), nil }
	lower := strings.ToLower(what)
	switch {
	case containsAny(lower, "create", "new issue", "new task"): return h.CreateIssue(ctx, req)
	case containsAny(lower, "risk"): return h.RecordRisk(ctx, req)
	case containsAny(lower, "decision"): return h.RecordDecision(ctx, req)
	case containsAny(lower, "block"): return h.RecordBlocker(ctx, req)
	case containsAny(lower, "context note"): return h.TeamContext(ctx, req)
	case containsAny(lower, "feedback", "gave feedback"): return h.GiveFeedback(ctx, req)
	case containsAny(lower, "learning", "tribal"): return h.RecordLearning(ctx, req)
	case containsAny(lower, "retro", "retrospective"): return h.RecordRetrospective(ctx, req)
	case containsAny(lower, "hypothesis", "experiment"): return h.RecordExperiment(ctx, req)
	case containsAny(lower, "kpi", "measurement"): return h.PMKPISnapshot(ctx, req)
	default: return errorResult("Try: create issue, record risk, record decision, record blocker, record feedback, record retro, record learning"), nil
	}
}

func (h *Handlers) PMReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t, _ := req.RequireString("type")
	switch strings.ToLower(t) {
	case "status": return h.PMDashboard(ctx, req)
	case "executive", "exec": return h.ExecutiveReport(ctx, req)
	case "release_notes": return h.GenerateReleaseNotes(ctx, req)
	case "weekly", "digest": return h.WeeklyDigest(ctx, req)
	case "health": return h.SprintHealthScore(ctx, req)
	case "velocity": return h.VelocityTrend(ctx, req)
	case "scorecard": return h.SprintScorecard(ctx, req)
	case "sentiment", "mood", "morale": return h.PMTeamPulseHistory(ctx, req)
	case "okr", "okr_health": return h.PMOKRHealth(ctx, req)
	case "kpi", "kpi_dashboard": return h.PMKPIDashboard(ctx, req)
	case "coaching": return h.CoachingAdvice(ctx, req)
	default: return errorResult("type: status, executive, release_notes, weekly, health, velocity, scorecard, sentiment, okr, kpi, coaching"), nil
	}
}

func (h *Handlers) PMTeam(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a, _ := req.RequireString("action")
	switch strings.ToLower(a) {
	case "workload": return h.Workload(ctx, req)
	case "1on1": return h.OneOnOnePrep(ctx, req)
	case "pulse": return h.PMTeamPulseHistory(ctx, req)
	case "radar": return h.PMTeamRadarHistory(ctx, req)
	case "health": return h.GetTeamHealth(ctx, req)
	default: return errorResult("action: workload, 1on1, pulse, radar, health"), nil
	}
}

func (h *Handlers) PMPlan(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a, _ := req.RequireString("action")
	switch strings.ToLower(a) {
	case "planning", "prep": return h.SprintPlanningSummary(ctx, req)
	case "capacity": return h.CapacityPlan(ctx, req)
	case "forecast": return h.MonteCarloForecast(ctx, req)
	case "ready": return h.CheckStoryReady(ctx, req)
	case "backlog": return h.BacklogGroom(ctx, req)
	case "goal": return h.GetSprintGoals(ctx, req)
	default: return errorResult("action: planning, capacity, forecast, ready, backlog, goal"), nil
	}
}

func (h *Handlers) PMRetro(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	a, _ := req.RequireString("action")
	switch strings.ToLower(a) {
	case "actions": return h.GetActionItems(ctx, req)
	case "experiments": return h.ReviewExperiments(ctx, req)
	case "anti_patterns": return h.DetectAntiPatterns(ctx, req)
	case "coaching": return h.CoachingAdvice(ctx, req)
	case "facilitate": return h.CeremonyFacilitator(ctx, req)
	default: return errorResult("action: actions, experiments, anti_patterns, coaching, facilitate"), nil
	}
}

func (h *Handlers) PMSearch(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	q, _ := req.RequireString("query")
	if q == "" { return errorResult("query required"), nil }
	lower := strings.ToLower(q)

	switch {
	case containsAny(lower, "decision", "decided", "keputusan"): return h.SearchDecisions(ctx, req)
	case containsAny(lower, "risk", "resiko", "mitigasi"): return h.GetRiskDashboard(ctx, req)
	case containsAny(lower, "blocker", "blocked", "impediment", "hambatan"): return h.GetBlockers(ctx, req)
	case containsAny(lower, "action", "todo", "pending", "tindak lanjut"): return h.GetActionItems(ctx, req)
	case containsAny(lower, "meeting", "rapat", "notulen"): return h.GetMeetings(ctx, req)
	case containsAny(lower, "kb", "knowledge", "learning", "tribal", "lesson", "pengetahuan"): return h.TeamKnowledgeBase(ctx, req)
	case containsAny(lower, "overdue", "stale", "kadaluarsa"): return h.Overdue(ctx, req)
	}

	// Fallback: AI converts to JQL, or raw Jira search
	if h.AI != nil { return h.NLToJQL(ctx, req) }
	return h.SearchIssues(ctx, req)
}
