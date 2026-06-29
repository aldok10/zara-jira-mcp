// Package mcputil provides shared MCP handler utilities used across all modules.
// This is part of the infrastructure layer — adapters for the MCP framework.
package mcputil

import (
	"fmt"
	"log/slog"
	"strings"
	"unicode"

	"github.com/mark3labs/mcp-go/mcp"
)

// TextResult creates a successful MCP tool result with formatted text.
func TextResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.TextContent{Type: "text", Text: text}},
	}
}

// ErrorResult creates an error MCP tool result with a user-safe message.
func ErrorResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		IsError: true,
		Content: []mcp.Content{mcp.TextContent{Type: "text", Text: msg}},
	}
}

// ErrJira returns a Jira API error with actionable guidance.
func ErrJira(action string, err error) *mcp.CallToolResult {
	slog.Error(action, "detail", err.Error())
	switch {
	case isAuthError(err):
		return ErrorResult("Jira auth failed. Check JIRA_EMAIL and JIRA_API_TOKEN — they may be expired or invalid.")
	case isRateLimit(err):
		return ErrorResult("Jira API rate limited. Wait 10-30 seconds, then retry. Rate limits reset per minute.")
	case isNotFound(err):
		return ErrorResult("Jira resource not found. Verify the ID/key and confirm you have project access.")
	case isNetworkError(err):
		return ErrorResult("Jira unreachable. Check your network, VPN, and that JIRA_BASE_URL is correct.")
	case isServerError(err):
		return ErrorResult("Jira server error (5xx). The service may be degraded — retry in a few minutes.")
	default:
		return ErrorResult(fmt.Sprintf("Jira operation failed: %s. Verify your configuration and try again.", action))
	}
}

// ErrAI returns an AI provider error with actionable guidance.
func ErrAI(action string, err error) *mcp.CallToolResult {
	slog.Error(action, "detail", err.Error())
	if isRateLimit(err) {
		return ErrorResult("AI provider rate limited. Wait a moment and retry. (10 req/min limit.)")
	}
	if isAuthError(err) {
		return ErrorResult("AI authentication failed. Check JIRA_AI_API_KEY and that the key hasn't expired.")
	}
	if isNetworkError(err) {
		return ErrorResult("AI provider unreachable. Check JIRA_AI_BASE_URL, network, and firewall.")
	}
	return ErrorResult("AI operation failed. Verify JIRA_AI_BASE_URL and JIRA_AI_API_KEY are configured correctly.")
}

// ErrInvalid returns a validation error with the specific issue.
func ErrInvalid(msg string) *mcp.CallToolResult {
	slog.Warn("validation failed", "message", msg)
	return ErrorResult(fmt.Sprintf("Invalid input: %s. Fix the value and retry.", msg))
}

// ErrInternal returns a generic internal error (no details leaked).
func ErrInternal(action string, err error) *mcp.CallToolResult {
	slog.Error(action, "detail", err)
	return ErrorResult(fmt.Sprintf("Internal error during '%s'. Check server logs for details.", action))
}

// ErrorHandler provides consistent error handling for MCP tools.
type ErrorHandler struct {
	logger *slog.Logger
}

// NewErrorHandler creates a new error handler.
func NewErrorHandler(logger *slog.Logger) *ErrorHandler {
	if logger == nil {
		logger = slog.Default()
	}
	return &ErrorHandler{logger: logger}
}

// Result creates a successful text result.
func (h *ErrorHandler) Result(text string) *mcp.CallToolResult {
	return TextResult(text)
}

// Error creates an error result.
func (h *ErrorHandler) Error(message string) *mcp.CallToolResult {
	return ErrorResult(message)
}

// Wrap wraps an error with context and returns an error result.
func (h *ErrorHandler) Wrap(operation string, err error) *mcp.CallToolResult {
	if err == nil {
		return h.Error("Jira operation failed: " + operation + ". Unknown error.")
	}
	return ErrJira(operation, err)
}

// WrapAI wraps an AI provider error with actionable guidance.
func (h *ErrorHandler) WrapAI(operation string, err error) *mcp.CallToolResult {
	if err == nil {
		return h.Error("AI operation failed: " + operation)
	}
	return ErrAI(operation, err)
}

// WrapNotFound wraps a "not found" error with the specific resource type.
func (h *ErrorHandler) WrapNotFound(kind, id string, err error) *mcp.CallToolResult {
	h.logger.Error("not found", "kind", kind, "id", id, "detail", err)
	return h.Error(kind + " " + id + " not found. Verify it exists and you have access.")
}

// WrapValidation wraps a validation error.
func (h *ErrorHandler) WrapValidation(msg string) *mcp.CallToolResult {
	h.logger.Warn("validation failed", "message", msg)
	return h.Error("Invalid input: " + msg)
}

// WrapInternal wraps an internal error without leaking details.
func (h *ErrorHandler) WrapInternal(operation string, err error) *mcp.CallToolResult {
	h.logger.Error(operation, "detail", err)
	return h.Error("Internal error. Check server logs for details.")
}

// IsSafeJQLValue checks that a string can be safely interpolated into JQL.
func IsSafeJQLValue(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != ' ' && r != '_' && r != '-' {
			return false
		}
	}
	return true
}

// TruncateStr truncates a string to maxLen chars, adding "..." if truncated.
func TruncateStr(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// --- private error classifiers ---

func isAuthError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "unauthorized") ||
		strings.Contains(msg, "forbidden") ||
		strings.Contains(msg, "401") ||
		strings.Contains(msg, "403") ||
		strings.Contains(msg, "invalid token") ||
		strings.Contains(msg, "api key")
}

func isRateLimit(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "rate limit") ||
		strings.Contains(msg, "too many requests") ||
		strings.Contains(msg, "429")
}

func isNotFound(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "not found") ||
		strings.Contains(msg, "404") ||
		strings.Contains(msg, "no results")
}

func isNetworkError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "connection") ||
		strings.Contains(msg, "timeout") ||
		strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "refused") ||
		func() bool {
			return strings.Contains(msg, "dial") || strings.Contains(msg, "eof") || strings.Contains(msg, "reset") ||
				strings.Contains(msg, "unreachable") || strings.Contains(msg, "lookup") ||
				strings.Contains(msg, "timeout") || strings.Contains(msg, "canceled")
		}() ||
		strings.Contains(err.Error(), "no address") // Avoid closure over err
}

func isServerError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "500") ||
		strings.Contains(msg, "502") ||
		strings.Contains(msg, "503") ||
		strings.Contains(msg, "internal server error") ||
		strings.Contains(msg, "service unavailable")
}
