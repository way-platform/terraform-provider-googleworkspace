package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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

// TestAccUser_RemoveCustomSchemas verifies that a custom schema removed from config is cleared with
// explicit nulls. The API keeps every custom field that an update leaves out.
func TestAccUser_RemoveCustomSchemas(t *testing.T) {
	var mu sync.Mutex
	customSchemas := map[string]map[string]any{
		"Mapbox":     {"role": "Admin"},
		"Contractor": {"contractor": true},
	}
	var updates []map[string]any

	userResponse := func() map[string]any {
		// A schema with every field cleared is returned as an empty object.
		schemas := map[string]any{}
		for name, values := range customSchemas {
			schemas[name] = values
		}
		return map[string]any{
			"kind":          "admin#directory#user",
			"id":            "user-remove",
			"primaryEmail":  "deepti@way.cloud",
			"name":          map[string]any{"givenName": "Deepti", "familyName": "Test"},
			"orgUnitPath":   "/",
			"suspended":     false,
			"archived":      false,
			"customSchemas": schemas,
		}
	}

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		switch {
		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/customer/C00000000/schemas/Mapbox"):
			jsonResponse(w, 200, schemaResponse("schema-mapbox", "Mapbox", "Mapbox", []map[string]any{
				{"fieldName": "role", "fieldType": "STRING", "multiValued": false},
			}))

		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/customer/C00000000/schemas/Contractor"):
			jsonResponse(w, 200, schemaResponse("schema-contractor", "Contractor", "Contractor", []map[string]any{
				{"fieldName": "contractor", "fieldType": "BOOL", "multiValued": false},
			}))

		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/users"):
			jsonResponse(w, 200, userResponse())

		case r.Method == "PUT" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/users/user-remove"):
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decode request: %v", err)
			}
			updates = append(updates, body)
			sent, _ := body["customSchemas"].(map[string]any)
			for name, values := range sent {
				for field, value := range values.(map[string]any) {
					if value == nil {
						delete(customSchemas[name], field)
					} else {
						customSchemas[name][field] = value
					}
				}
			}
			jsonResponse(w, 200, userResponse())

		case r.Method == "GET" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/users/user-remove"):
			jsonResponse(w, 200, userResponse())

		case r.Method == "DELETE" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/users/user-remove"):
			w.WriteHeader(204)

		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(500)
		}
	}))
	setupTestClient(t, server)

	lastUpdateSchemas := func() map[string]any {
		mu.Lock()
		defer mu.Unlock()
		if len(updates) == 0 {
			t.Fatal("no update request was sent")
		}
		schemas, ok := updates[len(updates)-1]["customSchemas"].(map[string]any)
		if !ok {
			t.Fatalf("expected customSchemas in the update, got %#v", updates[len(updates)-1])
		}
		return schemas
	}
	user := func(customSchemas string) string {
		return testProviderConfig + `
resource "googleworkspace_user" "test" {
  primary_email = "deepti@way.cloud"

  name {
    given_name  = "Deepti"
    family_name = "Test"
  }
` + customSchemas + `
}
`
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: user(`
  custom_schemas {
    schema_name   = "Contractor"
    schema_values = { contractor = jsonencode(true) }
  }
  custom_schemas {
    schema_name   = "Mapbox"
    schema_values = { role = jsonencode("Admin") }
  }
`),
			},
			{
				// Remove one schema and keep the other.
				Config: user(`
  custom_schemas {
    schema_name   = "Mapbox"
    schema_values = { role = jsonencode("Admin") }
  }
`),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_user.test", "custom_schemas.#", "1"),
					resource.TestCheckResourceAttr("googleworkspace_user.test", "custom_schemas.0.schema_name", "Mapbox"),
					func(*terraform.State) error {
						schemas := lastUpdateSchemas()
						contractor, ok := schemas["Contractor"].(map[string]any)
						if value, set := contractor["contractor"]; !ok || !set || value != nil {
							t.Errorf("expected Contractor.contractor=null, got %#v", schemas["Contractor"])
						}
						if mapbox, _ := schemas["Mapbox"].(map[string]any); mapbox["role"] != "Admin" {
							t.Errorf("expected Mapbox.role=Admin, got %#v", schemas["Mapbox"])
						}
						return nil
					},
				),
			},
			{
				// Remove the last schema.
				Config: user(""),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_user.test", "custom_schemas.#", "0"),
					func(*terraform.State) error {
						schemas := lastUpdateSchemas()
						mapbox, ok := schemas["Mapbox"].(map[string]any)
						if value, set := mapbox["role"]; !ok || !set || value != nil {
							t.Errorf("expected Mapbox.role=null, got %#v", schemas["Mapbox"])
						}
						return nil
					},
				),
			},
		},
	})
}
