package heygen_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/plexusone/heygen-go/heygen"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  heygen.Config
		wantErr bool
	}{
		{
			name:    "empty config without env var",
			config:  heygen.Config{},
			wantErr: true,
		},
		{
			name: "config with API key",
			config: heygen.Config{
				APIKey: "test-key",
			},
			wantErr: false,
		},
		{
			name: "config with custom base URL",
			config: heygen.Config{
				APIKey:  "test-key",
				BaseURL: "https://custom.api.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HEYGEN_API_KEY", "")

			_, err := heygen.NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClient_FromEnv(t *testing.T) {
	t.Setenv("HEYGEN_API_KEY", "env-api-key")

	client, err := heygen.NewClient(heygen.Config{})
	if err != nil {
		t.Errorf("NewClient() error = %v, want nil", err)
	}
	if client == nil {
		t.Error("NewClient() returned nil client")
	}
}

func TestClient_Get(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		response   map[string]any
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful GET",
			path: "/v3/test",
			response: map[string]any{
				"data": map[string]any{"id": "123"},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "unauthorized",
			path: "/v3/test",
			response: map[string]any{
				"error": map[string]any{
					"code":    "unauthorized",
					"message": "Invalid API key",
				},
			},
			statusCode: http.StatusUnauthorized,
			wantErr:    true,
		},
		{
			name: "not found",
			path: "/v3/missing",
			response: map[string]any{
				"error": map[string]any{
					"code":    "not_found",
					"message": "Resource not found",
				},
			},
			statusCode: http.StatusNotFound,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify headers
				if r.Header.Get("X-API-KEY") != "test-key" {
					t.Errorf("Get() X-API-KEY = %v, want test-key", r.Header.Get("X-API-KEY"))
				}
				if !strings.HasPrefix(r.Header.Get("User-Agent"), "heygen-go/") {
					t.Errorf("Get() User-Agent = %v, want heygen-go/*", r.Header.Get("User-Agent"))
				}

				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client, _ := heygen.NewClient(heygen.Config{
				APIKey:  "test-key",
				BaseURL: server.URL,
			})

			var result map[string]any
			err := client.Get(context.Background(), tt.path, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_Post(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		body       map[string]any
		response   map[string]any
		statusCode int
		wantErr    bool
	}{
		{
			name: "successful POST",
			path: "/v3/create",
			body: map[string]any{"name": "test"},
			response: map[string]any{
				"data": map[string]any{"id": "new-123"},
			},
			statusCode: http.StatusOK,
			wantErr:    false,
		},
		{
			name: "rate limited",
			path: "/v3/create",
			body: map[string]any{"name": "test"},
			response: map[string]any{
				"error": map[string]any{
					"code":    "rate_limited",
					"message": "Too many requests",
				},
			},
			statusCode: http.StatusTooManyRequests,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("Post() method = %v, want POST", r.Method)
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Post() Content-Type = %v, want application/json", r.Header.Get("Content-Type"))
				}

				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client, _ := heygen.NewClient(heygen.Config{
				APIKey:  "test-key",
				BaseURL: server.URL,
			})

			var result map[string]any
			err := client.Post(context.Background(), tt.path, tt.body, &result)

			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client, _ := heygen.NewClient(heygen.Config{
		APIKey:  "test-key",
		BaseURL: server.URL,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	var result map[string]any
	err := client.Get(ctx, "/v3/slow", &result)

	if err == nil {
		t.Error("Get() expected context cancellation error, got nil")
	}
}

func TestConfig_Defaults(t *testing.T) {
	if heygen.DefaultBaseURL != "https://api.heygen.com" {
		t.Errorf("DefaultBaseURL = %v, want https://api.heygen.com", heygen.DefaultBaseURL)
	}

	if heygen.EnvAPIKey != "HEYGEN_API_KEY" {
		t.Errorf("EnvAPIKey = %v, want HEYGEN_API_KEY", heygen.EnvAPIKey)
	}

	if heygen.DefaultTimeout != 30*time.Second {
		t.Errorf("DefaultTimeout = %v, want 30s", heygen.DefaultTimeout)
	}
}

func TestDefaultRetryConfig(t *testing.T) {
	cfg := heygen.DefaultRetryConfig()

	if cfg.MaxRetries != 2 {
		t.Errorf("MaxRetries = %v, want 2", cfg.MaxRetries)
	}
	if cfg.BaseDelay != 1*time.Second {
		t.Errorf("BaseDelay = %v, want 1s", cfg.BaseDelay)
	}
	if cfg.MaxDelay != 30*time.Second {
		t.Errorf("MaxDelay = %v, want 30s", cfg.MaxDelay)
	}
}
