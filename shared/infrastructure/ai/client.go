package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	domain "github.com/aldok10/zara-jira-mcp/shared/domain/ai"
	"github.com/aldok10/zara-jira-mcp/shared/infrastructure/config"
)

// aiRateLimiter is a token bucket rate limiter for AI API calls.
type aiRateLimiter struct {
	mu     sync.Mutex
	tokens float64
	max    float64
	last   time.Time
	refill float64 // tokens per second
}

func newAIRateLimiter(ratePerMinute int) *aiRateLimiter {
	return &aiRateLimiter{
		tokens: float64(ratePerMinute),
		max:    float64(ratePerMinute),
		last:   time.Now(),
		refill: float64(ratePerMinute) / 60.0,
	}
}

func (r *aiRateLimiter) wait(ctx context.Context) error {
	r.mu.Lock()
	now := time.Now()
	elapsed := now.Sub(r.last).Seconds()
	r.last = now
	r.tokens = math.Min(r.max, r.tokens+elapsed*r.refill)

	if r.tokens < 1 {
		// Need to wait for refill
		waitDur := time.Duration((1 - r.tokens) / r.refill * float64(time.Second))
		r.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(waitDur):
		}
		r.mu.Lock()
		r.tokens--
		r.mu.Unlock()
		return nil
	}

	r.tokens--
	r.mu.Unlock()
	return nil
}

// OpenAIClient supports OpenAI-compatible, Anthropic, and Google Gemini APIs.
// Provider is auto-detected from JIRA_AI_BASE_URL.
type OpenAIClient struct {
	baseURL    string
	apiKey     string
	model      string
	provider   string // "openai", "anthropic", "gemini"
	httpClient *http.Client
	limiter    *aiRateLimiter // 10 req/min by default
}

// Ensure OpenAIClient satisfies domain.Provider at compile time.
var _ domain.Provider = (*OpenAIClient)(nil)

func NewOpenAIClient(cfg *config.Config) *OpenAIClient {
	provider := detectProvider(cfg.AI.BaseURL)
	return &OpenAIClient{
		baseURL:    strings.TrimRight(cfg.AI.BaseURL, "/"),
		apiKey:     cfg.AI.APIKey,
		model:      cfg.AI.Model,
		provider:   provider,
		httpClient: &http.Client{Timeout: 60 * time.Second},
		limiter:    newAIRateLimiter(10), // 10 AI requests per minute
	}
}

func detectProvider(baseURL string) string {
	switch {
	case strings.Contains(baseURL, "anthropic"):
		return "anthropic"
	case strings.Contains(baseURL, "generativelanguage.googleapis.com"):
		return "gemini"
	default:
		return "openai"
	}
}

func (c *OpenAIClient) Complete(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	if c.baseURL == "" || c.apiKey == "" {
		return "", fmt.Errorf("AI not configured: set JIRA_AI_BASE_URL and JIRA_AI_API_KEY")
	}

	switch c.provider {
	case "anthropic":
		return c.completeAnthropic(ctx, systemPrompt, userPrompt)
	case "gemini":
		return c.completeGemini(ctx, systemPrompt, userPrompt)
	default:
		return c.completeOpenAI(ctx, systemPrompt, userPrompt)
	}
}

func (c *OpenAIClient) completeOpenAI(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	payload := map[string]any{
		"model": c.model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature": 0.3,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	respBody, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("AI returned no choices")
	}
	return result.Choices[0].Message.Content, nil
}

func (c *OpenAIClient) completeAnthropic(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	payload := map[string]any{
		"model":      c.model,
		"max_tokens": 4096,
		"system":     systemPrompt,
		"messages": []map[string]string{
			{"role": "user", "content": userPrompt},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	respBody, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	var result struct {
		Content []struct {
			Text string `json:"text"`
		} `json:"content"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	if len(result.Content) == 0 {
		return "", fmt.Errorf("AI returned no content")
	}
	return result.Content[0].Text, nil
}

func (c *OpenAIClient) completeGemini(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	payload := map[string]any{
		"contents": []map[string]any{
			{"role": "user", "parts": []map[string]string{{"text": userPrompt}}},
		},
		"systemInstruction": map[string]any{
			"parts": []map[string]string{{"text": systemPrompt}},
		},
		"generationConfig": map[string]any{
			"temperature": 0.3,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}
	url := fmt.Sprintf("%s/v1beta/models/%s:generateContent", c.baseURL, c.model)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", c.apiKey)

	respBody, err := c.doRequest(req)
	if err != nil {
		return "", err
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("AI returned no content")
	}
	return result.Candidates[0].Content.Parts[0].Text, nil
}

func (c *OpenAIClient) doRequest(req *http.Request) ([]byte, error) {
	// Apply AI rate limiter before making the request
	if err := c.limiter.wait(req.Context()); err != nil {
		return nil, fmt.Errorf("ai rate limit: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Limit response size to 10MB to prevent resource exhaustion
	body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("AI API error %d: %s", resp.StatusCode, string(body))
	}
	return body, nil
}
