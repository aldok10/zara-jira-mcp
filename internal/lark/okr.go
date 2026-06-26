package lark

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/aldok10/zara-jira-mcp/config"
)

// OKRClient talks to Lark OKR API via plain HTTP with tenant_access_token.
type OKRClient struct {
	appID      string
	appSecret  string
	httpClient *http.Client
	token      string
	tokenExp   time.Time
}

func NewOKRClient(cfg *config.Config) *OKRClient {
	return &OKRClient{
		appID:      cfg.Lark.AppID,
		appSecret:  cfg.Lark.AppSecret,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *OKRClient) Available() bool {
	return c.appID != "" && c.appSecret != ""
}

// OKR API models

type Period struct {
	ID     string `json:"period_id"`
	Name   string `json:"zh_name"`
	EnName string `json:"en_name"`
	Status int    `json:"status"` // 1=in progress, 2=not started, 3=ended
}

type Objective struct {
	ID         string      `json:"id"`
	Content    string      `json:"content"`
	Progress   int         `json:"progress_rate"`
	KeyResults []KeyResult `json:"kr_list"`
}

type KeyResult struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	Progress int    `json:"progress_rate"`
}

// ListPeriods returns OKR periods.
func (c *OKRClient) ListPeriods(ctx context.Context) ([]Period, error) {
	body, err := c.get(ctx, "https://open.larksuite.com/open-apis/okr/v1/periods?page_size=20")
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data struct {
			Items []Period `json:"items"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse periods: %w", err)
	}
	return resp.Data.Items, nil
}

// ListUserOKRs returns OKRs for a user in a period.
func (c *OKRClient) ListUserOKRs(ctx context.Context, userID, periodID string) ([]Objective, error) {
	url := fmt.Sprintf("https://open.larksuite.com/open-apis/okr/v1/users/%s/okrs?period_id=%s&user_id_type=open_id", userID, periodID)
	body, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data struct {
			OKRList []struct {
				ObjectiveList []Objective `json:"objective_list"`
			} `json:"okr_list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse user okrs: %w", err)
	}

	var objectives []Objective
	for _, okr := range resp.Data.OKRList {
		objectives = append(objectives, okr.ObjectiveList...)
	}
	return objectives, nil
}

// BatchGetOKRs gets OKR details by IDs.
func (c *OKRClient) BatchGetOKRs(ctx context.Context, okrIDs []string) ([]Objective, error) {
	url := fmt.Sprintf("https://open.larksuite.com/open-apis/okr/v1/okrs/batch_get?okr_ids=%s", strings.Join(okrIDs, "&okr_ids="))
	body, err := c.get(ctx, url)
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data struct {
			OKRList []struct {
				ObjectiveList []Objective `json:"objective_list"`
			} `json:"okr_list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse batch okrs: %w", err)
	}

	var objectives []Objective
	for _, okr := range resp.Data.OKRList {
		objectives = append(objectives, okr.ObjectiveList...)
	}
	return objectives, nil
}


// CreateProgressRecord posts a progress update to a Lark OKR Key Result.
func (c *OKRClient) CreateProgressRecord(ctx context.Context, krID string, content string) error {
	token, err := c.getTenantToken(ctx)
	if err != nil {
		return err
	}
	escaped := strings.ReplaceAll(content, `"`, `\"`)
	payload := fmt.Sprintf(`{"target_id":"%s","target_type":2,"content":{"blocks":[{"type":"paragraph","content":[{"type":"text","text":"%s"}]}]}}`, krID, escaped)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://open.larksuite.com/open-apis/okr/v1/progress_records", strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var base struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(body, &base); err == nil && base.Code != 0 {
		return fmt.Errorf("lark OKR progress error %d: %s", base.Code, base.Msg)
	}
	return nil
}
// Internal HTTP helpers

func (c *OKRClient) get(ctx context.Context, url string) ([]byte, error) {
	token, err := c.getTenantToken(ctx)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("lark OKR API %d: %s", resp.StatusCode, string(body))
	}

	// Check Lark-level error
	var base struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.Unmarshal(body, &base); err == nil && base.Code != 0 {
		return nil, fmt.Errorf("lark OKR error %d: %s", base.Code, base.Msg)
	}

	return body, nil
}

func (c *OKRClient) getTenantToken(ctx context.Context) (string, error) {
	if c.token != "" && time.Now().Before(c.tokenExp) {
		return c.token, nil
	}

	payload := fmt.Sprintf(`{"app_id":"%s","app_secret":"%s"}`, c.appID, c.appSecret)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://open.larksuite.com/open-apis/auth/v3/tenant_access_token/internal",
		strings.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Code              int    `json:"code"`
		Msg               string `json:"msg"`
		TenantAccessToken string `json:"tenant_access_token"`
		Expire            int    `json:"expire"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Code != 0 {
		return "", fmt.Errorf("tenant token error %d: %s", result.Code, result.Msg)
	}

	c.token = result.TenantAccessToken
	c.tokenExp = time.Now().Add(time.Duration(result.Expire-60) * time.Second)
	return c.token, nil
}
