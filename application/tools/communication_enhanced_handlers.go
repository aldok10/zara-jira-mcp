package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	memdom "github.com/aldok10/zara-jira-mcp/domain/memory"
	"github.com/mark3labs/mcp-go/mcp"
)

// PMCommunicate generates Minto Pyramid-structured messages for different audiences.
// If a 'message' param is provided, it rewrites it for the audience. Otherwise generates from topic + data.
func (h *Handlers) PMCommunicate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	topic, err := req.RequireString("topic")
	if err != nil {
		return errorResult("topic required"), nil
	}
	audience, err := req.RequireString("audience")
	if err != nil {
		return errorResult("audience required (exec/team/po/stakeholder)"), nil
	}
	boardID := req.GetInt("board_id", 0)
	existingMessage := req.GetString("message", "")

	systemPrompt := `You write structured PM updates using the Minto Pyramid Principle.
Rules:
1. Start with the conclusion/recommendation (1 sentence)
2. Then 2-3 key supporting arguments
3. Then supporting data only if needed
4. Adapt language to audience:
   - exec: no jargon, business impact, under 80 words
   - team: direct, actionable, technical OK, under 120 words
   - po: value/scope focus, feature names, under 100 words
   - stakeholder: timeline + budget + user impact, under 80 words`

	var userData strings.Builder
	userData.WriteString(fmt.Sprintf("Topic: %s\nAudience: %s\n", topic, audience))
	if existingMessage != "" {
		userData.WriteString(fmt.Sprintf("\nExisting message to rewrite:\n%s\n", existingMessage))
	}

	if boardID > 0 && h.Memory != nil {
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 1)
		if len(snaps) > 0 {
			s := snaps[0]
			userData.WriteString(fmt.Sprintf("Sprint context: %s | Done: %d | In Progress: %d | Blocked: %d\n",
				s.SprintName, s.Done, s.InProgress, s.Blocked))
		}
	}

	if h.AI == nil {
		return errorResult("AI provider not configured"), nil
	}

	result, err := h.aiComplete(ctx, systemPrompt, userData.String())
	if err != nil {
		return sanitizedError("ai analysis failed for enhanced comms", err), nil
	}
	return textResult(fmt.Sprintf("[Minto Pyramid | %s]\n\n%s", audience, result)), nil
}

// PMFeedbackPrep generates SBI-formatted feedback for a team member.
func (h *Handlers) PMFeedbackPrep(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	member, err := req.RequireString("member")
	if err != nil {
		return errorResult("member required"), nil
	}
	observation, err := req.RequireString("observation")
	if err != nil {
		return errorResult("observation required (what you observed)"), nil
	}
	fbType := req.GetString("type", "constructive")

	systemPrompt := fmt.Sprintf(`Generate %s feedback using the SBI model:
- Situation: When and where (specific time/meeting/event)
- Behavior: What the person specifically did (observable, not interpreted)
- Impact: Effect on team, project, or outcomes

Rules:
- Be specific and objective
- Describe behavior not personality
- No "always"/"never" language
- End with a forward-looking question
- Under 100 words`, fbType)

	var userData strings.Builder
	userData.WriteString(fmt.Sprintf("Team member: %s\nObservation: %s\n", member, observation))

	if h.Memory != nil {
		metrics, _ := h.Memory.GetTeamMetrics(ctx, member, 3)
		if len(metrics) > 0 {
			m := metrics[0]
			userData.WriteString(fmt.Sprintf("Recent data: %d assigned, %d done, %d carried over (sprint: %s)\n",
				m.IssuesAssigned, m.IssuesDone, m.CarryoverCount, m.SprintName))
		}
	}

	if h.AI == nil {
		return errorResult("AI provider not configured"), nil
	}

	result, err := h.aiComplete(ctx, systemPrompt, userData.String())
	if err != nil {
		return sanitizedError("ai analysis failed for enhanced comms", err), nil
	}
	return textResult(fmt.Sprintf("[SBI Feedback - %s for %s]\n\n%s", fbType, member, result)), nil
}

// PMEscalationDraft generates a pyramid-structured escalation message.
func (h *Handlers) PMEscalationDraft(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	issue, err := req.RequireString("issue")
	if err != nil {
		return errorResult("issue required"), nil
	}
	severity := req.GetString("severity", "high")
	deadline := req.GetString("deadline", "")

	if h.AI == nil {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("ESCALATION [%s]\n\n", strings.ToUpper(severity)))
		sb.WriteString(fmt.Sprintf("ASK: [action needed for] %s\n\n", issue))
		sb.WriteString("CONTEXT: [fill in]\n\nIMPACT IF NOT RESOLVED: [fill in]\n\nPROPOSED NEXT STEP: [fill in]\n")
		if deadline != "" {
			sb.WriteString(fmt.Sprintf("DEADLINE: %s\n", deadline))
		}
		return textResult(sb.String()), nil
	}

	systemPrompt := `Write an escalation message. Structure exactly:
1. ASK: One line - what you need (decision/resource/unblock)
2. CONTEXT: 2 sentences max - why this matters now
3. IMPACT: What happens if not resolved
4. PROPOSED NEXT STEP: Your recommendation
5. DEADLINE: When this needs resolution

Be direct. No filler. Under 100 words total.`

	userData := fmt.Sprintf("Issue: %s\nSeverity: %s\nDeadline: %s", issue, severity, deadline)

	result, err := h.aiComplete(ctx, systemPrompt, userData)
	if err != nil {
		return sanitizedError("ai analysis failed for enhanced comms", err), nil
	}
	return textResult(fmt.Sprintf("[Escalation - %s]\n\n%s", severity, result)), nil
}

// PMDecisionRecordEnhanced formats a decision as an Architecture Decision Record.
func (h *Handlers) PMDecisionRecordEnhanced(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	title, err := req.RequireString("title")
	if err != nil {
		return errorResult("title required"), nil
	}
	decision, err := req.RequireString("decision")
	if err != nil {
		return errorResult("decision required"), nil
	}
	contextStr := req.GetString("context", "")
	alternatives := req.GetString("alternatives", "")
	consequences := req.GetString("consequences", "")

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# ADR: %s\n\n", title))
	sb.WriteString("## Status\nAccepted\n\n")
	if contextStr != "" {
		sb.WriteString(fmt.Sprintf("## Context\n%s\n\n", contextStr))
	}
	sb.WriteString(fmt.Sprintf("## Decision\n%s\n\n", decision))
	if alternatives != "" {
		sb.WriteString(fmt.Sprintf("## Alternatives Considered\n%s\n\n", alternatives))
	}
	if consequences != "" {
		sb.WriteString(fmt.Sprintf("## Consequences\n%s\n\n", consequences))
	}

	// Store to memory if available
	if h.Memory != nil {
		d := &memdom.Decision{
			Title:    title,
			Decision: decision,
			Context:  contextStr,
			MadeAt:   time.Now(),
			Tags:     "adr",
		}
		_ = h.Memory.SaveDecision(ctx, d)
	}

	return textResult(sb.String()), nil
}
