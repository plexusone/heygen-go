// Package video provides access to HeyGen video generation APIs.
package video

import (
	"context"
	"fmt"

	"github.com/plexusone/heygen-go/heygen"
)

// Client provides access to video APIs.
type Client struct {
	client *heygen.Client
}

// NewClient creates a new video client.
func NewClient(client *heygen.Client) *Client {
	return &Client{client: client}
}

// GenerateRequest is the request to generate a video.
type GenerateRequest struct {
	// Title is the video title.
	Title string `json:"title,omitempty"`

	// VideoInputs defines the video content.
	VideoInputs []VideoInput `json:"video_inputs"`

	// Dimension specifies the video dimensions.
	Dimension *Dimension `json:"dimension,omitempty"`

	// AspectRatio specifies the aspect ratio (alternative to Dimension).
	AspectRatio string `json:"aspect_ratio,omitempty"`

	// Test generates a test video (watermarked, no credits used).
	Test bool `json:"test,omitempty"`

	// CallbackID is an optional callback identifier.
	CallbackID string `json:"callback_id,omitempty"`
}

// VideoInput defines a single video clip.
type VideoInput struct {
	// Character defines the avatar to use.
	Character Character `json:"character"`

	// Voice defines the voice and text.
	Voice VoiceInput `json:"voice"`

	// Background defines the background (optional).
	Background *Background `json:"background,omitempty"`
}

// Character defines the avatar character.
type Character struct {
	// Type is the character type (avatar, talking_photo).
	Type string `json:"type"`

	// AvatarID is the avatar identifier.
	AvatarID string `json:"avatar_id,omitempty"`

	// AvatarStyle is the avatar style (normal, circle, etc.).
	AvatarStyle string `json:"avatar_style,omitempty"`

	// TalkingPhotoID is the talking photo identifier.
	TalkingPhotoID string `json:"talking_photo_id,omitempty"`
}

// VoiceInput defines the voice and text for a video clip.
type VoiceInput struct {
	// Type is the voice input type (text, audio).
	Type string `json:"type"`

	// VoiceID is the voice identifier.
	VoiceID string `json:"voice_id,omitempty"`

	// InputText is the text to speak.
	InputText string `json:"input_text,omitempty"`

	// AudioURL is the URL of audio to use (for audio type).
	AudioURL string `json:"audio_url,omitempty"`

	// Speed is the speech speed (0.5 - 1.5).
	Speed float64 `json:"speed,omitempty"`

	// Pitch is the voice pitch adjustment.
	Pitch int `json:"pitch,omitempty"`

	// Emotion is the voice emotion.
	Emotion string `json:"emotion,omitempty"`
}

// Dimension specifies video dimensions.
type Dimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Background defines the video background.
type Background struct {
	// Type is the background type (color, image, video).
	Type string `json:"type"`

	// Value is the background value (color hex, URL, etc.).
	Value string `json:"value,omitempty"`
}

// GenerateResponse is the response from generating a video.
type GenerateResponse struct {
	// Data contains the video generation data.
	Data struct {
		VideoID string `json:"video_id"`
	} `json:"data"`
}

// Generate creates a new video.
func (c *Client) Generate(ctx context.Context, req GenerateRequest) (string, error) {
	var resp GenerateResponse
	if err := c.client.Post(ctx, "/v2/video/generate", req, &resp); err != nil {
		return "", err
	}
	return resp.Data.VideoID, nil
}

// Status represents the video generation status.
type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

// Video represents a generated video.
type Video struct {
	// VideoID is the unique video identifier.
	VideoID string `json:"video_id"`

	// Status is the current generation status.
	Status Status `json:"status"`

	// VideoURL is the URL of the completed video.
	VideoURL string `json:"video_url,omitempty"`

	// ThumbnailURL is the URL of the video thumbnail.
	ThumbnailURL string `json:"thumbnail_url,omitempty"`

	// Duration is the video duration in seconds.
	Duration float64 `json:"duration,omitempty"`

	// GifURL is the URL of the video as a GIF.
	GifURL string `json:"gif_url,omitempty"`

	// Error contains error details if generation failed.
	Error *VideoError `json:"error,omitempty"`

	// CallbackID is the callback identifier if provided.
	CallbackID string `json:"callback_id,omitempty"`
}

// VideoError contains error details for failed videos.
type VideoError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// StatusResponse is the response from getting video status.
type StatusResponse struct {
	// Data contains the video data.
	Data Video `json:"data"`
}

// GetStatus returns the status of a video.
func (c *Client) GetStatus(ctx context.Context, videoID string) (*Video, error) {
	var resp StatusResponse
	path := fmt.Sprintf("/v1/video_status.get?video_id=%s", videoID)
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// ListResponse is the response from listing videos.
type ListResponse struct {
	// Data contains the video list.
	Data struct {
		Videos []Video `json:"videos"`
	} `json:"data"`
}

// List returns all videos.
func (c *Client) List(ctx context.Context) ([]Video, error) {
	var resp ListResponse
	if err := c.client.Get(ctx, "/v1/video.list", &resp); err != nil {
		return nil, err
	}
	return resp.Data.Videos, nil
}

// Delete deletes a video.
func (c *Client) Delete(ctx context.Context, videoID string) error {
	path := fmt.Sprintf("/v1/video.delete?video_id=%s", videoID)
	return c.client.Delete(ctx, path, nil)
}
