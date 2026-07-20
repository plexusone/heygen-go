package avatar

import (
	"context"
	"strings"
)

// V2Avatar is an avatar returned by the v2 avatars API. Unlike the v3
// avatars API (which returns avatar groups), the ID here is directly
// usable as the avatar_id in v2 video generation.
type V2Avatar struct {
	// AvatarID is the identifier to use as avatar_id in video generation.
	AvatarID string `json:"avatar_id"`

	// AvatarName is the display name.
	AvatarName string `json:"avatar_name"`

	// Gender of the avatar.
	Gender string `json:"gender,omitempty"`

	// PreviewImageURL is the URL to a preview image.
	PreviewImageURL string `json:"preview_image_url,omitempty"`

	// PreviewVideoURL is the URL to a preview video.
	PreviewVideoURL string `json:"preview_video_url,omitempty"`
}

// TalkingPhoto is a talking-photo avatar returned by the v2 avatars API.
// The ID is usable as talking_photo_id in v2 video generation.
type TalkingPhoto struct {
	// TalkingPhotoID is the identifier to use as talking_photo_id.
	TalkingPhotoID string `json:"talking_photo_id"`

	// TalkingPhotoName is the display name.
	TalkingPhotoName string `json:"talking_photo_name"`

	// PreviewImageURL is the URL to a preview image.
	PreviewImageURL string `json:"preview_image_url,omitempty"`
}

// V2ListResponse is the response from the v2 avatars API.
type V2ListResponse struct {
	Data struct {
		Avatars       []V2Avatar     `json:"avatars"`
		TalkingPhotos []TalkingPhoto `json:"talking_photos"`
	} `json:"data"`
}

// ListV2 returns avatars and talking photos from the v2 avatars API,
// whose IDs are directly usable in v2 video generation (unlike the v3
// avatar-group IDs from List).
//
// Note: this endpoint returns the full catalog (often 1000+ avatars) in a
// single response and can be slow; configure the client with a generous
// timeout (see heygen.Config.HTTPClient).
func (c *Client) ListV2(ctx context.Context) (*V2ListResponse, error) {
	var resp V2ListResponse
	if err := c.client.Get(ctx, "/v2/avatars", &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SearchV2 returns v2 avatars whose ID or name contains the given term
// (case-insensitive). An empty term returns all avatars. Talking photos
// are not included.
func (c *Client) SearchV2(ctx context.Context, term string) ([]V2Avatar, error) {
	resp, err := c.ListV2(ctx)
	if err != nil {
		return nil, err
	}
	if term == "" {
		return resp.Data.Avatars, nil
	}

	term = strings.ToLower(term)
	var matches []V2Avatar
	for _, a := range resp.Data.Avatars {
		if strings.Contains(strings.ToLower(a.AvatarID), term) ||
			strings.Contains(strings.ToLower(a.AvatarName), term) {
			matches = append(matches, a)
		}
	}
	return matches, nil
}
