package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handlers) initOKRTables() {
	db := h.Memory.DB()
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS okr (
		id INTEGER PRIMARY KEY,
		level TEXT NOT NULL,
		title TEXT NOT NULL,
		description TEXT,
		owner TEXT,
		cycle TEXT,
		parent_id INTEGER,
		status TEXT DEFAULT 'active',
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS key_result (
		id INTEGER PRIMARY KEY,
		okr_id INTEGER NOT NULL,
		title TEXT NOT NULL,
		metric TEXT,
		start_value REAL DEFAULT 0,
		target_value REAL NOT NULL,
		current_value REAL DEFAULT 0,
		unit TEXT,
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
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS hypothesis (
		id INTEGER PRIMARY KEY,
		kr_id INTEGER,
		statement TEXT NOT NULL,
		metric TEXT,
		baseline TEXT,
		target TEXT,
		sprint_name TEXT,
		status TEXT DEFAULT 'open',
		result TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	_, _ = db.Exec(`CREATE TABLE IF NOT EXISTS kpi_definition (
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		formula TEXT,
		unit TEXT,
		target_value REAL,
		warning_threshold REAL,
		danger_threshold REAL,
		source TEXT DEFAULT 'jira',
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
		return errorResult("Failed to save OKR: " + execErr.Error()), nil
	}

	// Extract OKR ID from result
	type lastIDer interface{ LastInsertId() (int64, error) }
	var okrID int64
	if r, ok := result.(lastIDer); ok {
		okrID, _ = r.LastInsertId()
	}

	// Parse key results (newline-separated format: "KR text | target_value | unit")
	var krCount int
	if keyResults != "" && okrID > 0 {
		lines := strings.Split(keyResults, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.Split(line, "|")
			krTitle := strings.TrimSpace(parts[0])
			var targetVal float64 = 100
			var unit string
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

	msg := fmt.Sprintf("OKR created: '%s' (level: %s)", title, level)
	if krCount > 0 {
		msg += fmt.Sprintf("\n%d Key Results defined.", krCount)
	}
	if cycle != "" {
		msg += fmt.Sprintf("\nCycle: %s", cycle)
	}
	return textResult(msg), nil
}

// PMKRLink connects Jira issues/epics to a Key Result.
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

	// Verify KR exists
	row := db.QueryRow("SELECT title FROM key_result WHERE id = ?", krID)
	var krTitle string
	if scanErr := row.Scan(&krTitle); scanErr != nil {
		return errorResult(fmt.Sprintf("Key Result #%d not found", krID)), nil
	}

	keys := strings.Split(issueKeys, ",")
	linked := 0
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		_, _ = db.Exec("INSERT INTO kr_jira_link (kr_id, issue_key, link_type) VALUES (?, ?, ?)",
			krID, key, linkType)
		linked++
	}

	return textResult(fmt.Sprintf("Linked %d issues to KR '%s' (type: %s).", linked, krTitle, linkType)), nil
}

// PMKRProgress calculates Key Result progress from linked Jira issues.
func (h *Handlers) PMKRProgress(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	db := h.Memory.DB()

	// Get all OKRs with their KRs
	okrRows, err := db.Query(`SELECT o.id, o.title, o.level, o.cycle, o.status 
		FROM okr o WHERE o.status = 'active' ORDER BY o.created_at DESC`)
	if err != nil {
		return errorResult("Query failed: " + err.Error()), nil
	}
	defer okrRows.Close()

	type okrEntry struct {
		id     int64
		title  string
		level  string
		cycle  string
		status string
	}
	var okrs []okrEntry
	for okrRows.Next() {
		var o okrEntry
		if scanErr := okrRows.Scan(&o.id, &o.title, &o.level, &o.cycle, &o.status); scanErr != nil {
			continue
		}
		okrs = append(okrs, o)
	}

	if len(okrs) == 0 {
		return textResult("No active OKRs. Use pm_okr_define to create one."), nil
	}

	var sb strings.Builder
	sb.WriteString("OKR Progress Report\n")
	sb.WriteString(strings.Repeat("=", 40))
	sb.WriteString("\n\n")

	for _, o := range okrs {
		sb.WriteString(fmt.Sprintf("[%s] %s", strings.ToUpper(o.level), o.title))
		if o.cycle != "" {
			sb.WriteString(fmt.Sprintf(" (%s)", o.cycle))
		}
		sb.WriteString("\n")

		// Get KRs for this OKR
		krRows, krErr := db.Query(`SELECT kr.id, kr.title, kr.start_value, kr.target_value, kr.current_value, kr.unit
			FROM key_result kr WHERE kr.okr_id = ? AND kr.status = 'active'`, o.id)
		if krErr != nil {
			continue
		}

		krCount := 0
		totalProgress := 0.0
		for krRows.Next() {
			var krID int64
			var krTitle, unit string
			var startVal, targetVal, currentVal float64
			if scanErr := krRows.Scan(&krID, &krTitle, &startVal, &targetVal, &currentVal, &unit); scanErr != nil {
				continue
			}

			// Calculate progress from linked Jira issues
			linkRows, linkErr := db.Query("SELECT issue_key FROM kr_jira_link WHERE kr_id = ?", krID)
			if linkErr == nil {
				var issueKeys []string
				for linkRows.Next() {
					var key string
					if scanErr := linkRows.Scan(&key); scanErr == nil {
						issueKeys = append(issueKeys, key)
					}
				}
				linkRows.Close()

				// Query Jira for issue statuses
				if len(issueKeys) > 0 && h.Jira != nil {
					doneCount := 0
					totalIssues := len(issueKeys)
					for _, key := range issueKeys {
						issue, issErr := h.Jira.GetIssue(ctx, key)
						if issErr == nil && isDoneStatus(issue.Status) {
							doneCount++
						}
					}
					if totalIssues > 0 {
						jiraProgress := float64(doneCount) / float64(totalIssues) * targetVal
						currentVal = jiraProgress
						// Update current_value in DB
						_, _ = db.Exec("UPDATE key_result SET current_value = ? WHERE id = ?", currentVal, krID)
					}
				}
			}

			// Calculate percentage
			var progress float64
			if targetVal-startVal != 0 {
				progress = (currentVal - startVal) / (targetVal - startVal) * 100
			}
			if progress > 100 {
				progress = 100
			}
			if progress < 0 {
				progress = 0
			}
			totalProgress += progress
			krCount++

			// Progress bar
			bar := progressBar(progress)
			sb.WriteString(fmt.Sprintf("  KR%d: %s\n", krCount, krTitle))
			sb.WriteString(fmt.Sprintf("       %s %.0f%% (%.0f/%.0f %s)\n", bar, progress, currentVal, targetVal, unit))
		}
		krRows.Close()

		if krCount > 0 {
			avgProgress := totalProgress / float64(krCount)
			sb.WriteString(fmt.Sprintf("  Overall: %.0f%%\n", avgProgress))
		} else {
			sb.WriteString("  No Key Results defined.\n")
		}
		sb.WriteString("\n")
	}

	return textResult(sb.String()), nil
}

// PMHypothesis records a hypothesis for validation.
func (h *Handlers) PMHypothesis(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	statement, err := req.RequireString("statement")
	if err != nil {
		return errorResult("statement required (format: 'We believe [X] will result in [Y] as measured by [Z]')"), nil
	}
	metric := req.GetString("metric", "")
	baseline := req.GetString("baseline", "")
	target := req.GetString("target", "")
	sprintName := req.GetString("sprint_name", "")
	krID := req.GetString("kr_id", "")

	db := h.Memory.DB()

	var krIDVal *int64
	if krID != "" {
		var v int64
		fmt.Sscanf(krID, "%d", &v)
		if v > 0 {
			krIDVal = &v
		}
	}

	_, execErr := db.Exec(`INSERT INTO hypothesis (kr_id, statement, metric, baseline, target, sprint_name) 
		VALUES (?, ?, ?, ?, ?, ?)`, krIDVal, statement, metric, baseline, target, sprintName)
	if execErr != nil {
		return errorResult("Failed to save: " + execErr.Error()), nil
	}

	msg := fmt.Sprintf("Hypothesis recorded: '%s'", statement)
	if metric != "" {
		msg += fmt.Sprintf("\nMetric: %s (baseline: %s, target: %s)", metric, baseline, target)
	}
	return textResult(msg), nil
}

// PMHypothesisValidate validates or invalidates a hypothesis.
func (h *Handlers) PMHypothesisValidate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	id, err := req.RequireInt("id")
	if err != nil {
		return errorResult("id required"), nil
	}
	status, err := req.RequireString("status")
	if err != nil {
		return errorResult("status required (validated/invalidated/inconclusive)"), nil
	}
	result := req.GetString("result", "")

	db := h.Memory.DB()
	_, execErr := db.Exec("UPDATE hypothesis SET status = ?, result = ? WHERE id = ?", status, result, id)
	if execErr != nil {
		return errorResult("Failed to update: " + execErr.Error()), nil
	}

	return textResult(fmt.Sprintf("Hypothesis #%d marked as '%s'. %s", id, status, result)), nil
}

// PMOutcomeReview reviews if sprint work moved the OKR needle.
func (h *Handlers) PMOutcomeReview(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	boardID, err := req.RequireInt("board_id")
	if err != nil {
		return errorResult("board_id required"), nil
	}

	db := h.Memory.DB()
	var sb strings.Builder
	sb.WriteString("Outcome Review: Did Sprint Work Move the OKR Needle?\n")
	sb.WriteString(strings.Repeat("=", 50))
	sb.WriteString("\n\n")

	// Get latest outcome_map entries for this board
	rows, queryErr := db.Query(`SELECT sprint_name, objective, key_results FROM outcome_map 
		WHERE board_id = ? ORDER BY created_at DESC LIMIT 5`, boardID)
	if queryErr != nil {
		return errorResult("Query failed: " + queryErr.Error()), nil
	}
	defer rows.Close()

	type mapping struct {
		sprint    string
		objective string
		krs       string
	}
	var mappings []mapping
	for rows.Next() {
		var m mapping
		if scanErr := rows.Scan(&m.sprint, &m.objective, &m.krs); scanErr != nil {
			continue
		}
		mappings = append(mappings, m)
	}

	// Get active OKRs with progress
	okrRows, okrErr := db.Query(`SELECT o.title, o.level,
		(SELECT AVG(CASE WHEN kr.target_value - kr.start_value > 0 
			THEN (kr.current_value - kr.start_value) / (kr.target_value - kr.start_value) * 100 
			ELSE 0 END)
		 FROM key_result kr WHERE kr.okr_id = o.id AND kr.status = 'active') as avg_progress
		FROM okr o WHERE o.status = 'active'`)
	if okrErr == nil {
		defer okrRows.Close()
		sb.WriteString("ACTIVE OKR PROGRESS:\n")
		for okrRows.Next() {
			var title, level string
			var avgProgress *float64
			if scanErr := okrRows.Scan(&title, &level, &avgProgress); scanErr != nil {
				continue
			}
			prog := 0.0
			if avgProgress != nil {
				prog = *avgProgress
			}
			sb.WriteString(fmt.Sprintf("  [%s] %s — %.0f%%\n", strings.ToUpper(level), title, prog))
		}
		sb.WriteString("\n")
	}

	// Show sprint-to-objective mappings
	if len(mappings) > 0 {
		sb.WriteString("SPRINT → OBJECTIVE ALIGNMENT:\n")
		for _, m := range mappings {
			sb.WriteString(fmt.Sprintf("  %s → %s\n", m.sprint, m.objective))
		}
		sb.WriteString("\n")
	}

	// Hypotheses status
	hypRows, hypErr := db.Query(`SELECT statement, status, result, sprint_name 
		FROM hypothesis ORDER BY created_at DESC LIMIT 10`)
	if hypErr == nil {
		defer hypRows.Close()
		hasHyp := false
		for hypRows.Next() {
			if !hasHyp {
				sb.WriteString("HYPOTHESES:\n")
				hasHyp = true
			}
			var stmt, status, result, sprint string
			if scanErr := hypRows.Scan(&stmt, &status, &result, &sprint); scanErr != nil {
				continue
			}
			icon := "?"
			switch status {
			case "validated":
				icon = "+"
			case "invalidated":
				icon = "-"
			case "inconclusive":
				icon = "~"
			}
			sb.WriteString(fmt.Sprintf("  [%s] %s", icon, stmt))
			if sprint != "" {
				sb.WriteString(fmt.Sprintf(" (Sprint: %s)", sprint))
			}
			if result != "" {
				sb.WriteString(fmt.Sprintf("\n      Result: %s", result))
			}
			sb.WriteString("\n")
		}
		if hasHyp {
			sb.WriteString("\n")
		}
	}

	// AI summary if available
	if h.AI != nil {
		prompt := `Based on this OKR/outcome data, provide a brief 2-3 sentence assessment:
1. Are sprints clearly connected to business objectives?
2. Are Key Results actually moving?
3. What's the biggest gap between effort and outcome?
Be direct and actionable.`
		aiResult, aiErr := h.AI.Complete(ctx, prompt, sb.String())
		if aiErr == nil {
			sb.WriteString("AI ASSESSMENT:\n")
			sb.WriteString(aiResult)
			sb.WriteString("\n")
		}
	}

	return textResult(sb.String()), nil
}

// PMKPIDefine creates a KPI definition for tracking.
func (h *Handlers) PMKPIDefine(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	name, err := req.RequireString("name")
	if err != nil {
		return errorResult("name required"), nil
	}
	description := req.GetString("description", "")
	formula := req.GetString("formula", "")
	unit := req.GetString("unit", "")
	targetValue := req.GetFloat("target_value", 0)
	warningThreshold := req.GetFloat("warning_threshold", 0)
	dangerThreshold := req.GetFloat("danger_threshold", 0)

	db := h.Memory.DB()
	_, execErr := db.Exec(`INSERT INTO kpi_definition (name, description, formula, unit, target_value, warning_threshold, danger_threshold)
		VALUES (?, ?, ?, ?, ?, ?, ?)`, name, description, formula, unit, targetValue, warningThreshold, dangerThreshold)
	if execErr != nil {
		return errorResult("Failed to save: " + execErr.Error()), nil
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

	kpiID, err := req.RequireInt("kpi_id")
	if err != nil {
		// Try by name
		kpiName := req.GetString("kpi_name", "")
		if kpiName == "" {
			return errorResult("kpi_id or kpi_name required"), nil
		}
		db := h.Memory.DB()
		row := db.QueryRow("SELECT id FROM kpi_definition WHERE name = ?", kpiName)
		if scanErr := row.Scan(&kpiID); scanErr != nil {
			return errorResult(fmt.Sprintf("KPI '%s' not found", kpiName)), nil
		}
	}

	value, valErr := req.RequireFloat("value")
	if valErr != nil {
		return errorResult("value required"), nil
	}
	sprintName := req.GetString("sprint_name", "")
	notes := req.GetString("notes", "")

	db := h.Memory.DB()
	_, execErr := db.Exec("INSERT INTO kpi_snapshot (kpi_id, value, sprint_name, notes) VALUES (?, ?, ?, ?)",
		kpiID, value, sprintName, notes)
	if execErr != nil {
		return errorResult("Failed to save: " + execErr.Error()), nil
	}

	// Get KPI name and thresholds for status
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

// PMKPIDashboard shows all KPIs with latest values and trends.
func (h *Handlers) PMKPIDashboard(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	db := h.Memory.DB()
	rows, err := db.Query(`SELECT id, name, unit, target_value, warning_threshold, danger_threshold FROM kpi_definition ORDER BY name`)
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
		if scanErr := rows.Scan(&k.id, &k.name, &k.unit, &k.target, &k.warning, &k.danger); scanErr != nil {
			continue
		}
		kpis = append(kpis, k)
	}

	if len(kpis) == 0 {
		return textResult("No KPIs defined. Use pm_kpi_define to create one."), nil
	}

	var sb strings.Builder
	sb.WriteString("KPI Dashboard\n")
	sb.WriteString(strings.Repeat("=", 40))
	sb.WriteString("\n\n")

	for _, k := range kpis {
		// Get last 3 snapshots for trend
		snapRows, snapErr := db.Query(`SELECT value, sprint_name, created_at FROM kpi_snapshot 
			WHERE kpi_id = ? ORDER BY created_at DESC LIMIT 3`, k.id)
		if snapErr != nil {
			continue
		}

		var values []float64
		var latestSprint string
		for snapRows.Next() {
			var val float64
			var sprint, created string
			if scanErr := snapRows.Scan(&val, &sprint, &created); scanErr == nil {
				values = append(values, val)
				if latestSprint == "" {
					latestSprint = sprint
				}
			}
		}
		snapRows.Close()

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

			// Trend arrow
			if len(values) >= 2 {
				if values[0] > values[1] {
					sb.WriteString(" ^")
				} else if values[0] < values[1] {
					sb.WriteString(" v")
				} else {
					sb.WriteString(" =")
				}
			}
			if k.target > 0 {
				sb.WriteString(fmt.Sprintf(" (target: %.1f%s)", k.target, k.unit))
			}
		} else {
			sb.WriteString(": no data yet")
		}
		sb.WriteString("\n")
	}

	return textResult(sb.String()), nil
}

// PMOKRList shows all OKRs with hierarchy.
func (h *Handlers) PMOKRList(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	h.initOKRTables()

	status := req.GetString("status", "active")
	db := h.Memory.DB()

	rows, err := db.Query(`SELECT id, level, title, owner, cycle, status FROM okr WHERE status = ? ORDER BY level, created_at`, status)
	if err != nil {
		return errorResult("Query failed: " + err.Error()), nil
	}
	defer rows.Close()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("OKRs (status: %s)\n\n", status))

	count := 0
	for rows.Next() {
		var id int64
		var level, title, owner, cycle, st string
		if scanErr := rows.Scan(&id, &level, &title, &owner, &cycle, &st); scanErr != nil {
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

		// Show KRs
		krRows, _ := db.Query(`SELECT id, title, current_value, target_value, unit FROM key_result WHERE okr_id = ? AND status = 'active'`, id)
		if krRows != nil {
			for krRows.Next() {
				var krID int64
				var krTitle, unit string
				var current, target float64
				if scanErr := krRows.Scan(&krID, &krTitle, &current, &target, &unit); scanErr == nil {
					pct := 0.0
					if target > 0 {
						pct = current / target * 100
					}
					sb.WriteString(fmt.Sprintf("  KR#%d: %s (%.0f/%.0f%s = %.0f%%)\n", krID, krTitle, current, target, unit, pct))
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
