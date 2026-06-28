package persistence

import (
	"context"
	"time"

	domain "github.com/aldok10/zara-jira-mcp/modules/sprint/domain/memory"
)

func (s *SQLiteStore) migrateOKR() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS okr_signals (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		objective TEXT NOT NULL,
		key_result TEXT NOT NULL,
		signal_type TEXT NOT NULL DEFAULT 'pct_done',
		jql TEXT NOT NULL,
		formula TEXT NOT NULL DEFAULT 'pct_done',
		target_value REAL DEFAULT 100,
		current_value REAL DEFAULT 0,
		progress_pct REAL DEFAULT 0,
		lark_kr_id TEXT DEFAULT '',
		last_synced DATETIME,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_okr_signals_obj ON okr_signals(objective);
	`)
	return err
}

func (s *SQLiteStore) SaveOKRSignal(ctx context.Context, sig *domain.OKRSignal) error {
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO okr_signals (objective, key_result, signal_type, jql, formula, target_value, lark_kr_id, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		sig.Objective, sig.KeyResult, sig.SignalType, sig.JQL, sig.Formula, sig.TargetValue, sig.LarkKRID, time.Now())
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	sig.ID = id
	return nil
}

func (s *SQLiteStore) UpdateOKRSignalProgress(ctx context.Context, id int64, currentValue, progressPct float64) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE okr_signals SET current_value=?, progress_pct=?, last_synced=? WHERE id=?`,
		currentValue, progressPct, time.Now(), id)
	return err
}

func (s *SQLiteStore) GetOKRSignals(ctx context.Context) ([]domain.OKRSignal, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, objective, key_result, signal_type, jql, formula, target_value, current_value, progress_pct, lark_kr_id, last_synced, created_at
		FROM okr_signals ORDER BY objective, key_result`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.OKRSignal
	for rows.Next() {
		var sig domain.OKRSignal
		if err := rows.Scan(&sig.ID, &sig.Objective, &sig.KeyResult, &sig.SignalType, &sig.JQL, &sig.Formula, &sig.TargetValue, &sig.CurrentValue, &sig.ProgressPct, &sig.LarkKRID, &sig.LastSynced, &sig.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, sig)
	}
	return out, rows.Err()
}

func (s *SQLiteStore) DeleteOKRSignal(ctx context.Context, id int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM okr_signals WHERE id=?`, id)
	return err
}
