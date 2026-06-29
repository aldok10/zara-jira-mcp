package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/felixgeelhaar/jirasdk/core/agile"
	"github.com/felixgeelhaar/jirasdk/core/search"

	domain "github.com/aldok10/zara-jira-mcp/modules/jira/domain"
)

// GetActiveSprints returns active sprints for a board.
func (c *RestClient) GetActiveSprints(ctx context.Context, boardID int) ([]domain.Sprint, error) {
	sprints, err := c.sdk.Agile.GetBoardSprints(ctx, int64(boardID), &agile.SprintsOptions{
		State: "active",
	})
	if err != nil {
		return nil, fmt.Errorf("get active sprints: %w", err)
	}
	out := make([]domain.Sprint, len(sprints))
	for i, s := range sprints {
		out[i] = domain.Sprint{ID: int(s.ID), Name: s.Name, State: s.State, Goal: s.Goal}
	}
	return out, nil
}

// GetSprintIssues returns all issues in a sprint.
func (c *RestClient) GetSprintIssues(ctx context.Context, sprintID int) ([]domain.Issue, error) {
	jql := fmt.Sprintf("sprint = %d ORDER BY status ASC", sprintID)
	result, err := c.sdk.Search.SearchJQL(ctx, &search.SearchJQLOptions{
		JQL:        jql,
		Fields:     []string{"summary", "description", "status", "priority", "issuetype", "assignee", "reporter", "labels", "created", "updated", "duedate", "story_points"},
		MaxResults: 100,
	})
	if err != nil {
		return nil, fmt.Errorf("get sprint issues: %w", err)
	}
	out := make([]domain.Issue, 0, len(result.Issues))
	for _, i := range result.Issues {
		out = append(out, mapIssue(i))
	}
	return out, nil
}

// GetSprints returns sprints for a board, optionally filtered by state.
func (c *RestClient) GetSprints(ctx context.Context, boardID int, state string) ([]domain.Sprint, error) {
	path := fmt.Sprintf("/rest/agile/1.0/board/%d/sprint", boardID)
	if state != "" {
		path += "?state=" + state
	}
	data, _, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
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
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode sprints: %w", err)
	}
	out := make([]domain.Sprint, len(result.Values))
	for i, s := range result.Values {
		out[i] = domain.Sprint{ID: s.ID, Name: s.Name, State: s.State, Goal: s.Goal, StartDate: s.StartDate, EndDate: s.EndDate}
	}
	return out, nil
}

// CreateSprint creates a new sprint on a board.
func (c *RestClient) CreateSprint(ctx context.Context, boardID int, name, goal string) (*domain.Sprint, error) {
	payload := map[string]any{"name": name, "originBoardId": boardID}
	if goal != "" {
		payload["goal"] = goal
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}
	body, _, err := c.doRequest(ctx, http.MethodPost, "/rest/agile/1.0/sprint", data)
	if err != nil {
		return nil, err
	}
	var result struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode sprint: %w", err)
	}
	return &domain.Sprint{ID: result.ID, Name: result.Name, State: "future"}, nil
}

// StartSprint starts a sprint with start/end dates.
func (c *RestClient) StartSprint(ctx context.Context, sprintID int, startDate, endDate string) error {
	payload := map[string]any{"state": "active", "startDate": startDate, "endDate": endDate}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	_, _, err = c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/rest/agile/1.0/sprint/%d", sprintID), data)
	return err
}

// CloseSprint closes a sprint.
func (c *RestClient) CloseSprint(ctx context.Context, sprintID int) error {
	payload := map[string]any{"state": "closed"}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	_, _, err = c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/rest/agile/1.0/sprint/%d", sprintID), data)
	return err
}

// MoveIssuesToSprint moves issues into a sprint.
func (c *RestClient) MoveIssuesToSprint(ctx context.Context, sprintID int, issueKeys []string) error {
	payload := map[string]any{"issues": issueKeys}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	_, _, err = c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/rest/agile/1.0/sprint/%d/issue", sprintID), data)
	return err
}
