// Package avatar provides access to HeyGen avatar APIs.
package avatar

import (
	"context"
	"fmt"

	"github.com/plexusone/heygen-go/heygen"
)

// Client provides access to avatar APIs.
type Client struct {
	client *heygen.Client
}

// NewClient creates a new avatar client.
func NewClient(client *heygen.Client) *Client {
	return &Client{client: client}
}

// Avatar represents a HeyGen avatar group.
type Avatar struct {
	// ID is the unique avatar group identifier.
	ID string `json:"id"`

	// Name is the display name of the avatar.
	Name string `json:"name"`

	// Gender is the avatar's gender (male, female, Man, Woman).
	Gender string `json:"gender,omitempty"`

	// CreatedAt is the Unix timestamp when the avatar was created.
	CreatedAt int64 `json:"created_at"`

	// LooksCount is the number of looks (outfits/styles) available.
	LooksCount int `json:"looks_count"`

	// DefaultVoiceID is the default voice ID for this avatar.
	DefaultVoiceID string `json:"default_voice_id,omitempty"`

	// PreviewImageURL is the URL of the avatar's preview image.
	PreviewImageURL string `json:"preview_image_url,omitempty"`

	// PreviewVideoURL is the URL of the avatar's preview video.
	PreviewVideoURL string `json:"preview_video_url,omitempty"`

	// Status is the training status (processing, pending_consent, failed, completed).
	// Only present for private avatars.
	Status string `json:"status,omitempty"`
}

// ListResponse is the response from listing avatars (v3 API).
type ListResponse struct {
	// Data contains the avatar list.
	Data []Avatar `json:"data"`

	// HasMore indicates if more pages are available.
	HasMore bool `json:"has_more,omitempty"`

	// NextToken is the opaque cursor for the next page.
	NextToken string `json:"next_token,omitempty"`
}

// ListOptions configures the List request.
type ListOptions struct {
	// Limit is the maximum number of avatars to return.
	Limit int

	// Token is the pagination token from a previous request.
	Token string

	// Ownership filters by ownership (public, private, all).
	Ownership string
}

// List returns available avatars using the v3 API.
func (c *Client) List(ctx context.Context, opts *ListOptions) (*ListResponse, error) {
	path := "/v3/avatars"
	if opts != nil {
		sep := "?"
		if opts.Limit > 0 {
			path += fmt.Sprintf("%slimit=%d", sep, opts.Limit)
			sep = "&"
		}
		if opts.Token != "" {
			path += fmt.Sprintf("%stoken=%s", sep, opts.Token)
			sep = "&"
		}
		if opts.Ownership != "" {
			path += fmt.Sprintf("%sownership=%s", sep, opts.Ownership)
		}
	}

	var resp ListResponse
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Get returns details for a specific avatar group.
func (c *Client) Get(ctx context.Context, groupID string) (*Avatar, error) {
	path := fmt.Sprintf("/v3/avatars/%s", groupID)
	var resp struct {
		Data Avatar `json:"data"`
	}
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// Look represents an avatar look (outfit/style).
type Look struct {
	// ID is the unique look identifier. Use this as avatar_id in video generation.
	ID string `json:"id"`

	// Name is the display name of the look.
	Name string `json:"name"`

	// AvatarType determines engine and parameter compatibility.
	// Values: studio_avatar, digital_twin, photo_avatar
	AvatarType string `json:"avatar_type"`

	// GroupID is the ID of the avatar group this look belongs to.
	GroupID string `json:"group_id,omitempty"`

	// Gender of the avatar.
	Gender string `json:"gender,omitempty"`

	// DefaultVoiceID is the default voice for this look.
	DefaultVoiceID string `json:"default_voice_id,omitempty"`

	// PreviewImageURL is the URL to the look preview image.
	PreviewImageURL string `json:"preview_image_url,omitempty"`

	// PreviewVideoURL is the URL to the look preview video.
	PreviewVideoURL string `json:"preview_video_url,omitempty"`

	// SupportedAPIEngines lists engines this look supports.
	// Values: avatar_v, avatar_iv, avatar_iii
	SupportedAPIEngines []string `json:"supported_api_engines,omitempty"`

	// Status is the training status (processing, completed, failed).
	Status string `json:"status,omitempty"`
}

// LooksResponse is the response from listing avatar looks.
type LooksResponse struct {
	Data      []Look `json:"data"`
	HasMore   bool   `json:"has_more,omitempty"`
	NextToken string `json:"next_token,omitempty"`
}

// ListLooks returns the looks for an avatar group.
func (c *Client) ListLooks(ctx context.Context, groupID string, limit int) (*LooksResponse, error) {
	path := fmt.Sprintf("/v3/avatars/%s/looks", groupID)
	if limit > 0 {
		path += fmt.Sprintf("?limit=%d", limit)
	}

	var resp LooksResponse
	if err := c.client.Get(ctx, path, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
