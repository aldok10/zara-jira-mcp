package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	jira "github.com/felixgeelhaar/jirasdk"
	"github.com/felixgeelhaar/jirasdk/core/agile"
	"github.com/felixgeelhaar/jirasdk/core/issue"
	"github.com/felixgeelhaar/jirasdk/core/search"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
	domain "github.com/aldok10/zara-jira-mcp/modules/jira/domain"
)

type rateLimiter struct {
	mu     sync.Mutex
	tokens int
	max    int
	last   time.Time
}

func (r *rateLimiter) wait() {
	r.mu.Lock()
	defer r.mu.Unlock()
	refill := int(time.Since(r.last).Seconds() / 60.0 * float64(r.max))
	if refill > 0 {
		r.tokens += refill
		if r.tokens > r.max {
			r.tokens = r.max
		}
		r.last = time.Now()
	}
	if r.tokens <= 0 {
		time.Sleep(time.Second)
		r.tokens = 1
	}
	r.tokens--
}

// RestClient wraps jirasdk.Client and implements domain.Client.
type RestClient struct {
	sdk     *jira.Client
	baseURL string
	email   string
	token   string
	http    *http.Client
	limiter *rateLimiter
}

func NewRestClient(cfg *config.Config) (*RestClient, error) {
	// Secure HTTP client configuration with TLS verification
	secureTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		IdleConnTimeout:       30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		// Additional security hardening for enterprise environment
		DisableKeepAlives:     false,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
	}
	secureHTTPClient := &http.Client{
		Transport:  secureTransport,
		Timeout:    30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

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
		http:    secureHTTPClient,
		limiter: &rateLimiter{max: 60, tokens: 60, last: time.Now()},
	}, nil
}

func (c *RestClient) SearchIssues(ctx context.Context, jql string, maxResults int, startAt int) (*domain.SearchResult, error) {
	if maxResults <= 0 {
		maxResults = 50
	}
	result, err := c.sdk.Search.SearchJQL(ctx, &search.SearchJQLOptions{
		JQL:        jql,
		Fields:     []string{"summary", "description", "status", "priority", "issuetype", "assignee", "reporter", "labels", "created", "updated", "sprint", "story_points", "customfield_10016", "customfield_10028"},
		MaxResults: maxResults,
	})
	if err != nil {
		return nil, err
	}
	out := &domain.SearchResult{
		MaxResults: result.MaxResults,
		StartAt:    startAt,
	}
	for _, i := range result.Issues {
		out.Issues = append(out.Issues, mapIssue(i))
	}
	out.Total = len(out.Issues)
	out.HasMore = result.NextPageToken != ""
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
		Fields:     []string{"summary", "description", "status", "priority", "issuetype", "assignee", "reporter", "labels", "created", "updated", "story_points", "customfield_10016", "customfield_10028"},
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

func (c *RestClient) UpdateIssue(ctx context.Context, input *domain.UpdateIssueInput) error {
	fields := map[string]any{}
	if input.Summary != "" {
		fields["summary"] = input.Summary
	}
	if input.Description != "" {
		fields["description"] = map[string]any{
			"type": "doc", "version": 1,
			"content": []map[string]any{
				{"type": "paragraph", "content": []map[string]any{
					{"type": "text", "text": input.Description},
				}},
			},
		}
	}
	if input.Priority != "" {
		fields["priority"] = map[string]any{"name": input.Priority}
	}
	if input.Assignee != "" {
		fields["assignee"] = map[string]any{"accountId": input.Assignee}
	}
	if input.Labels != nil {
		fields["labels"] = input.Labels
	}
	if len(fields) == 0 {
		return nil
	}

	payload := map[string]any{"fields": fields}
	data, _ := json.Marshal(payload)

	path := fmt.Sprintf("%s/rest/api/3/issue/%s", c.baseURL, input.Key)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}
		return fmt.Errorf("update issue failed %d: %s", resp.StatusCode, string(b))
	}
	return nil
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
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}
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
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}
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

func (c *RestClient) AssignIssue(ctx context.Context, issueKey, accountID string) error {
	path := fmt.Sprintf("%s/rest/api/3/issue/%s/assignee", c.baseURL, issueKey)
	var body []byte
	if accountID == "" {
		body = []byte(`{"accountId":null}`)
	} else {
		payload := map[string]string{"accountId": accountID}
		body, _ = json.Marshal(payload)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}
		return fmt.Errorf("assign failed %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *RestClient) DeleteIssue(ctx context.Context, issueKey string) error {
	path := fmt.Sprintf("%s/rest/api/3/issue/%s", c.baseURL, issueKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.email, c.token)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}
		return fmt.Errorf("delete failed %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *RestClient) CreateSubtask(ctx context.Context, parentKey string, input *domain.CreateIssueInput) (*domain.Issue, error) {
	fields := map[string]any{
		"project":   map[string]string{"key": input.Project},
		"summary":   input.Summary,
		"issuetype": map[string]string{"name": "Sub-task"},
		"parent":    map[string]string{"key": parentKey},
	}
	if input.Description != "" {
		fields["description"] = map[string]any{
			"type": "doc", "version": 1,
			"content": []map[string]any{
				{"type": "paragraph", "content": []map[string]any{
					{"type": "text", "text": input.Description},
				}},
			},
		}
	}
	if input.Priority != "" {
		fields["priority"] = map[string]string{"name": input.Priority}
	}
	if input.Assignee != "" {
		fields["assignee"] = map[string]string{"accountId": input.Assignee}
	}

	payload := map[string]any{"fields": fields}
	data, _ := json.Marshal(payload)

	path := fmt.Sprintf("%s/rest/api/3/issue", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}
		return nil, fmt.Errorf("create subtask failed %d: %s", resp.StatusCode, string(b))
	}

	var result struct {
		Key string `json:"key"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &domain.Issue{Key: result.Key, Summary: input.Summary}, nil
}

func (c *RestClient) FindUser(ctx context.Context, query string) ([]domain.User, error) {
	path := fmt.Sprintf("%s/rest/api/3/user/search?query=%s", c.baseURL, url.QueryEscape(query))
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
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}
		return nil, fmt.Errorf("user search failed %d: %s", resp.StatusCode, string(b))
	}

	var users []struct {
		AccountID   string `json:"accountId"`
		DisplayName string `json:"displayName"`
		Email       string `json:"emailAddress"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	out := make([]domain.User, len(users))
	for i, u := range users {
		out[i] = domain.User{AccountID: u.AccountID, DisplayName: u.DisplayName, Email: u.Email}
	}
	return out, nil
}

func (c *RestClient) SetEpicLink(ctx context.Context, issueKey, epicKey string) error {
	path := fmt.Sprintf("%s/rest/api/3/issue/%s", c.baseURL, issueKey)
	payload := map[string]any{"fields": map[string]any{"parent": map[string]string{"key": epicKey}}}
	data, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}
		return fmt.Errorf("set epic link failed %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *RestClient) RemoveEpicLink(ctx context.Context, issueKey string) error {
	path := fmt.Sprintf("%s/rest/api/3/issue/%s", c.baseURL, issueKey)
	payload := map[string]any{"fields": map[string]any{"parent": nil}}
	data, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}
		return fmt.Errorf("remove epic link failed %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *RestClient) GetSprints(ctx context.Context, boardID int, state string) ([]domain.Sprint, error) {
	path := fmt.Sprintf("%s/rest/agile/1.0/board/%d/sprint", c.baseURL, boardID)
	if state != "" {
		path += "?state=" + state
	}
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
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}
		return nil, fmt.Errorf("get sprints failed %d: %s", resp.StatusCode, string(b))
	}
	var result struct {
		Values []struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			State     string `json:"state"`
			Goal      string `json:"goal"`
			StartDate string `json:"startDate"`
			EndDate   string `json:"endDate"`
		} `json:"values"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	out := make([]domain.Sprint, len(result.Values))
	for i, s := range result.Values {
		out[i] = domain.Sprint{ID: s.ID, Name: s.Name, State: s.State, Goal: s.Goal, StartDate: s.StartDate, EndDate: s.EndDate}
	}
	return out, nil
}

func (c *RestClient) CreateSprint(ctx context.Context, boardID int, name, goal string) (*domain.Sprint, error) {
	path := fmt.Sprintf("%s/rest/agile/1.0/sprint", c.baseURL)
	payload := map[string]any{"name": name, "originBoardId": boardID}
	if goal != "" {
		payload["goal"] = goal
	}
	data, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}
		return nil, fmt.Errorf("create sprint failed %d: %s", resp.StatusCode, string(b))
	}
	var result struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &domain.Sprint{ID: result.ID, Name: result.Name, State: "future"}, nil
}

func (c *RestClient) StartSprint(ctx context.Context, sprintID int, startDate, endDate string) error {
	path := fmt.Sprintf("%s/rest/agile/1.0/sprint/%d", c.baseURL, sprintID)
	payload := map[string]any{"state": "active", "startDate": startDate, "endDate": endDate}
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
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}
		return fmt.Errorf("start sprint failed %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *RestClient) CloseSprint(ctx context.Context, sprintID int) error {
	path := fmt.Sprintf("%s/rest/agile/1.0/sprint/%d", c.baseURL, sprintID)
	payload := map[string]any{"state": "closed"}
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
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}
		return fmt.Errorf("close sprint failed %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *RestClient) MoveIssuesToSprint(ctx context.Context, sprintID int, issueKeys []string) error {
	path := fmt.Sprintf("%s/rest/agile/1.0/sprint/%d/issue", c.baseURL, sprintID)
	payload := map[string]any{"issues": issueKeys}
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
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}
		return fmt.Errorf("move issues to sprint failed %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *RestClient) LinkIssues(ctx context.Context, inwardKey, outwardKey, linkType string) error {
	path := fmt.Sprintf("%s/rest/api/3/issueLink", c.baseURL)
	payload := map[string]any{
		"type":         map[string]string{"name": linkType},
		"inwardIssue":  map[string]string{"key": inwardKey},
		"outwardIssue": map[string]string{"key": outwardKey},
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
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}
		return fmt.Errorf("link issues failed %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *RestClient) GetLinkTypes(ctx context.Context) ([]domain.LinkType, error) {
	path := fmt.Sprintf("%s/rest/api/3/issueLinkType", c.baseURL)
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
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}
		return nil, fmt.Errorf("get link types failed %d: %s", resp.StatusCode, string(b))
	}
	var result struct {
		IssueLinkTypes []struct {
			Name    string `json:"name"`
			Inward  string `json:"inward"`
			Outward string `json:"outward"`
		} `json:"issueLinkTypes"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	out := make([]domain.LinkType, len(result.IssueLinkTypes))
	for i, lt := range result.IssueLinkTypes {
		out[i] = domain.LinkType{Name: lt.Name, Inward: lt.Inward, Outward: lt.Outward}
	}
	return out, nil
}

func (c *RestClient) AddWorklog(ctx context.Context, issueKey, timeSpent, comment string) error {
	path := fmt.Sprintf("%s/rest/api/3/issue/%s/worklog", c.baseURL, issueKey)
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
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}
		return fmt.Errorf("add worklog failed %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *RestClient) GetWorklogs(ctx context.Context, issueKey string) ([]domain.Worklog, error) {
	path := fmt.Sprintf("%s/rest/api/3/issue/%s/worklog", c.baseURL, issueKey)
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
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}
		return nil, fmt.Errorf("get worklogs failed %d: %s", resp.StatusCode, string(b))
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
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	out := make([]domain.Worklog, len(result.Worklogs))
	for i, w := range result.Worklogs {
		var comment string
		if w.Comment != nil && len(w.Comment.Content) > 0 && len(w.Comment.Content[0].Content) > 0 {
			comment = w.Comment.Content[0].Content[0].Text
		}
		out[i] = domain.Worklog{Author: w.Author.DisplayName, TimeSpent: w.TimeSpent, Started: w.Started, Comment: comment}
	}
	return out, nil
}

func (c *RestClient) AddWatcher(ctx context.Context, issueKey, accountID string) error {
	path := fmt.Sprintf("%s/rest/api/3/issue/%s/watchers", c.baseURL, issueKey)
	data, _ := json.Marshal(accountID)
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
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}
		return fmt.Errorf("add watcher failed %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *RestClient) GetWatchers(ctx context.Context, issueKey string) ([]domain.User, error) {
	path := fmt.Sprintf("%s/rest/api/3/issue/%s/watchers", c.baseURL, issueKey)
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
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}
		return nil, fmt.Errorf("get watchers failed %d: %s", resp.StatusCode, string(b))
	}
	var result struct {
		Watchers []struct {
			AccountID   string `json:"accountId"`
			DisplayName string `json:"displayName"`
		} `json:"watchers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	out := make([]domain.User, len(result.Watchers))
	for i, w := range result.Watchers {
		out[i] = domain.User{AccountID: w.AccountID, DisplayName: w.DisplayName}
	}
	return out, nil
}

func (c *RestClient) AddLabel(ctx context.Context, issueKey, label string) error {
	path := fmt.Sprintf("%s/rest/api/3/issue/%s", c.baseURL, issueKey)
	payload := map[string]any{"update": map[string]any{"labels": []map[string]string{{"add": label}}}}
	data, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body: %w", err)
		}
		return fmt.Errorf("add label failed %d: %s", resp.StatusCode, string(b))
	}
	return nil
}

func (c *RestClient) GetProjects(ctx context.Context) ([]domain.Project, error) {
	path := fmt.Sprintf("%s/rest/api/3/project/search?maxResults=50", c.baseURL)
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
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}
		return nil, fmt.Errorf("get projects failed %d: %s", resp.StatusCode, string(b))
	}
	var result struct {
		Values []struct {
			Key  string `json:"key"`
			Name string `json:"name"`
			Lead struct {
				DisplayName string `json:"displayName"`
			} `json:"lead"`
			ProjectTypeKey string `json:"projectTypeKey"`
		} `json:"values"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	out := make([]domain.Project, len(result.Values))
	for i, p := range result.Values {
		out[i] = domain.Project{Key: p.Key, Name: p.Name, Lead: p.Lead.DisplayName, Type: p.ProjectTypeKey}
	}
	return out, nil
}

func (c *RestClient) GetProject(ctx context.Context, key string) (*domain.ProjectDetail, error) {
	path := fmt.Sprintf("%s/rest/api/3/project/%s", c.baseURL, key)
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
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response body: %w", err)
		}
		return nil, fmt.Errorf("get project failed %d: %s", resp.StatusCode, string(b))
	}
	var raw struct {
		Key  string `json:"key"`
		Name string `json:"name"`
		Lead struct {
			DisplayName string `json:"displayName"`
		} `json:"lead"`
		ProjectTypeKey string `json:"projectTypeKey"`
		Description    string `json:"description"`
		Components     []struct {
			Name string `json:"name"`
		} `json:"components"`
		Versions []struct {
			Name string `json:"name"`
		} `json:"versions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	pd := &domain.ProjectDetail{
		Key:         raw.Key,
		Name:        raw.Name,
		Lead:        raw.Lead.DisplayName,
		Type:        raw.ProjectTypeKey,
		Description: raw.Description,
	}
	for _, c := range raw.Components {
		pd.Components = append(pd.Components, c.Name)
	}
	for _, v := range raw.Versions {
		pd.Versions = append(pd.Versions, v.Name)
	}
	return pd, nil
}

func (c *RestClient) RawRequest(ctx context.Context, method, path string, body []byte) ([]byte, int, error) {
	c.limiter.wait()
	fullURL := c.baseURL + path
	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return nil, 0, err
	}
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Accept", "application/json")
	if len(body) > 0 {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return nil, 0, fmt.Errorf("read response body: %w", err)
	}
	return respBody, resp.StatusCode, nil
}

func mapIssue(raw *issue.Issue) domain.Issue {
	if raw == nil {
		return domain.Issue{}
	}
	i := domain.Issue{
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
	// Extract story points from common custom field IDs
	if raw.Fields != nil && raw.Fields.Custom != nil {
		for _, fieldID := range storyPointFields {
			if sp, ok := raw.Fields.Custom.GetNumber(fieldID); ok && sp > 0 {
				i.StoryPoints = sp
				break
			}
		}
	}
	return i
}

// storyPointFields lists common Jira custom field IDs for story points.
// The first match wins. Covers: next-gen story_points, classic story points, and common variants.
var storyPointFields = []string{
	"story_points",      // next-gen projects
	"customfield_10016", // Jira Cloud default
	"customfield_10028", // common alternative
	"customfield_10004", // some instances
	"customfield_10014", // another variant
}

func (c *RestClient) GetAttachments(ctx context.Context, issueKey string) ([]domain.Attachment, error) {
	data, _, err := c.RawRequest(ctx, "GET", fmt.Sprintf("/rest/api/3/issue/%s?fields=attachment", issueKey), nil)
	if err != nil {
		return nil, err
	}
	var raw struct {
		Fields struct {
			Attachment []struct {
				ID       string `json:"id"`
				Filename string `json:"filename"`
				Size     int64  `json:"size"`
				MimeType string `json:"mimeType"`
				Author   struct {
					DisplayName string `json:"displayName"`
				} `json:"author"`
				Created string `json:"created"`
				Content string `json:"content"`
			} `json:"attachment"`
		} `json:"fields"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	out := make([]domain.Attachment, len(raw.Fields.Attachment))
	for i, a := range raw.Fields.Attachment {
		out[i] = domain.Attachment{
			ID: a.ID, Filename: a.Filename, Size: a.Size,
			MimeType: a.MimeType, Author: a.Author.DisplayName,
			Created: a.Created, URL: a.Content,
		}
	}
	return out, nil
}

func (c *RestClient) GetVersions(ctx context.Context, projectKey string) ([]domain.Version, error) {
	data, _, err := c.RawRequest(ctx, "GET", fmt.Sprintf("/rest/api/3/project/%s/versions", projectKey), nil)
	if err != nil {
		return nil, err
	}
	var raw []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Released    bool   `json:"released"`
		ReleaseDate string `json:"releaseDate"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	out := make([]domain.Version, len(raw))
	for i, v := range raw {
		out[i] = domain.Version{ID: v.ID, Name: v.Name, Description: v.Description, Released: v.Released, ReleaseDate: v.ReleaseDate}
	}
	return out, nil
}

func (c *RestClient) CreateVersion(ctx context.Context, projectKey, name, description string) (*domain.Version, error) {
	payload, _ := json.Marshal(map[string]string{"name": name, "description": description, "project": projectKey})
	data, _, err := c.RawRequest(ctx, "POST", "/rest/api/3/version", payload)
	if err != nil {
		return nil, err
	}
	var raw struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return &domain.Version{ID: raw.ID, Name: raw.Name, Description: description}, nil
}

func (c *RestClient) ReleaseVersion(ctx context.Context, versionID string) error {
	payload, _ := json.Marshal(map[string]any{"released": true, "releaseDate": time.Now().Format("2006-01-02")})
	_, _, err := c.RawRequest(ctx, "PUT", fmt.Sprintf("/rest/api/3/version/%s", versionID), payload)
	return err
}

func (c *RestClient) GetComponents(ctx context.Context, projectKey string) ([]domain.Component, error) {
	data, _, err := c.RawRequest(ctx, "GET", fmt.Sprintf("/rest/api/3/project/%s/components", projectKey), nil)
	if err != nil {
		return nil, err
	}
	var raw []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
		Lead *struct {
			DisplayName string `json:"displayName"`
		} `json:"lead"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	out := make([]domain.Component, len(raw))
	for i, comp := range raw {
		lead := ""
		if comp.Lead != nil {
			lead = comp.Lead.DisplayName
		}
		out[i] = domain.Component{ID: comp.ID, Name: comp.Name, Lead: lead}
	}
	return out, nil
}

func (c *RestClient) GetFields(ctx context.Context) ([]domain.Field, error) {
	data, _, err := c.RawRequest(ctx, "GET", "/rest/api/3/field", nil)
	if err != nil {
		return nil, err
	}
	var raw []struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Custom bool   `json:"custom"`
		Schema *struct {
			Type string `json:"type"`
		} `json:"schema"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	out := make([]domain.Field, len(raw))
	for i, f := range raw {
		ft := ""
		if f.Schema != nil {
			ft = f.Schema.Type
		}
		out[i] = domain.Field{ID: f.ID, Name: f.Name, Custom: f.Custom, Type: ft}
	}
	return out, nil
}

var _ domain.Client = (*RestClient)(nil)
