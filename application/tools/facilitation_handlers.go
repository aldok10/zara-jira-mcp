package tools

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// PMRetroFormat recommends a retro format based on team context.
func (h *Handlers) PMRetroFormat(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}

	boardID := req.GetInt("board_id", 0)
	var signals []string

	if boardID > 0 {
		// Sprint goal success
		goals, _ := h.Memory.GetGoalHistory(ctx, boardID, 5)
		if len(goals) > 0 {
			hit := 0
			for _, g := range goals {
				if g.Status == "achieved" {
					hit++
				}
			}
			signals = append(signals, fmt.Sprintf("Sprint goal success: %d/%d (%.0f%%)", hit, len(goals), float64(hit)/float64(len(goals))*100))
		}
	}

	// Blockers
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) > 0 {
		oldBlocker := 0
		for _, b := range blockers {
			if time.Since(b.BlockedSince).Hours() > 48 {
				oldBlocker++
			}
		}
		signals = append(signals, fmt.Sprintf("Blockers: %d active, %d older than 48h", len(blockers), oldBlocker))
	}

	// Pulse
	pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 30)
	if len(pulses) > 0 {
		low := 0
		for _, p := range pulses {
			if p.Score <= 2 {
				low++
			}
		}
		signals = append(signals, fmt.Sprintf("Team pulse: %d entries, %d low scores", len(pulses), low))
	}

	// Safety surveys
	h.initSafetyTables()
	var safeCount int
	if r := h.Memory.DB().QueryRow("SELECT COUNT(*) FROM safety_survey"); r != nil {
		_ = r.Scan(&safeCount)
	}
	signals = append(signals, fmt.Sprintf("Safety surveys: %d", safeCount))

	// Action items from last retro
	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 0 {
		signals = append(signals, fmt.Sprintf("Pending retro actions: %d", len(actions)))
	}

	// Retros
	retros, _ := h.Memory.GetRetrospectives(ctx, 3)
	emptyRetros := 0
	for _, r := range retros {
		if strings.TrimSpace(r.Improvements) == "" && strings.TrimSpace(r.WentWell) == "" {
			emptyRetros++
		}
	}
	signals = append(signals, fmt.Sprintf("Recent retros: %d, empty retros: %d", len(retros), emptyRetros))

	if len(signals) < 2 {
		return textResult("Not enough data. Record blockers, pulses, sprint goals, and retros first."), nil
	}

	if h.AI == nil {
		return textResult("Signals available (AI needed for format recommendation):\n" + strings.Join(signals, "\n")), nil
	}

	result, err := h.aiComplete(ctx,
		`You are a seasoned Agile coach. Based on the team signals below, recommend the BEST retro format for their CURRENT SITUATION.

Available formats:
- Start/Stop/Continue — Best for stable teams needing incremental improvement
- Sailboat (Wind/Anchor/Rocks/Island) — Great when there are clear blockers or challenges
- 4Ls (Liked/Learned/Lacked/Longed For) — Best for teams needing reflection depth
- Timeline (plot sprint as story) — Great after intense/eventful sprints
- Mad/Sad/Glad — Best when emotions are running high or low
- WWW (Wrong/What/Why) — Best after a failed sprint or incident
- Lean Coffee — Best when team has many topics to discuss
- Strengths-based (what went well, amplify) — Best when team morale is low

For each:
1. RECOMMENDED format (one only)
2. WHY it fits this team's current state
3. ONE facilitation tip to make it effective

Under 200 words. Be specific, not generic.`,
		strings.Join(signals, "\n"))
	if err != nil {
		return sanitizedError("AI failed", err), nil
	}
	return textResult("Retro Format Recommendation\n\n" + result), nil
}

// PMMeetingAudit assesses whether a meeting could be async.
func (h *Handlers) PMMeetingAudit(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	meetingName, _ := req.RequireString("meeting_name")
	meetingType, _ := req.RequireString("meeting_type")
	duration := req.GetInt("duration_minutes", 30)
	attendees := req.GetInt("attendees", 5)
	frequency := req.GetString("frequency", "weekly")
	agendaItems := req.GetInt("agenda_items", 3)

	// Calculate async score (0-100): higher = more async-friendly
	asyncScore := 0

	// Decision-heavy meetings should be sync
	switch meetingType {
	case "standup":
		asyncScore = 80
	case "status_update":
		asyncScore = 85
	case "planning":
		asyncScore = 20
	case "retro":
		asyncScore = 30
	case "grooming":
		asyncScore = 40
	case "review":
		asyncScore = 50
	case "brainstorming":
		asyncScore = 15
	case "decision":
		asyncScore = 10
	case "1on1":
		asyncScore = 25
	case "allhands":
		asyncScore = 70
	default:
		asyncScore = 50
	}

	// Adjustments
	if len(meetingName) > 0 {
		// "standup" in name pushes it toward async
		lower := strings.ToLower(meetingName)
		if strings.Contains(lower, "standup") || strings.Contains(lower, "daily") {
			asyncScore = (asyncScore + 80) / 2
		}
		if strings.Contains(lower, "decision") || strings.Contains(lower, "approve") {
			asyncScore = (asyncScore + 10) / 2
		}
	}

	// Long meetings with few agenda items = async-friendly
	if duration > 45 {
		asyncScore -= 10
	}
	if attendees > 10 {
		asyncScore += 15 // more people = harder to sync
	}
	if agendaItems <= 1 {
		asyncScore += 20 // single topic = email
	}
	if frequency == "daily" {
		asyncScore += 10
	}

	if asyncScore > 100 {
		asyncScore = 100
	}
	if asyncScore < 0 {
		asyncScore = 0
	}

	var verdict, action string
	switch {
	case asyncScore >= 70:
		verdict = "STRONGLY ASYNC-FRIENDLY"
		action = "Replace with async channel (Slack, email, doc). Schedule sync only if async fails."
	case asyncScore >= 45:
		verdict = "COULD BE PARTIALLY ASYNC"
		action = "Move status/broadcast items to async doc. Use sync time only for discussion/decisions."
	case asyncScore >= 25:
		verdict = "SYNC-RECOMMENDED"
		action = "Keep sync but tighten agenda. Set strict timebox per item."
	default:
		verdict = "MUST BE SYNC"
		action = "This meeting requires real-time collaboration. Protect the timebox and invite only essential people."
	}

	return textResult(fmt.Sprintf(`Meeting Audit: %s

Type: %s | Duration: %dm | Attendees: %d | Frequency: %s

Async Score: %d/100 — %s

Recommendation: %s

Tips:
- Share agenda 24h before
- Document decisions + action items during
- Timebox strictly`, meetingName, meetingType, duration, attendees, frequency, asyncScore, verdict, action)), nil
}
