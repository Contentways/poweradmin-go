// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import (
	"errors"
	"fmt"
	"strings"
)

// APIError represents an error returned by the Poweradmin API.
type APIError struct {
	StatusCode int
	Message    string
	Details    string
}

func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("poweradmin: HTTP %d: %s (%s)", e.StatusCode, e.Message, e.Details)
	}
	return fmt.Sprintf("poweradmin: HTTP %d: %s", e.StatusCode, e.Message)
}

// IsNotFound reports whether err is a 404 Not Found API error.
func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode == 404
	}
	// Fallback for string-wrapped errors.
	return strings.Contains(err.Error(), "HTTP 404")
}
