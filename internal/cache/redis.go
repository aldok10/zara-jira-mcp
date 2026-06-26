package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client wraps a Redis client. Nil-safe: all methods are no-ops when unavailable.
type Client struct {
	rdb *redis.Client
	ttl time.Duration
}

// NewClient connects to Redis using the given URL. Returns a nil-safe client
// if URL is empty or connection fails.
func NewClient(url string, ttl time.Duration) *Client {
	if url == "" {
		return &Client{}
	}
	opts, err := redis.ParseURL(url)
	if err != nil {
		return &Client{}
	}
	rdb := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if rdb.Ping(ctx).Err() != nil {
		return &Client{}
	}
	return &Client{rdb: rdb, ttl: ttl}
}

// Available reports whether a Redis connection is active.
func (c *Client) Available() bool {
	return c != nil && c.rdb != nil
}

// Get retrieves a string value by key.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	if !c.Available() {
		return "", redis.Nil
	}
	return c.rdb.Get(ctx, key).Result()
}

// Set stores a string value with the given TTL.
func (c *Client) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if !c.Available() {
		return nil
	}
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

// Delete removes a key.
func (c *Client) Delete(ctx context.Context, key string) error {
	if !c.Available() {
		return nil
	}
	return c.rdb.Del(ctx, key).Err()
}

// GetJSON retrieves and unmarshals a cached JSON value into dest.
func (c *Client) GetJSON(ctx context.Context, key string, dest any) error {
	val, err := c.Get(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// SetJSON marshals value to JSON and caches it with the given TTL.
func (c *Client) SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(ctx, key, string(data), ttl)
}
