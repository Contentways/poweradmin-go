// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import (
	"context"
	"net/http"
	"testing"
)

func TestZoneTemplateList(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/zone-templates" {
			t.Errorf("path = %s", r.URL.Path)
		}
		writeEnvelope(t, w, http.StatusOK, []map[string]any{
			{"id": 1, "name": "Default", "description": "Default tpl", "owner": 1, "is_global": false, "zones_linked": 3},
		})
	})
	templates, _, err := client.ZoneTemplate.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(templates) != 1 || templates[0].Name != "Default" || templates[0].ZonesLinked != 3 {
		t.Errorf("templates = %+v", templates)
	}
}

func TestZoneTemplateGetByID(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/zone-templates/7" {
			t.Errorf("path = %s", r.URL.Path)
		}
		writeEnvelope(t, w, http.StatusOK, map[string]any{
			"template": map[string]any{"id": 7, "name": "X", "description": "d"},
		})
	})
	tpl, _, err := client.ZoneTemplate.GetByID(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if tpl.ID != 7 || tpl.Name != "X" {
		t.Errorf("tpl = %+v", tpl)
	}
}

func TestZoneTemplateRecords(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/zone-templates/1/records" {
			t.Errorf("path = %s", r.URL.Path)
		}
		writeEnvelope(t, w, http.StatusOK, []map[string]any{
			{"id": 1, "name": "[ZONE]", "type": "SOA", "content": "[NS1] [HOSTMASTER] [SERIAL] 28800 7200 604800 86400", "ttl": 86400},
		})
	})
	records, _, err := client.ZoneTemplate.Records(context.Background(), 1)
	if err != nil {
		t.Fatalf("Records: %v", err)
	}
	if len(records) != 1 || records[0].Type != "SOA" {
		t.Errorf("records = %+v", records)
	}
}
