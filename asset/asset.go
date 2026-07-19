// Package asset provides access to the HeyGen asset upload API.
//
// Assets (audio, images, video) are uploaded to the dedicated upload
// service (upload.heygen.com, a different host from the main API) and
// return a hosted URL usable in other HeyGen APIs — for example as the
// audio source for video generation:
//
//	uploaded, err := assetClient.Upload(ctx, asset.ContentTypeMPEG, audioFile)
//	// use uploaded.URL as video.VoiceInput{Type: "audio", AudioURL: ...}
package asset

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/plexusone/heygen-go/heygen"
)

// DefaultBaseURL is the HeyGen asset upload service base URL.
// This is a different host from the main API (api.heygen.com).
const DefaultBaseURL = "https://upload.heygen.com"

// Supported asset content types per the HeyGen documentation.
const (
	// ContentTypeJPEG is a JPEG image asset.
	ContentTypeJPEG = "image/jpeg"

	// ContentTypePNG is a PNG image asset.
	ContentTypePNG = "image/png"

	// ContentTypeMP4 is an MP4 video asset.
	ContentTypeMP4 = "video/mp4"

	// ContentTypeWebM is a WebM video asset.
	ContentTypeWebM = "video/webm"

	// ContentTypeMPEG is an MPEG audio asset (e.g., MP3).
	ContentTypeMPEG = "audio/mpeg"
)

// Client provides access to the asset upload API.
type Client struct {
	client  *heygen.Client
	baseURL string
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL sets a custom upload service base URL.
// Default: DefaultBaseURL.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// NewClient creates a new asset client.
func NewClient(client *heygen.Client, opts ...Option) *Client {
	c := &Client{
		client:  client,
		baseURL: DefaultBaseURL,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// Asset represents an uploaded asset.
type Asset struct {
	// ID is the unique asset identifier.
	ID string `json:"id"`

	// Name is the asset name assigned by HeyGen.
	Name string `json:"name"`

	// FileType is the asset kind (e.g., "audio", "image", "video").
	FileType string `json:"file_type"`

	// URL is the hosted asset URL, usable in other HeyGen APIs.
	URL string `json:"url"`
}

// uploadResponse is the v1 API envelope for asset uploads.
type uploadResponse struct {
	Code int    `json:"code"`
	Data Asset  `json:"data"`
	Msg  string `json:"msg"`
}

// Upload uploads raw asset content with the given content type and
// returns the hosted asset. contentType must be one of the supported
// ContentType* values (the API rejects unsupported types).
func (c *Client) Upload(ctx context.Context, contentType string, r io.Reader) (*Asset, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("read asset content: %w", err)
	}

	var resp uploadResponse
	if err := c.client.RequestURL(ctx, http.MethodPost, c.baseURL+"/v1/asset", contentType, data, &resp); err != nil {
		return nil, err
	}
	if resp.Data.URL == "" {
		return nil, fmt.Errorf("heygen asset upload: response missing asset URL (code %d, msg %q)", resp.Code, resp.Msg)
	}
	return &resp.Data, nil
}
