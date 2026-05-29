// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestRetryOn5xxThenSuccess(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&attempts, 1)
		if n < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		writeEnvelope(t, w, http.StatusOK, map[string]any{"zones": []map[string]any{}})
	}))
	t.Cleanup(srv.Close)

	c, err := NewClient(
		WithBaseURL(srv.URL),
		WithAPIKey("k"),
		WithRetry(5),
		WithRetryBackoff(func(int) time.Duration { return time.Millisecond }),
	)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	if _, _, err := c.Zone.List(context.Background(), ListOpts{}); err != nil {
		t.Fatalf("List: %v", err)
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Errorf("attempts = %d, want 3", got)
	}
}

func TestRetryExhausted(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		w.WriteHeader(http.StatusBadGateway)
	}))
	t.Cleanup(srv.Close)

	c, _ := NewClient(
		WithBaseURL(srv.URL),
		WithAPIKey("k"),
		WithRetry(3),
		WithRetryBackoff(func(int) time.Duration { return time.Millisecond }),
	)
	_, _, err := c.Zone.List(context.Background(), ListOpts{})
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	if got := atomic.LoadInt32(&attempts); got != 3 {
		t.Errorf("attempts = %d, want 3", got)
	}
}

func TestRetryNotOn4xx(t *testing.T) {
	var attempts int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&attempts, 1)
		writeError(t, w, http.StatusBadRequest, "bad")
	}))
	t.Cleanup(srv.Close)

	c, _ := NewClient(
		WithBaseURL(srv.URL),
		WithAPIKey("k"),
		WithRetry(5),
		WithRetryBackoff(func(int) time.Duration { return time.Millisecond }),
	)
	_, _, err := c.Zone.List(context.Background(), ListOpts{})
	if err == nil {
		t.Fatal("expected error")
	}
	if got := atomic.LoadInt32(&attempts); got != 1 {
		t.Errorf("attempts = %d, want 1 (4xx must not retry)", got)
	}
}

func TestRetryRespectsContextCancellation(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	t.Cleanup(srv.Close)

	c, _ := NewClient(
		WithBaseURL(srv.URL),
		WithAPIKey("k"),
		WithRetry(5),
		WithRetryBackoff(func(int) time.Duration { return 100 * time.Millisecond }),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, _, err := c.Zone.List(ctx, ListOpts{})
	if err == nil {
		t.Fatal("expected error from context timeout")
	}
}

func TestDefaultBackoffShape(t *testing.T) {
	// Just sanity-check that DefaultBackoff produces strictly bounded values.
	for i := range 10 {
		d := DefaultBackoff(i)
		if d <= 0 || d > 11*time.Second {
			t.Errorf("attempt %d: backoff %s out of range", i, d)
		}
	}
}

func TestDebugWriter(t *testing.T) {
	var buf bytes.Buffer
	client, _ := newTestClientWithOpts(t,
		func(w http.ResponseWriter, r *http.Request) {
			writeEnvelope(t, w, http.StatusOK, map[string]any{"zones": []map[string]any{}})
		},
		WithDebugWriter(&buf),
	)
	if _, _, err := client.Zone.List(context.Background(), ListOpts{}); err != nil {
		t.Fatalf("List: %v", err)
	}
	got := buf.String()
	if !strings.Contains(got, "[poweradmin]") || !strings.Contains(got, "GET") || !strings.Contains(got, "200") {
		t.Errorf("debug output missing expected fields: %q", got)
	}
}
