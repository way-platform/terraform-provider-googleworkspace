package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDrivePermission_Create(t *testing.T) {
	var createBody []byte

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		// Create permission
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/files/") && strings.Contains(r.URL.Path, "/permissions"):
			body, _ := io.ReadAll(r.Body)
			createBody = body

			// Verify supportsAllDrives query param
			if r.URL.Query().Get("supportsAllDrives") != "true" {
				t.Error("expected supportsAllDrives=true query parameter")
			}

			jsonResponse(w, 200, map[string]any{
				"kind": "drive#permission",
				"id":   "perm-001",
			})

		// Read permission
		case r.Method == "GET" && strings.Contains(r.URL.Path, "/permissions/perm-001"):
			jsonResponse(w, 200, map[string]any{
				"kind":         "drive#permission",
				"id":           "perm-001",
				"emailAddress": "team@way.cloud",
				"role":         "writer",
				"type":         "group",
			})

		// Delete permission
		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/permissions/perm-001"):
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
resource "googleworkspace_drive_permission" "test" {
  file_id                = "drive-abc"
  email_address          = "team@way.cloud"
  role                   = "writer"
  type                   = "group"
  use_domain_admin_access = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_drive_permission.test", "permission_id", "perm-001"),
					resource.TestCheckResourceAttr("googleworkspace_drive_permission.test", "id", "drive-abc/perm-001"),
					resource.TestCheckResourceAttr("googleworkspace_drive_permission.test", "role", "writer"),
					resource.TestCheckResourceAttr("googleworkspace_drive_permission.test", "type", "group"),
				),
			},
		},
	})

	// Verify create request body
	var reqBody map[string]any
	if err := json.Unmarshal(createBody, &reqBody); err != nil {
		t.Fatalf("failed to parse create request body: %v", err)
	}
	if reqBody["role"] != "writer" {
		t.Errorf("expected role=writer, got %v", reqBody["role"])
	}
	if reqBody["type"] != "group" {
		t.Errorf("expected type=group, got %v", reqBody["type"])
	}
	if reqBody["emailAddress"] != "team@way.cloud" {
		t.Errorf("expected emailAddress=team@way.cloud, got %v", reqBody["emailAddress"])
	}
}

func TestAccDrivePermission_UpdateRole(t *testing.T) {
	var updateBody []byte
	step := 0

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/permissions"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#permission",
				"id":   "perm-002",
			})

		case r.Method == "PATCH" && strings.Contains(r.URL.Path, "/permissions/perm-002"):
			body, _ := io.ReadAll(r.Body)
			updateBody = body
			step++
			jsonResponse(w, 200, map[string]any{
				"kind":         "drive#permission",
				"id":           "perm-002",
				"emailAddress": "editor@way.cloud",
				"role":         "reader",
				"type":         "user",
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/permissions/perm-002"):
			role := "writer"
			if step > 0 {
				role = "reader"
			}
			jsonResponse(w, 200, map[string]any{
				"kind":         "drive#permission",
				"id":           "perm-002",
				"emailAddress": "editor@way.cloud",
				"role":         role,
				"type":         "user",
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/permissions/perm-002"):
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
resource "googleworkspace_drive_permission" "test" {
  file_id                = "drive-xyz"
  email_address          = "editor@way.cloud"
  role                   = "writer"
  type                   = "user"
  use_domain_admin_access = true
}
`,
			},
			{
				Config: testProviderConfig + `
resource "googleworkspace_drive_permission" "test" {
  file_id                = "drive-xyz"
  email_address          = "editor@way.cloud"
  role                   = "reader"
  type                   = "user"
  use_domain_admin_access = true
}
`,
				Check: resource.TestCheckResourceAttr("googleworkspace_drive_permission.test", "role", "reader"),
			},
		},
	})

	// Verify only role was sent in update body
	var reqBody map[string]any
	if err := json.Unmarshal(updateBody, &reqBody); err != nil {
		t.Fatalf("failed to parse update request body: %v", err)
	}
	if reqBody["role"] != "reader" {
		t.Errorf("expected role=reader in update, got %v", reqBody["role"])
	}
}

func TestAccDrivePermission_ReadNotFound(t *testing.T) {
	readCount := 0

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/permissions"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#permission",
				"id":   "perm-gone",
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/permissions/perm-gone"):
			readCount++
			if readCount > 1 {
				// Second read: simulate external deletion
				jsonResponse(w, 404, map[string]any{
					"error": map[string]any{
						"code":    404,
						"message": "Permission not found",
					},
				})
			} else {
				jsonResponse(w, 200, map[string]any{
					"kind":         "drive#permission",
					"id":           "perm-gone",
					"emailAddress": "gone@way.cloud",
					"role":         "writer",
					"type":         "user",
				})
			}

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/permissions/perm-gone"):
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
resource "googleworkspace_drive_permission" "test" {
  file_id                = "drive-xyz"
  email_address          = "gone@way.cloud"
  role                   = "writer"
  type                   = "user"
  use_domain_admin_access = true
}
`,
			},
			// The second step triggers a read that returns 404 (external deletion).
			// The framework removes the resource from state and plans recreation.
			{
				Config: testProviderConfig + `
resource "googleworkspace_drive_permission" "test" {
  file_id                = "drive-xyz"
  email_address          = "gone@way.cloud"
  role                   = "writer"
  type                   = "user"
  use_domain_admin_access = true
}
`,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
