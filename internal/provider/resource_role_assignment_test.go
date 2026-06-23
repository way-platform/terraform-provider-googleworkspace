package provider

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRoleAssignment_Create(t *testing.T) {
	var createBody []byte

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/roleassignments"):
			body, _ := io.ReadAll(r.Body)
			createBody = body
			jsonResponse(w, 200, map[string]any{
				"kind":             "admin#directory#roleAssignment",
				"roleAssignmentId": "12345",
				"roleId":           "67890",
				"assignedTo":       "user-abc",
				"scopeType":        "CUSTOMER",
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/roleassignments/12345"):
			jsonResponse(w, 200, map[string]any{
				"kind":             "admin#directory#roleAssignment",
				"roleAssignmentId": "12345",
				"roleId":           "67890",
				"assignedTo":       "user-abc",
				"scopeType":        "CUSTOMER",
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/roleassignments/12345"):
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
resource "googleworkspace_role_assignment" "test" {
  role_id     = "67890"
  assigned_to = "user-abc"
  scope_type  = "CUSTOMER"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_role_assignment.test", "id", "12345"),
					resource.TestCheckResourceAttr("googleworkspace_role_assignment.test", "role_id", "67890"),
					resource.TestCheckResourceAttr("googleworkspace_role_assignment.test", "assigned_to", "user-abc"),
					resource.TestCheckResourceAttr("googleworkspace_role_assignment.test", "scope_type", "CUSTOMER"),
				),
			},
		},
	})

	var reqBody map[string]any
	if err := json.Unmarshal(createBody, &reqBody); err != nil {
		t.Fatalf("failed to parse create request body: %v", err)
	}
	if reqBody["assignedTo"] != "user-abc" {
		t.Errorf("expected assignedTo=user-abc, got %v", reqBody["assignedTo"])
	}
	if reqBody["roleId"] != "67890" {
		t.Errorf("expected roleId=67890, got %v", reqBody["roleId"])
	}
}

func TestAccRoleAssignment_Import(t *testing.T) {
	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "POST" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/roleassignments"):
			jsonResponse(w, 200, map[string]any{
				"kind":             "admin#directory#roleAssignment",
				"roleAssignmentId": "99999",
				"roleId":           "11111",
				"assignedTo":       "user-xyz",
				"scopeType":        "CUSTOMER",
			})

		case r.Method == "GET" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/roleassignments/99999"):
			jsonResponse(w, 200, map[string]any{
				"kind":             "admin#directory#roleAssignment",
				"roleAssignmentId": "99999",
				"roleId":           "11111",
				"assignedTo":       "user-xyz",
				"scopeType":        "CUSTOMER",
			})

		case r.Method == "DELETE" && strings.Contains(r.URL.Path, "/admin/directory/v1/customer/C00000000/roleassignments/99999"):
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
resource "googleworkspace_role_assignment" "test" {
  role_id     = "11111"
  assigned_to = "user-xyz"
  scope_type  = "CUSTOMER"
}
`,
			},
			{
				ResourceName:      "googleworkspace_role_assignment.test",
				ImportState:       true,
				ImportStateId:     "99999",
				ImportStateVerify: true,
			},
		},
	})
}
