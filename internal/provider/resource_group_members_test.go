package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupMembers_Create(t *testing.T) {
	var mu sync.Mutex
	var insertedMembers []map[string]any

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		// Insert member
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/groups/group-001/members"):
			body, _ := io.ReadAll(r.Body)
			var member map[string]any
			if err := json.Unmarshal(body, &member); err != nil {
				t.Errorf("failed to unmarshal request: %v", err)
				w.WriteHeader(500)
				return
			}
			mu.Lock()
			insertedMembers = append(insertedMembers, member)
			mu.Unlock()
			jsonResponse(w, 200, map[string]any{
				"kind":  "admin#directory#member",
				"email": member["email"],
				"role":  member["role"],
				"type":  member["type"],
				"id":    "member-" + member["email"].(string),
			})

		// List members
		case r.Method == "GET" && strings.Contains(r.URL.Path, "/groups/group-001/members"):
			mu.Lock()
			members := make([]map[string]any, 0, len(insertedMembers))
			for _, m := range insertedMembers {
				members = append(members, map[string]any{
					"kind":  "admin#directory#member",
					"email": m["email"],
					"role":  m["role"],
					"type":  m["type"],
					"id":    "member-" + m["email"].(string),
				})
			}
			mu.Unlock()
			jsonResponse(w, 200, map[string]any{
				"kind":    "admin#directory#members",
				"members": members,
			})

		// Delete member
		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/groups/group-001/members/"):
			w.WriteHeader(204)

		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(500)
		}
	}))
	setupTestClient(t, server)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig + `
resource "googleworkspace_group_members" "test" {
  group_id = "group-001"

  members {
    email = "alice@way.cloud"
    role  = "MEMBER"
  }

  members {
    email = "bob@way.cloud"
    role  = "OWNER"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_group_members.test", "id", "group-001"),
					resource.TestCheckResourceAttr("googleworkspace_group_members.test", "members.#", "2"),
				),
			},
		},
	})

	mu.Lock()
	defer mu.Unlock()
	if len(insertedMembers) != 2 {
		t.Errorf("expected 2 member inserts, got %d", len(insertedMembers))
	}
}

func TestAccGroupMembers_Update(t *testing.T) {
	var mu sync.Mutex
	currentMembers := map[string]map[string]any{}
	var deletedEmails []string
	var updatedEmails []string
	var insertedEmails []string

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		switch {
		// Insert member
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/groups/group-002/members"):
			body, _ := io.ReadAll(r.Body)
			var member map[string]any
			if err := json.Unmarshal(body, &member); err != nil {
				t.Errorf("failed to unmarshal request: %v", err)
				w.WriteHeader(500)
				return
			}
			email := member["email"].(string)
			currentMembers[email] = member
			insertedEmails = append(insertedEmails, email)
			jsonResponse(w, 200, map[string]any{
				"kind":  "admin#directory#member",
				"email": email,
				"role":  member["role"],
				"type":  member["type"],
				"id":    "member-" + email,
			})

		// List members
		case r.Method == "GET" && strings.Contains(r.URL.Path, "/groups/group-002/members") && !strings.Contains(r.URL.Path, "/members/"):
			members := make([]map[string]any, 0, len(currentMembers))
			for email, m := range currentMembers {
				members = append(members, map[string]any{
					"kind":  "admin#directory#member",
					"email": email,
					"role":  m["role"],
					"type":  m["type"],
					"id":    "member-" + email,
				})
			}
			jsonResponse(w, 200, map[string]any{
				"kind":    "admin#directory#members",
				"members": members,
			})

		// Update member
		case r.Method == "PUT" && strings.Contains(r.URL.Path, "/groups/group-002/members/"):
			parts := strings.Split(r.URL.Path, "/members/")
			email := parts[len(parts)-1]
			body, _ := io.ReadAll(r.Body)
			var member map[string]any
			if err := json.Unmarshal(body, &member); err != nil {
				t.Errorf("failed to unmarshal request: %v", err)
				w.WriteHeader(500)
				return
			}
			if m, ok := currentMembers[email]; ok {
				m["role"] = member["role"]
				currentMembers[email] = m
			}
			updatedEmails = append(updatedEmails, email)
			jsonResponse(w, 200, map[string]any{
				"kind":  "admin#directory#member",
				"email": email,
				"role":  member["role"],
				"type":  "USER",
				"id":    "member-" + email,
			})

		// Delete member
		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/groups/group-002/members/"):
			parts := strings.Split(r.URL.Path, "/members/")
			email := parts[len(parts)-1]
			delete(currentMembers, email)
			deletedEmails = append(deletedEmails, email)
			w.WriteHeader(204)

		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(500)
		}
	}))
	setupTestClient(t, server)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig + `
resource "googleworkspace_group_members" "test" {
  group_id = "group-002"

  members {
    email = "alice@way.cloud"
    role  = "MEMBER"
  }

  members {
    email = "bob@way.cloud"
    role  = "MEMBER"
  }
}
`,
			},
			// Step 2: remove bob, add charlie, change alice's role
			{
				Config: testProviderConfig + `
resource "googleworkspace_group_members" "test" {
  group_id = "group-002"

  members {
    email = "alice@way.cloud"
    role  = "OWNER"
  }

  members {
    email = "charlie@way.cloud"
    role  = "MEMBER"
  }
}
`,
				Check: resource.TestCheckResourceAttr("googleworkspace_group_members.test", "members.#", "2"),
			},
		},
	})

	mu.Lock()
	defer mu.Unlock()

	// Verify bob was deleted (during update step, before destroy also deletes everyone)
	bobDeleted := false
	for _, email := range deletedEmails {
		if email == "bob@way.cloud" {
			bobDeleted = true
		}
	}
	if !bobDeleted {
		t.Error("expected bob@way.cloud to be deleted")
	}

	// Verify alice's role was updated
	aliceUpdated := false
	for _, email := range updatedEmails {
		if email == "alice@way.cloud" {
			aliceUpdated = true
		}
	}
	if !aliceUpdated {
		t.Error("expected alice@way.cloud to be updated")
	}

	// Verify charlie was inserted at some point (it may have been deleted during destroy)
	charlieInserted := false
	for _, email := range insertedEmails {
		if email == "charlie@way.cloud" {
			charlieInserted = true
		}
	}
	if !charlieInserted {
		t.Error("expected charlie@way.cloud to be inserted")
	}
}
