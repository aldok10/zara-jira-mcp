package mcp

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/mcputil"
)

// OnboardConfig holds config status for pm_onboard (no secrets exposed).
type OnboardConfig struct {
	JiraConfigured   bool
	JiraURL          string
	AIConfigured     bool
	AIProvider       string
	HasSlack         bool
	HasLark          bool
	HasDiscord       bool
	HasTelegram      bool
	HasTeams         bool
	HasEmail         bool
	HasGitHub        bool
	HasGitLab        bool
	HasNotion        bool
	HasConfluence    bool
	HasPagerDuty     bool
	MemoryDBPath     string
	MemoryDBExists   bool
	DashboardEnabled bool
	ToolsRegistered  int
}

// NewOnboardConfigFromConfig inspects a config and returns sanitized status.
func NewOnboardConfigFromConfig(cfg *config.Config, toolsCount int) *OnboardConfig {
	oc := &OnboardConfig{
		JiraConfigured:   cfg.Jira.BaseURL != "" && cfg.Jira.Email != "" && cfg.Jira.Token != "",
		JiraURL:          cfg.Jira.BaseURL,
		AIConfigured:     cfg.AI.BaseURL != "" && cfg.AI.APIKey != "",
		AIProvider:       detectProvider(cfg.AI.BaseURL),
		HasSlack:         cfg.Slack.BotToken != "" || cfg.Slack.WebhookURL != "",
		HasLark:          cfg.Lark.WebhookURL != "" || (cfg.Lark.AppID != "" && cfg.Lark.AppSecret != ""),
		HasDiscord:       cfg.Discord.BotToken != "" || cfg.Discord.WebhookURL != "",
		HasTelegram:      cfg.Telegram.BotToken != "",
		HasTeams:         cfg.Teams.WebhookURL != "",
		HasEmail:         cfg.Email.SMTPHost != "" && cfg.Email.Username != "",
		HasGitHub:        cfg.GitHub.Token != "",
		HasGitLab:        cfg.GitLab.Token != "",
		HasNotion:        cfg.Notion.APIKey != "",
		HasConfluence:    cfg.Confluence.BaseURL != "",
		HasPagerDuty:     cfg.PagerDuty.APIKey != "",
		MemoryDBPath:     cfg.Memory.DBPath,
		DashboardEnabled: cfg.Server.DashboardEnabled,
		ToolsRegistered:  toolsCount,
	}
	if cfg.Memory.DBPath != "" {
		if _, err := os.Stat(cfg.Memory.DBPath); err == nil {
			oc.MemoryDBExists = true
		}
	}
	return oc
}

func detectProvider(baseURL string) string {
	if strings.Contains(baseURL, "anthropic") {
		return "Anthropic"
	}
	if strings.Contains(baseURL, "generativelanguage.googleapis.com") {
		return "Gemini"
	}
	if baseURL != "" {
		return "OpenAI-compatible"
	}
	return "not detected"
}

// RegisterBackupTools registers the pm_backup tool.
func RegisterBackupTools(s *server.MCPServer, dbPath string) {
	s.AddTool(
		mcp.NewTool("pm_backup",
			mcp.WithDescription("Export PM memory to JSON backup."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleBackup(ctx, dbPath)
		},
	)
}

func handleBackup(ctx context.Context, dbPath string) (*mcp.CallToolResult, error) {
	if dbPath == "" {
		return mcputil.ErrorResult("PM memory DB not configured. Set PM_MEMORY_DB_PATH."), nil
	}

	db, err := sql.Open("sqlite3", dbPath+"?mode=ro&_journal_mode=WAL")
	if err != nil {
		return mcputil.ErrInternal("open db", err), nil
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		return mcputil.ErrInternal("query tables", err), nil
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return mcputil.ErrInternal("scan table name", err), nil
		}
		tables = append(tables, name)
	}
	if err := rows.Err(); err != nil {
		return mcputil.ErrInternal("rows iteration", err), nil
	}

	if len(tables) == 0 {
		return mcputil.TextResult("{\"tables\":[],\"note\":\"PM memory is empty — no tables found.\"}"), nil
	}

	export := make(map[string]any)
	for _, table := range tables {
		data, err := exportTable(ctx, db, table)
		if err != nil {
			export[table] = map[string]any{"error": err.Error()}
			continue
		}
		export[table] = data
	}

	export["_meta"] = map[string]any{
		"exported_at": "now",
		"table_count": len(tables),
	}

	result, err := json.MarshalIndent(export, "", "  ")
	if err != nil {
		return mcputil.ErrInternal("marshal backup", err), nil
	}

	return mcputil.TextResult(string(result)), nil
}

func exportTable(ctx context.Context, db *sql.DB, table string) (any, error) {
	rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT * FROM %q", table))
	if err != nil {
		return nil, fmt.Errorf("query %s: %w", table, err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("columns %s: %w", table, err)
	}

	var records []map[string]any
	for rows.Next() {
		vals := make([]any, len(columns))
		valPtrs := make([]any, len(columns))
		for i := range vals {
			valPtrs[i] = &vals[i]
		}

		if err := rows.Scan(valPtrs...); err != nil {
			return nil, fmt.Errorf("scan %s: %w", table, err)
		}

		record := make(map[string]any, len(columns))
		for i, col := range columns {
			record[col] = formatSQLValue(vals[i])
		}
		records = append(records, record)
	}

	return records, rows.Err()
}

func formatSQLValue(v any) any {
	switch val := v.(type) {
	case []byte:
		return string(val)
	case fmt.Stringer:
		return val.String()
	default:
		return val
	}
}

// BackupResult contains the output of a pm_backup.
type BackupResult struct {
	Tables  map[string]any `json:"tables"`
	Success bool           `json:"success"`
	Error   string         `json:"error,omitempty"`
}

// FormatBackupAsMarkdown formats backup data as readable markdown.
func FormatBackupAsMarkdown(data string) string {
	var sb strings.Builder
	sb.WriteString("## PM Memory Backup\n\n")
	sb.WriteString("```json\n")
	sb.WriteString(data)
	sb.WriteString("\n```")
	return sb.String()
}

// RegisterOnboardTool registers the pm_onboard first-run wizard.
func RegisterOnboardTool(s *server.MCPServer, cfg *config.Config, toolsCount int) {
	onboardCfg := NewOnboardConfigFromConfig(cfg, toolsCount)
	s.AddTool(
		mcp.NewTool("pm_onboard",
			mcp.WithDescription("First-run wizard: detect config, suggest next steps."),
		),
		func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return handleOnboard(ctx, onboardCfg)
		},
	)
}

func handleOnboard(_ context.Context, cfg *OnboardConfig) (*mcp.CallToolResult, error) {
	var sb strings.Builder
	sb.WriteString("# zara-jira-mcp Onboard\n\n")

	// Jira status
	if cfg.JiraConfigured {
		sb.WriteString("[OK] Jira: Configured\n")
		sb.WriteString(fmt.Sprintf("     URL: %s\n", cfg.JiraURL))
	} else {
		sb.WriteString("[  ] Jira: Not configured. Set JIRA_BASE_URL, JIRA_EMAIL, JIRA_API_TOKEN\n")
	}

	// AI status
	if cfg.AIConfigured {
		sb.WriteString(fmt.Sprintf("[OK] AI: Configured (%s)\n", cfg.AIProvider))
	} else {
		sb.WriteString("[  ] AI: Not configured. Set JIRA_AI_BASE_URL, JIRA_AI_API_KEY\n")
	}

	// Memory
	if cfg.MemoryDBPath != "" {
		if cfg.MemoryDBExists {
			sb.WriteString("[OK] PM Memory: Database exists\n")
		} else {
			sb.WriteString("[--] PM Memory: DB path configured but file not yet created\n")
		}
	} else {
		sb.WriteString("[--] PM Memory: Not configured. Set PM_MEMORY_DB_PATH\n")
	}

	// Notification channels
	sb.WriteString("\n## Notification Channels\n")
	channels := 0
	if cfg.HasSlack {
		sb.WriteString("  [OK] Slack\n")
		channels++
	}
	if cfg.HasLark {
		sb.WriteString("  [OK] Lark\n")
		channels++
	}
	if cfg.HasDiscord {
		sb.WriteString("  [OK] Discord\n")
		channels++
	}
	if cfg.HasTelegram {
		sb.WriteString("  [OK] Telegram\n")
		channels++
	}
	if cfg.HasTeams {
		sb.WriteString("  [OK] Microsoft Teams\n")
		channels++
	}
	if cfg.HasEmail {
		sb.WriteString("  [OK] Email\n")
		channels++
	}
	if channels == 0 {
		sb.WriteString("  [--] No notification channels configured\n")
	}

	// Integrations
	sb.WriteString("\n## Integrations\n")
	if cfg.HasGitHub {
		sb.WriteString("  [OK] GitHub\n")
	}
	if cfg.HasGitLab {
		sb.WriteString("  [OK] GitLab\n")
	}
	if cfg.HasNotion {
		sb.WriteString("  [OK] Notion\n")
	}
	if cfg.HasConfluence {
		sb.WriteString("  [OK] Confluence\n")
	}
	if cfg.HasPagerDuty {
		sb.WriteString("  [OK] PagerDuty\n")
	}

	// Tools registered
	sb.WriteString("\n## Status\n")
	sb.WriteString(fmt.Sprintf("  Tools registered: %d\n", cfg.ToolsRegistered))

	// Next steps
	sb.WriteString("\n## Next Steps\n")
	steps := []string{}
	if !cfg.JiraConfigured {
		steps = append(steps, "1. Set JIRA_BASE_URL, JIRA_EMAIL, JIRA_API_TOKEN")
	}
	if !cfg.AIConfigured {
		steps = append(steps, "2. Set JIRA_AI_BASE_URL, JIRA_AI_API_KEY")
	}
	if cfg.JiraConfigured {
		steps = append(steps, "1. Run `jira_boards` to discover your boards")
		steps = append(steps, "2. Run `jira_sprint_summary` with your board_id")
	}
	if len(steps) == 0 {
		sb.WriteString("  All systems ready. Try `pm` to get started.\n")
	} else {
		for _, step := range steps {
			sb.WriteString(fmt.Sprintf("  %s\n", step))
		}
	}

	return mcputil.TextResult(sb.String()), nil
}
