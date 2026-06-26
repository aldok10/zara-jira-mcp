package telegram

import (
	"context"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/aldok10/zara-jira-mcp/config"
)

// Client wraps telegram-bot-api.
type Client struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewClient(cfg *config.Config) *Client {
	c := &Client{}
	if cfg.Telegram.BotToken != "" {
		bot, err := tgbotapi.NewBotAPI(cfg.Telegram.BotToken)
		if err == nil {
			c.bot = bot
		}
	}
	if cfg.Telegram.ChatID != "" {
		c.chatID, _ = strconv.ParseInt(cfg.Telegram.ChatID, 10, 64)
	}
	return c
}

func (c *Client) Available() bool {
	return c.bot != nil
}

func (c *Client) SendMessage(ctx context.Context, chatID int64, text string) error {
	if c.bot == nil {
		return fmt.Errorf("telegram bot not configured")
	}
	if chatID == 0 {
		chatID = c.chatID
	}
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	_, err := c.bot.Send(msg)
	return err
}

func (c *Client) SendHTML(ctx context.Context, chatID int64, html string) error {
	if c.bot == nil {
		return fmt.Errorf("telegram bot not configured")
	}
	if chatID == 0 {
		chatID = c.chatID
	}
	msg := tgbotapi.NewMessage(chatID, html)
	msg.ParseMode = "HTML"
	_, err := c.bot.Send(msg)
	return err
}

func (c *Client) DefaultChatID() int64 {
	return c.chatID
}
