package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/felixgeelhaar/jirasdk/core/issue"

	domain "github.com/aldok10/zara-jira-mcp/modules/jira/domain"
)

// CreateIssue creates a new Jira issue.
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
		return nil, fmt.Errorf("create issue: %w", err)
	}
	mapped := mapIssue(created)
	return &mapped, nil
}

// UpdateIssue updates fields on an existing issue.
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
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	_, _, err = c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/rest/api/3/issue/%s", input.Key), data)
	return err
}

// DeleteIssue deletes an issue by key.
func (c *RestClient) DeleteIssue(ctx context.Context, issueKey string) error {
	_, _, err := c.doRequest(ctx, http.MethodDelete, fmt.Sprintf("/rest/api/3/issue/%s", issueKey), nil)
	return err
}

// AddComment adds a comment to an issue.
func (c *RestClient) AddComment(ctx context.Context, issueKey, body string) error {
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
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}

	_, _, err = c.doRequest(ctx, http.MethodPost, fmt.Sprintf("/rest/api/3/issue/%s/comment", issueKey), data)
	return err
}

// TransitionIssue transitions an issue to a new status.
func (c *RestClient) TransitionIssue(ctx context.Context, issueKey, transitionID string) error {
	return c.sdk.Issue.DoTransition(ctx, issueKey, &issue.TransitionInput{
		Transition: &issue.Transition{ID: transitionID},
	})
}

// GetTransitions returns available transitions for an issue.
func (c *RestClient) GetTransitions(ctx context.Context, issueKey string) ([]domain.Transition, error) {
	data, _, err := c.doRequest(ctx, http.MethodGet, fmt.Sprintf("/rest/api/3/issue/%s/transitions", issueKey), nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Transitions []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"transitions"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode transitions: %w", err)
	}

	out := make([]domain.Transition, len(result.Transitions))
	for i, t := range result.Transitions {
		out[i] = domain.Transition{ID: t.ID, Name: t.Name}
	}
	return out, nil
}

// AssignIssue assigns an issue to a user.
func (c *RestClient) AssignIssue(ctx context.Context, issueKey, accountID string) error {
	var body []byte
	if accountID == "" {
		body = []byte(`{"accountId":null}`)
	} else {
		var err error
		payload := map[string]string{"accountId": accountID}
		body, err = json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("marshal payload: %w", err)
		}
	}
	_, _, err := c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/rest/api/3/issue/%s/assignee", issueKey), body)
	return err
}

// CreateSubtask creates a sub-task under a parent issue.
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
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	body, _, err := c.doRequest(ctx, http.MethodPost, "/rest/api/3/issue", data)
	if err != nil {
		return nil, err
	}

	var result struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &domain.Issue{Key: result.Key, Summary: input.Summary}, nil
}

// AddLabel adds a label to an issue.
func (c *RestClient) AddLabel(ctx context.Context, issueKey, label string) error {
	payload := map[string]any{"update": map[string]any{"labels": []map[string]string{{"add": label}}}}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal payload: %w", err)
	}
	_, _, err = c.doRequest(ctx, http.MethodPut, fmt.Sprintf("/rest/api/3/issue/%s", issueKey), data)
	return err
}
