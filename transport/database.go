package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerDatabaseTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("db_query",
		mcp.WithDescription("Execute a read-only SQL query (SELECT only) against Postgres or MySQL."),
		mcp.WithString("query", mcp.Required(), mcp.Description("SQL query (SELECT/WITH/SHOW only)")),
		mcp.WithString("type", mcp.Description("postgres or mysql (auto-detect if empty)")),
		mcp.WithNumber("limit", mcp.Description("Max rows (default: 50)")),
	), h.DatabaseQuery)

	s.AddTool(mcp.NewTool("db_tables",
		mcp.WithDescription("List all tables in the configured SQL database."),
		mcp.WithString("type", mcp.Description("postgres or mysql")),
	), h.DatabaseListTables)

	s.AddTool(mcp.NewTool("mongo_query",
		mcp.WithDescription("Query a MongoDB collection with optional filter. Read-only."),
		mcp.WithString("collection", mcp.Required(), mcp.Description("Collection name")),
		mcp.WithString("filter", mcp.Description("MongoDB filter as JSON (default: {})")),
		mcp.WithNumber("limit", mcp.Description("Max documents (default: 20)")),
	), h.MongoQuery)

	s.AddTool(mcp.NewTool("mongo_collections",
		mcp.WithDescription("List all collections in the configured MongoDB database."),
	), h.MongoListCollections)
}
