package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerWhatNextTools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_what_next",
		mcp.WithDescription("AI-powered daily priority advisor. Tells SM/PM what to focus on RIGHT NOW based on blockers, risks, dependencies, and sprint state."),
		mcp.WithNumber("board_id", mcp.Description("Board ID for sprint context")),
	), h.WhatNext)

	s.AddTool(mcp.NewTool("pm_one_on_one_prep",
		mcp.WithDescription("Generate 1-on-1 prep notes for a team member. Shows their performance data, workload, patterns, and AI-suggested talking points."),
		mcp.WithString("member", mcp.Required(), mcp.Description("Team member name")),
	), h.OneOnOnePrep)

	s.AddTool(mcp.NewTool("pm_sprint_narrative",
		mcp.WithDescription("Generate sprint story/narrative for Review demo. Translates completed work into stakeholder-friendly business language."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.SprintNarrative)
}
