// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import "testing"

func TestListOptsValues(t *testing.T) {
	tests := []struct {
		name string
		opts ListOpts
		want string
	}{
		{"zero value", ListOpts{}, ""},
		{"page only", ListOpts{Page: 2}, "page=2"},
		{"per_page only", ListOpts{PerPage: 50}, "per_page=50"},
		{"both", ListOpts{Page: 3, PerPage: 25}, "page=3&per_page=25"},
		{"negative ignored", ListOpts{Page: -1, PerPage: 0}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.opts.values().Encode(); got != tt.want {
				t.Errorf("values() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAppendQuery(t *testing.T) {
	tests := []struct {
		name string
		path string
		opts ListOpts
		want string
	}{
		{"no opts", "zones", ListOpts{}, "zones"},
		{"with opts", "zones", ListOpts{Page: 2}, "zones?page=2"},
		{"existing query", "zones?filter=x", ListOpts{Page: 2}, "zones?filter=x&page=2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := appendQuery(tt.path, tt.opts.values()); got != tt.want {
				t.Errorf("appendQuery = %q, want %q", got, tt.want)
			}
		})
	}
}
