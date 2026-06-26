package dashboard

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/aldok10/zara-jira-mcp/config"
)

// Store manages config persistence to JSON file.
type Store struct {
	mu   sync.RWMutex
	cfg  *config.Config
	path string
}

func NewStore(cfg *config.Config) *Store {
	return &Store{
		cfg:  cfg,
		path: config.ConfigFilePath(),
	}
}

func (s *Store) Config() *config.Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.cfg
}

// MaskedConfig returns the config with sensitive fields masked.
func (s *Store) MaskedConfig() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Marshal to JSON then unmarshal to map for field iteration
	data, _ := json.Marshal(s.cfg)
	var m map[string]any
	json.Unmarshal(data, &m)

	// Remove server/memory from response (env-only)
	delete(m, "Server")
	delete(m, "Memory")

	// Mask secrets
	for section, val := range m {
		if obj, ok := val.(map[string]any); ok {
			for key, v := range obj {
				if str, ok := v.(string); ok && isSensitiveKey(key) && len(str) > 4 {
					obj[key] = "****" + str[len(str)-4:]
				}
			}
			m[section] = obj
		}
	}
	return m
}

// Update applies a partial FileConfig update, saves to disk, and updates in-memory config.
func (s *Store) Update(fc *config.FileConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Apply non-empty fields from fc to the live config
	if fc.Jira != nil {
		applyStr(&s.cfg.Jira.BaseURL, fc.Jira.BaseURL)
		applyStr(&s.cfg.Jira.Email, fc.Jira.Email)
		applyStr(&s.cfg.Jira.Token, fc.Jira.Token)
	}
	if fc.AI != nil {
		applyStr(&s.cfg.AI.BaseURL, fc.AI.BaseURL)
		applyStr(&s.cfg.AI.APIKey, fc.AI.APIKey)
		applyStr(&s.cfg.AI.Model, fc.AI.Model)
	}
	if fc.Lark != nil {
		applyStr(&s.cfg.Lark.WebhookURL, fc.Lark.WebhookURL)
		applyStr(&s.cfg.Lark.AppID, fc.Lark.AppID)
		applyStr(&s.cfg.Lark.AppSecret, fc.Lark.AppSecret)
		applyStr(&s.cfg.Lark.ChatID, fc.Lark.ChatID)
	}
	if fc.Slack != nil {
		applyStr(&s.cfg.Slack.BotToken, fc.Slack.BotToken)
		applyStr(&s.cfg.Slack.DefaultChannel, fc.Slack.DefaultChannel)
		applyStr(&s.cfg.Slack.WebhookURL, fc.Slack.WebhookURL)
	}
	if fc.Discord != nil {
		applyStr(&s.cfg.Discord.BotToken, fc.Discord.BotToken)
		applyStr(&s.cfg.Discord.ChannelID, fc.Discord.ChannelID)
		applyStr(&s.cfg.Discord.WebhookURL, fc.Discord.WebhookURL)
	}
	if fc.Telegram != nil {
		applyStr(&s.cfg.Telegram.BotToken, fc.Telegram.BotToken)
		applyStr(&s.cfg.Telegram.ChatID, fc.Telegram.ChatID)
	}
	if fc.Teams != nil {
		applyStr(&s.cfg.Teams.WebhookURL, fc.Teams.WebhookURL)
	}
	if fc.Email != nil {
		applyStr(&s.cfg.Email.SMTPHost, fc.Email.SMTPHost)
		applyStr(&s.cfg.Email.SMTPPort, fc.Email.SMTPPort)
		applyStr(&s.cfg.Email.Username, fc.Email.Username)
		applyStr(&s.cfg.Email.Password, fc.Email.Password)
		applyStr(&s.cfg.Email.From, fc.Email.From)
	}
	if fc.Confluence != nil {
		applyStr(&s.cfg.Confluence.BaseURL, fc.Confluence.BaseURL)
		applyStr(&s.cfg.Confluence.Email, fc.Confluence.Email)
		applyStr(&s.cfg.Confluence.Token, fc.Confluence.Token)
	}
	if fc.GoogleCalendar != nil {
		applyStr(&s.cfg.GoogleCalendar.APIKey, fc.GoogleCalendar.APIKey)
		applyStr(&s.cfg.GoogleCalendar.CalendarID, fc.GoogleCalendar.CalendarID)
	}
	if fc.GitHub != nil {
		applyStr(&s.cfg.GitHub.Token, fc.GitHub.Token)
		applyStr(&s.cfg.GitHub.Owner, fc.GitHub.Owner)
		applyStr(&s.cfg.GitHub.Repo, fc.GitHub.Repo)
	}
	if fc.GitLab != nil {
		applyStr(&s.cfg.GitLab.Token, fc.GitLab.Token)
		applyStr(&s.cfg.GitLab.BaseURL, fc.GitLab.BaseURL)
		applyStr(&s.cfg.GitLab.ProjectID, fc.GitLab.ProjectID)
	}
	if fc.Notion != nil {
		applyStr(&s.cfg.Notion.APIKey, fc.Notion.APIKey)
		applyStr(&s.cfg.Notion.DatabaseID, fc.Notion.DatabaseID)
	}
	if fc.Clockify != nil {
		applyStr(&s.cfg.Clockify.APIKey, fc.Clockify.APIKey)
		applyStr(&s.cfg.Clockify.WorkspaceID, fc.Clockify.WorkspaceID)
	}
	if fc.Linear != nil {
		applyStr(&s.cfg.Linear.APIKey, fc.Linear.APIKey)
	}
	if fc.PagerDuty != nil {
		applyStr(&s.cfg.PagerDuty.APIKey, fc.PagerDuty.APIKey)
	}
	if fc.GoogleSheets != nil {
		applyStr(&s.cfg.GoogleSheets.APIKey, fc.GoogleSheets.APIKey)
		applyStr(&s.cfg.GoogleSheets.SpreadsheetID, fc.GoogleSheets.SpreadsheetID)
	}

	return s.saveLocked()
}

func (s *Store) saveLocked() error {
	fc := &config.FileConfig{
		Jira:           &s.cfg.Jira,
		AI:             &s.cfg.AI,
		Lark:           &s.cfg.Lark,
		Slack:          &s.cfg.Slack,
		Discord:        &s.cfg.Discord,
		Telegram:       &s.cfg.Telegram,
		Teams:          &s.cfg.Teams,
		Email:          &s.cfg.Email,
		Confluence:     &s.cfg.Confluence,
		GoogleCalendar: &s.cfg.GoogleCalendar,
		GitHub:         &s.cfg.GitHub,
		GitLab:         &s.cfg.GitLab,
		Notion:         &s.cfg.Notion,
		Clockify:       &s.cfg.Clockify,
		Linear:         &s.cfg.Linear,
		PagerDuty:      &s.cfg.PagerDuty,
		GoogleSheets:   &s.cfg.GoogleSheets,
	}

	dir := filepath.Dir(s.path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(fc, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o600)
}

func applyStr(dst *string, src string) {
	if src != "" {
		*dst = src
	}
}

func isSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	for _, s := range []string{"token", "key", "secret", "password", "webhook"} {
		if strings.Contains(lower, s) {
			return true
		}
	}
	return false
}
