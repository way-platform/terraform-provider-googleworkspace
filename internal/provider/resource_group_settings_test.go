package provider

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGroupSettings_Create(t *testing.T) {
	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "PUT" && r.URL.Path == "/team@way.cloud":
			jsonResponse(w, 200, map[string]any{
				"kind":                 "groupsSettings#groups",
				"email":               "team@way.cloud",
				"whoCanViewMembership": "ALL_IN_DOMAIN_CAN_VIEW",
			})

		case r.Method == "GET" && r.URL.Path == "/team@way.cloud":
			jsonResponse(w, 200, map[string]any{
				"kind":                 "groupsSettings#groups",
				"email":               "team@way.cloud",
				"whoCanViewMembership": "ALL_IN_DOMAIN_CAN_VIEW",
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
resource "googleworkspace_group_settings" "test" {
  email                    = "team@way.cloud"
  who_can_view_membership = "ALL_IN_DOMAIN_CAN_VIEW"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_group_settings.test", "id", "team@way.cloud"),
					resource.TestCheckResourceAttr("googleworkspace_group_settings.test", "who_can_view_membership", "ALL_IN_DOMAIN_CAN_VIEW"),
				),
			},
		},
	})
}

func TestAccGroupSettings_Update(t *testing.T) {
	step := 0

	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "PUT" && r.URL.Path == "/team@way.cloud":
			step++
			membership := "ALL_IN_DOMAIN_CAN_VIEW"
			if step > 1 {
				membership = "ALL_MEMBERS_CAN_VIEW"
			}
			jsonResponse(w, 200, map[string]any{
				"kind":                 "groupsSettings#groups",
				"email":               "team@way.cloud",
				"whoCanViewMembership": membership,
			})

		case r.Method == "GET" && r.URL.Path == "/team@way.cloud":
			membership := "ALL_IN_DOMAIN_CAN_VIEW"
			if step > 1 {
				membership = "ALL_MEMBERS_CAN_VIEW"
			}
			jsonResponse(w, 200, map[string]any{
				"kind":                 "groupsSettings#groups",
				"email":               "team@way.cloud",
				"whoCanViewMembership": membership,
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
resource "googleworkspace_group_settings" "test" {
  email                    = "team@way.cloud"
  who_can_view_membership = "ALL_IN_DOMAIN_CAN_VIEW"
}
`,
				Check: resource.TestCheckResourceAttr("googleworkspace_group_settings.test", "who_can_view_membership", "ALL_IN_DOMAIN_CAN_VIEW"),
			},
			{
				Config: testProviderConfig + `
resource "googleworkspace_group_settings" "test" {
  email                    = "team@way.cloud"
  who_can_view_membership = "ALL_MEMBERS_CAN_VIEW"
}
`,
				Check: resource.TestCheckResourceAttr("googleworkspace_group_settings.test", "who_can_view_membership", "ALL_MEMBERS_CAN_VIEW"),
			},
		},
	})
}

func TestAccGroupSettings_Import(t *testing.T) {
	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "PUT" && r.URL.Path == "/imported@way.cloud":
			jsonResponse(w, 200, map[string]any{
				"kind":                 "groupsSettings#groups",
				"email":               "imported@way.cloud",
				"whoCanViewMembership": "ALL_IN_DOMAIN_CAN_VIEW",
			})

		case r.Method == "GET" && r.URL.Path == "/imported@way.cloud":
			jsonResponse(w, 200, map[string]any{
				"kind":                 "groupsSettings#groups",
				"email":               "imported@way.cloud",
				"whoCanViewMembership": "ALL_IN_DOMAIN_CAN_VIEW",
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
resource "googleworkspace_group_settings" "test" {
  email                    = "imported@way.cloud"
  who_can_view_membership = "ALL_IN_DOMAIN_CAN_VIEW"
}
`,
			},
			{
				ResourceName:      "googleworkspace_group_settings.test",
				ImportState:       true,
				ImportStateId:     "imported@way.cloud",
				ImportStateVerify: true,
			},
		},
	})
}
