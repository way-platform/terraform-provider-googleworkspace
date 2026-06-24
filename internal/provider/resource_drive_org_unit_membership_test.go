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

func TestAccDriveOrgUnitMembership_Create(t *testing.T) {
	var moveBody []byte

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/orgUnits/-/memberships/shared_drive;drive-abc:move"):
			body, _ := io.ReadAll(r.Body)
			moveBody = body
			jsonResponse(w, 200, map[string]any{
				"name": "operations/123",
				"done": true,
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/drives/drive-abc"):
			jsonResponse(w, 200, map[string]any{
				"kind":      "drive#drive",
				"id":        "drive-abc",
				"orgUnitId": "ou-123",
			})

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
resource "googleworkspace_drive_org_unit_membership" "test" {
  drive_id    = "drive-abc"
  org_unit_id = "ou-123"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_drive_org_unit_membership.test", "id", "drive-abc"),
					resource.TestCheckResourceAttr("googleworkspace_drive_org_unit_membership.test", "drive_id", "drive-abc"),
					resource.TestCheckResourceAttr("googleworkspace_drive_org_unit_membership.test", "org_unit_id", "ou-123"),
				),
			},
		},
	})

	var reqBody map[string]any
	if err := json.Unmarshal(moveBody, &reqBody); err != nil {
		t.Fatalf("failed to parse move request body: %v", err)
	}
	if reqBody["customer"] != "customers/C00000000" {
		t.Errorf("expected customer=customers/C00000000, got %v", reqBody["customer"])
	}
	if reqBody["destinationOrgUnit"] != "orgUnits/ou-123" {
		t.Errorf("expected destinationOrgUnit=orgUnits/ou-123, got %v", reqBody["destinationOrgUnit"])
	}
}

func TestAccDriveOrgUnitMembership_Update(t *testing.T) {
	var mu sync.Mutex
	currentOU := "ou-123"

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/orgUnits/-/memberships/shared_drive;drive-abc:move"):
			body, _ := io.ReadAll(r.Body)
			var reqBody map[string]any
			_ = json.Unmarshal(body, &reqBody)
			if dest, ok := reqBody["destinationOrgUnit"].(string); ok {
				mu.Lock()
				currentOU = strings.TrimPrefix(dest, "orgUnits/")
				mu.Unlock()
			}
			jsonResponse(w, 200, map[string]any{
				"name": "operations/123",
				"done": true,
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/drives/drive-abc"):
			mu.Lock()
			ou := currentOU
			mu.Unlock()
			jsonResponse(w, 200, map[string]any{
				"kind":      "drive#drive",
				"id":        "drive-abc",
				"orgUnitId": ou,
			})

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
resource "googleworkspace_drive_org_unit_membership" "test" {
  drive_id    = "drive-abc"
  org_unit_id = "ou-123"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_drive_org_unit_membership.test", "org_unit_id", "ou-123"),
					resource.TestCheckResourceAttr("googleworkspace_drive_org_unit_membership.test", "id", "drive-abc"),
				),
			},
			{
				Config: testProviderConfig + `
resource "googleworkspace_drive_org_unit_membership" "test" {
  drive_id    = "drive-abc"
  org_unit_id = "ou-456"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_drive_org_unit_membership.test", "org_unit_id", "ou-456"),
					resource.TestCheckResourceAttr("googleworkspace_drive_org_unit_membership.test", "id", "drive-abc"),
				),
			},
		},
	})
}

func TestAccDriveOrgUnitMembership_DriveNotFound(t *testing.T) {
	var mu sync.Mutex
	created := false

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/orgUnits/-/memberships/shared_drive;drive-gone:move"):
			mu.Lock()
			created = true
			mu.Unlock()
			jsonResponse(w, 200, map[string]any{
				"name": "operations/123",
				"done": true,
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/drives/drive-gone"):
			mu.Lock()
			c := created
			mu.Unlock()
			if c {
				jsonResponse(w, 404, map[string]any{
					"error": map[string]any{
						"code":    404,
						"message": "Shared drive not found: drive-gone",
						"status":  "NOT_FOUND",
					},
				})
				return
			}
			jsonResponse(w, 200, map[string]any{
				"kind":      "drive#drive",
				"id":        "drive-gone",
				"orgUnitId": "ou-123",
			})

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
resource "googleworkspace_drive_org_unit_membership" "test" {
  drive_id    = "drive-gone"
  org_unit_id = "ou-123"
}
`,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDriveOrgUnitMembership_Import(t *testing.T) {
	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/orgUnits/-/memberships/shared_drive;drive-imp:move"):
			jsonResponse(w, 200, map[string]any{
				"name": "operations/123",
				"done": true,
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/drives/drive-imp"):
			jsonResponse(w, 200, map[string]any{
				"kind":      "drive#drive",
				"id":        "drive-imp",
				"orgUnitId": "ou-789",
			})

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
resource "googleworkspace_drive_org_unit_membership" "test" {
  drive_id    = "drive-imp"
  org_unit_id = "ou-789"
}
`,
			},
			{
				ResourceName:      "googleworkspace_drive_org_unit_membership.test",
				ImportState:       true,
				ImportStateId:     "drive-imp",
				ImportStateVerify: true,
			},
		},
	})
}
