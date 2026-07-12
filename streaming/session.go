// Package streaming provides access to HeyGen streaming (LiveAvatar) APIs.
package streaming

import (
	"context"
	"fmt"

	"github.com/plexusone/heygen-go/heygen"
)

// Client provides access to streaming APIs.
type Client struct {
	client *heygen.Client
}

// NewClient creates a new streaming client.
func NewClient(client *heygen.Client) *Client {
	return &Client{client: client}
}

// Quality represents the streaming quality level.
type Quality string

const (
	QualityLow    Quality = "low"
	QualityMedium Quality = "medium"
	QualityHigh   Quality = "high"
)

// NewSessionRequest is the request to create a new streaming session.
type NewSessionRequest struct {
	// Quality is the video quality (low, medium, high).
	Quality Quality `json:"quality,omitempty"`

	// AvatarID is the avatar to use (optional, uses default if not specified).
	AvatarID string `json:"avatar_id,omitempty"`

	// VoiceID is the voice to use (optional).
	VoiceID string `json:"voice_id,omitempty"`

	// KnowledgeBaseID links a knowledge base for AI responses.
	KnowledgeBaseID string `json:"knowledge_base_id,omitempty"`

	// Language is the language code (e.g., "en").
	Language string `json:"language,omitempty"`

	// Version is the API version (default "v2").
	Version string `json:"version,omitempty"`
}

// Session represents a streaming session.
type Session struct {
	// SessionID is the unique session identifier.
	SessionID string `json:"session_id"`

	// AccessToken is the token for WebSocket authentication.
	AccessToken string `json:"access_token,omitempty"`

	// URL is the WebSocket URL for real-time streaming.
	URL string `json:"url,omitempty"`

	// RealtimeEndpoint is the endpoint for real-time communication.
	RealtimeEndpoint string `json:"realtime_endpoint,omitempty"`

	// SDPOffer is the SDP offer for WebRTC connection.
	SDPOffer string `json:"sdp_offer,omitempty"`

	// ICEServers contains ICE server configurations.
	ICEServers []ICEServer `json:"ice_servers,omitempty"`
}

// ICEServer represents a STUN/TURN server for WebRTC.
type ICEServer struct {
	URLs       []string `json:"urls"`
	Username   string   `json:"username,omitempty"`
	Credential string   `json:"credential,omitempty"`
}

// NewSessionResponse is the response from creating a new session.
type NewSessionResponse struct {
	Data Session `json:"data"`
}

// NewSession creates a new streaming session.
func (c *Client) NewSession(ctx context.Context, req NewSessionRequest) (*Session, error) {
	if req.Quality == "" {
		req.Quality = QualityMedium
	}
	if req.Version == "" {
		req.Version = "v2"
	}

	var resp NewSessionResponse
	if err := c.client.Post(ctx, "/v1/streaming.new", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

// SDP represents a Session Description Protocol message.
type SDP struct {
	Type string `json:"type"` // "offer" or "answer"
	SDP  string `json:"sdp"`
}

// StartRequest is the request to start a streaming session.
type StartRequest struct {
	// SessionID is the session to start.
	SessionID string `json:"session_id"`

	// SDP is the SDP offer for WebRTC connection.
	SDP SDP `json:"sdp"`
}

// StartResponse is the response from starting a session.
type StartResponse struct {
	Data struct {
		// SDP is the SDP answer for WebRTC connection.
		SDP SDP `json:"sdp"`
	} `json:"data"`
}

// Start starts a streaming session with WebRTC negotiation.
func (c *Client) Start(ctx context.Context, req StartRequest) (*SDP, error) {
	var resp StartResponse
	if err := c.client.Post(ctx, "/v1/streaming.start", req, &resp); err != nil {
		return nil, err
	}
	return &resp.Data.SDP, nil
}

// ICECandidate represents an ICE candidate for WebRTC.
type ICECandidate struct {
	Candidate        string `json:"candidate"`
	SDPMid           string `json:"sdpMid"`
	SDPMLineIndex    int    `json:"sdpMLineIndex"`
	UsernameFragment string `json:"usernameFragment,omitempty"`
}

// ICERequest is the request to send an ICE candidate.
type ICERequest struct {
	// SessionID is the session ID.
	SessionID string `json:"session_id"`

	// Candidate is the ICE candidate.
	Candidate ICECandidate `json:"candidate"`
}

// SendICE sends an ICE candidate for WebRTC connection.
func (c *Client) SendICE(ctx context.Context, req ICERequest) error {
	return c.client.Post(ctx, "/v1/streaming.ice", req, nil)
}

// TaskRequest is the request to send a task (speak text).
type TaskRequest struct {
	// SessionID is the session ID.
	SessionID string `json:"session_id"`

	// Text is the text for the avatar to speak.
	Text string `json:"text"`

	// TaskType specifies the task type (optional).
	TaskType string `json:"task_type,omitempty"`
}

// TaskResponse is the response from sending a task.
type TaskResponse struct {
	Data struct {
		TaskID string `json:"task_id"`
	} `json:"data"`
}

// SendTask sends text for the avatar to speak.
func (c *Client) SendTask(ctx context.Context, req TaskRequest) (string, error) {
	var resp TaskResponse
	if err := c.client.Post(ctx, "/v1/streaming.task", req, &resp); err != nil {
		return "", err
	}
	return resp.Data.TaskID, nil
}

// InterruptRequest is the request to interrupt the avatar.
type InterruptRequest struct {
	// SessionID is the session to interrupt.
	SessionID string `json:"session_id"`
}

// Interrupt stops the avatar from speaking.
func (c *Client) Interrupt(ctx context.Context, sessionID string) error {
	req := InterruptRequest{SessionID: sessionID}
	return c.client.Post(ctx, "/v1/streaming.interrupt", req, nil)
}

// StopRequest is the request to stop a session.
type StopRequest struct {
	// SessionID is the session to stop.
	SessionID string `json:"session_id"`
}

// Stop ends a streaming session.
func (c *Client) Stop(ctx context.Context, sessionID string) error {
	req := StopRequest{SessionID: sessionID}
	return c.client.Post(ctx, "/v1/streaming.stop", req, nil)
}

// ListResponse is the response from listing sessions.
type ListResponse struct {
	Data struct {
		Sessions []SessionInfo `json:"sessions"`
	} `json:"data"`
}

// SessionInfo contains basic session information.
type SessionInfo struct {
	SessionID string `json:"session_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// List returns all active streaming sessions.
func (c *Client) List(ctx context.Context) ([]SessionInfo, error) {
	var resp ListResponse
	if err := c.client.Get(ctx, "/v1/streaming.list", &resp); err != nil {
		return nil, err
	}
	return resp.Data.Sessions, nil
}

// CreateTokenResponse is the response from creating a session token.
type CreateTokenResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

// CreateToken creates a token for browser-based streaming.
func (c *Client) CreateToken(ctx context.Context) (string, error) {
	var resp CreateTokenResponse
	if err := c.client.Post(ctx, "/v1/streaming.create_token", struct{}{}, &resp); err != nil {
		return "", err
	}
	return resp.Data.Token, nil
}

// Speak is a convenience method to send text for the avatar to speak.
func (c *Client) Speak(ctx context.Context, sessionID, text string) (string, error) {
	return c.SendTask(ctx, TaskRequest{
		SessionID: sessionID,
		Text:      text,
	})
}

// Close is an alias for Stop for interface compatibility.
func (c *Client) Close(ctx context.Context, sessionID string) error {
	return c.Stop(ctx, sessionID)
}

// SessionConfig holds configuration for creating a LiveAvatar session.
type SessionConfig struct {
	// AvatarID is the avatar to use.
	AvatarID string

	// VoiceID is the voice to use.
	VoiceID string

	// Quality is the video quality.
	Quality Quality

	// Language is the language code.
	Language string
}

// CreateSession creates and returns a new session with the given configuration.
func (c *Client) CreateSession(ctx context.Context, cfg SessionConfig) (*Session, error) {
	req := NewSessionRequest{
		AvatarID: cfg.AvatarID,
		VoiceID:  cfg.VoiceID,
		Quality:  cfg.Quality,
		Language: cfg.Language,
	}
	return c.NewSession(ctx, req)
}

// String returns a string representation of the session.
func (s *Session) String() string {
	return fmt.Sprintf("Session{ID: %s}", s.SessionID)
}
