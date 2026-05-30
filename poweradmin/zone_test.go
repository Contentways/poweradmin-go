// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
)

func TestZoneGetByID(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("method = %s, want GET", r.Method)
		}
		if r.URL.Path != "/api/v2/zones/5" {
			t.Errorf("path = %s, want /api/v2/zones/5", r.URL.Path)
		}
		// Single-zone GET wraps the object under data.zone.
		writeEnvelope(t, w, http.StatusOK, map[string]any{
			"zone": map[string]any{
				"id":   5,
				"name": "example.com",
				"type": "MASTER",
			},
		})
	})
	z, resp, err := client.Zone.GetByID(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if resp == nil || resp.StatusCode != http.StatusOK {
		t.Errorf("response status = %v", resp)
	}
	if z.ID != 5 || z.Name != "example.com" || z.Type != ZoneTypeMaster {
		t.Errorf("zone = %+v", z)
	}
}

func TestZoneGetByIDNotFound(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeError(t, w, http.StatusNotFound, "no such zone")
	})
	_, _, err := client.Zone.GetByID(context.Background(), 99)
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound = false, want true (err=%v)", err)
	}
}

func TestZoneListSinglePage(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/zones" {
			t.Errorf("path = %s", r.URL.Path)
		}
		// List: data is a flat array, pagination is at envelope level.
		writeEnvelopeWithPagination(t, w, http.StatusOK,
			map[string]any{
				"zones": []map[string]any{
					{"id": 1, "name": "a.com", "type": "MASTER"},
					{"id": 2, "name": "b.com", "type": "SLAVE"},
				},
			},
			map[string]int{"current_page": 1, "per_page": 100, "total": 2, "last_page": 1},
		)
	})
	zones, resp, err := client.Zone.List(context.Background(), ListOpts{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(zones) != 2 {
		t.Fatalf("len(zones) = %d, want 2", len(zones))
	}
	if zones[0].Name != "a.com" || zones[1].Type != ZoneTypeSlave {
		t.Errorf("zones = %+v", zones)
	}
	if resp.Meta.Pagination == nil {
		t.Fatal("Meta.Pagination is nil")
	}
	if resp.Meta.Pagination.Total != 2 || resp.Meta.Pagination.LastPage != 1 {
		t.Errorf("Pagination = %+v", resp.Meta.Pagination)
	}
}

func TestZoneAllIteratesPages(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		writeEnvelopeWithPagination(t, w, http.StatusOK,
			map[string]any{
				"zones": []map[string]any{
					{"id": page, "name": "z" + strconv.Itoa(page) + ".com", "type": "MASTER"},
				},
			},
			map[string]int{"current_page": page, "per_page": 1, "total": 2, "last_page": 2},
		)
	})
	zones, err := client.Zone.All(context.Background())
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(zones) != 2 {
		t.Fatalf("len(zones) = %d, want 2", len(zones))
	}
	if zones[0].ID != 1 || zones[1].ID != 2 {
		t.Errorf("ids = [%d,%d], want [1,2]", zones[0].ID, zones[1].ID)
	}
}

func TestZoneCreate(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/api/v2/zones" {
			t.Errorf("path = %s", r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatalf("body: %v", err)
		}
		if got["name"] != "new.com" || got["type"] != "MASTER" {
			t.Errorf("body = %v", got)
		}
		writeEnvelope(t, w, http.StatusCreated, map[string]any{"zone_id": 42})
	})
	id, _, err := client.Zone.Create(context.Background(), ZoneCreateOpts{
		Name: "new.com",
		Type: ZoneTypeMaster,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if id != 42 {
		t.Errorf("id = %d, want 42", id)
	}
}

func TestZoneUpdate(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		_ = json.Unmarshal(body, &got)
		if got["description"] != "updated" {
			t.Errorf("description = %v, want 'updated'", got["description"])
		}
		if _, present := got["type"]; present {
			t.Errorf("type should be omitted in update")
		}
		writeEnvelope(t, w, http.StatusOK, map[string]any{
			"zone": map[string]any{"id": 5, "name": "x.com", "type": "MASTER", "description": "updated"},
		})
	})
	desc := "updated"
	z, _, err := client.Zone.Update(context.Background(), 5, ZoneUpdateOpts{Description: &desc})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if z.Description != "updated" {
		t.Errorf("description = %q", z.Description)
	}
}

func TestZoneDelete(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/api/v2/zones/7" {
			t.Errorf("path = %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	if _, err := client.Zone.Delete(context.Background(), 7); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

func TestZoneGetByName(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeEnvelopeWithPagination(t, w, http.StatusOK,
			map[string]any{
				"zones": []map[string]any{
					{"id": 1, "name": "a.com", "type": "MASTER"},
					{"id": 2, "name": "target.com", "type": "MASTER"},
				},
			},
			map[string]int{"current_page": 1, "per_page": 100, "total": 2, "last_page": 1},
		)
	})
	z, _, err := client.Zone.GetByName(context.Background(), "target.com")
	if err != nil {
		t.Fatalf("GetByName: %v", err)
	}
	if z.ID != 2 {
		t.Errorf("ID = %d, want 2", z.ID)
	}
}

func TestZoneGetByNameNotFound(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		writeEnvelopeWithPagination(t, w, http.StatusOK,
			map[string]any{"zones": []map[string]any{}},
			map[string]int{"current_page": 1, "per_page": 100, "total": 0, "last_page": 1},
		)
	})
	_, _, err := client.Zone.GetByName(context.Background(), "missing.com")
	if err == nil {
		t.Fatal("expected error")
	}
	if !IsNotFound(err) {
		t.Errorf("IsNotFound = false, want true (err=%v)", err)
	}
}

func TestZoneOwners(t *testing.T) {
	var sawAdd, sawDelete bool
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/zones/5/owners":
			writeEnvelope(t, w, http.StatusOK, map[string]any{"owners": []map[string]any{
				{"user_id": 1, "username": "alice", "fullname": "Alice A."},
			}})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/zones/5/owners":
			sawAdd = true
			body, _ := io.ReadAll(r.Body)
			var got map[string]any
			_ = json.Unmarshal(body, &got)
			if got["user_id"] == nil {
				t.Errorf("expected user_id in body, got %v", got)
			}
			writeEnvelope(t, w, http.StatusOK, nil)
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v2/zones/5/owners/7":
			sawDelete = true
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})
	owners, _, err := client.Zone.Owners(context.Background(), 5)
	if err != nil {
		t.Fatalf("Owners: %v", err)
	}
	if len(owners) != 1 || owners[0].Username != "alice" {
		t.Errorf("owners = %+v", owners)
	}
	if _, err := client.Zone.AddOwner(context.Background(), 5, 7); err != nil {
		t.Fatalf("AddOwner: %v", err)
	}
	if _, err := client.Zone.RemoveOwner(context.Background(), 5, 7); err != nil {
		t.Fatalf("RemoveOwner: %v", err)
	}
	if !sawAdd || !sawDelete {
		t.Errorf("sawAdd=%v sawDelete=%v", sawAdd, sawDelete)
	}
}

func TestZoneAddOwners(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v2/zones/5/owners" {
			t.Errorf("method/path = %s %s", r.Method, r.URL.Path)
		}
		body, _ := io.ReadAll(r.Body)
		var got map[string]any
		_ = json.Unmarshal(body, &got)
		ids, ok := got["user_ids"]
		if !ok {
			t.Errorf("expected user_ids in body, got %v", got)
		}
		// JSON numbers decode as float64
		if ids.([]any)[0].(float64) != 3 || ids.([]any)[1].(float64) != 4 {
			t.Errorf("user_ids = %v", ids)
		}
		writeEnvelope(t, w, http.StatusOK, nil)
	})
	if _, err := client.Zone.AddOwners(context.Background(), 5, []int{3, 4}); err != nil {
		t.Fatalf("AddOwners: %v", err)
	}
}
