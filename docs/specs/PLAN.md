# HeyGen Go SDK - Implementation Plan

## Phase 1: Project Setup & OpenAPI Generation

### 1.1 Project Initialization

- [x] Create go.mod with module path `github.com/plexusone/heygen-go`
- [x] Add .gitignore
- [x] Fetch OpenAPI spec from community source (reviewed, used manual implementation)
- [x] Validate and clean up OpenAPI spec if needed

### 1.2 OpenAPI Code Generation

- [ ] Install ogen: `go install github.com/ogen-go/ogen/cmd/ogen@latest`
- [ ] Configure ogen generation (ogen.yml)
- [ ] Generate client code to `ogen/` package
- [ ] Verify generated code compiles

> **Note**: Skipped ogen generation in favor of manual typed client. The community OpenAPI spec lacks response schemas, making generated code less useful. Manual implementation provides better type safety.

### 1.3 Core Client Wrapper

- [x] Create `heygen/client.go` with configuration
- [x] Implement API key authentication
- [ ] Add request/response logging hooks (deferred to v0.2.0)
- [x] Create `heygen/errors.go` with error types

## Phase 2: REST API Packages

### 2.1 Avatar Package

```go
// avatar/avatar.go
type Client struct { ... }
func (c *Client) List(ctx context.Context) ([]Avatar, error)
func (c *Client) ListStreaming(ctx context.Context) ([]StreamingAvatar, error)
```

- [x] Implement avatar listing
- [x] Implement streaming avatar listing
- [ ] Add unit tests

### 2.2 Voice Package

```go
// voice/voice.go
type Client struct { ... }
func (c *Client) List(ctx context.Context) ([]Voice, error)
func (c *Client) ListV1(ctx context.Context) ([]Voice, error)
```

- [x] Implement voice listing
- [ ] Implement voice details (not needed for v0.1.0)
- [ ] Add unit tests

### 2.3 Video Package

```go
// video/video.go
type Client struct { ... }
func (c *Client) Generate(ctx context.Context, req GenerateRequest) (string, error)
func (c *Client) GetStatus(ctx context.Context, id string) (*Video, error)
func (c *Client) List(ctx context.Context) ([]Video, error)
func (c *Client) Delete(ctx context.Context, id string) error
```

- [x] Implement video generation
- [x] Implement status polling
- [x] Implement video listing
- [x] Implement video deletion
- [ ] Add unit tests

## Phase 3: LiveAvatar Real-Time Streaming

### 3.1 Session Management

```go
// streaming/session.go
func (c *Client) NewSession(ctx, req) (*Session, error)
func (c *Client) Start(ctx, req) (*SDP, error)
func (c *Client) Stop(ctx, sessionID) error
func (c *Client) Interrupt(ctx, sessionID) error
func (c *Client) SendTask(ctx, req) (string, error)
```

- [x] Implement session creation via REST API
- [x] Implement session start/stop lifecycle
- [x] Handle session tokens and URLs
- [x] Implement ICE candidate exchange
- [x] Implement task/speak functionality
- [x] Implement interrupt functionality

### 3.2 WebSocket Client

- [ ] Implement WebSocket connection with nhooyr.io/websocket
- [ ] Implement message serialization (JSON)
- [ ] Handle reconnection logic
- [ ] Implement keepalive/ping-pong

> **Deferred to v0.2.0**: WebSocket-based audio streaming requires deeper protocol analysis of the Python starter.

### 3.3 Event Handling

- [ ] Define event types
- [ ] Implement event dispatching
- [ ] Add callback hooks

> **Deferred to v0.2.0**: Requires WebSocket implementation.

### 3.4 Audio Streaming

- [ ] Implement PCM frame encoding
- [ ] Handle sample rate conversion if needed
- [ ] Implement buffering for smooth streaming

> **Deferred to v0.2.0**: Requires WebSocket implementation.

## Phase 4: Integration & Documentation

### 4.1 High-Level Client

```go
// heygen.go (root package)
type Client struct {
    Avatar    *avatar.Client
    Voice     *voice.Client
    Video     *video.Client
    Streaming *streaming.Client
}

func New(apiKey string, opts ...Option) *Client
```

- [x] Create unified client entry point
- [x] Add functional options pattern
- [x] Export public API

### 4.2 Examples

- [x] `examples/list-avatars/main.go`
- [ ] `examples/generate-video/main.go`
- [ ] `examples/streaming-session/main.go`

### 4.3 Documentation

- [x] README.md with quick start
- [x] GoDoc comments on all public types
- [x] docs/specs/ specification documents
- [ ] docs/guides/ usage guides

## Phase 5: omni-livekit Integration

### 5.1 Provider Implementation

Location: `omni-livekit/avatar/heygen/`

```go
type Provider struct {
    client *heygen.Client
}

func (p *Provider) Name() string { return "heygen" }
func (p *Provider) CreateSession(ctx, cfg) (Session, error)
```

- [ ] Implement avatar.Provider interface
- [ ] Register in provider factory
- [ ] Add configuration via environment variables

### 5.2 Session Adapter

```go
type heygenSession struct {
    session *liveavatar.Session
}

func (s *heygenSession) Start(ctx context.Context) error
func (s *heygenSession) PushAudio(frame avatar.AudioFrame) error
func (s *heygenSession) Interrupt() error
func (s *heygenSession) Close() error
```

- [ ] Adapt heygen-go session to omni-livekit interface
- [ ] Handle LiveKit room joining
- [ ] Implement audio routing

## Verification Checklist

### Build

```bash
go build ./...      # ✓ Passing
go test ./...       # Tests not yet written
golangci-lint run   # ✓ 0 issues
```

### Integration Test

```bash
export HEYGEN_API_KEY="..."
go run ./examples/list-avatars
go run ./examples/streaming-session
```

### omni-livekit Test

```bash
export AVATAR_PROVIDER=heygen
export HEYGEN_API_KEY="..."
go run -tags opus ./cmd/livekit-agent
```

## File Creation Order

1. [x] `go.mod`, `.gitignore`
2. [x] `heygen/client.go`, `heygen/config.go`, `heygen/errors.go`
3. [x] `avatar/avatar.go`
4. [x] `voice/voice.go`
5. [x] `video/video.go`
6. [x] `streaming/session.go`
7. [ ] `streaming/websocket.go` (v0.2.0)
8. [ ] `streaming/events.go` (v0.2.0)
9. [ ] `streaming/audio.go` (v0.2.0)
10. [x] `heygen.go` (root)
11. [x] `examples/list-avatars/main.go`
12. [x] `README.md`
13. [x] `LICENSE`

## Risk Mitigation

| Risk | Mitigation | Status |
|------|------------|--------|
| OpenAPI spec incomplete | Manual typed client implementation | Resolved |
| WebSocket protocol undocumented | Study Python starter implementation | Pending |
| LiveKit integration complexity | Reuse patterns from Tavus/bitHuman providers | Pending |
| Rate limiting | Implemented exponential backoff | Resolved |
