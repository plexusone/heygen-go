# HeyGen Go SDK

[![Go Reference][docs-godoc-svg]][docs-godoc-url]
[![Go Report Card][goreport-svg]][goreport-url]
[![License][license-svg]][license-url]

 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/heygen-go
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/heygen-go
 [goreport-svg]: https://goreportcard.com/badge/github.com/plexusone/heygen-go
 [goreport-url]: https://goreportcard.com/report/github.com/plexusone/heygen-go
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/heygen-go/blob/main/LICENSE

Go SDK for the [HeyGen API](https://docs.heygen.com/).

## Features

- **Avatar Management** - List and retrieve avatar details
- **Voice Management** - List available voices
- **Video Generation** - Create AI-generated videos
- **LiveAvatar Streaming** - Real-time avatar streaming sessions
- **Retry Logic** - Automatic retries with exponential backoff
- **Error Handling** - Typed errors with request IDs

## Installation

```bash
go get github.com/plexusone/heygen-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    heygen "github.com/plexusone/heygen-go"
    "github.com/plexusone/heygen-go/avatar"
)

func main() {
    // Create client (reads HEYGEN_API_KEY from environment if not provided)
    client, err := heygen.New("your-api-key")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // List avatars (v3 API)
    resp, err := client.Avatar.List(ctx, &avatar.ListOptions{Limit: 10})
    if err != nil {
        log.Fatal(err)
    }
    for _, a := range resp.Data {
        fmt.Printf("Avatar: %s (%s) - %d looks\n", a.Name, a.ID, a.LooksCount)
    }
}
```

## Usage

### List Avatars

```go
import "github.com/plexusone/heygen-go/avatar"

// List avatars (v3 API with pagination)
resp, err := client.Avatar.List(ctx, &avatar.ListOptions{
    Limit:     20,
    Ownership: "public", // or "private", "all"
})
if err != nil {
    log.Fatal(err)
}

for _, a := range resp.Data {
    fmt.Printf("%s (%s) - %d looks\n", a.Name, a.ID, a.LooksCount)
}

// Get avatar looks (outfits/styles)
looks, err := client.Avatar.ListLooks(ctx, "avatar-group-id", 10)
```

### List Voices

```go
voices, err := client.Voice.List(ctx)
if err != nil {
    log.Fatal(err)
}
```

### Generate Video

```go
import "github.com/plexusone/heygen-go/video"

videoID, err := client.Video.Generate(ctx, video.GenerateRequest{
    Title: "My Video",
    Test:  true, // Use test mode (no credits)
    Dimension: &video.Dimension{
        Width:  1280,
        Height: 720,
    },
    VideoInputs: []video.VideoInput{
        {
            Character: video.Character{
                Type:        "avatar",
                AvatarID:    "avatar_id",
                AvatarStyle: "normal",
            },
            Voice: video.VoiceInput{
                Type:      "text",
                VoiceID:   "voice_id",
                InputText: "Hello from HeyGen!",
            },
        },
    },
})
if err != nil {
    log.Fatal(err)
}

// Poll for completion
for {
    status, err := client.Video.GetStatus(ctx, videoID)
    if err != nil {
        log.Fatal(err)
    }

    if status.Status == video.StatusCompleted {
        fmt.Printf("Video ready: %s\n", status.VideoURL)
        break
    }
    if status.Status == video.StatusFailed {
        log.Fatal("Video generation failed")
    }

    time.Sleep(5 * time.Second)
}
```

### LiveAvatar Streaming (LITE Mode)

LiveAvatar provides real-time avatar streaming. Use LITE mode for BYO AI stack
where your agent handles STT/LLM/TTS and LiveAvatar renders the avatar.

```go
import (
    "github.com/livekit/protocol/auth"
    "github.com/plexusone/heygen-go/liveavatar"
)

// Create LiveAvatar client (reads LIVEAVATAR_API_KEY from env)
client, err := liveavatar.NewClient(nil)
if err != nil {
    log.Fatal(err)
}

// Generate LiveKit token for the avatar agent
at := auth.NewAccessToken(liveKitAPIKey, liveKitAPISecret)
at.SetVideoGrant(&auth.VideoGrant{
    RoomJoin:       true,
    Room:           "my-room",
    CanPublish:     boolPtr(true),
    CanSubscribe:   boolPtr(true),
    CanPublishData: boolPtr(true),
}).SetIdentity("liveavatar-agent").SetValidFor(time.Hour)
lkToken, _ := at.ToJWT()

// Create session (use SandboxAvatarID for testing)
sessionResp, err := client.NewSession(ctx, &liveavatar.NewSessionRequest{
    Mode:      "LITE",
    AvatarID:  liveavatar.SandboxAvatarID, // or your avatar UUID
    IsSandbox: true, // 60s limit, no credits
    LiveKitConfig: &liveavatar.LiveKitConfig{
        LiveKitURL:         liveKitURL,
        LiveKitRoom:        "my-room",
        LiveKitClientToken: lkToken,
    },
})

// Start the session
startResp, err := client.StartSession(ctx, sessionResp.SessionID, sessionResp.SessionToken)
fmt.Printf("WebSocket URL: %s\n", startResp.WSURL)

// Connect to WebSocket and send agent.speak events with audio data
// See examples/liveavatar-session for full example

// Stop when done
client.StopSession(ctx, sessionResp.SessionID, sessionResp.SessionToken, liveavatar.StopReasonUserDisconnected)
```

## Configuration

### API Keys

**Important:** HeyGen and LiveAvatar use **separate API keys**:

| API | Environment Variable | Get Key From |
|-----|---------------------|--------------|
| HeyGen (video generation) | `HEYGEN_API_KEY` | [app.heygen.com/settings?nav=API](https://app.heygen.com/settings?nav=API) |
| LiveAvatar (real-time streaming) | `LIVEAVATAR_API_KEY` | [app.liveavatar.com/developers](https://app.liveavatar.com/developers) |

Set up a `.envrc` file for local development:

```bash
export HEYGEN_API_KEY="sk_..."        # For video generation
export LIVEAVATAR_API_KEY="la_..."    # For real-time streaming
```

Then load it with `source .envrc` or use [direnv](https://direnv.net/).

### Options

```go
client, err := heygen.New("api-key",
    heygen.WithBaseURL("https://custom-api.example.com"),
    heygen.WithRetry(3),
)
```

## Error Handling

```go
avatars, err := client.Avatar.List(ctx)
if err != nil {
    if heygen.IsUnauthorized(err) {
        log.Fatal("Invalid API key")
    }
    if heygen.IsRateLimited(err) {
        log.Fatal("Rate limited, try again later")
    }
    if heygen.IsNotFound(err) {
        log.Fatal("Resource not found")
    }

    // Get detailed error info
    var apiErr *heygen.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("Error: %s (code=%s, status=%d)\n",
            apiErr.Message, apiErr.Code, apiErr.StatusCode)
    }
}
```

## API Coverage

| API | Endpoints | Status |
|-----|-----------|--------|
| Avatars | `v3/avatars`, `v3/avatars/{id}`, `v3/avatars/{id}/looks` | Implemented |
| Voices | `v2/voices` | Planned |
| Videos | `v2/video/generate` | Planned |
| LiveAvatar | `v1/sessions/token`, `v1/sessions/start`, `v1/sessions/stop` | Implemented |

## Related Projects

- [omni-livekit](https://github.com/plexusone/omni-livekit) - LiveKit voice agents with avatar support
- [tavus-go](https://github.com/plexusone/tavus-go) - Go SDK for Tavus API
- [bithuman-go](https://github.com/plexusone/bithuman-go) - Go SDK for bitHuman API

## License

MIT License
