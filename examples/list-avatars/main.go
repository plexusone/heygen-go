// Example: List all available HeyGen avatars using the v3 API.
//
// Usage:
//
//	export HEYGEN_API_KEY="your-api-key"
//	go run ./examples/list-avatars
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	heygen "github.com/plexusone/heygen-go"
	"github.com/plexusone/heygen-go/avatar"
)

func main() {
	apiKey := os.Getenv("HEYGEN_API_KEY")
	if apiKey == "" {
		log.Fatal("HEYGEN_API_KEY environment variable is required")
	}

	client, err := heygen.New(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// List avatars (v3 API)
	fmt.Println("=== Avatars (v3 API) ===")
	resp, err := client.Avatar.List(ctx, &avatar.ListOptions{
		Limit: 10,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, a := range resp.Data {
		fmt.Printf("  %s (%s) - %s, %d looks\n", a.Name, a.ID, a.Gender, a.LooksCount)
	}
	fmt.Printf("Total: %d avatars (has_more: %v)\n", len(resp.Data), resp.HasMore)
}
