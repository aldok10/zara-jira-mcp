package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/mark3labs/mcp-go/mcp"
)

// ManageDoR manages Definition of Ready checklist (entry gate for sprint).
func (h *Handlers) ManageDoR(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	action := req.GetString("action", "list")
	project := req.GetString("project", "*")

	// Reuse DoDItem table with category prefix "dor:" to distinguish
	switch action {
	case "add":
		item := req.GetString("item", "")
		if item == "" {
			return errorResult("item required for add"), nil
		}
		d := &memdom.DoDItem{
			Project:  project,
			Item:     item,
			Category: "dor:" + req.GetString("category", "clarity"),
			OrderNum: req.GetInt("order", 0),
			Active:   true,
		}
		if err := h.Memory.SaveDoDItem(ctx, d); err != nil {
			return errorResult("Failed: " + err.Error()), nil
		}
		return textResult(fmt.Sprintf("DoR item added: %s", item)), nil

	case "remove":
		id := req.GetInt("item_id", 0)
		if id == 0 {
			return errorResult("item_id required"), nil
		}
		if err := h.Memory.DeleteDoDItem(ctx, int64(id)); err != nil {
			return errorResult("Failed: " + err.Error()), nil
		}
		return textResult(fmt.Sprintf("DoR item #%d removed.", id)), nil

	default:
		items, err := h.Memory.GetDoD(ctx, project)
		if err != nil {
			return errorResult("Failed: " + err.Error()), nil
		}

		// Filter DoR items (category starts with "dor:")
		var dorItems []memdom.DoDItem
		for _, item := range items {
			if strings.HasPrefix(item.Category, "dor:") {
				dorItems = append(dorItems, item)
			}
		}

		if len(dorItems) == 0 {
			return textResult(`No Definition of Ready configured. 

Suggested DoR items to add:
- Story has acceptance criteria (clarity)
- Story is estimated (estimation)
- Dependencies identified (dependencies)
- UX/design available if needed (design)
- Story fits in one sprint (size)
- Product Owner available for questions (support)

Use action="add" to build your DoR.`), nil
		}

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("Definition of Ready (project: %s):\n\n", project))
		for _, item := range dorItems {
			cat := strings.TrimPrefix(item.Category, "dor:")
			sb.WriteString(fmt.Sprintf("  #%d [%s] %s\n", item.ID, cat, item.Item))
		}
		return textResult(sb.String()), nil
	}
}

// CheckStoryReady evaluates if a Jira issue meets Definition of Ready / INVEST criteria.
func (h *Handlers) CheckStoryReady(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key required (Jira issue key)"), nil
	}

	issue, err := h.Jira.GetIssue(ctx, key)
	if err != nil {
		return errorResult("Failed to get issue: " + err.Error()), nil
	}

	// Get DoR items
	items, _ := h.Memory.GetDoD(ctx, "*")
	var dorItems []string
	for _, item := range items {
		if strings.HasPrefix(item.Category, "dor:") {
			dorItems = append(dorItems, item.Item)
		}
	}

	dorContext := "No DoR configured (using INVEST criteria as default)"
	if len(dorItems) > 0 {
		dorContext = "Definition of Ready:\n" + strings.Join(dorItems, "\n")
	}

	systemPrompt := `You are evaluating if a user story is READY to enter a sprint.

Check against:
1. INVEST criteria (Independent, Negotiable, Valuable, Estimable, Small, Testable)
2. The team's Definition of Ready (provided below)

For each criterion, give: PASS / FAIL / UNCLEAR with one-line reason.
End with: READY / NOT READY verdict + what's missing.

Be practical — a story doesn't need perfection, just enough clarity to start work without blocking questions.`

	issueData := fmt.Sprintf("Issue: %s\nSummary: %s\nType: %s\nDescription: %s\nLabels: %s\nPriority: %s\n\n%s",
		issue.Key, issue.Summary, issue.Type, issue.Description, strings.Join(issue.Labels, ","), issue.Priority, dorContext)

	analysis, err := h.aiComplete(ctx, systemPrompt, issueData)
	if err != nil {
		return errorResult("AI failed: " + err.Error()), nil
	}

	return textResult(fmt.Sprintf("Story Readiness Check: %s\n\n%s", key, analysis)), nil
}

// ManageAgreements manages team working agreements.
func (h *Handlers) ManageAgreements(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	action := req.GetString("action", "list")

	// Store agreements as decisions with tag "agreement"
	switch action {
	case "add":
		agreement := req.GetString("agreement", "")
		if agreement == "" {
			return errorResult("agreement required"), nil
		}
		d := &memdom.Decision{
			Title:     agreement,
			Decision:  agreement,
			Rationale: req.GetString("why", "team consensus"),
			MadeBy:    "team",
			MadeAt:    time.Now(),
			Tags:      "agreement",
		}
		if err := h.Memory.SaveDecision(ctx, d); err != nil {
			return errorResult("Failed: " + err.Error()), nil
		}
		return textResult(fmt.Sprintf("Agreement added: %s", agreement)), nil

	case "remove":
		// Agreements are immutable decisions — just acknowledge removal
		return textResult("Agreements are stored as decisions. Discuss in retro before removing. Search with pm_search_decisions(query:'agreement')."), nil

	default:
		decisions, err := h.Memory.SearchDecisions(ctx, "agreement")
		if err != nil {
			return errorResult("Failed: " + err.Error()), nil
		}

		if len(decisions) == 0 {
			return textResult(`No working agreements recorded.

Common agreements to consider:
- "We don't start new work when WIP > 3 per person"
- "All PRs reviewed within 24 hours"
- "Blockers raised same day, not in next standup"
- "Sprint scope locked after planning"
- "Tech debt gets 20% of sprint capacity"

Use action="add" to record agreements.`), nil
		}

		var sb strings.Builder
		sb.WriteString("Team Working Agreements:\n\n")
		for i, d := range decisions {
			sb.WriteString(fmt.Sprintf("%d. %s\n   (Why: %s | Since: %s)\n\n", i+1, d.Decision, d.Rationale, d.MadeAt.Format("2006-01-02")))
		}
		return textResult(sb.String()), nil
	}
}

// RecordExperiment tracks an improvement experiment from retro.
func (h *Handlers) RecordExperiment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	hypothesis, err := req.RequireString("hypothesis")
	if err != nil {
		return errorResult("hypothesis required (what do we think will improve)"), nil
	}

	action := req.GetString("action", "")
	if action == "" {
		return errorResult("action required (what will we try)"), nil
	}

	// Store as a decision with experiment metadata
	d := &memdom.Decision{
		Title:     "EXPERIMENT: " + hypothesis,
		Context:   req.GetString("context", ""),
		Decision:  action,
		Rationale: fmt.Sprintf("hypothesis: %s | measure: %s | duration: %s",
			hypothesis, req.GetString("measure", "observe"), req.GetString("duration", "1 sprint")),
		MadeBy: "team",
		MadeAt: time.Now(),
		Tags:   "experiment," + req.GetString("sprint_name", ""),
	}

	if err := h.Memory.SaveDecision(ctx, d); err != nil {
		return errorResult("Failed: " + err.Error()), nil
	}

	return textResult(fmt.Sprintf("Experiment recorded:\nHypothesis: %s\nAction: %s\nMeasure: %s\nDuration: %s",
		hypothesis, action, req.GetString("measure", "observe"), req.GetString("duration", "1 sprint"))), nil
}

// ReviewExperiments shows active experiments and checks outcomes.
func (h *Handlers) ReviewExperiments(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	decisions, err := h.Memory.SearchDecisions(ctx, "experiment")
	if err != nil {
		return errorResult("Failed: " + err.Error()), nil
	}

	if len(decisions) == 0 {
		return textResult("No experiments recorded. Use pm_experiment to start one from retro."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Improvement Experiments (%d):\n\n", len(decisions)))
	for _, d := range decisions {
		outcome := "PENDING"
		if d.Outcome != "" {
			outcome = d.Outcome
		}
		sb.WriteString(fmt.Sprintf("#%d [%s] %s\n", d.ID, outcome, strings.TrimPrefix(d.Title, "EXPERIMENT: ")))
		sb.WriteString(fmt.Sprintf("   Action: %s\n", d.Decision))
		sb.WriteString(fmt.Sprintf("   Details: %s\n\n", d.Rationale))
	}

	return textResult(sb.String()), nil
}

// SprintPlanningSummary generates a complete sprint planning prep package.
func (h *Handlers) SprintPlanningSummary(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	var sb strings.Builder
	sb.WriteString("=== SPRINT PLANNING PREP ===\n\n")

	// 1. Previous sprint outcome
	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 2)
	if len(snaps) > 0 {
		s := snaps[0]
		sb.WriteString(fmt.Sprintf("LAST SPRINT: %s\n", s.SprintName))
		sb.WriteString(fmt.Sprintf("  Velocity: %d | Completion: %.0f%% | Carryover: %d | Blocked: %d\n\n", s.Velocity, s.CompletionRate, s.Carryover, s.Blocked))
	}

	// 2. Capacity
	if len(snaps) >= 3 {
		var total int
		for _, s := range snaps {
			total += s.Velocity
		}
		avg := total / len(snaps)
		sb.WriteString(fmt.Sprintf("CAPACITY (avg %d sprints): %d points/sprint\n", len(snaps), avg))
		sb.WriteString(fmt.Sprintf("  Conservative (80%%): %d | Stretch: %d\n\n", int(float64(avg)*0.8), avg))
	}

	// 3. Carryover (unfinished from last sprint)
	sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
	if len(sprints) > 0 {
		issues, _ := h.Jira.GetSprintIssues(ctx, sprints[0].ID)
		var carryover int
		for _, issue := range issues {
			lower := strings.ToLower(issue.Status)
			if !strings.Contains(lower, "done") && !strings.Contains(lower, "closed") {
				carryover++
			}
		}
		if carryover > 0 {
			sb.WriteString(fmt.Sprintf("CARRYOVER: %d items still in progress\n\n", carryover))
		}
	}

	// 4. Open risks
	risks, _ := h.Memory.GetOpenRisks(ctx)
	if len(risks) > 0 {
		sb.WriteString(fmt.Sprintf("OPEN RISKS: %d\n", len(risks)))
		for _, r := range risks {
			sb.WriteString(fmt.Sprintf("  [%s] %s\n", r.Severity, r.Title))
		}
		sb.WriteString("\n")
	}

	// 5. Pending action items
	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 0 {
		sb.WriteString(fmt.Sprintf("PENDING RETRO ACTIONS: %d\n", len(actions)))
		for _, a := range actions {
			sb.WriteString(fmt.Sprintf("  - %s (%s)\n", a.Description, a.Owner))
		}
		sb.WriteString("\n")
	}

	// 6. Unresolved dependencies
	deps, _ := h.Memory.GetOpenDependencies(ctx)
	if len(deps) > 0 {
		sb.WriteString(fmt.Sprintf("OPEN DEPENDENCIES: %d\n", len(deps)))
		for _, d := range deps {
			sb.WriteString(fmt.Sprintf("  %s -> %s (%s)\n", d.FromIssueKey, d.ToIssueKey, d.DependencyType))
		}
		sb.WriteString("\n")
	}

	// 7. Active experiments
	experiments, _ := h.Memory.SearchDecisions(ctx, "experiment")
	activeExperiments := 0
	for _, e := range experiments {
		if e.Outcome == "" {
			activeExperiments++
		}
	}
	if activeExperiments > 0 {
		sb.WriteString(fmt.Sprintf("ACTIVE EXPERIMENTS: %d (review results before planning)\n\n", activeExperiments))
	}

	// 8. Working agreements reminder
	agreements, _ := h.Memory.SearchDecisions(ctx, "agreement")
	if len(agreements) > 0 {
		sb.WriteString(fmt.Sprintf("WORKING AGREEMENTS: %d active (review if still valid)\n\n", len(agreements)))
	}

	sb.WriteString("CHECKLIST:\n")
	sb.WriteString("  [ ] Review last sprint outcome\n")
	sb.WriteString("  [ ] Confirm team capacity (leaves, meetings)\n")
	sb.WriteString("  [ ] Address carryover items\n")
	sb.WriteString("  [ ] Set sprint goal\n")
	sb.WriteString("  [ ] Check dependencies\n")
	sb.WriteString("  [ ] Review experiments\n")
	sb.WriteString("  [ ] Confidence vote\n")

	return textResult(sb.String()), nil
}
