package tools

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// EmpathySystemPrompt is the personality DNA for all AI outputs in this system.
const EmpathySystemPrompt = `You are a smart, caring PM companion — not a reporting tool. You communicate like a trusted senior colleague who genuinely cares about the team.

Communication rules:
- Mix Bahasa Indonesia and English naturally (technical terms in English, relational framing in Bahasa)
- Be concise: 3-5 sentences for summaries, expand only when asked
- Use "Noticed..." / "Pattern shows..." not "Alert:" or "Warning:"
- Frame insights as observations + options, never commands
- Acknowledge difficulty before suggesting solutions
- Connect to historical patterns when relevant
- Show confidence level when data is thin
- Celebrate small wins naturally before surfacing issues
- When uncertain, say so honestly
- Never use emojis
- Keep language simple — avoid jargon unless the user uses it first
- Give 2-3 options, not mandates
- Use collaborative framing: "Worth exploring?" / "One approach..." not "You should..."
`

// EmpathyContext holds contextual signals for AI output calibration.
type EmpathyContext struct {
	TeamMood       string
	SprintPhase    string
	HistoryNote    string
	DataConfidence string
}

// aiComplete wraps h.AI.Complete with nil safety and empathetic system prompt.
// Prepends EmpathySystemPrompt to the caller's system for consistent warm tone.
func (h *Handlers) aiComplete(ctx context.Context, system, user string) (string, error) {
	if h.AI == nil {
		return "", fmt.Errorf("AI provider not configured")
	}
	fullSystem := EmpathySystemPrompt + "\n" + system
	return h.AI.Complete(ctx, fullSystem, user)
}

// aiCompleteStructured is like aiComplete but without the EmpathySystemPrompt,
// for calls that need raw structured output (e.g. JSON extraction).
func (h *Handlers) aiCompleteStructured(ctx context.Context, system, user string) (string, error) {
	if h.AI == nil {
		return "", fmt.Errorf("AI provider not configured")
	}
	return h.AI.Complete(ctx, system, user)
}

// enrichWithContext gathers team signals for empathy-aware AI responses.
func (h *Handlers) enrichWithContext(ctx context.Context, boardID int) EmpathyContext {
	ec := EmpathyContext{
		TeamMood:       "neutral",
		SprintPhase:    "mid",
		DataConfidence: "medium",
	}

	if h.Memory == nil {
		ec.DataConfidence = "low"
		return ec
	}

	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) >= 3 {
		ec.TeamMood = "stressed"
	}

	scores, _ := h.Memory.GetHealthScores(ctx, boardID, 1)
	if len(scores) > 0 && scores[0].OverallScore < 50 {
		ec.TeamMood = "stressed"
	}

	snap, _ := h.Memory.GetLatestSnapshot(ctx, boardID)
	if snap != nil {
		if snap.Carryover > 5 {
			ec.TeamMood = "frustrated"
		}
		if snap.Velocity > 0 && snap.Carryover == 0 {
			ec.TeamMood = "energized"
		}
	}

	if h.Jira != nil && boardID > 0 {
		sprints, err := h.Jira.GetActiveSprints(ctx, boardID)
		if err == nil && len(sprints) > 0 {
			if endDate, err := time.Parse("2006-01-02T15:04:05.000Z", sprints[0].EndDate); err == nil {
				remaining := int(time.Until(endDate).Hours() / 24)
				if remaining > 7 {
					ec.SprintPhase = "early"
				} else if remaining > 3 {
					ec.SprintPhase = "mid"
				} else {
					ec.SprintPhase = "late"
				}
			}
		}
	}

	snapshots, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
	if len(snapshots) >= 5 {
		ec.DataConfidence = "high"
	} else if len(snapshots) == 0 {
		ec.DataConfidence = "low"
	}

	return ec
}

// buildEmpathyPrompt returns context-aware additions to the system prompt.
func (h *Handlers) buildEmpathyPrompt(ec EmpathyContext) string {
	var parts []string

	switch ec.TeamMood {
	case "stressed":
		parts = append(parts, "Team is under pressure — be brief, acknowledge difficulty, focus on one actionable step.")
	case "frustrated":
		parts = append(parts, "Team may be frustrated (high carryover) — validate feeling, suggest scope adjustment.")
	case "energized":
		parts = append(parts, "Team is in good shape — match energy, stretch slightly.")
	}

	switch ec.SprintPhase {
	case "early":
		parts = append(parts, "Sprint just started — focus on clarity and early risk detection.")
	case "late":
		parts = append(parts, "Sprint ending soon — focus on completion and scope decisions.")
	}

	switch ec.DataConfidence {
	case "low":
		parts = append(parts, "Limited data — phrase as observations, not conclusions.")
	case "high":
		parts = append(parts, "Good data history — reference trends confidently.")
	}

	if ec.HistoryNote != "" {
		parts = append(parts, fmt.Sprintf("Context: %s", ec.HistoryNote))
	}

	if len(parts) == 0 {
		return ""
	}
	return "\nContext:\n" + strings.Join(parts, "\n")
}
