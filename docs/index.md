# heygen-go

Go SDK for the [HeyGen API](https://docs.heygen.com/) and [LiveAvatar API](https://docs.liveavatar.com/).

## Features

- **Avatar Management** - List and retrieve avatar details (v3 API)
- **Video Generation** - Create AI-generated avatar videos (v2 API)
- **Asset Upload** - Upload audio, images, and video for use in other APIs
- **LiveAvatar Streaming** - Real-time avatar sessions with LiveKit integration
- **Retry Logic** - Automatic retries with exponential backoff
- **Error Handling** - Typed errors with request IDs

## Installation

```bash
go get github.com/plexusone/heygen-go
```

## Quick Example

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
    client, err := heygen.New("your-api-key")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()
    resp, err := client.Avatar.List(ctx, &avatar.ListOptions{Limit: 10})
    if err != nil {
        log.Fatal(err)
    }

    for _, a := range resp.Data {
        fmt.Printf("Avatar: %s (%s) - %d looks\n", a.Name, a.ID, a.LooksCount)
    }
}
```

## API Keys

!!! warning "Separate API Keys Required"
    HeyGen and LiveAvatar use **different API keys** from different dashboards.

| API | Environment Variable | Dashboard |
|-----|---------------------|-----------|
| HeyGen (avatars, video) | `HEYGEN_API_KEY` | [app.heygen.com/settings?nav=API](https://app.heygen.com/settings?nav=API) |
| LiveAvatar (streaming) | `LIVEAVATAR_API_KEY` | [app.liveavatar.com/developers](https://app.liveavatar.com/developers) |

## Next Steps

- [Getting Started](getting-started.md) - Set up your environment
- [Avatars Guide](guides/avatars.md) - Work with HeyGen avatars
- [LiveAvatar Streaming](guides/liveavatar.md) - Real-time avatar sessions
