package memory

import (
	"context"
	"time"

	domain "github.com/aldok10/zara-jira-mcp/shared/domain/memory"
)

func (s *SQLiteStore) migrateAdvanced() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS dependencies (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		from_issue_key TEXT NOT NULL,
		to_issue_key TEXT NOT NULL,
		dependency_type TEXT NOT NULL DEFAULT 'blocks',
		description TEXT,
		status TEXT NOT NULL DEFAULT 'open',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		resolved_at DATETIME
	);

	CREATE TABLE IF NOT EXISTS meeting_notes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		meeting_type TEXT NOT NULL,
		date DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		attendees TEXT,
		notes TEXT,
		decisions TEXT,
		action_items TEXT,
		sprint_name TEXT
	);

	CREATE TABLE IF NOT EXISTS health_scores (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sprint_name TEXT NOT NULL,
		board_id INTEGER NOT NULL,
		computed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		overall_score INTEGER,
		velocity_score INTEGER,
		blocker_score INTEGER,
		scope_score INTEGER,
		team_score INTEGER,
		details TEXT
	);

	CREATE INDEX IF NOT EXISTS idx_deps_from ON dependencies(from_issue_key);
	CREATE INDEX IF NOT EXISTS idx_deps_to ON dependencies(to_issue_key);
	CREATE INDEX IF NOT EXISTS idx_deps_status ON dependencies(status);
	CREATE INDEX IF NOT EXISTS idx_meeting_type ON meeting_notes(meeting_type, date DESC);
	CREATE INDEX IF NOT EXISTS idx_health_board ON health_scores(board_id, computed_at DESC);
	`)
	return err
}

// Dependencies

func (s *SQLiteStore) SaveDependency(ctx context.Context, d *domain.Dependency) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO dependencies (from_issue_key, to_issue_key, dependency_type, description, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		d.FromIssueKey, d.ToIssueKey, d.DependencyType, d.Description, d.Status, d.CreatedAt)
	return err
}

func (s *SQLiteStore) ResolveDependency(ctx context.Context, id int64) error {
	now := time.Now()
	_, err := s.db.ExecContext(ctx, `UPDATE dependencies SET status='resolved', resolved_at=? WHERE id=?`, now, id)
	return err
}

func (s *SQLiteStore) GetOpenDependencies(ctx context.Context) ([]domain.Dependency, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, from_issue_key, to_issue_key, dependency_type, description, status, created_at, resolved_at
		FROM dependencies WHERE status='open' ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Dependency
	for rows.Next() {
		var d domain.Dependency
		if err := rows.Scan(&d.ID, &d.FromIssueKey, &d.ToIssueKey, &d.DependencyType, &d.Description, &d.Status, &d.CreatedAt, &d.ResolvedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) GetDependenciesForIssue(ctx context.Context, issueKey string) ([]domain.Dependency, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, from_issue_key, to_issue_key, dependency_type, description, status, created_at, resolved_at
		FROM dependencies WHERE from_issue_key=? OR to_issue_key=? ORDER BY created_at DESC`,
		issueKey, issueKey)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Dependency
	for rows.Next() {
		var d domain.Dependency
		if err := rows.Scan(&d.ID, &d.FromIssueKey, &d.ToIssueKey, &d.DependencyType, &d.Description, &d.Status, &d.CreatedAt, &d.ResolvedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

// Meeting Notes

func (s *SQLiteStore) SaveMeetingNote(ctx context.Context, m *domain.MeetingNote) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO meeting_notes (meeting_type, date, attendees, notes, decisions, action_items, sprint_name)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		m.MeetingType, m.Date, m.Attendees, m.Notes, m.Decisions, m.ActionItems, m.SprintName)
	return err
}

func (s *SQLiteStore) GetMeetingNotes(ctx context.Context, meetingType string, limit int) ([]domain.MeetingNote, error) {
	query := `SELECT id, meeting_type, date, attendees, notes, decisions, action_items, sprint_name FROM meeting_notes`
	var args []any
	if meetingType != "" {
		query += ` WHERE meeting_type=?`
		args = append(args, meetingType)
	}
	query += ` ORDER BY date DESC LIMIT ?`
	args = append(args, limit)

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.MeetingNote
	for rows.Next() {
		var m domain.MeetingNote
		if err := rows.Scan(&m.ID, &m.MeetingType, &m.Date, &m.Attendees, &m.Notes, &m.Decisions, &m.ActionItems, &m.SprintName); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

// Health Scores

func (s *SQLiteStore) SaveHealthScore(ctx context.Context, h *domain.HealthScore) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO health_scores (sprint_name, board_id, computed_at, overall_score, velocity_score, blocker_score, scope_score, team_score, details)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		h.SprintName, h.BoardID, h.ComputedAt, h.OverallScore, h.VelocityScore, h.BlockerScore, h.ScopeScore, h.TeamScore, h.Details)
	return err
}

func (s *SQLiteStore) GetHealthScores(ctx context.Context, boardID int, limit int) ([]domain.HealthScore, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, sprint_name, board_id, computed_at, overall_score, velocity_score, blocker_score, scope_score, team_score, details
		FROM health_scores WHERE board_id=? ORDER BY computed_at DESC LIMIT ?`, boardID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.HealthScore
	for rows.Next() {
		var h domain.HealthScore
		if err := rows.Scan(&h.ID, &h.SprintName, &h.BoardID, &h.ComputedAt, &h.OverallScore, &h.VelocityScore, &h.BlockerScore, &h.ScopeScore, &h.TeamScore, &h.Details); err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, rows.Err()
}
