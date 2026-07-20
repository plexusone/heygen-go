package omniavatar

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/plexusone/omniavatar-core/render"

	heygenasset "github.com/plexusone/heygen-go/asset"
	heygenavatar "github.com/plexusone/heygen-go/avatar"
	heygensdk "github.com/plexusone/heygen-go/heygen"
	heygenvideo "github.com/plexusone/heygen-go/video"
)

// RenderConfig configures the HeyGen render (batch video generation) provider.
//
// Note: batch video generation uses the HeyGen API (api.heygen.com) and
// the HEYGEN_API_KEY, which is distinct from the LiveAvatar API key used
// by the live provider.
type RenderConfig struct {
	// APIKey is the HeyGen API key.
	// Required.
	APIKey string

	// BaseURL is the HeyGen API base URL.
	// Default: https://api.heygen.com
	BaseURL string

	// AvatarID is the default HeyGen avatar used when
	// GenerateRequest.AvatarID is empty.
	AvatarID string

	// UploadBaseURL is the asset upload service base URL.
	// Default: https://upload.heygen.com
	UploadBaseURL string

	// HTTPClient is an optional custom HTTP client, used for both API
	// calls and video downloads.
	HTTPClient *http.Client
}

// RenderProvider implements render.Provider for HeyGen video generation.
// It also implements render.AudioUploader via the HeyGen asset upload API.
type RenderProvider struct {
	videos     *heygenvideo.Client
	assets     *heygenasset.Client
	avatarID   string
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Compile-time interface checks.
var (
	_ render.Provider      = (*RenderProvider)(nil)
	_ render.AudioUploader = (*RenderProvider)(nil)
)

// NewRenderProvider creates a HeyGen render provider.
func NewRenderProvider(cfg RenderConfig) (*RenderProvider, error) {
	if cfg.APIKey == "" {
		return nil, fmt.Errorf("%w: APIKey is required", render.ErrInvalidConfig)
	}

	client, err := heygensdk.NewClient(heygensdk.Config{
		APIKey:     cfg.APIKey,
		BaseURL:    cfg.BaseURL,
		HTTPClient: cfg.HTTPClient,
	})
	if err != nil {
		return nil, render.NewProviderError("heygen", "new_render_provider", err)
	}

	var assetOpts []heygenasset.Option
	if cfg.UploadBaseURL != "" {
		assetOpts = append(assetOpts, heygenasset.WithBaseURL(cfg.UploadBaseURL))
	}

	return &RenderProvider{
		videos:     heygenvideo.NewClient(client),
		assets:     heygenasset.NewClient(client, assetOpts...),
		avatarID:   cfg.AvatarID,
		apiKey:     cfg.APIKey,
		baseURL:    cfg.BaseURL,
		httpClient: cfg.HTTPClient,
	}, nil
}

// Compile-time check that the render provider can list avatars.
var _ render.AvatarLister = (*RenderProvider)(nil)

// ListAvatars implements render.AvatarLister using the HeyGen v2
// avatars API, whose IDs are directly usable as GenerateRequest.AvatarID
// (unlike the v3 avatar-group IDs returned elsewhere).
//
// The v2 avatars endpoint returns the full catalog in one response and can
// be slow, so this uses a dedicated client with a generous timeout unless
// the provider was configured with a custom HTTP client.
func (p *RenderProvider) ListAvatars(ctx context.Context, search string) ([]render.AvatarInfo, error) {
	httpClient := p.httpClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 3 * time.Minute}
	}

	client, err := heygensdk.NewClient(heygensdk.Config{
		APIKey:     p.apiKey,
		BaseURL:    p.baseURL,
		HTTPClient: httpClient,
	})
	if err != nil {
		return nil, render.NewProviderError("heygen", "list_avatars", err)
	}

	avatars, err := heygenavatar.NewClient(client).SearchV2(ctx, search)
	if err != nil {
		return nil, render.NewProviderError("heygen", "list_avatars", err)
	}

	infos := make([]render.AvatarInfo, len(avatars))
	for i, a := range avatars {
		infos[i] = render.AvatarInfo{ID: a.AvatarID, Name: a.AvatarName, Gender: a.Gender}
	}
	return infos, nil
}

// Name returns the provider name.
func (p *RenderProvider) Name() string { return "heygen" }

// Generate submits a video generation job to HeyGen.
//
// GenerateRequest.AvatarID maps to the HeyGen avatar ID; the
// "talking_photo_id" extension selects a talking photo instead.
// Background maps natively (color, image, video). Extensions:
// "avatar_style" (normal, circle, closeUp), "voice_id" (for Script
// input), "test" (bool, watermarked video without credits).
func (p *RenderProvider) Generate(ctx context.Context, req render.GenerateRequest) (*render.Job, error) {
	talkingPhotoID := req.GetString("talking_photo_id", "")
	if req.AvatarID == "" && talkingPhotoID == "" {
		req.AvatarID = p.avatarID
	}
	if talkingPhotoID != "" && req.AvatarID == "" {
		// Satisfy shared validation; character selection below uses the
		// talking photo, not the avatar ID.
		req.AvatarID = talkingPhotoID
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	character := heygenvideo.Character{
		Type:        "avatar",
		AvatarID:    req.AvatarID,
		AvatarStyle: req.GetString("avatar_style", ""),
	}
	if talkingPhotoID != "" {
		character = heygenvideo.Character{
			Type:           "talking_photo",
			TalkingPhotoID: talkingPhotoID,
		}
	}

	voice := heygenvideo.VoiceInput{Type: "audio", AudioURL: req.AudioURL}
	if req.Script != "" {
		voice = heygenvideo.VoiceInput{
			Type:      "text",
			InputText: req.Script,
			VoiceID:   req.GetString("voice_id", ""),
		}
	}

	input := heygenvideo.VideoInput{
		Character: character,
		Voice:     voice,
	}
	if req.Background != nil {
		input.Background = &heygenvideo.Background{
			Type:  req.Background.Type,
			Value: req.Background.Value,
		}
	}

	apiReq := heygenvideo.GenerateRequest{
		Title:       req.Title,
		VideoInputs: []heygenvideo.VideoInput{input},
		Test:        req.GetBool("test", false),
	}
	if req.Width > 0 && req.Height > 0 {
		apiReq.Dimension = &heygenvideo.Dimension{Width: req.Width, Height: req.Height}
	}

	videoID, err := p.videos.Generate(ctx, apiReq)
	if err != nil {
		return nil, render.NewProviderError("heygen", "generate", err)
	}

	return &render.Job{ID: videoID, Provider: "heygen"}, nil
}

// Status returns the current status of a generation job.
func (p *RenderProvider) Status(ctx context.Context, jobID string) (*render.JobStatus, error) {
	video, err := p.videos.GetStatus(ctx, jobID)
	if err != nil {
		return nil, render.NewProviderError("heygen", "status", err)
	}
	return videoToStatus(video), nil
}

// Download streams the completed video to dst. The status is re-fetched
// immediately before downloading because HeyGen video URLs are
// time-limited signed URLs.
func (p *RenderProvider) Download(ctx context.Context, jobID string, dst io.Writer) error {
	status, err := p.Status(ctx, jobID)
	if err != nil {
		return err
	}
	if status.State != render.JobStateCompleted || status.VideoURL == "" {
		return fmt.Errorf("%w: job %s is %s", render.ErrJobNotCompleted, jobID, status.State)
	}

	if err := render.DownloadURL(ctx, p.httpClient, status.VideoURL, dst); err != nil {
		return render.NewProviderError("heygen", "download", err)
	}
	return nil
}

// UploadAudio uploads audio content via the HeyGen asset upload API and
// returns a hosted URL usable as GenerateRequest.AudioURL.
//
// Note: HeyGen documents audio/mpeg (MP3) as the supported audio asset
// type; other formats may be rejected by the API.
func (p *RenderProvider) UploadAudio(ctx context.Context, filename string, r io.Reader) (string, error) {
	uploaded, err := p.assets.Upload(ctx, render.AudioContentType(filename), r)
	if err != nil {
		return "", render.NewProviderError("heygen", "upload_audio", err)
	}
	return uploaded.URL, nil
}

// videoToStatus converts a HeyGen Video to a normalized JobStatus.
func videoToStatus(video *heygenvideo.Video) *render.JobStatus {
	status := &render.JobStatus{
		ID:           video.VideoID,
		State:        mapVideoStatus(video.Status),
		RawStatus:    string(video.Status),
		VideoURL:     video.VideoURL,
		ThumbnailURL: video.ThumbnailURL,
		Duration:     video.Duration,
	}
	if video.Error != nil {
		status.ErrorCode = video.Error.Code
		status.ErrorMsg = video.Error.Message
	}
	return status
}

// mapVideoStatus maps HeyGen video statuses to normalized states.
func mapVideoStatus(s heygenvideo.Status) render.JobState {
	switch s {
	case heygenvideo.StatusPending:
		return render.JobStatePending
	case heygenvideo.StatusProcessing:
		return render.JobStateProcessing
	case heygenvideo.StatusCompleted:
		return render.JobStateCompleted
	case heygenvideo.StatusFailed:
		return render.JobStateFailed
	default:
		// Unknown states stay non-terminal so pollers keep waiting.
		return render.JobStateProcessing
	}
}
