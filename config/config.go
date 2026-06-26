package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all application configuration parsed from environment variables.
type Config struct {
	Jira       JiraConfig
	AI         AIConfig
	Lark       LarkConfig
	Slack      SlackConfig
	Discord    DiscordConfig
	Telegram   TelegramConfig
	Teams      TeamsConfig
	Email      EmailConfig
	Confluence ConfluenceConfig
	Memory     MemoryConfig
}

type DiscordConfig struct {
	BotToken  string // DISCORD_BOT_TOKEN
	ChannelID string // DISCORD_CHANNEL_ID (default channel)
	WebhookURL string // DISCORD_WEBHOOK_URL (fallback)
}

type TelegramConfig struct {
	BotToken string // TELEGRAM_BOT_TOKEN
	ChatID   string // TELEGRAM_CHAT_ID (default chat/group)
}

type TeamsConfig struct {
	WebhookURL string // TEAMS_WEBHOOK_URL (incoming webhook connector)
}

type EmailConfig struct {
	SMTPHost string // EMAIL_SMTP_HOST
	SMTPPort string // EMAIL_SMTP_PORT (default: 587)
	Username string // EMAIL_USERNAME
	Password string // EMAIL_PASSWORD
	From     string // EMAIL_FROM
}

type ConfluenceConfig struct {
	BaseURL string // CONFLUENCE_BASE_URL
	Email   string // CONFLUENCE_EMAIL
	Token   string // CONFLUENCE_API_TOKEN
}

type SlackConfig struct {
	BotToken       string // SLACK_BOT_TOKEN (xoxb-...)
	DefaultChannel string // SLACK_DEFAULT_CHANNEL (channel ID or name)
	WebhookURL     string // SLACK_WEBHOOK_URL (fallback, no token needed)
}

type MemoryConfig struct {
	DBPath string // PM_MEMORY_DB_PATH (default: ~/.zara-jira-mcp/pm_memory.db)
}

type LarkConfig struct {
	WebhookURL string // JIRA_LARK_WEBHOOK_URL (simple webhook, no app needed)
	AppID      string // LARK_APP_ID (for full SDK access)
	AppSecret  string // LARK_APP_SECRET
	ChatID     string // LARK_CHAT_ID (target chat for SDK messages)
}

type JiraConfig struct {
	BaseURL string // JIRA_BASE_URL (e.g. https://company.atlassian.net)
	Email   string // JIRA_EMAIL
	Token   string // JIRA_API_TOKEN
}

type AIConfig struct {
	BaseURL string // JIRA_AI_BASE_URL (OpenAI-compatible endpoint)
	APIKey  string // JIRA_AI_API_KEY
	Model   string // JIRA_AI_MODEL (default: gpt-4o-mini)
}

func Load() (*Config, error) {
	cfg := &Config{
		Jira: JiraConfig{
			BaseURL: strings.TrimRight(os.Getenv("JIRA_BASE_URL"), "/"),
			Email:   os.Getenv("JIRA_EMAIL"),
			Token:   os.Getenv("JIRA_API_TOKEN"),
		},
		AI: AIConfig{
			BaseURL: os.Getenv("JIRA_AI_BASE_URL"),
			APIKey:  os.Getenv("JIRA_AI_API_KEY"),
			Model:   os.Getenv("JIRA_AI_MODEL"),
		},
		Lark: LarkConfig{
			WebhookURL: os.Getenv("JIRA_LARK_WEBHOOK_URL"),
			AppID:      os.Getenv("LARK_APP_ID"),
			AppSecret:  os.Getenv("LARK_APP_SECRET"),
			ChatID:     os.Getenv("LARK_CHAT_ID"),
		},
		Slack: SlackConfig{
			BotToken:       os.Getenv("SLACK_BOT_TOKEN"),
			DefaultChannel: os.Getenv("SLACK_DEFAULT_CHANNEL"),
			WebhookURL:     os.Getenv("SLACK_WEBHOOK_URL"),
		},
		Discord: DiscordConfig{
			BotToken:   os.Getenv("DISCORD_BOT_TOKEN"),
			ChannelID:  os.Getenv("DISCORD_CHANNEL_ID"),
			WebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		},
		Telegram: TelegramConfig{
			BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
			ChatID:   os.Getenv("TELEGRAM_CHAT_ID"),
		},
		Teams: TeamsConfig{
			WebhookURL: os.Getenv("TEAMS_WEBHOOK_URL"),
		},
		Email: EmailConfig{
			SMTPHost: os.Getenv("EMAIL_SMTP_HOST"),
			SMTPPort: os.Getenv("EMAIL_SMTP_PORT"),
			Username: os.Getenv("EMAIL_USERNAME"),
			Password: os.Getenv("EMAIL_PASSWORD"),
			From:     os.Getenv("EMAIL_FROM"),
		},
		Confluence: ConfluenceConfig{
			BaseURL: strings.TrimRight(os.Getenv("CONFLUENCE_BASE_URL"), "/"),
			Email:   os.Getenv("CONFLUENCE_EMAIL"),
			Token:   os.Getenv("CONFLUENCE_API_TOKEN"),
		},
		Memory: MemoryConfig{
			DBPath: os.Getenv("PM_MEMORY_DB_PATH"),
		},
	}

	if cfg.Jira.BaseURL == "" {
		return nil, fmt.Errorf("JIRA_BASE_URL is required")
	}
	if cfg.Jira.Email == "" {
		return nil, fmt.Errorf("JIRA_EMAIL is required")
	}
	if cfg.Jira.Token == "" {
		return nil, fmt.Errorf("JIRA_API_TOKEN is required")
	}

	if cfg.AI.Model == "" {
		cfg.AI.Model = "gpt-4o-mini"
	}

	return cfg, nil
}
