package asset

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/plexusone/heygen-go/heygen"
)

func newTestClient(t *testing.T, baseURL string) *Client {
	t.Helper()
	client, err := heygen.NewClient(heygen.Config{APIKey: "test-key"})
	if err != nil {
		t.Fatalf("heygen.NewClient: %v", err)
	}
	return NewClient(client, WithBaseURL(baseURL))
}

func TestUpload(t *testing.T) {
	content := []byte("fake-mp3-bytes")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if r.URL.Path != "/v1/asset" {
			t.Errorf("path = %s, want /v1/asset", r.URL.Path)
		}
		if got := r.Header.Get("X-Api-Key"); got != "test-key" {
			t.Errorf("X-Api-Key = %q, want %q", got, "test-key")
		}
		if got := r.Header.Get("Content-Type"); got != ContentTypeMPEG {
			t.Errorf("Content-Type = %q, want %q", got, ContentTypeMPEG)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !bytes.Equal(body, content) {
			t.Errorf("body = %q, want %q", body, content)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"code": 100,
			"data": {
				"id": "asset-123",
				"name": "narration",
				"file_type": "audio",
				"url": "https://resource.heygen.ai/asset-123.mp3"
			},
			"msg": null
		}`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)

	uploaded, err := client.Upload(context.Background(), ContentTypeMPEG, bytes.NewReader(content))
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}
	if uploaded.ID != "asset-123" {
		t.Errorf("ID = %q, want %q", uploaded.ID, "asset-123")
	}
	if uploaded.FileType != "audio" {
		t.Errorf("FileType = %q, want %q", uploaded.FileType, "audio")
	}
	if uploaded.URL != "https://resource.heygen.ai/asset-123.mp3" {
		t.Errorf("URL = %q, want hosted URL", uploaded.URL)
	}
}

func TestUploadMissingURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"code": 400123, "data": null, "msg": "unsupported file type"}`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)

	if _, err := client.Upload(context.Background(), "application/zip", bytes.NewReader([]byte("x"))); err == nil {
		t.Fatal("Upload() error = nil, want error for missing asset URL")
	}
}

func TestUploadAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"code": "unauthorized", "message": "invalid api key"}`))
	}))
	defer srv.Close()

	client := newTestClient(t, srv.URL)

	_, err := client.Upload(context.Background(), ContentTypeMPEG, bytes.NewReader([]byte("x")))
	if err == nil {
		t.Fatal("Upload() error = nil, want APIError")
	}
	var apiErr *heygen.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("Upload() error = %v (%T), want *heygen.APIError", err, err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized {
		t.Errorf("StatusCode = %d, want 401", apiErr.StatusCode)
	}
	if !heygen.IsUnauthorized(err) {
		t.Error("IsUnauthorized(err) = false, want true")
	}
}
