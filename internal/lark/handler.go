package lark

import (
	"encoding/json"
	"io"
	"net/http"
)

type eventEnvelope struct {
	Schema    string      `json:"schema"`
	Header    eventHeader `json:"header"`
	Event     eventBody   `json:"event"`
	Type      string      `json:"type"`
	Challenge string      `json:"challenge"`
	Token     string      `json:"token"`
}

type eventHeader struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	Token     string `json:"token"`
	AppID     string `json:"app_id"`
}

type eventBody struct {
	Sender  eventSender  `json:"sender"`
	Message eventMessage `json:"message"`
}

type eventSender struct {
	SenderID   senderID `json:"sender_id"`
	SenderType string   `json:"sender_type"`
}

type senderID struct {
	UserID string `json:"user_id"`
	OpenID string `json:"open_id"`
}

type eventMessage struct {
	MessageID   string `json:"message_id"`
	ChatID      string `json:"chat_id"`
	ChatType    string `json:"chat_type"`
	MessageType string `json:"message_type"`
	Content     string `json:"content"`
}

type textContent struct {
	Text string `json:"text"`
}

func (b *BotHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	if env.Type == "url_verification" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"challenge": env.Challenge}) //nolint:errcheck
		return
	}
	if b.token != "" && env.Header.Token != b.token {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}
	if env.Header.EventType != "im.message.receive_v1" || env.Event.Message.MessageType != "text" {
		w.WriteHeader(http.StatusOK)
		return
	}
	var tc textContent
	if err := json.Unmarshal([]byte(env.Event.Message.Content), &tc); err != nil {
		b.logger.Error("parse message content", "err", err)
		w.WriteHeader(http.StatusOK)
		return
	}
	chatID := env.Event.Message.ChatID
	text := tc.Text
	go func() {
		ctx := r.Context()
		reply := b.handleMessage(ctx, text)
		if err := b.replyToChat(ctx, chatID, reply); err != nil {
			b.logger.Error("reply to chat failed", "chat_id", chatID, "err", err)
		}
	}()
	w.WriteHeader(http.StatusOK)
}
