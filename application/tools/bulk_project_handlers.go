package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// BulkTransition transitions multiple issues at once.
func (h *Handlers) BulkTransition(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	keysRaw, err := req.RequireString("issue_keys")
	if err != nil {
		return errorResult("issue_keys parameter is required"), nil
	}
	transitionID, err := req.RequireString("transition_id")
	if err != nil {
		return errorResult("transition_id parameter is required"), nil
	}

	keys := strings.Split(keysRaw, ",")
	var success, failed int
	var errors []string
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if err := h.Jira.TransitionIssue(ctx, key, transitionID); err != nil {
			failed++
			errors = append(errors, fmt.Sprintf("%s: %s", key, err.Error()))
		} else {
			success++
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Bulk transition: %d succeeded, %d failed\n", success, failed))
	for _, e := range errors {
		sb.WriteString(fmt.Sprintf("  - %s\n", e))
	}
	return textResult(sb.String()), nil
}

// BulkAssign assigns multiple issues to one person.
func (h *Handlers) BulkAssign(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	keysRaw, err := req.RequireString("issue_keys")
	if err != nil {
		return errorResult("issue_keys parameter is required"), nil
	}
	assigneeID, err := req.RequireString("assignee_id")
	if err != nil {
		return errorResult("assignee_id parameter is required"), nil
	}

	keys := strings.Split(keysRaw, ",")
	var success, failed int
	var errors []string
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if err := h.Jira.AssignIssue(ctx, key, assigneeID); err != nil {
			failed++
			errors = append(errors, fmt.Sprintf("%s: %s", key, err.Error()))
		} else {
			success++
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Bulk assign: %d succeeded, %d failed\n", success, failed))
	for _, e := range errors {
		sb.WriteString(fmt.Sprintf("  - %s\n", e))
	}
	return textResult(sb.String()), nil
}

// BulkLabel adds a label to multiple issues.
func (h *Handlers) BulkLabel(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	keysRaw, err := req.RequireString("issue_keys")
	if err != nil {
		return errorResult("issue_keys parameter is required"), nil
	}
	label, err := req.RequireString("label")
	if err != nil {
		return errorResult("label parameter is required"), nil
	}

	keys := strings.Split(keysRaw, ",")
	var success, failed int
	var errors []string
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if err := h.Jira.AddLabel(ctx, key, label); err != nil {
			failed++
			errors = append(errors, fmt.Sprintf("%s: %s", key, err.Error()))
		} else {
			success++
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Bulk label '%s': %d succeeded, %d failed\n", label, success, failed))
	for _, e := range errors {
		sb.WriteString(fmt.Sprintf("  - %s\n", e))
	}
	return textResult(sb.String()), nil
}

// ListProjects lists accessible Jira projects.
func (h *Handlers) ListProjects(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	projects, err := h.Jira.GetProjects(ctx)
	if err != nil {
		return sanitizedError("failed to list projects", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Projects (%d):\n\n", len(projects)))
	for _, p := range projects {
		sb.WriteString(fmt.Sprintf("- **%s** %s | Lead: %s | Type: %s\n", p.Key, p.Name, p.Lead, p.Type))
	}
	return textResult(sb.String()), nil
}

// ProjectDetail gets full project details.
func (h *Handlers) ProjectDetail(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	key, err := req.RequireString("key")
	if err != nil {
		return errorResult("key parameter is required"), nil
	}

	project, err := h.Jira.GetProject(ctx, key)
	if err != nil {
		return sanitizedError("failed to get project details", err), nil
	}

	data, _ := json.MarshalIndent(project, "", "  ")
	return textResult(string(data)), nil
}

// RawRequest makes an arbitrary Jira REST API call.
func (h *Handlers) RawRequest(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	method, err := req.RequireString("method")
	if err != nil {
		return errorResult("method parameter is required"), nil
	}
	path, err := req.RequireString("path")
	if err != nil {
		return errorResult("path parameter is required"), nil
	}

	method = strings.ToUpper(method)
	if method != "GET" && method != "POST" && method != "PUT" && method != "DELETE" {
		return errorResult("method must be GET, POST, PUT, or DELETE"), nil
	}

	var body []byte
	bodyStr := req.GetString("body", "")
	if bodyStr != "" {
		body = []byte(bodyStr)
	}

	respBody, statusCode, err := h.Jira.RawRequest(ctx, method, path, body)
	if err != nil {
		return sanitizedError("raw jira request failed", err), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Status: %d\n\n", statusCode))
	if len(respBody) > 0 {
		// Try to pretty-print JSON
		var raw json.RawMessage
		if json.Unmarshal(respBody, &raw) == nil {
			pretty, _ := json.MarshalIndent(raw, "", "  ")
			sb.Write(pretty)
		} else {
			sb.Write(respBody)
		}
	}
	return textResult(sb.String()), nil
}
