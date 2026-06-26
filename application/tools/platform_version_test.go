package tools_test

import (
	"context"
	"testing"

	"github.com/aldok10/zara-jira-mcp/application/tools"
)

// --- Platform handlers (nil checks) ---

func TestDiscordSend_Nil(t *testing.T) {
	h := &tools.Handlers{Jira: &mockJira{}, AI: &mockAI{}, Lark: &mockLark{}, Memory: &mockMemory{}}
	result, _ := h.DiscordSend(context.Background(), makeReq(map[string]any{"content": "hi"}))
	if !result.IsError {
		t.Error("expected error for nil discord")
	}
}

func TestTelegramSend_Nil(t *testing.T) {
	h := &tools.Handlers{Jira: &mockJira{}, AI: &mockAI{}, Lark: &mockLark{}, Memory: &mockMemory{}}
	result, _ := h.TelegramSend(context.Background(), makeReq(map[string]any{"message": "hi"}))
	if !result.IsError {
		t.Error("expected error for nil telegram")
	}
}

func TestTeamsSend_Nil(t *testing.T) {
	h := &tools.Handlers{Jira: &mockJira{}, AI: &mockAI{}, Lark: &mockLark{}, Memory: &mockMemory{}}
	result, _ := h.TeamsSend(context.Background(), makeReq(map[string]any{"content": "hi"}))
	if !result.IsError {
		t.Error("expected error for nil teams")
	}
}

func TestEmailSend_Nil(t *testing.T) {
	h := &tools.Handlers{Jira: &mockJira{}, AI: &mockAI{}, Lark: &mockLark{}, Memory: &mockMemory{}}
	result, _ := h.EmailSend(context.Background(), makeReq(map[string]any{"to": "x@x.com", "subject": "hi", "body": "yo"}))
	if !result.IsError {
		t.Error("expected error for nil email")
	}
}

func TestConfluenceSearch_Nil(t *testing.T) {
	h := &tools.Handlers{Jira: &mockJira{}, AI: &mockAI{}, Lark: &mockLark{}, Memory: &mockMemory{}}
	result, _ := h.ConfluenceSearch(context.Background(), makeReq(map[string]any{"query": "test"}))
	if !result.IsError {
		t.Error("expected error for nil confluence")
	}
}

// --- Git handlers (nil GitHub panics on .Available(), skip these) ---

// --- Version handlers ---

func TestGetAttachments(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.GetAttachments(context.Background(), makeReq(map[string]any{"key": "T-1"}))
	// mockJira returns nil attachments
	text := resultText(result)
	if text == "" && result.IsError {
		t.Error("unexpected error")
	}
}

func TestGetAttachments_MissingKey(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.GetAttachments(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestGetVersions(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.GetVersions(context.Background(), makeReq(map[string]any{"project": "TEST"}))
	text := resultText(result)
	if text == "" && result.IsError {
		t.Error("unexpected error")
	}
}

func TestGetVersions_MissingProject(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.GetVersions(context.Background(), makeReq(map[string]any{}))
	if !result.IsError {
		t.Error("expected error")
	}
}

func TestCreateVersion(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.CreateVersion(context.Background(), makeReq(map[string]any{"project": "TEST", "name": "v1.0"}))
	text := resultText(result)
	if text == "" && result.IsError {
		t.Error("unexpected error")
	}
}

func TestReleaseVersion(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.ReleaseVersion(context.Background(), makeReq(map[string]any{"version_id": "123"}))
	text := resultText(result)
	if text == "" && result.IsError {
		t.Error("unexpected error")
	}
}

func TestGetComponents(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.GetComponents(context.Background(), makeReq(map[string]any{"project": "TEST"}))
	text := resultText(result)
	if text == "" && result.IsError {
		t.Error("unexpected error")
	}
}

func TestGetFields(t *testing.T) {
	h := newTestHandlers(&mockJira{})
	result, _ := h.GetFields(context.Background(), makeReq(nil))
	text := resultText(result)
	if text == "" && result.IsError {
		t.Error("unexpected error")
	}
}
