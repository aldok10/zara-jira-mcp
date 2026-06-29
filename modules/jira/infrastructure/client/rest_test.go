package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func TestBackoff(t *testing.T) {
	tests := []struct {
		name    string
		attempt int
		wantMin time.Duration
		wantMax time.Duration
	}{
		{name: "attempt 0", attempt: 0, wantMin: 100 * time.Millisecond, wantMax: 100 * time.Millisecond},
		{name: "attempt 1", attempt: 1, wantMin: 200 * time.Millisecond, wantMax: 200 * time.Millisecond},
		{name: "attempt 2", attempt: 2, wantMin: 400 * time.Millisecond, wantMax: 400 * time.Millisecond},
		{name: "attempt 3", attempt: 3, wantMin: 800 * time.Millisecond, wantMax: 800 * time.Millisecond},
		{name: "attempt 4", attempt: 4, wantMin: 1600 * time.Millisecond, wantMax: 1600 * time.Millisecond},
		{name: "attempt 5 (capped)", attempt: 5, wantMin: maxWait, wantMax: maxWait},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := backoff(tt.attempt)
			if got < tt.wantMin || got > tt.wantMax {
				t.Errorf("backoff(%d) = %v, want between %v and %v", tt.attempt, got, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestRateLimiter_ConcurrentAccess(t *testing.T) {
	rl := &rateLimiter{max: 60, tokens: 60, last: time.Now()}

	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			rl.wait()
		}()
	}
	wg.Wait()

	rl.mu.Lock()
	tokens := rl.tokens
	rl.mu.Unlock()

	if tokens > 60 {
		t.Errorf("tokens exceeded max: %d", tokens)
	}
	if tokens < 55 {
		t.Errorf("tokens decreased too much: %d", tokens)
	}
}

func TestRateLimiter_Refill(t *testing.T) {
	rl := &rateLimiter{max: 60, tokens: 0, last: time.Now().Add(-61 * time.Second)}
	rl.wait()

	rl.mu.Lock()
	tokens := rl.tokens
	rl.mu.Unlock()

	if tokens < 58 {
		t.Errorf("expected tokens >= 58 after refill, got %d", tokens)
	}
}

func TestRateLimiter_MaxTokens(t *testing.T) {
	rl := &rateLimiter{max: 60, tokens: 60, last: time.Now().Add(-120 * time.Second)}
	rl.wait()

	rl.mu.Lock()
	tokens := rl.tokens
	rl.mu.Unlock()

	if tokens > 60 {
		t.Errorf("tokens exceeded max: %d", tokens)
	}
}

func TestDrainAndClose_NilResponse(t *testing.T) {
	drainAndClose(nil) // should not panic
}

func TestDrainAndClose_NilBody(t *testing.T) {
	drainAndClose(&http.Response{Body: nil}) // should not panic
}

func TestDoRequest_ContextCancelled(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := &RestClient{
		baseURL: ts.URL,
		http:    ts.Client(),
		limiter: &rateLimiter{max: 60, tokens: 60, last: time.Now()},
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, _, err := client.doRequest(ctx, http.MethodGet, "/test", nil)
	if err == nil {
		t.Error("expected error for cancelled context, got nil")
	}
}

func TestDoRequest_RetryOnServerError(t *testing.T) {
	var attempts int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := &RestClient{
		baseURL: ts.URL,
		http:    ts.Client(),
		limiter: &rateLimiter{max: 60, tokens: 60, last: time.Now()},
	}

	_, status, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if status != http.StatusOK {
		t.Errorf("expected 200, got %d", status)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestDoRequest_FatalOnClientError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	client := &RestClient{
		baseURL: ts.URL,
		http:    ts.Client(),
		limiter: &rateLimiter{max: 60, tokens: 60, last: time.Now()},
	}

	_, status, err := client.doRequest(context.Background(), http.MethodGet, "/test", nil)
	if err == nil {
		t.Error("expected error for 400, got nil")
	}
	if status != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", status)
	}
}
