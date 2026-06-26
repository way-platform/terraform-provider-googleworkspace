package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSchema_Create(t *testing.T) {
	var createBody []byte

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/customer/C00000000/schemas"):
			body, _ := io.ReadAll(r.Body)
			createBody = body
			jsonResponse(w, 200, schemaResponse("schema-123", "Mapbox", "Mapbox", []map[string]any{
				{
					"fieldName":      "role",
					"fieldType":      "STRING",
					"displayName":    "Role",
					"multiValued":    false,
					"indexed":        true,
					"readAccessType": "ADMINS_AND_SELF",
				},
			}))

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/schemas/schema-123"):
			jsonResponse(w, 200, schemaResponse("schema-123", "Mapbox", "Mapbox", []map[string]any{
				{
					"fieldName":      "role",
					"fieldType":      "STRING",
					"displayName":    "Role",
					"multiValued":    false,
					"indexed":        true,
					"readAccessType": "ADMINS_AND_SELF",
				},
			}))

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/schemas/schema-123"):
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
resource "googleworkspace_schema" "test" {
  schema_name  = "Mapbox"
  display_name = "Mapbox"

  field {
    field_name       = "role"
    field_type       = "STRING"
    display_name     = "Role"
    multi_valued     = false
    indexed          = true
    read_access_type = "ADMINS_AND_SELF"
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_schema.test", "id", "schema-123"),
					resource.TestCheckResourceAttr("googleworkspace_schema.test", "schema_name", "Mapbox"),
					resource.TestCheckResourceAttr("googleworkspace_schema.test", "display_name", "Mapbox"),
					resource.TestCheckResourceAttr("googleworkspace_schema.test", "field.0.field_name", "role"),
					resource.TestCheckResourceAttr("googleworkspace_schema.test", "field.0.field_type", "STRING"),
					resource.TestCheckResourceAttr("googleworkspace_schema.test", "field.0.multi_valued", "false"),
					resource.TestCheckResourceAttr("googleworkspace_schema.test", "field.0.indexed", "true"),
					resource.TestCheckResourceAttr("googleworkspace_schema.test", "field.0.read_access_type", "ADMINS_AND_SELF"),
				),
			},
		},
	})

	var reqBody map[string]any
	if err := json.Unmarshal(createBody, &reqBody); err != nil {
		t.Fatalf("failed to parse create request body: %v", err)
	}
	fields, ok := reqBody["fields"].([]any)
	if !ok || len(fields) != 1 {
		t.Fatalf("expected one field in create request, got %#v", reqBody["fields"])
	}
	field, ok := fields[0].(map[string]any)
	if !ok {
		t.Fatalf("expected field object, got %#v", fields[0])
	}
	if _, exists := field["multiValued"]; !exists {
		t.Error("ForceSendFields bug: multiValued not present in create request body")
	}
	if field["multiValued"] != false {
		t.Errorf("expected multiValued=false, got %v", field["multiValued"])
	}
	if field["indexed"] != true {
		t.Errorf("expected indexed=true, got %v", field["indexed"])
	}
}

func TestAccSchema_Update(t *testing.T) {
	step := 0
	var updateBody []byte

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/customer/C00000000/schemas"):
			jsonResponse(w, 200, schemaResponse("schema-456", "Mapbox", "Mapbox", []map[string]any{
				{
					"fieldName":   "role",
					"fieldType":   "STRING",
					"displayName": "Role",
					"multiValued": false,
					"indexed":     true,
				},
			}))

		case r.Method == "PUT" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/schemas/schema-456"):
			body, _ := io.ReadAll(r.Body)
			updateBody = body
			step++
			jsonResponse(w, 200, schemaResponse("schema-456", "Mapbox", "Mapbox Profile", []map[string]any{
				{
					"fieldName":   "role",
					"fieldType":   "STRING",
					"displayName": "Role",
					"multiValued": false,
					"indexed":     false,
				},
				{
					"fieldName":   "team",
					"fieldType":   "STRING",
					"displayName": "Team",
					"multiValued": true,
					"indexed":     true,
				},
			}))

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/schemas/schema-456"):
			displayName := "Mapbox"
			fields := []map[string]any{
				{
					"fieldName":   "role",
					"fieldType":   "STRING",
					"displayName": "Role",
					"multiValued": false,
					"indexed":     true,
				},
			}
			if step > 0 {
				displayName = "Mapbox Profile"
				fields = []map[string]any{
					{
						"fieldName":   "role",
						"fieldType":   "STRING",
						"displayName": "Role",
						"multiValued": false,
						"indexed":     false,
					},
					{
						"fieldName":   "team",
						"fieldType":   "STRING",
						"displayName": "Team",
						"multiValued": true,
						"indexed":     true,
					},
				}
			}
			jsonResponse(w, 200, schemaResponse("schema-456", "Mapbox", displayName, fields))

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/schemas/schema-456"):
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
resource "googleworkspace_schema" "test" {
  schema_name  = "Mapbox"
  display_name = "Mapbox"

  field {
    field_name   = "role"
    field_type   = "STRING"
    display_name = "Role"
  }
}
`,
			},
			{
				Config: testProviderConfig + `
resource "googleworkspace_schema" "test" {
  schema_name  = "Mapbox"
  display_name = "Mapbox Profile"

  field {
    field_name   = "role"
    field_type   = "STRING"
    display_name = "Role"
    indexed      = false
  }

  field {
    field_name   = "team"
    field_type   = "STRING"
    display_name = "Team"
    multi_valued = true
  }
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_schema.test", "display_name", "Mapbox Profile"),
					resource.TestCheckResourceAttr("googleworkspace_schema.test", "field.#", "2"),
					resource.TestCheckResourceAttr("googleworkspace_schema.test", "field.0.indexed", "false"),
					resource.TestCheckResourceAttr("googleworkspace_schema.test", "field.1.multi_valued", "true"),
				),
			},
		},
	})

	var reqBody map[string]any
	if err := json.Unmarshal(updateBody, &reqBody); err != nil {
		t.Fatalf("failed to parse update request body: %v", err)
	}
	fields, ok := reqBody["fields"].([]any)
	if !ok || len(fields) != 2 {
		t.Fatalf("expected two fields in update request, got %#v", reqBody["fields"])
	}
	first := fields[0].(map[string]any)
	if first["indexed"] != false {
		t.Errorf("expected indexed=false in update request, got %v", first["indexed"])
	}
}

func TestAccSchema_Import(t *testing.T) {
	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/customer/C00000000/schemas"):
			jsonResponse(w, 200, schemaResponse("schema-import", "Mapbox", "Mapbox", []map[string]any{
				{
					"fieldName":   "role",
					"fieldType":   "STRING",
					"displayName": "Role",
					"multiValued": false,
					"indexed":     true,
				},
			}))

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/schemas/schema-import"):
			jsonResponse(w, 200, schemaResponse("schema-import", "Mapbox", "Mapbox", []map[string]any{
				{
					"fieldName":   "role",
					"fieldType":   "STRING",
					"displayName": "Role",
					"multiValued": false,
					"indexed":     true,
				},
			}))

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/schemas/Mapbox"):
			jsonResponse(w, 200, schemaResponse("schema-import", "Mapbox", "Mapbox", []map[string]any{
				{
					"fieldName":   "role",
					"fieldType":   "STRING",
					"displayName": "Role",
					"multiValued": false,
					"indexed":     true,
				},
			}))

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/schemas/schema-import"):
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
resource "googleworkspace_schema" "test" {
  schema_name  = "Mapbox"
  display_name = "Mapbox"

  field {
    field_name   = "role"
    field_type   = "STRING"
    display_name = "Role"
  }
}
`,
			},
			{
				ResourceName:      "googleworkspace_schema.test",
				ImportState:       true,
				ImportStateId:     "Mapbox",
				ImportStateVerify: true,
			},
		},
	})
}

func schemaResponse(id, name, displayName string, fields []map[string]any) map[string]any {
	return map[string]any{
		"kind":        "admin#directory#schema",
		"schemaId":    id,
		"schemaName":  name,
		"displayName": displayName,
		"fields":      fields,
	}
}
