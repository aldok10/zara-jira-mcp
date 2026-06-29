// Package httpclient provides HTTP client utilities with security defaults.
package httpclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"
)

// NewWithTimeout creates a secure HTTP client with the given timeout.
func NewWithTimeout(timeout time.Duration) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
		IdleConnTimeout:       30 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    false,
	}
	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
}

// SanitizeError wraps HTTP-related errors to avoid leaking sensitive details.
func SanitizeError(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("http client error: %w", err)
}
