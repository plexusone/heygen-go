package omniavatar

import (
	"github.com/plexusone/omniavatar-core/registry"
	"github.com/plexusone/omniavatar-core/render"
)

// NewRenderProviderFromConfig creates a HeyGen render provider from
// registry config. It is the registry-friendly constructor the batteries
// omniavatar package registers as a render.Provider factory.
func NewRenderProviderFromConfig(cfg registry.ProviderConfig) (render.Provider, error) {
	return NewRenderProvider(RenderConfig{
		APIKey:        cfg.APIKey,
		BaseURL:       cfg.BaseURL,
		AvatarID:      cfg.GetString("avatar_id", ""),
		UploadBaseURL: cfg.GetString("upload_base_url", ""),
	})
}
