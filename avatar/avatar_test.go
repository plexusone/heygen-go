package avatar_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/plexusone/heygen-go/avatar"
	"github.com/plexusone/heygen-go/heygen"
)

func TestClient_List(t *testing.T) {
	tests := []struct {
		name      string
		opts      *avatar.ListOptions
		response  map[string]any
		wantPath  string
		wantCount int
		wantErr   bool
	}{
		{
			name: "list with no options",
			opts: nil,
			response: map[string]any{
				"data": []map[string]any{
					{"id": "abc123", "name": "Avatar1", "looks_count": 5},
					{"id": "def456", "name": "Avatar2", "looks_count": 3},
				},
				"has_more": false,
			},
			wantPath:  "/v3/avatars",
			wantCount: 2,
		},
		{
			name: "list with limit",
			opts: &avatar.ListOptions{Limit: 10},
			response: map[string]any{
				"data": []map[string]any{
					{"id": "abc123", "name": "Avatar1", "looks_count": 5},
				},
				"has_more":   true,
				"next_token": "token123",
			},
			wantPath:  "/v3/avatars?limit=10",
			wantCount: 1,
		},
		{
			name: "list with ownership filter",
			opts: &avatar.ListOptions{Ownership: "public"},
			response: map[string]any{
				"data":     []map[string]any{},
				"has_more": false,
			},
			wantPath:  "/v3/avatars?ownership=public",
			wantCount: 0,
		},
		{
			name: "list with all options",
			opts: &avatar.ListOptions{
				Limit:     20,
				Token:     "page2",
				Ownership: "private",
			},
			response: map[string]any{
				"data": []map[string]any{
					{"id": "xyz789", "name": "MyAvatar", "looks_count": 1},
				},
				"has_more": false,
			},
			wantPath:  "/v3/avatars?limit=20&token=page2&ownership=private",
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.String()
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			heygenClient, err := heygen.NewClient(heygen.Config{
				APIKey:  "test-key",
				BaseURL: server.URL,
			})
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			client := avatar.NewClient(heygenClient)
			resp, err := client.List(context.Background(), tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotPath != tt.wantPath {
				t.Errorf("List() path = %v, want %v", gotPath, tt.wantPath)
			}

			if len(resp.Data) != tt.wantCount {
				t.Errorf("List() count = %v, want %v", len(resp.Data), tt.wantCount)
			}
		})
	}
}

func TestClient_Get(t *testing.T) {
	tests := []struct {
		name     string
		groupID  string
		response map[string]any
		wantID   string
		wantName string
		wantErr  bool
	}{
		{
			name:    "get avatar by ID",
			groupID: "abc123",
			response: map[string]any{
				"data": map[string]any{
					"id":          "abc123",
					"name":        "TestAvatar",
					"gender":      "female",
					"looks_count": 10,
				},
			},
			wantID:   "abc123",
			wantName: "TestAvatar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			heygenClient, err := heygen.NewClient(heygen.Config{
				APIKey:  "test-key",
				BaseURL: server.URL,
			})
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			client := avatar.NewClient(heygenClient)
			avatar, err := client.Get(context.Background(), tt.groupID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			wantPath := "/v3/avatars/" + tt.groupID
			if gotPath != wantPath {
				t.Errorf("Get() path = %v, want %v", gotPath, wantPath)
			}

			if avatar.ID != tt.wantID {
				t.Errorf("Get() ID = %v, want %v", avatar.ID, tt.wantID)
			}

			if avatar.Name != tt.wantName {
				t.Errorf("Get() Name = %v, want %v", avatar.Name, tt.wantName)
			}
		})
	}
}

func TestClient_ListLooks(t *testing.T) {
	tests := []struct {
		name      string
		groupID   string
		limit     int
		response  map[string]any
		wantPath  string
		wantCount int
		wantErr   bool
	}{
		{
			name:    "list looks without limit",
			groupID: "abc123",
			limit:   0,
			response: map[string]any{
				"data": []map[string]any{
					{"id": "look1", "name": "Casual", "avatar_type": "studio_avatar"},
					{"id": "look2", "name": "Formal", "avatar_type": "studio_avatar"},
				},
				"has_more": false,
			},
			wantPath:  "/v3/avatars/abc123/looks",
			wantCount: 2,
		},
		{
			name:    "list looks with limit",
			groupID: "def456",
			limit:   5,
			response: map[string]any{
				"data": []map[string]any{
					{"id": "look1", "name": "Style1", "avatar_type": "digital_twin"},
				},
				"has_more":   true,
				"next_token": "more",
			},
			wantPath:  "/v3/avatars/def456/looks?limit=5",
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.String()
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(tt.response)
			}))
			defer server.Close()

			heygenClient, err := heygen.NewClient(heygen.Config{
				APIKey:  "test-key",
				BaseURL: server.URL,
			})
			if err != nil {
				t.Fatalf("NewClient() error = %v", err)
			}

			client := avatar.NewClient(heygenClient)
			resp, err := client.ListLooks(context.Background(), tt.groupID, tt.limit)

			if (err != nil) != tt.wantErr {
				t.Errorf("ListLooks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotPath != tt.wantPath {
				t.Errorf("ListLooks() path = %v, want %v", gotPath, tt.wantPath)
			}

			if len(resp.Data) != tt.wantCount {
				t.Errorf("ListLooks() count = %v, want %v", len(resp.Data), tt.wantCount)
			}
		})
	}
}

func TestClient_List_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"code":    "unauthorized",
				"message": "Invalid API key",
			},
		})
	}))
	defer server.Close()

	heygenClient, err := heygen.NewClient(heygen.Config{
		APIKey:  "bad-key",
		BaseURL: server.URL,
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}

	client := avatar.NewClient(heygenClient)
	_, err = client.List(context.Background(), nil)

	if err == nil {
		t.Error("List() expected error for unauthorized, got nil")
	}
}
