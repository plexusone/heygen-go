# OmniAvatar Adapter

The `omniavatar` subpackage lets HeyGen video generation be used behind the
provider-agnostic [OmniAvatar](https://github.com/plexusone/omniavatar)
**render** interfaces. It implements the `omniavatar-core` interfaces on top
of this SDK and depends only on
[`omniavatar-core`](https://github.com/plexusone/omniavatar-core) (interfaces
+ helpers — no LiveKit).

```go
import heygenomni "github.com/plexusone/heygen-go/omniavatar"
```

## Implemented interfaces

| Interface | Capability |
|-----------|------------|
| `render.Provider` | Generate / Status / Download (audio-driven talking-head MP4) |
| `render.AudioUploader` | Upload local narration via the asset API → hosted URL |
| `render.AvatarLister` | List generation-ready avatars via the v2 avatars API |

## Construct directly

```go
p, err := heygenomni.NewRenderProvider(heygenomni.RenderConfig{
    APIKey: os.Getenv("HEYGEN_API_KEY"),
})

// Upload local narration (AudioUploader)
url, err := p.UploadAudio(ctx, "narration.mp3", f)

// Generate, wait, download
job, err := p.Generate(ctx, render.GenerateRequest{AvatarID: avatarID, AudioURL: url})
status, err := render.Wait(ctx, p, job.ID, 5*time.Second)
err = p.Download(ctx, job.ID, out)

// Discover generation-ready avatar IDs (AvatarLister)
avatars, err := p.ListAvatars(ctx, "abigail")
```

## Via the batteries package

The [`omniavatar`](https://github.com/plexusone/omniavatar) package registers
this adapter by name, so you can construct it from configuration:

```go
import (
    "github.com/plexusone/omniavatar"
    _ "github.com/plexusone/omniavatar/providers/all"
)

p, err := omniavatar.GetRenderProvider("heygen",
    omniavatar.WithAPIKey(os.Getenv("HEYGEN_API_KEY")))
```

## Config

`RenderConfig` uses the **HeyGen API key** (`HEYGEN_API_KEY`), distinct from
the LiveAvatar key. Extensions supported via the batteries registry include
`avatar_id`, `talking_photo_id`, `avatar_style`, `voice_id`, `test`, and
`upload_base_url`.

!!! note "Live vs. render"
    This adapter covers the **render** (batch video generation) surface. The
    real-time **live** (LiveAvatar) adapter lives in the batteries-included
    `omniavatar` package, because its LiveKit integration lives there.
