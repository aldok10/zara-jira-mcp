package tools

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

var feedbackTableOnce sync.Once

func (h *Handlers) initFeedbackTables() {
	feedbackTableOnce.Do(func() {
		if h.Memory == nil || h.Memory.DB() == nil {
			return
		}
		db := h.Memory.DB()
		_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS feedback_log (
			id INTEGER PRIMARY KEY,
			person TEXT NOT NULL,
			topic TEXT NOT NULL,
			type TEXT DEFAULT 'constructive',
			follow_up_at DATETIME,
			status TEXT DEFAULT 'open',
			outcome TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
		_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS context_note (
			id INTEGER PRIMARY KEY,
			subject TEXT NOT NULL,
			note TEXT NOT NULL,
			sentiment TEXT DEFAULT 'neutral',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
	})
}

func (h *Handlers) CadenceCheck(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	var sb strings.Builder
	sb.WriteString("Communication Cadence Check:\n\n")

	retros, _ := h.Memory.GetRetrospectives(ctx, 3)
	if len(retros) == 0 {
		sb.WriteString("  [OVERDUE] No retrospectives recorded\n")
	} else if time.Since(retros[0].Date).Hours()/24 > 21 {
		sb.WriteString(fmt.Sprintf("  [OVERDUE] Last retro: %d days ago\n", int(time.Since(retros[0].Date).Hours()/24)))
	} else {
		sb.WriteString("  [OK] Retro recent\n")
	}

	pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 5)
	if len(pulses) == 0 {
		sb.WriteString("  [OVERDUE] No pulse recorded\n")
	} else if time.Since(pulses[0].CreatedAt).Hours()/24 > 14 {
		sb.WriteString(fmt.Sprintf("  [OVERDUE] Last pulse: %d days ago\n", int(time.Since(pulses[0].CreatedAt).Hours()/24)))
	} else {
		sb.WriteString("  [OK] Pulse recent\n")
	}

	decisions, _ := h.Memory.GetDecisions(ctx, 5)
	if len(decisions) > 0 && time.Since(decisions[0].MadeAt).Hours()/24 > 14 {
		sb.WriteString(fmt.Sprintf("  [INFO] Last decision: %d days ago\n", int(time.Since(decisions[0].MadeAt).Hours()/24)))
	} else if len(decisions) > 0 {
		sb.WriteString("  [OK] Decisions being recorded\n")
	}

	return textResult(sb.String()), nil
}

func (h *Handlers) CommsNudge(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	var nudges []string

	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	for _, b := range blockers {
		if int(time.Since(b.BlockedSince).Hours()/24) >= 3 {
			nudges = append(nudges, fmt.Sprintf("Escalate: '%s' blocked %d days", b.Description, int(time.Since(b.BlockedSince).Hours()/24)))
		}
	}

	actions, _ := h.Memory.GetPendingActionItems(ctx)
	if len(actions) > 3 {
		nudges = append(nudges, fmt.Sprintf("Follow up: %d pending retro actions", len(actions)))
	}

	pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 10)
	for _, p := range pulses {
		if p.Score <= 2 && time.Since(p.CreatedAt).Hours()/24 < 14 {
			nudges = append(nudges, fmt.Sprintf("Check in: %s scored %d/5 recently", p.Member, p.Score))
			break
		}
	}

	if len(nudges) == 0 {
		return textResult("No urgent nudges. Signals look stable."), nil
	}
	return textResult(fmt.Sprintf("Communication Nudges (%d):\n\n%s", len(nudges), strings.Join(nudges, "\n"))), nil
}

func (h *Handlers) CommsEffectiveness(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	score := 0
	decisions, _ := h.Memory.GetDecisions(ctx, 30)
	score += min(len(decisions)*5, 25)
	history, _ := h.Memory.GetBlockerHistory(ctx, 20)
	resolved := 0
	for _, b := range history {
		if b.ResolvedAt != nil {
			resolved++
		}
	}
	score += min(resolved*5, 25)
	pending, _ := h.Memory.GetPendingActionItems(ctx)
	if len(pending) == 0 {
		score += 25
	} else if len(pending) <= 3 {
		score += 15
	}
	pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 10)
	if len(pulses) >= 5 {
		score += 25
	} else if len(pulses) > 0 {
		score += 10
	}
	return textResult(fmt.Sprintf("Communication Effectiveness: %d/100", score)), nil
}

func (h *Handlers) ConversationPrep(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	convType, err := req.RequireString("type")
	if err != nil {
		return errorResult("type required"), nil
	}
	contextStr, err := req.RequireString("context")
	if err != nil {
		return errorResult("context required"), nil
	}
	person := req.GetString("person", "")

	frameworks := map[string]string{
		"performance":  "SBI: Situation-Behavior-Impact + forward question",
		"conflict":     "STATE: Share facts, Tell story, Ask theirs, Talk tentatively, Encourage",
		"scope":        "BLUF: Bottom line, impact data, options",
		"bad_news":     "Acknowledge difficulty, state facts, next steps",
		"recognition":  "Specific: what they did, impact, why it matters",
	}

	framework, ok := frameworks[convType]
	if !ok {
		return errorResult("type: performance, conflict, scope, bad_news, recognition"), nil
	}

	if h.AI == nil {
		return textResult(fmt.Sprintf("Prep (%s):\nFramework: %s\nContext: %s\nPerson: %s", convType, framework, contextStr, person)), nil
	}

	result, aiErr := h.aiComplete(ctx,
		fmt.Sprintf("Prepare talking points for a %s conversation. Framework: %s. Be direct, under 100 words.", convType, framework),
		fmt.Sprintf("Situation: %s\nPerson: %s", contextStr, person))
	if aiErr != nil {
		return errorResult("AI failed: " + aiErr.Error()), nil
	}
	return textResult(fmt.Sprintf("[%s prep]\n\n%s", convType, result)), nil
}

func (h *Handlers) FeedbackLog(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initFeedbackTables()
	person, _ := req.RequireString("person")
	topic, _ := req.RequireString("topic")
	fbType := req.GetString("type", "constructive")
	days := req.GetInt("follow_up_days", 7)
	followUp := time.Now().AddDate(0, 0, days)
	_, _ = h.Memory.DB().Exec("INSERT INTO feedback_log(person,topic,type,follow_up_at) VALUES(?,?,?,?)", person, topic, fbType, followUp)
	return textResult(fmt.Sprintf("Logged: %s to %s [%s]. Follow-up: %s", topic, person, fbType, followUp.Format("Jan 2"))), nil
}

func (h *Handlers) FeedbackDue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initFeedbackTables()
	rows, err := h.Memory.DB().Query("SELECT id,person,topic,type FROM feedback_log WHERE status='open' AND follow_up_at<=datetime('now')")
	if err != nil {
		return errorResult("query failed"), nil
	}
	defer rows.Close()
	var sb strings.Builder
	count := 0
	for rows.Next() {
		var id int
		var person, topic, fbType string
		_ = rows.Scan(&id, &person, &topic, &fbType)
		sb.WriteString(fmt.Sprintf("  #%d: %s — %s [%s]\n", id, person, topic, fbType))
		count++
	}
	if count == 0 {
		return textResult("No overdue follow-ups."), nil
	}
	return textResult(fmt.Sprintf("Overdue (%d):\n%s", count, sb.String())), nil
}

func (h *Handlers) FeedbackClose(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initFeedbackTables()
	id, _ := req.RequireInt("id")
	outcome, _ := req.RequireString("outcome")
	_, _ = h.Memory.DB().Exec("UPDATE feedback_log SET status='closed',outcome=? WHERE id=?", outcome, id)
	return textResult(fmt.Sprintf("Feedback #%d closed: %s", id, outcome)), nil
}

func (h *Handlers) PMSentiment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	var signals []string
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	for _, b := range blockers {
		signals = append(signals, fmt.Sprintf("[blocker %dd] %s", int(time.Since(b.BlockedSince).Hours()/24), b.Description))
	}
	retros, _ := h.Memory.GetRetrospectives(ctx, 3)
	for _, r := range retros {
		if r.Improvements != "" {
			signals = append(signals, "[improve] "+r.Improvements)
		}
	}
	pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 20)
	for _, p := range pulses {
		if p.Score <= 2 && p.Notes != "" {
			signals = append(signals, fmt.Sprintf("[low-pulse %s] %s", p.Member, p.Notes))
		}
	}
	if msg := req.GetString("message", ""); msg != "" {
		signals = append(signals, "[observation] "+msg)
	}
	if len(signals) == 0 {
		return textResult("Not enough data. Record blockers, retros, pulse first."), nil
	}
	if h.AI == nil {
		return textResult(fmt.Sprintf("Signals: %d blockers, %d retros (AI unavailable for analysis)", len(blockers), len(retros))), nil
	}
	result, err := h.aiComplete(ctx,
		`Empathetic team coach. Assess:
1. MOOD (one word)
2. ROOT CAUSE (1-2 sentences)
3. WHAT TO DO (1 suggestion, gentle)
4. WHAT TO SAY (short warm message PM could send)
Under 150 words.`, strings.Join(signals, "\n"))
	if err != nil {
		return sanitizedError("ai analysis failed for feedback comms", err), nil
	}
	return textResult("Team Sentiment\n\n" + result), nil
}

func (h *Handlers) PMContextNote(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initFeedbackTables()
	subject, _ := req.RequireString("subject")
	note, _ := req.RequireString("note")
	sentiment := req.GetString("sentiment", "neutral")
	_, _ = h.Memory.DB().Exec("INSERT INTO context_note(subject,note,sentiment) VALUES(?,?,?)", subject, note, sentiment)
	return textResult(fmt.Sprintf("Context noted: %s [%s]", subject, sentiment)), nil
}

func (h *Handlers) PMContextRecall(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initFeedbackTables()
	subject := req.GetString("subject", "")
	var query string
	var args []any
	if subject != "" {
		query = "SELECT subject,note,sentiment,created_at FROM context_note WHERE subject LIKE ? ORDER BY created_at DESC LIMIT 20"
		args = []any{"%" + subject + "%"}
	} else {
		query = "SELECT subject,note,sentiment,created_at FROM context_note ORDER BY created_at DESC LIMIT 20"
	}
	rows, err := h.Memory.DB().Query(query, args...)
	if err != nil {
		return errorResult("query failed"), nil
	}
	defer rows.Close()
	var sb strings.Builder
	count := 0
	for rows.Next() {
		var subj, note, sent, created string
		_ = rows.Scan(&subj, &note, &sent, &created)
		sb.WriteString(fmt.Sprintf("  [%s] %s: %s (%s)\n", sent, subj, note, created))
		count++
	}
	if count == 0 {
		return textResult("No context notes. Use pm_context_note to record."), nil
	}
	return textResult(fmt.Sprintf("Context Notes (%d):\n%s", count, sb.String())), nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
