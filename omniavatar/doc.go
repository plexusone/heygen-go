// Package omniavatar provides OmniAvatar provider implementations backed by
// the HeyGen API.
//
// It implements the omniavatar-core interfaces for both surfaces:
//
//   - live.Provider / live.Session — real-time LiveAvatar sessions (LITE
//     mode, LiveKit audio streaming), using the LiveAvatar API key
//   - render.Provider — asynchronous video generation (audio-driven
//     talking-head MP4), using the HeyGen API key; also implements
//     render.AudioUploader (asset API) and render.AvatarLister (v2 avatars)
//
// The adapter is constructor-based and depends only on omniavatar-core
// (not the batteries omniavatar package). The batteries package registers
// these constructors; import it to use them by name:
//
//	import (
//	    "github.com/plexusone/omniavatar"
//	    _ "github.com/plexusone/omniavatar/providers/all"
//	)
//
//	live, err := omniavatar.GetLiveProvider("heygen",
//	    omniavatar.WithAPIKey(os.Getenv("LIVEAVATAR_API_KEY")),
//	    omniavatar.WithExtension("avatar_id", avatarID))
//
//	render, err := omniavatar.GetRenderProvider("heygen",
//	    omniavatar.WithAPIKey(os.Getenv("HEYGEN_API_KEY")))
//
// Or construct directly:
//
//	p, err := heygenomni.NewRenderProvider(heygenomni.RenderConfig{
//	    APIKey: os.Getenv("HEYGEN_API_KEY"),
//	})
package omniavatar
