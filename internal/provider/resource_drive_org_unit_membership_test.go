package provider

import (
	"encoding/json"
	"fmt"
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

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/orgUnits/ou-123/memberships"):
			assertDriveOrgUnitMembershipListRequest(t, r)
			jsonResponse(w, 200, driveOrgUnitMembershipListResponse("ou-123", "drive-abc"))

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

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/orgUnits/ou-123/memberships"):
			mu.Lock()
			ou := currentOU
			mu.Unlock()
			assertDriveOrgUnitMembershipListRequest(t, r)
			if ou == "ou-123" {
				jsonResponse(w, 200, driveOrgUnitMembershipListResponse(ou, "drive-abc"))
				return
			}
			jsonResponse(w, 200, map[string]any{"orgMemberships": []any{}})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/orgUnits/ou-456/memberships"):
			mu.Lock()
			ou := currentOU
			mu.Unlock()
			assertDriveOrgUnitMembershipListRequest(t, r)
			if ou == "ou-456" {
				jsonResponse(w, 200, driveOrgUnitMembershipListResponse(ou, "drive-abc"))
				return
			}
			jsonResponse(w, 200, map[string]any{"orgMemberships": []any{}})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits"):
			jsonResponse(w, 200, map[string]any{
				"organizationUnits": []map[string]any{
					{"orgUnitId": "ou-123"},
					{"orgUnitId": "ou-456"},
				},
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

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/orgUnits/ou-123/memberships"):
			mu.Lock()
			c := created
			mu.Unlock()
			if c {
				jsonResponse(w, 200, map[string]any{"orgMemberships": []any{}})
				return
			}
			jsonResponse(w, 200, driveOrgUnitMembershipListResponse("ou-123", "drive-gone"))

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits"):
			jsonResponse(w, 200, map[string]any{
				"organizationUnits": []map[string]any{{"orgUnitId": "ou-123"}},
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

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/orgUnits/ou-789/memberships"):
			assertDriveOrgUnitMembershipListRequest(t, r)
			jsonResponse(w, 200, driveOrgUnitMembershipListResponse("ou-789", "drive-imp"))

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits"):
			jsonResponse(w, 200, map[string]any{
				"organizationUnits": []map[string]any{{"orgUnitId": "ou-789"}},
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

func TestAccDriveOrgUnitMembership_OrgUnitNotFound(t *testing.T) {
	var mu sync.Mutex
	created := false

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/orgUnits/-/memberships/shared_drive;drive-ou404:move"):
			mu.Lock()
			created = true
			mu.Unlock()
			jsonResponse(w, 200, map[string]any{
				"name": "operations/123",
				"done": true,
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/orgUnits/ou-deleted/memberships"):
			mu.Lock()
			c := created
			mu.Unlock()
			if c {
				jsonResponse(w, 404, map[string]any{
					"error": map[string]any{
						"code":    404,
						"message": "Org unit not found: ou-deleted",
						"status":  "NOT_FOUND",
					},
				})
				return
			}
			jsonResponse(w, 200, driveOrgUnitMembershipListResponse("ou-deleted", "drive-ou404"))

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
  drive_id    = "drive-ou404"
  org_unit_id = "ou-deleted"
}
`,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDriveOrgUnitMembership_Drift(t *testing.T) {
	var mu sync.Mutex
	actualOU := "ou-123"

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/orgUnits/-/memberships/shared_drive;drive-drift:move"):
			body, _ := io.ReadAll(r.Body)
			var reqBody map[string]any
			_ = json.Unmarshal(body, &reqBody)
			if dest, ok := reqBody["destinationOrgUnit"].(string); ok {
				mu.Lock()
				actualOU = strings.TrimPrefix(dest, "orgUnits/")
				mu.Unlock()
			}
			jsonResponse(w, 200, map[string]any{
				"name": "operations/123",
				"done": true,
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/orgUnits/ou-123/memberships"):
			mu.Lock()
			ou := actualOU
			mu.Unlock()
			assertDriveOrgUnitMembershipListRequest(t, r)
			if ou == "ou-123" {
				jsonResponse(w, 200, driveOrgUnitMembershipListResponse("ou-123", "drive-drift"))
				return
			}
			jsonResponse(w, 200, map[string]any{"orgMemberships": []any{}})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/orgUnits/ou-other/memberships"):
			mu.Lock()
			ou := actualOU
			mu.Unlock()
			assertDriveOrgUnitMembershipListRequest(t, r)
			if ou == "ou-other" {
				jsonResponse(w, 200, driveOrgUnitMembershipListResponse("ou-other", "drive-drift"))
				return
			}
			jsonResponse(w, 200, map[string]any{"orgMemberships": []any{}})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits"):
			jsonResponse(w, 200, map[string]any{
				"organizationUnits": []map[string]any{
					{"orgUnitId": "ou-123"},
					{"orgUnitId": "ou-other"},
				},
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
  drive_id    = "drive-drift"
  org_unit_id = "ou-123"
}
`,
				Check: resource.TestCheckResourceAttr("googleworkspace_drive_org_unit_membership.test", "org_unit_id", "ou-123"),
			},
			{
				PreConfig: func() {
					mu.Lock()
					actualOU = "ou-other"
					mu.Unlock()
				},
				Config: testProviderConfig + `
resource "googleworkspace_drive_org_unit_membership" "test" {
  drive_id    = "drive-drift"
  org_unit_id = "ou-123"
}
`,
				Check: resource.TestCheckResourceAttr("googleworkspace_drive_org_unit_membership.test", "org_unit_id", "ou-123"),
			},
		},
	})
}

func driveOrgUnitMembershipListResponse(orgUnitId, driveId string) map[string]any {
	return map[string]any{
		"orgMemberships": []map[string]any{
			{
				"name":      fmt.Sprintf("orgUnits/%s/memberships/shared_drive;%s", orgUnitId, driveId),
				"member":    fmt.Sprintf("//drive.googleapis.com/drives/%s", driveId),
				"memberUri": fmt.Sprintf("https://drive.googleapis.com/drive/v3/drives/%s", driveId),
				"type":      "SHARED_DRIVE",
			},
		},
	}
}

func assertDriveOrgUnitMembershipListRequest(t *testing.T, r *http.Request) {
	t.Helper()

	if got := r.URL.Query().Get("customer"); got != "customers/C00000000" {
		t.Errorf("expected customer=customers/C00000000, got %q", got)
	}
	if got := r.URL.Query().Get("filter"); got != "type == 'shared_drive'" {
		t.Errorf("expected shared drive filter, got %q", got)
	}
}
