package avatar_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	heygen "github.com/plexusone/heygen-go/heygen"

	"github.com/plexusone/heygen-go/avatar"
)

const v2AvatarsBody = `{
  "error": null,
  "data": {
    "avatars": [
      {"avatar_id": "Abigail_expressive_2024112501", "avatar_name": "Abigail (Upper Body)", "gender": "female"},
      {"avatar_id": "Marco_public_1", "avatar_name": "Marco in Suit", "gender": "male"}
    ],
    "talking_photos": [
      {"talking_photo_id": "tp_123", "talking_photo_name": "My Photo"}
    ]
  }
}`

func newTestClient(t *testing.T, baseURL string) *avatar.Client {
	t.Helper()
	core, err := heygen.NewClient(heygen.Config{APIKey: "test-key", BaseURL: baseURL})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	return avatar.NewClient(core)
}

func TestListV2(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/avatars" {
			t.Errorf("path = %s, want /v2/avatars", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(v2AvatarsBody))
	}))
	defer srv.Close()

	resp, err := newTestClient(t, srv.URL).ListV2(context.Background())
	if err != nil {
		t.Fatalf("ListV2() error = %v", err)
	}
	if len(resp.Data.Avatars) != 2 {
		t.Fatalf("avatars = %d, want 2", len(resp.Data.Avatars))
	}
	if resp.Data.Avatars[0].AvatarID != "Abigail_expressive_2024112501" {
		t.Errorf("avatar[0].AvatarID = %q", resp.Data.Avatars[0].AvatarID)
	}
	if len(resp.Data.TalkingPhotos) != 1 || resp.Data.TalkingPhotos[0].TalkingPhotoID != "tp_123" {
		t.Errorf("talking photos = %+v", resp.Data.TalkingPhotos)
	}
}

func TestSearchV2(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(v2AvatarsBody))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)

	all, err := client.SearchV2(context.Background(), "")
	if err != nil {
		t.Fatalf("SearchV2(\"\") error = %v", err)
	}
	if len(all) != 2 {
		t.Errorf("SearchV2(\"\") = %d avatars, want 2", len(all))
	}

	// Case-insensitive match on name.
	got, err := client.SearchV2(context.Background(), "abigail")
	if err != nil {
		t.Fatalf("SearchV2 error = %v", err)
	}
	if len(got) != 1 || got[0].AvatarID != "Abigail_expressive_2024112501" {
		t.Errorf("SearchV2(\"abigail\") = %+v, want the Abigail avatar", got)
	}

	// Match on ID substring.
	got, err = client.SearchV2(context.Background(), "marco_public")
	if err != nil {
		t.Fatalf("SearchV2 error = %v", err)
	}
	if len(got) != 1 || got[0].AvatarID != "Marco_public_1" {
		t.Errorf("SearchV2(\"marco_public\") = %+v, want Marco", got)
	}
}
