package heygen

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Client is the HeyGen API client.
type Client struct {
	config     Config
	httpClient *http.Client
}

// NewClient creates a new HeyGen API client.
func NewClient(cfg Config) (*Client, error) {
	cfg.applyDefaults()

	if cfg.APIKey == "" {
		return nil, ErrNoAPIKey
	}

	return &Client{
		config:     cfg,
		httpClient: cfg.HTTPClient,
	}, nil
}

// APIKey returns the configured API key.
func (c *Client) APIKey() string {
	return c.config.APIKey
}

// BaseURL returns the configured base URL.
func (c *Client) BaseURL() string {
	return c.config.BaseURL
}

// do executes an HTTP request with retry logic.
func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.Retry.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := c.calculateBackoff(attempt)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		attemptReq := req.Clone(ctx)
		if attempt > 0 && req.GetBody != nil {
			body, err := req.GetBody()
			if err != nil {
				return nil, fmt.Errorf("rewind request body for retry: %w", err)
			}
			attemptReq.Body = body
		}

		resp, err := c.httpClient.Do(attemptReq)
		if err != nil {
			lastErr = err
			if c.isRetryable(req.Method, 0) {
				continue
			}
			return nil, err
		}

		// Check for retryable status codes
		if c.isRetryable(req.Method, resp.StatusCode) {
			// Handle Retry-After header
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					delay := time.Duration(seconds) * time.Second
					if delay > c.config.Retry.MaxDelay {
						delay = c.config.Retry.MaxDelay
					}
					_ = resp.Body.Close()
					select {
					case <-ctx.Done():
						return nil, ctx.Err()
					case <-time.After(delay):
					}
					continue
				}
			}

			_ = resp.Body.Close()
			lastErr = fmt.Errorf("request failed with status %d", resp.StatusCode)
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// calculateBackoff returns the backoff duration for the given attempt.
func (c *Client) calculateBackoff(attempt int) time.Duration {
	delay := float64(c.config.Retry.BaseDelay) * math.Pow(2, float64(attempt-1))
	// Add jitter: multiply by random factor between 0.75 and 1.25
	jitter := 0.75 + rand.Float64()*0.5 //nolint:gosec // G404: weak random is fine for jitter
	delay *= jitter

	if delay > float64(c.config.Retry.MaxDelay) {
		delay = float64(c.config.Retry.MaxDelay)
	}
	return time.Duration(delay)
}

// isRetryable returns true if the request should be retried.
func (c *Client) isRetryable(method string, statusCode int) bool {
	// Always retry rate limits
	if statusCode == 429 {
		return true
	}

	// Only retry idempotent methods for server errors
	isIdempotent := method == http.MethodGet ||
		method == http.MethodHead ||
		method == http.MethodPut ||
		method == http.MethodDelete

	if isIdempotent {
		return statusCode == 0 || // Network error
			statusCode == 500 ||
			statusCode == 502 ||
			statusCode == 503 ||
			statusCode == 504
	}

	return false
}

// Request executes an API request and decodes the response.
func (c *Client) Request(ctx context.Context, method, path string, body, result any) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	url := c.config.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Api-Key", c.config.APIKey)
	req.Header.Set("User-Agent", c.config.UserAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	// Check for error responses
	if resp.StatusCode >= 400 {
		return c.parseError(resp, respBody)
	}

	// Decode successful response
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("unmarshal response: %w", err)
		}
	}

	return nil
}

// parseError parses an error response from the API.
func (c *Client) parseError(resp *http.Response, body []byte) error {
	requestID := resp.Header.Get("X-Request-Id")

	// Try to parse v3 API error envelope
	var envelope struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(body, &envelope); err == nil && envelope.Error.Code != "" {
		return &APIError{
			StatusCode: resp.StatusCode,
			Code:       envelope.Error.Code,
			Message:    envelope.Error.Message,
			RequestID:  requestID,
		}
	}

	// Try simpler error format
	var simpleErr struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &simpleErr); err == nil && simpleErr.Message != "" {
		return &APIError{
			StatusCode: resp.StatusCode,
			Code:       simpleErr.Code,
			Message:    simpleErr.Message,
			RequestID:  requestID,
		}
	}

	// Fallback to status-based error
	message := strings.TrimSpace(string(body))
	if message == "" {
		message = http.StatusText(resp.StatusCode)
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    message,
		RequestID:  requestID,
	}
}

// Get performs a GET request.
func (c *Client) Get(ctx context.Context, path string, result any) error {
	return c.Request(ctx, http.MethodGet, path, nil, result)
}

// Post performs a POST request.
func (c *Client) Post(ctx context.Context, path string, body, result any) error {
	return c.Request(ctx, http.MethodPost, path, body, result)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string, result any) error {
	return c.Request(ctx, http.MethodDelete, path, nil, result)
}
