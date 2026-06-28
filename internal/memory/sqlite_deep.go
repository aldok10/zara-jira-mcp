package memory

import (
	"context"

	domain "github.com/aldok10/zara-jira-mcp/shared/domain/memory"
)

func (s *SQLiteStore) migrateDeep() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS daily_progress (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sprint_name TEXT NOT NULL,
		board_id INTEGER NOT NULL,
		date DATE NOT NULL,
		total_issues INTEGER,
		done INTEGER,
		in_progress INTEGER,
		todo INTEGER,
		blocked INTEGER,
		points_done INTEGER DEFAULT 0,
		points_total INTEGER DEFAULT 0,
		UNIQUE(board_id, sprint_name, date)
	);

	CREATE TABLE IF NOT EXISTS sprint_goals (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sprint_name TEXT NOT NULL,
		board_id INTEGER NOT NULL,
		goal TEXT NOT NULL,
		key_results TEXT,
		status TEXT NOT NULL DEFAULT 'active',
		outcome TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		closed_at DATETIME
	);

	CREATE TABLE IF NOT EXISTS dod_items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project TEXT NOT NULL DEFAULT '*',
		item TEXT NOT NULL,
		category TEXT NOT NULL DEFAULT 'general',
		order_num INTEGER DEFAULT 0,
		active BOOLEAN NOT NULL DEFAULT 1
	);

	CREATE TABLE IF NOT EXISTS escalations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		reference_id INTEGER,
		title TEXT NOT NULL,
		severity TEXT,
		escalated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		channel TEXT NOT NULL DEFAULT 'lark',
		acknowledged BOOLEAN NOT NULL DEFAULT 0
	);

	CREATE INDEX IF NOT EXISTS idx_daily_board_sprint ON daily_progress(board_id, sprint_name, date);
	CREATE INDEX IF NOT EXISTS idx_goals_board ON sprint_goals(board_id, status);
	CREATE INDEX IF NOT EXISTS idx_dod_project ON dod_items(project, active);
	CREATE INDEX IF NOT EXISTS idx_escalations_time ON escalations(escalated_at DESC);
	`)
	return err
}

// Daily Progress

func (s *SQLiteStore) SaveDailyProgress(ctx context.Context, p *domain.DailyProgress) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT OR REPLACE INTO daily_progress (sprint_name, board_id, date, total_issues, done, in_progress, todo, blocked, points_done, points_total)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		p.SprintName, p.BoardID, p.Date, p.TotalIssues, p.Done, p.InProgress, p.Todo, p.Blocked, p.PointsDone, p.PointsTotal)
	return err
}

func (s *SQLiteStore) GetDailyProgress(ctx context.Context, boardID int, sprintName string) ([]domain.DailyProgress, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, sprint_name, board_id, date, total_issues, done, in_progress, todo, blocked, points_done, points_total
		FROM daily_progress WHERE board_id=? AND sprint_name=? ORDER BY date ASC`,
		boardID, sprintName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.DailyProgress
	for rows.Next() {
		var p domain.DailyProgress
		if err := rows.Scan(&p.ID, &p.SprintName, &p.BoardID, &p.Date, &p.TotalIssues, &p.Done, &p.InProgress, &p.Todo, &p.Blocked, &p.PointsDone, &p.PointsTotal); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// Sprint Goals

func (s *SQLiteStore) SaveSprintGoal(ctx context.Context, g *domain.SprintGoal) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO sprint_goals (sprint_name, board_id, goal, key_results, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		g.SprintName, g.BoardID, g.Goal, g.KeyResults, g.Status, g.CreatedAt)
	return err
}

func (s *SQLiteStore) UpdateSprintGoal(ctx context.Context, g *domain.SprintGoal) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE sprint_goals SET status=?, outcome=?, closed_at=? WHERE id=?`,
		g.Status, g.Outcome, g.ClosedAt, g.ID)
	return err
}

//nolint:rowserrcheck // scanGoals calls rows.Err() internally on the returned rows
func (s *SQLiteStore) GetActiveGoals(ctx context.Context, boardID int) ([]domain.SprintGoal, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, sprint_name, board_id, goal, key_results, status, outcome, created_at, closed_at
		FROM sprint_goals WHERE board_id=? AND status='active' ORDER BY created_at DESC`, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanGoals(rows)
}

//nolint:rowserrcheck // scanGoals calls rows.Err() internally on the returned rows
func (s *SQLiteStore) GetGoalHistory(ctx context.Context, boardID int, limit int) ([]domain.SprintGoal, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, sprint_name, board_id, goal, key_results, status, outcome, created_at, closed_at
		FROM sprint_goals WHERE board_id=? ORDER BY created_at DESC LIMIT ?`, boardID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanGoals(rows)
}

func scanGoals(rows interface {
	Next() bool
	Scan(...any) error
	Err() error
},
) ([]domain.SprintGoal, error) {
	var out []domain.SprintGoal
	for rows.Next() {
		var g domain.SprintGoal
		if err := rows.Scan(&g.ID, &g.SprintName, &g.BoardID, &g.Goal, &g.KeyResults, &g.Status, &g.Outcome, &g.CreatedAt, &g.ClosedAt); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

// Definition of Done

func (s *SQLiteStore) SaveDoDItem(ctx context.Context, item *domain.DoDItem) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO dod_items (project, item, category, order_num, active)
		VALUES (?, ?, ?, ?, ?)`,
		item.Project, item.Item, item.Category, item.OrderNum, item.Active)
	return err
}

func (s *SQLiteStore) GetDoD(ctx context.Context, project string) ([]domain.DoDItem, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, project, item, category, order_num, active
		FROM dod_items WHERE (project=? OR project='*') AND active=1 ORDER BY order_num ASC`,
		project)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.DoDItem
	for rows.Next() {
		var d domain.DoDItem
		if err := rows.Scan(&d.ID, &d.Project, &d.Item, &d.Category, &d.OrderNum, &d.Active); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) DeleteDoDItem(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `UPDATE dod_items SET active=0 WHERE id=?`, id)
	return err
}

// Escalations

func (s *SQLiteStore) SaveEscalation(ctx context.Context, e *domain.Escalation) error {
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO escalations (type, reference_id, title, severity, escalated_at, channel, acknowledged)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		e.Type, e.ReferenceID, e.Title, e.Severity, e.EscalatedAt, e.Channel, e.Acknowledged)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	e.ID = id
	return nil
}

func (s *SQLiteStore) GetRecentEscalations(ctx context.Context, limit int) ([]domain.Escalation, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, type, reference_id, title, severity, escalated_at, channel, acknowledged
		FROM escalations ORDER BY escalated_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Escalation
	for rows.Next() {
		var e domain.Escalation
		if err := rows.Scan(&e.ID, &e.Type, &e.ReferenceID, &e.Title, &e.Severity, &e.EscalatedAt, &e.Channel, &e.Acknowledged); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) AcknowledgeEscalation(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `UPDATE escalations SET acknowledged=1 WHERE id=?`, id)
	return err
}
