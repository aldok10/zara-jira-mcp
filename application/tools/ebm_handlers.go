package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// PMValueCheck helps distinguish outputs from outcomes.
func (h *Handlers) PMValueCheck(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}

	boardID := req.GetInt("board_id", 0)
	var sb strings.Builder

sb.WriteString("Value Check — Output vs Outcome\n\n")

	// Examples of outputs (what we produce) vs outcomes (what we achieve)

sb.WriteString("OUTPUTS (tangible deliverables):\n")

sb.WriteString("  - Code commits, story points, lines of code\n")

sb.WriteString("  - Number of meetings, tickets closed, documentation pages\n")

sb.WriteString("  - Sprint burndown charts, velocity graphs\n\n")


sb.WriteString("OUTCOMES (business impact):\n")

sb.WriteString("  - Customer satisfaction, user adoption rates\n")

sb.WriteString("  - Revenue generated, cost reduced, time saved\n")

sb.WriteString("  - Team morale, reduced technical debt, faster delivery\n\n")

	// Show examples from data
	if boardID > 0 {
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
		if len(snaps) > 0 {
			sb.WriteString("Recent Sprint Examples:\n")
			for _, s := range snaps {
				sb.WriteString(fmt.Sprintf("  Sprint %s: %d points, %d issues, %d blocked\n", s.SprintName, s.Velocity, s.TotalIssues, s.Blocked))
			}
		}
	}

	// AI-powered analysis if available
	if h.AI != nil {
		var signals []string
		pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 5)
		if len(pulses) > 0 {
			signals = append(signals, fmt.Sprintf("Team pulse: %.1f/5", float64(pulses[0].Score)))
		}
		goals, _ := h.Memory.GetActiveGoals(ctx, boardID)
		if len(goals) > 0 {
			signals = append(signals, fmt.Sprintf("Active goals: %d", len(goals)))
		}
		if len(signals) > 0 {
			result, err := h.aiComplete(ctx,
				`You are a business outcomes coach. Given the sprint data signals below, answer:

1. What are 2-3 outputs we're tracking?
2. What are 2-3 outcomes we're actually achieving?
3. Which outputs best predict our outcomes?
4. One recommendation: shift focus from output to outcome.

Keep it practical, not theoretical. Under 150 words.`,
				strings.Join(signals, "\n"))
			if err == nil {
				sb.WriteString("\nAI Analysis:\n" + result)
			} else {
				sb.WriteString("\nAI unavailable for analysis.")
			}
		}
	}

	sb.WriteString("\nNext steps:\n")
	sb.WriteString("  - Define outcome metrics for each KVA\n")
	sb.WriteString("  - Track outcomes weekly, not just outputs\n")
	sb.WriteString("  - Connect sprint work to business impact\n")

	return textResult(sb.String()), nil
}
