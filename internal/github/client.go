package github

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aldok10/zara-jira-mcp/config"
)

type Client struct {
	token      string
	owner      string
	repo       string
	httpClient *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		token:      cfg.GitHub.Token,
		owner:      cfg.GitHub.Owner,
		repo:       cfg.GitHub.Repo,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) Available() bool {
	return c.token != "" && c.owner != "" && c.repo != ""
}

type PullRequest struct {
	Number    int
	Title     string
	State     string
	Author    string
	User      string
	CreatedAt time.Time
	UpdatedAt time.Time
	Draft     bool
	Reviewers []string
	MergedTo  string
}

type Release struct {
	TagName     string
	Name        string
	PublishedAt time.Time
	Author      string
}

// ListPRs returns open pull requests.
func (c *Client) ListPRs(ctx context.Context, state string, limit int) ([]PullRequest, error) {
	if state == "" {
		state = "open"
	}
	if limit <= 0 {
		limit = 30
	}
	path := fmt.Sprintf("/repos/%s/%s/pulls?state=%s&per_page=%d&sort=updated&direction=desc",
		c.owner, c.repo, state, limit)
	body, err := c.doGet(ctx, path)
	if err != nil {
		return nil, err
	}

	var items []struct {
		Number    int    `json:"number"`
		Title     string `json:"title"`
		State     string `json:"state"`
		Draft     bool   `json:"draft"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		User      struct {
			Login string `json:"login"`
		} `json:"user"`
		RequestedReviewers []struct {
			Login string `json:"login"`
		} `json:"requested_reviewers"`
	}
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}

	prs := make([]PullRequest, 0, len(items))
	for _, item := range items {
		created, _ := time.Parse(time.RFC3339, item.CreatedAt)
		updated, _ := time.Parse(time.RFC3339, item.UpdatedAt)
		reviewers := make([]string, len(item.RequestedReviewers))
		for i, r := range item.RequestedReviewers {
			reviewers[i] = r.Login
		}
		prs = append(prs, PullRequest{
			Number:    item.Number,
			Title:     item.Title,
			State:     item.State,
			User:      item.User.Login,
			CreatedAt: created,
			UpdatedAt: updated,
			Draft:     item.Draft,
			Reviewers: reviewers,
		})
	}
	return prs, nil
}

// ListReleases returns recent releases.
func (c *Client) ListReleases(ctx context.Context, limit int) ([]Release, error) {
	if limit <= 0 {
		limit = 10
	}
	path := fmt.Sprintf("/repos/%s/%s/releases?per_page=%d", c.owner, c.repo, limit)
	body, err := c.doGet(ctx, path)
	if err != nil {
		return nil, err
	}

	var items []struct {
		TagName     string `json:"tag_name"`
		Name        string `json:"name"`
		PublishedAt string `json:"published_at"`
		Author      struct {
			Login string `json:"login"`
		} `json:"author"`
	}
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}

	releases := make([]Release, 0, len(items))
	for _, item := range items {
		pub, _ := time.Parse(time.RFC3339, item.PublishedAt)
		releases = append(releases, Release{
			TagName:     item.TagName,
			Name:        item.Name,
			PublishedAt: pub,
			Author:      item.Author.Login,
		})
	}
	return releases, nil
}

type RepoActivity struct {
	CommitCount    int
	PRsMerged      int
	IssuesClosed   int
}

// GetActivity returns repo activity summary for the last N days.
func (c *Client) GetActivity(ctx context.Context, days int) (*RepoActivity, error) {
	if days <= 0 {
		days = 7
	}
	since := time.Now().UTC().AddDate(0, 0, -days).Format(time.RFC3339)

	// Commits
	commitPath := fmt.Sprintf("/repos/%s/%s/commits?since=%s&per_page=100", c.owner, c.repo, since)
	commitBody, err := c.doGet(ctx, commitPath)
	if err != nil {
		return nil, err
	}
	var commits []json.RawMessage
	_ = json.Unmarshal(commitBody, &commits)

	// Merged PRs
	prPath := fmt.Sprintf("/repos/%s/%s/pulls?state=closed&sort=updated&direction=desc&per_page=100", c.owner, c.repo)
	prBody, err := c.doGet(ctx, prPath)
	if err != nil {
		return nil, err
	}
	var prs []struct {
		MergedAt string `json:"merged_at"`
	}
	_ = json.Unmarshal(prBody, &prs)
	sinceTime := time.Now().UTC().AddDate(0, 0, -days)
	merged := 0
	for _, pr := range prs {
		if pr.MergedAt != "" {
			t, _ := time.Parse(time.RFC3339, pr.MergedAt)
			if t.After(sinceTime) {
				merged++
			}
		}
	}

	// Closed issues
	issuePath := fmt.Sprintf("/repos/%s/%s/issues?state=closed&since=%s&per_page=100", c.owner, c.repo, since)
	issueBody, err := c.doGet(ctx, issuePath)
	if err != nil {
		return nil, err
	}
	var issues []struct {
		PullRequest *json.RawMessage `json:"pull_request"`
	}
	_ = json.Unmarshal(issueBody, &issues)
	closed := 0
	for _, iss := range issues {
		if iss.PullRequest == nil {
			closed++
		}
	}

	return &RepoActivity{
		CommitCount:  len(commits),
		PRsMerged:    merged,
		IssuesClosed: closed,
	}, nil
}

// SearchBranches finds branches matching a pattern (e.g. issue key).
func (c *Client) SearchBranches(ctx context.Context, pattern string) ([]Branch, error) {
	path := fmt.Sprintf("/repos/%s/%s/branches?per_page=100", c.owner, c.repo)
	body, err := c.doGet(ctx, path)
	if err != nil {
		return nil, err
	}
	var all []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(body, &all); err != nil {
		return nil, err
	}
	var matched []Branch
	lowerPattern := strings.ToLower(pattern)
	for _, b := range all {
		if strings.Contains(strings.ToLower(b.Name), lowerPattern) {
			matched = append(matched, Branch{Name: b.Name})
		}
	}
	return matched, nil
}

// SearchPRsByBranch finds PRs (open/closed/merged) for a branch.
func (c *Client) SearchPRsByBranch(ctx context.Context, branch string) ([]PullRequest, error) {
	// Search all states
	path := fmt.Sprintf("/repos/%s/%s/pulls?state=all&head=%s:%s&per_page=10",
		c.owner, c.repo, c.owner, branch)
	body, err := c.doGet(ctx, path)
	if err != nil {
		return nil, err
	}
	var items []struct {
		Number   int    `json:"number"`
		Title    string `json:"title"`
		State    string `json:"state"`
		MergedAt string `json:"merged_at"`
		Base     struct {
			Ref string `json:"ref"`
		} `json:"base"`
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	}
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, err
	}
	var prs []PullRequest
	for _, item := range items {
		pr := PullRequest{
			Number: item.Number,
			Title:  item.Title,
			State:  item.State,
			Author: item.User.Login,
		}
		if item.MergedAt != "" {
			pr.State = "merged"
			pr.MergedTo = item.Base.Ref
		}
		prs = append(prs, pr)
	}
	return prs, nil
}

// Branch represents a git branch.
type Branch struct {
	Name string
}

func (c *Client) doGet(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com"+path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("github API %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

func (c *Client) doPost(ctx context.Context, path string, payload any) ([]byte, error) {
	data, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.github.com"+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("github API %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

// Issue represents a GitHub issue.
type Issue struct {
	Number    int
	Title     string
	Body      string
	State     string
	Labels    []string
	Assignee  string
	Milestone string
	CreatedAt time.Time
}

// CreateIssue creates a GitHub issue.
func (c *Client) CreateIssue(ctx context.Context, title, body string, labels []string, assignees []string, milestone int) (*Issue, error) {
	payload := map[string]any{
		"title": title,
		"body":  body,
	}
	if len(labels) > 0 {
		payload["labels"] = labels
	}
	if len(assignees) > 0 {
		payload["assignees"] = assignees
	}
	if milestone > 0 {
		payload["milestone"] = milestone
	}

	path := fmt.Sprintf("/repos/%s/%s/issues", c.owner, c.repo)
	resp, err := c.doPost(ctx, path, payload)
	if err != nil {
		return nil, err
	}

	var result struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
		State  string `json:"state"`
	}
	_ = json.Unmarshal(resp, &result)
	return &Issue{Number: result.Number, Title: result.Title, State: result.State}, nil
}

// ListIssues lists open issues.
func (c *Client) ListIssues(ctx context.Context, state string, labels string, limit int) ([]Issue, error) {
	if state == "" {
		state = "open"
	}
	if limit <= 0 {
		limit = 30
	}
	path := fmt.Sprintf("/repos/%s/%s/issues?state=%s&per_page=%d", c.owner, c.repo, state, limit)
	if labels != "" {
		path += "&labels=" + labels
	}
	body, err := c.doGet(ctx, path)
	if err != nil {
		return nil, err
	}

	var items []struct {
		Number      int    `json:"number"`
		Title       string `json:"title"`
		Body        string `json:"body"`
		State       string `json:"state"`
		Labels      []struct{ Name string } `json:"labels"`
		Assignee    *struct{ Login string } `json:"assignee"`
		Milestone   *struct{ Title string } `json:"milestone"`
		CreatedAt   string `json:"created_at"`
		PullRequest *json.RawMessage `json:"pull_request"`
	}
	_ = json.Unmarshal(body, &items)

	var issues []Issue
	for _, item := range items {
		if item.PullRequest != nil {
			continue // skip PRs
		}
		var labels []string
		for _, l := range item.Labels {
			labels = append(labels, l.Name)
		}
		assignee := ""
		if item.Assignee != nil {
			assignee = item.Assignee.Login
		}
		milestone := ""
		if item.Milestone != nil {
			milestone = item.Milestone.Title
		}
		created, _ := time.Parse(time.RFC3339, item.CreatedAt)
		issues = append(issues, Issue{
			Number: item.Number, Title: item.Title, Body: item.Body,
			State: item.State, Labels: labels, Assignee: assignee,
			Milestone: milestone, CreatedAt: created,
		})
	}
	return issues, nil
}

// Milestone represents a GitHub milestone.
type Milestone struct {
	Number      int
	Title       string
	Description string
	State       string
	DueOn       string
	OpenIssues  int
	ClosedIssues int
}

// CreateMilestone creates a milestone.
func (c *Client) CreateMilestone(ctx context.Context, title, description, dueOn string) (*Milestone, error) {
	payload := map[string]any{"title": title}
	if description != "" {
		payload["description"] = description
	}
	if dueOn != "" {
		payload["due_on"] = dueOn + "T00:00:00Z"
	}
	path := fmt.Sprintf("/repos/%s/%s/milestones", c.owner, c.repo)
	resp, err := c.doPost(ctx, path, payload)
	if err != nil {
		return nil, err
	}
	var m struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	}
	_ = json.Unmarshal(resp, &m)
	return &Milestone{Number: m.Number, Title: m.Title, State: "open"}, nil
}

// ListMilestones lists milestones.
func (c *Client) ListMilestones(ctx context.Context, state string) ([]Milestone, error) {
	if state == "" {
		state = "open"
	}
	path := fmt.Sprintf("/repos/%s/%s/milestones?state=%s", c.owner, c.repo, state)
	body, err := c.doGet(ctx, path)
	if err != nil {
		return nil, err
	}
	var items []struct {
		Number       int    `json:"number"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		State        string `json:"state"`
		DueOn        string `json:"due_on"`
		OpenIssues   int    `json:"open_issues"`
		ClosedIssues int    `json:"closed_issues"`
	}
	_ = json.Unmarshal(body, &items)

	var milestones []Milestone
	for _, item := range items {
		milestones = append(milestones, Milestone{
			Number: item.Number, Title: item.Title, Description: item.Description,
			State: item.State, DueOn: item.DueOn, OpenIssues: item.OpenIssues, ClosedIssues: item.ClosedIssues,
		})
	}
	return milestones, nil
}

// GetFileContent reads a file from the repo.
func (c *Client) GetFileContent(ctx context.Context, path, ref string) (string, error) {
	apiPath := fmt.Sprintf("/repos/%s/%s/contents/%s", c.owner, c.repo, path)
	if ref != "" {
		apiPath += "?ref=" + ref
	}
	body, err := c.doGet(ctx, apiPath)
	if err != nil {
		return "", err
	}
	var file struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}
	_ = json.Unmarshal(body, &file)
	if file.Encoding == "base64" {
		decoded, err := base64.StdEncoding.DecodeString(file.Content)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	}
	return file.Content, nil
}

// ListFiles lists files/directories at a path.
func (c *Client) ListFiles(ctx context.Context, path, ref string) ([]string, error) {
	apiPath := fmt.Sprintf("/repos/%s/%s/contents/%s", c.owner, c.repo, path)
	if ref != "" {
		apiPath += "?ref=" + ref
	}
	body, err := c.doGet(ctx, apiPath)
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
		if item.Type == "dir" {
			suffix = "/"
		}
		files = append(files, item.Name+suffix)
	}
	return files, nil
}
