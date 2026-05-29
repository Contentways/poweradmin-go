// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import (
	"errors"
	"fmt"
	"testing"
)

func TestAPIErrorMessage(t *testing.T) {
	tests := []struct {
		name string
		err  *APIError
		want string
	}{
		{
			name: "without details",
			err:  &APIError{StatusCode: 404, Message: "not found"},
			want: "poweradmin: HTTP 404: not found",
		},
		{
			name: "with details",
			err:  &APIError{StatusCode: 422, Message: "invalid", Details: "bad name"},
			want: "poweradmin: HTTP 422: invalid (bad name)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
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
		{"nil", nil, false},
		{"plain 404", &APIError{StatusCode: 404}, true},
		{"plain 500", &APIError{StatusCode: 500}, false},
		{"wrapped 404", fmt.Errorf("wrap: %w", &APIError{StatusCode: 404}), true},
		{"non-api error", errors.New("oops"), false},
		{"string fallback HTTP 404", errors.New("got HTTP 404 oops"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotFound(tt.err); got != tt.want {
				t.Errorf("IsNotFound(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}
