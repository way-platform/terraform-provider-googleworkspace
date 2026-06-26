package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUser_Create(t *testing.T) {
	var createBody []byte

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		// Create user
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/users"):
			body, _ := io.ReadAll(r.Body)
			createBody = body
			jsonResponse(w, 200, map[string]any{
				"kind":         "admin#directory#user",
				"id":           "user-123",
				"primaryEmail": "oskari@way.cloud",
				"name": map[string]any{
					"givenName":  "Oskari",
					"familyName": "de Souza",
				},
				"orgUnitPath": "/Engineering",
				"suspended":   false,
				"archived":    false,
			})

		// Read user
		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/users/user-123"):
			jsonResponse(w, 200, map[string]any{
				"kind":         "admin#directory#user",
				"id":           "user-123",
				"primaryEmail": "oskari@way.cloud",
				"name": map[string]any{
					"givenName":  "Oskari",
					"familyName": "de Souza",
				},
				"orgUnitPath": "/Engineering",
				"suspended":   false,
				"archived":    false,
			})

		// Delete user
		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/users/user-123"):
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
resource "googleworkspace_user" "test" {
  primary_email = "oskari@way.cloud"
  org_unit_path = "/Engineering"

  name {
    given_name  = "Oskari"
    family_name = "de Souza"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_user.test", "id", "user-123"),
					resource.TestCheckResourceAttr("googleworkspace_user.test", "primary_email", "oskari@way.cloud"),
					resource.TestCheckResourceAttr("googleworkspace_user.test", "org_unit_path", "/Engineering"),
					resource.TestCheckResourceAttr("googleworkspace_user.test", "suspended", "false"),
					resource.TestCheckResourceAttr("googleworkspace_user.test", "archived", "false"),
				),
			},
		},
	})

	// Verify create request body includes password and ForceSendFields for booleans
	var reqBody map[string]any
	if err := json.Unmarshal(createBody, &reqBody); err != nil {
		t.Fatalf("failed to parse create request body: %v", err)
	}
	if reqBody["primaryEmail"] != "oskari@way.cloud" {
		t.Errorf("expected primaryEmail=oskari@way.cloud, got %v", reqBody["primaryEmail"])
	}
	if reqBody["password"] == nil || reqBody["password"] == "" {
		t.Error("expected password to be set on create")
	}
	if reqBody["changePasswordAtNextLogin"] != true {
		t.Error("expected changePasswordAtNextLogin=true on create")
	}
	// Verify ForceSendFields: suspended and archived should be explicitly present
	// even when false, thanks to ForceSendFields.
	if _, exists := reqBody["suspended"]; !exists {
		t.Error("ForceSendFields bug: suspended not present in create request body")
	}
	if _, exists := reqBody["archived"]; !exists {
		t.Error("ForceSendFields bug: archived not present in create request body")
	}
}

func TestAccUser_Update(t *testing.T) {
	var updateBody []byte
	step := 0

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/users"):
			jsonResponse(w, 200, map[string]any{
				"kind":         "admin#directory#user",
				"id":           "user-456",
				"primaryEmail": "alice@way.cloud",
				"name": map[string]any{
					"givenName":  "Alice",
					"familyName": "Smith",
				},
				"orgUnitPath": "/Engineering",
				"suspended":   false,
				"archived":    false,
			})

		case r.Method == "PUT" && strings.Contains(r.URL.Path, "/admin/directory/v1/users/user-456"):
			body, _ := io.ReadAll(r.Body)
			updateBody = body
			step++
			jsonResponse(w, 200, map[string]any{
				"kind":         "admin#directory#user",
				"id":           "user-456",
				"primaryEmail": "alice@way.cloud",
				"name": map[string]any{
					"givenName":  "Alice",
					"familyName": "Smith",
				},
				"orgUnitPath": "/Deactivated",
				"suspended":   true,
				"archived":    false,
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/users/user-456"):
			suspended := false
			orgUnit := "/Engineering"
			if step > 0 {
				suspended = true
				orgUnit = "/Deactivated"
			}
			jsonResponse(w, 200, map[string]any{
				"kind":         "admin#directory#user",
				"id":           "user-456",
				"primaryEmail": "alice@way.cloud",
				"name": map[string]any{
					"givenName":  "Alice",
					"familyName": "Smith",
				},
				"orgUnitPath": orgUnit,
				"suspended":   suspended,
				"archived":    false,
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/users/user-456"):
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
resource "googleworkspace_user" "test" {
  primary_email = "alice@way.cloud"
  org_unit_path = "/Engineering"

  name {
    given_name  = "Alice"
    family_name = "Smith"
  }
}
`,
			},
			{
				Config: testProviderConfig + `
resource "googleworkspace_user" "test" {
  primary_email = "alice@way.cloud"
  org_unit_path = "/Deactivated"
  suspended     = true

  name {
    given_name  = "Alice"
    family_name = "Smith"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_user.test", "suspended", "true"),
					resource.TestCheckResourceAttr("googleworkspace_user.test", "org_unit_path", "/Deactivated"),
				),
			},
		},
	})

	// Verify ForceSendFields in update: archived=false must be sent
	var reqBody map[string]any
	if err := json.Unmarshal(updateBody, &reqBody); err != nil {
		t.Fatalf("failed to parse update request body: %v", err)
	}
	if _, exists := reqBody["archived"]; !exists {
		t.Error("ForceSendFields bug: archived not present in update request body")
	}
	if reqBody["suspended"] != true {
		t.Errorf("expected suspended=true in update, got %v", reqBody["suspended"])
	}
}

func TestAccUser_CustomSchemas(t *testing.T) {
	var createBody []byte

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/schemas/Mapbox"):
			jsonResponse(w, 200, schemaResponse("schema-mapbox", "Mapbox", "Mapbox", []map[string]any{
				{
					"fieldName":   "role",
					"fieldType":   "STRING",
					"multiValued": false,
					"indexed":     true,
				},
				{
					"fieldName":   "teams",
					"fieldType":   "STRING",
					"multiValued": true,
					"indexed":     true,
				},
			}))

		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/users"):
			body, _ := io.ReadAll(r.Body)
			createBody = body
			jsonResponse(w, 200, map[string]any{
				"kind":         "admin#directory#user",
				"id":           "user-custom",
				"primaryEmail": "oskari@way.cloud",
				"name": map[string]any{
					"givenName":  "Oskari",
					"familyName": "Pétas",
				},
				"orgUnitPath": "/",
				"suspended":   false,
				"archived":    false,
				"customSchemas": map[string]any{
					"Mapbox": map[string]any{
						"role": "Root",
						"teams": []map[string]any{
							{"value": "Platform"},
							{"value": "Ops"},
						},
					},
				},
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/users/user-custom"):
			if got := r.URL.Query().Get("projection"); got != "full" {
				t.Errorf("expected projection=full, got %q", got)
			}
			jsonResponse(w, 200, map[string]any{
				"kind":         "admin#directory#user",
				"id":           "user-custom",
				"primaryEmail": "oskari@way.cloud",
				"name": map[string]any{
					"givenName":  "Oskari",
					"familyName": "Pétas",
				},
				"orgUnitPath": "/",
				"suspended":   false,
				"archived":    false,
				"customSchemas": map[string]any{
					"Mapbox": map[string]any{
						"role": "Root",
						"teams": []map[string]any{
							{"value": "Platform"},
							{"value": "Ops"},
						},
					},
				},
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/users/user-custom"):
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
resource "googleworkspace_user" "test" {
  primary_email = "oskari@way.cloud"

  name {
    given_name  = "Oskari"
    family_name = "Pétas"
  }

  custom_schemas {
    schema_name = "Mapbox"

    schema_values = {
      role  = jsonencode("Root")
      teams = jsonencode(["Platform", "Ops"])
    }
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_user.test", "id", "user-custom"),
					resource.TestCheckResourceAttr("googleworkspace_user.test", "custom_schemas.0.schema_name", "Mapbox"),
					resource.TestCheckResourceAttr("googleworkspace_user.test", "custom_schemas.0.schema_values.role", `"Root"`),
					resource.TestCheckResourceAttr("googleworkspace_user.test", "custom_schemas.0.schema_values.teams", `["Platform","Ops"]`),
				),
			},
		},
	})

	var reqBody map[string]any
	if err := json.Unmarshal(createBody, &reqBody); err != nil {
		t.Fatalf("failed to parse create request body: %v", err)
	}
	customSchemas, ok := reqBody["customSchemas"].(map[string]any)
	if !ok {
		t.Fatalf("expected customSchemas object, got %#v", reqBody["customSchemas"])
	}
	mapbox, ok := customSchemas["Mapbox"].(map[string]any)
	if !ok {
		t.Fatalf("expected Mapbox custom schema, got %#v", customSchemas["Mapbox"])
	}
	if mapbox["role"] != "Root" {
		t.Errorf("expected Mapbox.role=Root, got %v", mapbox["role"])
	}
	teams, ok := mapbox["teams"].([]any)
	if !ok || len(teams) != 2 {
		t.Fatalf("expected two Mapbox.teams values, got %#v", mapbox["teams"])
	}
	firstTeam, ok := teams[0].(map[string]any)
	if !ok || firstTeam["value"] != "Platform" {
		t.Errorf("expected first Mapbox.teams value=Platform, got %#v", teams[0])
	}
}

func TestAccUser_Import(t *testing.T) {
	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/users"):
			jsonResponse(w, 200, map[string]any{
				"kind":         "admin#directory#user",
				"id":           "user-import",
				"primaryEmail": "bob@way.cloud",
				"name": map[string]any{
					"givenName":  "Bob",
					"familyName": "Jones",
				},
				"orgUnitPath": "/",
				"suspended":   false,
				"archived":    false,
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/users/user-import"):
			jsonResponse(w, 200, map[string]any{
				"kind":         "admin#directory#user",
				"id":           "user-import",
				"primaryEmail": "bob@way.cloud",
				"name": map[string]any{
					"givenName":  "Bob",
					"familyName": "Jones",
				},
				"orgUnitPath": "/",
				"suspended":   false,
				"archived":    false,
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/users/user-import"):
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
resource "googleworkspace_user" "test" {
  primary_email = "bob@way.cloud"
  org_unit_path = "/"

  name {
    given_name  = "Bob"
    family_name = "Jones"
  }
}
`,
			},
			{
				ResourceName:      "googleworkspace_user.test",
				ImportState:       true,
				ImportStateId:     "user-import",
				ImportStateVerify: true,
			},
		},
	})
}
