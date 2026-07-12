package heygen

import (
	"errors"
	"fmt"
)

// Common errors returned by the client.
var (
	// ErrNoAPIKey is returned when no API key is configured.
	ErrNoAPIKey = errors.New("heygen: no API key configured")

	// ErrRateLimited is returned when the API rate limit is exceeded.
	ErrRateLimited = errors.New("heygen: rate limit exceeded")

	// ErrUnauthorized is returned when authentication fails.
	ErrUnauthorized = errors.New("heygen: unauthorized")
)

// APIError represents an error response from the HeyGen API.
type APIError struct {
	// StatusCode is the HTTP status code.
	StatusCode int

	// Code is the error code from the API.
	Code string

	// Message is the human-readable error message.
	Message string

	// RequestID is the request ID from the X-Request-Id header.
	RequestID string
}

// Error implements the error interface.
func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("heygen: %s (code=%s, status=%d, request_id=%s)",
			e.Message, e.Code, e.StatusCode, e.RequestID)
	}
	return fmt.Sprintf("heygen: %s (status=%d, request_id=%s)",
		e.Message, e.StatusCode, e.RequestID)
}

// IsNotFound returns true if the error is a not found error.
func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 404 || apiErr.Code == "not_found"
	}
	return false
}

// IsRateLimited returns true if the error is a rate limit error.
func IsRateLimited(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 429 || apiErr.Code == "rate_limit_exceeded"
	}
	return errors.Is(err, ErrRateLimited)
}

// IsUnauthorized returns true if the error is an authentication error.
func IsUnauthorized(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 401 || apiErr.Code == "unauthorized"
	}
	return errors.Is(err, ErrUnauthorized)
}
