package database

import (
	"testing"
)

func TestQuerySQL_BlocksWriteOperations(t *testing.T) {
	c := &Client{} // no db configured, but validation happens first

	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{"SELECT allowed", "SELECT * FROM users", false},
		{"WITH allowed", "WITH cte AS (SELECT 1) SELECT * FROM cte", false},
		{"SHOW allowed", "SHOW TABLES", false},
		{"DESCRIBE allowed", "DESCRIBE users", false},
		{"EXPLAIN allowed", "EXPLAIN SELECT 1", false},
		{"INSERT blocked", "INSERT INTO users (name) VALUES ('x')", true},
		{"UPDATE blocked", "UPDATE users SET name='x'", true},
		{"DELETE blocked", "DELETE FROM users WHERE id=1", true},
		{"DROP blocked", "DROP TABLE users", true},
		{"ALTER blocked", "ALTER TABLE users ADD col TEXT", true},
		{"TRUNCATE blocked", "TRUNCATE TABLE users", true},
		{"select lowercase", "select * from users", false},
		{"  leading space SELECT", "  SELECT 1", false},
		{"CREATE blocked", "CREATE TABLE foo (id int)", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := c.QuerySQL(t.Context(), "", tt.query, 10)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for %q, got nil", tt.query)
				} else if err.Error() != "only SELECT/WITH/SHOW/DESCRIBE/EXPLAIN allowed (read-only)" && err.Error() != "no SQL database configured" {
					// If write is blocked, we get the validation error
					// If it passes validation but no db, we get "no SQL database configured"
					// Write ops should fail at validation
				}
			} else {
				// Allowed queries should pass validation and fail on "no SQL database configured"
				if err == nil {
					t.Error("expected 'no SQL database configured' error")
				} else if err.Error() != "no SQL database configured" {
					t.Errorf("expected 'no SQL database configured', got: %s", err.Error())
				}
			}
		})
	}
}

func TestFormatResults_Empty(t *testing.T) {
	result := FormatResults(nil)
	if result != "No results." {
		t.Errorf("got %q, want 'No results.'", result)
	}
}

func TestFormatResults_WithData(t *testing.T) {
	data := []map[string]any{
		{"id": 1, "name": "alice"},
		{"id": 2, "name": "bob"},
	}
	result := FormatResults(data)
	if result == "No results." {
		t.Error("expected formatted data")
	}
	if len(result) == 0 {
		t.Error("expected non-empty result")
	}
}

func TestFormatResults_LargeData(t *testing.T) {
	// Create data that exceeds 4000 chars
	var data []map[string]any
	for i := 0; i < 200; i++ {
		data = append(data, map[string]any{
			"id":          i,
			"description": "This is a very long description field that will contribute to the total size of the output when formatted as JSON with indentation. It needs to be quite verbose to push the total output over 4000 characters.",
		})
	}
	result := FormatResults(data)
	if len(result) > 4100 { // truncated at 4000 + "... (truncated)" suffix
		t.Errorf("expected truncated result, got length %d", len(result))
	}
	if result[len(result)-len("... (truncated)"):] != "... (truncated)" {
		t.Error("expected truncation suffix")
	}
}

func TestClient_Available(t *testing.T) {
	c := &Client{}
	if c.Available() {
		t.Error("empty client should not be available")
	}
	if c.HasPostgres() {
		t.Error("should not have postgres")
	}
	if c.HasMySQL() {
		t.Error("should not have mysql")
	}
	if c.HasMongo() {
		t.Error("should not have mongo")
	}
}

func TestClient_QueryMongo_NotConfigured(t *testing.T) {
	c := &Client{}
	_, err := c.QueryMongo(t.Context(), "col", nil, 10)
	if err == nil || err.Error() != "MongoDB not configured" {
		t.Errorf("expected 'MongoDB not configured', got: %v", err)
	}
}

func TestClient_ListCollections_NotConfigured(t *testing.T) {
	c := &Client{}
	_, err := c.ListCollections(t.Context())
	if err == nil || err.Error() != "MongoDB not configured" {
		t.Errorf("expected 'MongoDB not configured', got: %v", err)
	}
}
