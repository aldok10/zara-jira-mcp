package memory

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	domain "github.com/aldok10/zara-jira-mcp/domain/memory"
)

// SQLiteStore implements domain.Store using SQLite.
type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	s := &SQLiteStore{db: db}
	if err := s.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}
	return s, nil
}

func (s *SQLiteStore) Close() error { return s.db.Close() }

func (s *SQLiteStore) migrate() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS sprint_snapshots (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sprint_name TEXT NOT NULL,
		board_id INTEGER NOT NULL,
		snapshot_date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		total_issues INTEGER,
		done INTEGER,
		in_progress INTEGER,
		todo INTEGER,
		blocked INTEGER,
		carryover INTEGER,
		velocity INTEGER,
		completion_rate REAL,
		notes TEXT
	);

	CREATE TABLE IF NOT EXISTS risks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		severity TEXT NOT NULL DEFAULT 'medium',
		status TEXT NOT NULL DEFAULT 'open',
		owner TEXT,
		mitigation TEXT,
		identified_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		resolved_at DATETIME,
		sprint_name TEXT
	);

	CREATE TABLE IF NOT EXISTS decisions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		context TEXT,
		decision TEXT NOT NULL,
		rationale TEXT,
		outcome TEXT,
		made_by TEXT,
		made_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		tags TEXT
	);

	CREATE TABLE IF NOT EXISTS blockers (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		issue_key TEXT,
		description TEXT NOT NULL,
		blocked_since DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		resolved_at DATETIME,
		resolution TEXT,
		owner TEXT,
		days_blocked INTEGER DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS team_metrics (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		member_name TEXT NOT NULL,
		sprint_name TEXT NOT NULL,
		recorded_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		issues_assigned INTEGER,
		issues_done INTEGER,
		blocker_count INTEGER,
		carryover_count INTEGER,
		notes TEXT
	);

	CREATE TABLE IF NOT EXISTS retrospectives (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sprint_name TEXT NOT NULL,
		date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		went_well TEXT,
		improvements TEXT,
		action_items TEXT,
		status TEXT NOT NULL DEFAULT 'open'
	);

	CREATE TABLE IF NOT EXISTS action_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		retro_id INTEGER REFERENCES retrospectives(id),
		description TEXT NOT NULL,
		owner TEXT,
		due_date DATETIME,
		status TEXT NOT NULL DEFAULT 'pending'
	);

	CREATE INDEX IF NOT EXISTS idx_snapshots_board ON sprint_snapshots(board_id, snapshot_date DESC);
	CREATE INDEX IF NOT EXISTS idx_risks_status ON risks(status);
	CREATE INDEX IF NOT EXISTS idx_blockers_resolved ON blockers(resolved_at);
	CREATE INDEX IF NOT EXISTS idx_team_metrics_member ON team_metrics(member_name, recorded_at DESC);
	CREATE INDEX IF NOT EXISTS idx_action_items_status ON action_items(status);
	`)
	return err
}

// Sprint Snapshots

func (s *SQLiteStore) SaveSprintSnapshot(ctx context.Context, snap *domain.SprintSnapshot) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO sprint_snapshots (sprint_name, board_id, snapshot_date, total_issues, done, in_progress, todo, blocked, carryover, velocity, completion_rate, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		snap.SprintName, snap.BoardID, snap.SnapshotDate, snap.TotalIssues, snap.Done, snap.InProgress, snap.Todo, snap.Blocked, snap.Carryover, snap.Velocity, snap.CompletionRate, snap.Notes)
	return err
}

func (s *SQLiteStore) GetSprintSnapshots(ctx context.Context, boardID int, limit int) ([]domain.SprintSnapshot, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, sprint_name, board_id, snapshot_date, total_issues, done, in_progress, todo, blocked, carryover, velocity, completion_rate, notes
		FROM sprint_snapshots WHERE board_id = ? ORDER BY snapshot_date DESC LIMIT ?`, boardID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSnapshots(rows)
}

func (s *SQLiteStore) GetLatestSnapshot(ctx context.Context, boardID int) (*domain.SprintSnapshot, error) {
	snaps, err := s.GetSprintSnapshots(ctx, boardID, 1)
	if err != nil {
		return nil, err
	}
	if len(snaps) == 0 {
		return nil, nil
	}
	return &snaps[0], nil
}

// Risks

func (s *SQLiteStore) SaveRisk(ctx context.Context, r *domain.Risk) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO risks (title, description, severity, status, owner, mitigation, identified_at, sprint_name)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		r.Title, r.Description, r.Severity, r.Status, r.Owner, r.Mitigation, r.IdentifiedAt, r.SprintName)
	return err
}

func (s *SQLiteStore) UpdateRisk(ctx context.Context, r *domain.Risk) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE risks SET title=?, description=?, severity=?, status=?, owner=?, mitigation=?, resolved_at=?, sprint_name=?
		WHERE id=?`,
		r.Title, r.Description, r.Severity, r.Status, r.Owner, r.Mitigation, r.ResolvedAt, r.SprintName, r.ID)
	return err
}

func (s *SQLiteStore) GetOpenRisks(ctx context.Context) ([]domain.Risk, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, description, severity, status, owner, mitigation, identified_at, resolved_at, sprint_name
		FROM risks WHERE status IN ('open', 'mitigating') ORDER BY
		CASE severity WHEN 'critical' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 ELSE 4 END`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRisks(rows)
}

func (s *SQLiteStore) GetAllRisks(ctx context.Context, limit int) ([]domain.Risk, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, description, severity, status, owner, mitigation, identified_at, resolved_at, sprint_name
		FROM risks ORDER BY identified_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRisks(rows)
}

// Decisions

func (s *SQLiteStore) SaveDecision(ctx context.Context, d *domain.Decision) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO decisions (title, context, decision, rationale, outcome, made_by, made_at, tags)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		d.Title, d.Context, d.Decision, d.Rationale, d.Outcome, d.MadeBy, d.MadeAt, d.Tags)
	return err
}

func (s *SQLiteStore) GetDecisions(ctx context.Context, limit int) ([]domain.Decision, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, context, decision, rationale, outcome, made_by, made_at, tags
		FROM decisions ORDER BY made_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDecisions(rows)
}

func (s *SQLiteStore) SearchDecisions(ctx context.Context, query string) ([]domain.Decision, error) {
	like := "%" + query + "%"
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, title, context, decision, rationale, outcome, made_by, made_at, tags
		FROM decisions WHERE title LIKE ? OR decision LIKE ? OR tags LIKE ? ORDER BY made_at DESC LIMIT 20`,
		like, like, like)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDecisions(rows)
}

// Blockers

func (s *SQLiteStore) SaveBlocker(ctx context.Context, b *domain.Blocker) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO blockers (issue_key, description, blocked_since, owner)
		VALUES (?, ?, ?, ?)`,
		b.IssueKey, b.Description, b.BlockedSince, b.Owner)
	return err
}

func (s *SQLiteStore) ResolveBlocker(ctx context.Context, id int64, resolution string) error {
	now := time.Now()
	_, err := s.db.ExecContext(ctx, `
		UPDATE blockers SET resolved_at=?, resolution=?,
		days_blocked = CAST((julianday(?) - julianday(blocked_since)) AS INTEGER)
		WHERE id=?`, now, resolution, now, id)
	return err
}

func (s *SQLiteStore) GetActiveBlockers(ctx context.Context) ([]domain.Blocker, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, issue_key, description, blocked_since, resolved_at, resolution, owner, days_blocked
		FROM blockers WHERE resolved_at IS NULL ORDER BY blocked_since ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBlockers(rows)
}

func (s *SQLiteStore) GetBlockerHistory(ctx context.Context, limit int) ([]domain.Blocker, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, issue_key, description, blocked_since, resolved_at, resolution, owner, days_blocked
		FROM blockers ORDER BY blocked_since DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanBlockers(rows)
}

// Team Metrics

func (s *SQLiteStore) SaveTeamMetric(ctx context.Context, m *domain.TeamMetric) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO team_metrics (member_name, sprint_name, recorded_at, issues_assigned, issues_done, blocker_count, carryover_count, notes)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		m.MemberName, m.SprintName, m.RecordedAt, m.IssuesAssigned, m.IssuesDone, m.BlockerCount, m.CarryoverCount, m.Notes)
	return err
}

func (s *SQLiteStore) GetTeamMetrics(ctx context.Context, memberName string, limit int) ([]domain.TeamMetric, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, member_name, sprint_name, recorded_at, issues_assigned, issues_done, blocker_count, carryover_count, notes
		FROM team_metrics WHERE member_name = ? ORDER BY recorded_at DESC LIMIT ?`, memberName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTeamMetrics(rows)
}

func (s *SQLiteStore) GetTeamOverview(ctx context.Context, sprintName string) ([]domain.TeamMetric, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, member_name, sprint_name, recorded_at, issues_assigned, issues_done, blocker_count, carryover_count, notes
		FROM team_metrics WHERE sprint_name = ? ORDER BY member_name`, sprintName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTeamMetrics(rows)
}

// Retrospectives

func (s *SQLiteStore) SaveRetrospective(ctx context.Context, r *domain.Retrospective) error {
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO retrospectives (sprint_name, date, went_well, improvements, action_items, status)
		VALUES (?, ?, ?, ?, ?, ?)`,
		r.SprintName, r.Date, r.WentWell, r.Improvements, r.ActionItems, r.Status)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	r.ID = id
	return nil
}

func (s *SQLiteStore) GetRetrospectives(ctx context.Context, limit int) ([]domain.Retrospective, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, sprint_name, date, went_well, improvements, action_items, status
		FROM retrospectives ORDER BY date DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Retrospective
	for rows.Next() {
		var r domain.Retrospective
		if err := rows.Scan(&r.ID, &r.SprintName, &r.Date, &r.WentWell, &r.Improvements, &r.ActionItems, &r.Status); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) SaveActionItem(ctx context.Context, a *domain.ActionItem) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO action_items (retro_id, description, owner, due_date, status)
		VALUES (?, ?, ?, ?, ?)`,
		a.RetroID, a.Description, a.Owner, a.DueDate, a.Status)
	return err
}

func (s *SQLiteStore) GetPendingActionItems(ctx context.Context) ([]domain.ActionItem, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, retro_id, description, owner, due_date, status
		FROM action_items WHERE status = 'pending' ORDER BY due_date ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.ActionItem
	for rows.Next() {
		var a domain.ActionItem
		if err := rows.Scan(&a.ID, &a.RetroID, &a.Description, &a.Owner, &a.DueDate, &a.Status); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) CompleteActionItem(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `UPDATE action_items SET status = 'done' WHERE id = ?`, id)
	return err
}

// Scanners

func scanSnapshots(rows *sql.Rows) ([]domain.SprintSnapshot, error) {
	var out []domain.SprintSnapshot
	for rows.Next() {
		var s domain.SprintSnapshot
		if err := rows.Scan(&s.ID, &s.SprintName, &s.BoardID, &s.SnapshotDate, &s.TotalIssues, &s.Done, &s.InProgress, &s.Todo, &s.Blocked, &s.Carryover, &s.Velocity, &s.CompletionRate, &s.Notes); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func scanRisks(rows *sql.Rows) ([]domain.Risk, error) {
	var out []domain.Risk
	for rows.Next() {
		var r domain.Risk
		if err := rows.Scan(&r.ID, &r.Title, &r.Description, &r.Severity, &r.Status, &r.Owner, &r.Mitigation, &r.IdentifiedAt, &r.ResolvedAt, &r.SprintName); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

func scanDecisions(rows *sql.Rows) ([]domain.Decision, error) {
	var out []domain.Decision
	for rows.Next() {
		var d domain.Decision
		if err := rows.Scan(&d.ID, &d.Title, &d.Context, &d.Decision, &d.Rationale, &d.Outcome, &d.MadeBy, &d.MadeAt, &d.Tags); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func scanBlockers(rows *sql.Rows) ([]domain.Blocker, error) {
	var out []domain.Blocker
	for rows.Next() {
		var b domain.Blocker
		if err := rows.Scan(&b.ID, &b.IssueKey, &b.Description, &b.BlockedSince, &b.ResolvedAt, &b.Resolution, &b.Owner, &b.DaysBlocked); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

func scanTeamMetrics(rows *sql.Rows) ([]domain.TeamMetric, error) {
	var out []domain.TeamMetric
	for rows.Next() {
		var m domain.TeamMetric
		if err := rows.Scan(&m.ID, &m.MemberName, &m.SprintName, &m.RecordedAt, &m.IssuesAssigned, &m.IssuesDone, &m.BlockerCount, &m.CarryoverCount, &m.Notes); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// Ensure interface compliance.
var _ domain.Store = (*SQLiteStore)(nil)
