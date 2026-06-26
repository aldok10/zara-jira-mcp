package clockify

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
	apiKey      string
	workspaceID string
	httpClient  *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		apiKey:      cfg.Clockify.APIKey,
		workspaceID: cfg.Clockify.WorkspaceID,
		httpClient:  &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) Available() bool {
	return c.apiKey != "" && c.workspaceID != ""
}

type TimeEntry struct {
	ID          string
	Description string
	UserName    string
	ProjectName string
	Start       time.Time
	End         time.Time
	Duration    time.Duration
}

func (c *Client) GetTimeEntries(ctx context.Context, start, end time.Time) ([]TimeEntry, error) {
	// Use the reports endpoint which gives all workspace entries
	url := fmt.Sprintf("https://reports.api.clockify.me/v1/workspaces/%s/reports/detailed", c.workspaceID)

	body, err := c.doPost(ctx, url, fmt.Sprintf(`{"dateRangeStart":"%s","dateRangeEnd":"%s","detailedFilter":{"page":1,"pageSize":200}}`,
		start.Format(time.RFC3339), end.Format(time.RFC3339)))
	if err != nil {
		return nil, err
	}

	var resp struct {
		TimeEntries []struct {
			ID          string `json:"_id"`
			Description string `json:"description"`
			UserName    string `json:"userName"`
			ProjectName string `json:"projectName"`
			TimeInterval struct {
				Start    string `json:"start"`
				End      string `json:"end"`
				Duration int64  `json:"duration"`
			} `json:"timeInterval"`
		} `json:"timeentries"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	entries := make([]TimeEntry, 0, len(resp.TimeEntries))
	for _, e := range resp.TimeEntries {
		s, _ := time.Parse(time.RFC3339, e.TimeInterval.Start)
		en, _ := time.Parse(time.RFC3339, e.TimeInterval.End)
		entries = append(entries, TimeEntry{
			ID:          e.ID,
			Description: e.Description,
			UserName:    e.UserName,
			ProjectName: e.ProjectName,
			Start:       s,
			End:         en,
			Duration:    time.Duration(e.TimeInterval.Duration) * time.Second,
		})
	}
	return entries, nil
}

func (c *Client) GetSummaryReport(ctx context.Context, start, end time.Time) (map[string]map[string]time.Duration, error) {
	url := fmt.Sprintf("https://reports.api.clockify.me/v1/workspaces/%s/reports/summary", c.workspaceID)

	body, err := c.doPost(ctx, url, fmt.Sprintf(`{"dateRangeStart":"%s","dateRangeEnd":"%s","summaryFilter":{"groups":["USER","PROJECT"]}}`,
		start.Format(time.RFC3339), end.Format(time.RFC3339)))
	if err != nil {
		return nil, err
	}

	var resp struct {
		GroupOne []struct {
			Name     string `json:"name"`
			Duration int64  `json:"duration"`
			Children []struct {
				Name     string `json:"name"`
				Duration int64  `json:"duration"`
			} `json:"children"`
		} `json:"groupOne"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	// user -> project -> duration
	result := make(map[string]map[string]time.Duration)
	for _, user := range resp.GroupOne {
		projects := make(map[string]time.Duration)
		for _, proj := range user.Children {
			projects[proj.Name] = time.Duration(proj.Duration) * time.Second
		}
		result[user.Name] = projects
	}
	return result, nil
}

func (c *Client) doPost(ctx context.Context, url, payload string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBufferString(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("clockify API %d: %s", resp.StatusCode, string(b))
	}
	return b, nil
}
