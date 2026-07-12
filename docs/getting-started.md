# Getting Started

This guide walks you through setting up heygen-go and making your first API calls.

## Prerequisites

- Go 1.21 or later
- A HeyGen API key (for avatars/video)
- A LiveAvatar API key (for real-time streaming)

## Installation

```bash
go get github.com/plexusone/heygen-go
```

## Configuration

### Environment Variables

Set up your API keys as environment variables:

```bash
# For avatar listing and video generation
export HEYGEN_API_KEY="sk_..."

# For real-time streaming (LiveAvatar)
export LIVEAVATAR_API_KEY="..."

# For LiveKit integration (required for LiveAvatar LITE mode)
export LIVEKIT_URL="wss://your-project.livekit.cloud"
export LIVEKIT_API_KEY="..."
export LIVEKIT_API_SECRET="..."
```

!!! tip "Using direnv"
    Create a `.envrc` file in your project and use [direnv](https://direnv.net/) to auto-load environment variables.

### Creating a Client

```go
import heygen "github.com/plexusone/heygen-go"

// Option 1: Pass API key directly
client, err := heygen.New("your-api-key")

// Option 2: Read from HEYGEN_API_KEY environment variable
client, err := heygen.New("")
```

## Your First API Call

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
    client, err := heygen.New("")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // List available avatars
    resp, err := client.Avatar.List(ctx, &avatar.ListOptions{
        Limit:     10,
        Ownership: "public",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found %d avatars\n", len(resp.Data))
    for _, a := range resp.Data {
        fmt.Printf("  - %s (%s): %d looks\n", a.Name, a.ID, a.LooksCount)
    }
}
```

## Next Steps

- [Avatars Guide](guides/avatars.md) - Learn about avatar groups and looks
- [LiveAvatar Streaming](guides/liveavatar.md) - Set up real-time avatar sessions
