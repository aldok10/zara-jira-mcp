package domain

import "strings"

// StatusClassifier provides board-aware status classification.
// When a board configuration is available, it maps status names to
// generalized categories using the board's column layout.
// Falls back to heuristic matching when no configuration is provided.
type StatusClassifier struct {
	configs map[int]*BoardConfiguration // boardID → configuration
}

// NewStatusClassifier creates a classifier for the given board configurations.
func NewStatusClassifier(configs map[int]*BoardConfiguration) *StatusClassifier {
	if configs == nil {
		configs = make(map[int]*BoardConfiguration)
	}
	return &StatusClassifier{configs: configs}
}

// AddConfig adds or updates a board configuration.
func (sc *StatusClassifier) AddConfig(boardID int, cfg *BoardConfiguration) {
	if sc.configs == nil {
		sc.configs = make(map[int]*BoardConfiguration)
	}
	sc.configs[boardID] = cfg
}

// GetConfig returns the configuration for a board, if available.
func (sc *StatusClassifier) GetConfig(boardID int) *BoardConfiguration {
	if sc.configs == nil {
		return nil
	}
	return sc.configs[boardID]
}

// Classify returns the category for a status within a specific board context.
// Categories: "todo", "progress", "blocked", "done".
// Falls back to heuristic matching when board config is unavailable.
func (sc *StatusClassifier) Classify(boardID int, statusName string) string {
	if cfg := sc.GetConfig(boardID); cfg != nil {
		return cfg.StatusCategory(statusName)
	}
	return HeuristicClassify(statusName)
}

// ToDone returns the done statuses for a board. Useful for JQL filters.
func (sc *StatusClassifier) DoneStatuses(boardID int) []string {
	cfg := sc.GetConfig(boardID)
	if cfg == nil {
		return nil
	}
	var statuses []string
	for _, col := range cfg.ColumnConfig.Columns {
		cat := columnNameToCategory(col.Name)
		if cat == "done" {
			for _, s := range col.Statuses {
				statuses = append(statuses, s.Name)
			}
		}
	}
	return statuses
}

// HeuristicClassify is the fallback classifier used when no board config is available.
func HeuristicClassify(statusName string) string {
	lower := strings.ToLower(statusName)
	switch {
	case containsFold(lower, "done"),
		containsFold(lower, "closed"),
		containsFold(lower, "resolved"),
		containsFold(lower, "complete"):
		return "done"
	case containsFold(lower, "blocked"):
		return "blocked"
	case containsFold(lower, "progress"),
		containsFold(lower, "review"),
		containsFold(lower, "testing"),
		containsFold(lower, "dev"):
		return "progress"
	default:
		return "todo"
	}
}
