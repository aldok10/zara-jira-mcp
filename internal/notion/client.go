package notion

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
	databaseID string
	httpClient *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		apiKey:     cfg.Notion.APIKey,
		databaseID: cfg.Notion.DatabaseID,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) Available() bool {
	return c.apiKey != ""
}

func (c *Client) DefaultDatabaseID() string {
	return c.databaseID
}

type SearchResult struct {
	ID    string
	Title string
	Type  string // "page" or "database"
	URL   string
}

// Search searches Notion by keyword.
func (c *Client) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		limit = 10
	}
	payload := map[string]any{
		"query":    query,
		"page_size": limit,
	}
	body, err := c.doPost(ctx, "/v1/search", payload)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Results []struct {
			ID         string `json:"id"`
			Object     string `json:"object"`
			URL        string `json:"url"`
			Properties map[string]struct {
				Title []struct {
					PlainText string `json:"plain_text"`
				} `json:"title"`
			} `json:"properties"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	results := make([]SearchResult, 0, len(resp.Results))
	for _, r := range resp.Results {
		title := ""
		for _, prop := range r.Properties {
			if len(prop.Title) > 0 {
				title = prop.Title[0].PlainText
				break
			}
		}
		results = append(results, SearchResult{
			ID:    r.ID,
			Title: title,
			Type:  r.Object,
			URL:   r.URL,
		})
	}
	return results, nil
}

// CreatePage creates a page in a parent (database or page).
func (c *Client) CreatePage(ctx context.Context, parentID, title, content string) (*SearchResult, error) {
	payload := map[string]any{
		"parent": map[string]string{"database_id": parentID},
		"properties": map[string]any{
			"title": map[string]any{
				"title": []map[string]any{
					{"text": map[string]string{"content": title}},
				},
			},
		},
	}
	if content != "" {
		payload["children"] = []map[string]any{
			{
				"object": "block",
				"type":   "paragraph",
				"paragraph": map[string]any{
					"rich_text": []map[string]any{
						{"type": "text", "text": map[string]string{"content": content}},
					},
				},
			},
		}
	}

	body, err := c.doPost(ctx, "/v1/pages", payload)
	if err != nil {
		return nil, err
	}

	var resp struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &SearchResult{ID: resp.ID, Title: title, Type: "page", URL: resp.URL}, nil
}

// QueryDatabase queries a Notion database with optional filter.
func (c *Client) QueryDatabase(ctx context.Context, databaseID string, filterJSON string, limit int) ([]map[string]any, error) {
	if limit <= 0 {
		limit = 20
	}
	payload := map[string]any{"page_size": limit}
	if filterJSON != "" {
		var filter any
		if err := json.Unmarshal([]byte(filterJSON), &filter); err == nil {
			payload["filter"] = filter
		}
	}

	path := fmt.Sprintf("/v1/databases/%s/query", databaseID)
	body, err := c.doPost(ctx, path, payload)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Results []map[string]any `json:"results"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

func (c *Client) doPost(ctx context.Context, path string, payload any) ([]byte, error) {
	data, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.notion.com"+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("notion API %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}
