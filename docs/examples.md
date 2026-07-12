# Examples

This page provides complete, runnable examples for common use cases.

## List Avatars

List available HeyGen avatars using the v3 API.

```go
// examples/list-avatars/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    heygen "github.com/plexusone/heygen-go"
    "github.com/plexusone/heygen-go/avatar"
)

func main() {
    apiKey := os.Getenv("HEYGEN_API_KEY")
    if apiKey == "" {
        log.Fatal("HEYGEN_API_KEY environment variable is required")
    }

    client, err := heygen.New(apiKey)
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    fmt.Println("=== Avatars (v3 API) ===")
    resp, err := client.Avatar.List(ctx, &avatar.ListOptions{
        Limit: 10,
    })
    if err != nil {
        log.Fatal(err)
    }

    for _, a := range resp.Data {
        fmt.Printf("  %s (%s) - %s, %d looks\n", a.Name, a.ID, a.Gender, a.LooksCount)
    }
    fmt.Printf("Total: %d avatars (has_more: %v)\n", len(resp.Data), resp.HasMore)
}
```

**Run:**

```bash
export HEYGEN_API_KEY="sk_..."
go run ./examples/list-avatars
```

## LiveAvatar Session

Create a real-time avatar streaming session with LiveKit.

```go
// examples/liveavatar-session/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "time"

    "github.com/livekit/protocol/auth"
    "github.com/plexusone/heygen-go/liveavatar"
)

func main() {
    ctx := context.Background()

    // Validate required environment variables
    liveKitURL := os.Getenv("LIVEKIT_URL")
    liveKitAPIKey := os.Getenv("LIVEKIT_API_KEY")
    liveKitAPISecret := os.Getenv("LIVEKIT_API_SECRET")

    if liveKitURL == "" || liveKitAPIKey == "" || liveKitAPISecret == "" {
        log.Fatal("LIVEKIT_URL, LIVEKIT_API_KEY, and LIVEKIT_API_SECRET are required")
    }

    // Create LiveAvatar client (reads LIVEAVATAR_API_KEY from env)
    client, err := liveavatar.NewClient(nil)
    if err != nil {
        log.Fatal(err)
    }

    // Generate LiveKit token for the avatar agent
    roomName := "liveavatar-demo"
    identity := "liveavatar-avatar-agent"

    at := auth.NewAccessToken(liveKitAPIKey, liveKitAPISecret)
    at.SetVideoGrant(&auth.VideoGrant{
        RoomJoin:       true,
        Room:           roomName,
        CanPublish:     boolPtr(true),
        CanSubscribe:   boolPtr(true),
        CanPublishData: boolPtr(true),
    }).
        SetIdentity(identity).
        SetValidFor(time.Hour)

    lkToken, err := at.ToJWT()
    if err != nil {
        log.Fatal("Failed to generate LiveKit token:", err)
    }

    // Create streaming session (using sandbox mode for testing)
    fmt.Println("Creating LiveAvatar session...")
    sessionResp, err := client.NewSession(ctx, &liveavatar.NewSessionRequest{
        Mode:      "LITE",
        AvatarID:  liveavatar.SandboxAvatarID,
        IsSandbox: true,
        LiveKitConfig: &liveavatar.LiveKitConfig{
            LiveKitURL:         liveKitURL,
            LiveKitRoom:        roomName,
            LiveKitClientToken: lkToken,
        },
    })
    if err != nil {
        log.Fatal("Failed to create session:", err)
    }

    fmt.Printf("Session created: %s\n", sessionResp.SessionID)

    // Start the session
    fmt.Println("Starting session...")
    startResp, err := client.StartSession(ctx, sessionResp.SessionID, sessionResp.SessionToken)
    if err != nil {
        log.Fatal("Failed to start session:", err)
    }

    fmt.Printf("Session started!\n")
    fmt.Printf("  Session ID: %s\n", startResp.SessionID)
    fmt.Printf("  WebSocket URL: %s\n", startResp.WSURL)
    fmt.Printf("  Max Duration: %d seconds\n", startResp.MaxSessionDuration)

    // In a real application, connect to WebSocket and stream audio
    time.Sleep(5 * time.Second)

    // Stop the session
    fmt.Println("Stopping session...")
    if err := client.StopSession(ctx, sessionResp.SessionID, sessionResp.SessionToken,
        liveavatar.StopReasonUserDisconnected); err != nil {
        log.Fatal("Failed to stop session:", err)
    }

    fmt.Println("Session stopped successfully!")
}

func boolPtr(b bool) *bool { return &b }
```

**Run:**

```bash
export LIVEAVATAR_API_KEY="..."
export LIVEKIT_URL="wss://project.livekit.cloud"
export LIVEKIT_API_KEY="..."
export LIVEKIT_API_SECRET="..."
go run ./examples/liveavatar-session
```

## Running Examples

All examples are in the `examples/` directory:

```bash
# Clone the repo
git clone https://github.com/plexusone/heygen-go
cd heygen-go

# Set up environment
cp .envrc.example .envrc  # Edit with your keys
source .envrc

# Run an example
go run ./examples/list-avatars
go run ./examples/liveavatar-session
```
