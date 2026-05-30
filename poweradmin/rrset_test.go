// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestRRSetGet(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/zones/5/rrsets/www/A" {
			t.Errorf("path = %s, want /api/v2/zones/5/rrsets/www/A", r.URL.Path)
		}
		// Single RRSet GET: the RRSet is wrapped under data.rrset.
		writeEnvelope(t, w, http.StatusOK, map[string]any{
			"rrset": map[string]any{
				"name": "www",
				"type": "A",
				"ttl":  3600,
				"records": []map[string]any{
					{"content": "1.2.3.4", "disabled": false},
					{"content": "1.2.3.5", "disabled": true},
				},
			},
		})
	})
	rr, _, err := client.RRSet.Get(context.Background(), 5, "www", "A")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rr.Name != "www" || rr.Type != "A" || rr.TTL != 3600 {
		t.Errorf("RRSet = %+v", rr)
	}
	if len(rr.Records) != 2 || rr.Records[0].Content != "1.2.3.4" || !rr.Records[1].Disabled {
		t.Errorf("Records = %+v", rr.Records)
	}
}

func TestRRSetList(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/zones/5/rrsets" {
			t.Errorf("path = %s", r.URL.Path)
		}
		writeEnvelope(t, w, http.StatusOK, map[string]any{"rrsets": []map[string]any{
			{
				"name":    "www.example.com",
				"type":    "A",
				"ttl":     3600,
				"records": []map[string]any{{"content": "1.2.3.4"}},
			},
		}})
	})
	rrsets, _, err := client.RRSet.List(context.Background(), 5)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(rrsets) != 1 || rrsets[0].Name != "www.example.com" {
		t.Errorf("rrsets = %+v", rrsets)
	}
}

func TestRRSetUpsert(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s, want PUT", r.Method)
		}
		if r.URL.Path != "/api/v2/zones/5/rrsets" {
			t.Errorf("path = %s, want /api/v2/zones/5/rrsets", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		_ = json.Unmarshal(body, &got)
		if got["name"] != "www" || got["type"] != "A" {
			t.Errorf("body = %v", got)
		}
		records, ok := got["records"].([]any)
		if !ok || len(records) != 2 {
			t.Errorf("records = %v", got["records"])
		}
		writeEnvelope(t, w, http.StatusOK, nil)
	})
	_, err := client.RRSet.Upsert(context.Background(), 5, RRSetUpsertOpts{
		Name: "www",
		Type: "A",
		TTL:  3600,
		Records: []RRSetRecord{
			{Content: "1.2.3.4"},
			{Content: "1.2.3.5", Disabled: true},
		},
	})
	if err != nil {
		t.Fatalf("Upsert: %v", err)
	}
}

func TestRRSetDelete(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/api/v2/zones/5/rrsets/mail/MX" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if _, err := client.RRSet.Delete(context.Background(), 5, "mail", "MX"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}
