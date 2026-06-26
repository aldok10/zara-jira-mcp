package tools

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

var safetyTablesOnce sync.Once

func (h *Handlers) initSafetyTables() {
	safetyTablesOnce.Do(func() {
		if h.Memory == nil || h.Memory.DB() == nil {
			return
		}
		_, _ = h.Memory.DB().Exec(`CREATE TABLE IF NOT EXISTS safety_survey (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			member TEXT NOT NULL,
			sprint TEXT NOT NULL,
			q1 INTEGER NOT NULL DEFAULT 3,
			q2 INTEGER NOT NULL DEFAULT 3,
			q3 INTEGER NOT NULL DEFAULT 3,
			q4 INTEGER NOT NULL DEFAULT 3,
			q5 INTEGER NOT NULL DEFAULT 3,
			q6 INTEGER NOT NULL DEFAULT 3,
			q7 INTEGER NOT NULL DEFAULT 3,
			notes TEXT DEFAULT '',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
	})
}

// safetyQuestions is the 7-item Edmondson Psychological Safety Scale.
// 1=Strongly Disagree, 5=Strongly Agree.
// q1, q3, q5 are reverse-scored (6 - raw) before averaging.
var safetyQuestions = []struct {
	ID      int
	Text    string
	Reverse bool
}{
	{1, "If I make a mistake on this team, it is held against me", true},
	{2, "Members of this team are able to bring up problems and tough issues", false},
	{3, "People on this team sometimes reject others for being different", true},
	{4, "It is safe to take a risk on this team", false},
	{5, "It is difficult to ask other members of this team for help", true},
	{6, "No one on this team would deliberately undermine my efforts", false},
	{7, "My unique skills and talents are valued and utilized", false},
}

func scoreLabel(avg float64) string {
	switch {
	case avg >= 4.5:
		return "excellent"
	case avg >= 3.5:
		return "good"
	case avg >= 2.5:
		return "moderate"
	case avg >= 1.5:
		return "low"
	default:
		return "critical"
	}
}

func formatScore(s int) string {
	switch s {
	case 1:
		return "Strongly Disagree"
	case 2:
		return "Disagree"
	case 3:
		return "Neutral"
	case 4:
		return "Agree"
	case 5:
		return "Strongly Agree"
	default:
		return fmt.Sprintf("%d", s)
	}
}

// adjustScore applies reverse scoring for q1, q3, q5.
func adjustScore(raw int, reverse bool) int {
	if raw < 1 || raw > 5 {
		raw = 3
	}
	if reverse {
		return 6 - raw
	}
	return raw
}

// PMSafetySurvey records a 7-item psychological safety survey for a member.
func (h *Handlers) PMSafetySurvey(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initSafetyTables()

	member, _ := req.RequireString("member")
	sprint, _ := req.RequireString("sprint")

	scores := make([]int, 7)
	for i := 0; i < 7; i++ {
		qName := fmt.Sprintf("q%d", i+1)
		s := req.GetInt(qName, 3)
		if s < 1 || s > 5 {
			return errorResult(fmt.Sprintf("%s must be 1-5, got %d", qName, s)), nil
		}
		scores[i] = s
	}
	notes := req.GetString("notes", "")

	_, _ = h.Memory.DB().Exec(
		"INSERT INTO safety_survey(member,sprint,q1,q2,q3,q4,q5,q6,q7,notes) VALUES(?,?,?,?,?,?,?,?,?,?)",
		member, sprint, scores[0], scores[1], scores[2], scores[3], scores[4], scores[5], scores[6], notes)

	var total float64
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Safety Survey — %s (sprint: %s)\n\n", member, sprint))
	for i, q := range safetyQuestions {
		raw := scores[i]
		adjusted := adjustScore(raw, q.Reverse)
		total += float64(adjusted)
		arrow := ""
		if q.Reverse {
			arrow = fmt.Sprintf(" [raw %d -> adj %d]", raw, adjusted)
		}
		sb.WriteString(fmt.Sprintf("  q%d: %s%s\n", q.ID, formatScore(raw), arrow))
	}

	avg := total / 7.0
	sb.WriteString(fmt.Sprintf("\n  Average (adjusted): %.1f/5 - %s\n", avg, scoreLabel(avg)))
	return textResult(sb.String()), nil
}

// PMSafetyTrend shows safety score trends across sprints.
func (h *Handlers) PMSafetyTrend(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initSafetyTables()

	member := req.GetString("member", "")
	query := "SELECT sprint, member, q1,q2,q3,q4,q5,q6,q7 FROM safety_survey"
	args := []any{}
	if member != "" {
		query += " WHERE member=?"
		args = append(args, member)
	}
	query += " ORDER BY sprint"

	rows, err := h.Memory.DB().Query(query, args...)
	if err != nil {
		return errorResult("query failed"), nil
	}
	defer rows.Close()

	// Aggregate in Go: map[sprint] -> {count, sumAdjusted}
	type sprintAgg struct {
		count   int
		sum     float64
		members map[string]bool
	}
	agg := map[string]*sprintAgg{}
	order := []string{}

	for rows.Next() {
		var sprint, member string
		var q [7]int
		_ = rows.Scan(&sprint, &member, &q[0], &q[1], &q[2], &q[3], &q[4], &q[5], &q[6])
		if _, ok := agg[sprint]; !ok {
			order = append(order, sprint)
			agg[sprint] = &sprintAgg{members: map[string]bool{}}
		}
		a := agg[sprint]
		a.count++
		a.members[member] = true
		for i := 0; i < 7; i++ {
			a.sum += float64(adjustScore(q[i], safetyQuestions[i].Reverse))
		}
	}

	var sb strings.Builder
	sb.WriteString("Safety Trend\n\n")
	prev := -1.0
	for _, sprint := range order {
		a := agg[sprint]
		avg := a.sum / (7.0 * float64(a.count))
		dir := "->"
		if prev >= 0 {
			if avg > prev+0.2 {
				dir = "+"
			} else if avg < prev-0.2 {
				dir = "-"
			}
		}
		sb.WriteString(fmt.Sprintf("  %s %s: %.1f (n=%d, members=%d) %s\n",
			sprint, dir, avg, a.count, len(a.members), scoreLabel(avg)))
		prev = avg
	}
	if len(order) == 0 {
		return textResult("No survey data yet. Use pm_safety_survey to record."), nil
	}
	return textResult(sb.String()), nil
}

// PMTeamAristotle runs a Google Project Aristotle 5-pillar assessment using AI.
func (h *Handlers) PMTeamAristotle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}

	boardID, _ := req.RequireInt("board_id")

	var signals []string

	// Blockers
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	if len(blockers) > 0 {
		signals = append(signals, fmt.Sprintf("Active blockers (%d):", len(blockers)))
		for _, b := range blockers {
			age := int(time.Since(b.BlockedSince).Hours() / 24)
			signals = append(signals, fmt.Sprintf("  - [%dd] %s", age, b.Description))
		}
	}

	// Retros
	retros, _ := h.Memory.GetRetrospectives(ctx, 5)
	if len(retros) > 0 {
		emptyImprove := 0
		for _, r := range retros {
			if strings.TrimSpace(r.Improvements) == "" {
				emptyImprove++
			}
		}
		signals = append(signals, fmt.Sprintf("Retros: %d in range, %d with empty improvements", len(retros), emptyImprove))
	}

	// Pulses
	pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 50)
	if len(pulses) > 0 {
		lowCount := 0
		totalScore := 0
		for _, p := range pulses {
			totalScore += p.Score
			if p.Score <= 2 {
				lowCount++
			}
		}
		avgPulse := float64(totalScore) / float64(len(pulses))
		signals = append(signals, fmt.Sprintf("Pulse: avg %.1f/5, %d low scores", avgPulse, lowCount))
	}

	// Safety surveys
	h.initSafetyTables()
	var surveyCount int
	if row := h.Memory.DB().QueryRow("SELECT COUNT(*) FROM safety_survey"); row != nil {
		_ = row.Scan(&surveyCount)
	}
	signals = append(signals, fmt.Sprintf("Safety surveys recorded: %d", surveyCount))

	// Decisions
	decisions, _ := h.Memory.GetDecisions(ctx, 20)
	signals = append(signals, fmt.Sprintf("Decisions recorded: %d", len(decisions)))

	// Dependencies
	deps, _ := h.Memory.GetOpenDependencies(ctx)
	signals = append(signals, fmt.Sprintf("Open dependencies: %d", len(deps)))

	// Goals
	goals, _ := h.Memory.GetGoalHistory(ctx, boardID, 10)
	signals = append(signals, fmt.Sprintf("Goals tracked: %d", len(goals)))

	if len(signals) < 3 {
		return textResult("Not enough data for Aristotle assessment. Record blockers, retros, pulses, and safety surveys first."), nil
	}

	if h.AI == nil {
		return textResult("AI unavailable. Install an AI provider or review signals manually:\n" + strings.Join(signals, "\n")), nil
	}

	result, err := h.aiComplete(ctx,
		`You are an empathetic team coach trained on Google's Project Aristotle research (2012-2017).
Assess this team across the 5 pillars that predict high-performing teams:

1. PSYCHOLOGICAL SAFETY — Can members take risks, speak up, admit mistakes?
2. DEPENDABILITY — Do members deliver quality work on time?
3. STRUCTURE & CLARITY — Are goals, roles, and plans clear?
4. MEANING — Does the work matter personally to members?
5. IMPACT — Does the team believe their work creates change?

For each pillar:
- Score 1-5 (with evidence from the data)
- 1-2 sentence assessment
- 1 actionable recommendation

End with a one-sentence overall verdict.
Use the team signals below as evidence. Be honest but constructive. Under 300 words.`,
		strings.Join(signals, "\n"))
	if err != nil {
		return errorResult("AI failed: "+err.Error()), nil
	}
	return textResult("Project Aristotle Assessment\n\n" + result), nil
}
