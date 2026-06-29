package domain

// Board represents a Jira board.
type Board struct {
	ID   int
	Name string
	Type string
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
