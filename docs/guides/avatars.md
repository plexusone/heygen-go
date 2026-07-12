# Avatars Guide

HeyGen provides a library of AI avatars for video generation and real-time streaming. This guide covers how to work with avatars using the v3 API.

## Concepts

### Avatar Groups

An avatar group represents a single avatar character with multiple "looks" (outfits, poses, or styles). Each avatar has:

- **ID** - Unique identifier (hex string like `e0e84faea390465896db75a83be45085`)
- **Name** - Display name (e.g., "Annie", "Brandon")
- **Gender** - Male, Female, Man, Woman
- **LooksCount** - Number of available looks

### Avatar Looks

A look is a specific version of an avatar with a particular outfit or style. Use the look ID (not the group ID) when generating videos.

Look properties:

- **ID** - Unique look identifier (use this for video generation)
- **AvatarType** - Engine compatibility: `studio_avatar`, `digital_twin`, `photo_avatar`
- **SupportedAPIEngines** - List of supported engines: `avatar_v`, `avatar_iv`, `avatar_iii`

## Listing Avatars

```go
import "github.com/plexusone/heygen-go/avatar"

// List public avatars
resp, err := client.Avatar.List(ctx, &avatar.ListOptions{
    Limit:     20,
    Ownership: "public",
})

// List your private avatars
resp, err := client.Avatar.List(ctx, &avatar.ListOptions{
    Ownership: "private",
})

// List all avatars (public + private)
resp, err := client.Avatar.List(ctx, &avatar.ListOptions{
    Ownership: "all",
})
```

### Pagination

Use the `Token` field to paginate through results:

```go
var allAvatars []avatar.Avatar

opts := &avatar.ListOptions{Limit: 100}
for {
    resp, err := client.Avatar.List(ctx, opts)
    if err != nil {
        return err
    }

    allAvatars = append(allAvatars, resp.Data...)

    if !resp.HasMore {
        break
    }
    opts.Token = resp.NextToken
}
```

## Getting Avatar Details

```go
// Get a specific avatar group
avatar, err := client.Avatar.Get(ctx, "e0e84faea390465896db75a83be45085")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Avatar: %s\n", avatar.Name)
fmt.Printf("Gender: %s\n", avatar.Gender)
fmt.Printf("Looks: %d\n", avatar.LooksCount)
```

## Listing Avatar Looks

```go
// Get looks for an avatar group
looks, err := client.Avatar.ListLooks(ctx, "e0e84faea390465896db75a83be45085", 50)
if err != nil {
    log.Fatal(err)
}

for _, look := range looks.Data {
    fmt.Printf("Look: %s (%s)\n", look.Name, look.ID)
    fmt.Printf("  Type: %s\n", look.AvatarType)
    fmt.Printf("  Engines: %v\n", look.SupportedAPIEngines)
}
```

## Avatar Types

| Type | Description | Use Case |
|------|-------------|----------|
| `studio_avatar` | Professional studio-quality avatars | Production videos |
| `digital_twin` | Custom-trained avatars | Personalized content |
| `photo_avatar` | Generated from a single photo | Quick prototyping |

## Best Practices

1. **Cache avatar lists** - Avatar data changes infrequently; cache for 1-24 hours
2. **Use look IDs for video** - Always use the specific look ID, not the group ID
3. **Check engine compatibility** - Verify the look supports your target API engine
4. **Handle pagination** - Large avatar libraries may require multiple requests
