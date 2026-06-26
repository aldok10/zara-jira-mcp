package jira

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/aldok10/zara-jira-mcp/config"
	domain "github.com/aldok10/zara-jira-mcp/domain/jira"
)

// RestClient implements domain.Client using Jira REST API v3.
type RestClient struct {
	baseURL    string
	httpClient *http.Client
	email      string
	token      string
}

func NewRestClient(cfg *config.Config) *RestClient {
	return &RestClient{
		baseURL:    cfg.Jira.BaseURL,
		email:      cfg.Jira.Email,
		token:      cfg.Jira.Token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *RestClient) doRequest(ctx context.Context, method, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("jira API error %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *RestClient) SearchIssues(ctx context.Context, jql string, maxResults int) (*domain.SearchResult, error) {
	if maxResults <= 0 {
		maxResults = 50
	}

	path := fmt.Sprintf("/rest/api/3/search?jql=%s&maxResults=%d&fields=summary,description,status,priority,issuetype,assignee,reporter,labels,created,updated,sprint",
		url.QueryEscape(jql), maxResults)

	body, err := c.doRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	var raw searchResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parse search response: %w", err)
	}

	result := &domain.SearchResult{
		Total:      raw.Total,
		MaxResults: raw.MaxResults,
	}

	for _, ri := range raw.Issues {
		result.Issues = append(result.Issues, mapIssue(ri))
	}

	return result, nil
}

func (c *RestClient) GetIssue(ctx context.Context, key string) (*domain.Issue, error) {
	path := fmt.Sprintf("/rest/api/3/issue/%s?fields=summary,description,status,priority,issuetype,assignee,reporter,labels,created,updated,sprint", key)

	body, err := c.doRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	var raw rawIssue
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("parse issue response: %w", err)
	}

	issue := mapIssue(raw)
	return &issue, nil
}

func (c *RestClient) GetBoards(ctx context.Context) ([]domain.Board, error) {
	body, err := c.doRequest(ctx, http.MethodGet, "/rest/agile/1.0/board")
	if err != nil {
		return nil, err
	}

	var raw struct {
		Values []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
			Type string `json:"type"`
		} `json:"values"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	boards := make([]domain.Board, len(raw.Values))
	for i, v := range raw.Values {
		boards[i] = domain.Board{ID: v.ID, Name: v.Name, Type: v.Type}
	}
	return boards, nil
}

func (c *RestClient) GetActiveSprints(ctx context.Context, boardID int) ([]domain.Sprint, error) {
	path := fmt.Sprintf("/rest/agile/1.0/board/%d/sprint?state=active", boardID)
	body, err := c.doRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	var raw struct {
		Values []struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			State string `json:"state"`
			Goal  string `json:"goal"`
		} `json:"values"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	sprints := make([]domain.Sprint, len(raw.Values))
	for i, v := range raw.Values {
		sprints[i] = domain.Sprint{ID: v.ID, Name: v.Name, State: v.State, Goal: v.Goal}
	}
	return sprints, nil
}

func (c *RestClient) GetSprintIssues(ctx context.Context, sprintID int) ([]domain.Issue, error) {
	path := fmt.Sprintf("/rest/agile/1.0/sprint/%d/issue?fields=summary,description,status,priority,issuetype,assignee,reporter,labels,created,updated", sprintID)
	body, err := c.doRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}

	var raw searchResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, err
	}

	issues := make([]domain.Issue, len(raw.Issues))
	for i, ri := range raw.Issues {
		issues[i] = mapIssue(ri)
	}
	return issues, nil
}

// JSON mapping types for Jira API responses.

type searchResponse struct {
	Total      int        `json:"total"`
	MaxResults int        `json:"maxResults"`
	Issues     []rawIssue `json:"issues"`
}

type rawIssue struct {
	Key    string `json:"key"`
	Fields struct {
		Summary     string `json:"summary"`
		Description *struct {
			Content []struct {
				Content []struct {
					Text string `json:"text"`
				} `json:"content"`
			} `json:"content"`
		} `json:"description"`
		Status struct {
			Name string `json:"name"`
		} `json:"status"`
		Priority *struct {
			Name string `json:"name"`
		} `json:"priority"`
		IssueType struct {
			Name string `json:"name"`
		} `json:"issuetype"`
		Assignee *struct {
			DisplayName string `json:"displayName"`
		} `json:"assignee"`
		Reporter *struct {
			DisplayName string `json:"displayName"`
		} `json:"reporter"`
		Labels  []string `json:"labels"`
		Created string   `json:"created"`
		Updated string   `json:"updated"`
		Sprint  *struct {
			Name string `json:"name"`
		} `json:"sprint"`
	} `json:"fields"`
}

func mapIssue(raw rawIssue) domain.Issue {
	issue := domain.Issue{
		Key:     raw.Key,
		Summary: raw.Fields.Summary,
		Status:  raw.Fields.Status.Name,
		Type:    raw.Fields.IssueType.Name,
		Labels:  raw.Fields.Labels,
	}

	if raw.Fields.Description != nil {
		for _, block := range raw.Fields.Description.Content {
			for _, inline := range block.Content {
				if inline.Text != "" {
					issue.Description += inline.Text + "\n"
				}
			}
		}
	}

	if raw.Fields.Priority != nil {
		issue.Priority = raw.Fields.Priority.Name
	}
	if raw.Fields.Assignee != nil {
		issue.Assignee = raw.Fields.Assignee.DisplayName
	}
	if raw.Fields.Reporter != nil {
		issue.Reporter = raw.Fields.Reporter.DisplayName
	}
	if raw.Fields.Sprint != nil {
		issue.SprintName = raw.Fields.Sprint.Name
	}

	issue.Created, _ = parseJiraTime(raw.Fields.Created)
	issue.Updated, _ = parseJiraTime(raw.Fields.Updated)

	return issue
}

func parseJiraTime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	// Jira uses ISO 8601 with timezone offset
	return time.Parse("2006-01-02T15:04:05.000-0700", s)
}

// Ensure RestClient implements the domain interface.
var _ domain.Client = (*RestClient)(nil)
