package tools

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// LarkOKRPeriods lists available OKR periods from Lark.
func (h *Handlers) LarkOKRPeriods(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.OKR == nil || !h.OKR.Available() {
		return errorResult("Lark OKR not configured (need LARK_APP_ID + LARK_APP_SECRET)"), nil
	}

	periods, err := h.OKR.ListPeriods(ctx)
	if err != nil {
		return sanitizedError("Lark OKR: failed to fetch periods", err), nil
	}
	if len(periods) == 0 {
		return textResult("No OKR periods found."), nil
	}

	var sb strings.Builder
	sb.WriteString("Lark OKR Periods:\n\n")
	for _, p := range periods {
		status := "unknown"
		switch p.Status {
		case 1:
			status = "in progress"
		case 2:
			status = "not started"
		case 3:
			status = "ended"
		}
		name := p.EnName
		if name == "" {
			name = p.Name
		}
		sb.WriteString(fmt.Sprintf("  [%s] %s (ID: %s)\n", status, name, p.ID))
	}
	return textResult(sb.String()), nil
}

// LarkOKRPull fetches a user's OKRs from Lark for a given period.
func (h *Handlers) LarkOKRPull(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.OKR == nil || !h.OKR.Available() {
		return errorResult("Lark OKR not configured"), nil
	}

	userID, err := req.RequireString("user_id")
	if err != nil {
		return errorResult("user_id required (Lark open_id)"), nil
	}
	if strings.ContainsAny(userID, "/:\\?#") {
		return errorResult("invalid user_id format"), nil
	}
	periodID, err := req.RequireString("period_id")
	if err != nil {
		return errorResult("period_id required (use lark_okr_periods to find)"), nil
	}

	objectives, err := h.OKR.ListUserOKRs(ctx, userID, periodID)
	if err != nil {
		return sanitizedError("Lark OKR: failed to fetch user OKRs", err), nil
	}
	if len(objectives) == 0 {
		return textResult("No OKRs found for this user/period."), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("OKRs for user %s:\n\n", userID))
	for i, obj := range objectives {
		sb.WriteString(fmt.Sprintf("O%d: %s (%d%%)\n", i+1, obj.Content, obj.Progress))
		for j, kr := range obj.KeyResults {
			sb.WriteString(fmt.Sprintf("  KR%d.%d: %s (%d%%)\n", i+1, j+1, kr.Content, kr.Progress))
		}
		sb.WriteString("\n")
	}
	return textResult(sb.String()), nil
}

// LarkOKRSyncProgress pushes a progress update to a Lark OKR key result.
func (h *Handlers) LarkOKRSyncProgress(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	if h.OKR == nil || !h.OKR.Available() {
		return errorResult("Lark OKR not configured"), nil
	}

	krID, err := req.RequireString("kr_id")
	if err != nil {
		return errorResult("kr_id required (Lark key result ID)"), nil
	}
	content, err := req.RequireString("content")
	if err != nil {
		return errorResult("content required (progress update text)"), nil
	}
	dryRun := req.GetString("dry_run", "true")

	if dryRun == "true" {
		return textResult(fmt.Sprintf("DRY RUN — Would push to KR %s:\n\n%s\n\nRun with dry_run:false to push.", krID, content)), nil
	}

	if err := h.OKR.CreateProgressRecord(ctx, krID, content); err != nil {
		return sanitizedError("Lark OKR: failed to sync progress", err), nil
	}
	return textResult(fmt.Sprintf("Progress record pushed to Lark OKR KR %s.", krID)), nil
}
