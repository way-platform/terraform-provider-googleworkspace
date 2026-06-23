package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDrive_CreateWithRestrictions(t *testing.T) {
	var createCalled atomic.Bool
	var updateBody []byte

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		// Create drive
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/drives"):
			createCalled.Store(true)
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-123",
				"name": "Engineering",
			})

		// Update drive (set restrictions after create)
		case r.Method == "PATCH" && strings.Contains(r.URL.Path, "/drives/drive-123"):
			body, _ := io.ReadAll(r.Body)
			updateBody = body
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-123",
				"name": "Engineering",
				"restrictions": map[string]any{
					"adminManagedRestrictions":                  true,
					"copyRequiresWriterPermission":              false,
					"domainUsersOnly":                           true,
					"driveMembersOnly":                          false,
					"sharingFoldersRequiresOrganizerPermission": false,
				},
			})

		// Read drive
		case r.Method == "GET" && strings.Contains(r.URL.Path, "/drives/drive-123"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-123",
				"name": "Engineering",
				"restrictions": map[string]any{
					"adminManagedRestrictions":                  true,
					"copyRequiresWriterPermission":              false,
					"domainUsersOnly":                           true,
					"driveMembersOnly":                          false,
					"sharingFoldersRequiresOrganizerPermission": false,
				},
			})

		// Delete drive
		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/drives/drive-123"):
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
resource "googleworkspace_drive" "test" {
  name                   = "Engineering"
  use_domain_admin_access = true

  restrictions {
    admin_managed_restrictions = true
    domain_users_only          = true
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_drive.test", "id", "drive-123"),
					resource.TestCheckResourceAttr("googleworkspace_drive.test", "name", "Engineering"),
					resource.TestCheckResourceAttr("googleworkspace_drive.test", "restrictions.admin_managed_restrictions", "true"),
					resource.TestCheckResourceAttr("googleworkspace_drive.test", "restrictions.domain_users_only", "true"),
					resource.TestCheckResourceAttr("googleworkspace_drive.test", "restrictions.drive_members_only", "false"),
				),
			},
		},
	})

	if !createCalled.Load() {
		t.Error("expected POST /drives to be called")
	}

	// Verify ForceSendFields: the update body should include false-valued restriction fields.
	var reqBody map[string]any
	if err := json.Unmarshal(updateBody, &reqBody); err != nil {
		t.Fatalf("failed to parse update request body: %v", err)
	}
	restrictions, ok := reqBody["restrictions"].(map[string]any)
	if !ok {
		t.Fatal("expected restrictions in update request body")
	}
	// These are explicitly set to false in the config (via defaults).
	// Without ForceSendFields they would be omitted due to omitempty.
	for _, field := range []string{"copyRequiresWriterPermission", "driveMembersOnly", "sharingFoldersRequiresOrganizerPermission"} {
		val, exists := restrictions[field]
		if !exists {
			t.Errorf("ForceSendFields bug: %q not present in request body (omitempty dropped it)", field)
		} else if val != false {
			t.Errorf("expected %q to be false, got %v", field, val)
		}
	}
}

func TestAccDrive_CreateWithRestrictionsRetriesWithDomainAdminAccess(t *testing.T) {
	origDelay := driveRestrictionRetryDelay
	driveRestrictionRetryDelay = func(int) time.Duration { return 0 }
	defer func() { driveRestrictionRetryDelay = origDelay }()

	var patchCount atomic.Int32
	var missingDomainAdminAccess atomic.Bool

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/drives"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-retry",
				"name": "Retry",
			})

		case r.Method == "PATCH" && strings.Contains(r.URL.Path, "/drives/drive-retry"):
			if r.URL.Query().Get("useDomainAdminAccess") != "true" {
				missingDomainAdminAccess.Store(true)
			}
			if patchCount.Add(1) == 1 {
				jsonResponse(w, 404, map[string]any{
					"error": map[string]any{
						"code":    404,
						"message": "Shared drive not found: drive-retry",
						"status":  "NOT_FOUND",
					},
				})
				return
			}
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-retry",
				"name": "Retry",
				"restrictions": map[string]any{
					"adminManagedRestrictions":                  true,
					"copyRequiresWriterPermission":              false,
					"domainUsersOnly":                           false,
					"driveMembersOnly":                          true,
					"sharingFoldersRequiresOrganizerPermission": false,
				},
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/drives/drive-retry"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-retry",
				"name": "Retry",
				"restrictions": map[string]any{
					"adminManagedRestrictions":                  true,
					"copyRequiresWriterPermission":              false,
					"domainUsersOnly":                           false,
					"driveMembersOnly":                          true,
					"sharingFoldersRequiresOrganizerPermission": false,
				},
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/drives/drive-retry"):
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
resource "googleworkspace_drive" "test" {
  name                    = "Retry"
  use_domain_admin_access = true

  restrictions {
    admin_managed_restrictions = true
    drive_members_only         = true
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_drive.test", "id", "drive-retry"),
					resource.TestCheckResourceAttr("googleworkspace_drive.test", "restrictions.admin_managed_restrictions", "true"),
					resource.TestCheckResourceAttr("googleworkspace_drive.test", "restrictions.drive_members_only", "true"),
				),
			},
		},
	})

	if patchCount.Load() != 2 {
		t.Fatalf("expected 2 restriction update attempts, got %d", patchCount.Load())
	}
	if missingDomainAdminAccess.Load() {
		t.Fatal("create-time restriction update should use domain admin access")
	}
}

func TestAccDrive_CreateWithRestrictionsExhaustsRetryAfterSavingState(t *testing.T) {
	origDelay := driveRestrictionRetryDelay
	driveRestrictionRetryDelay = func(int) time.Duration { return 0 }
	defer func() { driveRestrictionRetryDelay = origDelay }()

	var patchCount atomic.Int32
	var deleteCount atomic.Int32
	var missingDomainAdminAccess atomic.Bool

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/drives"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-partial",
				"name": "Partial",
			})

		case r.Method == "PATCH" && strings.Contains(r.URL.Path, "/drives/drive-partial"):
			if r.URL.Query().Get("useDomainAdminAccess") != "true" {
				missingDomainAdminAccess.Store(true)
			}
			patchCount.Add(1)
			jsonResponse(w, 404, map[string]any{
				"error": map[string]any{
					"code":    404,
					"message": "Shared drive not found: drive-partial",
					"status":  "NOT_FOUND",
				},
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/drives/drive-partial"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-partial",
				"name": "Partial",
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/drives/drive-partial"):
			deleteCount.Add(1)
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
resource "googleworkspace_drive" "test" {
  name                    = "Partial"
  use_domain_admin_access = true

  restrictions {
    drive_members_only = true
  }
}
`,
				ExpectError: regexp.MustCompile(`Unable to update drive restrictions`),
			},
		},
	})

	if patchCount.Load() != 6 {
		t.Fatalf("expected 6 restriction update attempts, got %d", patchCount.Load())
	}
	if deleteCount.Load() != 1 {
		t.Fatalf("expected cleanup delete for partial create, got %d", deleteCount.Load())
	}
	if missingDomainAdminAccess.Load() {
		t.Fatal("create-time restriction update should use domain admin access")
	}
}

func TestAccDrive_UpdateRestrictions(t *testing.T) {
	var updateCount atomic.Int32
	var lastUpdateBody []byte

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/drives"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-456",
				"name": "Sales",
			})

		case r.Method == "PATCH" && strings.Contains(r.URL.Path, "/drives/drive-456"):
			body, _ := io.ReadAll(r.Body)
			lastUpdateBody = body
			updateCount.Add(1)

			// Return whatever restrictions are in the request
			var req map[string]any
			if err := json.Unmarshal(body, &req); err != nil {
				t.Errorf("failed to unmarshal request: %v", err)
				w.WriteHeader(500)
				return
			}
			resp := map[string]any{
				"kind": "drive#drive",
				"id":   "drive-456",
				"name": "Sales",
			}
			if restrictions, ok := req["restrictions"]; ok {
				resp["restrictions"] = restrictions
			} else {
				resp["restrictions"] = map[string]any{
					"adminManagedRestrictions":                  false,
					"copyRequiresWriterPermission":              false,
					"domainUsersOnly":                           false,
					"driveMembersOnly":                          false,
					"sharingFoldersRequiresOrganizerPermission": false,
				}
			}
			jsonResponse(w, 200, resp)

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/drives/drive-456"):
			// Return current state based on update count
			restrictions := map[string]any{
				"adminManagedRestrictions":                  false,
				"copyRequiresWriterPermission":              false,
				"domainUsersOnly":                           false,
				"driveMembersOnly":                          false,
				"sharingFoldersRequiresOrganizerPermission": false,
			}
			if updateCount.Load() >= 2 {
				// After the second step's update
				restrictions["driveMembersOnly"] = true
				restrictions["domainUsersOnly"] = false
			} else if updateCount.Load() >= 1 {
				// After the first step's update
				restrictions["domainUsersOnly"] = true
			}
			jsonResponse(w, 200, map[string]any{
				"kind":         "drive#drive",
				"id":           "drive-456",
				"name":         "Sales",
				"restrictions": restrictions,
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/drives/drive-456"):
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
resource "googleworkspace_drive" "test" {
  name                   = "Sales"
  use_domain_admin_access = true

  restrictions {
    domain_users_only = true
  }
}
`,
				Check: resource.TestCheckResourceAttr("googleworkspace_drive.test", "restrictions.domain_users_only", "true"),
			},
			// Step 2: change restrictions (set domain_users_only=false, drive_members_only=true)
			{
				Config: testProviderConfig + `
resource "googleworkspace_drive" "test" {
  name                   = "Sales"
  use_domain_admin_access = true

  restrictions {
    domain_users_only  = false
    drive_members_only = true
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_drive.test", "restrictions.domain_users_only", "false"),
					resource.TestCheckResourceAttr("googleworkspace_drive.test", "restrictions.drive_members_only", "true"),
				),
			},
		},
	})

	// Verify that the final update body includes domain_users_only=false via ForceSendFields.
	var reqBody map[string]any
	if err := json.Unmarshal(lastUpdateBody, &reqBody); err != nil {
		t.Fatalf("failed to parse update request body: %v", err)
	}
	restrictions, ok := reqBody["restrictions"].(map[string]any)
	if !ok {
		t.Fatal("expected restrictions in update request body")
	}
	if val, exists := restrictions["domainUsersOnly"]; !exists {
		t.Error("ForceSendFields bug: domainUsersOnly not sent when toggling to false")
	} else if val != false {
		t.Errorf("expected domainUsersOnly=false, got %v", val)
	}
}

func TestAccDrive_Import(t *testing.T) {
	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/drives"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-import",
				"name": "Imported",
			})

		case r.Method == "PATCH" && strings.Contains(r.URL.Path, "/drives/drive-import"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-import",
				"name": "Imported",
				"restrictions": map[string]any{
					"adminManagedRestrictions":                  false,
					"copyRequiresWriterPermission":              false,
					"domainUsersOnly":                           false,
					"driveMembersOnly":                          false,
					"sharingFoldersRequiresOrganizerPermission": false,
				},
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/drives/drive-import"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-import",
				"name": "Imported",
				"restrictions": map[string]any{
					"adminManagedRestrictions":                  false,
					"copyRequiresWriterPermission":              false,
					"domainUsersOnly":                           false,
					"driveMembersOnly":                          false,
					"sharingFoldersRequiresOrganizerPermission": false,
				},
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/drives/drive-import"):
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
resource "googleworkspace_drive" "test" {
  name                   = "Imported"
  use_domain_admin_access = true

  restrictions {}
}
`,
			},
			{
				ResourceName:            "googleworkspace_drive.test",
				ImportState:             true,
				ImportStateId:           "true,drive-import",
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"restrictions"},
			},
		},
	})
}

func TestAccDrive_DeleteNotFound(t *testing.T) {
	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/drives"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-gone",
				"name": "Gone",
			})

		case r.Method == "PATCH" && strings.Contains(r.URL.Path, "/drives/drive-gone"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-gone",
				"name": "Gone",
				"restrictions": map[string]any{
					"adminManagedRestrictions":                  false,
					"copyRequiresWriterPermission":              false,
					"domainUsersOnly":                           false,
					"driveMembersOnly":                          false,
					"sharingFoldersRequiresOrganizerPermission": false,
				},
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/drives/drive-gone"):
			// Simulate the drive being deleted externally: return 404 on read.
			jsonResponse(w, 404, map[string]any{
				"error": map[string]any{
					"code":    404,
					"message": "File not found",
				},
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/drives/drive-gone"):
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
resource "googleworkspace_drive" "test" {
  name                   = "Gone"
  use_domain_admin_access = true

  restrictions {}
}
`,
				// After apply, the post-apply refresh Read returns 404, removing the
				// resource from state. The refresh plan shows a create is needed.
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccDrive_NoRestrictionsBlock(t *testing.T) {
	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/drives"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-norestr",
				"name": "NoRestrictions",
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/drives/drive-norestr"):
			jsonResponse(w, 200, map[string]any{
				"kind": "drive#drive",
				"id":   "drive-norestr",
				"name": "NoRestrictions",
				"restrictions": map[string]any{
					"adminManagedRestrictions":                  false,
					"copyRequiresWriterPermission":              false,
					"domainUsersOnly":                           false,
					"driveMembersOnly":                          false,
					"sharingFoldersRequiresOrganizerPermission": false,
				},
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/drives/drive-norestr"):
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
resource "googleworkspace_drive" "test" {
  name = "NoRestrictions"
}
`,
				Check: resource.TestCheckResourceAttr("googleworkspace_drive.test", "id", "drive-norestr"),
			},
		},
	})
}
