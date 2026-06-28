package database

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver registration (init)
	_ "github.com/lib/pq"              // PostgreSQL driver registration (init)
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
)

// readOnlyRE matches SQL statements that are read-only (case-insensitive).
var readOnlyRE = regexp.MustCompile(`(?i)^\s*(SELECT|WITH|EXPLAIN|SHOW|DESCRIBE|DESC)\b`)

type Client struct {
	postgres *sql.DB
	mysql    *sql.DB
	mongo    *mongo.Client
	mongoDB  string
}

func NewClient(cfg *config.Config) *Client {
	c := &Client{}
	if cfg.Database.PostgresDSN != "" {
		if db, err := sql.Open("postgres", cfg.Database.PostgresDSN); err == nil {
			db.SetMaxOpenConns(5)
			db.SetConnMaxLifetime(5 * time.Minute)
			c.postgres = db
		}
	}
	if cfg.Database.MySQLDSN != "" {
		if db, err := sql.Open("mysql", cfg.Database.MySQLDSN); err == nil {
			db.SetMaxOpenConns(5)
			db.SetConnMaxLifetime(5 * time.Minute)
			c.mysql = db
		}
	}
	if cfg.Database.MongoURI != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Database.MongoURI)); err == nil {
			c.mongo = client
			parts := strings.Split(cfg.Database.MongoURI, "/")
			if len(parts) > 3 {
				db := parts[len(parts)-1]
				if idx := strings.Index(db, "?"); idx > 0 {
					db = db[:idx]
				}
				c.mongoDB = db
			}
		}
	}
	return c
}

func (c *Client) Available() bool   { return c.postgres != nil || c.mysql != nil || c.mongo != nil }
func (c *Client) HasPostgres() bool { return c.postgres != nil }
func (c *Client) HasMySQL() bool    { return c.mysql != nil }
func (c *Client) HasMongo() bool    { return c.mongo != nil }

func (c *Client) QuerySQL(ctx context.Context, dbType, query string, limit int) ([]map[string]any, error) {
	if limit <= 0 {
		limit = 50
	}

	// Security: validate query is read-only single statement
	q, err := sanitizeSQL(query)
	if err != nil {
		return nil, err
	}

	var db *sql.DB
	switch dbType {
	case "postgres":
		db = c.postgres
	case "mysql":
		db = c.mysql
	default:
		if c.postgres != nil {
			db = c.postgres
		} else {
			db = c.mysql
		}
	}
	if db == nil {
		return nil, fmt.Errorf("no SQL database configured")
	}
	if !strings.Contains(strings.ToUpper(q), "LIMIT") {
		q = fmt.Sprintf("%s LIMIT %d", q, limit)
	}
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("get columns: %w", err)
	}
	var results []map[string]any
	for rows.Next() {
		values := make([]any, len(columns))
		ptrs := make([]any, len(columns))
		for i := range values {
			ptrs[i] = &values[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := make(map[string]any)
		for i, col := range columns {
			if b, ok := values[i].([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = values[i]
			}
		}
		results = append(results, row)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// sanitizeSQL validates that the query is a read-only, single-statement SQL query.
// It prevents SQL injection via multi-statement and restricts to read-only operations.
func sanitizeSQL(query string) (string, error) {
	s := strings.TrimSpace(query)
	if s == "" {
		return "", fmt.Errorf("query is empty")
	}

	// Reject multi-statement: if there's a semicolon with non-whitespace after it, reject.
	if idx := strings.IndexByte(s, ';'); idx >= 0 {
		after := strings.TrimSpace(s[idx+1:])
		if after != "" {
			return "", fmt.Errorf("multi-statement queries are not allowed")
		}
		s = strings.TrimSpace(s[:idx])
	}

	// Reject if not a read-only statement
	if !readOnlyRE.MatchString(s) {
		return "", fmt.Errorf("only SELECT/WITH/SHOW/DESCRIBE/EXPLAIN allowed (read-only)")
	}

	return s, nil
}

func (c *Client) QueryMongo(ctx context.Context, collection string, filter map[string]any, limit int) ([]map[string]any, error) {
	if c.mongo == nil {
		return nil, fmt.Errorf("MongoDB not configured")
	}
	if limit <= 0 {
		limit = 50
	}
	coll := c.mongo.Database(c.mongoDB).Collection(collection)
	bsonFilter := bson.M{}
	for k, v := range filter {
		bsonFilter[k] = v
	}
	cursor, err := coll.Find(ctx, bsonFilter, options.Find().SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var results []map[string]any
	for cursor.Next(ctx) {
		var doc bson.M
		if err := cursor.Decode(&doc); err != nil {
			continue
		}
		row := make(map[string]any)
		for k, v := range doc {
			row[k] = v
		}
		results = append(results, row)
	}
	return results, nil
}

func (c *Client) ListCollections(ctx context.Context) ([]string, error) {
	if c.mongo == nil {
		return nil, fmt.Errorf("MongoDB not configured")
	}
	return c.mongo.Database(c.mongoDB).ListCollectionNames(ctx, bson.M{})
}

func (c *Client) ListTables(ctx context.Context, dbType string) ([]string, error) {
	var query string
	if dbType == "postgres" {
		query = "SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'"
	} else {
		query = "SHOW TABLES"
	}
	rows, err := c.QuerySQL(ctx, dbType, query, 200)
	if err != nil {
		return nil, err
	}
	var tables []string
	for _, row := range rows {
		for _, v := range row {
			if s, ok := v.(string); ok {
				tables = append(tables, s)
			}
		}
	}
	return tables, nil
}

func FormatResults(results []map[string]any) string {
	if len(results) == 0 {
		return "No results."
	}
	data, _ := json.MarshalIndent(results, "", "  ")
	if len(data) > 4000 {
		return string(data[:4000]) + "\n... (truncated)"
	}
	return string(data)
}
