package dashboard

import (
	"context"
	"embed"
	"encoding/json"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/aldok10/zara-jira-mcp/config"
)

//go:embed index.html
var static embed.FS

type Server struct {
	cfg    *config.Config
	logger *slog.Logger
	srv    *http.Server
}

func New(cfg *config.Config, logger *slog.Logger) *Server {
	return &Server{cfg: cfg, logger: logger}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", s.index)
	mux.HandleFunc("GET /api/config", s.getConfig)
	mux.HandleFunc("GET /api/status", s.getStatus)
	addr := ":" + s.cfg.Server.DashboardPort
	s.srv = &http.Server{Addr: addr, Handler: mux}
	s.logger.Info("dashboard", "url", "http://localhost"+addr)
	ln, err := net.Listen("tcp", addr)
	if err != nil { return err }
	go s.srv.Serve(ln)
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.srv != nil { return s.srv.Shutdown(ctx) }
	return nil
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	d, _ := static.ReadFile("index.html")
	w.Header().Set("Content-Type", "text/html")
	w.Write(d)
}

func (s *Server) getConfig(w http.ResponseWriter, r *http.Request) {
	out := map[string]any{
		"jira":       map[string]string{"base_url": s.cfg.Jira.BaseURL, "email": s.cfg.Jira.Email, "token": mask(s.cfg.Jira.Token)},
		"ai":         map[string]string{"base_url": s.cfg.AI.BaseURL, "api_key": mask(s.cfg.AI.APIKey), "model": s.cfg.AI.Model},
		"slack":      map[string]string{"bot_token": mask(s.cfg.Slack.BotToken), "channel": s.cfg.Slack.DefaultChannel},
		"lark":       map[string]string{"webhook": mask(s.cfg.Lark.WebhookURL), "app_id": s.cfg.Lark.AppID},
		"discord":    map[string]string{"bot_token": mask(s.cfg.Discord.BotToken), "channel_id": s.cfg.Discord.ChannelID},
		"telegram":   map[string]string{"bot_token": mask(s.cfg.Telegram.BotToken), "chat_id": s.cfg.Telegram.ChatID},
		"teams":      map[string]string{"webhook": mask(s.cfg.Teams.WebhookURL)},
		"email":      map[string]string{"host": s.cfg.Email.SMTPHost, "port": s.cfg.Email.SMTPPort, "user": s.cfg.Email.Username},
		"confluence": map[string]string{"base_url": s.cfg.Confluence.BaseURL, "token": mask(s.cfg.Confluence.Token)},
		"github":     map[string]string{"token": mask(s.cfg.GitHub.Token), "owner": s.cfg.GitHub.Owner, "repo": s.cfg.GitHub.Repo},
		"notion":     map[string]string{"api_key": mask(s.cfg.Notion.APIKey)},
		"calendar":   map[string]string{"api_key": mask(s.cfg.GoogleCalendar.APIKey)},
		"linear":     map[string]string{"api_key": mask(s.cfg.Linear.APIKey)},
		"pagerduty":  map[string]string{"api_key": mask(s.cfg.PagerDuty.APIKey)},
		"clockify":   map[string]string{"api_key": mask(s.cfg.Clockify.APIKey)},
		"sheets":     map[string]string{"api_key": mask(s.cfg.GoogleSheets.APIKey)},
		"server":     map[string]string{"transport": s.cfg.Server.Transport, "port": s.cfg.Server.Port, "profile": s.cfg.Server.Profile},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

type status struct {
	Name       string `json:"name"`
	Configured bool   `json:"configured"`
}

func (s *Server) getStatus(w http.ResponseWriter, r *http.Request) {
	checks := []status{
		{"Jira", s.cfg.Jira.Token != ""},
		{"AI", s.cfg.AI.APIKey != ""},
		{"Slack", s.cfg.Slack.BotToken != "" || s.cfg.Slack.WebhookURL != ""},
		{"Lark", s.cfg.Lark.WebhookURL != "" || s.cfg.Lark.AppID != ""},
		{"Discord", s.cfg.Discord.BotToken != ""},
		{"Telegram", s.cfg.Telegram.BotToken != ""},
		{"Teams", s.cfg.Teams.WebhookURL != ""},
		{"Email", s.cfg.Email.SMTPHost != ""},
		{"Confluence", s.cfg.Confluence.Token != ""},
		{"GitHub", s.cfg.GitHub.Token != ""},
		{"Notion", s.cfg.Notion.APIKey != ""},
		{"Calendar", s.cfg.GoogleCalendar.APIKey != ""},
		{"Linear", s.cfg.Linear.APIKey != ""},
		{"PagerDuty", s.cfg.PagerDuty.APIKey != ""},
		{"Clockify", s.cfg.Clockify.APIKey != ""},
		{"Sheets", s.cfg.GoogleSheets.APIKey != ""},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checks)
}

func mask(s string) string {
	if s == "" { return "" }
	if len(s) <= 4 { return "****" }
	return strings.Repeat("*", len(s)-4) + s[len(s)-4:]
}
