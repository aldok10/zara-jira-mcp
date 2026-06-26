package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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

type ServerConfig struct {
	Transport        string // MCP_TRANSPORT: stdio (default), sse, http
	Port             string // MCP_PORT (default: 8080)
	Profile          string // PM_PROFILE: unified, lite, standard, full (default: all)
	DashboardEnabled bool   // MCP_DASHBOARD: true to enable web dashboard
	DashboardPort    string // MCP_DASHBOARD_PORT (default: 9090)
}

type LarkConfig struct {
	BotEnabled        bool   // LARK_BOT_ENABLED
	BotPort           string // LARK_BOT_PORT (default: 9091)
	VerificationToken string // LARK_VERIFICATION_TOKEN
	EncryptKey        string // LARK_ENCRYPT_KEY
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

// FileConfig represents the subset of Config that can be persisted to disk.
// Server and Memory settings remain env-only.
type FileConfig struct {
	Jira           *JiraConfig           `json:"jira,omitempty"`
	AI             *AIConfig             `json:"ai,omitempty"`
	Lark           *LarkConfig           `json:"lark,omitempty"`
	Slack          *SlackConfig          `json:"slack,omitempty"`
	Discord        *DiscordConfig        `json:"discord,omitempty"`
	Telegram       *TelegramConfig       `json:"telegram,omitempty"`
	Teams          *TeamsConfig          `json:"teams,omitempty"`
	Email          *EmailConfig          `json:"email,omitempty"`
	Confluence     *ConfluenceConfig     `json:"confluence,omitempty"`
	GoogleCalendar *GoogleCalendarConfig `json:"google_calendar,omitempty"`
	GitHub         *GitHubConfig         `json:"github,omitempty"`
	GitLab         *GitLabConfig         `json:"gitlab,omitempty"`
	Notion         *NotionConfig         `json:"notion,omitempty"`
	Clockify       *ClockifyConfig       `json:"clockify,omitempty"`
	Linear         *LinearConfig         `json:"linear,omitempty"`
	PagerDuty      *PagerDutyConfig      `json:"pagerduty,omitempty"`
	GoogleSheets   *GoogleSheetsConfig   `json:"google_sheets,omitempty"`
}

// ConfigFilePath returns the path to the config JSON file.
func ConfigFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".zara-jira-mcp", "config.json")
}

func loadFileConfig() *FileConfig {
	data, err := os.ReadFile(ConfigFilePath())
	if err != nil {
		return nil
	}
	var fc FileConfig
	if json.Unmarshal(data, &fc) != nil {
		return nil
	}
	return &fc
}

// envOrFile returns env if non-empty, else file value.
func envOrFile(envVal, fileVal string) string {
	if envVal != "" {
		return envVal
	}
	return fileVal
}

func Load() (*Config, error) {
	fc := loadFileConfig()
	if fc == nil {
		fc = &FileConfig{}
	}
	// Helper to get file value or empty
	fJira := fc.Jira
	if fJira == nil {
		fJira = &JiraConfig{}
	}
	fAI := fc.AI
	if fAI == nil {
		fAI = &AIConfig{}
	}
	fLark := fc.Lark
	if fLark == nil {
		fLark = &LarkConfig{}
	}
	fSlack := fc.Slack
	if fSlack == nil {
		fSlack = &SlackConfig{}
	}
	fDiscord := fc.Discord
	if fDiscord == nil {
		fDiscord = &DiscordConfig{}
	}
	fTelegram := fc.Telegram
	if fTelegram == nil {
		fTelegram = &TelegramConfig{}
	}
	fTeams := fc.Teams
	if fTeams == nil {
		fTeams = &TeamsConfig{}
	}
	fEmail := fc.Email
	if fEmail == nil {
		fEmail = &EmailConfig{}
	}
	fConfluence := fc.Confluence
	if fConfluence == nil {
		fConfluence = &ConfluenceConfig{}
	}
	fCal := fc.GoogleCalendar
	if fCal == nil {
		fCal = &GoogleCalendarConfig{}
	}
	fGH := fc.GitHub
	if fGH == nil {
		fGH = &GitHubConfig{}
	}
	fGL := fc.GitLab
	if fGL == nil {
		fGL = &GitLabConfig{}
	}
	fNotion := fc.Notion
	if fNotion == nil {
		fNotion = &NotionConfig{}
	}
	fClockify := fc.Clockify
	if fClockify == nil {
		fClockify = &ClockifyConfig{}
	}
	fLinear := fc.Linear
	if fLinear == nil {
		fLinear = &LinearConfig{}
	}
	fPD := fc.PagerDuty
	if fPD == nil {
		fPD = &PagerDutyConfig{}
	}
	fSheets := fc.GoogleSheets
	if fSheets == nil {
		fSheets = &GoogleSheetsConfig{}
	}

	cfg := &Config{
		Server: ServerConfig{
			Transport:        os.Getenv("MCP_TRANSPORT"),
			Port:             os.Getenv("MCP_PORT"),
			Profile:          os.Getenv("PM_PROFILE"),
			DashboardEnabled: os.Getenv("MCP_DASHBOARD") == "true",
			DashboardPort:    os.Getenv("MCP_DASHBOARD_PORT"),
		},
		Jira: JiraConfig{
			BaseURL: strings.TrimRight(envOrFile(os.Getenv("JIRA_BASE_URL"), fJira.BaseURL), "/"),
			Email:   envOrFile(os.Getenv("JIRA_EMAIL"), fJira.Email),
			Token:   envOrFile(os.Getenv("JIRA_API_TOKEN"), fJira.Token),
		},
		AI: AIConfig{
			BaseURL: envOrFile(os.Getenv("JIRA_AI_BASE_URL"), fAI.BaseURL),
			APIKey:  envOrFile(os.Getenv("JIRA_AI_API_KEY"), fAI.APIKey),
			Model:   envOrFile(os.Getenv("JIRA_AI_MODEL"), fAI.Model),
		},
		Lark: LarkConfig{
			WebhookURL: envOrFile(os.Getenv("JIRA_LARK_WEBHOOK_URL"), fLark.WebhookURL),
			AppID:      envOrFile(os.Getenv("LARK_APP_ID"), fLark.AppID),
			AppSecret:  envOrFile(os.Getenv("LARK_APP_SECRET"), fLark.AppSecret),
			ChatID:     envOrFile(os.Getenv("LARK_CHAT_ID"), fLark.ChatID),
			BotEnabled:        os.Getenv("LARK_BOT_ENABLED") == "true",
			BotPort:           os.Getenv("LARK_BOT_PORT"),
			VerificationToken: os.Getenv("LARK_VERIFICATION_TOKEN"),
			EncryptKey:        os.Getenv("LARK_ENCRYPT_KEY"),
		},
		Slack: SlackConfig{
			BotToken:       envOrFile(os.Getenv("SLACK_BOT_TOKEN"), fSlack.BotToken),
			DefaultChannel: envOrFile(os.Getenv("SLACK_DEFAULT_CHANNEL"), fSlack.DefaultChannel),
			WebhookURL:     envOrFile(os.Getenv("SLACK_WEBHOOK_URL"), fSlack.WebhookURL),
		},
		Discord: DiscordConfig{
			BotToken:   envOrFile(os.Getenv("DISCORD_BOT_TOKEN"), fDiscord.BotToken),
			ChannelID:  envOrFile(os.Getenv("DISCORD_CHANNEL_ID"), fDiscord.ChannelID),
			WebhookURL: envOrFile(os.Getenv("DISCORD_WEBHOOK_URL"), fDiscord.WebhookURL),
		},
		Telegram: TelegramConfig{
			BotToken: envOrFile(os.Getenv("TELEGRAM_BOT_TOKEN"), fTelegram.BotToken),
			ChatID:   envOrFile(os.Getenv("TELEGRAM_CHAT_ID"), fTelegram.ChatID),
		},
		Teams: TeamsConfig{
			WebhookURL: envOrFile(os.Getenv("TEAMS_WEBHOOK_URL"), fTeams.WebhookURL),
		},
		Email: EmailConfig{
			SMTPHost: envOrFile(os.Getenv("EMAIL_SMTP_HOST"), fEmail.SMTPHost),
			SMTPPort: envOrFile(os.Getenv("EMAIL_SMTP_PORT"), fEmail.SMTPPort),
			Username: envOrFile(os.Getenv("EMAIL_USERNAME"), fEmail.Username),
			Password: envOrFile(os.Getenv("EMAIL_PASSWORD"), fEmail.Password),
			From:     envOrFile(os.Getenv("EMAIL_FROM"), fEmail.From),
		},
		Confluence: ConfluenceConfig{
			BaseURL: strings.TrimRight(envOrFile(os.Getenv("CONFLUENCE_BASE_URL"), fConfluence.BaseURL), "/"),
			Email:   envOrFile(os.Getenv("CONFLUENCE_EMAIL"), fConfluence.Email),
			Token:   envOrFile(os.Getenv("CONFLUENCE_API_TOKEN"), fConfluence.Token),
		},
		Memory: MemoryConfig{
			DBPath: os.Getenv("PM_MEMORY_DB_PATH"),
		},
		GoogleCalendar: GoogleCalendarConfig{
			APIKey:     envOrFile(os.Getenv("GOOGLE_CALENDAR_API_KEY"), fCal.APIKey),
			CalendarID: envOrFile(os.Getenv("GOOGLE_CALENDAR_ID"), fCal.CalendarID),
		},
		GitHub: GitHubConfig{
			Token: envOrFile(os.Getenv("GITHUB_TOKEN"), fGH.Token),
			Owner: envOrFile(os.Getenv("GITHUB_OWNER"), fGH.Owner),
			Repo:  envOrFile(os.Getenv("GITHUB_REPO"), fGH.Repo),
		},
		GitLab: GitLabConfig{
			Token:     envOrFile(os.Getenv("GITLAB_TOKEN"), fGL.Token),
			BaseURL:   envOrFile(os.Getenv("GITLAB_BASE_URL"), fGL.BaseURL),
			ProjectID: envOrFile(os.Getenv("GITLAB_PROJECT_ID"), fGL.ProjectID),
		},
		Notion: NotionConfig{
			APIKey:     envOrFile(os.Getenv("NOTION_API_KEY"), fNotion.APIKey),
			DatabaseID: envOrFile(os.Getenv("NOTION_DATABASE_ID"), fNotion.DatabaseID),
		},
		Clockify: ClockifyConfig{
			APIKey:      envOrFile(os.Getenv("CLOCKIFY_API_KEY"), fClockify.APIKey),
			WorkspaceID: envOrFile(os.Getenv("CLOCKIFY_WORKSPACE_ID"), fClockify.WorkspaceID),
		},
		Linear: LinearConfig{
			APIKey: envOrFile(os.Getenv("LINEAR_API_KEY"), fLinear.APIKey),
		},
		PagerDuty: PagerDutyConfig{
			APIKey: envOrFile(os.Getenv("PAGERDUTY_API_KEY"), fPD.APIKey),
		},
		GoogleSheets: GoogleSheetsConfig{
			APIKey:        envOrFile(os.Getenv("GOOGLE_SHEETS_API_KEY"), fSheets.APIKey),
			SpreadsheetID: envOrFile(os.Getenv("GOOGLE_SHEETS_SPREADSHEET_ID"), fSheets.SpreadsheetID),
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

	return cfg, nil
}
