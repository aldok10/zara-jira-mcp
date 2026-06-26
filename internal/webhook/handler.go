package webhook

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/aldok10/zara-jira-mcp/domain/memory"
)

// Event represents a parsed Jira webhook event.
type Event struct {
	WebhookEvent string `json:"webhookEvent"`
	Issue        struct {
		Key    string `json:"key"`
		Fields struct {
			Summary string `json:"summary"`
			Status  struct {
				Name string `json:"name"`
			} `json:"status"`
			Assignee struct {
				DisplayName string `json:"displayName"`
			} `json:"assignee"`
		} `json:"fields"`
	} `json:"issue"`
	Sprint struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		State string `json:"state"`
	} `json:"sprint"`
	Changelog struct {
		Items []struct {
			Field      string `json:"field"`
			FromString string `json:"fromString"`
			ToString   string `json:"toString"`
		} `json:"items"`
	} `json:"changelog"`
}

// Handler processes incoming Jira webhook events.
type Handler struct {
	secret string
	memory memory.Store
	logger *slog.Logger
}

// NewHandler creates a Jira webhook handler.
func NewHandler(secret string, mem memory.Store, logger *slog.Logger) *Handler {
	return &Handler{secret: secret, memory: mem, logger: logger}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1MB limit
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if h.secret != "" && !h.verifySignature(r, body) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var event Event
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	h.processEvent(r.Context(), &event)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) verifySignature(r *http.Request, body []byte) bool {
	// Jira Cloud/Server: check X-Atlassian-Webhook-Signature first, fall back to X-Hub-Signature
	sig := r.Header.Get("X-Atlassian-Webhook-Signature")
	if sig == "" {
		sig = r.Header.Get("X-Hub-Signature")
	}
	if sig == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(h.secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(sig))
}

func (h *Handler) processEvent(ctx context.Context, event *Event) {
	switch {
	case event.WebhookEvent == "jira:issue_updated":
		h.handleIssueUpdated(ctx, event)
	case event.WebhookEvent == "jira:issue_created":
		h.logger.Info("webhook: issue created", "key", event.Issue.Key, "summary", event.Issue.Fields.Summary)
	case event.WebhookEvent == "sprint_started":
		h.logger.Info("webhook: sprint started", "name", event.Sprint.Name)
	case event.WebhookEvent == "sprint_closed":
		h.handleSprintClosed(ctx, event)
	}
}

func (h *Handler) handleIssueUpdated(ctx context.Context, event *Event) {
	for _, item := range event.Changelog.Items {
		if item.Field == "status" && strings.EqualFold(item.ToString, "Blocked") {
			b := &memory.Blocker{
				IssueKey:     event.Issue.Key,
				Description:  fmt.Sprintf("Auto-detected: %s (%s) moved to Blocked", event.Issue.Key, event.Issue.Fields.Summary),
				BlockedSince: time.Now(),
				Owner:        event.Issue.Fields.Assignee.DisplayName,
			}
			if err := h.memory.SaveBlocker(ctx, b); err != nil {
				h.logger.Error("webhook: failed to record blocker", "err", err)
			} else {
				h.logger.Info("webhook: auto-recorded blocker", "key", event.Issue.Key)
			}
		}
	}
}

func (h *Handler) handleSprintClosed(ctx context.Context, event *Event) {
	h.logger.Info("webhook: sprint closed", "name", event.Sprint.Name, "id", event.Sprint.ID)
	note := &memory.MeetingNote{
		MeetingType: "adhoc",
		Date:        time.Now(),
		Notes:       fmt.Sprintf("Sprint '%s' was closed via webhook. PM should run pm_snapshot_sprint.", event.Sprint.Name),
		SprintName:  event.Sprint.Name,
	}
	if err := h.memory.SaveMeetingNote(ctx, note); err != nil {
		h.logger.Error("webhook: failed to save sprint close note", "err", err)
	}
}
