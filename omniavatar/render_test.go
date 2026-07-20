package omniavatar

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	heygenvideo "github.com/plexusone/heygen-go/video"

	"github.com/plexusone/omniavatar-core/render"
)

func TestMapVideoStatus(t *testing.T) {
	tests := []struct {
		in   heygenvideo.Status
		want render.JobState
	}{
		{heygenvideo.StatusPending, render.JobStatePending},
		{heygenvideo.StatusProcessing, render.JobStateProcessing},
		{heygenvideo.StatusCompleted, render.JobStateCompleted},
		{heygenvideo.StatusFailed, render.JobStateFailed},
		{heygenvideo.Status("unknown"), render.JobStateProcessing},
	}
	for _, tt := range tests {
		if got := mapVideoStatus(tt.in); got != tt.want {
			t.Errorf("mapVideoStatus(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestVideoToStatus(t *testing.T) {
	video := &heygenvideo.Video{
		VideoID:  "vid-1",
		Status:   heygenvideo.StatusCompleted,
		VideoURL: "https://x/v.mp4",
		Duration: 12.5,
	}
	status := videoToStatus(video)
	if status.State != render.JobStateCompleted {
		t.Errorf("State = %q, want %q", status.State, render.JobStateCompleted)
	}
	if status.VideoURL != "https://x/v.mp4" || status.Duration != 12.5 {
		t.Errorf("videoToStatus() = %+v, want URL and duration mapped", status)
	}
}

func TestUploadAudio(t *testing.T) {
	content := []byte("fake-mp3-bytes")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/asset" {
			t.Errorf("path = %s, want /v1/asset", r.URL.Path)
		}
		if got := r.Header.Get("Content-Type"); got != "audio/mpeg" {
			t.Errorf("Content-Type = %q, want %q", got, "audio/mpeg")
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !bytes.Equal(body, content) {
			t.Errorf("body = %q, want %q", body, content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code": 100, "data": {"id": "a1", "file_type": "audio", "url": "https://resource.heygen.ai/a1.mp3"}}`))
	}))
	defer srv.Close()

	provider, err := NewRenderProvider(RenderConfig{
		APIKey:        "test-key",
		UploadBaseURL: srv.URL,
	})
	if err != nil {
		t.Fatalf("NewRenderProvider() error = %v", err)
	}

	url, err := provider.UploadAudio(context.Background(), "narration.mp3", bytes.NewReader(content))
	if err != nil {
		t.Fatalf("UploadAudio() error = %v", err)
	}
	if url != "https://resource.heygen.ai/a1.mp3" {
		t.Errorf("UploadAudio() = %q, want hosted URL", url)
	}
}

func TestRenderProviderImplementsAudioUploader(t *testing.T) {
	provider, err := NewRenderProvider(RenderConfig{APIKey: "test-key"})
	if err != nil {
		t.Fatalf("NewRenderProvider() error = %v", err)
	}
	if _, ok := any(provider).(render.AudioUploader); !ok {
		t.Error("RenderProvider does not implement render.AudioUploader")
	}
}

func TestListAvatars(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/avatars" {
			t.Errorf("path = %s, want /v2/avatars", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"avatars":[
			{"avatar_id":"Abigail_expressive_2024112501","avatar_name":"Abigail (Upper Body)","gender":"female"},
			{"avatar_id":"Marco_public_1","avatar_name":"Marco","gender":"male"}
		],"talking_photos":[]}}`))
	}))
	defer srv.Close()

	provider, err := NewRenderProvider(RenderConfig{APIKey: "test-key", BaseURL: srv.URL})
	if err != nil {
		t.Fatalf("NewRenderProvider() error = %v", err)
	}

	all, err := provider.ListAvatars(context.Background(), "")
	if err != nil {
		t.Fatalf("ListAvatars() error = %v", err)
	}
	if len(all) != 2 {
		t.Fatalf("ListAvatars(\"\") = %d, want 2", len(all))
	}
	if all[0].ID != "Abigail_expressive_2024112501" || all[0].Gender != "female" {
		t.Errorf("avatar[0] = %+v", all[0])
	}

	got, err := provider.ListAvatars(context.Background(), "marco")
	if err != nil {
		t.Fatalf("ListAvatars(\"marco\") error = %v", err)
	}
	if len(got) != 1 || got[0].ID != "Marco_public_1" {
		t.Errorf("ListAvatars(\"marco\") = %+v, want Marco only", got)
	}
}

func TestVideoToStatusFailed(t *testing.T) {
	video := &heygenvideo.Video{
		VideoID: "vid-1",
		Status:  heygenvideo.StatusFailed,
		Error:   &heygenvideo.VideoError{Code: "E1", Message: "boom"},
	}
	status := videoToStatus(video)
	if status.State != render.JobStateFailed {
		t.Errorf("State = %q, want %q", status.State, render.JobStateFailed)
	}
	if status.ErrorCode != "E1" || status.ErrorMsg != "boom" {
		t.Errorf("videoToStatus() error fields = %q/%q, want E1/boom", status.ErrorCode, status.ErrorMsg)
	}
}
