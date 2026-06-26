package tools

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

var v6TablesOnce sync.Once

func (h *Handlers) initV6Tables() {
	v6TablesOnce.Do(func() {
		if h.Memory == nil || h.Memory.DB() == nil {
			return
		}
		_, _ = h.Memory.DB().Exec(`CREATE TABLE IF NOT EXISTS hypothesis (
			id INTEGER PRIMARY KEY, belief TEXT NOT NULL, expected_outcome TEXT NOT NULL,
			measure TEXT DEFAULT 'observe', duration TEXT DEFAULT '1 sprint',
			status TEXT DEFAULT 'active', actual_outcome TEXT, sprint_name TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP, closed_at DATETIME)`)
	})
}

func (h *Handlers) PMHypothesis(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initV6Tables()
	belief, _ := req.RequireString("belief")
	expected, _ := req.RequireString("expected_outcome")
	measure := req.GetString("measure", "observe")
	duration := req.GetString("duration", "1 sprint")
	sprint := req.GetString("sprint_name", "")
	_, _ = h.Memory.DB().Exec("INSERT INTO hypothesis(belief,expected_outcome,measure,duration,sprint_name) VALUES(?,?,?,?,?)", belief, expected, measure, duration, sprint)
	return textResult(fmt.Sprintf("Hypothesis recorded:\n  If: %s\n  Then: %s\n  Measure: %s (%s)", belief, expected, measure, duration)), nil
}

func (h *Handlers) PMHypothesisReview(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initV6Tables()
	rows, err := h.Memory.DB().Query("SELECT id,belief,expected_outcome,measure,status,actual_outcome FROM hypothesis ORDER BY created_at DESC LIMIT 20")
	if err != nil {
		return errorResult("query failed"), nil
	}
	defer rows.Close()
	var sb strings.Builder
	sb.WriteString("Hypotheses:\n\n")
	count := 0
	for rows.Next() {
		var id int
		var belief, expected, measure, status, actual string
		_ = rows.Scan(&id, &belief, &expected, &measure, &status, &actual)
		icon := "?"
		if status == "validated" {
			icon = "Y"
		} else if status == "invalidated" {
			icon = "N"
		}
		sb.WriteString(fmt.Sprintf("  #%d [%s] %s -> %s (measure: %s)\n", id, icon, belief, expected, measure))
		if actual != "" {
			sb.WriteString(fmt.Sprintf("       Actual: %s\n", actual))
		}
		count++
	}
	if count == 0 {
		return textResult("No hypotheses yet."), nil
	}
	return textResult(sb.String()), nil
}

func (h *Handlers) PMHypothesisClose(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initV6Tables()
	id, _ := req.RequireInt("id")
	status, _ := req.RequireString("status")
	actual := req.GetString("actual_outcome", "")
	_, _ = h.Memory.DB().Exec("UPDATE hypothesis SET status=?,actual_outcome=?,closed_at=? WHERE id=?", status, actual, time.Now(), id)
	return textResult(fmt.Sprintf("Hypothesis #%d: %s", id, status)), nil
}

func (h *Handlers) PMEstimationAccuracy(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	boardID, _ := req.RequireInt("board_id")
	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 10)
	if len(snaps) < 3 {
		return textResult("Need 3+ snapshots."), nil
	}
	var sb strings.Builder
	sb.WriteString("Estimation Accuracy:\n\n")
	totalC, totalD := 0, 0
	for _, s := range snaps {
		if s.TotalIssues == 0 {
			continue
		}
		pct := float64(s.Done) / float64(s.TotalIssues) * 100
		totalC += s.TotalIssues
		totalD += s.Done
		sb.WriteString(fmt.Sprintf("  %s: %d/%d (%.0f%%)\n", s.SprintName, s.Done, s.TotalIssues, pct))
	}
	if totalC > 0 {
		o := float64(totalD) / float64(totalC) * 100
		sb.WriteString(fmt.Sprintf("\nDelivery rate: %.0f%%\n", o))
	}
	return textResult(sb.String()), nil
}

func (h *Handlers) PMSpaceMetrics(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	boardID := req.GetInt("board_id", 0)
	var sb strings.Builder
	sb.WriteString("SPACE:\n")
	pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 20)
	if len(pulses) > 0 {
		t := 0
		for _, p := range pulses {
			t += p.Score
		}
		sb.WriteString(fmt.Sprintf("  S: %.1f/5\n", float64(t)/float64(len(pulses))))
	}
	if boardID > 0 {
		goals, _ := h.Memory.GetGoalHistory(ctx, boardID, 10)
		hit := 0
		for _, g := range goals {
			if g.Status == "achieved" {
				hit++
			}
		}
		if len(goals) > 0 {
			sb.WriteString(fmt.Sprintf("  P: %d/%d goals\n", hit, len(goals)))
		}
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
		if len(snaps) > 0 {
			t := 0
			for _, s := range snaps {
				t += s.Done
			}
			sb.WriteString(fmt.Sprintf("  A: %.0f items/sprint\n", float64(t)/float64(len(snaps))))
		}
	}
	decisions, _ := h.Memory.GetDecisions(ctx, 30)
	sb.WriteString(fmt.Sprintf("  C: %d decisions\n", len(decisions)))
	history, _ := h.Memory.GetBlockerHistory(ctx, 20)
	r, d := 0, 0.0
	for _, b := range history {
		if b.ResolvedAt != nil {
			r++
			d += b.ResolvedAt.Sub(b.BlockedSince).Hours() / 24
		}
	}
	if r > 0 {
		sb.WriteString(fmt.Sprintf("  E: %.1fd blocker avg\n", d/float64(r)))
	}
	return textResult(sb.String()), nil
}

func (h *Handlers) PMEBMDashboard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	boardID := req.GetInt("board_id", 0)
	var sb strings.Builder
	sb.WriteString("EBM (4 KVAs):\n\n")
	sb.WriteString("1. CURRENT VALUE\n")
	pulses, _ := h.Memory.GetTeamPulseHistory(ctx, 20)
	if len(pulses) > 0 {
		t := 0
		for _, p := range pulses {
			t += p.Score
		}
		sb.WriteString(fmt.Sprintf("   Satisfaction: %.1f/5\n", float64(t)/float64(len(pulses))))
	}
	sb.WriteString("\n2. UNREALIZED VALUE\n")
	if boardID > 0 {
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 5)
		if len(snaps) >= 2 {
			sb.WriteString(fmt.Sprintf("   Scope: %d -> %d items\n", snaps[len(snaps)-1].TotalIssues, snaps[0].TotalIssues))
		}
	}
	sb.WriteString("\n3. ABILITY TO INNOVATE\n")
	var debt int
	if h.Memory.DB() != nil {
		row := h.Memory.DB().QueryRow("SELECT COUNT(*) FROM tech_debt WHERE status='open'")
		_ = row.Scan(&debt)
	}
	sb.WriteString(fmt.Sprintf("   Tech debt: %d open\n", debt))
	sb.WriteString("\n4. TIME TO MARKET\n")
	if boardID > 0 {
		snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 3)
		if len(snaps) > 0 {
			sb.WriteString(fmt.Sprintf("   Velocity: %d pts\n", snaps[0].Velocity))
		}
	}
	return textResult(sb.String()), nil
}
