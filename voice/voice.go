// Package voice provides access to HeyGen voice APIs.
package voice

import (
	"context"

	"github.com/plexusone/heygen-go/heygen"
)

// Client provides access to voice APIs.
type Client struct {
	client *heygen.Client
}

// NewClient creates a new voice client.
func NewClient(client *heygen.Client) *Client {
	return &Client{client: client}
}

// Voice represents a HeyGen voice.
type Voice struct {
	// VoiceID is the unique identifier for the voice.
	VoiceID string `json:"voice_id"`

	// Name is the display name of the voice.
	Name string `json:"name"`

	// Language is the voice's language (e.g., "en-US").
	Language string `json:"language"`

	// Gender is the voice's gender (male, female).
	Gender string `json:"gender"`

	// PreviewAudio is the URL of the voice's preview audio.
	PreviewAudio string `json:"preview_audio"`

	// SupportPause indicates if the voice supports pause markers.
	SupportPause bool `json:"support_pause"`

	// EmotionSupport indicates if the voice supports emotion.
	EmotionSupport bool `json:"emotion_support"`

	// Type is the voice type (heygen, elevenlabs, etc.).
	Type string `json:"type"`
}

// ListResponse is the response from listing voices.
type ListResponse struct {
	// Data contains the voice data.
	Data struct {
		Voices []Voice `json:"voices"`
	} `json:"data"`
}

// List returns all available voices.
func (c *Client) List(ctx context.Context) ([]Voice, error) {
	var resp ListResponse
	if err := c.client.Get(ctx, "/v2/voices", &resp); err != nil {
		return nil, err
	}
	return resp.Data.Voices, nil
}

// ListV1Response is the response from the v1 voice list endpoint.
type ListV1Response struct {
	// Data contains the voice list.
	Data struct {
		Voices []Voice `json:"voices"`
	} `json:"data"`
}

// ListV1 returns all available voices using the v1 endpoint.
func (c *Client) ListV1(ctx context.Context) ([]Voice, error) {
	var resp ListV1Response
	if err := c.client.Get(ctx, "/v1/voice.list", &resp); err != nil {
		return nil, err
	}
	return resp.Data.Voices, nil
}
