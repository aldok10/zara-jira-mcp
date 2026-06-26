package tools

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/mark3labs/mcp-go/mcp"
)

var okrTablesOnce sync.Once

func (h *Handlers) initOKRTables() {
	okrTablesOnce.Do(func() {
		db := h.Memory.DB()
		_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS okr (
			id INTEGER PRIMARY KEY,
			level TEXT NOT NULL DEFAULT 'team',
			title TEXT NOT NULL,
			description TEXT,
			owner TEXT,
			cycle TEXT,
			board_id INTEGER,
			status TEXT DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
		_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS key_result (
			id INTEGER PRIMARY KEY,
			okr_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			start_value REAL DEFAULT 0,
			target_value REAL DEFAULT 100,
			current_value REAL DEFAULT 0,
			unit TEXT DEFAULT '%',
			status TEXT DEFAULT 'active',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (okr_id) REFERENCES okr(id)
		)`)
		_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS kr_jira_link (
			id INTEGER PRIMARY KEY,
			kr_id INTEGER NOT NULL,
			issue_key TEXT NOT NULL,
			link_type TEXT DEFAULT 'contributes',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (kr_id) REFERENCES key_result(id)
		)`)
		_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS kpi_definition (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			formula TEXT,
			unit TEXT DEFAULT '%',
			target_value REAL,
			warning_threshold REAL,
			danger_threshold REAL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`)
		_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS kpi_snapshot (
			id INTEGER PRIMARY KEY,
			kpi_id INTEGER NOT NULL,
			value REAL NOT NULL,
			sprint_name TEXT,
			notes TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (kpi_id) REFERENCES kpi_definition(id)
		)`)
	})
}

// PMOKRDefine creates an OKR with key results.
func (h *Handlers) PMOKRDefine(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	title, err := req.RequireString("title")
	if err != nil {
		return errorResult("title required"), nil
	}
	level := req.GetString("level", "team")
	owner := req.GetString("owner", "")
	cycle := req.GetString("cycle", "")
	description := req.GetString("description", "")
	keyResults := req.GetString("key_results", "")

	db := h.Memory.DB()
	result, execErr := db.Exec("INSERT INTO okr (level, title, description, owner, cycle) VALUES (?, ?, ?, ?, ?)",
		level, title, description, owner, cycle)
	if execErr != nil {
		return errorResult("Failed: " + execErr.Error()), nil
	}

	type lastIDer interface{ LastInsertId() (int64, error) }
	var okrID int64
	if r, ok := result.(lastIDer); ok {
		okrID, _ = r.LastInsertId()
	}

	// Parse key results: "KR title | target_value | unit" per line
	var krCount int
	if keyResults != "" && okrID > 0 {
		for _, line := range strings.Split(keyResults, "\n") {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.Split(line, "|")
			krTitle := strings.TrimSpace(parts[0])
			var targetVal float64 = 100
			unit := "%"
			if len(parts) > 1 {
				fmt.Sscanf(strings.TrimSpace(parts[1]), "%f", &targetVal)
			}
			if len(parts) > 2 {
				unit = strings.TrimSpace(parts[2])
			}
			_, _ = db.Exec("INSERT INTO key_result (okr_id, title, target_value, unit) VALUES (?, ?, ?, ?)",
				okrID, krTitle, targetVal, unit)
			krCount++
		}
	}

	msg := fmt.Sprintf("OKR #%d created: '%s' (level: %s)", okrID, title, level)
	if krCount > 0 {
		msg += fmt.Sprintf("\n%d Key Results defined.", krCount)
	}
	return textResult(msg), nil
}

// PMOKRList shows all OKRs with KR progress.
func (h *Handlers) PMOKRList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	status := req.GetString("status", "active")
	db := h.Memory.DB()

	rows, err := db.Query("SELECT id, level, title, owner, cycle FROM okr WHERE status = ? ORDER BY level, created_at", status)
	if err != nil {
		return errorResult("Query failed: " + err.Error()), nil
	}
	defer rows.Close()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("OKRs (%s)\n\n", status))

	count := 0
	for rows.Next() {
		var id int64
		var level, title, owner, cycle string
		if err := rows.Scan(&id, &level, &title, &owner, &cycle); err != nil {
			continue
		}
		count++
		sb.WriteString(fmt.Sprintf("#%d [%s] %s", id, strings.ToUpper(level), title))
		if owner != "" {
			sb.WriteString(fmt.Sprintf(" (@%s)", owner))
		}
		if cycle != "" {
			sb.WriteString(fmt.Sprintf(" [%s]", cycle))
		}
		sb.WriteString("\n")

		krRows, _ := db.Query("SELECT id, title, current_value, target_value, unit FROM key_result WHERE okr_id = ? AND status = 'active'", id)
		if krRows != nil {
			for krRows.Next() {
				var krID int64
				var krTitle, unit string
				var current, target float64
				if err := krRows.Scan(&krID, &krTitle, &current, &target, &unit); err == nil {
					pct := 0.0
					if target > 0 {
						pct = current / target * 100
					}
					sb.WriteString(fmt.Sprintf("  KR#%d: %s — %.0f/%.0f%s (%.0f%%)\n", krID, krTitle, current, target, unit, pct))
				}
			}
			krRows.Close()
		}
	}

	if count == 0 {
		return textResult("No OKRs found. Use pm_okr_define to create one."), nil
	}
	return textResult(sb.String()), nil
}

// PMKRLink links Jira issues to a Key Result.
func (h *Handlers) PMKRLink(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	krID, err := req.RequireInt("kr_id")
	if err != nil {
		return errorResult("kr_id required"), nil
	}
	issueKeys, err := req.RequireString("issue_keys")
	if err != nil {
		return errorResult("issue_keys required (comma-separated)"), nil
	}
	linkType := req.GetString("link_type", "contributes")

	db := h.Memory.DB()
	row := db.QueryRow("SELECT title FROM key_result WHERE id = ?", krID)
	var krTitle string
	if scanErr := row.Scan(&krTitle); scanErr != nil {
		return errorResult(fmt.Sprintf("Key Result #%d not found", krID)), nil
	}

	linked := 0
	for _, key := range strings.Split(issueKeys, ",") {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		_, _ = db.Exec("INSERT INTO kr_jira_link (kr_id, issue_key, link_type) VALUES (?, ?, ?)", krID, key, linkType)
		linked++
	}

	return textResult(fmt.Sprintf("Linked %d issues to KR#%d '%s'.", linked, krID, krTitle)), nil
}

// PMKRProgress calculates KR progress from Jira data.
func (h *Handlers) PMKRProgress(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	db := h.Memory.DB()
	okrRows, err := db.Query("SELECT id, title, level FROM okr WHERE status = 'active' ORDER BY created_at DESC")
	if err != nil {
		return errorResult("Query failed: " + err.Error()), nil
	}
	defer okrRows.Close()

	type okrEntry struct{ id int64; title, level string }
	var okrs []okrEntry
	for okrRows.Next() {
		var o okrEntry
		if err := okrRows.Scan(&o.id, &o.title, &o.level); err == nil {
			okrs = append(okrs, o)
		}
	}
	if len(okrs) == 0 {
		return textResult("No active OKRs. Use pm_okr_define to create one."), nil
	}

	var sb strings.Builder
	sb.WriteString("OKR Progress Report\n\n")

	for _, o := range okrs {
		sb.WriteString(fmt.Sprintf("[%s] %s\n", strings.ToUpper(o.level), o.title))

		krRows, _ := db.Query("SELECT id, title, start_value, target_value, current_value, unit FROM key_result WHERE okr_id = ? AND status = 'active'", o.id)
		if krRows == nil {
			continue
		}
		krCount := 0
		totalPct := 0.0
		for krRows.Next() {
			var krID int64
			var krTitle, unit string
			var startVal, targetVal, currentVal float64
			if err := krRows.Scan(&krID, &krTitle, &startVal, &targetVal, &currentVal, &unit); err != nil {
				continue
			}

			// Auto-calculate from linked Jira issues
			linkRows, _ := db.Query("SELECT issue_key FROM kr_jira_link WHERE kr_id = ?", krID)
			if linkRows != nil {
				var keys []string
				for linkRows.Next() {
					var k string
					if err := linkRows.Scan(&k); err == nil {
						keys = append(keys, k)
					}
				}
				linkRows.Close()

				if len(keys) > 0 && h.Jira != nil {
					done := 0
					for _, k := range keys {
						issue, err := h.Jira.GetIssue(ctx, k)
						if err == nil && isDoneStatus(issue.Status) {
							done++
						}
					}
					currentVal = float64(done) / float64(len(keys)) * targetVal
					_, _ = db.Exec("UPDATE key_result SET current_value = ? WHERE id = ?", currentVal, krID)
				}
			}

			pct := 0.0
			if targetVal-startVal > 0 {
				pct = (currentVal - startVal) / (targetVal - startVal) * 100
			}
			if pct > 100 { pct = 100 }
			if pct < 0 { pct = 0 }
			totalPct += pct
			krCount++

			bar := progressBar(pct)
			sb.WriteString(fmt.Sprintf("  KR#%d: %s\n       %s %.0f%%\n", krID, krTitle, bar, pct))
		}
		krRows.Close()

		if krCount > 0 {
			sb.WriteString(fmt.Sprintf("  Overall: %.0f%%\n", totalPct/float64(krCount)))
		}
		sb.WriteString("\n")
	}

	return textResult(sb.String()), nil
}

// PMOutcomeReview uses AI to assess if sprint work moved OKRs.
func (h *Handlers) PMOutcomeReview(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	db := h.Memory.DB()
	var sb strings.Builder

	// Gather OKR progress data
	okrRows, _ := db.Query("SELECT title, level FROM okr WHERE status = 'active'")
	if okrRows != nil {
		sb.WriteString("Active OKRs:\n")
		for okrRows.Next() {
			var title, level string
			if err := okrRows.Scan(&title, &level); err == nil {
				sb.WriteString(fmt.Sprintf("  [%s] %s\n", level, title))
			}
		}
		okrRows.Close()
	}

	// Gather outcome_map data
	mapRows, _ := db.Query("SELECT sprint_name, objective FROM outcome_map WHERE board_id = ? ORDER BY created_at DESC LIMIT 5", boardID)
	if mapRows != nil {
		sb.WriteString("\nSprint-Objective Mappings:\n")
		for mapRows.Next() {
			var sprint, obj string
			if err := mapRows.Scan(&sprint, &obj); err == nil {
				sb.WriteString(fmt.Sprintf("  %s → %s\n", sprint, obj))
			}
		}
		mapRows.Close()
	}

	// Sprint goals
	goals, _ := h.Memory.GetGoalHistory(ctx, boardID, 3)
	if len(goals) > 0 {
		sb.WriteString("\nRecent Sprint Goals:\n")
		for _, g := range goals {
			sb.WriteString(fmt.Sprintf("  %s: %s [%s]\n", g.SprintName, g.Goal, g.Status))
		}
	}

	if h.AI == nil {
		return textResult(sb.String() + "\n(AI not configured — cannot generate assessment)"), nil
	}

	prompt := `Based on this sprint/OKR data, answer concisely:
1. Are sprints clearly connected to business objectives?
2. Are KRs actually progressing?
3. What's the biggest gap between effort and outcome?
Be direct. 3-5 sentences max.`

	aiResult, aiErr := h.AI.Complete(ctx, prompt, sb.String())
	if aiErr != nil {
		return textResult(sb.String() + "\n(AI assessment failed)"), nil
	}

	return textResult(sb.String() + "\n\nAI Assessment:\n" + aiResult), nil
}

// PMKPIDefine creates a KPI definition.
func (h *Handlers) PMKPIDefine(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	name, err := req.RequireString("name")
	if err != nil {
		return errorResult("name required"), nil
	}
	description := req.GetString("description", "")
	formula := req.GetString("formula", "")
	unit := req.GetString("unit", "%")
	targetValue := req.GetFloat("target_value", 0)
	warningThreshold := req.GetFloat("warning_threshold", 0)
	dangerThreshold := req.GetFloat("danger_threshold", 0)

	db := h.Memory.DB()
	_, execErr := db.Exec(`INSERT INTO kpi_definition (name, description, formula, unit, target_value, warning_threshold, danger_threshold)
		VALUES (?, ?, ?, ?, ?, ?, ?)`, name, description, formula, unit, targetValue, warningThreshold, dangerThreshold)
	if execErr != nil {
		return errorResult("Failed: " + execErr.Error()), nil
	}

	msg := fmt.Sprintf("KPI defined: '%s'", name)
	if targetValue > 0 {
		msg += fmt.Sprintf(" (target: %.1f%s)", targetValue, unit)
	}
	return textResult(msg), nil
}

// PMKPISnapshot records a KPI measurement.
func (h *Handlers) PMKPISnapshot(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	kpiID := req.GetInt("kpi_id", 0)
	if kpiID == 0 {
		kpiName := req.GetString("kpi_name", "")
		if kpiName == "" {
			return errorResult("kpi_id or kpi_name required"), nil
		}
		db := h.Memory.DB()
		row := db.QueryRow("SELECT id FROM kpi_definition WHERE name = ?", kpiName)
		if err := row.Scan(&kpiID); err != nil {
			return errorResult(fmt.Sprintf("KPI '%s' not found", kpiName)), nil
		}
	}

	value := req.GetFloat("value", 0)
	sprintName := req.GetString("sprint_name", "")
	notes := req.GetString("notes", "")

	db := h.Memory.DB()
	_, execErr := db.Exec("INSERT INTO kpi_snapshot (kpi_id, value, sprint_name, notes) VALUES (?, ?, ?, ?)",
		kpiID, value, sprintName, notes)
	if execErr != nil {
		return errorResult("Failed: " + execErr.Error()), nil
	}

	// Get name + thresholds for status
	var name, unit string
	var target, warning, danger float64
	row := db.QueryRow("SELECT name, unit, target_value, warning_threshold, danger_threshold FROM kpi_definition WHERE id = ?", kpiID)
	_ = row.Scan(&name, &unit, &target, &warning, &danger)

	status := "OK"
	if danger > 0 && value <= danger {
		status = "DANGER"
	} else if warning > 0 && value <= warning {
		status = "WARNING"
	} else if target > 0 && value >= target {
		status = "ON TARGET"
	}

	return textResult(fmt.Sprintf("KPI '%s' = %.1f%s [%s]", name, value, unit, status)), nil
}

// PMKPIDashboard shows all KPIs with trends.
func (h *Handlers) PMKPIDashboard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	db := h.Memory.DB()
	rows, err := db.Query("SELECT id, name, unit, target_value, warning_threshold, danger_threshold FROM kpi_definition ORDER BY name")
	if err != nil {
		return errorResult("Query failed: " + err.Error()), nil
	}
	defer rows.Close()

	type kpiDef struct {
		id      int64
		name    string
		unit    string
		target  float64
		warning float64
		danger  float64
	}
	var kpis []kpiDef
	for rows.Next() {
		var k kpiDef
		if err := rows.Scan(&k.id, &k.name, &k.unit, &k.target, &k.warning, &k.danger); err == nil {
			kpis = append(kpis, k)
		}
	}
	if len(kpis) == 0 {
		return textResult("No KPIs defined. Use pm_kpi_define to create one."), nil
	}

	var sb strings.Builder
	sb.WriteString("KPI Dashboard\n\n")

	for _, k := range kpis {
		snapRows, _ := db.Query("SELECT value FROM kpi_snapshot WHERE kpi_id = ? ORDER BY created_at DESC LIMIT 3", k.id)
		var values []float64
		if snapRows != nil {
			for snapRows.Next() {
				var v float64
				if err := snapRows.Scan(&v); err == nil {
					values = append(values, v)
				}
			}
			snapRows.Close()
		}

		sb.WriteString(fmt.Sprintf("  %s", k.name))
		if len(values) > 0 {
			current := values[0]
			status := "OK"
			if k.danger > 0 && current <= k.danger {
				status = "DANGER"
			} else if k.warning > 0 && current <= k.warning {
				status = "WARNING"
			} else if k.target > 0 && current >= k.target {
				status = "ON TARGET"
			}
			sb.WriteString(fmt.Sprintf(": %.1f%s [%s]", current, k.unit, status))
			if len(values) >= 2 {
				if values[0] > values[1] {
					sb.WriteString(" ^")
				} else if values[0] < values[1] {
					sb.WriteString(" v")
				} else {
					sb.WriteString(" =")
				}
			}
		} else {
			sb.WriteString(": no data")
		}
		sb.WriteString("\n")
	}
	return textResult(sb.String()), nil
}

// PMGoalHitRate calculates sprint goal success rate.
func (h *Handlers) PMGoalHitRate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}
	limit := req.GetInt("limit", 10)

	goals, _ := h.Memory.GetGoalHistory(ctx, boardID, limit)
	if len(goals) == 0 {
		return textResult("No sprint goal history. Use pm_set_sprint_goal + pm_close_sprint_goal to track."), nil
	}

	achieved, partial, missed := 0, 0, 0
	var sb strings.Builder
	sb.WriteString("Sprint Goal Hit Rate\n\n")

	for _, g := range goals {
		icon := " "
		switch g.Status {
		case "achieved":
			achieved++
			icon = "+"
		case "partially_achieved":
			partial++
			icon = "~"
		case "missed":
			missed++
			icon = "-"
		case "active":
			icon = ">"
		}
		sb.WriteString(fmt.Sprintf("  [%s] %s: %s\n", icon, g.SprintName, g.Goal))
	}

	closed := achieved + partial + missed
	if closed > 0 {
		hitRate := float64(achieved) / float64(closed) * 100
		sb.WriteString(fmt.Sprintf("\nAchieved: %d/%d (%.0f%%)\n", achieved, closed, hitRate))
		if hitRate >= 70 {
			sb.WriteString("Above industry avg (52%). Strong.")
		} else if hitRate >= 52 {
			sb.WriteString("Around industry avg (52%).")
		} else {
			sb.WriteString("Below industry avg (52%). Goals too ambitious or execution gaps?")
		}
	}
	return textResult(sb.String()), nil
}

// PMOKRHealth assesses OKR risk: time elapsed vs progress.
func (h *Handlers) PMOKRHealth(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	db := h.Memory.DB()
	rows, err := db.Query(`SELECT o.id, o.title, o.cycle, o.created_at,
		AVG(CASE WHEN kr.target_value > 0 THEN kr.current_value / kr.target_value * 100 ELSE 0 END) as avg_progress
		FROM okr o
		LEFT JOIN key_result kr ON kr.okr_id = o.id AND kr.status = 'active'
		WHERE o.status = 'active'
		GROUP BY o.id`)
	if err != nil {
		return errorResult("Query failed: " + err.Error()), nil
	}
	defer rows.Close()

	var sb strings.Builder
	sb.WriteString("OKR Health Check\n\n")

	count := 0
	for rows.Next() {
		var id int64
		var title, cycle, createdAt string
		var avgProgress *float64
		if err := rows.Scan(&id, &title, &cycle, &createdAt, &avgProgress); err != nil {
			continue
		}
		count++
		prog := 0.0
		if avgProgress != nil {
			prog = *avgProgress
		}

		// Simple risk assessment
		status := "ON TRACK"
		if prog < 25 {
			status = "AT RISK"
		} else if prog < 50 {
			status = "NEEDS ATTENTION"
		}

		sb.WriteString(fmt.Sprintf("  #%d %s\n", id, title))
		sb.WriteString(fmt.Sprintf("     Progress: %.0f%% | Status: %s\n", prog, status))
		if cycle != "" {
			sb.WriteString(fmt.Sprintf("     Cycle: %s\n", cycle))
		}
	}

	if count == 0 {
		return textResult("No active OKRs to assess."), nil
	}
	return textResult(sb.String()), nil
}

// helpers

func isDoneStatus(status string) bool {
	s := strings.ToLower(status)
	return s == "done" || s == "closed" || s == "resolved" || s == "completed"
}

func progressBar(pct float64) string {
	filled := int(pct / 10)
	if filled > 10 {
		filled = 10
	}
	if filled < 0 {
		filled = 0
	}
	return "[" + strings.Repeat("#", filled) + strings.Repeat("-", 10-filled) + "]"
}
