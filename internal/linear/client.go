package linear

import (
	"bytes"
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
		apiKey:     cfg.Linear.APIKey,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) Available() bool {
	return c.apiKey != ""
}

type Issue struct {
	ID       string
	Title    string
	State    string
	Assignee string
	Priority int
	Team     string
}

type Cycle struct {
	ID        string
	Name      string
	Number    int
	StartsAt  string
	EndsAt    string
	Progress  float64
	IssueCount int
}

type Activity struct {
	Type      string
	CreatedAt string
	Issue     string
	Actor     string
}

func (c *Client) ListIssues(ctx context.Context, teamKey, state string) ([]Issue, error) {
	filter := ""
	if teamKey != "" {
		filter += fmt.Sprintf(`, filter: { team: { key: { eq: "%s" } }`, teamKey)
		if state != "" {
			filter += fmt.Sprintf(`, state: { name: { eq: "%s" } }`, state)
		}
		filter += " }"
	} else if state != "" {
		filter += fmt.Sprintf(`, filter: { state: { name: { eq: "%s" } } }`, state)
	}

	query := fmt.Sprintf(`{ "query": "{ issues(first: 50%s) { nodes { id title state { name } assignee { name } priority priorityLabel team { name } } } }" }`, filter)

	body, err := c.doGraphQL(ctx, query)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data struct {
			Issues struct {
				Nodes []struct {
					ID       string `json:"id"`
					Title    string `json:"title"`
					State    struct{ Name string } `json:"state"`
					Assignee *struct{ Name string } `json:"assignee"`
					Priority int    `json:"priority"`
					Team     struct{ Name string } `json:"team"`
				} `json:"nodes"`
			} `json:"issues"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	issues := make([]Issue, 0, len(resp.Data.Issues.Nodes))
	for _, n := range resp.Data.Issues.Nodes {
		assignee := ""
		if n.Assignee != nil {
			assignee = n.Assignee.Name
		}
		issues = append(issues, Issue{
			ID:       n.ID,
			Title:    n.Title,
			State:    n.State.Name,
			Assignee: assignee,
			Priority: n.Priority,
			Team:     n.Team.Name,
		})
	}
	return issues, nil
}

func (c *Client) ListCycles(ctx context.Context) ([]Cycle, error) {
	query := `{ "query": "{ cycles(first: 10, orderBy: createdAt) { nodes { id name number startsAt endsAt progress completedIssueCountAfterEachDay } } }" }`

	body, err := c.doGraphQL(ctx, query)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data struct {
			Cycles struct {
				Nodes []struct {
					ID       string  `json:"id"`
					Name     string  `json:"name"`
					Number   int     `json:"number"`
					StartsAt string  `json:"startsAt"`
					EndsAt   string  `json:"endsAt"`
					Progress float64 `json:"progress"`
				} `json:"nodes"`
			} `json:"cycles"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	cycles := make([]Cycle, 0, len(resp.Data.Cycles.Nodes))
	for _, n := range resp.Data.Cycles.Nodes {
		cycles = append(cycles, Cycle{
			ID:       n.ID,
			Name:     n.Name,
			Number:   n.Number,
			StartsAt: n.StartsAt,
			EndsAt:   n.EndsAt,
			Progress: n.Progress,
		})
	}
	return cycles, nil
}

func (c *Client) RecentActivity(ctx context.Context) ([]Activity, error) {
	query := `{ "query": "{ issueHistory(first: 20) { nodes { createdAt fromState { name } toState { name } issue { title identifier } actor { name } } } }" }`

	body, err := c.doGraphQL(ctx, query)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data struct {
			IssueHistory struct {
				Nodes []struct {
					CreatedAt string `json:"createdAt"`
					FromState *struct{ Name string } `json:"fromState"`
					ToState   *struct{ Name string } `json:"toState"`
					Issue     struct {
						Title      string `json:"title"`
						Identifier string `json:"identifier"`
					} `json:"issue"`
					Actor *struct{ Name string } `json:"actor"`
				} `json:"nodes"`
			} `json:"issueHistory"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	activities := make([]Activity, 0, len(resp.Data.IssueHistory.Nodes))
	for _, n := range resp.Data.IssueHistory.Nodes {
		actType := "updated"
		if n.FromState != nil && n.ToState != nil {
			actType = n.FromState.Name + " -> " + n.ToState.Name
		}
		actor := ""
		if n.Actor != nil {
			actor = n.Actor.Name
		}
		activities = append(activities, Activity{
			Type:      actType,
			CreatedAt: n.CreatedAt,
			Issue:     n.Issue.Identifier + " " + n.Issue.Title,
			Actor:     actor,
		})
	}
	return activities, nil
}

func (c *Client) doGraphQL(ctx context.Context, query string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.linear.app/graphql", bytes.NewBufferString(query))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("linear API %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}
