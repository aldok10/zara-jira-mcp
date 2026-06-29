package persistence

import (
	"context"

	domain "github.com/aldok10/zara-jira-mcp/modules/sprint/domain/memory"
)

func (s *SQLiteStore) migrateWorkflow() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS workflow_patterns (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		board_id INTEGER NOT NULL,
		status_name TEXT NOT NULL,
		classification TEXT NOT NULL,
		pattern TEXT,
		is_auto INTEGER DEFAULT 1,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(board_id, status_name)
	);
	`)
	return err
}

func (s *SQLiteStore) SaveWorkflowPattern(ctx context.Context, p *domain.WorkflowPattern) error {
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO workflow_patterns (board_id, status_name, classification, pattern, is_auto, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		p.BoardID, p.StatusName, p.Classification, p.Pattern, b2i(p.IsAuto))
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	p.ID = id
	return nil
}

func (s *SQLiteStore) UpsertWorkflowPattern(ctx context.Context, p *domain.WorkflowPattern) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO workflow_patterns (board_id, status_name, classification, pattern, is_auto, updated_at)
		VALUES (?, ?, ?, ?, ?, CURRENT_TIMESTAMP)
		ON CONFLICT(board_id, status_name) DO UPDATE SET
			classification=excluded.classification,
			pattern=excluded.pattern,
			is_auto=excluded.is_auto,
			updated_at=CURRENT_TIMESTAMP`,
		p.BoardID, p.StatusName, p.Classification, p.Pattern, b2i(p.IsAuto))
	return err
}

func (s *SQLiteStore) GetWorkflowPatterns(ctx context.Context, boardID int) ([]domain.WorkflowPattern, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, board_id, status_name, classification, pattern, is_auto, created_at, updated_at
		FROM workflow_patterns WHERE board_id=? ORDER BY status_name`, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.WorkflowPattern
	for rows.Next() {
		var p domain.WorkflowPattern
		var isAuto int
		if err := rows.Scan(&p.ID, &p.BoardID, &p.StatusName, &p.Classification, &p.Pattern, &isAuto, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.IsAuto = isAuto != 0
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) DeleteWorkflowPattern(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM workflow_patterns WHERE id=?`, id)
	return err
}

func (s *SQLiteStore) DeleteWorkflowPatternsByBoard(ctx context.Context, boardID int) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM workflow_patterns WHERE board_id=?`, boardID)
	return err
}

// b2i converts bool to int for SQLite.
func b2i(b bool) int {
	if b {
		return 1
	}
	return 0
}
