package config

import (
	"fmt"
	"os"
	"strings"
)

// Config holds all application configuration parsed from environment variables.
type Config struct {
	Jira JiraConfig
	AI   AIConfig
	Lark LarkConfig
}

type LarkConfig struct {
	WebhookURL string // JIRA_LARK_WEBHOOK_URL
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
