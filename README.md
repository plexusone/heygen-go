# HeyGen Go SDK

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/plexusone/heygen-go/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/plexusone/heygen-go/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/plexusone/heygen-go/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/plexusone/heygen-go/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/plexusone/heygen-go/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/plexusone/heygen-go/actions/workflows/go-sast-codeql.yaml
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/plexusone/heygen-go
 [docs-godoc-url]: https://pkg.go.dev/github.com/plexusone/heygen-go
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://plexusone.dev/heygen-go
 [viz-svg]: https://img.shields.io/badge/Go-visualizaton-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=plexusone%2Fheygen-go
 [loc-svg]: https://tokei.rs/b1/github/plexusone/heygen-go
 [repo-url]: https://github.com/plexusone/heygen-go
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/plexusone/heygen-go/blob/main/LICENSE

Go SDK for the [HeyGen API](https://docs.heygen.com/).

## Features

- 🎭 **Avatar Management** - List avatars (v3 groups and generation-ready v2 IDs) and a built-in public-avatar catalog
- 🎙️ **Voice Management** - List available voices
- 🎬 **Video Generation** - Create AI-generated videos
- 📤 **Asset Upload** - Upload audio, images, and video for use in other APIs
- 📡 **LiveAvatar Streaming** - Real-time avatar streaming sessions
- 🤖 **OmniAvatar Adapter** - Use HeyGen render behind the provider-agnostic [OmniAvatar](https://github.com/plexusone/omniavatar) interfaces
- 🔄 **Retry Logic** - Automatic retries with exponential backoff
- ⚠️ **Error Handling** - Typed errors with request IDs

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

### Upload Assets

Upload audio, images, or video to HeyGen's asset service (upload.heygen.com)
and use the hosted URL in other APIs — for example, driving avatar lip-sync
from your own narration audio:

```go
import "github.com/plexusone/heygen-go/asset"

f, err := os.Open("narration.mp3")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

uploaded, err := client.Asset.Upload(ctx, asset.ContentTypeMPEG, f)
if err != nil {
    log.Fatal(err)
}

// Use the hosted URL as the audio source for video generation
videoID, err := client.Video.Generate(ctx, video.GenerateRequest{
    VideoInputs: []video.VideoInput{
        {
            Character: video.Character{Type: "avatar", AvatarID: "avatar_id"},
            Voice:     video.VoiceInput{Type: "audio", AudioURL: uploaded.URL},
        },
    },
})
```

Supported content types: `asset.ContentTypeJPEG`, `ContentTypePNG`,
`ContentTypeMP4`, `ContentTypeWebM`, `ContentTypeMPEG` (MP3 audio).

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
| Avatars | `v3/avatars`, `v3/avatars/{id}`, `v3/avatars/{id}/looks`, `v2/avatars` | Implemented |
| Voices | `v2/voices` | Planned |
| Videos | `v2/video/generate` | Planned |
| Assets | `v1/asset` (upload.heygen.com) | Implemented |
| LiveAvatar | `v1/sessions/token`, `v1/sessions/start`, `v1/sessions/stop` | Implemented |

## OmniAvatar Integration

The [`omniavatar`](omniavatar/) subpackage implements the provider-agnostic
[OmniAvatar](https://github.com/plexusone/omniavatar) **render** interfaces
(`render.Provider`, `render.AudioUploader`, `render.AvatarLister`) on top of
this SDK, so HeyGen video generation can be used behind the OmniAvatar
abstraction. It depends only on
[`omniavatar-core`](https://github.com/plexusone/omniavatar-core) (interfaces
only — no LiveKit).

```go
import heygenomni "github.com/plexusone/heygen-go/omniavatar"

p, err := heygenomni.NewRenderProvider(heygenomni.RenderConfig{
    APIKey: os.Getenv("HEYGEN_API_KEY"),
})
```

Provider adapters live in the provider SDK repos so all HeyGen-specific
knowledge stays here. The real-time **live** (LiveAvatar) adapter, which
requires LiveKit, lives in the batteries-included
[`omniavatar`](https://github.com/plexusone/omniavatar) package instead.

## Related Projects

- [omni-livekit](https://github.com/plexusone/omni-livekit) - LiveKit voice agents with avatar support
- [tavus-go](https://github.com/plexusone/tavus-go) - Go SDK for Tavus API
- [bithuman-go](https://github.com/plexusone/bithuman-go) - Go SDK for bitHuman API

## License

MIT License
