package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	domain "github.com/aldok10/zara-jira-mcp/modules/jira/domain"
)

// FindUser searches users in Jira.
func (c *RestClient) FindUser(ctx context.Context, query string) ([]domain.User, error) {
	data, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/rest/api/3/user/search?query=%s", url.QueryEscape(query)), nil)
	if err != nil {
		return nil, err
	}
	var users []struct {
		AccountID   string `json:"accountId"`
		DisplayName string `json:"displayName"`
		Email       string `json:"emailAddress"`
	}
	if err := json.Unmarshal(data, &users); err != nil {
		return nil, fmt.Errorf("decode users: %w", err)
	}
	out := make([]domain.User, len(users))
	for i, u := range users {
		out[i] = domain.User{AccountID: u.AccountID, DisplayName: u.DisplayName, Email: u.Email}
	}
	return out, nil
}

// SetEpicLink sets an epic as parent of an issue.
func (c *RestClient) SetEpicLink(ctx context.Context, issueKey, epicKey string) error {
	payload := map[string]any{"fields": map[string]any{"parent": map[string]string{"key": epicKey}}}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	_, _, err = c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/rest/api/3/issue/%s", issueKey), data)
	return err
}

// RemoveEpicLink removes epic link from an issue.
func (c *RestClient) RemoveEpicLink(ctx context.Context, issueKey string) error {
	payload := map[string]any{"fields": map[string]any{"parent": nil}}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	_, _, err = c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/rest/api/3/issue/%s", issueKey), data)
	return err
}

// LinkIssues links two issues together.
func (c *RestClient) LinkIssues(ctx context.Context, inwardKey, outwardKey, linkType string) error {
	payload := map[string]any{
		"type":         map[string]string{"name": linkType},
		"inwardIssue":  map[string]string{"key": inwardKey},
		"outwardIssue": map[string]string{"key": outwardKey},
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	_, _, err = c.doRequest(ctx, http.MethodPost, "/rest/api/3/issueLink", data)
	return err
}

// GetLinkTypes retrieves available link types.
func (c *RestClient) GetLinkTypes(ctx context.Context) ([]domain.LinkType, error) {
	data, _, err := c.doRequest(ctx, http.MethodGet, "/rest/api/3/issueLinkType", nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		IssueLinkTypes []struct {
			Name    string `json:"name"`
			Inward  string `json:"inward"`
			Outward string `json:"outward"`
		} `json:"issueLinkTypes"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode link types: %w", err)
	}
	out := make([]domain.LinkType, len(result.IssueLinkTypes))
	for i, lt := range result.IssueLinkTypes {
		out[i] = domain.LinkType{Name: lt.Name, Inward: lt.Inward, Outward: lt.Outward}
	}
	return out, nil
}

// AddWorklog adds a worklog entry.
func (c *RestClient) AddWorklog(ctx context.Context, issueKey, timeSpent, comment string) error {
	payload := map[string]any{"timeSpent": timeSpent}
	if comment != "" {
		payload["comment"] = map[string]any{
			"type": "doc", "version": 1,
			"content": []map[string]any{
				{"type": "paragraph", "content": []map[string]any{
					{"type": "text", "text": comment},
				}},
			},
		}
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	_, _, err = c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/rest/api/3/issue/%s/worklog", issueKey), data)
	return err
}

// GetWorklogs retrieves worklogs for an issue.
func (c *RestClient) GetWorklogs(ctx context.Context, issueKey string) ([]domain.Worklog, error) {
	data, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/rest/api/3/issue/%s/worklog", issueKey), nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Worklogs []struct {
			Author    struct{ DisplayName string } `json:"author"`
			TimeSpent string                       `json:"timeSpent"`
			Started   string                       `json:"started"`
			Comment   *struct {
				Content []struct {
					Content []struct {
						Text string `json:"text"`
					} `json:"content"`
				} `json:"content"`
			} `json:"comment"`
		} `json:"worklogs"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode worklogs: %w", err)
	}
	out := make([]domain.Worklog, len(result.Worklogs))
	for i, w := range result.Worklogs {
		var commentStr string
		if w.Comment != nil && len(w.Comment.Content) > 0 && len(w.Comment.Content[0].Content) > 0 {
			commentStr = w.Comment.Content[0].Content[0].Text
		}
		out[i] = domain.Worklog{Author: w.Author.DisplayName, TimeSpent: w.TimeSpent, Started: w.Started, Comment: commentStr}
	}
	return out, nil
}

// AddWatcher adds a watcher to an issue.
func (c *RestClient) AddWatcher(ctx context.Context, issueKey, accountID string) error {
	body := accountID
	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	_, _, err = c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/rest/api/3/issue/%s/watchers", issueKey), data)
	return err
}

// GetWatchers retrieves watchers for an issue.
func (c *RestClient) GetWatchers(ctx context.Context, issueKey string) ([]domain.User, error) {
	data, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/rest/api/3/issue/%s/watchers", issueKey), nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Watchers []struct {
			AccountID   string `json:"accountId"`
			DisplayName string `json:"displayName"`
		} `json:"watchers"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode watchers: %w", err)
	}
	out := make([]domain.User, len(result.Watchers))
	for i, w := range result.Watchers {
		out[i] = domain.User{AccountID: w.AccountID, DisplayName: w.DisplayName}
	}
	return out, nil
}

// List of allowed path prefixes for RawRequest. Prevents SSRF-style attacks
// and restricts access to unintended Jira API endpoints.
var rawRequestAllowedPrefixes = []string{
	"/rest/api/3/",
	"/rest/agile/1.0/",
}

// rawRequestSensitivePrefixes blocks destructive operations on critical resources.
// DELETE requests targeting these prefixes are rejected outright.
var rawRequestSensitivePrefixes = []string{
	"/rest/api/3/project",
	"/rest/api/3/user",
	"/rest/api/3/applicationrole",
	"/rest/api/3/configuration",
	"/rest/api/3/webhook",
	"/rest/api/3/field",
}

// rawRequestBlockedMethods are HTTP methods that are never allowed.
var rawRequestBlockedMethods = []string{
	http.MethodDelete,
}

// RawRequest executes a raw HTTP request for advanced operations with security
// restrictions applied:
//   - Only predefined path prefixes are allowed.
//   - DELETE requests to sensitive endpoints are rejected.
//   - Context cancellation is respected.
func (c *RestClient) RawRequest(ctx context.Context, method, path string, body []byte) ([]byte, int, error) {
	if err := validateRawRequest(method, path); err != nil {
		return nil, 0, fmt.Errorf("raw request denied: %w", err)
	}
	return c.doRequest(ctx, method, path, body)
}

// GetBoardConfiguration fetches the full board configuration including column layout
// and status mappings from the Jira Agile API.
func (c *RestClient) GetBoardConfiguration(ctx context.Context, boardID int) (*domain.BoardConfiguration, error) {
	data, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/rest/agile/1.0/board/%d/configuration", boardID), nil)
	if err != nil {
		return nil, fmt.Errorf("get board %d config: %w", boardID, err)
	}

	var cfg domain.BoardConfiguration
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("decode board %d config: %w", boardID, err)
	}
	return &cfg, nil
}

func validateRawRequest(method, path string) error {
	// Check for blocked methods first (e.g. DELETE always blocked).
	for _, blocked := range rawRequestBlockedMethods {
		if method == blocked {
			return fmt.Errorf("method %s is not allowed on raw requests", method)
		}
	}

	// Verify path is under an allowed prefix.
	allowed := false
	for _, prefix := range rawRequestAllowedPrefixes {
		if strings.HasPrefix(path, prefix) {
			allowed = true
			break
		}
	}
	if !allowed {
		return fmt.Errorf("path %q is not in the allowed prefix list", path)
	}

	// Block requests to sensitive resources.
	for _, sens := range rawRequestSensitivePrefixes {
		if strings.HasPrefix(path, sens) {
			return fmt.Errorf("path %q is a sensitive endpoint and cannot be accessed via raw request", path)
		}
	}

	return nil
}
