package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroup_Create(t *testing.T) {
	var createBody []byte

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/groups"):
			body, _ := io.ReadAll(r.Body)
			createBody = body
			jsonResponse(w, 200, map[string]any{
				"kind":  "admin#directory#group",
				"id":    "group-abc",
				"email": "builders@way.cloud",
				"name":  "Builders",
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/groups/group-abc"):
			jsonResponse(w, 200, map[string]any{
				"kind":        "admin#directory#group",
				"id":          "group-abc",
				"email":       "builders@way.cloud",
				"name":        "Builders",
				"description": "The builders team",
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/groups/group-abc"):
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
resource "googleworkspace_group" "test" {
  email       = "builders@way.cloud"
  name        = "Builders"
  description = "The builders team"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_group.test", "id", "group-abc"),
					resource.TestCheckResourceAttr("googleworkspace_group.test", "email", "builders@way.cloud"),
					resource.TestCheckResourceAttr("googleworkspace_group.test", "name", "Builders"),
					resource.TestCheckResourceAttr("googleworkspace_group.test", "description", "The builders team"),
				),
			},
		},
	})

	var reqBody map[string]any
	if err := json.Unmarshal(createBody, &reqBody); err != nil {
		t.Fatalf("failed to parse create request body: %v", err)
	}
	if reqBody["email"] != "builders@way.cloud" {
		t.Errorf("expected email=builders@way.cloud, got %v", reqBody["email"])
	}
	if reqBody["name"] != "Builders" {
		t.Errorf("expected name=Builders, got %v", reqBody["name"])
	}
}

func TestAccGroup_Update(t *testing.T) {
	step := 0

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/groups"):
			jsonResponse(w, 200, map[string]any{
				"kind":  "admin#directory#group",
				"id":    "group-upd",
				"email": "merchants@way.cloud",
				"name":  "Merchants",
			})

		case r.Method == "PUT" && strings.Contains(r.URL.Path, "/admin/directory/v1/groups/group-upd"):
			step++
			jsonResponse(w, 200, map[string]any{
				"kind":        "admin#directory#group",
				"id":          "group-upd",
				"email":       "merchants@way.cloud",
				"name":        "Merchants Team",
				"description": "Updated description",
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/groups/group-upd"):
			name := "Merchants"
			desc := ""
			if step > 0 {
				name = "Merchants Team"
				desc = "Updated description"
			}
			resp := map[string]any{
				"kind":  "admin#directory#group",
				"id":    "group-upd",
				"email": "merchants@way.cloud",
				"name":  name,
			}
			if desc != "" {
				resp["description"] = desc
			}
			jsonResponse(w, 200, resp)

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/groups/group-upd"):
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
resource "googleworkspace_group" "test" {
  email = "merchants@way.cloud"
  name  = "Merchants"
}
`,
				Check: resource.TestCheckResourceAttr("googleworkspace_group.test", "name", "Merchants"),
			},
			{
				Config: testProviderConfig + `
resource "googleworkspace_group" "test" {
  email       = "merchants@way.cloud"
  name        = "Merchants Team"
  description = "Updated description"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_group.test", "name", "Merchants Team"),
					resource.TestCheckResourceAttr("googleworkspace_group.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccGroup_Import(t *testing.T) {
	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.HasSuffix(r.URL.Path, "/admin/directory/v1/groups"):
			jsonResponse(w, 200, map[string]any{
				"kind":  "admin#directory#group",
				"id":    "group-imp",
				"email": "imported@way.cloud",
				"name":  "Imported",
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/groups/group-imp"):
			jsonResponse(w, 200, map[string]any{
				"kind":  "admin#directory#group",
				"id":    "group-imp",
				"email": "imported@way.cloud",
				"name":  "Imported",
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/groups/group-imp"):
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
resource "googleworkspace_group" "test" {
  email = "imported@way.cloud"
  name  = "Imported"
}
`,
			},
			{
				ResourceName:      "googleworkspace_group.test",
				ImportState:       true,
				ImportStateId:     "group-imp",
				ImportStateVerify: true,
			},
		},
	})
}
