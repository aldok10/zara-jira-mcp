package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aldok10/zara-jira-mcp/config"
)

// IntegrationStatus represents the connectivity status of one integration.
type IntegrationStatus struct {
	Configured bool   `json:"configured"`
	Connected  bool   `json:"connected"`
	Error      string `json:"error,omitempty"`
}

// CheckAllStatus runs lightweight connectivity checks for all integrations.
func CheckAllStatus(cfg *config.Config) map[string]*IntegrationStatus {
	results := make(map[string]*IntegrationStatus)

	results["jira"] = checkHTTP("GET", cfg.Jira.BaseURL+"/rest/api/2/myself", "", cfg.Jira.Email, cfg.Jira.Token, cfg.Jira.BaseURL != "")
	results["github"] = checkHTTP("GET", "https://api.github.com/user", cfg.GitHub.Token, "", "", cfg.GitHub.Token != "")
	results["slack"] = checkSlack(cfg.Slack.BotToken)
	results["notion"] = checkHTTP("GET", "https://api.notion.com/v1/users/me", cfg.Notion.APIKey, "", "", cfg.Notion.APIKey != "")
	results["linear"] = checkLinear(cfg.Linear.APIKey)
	results["pagerduty"] = checkHTTP("GET", "https://api.pagerduty.com/abilities", cfg.PagerDuty.APIKey, "", "", cfg.PagerDuty.APIKey != "")
	results["clockify"] = checkHTTP("GET", "https://api.clockify.me/api/v1/user", cfg.Clockify.APIKey, "", "", cfg.Clockify.APIKey != "")
	results["google_calendar"] = checkCalendar(cfg.GoogleCalendar)
	results["google_sheets"] = configuredOnly(cfg.GoogleSheets.APIKey != "")
	results["ai"] = configuredOnly(cfg.AI.APIKey != "")
	results["lark"] = configuredOnly(cfg.Lark.WebhookURL != "" || cfg.Lark.AppID != "")
	results["discord"] = configuredOnly(cfg.Discord.BotToken != "" || cfg.Discord.WebhookURL != "")
	results["telegram"] = configuredOnly(cfg.Telegram.BotToken != "")
	results["teams"] = configuredOnly(cfg.Teams.WebhookURL != "")
	results["email"] = configuredOnly(cfg.Email.SMTPHost != "")
	results["confluence"] = checkHTTP("GET", cfg.Confluence.BaseURL+"/wiki/rest/api/user/current", "", cfg.Confluence.Email, cfg.Confluence.Token, cfg.Confluence.BaseURL != "")
	results["gitlab"] = checkHTTP("GET", gitlabURL(cfg.GitLab.BaseURL)+"/api/v4/user", cfg.GitLab.Token, "", "", cfg.GitLab.Token != "")

	return results
}

func checkHTTP(method, url, bearerToken, basicUser, basicPass string, configured bool) *IntegrationStatus {
	s := &IntegrationStatus{Configured: configured}
	if !configured {
		return s
	}
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		s.Error = err.Error()
		return s
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	} else if basicUser != "" {
		req.SetBasicAuth(basicUser, basicPass)
	}
	if strings.Contains(url, "notion.com") {
		req.Header.Set("Notion-Version", "2022-06-28")
	}
	if strings.Contains(url, "pagerduty.com") {
		req.Header.Set("Authorization", "Token token="+bearerToken)
	}
	resp, err := client.Do(req)
	if err != nil {
		s.Error = err.Error()
		return s
	}
	resp.Body.Close()
	if resp.StatusCode < 300 {
		s.Connected = true
	} else {
		s.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
	return s
}

func checkSlack(token string) *IntegrationStatus {
	s := &IntegrationStatus{Configured: token != ""}
	if !s.Configured {
		return s
	}
	client := &http.Client{Timeout: 5 * time.Second}
	req, _ := http.NewRequest("POST", "https://slack.com/api/auth.test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		s.Error = err.Error()
		return s
	}
	defer resp.Body.Close()
	var result struct {
		OK    bool   `json:"ok"`
		Error string `json:"error"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	if result.OK {
		s.Connected = true
	} else {
		s.Error = result.Error
	}
	return s
}

func checkLinear(apiKey string) *IntegrationStatus {
	s := &IntegrationStatus{Configured: apiKey != ""}
	if !s.Configured {
		return s
	}
	client := &http.Client{Timeout: 5 * time.Second}
	body := strings.NewReader(`{"query":"{ viewer { id } }"}`)
	req, _ := http.NewRequest("POST", "https://api.linear.app/graphql", body)
	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		s.Error = err.Error()
		return s
	}
	resp.Body.Close()
	if resp.StatusCode < 300 {
		s.Connected = true
	} else {
		s.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
	return s
}

func checkCalendar(cfg config.GoogleCalendarConfig) *IntegrationStatus {
	s := &IntegrationStatus{Configured: cfg.APIKey != "" && cfg.CalendarID != ""}
	if !s.Configured {
		return s
	}
	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s?key=%s", cfg.CalendarID, cfg.APIKey)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		s.Error = err.Error()
		return s
	}
	resp.Body.Close()
	if resp.StatusCode < 300 {
		s.Connected = true
	} else {
		s.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
	return s
}

func configuredOnly(configured bool) *IntegrationStatus {
	return &IntegrationStatus{Configured: configured, Connected: configured}
}

func gitlabURL(baseURL string) string {
	if baseURL != "" {
		return strings.TrimRight(baseURL, "/")
	}
	return "https://gitlab.com"
}
