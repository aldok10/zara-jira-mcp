package tools

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

var empathyTableOnce sync.Once

func (h *Handlers) initEmpathyTables() {
	empathyTableOnce.Do(func() {
		if h.Memory == nil || h.Memory.DB() == nil {
			return
		}
		db := h.Memory.DB()
		_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS team_context (id INTEGER PRIMARY KEY, person TEXT NOT NULL, context_type TEXT NOT NULL, content TEXT NOT NULL, sentiment TEXT DEFAULT 'neutral', created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`)
		_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS learning_log (id INTEGER PRIMARY KEY, category TEXT NOT NULL, observation TEXT NOT NULL, lesson TEXT, applied INTEGER DEFAULT 0, created_at DATETIME DEFAULT CURRENT_TIMESTAMP)`)
		_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_team_context_person ON team_context(person)`)
	})
}

func (h *Handlers) TeamContext(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initEmpathyTables()
	person, err := req.RequireString("person")
	if err != nil {
		return errorResult("person required"), nil
	}
	content, err := req.RequireString("note")
	if err != nil {
		return errorResult("note required"), nil
	}
	ctxType := req.GetString("type", "observation")
	sentiment := req.GetString("sentiment", "neutral")
	db := h.Memory.DB()
	_, _ = db.Exec("INSERT INTO team_context (person, context_type, content, sentiment) VALUES (?, ?, ?, ?)", person, ctxType, content, sentiment)
	return textResult(fmt.Sprintf("Noted about %s: %s", person, content)), nil
}

func (h *Handlers) TeamRecall(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initEmpathyTables()
	person, err := req.RequireString("person")
	if err != nil {
		return errorResult("person required"), nil
	}
	db := h.Memory.DB()
	rows, qErr := db.Query("SELECT context_type, content, sentiment, created_at FROM team_context WHERE person = ? ORDER BY created_at DESC LIMIT 20", person)
	if qErr != nil {
		return errorResult("Query failed"), nil
	}
	defer rows.Close()
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Context about %s:\n\n", person))
	count := 0
	for rows.Next() {
		var ctxType, content, sentiment, created string
		if err := rows.Scan(&ctxType, &content, &sentiment, &created); err != nil {
			continue
		}
		count++
		date := created
		if len(created) >= 10 {
			date = created[:10]
		}
		sb.WriteString(fmt.Sprintf("  [%s] %s (%s, %s)\n", ctxType, content, sentiment, date))
	}
	if count == 0 {
		return textResult(fmt.Sprintf("No notes about %s yet.", person)), nil
	}
	return textResult(sb.String()), nil
}

func (h *Handlers) SentimentCheck(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	message, err := req.RequireString("message")
	if err != nil {
		return errorResult("message required"), nil
	}
	if h.AI == nil {
		return textResult(fmt.Sprintf("Sentiment: %s (keyword-based, AI unavailable)", detectSentimentSimple(message))), nil
	}
	systemPrompt := "Analyze tone and sentiment. Output:\nSENTIMENT: [positive/neutral/negative/frustrated/anxious/excited]\nTONE: [supportive/demanding/passive-aggressive/direct/tentative]\nIMPACT: How the receiver likely feels\nSUGGESTION: One alternative if needed. If good, say \"Good as-is.\"\nMax 80 words. Direct. Never patronize."
	result, aiErr := h.AI.Complete(ctx, systemPrompt, message)
	if aiErr != nil {
		return errorResult("AI failed: " + aiErr.Error()), nil
	}
	return textResult(result), nil
}

func (h *Handlers) LearnFromSprint(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initEmpathyTables()
	category, err := req.RequireString("category")
	if err != nil {
		return errorResult("category required: process, communication, technical, people, delivery"), nil
	}
	observation, err := req.RequireString("observation")
	if err != nil {
		return errorResult("observation required"), nil
	}
	lesson := req.GetString("lesson", "")
	db := h.Memory.DB()
	_, _ = db.Exec("INSERT INTO learning_log (category, observation, lesson) VALUES (?, ?, ?)", category, observation, lesson)
	return textResult(fmt.Sprintf("Learned: [%s] %s", category, observation)), nil
}

func (h *Handlers) WisdomRecall(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initEmpathyTables()
	category := req.GetString("category", "")
	limit := req.GetInt("limit", 10)
	db := h.Memory.DB()
	var rows interface{ Next() bool; Scan(...any) error; Close() error }
	var err error
	if category != "" {
		rows, err = db.Query("SELECT category, observation, lesson, created_at FROM learning_log WHERE category = ? ORDER BY created_at DESC LIMIT ?", category, limit)
	} else {
		rows, err = db.Query("SELECT category, observation, lesson, created_at FROM learning_log ORDER BY created_at DESC LIMIT ?", limit)
	}
	if err != nil {
		return errorResult("Query failed"), nil
	}
	defer rows.Close()
	var sb strings.Builder
	sb.WriteString("Accumulated Wisdom:\n\n")
	count := 0
	for rows.Next() {
		var cat, obs, lesson, created string
		if err := rows.Scan(&cat, &obs, &lesson, &created); err != nil {
			continue
		}
		count++
		sb.WriteString(fmt.Sprintf("  [%s] %s\n", cat, obs))
		if lesson != "" {
			sb.WriteString(fmt.Sprintf("    -> %s\n", lesson))
		}
	}
	if count == 0 {
		return textResult("No wisdom recorded yet. Use pm_learn after each sprint."), nil
	}
	return textResult(sb.String()), nil
}

func (h *Handlers) TeamMood(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initEmpathyTables()
	person, err := req.RequireString("person")
	if err != nil {
		return errorResult("person required"), nil
	}
	mood := req.GetString("mood", "")
	if mood == "" {
		return errorResult("mood required: happy, neutral, stressed, frustrated, excited, tired, focused"), nil
	}
	note := req.GetString("note", "")
	sentiment := "neutral"
	switch mood {
	case "happy", "excited", "focused":
		sentiment = "positive"
	case "stressed", "frustrated", "tired":
		sentiment = "negative"
	}
	content := "Mood: " + mood
	if note != "" {
		content += " - " + note
	}
	db := h.Memory.DB()
	_, _ = db.Exec("INSERT INTO team_context (person, context_type, content, sentiment) VALUES (?, ?, ?, ?)", person, "mood", content, sentiment)
	var negCount int
	cutoff := time.Now().AddDate(0, 0, -14)
	row := db.QueryRow("SELECT COUNT(*) FROM team_context WHERE person = ? AND sentiment = 'negative' AND created_at > ?", person, cutoff.Format("2006-01-02"))
	_ = row.Scan(&negCount)
	response := fmt.Sprintf("Recorded: %s feeling %s.", person, mood)
	if negCount >= 3 {
		response += fmt.Sprintf("\n\nHeads up: %s has shown %d negative signals in 2 weeks. Consider a 1-on-1.", person, negCount)
	}
	return textResult(response), nil
}

func detectSentimentSimple(msg string) string {
	lower := strings.ToLower(msg)
	neg := []string{"frustrated", "angry", "stuck", "blocked", "impossible", "stressed", "tired", "annoyed", "overwhelmed"}
	pos := []string{"great", "awesome", "done", "shipped", "fixed", "happy", "excited", "progress", "resolved"}
	n, p := 0, 0
	for _, w := range neg {
		if strings.Contains(lower, w) {
			n++
		}
	}
	for _, w := range pos {
		if strings.Contains(lower, w) {
			p++
		}
	}
	if n > p {
		return "negative"
	}
	if p > n {
		return "positive"
	}
	return "neutral"
}
