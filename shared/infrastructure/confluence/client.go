package confluence

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
)

// Client wraps Confluence REST API v2.
type Client struct {
	baseURL    string
	email      string
	token      string
	httpClient *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL:    cfg.Confluence.BaseURL,
		email:      cfg.Confluence.Email,
		token:      cfg.Confluence.Token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *Client) Available() bool {
	return c.baseURL != "" && c.token != ""
}

// SearchPages searches Confluence by CQL.
func (c *Client) SearchPages(ctx context.Context, cql string, limit int) ([]Page, error) {
	if limit <= 0 {
		limit = 10
	}
	path := fmt.Sprintf("/wiki/rest/api/content/search?cql=%s&limit=%d", cql, limit)
	body, err := c.doGet(ctx, path)
	if err != nil {
		return nil, err
	}
	var result struct {
		Results []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			Type  string `json:"type"`
			Links struct {
				WebUI string `json:"webui"`
			} `json:"_links"`
			Space struct {
				Key  string `json:"key"`
				Name string `json:"name"`
			} `json:"space"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	pages := make([]Page, len(result.Results))
	for i, r := range result.Results {
		pages[i] = Page{ID: r.ID, Title: r.Title, Type: r.Type, SpaceKey: r.Space.Key, WebURL: c.baseURL + r.Links.WebUI}
	}
	return pages, nil
}

// GetPage gets a page by ID with body content.
func (c *Client) GetPage(ctx context.Context, pageID string) (*PageDetail, error) {
	path := fmt.Sprintf("/wiki/rest/api/content/%s?expand=body.storage,version,space", pageID)
	body, err := c.doGet(ctx, path)
	if err != nil {
		return nil, err
	}
	var result struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Space struct {
			Key string `json:"key"`
		} `json:"space"`
		Version struct {
			Number int `json:"number"`
		} `json:"version"`
		Body struct {
			Storage struct {
				Value string `json:"value"`
			} `json:"storage"`
		} `json:"body"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &PageDetail{
		ID: result.ID, Title: result.Title, SpaceKey: result.Space.Key,
		Version: result.Version.Number, Body: result.Body.Storage.Value,
	}, nil
}

// CreatePage creates a new Confluence page.
func (c *Client) CreatePage(ctx context.Context, spaceKey, title, body string, parentID string) (*Page, error) {
	payload := map[string]any{
		"type":  "page",
		"title": title,
		"space": map[string]string{"key": spaceKey},
		"body": map[string]any{
			"storage": map[string]string{"value": body, "representation": "storage"},
		},
	}
	if parentID != "" {
		payload["ancestors"] = []map[string]string{{"id": parentID}}
	}
	data, _ := json.Marshal(payload)

	respBody, err := c.doPost(ctx, "/wiki/rest/api/content", data)
	if err != nil {
		return nil, err
	}
	var result struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Links struct {
			WebUI string `json:"webui"`
		} `json:"_links"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}
	return &Page{ID: result.ID, Title: result.Title, SpaceKey: spaceKey, WebURL: c.baseURL + result.Links.WebUI}, nil
}

// Page is a simplified Confluence page.
type Page struct {
	ID       string
	Title    string
	Type     string
	SpaceKey string
	WebURL   string
}

// PageDetail includes body content.
type PageDetail struct {
	ID       string
	Title    string
	SpaceKey string
	Version  int
	Body     string
}

func (c *Client) doGet(ctx context.Context, path string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Accept", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("confluence API %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}

func (c *Client) doPost(ctx context.Context, path string, data []byte) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+path, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.email, c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("confluence API %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}
