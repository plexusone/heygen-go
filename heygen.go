// Package heygen provides a Go SDK for the HeyGen API.
//
// The SDK provides access to HeyGen's video generation and real-time
// streaming (LiveAvatar) capabilities.
//
// # Quick Start
//
// Create a client with your API key:
//
//	client, err := heygen.New("your-api-key")
//	if err != nil {
//		log.Fatal(err)
//	}
//
// List available avatars:
//
//	avatars, err := client.Avatar.List(ctx)
//
// Generate a video:
//
//	videoID, err := client.Video.Generate(ctx, video.GenerateRequest{
//		Title: "My Video",
//		VideoInputs: []video.VideoInput{
//			{
//				Character: video.Character{
//					Type:     "avatar",
//					AvatarID: "avatar_id",
//				},
//				Voice: video.VoiceInput{
//					Type:      "text",
//					VoiceID:   "voice_id",
//					InputText: "Hello, world!",
//				},
//			},
//		},
//	})
//
// Create a streaming session:
//
//	session, err := client.Streaming.NewSession(ctx, streaming.NewSessionRequest{
//		Quality:  streaming.QualityHigh,
//		AvatarID: "avatar_id",
//	})
//
// Upload an asset (audio, image, video) for use in other APIs:
//
//	uploaded, err := client.Asset.Upload(ctx, asset.ContentTypeMPEG, audioFile)
//	// uploaded.URL is usable as video.VoiceInput{Type: "audio", AudioURL: ...}
package heygen

import (
	"github.com/plexusone/heygen-go/asset"
	"github.com/plexusone/heygen-go/avatar"
	"github.com/plexusone/heygen-go/heygen"
	"github.com/plexusone/heygen-go/streaming"
	"github.com/plexusone/heygen-go/video"
	"github.com/plexusone/heygen-go/voice"
)

// Client is the HeyGen SDK client providing access to all APIs.
type Client struct {
	// Avatar provides access to avatar APIs.
	Avatar *avatar.Client

	// Voice provides access to voice APIs.
	Voice *voice.Client

	// Video provides access to video generation APIs.
	Video *video.Client

	// Streaming provides access to streaming (LiveAvatar) APIs.
	Streaming *streaming.Client

	// Asset provides access to the asset upload API.
	Asset *asset.Client

	// core is the underlying HTTP client.
	core *heygen.Client
}

// New creates a new HeyGen SDK client with the given API key.
func New(apiKey string, opts ...Option) (*Client, error) {
	cfg := heygen.Config{
		APIKey: apiKey,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	core, err := heygen.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		Avatar:    avatar.NewClient(core),
		Voice:     voice.NewClient(core),
		Video:     video.NewClient(core),
		Streaming: streaming.NewClient(core),
		Asset:     asset.NewClient(core),
		core:      core,
	}, nil
}

// Option configures the client.
type Option func(*heygen.Config)

// WithBaseURL sets a custom base URL.
func WithBaseURL(url string) Option {
	return func(cfg *heygen.Config) {
		cfg.BaseURL = url
	}
}

// WithRetry configures retry behavior.
func WithRetry(maxRetries int) Option {
	return func(cfg *heygen.Config) {
		cfg.Retry.MaxRetries = maxRetries
	}
}

// Re-export common types for convenience
type (
	// Config is the client configuration.
	Config = heygen.Config

	// APIError represents an API error.
	APIError = heygen.APIError

	// RetryConfig configures retry behavior.
	RetryConfig = heygen.RetryConfig
)

// Re-export error checking functions
var (
	IsNotFound     = heygen.IsNotFound
	IsRateLimited  = heygen.IsRateLimited
	IsUnauthorized = heygen.IsUnauthorized
)

// Re-export errors
var (
	ErrNoAPIKey    = heygen.ErrNoAPIKey
	ErrRateLimited = heygen.ErrRateLimited
)

// Version returns the SDK version.
func Version() string {
	return heygen.Version
}
