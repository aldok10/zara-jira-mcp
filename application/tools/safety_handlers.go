package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handlers) initSafetyTable() {
	if h.Memory == nil {
		return
	}
	db := h.Memory.DB()
	if db == nil {
		return
	}
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS safety_survey (
		id INTEGER PRIMARY KEY,
		sprint_name TEXT,
		member TEXT,
		q1 INTEGER, q2 INTEGER, q3 INTEGER, q4 INTEGER, q5 INTEGER, q6 INTEGER, q7 INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
}

// PMSafetySurvey records psychological safety survey responses.
func (h *Handlers) PMSafetySurvey(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	sprintName, err := req.RequireString("sprint_name")
	if err != nil {
		return errorResult("sprint_name required"), nil
	}
	responsesStr, err := req.RequireString("responses")
	if err != nil {
		return errorResult("responses required (JSON: {\"member\": {\"q1\": 4, \"q2\": 3, ...}})"), nil
	}

	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}

	h.initSafetyTable()

	var responses map[string]map[string]int
	if err := json.Unmarshal([]byte(responsesStr), &responses); err != nil {
		return errorResult("invalid JSON for responses: " + err.Error()), nil
	}

	db := h.Memory.DB()
	var totalScore, count int

	for member, scores := range responses {
		q1 := scores["q1"]
		q2 := scores["q2"]
		q3 := scores["q3"]
		q4 := scores["q4"]
		q5 := scores["q5"]
		q6 := scores["q6"]
		q7 := scores["q7"]

		_, err := db.Exec(
			`INSERT INTO safety_survey (sprint_name, member, q1, q2, q3, q4, q5, q6, q7) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			sprintName, member, q1, q2, q3, q4, q5, q6, q7,
		)
		if err != nil {
			return errorResult("failed to store survey: " + err.Error()), nil
		}

		memberTotal := q1 + q2 + q3 + q4 + q5 + q6 + q7
		totalScore += memberTotal
		count++
	}

	avg := float64(0)
	if count > 0 {
		avg = float64(totalScore) / float64(count*7)
	}

	return textResult(fmt.Sprintf("Safety survey recorded for sprint %s\nResponses: %d members\nAverage score: %.1f/5\n\nQuestions (Project Aristotle):\n1. Safe to take risks\n2. Can depend on each other\n3. Clear roles and goals\n4. Work is meaningful\n5. Work has impact\n6. Can raise problems\n7. No blame for mistakes",
		sprintName, count, avg)), nil
}

// PMSafetyTrend shows psychological safety trends over time.
func (h *Handlers) PMSafetyTrend(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}

	h.initSafetyTable()
	db := h.Memory.DB()

	rows, err := db.Query(`SELECT sprint_name, AVG(q1+q2+q3+q4+q5+q6+q7)/7.0 as avg_score, COUNT(*) as responses
		FROM safety_survey GROUP BY sprint_name ORDER BY created_at DESC LIMIT 10`)
	if err != nil {
		return errorResult("query failed: " + err.Error()), nil
	}
	defer rows.Close()

	var sb strings.Builder
	sb.WriteString("Psychological Safety Trend\n\n")
	sb.WriteString("Sprint | Avg Score | Responses\n")
	sb.WriteString("-------|-----------|----------\n")

	found := false
	for rows.Next() {
		var sprint string
		var avg float64
		var count int
		if err := rows.Scan(&sprint, &avg, &count); err != nil {
			continue
		}
		found = true
		sb.WriteString(fmt.Sprintf("%-12s | %.1f/5 | %d\n", sprint, avg, count))
	}

	if !found {
		return textResult("No safety survey data yet. Use pm_safety_survey to record responses."), nil
	}

	return textResult(sb.String()), nil
}

// PMTeamAristotle performs a full 5-pillar assessment based on Google's Project Aristotle.
func (h *Handlers) PMTeamAristotle(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	if h.AI == nil {
		return errorResult("AI provider not configured"), nil
	}

	var data strings.Builder
	data.WriteString("5 Pillars Assessment Data:\n\n")

	// 1. Psychological Safety - from survey
	if h.Memory != nil {
		h.initSafetyTable()
		db := h.Memory.DB()
		row := db.QueryRow(`SELECT AVG(q1+q2+q3+q4+q5+q6+q7)/7.0, COUNT(*) FROM safety_survey`)
		var avg float64
		var count int
		if err := row.Scan(&avg, &count); err == nil && count > 0 {
			data.WriteString(fmt.Sprintf("1. SAFETY: Avg score %.1f/5 from %d responses\n", avg, count))
		} else {
			data.WriteString("1. SAFETY: No survey data\n")
		}

		// 2. Dependability - from team metrics (completion rates)
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
		if len(snaps) > 0 {
			var totalRate float64
			for _, s := range snaps {
				totalRate += s.CompletionRate
			}
			data.WriteString(fmt.Sprintf("2. DEPENDABILITY: Avg completion %.0f%% over %d sprints\n", totalRate/float64(len(snaps))*100, len(snaps)))
		} else {
			data.WriteString("2. DEPENDABILITY: No sprint data\n")
		}

		// 3. Structure & Clarity - sprint goals
		goals, _ := h.Memory.GetGoalHistory(ctx, boardID, 5)
		achieved := 0
		for _, g := range goals {
			if g.Status == "achieved" {
				achieved++
			}
		}
		if len(goals) > 0 {
			data.WriteString(fmt.Sprintf("3. CLARITY: %d/%d sprint goals achieved\n", achieved, len(goals)))
		} else {
			data.WriteString("3. CLARITY: No sprint goals recorded\n")
		}

		// 4. Meaning - team pulse
		pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 10)
		if len(pulses) > 0 {
			var total int
			for _, p := range pulses {
				total += p.Score
			}
			data.WriteString(fmt.Sprintf("4. MEANING: Team pulse avg %.1f/5 (%d entries)\n", float64(total)/float64(len(pulses)), len(pulses)))
		} else {
			data.WriteString("4. MEANING: No pulse data\n")
		}

		// 5. Impact - health scores
		health, _ := h.Memory.GetHealthScores(ctx, boardID, 5)
		if len(health) > 0 {
			var total int
			for _, hs := range health {
				total += hs.OverallScore
			}
			data.WriteString(fmt.Sprintf("5. IMPACT: Avg health score %d/100 over %d sprints\n", total/len(health), len(health)))
		} else {
			data.WriteString("5. IMPACT: No health score data\n")
		}
	}

	systemPrompt := `Assess this team using Google's Project Aristotle 5 pillars:
1. Psychological Safety - safe to take risks?
2. Dependability - can count on each other?
3. Structure & Clarity - clear goals and roles?
4. Meaning - work is personally meaningful?
5. Impact - work matters?

For each pillar give a score (1-5) based on available data, a brief assessment, and one recommendation.
End with an overall team health verdict and top 2 actions.
Under 250 words.`

	result, err := h.aiComplete(ctx, systemPrompt, data.String())
	if err != nil {
		return errorResult("AI failed: " + err.Error()), nil
	}
	return textResult(fmt.Sprintf("[Project Aristotle Assessment]\n\n%s", result)), nil
}
