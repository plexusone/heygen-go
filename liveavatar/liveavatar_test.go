//nolint:gosec // G101: test data contains fake tokens/credentials
package liveavatar_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/plexusone/heygen-go/liveavatar"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *liveavatar.Config
		wantErr bool
	}{
		{
			name:    "nil config without env var",
			config:  nil,
			wantErr: true, // no API key
		},
		{
			name: "config with API key",
			config: &liveavatar.Config{
				APIKey: "test-key",
			},
			wantErr: false,
		},
		{
			name: "config with custom base URL",
			config: &liveavatar.Config{
				APIKey:  "test-key",
				BaseURL: "https://custom.api.com",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear env var for consistent testing
			t.Setenv("LIVEAVATAR_API_KEY", "")

			_, err := liveavatar.NewClient(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewClient_FromEnv(t *testing.T) {
	t.Setenv("LIVEAVATAR_API_KEY", "env-api-key")

	client, err := liveavatar.NewClient(nil)
	if err != nil {
		t.Errorf("NewClient() error = %v, want nil", err)
	}
	if client == nil {
		t.Error("NewClient() returned nil client")
	}
}

func TestClient_NewSession(t *testing.T) {
	tests := []struct {
		name        string
		request     *liveavatar.NewSessionRequest
		response    map[string]any
		statusCode  int
		wantErr     bool
		wantSession string
	}{
		{
			name: "successful session creation",
			request: &liveavatar.NewSessionRequest{
				Mode:      "LITE",
				AvatarID:  "65f9e3c9-d48b-4118-b73a-4ae2e3cbb8f0",
				IsSandbox: true,
				LiveKitConfig: &liveavatar.LiveKitConfig{
					LiveKitURL:         "wss://test.livekit.cloud",
					LiveKitRoom:        "test-room",
					LiveKitClientToken: "test-token",
				},
			},
			response: map[string]any{
				"code": 1000,
				"data": map[string]any{
					"session_id":    "session-123",
					"session_token": "jwt-token-here",
				},
				"message": "Session token created successfully",
			},
			statusCode:  http.StatusOK,
			wantErr:     false,
			wantSession: "session-123",
		},
		{
			name: "validation error",
			request: &liveavatar.NewSessionRequest{
				Mode:     "LITE",
				AvatarID: "invalid",
			},
			response: map[string]any{
				"code":    4000,
				"data":    nil,
				"message": "Request validation errors",
			},
			statusCode: http.StatusOK,
			wantErr:    true,
		},
		{
			name: "unauthorized",
			request: &liveavatar.NewSessionRequest{
				Mode:      "LITE",
				AvatarID:  "65f9e3c9-d48b-4118-b73a-4ae2e3cbb8f0",
				IsSandbox: true,
			},
			response: map[string]any{
				"code":    4010,
				"data":    nil,
				"message": "Invalid credentials",
			},
			statusCode: http.StatusOK,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				if r.Method != "POST" {
					t.Errorf("NewSession() method = %v, want POST", r.Method)
				}
				if r.URL.Path != "/v1/sessions/token" {
					t.Errorf("NewSession() path = %v, want /v1/sessions/token", r.URL.Path)
				}
				if r.Header.Get("X-API-KEY") == "" {
					t.Error("NewSession() missing X-API-KEY header")
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("NewSession() Content-Type = %v, want application/json", r.Header.Get("Content-Type"))
				}

				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client, err := liveavatar.NewClient(&liveavatar.Config{
				APIKey:  "test-key",
				BaseURL: server.URL,
			})
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			resp, err := client.NewSession(context.Background(), tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && resp.SessionID != tt.wantSession {
				t.Errorf("NewSession() SessionID = %v, want %v", resp.SessionID, tt.wantSession)
			}
		})
	}
}

func TestClient_StartSession(t *testing.T) {
	tests := []struct {
		name      string
		sessionID string
		token     string
		response  map[string]any
		wantWSURL string
		wantErr   bool
	}{
		{
			name:      "successful start",
			sessionID: "session-123",
			token:     "jwt-token",
			response: map[string]any{
				"code": 1000,
				"data": map[string]any{
					"session_id":           "session-123",
					"ws_url":               "wss://webrtc.heygen.io/session/abc",
					"max_session_duration": 60,
				},
				"message": "Session created successfully",
			},
			wantWSURL: "wss://webrtc.heygen.io/session/abc",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("StartSession() method = %v, want POST", r.Method)
				}
				if r.URL.Path != "/v1/sessions/start" {
					t.Errorf("StartSession() path = %v, want /v1/sessions/start", r.URL.Path)
				}
				if r.Header.Get("Authorization") != "Bearer "+tt.token {
					t.Errorf("StartSession() Authorization = %v, want Bearer %v", r.Header.Get("Authorization"), tt.token)
				}

				_ = json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client, _ := liveavatar.NewClient(&liveavatar.Config{
				APIKey:  "test-key",
				BaseURL: server.URL,
			})

			resp, err := client.StartSession(context.Background(), tt.sessionID, tt.token)

			if (err != nil) != tt.wantErr {
				t.Errorf("StartSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && resp.WSURL != tt.wantWSURL {
				t.Errorf("StartSession() WSURL = %v, want %v", resp.WSURL, tt.wantWSURL)
			}
		})
	}
}

func TestClient_StopSession(t *testing.T) {
	tests := []struct {
		name      string
		sessionID string
		token     string
		reason    liveavatar.StopReason
		response  map[string]any
		wantErr   bool
	}{
		{
			name:      "successful stop",
			sessionID: "session-123",
			token:     "jwt-token",
			reason:    liveavatar.StopReasonUserDisconnected,
			response: map[string]any{
				"code":    1000,
				"data":    nil,
				"message": "Successfully stopped session",
			},
			wantErr: false,
		},
		{
			name:      "stop with session ended reason",
			sessionID: "session-456",
			token:     "jwt-token-2",
			reason:    liveavatar.StopReasonSessionEnded,
			response: map[string]any{
				"code":    1000,
				"data":    nil,
				"message": "Successfully stopped session",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotBody map[string]string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Errorf("StopSession() method = %v, want POST", r.Method)
				}
				if r.URL.Path != "/v1/sessions/stop" {
					t.Errorf("StopSession() path = %v, want /v1/sessions/stop", r.URL.Path)
				}

				_ = json.NewDecoder(r.Body).Decode(&gotBody)
				_ = json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			client, _ := liveavatar.NewClient(&liveavatar.Config{
				APIKey:  "test-key",
				BaseURL: server.URL,
			})

			err := client.StopSession(context.Background(), tt.sessionID, tt.token, tt.reason)

			if (err != nil) != tt.wantErr {
				t.Errorf("StopSession() error = %v, wantErr %v", err, tt.wantErr)
			}

			if gotBody["session_id"] != tt.sessionID {
				t.Errorf("StopSession() session_id = %v, want %v", gotBody["session_id"], tt.sessionID)
			}
			if gotBody["reason"] != string(tt.reason) {
				t.Errorf("StopSession() reason = %v, want %v", gotBody["reason"], tt.reason)
			}
		})
	}
}

func TestConstants(t *testing.T) {
	if liveavatar.DefaultBaseURL != "https://api.liveavatar.com" {
		t.Errorf("DefaultBaseURL = %v, want https://api.liveavatar.com", liveavatar.DefaultBaseURL)
	}

	if liveavatar.EnvAPIKey != "LIVEAVATAR_API_KEY" {
		t.Errorf("EnvAPIKey = %v, want LIVEAVATAR_API_KEY", liveavatar.EnvAPIKey)
	}

	if liveavatar.SandboxAvatarID != "65f9e3c9-d48b-4118-b73a-4ae2e3cbb8f0" {
		t.Errorf("SandboxAvatarID = %v, want 65f9e3c9-d48b-4118-b73a-4ae2e3cbb8f0", liveavatar.SandboxAvatarID)
	}
}

func TestVideoQuality(t *testing.T) {
	qualities := []liveavatar.VideoQuality{
		liveavatar.QualityVeryHigh,
		liveavatar.QualityHigh,
		liveavatar.QualityMedium,
		liveavatar.QualityLow,
	}

	expected := []string{"very_high", "high", "medium", "low"}

	for i, q := range qualities {
		if string(q) != expected[i] {
			t.Errorf("VideoQuality %d = %v, want %v", i, q, expected[i])
		}
	}
}
