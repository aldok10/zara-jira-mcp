package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// CommunicateDecision formats a decision announcement using the DACI framework.
func (h *Handlers) CommunicateDecision(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	decision, err := req.RequireString("decision")
	if err != nil {
		return errorResult("decision required (what was decided)"), nil
	}
	contextStr := req.GetString("context", "")
	driver := req.GetString("driver", "")
	approver := req.GetString("approver", "")
	contributors := req.GetString("contributors", "")
	informed := req.GetString("informed", "")

	systemPrompt := `Format this decision announcement using the DACI framework.
Structure:
DECISION: [one clear sentence]
CONTEXT: [why this decision was needed - 1-2 sentences]
ROLES:
  Driver: [who drove it]
  Approver: [who approved]
  Contributors: [who was consulted]
  Informed: [who needs to know]
RATIONALE: [key reasons, max 3 bullets]
IMPACT: [what changes for the team]
NEXT STEPS: [immediate actions]

Keep it under 150 words. Professional but human.`

	data := fmt.Sprintf("Decision: %s\nContext: %s\nDriver: %s\nApprover: %s\nContributors: %s\nInformed: %s",
		decision, contextStr, driver, approver, contributors, informed)

	result, err := h.aiComplete(ctx, systemPrompt, data)
	if err != nil {
		return textResult("DECISION: " + decision + "\nDriver: " + driver + "\nApprover: " + approver), nil
	}
	return textResult(result), nil
}

// EscalateWithSCQA formats an escalation using SCQA (Situation, Complication, Question, Answer).
func (h *Handlers) EscalateWithSCQA(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	situation, err := req.RequireString("situation")
	if err != nil {
		return errorResult("situation required"), nil
	}
	complication, err := req.RequireString("complication")
	if err != nil {
		return errorResult("complication required (what changed/went wrong)"), nil
	}

	question := req.GetString("question", "")
	answer := req.GetString("answer", "")

	var sb strings.Builder
	sb.WriteString("ESCALATION (SCQA Framework)\n\n")
	sb.WriteString(fmt.Sprintf("SITUATION: %s\n\n", situation))
	sb.WriteString(fmt.Sprintf("COMPLICATION: %s\n\n", complication))
	if question != "" {
		sb.WriteString(fmt.Sprintf("QUESTION: %s\n\n", question))
	} else {
		sb.WriteString("QUESTION: What should we do?\n\n")
	}
	if answer != "" {
		sb.WriteString(fmt.Sprintf("RECOMMENDED ANSWER: %s\n", answer))
	} else {
		sb.WriteString("RECOMMENDED ANSWER: [Awaiting management input]\n")
	}
	sb.WriteString(fmt.Sprintf("\nEscalated: %s\n", time.Now().Format("2006-01-02 15:04")))

	return textResult(sb.String()), nil
}

// AdaptMessage rewrites a message for a specific audience using appropriate frameworks.
func (h *Handlers) AdaptMessage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	message, err := req.RequireString("message")
	if err != nil {
		return errorResult("message required (the raw information to communicate)"), nil
	}
	audience, err := req.RequireString("audience")
	if err != nil {
		return errorResult("audience required (executive, po, team, engineering, stakeholder)"), nil
	}

	prompts := map[string]string{
		"executive":   `Rewrite for a VP/CTO with 30 seconds. Apply the Minto Pyramid: lead with the conclusion, then 2-3 supporting points. NO technical jargon, NO story points, NO methodology terms. Business outcomes only. Under 80 words.`,
		"po":          `Rewrite for a Product Owner. Focus on: value delivered, scope decisions needed, user impact. Include specific feature names. Under 100 words.`,
		"team":        `Rewrite for the development team. Be direct and specific. Include: what's changing, why, what they need to do differently. Technical terms OK. Under 100 words.`,
		"engineering": `Rewrite for an engineering audience. Include: technical context, architecture impact, what to watch for. Code-level specifics welcome. Under 120 words.`,
		"stakeholder": `Rewrite for a non-technical business stakeholder. Focus on: timeline impact, budget implications, user-facing changes. No implementation details. Under 80 words.`,
	}

	prompt, ok := prompts[strings.ToLower(audience)]
	if !ok {
		return errorResult("audience must be one of: executive, po, team, engineering, stakeholder"), nil
	}

	result, err := h.aiComplete(ctx, prompt, message)
	if err != nil {
		return textResult(message), nil
	}
	return textResult(fmt.Sprintf("[Adapted for: %s]\n\n%s", audience, result)), nil
}

// GiveFeedback structures feedback using the SBI model (Situation-Behavior-Impact).
func (h *Handlers) GiveFeedback(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	situation, err := req.RequireString("situation")
	if err != nil {
		return errorResult("situation required (when and where)"), nil
	}
	behavior, err := req.RequireString("behavior")
	if err != nil {
		return errorResult("behavior required (what they specifically did)"), nil
	}
	impact, err := req.RequireString("impact")
	if err != nil {
		return errorResult("impact required (effect on team/project/outcome)"), nil
	}

	feedbackType := req.GetString("type", "constructive")

	systemPrompt := fmt.Sprintf(`Format this as %s feedback using the SBI model (Situation-Behavior-Impact).
Rules:
- Be specific and objective (no "always" / "never")
- Describe behavior, not personality
- Focus on impact, not intent
- End with a forward-looking question or suggestion
- Keep warm but direct (Radical Candor: care personally + challenge directly)
- Under 80 words`, feedbackType)

	data := fmt.Sprintf("Situation: %s\nBehavior: %s\nImpact: %s", situation, behavior, impact)

	result, err := h.aiComplete(ctx, systemPrompt, data)
	if err != nil {
		return textResult(fmt.Sprintf("FEEDBACK (SBI)\n\nSituation: %s\nBehavior: %s\nImpact: %s", situation, behavior, impact)), nil
	}
	return textResult(result), nil
}

// WriteUpdate generates a structured status update pulling live Jira data and memory.
func (h *Handlers) WriteUpdate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID := req.GetInt("board_id", 0)
	cadence := req.GetString("cadence", "weekly")
	audience := req.GetString("audience", "team")

	var data strings.Builder
	data.WriteString(fmt.Sprintf("Cadence: %s | Audience: %s\n\n", cadence, audience))

	if boardID > 0 {
		sprints, _ := h.Jira.GetActiveSprints(ctx, boardID)
		if len(sprints) > 0 {
			sprint := sprints[0]
			issues, _ := h.Jira.GetSprintIssues(ctx, sprint.ID)
			var done, total, blocked int
			total = len(issues)
			for _, i := range issues {
				lower := strings.ToLower(i.Status)
				if strings.Contains(lower, "done") || strings.Contains(lower, "closed") {
					done++
				}
				if strings.Contains(lower, "block") {
					blocked++
				}
			}
			pct := float64(0)
			if total > 0 {
				pct = float64(done) / float64(total) * 100
			}
			data.WriteString(fmt.Sprintf("Sprint: %s (Goal: %s)\n", sprint.Name, sprint.Goal))
			data.WriteString(fmt.Sprintf("Progress: %d/%d (%.0f%%) | Blocked: %d\n\n", done, total, pct, blocked))
		}
	}

	risks, _ := h.Memory.GetOpenRisks(ctx)
	if len(risks) > 0 {
		data.WriteString(fmt.Sprintf("Risks: %d open\n", len(risks)))
		for _, r := range risks {
			if r.Severity == "critical" || r.Severity == "high" {
				data.WriteString(fmt.Sprintf("  [%s] %s\n", r.Severity, r.Title))
			}
		}
		data.WriteString("\n")
	}

	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) > 0 {
		data.WriteString(fmt.Sprintf("Blockers: %d active\n", len(blockers)))
	}

	prompts := map[string]string{
		"team":        `Write a team status update. Structure: What shipped > What's in progress > Blockers > Next focus. Casual, direct, under 150 words.`,
		"management":  `Write a management status update using Minto Pyramid: Lead with the conclusion (on track/at risk), then key supporting points, then details only if critical. No jargon. Under 120 words.`,
		"stakeholder": `Write a stakeholder update. Focus on: business progress, timeline confidence, decisions needed. No technical details. Under 100 words.`,
	}

	prompt := prompts[audience]
	if prompt == "" {
		prompt = prompts["team"]
	}

	result, err := h.aiComplete(ctx, prompt, data.String())
	if err != nil {
		return textResult(data.String()), nil
	}
	return textResult(result), nil
}

// CommunicationPlan generates a communication plan for a decision or change.
func (h *Handlers) CommunicationPlan(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topic, err := req.RequireString("topic")
	if err != nil {
		return errorResult("topic required (what needs to be communicated)"), nil
	}
	stakeholders := req.GetString("stakeholders", "")
	urgency := req.GetString("urgency", "normal")

	systemPrompt := `Create a communication plan for this topic. Structure:

KEY MESSAGE: [one sentence everyone should walk away with]

AUDIENCE MAP:
| Who | What They Need | When | Channel | Format |
|-----|---------------|------|---------|--------|

SEQUENCE:
1. First inform: [who, when, how]
2. Then: [broader audience]
3. Finally: [general team/org]

TALKING POINTS: [3-4 bullets for verbal communication]

RISKS: [what could go wrong if poorly communicated]

Keep practical. Under 200 words.`

	inputData := fmt.Sprintf("Topic: %s\nStakeholders: %s\nUrgency: %s", topic, stakeholders, urgency)

	result, err := h.aiComplete(ctx, systemPrompt, inputData)
	if err != nil {
		return sanitizedError("AI failed", err), nil
	}
	return textResult(result), nil
}
