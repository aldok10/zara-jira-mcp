// Package main is the entry point for zara-jira-mcp MCP server.
package main

import (
	"log"

	"github.com/aldok10/zara-jira-mcp/apps/api/internal/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
