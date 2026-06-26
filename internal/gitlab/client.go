package gitlab

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/aldok10/zara-jira-mcp/config"
)

type Client struct {
	token      string
	baseURL    string
	projectID  string
	httpClient *http.Client
}

func NewClient(cfg *config.Config) *Client {
	baseURL := cfg.GitLab.BaseURL
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}
	return &Client{
		token:      cfg.GitLab.Token,
		baseURL:    baseURL,
		projectID:  cfg.GitLab.ProjectID,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) Available() bool {
	return c.token != "" && c.projectID != ""
}

type Issue struct {
	IID       int
	Title     string
	Description string
	State     string
	Labels    []string
	Assignee  string
	Milestone string
	CreatedAt time.Time
	WebURL    string
}

type Milestone struct {
	ID          int
	IID         int
	Title       string
	Description string
	State       string
	DueDate     string
}

type MergeRequest struct {
	IID          int
	Title        string
	State        string
	Author       string
	TargetBranch string
	CreatedAt    time.Time
	Draft        bool
}

// CreateIssue creates a GitLab issue.
func (c *Client) CreateIssue(ctx context.Context, title, description string, labels []string, assigneeID int, milestoneID int) (*Issue, error) {
	payload := map[string]any{"title": title}
	if description != "" {
		payload["description"] = description
	}
	if len(labels) > 0 {
		payload["labels"] = joinLabels(labels)
	}
	if assigneeID > 0 {
		payload["assignee_ids"] = []int{assigneeID}
	}
	if milestoneID > 0 {
		payload["milestone_id"] = milestoneID
	}

	path := fmt.Sprintf("/api/v4/projects/%s/issues", url.PathEscape(c.projectID))
	body, err := c.doPost(ctx, path, payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		IID    int    `json:"iid"`
		Title  string `json:"title"`
		State  string `json:"state"`
		WebURL string `json:"web_url"`
	}
	_ = json.Unmarshal(body, &result)
	return &Issue{IID: result.IID, Title: result.Title, State: result.State, WebURL: result.WebURL}, nil
}

// ListIssues lists project issues.
func (c *Client) ListIssues(ctx context.Context, state string, labels string, limit int) ([]Issue, error) {
	if limit <= 0 {
		limit = 20
	}
	path := fmt.Sprintf("/api/v4/projects/%s/issues?per_page=%d", url.PathEscape(c.projectID), limit)
	if state != "" {
		path += "&state=" + state
	}
	if labels != "" {
		path += "&labels=" + url.QueryEscape(labels)
	}

	body, err := c.doGet(ctx, path)
	if err != nil {
		return nil, err
	}

	var items []struct {
		IID         int      `json:"iid"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		State       string   `json:"state"`
		Labels      []string `json:"labels"`
		Assignee    *struct{ Name string } `json:"assignee"`
		Milestone   *struct{ Title string } `json:"milestone"`
		CreatedAt   string `json:"created_at"`
		WebURL      string `json:"web_url"`
	}
	_ = json.Unmarshal(body, &items)

	var issues []Issue
	for _, item := range items {
		assignee := ""
		if item.Assignee != nil {
			assignee = item.Assignee.Name
		}
		milestone := ""
		if item.Milestone != nil {
			milestone = item.Milestone.Title
		}
		created, _ := time.Parse(time.RFC3339, item.CreatedAt)
		issues = append(issues, Issue{
			IID: item.IID, Title: item.Title, Description: item.Description,
			State: item.State, Labels: item.Labels, Assignee: assignee,
			Milestone: milestone, CreatedAt: created, WebURL: item.WebURL,
		})
	}
	return issues, nil
}

// CreateMilestone creates a GitLab milestone.
func (c *Client) CreateMilestone(ctx context.Context, title, description, dueDate string) (*Milestone, error) {
	payload := map[string]any{"title": title}
	if description != "" {
		payload["description"] = description
	}
	if dueDate != "" {
		payload["due_date"] = dueDate
	}

	path := fmt.Sprintf("/api/v4/projects/%s/milestones", url.PathEscape(c.projectID))
	body, err := c.doPost(ctx, path, payload)
	if err != nil {
		return nil, err
	}

	var m struct {
		ID    int    `json:"id"`
		IID   int    `json:"iid"`
		Title string `json:"title"`
	}
	_ = json.Unmarshal(body, &m)
	return &Milestone{ID: m.ID, IID: m.IID, Title: m.Title, State: "active"}, nil
}

// ListMilestones lists project milestones.
func (c *Client) ListMilestones(ctx context.Context, state string) ([]Milestone, error) {
	path := fmt.Sprintf("/api/v4/projects/%s/milestones", url.PathEscape(c.projectID))
	if state != "" {
		path += "?state=" + state
	}
	body, err := c.doGet(ctx, path)
	if err != nil {
		return nil, err
	}

	var items []struct {
		ID          int    `json:"id"`
		IID         int    `json:"iid"`
		Title       string `json:"title"`
		Description string `json:"description"`
		State       string `json:"state"`
		DueDate     string `json:"due_date"`
	}
	_ = json.Unmarshal(body, &items)

	var milestones []Milestone
	for _, item := range items {
		milestones = append(milestones, Milestone{
			ID: item.ID, IID: item.IID, Title: item.Title,
			Description: item.Description, State: item.State, DueDate: item.DueDate,
		})
	}
	return milestones, nil
}

// ListMRs lists merge requests.
func (c *Client) ListMRs(ctx context.Context, state string, limit int) ([]MergeRequest, error) {
	if limit <= 0 {
		limit = 20
	}
	path := fmt.Sprintf("/api/v4/projects/%s/merge_requests?per_page=%d", url.PathEscape(c.projectID), limit)
	if state != "" {
		path += "&state=" + state
	}
	body, err := c.doGet(ctx, path)
	if err != nil {
		return nil, err
	}

	var items []struct {
		IID       int    `json:"iid"`
		Title     string `json:"title"`
		State     string `json:"state"`
		Author    struct{ Name string } `json:"author"`
		CreatedAt string `json:"created_at"`
		Draft     bool   `json:"draft"`
	}
	_ = json.Unmarshal(body, &items)

	var mrs []MergeRequest
	for _, item := range items {
		created, _ := time.Parse(time.RFC3339, item.CreatedAt)
		mrs = append(mrs, MergeRequest{
			IID: item.IID, Title: item.Title, State: item.State,
			Author: item.Author.Name, CreatedAt: created, Draft: item.Draft,
		})
	}
	return mrs, nil
}

// GetFileContent reads a file from the repo.
func (c *Client) GetFileContent(ctx context.Context, filePath, ref string) (string, error) {
	if ref == "" {
		ref = "main"
	}
	path := fmt.Sprintf("/api/v4/projects/%s/repository/files/%s/raw?ref=%s",
		url.PathEscape(c.projectID), url.PathEscape(filePath), ref)
	body, err := c.doGet(ctx, path)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// ListFiles lists repository tree.
func (c *Client) ListFiles(ctx context.Context, dirPath, ref string) ([]string, error) {
	if ref == "" {
		ref = "main"
	}
	path := fmt.Sprintf("/api/v4/projects/%s/repository/tree?ref=%s&per_page=100",
		url.PathEscape(c.projectID), ref)
	if dirPath != "" {
		path += "&path=" + url.QueryEscape(dirPath)
	}
	body, err := c.doGet(ctx, path)
	if err != nil {
		return nil, err
	}

	var items []struct {
		Name string `json:"name"`
		Type string `json:"type"`
	}
	_ = json.Unmarshal(body, &items)

	var files []string
	for _, item := range items {
		suffix := ""
		if item.Type == "tree" {
			suffix = "/"
		}
		files = append(files, item.Name+suffix)
	}
	return files, nil
}

// SearchBranches finds branches matching a pattern.
func (c *Client) SearchBranches(ctx context.Context, pattern string) ([]Branch, error) {
	path := fmt.Sprintf("/api/v4/projects/%s/repository/branches?search=%s", c.projectID, pattern)
	body, err := c.doGet(ctx, path)
	if err != nil {
		return nil, err
	}
	var items []struct {
		Name   string `json:"name"`
		Merged bool   `json:"merged"`
	}
	_ = json.Unmarshal(body, &items)
	out := make([]Branch, len(items))
	for i, item := range items {
		out[i] = Branch{Name: item.Name, Merged: item.Merged}
	}
	return out, nil
}

// SearchMRsByBranch finds merge requests for a source branch.
func (c *Client) SearchMRsByBranch(ctx context.Context, branch string) ([]MergeRequest, error) {
	path := fmt.Sprintf("/api/v4/projects/%s/merge_requests?source_branch=%s&state=all&per_page=10", c.projectID, branch)
	body, err := c.doGet(ctx, path)
	if err != nil {
		return nil, err
	}
	var items []struct {
		IID          int    `json:"iid"`
		Title        string `json:"title"`
		State        string `json:"state"`
		TargetBranch string `json:"target_branch"`
		Author       struct {
			Username string `json:"username"`
		} `json:"author"`
		MergedAt string `json:"merged_at"`
	}
	_ = json.Unmarshal(body, &items)
	var mrs []MergeRequest
	for _, item := range items {
		mrs = append(mrs, MergeRequest{
			IID:          item.IID,
			Title:        item.Title,
			State:        item.State,
			TargetBranch: item.TargetBranch,
			Author:       item.Author.Username,
		})
	}
	return mrs, nil
}

// Branch represents a git branch.
type Branch struct {
	Name   string
	Merged bool
}

// HTTP helpers

func (c *Client) doGet(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("PRIVATE-TOKEN", c.token)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("gitlab API %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

func (c *Client) doPost(ctx context.Context, path string, payload any) ([]byte, error) {
	data, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("PRIVATE-TOKEN", c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("gitlab API %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

func joinLabels(labels []string) string {
	result := ""
	for i, l := range labels {
		if i > 0 {
			result += ","
		}
		result += l
	}
	return result
}
