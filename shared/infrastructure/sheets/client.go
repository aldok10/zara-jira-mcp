package sheets

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		apiKey:     cfg.GoogleSheets.APIKey,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *Client) Available() bool {
	return c.apiKey != ""
}

// ReadRange reads a range from a Google Sheet (must be publicly shared or shared with the API key's project).
func (c *Client) ReadRange(ctx context.Context, spreadsheetID, rangeStr string) ([][]string, error) {
	u := fmt.Sprintf("https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s?key=%s",
		url.PathEscape(spreadsheetID), url.PathEscape(rangeStr), c.apiKey)

	body, err := c.doGet(ctx, u)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Values [][]string `json:"values"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return resp.Values, nil
}

func (c *Client) doGet(ctx context.Context, u string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
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
		return nil, fmt.Errorf("google sheets API %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}
