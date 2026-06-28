// Package validate provides input validation for MCP tool parameters.
//
// All validators return nil on success or an error describing what's wrong.
// Use them before passing user-supplied values to downstream services.
//
// Guiding principles:
//   - Reject early: catch bad input before any API/DB calls.
//   - Actionable errors: tell the user what format is expected, not just "bad input".
//   - Stdlib only: no external validation libraries.
//   - No panics: every function returns an error.
package validate

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	issueKeyRE    = regexp.MustCompile(`^[A-Z][A-Z0-9_]*-\d+$`)
	internalIPRE  = regexp.MustCompile(`(?:10\.\d{1,3}\.\d{1,3}\.\d{1,3}|127\.\d{1,3}\.\d{1,3}\.\d{1,3}|172\.(?:1[6-9]|2[0-9]|3[0-1])\.\d{1,3}\.\d{1,3}|192\.168\.\d{1,3}\.\d{1,3}|169\.254\.\d{1,3}\.\d{1,3})`)
)

// common Jira field error messages.
const (
	ErrRequired = "is required"
)

// Required checks that s is non-empty.
func Required(name, s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("%s %s", name, ErrRequired)
	}
	return nil
}

// IssueKey validates a Jira issue key (e.g. PROJ-123, ENGINE-42).
func IssueKey(key string) error {
	if err := Required("issue key", key); err != nil {
		return err
	}
	if !issueKeyRE.MatchString(key) {
		return fmt.Errorf("invalid issue key %q: expected format PROJECT-123 (e.g. PROJ-42)", key)
	}
	return nil
}

// IssueKeys validates a comma-separated list of issue keys.
func IssueKeys(keys string) ([]string, error) {
	if err := Required("issue keys", keys); err != nil {
		return nil, err
	}
	parts := strings.Split(keys, ",")
	seen := make(map[string]bool, len(parts))
	for i, k := range parts {
		k = strings.TrimSpace(k)
		if k == "" {
			return nil, fmt.Errorf("issue key at position %d is empty", i+1)
		}
		if !issueKeyRE.MatchString(k) {
			return nil, fmt.Errorf("invalid issue key %q at position %d: expected format PROJECT-123", k, i+1)
		}
		if seen[k] {
			return nil, fmt.Errorf("duplicate issue key: %s", k)
		}
		seen[k] = true
		parts[i] = k
	}
	return parts, nil
}

// ProjectKey validates a Jira project key (e.g. PROJ, ENGINE).
func ProjectKey(key string) error {
	if err := Required("project key", key); err != nil {
		return err
	}
	if len(key) < 2 || len(key) > 10 {
		return fmt.Errorf("project key %q must be 2–10 characters", key)
	}
	for _, r := range key {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return fmt.Errorf("project key %q must contain only letters and digits", key)
		}
	}
	return nil
}

// BoardID validates that id is positive.
func BoardID(id int) error {
	if id <= 0 {
		return errors.New("board_id must be a positive number")
	}
	return nil
}

// SprintID validates that id is positive.
func SprintID(id int) error {
	if id <= 0 {
		return errors.New("sprint_id must be a positive number")
	}
	return nil
}

// JQL performs a basic safety check on a JQL query.
// It does NOT validate JQL syntax — only blocks destructive patterns.
func JQL(jql string) error {
	if err := Required("JQL query", jql); err != nil {
		return err
	}
	lower := strings.ToLower(jql)
	for _, kw := range []string{"delete", "drop ", "truncate", "alter ", ";--", "';"} {
		if strings.Contains(lower, kw) {
			return fmt.Errorf("JQL contains forbidden pattern: %q", kw)
		}
	}
	return nil
}

// Email validates a basic email format.
func Email(email string) error {
	if err := Required("email", email); err != nil {
		return err
	}
	at := strings.LastIndex(email, "@")
	if at < 1 || at >= len(email)-4 {
		return fmt.Errorf("invalid email address: missing or misplaced '@'")
	}
	domain := email[at+1:]
	dot := strings.LastIndex(domain, ".")
	if dot < 1 || dot >= len(domain)-1 {
		return fmt.Errorf("invalid email address: incomplete domain")
	}
	return nil
}

// URL validates that s is a plausible http/https URL.
func URL(s string) error {
	if err := Required("URL", s); err != nil {
		return err
	}
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		return fmt.Errorf("URL must start with http:// or https://")
	}
	if len(s) < 10 {
		return fmt.Errorf("URL is too short")
	}
	return nil
}

// JiraBaseURL validates a Jira instance URL with security checks.
func JiraBaseURL(s string) error {
	if err := Required("Jira base URL", s); err != nil {
		return err
	}
	if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
		return fmt.Errorf("Jira URL must start with http:// or https://")
	}
	if strings.HasPrefix(s, "http://") {
		return fmt.Errorf("Jira URL must use https:// (http is insecure)")
	}
	return nil
}

// Priority validates Jira priority values.
func Priority(p string) error {
	switch p {
	case "", "Highest", "High", "Medium", "Low", "Lowest":
		return nil
	default:
		return fmt.Errorf("invalid priority %q: expected one of: Highest, High, Medium, Low, Lowest", p)
	}
}

// IssueType validates Jira issue type names.
func IssueType(t string) error {
	switch t {
	case "", "Task", "Bug", "Story", "Epic", "Subtask", "Improvement", "New Feature":
		return nil
	default:
		// Issue types are customizable in Jira, so only warn.
		return nil
	}
}

// TransitionID validates a transition ID format.
func TransitionID(id string) error {
	if err := Required("transition_id", id); err != nil {
		return err
	}
	return nil
}

// AccountID validates a Jira account ID format (typically a hash).
func AccountID(id string) error {
	if err := Required("account_id", id); err != nil {
		return err
	}
	if len(id) < 10 {
		return fmt.Errorf("account_id looks too short (%d chars)", len(id))
	}
	return nil
}

// Labels validates comma-separated labels.
func Labels(s string) ([]string, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	for i, l := range parts {
		l = strings.TrimSpace(l)
		if l == "" {
			return nil, fmt.Errorf("label at position %d is empty", i+1)
		}
		parts[i] = l
	}
	return parts, nil
}

// SlackChannel validates a Slack channel ID or name.
// Accepts channel IDs (C...) or channel names (#...).
func SlackChannel(ch string) error {
	if ch == "" {
		return nil // empty means use default
	}
	if len(ch) < 2 {
		return fmt.Errorf("channel %q is too short", ch)
	}
	if ch[0] == '#' {
		return nil // channel name like #general
	}
	if ch[0] == 'C' || ch[0] == 'G' || ch[0] == 'D' {
		for _, r := range ch[1:] {
			if !(r >= 'A' && r <= 'Z') && !(r >= '0' && r <= '9') {
				return fmt.Errorf("invalid Slack channel ID %q: expected alphanumeric", ch)
			}
		}
		return nil
	}
	return fmt.Errorf("invalid Slack channel %q: expected #name or C... ID", ch)
}

// DiscordChannel validates a Discord snowflake channel ID.
func DiscordChannel(ch string) error {
	if ch == "" {
		return nil
	}
	for _, r := range ch {
		if r < '0' || r > '9' {
			return fmt.Errorf("invalid Discord channel %q: must be numeric snowflake ID", ch)
		}
	}
	if len(ch) < 10 {
		return fmt.Errorf("Discord channel ID %q is too short (min 10 digits)", ch)
	}
	return nil
}
