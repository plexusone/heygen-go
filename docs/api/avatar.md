# Avatar API Reference

Package `avatar` provides access to HeyGen avatar APIs (v3).

```go
import "github.com/plexusone/heygen-go/avatar"
```

## Types

### Avatar

Represents a HeyGen avatar group.

```go
type Avatar struct {
    ID              string `json:"id"`
    Name            string `json:"name"`
    Gender          string `json:"gender,omitempty"`
    CreatedAt       int64  `json:"created_at"`
    LooksCount      int    `json:"looks_count"`
    DefaultVoiceID  string `json:"default_voice_id,omitempty"`
    PreviewImageURL string `json:"preview_image_url,omitempty"`
    PreviewVideoURL string `json:"preview_video_url,omitempty"`
    Status          string `json:"status,omitempty"`
}
```

| Field | Description |
|-------|-------------|
| `ID` | Unique avatar group identifier (hex string) |
| `Name` | Display name |
| `Gender` | Avatar gender (male, female, Man, Woman) |
| `CreatedAt` | Unix timestamp of creation |
| `LooksCount` | Number of available looks |
| `DefaultVoiceID` | Default voice for this avatar |
| `PreviewImageURL` | URL to preview image |
| `PreviewVideoURL` | URL to preview video |
| `Status` | Training status (private avatars only) |

### Look

Represents an avatar look (outfit/style).

```go
type Look struct {
    ID                  string   `json:"id"`
    Name                string   `json:"name"`
    AvatarType          string   `json:"avatar_type"`
    GroupID             string   `json:"group_id,omitempty"`
    Gender              string   `json:"gender,omitempty"`
    DefaultVoiceID      string   `json:"default_voice_id,omitempty"`
    PreviewImageURL     string   `json:"preview_image_url,omitempty"`
    PreviewVideoURL     string   `json:"preview_video_url,omitempty"`
    SupportedAPIEngines []string `json:"supported_api_engines,omitempty"`
    Status              string   `json:"status,omitempty"`
}
```

| Field | Description |
|-------|-------------|
| `ID` | Unique look identifier (use this for video generation) |
| `Name` | Display name |
| `AvatarType` | Engine type: `studio_avatar`, `digital_twin`, `photo_avatar` |
| `GroupID` | Parent avatar group ID |
| `SupportedAPIEngines` | Compatible engines: `avatar_v`, `avatar_iv`, `avatar_iii` |
| `Status` | Training status: `processing`, `completed`, `failed` |

### ListOptions

Options for listing avatars.

```go
type ListOptions struct {
    Limit     int    // Max avatars to return
    Token     string // Pagination token
    Ownership string // Filter: "public", "private", "all"
}
```

### ListResponse

Response from listing avatars.

```go
type ListResponse struct {
    Data      []Avatar `json:"data"`
    HasMore   bool     `json:"has_more,omitempty"`
    NextToken string   `json:"next_token,omitempty"`
}
```

### LooksResponse

Response from listing looks.

```go
type LooksResponse struct {
    Data      []Look `json:"data"`
    HasMore   bool   `json:"has_more,omitempty"`
    NextToken string `json:"next_token,omitempty"`
}
```

## Methods

### List

Lists available avatars.

```go
func (c *Client) List(ctx context.Context, opts *ListOptions) (*ListResponse, error)
```

**Parameters:**

- `ctx` - Context for cancellation
- `opts` - List options (nil for defaults)

**Example:**

```go
resp, err := client.Avatar.List(ctx, &avatar.ListOptions{
    Limit:     20,
    Ownership: "public",
})
```

### Get

Gets details for a specific avatar group.

```go
func (c *Client) Get(ctx context.Context, groupID string) (*Avatar, error)
```

**Parameters:**

- `ctx` - Context for cancellation
- `groupID` - Avatar group ID

**Example:**

```go
avatar, err := client.Avatar.Get(ctx, "e0e84faea390465896db75a83be45085")
```

### ListLooks

Lists looks for an avatar group.

```go
func (c *Client) ListLooks(ctx context.Context, groupID string, limit int) (*LooksResponse, error)
```

**Parameters:**

- `ctx` - Context for cancellation
- `groupID` - Avatar group ID
- `limit` - Max looks to return (0 for default)

**Example:**

```go
looks, err := client.Avatar.ListLooks(ctx, "e0e84faea390465896db75a83be45085", 50)
```

## v2 Avatars (generation-ready IDs)

The `List`/`Get`/`ListLooks` methods above use the **v3** avatars API, which
returns avatar *groups*. The **v2** avatars API returns individual avatars
whose IDs are directly usable as `avatar_id` in video generation.

!!! warning "Use v2 IDs for generation"
    A v3 avatar-group ID (e.g. `9fc2a78e642547b5a21ea8cf06a953d4`) is
    **rejected** by the video-generation endpoint (`avatar look not found`).
    Use a v2 avatar ID (e.g. `Abigail_expressive_2024112501`) from `ListV2`
    or `SearchV2`.

### V2Avatar / TalkingPhoto

```go
type V2Avatar struct {
    AvatarID        string // use as avatar_id in generation
    AvatarName      string
    Gender          string
    PreviewImageURL string
    PreviewVideoURL string
}

type TalkingPhoto struct {
    TalkingPhotoID  string // use as talking_photo_id in generation
    TalkingPhotoName string
    PreviewImageURL string
}
```

### ListV2

```go
func (c *Client) ListV2(ctx context.Context) (*V2ListResponse, error)
```

Returns avatars and talking photos from the v2 avatars API. This endpoint
returns the full catalog (often 1000+ avatars) in one response and can be
slow, so configure the client with a generous timeout.

### SearchV2

```go
func (c *Client) SearchV2(ctx context.Context, term string) ([]V2Avatar, error)
```

Filters the v2 catalog by ID or name (case-insensitive substring; an empty
term returns all).

**Example:**

```go
avatars, err := client.Avatar.SearchV2(ctx, "abigail")
for _, a := range avatars {
    fmt.Printf("%s  %s (%s)\n", a.AvatarID, a.AvatarName, a.Gender)
}
```
