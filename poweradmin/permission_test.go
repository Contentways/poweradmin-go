// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import (
	"context"
	"net/http"
	"strconv"
	"testing"
)

func TestPermissionGetByID(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v2/permissions/7" {
			t.Errorf("method/path = %s %s", r.Method, r.URL.Path)
		}
		writeEnvelope(t, w, http.StatusOK, map[string]any{
			"permission": map[string]any{
				"id":    7,
				"name":  "zone_master_add",
				"descr": "Add new master zones",
			},
		})
	})
	perm, _, err := client.Permission.GetByID(context.Background(), 7)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if perm.ID != 7 || perm.Name != "zone_master_add" {
		t.Errorf("perm = %+v", perm)
	}
}

func TestPermissionList(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/permissions" {
			t.Errorf("path = %s", r.URL.Path)
		}
		writeEnvelopeWithPagination(t, w, http.StatusOK,
			map[string]any{"permissions": []map[string]any{
				{"id": 1, "name": "zone_master_add", "descr": "Add master zones"},
				{"id": 2, "name": "zone_slave_add", "descr": "Add slave zones"},
			}},
			map[string]int{"current_page": 1, "per_page": 100, "total": 2, "last_page": 1},
		)
	})
	perms, _, err := client.Permission.List(context.Background(), ListOpts{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(perms) != 2 {
		t.Errorf("len(perms) = %d, want 2", len(perms))
	}
}

func TestPermissionAll(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		writeEnvelopeWithPagination(t, w, http.StatusOK,
			map[string]any{"permissions": []map[string]any{
				{"id": page, "name": "perm" + strconv.Itoa(page), "descr": ""},
			}},
			map[string]int{"current_page": page, "per_page": 1, "total": 2, "last_page": 2},
		)
	})
	perms, err := client.Permission.All(context.Background())
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(perms) != 2 {
		t.Errorf("len(perms) = %d, want 2", len(perms))
	}
}
