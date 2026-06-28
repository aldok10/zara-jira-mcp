package memory

import (
	"context"

	domain "github.com/aldok10/zara-jira-mcp/shared/domain/memory"
)

func (s *SQLiteStore) migrateCoaching() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS team_pulse (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sprint_name TEXT NOT NULL,
		member TEXT NOT NULL,
		score INTEGER NOT NULL,
		notes TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS meeting_effectiveness (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		ceremony TEXT NOT NULL,
		duration_minutes INTEGER NOT NULL,
		score INTEGER NOT NULL,
		notes TEXT,
		sprint_name TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS team_radar (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sprint_name TEXT NOT NULL,
		dimension TEXT NOT NULL,
		score INTEGER NOT NULL,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_pulse_sprint ON team_pulse(sprint_name);
	CREATE INDEX IF NOT EXISTS idx_meeting_eff_ceremony ON meeting_effectiveness(ceremony, created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_radar_sprint ON team_radar(sprint_name);
	`)
	return err
}

// Team Pulse

func (s *SQLiteStore) SaveTeamPulse(ctx context.Context, p *domain.TeamPulse) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO team_pulse (sprint_name, member, score, notes) VALUES (?, ?, ?, ?)`,
		p.SprintName, p.Member, p.Score, p.Notes)
	return err
}

func (s *SQLiteStore) GetTeamPulseHistory(ctx context.Context, limit int) ([]domain.TeamPulse, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, sprint_name, member, score, notes, created_at
		FROM team_pulse ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.TeamPulse
	for rows.Next() {
		var p domain.TeamPulse
		if err := rows.Scan(&p.ID, &p.SprintName, &p.Member, &p.Score, &p.Notes, &p.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, p)
	}
	return results, rows.Err()
}

// Meeting Effectiveness

func (s *SQLiteStore) SaveMeetingEffectiveness(ctx context.Context, m *domain.MeetingEffectiveness) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO meeting_effectiveness (ceremony, duration_minutes, score, notes, sprint_name) VALUES (?, ?, ?, ?, ?)`,
		m.Ceremony, m.DurationMinutes, m.Score, m.Notes, m.SprintName)
	return err
}

func (s *SQLiteStore) GetMeetingEffectivenessHistory(ctx context.Context, ceremony string, limit int) ([]domain.MeetingEffectiveness, error) {
	var query string
	var args []any
	if ceremony != "" {
		query = `SELECT id, ceremony, duration_minutes, score, notes, sprint_name, created_at
			FROM meeting_effectiveness WHERE ceremony = ? ORDER BY created_at DESC LIMIT ?`
		args = []any{ceremony, limit}
	} else {
		query = `SELECT id, ceremony, duration_minutes, score, notes, sprint_name, created_at
			FROM meeting_effectiveness ORDER BY created_at DESC LIMIT ?`
		args = []any{limit}
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.MeetingEffectiveness
	for rows.Next() {
		var m domain.MeetingEffectiveness
		if err := rows.Scan(&m.ID, &m.Ceremony, &m.DurationMinutes, &m.Score, &m.Notes, &m.SprintName, &m.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, m)
	}
	return results, rows.Err()
}

// Team Radar

func (s *SQLiteStore) SaveTeamRadar(ctx context.Context, r *domain.TeamRadar) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO team_radar (sprint_name, dimension, score) VALUES (?, ?, ?)`,
		r.SprintName, r.Dimension, r.Score)
	return err
}

func (s *SQLiteStore) GetTeamRadarHistory(ctx context.Context, limit int) ([]domain.TeamRadar, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, sprint_name, dimension, score, created_at
		FROM team_radar ORDER BY created_at DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.TeamRadar
	for rows.Next() {
		var r domain.TeamRadar
		if err := rows.Scan(&r.ID, &r.SprintName, &r.Dimension, &r.Score, &r.CreatedAt); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}
