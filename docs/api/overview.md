# API Reference

This section provides detailed API documentation for all heygen-go packages.

## Packages

| Package | Description |
|---------|-------------|
| [`heygen`](https://pkg.go.dev/github.com/plexusone/heygen-go) | Main client and configuration |
| [`avatar`](avatar.md) | Avatar listing and details (v3 groups + v2 generation-ready IDs) |
| [`asset`](asset.md) | Asset upload (audio, images, video) |
| [`liveavatar`](liveavatar.md) | Real-time streaming sessions + public avatar catalog |
| [`omniavatar`](omniavatar.md) | OmniAvatar render adapter |

## Client Architecture

```go
import heygen "github.com/plexusone/heygen-go"

client, err := heygen.New("api-key")
// client.Avatar  → avatar.Client
// client.Voice   → voice.Client (planned)
// client.Video   → video.Client
// client.Asset   → asset.Client
```

## Configuration Options

```go
client, err := heygen.New("api-key",
    heygen.WithBaseURL("https://custom-api.example.com"),
    heygen.WithRetry(3),
    heygen.WithHTTPClient(customHTTPClient),
)
```

## Error Handling

All API errors are wrapped in `heygen.APIError`:

```go
import "errors"

var apiErr *heygen.APIError
if errors.As(err, &apiErr) {
    fmt.Printf("Code: %s\n", apiErr.Code)
    fmt.Printf("Message: %s\n", apiErr.Message)
    fmt.Printf("Status: %d\n", apiErr.StatusCode)
    fmt.Printf("RequestID: %s\n", apiErr.RequestID)
}

// Helper functions
if heygen.IsUnauthorized(err) { /* 401 */ }
if heygen.IsRateLimited(err) { /* 429 */ }
if heygen.IsNotFound(err) { /* 404 */ }
```

## Retry Behavior

The client automatically retries on:

- HTTP 429 (Rate Limited)
- HTTP 500+ (Server Errors)
- Network timeouts

Default retry configuration:

- Max retries: 2
- Base delay: 1 second
- Max delay: 30 seconds
- Backoff: Exponential with jitter
