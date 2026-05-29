// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import (
	"math/rand/v2"
	"net/http"
	"time"
)

// Backoff returns the delay before the next retry attempt.
// attempt is the zero-based index of the *next* attempt (0 = first retry).
type Backoff func(attempt int) time.Duration

// retryConfig holds retry behavior; nil means "no retries".
type retryConfig struct {
	maxAttempts int
	backoff     Backoff
}

// DefaultBackoff is exponential with ±25% jitter, hard-capped at 10s.
// attempt 0 → ~100ms, 1 → ~200ms, 2 → ~400ms, 3 → ~800ms, ..., ≤ 10s.
func DefaultBackoff(attempt int) time.Duration {
	const (
		base       = 100 * time.Millisecond
		maxBackoff = 10 * time.Second
	)
	shift := min(attempt,
		// base << 16 ≈ 109min — well beyond max
		16)
	d := base << shift
	jitterRange := int64(d / 2)
	if jitterRange > 0 {
		d += time.Duration(rand.Int64N(jitterRange)) - d/4
	}
	if d > maxBackoff {
		return maxBackoff
	}
	if d < base {
		return base
	}
	return d
}

// shouldRetryStatus reports whether an HTTP status code is retryable.
// 429 (Too Many Requests) and 5xx are retried.
func shouldRetryStatus(code int) bool {
	return code == http.StatusTooManyRequests || (code >= 500 && code < 600)
}
