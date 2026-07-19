package heygen_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/plexusone/heygen-go/heygen"
)

func TestClient_RequestURL(t *testing.T) {
	content := []byte("raw-binary-content")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Content-Type"); got != "audio/mpeg" {
			t.Errorf("Content-Type = %q, want %q", got, "audio/mpeg")
		}
		if got := r.Header.Get("X-Api-Key"); got != "test-key" {
			t.Errorf("X-Api-Key = %q, want %q", got, "test-key")
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !bytes.Equal(body, content) {
			t.Errorf("body = %q, want %q", body, content)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"ok": true}`))
	}))
	defer srv.Close()

	client, err := heygen.NewClient(heygen.Config{APIKey: "test-key"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	var result struct {
		OK bool `json:"ok"`
	}
	err = client.RequestURL(context.Background(), http.MethodPost, srv.URL+"/v1/asset", "audio/mpeg", content, &result)
	if err != nil {
		t.Fatalf("RequestURL() error = %v", err)
	}
	if !result.OK {
		t.Error("result.OK = false, want true")
	}
}

func TestClient_RetryRewindsBody(t *testing.T) {
	content := []byte("must-arrive-intact-on-retry")
	var attempts atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt := attempts.Add(1)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read body: %v", err)
		}
		if !bytes.Equal(body, content) {
			t.Errorf("attempt %d: body = %q, want %q", attempt, body, content)
		}
		if attempt == 1 {
			// 429 is retryable for all methods, including POST.
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client, err := heygen.NewClient(heygen.Config{APIKey: "test-key"})
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	err = client.RequestURL(context.Background(), http.MethodPost, srv.URL+"/v1/asset", "audio/mpeg", content, nil)
	if err != nil {
		t.Fatalf("RequestURL() error = %v", err)
	}
	if got := attempts.Load(); got != 2 {
		t.Errorf("attempts = %d, want 2 (one 429 retry)", got)
	}
}
