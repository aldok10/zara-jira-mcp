package domain

import "strings"

// Board represents a Jira board.
type Board struct {
	ID   int
	Name string
	Type string
}

// ColumnStatus represents a single Jira status mapped to a board column.
type ColumnStatus struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// BoardColumn represents a column on a Jira board with its mapped statuses.
type BoardColumn struct {
	Name     string          `json:"name"`
	Statuses []ColumnStatus  `json:"statuses"`
}

// ColumnConfig holds the column layout for a board.
type ColumnConfig struct {
	Columns []BoardColumn `json:"columns"`
}

// BoardConfiguration holds the full board configuration including column layout
// and status mappings. Used by StatusClassifier to understand board-specific statuses.
type BoardConfiguration struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Type         string        `json:"type"`
	ColumnConfig ColumnConfig  `json:"columnConfig"`
}

// StatusCategory returns the board-aware classification for a status name.
// Maps a status string to: "todo", "progress", "done", or "blocked".
func (bc *BoardConfiguration) StatusCategory(statusName string) string {
	if statusName == "" {
		return "todo"
	}

	for _, col := range bc.ColumnConfig.Columns {
		for _, s := range col.Statuses {
			if strings.EqualFold(s.Name, statusName) {
				return columnNameToCategory(col.Name)
			}
		}
	}
	return "todo"
}

// AllStatusNames returns all status names defined in this board configuration.
func (bc *BoardConfiguration) AllStatusNames() []string {
	var names []string
	for _, col := range bc.ColumnConfig.Columns {
		for _, s := range col.Statuses {
			names = append(names, s.Name)
		}
	}
	return names
}

// columnNameToCategory maps a Jira board column name to a generalized category.
// Uses heuristic matching on column names.
func columnNameToCategory(colName string) string {
	lower := strings.ToLower(colName)
	switch {
	case containsFold(lower, "done"),
		containsFold(lower, "closed"),
		containsFold(lower, "resolved"),
		containsFold(lower, "complete"),
		containsFold(lower, "released"):
		return "done"
	case containsFold(lower, "blocked"),
		containsFold(lower, "stalled"),
		containsFold(lower, "waiting"),
		containsFold(lower, "impediment"):
		return "blocked"
	case containsFold(lower, "in progress"),
		containsFold(lower, "review"),
		containsFold(lower, "testing"),
		containsFold(lower, "dev"),
		containsFold(lower, "selected for development"):
		return "progress"
	default:
		return "todo"
	}
}

// containsFold returns true if s contains substr (case-insensitive).
func containsFold(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// Sprint represents a Jira sprint.
type Sprint struct {
	ID        int
	Name      string
	State     string
	Goal      string
	StartDate string
	EndDate   string
}

// SearchResult holds paginated search results.
type SearchResult struct {
	Issues     []Issue
	Total      int
	StartAt    int
	MaxResults int
	HasMore    bool
}
