package heygen_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/plexusone/heygen-go/heygen"
)

func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name        string
		err         *heygen.APIError
		wantContain string
	}{
		{
			name: "error with code",
			err: &heygen.APIError{
				Code:       "unauthorized",
				Message:    "Invalid API key",
				StatusCode: 401,
			},
			wantContain: "unauthorized",
		},
		{
			name: "error with request ID",
			err: &heygen.APIError{
				Code:       "rate_limited",
				Message:    "Too many requests",
				StatusCode: 429,
				RequestID:  "req-123",
			},
			wantContain: "req-123",
		},
		{
			name: "error without code",
			err: &heygen.APIError{
				Message:    "Server error",
				StatusCode: 500,
			},
			wantContain: "500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if !strings.Contains(got, tt.wantContain) {
				t.Errorf("APIError.Error() = %v, want to contain %v", got, tt.wantContain)
			}
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "401 error",
			err:  &heygen.APIError{StatusCode: 401},
			want: true,
		},
		{
			name: "unauthorized code",
			err:  &heygen.APIError{Code: "unauthorized", StatusCode: 400},
			want: true,
		},
		{
			name: "200 error",
			err:  &heygen.APIError{StatusCode: 200},
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "ErrUnauthorized sentinel",
			err:  heygen.ErrUnauthorized,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := heygen.IsUnauthorized(tt.err); got != tt.want {
				t.Errorf("IsUnauthorized() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsRateLimited(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "429 error",
			err:  &heygen.APIError{StatusCode: 429},
			want: true,
		},
		{
			name: "rate_limit_exceeded code",
			err:  &heygen.APIError{Code: "rate_limit_exceeded", StatusCode: 400},
			want: true,
		},
		{
			name: "401 error",
			err:  &heygen.APIError{StatusCode: 401},
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
		{
			name: "ErrRateLimited sentinel",
			err:  heygen.ErrRateLimited,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := heygen.IsRateLimited(tt.err); got != tt.want {
				t.Errorf("IsRateLimited() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "404 error",
			err:  &heygen.APIError{StatusCode: 404},
			want: true,
		},
		{
			name: "not_found code",
			err:  &heygen.APIError{Code: "not_found", StatusCode: 400},
			want: true,
		},
		{
			name: "500 error",
			err:  &heygen.APIError{StatusCode: 500},
			want: false,
		},
		{
			name: "non-API error",
			err:  errors.New("some error"),
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := heygen.IsNotFound(tt.err); got != tt.want {
				t.Errorf("IsNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSentinelErrors(t *testing.T) {
	if heygen.ErrNoAPIKey == nil {
		t.Error("ErrNoAPIKey should not be nil")
	}
	if heygen.ErrRateLimited == nil {
		t.Error("ErrRateLimited should not be nil")
	}
	if heygen.ErrUnauthorized == nil {
		t.Error("ErrUnauthorized should not be nil")
	}
}
