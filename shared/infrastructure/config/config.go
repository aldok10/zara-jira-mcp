package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all application configuration parsed from environment variables.
type Config struct {
	Server         ServerConfig
	Jira           JiraConfig
	AI             AIConfig
	Lark           LarkConfig
	Slack          SlackConfig
	Discord        DiscordConfig
	Telegram       TelegramConfig
	Teams          TeamsConfig
	Email          EmailConfig
	Confluence     ConfluenceConfig
	Memory         MemoryConfig
	GoogleCalendar GoogleCalendarConfig
	GitHub         GitHubConfig
	GitLab         GitLabConfig
	Notion         NotionConfig
	Clockify       ClockifyConfig
	Linear         LinearConfig
	PagerDuty      PagerDutyConfig
	GoogleSheets   GoogleSheetsConfig
	Database       DatabaseConfig
	Redis          RedisConfig
	Webhook        WebhookConfig
}

type DatabaseConfig struct {
	PostgresDSN string
	MySQLDSN    string
	MongoURI    string
}

type RedisConfig struct {
	URL string // REDIS_URL (e.g. redis://localhost:6379)
	TTL string // REDIS_TTL (default: 5m)
}

type WebhookConfig struct {
	Enabled bool   // WEBHOOK_ENABLED
	Port    string // WEBHOOK_PORT (default: 8081)
	Secret  string // WEBHOOK_SECRET (for signature verification)
}

type GoogleCalendarConfig struct {
	APIKey     string // GOOGLE_CALENDAR_API_KEY
	CalendarID string // GOOGLE_CALENDAR_ID
}

type GitHubConfig struct {
	Token string // GITHUB_TOKEN
	Owner string // GITHUB_OWNER
	Repo  string // GITHUB_REPO
}

type GitLabConfig struct {
	Token     string // GITLAB_TOKEN
	BaseURL   string // GITLAB_BASE_URL (default: https://gitlab.com)
	ProjectID string // GITLAB_PROJECT_ID
}

type NotionConfig struct {
	APIKey     string // NOTION_API_KEY
	DatabaseID string // NOTION_DATABASE_ID (optional)
}

type ClockifyConfig struct {
	APIKey      string // CLOCKIFY_API_KEY
	WorkspaceID string // CLOCKIFY_WORKSPACE_ID
}

type LinearConfig struct {
	APIKey string // LINEAR_API_KEY
}

type PagerDutyConfig struct {
	APIKey string // PAGERDUTY_API_KEY
}

type GoogleSheetsConfig struct {
	APIKey        string // GOOGLE_SHEETS_API_KEY
	SpreadsheetID string // GOOGLE_SHEETS_SPREADSHEET_ID
}

type DiscordConfig struct {
	BotToken   string // DISCORD_BOT_TOKEN
	ChannelID  string // DISCORD_CHANNEL_ID (default channel)
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

type ServerConfig struct {
	Transport               string // MCP_TRANSPORT: stdio (default), sse, http
	Port                    string // MCP_PORT (default: 8080)
	Profile                 string // PM_PROFILE: unified, lite, standard, full (default: all)
	DashboardEnabled        bool   // MCP_DASHBOARD: true to enable web dashboard
	DashboardPort           string // MCP_DASHBOARD_PORT (default: 9090)
	NotificationDailyBudget int    // PM_NOTIFICATION_DAILY_BUDGET: max notifications per day (default: 10)
}

type LarkConfig struct {
	BotEnabled        bool   // LARK_BOT_ENABLED
	BotPort           string // LARK_BOT_PORT (default: 9091)
	VerificationToken string // LARK_VERIFICATION_TOKEN
	EncryptKey        string // LARK_ENCRYPT_KEY
	WebhookURL        string // JIRA_LARK_WEBHOOK_URL (simple webhook, no app needed)
	AppID             string // LARK_APP_ID (for full SDK access)
	AppSecret         string // LARK_APP_SECRET
	ChatID            string // LARK_CHAT_ID (target chat for SDK messages)
}

type JiraConfig struct {
	BaseURL string // JIRA_BASE_URL (e.g. https://company.atlassian.net)
	Email   string // JIRA_EMAIL
	Token   string // JIRA_API_TOKEN (store securely, prefer OAuth or JWT)
}

type AIConfig struct {
	BaseURL string // JIRA_AI_BASE_URL (OpenAI-compatible endpoint)
	APIKey  string // JIRA_AI_API_KEY
	Model   string // JIRA_AI_MODEL (default: gpt-4o-mini)
}

// Load reads configuration exclusively from environment variables.
// No credential data is persisted to disk — secrets are memory-only.
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Transport:               os.Getenv("MCP_TRANSPORT"),
			Port:                    os.Getenv("MCP_PORT"),
			Profile:                 os.Getenv("PM_PROFILE"),
			DashboardEnabled:        os.Getenv("MCP_DASHBOARD") == "true",
			DashboardPort:           os.Getenv("MCP_DASHBOARD_PORT"),
			NotificationDailyBudget: 10,
		},
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
			BotEnabled:        os.Getenv("LARK_BOT_ENABLED") == "true",
			BotPort:           os.Getenv("LARK_BOT_PORT"),
			VerificationToken: os.Getenv("LARK_VERIFICATION_TOKEN"),
			EncryptKey:        os.Getenv("LARK_ENCRYPT_KEY"),
			WebhookURL:        os.Getenv("JIRA_LARK_WEBHOOK_URL"),
			AppID:             os.Getenv("LARK_APP_ID"),
			AppSecret:         os.Getenv("LARK_APP_SECRET"),
			ChatID:            os.Getenv("LARK_CHAT_ID"),
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
		GoogleCalendar: GoogleCalendarConfig{
			APIKey:     os.Getenv("GOOGLE_CALENDAR_API_KEY"),
			CalendarID: os.Getenv("GOOGLE_CALENDAR_ID"),
		},
		GitHub: GitHubConfig{
			Token: os.Getenv("GITHUB_TOKEN"),
			Owner: os.Getenv("GITHUB_OWNER"),
			Repo:  os.Getenv("GITHUB_REPO"),
		},
		GitLab: GitLabConfig{
			Token:     os.Getenv("GITLAB_TOKEN"),
			BaseURL:   os.Getenv("GITLAB_BASE_URL"),
			ProjectID: os.Getenv("GITLAB_PROJECT_ID"),
		},
		Notion: NotionConfig{
			APIKey:     os.Getenv("NOTION_API_KEY"),
			DatabaseID: os.Getenv("NOTION_DATABASE_ID"),
		},
		Clockify: ClockifyConfig{
			APIKey:      os.Getenv("CLOCKIFY_API_KEY"),
			WorkspaceID: os.Getenv("CLOCKIFY_WORKSPACE_ID"),
		},
		Linear: LinearConfig{
			APIKey: os.Getenv("LINEAR_API_KEY"),
		},
		PagerDuty: PagerDutyConfig{
			APIKey: os.Getenv("PAGERDUTY_API_KEY"),
		},
		GoogleSheets: GoogleSheetsConfig{
			APIKey:        os.Getenv("GOOGLE_SHEETS_API_KEY"),
			SpreadsheetID: os.Getenv("GOOGLE_SHEETS_SPREADSHEET_ID"),
		},
		Database: DatabaseConfig{
			PostgresDSN: os.Getenv("DATABASE_POSTGRES_DSN"),
			MySQLDSN:    os.Getenv("DATABASE_MYSQL_DSN"),
			MongoURI:    os.Getenv("DATABASE_MONGO_URI"),
		},
		Redis: RedisConfig{
			URL: os.Getenv("REDIS_URL"),
			TTL: os.Getenv("REDIS_TTL"),
		},
		Webhook: WebhookConfig{
			Enabled: os.Getenv("WEBHOOK_ENABLED") == "true",
			Port:    os.Getenv("WEBHOOK_PORT"),
			Secret:  os.Getenv("WEBHOOK_SECRET"),
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

	if cfg.Server.Transport == "" {
		cfg.Server.Transport = "stdio"
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}
	if cfg.Server.DashboardPort == "" {
		cfg.Server.DashboardPort = "9090"
	}
	if cfg.Lark.BotPort == "" {
		cfg.Lark.BotPort = "9091"
	}
	if cfg.Webhook.Port == "" {
		cfg.Webhook.Port = "8081"
	}

	return cfg, nil
}
