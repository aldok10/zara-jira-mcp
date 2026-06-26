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

func (h *Handlers) initFeedbackTable() {
	feedbackTableOnce.Do(func() {
		if h.Memory == nil || h.Memory.DB() == nil {
			return
		}
		db := h.Memory.DB()
		_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS feedback_log (id INTEGER PRIMARY KEY, person TEXT NOT NULL, type TEXT DEFAULT 'constructive', topic TEXT NOT NULL, follow_up_date TEXT, followed_up INTEGER DEFAULT 0, outcome TEXT, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`)
	})
}

func (h *Handlers) CadenceCheck(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	db := h.Memory.DB()
	if db == nil {
		return errorResult("database not available"), nil
	}
	now := time.Now()
	type item struct{ name, query string; targetD int }
	items := []item{
		{"Escalation", "SELECT MAX(created_at) FROM escalations", 14},
		{"Decision", "SELECT MAX(made_at) FROM decisions", 14},
		{"Sprint Snapshot", "SELECT MAX(created_at) FROM sprint_snapshots", 21},
		{"Retrospective", "SELECT MAX(date) FROM retrospectives", 21},
		{"Risk Scan", "SELECT MAX(identified_at) FROM risks", 7},
	}
	var sb strings.Builder
	sb.WriteString("Communication Cadence Status:\n\n")
	overdue := 0
	for _, it := range items {
		var lastStr *string
		row := db.QueryRow(it.query)
		_ = row.Scan(&lastStr)
		if lastStr == nil || *lastStr == "" {
			sb.WriteString(fmt.Sprintf("  %-18s Never recorded [NO DATA]\n", it.name+":"))
			overdue++
			continue
		}
		last, err := parseDateFlex(*lastStr)
		if err != nil {
			continue
		}
		days := int(now.Sub(last).Hours() / 24)
		status := "OK"
		if days > it.targetD {
			status = fmt.Sprintf("OVERDUE (target: %dd)", it.targetD)
			overdue++
		}
		sb.WriteString(fmt.Sprintf("  %-18s %d days ago [%s]\n", it.name+":", days, status))
	}
	sb.WriteString(fmt.Sprintf("\n%d/%d overdue.", overdue, len(items)))
	if overdue > 0 {
		sb.WriteString(" Address these to maintain trust.")
	}
	return textResult(sb.String()), nil
}

func (h *Handlers) KPITrend(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initOKRTables()
	kpiName := req.GetString("name", "")
	kpiID := req.GetInt("kpi_id", 0)
	limit := req.GetInt("limit", 10)
	db := h.Memory.DB()
	var rows interface{ Next() bool; Scan(...any) error; Close() error }
	var err error
	if kpiID > 0 {
		rows, err = db.Query("SELECT s.value, s.sprint_name, s.created_at, d.name, d.unit, d.target_value FROM kpi_snapshot s JOIN kpi_definition d ON d.id = s.kpi_id WHERE s.kpi_id = ? ORDER BY s.created_at DESC LIMIT ?", kpiID, limit)
	} else if kpiName != "" {
		rows, err = db.Query("SELECT s.value, s.sprint_name, s.created_at, d.name, d.unit, d.target_value FROM kpi_snapshot s JOIN kpi_definition d ON d.id = s.kpi_id WHERE d.name LIKE ? ORDER BY s.created_at DESC LIMIT ?", "%"+kpiName+"%", limit)
	} else {
		return errorResult("provide kpi_id or name"), nil
	}
	if err != nil {
		return errorResult("Query failed: " + err.Error()), nil
	}
	defer rows.Close()
	type point struct{ value, target float64; sprint, date, name, unit string }
	var points []point
	for rows.Next() {
		var p point
		var target *float64
		if err := rows.Scan(&p.value, &p.sprint, &p.date, &p.name, &p.unit, &target); err != nil {
			continue
		}
		if target != nil {
			p.target = *target
		}
		points = append(points, p)
	}
	if len(points) == 0 {
		return textResult("No KPI snapshots. Use pm_kpi_snapshot to record data."), nil
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("KPI Trend: %s (%s)\n", points[0].name, points[0].unit))
	if points[0].target > 0 {
		sb.WriteString(fmt.Sprintf("Target: %.1f%s\n\n", points[0].target, points[0].unit))
	}
	for i := len(points) - 1; i >= 0; i-- {
		p := points[i]
		trend := " "
		if i < len(points)-1 {
			if p.value > points[i+1].value { trend = "^" } else if p.value < points[i+1].value { trend = "v" } else { trend = "=" }
		}
		label := p.sprint
		if label == "" && len(p.date) >= 10 { label = p.date[:10] }
		sb.WriteString(fmt.Sprintf("  %s %s %.1f%s\n", trend, label, p.value, p.unit))
	}
	if len(points) >= 2 {
		first, last := points[len(points)-1].value, points[0].value
		dir := "stable"
		if last > first { dir = "improving" } else if last < first { dir = "declining" }
		sb.WriteString(fmt.Sprintf("\nDirection: %s (%.1f -> %.1f)", dir, first, last))
	}
	return textResult(sb.String()), nil
}

func (h *Handlers) CommsNudge(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	var nudges []string
	now := time.Now()
	blockers, _ := h.Memory.GetActiveBlockers(ctx)
	for _, b := range blockers {
		if days := int(now.Sub(b.BlockedSince).Hours() / 24); days >= 3 {
			nudges = append(nudges, fmt.Sprintf("[ESCALATION] '%s' blocked %d days. Consider escalating.", b.Description, days))
			break
		}
	}
	pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 10)
	if len(pulses) >= 2 && pulses[0].Score < pulses[1].Score && pulses[0].Score <= 3 {
		nudges = append(nudges, fmt.Sprintf("[TEAM] Pulse dropped to %.1f. Might need a sync.", pulses[0].Score))
	}
	pending, _ := h.Memory.GetPendingActionItems(ctx)
	if len(pending) >= 5 {
		nudges = append(nudges, fmt.Sprintf("[FOLLOW-UP] %d action items pending. Still relevant?", len(pending)))
	}
	h.initFeedbackTable()
	if db := h.Memory.DB(); db != nil {
		var dueCount int
		row := db.QueryRow("SELECT COUNT(*) FROM feedback_log WHERE followed_up = 0 AND follow_up_date <= ?", now.Format("2006-01-02"))
		_ = row.Scan(&dueCount)
		if dueCount > 0 {
			nudges = append(nudges, fmt.Sprintf("[FEEDBACK] %d follow-ups overdue.", dueCount))
		}
	}
	decisions, _ := h.Memory.GetDecisions(ctx, 5)
	if len(decisions) > 0 {
		if days := int(now.Sub(decisions[0].MadeAt).Hours() / 24); days > 14 {
			nudges = append(nudges, fmt.Sprintf("[DECISION] No decisions documented in %d days.", days))
		}
	}
	if len(nudges) == 0 {
		return textResult("All clear. No communication actions needed right now."), nil
	}
	return textResult(fmt.Sprintf("Nudges (%d):\n\n%s", len(nudges), strings.Join(nudges, "\n\n"))), nil
}

func (h *Handlers) ConversationPrep(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	convType, err := req.RequireString("type")
	if err != nil {
		return errorResult("type required: performance, conflict, scope_negotiation, bad_news, recognition"), nil
	}
	situation := req.GetString("context", "")
	person := req.GetString("person", "")
	boardID := req.GetInt("board_id", 0)
	var data strings.Builder
	data.WriteString(fmt.Sprintf("Type: %s\nContext: %s\n", convType, situation))
	if person != "" { data.WriteString(fmt.Sprintf("Person: %s\n", person)) }
	if h.Memory != nil && boardID > 0 {
		scores, _ := h.Memory.GetHealthScores(ctx, boardID, 3)
		if len(scores) > 0 { data.WriteString(fmt.Sprintf("Sprint health: %d/100\n", scores[0].OverallScore)) }
	}
	frameworks := map[string]string{"performance": "SBI + Radical Candor", "conflict": "NVC + SCARF", "scope_negotiation": "STATE path (Crucial Conversations)", "bad_news": "Pyramid Principle + SCARF", "recognition": "Specific praise + growth connection"}
	fw := frameworks[convType]
	if fw == "" { fw = frameworks["performance"] }
	if h.AI == nil {
		return textResult(fmt.Sprintf("Conversation Prep (%s)\nFramework: %s\n\nData:\n%s\n\n1. Facts (observable)\n2. Intent (for them, for team)\n3. SCARF risks\n4. Opening line\n5. If defensive: restore safety", convType, fw, data.String())), nil
	}
	prompt := fmt.Sprintf("Prepare a %s conversation using %s.\n\nOutput:\nINTENT: What you want for them and the team\nSCARF RISK: Which domain threatened\nFACTS: Observable facts from data\nOPENING: 2 options (tentative, caring)\nCORE: Message using framework\nIF DEFENSIVE: Restore safety\nCLOSE: Agreement + timeline\n\nUnder 200 words. Direct. Partner tone, never lecture.", convType, fw)
	result, aiErr := h.AI.Complete(ctx, prompt, data.String())
	if aiErr != nil {
		return errorResult("AI failed: " + aiErr.Error()), nil
	}
	return textResult(result), nil
}

func (h *Handlers) FeedbackLog(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil { return errorResult("memory not configured"), nil }
	h.initFeedbackTable()
	person, err := req.RequireString("person")
	if err != nil { return errorResult("person required"), nil }
	topic, err := req.RequireString("topic")
	if err != nil { return errorResult("topic required"), nil }
	fbType := req.GetString("type", "constructive")
	followUp := req.GetString("follow_up_date", "")
	if followUp == "" { followUp = time.Now().AddDate(0, 0, 7).Format("2006-01-02") }
	db := h.Memory.DB()
	_, _ = db.Exec("INSERT INTO feedback_log (person, type, topic, follow_up_date) VALUES (?, ?, ?, ?)", person, fbType, topic, followUp)
	return textResult(fmt.Sprintf("Feedback logged: %s to %s\nTopic: %s\nFollow-up: %s", fbType, person, topic, followUp)), nil
}

func (h *Handlers) FeedbackDue(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil { return errorResult("memory not configured"), nil }
	h.initFeedbackTable()
	db := h.Memory.DB()
	rows, err := db.Query("SELECT id, person, type, topic, follow_up_date FROM feedback_log WHERE followed_up = 0 ORDER BY follow_up_date ASC")
	if err != nil { return errorResult("Query failed"), nil }
	defer rows.Close()
	now := time.Now()
	var sb strings.Builder
	sb.WriteString("Feedback Follow-ups:\n\n")
	count := 0
	for rows.Next() {
		var id int64
		var person, fbType, topic, followUp string
		if err := rows.Scan(&id, &person, &fbType, &topic, &followUp); err != nil { continue }
		count++
		status := "pending"
		if fDate, err := time.Parse("2006-01-02", followUp); err == nil && now.After(fDate) {
			status = fmt.Sprintf("OVERDUE %d days", int(now.Sub(fDate).Hours()/24))
		}
		sb.WriteString(fmt.Sprintf("  #%d [%s] %s -> %s | %s\n", id, status, fbType, person, topic))
	}
	if count == 0 { return textResult("No pending follow-ups."), nil }
	return textResult(sb.String()), nil
}

func (h *Handlers) FeedbackClose(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil { return errorResult("memory not configured"), nil }
	h.initFeedbackTable()
	id, err := req.RequireInt("id")
	if err != nil { return errorResult("id required"), nil }
	outcome := req.GetString("outcome", "acknowledged")
	db := h.Memory.DB()
	_, _ = db.Exec("UPDATE feedback_log SET followed_up = 1, outcome = ? WHERE id = ?", outcome, id)
	return textResult(fmt.Sprintf("Feedback #%d closed. Outcome: %s", id, outcome)), nil
}

func (h *Handlers) CommsEffectiveness(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil { return errorResult("memory not configured"), nil }
	h.initFeedbackTable()
	var sb strings.Builder
	sb.WriteString("Communication Effectiveness\n\n")
	db := h.Memory.DB()
	var total, followed int
	row := db.QueryRow("SELECT COUNT(*), COALESCE(SUM(followed_up),0) FROM feedback_log")
	_ = row.Scan(&total, &followed)
	if total > 0 {
		sb.WriteString(fmt.Sprintf("Feedback Follow-through: %.0f%% (%d/%d)\n", float64(followed)/float64(total)*100, followed, total))
	} else {
		sb.WriteString("Feedback: no data yet\n")
	}
	escalations, _ := h.Memory.GetRecentEscalations(ctx, 20)
	if len(escalations) > 0 {
		acked := 0
		for _, e := range escalations { if e.Acknowledged { acked++ } }
		sb.WriteString(fmt.Sprintf("Escalation Response: %.0f%% (%d/%d)\n", float64(acked)/float64(len(escalations))*100, acked, len(escalations)))
	}
	decisions, _ := h.Memory.GetDecisions(ctx, 100)
	recent := 0
	cutoff := time.Now().AddDate(0, 0, -30)
	for _, d := range decisions { if d.MadeAt.After(cutoff) { recent++ } }
	sb.WriteString(fmt.Sprintf("Decisions (30d): %d\n", recent))
	pending, _ := h.Memory.GetPendingActionItems(ctx)
	var completed int
	row2 := db.QueryRow("SELECT COUNT(*) FROM action_items WHERE status='done'")
	_ = row2.Scan(&completed)
	tot := len(pending) + completed
	if tot > 0 {
		sb.WriteString(fmt.Sprintf("Action Completion: %.0f%% (%d/%d)\n", float64(completed)/float64(tot)*100, completed, tot))
	}
	pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 5)
	if len(pulses) >= 2 {
		dir := "stable"
		if pulses[0].Score > pulses[len(pulses)-1].Score { dir = "improving" } else if pulses[0].Score < pulses[len(pulses)-1].Score { dir = "declining" }
		sb.WriteString(fmt.Sprintf("Team Pulse: %s (%.1f)\n", dir, pulses[0].Score))
	}
	return textResult(sb.String()), nil
}

func parseDateFlex(s string) (time.Time, error) {
	for _, f := range []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", time.RFC3339, "2006-01-02"} {
		if t, err := time.Parse(f, s); err == nil { return t, nil }
	}
	return time.Time{}, fmt.Errorf("cannot parse: %s", s)
}
