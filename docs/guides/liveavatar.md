# LiveAvatar Streaming Guide

LiveAvatar provides real-time avatar streaming for live video conversations. This guide covers LITE mode integration with LiveKit.

## Overview

LiveAvatar is HeyGen's real-time streaming product. In **LITE mode**, you bring your own AI stack (STT/LLM/TTS) and LiveAvatar handles avatar rendering.

```
┌─────────────────┐     ┌──────────────────┐     ┌────────────────┐
│   Your Agent    │────▶│    LiveAvatar    │────▶│   LiveKit      │
│  (STT/LLM/TTS)  │     │  (Avatar Render) │     │   (WebRTC)     │
└─────────────────┘     └──────────────────┘     └────────────────┘
         │                       │                       │
         └───────────────────────┴───────────────────────┘
                          Audio + Video
```

## Prerequisites

1. **LiveAvatar API Key** - From [app.liveavatar.com/developers](https://app.liveavatar.com/developers)
2. **LiveKit Credentials** - URL, API key, and secret from [livekit.io](https://livekit.io)
3. **Avatar ID** - UUID of the avatar to use (see [Sandbox Mode](#sandbox-mode))

## Session Flow

### 1. Create a Session Token

```go
import (
    "github.com/livekit/protocol/auth"
    "github.com/plexusone/heygen-go/liveavatar"
)

// Create LiveAvatar client
client, err := liveavatar.NewClient(nil) // reads LIVEAVATAR_API_KEY

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

// Create session
sessionResp, err := client.NewSession(ctx, &liveavatar.NewSessionRequest{
    Mode:      "LITE",
    AvatarID:  "65f9e3c9-d48b-4118-b73a-4ae2e3cbb8f0",
    IsSandbox: true,
    LiveKitConfig: &liveavatar.LiveKitConfig{
        LiveKitURL:         liveKitURL,
        LiveKitRoom:        "my-room",
        LiveKitClientToken: lkToken,
    },
})
```

### 2. Start the Session

```go
startResp, err := client.StartSession(ctx, sessionResp.SessionID, sessionResp.SessionToken)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("WebSocket URL: %s\n", startResp.WSURL)
fmt.Printf("Max Duration: %d seconds\n", startResp.MaxSessionDuration)
```

### 3. Stream Audio via WebSocket

Connect to the WebSocket URL and send audio events:

```go
// Connect to WebSocket
conn, _, err := websocket.DefaultDialer.Dial(startResp.WSURL, nil)

// Send audio data
msg := map[string]interface{}{
    "type":     "agent.speak",
    "event_id": uuid.New().String(),
    "audio":    base64.StdEncoding.EncodeToString(audioData),
}
conn.WriteJSON(msg)

// Signal end of speech
conn.WriteJSON(map[string]interface{}{
    "type":     "agent.speak_end",
    "event_id": uuid.New().String(),
})
```

### 4. Stop the Session

```go
err := client.StopSession(ctx, sessionResp.SessionID, sessionResp.SessionToken,
    liveavatar.StopReasonUserDisconnected)
```

## Sandbox Mode

Use sandbox mode for development and testing:

- **No credit usage** - Free to use
- **60-second limit** - Sessions auto-terminate after ~1 minute
- **Limited avatars** - Only specific avatars are available

```go
sessionResp, err := client.NewSession(ctx, &liveavatar.NewSessionRequest{
    Mode:      "LITE",
    AvatarID:  liveavatar.SandboxAvatarID, // Pre-defined test avatar
    IsSandbox: true,
    // ...
})
```

## WebSocket Events

### Outgoing Events (Agent → LiveAvatar)

| Event | Description |
|-------|-------------|
| `agent.speak` | Send audio data (base64 PCM 24kHz mono) |
| `agent.speak_end` | Signal end of current speech |
| `agent.interrupt` | Interrupt current avatar speech |
| `agent.start_listening` | Indicate agent is listening |
| `agent.stop_listening` | Indicate agent stopped listening |
| `session.keep_alive` | Keep session alive (send every 60s) |

### Incoming Events (LiveAvatar → Agent)

| Event | Description |
|-------|-------------|
| `session.state_updated` | Session state changed (connected, etc.) |
| `agent.speak_started` | Avatar started speaking |
| `agent.speak_ended` | Avatar finished speaking |
| `agent.speak_interrupted` | Avatar speech was interrupted |

## Audio Format

- **Sample Rate**: 24000 Hz
- **Channels**: Mono (1)
- **Format**: PCM signed 16-bit little-endian
- **Encoding**: Base64 for WebSocket transport

## LiveKit Token Requirements

The LiveKit token must include these permissions:

| Permission | Required | Purpose |
|------------|----------|---------|
| `room_join` | Yes | Join the room |
| `can_publish` | Yes | Publish avatar video |
| `can_subscribe` | Yes | Subscribe to participant audio |
| `can_publish_data` | Yes | Send data messages |

## Error Handling

```go
sessionResp, err := client.NewSession(ctx, req)
if err != nil {
    // Check for specific error codes
    if strings.Contains(err.Error(), "4000") {
        log.Fatal("Validation error - check request format")
    }
    if strings.Contains(err.Error(), "4010") {
        log.Fatal("Invalid API key")
    }
    log.Fatal(err)
}
```

## Best Practices

1. **Reuse sessions** - Don't create a new session for each interaction
2. **Handle reconnection** - WebSocket may disconnect; implement reconnection logic
3. **Buffer audio** - Send 600ms-1s chunks for optimal lip sync
4. **Use sandbox first** - Develop with sandbox mode, then switch to production
5. **Clean up sessions** - Always call StopSession when done
