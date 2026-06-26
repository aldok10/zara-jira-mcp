package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aldok10/zara-jira-mcp/internal/database"
	"github.com/mark3labs/mcp-go/mcp"
)

func (h *Handlers) DatabaseQuery(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Database == nil || !h.Database.Available() {
		return errorResult("No database configured. Set DATABASE_POSTGRES_DSN or DATABASE_MYSQL_DSN."), nil
	}
	query, _ := req.RequireString("query")
	results, err := h.Database.QuerySQL(ctx, req.GetString("type", ""), query, req.GetInt("limit", 50))
	if err != nil {
		return sanitizedError("Query failed", err), nil
	}
	return textResult(fmt.Sprintf("Results (%d rows):\n\n%s", len(results), database.FormatResults(results))), nil
}

func (h *Handlers) DatabaseListTables(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Database == nil || !h.Database.Available() {
		return errorResult("No database configured."), nil
	}
	dbType := req.GetString("type", "")
	if dbType == "" && h.Database.HasPostgres() {
		dbType = "postgres"
	} else if dbType == "" {
		dbType = "mysql"
	}
	tables, err := h.Database.ListTables(ctx, dbType)
	if err != nil {
		return sanitizedError("Failed", err), nil
	}
	var sb strings.Builder
	for _, t := range tables {
		sb.WriteString("- " + t + "\n")
	}
	return textResult(fmt.Sprintf("Tables (%d):\n\n%s", len(tables), sb.String())), nil
}

func (h *Handlers) MongoQuery(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Database == nil || !h.Database.HasMongo() {
		return errorResult("MongoDB not configured. Set DATABASE_MONGO_URI."), nil
	}
	collection, _ := req.RequireString("collection")
	var filter map[string]any
	_ = json.Unmarshal([]byte(req.GetString("filter", "{}")), &filter)
	if filter == nil {
		filter = map[string]any{}
	}
	results, err := h.Database.QueryMongo(ctx, collection, filter, req.GetInt("limit", 20))
	if err != nil {
		return sanitizedError("MongoDB query failed", err), nil
	}
	return textResult(fmt.Sprintf("Results (%d docs):\n\n%s", len(results), database.FormatResults(results))), nil
}

func (h *Handlers) MongoListCollections(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.Database == nil || !h.Database.HasMongo() {
		return errorResult("MongoDB not configured."), nil
	}
	collections, err := h.Database.ListCollections(ctx)
	if err != nil {
		return sanitizedError("Failed", err), nil
	}
	var sb strings.Builder
	for _, c := range collections {
		sb.WriteString("- " + c + "\n")
	}
	return textResult(fmt.Sprintf("Collections (%d):\n\n%s", len(collections), sb.String())), nil
}
