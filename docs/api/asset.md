# Asset API

The `asset` package uploads audio, images, and video to HeyGen's asset
service and returns a hosted URL usable in other HeyGen APIs — most
notably as the audio source for audio-driven video generation.

!!! note "Separate Upload Host"
    Assets are uploaded to `upload.heygen.com`, a different host from the
    main API (`api.heygen.com`). The same `HEYGEN_API_KEY` is used.

## Upload an Asset

```go
import (
    "context"
    "log"
    "os"

    heygen "github.com/plexusone/heygen-go"
    "github.com/plexusone/heygen-go/asset"
)

func main() {
    client, err := heygen.New(os.Getenv("HEYGEN_API_KEY"))
    if err != nil {
        log.Fatal(err)
    }

    f, err := os.Open("narration.mp3")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    uploaded, err := client.Asset.Upload(context.Background(), asset.ContentTypeMPEG, f)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("asset %s hosted at %s", uploaded.ID, uploaded.URL)
}
```

## Audio-Driven Video Generation

The primary use case: drive avatar lip-sync from your own narration audio
instead of HeyGen TTS.

```go
import "github.com/plexusone/heygen-go/video"

videoID, err := client.Video.Generate(ctx, video.GenerateRequest{
    VideoInputs: []video.VideoInput{
        {
            Character: video.Character{Type: "avatar", AvatarID: avatarID},
            Voice:     video.VoiceInput{Type: "audio", AudioURL: uploaded.URL},
        },
    },
})
```

## Supported Content Types

| Constant | MIME Type | Asset Kind |
|----------|-----------|------------|
| `asset.ContentTypeJPEG` | `image/jpeg` | Image |
| `asset.ContentTypePNG` | `image/png` | Image |
| `asset.ContentTypeMP4` | `video/mp4` | Video |
| `asset.ContentTypeWebM` | `video/webm` | Video |
| `asset.ContentTypeMPEG` | `audio/mpeg` | Audio (MP3) |

!!! warning "Audio Format"
    MP3 (`audio/mpeg`) is the documented audio asset type. Other audio
    formats (WAV, OGG) may be rejected by the API — transcode to MP3
    first if needed.

## Types

```go
// Asset represents an uploaded asset.
type Asset struct {
    ID       string // unique asset identifier
    Name     string // asset name assigned by HeyGen
    FileType string // "audio", "image", or "video"
    URL      string // hosted URL usable in other HeyGen APIs
}
```

## Options

```go
// Custom upload endpoint (testing, proxies)
assetClient := asset.NewClient(coreClient, asset.WithBaseURL("https://upload.example.com"))
```

## Error Handling

Upload failures return `*heygen.APIError` like all other endpoints:

```go
_, err := client.Asset.Upload(ctx, asset.ContentTypeMPEG, f)
if heygen.IsUnauthorized(err) {
    // invalid API key
}
```
