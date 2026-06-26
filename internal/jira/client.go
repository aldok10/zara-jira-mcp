package jira

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	jira "github.com/felixgeelhaar/jirasdk"
	"github.com/felixgeelhaar/jirasdk/core/agile"
	"github.com/felixgeelhaar/jirasdk/core/issue"
	"github.com/felixgeelhaar/jirasdk/core/search"

	"github.com/aldok10/zara-jira-mcp/config"
	domain "github.com/aldok10/zara-jira-mcp/domain/jira"
)

// RestClient wraps jirasdk.Client and implements domain.Client.
type RestClient struct {
	sdk     *jira.Client
	baseURL string
	email   string
	token   string
	http    *http.Client
}

func NewRestClient(cfg *config.Config) (*RestClient, error) {
	client, err := jira.NewClient(
		jira.WithBaseURL(cfg.Jira.BaseURL),
		jira.WithAPIToken(cfg.Jira.Email, cfg.Jira.Token),
		jira.WithTimeout(30*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("create jira client: %w", err)
	}
	return &RestClient{
		sdk:     client,
		baseURL: cfg.Jira.BaseURL,
		email:   cfg.Jira.Email,
		token:   cfg.Jira.Token,
		http:    &http.Client{Timeout: 30 * time.Second},
	}, nil
}

func (c *RestClient) SearchIssues(ctx context.Context, jql string, maxResults int) (*domain.SearchResult, error) {
	if maxResults <= 0 {
		maxResults = 50
	}
	result, err := c.sdk.Search.SearchJQL(ctx, &search.SearchJQLOptions{
		JQL:        jql,
		Fields:     []string{"summary", "description", "status", "priority", "issuetype", "assignee", "reporter", "labels", "created", "updated", "sprint"},
		MaxResults: maxResults,
	})
	if err != nil {
		return nil, err
	}
	out := &domain.SearchResult{MaxResults: result.MaxResults}
	for _, i := range result.Issues {
		out.Issues = append(out.Issues, mapIssue(i))
	}
	out.Total = len(out.Issues)
	return out, nil
}

func (c *RestClient) GetIssue(ctx context.Context, key string) (*domain.Issue, error) {
	i, err := c.sdk.Issue.Get(ctx, key, nil)
	if err != nil {
		return nil, err
	}
	mapped := mapIssue(i)
	return &mapped, nil
}

func (c *RestClient) GetBoards(ctx context.Context) ([]domain.Board, error) {
	boards, err := c.sdk.Agile.GetBoards(ctx, nil)
	if err != nil {
		return nil, err
	}
	out := make([]domain.Board, len(boards))
	for i, b := range boards {
		out[i] = domain.Board{ID: int(b.ID), Name: b.Name, Type: b.Type}
	}
	return out, nil
}

func (c *RestClient) GetActiveSprints(ctx context.Context, boardID int) ([]domain.Sprint, error) {
	sprints, err := c.sdk.Agile.GetBoardSprints(ctx, int64(boardID), &agile.SprintsOptions{
		State: "active",
	})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Sprint, len(sprints))
	for i, s := range sprints {
		out[i] = domain.Sprint{ID: int(s.ID), Name: s.Name, State: s.State, Goal: s.Goal}
	}
	return out, nil
}

func (c *RestClient) GetSprintIssues(ctx context.Context, sprintID int) ([]domain.Issue, error) {
	jql := fmt.Sprintf("sprint = %d ORDER BY status ASC", sprintID)
	result, err := c.sdk.Search.SearchJQL(ctx, &search.SearchJQLOptions{
		JQL:        jql,
		Fields:     []string{"summary", "description", "status", "priority", "issuetype", "assignee", "reporter", "labels", "created", "updated"},
		MaxResults: 100,
	})
	if err != nil {
		return nil, err
	}
	out := make([]domain.Issue, 0, len(result.Issues))
	for _, i := range result.Issues {
		out = append(out, mapIssue(i))
	}
	return out, nil
}

func (c *RestClient) CreateIssue(ctx context.Context, input *domain.CreateIssueInput) (*domain.Issue, error) {
	fields := &issue.IssueFields{
		Project:   &issue.Project{Key: input.Project},
		Summary:   input.Summary,
		IssueType: &issue.IssueType{Name: input.IssueType},
		Labels:    input.Labels,
	}
	if input.Description != "" {
		fields.SetDescriptionText(input.Description)
	}
	if input.Priority != "" {
		fields.Priority = &issue.Priority{Name: input.Priority}
	}
	if input.Assignee != "" {
		fields.Assignee = &issue.User{AccountID: input.Assignee}
	}

	created, err := c.sdk.Issue.Create(ctx, &issue.CreateInput{Fields: fields})
	if err != nil {
		return nil, err
	}
	mapped := mapIssue(created)
	return &mapped, nil
}

func (c *RestClient) AddComment(ctx context.Context, issueKey, body string) error {
	path := fmt.Sprintf("%s/rest/api/3/issue/%s/comment", c.baseURL, issueKey)
	payload := map[string]any{
		"body": map[string]any{
			"type":    "doc",
			"version": 1,
			"content": []map[string]any{
				{"type": "paragraph", "content": []map[string]any{
					{"type": "text", "text": body},
				}},
			},
		},
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("add comment failed %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *RestClient) TransitionIssue(ctx context.Context, issueKey, transitionID string) error {
	return c.sdk.Issue.DoTransition(ctx, issueKey, &issue.TransitionInput{
		Transition: &issue.Transition{ID: transitionID},
	})
}

func (c *RestClient) GetTransitions(ctx context.Context, issueKey string) ([]domain.Transition, error) {
	path := fmt.Sprintf("%s/rest/api/3/issue/%s/transitions", c.baseURL, issueKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("get transitions failed %d: %s", resp.StatusCode, string(b))
	}

	var result struct {
		Transitions []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"transitions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	out := make([]domain.Transition, len(result.Transitions))
	for i, t := range result.Transitions {
		out[i] = domain.Transition{ID: t.ID, Name: t.Name}
	}
	return out, nil
}

func mapIssue(raw *issue.Issue) domain.Issue {
	if raw == nil {
		return domain.Issue{}
	}
	return domain.Issue{
		Key:         raw.Key,
		Summary:     raw.GetSummary(),
		Description: raw.GetDescriptionText(),
		Status:      raw.GetStatusName(),
		Priority:    raw.GetPriorityName(),
		Type:        raw.GetIssueTypeName(),
		Assignee:    raw.GetAssigneeName(),
		Reporter:    raw.GetReporterName(),
		Labels:      raw.GetLabels(),
		Created:     raw.GetCreatedTime(),
		Updated:     raw.GetUpdatedTime(),
	}
}

var _ domain.Client = (*RestClient)(nil)
