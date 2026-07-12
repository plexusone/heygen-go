// Package heygen provides a Go client for the HeyGen API.
package heygen

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

const (
	// DefaultBaseURL is the default HeyGen API base URL.
	DefaultBaseURL = "https://api.heygen.com"

	// LiveAvatarBaseURL is the LiveAvatar API base URL.
	// LiveAvatar is a separate service with its own API key.
	LiveAvatarBaseURL = "https://api.liveavatar.com"

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second

	// EnvAPIKey is the environment variable for the HeyGen API key.
	EnvAPIKey = "HEYGEN_API_KEY"

	// EnvLiveAvatarAPIKey is the environment variable for the LiveAvatar API key.
	// Note: LiveAvatar and HeyGen use DIFFERENT API keys.
	// Get your LiveAvatar key from: app.liveavatar.com/developers
	EnvLiveAvatarAPIKey = "LIVEAVATAR_API_KEY"
)

// Config holds configuration for the HeyGen client.
type Config struct {
	// APIKey is the HeyGen API key for authentication.
	// If empty, reads from HEYGEN_API_KEY environment variable.
	APIKey string

	// BaseURL is the API base URL.
	// Defaults to https://api.heygen.com
	BaseURL string

	// HTTPClient is the HTTP client to use for requests.
	// If nil, a default client with timeout is created.
	HTTPClient *http.Client

	// Logger is the structured logger for debug output.
	// If nil, logging is disabled.
	Logger *slog.Logger

	// Retry configures retry behavior for transient failures.
	Retry RetryConfig

	// UserAgent is the User-Agent header value.
	// Defaults to "heygen-go/VERSION".
	UserAgent string
}

// RetryConfig configures retry behavior.
type RetryConfig struct {
	// MaxRetries is the maximum number of retry attempts.
	// Default: 2
	MaxRetries int

	// BaseDelay is the initial delay between retries.
	// Default: 1s
	BaseDelay time.Duration

	// MaxDelay is the maximum delay between retries.
	// Default: 30s
	MaxDelay time.Duration
}

// DefaultRetryConfig returns the default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries: 2,
		BaseDelay:  1 * time.Second,
		MaxDelay:   30 * time.Second,
	}
}

// applyDefaults fills in default values for empty config fields.
func (c *Config) applyDefaults() {
	if c.APIKey == "" {
		c.APIKey = os.Getenv(EnvAPIKey)
	}
	if c.BaseURL == "" {
		c.BaseURL = DefaultBaseURL
	}
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{
			Timeout: DefaultTimeout,
		}
	}
	if c.UserAgent == "" {
		c.UserAgent = "heygen-go/" + Version
	}
	if c.Retry.MaxRetries == 0 {
		c.Retry = DefaultRetryConfig()
	}
}

// Version is the SDK version.
const Version = "0.1.0"
