// Example: Create a LiveAvatar streaming session with LiveKit.
//
// This example demonstrates the LITE mode flow where your agent handles
// STT/LLM/TTS and LiveAvatar provides the avatar rendering.
//
// Required environment variables:
//
//	export LIVEAVATAR_API_KEY="your-liveavatar-api-key"  # From app.liveavatar.com/developers
//	export LIVEKIT_URL="wss://your-project.livekit.cloud"
//	export LIVEKIT_API_KEY="your-livekit-api-key"
//	export LIVEKIT_API_SECRET="your-livekit-api-secret"
//
// Usage:
//
//	go run ./examples/liveavatar-session
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

	// In a real application, you would:
	// 1. Connect to the WebSocket URL
	// 2. Send agent.speak events with audio data
	// 3. Handle session.state_updated events
	//
	// See the Python livekit-plugins-liveavatar for the WebSocket protocol.

	fmt.Println("\nSession will auto-stop in sandbox mode after ~60 seconds.")
	fmt.Println("Press Ctrl+C to stop early, or wait for auto-stop.")

	// Wait a bit then stop (in real usage, you'd wait for user disconnect)
	time.Sleep(5 * time.Second)

	// Stop the session
	fmt.Println("Stopping session...")
	if err := client.StopSession(ctx, sessionResp.SessionID, sessionResp.SessionToken, liveavatar.StopReasonUserDisconnected); err != nil {
		log.Fatal("Failed to stop session:", err)
	}

	fmt.Println("Session stopped successfully!")
}

func boolPtr(b bool) *bool { return &b }
