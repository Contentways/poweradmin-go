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
		writeEnvelope(t, w, http.StatusOK, map[string]any{"templates": []map[string]any{
			{"id": 1, "name": "Zone Admin", "descr": "Admin for zones"},
		}})
	})
	templates, _, err := client.PermissionTemplate.List(context.Background())
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(templates) != 1 || templates[0].Name != "Zone Admin" {
		t.Errorf("templates = %+v", templates)
	}
}

func TestPermissionTemplateCreate(t *testing.T) {
	// Create returns an empty body, so the client looks the template up by name.
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/permission-templates":
			writeEnvelope(t, w, http.StatusCreated, nil)
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/permission-templates":
			writeEnvelope(t, w, http.StatusOK, map[string]any{"templates": []map[string]any{
				{"id": 8, "name": "Editors", "descr": "Can edit"},
			}})
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
	})
	tpl, _, err := client.PermissionTemplate.Create(context.Background(), PermissionTemplateOpts{
		Name: "Editors", Descr: "Can edit",
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if tpl.ID != 8 || tpl.Name != "Editors" {
		t.Errorf("tpl = %+v", tpl)
	}
}

func TestPermissionTemplateUpdate(t *testing.T) {
	// Update returns an empty body, so the client reads the template back by ID.
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPut && r.URL.Path == "/api/v2/permission-templates/8":
			writeEnvelope(t, w, http.StatusOK, nil)
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/permission-templates/8":
			writeEnvelope(t, w, http.StatusOK, map[string]any{
				"template": map[string]any{
					"id": 8, "name": "Editors", "descr": "Updated",
					"permissions": []map[string]any{{"id": 10, "name": "zone_master_add", "descr": "..."}},
				},
			})
		default:
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
	})
	tpl, _, err := client.PermissionTemplate.Update(context.Background(), 8, PermissionTemplateOpts{
		Name: "Editors", Descr: "Updated",
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if tpl.Descr != "Updated" || len(tpl.Permissions) != 1 {
		t.Errorf("tpl = %+v", tpl)
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
