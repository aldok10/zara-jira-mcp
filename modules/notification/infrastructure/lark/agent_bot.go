package lark

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

// AgentBotHandler is a Lark webhook handler that routes messages through the
// OpenClaw/Hermes-style agent gateway instead of the old simple AI complete.
// This is the upgraded replacement for BotHandler.
type AgentBotHandler struct {
	gateway MessageRouter
	token   string
	client  *WebhookClient
	logger  *slog.Logger
}

// MessageRouter is the interface the agent gateway manager implements.
type MessageRouter interface {
	HandleIncoming(ctx context.Context, channel, channelID, userID, userName, message string)
}

// NewAgentBotHandler creates a new Lark agent bot handler.
func NewAgentBotHandler(client *WebhookClient, router MessageRouter, verificationToken string, logger *slog.Logger) *AgentBotHandler {
	return &AgentBotHandler{
		client:  client,
		gateway: router,
		token:   verificationToken,
		logger:  logger,
	}
}

// ServeHTTP handles Lark webhook events.
// Same webhook interface but routes through the agent gateway.
func (h *AgentBotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "read body failed", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var env eventEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	// Handle URL verification (Lark initial setup)
	if env.Type == "url_verification" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"challenge": env.Challenge}) //nolint:errcheck
		return
	}

	// Verify token if configured
	if h.token != "" && env.Header.Token != h.token {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// Only handle text messages
	if env.Header.EventType != "im.message.receive_v1" || env.Event.Message.MessageType != "text" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var tc textContent
	if err := json.Unmarshal([]byte(env.Event.Message.Content), &tc); err != nil {
		h.logger.Error("parse message content", "err", err)
		w.WriteHeader(http.StatusOK)
		return
	}

	chatID := env.Event.Message.ChatID
	userID := env.Event.Sender.SenderID.UserID
	userOpenID := env.Event.Sender.SenderID.OpenID
	text := strings.TrimSpace(stripMention(tc.Text))

	// Use the open_id as the user identifier
	senderID := userOpenID
	if senderID == "" {
		senderID = userID
	}

	h.logger.Info("lark agent message",
		"chat", chatID,
		"user", senderID,
		"text_preview", truncateStr(text, 80),
	)

	// Route through agent gateway (async to free webhook)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				h.logger.Error("agent gateway panic", "chat", chatID, "recover", r)
			}
		}()
		ctx := context.Background()
		h.gateway.HandleIncoming(ctx, "lark", chatID, senderID, senderID, text)
	}()

	w.WriteHeader(http.StatusOK)
}

func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
