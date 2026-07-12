// Package liveavatar provides access to the LiveAvatar real-time streaming API.
//
// LiveAvatar is HeyGen's real-time avatar streaming product for live video
// conversations with AI avatars. It requires a separate API key from HeyGen.
//
// Get your API key from: https://app.liveavatar.com/developers
package liveavatar

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// DefaultBaseURL is the LiveAvatar API base URL.
	DefaultBaseURL = "https://api.liveavatar.com"

	// EnvAPIKey is the environment variable for the LiveAvatar API key.
	EnvAPIKey = "LIVEAVATAR_API_KEY" //nolint:gosec // G101: not a credential, just env var name

	// DefaultTimeout is the default HTTP client timeout.
	DefaultTimeout = 30 * time.Second

	// SandboxAvatarID is a public sample avatar ID for testing.
	// Sandbox sessions are limited to 60 seconds and don't consume credits.
	SandboxAvatarID = "65f9e3c9-d48b-4118-b73a-4ae2e3cbb8f0"
)

// VideoQuality specifies the avatar video quality.
type VideoQuality string

const (
	QualityVeryHigh VideoQuality = "very_high"
	QualityHigh     VideoQuality = "high"
	QualityMedium   VideoQuality = "medium"
	QualityLow      VideoQuality = "low"
)

// Client provides access to LiveAvatar APIs.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Config holds configuration for the LiveAvatar client.
type Config struct {
	// APIKey is the LiveAvatar API key for authentication.
	// If empty, reads from LIVEAVATAR_API_KEY environment variable.
	APIKey string

	// BaseURL is the API base URL.
	// Defaults to https://api.liveavatar.com
	BaseURL string

	// HTTPClient is the HTTP client to use for requests.
	// If nil, a default client with timeout is created.
	HTTPClient *http.Client
}

// NewClient creates a new LiveAvatar client.
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	apiKey := cfg.APIKey
	if apiKey == "" {
		apiKey = os.Getenv(EnvAPIKey)
	}
	if apiKey == "" {
		return nil, fmt.Errorf("liveavatar: API key required (set %s or pass APIKey in config)", EnvAPIKey)
	}

	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: DefaultTimeout}
	}

	return &Client{
		apiKey:     apiKey,
		baseURL:    baseURL,
		httpClient: httpClient,
	}, nil
}

// LiveKitConfig contains LiveKit room configuration for LITE mode.
type LiveKitConfig struct {
	// LiveKitURL is the LiveKit server URL (wss://...).
	LiveKitURL string `json:"livekit_url"`

	// LiveKitRoom is the room name the avatar should join.
	LiveKitRoom string `json:"livekit_room"`

	// LiveKitClientToken is a LiveKit JWT token with permissions:
	// room_join=true, can_publish=true, can_subscribe=true, can_publish_data=true
	LiveKitClientToken string `json:"livekit_client_token"`
}

// NewSessionRequest contains parameters for creating a streaming session.
type NewSessionRequest struct {
	// Mode is the session mode. Use "LITE" for BYO AI stack with LiveKit.
	Mode string `json:"mode"`

	// AvatarID is the UUID of the avatar to use.
	// Use SandboxAvatarID for testing without credits.
	AvatarID string `json:"avatar_id"`

	// IsSandbox enables sandbox mode (60s limit, no credit usage).
	IsSandbox bool `json:"is_sandbox,omitempty"`

	// VideoQuality sets the avatar video quality.
	VideoQuality VideoQuality `json:"video_quality,omitempty"`

	// LiveKitConfig contains LiveKit room configuration (required for LITE mode).
	LiveKitConfig *LiveKitConfig `json:"livekit_config,omitempty"`
}

// NewSessionResponse contains the session token response.
type NewSessionResponse struct {
	// SessionID is the unique session identifier.
	SessionID string `json:"session_id"`

	// SessionToken is the JWT token for starting/stopping the session.
	SessionToken string `json:"session_token"`
}

// StartSessionResponse contains the session start response.
type StartSessionResponse struct {
	// SessionID is the unique session identifier.
	SessionID string `json:"session_id"`

	// WSURL is the WebSocket URL for streaming audio/events.
	WSURL string `json:"ws_url"`

	// MaxSessionDuration is the maximum session duration in seconds.
	MaxSessionDuration int `json:"max_session_duration"`
}

// apiResponse wraps API responses.
type apiResponse struct {
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

// NewSession creates a new streaming session and returns session credentials.
func (c *Client) NewSession(ctx context.Context, req *NewSessionRequest) (*NewSessionResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("liveavatar: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/sessions/token", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("liveavatar: create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-KEY", c.apiKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("liveavatar: send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("liveavatar: read response: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("liveavatar: unmarshal response: %w", err)
	}

	if apiResp.Code != 1000 {
		return nil, fmt.Errorf("liveavatar: API error (code=%d): %s", apiResp.Code, apiResp.Message)
	}

	var result NewSessionResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("liveavatar: unmarshal session data: %w", err)
	}

	return &result, nil
}

// StartSession starts a streaming session and returns the WebSocket URL.
func (c *Client) StartSession(ctx context.Context, sessionID, sessionToken string) (*StartSessionResponse, error) {
	body, err := json.Marshal(map[string]string{"session_id": sessionID})
	if err != nil {
		return nil, fmt.Errorf("liveavatar: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/sessions/start", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("liveavatar: create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+sessionToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("liveavatar: send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("liveavatar: read response: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("liveavatar: unmarshal response: %w", err)
	}

	if apiResp.Code != 1000 {
		return nil, fmt.Errorf("liveavatar: API error (code=%d): %s", apiResp.Code, apiResp.Message)
	}

	var result StartSessionResponse
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("liveavatar: unmarshal session data: %w", err)
	}

	return &result, nil
}

// StopReason specifies why a session was stopped.
type StopReason string

const (
	StopReasonUserDisconnected StopReason = "USER_DISCONNECTED"
	StopReasonSessionEnded     StopReason = "SESSION_ENDED"
)

// StopSession stops a streaming session.
func (c *Client) StopSession(ctx context.Context, sessionID, sessionToken string, reason StopReason) error {
	body, err := json.Marshal(map[string]string{
		"session_id": sessionID,
		"reason":     string(reason),
	})
	if err != nil {
		return fmt.Errorf("liveavatar: marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/sessions/stop", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("liveavatar: create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+sessionToken)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("liveavatar: send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("liveavatar: read response: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("liveavatar: unmarshal response: %w", err)
	}

	if apiResp.Code != 1000 {
		return fmt.Errorf("liveavatar: API error (code=%d): %s", apiResp.Code, apiResp.Message)
	}

	return nil
}
