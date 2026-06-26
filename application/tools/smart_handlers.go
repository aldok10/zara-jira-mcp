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
	case containsAny(lower, "help"): return h.PMHelp(ctx, req)
	default:
		if h.AI != nil {
			r, e := h.AI.Complete(ctx, "Concise PM assistant. 2-3 sentences.", ask)
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
	default: return errorResult("Try: create issue, record risk, record decision, record blocker"), nil
	}
}

func (h *Handlers) PMReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	t, _ := req.RequireString("type")
	switch strings.ToLower(t) {
	case "status": return h.PMDashboard(ctx, req)
	case "executive", "exec": return h.ExecutiveReport(ctx, req)
	case "release_notes": return h.GenerateReleaseNotes(ctx, req)
	case "weekly": return h.WeeklyDigest(ctx, req)
	case "health": return h.SprintHealthScore(ctx, req)
	case "velocity": return h.VelocityTrend(ctx, req)
	case "scorecard": return h.SprintScorecard(ctx, req)
	default: return errorResult("type: status, executive, release_notes, weekly, health, velocity, scorecard"), nil
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
	if containsAny(strings.ToLower(q), "decision") { return h.SearchDecisions(ctx, req) }
	if h.AI != nil { return h.NLToJQL(ctx, req) }
	return h.SearchIssues(ctx, req)
}
