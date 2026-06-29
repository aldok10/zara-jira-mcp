package lark

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	domain "github.com/aldok10/zara-jira-mcp/modules/sprint/domain/memory"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

type AIProvider interface {
	Complete(ctx context.Context, system, user string) (string, error)
}

type MemoryProvider interface {
	GetActiveBlockers(ctx context.Context) ([]domain.Blocker, error)
	GetOpenRisks(ctx context.Context) ([]domain.Risk, error)
	GetPendingActionItems(ctx context.Context) ([]domain.ActionItem, error)
}

type BotHandler struct {
	client *WebhookClient
	ai     AIProvider
	memory MemoryProvider
	token  string
	logger *slog.Logger
}

func NewBotHandler(client *WebhookClient, ai AIProvider, mem MemoryProvider, verificationToken string, logger *slog.Logger) *BotHandler {
	return &BotHandler{client: client, ai: ai, memory: mem, token: verificationToken, logger: logger}
}

const botSystemPrompt = `You are a PM assistant in a Lark group chat. You help the team with sprint status, blockers, risks, and action items. Answer concisely in 2-3 sentences max. Use bullet points for lists. No walls of text.`

func (b *BotHandler) handleMessage(ctx context.Context, text string) string {
	text = stripMention(text)
	text = strings.TrimSpace(text)
	if text == "" {
		return "Hi! I am your PM assistant. Try /help to see what I can do."
	}
	if strings.HasPrefix(text, "/") {
		return b.handleCommand(ctx, text)
	}
	resp, err := b.ai.Complete(ctx, botSystemPrompt, text)
	if err != nil {
		b.logger.Error("AI complete failed", "err", err)
		return "Sorry, I could not process that right now."
	}
	return resp
}

func (b *BotHandler) replyToChat(ctx context.Context, chatID, text string) error {
	if b.client.sdk == nil {
		return fmt.Errorf("lark SDK not configured, cannot reply")
	}
	content := larkim.NewTextMsgBuilder().Text(text).Build()
	resp, err := b.client.sdk.Im.Message.Create(ctx, larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(larkim.CreateMessageV1ReceiveIDTypeChatId).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			MsgType(larkim.MsgTypeText).
			ReceiveId(chatID).
			Content(content).
			Build()).
		Build())
	if err != nil {
		return err
	}
	if !resp.Success() {
		return fmt.Errorf("lark reply error %d: %s", resp.Code, resp.Msg)
	}
	return nil
}

func stripMention(text string) string {
	if idx := strings.Index(text, " "); idx > 0 && strings.HasPrefix(text, "@") {
		return text[idx+1:]
	}
	return text
}
