package provider

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// fakeGroupMembers is an in-memory Directory API members collection for group-001.
type fakeGroupMembers struct {
	t        *testing.T
	mu       sync.Mutex
	members  map[string]map[string]any
	requests []string
}

func newFakeGroupMembers(t *testing.T, members ...map[string]any) *fakeGroupMembers {
	f := &fakeGroupMembers{t: t, members: map[string]map[string]any{}}
	for _, m := range members {
		f.members[m["email"].(string)] = m
	}
	return f
}

func (f *fakeGroupMembers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f.mu.Lock()
	defer f.mu.Unlock()
	const prefix = "/admin/directory/v1/groups/group-001/members"
	if !strings.HasPrefix(r.URL.Path, prefix) {
		f.t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(500)
		return
	}
	email := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, prefix), "/")
	f.requests = append(f.requests, r.Method+" "+email)
	switch {
	case r.Method == "POST" && email == "":
		var member map[string]any
		if err := json.NewDecoder(r.Body).Decode(&member); err != nil {
			f.t.Errorf("decode request: %v", err)
		}
		key := member["email"].(string)
		if _, ok := f.members[key]; ok {
			jsonResponse(w, 409, map[string]any{"error": map[string]any{"code": 409, "message": "Member already exists."}})
			return
		}
		f.members[key] = member
		jsonResponse(w, 200, member)
	case r.Method == "GET" && email != "":
		member, ok := f.members[email]
		if !ok {
			jsonResponse(w, 404, map[string]any{"error": map[string]any{"code": 404, "message": "Resource Not Found: memberKey"}})
			return
		}
		jsonResponse(w, 200, member)
	case r.Method == "PUT" && email != "":
		var update map[string]any
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			f.t.Errorf("decode request: %v", err)
		}
		f.members[email]["role"] = update["role"]
		jsonResponse(w, 200, f.members[email])
	case r.Method == "DELETE" && email != "":
		delete(f.members, email)
		w.WriteHeader(204)
	default:
		f.t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(500)
	}
}

// touched reports whether any request named the given member.
func (f *fakeGroupMembers) touched(email string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, req := range f.requests {
		if strings.HasSuffix(req, " "+email) {
			return true
		}
	}
	return false
}

func (f *fakeGroupMembers) has(email string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, ok := f.members[email]
	return ok
}

// TestAccGroupMember_Lifecycle verifies create, role update, import and destroy, and that a member
// added outside Terraform is never read or removed.
func TestAccGroupMember_Lifecycle(t *testing.T) {
	fake := newFakeGroupMembers(t, map[string]any{
		"email": "customer@example.com", "role": "MEMBER", "type": "USER",
	})
	server := setupTestServer(t, fake)
	setupTestClient(t, server)

	unmanagedUntouched := func(*terraform.State) error {
		if fake.touched("customer@example.com") {
			t.Error("a request named the unmanaged member")
		}
		if !fake.has("customer@example.com") {
			t.Error("the unmanaged member was removed")
		}
		return nil
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			if fake.has("alice@way.cloud") {
				t.Error("the managed member was not removed on destroy")
			}
			return unmanagedUntouched(s)
		},
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig + `
resource "googleworkspace_group_member" "test" {
  group_id = "group-001"
  email    = "alice@way.cloud"
  role     = "MEMBER"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_group_member.test", "id", "group-001/alice@way.cloud"),
					resource.TestCheckResourceAttr("googleworkspace_group_member.test", "role", "MEMBER"),
					resource.TestCheckResourceAttr("googleworkspace_group_member.test", "type", "USER"),
					unmanagedUntouched,
				),
			},
			{
				Config: testProviderConfig + `
resource "googleworkspace_group_member" "test" {
  group_id = "group-001"
  email    = "alice@way.cloud"
  role     = "MANAGER"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("googleworkspace_group_member.test", "role", "MANAGER"),
					unmanagedUntouched,
				),
			},
			{
				ResourceName:      "googleworkspace_group_member.test",
				ImportState:       true,
				ImportStateId:     "groups/group-001/alice@way.cloud",
				ImportStateVerify: true,
			},
		},
	})
}

// TestAccGroupMember_AlreadyExists verifies that adding a member who is already in the group fails
// with a hint to import it, rather than taking it over.
func TestAccGroupMember_AlreadyExists(t *testing.T) {
	fake := newFakeGroupMembers(t, map[string]any{
		"email": "alice@way.cloud", "role": "MEMBER", "type": "USER",
	})
	server := setupTestServer(t, fake)
	setupTestClient(t, server)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig + `
resource "googleworkspace_group_member" "test" {
  group_id = "group-001"
  email    = "alice@way.cloud"
  role     = "MANAGER"
}
`,
				ExpectError: regexp.MustCompile(`Import it\s+with\s+ID\s+"group-001/alice@way.cloud"`),
			},
		},
	})
}

// TestAccGroupMember_RemovedOutsideTerraform verifies that a member removed outside Terraform is
// planned for creation again.
func TestAccGroupMember_RemovedOutsideTerraform(t *testing.T) {
	fake := newFakeGroupMembers(t)
	server := setupTestServer(t, fake)
	setupTestClient(t, server)

	config := testProviderConfig + `
resource "googleworkspace_group_member" "test" {
  group_id = "group-001"
  email    = "alice@way.cloud"
  role     = "MEMBER"
}
`
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{Config: config},
			{
				PreConfig: func() {
					fake.mu.Lock()
					delete(fake.members, "alice@way.cloud")
					fake.mu.Unlock()
				},
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccGroupMember_InvalidImportId(t *testing.T) {
	server := setupTestServer(t, newFakeGroupMembers(t))
	setupTestClient(t, server)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testProviderConfig + `
resource "googleworkspace_group_member" "test" {
  group_id = "group-001"
  email    = "alice@way.cloud"
  role     = "MEMBER"
}
`,
				ResourceName:  "googleworkspace_group_member.test",
				ImportState:   true,
				ImportStateId: "group-001",
				ExpectError:   regexp.MustCompile(`Expected format: group_id/email`),
			},
		},
	})
}
