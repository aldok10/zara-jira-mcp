package transport

import (
	"github.com/aldok10/zara-jira-mcp/application/tools"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerV6Tools(s *server.MCPServer, h *tools.Handlers) {
	s.AddTool(mcp.NewTool("pm_hypothesis",
		mcp.WithDescription("Record improvement hypothesis: belief + expected outcome."),
		mcp.WithString("belief", mcp.Required(), mcp.Description("What you think will improve")),
		mcp.WithString("expected_outcome", mcp.Required(), mcp.Description("Expected result")),
		mcp.WithString("measure", mcp.Description("How to verify (default: observe)")),
		mcp.WithString("duration", mcp.Description("Duration (default: 1 sprint)")),
		mcp.WithString("sprint_name", mcp.Description("Sprint context")),
	), h.PMHypothesis)

	s.AddTool(mcp.NewTool("pm_hypothesis_review",
		mcp.WithDescription("Show all hypotheses and validation status."),
	), h.PMHypothesisReview)

	s.AddTool(mcp.NewTool("pm_hypothesis_close",
		mcp.WithDescription("Validate or invalidate a hypothesis."),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("Hypothesis ID")),
		mcp.WithString("status", mcp.Required(), mcp.Description("validated or invalidated")),
		mcp.WithString("actual_outcome", mcp.Description("What happened")),
	), h.PMHypothesisClose)

	s.AddTool(mcp.NewTool("pm_estimation_accuracy",
		mcp.WithDescription("Committed vs delivered across sprints."),
		mcp.WithNumber("board_id", mcp.Required(), mcp.Description("Board ID")),
	), h.PMEstimationAccuracy)

	s.AddTool(mcp.NewTool("pm_space",
		mcp.WithDescription("SPACE: Satisfaction, Performance, Activity, Comms, Efficiency."),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
	), h.PMSpaceMetrics)

	s.AddTool(mcp.NewTool("pm_ebm",
		mcp.WithDescription("Evidence-Based Management: 4 Key Value Areas."),
		mcp.WithNumber("board_id", mcp.Description("Board ID")),
	), h.PMEBMDashboard)
}
