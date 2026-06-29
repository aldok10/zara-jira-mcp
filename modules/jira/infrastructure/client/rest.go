package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	jira "github.com/felixgeelhaar/jirasdk"
	"github.com/felixgeelhaar/jirasdk/core/issue"

	domain "github.com/aldok10/zara-jira-mcp/modules/jira/domain"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
)

// Retry configuration
const (
	maxRetries  = 3
	initialWait = 100 * time.Millisecond
	maxWait     = 2 * time.Second
)

type rateLimiter struct {
	mu     sync.Mutex
	tokens int
	max    int
	last   time.Time
}

func (r *rateLimiter) wait() {
	r.mu.Lock()
	defer r.mu.Unlock()
	refill := int(time.Since(r.last).Seconds() / 60.0 * float64(r.max))
	if refill > 0 {
		r.tokens += refill
		if r.tokens > r.max {
			r.tokens = r.max
		}
		r.last = time.Now()
	}
	if r.tokens <= 0 {
		time.Sleep(time.Second)
		r.tokens = 1
	}
	r.tokens--
}

// RestClient wraps jirasdk.Client and implements domain.Client.
type RestClient struct {
	sdk     *jira.Client
	baseURL string
	email   string
	token   string
	http    *http.Client
	limiter *rateLimiter
}

// drainAndClose ensures response body is properly drained and closed.
func drainAndClose(resp *http.Response) {
	if resp == nil || resp.Body == nil {
		return
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}

// doRequest is the unified HTTP request helper with retry, rate limiting, and error handling.
// Returns response body, status code, or error.
func (c *RestClient) doRequest(ctx context.Context, method, path string, payload []byte) ([]byte, int, error) {
	// Check context cancellation first
	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	default:
	}

	c.limiter.wait()
	fullURL := c.baseURL + path

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, 0, ctx.Err()
			case <-time.After(backoff(attempt)):
			}
		}

		var bodyReader io.Reader
		if len(payload) > 0 {
			bodyReader = bytes.NewReader(payload)
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
		if err != nil {
			return nil, 0, fmt.Errorf("create request: %w", err)
		}
		req.SetBasicAuth(c.email, c.token)
		req.Header.Set("Accept", "application/json")
		if len(payload) > 0 {
			req.Header.Set("Content-Type", "application/json")
		}

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("%s %s: %w", method, path, err)
			slog.Warn("request failed, retrying", "method", method, "path", path, "attempt", attempt+1, "error", err)
			continue
		}

		if resp.StatusCode < 400 || resp.StatusCode == 429 {
			// Read body (limit to 10MB)
			body, rerr := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
			drainAndClose(resp)
			if rerr != nil {
				return nil, 0, fmt.Errorf("read response: %w", rerr)
			}
			return body, resp.StatusCode, nil
		}

		// Error responses
		respBody, rerr := io.ReadAll(io.LimitReader(resp.Body, 1<<16))
		drainAndClose(resp)
		if rerr != nil {
			return nil, 0, fmt.Errorf("%s %s: read error body: %w", method, path, rerr)
		}

		if resp.StatusCode >= 500 || resp.StatusCode == 429 {
			lastErr = fmt.Errorf("%s %s: status %d: %s", method, path, resp.StatusCode, string(respBody))
			slog.Warn("server error, retrying", "method", method, "path", path, "attempt", attempt+1, "status", resp.StatusCode)
			continue
		}

		// 4xx (except 429) are fatal — no retry
		return nil, resp.StatusCode, fmt.Errorf("%s %s: status %d: %s", method, path, resp.StatusCode, string(respBody))
	}

	return nil, 0, fmt.Errorf("%s %s: max retries exceeded: %w", method, path, lastErr)
}

// backoff computes exponential backoff duration for retry attempt.
func backoff(attempt int) time.Duration {
	d := initialWait * (1 << attempt)
	if d > maxWait {
		d = maxWait
	}
	return d
}

// NewRestClient creates a new RestClient with secure defaults.
func NewRestClient(cfg *config.Config) (*RestClient, error) {
	secureTransport := &http.Transport{
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: false},
		IdleConnTimeout:       30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableKeepAlives:     false,
		MaxIdleConns:          100,
		MaxIdleConnsPerHost:   10,
	}
	secureHTTPClient := &http.Client{
		Transport: secureTransport,
		Timeout:   30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 5 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	client, err := jira.NewClient(
		jira.WithBaseURL(cfg.Jira.BaseURL),
		jira.WithAPIToken(cfg.Jira.Email, cfg.Jira.Token),
		jira.WithTimeout(30*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("create jira client: %w", err)
	}
	return &RestClient{
		sdk:     client,
		baseURL: cfg.Jira.BaseURL,
		email:   cfg.Jira.Email,
		token:   cfg.Jira.Token,
		http:    secureHTTPClient,
		limiter: &rateLimiter{max: 60, tokens: 60, last: time.Now()},
	}, nil
}

// mapIssue converts a jirasdk issue to domain.Issue.
func mapIssue(raw *issue.Issue) domain.Issue {
	if raw == nil {
		return domain.Issue{}
	}
	i := domain.Issue{
		Key:         raw.Key,
		Summary:     raw.GetSummary(),
		Description: raw.GetDescriptionText(),
		Status:      raw.GetStatusName(),
		Priority:    raw.GetPriorityName(),
		Type:        raw.GetIssueTypeName(),
		Assignee:    raw.GetAssigneeName(),
		Reporter:    raw.GetReporterName(),
		Labels:      raw.GetLabels(),
		Created:     raw.GetCreatedTime(),
		Updated:     raw.GetUpdatedTime(),
	}
	// Extract story points from common custom field IDs
	if raw.Fields != nil && raw.Fields.Custom != nil {
		for _, fieldID := range storyPointFields {
			if sp, ok := raw.Fields.Custom.GetNumber(fieldID); ok && sp > 0 {
				i.StoryPoints = sp
				break
			}
		}
	}
	return i
}

// storyPointFields lists common Jira custom field IDs for story points.
// The first match wins.
var storyPointFields = []string{
	"story_points",      // next-gen projects
	"customfield_10016", // Jira Cloud default
	"customfield_10028", // common alternative
	"customfield_10004", // some instances
	"customfield_10014", // another variant
}

// Compile-time interface compliance check.
var _ domain.Client = (*RestClient)(nil)
