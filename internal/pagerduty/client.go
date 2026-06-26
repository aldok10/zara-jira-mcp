package pagerduty

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aldok10/zara-jira-mcp/config"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		apiKey:     cfg.PagerDuty.APIKey,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) Available() bool {
	return c.apiKey != ""
}

type Incident struct {
	ID          string
	Title       string
	Status      string
	Urgency     string
	CreatedAt   string
	Assignee    string
	ServiceName string
}

type OnCall struct {
	UserName   string
	Schedule   string
	Start      string
	End        string
}

func (c *Client) ListIncidents(ctx context.Context, status string) ([]Incident, error) {
	url := "https://api.pagerduty.com/incidents?sort_by=created_at:desc&limit=25"
	if status != "" {
		url += "&statuses[]=" + status
	}

	body, err := c.doGet(ctx, url)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Incidents []struct {
			ID        string `json:"id"`
			Title     string `json:"title"`
			Status    string `json:"status"`
			Urgency   string `json:"urgency"`
			CreatedAt string `json:"created_at"`
			Assignments []struct {
				Assignee struct{ Summary string } `json:"assignee"`
			} `json:"assignments"`
			Service struct{ Summary string } `json:"service"`
		} `json:"incidents"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	incidents := make([]Incident, 0, len(resp.Incidents))
	for _, i := range resp.Incidents {
		assignee := ""
		if len(i.Assignments) > 0 {
			assignee = i.Assignments[0].Assignee.Summary
		}
		incidents = append(incidents, Incident{
			ID:          i.ID,
			Title:       i.Title,
			Status:      i.Status,
			Urgency:     i.Urgency,
			CreatedAt:   i.CreatedAt,
			Assignee:    assignee,
			ServiceName: i.Service.Summary,
		})
	}
	return incidents, nil
}

func (c *Client) GetOnCalls(ctx context.Context) ([]OnCall, error) {
	body, err := c.doGet(ctx, "https://api.pagerduty.com/oncalls?limit=25")
	if err != nil {
		return nil, err
	}

	var resp struct {
		Oncalls []struct {
			User     struct{ Summary string } `json:"user"`
			Schedule struct{ Summary string } `json:"schedule"`
			Start    string `json:"start"`
			End      string `json:"end"`
		} `json:"oncalls"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	oncalls := make([]OnCall, 0, len(resp.Oncalls))
	for _, o := range resp.Oncalls {
		oncalls = append(oncalls, OnCall{
			UserName: o.User.Summary,
			Schedule: o.Schedule.Summary,
			Start:    o.Start,
			End:      o.End,
		})
	}
	return oncalls, nil
}

func (c *Client) doGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Token token="+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("pagerduty API %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}
