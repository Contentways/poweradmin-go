// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import (
	"context"
	"net/http"
	"strconv"
	"testing"
)

func TestUserGetByID(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v2/users/5" {
			t.Errorf("method/path = %s %s", r.Method, r.URL.Path)
		}
		writeEnvelope(t, w, http.StatusOK, map[string]any{
			"user": map[string]any{
				"user_id":  5,
				"username": "alice",
				"fullname": "Alice A.",
				"email":    "alice@example.com",
				"active":   true,
			},
		})
	})
	user, _, err := client.User.GetByID(context.Background(), 5)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if user.ID != 5 || user.Username != "alice" {
		t.Errorf("user = %+v", user)
	}
}

func TestUserList(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/users" {
			t.Errorf("path = %s", r.URL.Path)
		}
		writeEnvelopeWithPagination(t, w, http.StatusOK,
			map[string]any{"users": []map[string]any{
				{"user_id": 1, "username": "alice"},
				{"user_id": 2, "username": "bob"},
			}},
			map[string]int{"current_page": 1, "per_page": 100, "total": 2, "last_page": 1},
		)
	})
	users, _, err := client.User.List(context.Background(), ListOpts{})
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("len(users) = %d, want 2", len(users))
	}
}

func TestUserAll(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page == 0 {
			page = 1
		}
		writeEnvelopeWithPagination(t, w, http.StatusOK,
			map[string]any{"users": []map[string]any{
				{"user_id": page, "username": "user" + strconv.Itoa(page)},
			}},
			map[string]int{"current_page": page, "per_page": 1, "total": 2, "last_page": 2},
		)
	})
	users, err := client.User.All(context.Background())
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(users) != 2 {
		t.Errorf("len(users) = %d, want 2", len(users))
	}
}

func TestUserCreate(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v2/users" {
			t.Errorf("method/path = %s %s", r.Method, r.URL.Path)
		}
		writeEnvelope(t, w, http.StatusCreated, map[string]any{"user_id": 42})
	})
	id, _, err := client.User.Create(context.Background(), UserCreateOpts{
		Username: "charlie", Password: "secret", Email: "charlie@example.com", Active: true,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if id != 42 {
		t.Errorf("id = %d, want 42", id)
	}
}

func TestUserUpdate(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut || r.URL.Path != "/api/v2/users/5" {
			t.Errorf("method/path = %s %s", r.Method, r.URL.Path)
		}
		writeEnvelope(t, w, http.StatusOK, map[string]any{
			"user": map[string]any{
				"user_id":  5,
				"username": "alice",
				"email":    "newalice@example.com",
			},
		})
	})
	active := true
	user, _, err := client.User.Update(context.Background(), 5, UserUpdateOpts{
		Username: "alice", Email: "newalice@example.com", Active: &active,
	})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if user.Email != "newalice@example.com" {
		t.Errorf("user = %+v", user)
	}
}

func TestUserDelete(t *testing.T) {
	deleted := false
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || r.URL.Path != "/api/v2/users/5" {
			t.Errorf("method/path = %s %s", r.Method, r.URL.Path)
		}
		deleted = true
		writeEnvelope(t, w, http.StatusOK, nil)
	})
	_, err := client.User.Delete(context.Background(), 5)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if !deleted {
		t.Error("DELETE was not called")
	}
}

func TestUserSetPermissionTemplate(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch || r.URL.Path != "/api/v2/users/5" {
			t.Errorf("method/path = %s %s", r.Method, r.URL.Path)
		}
		writeEnvelope(t, w, http.StatusOK, nil)
	})
	_, err := client.User.SetPermissionTemplate(context.Background(), 5, 3)
	if err != nil {
		t.Fatalf("SetPermissionTemplate: %v", err)
	}
}
