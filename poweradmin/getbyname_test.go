// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import (
	"context"
	"net/http"
	"testing"
)

func TestUserGetByName(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeEnvelopeWithPagination(t, w, http.StatusOK,
			[]map[string]any{
				{"user_id": 1, "username": "alice"},
				{"user_id": 2, "username": "bob"},
			},
			map[string]int{"current_page": 1, "per_page": 100, "total": 2, "last_page": 1},
		)
	})
	u, _, err := client.User.GetByName(context.Background(), "bob")
	if err != nil {
		t.Fatalf("GetByName: %v", err)
	}
	if u.ID != 2 || u.Username != "bob" {
		t.Errorf("user = %+v", u)
	}
}

func TestUserGetByNameNotFound(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeEnvelopeWithPagination(t, w, http.StatusOK,
			[]map[string]any{{"user_id": 1, "username": "alice"}},
			map[string]int{"current_page": 1, "per_page": 100, "total": 1, "last_page": 1},
		)
	})
	_, _, err := client.User.GetByName(context.Background(), "missing")
	if !IsNotFound(err) {
		t.Errorf("IsNotFound = false, err = %v", err)
	}
}

func TestGroupGetByName(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeEnvelopeWithPagination(t, w, http.StatusOK,
			[]map[string]any{
				{"id": 1, "name": "ops"},
				{"id": 2, "name": "dev"},
			},
			map[string]int{"current_page": 1, "per_page": 100, "total": 2, "last_page": 1},
		)
	})
	g, _, err := client.Group.GetByName(context.Background(), "dev")
	if err != nil {
		t.Fatalf("GetByName: %v", err)
	}
	if g.ID != 2 {
		t.Errorf("group = %+v", g)
	}
}

func TestZoneTemplateGetByName(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeEnvelope(t, w, http.StatusOK, []map[string]any{
			{"id": 1, "name": "Default", "description": "x"},
			{"id": 2, "name": "Custom", "description": "y"},
		})
	})
	tpl, _, err := client.ZoneTemplate.GetByName(context.Background(), "Custom")
	if err != nil {
		t.Fatalf("GetByName: %v", err)
	}
	if tpl.ID != 2 {
		t.Errorf("template = %+v", tpl)
	}
}

func TestPermissionTemplateGetByName(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeEnvelope(t, w, http.StatusOK, []map[string]any{
			{"id": 1, "name": "Admin", "descr": "x"},
			{"id": 2, "name": "Viewer", "descr": "y"},
		})
	})
	tpl, _, err := client.PermissionTemplate.GetByName(context.Background(), "Viewer")
	if err != nil {
		t.Fatalf("GetByName: %v", err)
	}
	if tpl.ID != 2 {
		t.Errorf("template = %+v", tpl)
	}
}
