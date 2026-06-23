package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOrgUnit_Create(t *testing.T) {
	var createBody []byte

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits"):
			body, _ := io.ReadAll(r.Body)
			createBody = body
			jsonResponse(w, 200, map[string]any{
				"kind":              "admin#directory#orgUnit",
				"orgUnitId":        "id:ou-001",
				"name":             "Engineering",
				"orgUnitPath":      "/Engineering",
				"parentOrgUnitPath": "/",
				"description":      "Engineering team",
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits/id:ou-001"):
			jsonResponse(w, 200, map[string]any{
				"kind":              "admin#directory#orgUnit",
				"orgUnitId":        "id:ou-001",
				"name":             "Engineering",
				"orgUnitPath":      "/Engineering",
				"parentOrgUnitPath": "/",
				"description":      "Engineering team",
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits/id:ou-001"):
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
resource "googleworkspace_org_unit" "test" {
  name                 = "Engineering"
  parent_org_unit_path = "/"
  description          = "Engineering team"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_org_unit.test", "id", "ou-001"),
					resource.TestCheckResourceAttr("googleworkspace_org_unit.test", "name", "Engineering"),
					resource.TestCheckResourceAttr("googleworkspace_org_unit.test", "org_unit_path", "/Engineering"),
					resource.TestCheckResourceAttr("googleworkspace_org_unit.test", "parent_org_unit_path", "/"),
					resource.TestCheckResourceAttr("googleworkspace_org_unit.test", "description", "Engineering team"),
				),
			},
		},
	})

	var reqBody map[string]any
	if err := json.Unmarshal(createBody, &reqBody); err != nil {
		t.Fatalf("failed to parse create request body: %v", err)
	}
	if reqBody["name"] != "Engineering" {
		t.Errorf("expected name=Engineering, got %v", reqBody["name"])
	}
	if reqBody["parentOrgUnitPath"] != "/" {
		t.Errorf("expected parentOrgUnitPath=/, got %v", reqBody["parentOrgUnitPath"])
	}
}

func TestAccOrgUnit_Update(t *testing.T) {
	step := 0

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits"):
			jsonResponse(w, 200, map[string]any{
				"kind":              "admin#directory#orgUnit",
				"orgUnitId":        "id:ou-002",
				"name":             "Platform",
				"orgUnitPath":      "/Platform",
				"parentOrgUnitPath": "/",
			})

		case r.Method == "PUT" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits/id:ou-002"):
			step++
			jsonResponse(w, 200, map[string]any{
				"kind":              "admin#directory#orgUnit",
				"orgUnitId":        "id:ou-002",
				"name":             "Platform Engineering",
				"orgUnitPath":      "/Platform Engineering",
				"parentOrgUnitPath": "/",
				"description":      "Platform eng",
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits/id:ou-002"):
			name := "Platform"
			ouPath := "/Platform"
			desc := ""
			if step > 0 {
				name = "Platform Engineering"
				ouPath = "/Platform Engineering"
				desc = "Platform eng"
			}
			resp := map[string]any{
				"kind":              "admin#directory#orgUnit",
				"orgUnitId":        "id:ou-002",
				"name":             name,
				"orgUnitPath":      ouPath,
				"parentOrgUnitPath": "/",
			}
			if desc != "" {
				resp["description"] = desc
			}
			jsonResponse(w, 200, resp)

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits/id:ou-002"):
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
resource "googleworkspace_org_unit" "test" {
  name                 = "Platform"
  parent_org_unit_path = "/"
}
`,
				Check: resource.TestCheckResourceAttr("googleworkspace_org_unit.test", "name", "Platform"),
			},
			{
				Config: testProviderConfig + `
resource "googleworkspace_org_unit" "test" {
  name                 = "Platform Engineering"
  parent_org_unit_path = "/"
  description          = "Platform eng"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_org_unit.test", "name", "Platform Engineering"),
					resource.TestCheckResourceAttr("googleworkspace_org_unit.test", "org_unit_path", "/Platform Engineering"),
					resource.TestCheckResourceAttr("googleworkspace_org_unit.test", "description", "Platform eng"),
				),
			},
		},
	})
}

func TestAccOrgUnit_Import(t *testing.T) {
	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits"):
			jsonResponse(w, 200, map[string]any{
				"kind":              "admin#directory#orgUnit",
				"orgUnitId":        "id:ou-imp",
				"name":             "Imported",
				"orgUnitPath":      "/Imported",
				"parentOrgUnitPath": "/",
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits/id:ou-imp"):
			jsonResponse(w, 200, map[string]any{
				"kind":              "admin#directory#orgUnit",
				"orgUnitId":        "id:ou-imp",
				"name":             "Imported",
				"orgUnitPath":      "/Imported",
				"parentOrgUnitPath": "/",
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits/id:ou-imp"):
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
resource "googleworkspace_org_unit" "test" {
  name                 = "Imported"
  parent_org_unit_path = "/"
}
`,
			},
			{
				ResourceName:      "googleworkspace_org_unit.test",
				ImportState:       true,
				ImportStateId:     "ou-imp",
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccOrgUnit_ClearDescription(t *testing.T) {
	step := 0

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits"):
			jsonResponse(w, 200, map[string]any{
				"kind":              "admin#directory#orgUnit",
				"orgUnitId":        "id:ou-desc",
				"name":             "DescOU",
				"orgUnitPath":      "/DescOU",
				"parentOrgUnitPath": "/",
				"description":      "Initial",
			})

		case r.Method == "PUT" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits/id:ou-desc"):
			step++
			jsonResponse(w, 200, map[string]any{
				"kind":              "admin#directory#orgUnit",
				"orgUnitId":        "id:ou-desc",
				"name":             "DescOU",
				"orgUnitPath":      "/DescOU",
				"parentOrgUnitPath": "/",
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits/id:ou-desc"):
			resp := map[string]any{
				"kind":              "admin#directory#orgUnit",
				"orgUnitId":        "id:ou-desc",
				"name":             "DescOU",
				"orgUnitPath":      "/DescOU",
				"parentOrgUnitPath": "/",
			}
			if step == 0 {
				resp["description"] = "Initial"
			}
			jsonResponse(w, 200, resp)

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/orgunits/id:ou-desc"):
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
resource "googleworkspace_org_unit" "test" {
  name                 = "DescOU"
  parent_org_unit_path = "/"
  description          = "Initial"
}
`,
				Check: resource.TestCheckResourceAttr("googleworkspace_org_unit.test", "description", "Initial"),
			},
			{
				Config: testProviderConfig + `
resource "googleworkspace_org_unit" "test" {
  name                 = "DescOU"
  parent_org_unit_path = "/"
}
`,
				Check: resource.TestCheckNoResourceAttr("googleworkspace_org_unit.test", "description"),
			},
		},
	})
}
