package client

import (
	"context"
	"fmt"

	"github.com/felixgeelhaar/jirasdk/core/search"

	domain "github.com/aldok10/zara-jira-mcp/modules/jira/domain"
)

// SearchIssues searches Jira issues using JQL.
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
		return nil, fmt.Errorf("search issues: %w", err)
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

// GetIssue retrieves a single Jira issue by key.
func (c *RestClient) GetIssue(ctx context.Context, key string) (*domain.Issue, error) {
	i, err := c.sdk.Issue.Get(ctx, key, nil)
	if err != nil {
		return nil, fmt.Errorf("get issue %s: %w", key, err)
	}
	mapped := mapIssue(i)
	return &mapped, nil
}

// GetBoards returns all accessible boards.
func (c *RestClient) GetBoards(ctx context.Context) ([]domain.Board, error) {
	boards, err := c.sdk.Agile.GetBoards(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("get boards: %w", err)
	}
	out := make([]domain.Board, len(boards))
	for i, b := range boards {
		out[i] = domain.Board{ID: int(b.ID), Name: b.Name, Type: b.Type}
	}
	return out, nil
}
