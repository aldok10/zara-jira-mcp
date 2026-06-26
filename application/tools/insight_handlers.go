package tools

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
)

var insightTablesOnce sync.Once

func (h *Handlers) initInsightTables() {
	insightTablesOnce.Do(func() {
		if h.Memory == nil || h.Memory.DB() == nil {
			return
		}
		db := h.Memory.DB()
		_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS tool_usage (
			id INTEGER PRIMARY KEY,
			tool_name TEXT NOT NULL,
			called_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
		_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS notification_log (
			id INTEGER PRIMARY KEY,
			channel TEXT NOT NULL,
			severity TEXT NOT NULL,
			title_hash TEXT,
			sent_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
		_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_tool_usage_name ON tool_usage(tool_name)`)
		_, _ = db.Exec(`CREATE INDEX IF NOT EXISTS idx_notification_log_day ON notification_log(sent_at)`)
	})
}

// TrackToolUsage records a tool call.
func (h *Handlers) TrackToolUsage(toolName string) {
	if h.Memory == nil || h.Memory.DB() == nil {
		return
	}
	h.initInsightTables()
	_, _ = h.Memory.DB().Exec("INSERT INTO tool_usage(tool_name) VALUES(?)", toolName)
}

// PMToolUsage shows which tools are used most/least.
func (h *Handlers) PMToolUsage(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initInsightTables()
	days := req.GetInt("days", 30)

	rows, err := h.Memory.DB().Query(`
		SELECT tool_name, COUNT(*) as cnt 
		FROM tool_usage 
		WHERE called_at > datetime('now', '-' || ? || ' days')
		GROUP BY tool_name ORDER BY cnt DESC LIMIT 20`, days)
	if err != nil {
		return sanitizedError("insight query failed", err), nil
	}
	defer rows.Close()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Tool Usage (last %d days):\n\n", days))
	total := 0
	for rows.Next() {
		var name string
		var cnt int
		_ = rows.Scan(&name, &cnt)
		sb.WriteString(fmt.Sprintf("  %s: %d calls\n", name, cnt))
		total += cnt
	}
	if total == 0 {
		return textResult("No usage data yet. Tools are tracked automatically."), nil
	}
	sb.WriteString(fmt.Sprintf("\nTotal: %d calls", total))
	return textResult(sb.String()), nil
}

// PMCalibrationReport compares forecast accuracy vs actual outcomes.
func (h *Handlers) PMCalibrationReport(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	snaps, _ := h.Memory.GetSprintSnapshots(ctx, boardID, 10)
	if len(snaps) < 3 {
		return textResult("Need 3+ sprint snapshots for calibration."), nil
	}

	var sb strings.Builder
	sb.WriteString("Forecast Calibration Report:\n\n")
	totalCommitted, totalDelivered := 0, 0

	for _, s := range snaps {
		if s.TotalIssues == 0 {
			continue
		}
		accuracy := float64(s.Done) / float64(s.TotalIssues) * 100
		totalCommitted += s.TotalIssues
		totalDelivered += s.Done
		sb.WriteString(fmt.Sprintf("  %s: %d/%d (%.0f%%)\n", s.SprintName, s.Done, s.TotalIssues, accuracy))
	}

	if totalCommitted > 0 {
		overall := float64(totalDelivered) / float64(totalCommitted) * 100
		sb.WriteString(fmt.Sprintf("\nOverall Predictability: %.0f%%\n", overall))
		if overall >= 80 {
			sb.WriteString("Rating: HIGH\n")
		} else if overall >= 60 {
			sb.WriteString("Rating: MEDIUM — reduce scope 10-20%%\n")
		} else {
			sb.WriteString("Rating: LOW — use pm_capacity_plan\n")
		}
	}
	return textResult(sb.String()), nil
}

// PMMeetingROI tracks meeting effectiveness.
func (h *Handlers) PMMeetingROI(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}

	notes, _ := h.Memory.GetMeetingNotes(ctx, "", 20)
	if len(notes) == 0 {
		return textResult("No meeting data. Use pm_record_meeting."), nil
	}

	var sb strings.Builder
	sb.WriteString("Meeting ROI:\n\n")

	type stats struct{ count, decisions, actions int }
	typeStats := map[string]*stats{}

	for _, n := range notes {
		if typeStats[n.MeetingType] == nil {
			typeStats[n.MeetingType] = &stats{}
		}
		s := typeStats[n.MeetingType]
		s.count++
		if n.Decisions != "" {
			s.decisions += len(strings.Split(n.Decisions, "\n"))
		}
		if n.ActionItems != "" {
			s.actions += len(strings.Split(n.ActionItems, "\n"))
		}
	}

	for mType, s := range typeStats {
		output := float64(s.decisions+s.actions) / float64(s.count)
		rating := "low"
		if output >= 3 {
			rating = "high"
		} else if output >= 1.5 {
			rating = "medium"
		}
		sb.WriteString(fmt.Sprintf("  %s: %.1f outputs/meeting [%s]\n", mType, output, rating))
	}
	sb.WriteString("\nLow = consider async.")
	return textResult(sb.String()), nil
}

// PMNotificationBudgetCheck shows notification budget status.
func (h *Handlers) PMNotificationBudgetCheck(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Memory == nil {
		return errorResult("memory not configured"), nil
	}
	h.initInsightTables()

	var today, week int
	row := h.Memory.DB().QueryRow("SELECT COUNT(*) FROM notification_log WHERE sent_at > datetime('now', '-1 day')")
	_ = row.Scan(&today)
	row = h.Memory.DB().QueryRow("SELECT COUNT(*) FROM notification_log WHERE sent_at > datetime('now', '-7 days')")
	_ = row.Scan(&week)

	remaining := 5 - today
	if remaining < 0 {
		remaining = 0
	}
	status := "within budget"
	if today >= 5 {
		status = "OVER BUDGET"
	}
	return textResult(fmt.Sprintf("Notifications: %d/5 today [%s], %d this week, %d remaining", today, status, week, remaining)), nil
}

// CheckNotificationBudget returns true if within daily budget.
func (h *Handlers) CheckNotificationBudget() bool {
	if h.Memory == nil || h.Memory.DB() == nil {
		return true
	}
	h.initInsightTables()
	var count int
	row := h.Memory.DB().QueryRow("SELECT COUNT(*) FROM notification_log WHERE sent_at > datetime('now', '-1 day')")
	_ = row.Scan(&count)
	return count < 5
}

// LogNotification records a sent notification.
func (h *Handlers) LogNotification(channel, severity, title string) {
	if h.Memory == nil || h.Memory.DB() == nil {
		return
	}
	h.initInsightTables()
	_, _ = h.Memory.DB().Exec("INSERT INTO notification_log(channel, severity, title_hash) VALUES(?, ?, ?)", channel, severity, title)
}
