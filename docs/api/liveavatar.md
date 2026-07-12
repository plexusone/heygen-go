# LiveAvatar API Reference

Package `liveavatar` provides access to the LiveAvatar real-time streaming API.

```go
import "github.com/plexusone/heygen-go/liveavatar"
```

## Constants

```go
const (
    DefaultBaseURL  = "https://api.liveavatar.com"
    EnvAPIKey       = "LIVEAVATAR_API_KEY"
    DefaultTimeout  = 30 * time.Second
    SandboxAvatarID = "65f9e3c9-d48b-4118-b73a-4ae2e3cbb8f0"
)
```

### Video Quality

```go
type VideoQuality string

const (
    QualityVeryHigh VideoQuality = "very_high"
    QualityHigh     VideoQuality = "high"
    QualityMedium   VideoQuality = "medium"
    QualityLow      VideoQuality = "low"
)
```

### Stop Reasons

```go
type StopReason string

const (
    StopReasonUserDisconnected StopReason = "USER_DISCONNECTED"
    StopReasonSessionEnded     StopReason = "SESSION_ENDED"
)
```

## Types

### Config

Configuration for the LiveAvatar client.

```go
type Config struct {
    APIKey     string       // API key (or LIVEAVATAR_API_KEY env var)
    BaseURL    string       // API base URL (default: https://api.liveavatar.com)
    HTTPClient *http.Client // Custom HTTP client
}
```

### LiveKitConfig

LiveKit room configuration for LITE mode.

```go
type LiveKitConfig struct {
    LiveKitURL         string `json:"livekit_url"`
    LiveKitRoom        string `json:"livekit_room"`
    LiveKitClientToken string `json:"livekit_client_token"`
}
```

| Field | Description |
|-------|-------------|
| `LiveKitURL` | LiveKit server URL (wss://...) |
| `LiveKitRoom` | Room name for the avatar to join |
| `LiveKitClientToken` | JWT token with room_join, can_publish, can_subscribe, can_publish_data |

### NewSessionRequest

Parameters for creating a streaming session.

```go
type NewSessionRequest struct {
    Mode          string         `json:"mode"`
    AvatarID      string         `json:"avatar_id"`
    IsSandbox     bool           `json:"is_sandbox,omitempty"`
    VideoQuality  VideoQuality   `json:"video_quality,omitempty"`
    LiveKitConfig *LiveKitConfig `json:"livekit_config,omitempty"`
}
```

| Field | Description |
|-------|-------------|
| `Mode` | Session mode: `"LITE"` for BYO AI stack |
| `AvatarID` | UUID of the avatar to use |
| `IsSandbox` | Enable sandbox mode (60s limit, no credits) |
| `VideoQuality` | Avatar video quality |
| `LiveKitConfig` | LiveKit configuration (required for LITE mode) |

### NewSessionResponse

Response from creating a session.

```go
type NewSessionResponse struct {
    SessionID    string `json:"session_id"`
    SessionToken string `json:"session_token"`
}
```

### StartSessionResponse

Response from starting a session.

```go
type StartSessionResponse struct {
    SessionID          string `json:"session_id"`
    WSURL              string `json:"ws_url"`
    MaxSessionDuration int    `json:"max_session_duration"`
}
```

| Field | Description |
|-------|-------------|
| `SessionID` | Unique session identifier |
| `WSURL` | WebSocket URL for streaming audio/events |
| `MaxSessionDuration` | Maximum session duration in seconds |

## Functions

### NewClient

Creates a new LiveAvatar client.

```go
func NewClient(cfg *Config) (*Client, error)
```

**Example:**

```go
// Using environment variable
client, err := liveavatar.NewClient(nil)

// With explicit API key
client, err := liveavatar.NewClient(&liveavatar.Config{
    APIKey: "your-api-key",
})
```

## Methods

### NewSession

Creates a new streaming session and returns session credentials.

```go
func (c *Client) NewSession(ctx context.Context, req *NewSessionRequest) (*NewSessionResponse, error)
```

**Example:**

```go
resp, err := client.NewSession(ctx, &liveavatar.NewSessionRequest{
    Mode:      "LITE",
    AvatarID:  liveavatar.SandboxAvatarID,
    IsSandbox: true,
    LiveKitConfig: &liveavatar.LiveKitConfig{
        LiveKitURL:         "wss://project.livekit.cloud",
        LiveKitRoom:        "my-room",
        LiveKitClientToken: lkToken,
    },
})
```

### StartSession

Starts a streaming session and returns the WebSocket URL.

```go
func (c *Client) StartSession(ctx context.Context, sessionID, sessionToken string) (*StartSessionResponse, error)
```

**Example:**

```go
startResp, err := client.StartSession(ctx, resp.SessionID, resp.SessionToken)
fmt.Printf("WebSocket: %s\n", startResp.WSURL)
```

### StopSession

Stops a streaming session.

```go
func (c *Client) StopSession(ctx context.Context, sessionID, sessionToken string, reason StopReason) error
```

**Example:**

```go
err := client.StopSession(ctx, sessionID, sessionToken, liveavatar.StopReasonUserDisconnected)
```

## Error Codes

| Code | Description |
|------|-------------|
| 1000 | Success |
| 4000 | Validation error (check request format) |
| 4010 | Invalid credentials (check API key) |
