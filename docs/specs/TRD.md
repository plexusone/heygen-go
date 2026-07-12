# HeyGen Go SDK - Technical Requirements Document

## Architecture

### Package Structure

```
heygen-go/
├── heygen/                 # Core client and shared types
│   ├── client.go           # HTTP client with auth
│   ├── config.go           # Configuration options
│   └── errors.go           # Error types
├── ogen/                   # Generated OpenAPI client
│   └── ...                 # ogen-generated code
├── avatar/                 # Avatar management
│   └── avatar.go           # List, get avatars
├── voice/                  # Voice management
│   └── voice.go            # List, get voices
├── video/                  # Video generation
│   └── video.go            # Create, get videos
├── liveavatar/             # Real-time streaming
│   ├── session.go          # Session lifecycle
│   ├── websocket.go        # WebSocket client
│   ├── events.go           # Event types
│   └── audio.go            # Audio frame handling
└── docs/
    └── specs/
```

### Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        heygen-go                            │
├─────────────────────────────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌─────────────┐  │
│  │  avatar  │  │  voice   │  │  video   │  │ liveavatar  │  │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └──────┬──────┘  │
│       │             │             │               │         │
│       └─────────────┴─────────────┴───────┬───────┘         │
│                                           │                 │
│  ┌────────────────────────────────────────┴──────────────┐  │
│  │                     heygen (core)                     │  │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────────────────┐ │  │
│  │  │  client  │  │  config  │  │  ogen (generated)    │ │  │
│  │  └──────────┘  └──────────┘  └──────────────────────┘ │  │
│  └───────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    HeyGen REST + WebSocket APIs
```

## Technical Specifications

### Authentication

HeyGen uses API key authentication via header:

```
X-Api-Key: <api_key>
```

Configuration:

```go
type Config struct {
    APIKey     string
    BaseURL    string        // Default: https://api.heygen.com
    HTTPClient *http.Client  // Optional custom client
    Logger     *slog.Logger  // Optional logger
}
```

### REST API Client

Base URL: `https://api.heygen.com`

Generated via ogen from OpenAPI spec:

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/avatar.list` | GET | List available avatars |
| `/v1/voice.list` | GET | List available voices |
| `/v2/video/generate` | POST | Generate video |
| `/v1/video_status.get` | GET | Get video status |
| `/v1/streaming.new` | POST | Create LiveAvatar session |
| `/v1/streaming.start` | POST | Start streaming session |
| `/v1/streaming.stop` | POST | Stop streaming session |

### LiveAvatar WebSocket Protocol

Connection URL: `wss://liveavatar.heygen.com/v1/ws`

#### Message Types (Client → Server)

```go
// Audio frame
type AudioMessage struct {
    Type      string `json:"type"`      // "audio"
    SessionID string `json:"session_id"`
    Audio     string `json:"audio"`     // base64 PCM
    SampleRate int   `json:"sample_rate"`
}

// Interrupt current speech
type InterruptMessage struct {
    Type      string `json:"type"`      // "interrupt"
    SessionID string `json:"session_id"`
}

// End session
type CloseMessage struct {
    Type      string `json:"type"`      // "close"
    SessionID string `json:"session_id"`
}
```

#### Message Types (Server → Client)

```go
// Session ready
type ReadyEvent struct {
    Type      string `json:"type"`      // "ready"
    SessionID string `json:"session_id"`
}

// Avatar started speaking
type SpeakingStartEvent struct {
    Type      string `json:"type"`      // "speaking_start"
    SessionID string `json:"session_id"`
}

// Avatar finished speaking
type SpeakingEndEvent struct {
    Type      string `json:"type"`      // "speaking_end"
    SessionID string `json:"session_id"`
}

// Error occurred
type ErrorEvent struct {
    Type    string `json:"type"`        // "error"
    Code    string `json:"code"`
    Message string `json:"message"`
}
```

### Audio Format

LiveAvatar expects:

- Format: PCM (raw audio samples)
- Sample rate: 16000 Hz or 24000 Hz
- Channels: 1 (mono)
- Bit depth: 16-bit signed little-endian

```go
type AudioFrame struct {
    PCM        []byte
    SampleRate int  // 16000 or 24000
    Channels   int  // 1
}
```

### Session Lifecycle

```
┌─────────────────────────────────────────────────────────────┐
│                    LiveAvatar Session                       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  POST /streaming.new ──► SessionID + WebSocket URL          │
│          │                                                  │
│          ▼                                                  │
│  Connect WebSocket ──► Wait for "ready" event               │
│          │                                                  │
│          ▼                                                  │
│  POST /streaming.start ──► Avatar joins LiveKit room        │
│          │                                                  │
│          ▼                                                  │
│  Stream audio frames ──► Avatar lip-syncs and speaks        │
│          │                                                  │
│          ▼                                                  │
│  POST /streaming.stop ──► Session ends                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### Error Handling

All errors implement the standard error interface with additional context:

```go
type APIError struct {
    StatusCode int
    Code       string
    Message    string
    RequestID  string
}

func (e *APIError) Error() string {
    return fmt.Sprintf("heygen: %s (code=%s, status=%d)", e.Message, e.Code, e.StatusCode)
}
```

### Rate Limiting

HeyGen enforces rate limits. The SDK should:

1. Parse `X-RateLimit-*` headers
2. Implement exponential backoff on 429 responses
3. Expose rate limit info to callers

### Logging

Use `log/slog` for structured logging:

```go
logger.Info("session created",
    "session_id", sessionID,
    "avatar_id", avatarID,
)
```

## Integration with omni-livekit

The SDK enables an omni-livekit HeyGen provider:

```go
// In omni-livekit/avatar/heygen/provider.go

type Provider struct {
    client *heygen.Client
}

func (p *Provider) CreateSession(ctx context.Context, cfg avatar.Config) (avatar.Session, error) {
    // Use heygen-go liveavatar package
    session, err := liveavatar.NewSession(p.client, liveavatar.SessionConfig{
        AvatarID:   cfg.AvatarID,
        VoiceID:    cfg.VoiceID,
        LiveKitURL: cfg.LiveKitURL,
        RoomName:   cfg.RoomName,
    })
    if err != nil {
        return nil, err
    }
    return &heygenSession{session: session}, nil
}
```

## Testing Strategy

| Type | Coverage | Tools |
|------|----------|-------|
| Unit tests | Core logic, parsing | go test |
| Integration tests | API calls (with mocks) | httptest |
| E2E tests | Real HeyGen API | Manual/CI with API key |

## Security Considerations

1. **API Key Storage** - Never log or expose API keys
2. **TLS** - All connections use HTTPS/WSS
3. **Input Validation** - Validate all user inputs before API calls

## Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/ogen-go/ogen` | latest | OpenAPI generation |
| `nhooyr.io/websocket` | v1.8.x | WebSocket client |
| `golang.org/x/sync` | latest | Concurrency primitives |

## Compatibility

- Go 1.22+
- HeyGen API v1/v2
