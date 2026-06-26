package config

import (
	"os"
	"testing"
)

func TestLoad_RequiredFields(t *testing.T) {
	// Clear env
	os.Clearenv()

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when JIRA_BASE_URL missing")
	}

	os.Setenv("JIRA_BASE_URL", "https://test.atlassian.net")
	_, err = Load()
	if err == nil {
		t.Fatal("expected error when JIRA_EMAIL missing")
	}

	os.Setenv("JIRA_EMAIL", "test@test.com")
	_, err = Load()
	if err == nil {
		t.Fatal("expected error when JIRA_API_TOKEN missing")
	}

	os.Setenv("JIRA_API_TOKEN", "token123")
	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Jira.BaseURL != "https://test.atlassian.net" {
		t.Errorf("got %s, want https://test.atlassian.net", cfg.Jira.BaseURL)
	}
}

func TestLoad_Defaults(t *testing.T) {
	os.Clearenv()
	os.Setenv("JIRA_BASE_URL", "https://x.atlassian.net/")
	os.Setenv("JIRA_EMAIL", "a@b.com")
	os.Setenv("JIRA_API_TOKEN", "tok")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Trailing slash stripped
	if cfg.Jira.BaseURL != "https://x.atlassian.net" {
		t.Errorf("trailing slash not trimmed: %s", cfg.Jira.BaseURL)
	}
	// Default AI model
	if cfg.AI.Model != "gpt-4o-mini" {
		t.Errorf("default model not set: %s", cfg.AI.Model)
	}
}

func TestLoad_OptionalPlatforms(t *testing.T) {
	os.Clearenv()
	os.Setenv("JIRA_BASE_URL", "https://x.atlassian.net")
	os.Setenv("JIRA_EMAIL", "a@b.com")
	os.Setenv("JIRA_API_TOKEN", "tok")
	os.Setenv("SLACK_BOT_TOKEN", "xoxb-test")
	os.Setenv("DISCORD_BOT_TOKEN", "discord-tok")
	os.Setenv("TELEGRAM_BOT_TOKEN", "tg-tok")
	os.Setenv("TEAMS_WEBHOOK_URL", "https://teams.webhook")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Slack.BotToken != "xoxb-test" {
		t.Error("slack token not loaded")
	}
	if cfg.Discord.BotToken != "discord-tok" {
		t.Error("discord token not loaded")
	}
	if cfg.Telegram.BotToken != "tg-tok" {
		t.Error("telegram token not loaded")
	}
	if cfg.Teams.WebhookURL != "https://teams.webhook" {
		t.Error("teams webhook not loaded")
	}
}
