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
				"kind":                   "groupsSettings#groups",
				"email":                  "team@way.cloud",
				"whoCanViewMembership":   "ALL_IN_DOMAIN_CAN_VIEW",
				"whoCanPostMessage":      "ANYONE_CAN_POST",
				"messageModerationLevel": "MODERATE_NONE",
			})

		case r.Method == "GET" && r.URL.Path == "/team@way.cloud":
			jsonResponse(w, 200, map[string]any{
				"kind":                   "groupsSettings#groups",
				"email":                  "team@way.cloud",
				"whoCanViewMembership":   "ALL_IN_DOMAIN_CAN_VIEW",
				"whoCanPostMessage":      "ANYONE_CAN_POST",
				"messageModerationLevel": "MODERATE_NONE",
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
  who_can_view_membership  = "ALL_IN_DOMAIN_CAN_VIEW"
  who_can_post_message     = "ANYONE_CAN_POST"
  message_moderation_level = "MODERATE_NONE"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_group_settings.test", "id", "team@way.cloud"),
					resource.TestCheckResourceAttr("googleworkspace_group_settings.test", "who_can_view_membership", "ALL_IN_DOMAIN_CAN_VIEW"),
					resource.TestCheckResourceAttr("googleworkspace_group_settings.test", "who_can_post_message", "ANYONE_CAN_POST"),
					resource.TestCheckResourceAttr("googleworkspace_group_settings.test", "message_moderation_level", "MODERATE_NONE"),
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
			post := "ANYONE_CAN_POST"
			if step > 1 {
				membership = "ALL_MEMBERS_CAN_VIEW"
				post = "ALL_MEMBERS_CAN_POST"
			}
			jsonResponse(w, 200, map[string]any{
				"kind":                 "groupsSettings#groups",
				"email":                "team@way.cloud",
				"whoCanViewMembership": membership,
				"whoCanPostMessage":    post,
			})

		case r.Method == "GET" && r.URL.Path == "/team@way.cloud":
			membership := "ALL_IN_DOMAIN_CAN_VIEW"
			post := "ANYONE_CAN_POST"
			if step > 1 {
				membership = "ALL_MEMBERS_CAN_VIEW"
				post = "ALL_MEMBERS_CAN_POST"
			}
			jsonResponse(w, 200, map[string]any{
				"kind":                 "groupsSettings#groups",
				"email":                "team@way.cloud",
				"whoCanViewMembership": membership,
				"whoCanPostMessage":    post,
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
  email                   = "team@way.cloud"
  who_can_view_membership = "ALL_IN_DOMAIN_CAN_VIEW"
  who_can_post_message    = "ANYONE_CAN_POST"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_group_settings.test", "who_can_view_membership", "ALL_IN_DOMAIN_CAN_VIEW"),
					resource.TestCheckResourceAttr("googleworkspace_group_settings.test", "who_can_post_message", "ANYONE_CAN_POST"),
				),
			},
			{
				Config: testProviderConfig + `
resource "googleworkspace_group_settings" "test" {
  email                   = "team@way.cloud"
  who_can_view_membership = "ALL_MEMBERS_CAN_VIEW"
  who_can_post_message    = "ALL_MEMBERS_CAN_POST"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_group_settings.test", "who_can_view_membership", "ALL_MEMBERS_CAN_VIEW"),
					resource.TestCheckResourceAttr("googleworkspace_group_settings.test", "who_can_post_message", "ALL_MEMBERS_CAN_POST"),
				),
			},
		},
	})
}

func TestAccGroupSettings_Import(t *testing.T) {
	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "PUT" && r.URL.Path == "/imported@way.cloud":
			jsonResponse(w, 200, map[string]any{
				"kind":                   "groupsSettings#groups",
				"email":                  "imported@way.cloud",
				"whoCanViewMembership":   "ALL_IN_DOMAIN_CAN_VIEW",
				"whoCanPostMessage":      "ANYONE_CAN_POST",
				"messageModerationLevel": "MODERATE_NONE",
			})

		case r.Method == "GET" && r.URL.Path == "/imported@way.cloud":
			jsonResponse(w, 200, map[string]any{
				"kind":                   "groupsSettings#groups",
				"email":                  "imported@way.cloud",
				"whoCanViewMembership":   "ALL_IN_DOMAIN_CAN_VIEW",
				"whoCanPostMessage":      "ANYONE_CAN_POST",
				"messageModerationLevel": "MODERATE_NONE",
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
  who_can_view_membership  = "ALL_IN_DOMAIN_CAN_VIEW"
  who_can_post_message     = "ANYONE_CAN_POST"
  message_moderation_level = "MODERATE_NONE"
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

// TestAccGroupSettings_ComputedDefaults verifies that when who_can_post_message is left unset, it is
// populated from the server response (Computed) rather than producing a perpetual diff, and that
// message_moderation_level falls back to its MODERATE_NONE default.
func TestAccGroupSettings_ComputedDefaults(t *testing.T) {
	server := setupTestServer(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == "PUT" && r.URL.Path == "/team@way.cloud",
			r.Method == "GET" && r.URL.Path == "/team@way.cloud":
			// who_can_post_message is unset in config; the API assigns a value server-side.
			jsonResponse(w, 200, map[string]any{
				"kind":                   "groupsSettings#groups",
				"email":                  "team@way.cloud",
				"whoCanViewMembership":   "ALL_IN_DOMAIN_CAN_VIEW",
				"whoCanPostMessage":      "ALL_MANAGERS_CAN_POST",
				"messageModerationLevel": "MODERATE_NONE",
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
  email                   = "team@way.cloud"
  who_can_view_membership = "ALL_IN_DOMAIN_CAN_VIEW"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Computed from the server response despite being absent from config.
					resource.TestCheckResourceAttr("googleworkspace_group_settings.test", "who_can_post_message", "ALL_MANAGERS_CAN_POST"),
					// Falls back to the schema default.
					resource.TestCheckResourceAttr("googleworkspace_group_settings.test", "message_moderation_level", "MODERATE_NONE"),
				),
			},
		},
	})
}
