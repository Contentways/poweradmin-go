// Copyright (c) 2026 Contentways
// SPDX-License-Identifier: MIT
package poweradmin

import (
	"context"
	"net/http"
	"testing"
)

func TestGroupMembers(t *testing.T) {
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/groups/3/members" {
			t.Errorf("path = %s", r.URL.Path)
		}
		writeEnvelope(t, w, http.StatusOK, []map[string]any{
			{"user_id": 1, "username": "alice", "fullname": "Alice A."},
			{"user_id": 2, "username": "bob"},
		})
	})
	members, _, err := client.Group.Members(context.Background(), 3)
	if err != nil {
		t.Fatalf("Members: %v", err)
	}
	if len(members) != 2 || members[0].Username != "alice" {
		t.Errorf("members = %+v", members)
	}
}

func TestGroupAddRemoveMember(t *testing.T) {
	var sawAdd, sawRemove bool
	client, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/groups/3/members":
			sawAdd = true
			writeEnvelope(t, w, http.StatusOK, nil)
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v2/groups/3/members/9":
			sawRemove = true
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	})
	if _, err := client.Group.AddMember(context.Background(), 3, 9); err != nil {
		t.Fatalf("AddMember: %v", err)
	}
	if _, err := client.Group.RemoveMember(context.Background(), 3, 9); err != nil {
		t.Fatalf("RemoveMember: %v", err)
	}
	if !sawAdd || !sawRemove {
		t.Errorf("sawAdd=%v sawRemove=%v", sawAdd, sawRemove)
	}
}
