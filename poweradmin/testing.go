// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

// NewTestClient creates a [Client] with injected mock clients for testing.
// Only the provided clients are set — remaining fields are nil.
// Use this in unit tests to avoid real HTTP calls.
func NewTestClient(zone IZoneClient, record IRecordClient) *Client {
	return &Client{
		Zone:   zone,
		Record: record,
	}
}
