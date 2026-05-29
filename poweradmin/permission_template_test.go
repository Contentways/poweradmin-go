// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import (
	"context"
	"net/http"
	"testing"
)

func TestPermissionTemplateList(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/permission-templates" {
			t.Errorf("path = %s", r.URL.Path)
		}
		writeEnvelope(t, w, http.StatusOK, []map[string]any{
			{"id": 1, "name": "Zone Admin", "descr": "Admin for zones"},
		})
	})
	templates, _, err := client.PermissionTemplate.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(templates) != 1 || templates[0].Name != "Zone Admin" {
		t.Errorf("templates = %+v", templates)
	}
}

func TestPermissionTemplateGetByID(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/permission-templates/3" {
			t.Errorf("path = %s", r.URL.Path)
		}
		writeEnvelope(t, w, http.StatusOK, map[string]any{
			"template": map[string]any{
				"id":    3,
				"name":  "Zone Admin",
				"descr": "Admin for zones",
				"permissions": []map[string]any{
					{"id": 10, "name": "zone_master_add", "descr": "..."},
				},
			},
		})
	})
	tpl, _, err := client.PermissionTemplate.GetByID(context.Background(), 3)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if tpl.ID != 3 || len(tpl.Permissions) != 1 || tpl.Permissions[0].Name != "zone_master_add" {
		t.Errorf("tpl = %+v", tpl)
	}
}
